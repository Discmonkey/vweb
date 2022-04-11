package source

import (
	"encoding/json"
	"github.com/discmonkey/vweb/pkg/endpoints/utils"
	"github.com/discmonkey/vweb/pkg/swagger"
	"net/http"
	"sync"
)

type Manager struct {
	m sync.Mutex
}

func post(res http.ResponseWriter, req *http.Request) {
	var s swagger.Source
	if utils.HttpNotOk(400, res, "decode error", json.NewDecoder(req.Body).Decode(&s)) {
		return
	}

}

func get(w http.ResponseWriter, r *http.Request) {

}
