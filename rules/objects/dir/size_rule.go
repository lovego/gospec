package dirpkg

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"

	"github.com/lovego/gospec/problems"
)

type sizeRule struct {
	MaxEntries uint `yaml:"maxEntries"`
}

func (r sizeRule) check(path, key string) {
	count := entriesCount(path)
	if count <= r.MaxEntries {
		return
	}
	problems.Add(
		token.Position{Filename: path},
		fmt.Sprintf(`dir %s size: %d entries, limit: %d`, filepath.Base(path), count, r.MaxEntries),
		key+`.maxEntries`,
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
