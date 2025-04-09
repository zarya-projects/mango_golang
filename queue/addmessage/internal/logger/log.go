package sl

import (
	"log"
	"log/slog"
	"os"
)

func MustLoad() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
}

func openFile(path string) *os.File {

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Panic(err)
	}

	go func(file *os.File) {
		defer file.Close()
	}(file)

	return file
}
