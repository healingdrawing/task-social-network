<template>
  <div>
    <!-- SECTION 1 - create new group -->
    <h2>Create new group</h2>
    <form @submit.prevent="createGroup">
      <label>
        Group title:
        <input type="text" v-model="title" required>
      </label>
      <br>
      <label>
        Group description:
        <textarea v-model="description"></textarea>
      </label>
      <br>
      <label>
        Invite users:
        <select v-model="invitedUsers" multiple>
          <option v-for="user in allUsers" :key="user.id" :value="user.id">{{ user.fullName }}</option>
        </select>
      </label>
      <br>
      <button type="submit">Create group</button>
    </form>

    <!-- SECTION 2 - groups list -->
    <h2><router-link :to="{ name: 'groupsAll' }">Browse All Groups</router-link></h2>
    <h2>Groups list with membership:</h2>
    <ul>
      <li v-for="group in groupsList" :key="group.id">
        <hr>
        <router-link :to="{ name: 'group' }" @click="piniaManageData(group)">
          group title:{{ group.title }}
          <br>
          ({{ group.members.length }} members)
          <br> (at the moment it is just dummy invitations number, not real membership)
          <br>id:{{ group.id }}
          <br>
          group description:{{ group.description }}
        </router-link>
      </li>
    </ul>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { useGroupStore } from '@/store/group';

interface User {
  id: number;
  fullName: string;
}

interface Group {
  id: number;
  title: string;
  description: string;
  members: User[];
}

const title = ref('');
const description = ref('');
const invitedUsers = ref<number[]>([]);
const allUsers: User[] = [
  { id: 1, fullName: 'Alice Dummy' },
  { id: 2, fullName: 'Bob Dummy' },
  { id: 3, fullName: 'Charlie Dummy' },
  { id: 4, fullName: 'David Dummy' },
];
const groupsList = ref<Group[]>([]);

const createGroup = () => {
  const group: Group = {
    id: Date.now(),
    title: title.value,
    description: description.value,
    members: invitedUsers.value.map(id => allUsers.find(user => user.id === id)!)
  };
  groupsList.value.unshift(group);
  title.value = '';
  description.value = '';
  invitedUsers.value = [];
};

function updateGroupsList() {
  groupsList.value = [
    { id: 1, title: 'Group 1', description: 'Description 1', members: [allUsers[0], allUsers[1]] },
    { id: 2, title: 'Group 2', description: 'Description 2', members: [allUsers[2], allUsers[3]] },
  ]
}

const groupStore = useGroupStore();
const piniaManageData = (group: Group) => {
  groupStore.setGroupId(group.id);
};

onMounted(() => {
  updateGroupsList();
});

</script>
