<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-text-field
          v-model="videoUrl"
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
        <video-player v-if="playing" :url="videoUrl"></video-player>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
  import Vue from 'vue'
  import VideoPlayer from "@/components/VideoPlayer.vue";

  export default Vue.extend({
    name: 'HelloWorld',
    components: {VideoPlayer},
    computed: {
      playing() {
        return this.videoUrl !== "" && this.playRequested;
      }
    },
    data: () => ({
      videoUrl: "",
      playRequested: false,
      ecosystem: [
        {
          text: 'vuetify-loader',
          href: 'https://github.com/vuetifyjs/vuetify-loader',
        },
        {
          text: 'github',
          href: 'https://github.com/vuetifyjs/vuetify',
        },
        {
          text: 'awesome-vuetify',
          href: 'https://github.com/vuetifyjs/awesome-vuetify',
        },
      ],
      importantLinks: [
        {
          text: 'Documentation',
          href: 'https://vuetifyjs.com',
        },
        {
          text: 'Chat',
          href: 'https://community.vuetifyjs.com',
        },
        {
          text: 'Made with Vuetify',
          href: 'https://madewithvuejs.com/vuetify',
        },
        {
          text: 'Twitter',
          href: 'https://twitter.com/vuetifyjs',
        },
        {
          text: 'Articles',
          href: 'https://medium.com/vuetify',
        },
      ],
      whatsNext: [
        {
          text: 'Explore components',
          href: 'https://vuetifyjs.com/components/api-explorer',
        },
        {
          text: 'Select a layout',
          href: 'https://vuetifyjs.com/getting-started/pre-made-layouts',
        },
        {
          text: 'Frequently Asked Questions',
          href: 'https://vuetifyjs.com/getting-started/frequently-asked-questions',
        },
      ],
    }),
  })
</script>
