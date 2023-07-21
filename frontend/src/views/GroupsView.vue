<template>
  <div>
    <h2>Create Group</h2>
    <form @submit.prevent="createGroup">
      <label>
        Group title:
        <br> <input type="text" v-model="name" required>
      </label>
      <br>
      <label>
        Group description:
        <br> <textarea v-model="description" required></textarea>
      </label>
      <div v-if="followersList && followersList.length > 0">
        <label>
          Invite followers:
          <br>
          <select multiple v-model="selectedFollowers" class="users_list_with_scroll">
          <option v-for="follower in followersList" :key="follower.email" :value="follower.email">{{ follower.first_name }} {{ follower.last_name }} ({{ follower.email }})</option>
        </select>
        </label>
      </div>
      <br>
      <button type="submit">Create Group</button>
    </form>

    
    <router-link :to="{ name: 'groups_all' }">
      <div class="router_link_box">
        Browse All Groups
      </div>
    </router-link>
  
    <h2>Groups With Membership:</h2>
    <div v-for="group in groupsList" :key="group.id">
      <div class="single_div_box">
        <br>
        <h3> Group name: </h3> <p> {{ group.name }} </p>
        <h3> Group description: </h3> <p> {{ group.description }} </p>
        <h3> Group created: </h3> <p> {{ group.created_at }} </p>
        <router-link
          :Title = "group.first_name + '\n' + group.last_name + '\n' + group.email"
          :to="{ name: 'target' }" @click="piniaManageDataProfile(group.email)">
          <div class="router_link_box">
            visit creator profile
          </div>
        </router-link>
        <br>
        <router-link :to="{ name: 'group' }" @click="piniaManageDataGroup(group)">
          <div class="router_link_box">
            visit group {{ group.id }}
          </div>
        </router-link>
      </div>
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

onMounted(async () => {
  wss.refresh_websocket()
  await wss.waitForConnection();
  updateFollowersList();
  updateGroupsList(); // with membership
});

</script>
