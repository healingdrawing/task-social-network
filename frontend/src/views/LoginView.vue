<template>
  <div class="login">
    <img alt="Vue logo" src="../assets/logo.png">
    <div> "LoginView.vue . Welcome to Your Vue.js + TypeScript App"</div>
  </div>
  <div>
    <h1>Log in</h1>
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
// import { useRouter } from 'vue-router' // TODO: remove this not working shit
import router from '@/router/index' /* this works,
some fail with .ts at the end, so it was removed.
Looks like once created router(inside "router/index.ts"),
must be used everywhere
*/

const login = async () => {
  // const router = useRouter() //todo: remove this not working, the result of shit above

  /* shoud happens only if login is successful */
  router.push('/profile')

  /* todo: also manual open any routes from browser url window,
  should logout the user, and redirect to login page again
  to prevent any not desirable experimenting, and decrease
  chances of pretension in audit process, from experimentators */

}
</script>