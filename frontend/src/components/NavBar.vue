<template>
  <div>
    <div class="logo">
      <h1> Royal Resort</h1>
      <img src="../assets/logo.png" alt="logo" />
    </div>
    <div>
      <router-link
        v-if="wss.bellsList.length < 1"
        to="/bell"
        @click="wss.facepalm()"
      >
        <div class="router_link_box">
          Express Royal Will
        </div>
      </router-link>
      <router-link
        v-else
        to="/bell"
        @click="wss.facepalm()"
        :class="{ 'fade-in': showLink, 'fade-out': !showLink }"
      >
        <div class="router_link_box">
          Express Royal Will
        </div>
      </router-link>
      <br>
      <br>
      <br>
      <router-link to="/profile">
        <div class="router_link_box">
          Profile
        </div>
      </router-link>
      <router-link to="/posts">
        <div class="router_link_box">
          Posts
        </div>
      </router-link>
      <router-link to="/chats">
        <div class="router_link_box">
          Chats
        </div>
      </router-link>
      <router-link to="/groups">
        <div class="router_link_box">
          Groups
        </div>
      </router-link>
      <router-link to="/" @click="logout()">
        <div class="router_link_box">
          Logout
        </div>
      </router-link>
    </div>
    <router-view/>
  </div>
  <div v-if="logoutError">logout error: {{ logoutError }}</div>
</template>

<style scoped>
.fade-in {
  animation: fade-in 0.5s ease-in;
}

.fade-out {
  animation: fade-out 0.5s ease-out;
}

@keyframes fade-in {
  from {
    opacity: 0;
  }

  to {
    opacity: 1;
  }
}

@keyframes fade-out {
  from {
    opacity: 1;
  }

  to {
    opacity: 0;
  }
}
</style>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { ErrorResponse, BellRequest, WSMessageType } from '@/api/types';

import { useUUIDStore } from '@/store/uuid';
import { useLoginStore } from '@/store/login';
import { useSignupStore } from '@/store/signup';
import { useWebSocketStore } from '@/store/websocket';

const logoutError = ref('');

const UUIDStore = useUUIDStore();
const loginStore = useLoginStore();
const signupStore = useSignupStore();
const wss = useWebSocketStore();

//todo: reset all pinia stores. Add more later if needed
function resetPiniaStores() {
  UUIDStore.$reset();
  loginStore.$reset();
  signupStore.$reset();
  window.dispatchEvent(new Event('beforeunload'));
}

function disconnectWebSocket() {
  wss.disconnect();
}

async function logout() {
  console.log("= logout =") //todo: remove debug
  try {
    const bodyJson = JSON.stringify(useUUIDStore().getUUID);
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
    const data = await response.json();
    if (data.error) {
      throw new Error(data.error as string + "problem with json parsing of response");
    }

    console.log(data);
    logoutError.value = '';

    clearInterval(updateInterval.value);
    disconnectWebSocket();
    resetPiniaStores();

  } catch (error) {
    const errorResponse = error as ErrorResponse;
    logoutError.value = errorResponse.message;
  } finally {
    console.log("= logout = 'finally' fired") //todo: remove debug
  }
}

//fade in/out effect for link
let showLink = ref(true);

const updateInterval = ref(0);

onMounted(async () => {
  wss.refresh_websocket()
  await wss.waitForConnection();

  // todo: UNCOMMENT to monitor new notifications
  setInterval(() => {
    showLink.value = !showLink.value;
  }, 1000); // Adjust the interval duration as needed

  updateInterval.value = setInterval(() => {
    updateBells(); // Call the update function
  }, 10000); // Repeat every 10 seconds

  // Clear the interval when the component is unmounted
  onUnmounted(() => {
    clearInterval(updateInterval.value);
  });

  updateBells(); // after success login call the update function once
});

function updateBells() {
  // add x4 cases for each type of bell/notification
  wss.sendMessage({
    type: WSMessageType.FOLLOW_REQUESTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
    } as BellRequest,
  })
  wss.sendMessage({
    type: WSMessageType.GROUP_REQUESTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
    } as BellRequest,
  })
  wss.sendMessage({
    type: WSMessageType.GROUP_INVITES_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
    } as BellRequest,
  })
  wss.sendMessage({
    type: WSMessageType.USER_GROUPS_FRESH_EVENTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
    } as BellRequest,
  })
}

</script>
