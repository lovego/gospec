package name

import (
	"fmt"
	"go/token"
	"regexp"
	"strings"

	problemsPkg "github.com/lovego/gospec/problems"
)

var lowercaseUnderscore = regexp.MustCompile(`^(_?[a-z0-9]+)+$`)
var lowercaseDash = regexp.MustCompile(`^(-?[a-z0-9]+)+$`)
var camelcase = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
var lowerCamelCase = regexp.MustCompile(`^[a-z][a-zA-Z0-9]*$`)

type Rule struct {
	MaxLen uint   `yaml:"maxLen"` // max length of a name
	Style  string `yaml:"style"`  // style of a name
}

func (r Rule) Exec(name, thing, key string, pos token.Position) {
	problems, rules := r.check(name)
	if problems == "" {
		return
	}
	problemsPkg.Add(pos, fmt.Sprintf("%s name %s %s", thing, name, problems), key+"."+rules)
}

func (r Rule) check(name string) (string, string) {
	if name == "_" {
		return "", ""
	}
	problems, rules := make([]string, 0, 2), make([]string, 0, 2)
	if uint(len(name)) > r.MaxLen {
		problems = append(problems, fmt.Sprintf("%d chars long, limit: %d", len(name), r.MaxLen))
		rules = append(rules, "maxLen")
	}
	if !r.checkStyle(name) {
		problems = append(problems, fmt.Sprintf("should be %s style", r.Style))
		rules = append(rules, "style")
	}
	rulesStr := strings.Join(rules, ", ")
	if len(rules) > 1 {
		rulesStr = "{" + rulesStr + "}"
	}
	return strings.Join(problems, " and "), rulesStr
}

func (r Rule) checkStyle(name string) bool {
	if r.Style == `` {
		return true
	}
	switch r.Style {
	case `lower_case`:
		return lowercaseUnderscore.MatchString(name)
	case `lower-case`:
		return lowercaseDash.MatchString(name)
	case `camelCase`:
		return camelcase.MatchString(name)
	case `lowerCamelCase`:
		return lowerCamelCase.MatchString(name)
	default:
		panic(fmt.Sprintf(`unknown style config: "%s".`, r.Style))
	}
}

func (r Rule) valid() bool {
	return r.Style != ""
}
