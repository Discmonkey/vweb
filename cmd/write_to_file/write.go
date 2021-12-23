package main

import (
	"fmt"
	"github.com/discmonkey/vweb/internal/ffmpeg"
	"os"
)

func main() {
	reader, err := ffmpeg.NewPlayer("/home/max/go/src/vweb/test/data/output.ts")
	if err != nil {
		return
	}

	f, err := os.Create("out.ts")
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	if err != nil {
		return
	}
	for {
		frame, _, err := reader.Next()
		if err != nil {
			fmt.Println(err)
			break
		} else if frame.IsKey() {
			_, _ = f.Write(frame.Bytes())
		} else {
			bytes := frame.Bytes()
			_, _ = f.Write(bytes)
		}
	}

}
