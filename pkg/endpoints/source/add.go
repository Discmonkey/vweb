package source

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/discmonkey/vweb/pkg/android"
	"github.com/discmonkey/vweb/pkg/swagger"
	"github.com/discmonkey/vweb/pkg/utils"
	"github.com/discmonkey/vweb/pkg/video"
	"log"
	"net/http"
	"sync"
)

type Manager struct {
	m       sync.Mutex
	players map[string]playerContext
}

type playerContext struct {
	player video.Player
	cancel context.CancelFunc
}

func post(m *Manager, res http.ResponseWriter, req *http.Request) {
	var s swagger.Source
	if utils.HttpNotOk(400, res, "decode error", json.NewDecoder(req.Body).Decode(&s)) {
		return
	}

	m.m.Lock()
	defer m.m.Unlock()

	udp, port, err := utils.NewRandomUdpConn()
	if utils.HttpNotOk(500, res, "error assigning udp listener", err) {
		return
	}

	// TODO(max) save the cancel somewhere
	player, cancel := android.NewPlayer(udp)

	m.players[s.Name] = playerContext{
		player: player,
		cancel: func() {
			cancel()
			if err := udp.Close(); err != nil {
				log.Println(err)
			}
		},
	}

	json.NewEncoder(res).Encode(swagger.Address{
		Ip:   "",
		Port: int32(port),
	})

}

func get(m *Manager, w http.ResponseWriter) {
	m.m.Lock()
	defer m.m.Unlock()

	sources := make([]swagger.Source, 0, len(m.players))
	for k, v := range m.players {
		sources = append(sources, swagger.Source{
			Codec: v.player.Type(), Name: k,
		})
	}

	err := json.NewEncoder(w).Encode(sources)
	if err != nil {
		log.Print(err)
	}
}

func Source(m *Manager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			get(m, w)
		case "POST":
			post(m, w, req)
		default:
			http.Error(w,
				fmt.Sprintf("request %s not supported", req.Method),
				404)
		}
	}
}
