package spec

type configStruct struct {
	Name configNames
	Size configSizes
}

type configSizes struct {
	Dir, File, Row, Func uint
}

type configNames struct {
	Dir, File, Pkg, Func, Const, Var string
}

type configName struct {
	Style  string
	MaxLen uint8
}
