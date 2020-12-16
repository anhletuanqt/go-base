package config

type Config struct {
	Mongo Mong
	Server
}

type Server struct {
	Port string
}

type Mong struct {
	URL    string
	DBName string
}

func New() *Config {
	mongo := Mong{
		URL: "mongodb://localhost:27017",
		// URL:    "mongodb+srv://username:password@cluster0.nhxmu.mongodb.net/test?retryWrites=true&w=majority",
		DBName: "test",
	}
	server := Server{
		Port: "3000",
	}
	return &Config{
		Mongo:  mongo,
		Server: server,
	}
}
