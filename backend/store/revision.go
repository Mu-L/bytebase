package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/bytebase/bytebase/backend/common"
	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
)

type RevisionMessage struct {
	InstanceID   string
	DatabaseName string
	Version      string
	Payload      *storepb.RevisionPayload

	// output only
	UID        int64
	CreatedAt  time.Time
	DeleterUID *int
	DeletedAt  *time.Time
}

type FindRevisionMessage struct {
	UID          *int64
	InstanceID   *string
	DatabaseName *string

	Version  *string
	Versions *[]string

	Limit  *int
	Offset *int

	ShowDeleted bool
}

func (s *Store) ListRevisions(ctx context.Context, find *FindRevisionMessage) ([]*RevisionMessage, error) {
	where, args := []string{"TRUE"}, []any{}

	if v := find.UID; v != nil {
		where = append(where, fmt.Sprintf("id = $%d", len(args)+1))
		args = append(args, *v)
	}
	if v := find.InstanceID; v != nil {
		where = append(where, fmt.Sprintf("instance = $%d", len(args)+1))
		args = append(args, *v)
	}
	if v := find.DatabaseName; v != nil {
		where = append(where, fmt.Sprintf("db_name = $%d", len(args)+1))
		args = append(args, *v)
	}
	if v := find.Version; v != nil {
		where = append(where, fmt.Sprintf("version = $%d", len(args)+1))
		args = append(args, *v)
	}
	if v := find.Versions; v != nil {
		where = append(where, fmt.Sprintf("version = ANY($%d)", len(args)+1))
		args = append(args, *v)
	}
	if !find.ShowDeleted {
		where = append(where, "deleted_at IS NULL")
	}

	limitOffsetClause := ""
	if v := find.Limit; v != nil {
		limitOffsetClause += fmt.Sprintf(" LIMIT %d", *v)
	}
	if v := find.Offset; v != nil {
		limitOffsetClause += fmt.Sprintf(" OFFSET %d", *v)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			instance,
			db_name,
			created_at,
			deleter_id,
			deleted_at,
			version,
			payload
		FROM revision
		WHERE %s
		ORDER BY version DESC
		%s
	`, strings.Join(where, " AND "), limitOffsetClause)

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to begin tx")
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query context")
	}
	defer rows.Close()

	var revisions []*RevisionMessage
	for rows.Next() {
		r := RevisionMessage{
			Payload: &storepb.RevisionPayload{},
		}
		var p []byte
		if err := rows.Scan(
			&r.UID,
			&r.InstanceID,
			&r.DatabaseName,
			&r.CreatedAt,
			&r.DeleterUID,
			&r.DeletedAt,
			&r.Version,
			&p,
		); err != nil {
			return nil, errors.Wrapf(err, "failed to scan")
		}

		if err := common.ProtojsonUnmarshaler.Unmarshal(p, r.Payload); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal")
		}

		revisions = append(revisions, &r)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "rows err")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "failed to commit tx")
	}

	return revisions, nil
}

func (s *Store) GetRevision(ctx context.Context, uid int64) (*RevisionMessage, error) {
	revisions, err := s.ListRevisions(ctx, &FindRevisionMessage{UID: &uid, ShowDeleted: true})
	if err != nil {
		return nil, err
	}
	if len(revisions) == 0 {
		return nil, errors.Errorf("revision not found: %d", uid)
	}
	if len(revisions) > 1 {
		return nil, errors.Errorf("found multiple revisions for uid: %d", uid)
	}
	return revisions[0], nil
}

func (s *Store) CreateRevision(ctx context.Context, revision *RevisionMessage) (*RevisionMessage, error) {
	query := `
		INSERT INTO revision (
			instance,
			db_name,
			version,
			payload
		) VALUES (
		 	$1,
			$2,
			$3,
			$4
		)
		RETURNING id, created_at
	`

	p, err := protojson.Marshal(revision.Payload)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal revision payload")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to begin tx")
	}
	defer tx.Rollback()

	var id int64
	if err := tx.QueryRowContext(ctx, query,
		revision.InstanceID,
		revision.DatabaseName,
		revision.Version,
		p,
	).Scan(&id, &revision.CreatedAt); err != nil {
		return nil, errors.Wrapf(err, "failed to query and scan")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "failed to commit tx")
	}

	revision.UID = id

	return revision, nil
}

func (s *Store) DeleteRevision(ctx context.Context, uid int64, deleterUID int) error {
	query :=
		`UPDATE revision
		SET deleter_id = $1, deleted_at = now()
		WHERE id = $2`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to begin tx")
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, query, deleterUID, uid); err != nil {
		return errors.Wrapf(err, "failed to exec")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, "failed to commit tx")
	}

	return nil
}
