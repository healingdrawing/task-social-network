<template>
  <div class="login">
    <img alt="Vue logo" src="../assets/logo.png">
    <div> "LoginView.vue . Welcome to Your Vue.js + TypeScript App"</div>
  </div>
  <div>
    <h1>Log in</h1>
     <!-- todo: add required keyword into tags input(email and password) later, to prevent empty input. Login uses only email as nickname, because the "Nickname" property is optional in registration, so ignore it in login implementation , because unsafe to use -->
    <form @submit.prevent="login"> <!-- login means the function name will be call after the default behavior was prevented -->
      <div>
        <label for="email">Nickname:</label>
        <input type="email" id="email"> v-model="email">
      </div>
      <div>
        <label for="password">Password:</label>
        <input type="password" id="password"> v-model="password">
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

import { useProfileStore } from '@/store/profile';

const profileStore = useProfileStore();

const login = async () => {
  // backend api call
  const dummyUserId = 88;
  /* todo: shoud happens only if login is successful */
  profileStore.setUserId(dummyUserId); // set once in login process. Perhaps need clear before logout etc.
  router.push('/profile')

  /* todo: also manual open any routes from browser url window,
  should logout the user, and redirect to login page again
  to prevent any not desirable experimenting, and decrease
  chances of pretension in audit process, from experimentators */

}
</script>