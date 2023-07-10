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
    <div v-if="followersList && followersList.length > 0">
      <br>
      <label>
        Invite followers:
        <select multiple v-model="selectedFollowers">
          <option v-for="follower in followersList" :key="follower.email" :value="follower.email">{{ follower.first_name }} {{ follower.last_name }} ({{ follower.email }})</option>
        </select>
      </label>
      <br>
      <button @click="inviteUsers">Submit</button>
    </div>
    <div v-else>
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
      group_id: groupStore.getGroupId,
      invited_emails: selectedFollowers.value.join(' '),
    } as GroupInvitesSubmit,
  });
  // go back to group
  router.push({ name: 'group' }); // the group id is same, so no need any changes in pinia, before router push
};

onMounted(() => {
  updateFollowersList();
});
</script>