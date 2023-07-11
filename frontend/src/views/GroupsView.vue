<template>
  <div>
    <h2>Create Group</h2>
    <form @submit.prevent="createGroup">
      <label>
        Group title:
        <input type="text" v-model="name" required>
      </label>
      <br>
      <label>
        Group description:
        <textarea v-model="description" required></textarea>
      </label>
      <div v-if="followersList && followersList.length > 0">
        <br>
        <label>
          Invite followers:
          <select multiple v-model="selectedFollowers">
          <option v-for="follower in followersList" :key="follower.email" :value="follower.email">{{ follower.first_name }} {{ follower.last_name }} ({{ follower.email }})</option>
        </select>
        </label>
      </div>
      <br>
      <button type="submit">Create Group</button>
    </form>

    <h2><router-link :to="{ name: 'groups_all' }">Browse All Groups</router-link></h2>
    <h2>Groups With Membership:</h2>
    <div v-for="group in groupsList" :key="group.id">
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
    </div>
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
const UUIDStore = useUUIDStore();
const profileStore = useProfileStore();
const followersList = computed(() => wss.userFollowersList);
function updateFollowersList() {
  wss.sendMessage({
    type: WSMessageType.USER_FOLLOWERS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      target_email: profileStore.getUserEmail,
    } as TargetProfileRequest,
  })
}

const groupsList = computed(() => wss.groupsList);
function updateGroupsList() {
  wss.sendMessage({
    type: WSMessageType.GROUPS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
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
      user_uuid: UUIDStore.getUUID,
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
  groupStore.setGroup(group);
};

const piniaManageDataProfile = (email: string) => {
  profileStore.setTargetUserEmail(email);
};

onMounted(() => {
  updateFollowersList();
  updateGroupsList(); // with membership
});

</script>
