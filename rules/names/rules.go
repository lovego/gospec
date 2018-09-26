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
	FuncInTest:  rule{Style: `camelCaseInTest`, MaxLen: 50},

	// rule for ast.Ident
	Pkg:  rule{Style: `lower_case`, MaxLen: 20},
	Func: rule{Style: `camelCase`, MaxLen: 30},

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
	for i := 0; i < value.NumField(); i++ {
		key := value.Type().FieldByIndex([]int{i}).Name
		field := value.FieldByIndex([]int{i})
		field.FieldByName("Key").SetString(lowercaseFirstChar(key))
		field.FieldByName("Desc").SetString(camelcaseToLower(key))
	}
	Rules.FuncInTest.Desc = "func"
}

func getRuleForIdent(ident *ast.Ident, local bool, fileSet *token.FileSet) rule {
	switch ident.Obj.Kind {
	case ast.Pkg:
		return Rules.Pkg
	case ast.Fun:
		return Rules.Func
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

/* 单词边界有两种
1. 非大写字符，且下一个是大写字符
2. 大写字符，且下一个是大写字符，且下下一个是非大写字符
*/
func camelcaseToLower(str string) string {
	var slice []string
	start := 0
	for end, char := range str {
		if end+1 < len(str) {
			next := str[end+1]
			if char < 'A' || char > 'Z' {
				if next >= 'A' && next <= 'Z' { // 非大写下一个是大写
					slice = append(slice, str[start:end+1])
					start, end = end+1, end+1
				}
			} else if end+2 < len(str) && (next >= 'A' && next <= 'Z') {
				if next2 := str[end+2]; next2 < 'A' || next2 > 'Z' {
					slice = append(slice, str[start:end+1])
					start, end = end+1, end+1
				}
			}
		} else {
			slice = append(slice, str[start:end+1])
		}
	}
	return strings.ToLower(strings.Join(slice, " "))
}
