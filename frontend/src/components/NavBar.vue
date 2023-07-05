<template>
  <div class="nav-bar">
    <div class="nav-bar__logo">
      <img src="../assets/logo.png" alt="Vue logo" />
    </div>
    <div class="nav-bar__links">
      <!-- todo: add implementation to mark/hightlight "Profile" and "Chats". Manage the "hasNewsBells" and "hasNewMessages" using backend data through the pinia store methods -->
      <router-link :class="{ highlighted: hasNewBells }" to="/profile">Profile</router-link> |
      <router-link to="/posts" @click="piniaManageData()">Posts</router-link> |
      <router-link :class="{ highlighted: hasNewMessages }" to="/chats">Chats</router-link> |
      <router-link to="/groups">Groups</router-link> |
      <router-link to="/" @click="logout()">Logout</router-link>
    </div>
    <router-view/>
  </div>
  <div v-if="logoutError">logout error: {{ logoutError }}</div>
</template>

<style scoped>
.highlighted {
  background-color: gold;
}
</style>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { ErrorResponse } from '@/api/types';
import { useBellStore } from '@/store/bell';
import { useChatsStore } from '@/store/chats';
import { useGroupStore } from '@/store/group';

import { useUUIDStore } from '@/store/uuid';
import { useLoginStore } from '@/store/login';
import { useSignupStore } from '@/store/signup';
import { useWebSocketStore } from '@/store/websocket';

const logoutError = ref('');

const uuidStore = useUUIDStore();
const loginStore = useLoginStore();
const signupStore = useSignupStore();
const wss = useWebSocketStore();

//todo: reset all pinia stores. Add more later if needed
function resetPiniaStores() {
  uuidStore.$reset();
  loginStore.$reset();
  signupStore.$reset();
  window.dispatchEvent(new Event('beforeunload'));
}

function disconnectWebSocket() {
  wss.disconnect();
}

async function logout() {
  console.log("stage 0")
  try {
    const bodyJson = JSON.stringify(useUUIDStore().getUUID); //todo: perhaps remove/replace later
    const response = await fetch('http://localhost:8080/api/user/logout', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Origin': 'http://localhost:8080'
      },
      body: bodyJson,
      mode: 'cors',
      // credentials: 'omit' // when commented, includes cookie for logout procedure on backend
    });
    console.log("stage 1") //todo: clean up later
    const data = await response.json();
    if (data.error) {
      throw new Error(data.error as string + "problem with json parsing of response");
    }
    console.log("stage 3")



    console.log(data);
    logoutError.value = '';

    disconnectWebSocket();
    resetPiniaStores();

  } catch (error) {
    const errorResponse = error as ErrorResponse;
    logoutError.value = errorResponse.message;
  } finally {
    console.log("stage 4")
  }
}

// when "posts" click happens, reset group id to -1 or 0, to prevent backend filtering of the posts to not show group only posts, but show all
const groupStore = useGroupStore();
function piniaManageData() {
  groupStore.setGroupId(-1); //todo: implement on backend. Now posts will be not filtered by group id. But only filtered by date, fresh first
}

const bellStore = useBellStore();
const hasNewBells = computed(() => bellStore.bells.length > 0);

const chatsStore = useChatsStore();
const hasNewMessages = computed(() => chatsStore.hasNewMessages);
</script>
