package calculation

import (
	"errors"
	"fmt"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
)

type closureDetails struct {
	Parameters        []string
	HasSingleReturn   bool
	HasArrayReturn    bool
	KeyedReturnFields []string
}

// parseAndWalk examines a javascript arrow function to extract its parameters and possible returns.
func parseAndWalk(src string) (*closureDetails, error) {
	p, err := parser.ParseFile(nil, "", src, 0)
	if err != nil {
		return nil, fmt.Errorf("parsing source: %w", err)
	}

	// find arrow function
	var arrowFn *ast.ArrowFunctionLiteral
	for _, v := range p.Body {
		expr, ok := v.(*ast.ExpressionStatement)
		if !ok {
			continue
		}

		candidate, ok := expr.Expression.(*ast.ArrowFunctionLiteral)
		if !ok {
			continue
		}

		if candidate.Body == nil {
			continue
		}

		arrowFn = candidate
		break
	}

	if arrowFn == nil {
		return nil, errors.New("no arrow function declared in source")
	}

	cd := closureDetails{}

	// get parameters
	cd.Parameters = make([]string, 0)
	for _, param := range arrowFn.ParameterList.List {
		id, ok := param.Target.(*ast.Identifier)
		if !ok {
			continue
		}
		cd.Parameters = append(cd.Parameters, id.Name.String())
	}

	// get outputs (note: limited subset of js features explicitly supported here)
	keyTracking := make(map[string]struct{})
	nodes := []ast.Node{arrowFn.Body}
	for len(nodes) > 0 {
		next := nodes[0]
		nodes[0] = nil
		nodes = nodes[1:]

		switch line := next.(type) {
		case *ast.ExpressionBody:
			// cases like (a, b) => a+b
			// or (a, b) => ({c: a+b})
			nodes = append(nodes, line.Expression)
		case *ast.BlockStatement:
			// cases like (a, b) => { return {a+b}; }
			// or (a, b) => { return ({a, b}); }
			for _, stmt := range line.List {
				nodes = append(nodes, stmt)
			}
		case *ast.ReturnStatement:
			nodes = append(nodes, line.Argument)
		case *ast.ArrayLiteral:
			cd.HasArrayReturn = true
		case *ast.ObjectLiteral:
			for _, prop := range line.Value {
				switch pt := prop.(type) {
				case *ast.PropertyKeyed:
					switch kt := pt.Key.(type) {
					case *ast.Identifier:
						keyTracking[kt.Name.String()] = struct{}{}
					case *ast.StringLiteral:
						keyTracking[kt.Value.String()] = struct{}{}
					case *ast.NumberLiteral:
						keyTracking[kt.Literal] = struct{}{}
					}
				case *ast.PropertyShort:
					keyTracking[pt.Name.Name.String()] = struct{}{}
				}
			}
		case *ast.IfStatement:
			nodes = append(nodes, line.Consequent, line.Alternate)
		case *ast.SwitchStatement:
			for _, sb := range line.Body {
				for _, cs := range sb.Consequent {
					nodes = append(nodes, cs)
				}
			}
		case *ast.StringLiteral,
			*ast.BooleanLiteral,
			*ast.NullLiteral,
			*ast.NumberLiteral,
			*ast.Identifier,
			*ast.BinaryExpression:
			// result of `() => something`
			cd.HasSingleReturn = true
		case nil,
			*ast.VariableStatement,
			*ast.VariableDeclaration:
			// implicit null returns and other lines we don't care about
			continue
		default:
			// TODO maybe this should be a list of unsupported features for users to consider
			return nil, fmt.Errorf("unsupported type %T", line)
		}
	}

	cd.KeyedReturnFields = make([]string, 0, len(keyTracking))
	for key := range keyTracking {
		cd.KeyedReturnFields = append(cd.KeyedReturnFields, key)
	}

	return &cd, nil
}
