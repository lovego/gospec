package sizes

import (
	"fmt"
)

type TConfig struct {
	Dir, File, Line, Func int
}

var Config = TConfig{Dir: 20, File: 200, Line: 100, Func: 20}
var descs struct {
	Dir, File, Line, Func string
}

func Setup() {
	descs.Dir = fmt.Sprintf(`dir should not have more than %d items`, Config.Dir)
	descs.File = fmt.Sprintf(`file should not have more than %d lines`, Config.File)
	descs.Line = fmt.Sprintf(`line should not have more than %d bytes`, Config.Line)
	descs.Func = fmt.Sprintf(`file should not have more than %d lines`, Config.Func)
}
