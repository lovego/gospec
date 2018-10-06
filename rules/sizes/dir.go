package sizes

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"

	"github.com/lovego/gospec/problems"
)

func checkDir(dir string) {
	if dir == `` {
		return
	}
	count := entriesCount(dir)
	if count <= Rules.Dir {
		return
	}
	problems.Add(
		token.Position{Filename: dir},
		fmt.Sprintf(`dir %s size: %d entries, limit: %d`, filepath.Base(dir), count, Rules.Dir),
		`sizes.dir`,
	)
}

func entriesCount(dir string) int {
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	var count int
	for _, name := range names {
		if name[0] != '.' {
			count++
		}
	}
	return count
}
