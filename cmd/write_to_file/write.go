package main

import (
	"github.com/discmonkey/vweb/internal/nal"
	"github.com/discmonkey/vweb/pkg/android"
	"os"
)

func main() {
	reader, cancelF, err := android.NewPlayer(9000)
	if err != nil {
		return
	}

	out, _, err := reader.Play()
	if err != nil {
		return
	}

	file, err := os.Create("out.ts")
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		return
	}
	for i := 0; i < 100; i++ {
		f := <-out
		bytes, err := f.Bytes()
		if err != nil {
			return
		}
		_, _ = file.WriteString(nal.ToString(nal.Type(bytes[0])) + "\n")
		if i > 1000 {
			cancelF()
			break
		}
	}

}
