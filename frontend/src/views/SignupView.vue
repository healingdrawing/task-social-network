<template>
  <div class="signup">
    <img alt="Vue logo" src="../assets/logo.png">
    <div>"SignupView.vue . Welcome to Your Vue.js + TypeScript App"</div>
  </div>
  <div>
    <h1>Sign up</h1>
    
    <div><hr><button type="button" @click="crap" title="remove in production">Fill Debug / remove later</button><hr></div> <!-- todo: remove later -->
    
    <form @submit.prevent="signup"> 
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
  <div v-if="error">{{ error }}</div>
  </div>
</template>

<script lang="ts" setup>
import { Ref, ref } from 'vue';
import router from '@/router/index'
import { useWebSocketStore } from '@/store/pinia';
import { useUUIDStore } from '@/store/pinia';
import { useSignupStore } from '@/store/pinia';
import { useProfileStore } from '@/store/pinia';
import { useAvatarStore } from '@/store/pinia';


const wss = useWebSocketStore();
const UUIDStore = useUUIDStore();
const signupStore = useSignupStore();
const profileStore = useProfileStore();
const error = signupStore.error

function resetPiniaStores() {
  signupStore.$reset();
  UUIDStore.$reset();
  wss.$reset();
}

const email = ref('');
const password = ref('');
const firstName = ref('');
const lastName = ref('');
const dob = ref('');
const avatar: Ref<Blob | null> = ref(null); //todo: chat gpt solution, to fix null value case, because field is optional
const nickname = ref('');
const aboutMe = ref('');


const avatarStore = useAvatarStore();
function handleAvatarChange(event: Event) {
  avatarStore.handleAvatarUpload(event);
  avatar.value = (event.target as HTMLInputElement).files?.[0] ?? null;
}

const signup = async () => {
  resetPiniaStores();
  try {
    /* todo: should happen only if signup is successful */
    await signupStore.fetchData({
      email: email.value,
      password: password.value,
      first_name: firstName.value,
      last_name: lastName.value,
      dob: dob.value,
      avatar: avatarStore.getAvatarString,
      nickname: nickname.value,
      about_me: aboutMe.value,
      public: false
    });


    // Result storage logic here
    // For example, you can store the result in a Vuex store or any other storage mechanism
    // Assuming you have a Vuex store setup, you can dispatch an action to store the result
    // Example:
    // await store.dispatch('storeSignupResult', result);

    if (signupStore.getData.UUID === undefined) {
      signupStore.error = "Error: UUID is undefined. Signup failed.";
      throw new Error(signupStore.error);
    } else {
      console.log("UUID: " + UUIDStore.getUUID);
      UUIDStore.setUUID(signupStore.getData.UUID)
      profileStore.setUserEmail(signupStore.getData.email);
      wss.connect(UUIDStore.getUUID);

      // slowdown this... masterpeice ... to wait for websocket to establish connection
      if (wss.socket) {
        const socket = wss.socket
        await new Promise((resolve) => {
          socket.onopen = resolve;
        });
      } else {
        throw new Error('WebSocket connection is null');
      }

      router.push('/profile');
    }

  } catch (error) {
    // Error handling logic here
    // For example, you can display the error message or log it
    // Assuming you have a Vuex store setup for error handling, you can dispatch an action to handle the error
    // Example:
    // await store.dispatch('handleSignupError', error);
    console.error(error);
  } finally {
    // Finally logic here
    console.log('finally');
    // print the result to the console
    // router.push('/profile');
  }


  // downscale avatar image and convert it into Blob string of bytes
  // signupUser({
  //   email: email.value,
  //   password: password.value,
  //   firstName: firstName.value,
  //   lastName: lastName.value,
  //   dob: dob.value,
  //   avatar: avatar.value,
  //   nickname: nickname.value,
  //   aboutMe: aboutMe.value,
  //   public: false
  // })
};

// todo: remove later
const crap = () => {
  email.value = 'dummy@mail.com'
  password.value = '123456'
  firstName.value = 'John'
  lastName.value = 'Doe'
  dob.value = '1990-01-01'
}
</script>
