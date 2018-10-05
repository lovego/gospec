package names

import (
	"fmt"
	"go/token"
	"regexp"
	"strings"

	"github.com/lovego/gospec/problems"
)

var lowercaseUnderscore = regexp.MustCompile(`^(_?[a-z0-9]+)+$`)
var lowercaseDash = regexp.MustCompile(`^(-?[a-z0-9]+)+$`)
var camelcase = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
var lowerCamelCase = regexp.MustCompile(`^[a-z][a-zA-Z0-9]*$`)
var exampleTestCase = regexp.MustCompile(`^Example[_A-Z]`)

type rule struct {
	Style  string
	MaxLen int
	Key    string
	Desc   string
}

func (r *rule) Exec(thing, name string, pos token.Position) {
	desc := r.check(name)
	if desc == `` {
		return
	}
	if thing == "" {
		thing = r.Desc
	}
	problems.Add(pos, fmt.Sprintf("%s name %s %s", thing, name, desc), "names."+r.Key)
}

func (r *rule) check(name string) string {
	if name == `_` {
		return ``
	}
	desc := []string{}
	if len(name) > r.MaxLen {
		desc = append(desc, fmt.Sprintf(`%d chars long, limit: %d`, len(name), r.MaxLen))
	}
	if !r.checkStyle(name) {
		desc = append(desc, fmt.Sprintf(`should be %s style`, r.Style))
	}
	return strings.Join(desc, ` and `)
}

func (r *rule) checkStyle(name string) bool {
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
	case `camelCaseInTest`:
		return camelcase.MatchString(name) || exampleTestCase.MatchString(name)
	case `lowerCamelCase`:
		return lowerCamelCase.MatchString(name)
	default:
		panic(fmt.Sprintf(`unknown style config: "%s".`, r.Style))
	}
}

func (r *rule) valid() bool {
	return r.Key != ""
}
