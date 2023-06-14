<template>
  <!-- todo: implement view "GroupInviteView.vue" -->
  <pre style="text-align: left;">
    GroupInviteView.vue , can be opened from GroupView.vue,
    by users with membership in group.
    Includes:
    - router-link -> "back to group"
    - multisilect choice list of users to invite, from users who are not in group yet and do not have pending request to join group. Collected from backend.
    - submit button to send invitations, and go "back to group"
  </pre>
  <div>
    <router-link to="/group">Back to Group</router-link>
    <br>
    <label>
      Invite users:
      <select v-model="invitedUsers" multiple>
        <option v-for="user in allUsers" :key="user.id" :value="user.id">{{ user.fullName }} ({{ user.email }})</option>
      </select>
    </label>
    <br>
    <button @click="inviteUsers">Submit</button>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import router from '@/router/index';


interface User {
  id: number
  email: string
  fullName: string
}

const invitedUsers = ref<number[]>([]);
const allUsers: User[] = [
  { id: 1, email: 'alice@mail.com', fullName: 'Alice Dummy' },
  { id: 2, email: 'bob@mail.com', fullName: 'Bob Dummy' },
  { id: 3, email: 'charlie@mail.com', fullName: 'Charlie Dummy' },
  { id: 4, email: 'david@mail.com', fullName: 'David Dummy' },
];

const inviteUsers = () => {
  //todo: Send invitations 
  // go back to group
  router.push({ name: 'group' }); // the group id is same, so no need any changes in pinia, before router push
};

onMounted(() => {
  //todo: Fetch users from backend, and replace dummy data syntax above
});
</script>