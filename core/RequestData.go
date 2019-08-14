package schmokin

type RequestData struct {
	Data    []byte
	Type    string            `yaml:"type" json:"type"`
	Method  string            `yaml:"method" json:"method"`
	URL     string            `yaml:"url" json:"url"`
	Headers map[string]string `yaml:"headers"`
	Verify  bool              `yaml:"verify"`
	Pretty  bool              `yaml:"pretty"`
	Before  []string          `yaml:"before"`
	Body    []byte
}
