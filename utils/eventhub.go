package utils

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type EventCallback func(ctx context.Context) error

type EventHub struct {
	store  *sync.Map
	logger *logrus.Entry
	wg     *sync.WaitGroup
}

func NewEventHub(namespace string) *EventHub {
	return &EventHub{
		store:  &sync.Map{},
		logger: GetLogger().WithField("service", "eventhub-"+namespace),
		wg:     &sync.WaitGroup{},
	}
}

// register an event callback to an event
func (eh *EventHub) Register(ev string, cb EventCallback) {
	_cbs, exist := eh.store.Load(ev)
	if !exist {
		newcbs := make([]EventCallback, 0)
		newcbs = append(newcbs, cb)
		eh.store.Store(ev, newcbs)

		return
	}

	cbs := _cbs.([]EventCallback)
	cbs = append(cbs, cb)

	eh.store.Store(ev, cbs)

}

// publish an event, event callbacks related parameters are passed by context
func (eh *EventHub) Notify(ev string, ctx context.Context) {
	_cbs, exist := eh.store.Load(ev)
	if !exist {
		return
	}

	for _, cb := range _cbs.([]EventCallback) {
		eh.wg.Add(1)
		go func(cb EventCallback) {
			defer eh.wg.Done()
			err := cb(ctx)
			if err != nil {
				eh.logger.Error(err)
			}
		}(cb)

	}
}



// wait for event callbacks finish
func (eh *EventHub) Wait() {
	eh.wg.Wait()
}
