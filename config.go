package spec

type configStruct struct {
	Size sizesConfig
	Name namesConfig
}

type sizesConfig struct {
	Dir, File, Row, Func uint
}

type namesConfig struct {
	Dir, File, Pkg, Func, Const, Var, LocalConst, LocalVar nameConfig
}

type nameConfig struct {
	Style  string
	MaxLen uint8
}

var config = configStruct{
	Size: sizesConfig{Dir: 20, File: 200, Row: 100, Func: 20},
	Name: namesConfig{
		Dir:        nameConfig{Style: `lower`, MaxLen: 10},
		File:       nameConfig{Style: `lower`, MaxLen: 10},
		Pkg:        nameConfig{Style: `lower`, MaxLen: 10},
		Func:       nameConfig{Style: `camel`, MaxLen: 20},
		Const:      nameConfig{Style: `camel`, MaxLen: 20},
		Var:        nameConfig{Style: `camel`, MaxLen: 20},
		LocalConst: nameConfig{Style: `camel`, MaxLen: 10},
		LocalVar:   nameConfig{Style: `camel`, MaxLen: 10},
	},
}
