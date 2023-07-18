package config

type Handler struct {
	Slack Slack `json:"slack"`
}

type Slack struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
	Title   string `json:"title"`
}

type Resource struct {
	Deployment  bool `json:"deployment"`
	DaemonSet   bool `json:"daemonset"`
	StatefulSet bool `json:"statefulset"`
}

type Logging struct {
	Level string `json:"level"`
}

type Config struct {
	Handler   Handler  `json:"handler"`
	Resource  Resource `json:"resource"`
	Logging   Logging  `json:"logging"`
	Namespace string   `json:"namespace,omitempty"`
}

func New() *Config {
	c := &Config{}

	return c
}
