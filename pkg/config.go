package pkg

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Port           int
	Delay          int
	VcgencmdBinary string
}

func LoadConfig() *Config {
	port := flag.Int("port", 2113, "The port that the health server will bind to")
	delay := flag.Int("delay", 10, "Delay, in seconds, between vcgencmd readings")
	vcgencmdBinary := flag.String("vcgencmd_binary", "/usr/bin/vcgencmd", "Location of the vcgencmd binary")

	flag.Parse()

	return &Config{
		Port:           intValue("PORT", *port),
		Delay:          intValue("DELAY", *delay),
		VcgencmdBinary: stringValue("VCGENCMD_BINARY", *vcgencmdBinary),
	}
}

func stringValue(name string, fallback string) string {
	val := os.Getenv(name)
	if len(val) > 0 {
		return val
	}
	return fallback
}

func intValue(name string, fallback int) int {
	val := os.Getenv(name)
	if len(val) > 0 {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal >= 0 {
			return parsedVal
		}
	}
	return fallback
}
