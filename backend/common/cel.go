//nolint:revive
package common

import (
	"time"

	"connectrpc.com/connect"
	"github.com/google/cel-go/cel"
	celtypes "github.com/google/cel-go/common/types"
	"github.com/pkg/errors"
	exprproto "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/genproto/googleapis/type/expr"
)

const celLimit = 1024 * 1024

// RiskFactors are the variables when evaluating the risk level.
var RiskFactors = []cel.EnvOption{
	cel.Variable("environment_id", cel.StringType), // use environment.resource_id
	cel.Variable("project_id", cel.StringType),     // use project.resource_id
	cel.Variable("db_engine", cel.StringType),

	cel.Variable("database_name", cel.StringType),
	cel.Variable("schema_name", cel.StringType),
	cel.Variable("table_name", cel.StringType),

	cel.Variable("affected_rows", cel.IntType),
	cel.Variable("table_rows", cel.IntType),
	cel.Variable("sql_type", cel.StringType),
	cel.Variable("sql_statement", cel.StringType),

	cel.Variable("expiration_days", cel.IntType),
	cel.Variable("export_rows", cel.IntType),
	cel.Variable("role", cel.StringType),
}

// ApprovalFactors are the variables when finding the approval template.
var ApprovalFactors = []cel.EnvOption{
	cel.Variable("level", cel.IntType),
	cel.Variable("source", cel.StringType),
	cel.ParserExpressionSizeLimit(celLimit),
}

// IAMPolicyConditionCELAttributes are the variables when evaluating IAM policy condition.
var IAMPolicyConditionCELAttributes = []cel.EnvOption{
	cel.Variable("resource.environment_name", cel.StringType),
	cel.Variable("resource.database", cel.StringType),
	cel.Variable("resource.schema", cel.StringType),
	cel.Variable("resource.table", cel.StringType),
	cel.Variable("request.row_limit", cel.IntType),
	cel.Variable("request.time", cel.TimestampType),
	cel.ParserExpressionSizeLimit(celLimit),
}

// MaskingRulePolicyCELAttributes are the variables when evaluating masking rule.
var MaskingRulePolicyCELAttributes = []cel.EnvOption{
	cel.Variable("environment_id", cel.StringType),
	cel.Variable("project_id", cel.StringType),
	cel.Variable("instance_id", cel.StringType),
	cel.Variable("database_name", cel.StringType),
	cel.Variable("schema_name", cel.StringType),
	cel.Variable("table_name", cel.StringType),
	cel.Variable("column_name", cel.StringType),
	cel.Variable("classification_level", cel.StringType),
	cel.ParserExpressionSizeLimit(celLimit),
}

// MaskingExceptionPolicyCELAttributes are the variables when evaluating masking exception.
var MaskingExceptionPolicyCELAttributes = []cel.EnvOption{
	cel.Variable("resource.instance_id", cel.StringType),
	cel.Variable("resource.database_name", cel.StringType),
	cel.Variable("resource.table_name", cel.StringType),
	cel.Variable("resource.schema_name", cel.StringType),
	cel.Variable("resource.column_name", cel.StringType),
	cel.Variable("request.time", cel.TimestampType),
	cel.ParserExpressionSizeLimit(celLimit),
}

// ConvertUnparsedRisk converts unparsed risk to parsed format.
func ConvertUnparsedRisk(expression *expr.Expr) (*exprproto.ParsedExpr, error) {
	if expression == nil || expression.Expression == "" {
		return nil, nil
	}
	e, err := cel.NewEnv(RiskFactors...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cel env")
	}

	ast, issues := e.Parse(expression.Expression)
	if issues != nil && issues.Err() != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Errorf("failed to parse expression: %v", issues.Err()))
	}
	expr, err := cel.AstToParsedExpr(ast)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Errorf("failed to convert ast to parsed expression: %v", err))
	}
	return expr, nil
}

// ConvertUnparsedApproval converts unparsed approval to parsed format.
func ConvertUnparsedApproval(expression *expr.Expr) (*exprproto.ParsedExpr, error) {
	if expression == nil || expression.Expression == "" {
		return nil, nil
	}
	e, err := cel.NewEnv(ApprovalFactors...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cel env")
	}

	ast, issues := e.Parse(expression.Expression)
	if issues != nil && issues.Err() != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Errorf("failed to parse expression: %v", issues.Err()))
	}
	expr, err := cel.AstToParsedExpr(ast)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.Errorf("failed to convert ast to parsed expression: %v", err))
	}
	return expr, nil
}

// ValidateGroupCELExpr validates group expr.
func ValidateGroupCELExpr(expr string) (cel.Program, error) {
	e, err := cel.NewEnv(
		cel.Variable("resource", cel.MapType(cel.StringType, cel.AnyType)),
	)
	if err != nil {
		return nil, err
	}
	ast, issues := e.Compile(expr)
	if issues != nil && issues.Err() != nil {
		return nil, issues.Err()
	}
	prog, err := e.Program(ast)
	if err != nil {
		return nil, err
	}
	return prog, nil
}

// ValidateMaskingRuleCELExpr validates masking rule expr.
func ValidateMaskingRuleCELExpr(expr string) (cel.Program, error) {
	e, err := cel.NewEnv(
		MaskingRulePolicyCELAttributes...,
	)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	ast, issues := e.Compile(expr)
	if issues != nil && issues.Err() != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, issues.Err())
	}
	prog, err := e.Program(ast)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return prog, nil
}

// ValidateMaskingExceptionCELExpr validates masking exception expr.
func ValidateMaskingExceptionCELExpr(expression *expr.Expr) (cel.Program, error) {
	return validateCELExpr(expression, MaskingExceptionPolicyCELAttributes)
}

func ValidateProjectMemberCELExpr(expression *expr.Expr) (cel.Program, error) {
	return validateCELExpr(expression, IAMPolicyConditionCELAttributes)
}

func validateCELExpr(expression *expr.Expr, conditionCELAttributes []cel.EnvOption) (cel.Program, error) {
	if expression == nil || expression.Expression == "" {
		return nil, nil
	}
	e, err := cel.NewEnv(
		conditionCELAttributes...,
	)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	ast, issues := e.Compile(expression.Expression)
	if issues != nil && issues.Err() != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, issues.Err())
	}
	prog, err := e.Program(ast)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return prog, nil
}

// QueryExportFactors is the factors for query and export.
type QueryExportFactors struct {
	Databases  []string
	ExportRows int64
}

// GetQueryExportFactors is used to get risk factors from query and export expressions.
func GetQueryExportFactors(expression string) (*QueryExportFactors, error) {
	// If the expression is empty, return an empty struct.
	if expression == "" {
		return &QueryExportFactors{}, nil
	}

	factors := &QueryExportFactors{}
	e, err := cel.NewEnv(IAMPolicyConditionCELAttributes...)
	if err != nil {
		return nil, err
	}
	ast, issues := e.Compile(expression)
	if issues != nil {
		return nil, errors.Errorf("found issue %v", issues)
	}
	parsedExpr, err := cel.AstToParsedExpr(ast)
	if err != nil {
		return nil, err
	}
	callExpr := parsedExpr.Expr.GetCallExpr()
	findField(callExpr, factors)
	return factors, nil
}

func findField(callExpr *exprproto.Expr_Call, factors *QueryExportFactors) {
	if callExpr == nil {
		return
	}
	if len(callExpr.Args) == 2 {
		idExpr := callExpr.Args[0].GetIdentExpr()
		if idExpr != nil {
			if idExpr.Name == "request.row_limit" {
				factors.ExportRows = callExpr.Args[1].GetConstExpr().GetInt64Value()
			}
			if idExpr.Name == "resource.database" && callExpr.Function == "_==_" {
				factors.Databases = append(factors.Databases, callExpr.Args[1].GetConstExpr().GetStringValue())
			}
			if idExpr.Name == "resource.database" && callExpr.Function == "@in" {
				list := callExpr.Args[1].GetListExpr()
				for _, element := range list.Elements {
					factors.Databases = append(factors.Databases, element.GetConstExpr().GetStringValue())
				}
			}
			return
		}
	}
	for _, arg := range callExpr.Args {
		callExpr := arg.GetCallExpr()
		findField(callExpr, factors)
	}
}

func EvalBindingCondition(expr string, requestTime time.Time) (bool, error) {
	input := map[string]any{
		"request.time": requestTime,
	}
	return doEvalBindingCondition(expr, input)
}

func doEvalBindingCondition(expr string, input map[string]any) (bool, error) {
	if expr == "" {
		return true, nil
	}

	e, err := cel.NewEnv(IAMPolicyConditionCELAttributes...)
	if err != nil {
		return false, errors.Wrapf(err, "failed to new cel env")
	}
	ast, iss := e.Compile(expr)
	if iss != nil && iss.Err() != nil {
		return false, errors.Wrapf(iss.Err(), "failed to compile expr %q", expr)
	}
	// enable partial evaluation because the input only has request.time
	// but the expression can have more.
	prg, err := e.Program(ast, cel.EvalOptions(cel.OptPartialEval))
	if err != nil {
		return false, errors.Wrapf(err, "failed to construct program")
	}
	vars, err := e.PartialVars(input)
	if err != nil {
		return false, errors.Wrapf(err, "failed to get vars")
	}
	out, _, err := prg.Eval(vars)
	if err != nil {
		return false, errors.Wrapf(err, "failed to eval cel expr")
	}
	// `out` is one of
	// - True
	// - False
	// - a residual expression.

	// return true if the result is a residual expression
	// which means that it passes "the request.time < xxx" check.
	if !celtypes.IsBool(out) {
		return true, nil
	}

	res, ok := out.Equal(celtypes.True).Value().(bool)
	if !ok {
		return false, errors.Errorf("failed to convert cel result to bool")
	}
	return res, nil
}
