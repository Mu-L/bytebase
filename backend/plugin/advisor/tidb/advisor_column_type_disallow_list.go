package tidb

// Framework code is generated by the generator.

import (
	"context"
	"fmt"
	"strings"

	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/common"
	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
	"github.com/bytebase/bytebase/backend/plugin/advisor"
)

var (
	_ advisor.Advisor = (*ColumnTypeDisallowListAdvisor)(nil)
	_ ast.Visitor     = (*columnTypeDisallowListChecker)(nil)
)

func init() {
	advisor.Register(storepb.Engine_TIDB, advisor.MySQLColumnTypeDisallowList, &ColumnTypeDisallowListAdvisor{})
}

// ColumnTypeDisallowListAdvisor is the advisor checking for column type restriction.
type ColumnTypeDisallowListAdvisor struct {
}

// Check checks for column type restriction.
func (*ColumnTypeDisallowListAdvisor) Check(_ context.Context, checkCtx advisor.Context) ([]*storepb.Advice, error) {
	stmtList, ok := checkCtx.AST.([]ast.StmtNode)
	if !ok {
		return nil, errors.Errorf("failed to convert to StmtNode")
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(checkCtx.Rule.Level)
	if err != nil {
		return nil, err
	}
	paylaod, err := advisor.UnmarshalStringArrayTypeRulePayload(checkCtx.Rule.Payload)
	if err != nil {
		return nil, err
	}
	checker := &columnTypeDisallowListChecker{
		level:           level,
		title:           string(checkCtx.Rule.Type),
		typeRestriction: make(map[string]bool),
	}
	for _, tp := range paylaod.List {
		checker.typeRestriction[strings.ToUpper(tp)] = true
	}

	for _, stmt := range stmtList {
		checker.text = stmt.Text()
		checker.line = stmt.OriginTextPosition()
		(stmt).Accept(checker)
	}

	return checker.adviceList, nil
}

type columnTypeDisallowListChecker struct {
	adviceList      []*storepb.Advice
	level           storepb.Advice_Status
	title           string
	text            string
	line            int
	typeRestriction map[string]bool
}

type columnTypeData struct {
	table  string
	column string
	tp     string
	line   int
}

// Enter implements the ast.Visitor interface.
func (checker *columnTypeDisallowListChecker) Enter(in ast.Node) (ast.Node, bool) {
	var columnList []columnTypeData
	switch node := in.(type) {
	case *ast.CreateTableStmt:
		for _, column := range node.Cols {
			if _, exist := checker.typeRestriction[strings.ToUpper(column.Tp.CompactStr())]; exist {
				columnList = append(columnList, columnTypeData{
					table:  node.Table.Name.O,
					column: column.Name.Name.O,
					tp:     strings.ToUpper(column.Tp.CompactStr()),
					line:   column.OriginTextPosition(),
				})
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range node.Specs {
			switch spec.Tp {
			case ast.AlterTableAddColumns:
				for _, column := range spec.NewColumns {
					if _, exist := checker.typeRestriction[strings.ToUpper(column.Tp.CompactStr())]; exist {
						columnList = append(columnList, columnTypeData{
							table:  node.Table.Name.O,
							column: column.Name.Name.O,
							tp:     strings.ToUpper(column.Tp.CompactStr()),
							line:   node.OriginTextPosition(),
						})
					}
				}
			case ast.AlterTableChangeColumn, ast.AlterTableModifyColumn:
				column := spec.NewColumns[0]
				if _, exist := checker.typeRestriction[strings.ToUpper(column.Tp.CompactStr())]; exist {
					columnList = append(columnList, columnTypeData{
						table:  node.Table.Name.O,
						column: column.Name.Name.O,
						tp:     strings.ToUpper(column.Tp.CompactStr()),
						line:   node.OriginTextPosition(),
					})
				}
			default:
			}
		}
	}

	for _, column := range columnList {
		checker.adviceList = append(checker.adviceList, &storepb.Advice{
			Status:        checker.level,
			Code:          advisor.DisabledColumnType.Int32(),
			Title:         checker.title,
			Content:       fmt.Sprintf("Disallow column type %s but column `%s`.`%s` is", column.tp, column.table, column.column),
			StartPosition: common.ConvertANTLRLineToPosition(column.line),
		})
	}

	return in, false
}

// Leave implements the ast.Visitor interface.
func (*columnTypeDisallowListChecker) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
