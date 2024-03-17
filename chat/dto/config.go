package dto

type (
	Config struct {
		Environment string
		ApiName     string
		ApiVersion  string
		Server      ConfigServer
		Database    ConfigDatabase
	}

	ConfigServer struct {
		Debug   bool
		Port    int
		BaseUrl string
		Log     string
	}

	ConfigDatabase struct {
		Host string
		Port int
		Name string
		Log  string
	}
)
