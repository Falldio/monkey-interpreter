package evaluator

import (
	"monkey-interpreter/pkg/ast"
	"monkey-interpreter/pkg/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, stmt := range program.Statements {
		if isMacroDefinition(stmt) {
			addMacro(stmt, env)
			definitions = append(definitions, i)
		}
	}

	// Remove macro definitions from the program AST
	for i := len(definitions) - 1; i >= 0; i-- {
		defIdx := definitions[i]
		program.Statements = append(program.Statements[:defIdx], program.Statements[defIdx+1:]...)
	}
}

func isMacroDefinition(stmt ast.Statement) bool {
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStmt.Value.(*ast.MacroLiteral)
	return ok
}

func addMacro(stmt ast.Statement, env *object.Environment) {
	letStmt := stmt.(*ast.LetStatement)
	macro := letStmt.Value.(*ast.MacroLiteral)
	obj := &object.Macro{Parameters: macro.Parameters, Env: env, Body: macro.Body}
	env.Set(letStmt.Name.Value, obj)
}

func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		macro, ok := isMacroCall(call, env)
		if !ok {
			return node
		}

		args := quoteArgs(call)
		evalEnv := extendMacroEnv(macro, args)

		evaluated := Eval(macro.Body, evalEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST nodes from macros")
		}

		return quote.Node
	})
}

func isMacroCall(exp *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	ident, ok := exp.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}

	obj, ok := env.Get(ident.Value)
	if !ok {
		return nil, false
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}

	return macro, true
}

func quoteArgs(exp *ast.CallExpression) []*object.Quote {
	args := []*object.Quote{}

	for _, arg := range exp.Arguments {
		args = append(args, &object.Quote{Node: arg})
	}

	return args
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	env := object.NewEnclosedEnvironment(macro.Env)

	for paramIdx, param := range macro.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}
