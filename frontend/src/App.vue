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
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
}

nav {
  padding: 30px;
}

/* styles */
:root {
  --yellow_dark: #fec722;
  --yellow_light: #fde056;
  --pink_dark: #511162;
  --pink_bright: #bd08e3;
  --pink_light: #b74cc0;
  --gradient: linear-gradient(to right,
      var(--pink_dark),
      var(--pink_bright) 40%,
      var(--pink_bright) 60%,
      var(--pink_dark));
  --g_horizontal_white_glass: linear-gradient(to right,
      hsla(0, 0%, 100%, 0.01),
      #ffffff00 30%,
      #ffffff00 70%,
      hsla(0, 0%, 100%, 0.10));
  --g_button_white_glass: linear-gradient(to right,
      hsla(0, 0%, 100%, 0.10),
      #ffffff00 30%,
      #ffffff00 70%,
      hsla(0, 0%, 100%, 0.10));

  --g_button_hover_white_glass: linear-gradient(to right,
      hsla(0, 0%, 100%, 0.30),
      #ffffff00 30%,
      #ffffff00 70%,
      hsla(0, 0%, 100%, 0.30));

  --g_active_router_link: linear-gradient(to right,
      white,
      var(--pink_dark) 20px,
      var(--pink_dark) calc(100% - 20px),
      white);
}


body {
  background-color: var(--pink_bright);
  background: var(--gradient);
  background-size: 100% 100vh;
  color: var(--yellow_dark);
}

.logo img {
  position: relative;
  margin-top: -50px;
}

:focus {
  outline-offset: -24px;
  outline-color: rgba(255, 255, 255, 1);
}

// headache

.users_list_with_scroll {
  background: var(--g_horizontal_white_glass);
  height: 100px;
  overflow-y: scroll;
  overflow-x: hidden;
  scrollbar-color: white hsla(0, 0%, 100%, 0.5);
  scrollbar-width: thin;
}

.users_list_with_scroll::-webkit-scrollbar {
  width: 2px;
}

.users_list_with_scroll::-webkit-scrollbar-thumb {
  background-color: white;
}

.users_list_with_scroll::-webkit-scrollbar-track {
  background-color: hsla(0, 0%, 100%, 0);
}

textarea {
  color: var(--yellow_dark);
  scrollbar-color: white hsla(0, 0%, 100%, 0);
  scrollbar-width: thin;
}

textarea::-webkit-scrollbar {
  width: 2px;
}

textarea::-webkit-scrollbar-thumb {
  background-color: white;
}

textarea::-webkit-scrollbar-track {
  background-color: hsla(0, 0%, 100%, 0);
  ;
}

// end of headache

.label_file_upload input[type="file"] {
  /* Hide the default file input */
  display: none;
}

.label_file_upload,
button {
  background: var(--g_button_white_glass);
  border-style: solid;
  border-width: 0 2px;
  border-color: white;
  border-radius: 24px;
  color: var(--yellow_light);
  padding: 24px;
  margin: 24px;
  word-wrap: break-word;
  overflow-wrap: break-word;
  display: inline-block;
  cursor: pointer;

}

.label_file_upload:hover,
button:hover {
  background: var(--g_button_hover_white_glass);
  text-shadow: 0px 0px 2px white;
}

a {
  text-decoration: none;
  color: var(--yellow_light);
  cursor: pointer;
}

a.router-link-exact-active {
  background: var(--g_active_router_link);
  // #003FAE
}

a:hover {
  text-shadow: 0px 0px 2px white;
}

h1,
h2,
h3,
h4,
h5,
h6 {
  text-shadow: var(--pink_dark) 0px 0px 10px;
  color: hsl(295, 71%, 63%);
}

.single_div_box {
  border-radius: 24px;
  border-left: 2px solid #fde056;
  border-right: 2px solid #fde056;

  word-wrap: break-word;
  overflow-wrap: break-word;
  padding: 24px;
}

.router_link_box {
  background: var(--g_button_white_glass);
  border-left: 2px solid white;
  border-right: 2px solid white;
  border-radius: 24px;
  padding: 24px;
  margin-bottom: -24px;
  text-decoration: none;
  word-wrap: break-word;
  overflow-wrap: break-word;
  display: inline-block;
}

.router_link_box:hover {
  background: var(--g_button_hover_white_glass);
}

input {
  background: var(--g_horizontal_white_glass);
  color: var(--yellow_dark);
  border-style: solid;
  border-width: 0 2px;
  border-color: var(--yellow_light);
  word-wrap: break-word;
  overflow-wrap: break-word;
}

input:focus {
  outline: none;
  border-color: var(--yellow_dark);
  background: var(--pink_dark);
}

select {
  overflow: hidden;
  border: none;
  color: var(--yellow_light);
}

select:focus {
  outline: none;
  background: var(--pink_dark);
}

textarea {
  background: var(--g_horizontal_white_glass);
  border-style: solid;
  border-width: 0 2px;
  border-color: var(--yellow_light);
  word-wrap: break-word;
  overflow-wrap: break-word;
}

textarea:focus {
  outline: none;
  border-color: var(--yellow_dark);
  background: var(--pink_dark);
}
</style>
