package utils

import (
	"fmt"
	"testing"
)

func TestPortFromString(t *testing.T) {
	p, err := Port("0.0.0.0:9000")
	if err != nil {
		t.Fatal(err)
	}

	if p != 9000 {
		t.Fatal(p)
	}
}

func TestRandomPort(t *testing.T) {
	_, p, err := NewRandomUdpConn()
	if err != nil {
		t.Fatal(err)
	}

	if p <= 200 {
		t.Fatal(p)
	}
	fmt.Println(p)
}
