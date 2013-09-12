package godrone

var DefaultConfig = Config{
	LogLevel:      "debug",
	LogTimeFormat: "15:04:05.999999",
	MotorboardTTY: "/dev/ttyO0",
	NavboardTTY:   "/dev/ttyO1",
}

type Config struct {
	LogLevel      string
	LogTimeFormat string
	MotorboardTTY string
	NavboardTTY   string
}
