package main

import (
	"fmt"
	"github.com/discmonkey/vweb/internal/vars"
	"github.com/discmonkey/vweb/pkg/android"
	"github.com/discmonkey/vweb/pkg/endpoints/open"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir(vars.HttpStaticDir()))
	player, _, err := android.NewPlayer(9000)
	if err != nil {
		log.Fatalln(err)
	}
	http.Handle("/", fs)
	http.HandleFunc("/open", open.VideoEndpoint(player))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", vars.HttpServerPort()), nil); err != nil {
		panic(err)
	}
}
