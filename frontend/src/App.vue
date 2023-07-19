<template>
   <div id="app">
    <component :is="navComponent"></component>
  </div>
</template>

<script lang="ts" setup>
import { computed, onBeforeUnmount } from 'vue'
import { useRoute } from 'vue-router'

import NavBar from './components/NavBar.vue'
import NavBarGuest from './components/NavBarGuest.vue'
import { useWebSocketStore } from './store/websocket';

const wss = useWebSocketStore();

window.onbeforeunload = function () {
  alert("Your prestige is falling!!! See you soon >:<\n... after your kingdom has been conquered o:)")
  console.log("====================================")
  console.log("onbeforeunload event listener triggered")
  console.log("====================================")
  // wss.killThemAll();
}

// Unregister the event listener when the component is unmounted
onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', () => {
    alert("wtf are you doing? Cant you use buttons like a normal person? >:(")
  });
});

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
  // color: #2c3e50;
}

nav {
  padding: 30px;

  // a {
  //   font-weight: bold;
  //   color: #2c3e50;

  //   &.router-link-exact-active {
  //     color: #42b983;
  //   }
  // }
}

/* styles */
:root {
  --gradient: linear-gradient(to right, #b74cc0, #bd08e3 30%, #bd08e3 70%, #b74cc0);
  --goldenGradient: linear-gradient(to bottom, #fde056, #ffffff 30%, #fec722);
  --yellowPink: linear-gradient(to bottom, #fde056, #ffffff 30%, #b74cc0);
}


body {
  background-color: #bd08e3 !important;
  background: var(--gradient);
  background-size: 100% 100vh;
  color: white;
}

.label_file_upload input[type="file"] {
  /* Hide the default file input */
  display: none;
}

.label_file_upload,
button {
  background: var(--goldenGradient);
  color: #b74cc0;
  border: 2px solid white;
  padding: 24px;
  margin: 24px;
  text-decoration: none;
  border-radius: 100%;
  word-wrap: break-word;
  overflow-wrap: break-word;
  display: inline-block;
  cursor: pointer;
}

a {
  text-decoration: none;
  color: #b74cc0;
  cursor: pointer;
}

a.router-link-exact-active {
  color: #003FAE;
}

a:hover {
  color: #511162;
  font-size: 110%;
}

h1,
h2,
h3,
h4,
h5,
h6 {
  text-shadow:
    #511162 0px 0px 10px;
  color: #b74cc0;
}

.single_div_box {
  border-radius: 24px;
  border-left: 2px solid #fde056;
  border-right: 2px solid #fec722;

  word-wrap: break-word;
  overflow-wrap: break-word;
  padding: 24px;
}

.router_link_box {
  background: var(--goldenGradient);
  border: 2px solid white;
  padding: 24px;
  margin-bottom: -24px;
  text-decoration: none;
  border-radius: 100%;
  word-wrap: break-word;
  overflow-wrap: break-word;
  display: inline-block;
}
</style>
