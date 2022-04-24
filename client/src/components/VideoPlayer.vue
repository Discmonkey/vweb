<template>
  <video ref="video" muted autoplay controls></video>
</template>

<script>

import 'axios';
import axios from "axios";

const iceServers = [
  {
    urls: 'stun:stun.l.google.com:19302'
  }
]
export default {
  name: "VideoPlayer",
  props: [name],
  data: () => (
    {
      pc: new RTCPeerConnection({
        iceServers
      })
    }),

  async mounted() {
    this.pc.oniceconnectionstatechange = e => console.log(e, this.pc.iceConnectionState);
    this.pc.addTransceiver('video', {'direction': 'sendrecv'});
    this.pc.addTransceiver('audio', {'direction': 'sendrecv'})
    const offer = await this.pc.createOffer();
    const response = await axios.post("play", {
      stream: {name: this.url},
      sdp: btoa(JSON.stringify(offer))
    });

    this.pc.ontrack = (event) => {
      this.$refs.video.srcObject = event.streams[0];
    }

    await this.pc.setLocalDescription(offer);
    await this.pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(response.data.sdp))))
  }

}
</script>

<style scoped>

</style>