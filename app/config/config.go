package config

type Config struct {
	DBPath       string
	Port         int
	EmployeesURL string
}

func NewConfig() *Config {
	return &Config{
		DBPath:       "employees.db",
		Port:         8080,
		EmployeesURL: "https://gist.githubusercontent.com/chancock09/6d2a5a4436dcd488b8287f3e3e4fc73d/raw/fa47d64c6d5fc860fabd3033a1a4e3c59336324e/employees.json",
	}
} 