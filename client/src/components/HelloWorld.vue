<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-select
          v-model="name"
          :items="videos"
          label="Video URL"
        />
        <v-btn v-if="!playing" color="green" @click="playRequested = true">
          Play
        </v-btn>

        <v-btn v-else color="red" @click="playRequested = false">
          Stop
        </v-btn>
      </v-col>
    </v-row>
    <v-row class="text-center">
      <v-col cols="12">
        <video-player v-if="playing" :name="name"></video-player>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
  import axios from "axios";
  import Vue from 'vue'
  import VideoPlayer from "@/components/VideoPlayer.vue";
  import {Source} from "@/swagger/model/source";
  export default Vue.extend({
    name: 'HelloWorld',
    components: {VideoPlayer},
    computed: {
      playing() {
        return this.name !== "" && this.playRequested;
      }
    },
    async mounted() {
      const videos: [Source]= (await axios.get("/source")).data;
      videos.forEach((item: Source) => {
        // @ts-ignore
        this.videos.push(item.name);
      })

    },
    data: () => ({
      name: "",
      playRequested: false,
      videos: [],
    }),
  })
</script>
