package pg

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/common"
	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
	"github.com/bytebase/bytebase/backend/plugin/advisor"
	"github.com/bytebase/bytebase/backend/plugin/parser/sql/ast"
)

var (
	_ advisor.Advisor = (*InsertRowLimitAdvisor)(nil)
	_ ast.Visitor     = (*insertRowLimitChecker)(nil)
)

func init() {
	advisor.Register(storepb.Engine_POSTGRES, advisor.PostgreSQLInsertRowLimit, &InsertRowLimitAdvisor{})
}

// InsertRowLimitAdvisor is the advisor checking for to limit INSERT rows.
type InsertRowLimitAdvisor struct {
}

// Check checks for the WHERE clause requirement.
func (*InsertRowLimitAdvisor) Check(ctx context.Context, checkCtx advisor.Context) ([]*storepb.Advice, error) {
	stmts, ok := checkCtx.AST.([]ast.Node)
	if !ok {
		return nil, errors.Errorf("failed to convert to Node")
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(checkCtx.Rule.Level)
	if err != nil {
		return nil, err
	}
	payload, err := advisor.UnmarshalNumberTypeRulePayload(checkCtx.Rule.Payload)
	if err != nil {
		return nil, err
	}
	checker := &insertRowLimitChecker{
		level:                    level,
		title:                    string(checkCtx.Rule.Type),
		maxRow:                   payload.Number,
		driver:                   checkCtx.Driver,
		ctx:                      ctx,
		UsePostgresDatabaseOwner: checkCtx.UsePostgresDatabaseOwner,
	}

	if payload.Number > 0 {
		for _, stmt := range stmts {
			checker.text = stmt.Text()
			checker.line = stmt.LastLine()
			ast.Walk(checker, stmt)
			if checker.explainCount >= common.MaximumLintExplainSize {
				break
			}
		}
	}

	return checker.adviceList, nil
}

type insertRowLimitChecker struct {
	adviceList               []*storepb.Advice
	level                    storepb.Advice_Status
	title                    string
	text                     string
	line                     int
	maxRow                   int
	driver                   *sql.DB
	ctx                      context.Context
	explainCount             int
	setRoles                 []string
	UsePostgresDatabaseOwner bool
}

// Visit implements the ast.Visitor interface.
func (checker *insertRowLimitChecker) Visit(node ast.Node) ast.Visitor {
	code := advisor.Ok
	rows := int64(0)

	switch n := node.(type) {
	case *ast.VariableSetStmt:
		if n.Name == "role" {
			checker.setRoles = append(checker.setRoles, node.Text())
		}
	case *ast.InsertStmt:
		if len(n.ValueList) > 0 {
			// For INSERT INTO ... VALUES ... statements, use parser only.
			if len(n.ValueList) > checker.maxRow {
				code = advisor.InsertTooManyRows
				rows = int64(len(n.ValueList))
			}
		} else if checker.driver != nil {
			// For INSERT INTO ... SELECT statements, use EXPLAIN.
			checker.explainCount++
			res, err := advisor.Query(checker.ctx, advisor.QueryContext{
				UsePostgresDatabaseOwner: checker.UsePostgresDatabaseOwner,
				PreExecutions:            checker.setRoles,
			}, checker.driver, storepb.Engine_POSTGRES, fmt.Sprintf("EXPLAIN %s", node.Text()))
			if err != nil {
				checker.adviceList = append(checker.adviceList, &storepb.Advice{
					Status:        checker.level,
					Code:          advisor.InsertTooManyRows.Int32(),
					Title:         checker.title,
					Content:       fmt.Sprintf("\"%s\" dry runs failed: %s", checker.text, err.Error()),
					StartPosition: common.ConvertPGParserLineToPosition(checker.line),
				})
			} else {
				rowCount, err := getAffectedRows(res)
				if err != nil {
					checker.adviceList = append(checker.adviceList, &storepb.Advice{
						Status:        checker.level,
						Code:          advisor.Internal.Int32(),
						Title:         checker.title,
						Content:       fmt.Sprintf("failed to get row count for \"%s\": %s", checker.text, err.Error()),
						StartPosition: common.ConvertPGParserLineToPosition(checker.line),
					})
				} else if rowCount > int64(checker.maxRow) {
					code = advisor.InsertTooManyRows
					rows = rowCount
				}
			}
		}
	}

	if code != advisor.Ok {
		checker.adviceList = append(checker.adviceList, &storepb.Advice{
			Status:        checker.level,
			Code:          code.Int32(),
			Title:         checker.title,
			Content:       fmt.Sprintf("The statement \"%s\" inserts %d rows. The count exceeds %d.", checker.text, rows, checker.maxRow),
			StartPosition: common.ConvertPGParserLineToPosition(checker.line),
		})
	}
	return checker
}

func getAffectedRows(res []any) (int64, error) {
	// the res struct is []any{columnName, columnTable, rowDataList}
	if len(res) != 3 {
		return 0, errors.Errorf("expected 3 but got %d", len(res))
	}
	rowList, ok := res[2].([]any)
	if !ok {
		return 0, errors.Errorf("expected []any but got %t", res[2])
	}
	// test-bb=# EXPLAIN INSERT INTO t SELECT * FROM t;
	// QUERY PLAN
	// -------------------------------------------------------------
	//  Insert on t  (cost=0.00..1.03 rows=0 width=0)
	//    ->  Seq Scan on t t_1  (cost=0.00..1.03 rows=3 width=520)
	// (2 rows)
	if len(rowList) < 2 {
		return 0, errors.Errorf("not found any data")
	}
	// We need the row 2.
	rowTwo, ok := rowList[1].([]any)
	if !ok {
		return 0, errors.Errorf("expected []any but got %t", rowList[0])
	}
	// PostgreSQL EXPLAIN statement result has one column.
	if len(rowTwo) != 1 {
		return 0, errors.Errorf("expected one but got %d", len(rowTwo))
	}
	// Get the string value.
	text, ok := rowTwo[0].(string)
	if !ok {
		return 0, errors.Errorf("expected string but got %t", rowTwo[0])
	}

	rowsRegexp := regexp.MustCompile("rows=([0-9]+)")
	matches := rowsRegexp.FindStringSubmatch(text)
	if len(matches) != 2 {
		return 0, errors.Errorf("failed to find rows in %q", text)
	}
	value, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, errors.Errorf("failed to get integer from %q", matches[1])
	}
	return value, nil
}
