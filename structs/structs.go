package structs

type Action struct {
	Type   string            `yaml:"type"`
	Params map[string]string `yaml:"params"`
}

type Compare struct {
	Operator string      `yaml:"operator"`
	Value    interface{} `yaml:"value"` // Can be either string of int
}

type Monitor struct {
	Name     string   `yaml:"name"`
	Bash     string   `yaml:"bash"` // Bash command
	Compare  Compare  `yaml:"compare"`
	Actions  []Action `yaml:"actions"`
	Interval string   `yaml:"interval,omitempty"` // Interval: format - "1s", 2, "4m" etc
	Timeout  string   `yaml:"timeout,omitempty"`  // Timeout: format - "1s", 2, "4m" etc
	Retries  int      `yaml:"retries,omitempty"`  // Amount of retries
}
