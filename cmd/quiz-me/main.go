package main

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	quizme "github.com/JDinABox/quiz-me"
)

func main() {
	configs := []quizme.Option{}

	if listenAddr := getEnvTrim("LISTEN"); listenAddr != "" {
		configs = append(configs, quizme.WithListenAddr(listenAddr))
	}
	if loggingEnv := getEnvTrim("LOGGING"); loggingEnv != "" {
		logging, err := strconv.ParseBool(loggingEnv)
		if err != nil {
			slog.Error("failed to parse LOGGING environment variable", "error", err)
			os.Exit(1)
		}
		configs = append(configs, quizme.WithLogging(logging))
	}

	if err := quizme.Start(configs...); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}

func getEnvTrim(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}
