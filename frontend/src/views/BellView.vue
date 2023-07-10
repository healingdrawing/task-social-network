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
          Your Majesty! Some noise is heard from the castle walls.
          It is about:
          <br> " {{ bell.event_name }} "
          <br> from the:
          <br> " {{ bell.group_name }} "
          <br> <button title="Discover" @click="openGroup(bell)">
            This could be a matter of extreme importance 😤 !
            <br> Prepare my royal horse ! Open the gate !
          </button>
          <button title="Close Event Reminder" @click="removeBell(bell)">
            Again ?! 😒 Boring! Prepare my dolphins 🥹 !
            <br> I am full of spirit today. Move me to the pool.
          </button>
          <h6>
            type: {{ bell.type }} notification
            <br> name: {{ bell.event_name }}
            <br> group: {{ bell.group_name }}
          </h6>
        </div>
        <div v-else-if="bell.type === BellType.FOLLOWING">
          Your Majesty! A peasant named 
          <br> {{ bell.first_name }} {{ bell.last_name }} ({{ bell.email }})
          <br> is in revolt.
          <br> Says that a member of the royal family
          <br> from a neighboring kingdom.
          <br> Also says there is not enough snow in their market.
          <br> <button title="Accept Following Request" @click="acceptFollowRequest(bell)">
            😳 Outrageous! Open the gate!
            <br> A matter of extreme importance!
            <br> So my majesty should
            <br> powder the nose first...
          </button>
          <button title="Reject Following Request" @click="rejectFollowRequest(bell)">
            Terrible! Can't you see I'm eating!
            <br> In shock, I spilled the spirit on my pants.
            <br> Bring me the head 😌 of this poor peasant.
            <br> I want to look into those dishonest eyes.
          </button>
          <h6>
            type: {{ bell.type }} request
            <br> from: {{ bell.first_name }} {{ bell.last_name }} ({{ bell.email }})
          </h6>
        </div>
        <div v-else-if="bell.type === BellType.INVITATION">
          Ambassador of an international organization called
          <br> " {{ bell.group_name }} "
          <br> respectfully invites Your Majesty to join the Board of Governors.
          <br> Says their market trades more snow than Your Majesty's market.
          <br> Says can prove 😏.
          <br> <button title="Visit" @click="openGroup(bell)">
            🤯 Outrageous!!! More snow than in my market!
            <br> Alert my personal leprechaun squad, mobile amusement park
            <br> with blackjack and ... and a swimming pool with trained dolphins 🧐 !
            <br> We are moving out now ! Open the gate !
          </button>
          <br> <button title="Accept Invitation" @click="acceptInvitation(bell)">
            🤔 can prove ... 😳!
            <br> Execute a royal decree 😤 !
            <br> Prepare a banquet hall and
            <br> a trained dolphin 🧐 with soy sauce.
            <br> Bring me this gorgeous person 🥹 !
            <br> This could be a matter of extreme importance 😤 !
            <br> Also make today the annual official holiday of snow 🥹 !
            <br> It's not every day you meet a person
            <br> who can prove 🥹 for free.
          </button>
          <button title="Reject Invitation" @click="rejectInvitation(bell)">
            😠 Don't you see how I'm suffering 🥺 ?
            <br> They said I should stop my spirit diet to boost my spirit.
            <br> But it sounds stupid and works stupidly.
            <br> I know how the universe works.
            <br> More spirit equals more spirit. It is obvious! I'll prove! 
            <br> Why should I suffer alone 🥺 ? Execute a royal decree 😤 !
            <br> Stop the spirit diet of my trained dolphins!
            <br> Attach a laser blaster to the head of each dolphin and
            <br> teleport them to the headquarters of this organization!
          </button>
          <h6>
            type: {{ bell.type }} to join
            <br> group: " {{ bell.group_name }} "
            <br> from: {{ bell.first_name }} {{ bell.last_name }} ({{ bell.email }})
          </h6>
        </div>
        <div v-else-if="bell.type === BellType.REQUEST">
          Your Majesty! The spy 🕵️ is caught outside the castle walls!
          <br> Says that he brought, in a bag, snow
          <br> from a neighboring kingdom, for research by your scientists.
          <br> Also says that wants to join 🤩 the organization
          <br> " {{ bell.group_name }} "
          <br> created by Your Majesty.
          <br> <button title="Accept Request" @click="acceptJoinRequest(bell)">
            Perfect! The Kingdom needs environmentalists!
            <br> Appointing him as a florist 🧐 in my poppy fields.
            <br> Fine and red 🥴 is not bad. And now it's lunch time!
            <br> But first 😏 My Majesty will powder the nose!
          </button>
          <button title="Reject Request" @click="rejectJoinRequest(bell)">
            My Majesty 🧐 grants him freedom!
            <br> 😳 Bring me all his snow!
            <br> Pour the scientists 🥴 10%!
            <br> The rest I'll research 😤 personally!
          </button>
          <h6>
            type: {{ bell.type }} to join
            <br> group: " {{ bell.group_name }} "
            <br> from: {{ bell.first_name }} {{ bell.last_name }} ({{ bell.email }})
          </h6>
        </div>
      </li>
    </ul>
  </div>

</template>

<script lang="ts" setup>
import { computed, onMounted } from 'vue';
import { useWebSocketStore } from '@/store/websocket';
import router from '@/router';
import { useGroupStore } from '@/store/group';
import { BellType, Bell, TargetProfileRequest, WSMessageType, GroupRequestActionSubmit, GroupVisitorStatusRequest } from '@/api/types';
import { useUUIDStore } from '@/store/uuid';
import { useProfileStore } from '@/store/profile';

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
  wss.sendMessage({
    type: WSMessageType.GROUP_INVITE_ACCEPT,
    data: {
      user_uuid: UUIDStore.getUUID,
      group_id: bell.group_id,
      requester_email: bell.email,
    } as GroupVisitorStatusRequest,
  })
  updateBells();
}

function rejectInvitation(bell: Bell) {
  wss.sendMessage({
    type: WSMessageType.GROUP_INVITE_REJECT,
    data: {
      user_uuid: UUIDStore.getUUID,
      group_id: bell.group_id,
      requester_email: bell.email,
    } as GroupVisitorStatusRequest,
  })
  updateBells();
}

function acceptJoinRequest(bell: Bell) {
  wss.sendMessage({
    type: WSMessageType.GROUP_REQUEST_ACCEPT,
    data: {
      user_uuid: UUIDStore.getUUID,
      group_id: bell.group_id,
      requester_email: bell.email,
    } as GroupRequestActionSubmit,
  })
  updateBells();
}

function rejectJoinRequest(bell: Bell) {
  wss.sendMessage({
    type: WSMessageType.GROUP_REQUEST_REJECT,
    data: {
      user_uuid: UUIDStore.getUUID,
      group_id: bell.group_id,
      requester_email: bell.email,
    } as GroupRequestActionSubmit,
  })
  updateBells();
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
  wss.sendMessage({
    type: WSMessageType.GROUP_REQUESTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      target_email: profileStore.getTargetUserEmail,
    } as TargetProfileRequest,
  })
  wss.sendMessage({
    type: WSMessageType.GROUP_INVITES_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      target_email: profileStore.getTargetUserEmail,
    } as TargetProfileRequest,
  })

}

onMounted(() => {
  updateBells();
});

</script>