package bot

import (
	"errors"
	"strings"
	"sync"

	"github.com/0w0mewo/mcrc_tgbot/config"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/telebot.v3"
)

var ErrUnregisterMod = errors.New("mod is not registered")

var TgModRegister = newTgModMan()
var DcModRegister = newDcModMan()

type WrappedBotType interface {
	TelegramBot | DiscordBot
}

// bot module
type BotModule[T WrappedBotType] interface {
	Name() string
	Start(b *T)
	Stop(b *T)
	Reload()
}

type modMan[T WrappedBotType] struct {
	tghlock, dchlock *sync.RWMutex
	mods             map[string]BotModule[T]
	tghandlers       map[string][]telebot.HandlerFunc                         // registered telegram handlers
	dchandlers       []func(s *discordgo.Session, m *discordgo.MessageCreate) // registered discord message handlers
}

// create a discord bot module manager
func newDcModMan() *modMan[DiscordBot] {
	ret := &modMan[DiscordBot]{
		mods:       make(map[string]BotModule[DiscordBot]),
		tghlock:    &sync.RWMutex{},
		dchlock:    &sync.RWMutex{},
		dchandlers: make([]func(s *discordgo.Session, m *discordgo.MessageCreate), 0),
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
		tghlock:    &sync.RWMutex{},
		tghandlers: make(map[string][]telebot.HandlerFunc),
		dchlock:    &sync.RWMutex{},
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
	mm.tghlock.Lock()
	defer mm.tghlock.Unlock()

	for _, mod := range mm.mods {
		mod.Reload()
	}
}

func (mm *modMan[T]) AddTgEventHandler(_type string, handler telebot.HandlerFunc) {
	mm.tghlock.Lock()
	defer mm.tghlock.Unlock()

	// make sure the event space exist
	if _, exist := mm.tghandlers[_type]; !exist {
		mm.tghandlers[_type] = make([]telebot.HandlerFunc, 0)
		listenTo = append(listenTo, _type)
	}

	if handler == nil {
		handler = defaultTgHandler
	}

	mm.tghandlers[_type] = append(mm.tghandlers[_type], handler)
}

func (mm *modMan[T]) GetTgEventHandlers(_type string) []telebot.HandlerFunc {
	mm.tghlock.RLock()
	defer mm.tghlock.RUnlock()

	// make sure the event handler is not null
	if _, exist := mm.tghandlers[_type]; !exist {
		return []telebot.HandlerFunc{defaultTgHandler}
	}

	return mm.tghandlers[_type]
}

// add discord command handler
// unlike telebot, discord bot have proper handlers manager
func (mm *modMan[T]) AddDcHandler(handler func(s *discordgo.Session, m *discordgo.MessageCreate)) {
	mm.dchlock.Lock()
	defer mm.dchlock.Unlock()

	mm.dchandlers = append(mm.dchandlers, handler)
}

// get discord command handler
func (mm *modMan[T]) GetDcHandlers() []func(s *discordgo.Session, m *discordgo.MessageCreate) {
	mm.dchlock.RLock()
	defer mm.dchlock.RUnlock()

	return mm.dchandlers
}

func (mm *modMan[T]) RegistryMod(mod BotModule[T]) {
	modName := mod.Name()

	mm.tghlock.Lock()
	defer mm.tghlock.Unlock()

	mm.mods[modName] = mod
}

func (mm *modMan[T]) GetModules() []BotModule[T] {
	mm.tghlock.RLock()
	defer mm.tghlock.RUnlock()

	mods := make([]BotModule[T], 0, len(mm.mods))

	for _, mod := range mm.mods {
		mods = append(mods, mod)
	}

	return mods
}

func defaultTgHandler(ctx telebot.Context) error {
	return nil
}

func WrappedDiscordCmdHandler(cmd string, next DcMsgHandler) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if !strings.Contains(m.Content, ">"+cmd) {
			return
		}

		next(s, m)
	}
}
