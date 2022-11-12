package utils

import (
	"sync"

	"github.com/agoalofalife/event"
	"github.com/sirupsen/logrus"
)

// various global default utils instance
var once sync.Once
var logger *logrus.Logger
var eventHub *event.Dispatcher
var scheduledTasks *ScheduledTaskGroup

func init() {
	once.Do(func() {
		logger = logrus.New()
		logger.Formatter = &logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		}

		eventHub = event.New()

		scheduledTasks = NewScheduledTaskGroup("default")
	})

}

// globally available logger
func GetLogger() *logrus.Logger {
	return logger

}

// globally available event hub
func GetDefaultEventHub() *event.Dispatcher {
	return eventHub
}

// globally available scheduler
func GetDefaultScheduledTasksGrp() *ScheduledTaskGroup {
	return scheduledTasks
}
