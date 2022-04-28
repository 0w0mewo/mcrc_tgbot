package bot

import (
	"errors"
	"sync"
	"time"

	constant "github.com/0w0mewo/mcrc_tgbot/const"
)

var ErrUnregisterMod = errors.New("mod is not registered")

var ModRegister = newModMan()

type Module interface {
	Name() string
	Start(b *Bot)
	Stop(b *Bot)
	Reload()
}

type modMan struct {
	lock       *sync.RWMutex
	mods       map[string]Module
	reloadTick *time.Ticker
}

func newModMan() *modMan {
	ret := &modMan{
		mods:       make(map[string]Module),
		lock:       &sync.RWMutex{},
		reloadTick: time.NewTicker(constant.MOD_RELOAD_INTERVAL),
	}

	go func(mods map[string]Module) {
		for range ret.reloadTick.C {
			for _, mod := range mods {
				mod.Reload()
			}
		}
	}(ret.mods)

	return ret
}

func (mm *modMan) Registry(mod Module) {
	modName := mod.Name()

	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.mods[modName] = mod
}

func (mm *modMan) Get() []Module {
	mm.lock.RLock()
	defer mm.lock.RUnlock()

	mods := make([]Module, 0, len(mm.mods))

	for _, mod := range mm.mods {
		mods = append(mods, mod)
	}

	return mods
}
