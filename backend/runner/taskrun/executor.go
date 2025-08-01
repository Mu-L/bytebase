package taskrun

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/bytebase/bytebase/backend/common"
	"github.com/bytebase/bytebase/backend/common/log"
	"github.com/bytebase/bytebase/backend/component/config"
	"github.com/bytebase/bytebase/backend/component/dbfactory"
	"github.com/bytebase/bytebase/backend/component/state"
	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
	"github.com/bytebase/bytebase/backend/plugin/db"
	"github.com/bytebase/bytebase/backend/runner/schemasync"
	"github.com/bytebase/bytebase/backend/store"
	"github.com/bytebase/bytebase/backend/store/model"
)

// Executor is the task executor.
type Executor interface {
	// RunOnce will be called periodically by the scheduler until terminated is true.
	//
	// NOTE
	//
	// 1. It's possible that err could be non-nil while terminated is false, which
	// usually indicates a transient error and will make scheduler retry later.
	// 2. If err is non-nil, then the detail field will be ignored since info is provided in the err.
	// driverCtx is used by the database driver so that we can cancel the query
	// while have the ability to cleanup migration history etc.
	RunOnce(ctx context.Context, driverCtx context.Context, task *store.TaskMessage, taskRunUID int) (terminated bool, result *storepb.TaskRunResult, err error)
}

type execFuncType func(context.Context, string) error

// RunExecutorOnce wraps a TaskExecutor.RunOnce call with panic recovery.
func RunExecutorOnce(ctx context.Context, driverCtx context.Context, exec Executor, task *store.TaskMessage, taskRunUID int) (terminated bool, result *storepb.TaskRunResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			panicErr, ok := r.(error)
			if !ok {
				panicErr = errors.Errorf("%v", r)
			}
			slog.Error("TaskExecutor PANIC RECOVER", log.BBError(panicErr), log.BBStack("panic-stack"))
			terminated = true
			result = nil
			err = errors.Errorf("TaskExecutor PANIC RECOVER, err: %v", panicErr)
		}
	}()

	return exec.RunOnce(ctx, driverCtx, task, taskRunUID)
}

// Pointer fields are not nullable unless mentioned otherwise.
type migrateContext struct {
	syncer    *schemasync.Syncer
	profile   *config.Profile
	dbFactory *dbfactory.DBFactory

	instance *store.InstanceMessage
	database *store.DatabaseMessage
	// nullable if type=baseline
	sheet *store.SheetMessage
	// empty if type=baseline
	sheetName string

	task        *store.TaskMessage
	taskRunUID  int
	taskRunName string
	issueName   string

	version string

	release struct {
		// The release
		// Format: projects/{project}/releases/{release}
		release string
		// The file
		// Format: projects/{project}/releases/{release}/files/{id}
		file string
	}

	// mutable
	// changelog uid
	changelog int64

	migrateType          db.MigrationType
	changeHistoryPayload *storepb.InstanceChangeHistoryPayload
}

func getMigrationInfo(ctx context.Context, stores *store.Store, profile *config.Profile, syncer *schemasync.Syncer, task *store.TaskMessage, migrationType db.MigrationType, schemaVersion string, sheetID *int, taskRunUID int, dbFactory *dbfactory.DBFactory) (*migrateContext, error) {
	instance, err := stores.GetInstanceV2(ctx, &store.FindInstanceMessage{ResourceID: &task.InstanceID})
	if err != nil {
		return nil, err
	}
	database, err := stores.GetDatabaseV2(ctx, &store.FindDatabaseMessage{InstanceID: &task.InstanceID, DatabaseName: task.DatabaseName})
	if err != nil {
		return nil, err
	}
	if database == nil {
		return nil, errors.Errorf("database not found")
	}
	pipeline, err := stores.GetPipelineV2ByID(ctx, task.PipelineID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get pipeline")
	}
	if pipeline == nil {
		return nil, errors.Errorf("pipeline %v not found", task.PipelineID)
	}

	mc := &migrateContext{
		syncer:      syncer,
		profile:     profile,
		dbFactory:   dbFactory,
		instance:    instance,
		database:    database,
		task:        task,
		version:     schemaVersion,
		taskRunName: common.FormatTaskRun(pipeline.ProjectID, task.PipelineID, task.Environment, task.ID, taskRunUID),
		taskRunUID:  taskRunUID,
		migrateType: migrationType,
	}

	if sheetID != nil {
		sheet, err := stores.GetSheet(ctx, &store.FindSheetMessage{UID: sheetID})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get sheet")
		}
		if sheet == nil {
			return nil, errors.Errorf("sheet not found")
		}
		mc.sheet = sheet
		mc.sheetName = common.FormatSheet(pipeline.ProjectID, sheet.UID)
	}

	if isChangeDatabaseTask(task.Type) {
		if f := task.Payload.GetTaskReleaseSource().GetFile(); f != "" {
			project, release, _, err := common.GetProjectReleaseUIDFile(f)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse file %s", f)
			}
			mc.release.release = common.FormatReleaseName(project, release)
			mc.release.file = f
		}
	}

	plans, err := stores.ListPlans(ctx, &store.FindPlanMessage{PipelineID: &task.PipelineID})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list plans")
	}
	if len(plans) == 1 {
		planTypes := []store.PlanCheckRunType{store.PlanCheckDatabaseStatementSummaryReport}
		status := []store.PlanCheckRunStatus{store.PlanCheckRunStatusDone}
		runs, err := stores.ListPlanCheckRuns(ctx, &store.FindPlanCheckRunMessage{
			PlanUID: &plans[0].UID,
			Type:    &planTypes,
			Status:  &status,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list plan check runs")
		}
		slices.SortFunc(runs, func(a, b *store.PlanCheckRunMessage) int {
			if a.UID > b.UID {
				return -1
			} else if a.UID < b.UID {
				return 1
			}
			return 0
		})
		foundChangedResources := false
		for _, run := range runs {
			if foundChangedResources {
				break
			}
			taskInstance, err := stores.GetInstanceV2(ctx, &store.FindInstanceMessage{ResourceID: &task.InstanceID})
			if err != nil {
				return nil, err
			}
			if taskInstance == nil {
				return nil, errors.Errorf("task %d instance not found", task.ID)
			}

			if run.Config.InstanceId != taskInstance.ResourceID {
				continue
			}
			if run.Config.DatabaseName != database.DatabaseName {
				continue
			}
			if sheetID != nil && run.Config.SheetUid != int32(*sheetID) {
				continue
			}
			if run.Result == nil {
				continue
			}
			for _, result := range run.Result.Results {
				if result.Status != storepb.PlanCheckRunResult_Result_SUCCESS {
					continue
				}
				if report := result.GetSqlSummaryReport(); report != nil {
					mc.changeHistoryPayload = &storepb.InstanceChangeHistoryPayload{ChangedResources: report.ChangedResources}
					foundChangedResources = true
					break
				}
			}
		}
	}

	issue, err := stores.GetIssueV2(ctx, &store.FindIssueMessage{PipelineID: &task.PipelineID})
	if err != nil {
		slog.Error("failed to find containing issue", log.BBError(err))
	}
	if issue != nil {
		mc.issueName = common.FormatIssue(issue.Project.ResourceID, issue.UID)
	}
	return mc, nil
}

func getCreateTaskRunLog(ctx context.Context, taskRunUID int, s *store.Store, profile *config.Profile) func(t time.Time, e *storepb.TaskRunLog) error {
	return func(t time.Time, e *storepb.TaskRunLog) error {
		return s.CreateTaskRunLog(ctx, taskRunUID, t.UTC(), profile.DeployID, e)
	}
}

func getUseDatabaseOwner(ctx context.Context, stores *store.Store, instance *store.InstanceMessage, database *store.DatabaseMessage) (bool, error) {
	if instance.Metadata.GetEngine() != storepb.Engine_POSTGRES {
		return false, nil
	}

	// Check the project setting to see if we should use the database owner.
	project, err := stores.GetProjectV2(ctx, &store.FindProjectMessage{ResourceID: &database.ProjectID})
	if err != nil {
		return false, errors.Wrapf(err, "failed to get project")
	}

	if project.Setting == nil {
		return false, nil
	}

	return project.Setting.PostgresDatabaseTenantMode, nil
}

func doMigrationWithFunc(
	ctx context.Context,
	driverCtx context.Context,
	stores *store.Store,
	stateCfg *state.State,
	profile *config.Profile,
	statement string,
	mc *migrateContext,
	execFunc execFuncType,
) (bool, error) {
	instance := mc.instance
	database := mc.database

	useDBOwner, err := getUseDatabaseOwner(ctx, stores, instance, database)
	if err != nil {
		return false, errors.Wrapf(err, "failed to check if we should use database owner")
	}
	driver, err := mc.dbFactory.GetAdminDatabaseDriver(ctx, instance, database, db.ConnectionContext{
		UseDatabaseOwner: useDBOwner,
	})
	if err != nil {
		return false, errors.Wrapf(err, "failed to get driver connection for instance %q", instance.ResourceID)
	}
	defer driver.Close(ctx)

	statementRecord, _ := common.TruncateString(statement, common.MaxSheetSize)
	slog.Debug("Start migration...",
		slog.String("instance", instance.ResourceID),
		slog.String("database", database.DatabaseName),
		slog.String("type", string(mc.migrateType)),
		slog.String("statement", statementRecord),
	)

	opts := db.ExecuteOptions{}

	project, err := stores.GetProjectV2(ctx, &store.FindProjectMessage{ResourceID: &database.ProjectID})
	if err != nil {
		return false, errors.Wrapf(err, "failed to get project %v for database %v", database.ProjectID, database.DatabaseName)
	}
	if project != nil && project.Setting != nil {
		opts.MaximumRetries = int(project.Setting.GetExecutionRetryPolicy().GetMaximumRetries())
	}

	opts.SetConnectionID = func(id string) {
		stateCfg.TaskRunConnectionID.Store(mc.taskRunUID, id)
	}
	opts.DeleteConnectionID = func() {
		stateCfg.TaskRunConnectionID.Delete(mc.taskRunUID)
	}

	if stateCfg != nil {
		switch mc.task.Type {
		case storepb.Task_DATABASE_SCHEMA_UPDATE, storepb.Task_DATABASE_DATA_UPDATE:
			switch instance.Metadata.GetEngine() {
			case storepb.Engine_MYSQL, storepb.Engine_TIDB, storepb.Engine_OCEANBASE,
				storepb.Engine_STARROCKS, storepb.Engine_DORIS, storepb.Engine_POSTGRES,
				storepb.Engine_REDSHIFT, storepb.Engine_RISINGWAVE, storepb.Engine_ORACLE,
				storepb.Engine_DM, storepb.Engine_OCEANBASE_ORACLE, storepb.Engine_MSSQL,
				storepb.Engine_DYNAMODB:
				opts.CreateTaskRunLog = getCreateTaskRunLog(ctx, mc.taskRunUID, stores, profile)
			default:
				// do nothing
			}
		default:
			// No task run log needed for other task types
		}
	}

	if execFunc == nil {
		execFunc = func(ctx context.Context, execStatement string) error {
			if _, err := driver.Execute(ctx, execStatement, opts); err != nil {
				return err
			}
			return nil
		}
	}
	return executeMigrationWithFunc(ctx, driverCtx, stores, mc, statement, execFunc, opts)
}

func postMigration(ctx context.Context, stores *store.Store, mc *migrateContext, skipped bool) (bool, *storepb.TaskRunResult, error) {
	if skipped {
		return true, &storepb.TaskRunResult{
			Detail: fmt.Sprintf("Task skipped because version %s has been applied", mc.version),
		}, nil
	}

	instance := mc.instance
	database := mc.database

	slog.Debug("Post migration...",
		slog.String("instance", instance.ResourceID),
		slog.String("database", database.DatabaseName),
	)

	// Remove schema drift.
	metadata, ok := proto.Clone(database.Metadata).(*storepb.DatabaseMetadata)
	if !ok {
		return false, nil, errors.Errorf("failed to convert database metadata type")
	}
	metadata.Drifted = false
	if _, err := stores.UpdateDatabase(ctx, &store.UpdateDatabaseMessage{
		InstanceID:   database.InstanceID,
		DatabaseName: database.DatabaseName,
		Metadata:     metadata,
	}); err != nil {
		return false, nil, errors.Wrapf(err, "failed to update database %q for instance %q", database.DatabaseName, database.InstanceID)
	}

	detail := fmt.Sprintf("Applied migration version %s to database %q.", mc.version, database.DatabaseName)
	if mc.migrateType == db.Baseline {
		detail = fmt.Sprintf("Established baseline version %s for database %q.", mc.version, database.DatabaseName)
	}

	return true, &storepb.TaskRunResult{
		Detail:    detail,
		Changelog: common.FormatChangelog(instance.ResourceID, database.DatabaseName, mc.changelog),
		Version:   mc.version,
	}, nil
}

func runMigrationWithFunc(ctx context.Context, driverCtx context.Context, store *store.Store, dbFactory *dbfactory.DBFactory, stateCfg *state.State, syncer *schemasync.Syncer, profile *config.Profile, task *store.TaskMessage, taskRunUID int, migrationType db.MigrationType, statement string, schemaVersion string, sheetID *int, execFunc execFuncType) (terminated bool, result *storepb.TaskRunResult, err error) {
	mc, err := getMigrationInfo(ctx, store, profile, syncer, task, migrationType, schemaVersion, sheetID, taskRunUID, dbFactory)
	if err != nil {
		return true, nil, err
	}

	skipped, err := doMigrationWithFunc(ctx, driverCtx, store, stateCfg, profile, statement, mc, execFunc)
	if err != nil {
		return true, nil, err
	}
	return postMigration(ctx, store, mc, skipped)
}

func runMigration(ctx context.Context, driverCtx context.Context, store *store.Store, dbFactory *dbfactory.DBFactory, stateCfg *state.State, syncer *schemasync.Syncer, profile *config.Profile, task *store.TaskMessage, taskRunUID int, migrationType db.MigrationType, statement string, schemaVersion string, sheetID *int) (terminated bool, result *storepb.TaskRunResult, err error) {
	return runMigrationWithFunc(ctx, driverCtx, store, dbFactory, stateCfg, syncer, profile, task, taskRunUID, migrationType, statement, schemaVersion, sheetID, nil /* default */)
}

// executeMigrationWithFunc executes the migration with custom migration function.
func executeMigrationWithFunc(ctx context.Context, driverCtx context.Context, s *store.Store, mc *migrateContext, statement string, execFunc func(ctx context.Context, execStatement string) error, opts db.ExecuteOptions) (skipped bool, resErr error) {
	// Phase 1 - Dump before migration.
	// Check if versioned is already applied.
	skipExecution, err := beginMigration(ctx, s, mc, opts)
	if err != nil {
		return false, errors.Wrapf(err, "failed to begin migration")
	}
	if skipExecution {
		return true, nil
	}

	defer func() {
		// Phase 3 - Dump after migration.
		// Insert revision for versioned.
		if err := endMigration(ctx, s, mc, resErr == nil /* isDone */); err != nil {
			slog.Error("failed to end migration",
				log.BBError(err),
			)
		}
	}()

	// Phase 2 - Executing migration.
	if mc.migrateType != db.Baseline {
		// Database secrets feature has been removed
		// To avoid leak the rendered statement, the error message should use the original statement and not the rendered statement.
		// Database secrets feature removed - using original statement directly
		if err := execFunc(driverCtx, statement); err != nil {
			return false, err
		}
	}

	return false, nil
}

// beginMigration checks before executing migration and inserts a migration history record with pending status.
func beginMigration(ctx context.Context, stores *store.Store, mc *migrateContext, opts db.ExecuteOptions) (bool, error) {
	// list revisions and see if it has been applied
	// we can do this because
	// versioned migrations are executed one by one
	// so no other migrations can insert revisions
	//
	// users can create revisions though via API
	// however we can warn users not to unless they know
	// what they are doing
	if mc.version != "" {
		list, err := stores.ListRevisions(ctx, &store.FindRevisionMessage{
			InstanceID:   &mc.database.InstanceID,
			DatabaseName: &mc.database.DatabaseName,
			Version:      &mc.version,
		})
		if err != nil {
			return false, errors.Wrapf(err, "failed to list revisions")
		}
		if len(list) > 0 {
			// This version has been executed.
			// skip execution.
			return true, nil
		}
	}

	// sync history
	var syncHistoryPrevUID *int64
	if mc.migrateType.NeedDump() {
		opts.LogDatabaseSyncStart()
		syncHistoryPrev, err := mc.syncer.SyncDatabaseSchemaToHistory(ctx, mc.database)
		if err != nil {
			opts.LogDatabaseSyncEnd(err.Error())
			return false, errors.Wrapf(err, "failed to sync database metadata and schema")
		}
		opts.LogDatabaseSyncEnd("")
		syncHistoryPrevUID = &syncHistoryPrev
	}

	// create pending changelog
	changelogUID, err := stores.CreateChangelog(ctx, &store.ChangelogMessage{
		InstanceID:         mc.database.InstanceID,
		DatabaseName:       mc.database.DatabaseName,
		Status:             store.ChangelogStatusPending,
		PrevSyncHistoryUID: syncHistoryPrevUID,
		SyncHistoryUID:     nil,
		Payload: &storepb.ChangelogPayload{
			TaskRun:          mc.taskRunName,
			Issue:            mc.issueName,
			Revision:         0,
			ChangedResources: mc.changeHistoryPayload.GetChangedResources(),
			Sheet:            mc.sheetName,
			Version:          mc.version,
			Type:             convertTaskType(mc.task.Type),
			GitCommit:        mc.profile.GitCommit,
		}})
	if err != nil {
		return false, errors.Wrapf(err, "failed to create changelog")
	}
	mc.changelog = changelogUID

	return false, nil
}

// endMigration updates the migration history record to DONE or FAILED depending on migration is done or not.
func endMigration(ctx context.Context, storeInstance *store.Store, mc *migrateContext, isDone bool) error {
	update := &store.UpdateChangelogMessage{
		UID: mc.changelog,
	}

	if mc.migrateType.NeedDump() {
		syncHistory, err := mc.syncer.SyncDatabaseSchemaToHistory(ctx, mc.database)
		if err != nil {
			return errors.Wrapf(err, "failed to sync database metadata and schema")
		}
		update.SyncHistoryUID = &syncHistory
	}

	if isDone {
		// if isDone, record in revision
		if mc.version != "" {
			r := &store.RevisionMessage{
				InstanceID:   mc.database.InstanceID,
				DatabaseName: mc.database.DatabaseName,
				Version:      mc.version,
				Payload: &storepb.RevisionPayload{
					Release:     mc.release.release,
					File:        mc.release.file,
					Sheet:       "",
					SheetSha256: "",
					TaskRun:     mc.taskRunName,
				},
			}
			if mc.sheet != nil {
				r.Payload.Sheet = mc.sheetName
				r.Payload.SheetSha256 = mc.sheet.GetSha256Hex()
			}

			revision, err := storeInstance.CreateRevision(ctx, r)
			if err != nil {
				return errors.Wrapf(err, "failed to create revision")
			}
			update.RevisionUID = &revision.UID

			// Update database metadata with the version only if the new version is greater
			metadata := mc.database.Metadata
			if metadata == nil {
				metadata = &storepb.DatabaseMetadata{}
			}

			if shouldUpdateVersion(metadata.Version, mc.version) {
				metadata.Version = mc.version
				if _, err := storeInstance.UpdateDatabase(ctx, &store.UpdateDatabaseMessage{
					InstanceID:   mc.database.InstanceID,
					DatabaseName: mc.database.DatabaseName,
					Metadata:     metadata,
				}); err != nil {
					return errors.Wrapf(err, "failed to update database metadata with version")
				}
			}
		}
		status := store.ChangelogStatusDone
		update.Status = &status
	} else {
		status := store.ChangelogStatusFailed
		update.Status = &status
	}

	if err := storeInstance.UpdateChangelog(ctx, update); err != nil {
		return errors.Wrapf(err, "failed to update changelog")
	}

	return nil
}

// shouldUpdateVersion checks if newVersion is greater than currentVersion.
// Returns true if:
// - currentVersion is empty
// - currentVersion is invalid
// - newVersion is greater than currentVersion
func shouldUpdateVersion(currentVersion, newVersion string) bool {
	if currentVersion == "" {
		// If no current version, always update
		return true
	}
	current, err := model.NewVersion(currentVersion)
	if err != nil {
		// If current version is invalid, update with new version
		return true
	}

	nv, err := model.NewVersion(newVersion)
	if err != nil {
		// If new version is invalid, don't update
		return false
	}
	return current.LessThan(nv)
}

func convertTaskType(t storepb.Task_Type) storepb.ChangelogPayload_Type {
	switch t {
	case storepb.Task_DATABASE_DATA_UPDATE:
		return storepb.ChangelogPayload_DATA
	case storepb.Task_DATABASE_SCHEMA_UPDATE:
		return storepb.ChangelogPayload_MIGRATE
	case storepb.Task_DATABASE_SCHEMA_UPDATE_GHOST:
		return storepb.ChangelogPayload_MIGRATE_GHOST

	case storepb.Task_DATABASE_CREATE:
	case storepb.Task_DATABASE_EXPORT:
	default:
		// Return unspecified for other task types
	}
	return storepb.ChangelogPayload_TYPE_UNSPECIFIED
}

// isChangeDatabaseTask returns whether the task involves changing a database.
func isChangeDatabaseTask(taskType storepb.Task_Type) bool {
	switch taskType {
	case storepb.Task_DATABASE_DATA_UPDATE,
		storepb.Task_DATABASE_SCHEMA_UPDATE,
		storepb.Task_DATABASE_SCHEMA_UPDATE_GHOST:
		return true
	case storepb.Task_DATABASE_CREATE,
		storepb.Task_DATABASE_EXPORT:
		return false
	default:
		return false
	}
}
