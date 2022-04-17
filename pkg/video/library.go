package video

import (
	"context"
	"errors"
	"github.com/discmonkey/vweb/pkg/swagger"
	"sync"
)

type Library struct {
	m       sync.Mutex
	players map[string]entry
}

func NewLibrary() *Library {
	return &Library{
		m:       sync.Mutex{},
		players: make(map[string]entry),
	}
}

type entry struct {
	player Player
	cancel context.CancelFunc
}

func (l *Library) Add(title string, player Player, callback context.CancelFunc) {
	l.m.Lock()
	defer l.m.Unlock()

	l.players[title] = entry{
		player: player,
		cancel: callback,
	}
}

func (l *Library) Remove(title string) error {
	l.m.Lock()
	l.m.Unlock()
	entry, ok := l.players[title]
	if !ok {
		return errors.New("video with " + title + " does not exist")
	}

	entry.cancel()
	delete(l.players, title)

	return nil
}

func (l *Library) PlayTitle(title string) (chan Frame, context.CancelFunc, error) {
	l.m.Lock()
	defer l.m.Unlock()
	item, okay := l.players[title]
	if !okay {
		return nil, nil, errors.New("source with given name does not exist")
	}

	return item.player.Play()
}

func (l *Library) DescribeTitle(title string) (Type, error) {
	l.m.Lock()
	defer l.m.Unlock()
	if item, okay := l.players[title]; okay {
		return item.player.Type(), nil
	} else {
		return "", errors.New("source with title does not exist")
	}
}

func (l *Library) List() []swagger.Source {
	l.m.Lock()
	defer l.m.Unlock()

	var sources []swagger.Source
	for k, v := range l.players {
		sources = append(sources, swagger.Source{
			Codec: v.player.Type(), Name: k,
		})
	}

	return sources
}
