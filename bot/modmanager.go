package bot

import (
	"errors"
	"sync"
	"time"

	constant "github.com/0w0mewo/mcrc_tgbot/const"
	"gopkg.in/telebot.v3"
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
	handlers   map[string][]telebot.HandlerFunc // registered handlers
}

func newModMan() *modMan {
	ret := &modMan{
		mods:       make(map[string]Module),
		lock:       &sync.RWMutex{},
		reloadTick: time.NewTicker(constant.MOD_RELOAD_INTERVAL),
		handlers:   make(map[string][]telebot.HandlerFunc),
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

func (mm *modMan) AddTgEventHandler(_type string, handler telebot.HandlerFunc) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	// make sure the event space exist
	if _, exist := mm.handlers[_type]; !exist {
		mm.handlers[_type] = make([]telebot.HandlerFunc, 0)
		listenTo = append(listenTo, _type)
	}

	if handler == nil {
		mm.handlers[_type] = append(mm.handlers[_type], defaultHandler)
		return
	}

	mm.handlers[_type] = append(mm.handlers[_type], handler)
}

func (mm *modMan) GetTgEventHandlers(_type string) []telebot.HandlerFunc {
	mm.lock.RLock()
	defer mm.lock.RUnlock()

	// make sure the event handler is not null
	if _, exist := mm.handlers[_type]; !exist {
		return []telebot.HandlerFunc{defaultHandler}
	}

	return mm.handlers[_type]
}

func (mm *modMan) RegistryMod(mod Module) {
	modName := mod.Name()

	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.mods[modName] = mod
}

func (mm *modMan) GetModules() []Module {
	mm.lock.RLock()
	defer mm.lock.RUnlock()

	mods := make([]Module, 0, len(mm.mods))

	for _, mod := range mm.mods {
		mods = append(mods, mod)
	}

	return mods
}

func defaultHandler(ctx telebot.Context) error {
	return nil
}
