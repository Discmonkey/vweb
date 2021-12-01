package main

import (
	"fmt"
	"github.com/discmonkey/vweb/internal/vars"
	"github.com/discmonkey/vweb/pkg/endpoints/open"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir(vars.HttpStaticDir()))
	http.Handle("/", fs)
	http.HandleFunc("/open", open.Open)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", vars.HttpServerPort()), nil); err != nil {
		panic(err)
	}
}
