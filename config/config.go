package config

type AppConfig struct {
	Neo4jConfig *Neo4jConfig `yaml:"neo"`
	AuthConfig  *AuthConfig `yaml:"auth"`
}

type Neo4jConfig struct {
	ServerURL string `yaml:"serverURL"`
	Port      string `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type AuthConfig struct {
	JwtToken  string `yaml:"jwt"`
}