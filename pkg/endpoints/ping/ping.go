package ping

import "net/http"

func Ping(res http.ResponseWriter, req *http.Request) {
	_, _ = res.Write([]byte("pong"))
}
