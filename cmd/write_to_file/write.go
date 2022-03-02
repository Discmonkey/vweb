package main

import (
	"fmt"
	"github.com/discmonkey/vweb/internal/utils"
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

	timer := utils.FpsTimer{}
	timer.Start()
	if err != nil {
		return
	}
	for i := 0; ; i++ {
		f := <-out
		bytes, err := f.Bytes()
		fmt.Println(len(bytes), timer.Tick())
		if err != nil {
			return
		}

		_, _ = file.Write(bytes)
		if i > 1000 {
			cancelF()
			break
		}
	}

}
