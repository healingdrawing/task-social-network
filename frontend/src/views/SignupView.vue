<template>
  <div class="signup">
    <img alt="Vue logo" src="../assets/logo.png">
    <h1>
      The Head of the Royal Family of these lands
      <br> informs you all, that from now on,
      <br> these lands belong to My Will.
      <br> Call me "Your Majesty"! 
      <br> Now Let is continue ...
    </h1>
  </div>
  <div>
    <h3>Sign up</h3>
    
    <div><hr><button type="button" @click="crap" title="remove in production">Fill Debug / remove later</button><hr></div> <!-- todo: remove later -->
    
    <form @submit.prevent="signup"> 
    <div>
      <label for="email">Email:</label>
      <input type="email" id="email" v-model="email" required>
    </div>
    <div>
      <label for="password">Password:</label>
      <input title="6-15 english letters/digits. No spaces" type="password" id="password" v-model="password" required minlength="6" maxlength="15" pattern="[a-zA-Z0-9]+">
    </div>
    <div>
      <label for="firstName">First Name:</label>
      <input title="1-32 characters" type="text" id="firstName" v-model="firstName" required minlength="1" maxlength="32">
    </div>
    <div>
      <label for="lastName">Last Name:</label>
      <input title="1-32 characters" type="text" id="lastName" v-model="lastName" required minlength="1" maxlength="32">
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
      <input title="1-15 characters" type="text" id="nickname" v-model="nickname" maxlength="15">
      <div class="optional">(optional)</div>
    </div>
    <div>
      <label for="aboutMe">About Me:</label>
      <textarea title="1-300 characters" id="aboutMe" v-model="aboutMe" maxlength="300"></textarea>
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
  console.log('== signup ==');
  try {
    /* todo: should happen only if signup is successful */
    await signupStore.fetchData({
      email: email.value,
      password: password.value,
      first_name: firstName.value,
      last_name: lastName.value,
      dob: dob.value,
      avatar: avatarStore.getAvatarBase64String,
      nickname: nickname.value,
      about_me: aboutMe.value,
      public: false
    });

    if (signupStore.getData.UUID === undefined) {
      signupStore.error = "Error: UUID is undefined. Signup failed.";
      throw new Error(signupStore.error);
    } else {
      console.log("UUID: " + UUIDStore.getUUID);
      UUIDStore.setUUID(signupStore.getData.UUID)
      profileStore.setUserEmail(signupStore.getData.email);
      profileStore.setTargetUserEmail(signupStore.getData.email);
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
    alert("Are you sure you are from Royal Family? :|\nWe expect you at least:\n- can read and write o:)\n- know other Royal Families\n- will not try to use their family emails")
    console.error(error);
  } finally {
    console.log("== signup == 'finally' fired");
  }
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
