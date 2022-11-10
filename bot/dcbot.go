package bot

import (
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
)

// TODO
type DiscordBot struct {
	logger            *logrus.Entry
	eventHandlersPool *ants.Pool
}
