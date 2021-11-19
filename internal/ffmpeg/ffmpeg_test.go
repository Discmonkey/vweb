package ffmpeg

import "testing"

func TestPlayer_Next(t *testing.T) {
	player, err := NewPlayer("/home/max/go/src/vweb/test/data/big_buck_bunny_1080_10s_1mb_h264.mp4")
	if err != nil {
		t.Fatalf(err.Error())
	}

	frame, count, err := player.Next()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if len(frame.Bytes()) <= 0 {
		t.Fatalf("corruped packet")
	}

	frame, count2, err := player.Next()

	if count2 <= count {
		t.Fatalf("out of order count")
	}

}
