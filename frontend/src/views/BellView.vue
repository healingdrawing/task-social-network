<template>
  <pre style="text-align: left;">
    //todo: add implementation
    The BellView.vue which displays the list of bells for the user.
    There are x4 types of bells:
    1. "event" - when an event is created in a group the user is a member of. No action required. But it must have buttons :
      - "Open Group" - to open group which event came from.
      - "Close" - to remove the bell from the list.
    2. "following" - when other user sends a follow request to the user , and user profile is private. Must have two buttons:
      - "Accept" - to accept the follow request, and remove the bell from the list.
      - "Reject" - to reject the follow request, and remove the bell from the list.
    3. "invitation" - when other user invites the user to a group(include case of invitation from group the creator in time of group creation). Must have buttons:
      - "Accept" - to join the group, and remove the bell from the list.
      - "Reject" - to reject the invitation to join the group, and remove the bell from the list.
    4. "request" - when other user sends a request to join a group the user is a creator of. Must have buttons:
      - "Accept" - to allow other user to join the group, and remove the bell from the list.
      - "Reject" - to reject the request to join the group, and remove the bell from the list.
    
  </pre>

  <div v-if="bells.length > 0">
    <div>
      bells: {{ bells }} // todo: remove debug
    </div>
    <h1>
      Your Majesty! The streets are not calm again.
      <br> Intervention of Your Majesty is required!
    </h1>
    <h1>Bells</h1>
    <ul>
      <li v-for="(bell, index) in bells" :key="index">
        <hr>
        <div v-if="bell.type === BellType.EVENT">
          type: {{ bell.type }} | {{ bell.event_name }}
          <br> group: {{ bell.group_name }}
          <br> <button @click="openGroup(bell)">Open Group</button>
          <button @click="removeBell(bell)">Close</button>
        </div>
        <div v-else-if="bell.type === BellType.FOLLOWING">
          Your Majesty! A peasant named 
          <br> {{ bell.first_name }} {{ bell.last_name }} ({{ bell.email }})
          <br> is in revolt.
          <br> Says that a member of the royal family
          <br> from a neighboring kingdom.
          <br> Also says there is not enough snow in their market.
          <br> <button title="Accept" @click="acceptFollowRequest(bell)">
            Outrageous! Open the gate!
            <br> A matter of extreme importance!
            <br> So my majesty should
            <br> powder his nose first...
          </button>
          <button title="Reject" @click="rejectFollowRequest(bell)">
            Terrible! Can't you see I'm eating!
            <br> In shock, I spilled the spirit on my pants.
            <br> Bring me the head of this poor peasant.
            <br> I want to look into those dishonest eyes.
          </button>
          <h6>
            type: {{ bell.type }} | {{ bell.first_name }} {{ bell.last_name }} ({{ bell.email }})
          </h6>
        </div>
        <div v-else-if="bell.type === BellType.INVITATION">
          type: {{ bell.type }} | {{ bell.group_name }}
          <br> <button @click="openGroup(bell)">Open Group</button>
          <br> <button @click="acceptInvitation(bell)">Accept</button>
          <button @click="rejectInvitation(bell)">Reject</button>
        </div>
        <div v-else-if="bell.type === BellType.REQUEST">
          type: {{ bell.type }} | {{ bell.group_name }}
          <button @click="allowJoinRequest(bell)">Accept</button>
          <button @click="rejectJoinRequest(bell)">Reject</button>
        </div>
      </li>
    </ul>
  </div>

</template>

<script lang="ts" setup>
import { computed, onBeforeMount, onMounted } from 'vue';
import { useWebSocketStore } from '@/store/websocket';
import router from '@/router';
import { useGroupStore } from '@/store/group';
import { BellType, Bell, TargetProfileRequest, WSMessageType, WSMessage } from '@/api/types';
import { useUUIDStore } from '@/store/uuid';
import { useProfileStore } from '@/store/profile';
import { mapGetters } from 'pinia';

const wss = useWebSocketStore()
const bells = computed(() => wss.bellsList);

const groupStore = useGroupStore();
function openGroup(bell: Bell) {
  // code to open group
  groupStore.setGroupId(bell.group_id);
  router.push({ name: 'group' });
}

function removeBell(bell: Bell) {
  bell.status = 'hidden'
}

function acceptFollowRequest(bell: Bell) {
  wss.sendMessage({
    type: WSMessageType.FOLLOW_REQUEST_ACCEPT,
    data: {
      user_uuid: UUIDStore.getUUID,
      target_email: bell.email,
    } as TargetProfileRequest,
  })
  updateBells();
}

function rejectFollowRequest(bell: Bell) {
  wss.sendMessage({
    type: WSMessageType.FOLLOW_REQUEST_REJECT,
    data: {
      user_uuid: UUIDStore.getUUID,
      target_email: bell.email,
    } as TargetProfileRequest,
  })
  updateBells();
}

function acceptInvitation(bell: Bell) {
  // code to accept invitation
}

function rejectInvitation(bell: Bell) {
  //gap
}

function allowJoinRequest(bell: Bell) {
  // code to allow join request
}

function rejectJoinRequest(bell: Bell) {
  //gap
}

const UUIDStore = useUUIDStore();
const profileStore = useProfileStore();
function updateBells() {
  // todo: add x4 cases for each type of bell
  wss.sendMessage({
    type: WSMessageType.FOLLOW_REQUESTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      target_email: profileStore.getTargetUserEmail,
    } as TargetProfileRequest,
  })
  // wss.sendMessage({
  //   type: WSMessageType.GROUP_REQUESTS_LIST,
  //   data: {
  //     user_uuid:UUIDStore.getUUID,
  //     target_email: profileStore.getTargetUserEmail,
  //   } as TargetProfileRequest,
  // })

}

onMounted(() => {
  updateBells();
});

</script>