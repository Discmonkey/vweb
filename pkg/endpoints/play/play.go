package play

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/discmonkey/vweb/pkg/swagger"
	"github.com/discmonkey/vweb/pkg/utils"
	"github.com/discmonkey/vweb/pkg/video"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"net/http"
	"time"
)

// VideoEndpoint TODO(max) refactor and productionalize this method
// VideoEndpoint plays a given media over the network
func VideoEndpoint(l *video.Library) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req, res swagger.Session
		if utils.HttpNotOk(404, w, "bad request",
			json.NewDecoder(r.Body).Decode(&req)) {
			return
		}

		codec, err := l.DescribeTitle(req.Stream.Name)
		if utils.HttpNotOk(404, w, "stream not found", err) {
			return
		}

		peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
		})
		if utils.HttpNotOk(404, w, "could not create peer connection", err) {
			return
		}

		iceConnectedCtx, iceConnectedCtxCancel := context.WithCancel(context.Background())

		videoTrack, err := webrtc.NewTrackLocalStaticSample(
			webrtc.RTPCodecCapability{MimeType: codec, ClockRate: 90000}, "video", "pion")
		if utils.HttpNotOk(500, w, "could not start track", err) {
			iceConnectedCtxCancel()
			return
		}

		rtpSender, err := peerConnection.AddTrack(videoTrack)
		if utils.HttpNotOk(500, w, "could not add track", err) {
			iceConnectedCtxCancel()
			return
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
			stream, _, err := l.PlayTitle(req.Stream.Name)

			if utils.HttpNotOk(400, w, "could not stream contents", err) {
				return
			}

			// Wait for connection established
			<-iceConnectedCtx.Done()

			// Send our video file frame at a time. Pace our sending so we send it at the same speed it should be played back as.
			// This isn't required since the video is timestamped, but we will such much higher loss if we send all at once.
			for frame := range stream {
				bytes, err := frame.Bytes()
				if err != nil {
					fmt.Println(err)
					break
				}

				if err = videoTrack.WriteSample(media.Sample{Data: bytes, Duration: time.Second}); err != nil {
					panic(err)
				}
			}
		}()

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
				// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
				// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
				// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
				fmt.Println("Peer Connection has gone to failed exiting")
			}
		})

		offer := webrtc.SessionDescription{}
		decoded, err := base64.StdEncoding.DecodeString(req.Sdp)
		if utils.HttpNotOk(400, w, "could not decode sdp", err) {
			return
		}

		if utils.HttpNotOk(400, w, "could not parse sdp",
			json.Unmarshal(decoded, &offer)) {
			return
		}

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
		localDescription := *peerConnection.LocalDescription()
		marshalled, err := json.Marshal(localDescription)
		if utils.HttpNotOk(400, w, "could not get local description", err) {
			return
		}
		res.Sdp = base64.StdEncoding.EncodeToString(marshalled)

		utils.LogIf(json.NewEncoder(w).Encode(res))
	}
}
