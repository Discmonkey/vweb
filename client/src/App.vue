<template>
  <v-app>
    <v-app-bar
      app
      color="primary"
      dark
    >
      <h3>ReWinder</h3>

    </v-app-bar>

    <v-main>
      <router-view/>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  name: 'App',

  data: () => ({
    //
  }),
});
</script>
