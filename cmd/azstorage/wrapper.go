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

		logFile, err := openLogFile()
		if err != nil {
			return nil
		}
		defer logFile.Close()

		configureLogger(logFile)

		result := f(cmd, args)
		if result != nil {
			log.Println(result)
		}

		fmt.Println("Log file was saved at: " + logFile.Name())

		return result
	}
}

func openLogFile() (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	logDir := filepath.Join(homeDir, ".azstorage")
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	logPath := filepath.Join(logDir, generateLogFilename())
	return os.OpenFile(logPath, os.O_RDWR|os.O_CREATE, 0644)
}

func generateLogFilename() string {
	t := time.Now()
	return t.Format("20060102T150405.000") + ".log"
}

func configureLogger(logFile *os.File) {
	log.Default().SetOutput(io.MultiWriter(os.Stdout, logFile))
}
