package names

import (
	"go/ast"
	"go/token"
	"reflect"
	"strings"
)

type rules struct {
	Dir, File                       rule
	Pkg, Func, FuncInTest, Label    rule
	Const, Var, Type                rule
	LocalConst, LocalVar, LocalType rule
	StructField                     rule
}

var Rules = rules{
	Dir:  rule{Style: `lower_case`, MaxLen: 20},
	File: rule{Style: `lower_case`, MaxLen: 20},

	StructField: rule{Style: `camelCase`, MaxLen: 30},

	// rule for ast.Ident
	Pkg:        rule{Style: `lower_case`, MaxLen: 20},
	Func:       rule{Style: `camelCase`, MaxLen: 30},
	FuncInTest: rule{Style: `camelCaseInTest`, MaxLen: 50},

	Const: rule{Style: `camelCase`, MaxLen: 30},
	Var:   rule{Style: `camelCase`, MaxLen: 40},
	Type:  rule{Style: `camelCase`, MaxLen: 30},

	LocalConst: rule{Style: `lowerCamelCase`, MaxLen: 20},
	LocalVar:   rule{Style: `lowerCamelCase`, MaxLen: 30},
	LocalType:  rule{Style: `lowerCamelCase`, MaxLen: 20},
	Label:      rule{Style: `lowerCamelCase`, MaxLen: 20},
}

func init() {
	value := reflect.ValueOf(&Rules).Elem()
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		value.FieldByIndex([]int{i}).FieldByName("Name").SetString(
			lowercaseFirstChar(typ.FieldByIndex([]int{i}).Name),
		)
	}
}

func getRuleForIdent(ident *ast.Ident, local bool, fileSet *token.FileSet) rule {
	switch ident.Obj.Kind {
	case ast.Pkg:
		return Rules.Pkg
	case ast.Fun:
		if strings.HasSuffix(fileSet.Position(ident.Pos()).Filename, "_test.go") {
			return Rules.FuncInTest
		} else {
			return Rules.Func
		}
	case ast.Lbl:
		return Rules.Label
	}
	if local {
		switch ident.Obj.Kind {
		case ast.Con:
			return Rules.LocalConst
		case ast.Var:
			return Rules.LocalVar
		case ast.Typ:
			return Rules.LocalType
		}
	} else {
		switch ident.Obj.Kind {
		case ast.Con:
			return Rules.Const
		case ast.Var:
			return Rules.Var
		case ast.Typ:
			return Rules.Type
		}
	}
	return rule{}
}

func lowercaseFirstChar(s string) string {
	b := []byte(s)
	b[0] += 'a' - 'A'
	return string(b)
}
