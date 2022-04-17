package main

import (
	"fmt"
	"github.com/discmonkey/vweb/internal/vars"
	"github.com/discmonkey/vweb/pkg/endpoints/play"
	"github.com/discmonkey/vweb/pkg/endpoints/source"
	"github.com/discmonkey/vweb/pkg/video"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir(vars.HttpStaticDir()))
	l := video.NewLibrary()
	http.Handle("/", fs)
	http.HandleFunc("/source", source.Source(l))
	http.HandleFunc("/play", play.VideoEndpoint(l))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", vars.HttpServerPort()), nil); err != nil {
		panic(err)
	}
}
