package structs

type Action struct {
	Type   string            `yaml:"type"`
	Params map[string]string `yaml:"params"`
}

type Compare struct {
	Operator string      `yaml:"operator"`
	Value    interface{} `yaml:"value"` // Может быть int/string
}

type Monitor struct {
	Name     string   `yaml:"name"`
	Bash     string   `yaml:"bash"`
	Compare  Compare  `yaml:"compare"`
	Actions  []Action `yaml:"actions"`
	Interval int      `yaml:"interval,omitempty"` // Интервал в секундах (опционально)
}
