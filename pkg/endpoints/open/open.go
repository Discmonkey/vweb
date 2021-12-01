package open

import (
	"encoding/json"
	"fmt"
	"github.com/discmonkey/vweb/pkg/endpoints/utils"
	"net/http"
)

type Request struct {
	SDP string `json:"sdp"`
	URL string `json:"url"`
}

func Open(w http.ResponseWriter, r *http.Request) {
	var req Request
	if utils.HttpNotOk(404, w, "bad request",
		json.NewDecoder(r.Body).Decode(&req)) {
		return
	}

	fmt.Println(req)
}
