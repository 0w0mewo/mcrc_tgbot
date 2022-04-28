package commands

import (
	"sync"

	"gopkg.in/telebot.v3"
)

type handlerFunc func(c telebot.Context) error

type cmd struct {
	name    string
	handler handlerFunc
}

type cmdMan struct {
	cmtable map[string]handlerFunc
	lock    *sync.RWMutex
}

func newCmdManager() *cmdMan {
	return &cmdMan{
		cmtable: make(map[string]handlerFunc),
		lock:    &sync.RWMutex{},
	}
}

func (cm *cmdMan) Registry(cmd string, handler handlerFunc) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	cm.cmtable[cmd] = handler
}

func (cm *cmdMan) GetRegTable() map[string]handlerFunc {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	return cm.cmtable
}
