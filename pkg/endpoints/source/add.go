package source

import (
	"net/http"
	"sync"
)

type Manager struct {
	m sync.Mutex
}

type Source struct {
}

func Add(http.ResponseWriter, *http.Request) {

}

func Get(w http.ResponseWriter, r *http.Request) {

}
