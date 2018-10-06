package dir

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/name"
)

var Rule = RuleT{
	Name: name.Rule{MaxLen: 20, Style: "lower_case"},
	Size: sizeRule{MaxEntries: 20},
}

type RuleT struct {
	Name name.Rule
	Size sizeRule
}
type sizeRule struct {
	MaxEntries uint `yaml:"maxEntries"`
}

func Check(path string) {
	name := filepath.Base(path)
	checkName(name, path)
	checkSize(name, path)
}

func checkName(name, path string) {
	if path == `.` || path == `..` || path == `/` {
		return
	}
	Rule.Name.Exec(name, "dir", "dir.name", token.Position{Filename: path})
}

func checkSize(name, path string) {
	count := entriesCount(path)
	if count <= Rule.Size.MaxEntries {
		return
	}
	problems.Add(
		token.Position{Filename: path},
		fmt.Sprintf(`dir %s size: %d entries, limit: %d`, name, count, Rule.Size.MaxEntries),
		`dir.size.maxEntries`,
	)
}

func entriesCount(dir string) uint {
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	var count uint
	for _, name := range names {
		if name[0] != '.' {
			count++
		}
	}
	return count
}
