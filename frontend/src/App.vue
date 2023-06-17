<template>
   <div id="app">
    <component :is="navComponent"></component>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import NavBar from './components/NavBar.vue'
import NavBarGuest from './components/NavBarGuest.vue'

const route = useRoute()

const navComponent = computed(() => {
  if (route.path === '/' || route.path === '/signup') {
    return NavBarGuest
  } else {
    return NavBar
  }
})
</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
