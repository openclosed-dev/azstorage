package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type runCommand func(cmd *cobra.Command, args []string) error

func runWithLog(f runCommand) runCommand {
	return func(cmd *cobra.Command, args []string) error {

		startTime := time.Now()

		logFile, err := openLogFile(&startTime)
		if err != nil {
			return nil
		}
		defer logFile.Close()

		configureLogger(logFile)

		result := f(cmd, args)
		if result != nil {
			log.Println(result)
		}

		elapsedTime := time.Since(startTime)
		log.Printf("Command execution took %s long", elapsedTime)

		fmt.Println("Log file was saved at: " + logFile.Name())

		return result
	}
}

func openLogFile(now *time.Time) (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	logDir := filepath.Join(homeDir, ".azstorage")
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	logPath := filepath.Join(logDir, generateLogFilename(now))
	return os.OpenFile(logPath, os.O_RDWR|os.O_CREATE, 0644)
}

func generateLogFilename(now *time.Time) string {
	return now.Format("20060102T150405.000") + ".log"
}

func configureLogger(logFile *os.File) {
	log.Default().SetOutput(io.MultiWriter(os.Stdout, logFile))
}
