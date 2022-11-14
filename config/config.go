package config

import (
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type ConfigChangedHandler func()

type ConfigType map[string]any

type configManager struct {
	configwatcher   *fsnotify.Watcher
	cfgFile         string
	cfg             ConfigType         // config file map
	regTable        map[string]IConfig // registered configuration
	rwlock          sync.RWMutex
	changed         chan bool
	changedhandlers []ConfigChangedHandler
}

func newConfigManager(cfgfile string) *configManager {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	err = watcher.Add(cfgfile)
	if err != nil {
		panic(err)
	}

	cm := &configManager{
		configwatcher:   watcher,
		regTable:        make(map[string]IConfig),
		cfg:             make(ConfigType),
		cfgFile:         cfgfile,
		changed:         make(chan bool),
		changedhandlers: make([]ConfigChangedHandler, 0),
	}

	cm.loadConfig()

	// file changed signaling thread
	go func(ev chan fsnotify.Event) {
		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op == fsnotify.Write {
					cm.Reload()
					cm.changed <- true
				}
			}
		}
	}(watcher.Events)

	// thread for handle file changed
	go func() {
		for range cm.changed {
			for _, handler := range cm.changedhandlers {
				if handler != nil {
					handler()
				}

			}
		}
	}()

	return cm
}

func (cm *configManager) ConfigChanged() chan bool {
	return cm.changed
}

func (cm *configManager) OnConfigChanged(cb ConfigChangedHandler) {
	cm.changedhandlers = append(cm.changedhandlers, cb)
}

func (cm *configManager) GetConfigFile() ConfigType {
	return cm.cfg
}

func (cm *configManager) RegisterModuleConfig(module string, c IConfig) {
	cm.rwlock.Lock()
	defer cm.rwlock.Unlock()

	cm.regTable[module] = c
}

func (cm *configManager) GetModuleConfig(module string) IConfig {
	cm.rwlock.RLock()
	defer cm.rwlock.RUnlock()

	return cm.regTable[module]
}

func (cm *configManager) Reload() {
	err := cm.loadConfig()
	if err != nil {
		log.Println(err)
	}

	cm.rwlock.RLock()
	defer cm.rwlock.RUnlock()

	for _, config := range cm.regTable {
		config.Reload()
	}
}

func (cm *configManager) Table() map[string]IConfig {
	return cm.regTable
}

type IConfig interface {
	Reload()
}

func (cm *configManager) loadConfig() error {
	fd, err := os.Open(cm.cfgFile)
	if err != nil {
		return err
	}
	defer fd.Close()

	err = yaml.NewDecoder(fd).Decode(cm.cfg)

	return err
}
