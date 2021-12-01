<template>
  <video ref="video"></video>
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
  props: ['url'],
  data: () => (
    {
      pc: new RTCPeerConnection({
        iceServers
      })
    }),
  async mounted() {
    this.pc.oniceconnectionstatechange = e => console.log(e, this.pc.iceConnectionState);
    this.pc.addTransceiver('video', {'direction': 'sendrecv'});

    const offer = await this.pc.createOffer();
    console.log(offer);

    const response = await axios.post("open", {
      url: this.url,
      sdp: offer
    });

    console.log(response);
    await this.pc.setLocalDescription(offer);
    await this.pc.setLocalDescription(response.sdp)

    console.log("mounted");
  }

}
</script>

<style scoped>

</style>