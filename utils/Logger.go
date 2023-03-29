/*
 * Copyright (c) Motadata 2022.  All rights reserved
 */

package utils

import (
	"os"
	"strings"
	"time"
)

type Logger struct {
	component string
}

func NewLogger(component string) Logger {

	return Logger{component: component}
}

func (logger *Logger) Trace(message string) {

	logger.format(&message, "TRACE")

	logger.write(&message)
}

func (logger *Logger) Debug(message string) {

	logger.format(&message, "DEBUG")

	logger.write(&message)
}

func (logger *Logger) Info(message string) {

	logger.format(&message, "INFO")

	logger.write(&message)

}

func (logger *Logger) Warn(message string) {

	logger.format(&message, "WARN")

	logger.write(&message)

}

func (logger *Logger) Fatal(message string) {

	logger.format(&message, "FATAL")

	logger.write(&message)

}

func (logger *Logger) Error(message string) {

	logger.format(&message, "ERROR")

	logger.write(&message)

}

func (logger *Logger) format(message *string, level string) {

	currentDate := time.Now().Format("02-January-2006 15")

	currentTime := time.Now().Format("03:04:05.000000 PM")

	*message = currentDate + " " + currentTime + " [" + level + "] " + *message + "\n"

}

func (logger *Logger) write(message *string) {

	logDir := "./logs"

	_, err := os.Stat(logDir)

	if os.IsNotExist(err) {

		os.MkdirAll(logDir, 0755)

	}

	logFile := logDir + "/"

	currentDate := time.Now().Format("02-January-2006 15")

	logFile = logFile + strings.ReplaceAll("@@@-###.log", "@@@", currentDate)

	logFile = strings.ReplaceAll(logFile, "###", logger.component)

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err == nil {

		defer file.Close()

		file.WriteString(*message)
	}

}
