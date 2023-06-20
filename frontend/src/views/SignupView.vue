<template>
  <div class="signup">
    <img alt="Vue logo" src="../assets/logo.png">
    <div>"SignupView.vue . Welcome to Your Vue.js + TypeScript App"</div>
  </div>
  <div>
    <h1>Sign up</h1>
    <form @submit.prevent="signup">
    
    <div><hr><button type="button" @click="crap" title="remove in production">Fill Debug / remove later</button><hr></div> <!-- todo: remove later -->
      
    <div>
      <label for="email">Email:</label>
      <input type="email" id="email" v-model="email" required>
    </div>
    <div>
      <label for="password">Password:</label>
      <input type="password" id="password" v-model="password" required>
    </div>
    <div>
      <label for="firstName">First Name:</label>
      <input type="text" id="firstName" v-model="firstName" required>
    </div>
    <div>
      <label for="lastName">Last Name:</label>
      <input type="text" id="lastName" v-model="lastName" required>
    </div>
    <div>
      <label for="dob">Date of Birth:</label>
      <input type="date" id="dob" v-model="dob" required>
    </div>
    <div>
      <label for="avatar">Avatar:</label>
      <input type="file" id="avatar" accept="image/jpeg, image/png, image/gif" @change="handleAvatarChange">
      <div class="optional">(optional)</div>
    </div>
    <div>
      <label for="nickname">Nickname:</label>
      <input type="text" id="nickname" v-model="nickname">
      <div class="optional">(optional)</div>
    </div>
    <div>
      <label for="aboutMe">About Me:</label>
      <textarea id="aboutMe" v-model="aboutMe"></textarea>
      <div class="optional">(optional)</div>
    </div>
    <button type="submit">Submit</button>
  </form>
  <div v-if="avatarStore.avatarError">{{ avatarStore.avatarError }}</div>
  </div>
</template>

<script lang="ts" setup>
import { Ref, ref } from 'vue';
import router from '@/router/index'
import { signupUser } from '@/api/methods'
import { useAvatarStore } from '@/store/pinia';

const email = ref('');
const password = ref('');
const firstName = ref('');
const lastName = ref('');
const dob = ref('');
const avatar: Ref<File | null> = ref(null); //todo: chat gpt solution, to fix null value case, because field is optional
const nickname = ref('');
const aboutMe = ref('');


const avatarStore = useAvatarStore();
function handleAvatarChange(event: Event) {
  avatarStore.handleAvatarUpload(event);
  avatar.value = (event.target as HTMLInputElement).files?.[0] ?? null;
}

const signup = async () => {
  /* todo: shoud happens only if signup is successful */
  // downscale avatar image and convert it into Blob string of bytes
  signupUser({
    email: email.value,
    password: password.value,
    firstName: firstName.value,
    lastName: lastName.value,
    dob: dob.value,
    avatar: avatar.value,
    nickname: nickname.value,
    aboutMe: aboutMe.value,
    public: false
  })
  router.push('/') /* go back to login, after signup successful */
}

// todo: remove later
const crap = () => {
  email.value = 'dummy@mail.com'
  password.value = '123456'
  firstName.value = 'John'
  lastName.value = 'Doe'
  dob.value = '1990-01-01'
}
</script>
