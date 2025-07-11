package mysql

// Framework code is generated by the generator.

import (
	"context"
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/pkg/errors"

	mysql "github.com/bytebase/mysql-parser"

	"github.com/bytebase/bytebase/backend/common"
	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
	"github.com/bytebase/bytebase/backend/plugin/advisor"
	mysqlparser "github.com/bytebase/bytebase/backend/plugin/parser/mysql"
)

var (
	_ advisor.Advisor = (*ColumnDisallowSetCharsetAdvisor)(nil)
)

func init() {
	advisor.Register(storepb.Engine_MYSQL, advisor.MySQLDisallowSetColumnCharset, &ColumnDisallowSetCharsetAdvisor{})
	advisor.Register(storepb.Engine_MARIADB, advisor.MySQLDisallowSetColumnCharset, &ColumnDisallowSetCharsetAdvisor{})
	advisor.Register(storepb.Engine_OCEANBASE, advisor.MySQLDisallowSetColumnCharset, &ColumnDisallowSetCharsetAdvisor{})
}

// ColumnDisallowSetCharsetAdvisor is the advisor checking for disallow set column charset.
type ColumnDisallowSetCharsetAdvisor struct {
}

// Check checks for disallow set column charset.
func (*ColumnDisallowSetCharsetAdvisor) Check(_ context.Context, checkCtx advisor.Context) ([]*storepb.Advice, error) {
	stmtList, ok := checkCtx.AST.([]*mysqlparser.ParseResult)
	if !ok {
		return nil, errors.Errorf("failed to convert to mysql parse result")
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(checkCtx.Rule.Level)
	if err != nil {
		return nil, err
	}
	checker := &columnDisallowSetCharsetChecker{
		level: level,
		title: string(checkCtx.Rule.Type),
	}

	for _, stmt := range stmtList {
		checker.baseLine = stmt.BaseLine
		antlr.ParseTreeWalkerDefault.Walk(checker, stmt.Tree)
	}

	return checker.adviceList, nil
}

type columnDisallowSetCharsetChecker struct {
	*mysql.BaseMySQLParserListener

	baseLine   int
	adviceList []*storepb.Advice
	level      storepb.Advice_Status
	title      string
	text       string
}

func (checker *columnDisallowSetCharsetChecker) EnterQuery(ctx *mysql.QueryContext) {
	checker.text = ctx.GetParser().GetTokenStream().GetTextFromRuleContext(ctx)
}

func (checker *columnDisallowSetCharsetChecker) EnterCreateTable(ctx *mysql.CreateTableContext) {
	if !mysqlparser.IsTopMySQLRule(&ctx.BaseParserRuleContext) {
		return
	}
	if ctx.TableElementList() == nil || ctx.TableName() == nil {
		return
	}

	for _, tableElement := range ctx.TableElementList().AllTableElement() {
		if tableElement.ColumnDefinition() == nil {
			continue
		}
		if tableElement.ColumnDefinition().FieldDefinition() == nil {
			continue
		}
		if tableElement.ColumnDefinition().FieldDefinition().DataType() == nil {
			continue
		}
		charset := checker.getCharSet(tableElement.ColumnDefinition().FieldDefinition().DataType())
		if !checker.checkCharset(charset) {
			checker.adviceList = append(checker.adviceList, &storepb.Advice{
				Status:        checker.level,
				Code:          advisor.SetColumnCharset.Int32(),
				Title:         checker.title,
				Content:       fmt.Sprintf("Disallow set column charset but \"%s\" does", checker.text),
				StartPosition: common.ConvertANTLRLineToPosition(checker.baseLine + ctx.GetStart().GetLine()),
			})
		}
	}
}

func (checker *columnDisallowSetCharsetChecker) EnterAlterTable(ctx *mysql.AlterTableContext) {
	if !mysqlparser.IsTopMySQLRule(&ctx.BaseParserRuleContext) {
		return
	}
	if ctx.AlterTableActions() == nil {
		return
	}
	if ctx.AlterTableActions().AlterCommandList() == nil {
		return
	}
	if ctx.AlterTableActions().AlterCommandList().AlterList() == nil {
		return
	}

	_, tableName := mysqlparser.NormalizeMySQLTableRef(ctx.TableRef())
	if tableName == "" {
		return
	}
	// alter table add column, change column, modify column.
	for _, item := range ctx.AlterTableActions().AlterCommandList().AlterList().AllAlterListItem() {
		if item == nil {
			continue
		}

		var charsetList []string
		switch {
		// add column.
		case item.ADD_SYMBOL() != nil:
			switch {
			case item.Identifier() != nil && item.FieldDefinition() != nil:
				if item.FieldDefinition().DataType() == nil {
					continue
				}

				charsetName := checker.getCharSet(item.FieldDefinition().DataType())
				charsetList = append(charsetList, charsetName)
			case item.OPEN_PAR_SYMBOL() != nil && item.TableElementList() != nil:
				for _, tableElement := range item.TableElementList().AllTableElement() {
					if tableElement.ColumnDefinition() == nil {
						continue
					}
					if tableElement.ColumnDefinition().FieldDefinition() == nil {
						continue
					}
					if tableElement.ColumnDefinition().FieldDefinition().DataType() == nil {
						continue
					}

					charsetName := checker.getCharSet(tableElement.ColumnDefinition().FieldDefinition().DataType())
					charsetList = append(charsetList, charsetName)
				}
			}
		// change column.
		case item.CHANGE_SYMBOL() != nil && item.ColumnInternalRef() != nil && item.Identifier() != nil && item.FieldDefinition() != nil:
			charsetName := checker.getCharSet(item.FieldDefinition().DataType())
			charsetList = append(charsetList, charsetName)
		// modify column.
		case item.MODIFY_SYMBOL() != nil && item.ColumnInternalRef() != nil && item.FieldDefinition() != nil:
			charsetName := checker.getCharSet(item.FieldDefinition().DataType())
			charsetList = append(charsetList, charsetName)
		default:
			continue
		}

		for _, charsetName := range charsetList {
			if !checker.checkCharset(charsetName) {
				checker.adviceList = append(checker.adviceList, &storepb.Advice{
					Status:        checker.level,
					Code:          advisor.SetColumnCharset.Int32(),
					Title:         checker.title,
					Content:       fmt.Sprintf("Disallow set column charset but \"%s\" does", checker.text),
					StartPosition: common.ConvertANTLRLineToPosition(checker.baseLine + ctx.GetStart().GetLine()),
				})
			}
		}
	}
}

func (*columnDisallowSetCharsetChecker) getCharSet(ctx mysql.IDataTypeContext) string {
	if ctx.CharsetWithOptBinary() == nil {
		return ""
	}
	charset := mysqlparser.NormalizeMySQLCharsetName(ctx.CharsetWithOptBinary().CharsetName())
	return charset
}

func (*columnDisallowSetCharsetChecker) checkCharset(charset string) bool {
	switch charset {
	// empty charset or binary for JSON.
	case "", "binary":
		return true
	default:
		return false
	}
}
