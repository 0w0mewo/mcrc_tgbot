package utils

import (
	"sync"

	"github.com/sirupsen/logrus"
)

// various global default utils instance
var once sync.Once
var logger *logrus.Logger
var eventHub *EventHub
var scheduledTasks *ScheduledTaskGroup

func init() {
	once.Do(func() {
		logger = logrus.New()
		logger.Formatter = &logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		}

		eventHub = NewEventHub("default")

		scheduledTasks = NewScheduledTaskGroup("default")
	})

}

// globally available logger
func GetLogger() *logrus.Logger {
	return logger

}

// globally available event hub
func GetDefaultEventHub() *EventHub {
	return eventHub
}

// globally available scheduler
func GetDefaultScheduledTasksGrp() *ScheduledTaskGroup {
	return scheduledTasks
}
