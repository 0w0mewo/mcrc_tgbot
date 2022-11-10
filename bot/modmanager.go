package bot

import (
	"errors"
	"sync"

	"github.com/0w0mewo/mcrc_tgbot/config"
	"gopkg.in/telebot.v3"
)

var ErrUnregisterMod = errors.New("mod is not registered")

var TgModRegister = newTgModMan()
var DcModRegister = newDcModMan()

type BotType interface {
	TelegramBot | DiscordBot
}

// bot module
type BotModule[T BotType] interface {
	Name() string
	Start(b *T)
	Stop(b *T)
	Reload()
}

type modMan[T BotType] struct {
	lock       *sync.RWMutex
	mods       map[string]BotModule[T]
	tghandlers map[string][]telebot.HandlerFunc // registered handlers
}

// create a discord bot module manager
func newDcModMan() *modMan[DiscordBot] {
	ret := &modMan[DiscordBot]{
		mods: make(map[string]BotModule[DiscordBot]),
		lock: &sync.RWMutex{},
	}

	// reload modules when config file changed
	go func() {
		for range config.ConfigChanged() {
			ret.ReloadModules()
		}
	}()

	return ret
}

// create a telegram bot module manager
func newTgModMan() *modMan[TelegramBot] {
	ret := &modMan[TelegramBot]{
		mods:       make(map[string]BotModule[TelegramBot]),
		lock:       &sync.RWMutex{},
		tghandlers: make(map[string][]telebot.HandlerFunc),
	}

	// reload modules when config file changed
	go func() {
		for range config.ConfigChanged() {
			ret.ReloadModules()
		}
	}()

	return ret
}

func (mm *modMan[T]) ReloadModules() {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	for _, mod := range mm.mods {
		mod.Reload()
	}
}

func (mm *modMan[T]) AddTgEventHandler(_type string, handler telebot.HandlerFunc) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	// make sure the event space exist
	if _, exist := mm.tghandlers[_type]; !exist {
		mm.tghandlers[_type] = make([]telebot.HandlerFunc, 0)
		listenTo = append(listenTo, _type)
	}

	if handler == nil {
		handler = defaultHandler
	}

	mm.tghandlers[_type] = append(mm.tghandlers[_type], handler)
}

func (mm *modMan[T]) GetTgEventHandlers(_type string) []telebot.HandlerFunc {
	mm.lock.RLock()
	defer mm.lock.RUnlock()

	// make sure the event handler is not null
	if _, exist := mm.tghandlers[_type]; !exist {
		return []telebot.HandlerFunc{defaultHandler}
	}

	return mm.tghandlers[_type]
}

func (mm *modMan[T]) RegistryMod(mod BotModule[T]) {
	modName := mod.Name()

	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.mods[modName] = mod
}

func (mm *modMan[T]) GetModules() []BotModule[T] {
	mm.lock.RLock()
	defer mm.lock.RUnlock()

	mods := make([]BotModule[T], 0, len(mm.mods))

	for _, mod := range mm.mods {
		mods = append(mods, mod)
	}

	return mods
}

func defaultHandler(ctx telebot.Context) error {
	return nil
}
