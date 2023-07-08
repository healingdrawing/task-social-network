<template>
  <div>
    <!-- SECTION 1 - create new group -->
    <h2>Create new group</h2>
    <form @submit.prevent="createGroup">
      <label>
        Group title:
        <input type="text" v-model="name" required>
      </label>
      <br>
      <label>
        Group description:
        <textarea v-model="description"></textarea>
      </label>
      <br>
      <label>
        Invite users:
        <select v-model="selectedFollowers" multiplev-model="selectedFollowers">
        <option v-for="follower in followersList" :key="follower.email" :value="follower.email">{{ follower.first_name }} {{ follower.last_name }} ({{ follower.email }})</option>
      </select>
      </label>
      <br>
      <button type="submit">Create Group</button>
    </form>

    <!-- SECTION 2 - groups list -->
    <h2><router-link :to="{ name: 'groupsAll' }">Browse All Groups</router-link></h2>
    <h2>Groups list with membership:</h2>
    <ul>
      <li v-for="group in groupsList" :key="group.id">
        <hr>
        <router-link :to="{ name: 'group' }" @click="piniaManageDataGroup(group)">
          group id: {{ group.id }}
          <br> group name: {{ group.name }}
          <br> group description: {{ group.description }}
          <br> group created: {{ group.created_at }}
        </router-link>
        <router-link :to="{ name: 'target' }" @click="piniaManageDataProfile(group.email)">
          <br> group creator: {{ group.first_name }} {{ group.last_name }} ({{ group.email }})
        </router-link>
      </li>
    </ul>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue';
import { useWebSocketStore } from '@/store/pinia';
import { useUUIDStore } from '@/store/pinia';
import { useProfileStore } from '@/store/pinia';
import { useGroupStore } from '@/store/group';
import { TargetProfileRequest, WSMessage, WSMessageType, GroupSubmit, Group } from '@/api/types';

const wss = useWebSocketStore();
const storeUUID = useUUIDStore();
const profileStore = useProfileStore();
const followersList = computed(() => wss.userFollowersList);
function updateFollowersList() {
  wss.sendMessage({
    type: WSMessageType.USER_FOLLOWERS_LIST,
    data: {
      user_uuid: storeUUID.getUUID,
      target_email: profileStore.getUserEmail,
    } as TargetProfileRequest,
  })
}

const groupsList = computed(() => wss.groupsList);
function updateGroupsList() {
  wss.sendMessage({
    type: WSMessageType.GROUPS_LIST,
    data: {
      user_uuid: storeUUID.getUUID,
      target_email: profileStore.getUserEmail,
    } as TargetProfileRequest,
  })
}


const name = ref('');
const description = ref('');
const selectedFollowers = ref<string[]>([]);

const createGroup = () => {
  const message: WSMessage = {
    type: WSMessageType.GROUP_SUBMIT,
    data: {
      user_uuid: storeUUID.getUUID,
      name: name.value,
      description: description.value,
      invited_emails: selectedFollowers.value.join(' '),
    } as GroupSubmit,
  };
  wss.sendMessage(message);

  // reset form
  name.value = '';
  description.value = '';
  selectedFollowers.value = [];
};



const groupStore = useGroupStore();
const piniaManageDataGroup = (group: Group) => {
  groupStore.setGroupId(group.id);
};

const piniaManageDataProfile = (email: string) => {
  profileStore.setTargetUserEmail(email);
};

onMounted(() => {
  updateFollowersList();
  updateGroupsList(); // with membership
});

</script>
