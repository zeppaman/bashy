package model

type Parameter struct {
	Name     string
	Desc     string
	Type     string
	Default  string
	Required bool
}

type Script struct {
	Cmds        []string // map[string][]string
	Cmd         string
	Name        string
	ArgUsage    string
	Params      []Parameter
	Description string
	Script string
}
