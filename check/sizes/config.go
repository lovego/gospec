package sizes

import (
	"fmt"
)

type TConfig struct {
	Dir, File, Line, Func int
}

var Config = TConfig{Dir: 20, File: 200, Line: 100, Func: 20}
var rules struct {
	Dir, File, Line, Func [2]string
}

func Setup() {
	rules.Dir = [2]string{`sizes.dir`, fmt.Sprintf(`dir should not have more than %d items`, Config.Dir)}
	rules.File = [2]string{`sizes.file`, fmt.Sprintf(`file should not have more than %d lines`, Config.File)}
	rules.Line = [2]string{`sizes.line`, fmt.Sprintf(`line should not have more than %d bytes`, Config.Line)}
	rules.Func = [2]string{`sizes.func`, fmt.Sprintf(`func should not have more than %d lines`, Config.Func)}
}
