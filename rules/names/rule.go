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

var hasUppercase = regexp.MustCompile(`[A-Z]+`)
var hasUppercaseOrUnderscore = regexp.MustCompile(`[A-Z_]+`)
var exampleTestCase = regexp.MustCompile(`^Example[_A-Z]`)

type rule struct {
	Style  string
	MaxLen int
	Key    string
	Desc   string
}

func (r *rule) Exec(thing, name string, pos token.Position) {
	desc := r.check(name, false)
	if desc == `` {
		return
	}
	if thing == "" {
		thing = r.Desc
	}
	problems.Add(pos, fmt.Sprintf("%s name %s %s", thing, name, desc), "names."+r.Key)
}

func (r *rule) check(name string, loose bool) string {
	if name == `_` {
		return ``
	}
	desc := []string{}
	if len(name) > r.MaxLen {
		desc = append(desc, fmt.Sprintf(`%d chars long, limits %d`, len(name), r.MaxLen))
	}
	if !r.checkStyle(name, loose) {
		desc = append(desc, fmt.Sprintf(`should be %s style`, r.Style))
	}
	return strings.Join(desc, ` and `)
}

func (r *rule) checkStyle(name string, loose bool) bool {
	if r.Style == `` {
		return true
	}
	switch r.Style {
	case `lower_case`:
		return loose && !hasUppercase.MatchString(name) || lowercaseUnderscore.MatchString(name)
	case `lower-case`:
		return loose && !hasUppercaseOrUnderscore.MatchString(name) ||
			lowercaseDash.MatchString(name)
	case `camelCase`:
		return loose && strings.IndexByte(name, '_') < 0 || camelcase.MatchString(name)
	case `camelCaseInTest`:
		return loose && strings.IndexByte(name, '_') < 0 || camelcase.MatchString(name) ||
			exampleTestCase.MatchString(name)
	case `lowerCamelCase`:
		return loose && strings.IndexByte(name, '_') < 0 && (name[0] < 'a' || name[0] > 'z') ||
			lowerCamelCase.MatchString(name)
	default:
		panic(fmt.Sprintf(`unknown style config: "%s".`, r.Style))
	}
}

func (r *rule) valid() bool {
	return r.Key != ""
}
