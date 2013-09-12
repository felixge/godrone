package godrone

var DefaultConfig = Config{
	LogLevel:      "debug",
	LogTimeFormat: "15:04:05.999999",
}

type Config struct {
	LogLevel      string
	LogTimeFormat string
}
