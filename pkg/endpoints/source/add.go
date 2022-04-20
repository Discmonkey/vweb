package source

import (
	"encoding/json"
	"fmt"
	"github.com/discmonkey/vweb/pkg/android"
	"github.com/discmonkey/vweb/pkg/swagger"
	"github.com/discmonkey/vweb/pkg/utils"
	"github.com/discmonkey/vweb/pkg/video"
	"log"
	"net/http"
)

func post(l *video.Library, res http.ResponseWriter, req *http.Request) {
	var s swagger.Source
	if utils.HttpNotOk(400, res, "decode error", json.NewDecoder(req.Body).Decode(&s)) {
		return
	}

	udp, port, err := utils.NewRandomUdpConn()
	if utils.HttpNotOk(500, res, "error assigning udp listener", err) {
		return
	}

	player, cancel := android.NewPlayer(udp)

	// add our video to the library
	l.Add(s.Name, player, func() {
		cancel()
		if err := udp.Close(); err != nil {
			log.Println(err)
		}
	})

	utils.LogIf(json.NewEncoder(res).Encode(swagger.Address{
		Ip:   "",
		Port: int32(port),
	}))
}

func get(l *video.Library, w http.ResponseWriter) {
	utils.LogIf(json.NewEncoder(w).Encode(l.List()))
}

func Source(l *video.Library) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Println("source")
		switch req.Method {
		case "GET":
			get(l, w)
		case "POST":
			post(l, w, req)
		default:
			http.Error(w,
				fmt.Sprintf("request %s not supported", req.Method),
				404)
		}
	}
}
