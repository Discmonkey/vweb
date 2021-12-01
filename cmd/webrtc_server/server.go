package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/discmonkey/vweb/internal/ffmpeg"
	"os"
	"syscall"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
)

const (
	videoFileName = "/home/max/go/src/vweb/test/data/big_buck_bunny_1080_10s_1mb_h264.mp4"
)

func MustReadStdin() string {
	return "eyJ0eXBlIjoib2ZmZXIiLCJzZHAiOiJ2PTBcclxubz0tIDM4ODQxNjQ4Mjk4MjE1MDc2NDEgMiBJTiBJUDQgMTI3LjAuMC4xXHJcbnM9LVxyXG50PTAgMFxyXG5hPWdyb3VwOkJVTkRMRSAwIDFcclxuYT1leHRtYXAtYWxsb3ctbWl4ZWRcclxuYT1tc2lkLXNlbWFudGljOiBXTVNcclxubT12aWRlbyA0Nzk0NyBVRFAvVExTL1JUUC9TQVZQRiA5NiA5NyA5OCA5OSAxMDAgMTAxIDEwMiAxMjIgMTI3IDEyMSAxMjUgMTA3IDEwOCAxMDkgMzUgMzYgMTIwIDExOSAxMjRcclxuYz1JTiBJUDQgNjQuMTIxLjIxNS43XHJcbmE9cnRjcDo5IElOIElQNCAwLjAuMC4wXHJcbmE9Y2FuZGlkYXRlOjQwNzc1Njc3MjAgMSB1ZHAgMjExMzkzNzE1MSAzZmMxMGExYy02ZWE3LTQzMzgtYTkyNC00OGY5M2ZjOGYwNjYubG9jYWwgNDc5NDcgdHlwIGhvc3QgZ2VuZXJhdGlvbiAwIG5ldHdvcmstY29zdCA5OTlcclxuYT1jYW5kaWRhdGU6ODQyMTYzMDQ5IDEgdWRwIDE2Nzc3Mjk1MzUgNjQuMTIxLjIxNS43IDQ3OTQ3IHR5cCBzcmZseCByYWRkciAwLjAuMC4wIHJwb3J0IDAgZ2VuZXJhdGlvbiAwIG5ldHdvcmstY29zdCA5OTlcclxuYT1pY2UtdWZyYWc6a2lFQlxyXG5hPWljZS1wd2Q6WW9YeTlncG1WYUtLRXhOQnlGak1yQUhyXHJcbmE9aWNlLW9wdGlvbnM6dHJpY2tsZVxyXG5hPWZpbmdlcnByaW50OnNoYS0yNTYgN0Q6RkU6M0M6OTY6NDA6OTg6NTY6NDk6Q0U6NEE6QjI6MkY6RTA6Njk6OEM6RTI6MDA6Qzg6QUY6N0U6Mzc6QzY6QTQ6NzI6RkQ6NTY6RjE6MDY6NkI6NDA6ODI6NUZcclxuYT1zZXR1cDphY3RwYXNzXHJcbmE9bWlkOjBcclxuYT1leHRtYXA6MSB1cm46aWV0ZjpwYXJhbXM6cnRwLWhkcmV4dDp0b2Zmc2V0XHJcbmE9ZXh0bWFwOjIgaHR0cDovL3d3dy53ZWJydGMub3JnL2V4cGVyaW1lbnRzL3J0cC1oZHJleHQvYWJzLXNlbmQtdGltZVxyXG5hPWV4dG1hcDozIHVybjozZ3BwOnZpZGVvLW9yaWVudGF0aW9uXHJcbmE9ZXh0bWFwOjQgaHR0cDovL3d3dy5pZXRmLm9yZy9pZC9kcmFmdC1ob2xtZXItcm1jYXQtdHJhbnNwb3J0LXdpZGUtY2MtZXh0ZW5zaW9ucy0wMVxyXG5hPWV4dG1hcDo1IGh0dHA6Ly93d3cud2VicnRjLm9yZy9leHBlcmltZW50cy9ydHAtaGRyZXh0L3BsYXlvdXQtZGVsYXlcclxuYT1leHRtYXA6NiBodHRwOi8vd3d3LndlYnJ0Yy5vcmcvZXhwZXJpbWVudHMvcnRwLWhkcmV4dC92aWRlby1jb250ZW50LXR5cGVcclxuYT1leHRtYXA6NyBodHRwOi8vd3d3LndlYnJ0Yy5vcmcvZXhwZXJpbWVudHMvcnRwLWhkcmV4dC92aWRlby10aW1pbmdcclxuYT1leHRtYXA6OCBodHRwOi8vd3d3LndlYnJ0Yy5vcmcvZXhwZXJpbWVudHMvcnRwLWhkcmV4dC9jb2xvci1zcGFjZVxyXG5hPWV4dG1hcDo5IHVybjppZXRmOnBhcmFtczpydHAtaGRyZXh0OnNkZXM6bWlkXHJcbmE9ZXh0bWFwOjEwIHVybjppZXRmOnBhcmFtczpydHAtaGRyZXh0OnNkZXM6cnRwLXN0cmVhbS1pZFxyXG5hPWV4dG1hcDoxMSB1cm46aWV0ZjpwYXJhbXM6cnRwLWhkcmV4dDpzZGVzOnJlcGFpcmVkLXJ0cC1zdHJlYW0taWRcclxuYT1zZW5kcmVjdlxyXG5hPW1zaWQ6LSA4NjRhYzRkMi05ZDE0LTQyYmYtYmY0YS00NmNkOWIxYTRiMTFcclxuYT1ydGNwLW11eFxyXG5hPXJ0Y3AtcnNpemVcclxuYT1ydHBtYXA6OTYgVlA4LzkwMDAwXHJcbmE9cnRjcC1mYjo5NiBnb29nLXJlbWJcclxuYT1ydGNwLWZiOjk2IHRyYW5zcG9ydC1jY1xyXG5hPXJ0Y3AtZmI6OTYgY2NtIGZpclxyXG5hPXJ0Y3AtZmI6OTYgbmFja1xyXG5hPXJ0Y3AtZmI6OTYgbmFjayBwbGlcclxuYT1ydHBtYXA6OTcgcnR4LzkwMDAwXHJcbmE9Zm10cDo5NyBhcHQ9OTZcclxuYT1ydHBtYXA6OTggVlA5LzkwMDAwXHJcbmE9cnRjcC1mYjo5OCBnb29nLXJlbWJcclxuYT1ydGNwLWZiOjk4IHRyYW5zcG9ydC1jY1xyXG5hPXJ0Y3AtZmI6OTggY2NtIGZpclxyXG5hPXJ0Y3AtZmI6OTggbmFja1xyXG5hPXJ0Y3AtZmI6OTggbmFjayBwbGlcclxuYT1mbXRwOjk4IHByb2ZpbGUtaWQ9MFxyXG5hPXJ0cG1hcDo5OSBydHgvOTAwMDBcclxuYT1mbXRwOjk5IGFwdD05OFxyXG5hPXJ0cG1hcDoxMDAgVlA5LzkwMDAwXHJcbmE9cnRjcC1mYjoxMDAgZ29vZy1yZW1iXHJcbmE9cnRjcC1mYjoxMDAgdHJhbnNwb3J0LWNjXHJcbmE9cnRjcC1mYjoxMDAgY2NtIGZpclxyXG5hPXJ0Y3AtZmI6MTAwIG5hY2tcclxuYT1ydGNwLWZiOjEwMCBuYWNrIHBsaVxyXG5hPWZtdHA6MTAwIHByb2ZpbGUtaWQ9MlxyXG5hPXJ0cG1hcDoxMDEgcnR4LzkwMDAwXHJcbmE9Zm10cDoxMDEgYXB0PTEwMFxyXG5hPXJ0cG1hcDoxMDIgSDI2NC85MDAwMFxyXG5hPXJ0Y3AtZmI6MTAyIGdvb2ctcmVtYlxyXG5hPXJ0Y3AtZmI6MTAyIHRyYW5zcG9ydC1jY1xyXG5hPXJ0Y3AtZmI6MTAyIGNjbSBmaXJcclxuYT1ydGNwLWZiOjEwMiBuYWNrXHJcbmE9cnRjcC1mYjoxMDIgbmFjayBwbGlcclxuYT1mbXRwOjEwMiBsZXZlbC1hc3ltbWV0cnktYWxsb3dlZD0xO3BhY2tldGl6YXRpb24tbW9kZT0xO3Byb2ZpbGUtbGV2ZWwtaWQ9NDIwMDFmXHJcbmE9cnRwbWFwOjEyMiBydHgvOTAwMDBcclxuYT1mbXRwOjEyMiBhcHQ9MTAyXHJcbmE9cnRwbWFwOjEyNyBIMjY0LzkwMDAwXHJcbmE9cnRjcC1mYjoxMjcgZ29vZy1yZW1iXHJcbmE9cnRjcC1mYjoxMjcgdHJhbnNwb3J0LWNjXHJcbmE9cnRjcC1mYjoxMjcgY2NtIGZpclxyXG5hPXJ0Y3AtZmI6MTI3IG5hY2tcclxuYT1ydGNwLWZiOjEyNyBuYWNrIHBsaVxyXG5hPWZtdHA6MTI3IGxldmVsLWFzeW1tZXRyeS1hbGxvd2VkPTE7cGFja2V0aXphdGlvbi1tb2RlPTA7cHJvZmlsZS1sZXZlbC1pZD00MjAwMWZcclxuYT1ydHBtYXA6MTIxIHJ0eC85MDAwMFxyXG5hPWZtdHA6MTIxIGFwdD0xMjdcclxuYT1ydHBtYXA6MTI1IEgyNjQvOTAwMDBcclxuYT1ydGNwLWZiOjEyNSBnb29nLXJlbWJcclxuYT1ydGNwLWZiOjEyNSB0cmFuc3BvcnQtY2NcclxuYT1ydGNwLWZiOjEyNSBjY20gZmlyXHJcbmE9cnRjcC1mYjoxMjUgbmFja1xyXG5hPXJ0Y3AtZmI6MTI1IG5hY2sgcGxpXHJcbmE9Zm10cDoxMjUgbGV2ZWwtYXN5bW1ldHJ5LWFsbG93ZWQ9MTtwYWNrZXRpemF0aW9uLW1vZGU9MTtwcm9maWxlLWxldmVsLWlkPTQyZTAxZlxyXG5hPXJ0cG1hcDoxMDcgcnR4LzkwMDAwXHJcbmE9Zm10cDoxMDcgYXB0PTEyNVxyXG5hPXJ0cG1hcDoxMDggSDI2NC85MDAwMFxyXG5hPXJ0Y3AtZmI6MTA4IGdvb2ctcmVtYlxyXG5hPXJ0Y3AtZmI6MTA4IHRyYW5zcG9ydC1jY1xyXG5hPXJ0Y3AtZmI6MTA4IGNjbSBmaXJcclxuYT1ydGNwLWZiOjEwOCBuYWNrXHJcbmE9cnRjcC1mYjoxMDggbmFjayBwbGlcclxuYT1mbXRwOjEwOCBsZXZlbC1hc3ltbWV0cnktYWxsb3dlZD0xO3BhY2tldGl6YXRpb24tbW9kZT0wO3Byb2ZpbGUtbGV2ZWwtaWQ9NDJlMDFmXHJcbmE9cnRwbWFwOjEwOSBydHgvOTAwMDBcclxuYT1mbXRwOjEwOSBhcHQ9MTA4XHJcbmE9cnRwbWFwOjM1IEFWMS85MDAwMFxyXG5hPXJ0Y3AtZmI6MzUgZ29vZy1yZW1iXHJcbmE9cnRjcC1mYjozNSB0cmFuc3BvcnQtY2NcclxuYT1ydGNwLWZiOjM1IGNjbSBmaXJcclxuYT1ydGNwLWZiOjM1IG5hY2tcclxuYT1ydGNwLWZiOjM1IG5hY2sgcGxpXHJcbmE9cnRwbWFwOjM2IHJ0eC85MDAwMFxyXG5hPWZtdHA6MzYgYXB0PTM1XHJcbmE9cnRwbWFwOjEyMCByZWQvOTAwMDBcclxuYT1ydHBtYXA6MTE5IHJ0eC85MDAwMFxyXG5hPWZtdHA6MTE5IGFwdD0xMjBcclxuYT1ydHBtYXA6MTI0IHVscGZlYy85MDAwMFxyXG5hPXNzcmMtZ3JvdXA6RklEIDE4OTg5NzI2OSAzMDAxMzI4MjYwXHJcbmE9c3NyYzoxODk4OTcyNjkgY25hbWU6WStUYXR2M0s4NWJhV2tkRFxyXG5hPXNzcmM6MTg5ODk3MjY5IG1zaWQ6LSA4NjRhYzRkMi05ZDE0LTQyYmYtYmY0YS00NmNkOWIxYTRiMTFcclxuYT1zc3JjOjE4OTg5NzI2OSBtc2xhYmVsOi1cclxuYT1zc3JjOjE4OTg5NzI2OSBsYWJlbDo4NjRhYzRkMi05ZDE0LTQyYmYtYmY0YS00NmNkOWIxYTRiMTFcclxuYT1zc3JjOjMwMDEzMjgyNjAgY25hbWU6WStUYXR2M0s4NWJhV2tkRFxyXG5hPXNzcmM6MzAwMTMyODI2MCBtc2lkOi0gODY0YWM0ZDItOWQxNC00MmJmLWJmNGEtNDZjZDliMWE0YjExXHJcbmE9c3NyYzozMDAxMzI4MjYwIG1zbGFiZWw6LVxyXG5hPXNzcmM6MzAwMTMyODI2MCBsYWJlbDo4NjRhYzRkMi05ZDE0LTQyYmYtYmY0YS00NmNkOWIxYTRiMTFcclxubT1hdWRpbyA1MDc4NyBVRFAvVExTL1JUUC9TQVZQRiAxMTEgNjMgMTAzIDEwNCA5IDAgOCAxMDYgMTA1IDEzIDExMCAxMTIgMTEzIDEyNlxyXG5jPUlOIElQNCA2NC4xMjEuMjE1LjdcclxuYT1ydGNwOjkgSU4gSVA0IDAuMC4wLjBcclxuYT1jYW5kaWRhdGU6NDA3NzU2NzcyMCAxIHVkcCAyMTEzOTM3MTUxIDNmYzEwYTFjLTZlYTctNDMzOC1hOTI0LTQ4ZjkzZmM4ZjA2Ni5sb2NhbCA1MDc4NyB0eXAgaG9zdCBnZW5lcmF0aW9uIDAgbmV0d29yay1jb3N0IDk5OVxyXG5hPWNhbmRpZGF0ZTo4NDIxNjMwNDkgMSB1ZHAgMTY3NzcyOTUzNSA2NC4xMjEuMjE1LjcgNTA3ODcgdHlwIHNyZmx4IHJhZGRyIDAuMC4wLjAgcnBvcnQgMCBnZW5lcmF0aW9uIDAgbmV0d29yay1jb3N0IDk5OVxyXG5hPWljZS11ZnJhZzpraUVCXHJcbmE9aWNlLXB3ZDpZb1h5OWdwbVZhS0tFeE5CeUZqTXJBSHJcclxuYT1pY2Utb3B0aW9uczp0cmlja2xlXHJcbmE9ZmluZ2VycHJpbnQ6c2hhLTI1NiA3RDpGRTozQzo5Njo0MDo5ODo1Njo0OTpDRTo0QTpCMjoyRjpFMDo2OTo4QzpFMjowMDpDODpBRjo3RTozNzpDNjpBNDo3MjpGRDo1NjpGMTowNjo2Qjo0MDo4Mjo1RlxyXG5hPXNldHVwOmFjdHBhc3NcclxuYT1taWQ6MVxyXG5hPWV4dG1hcDoxNCB1cm46aWV0ZjpwYXJhbXM6cnRwLWhkcmV4dDpzc3JjLWF1ZGlvLWxldmVsXHJcbmE9ZXh0bWFwOjIgaHR0cDovL3d3dy53ZWJydGMub3JnL2V4cGVyaW1lbnRzL3J0cC1oZHJleHQvYWJzLXNlbmQtdGltZVxyXG5hPWV4dG1hcDo0IGh0dHA6Ly93d3cuaWV0Zi5vcmcvaWQvZHJhZnQtaG9sbWVyLXJtY2F0LXRyYW5zcG9ydC13aWRlLWNjLWV4dGVuc2lvbnMtMDFcclxuYT1leHRtYXA6OSB1cm46aWV0ZjpwYXJhbXM6cnRwLWhkcmV4dDpzZGVzOm1pZFxyXG5hPWV4dG1hcDoxMCB1cm46aWV0ZjpwYXJhbXM6cnRwLWhkcmV4dDpzZGVzOnJ0cC1zdHJlYW0taWRcclxuYT1leHRtYXA6MTEgdXJuOmlldGY6cGFyYW1zOnJ0cC1oZHJleHQ6c2RlczpyZXBhaXJlZC1ydHAtc3RyZWFtLWlkXHJcbmE9c2VuZHJlY3ZcclxuYT1tc2lkOi0gOGUxNTM3NTUtYWU5ZC00YTc4LWFjZDgtNzE1YWZjMTAwYjNjXHJcbmE9cnRjcC1tdXhcclxuYT1ydHBtYXA6MTExIG9wdXMvNDgwMDAvMlxyXG5hPXJ0Y3AtZmI6MTExIHRyYW5zcG9ydC1jY1xyXG5hPWZtdHA6MTExIG1pbnB0aW1lPTEwO3VzZWluYmFuZGZlYz0xXHJcbmE9cnRwbWFwOjYzIHJlZC80ODAwMC8yXHJcbmE9Zm10cDo2MyAxMTEvMTExXHJcbmE9cnRwbWFwOjEwMyBJU0FDLzE2MDAwXHJcbmE9cnRwbWFwOjEwNCBJU0FDLzMyMDAwXHJcbmE9cnRwbWFwOjkgRzcyMi84MDAwXHJcbmE9cnRwbWFwOjAgUENNVS84MDAwXHJcbmE9cnRwbWFwOjggUENNQS84MDAwXHJcbmE9cnRwbWFwOjEwNiBDTi8zMjAwMFxyXG5hPXJ0cG1hcDoxMDUgQ04vMTYwMDBcclxuYT1ydHBtYXA6MTMgQ04vODAwMFxyXG5hPXJ0cG1hcDoxMTAgdGVsZXBob25lLWV2ZW50LzQ4MDAwXHJcbmE9cnRwbWFwOjExMiB0ZWxlcGhvbmUtZXZlbnQvMzIwMDBcclxuYT1ydHBtYXA6MTEzIHRlbGVwaG9uZS1ldmVudC8xNjAwMFxyXG5hPXJ0cG1hcDoxMjYgdGVsZXBob25lLWV2ZW50LzgwMDBcclxuYT1zc3JjOjQzMjk3NzM0NCBjbmFtZTpZK1RhdHYzSzg1YmFXa2REXHJcbmE9c3NyYzo0MzI5NzczNDQgbXNpZDotIDhlMTUzNzU1LWFlOWQtNGE3OC1hY2Q4LTcxNWFmYzEwMGIzY1xyXG5hPXNzcmM6NDMyOTc3MzQ0IG1zbGFiZWw6LVxyXG5hPXNzcmM6NDMyOTc3MzQ0IGxhYmVsOjhlMTUzNzU1LWFlOWQtNGE3OC1hY2Q4LTcxNWFmYzEwMGIzY1xyXG4ifQ=="
}

// Encode encodes the input in base64
// It can optionally zip the input before encoding
func Encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}

// Decode decodes the input from base64
// It can optionally unzip the input after decoding
func Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Assert that we have an audio or video file
	_, err := os.Stat(videoFileName)
	if os.IsNotExist(err) {
		fmt.Println("file not found: ", videoFileName)
		syscall.Exit(1)
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

	iceConnectedCtx, iceConnectedCtxCancel := context.WithCancel(context.Background())

	// Create a video track
	videoTrack, videoTrackErr :=
		webrtc.NewTrackLocalStaticSample(
			webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264},
			"video", "h264test")

	if videoTrackErr != nil {
		panic(videoTrackErr)
	}

	rtpSender, videoTrackErr := peerConnection.AddTrack(videoTrack)
	if videoTrackErr != nil {
		panic(videoTrackErr)
	}

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	go func() {
		player, err := ffmpeg.NewPlayer(videoFileName)
		if err != nil {
			fmt.Println("ffmpeg open error", err)
			syscall.Exit(1)
		}
		// Wait for connection established
		<-iceConnectedCtx.Done()

		sleepTime := time.Millisecond * time.Duration(30)
		for {
			frame, _, err := player.Next()
			if err != nil {
				panic(err)
			}

			time.Sleep(sleepTime)
			err = videoTrack.WriteSample(
				media.Sample{Data: frame.Bytes(), Duration: time.Second})

			if err != nil {
				panic(err)
			}
		}
	}()

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
		if connectionState == webrtc.ICEConnectionStateConnected {
			iceConnectedCtxCancel()
		}
	})

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity
			// for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are
			// interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("Peer Connection has gone to failed exiting")
			os.Exit(0)
		}
	})

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	Decode(MustReadStdin(), &offer)

	// Set the remote SessionDescription
	if err = peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	fmt.Println(Encode(*peerConnection.LocalDescription()))

	// Block forever
	select {}
}
