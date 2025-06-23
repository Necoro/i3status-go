package parser

//go:generate pigeon -o config.pigeon.go -optimize-parser config.pigeon.peg

type Config struct {
	Global   Params
	Sections []Section
}

type Params map[string]string

type Section struct {
	Name      string
	Qualifier string
	Params    Params
}

func (s Section) FullName() string {
	return s.Name + "." + s.Qualifier
}
