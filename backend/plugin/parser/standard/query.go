package standard

import (
	"log/slog"
	"regexp"
	"strings"

	"github.com/bytebase/bytebase/backend/common/log"
	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
	"github.com/bytebase/bytebase/backend/plugin/parser/base"
	"github.com/bytebase/bytebase/backend/plugin/parser/tokenizer"
)

func init() {
	base.RegisterQueryValidator(storepb.Engine_CLICKHOUSE, ValidateSQLForEditor)
	base.RegisterQueryValidator(storepb.Engine_SQLITE, ValidateSQLForEditor)
	base.RegisterQueryValidator(storepb.Engine_SPANNER, ValidateSQLForEditor)
	base.RegisterQueryValidator(storepb.Engine_HIVE, ValidateSQLForEditor)
	base.RegisterQueryValidator(storepb.Engine_BIGQUERY, ValidateSQLForEditor)
	base.RegisterQueryValidator(storepb.Engine_CASSANDRA, ValidateSQLForEditor)
}

// ValidateSQLForEditor validates the SQL statement for SQL editor.
// We validate the statement by following steps:
// 1. Remove all quoted text(quoted identifier, string literal) and comments from the statement.
// 2. Use regexp to check if the statement is a normal SELECT statement and EXPLAIN statement.
// 3. For CTE, use regexp to check if the statement has UPDATE, DELETE and INSERT statements.
func ValidateSQLForEditor(statement string) (bool, bool, error) {
	textWithoutQuotedAndComment, err := tokenizer.StandardRemoveQuotedTextAndComment(statement)
	if err != nil {
		slog.Debug("Failed to remove quoted text and comment", slog.String("statement", statement), log.BBError(err))
		return false, false, err
	}

	ok := CheckStatementWithoutQuotedTextAndComment(textWithoutQuotedAndComment)
	return ok, ok, nil
}

func CheckStatementWithoutQuotedTextAndComment(statement string) bool {
	formattedStr := strings.ToUpper(strings.TrimSpace(statement))
	if isSelect, _ := regexp.MatchString(`^SELECT\s+?`, formattedStr); isSelect {
		return true
	}

	if isSelect, _ := regexp.MatchString(`^SELECT\*\s+?`, formattedStr); isSelect {
		return true
	}

	if isExplain, _ := regexp.MatchString(`^EXPLAIN\s+?`, formattedStr); isExplain {
		if isExplainAnalyze, _ := regexp.MatchString(`^EXPLAIN\s+ANALYZE\s+?`, formattedStr); isExplainAnalyze {
			return false
		}
		return true
	}

	cteRegex := regexp.MustCompile(`^WITH\s+?`)
	if matchResult := cteRegex.MatchString(formattedStr); matchResult {
		dmlRegs := []string{`\bINSERT\b`, `\bUPDATE\b`, `\bDELETE\b`}
		for _, reg := range dmlRegs {
			if matchResult, _ := regexp.MatchString(reg, formattedStr); matchResult {
				return false
			}
		}
		return true
	}

	return false
}
