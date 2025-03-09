package entity

type Env struct {
	Name  string
	Value string
}

func (e Env) String() string {
	return e.Name + "=" + e.Value
}

type EnvEntry struct {
	Num       int
	GroupName string
	Envs      []Env
}

type EnvGroup struct {
	Num       int
	Name      string
	Note      string
	Entrys    []EnvEntry
	SubGroups []EnvGroup
}
