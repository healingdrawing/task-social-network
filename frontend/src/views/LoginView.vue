<template>
  <div class="login">
    <img alt="Vue logo" src="../assets/logo.png">
    <div> "LoginView.vue . Welcome to Your Vue.js + TypeScript App"</div>
  </div>
  <div>
    <h1>Log in</h1>

    <div><hr><button type="button" @click="crap" title="remove in production">Fill Debug / remove later</button><hr></div> <!-- todo: remove later -->

    <form @submit.prevent="login"> <!-- login means the function name will be call after the default behavior was prevented -->
      <div>
        <label for="email">Email:</label>
        <input type="email" id="email" v-model="email" required>
      </div>
      <div>
        <label for="password">Password:</label>
        <input type="password" id="password" v-model="password" required>
      </div>
      <button type="submit">Submit</button>
    </form>
  </div>
</template>

<script lang="ts" setup>
import router from '@/router/index' /* this works.
Looks like once created router(inside "router/index.ts"),
must be used everywhere
*/

import { ref } from 'vue';
import { useLoginStore } from '@/store/login';
import { useUUIDStore } from '@/store/uuid';
import { useProfileStore } from '@/store/profile';
import { useWebSocketStore } from '@/store/websocket';

const store = useLoginStore();
const storeUUID = useUUIDStore();
const profileStore = useProfileStore();
const webSocketStore = useWebSocketStore();

const email = ref('');
const password = ref('');

const login = async () => {

  try {
    /* todo: should happen only if signup is successful */
    await store.fetchData({
      email: email.value,
      password: password.value,
    });

    if (store.getData.UUID === undefined) {
      store.error = "Error: UUID is undefined. Signup failed.";
      throw new Error(store.error);
    } else {
      console.log("UUID: " + storeUUID.getUUID);
      storeUUID.setUUID(store.getData.UUID)
      webSocketStore.connect(storeUUID.getUUID);
      router.push('/profile');
    }

  } catch (error) {
    console.error(error);
  } finally {
    console.log('finally');
  }


  /* todo: also manual open any routes from browser url window,
  should logout the user, and redirect to login page again
  to prevent any not desirable experimenting, and decrease
  chances of pretension in audit process, from experimentators */

}

const crap = () => {
  email.value = 'dummy@mail.com'
  password.value = '123456'
}
</script>