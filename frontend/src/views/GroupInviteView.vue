<template>
  <div>
    <br>
    <router-link to="/group">
      <div class="router_link_box">
        Back to Group
      </div>
    </router-link>
    <br>
    <div v-if="followersList && followersList.length > 0">
      <br> <br>
      <label>
        Invite followers:
        <br>
        <select multiple v-model="selectedFollowers" class="users_list_with_scroll">
          <option v-for="follower in followersList" :key="follower.email" :value="follower.email">{{ follower.first_name }} {{ follower.last_name }} ({{ follower.email }})</option>
        </select>
      </label>
      <br>
      <button @click="inviteUsers">Submit</button>
    </div>
    <div v-else>
      <br> <br>
      No followers to invite.
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue';
import router from '@/router/index';
import { useWebSocketStore } from '@/store/pinia';
import { useUUIDStore } from '@/store/pinia';
import { useProfileStore } from '@/store/pinia';
import { useGroupStore } from '@/store/group';
import { WSMessageType, TargetProfileRequest } from '@/api/types';
import { GroupInvitesSubmit } from '@/api/types';

const selectedFollowers = ref<string[]>([]);

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

const groupStore = useGroupStore();
const inviteUsers = () => {
  wss.sendMessage({
    type: WSMessageType.GROUP_INVITES_SUBMIT,
    data: {
      user_uuid: UUIDStore.getUUID,
      group_id: groupStore.getGroup.id,
      invited_emails: selectedFollowers.value.join(' '),
    } as GroupInvitesSubmit,
  });
  // go back to group
  router.push({ name: 'group' }); // the group id is same, so no need any changes in pinia, before router push
};

onMounted(async () => {
  wss.refresh_websocket()
  await wss.waitForConnection();
  updateFollowersList();
});
</script>