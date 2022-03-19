package models

type Config struct {
	Name    string `yaml:"name"`
	Base    string `yaml:"base"`
	Repo    string `yaml:"repo"`
	Context struct {
		ServiceName string   `yaml:"serviceName"`
		ProductName string   `yaml:"productName"`
		Port        int      `yaml:"port"`
		Volumes     []string `yaml:"volumes"`
	}
}

type Manifest map[string]string
