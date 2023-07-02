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

  <div>
    <h1>Bells</h1>
    <div v-if="paginatedBells.length > 0">
      <button @click="previousPage" :disabled="currentPage === 1">Previous Page</button>
      <button @click="nextPage" :disabled="currentPage === totalPages">Next Page</button>
      <button @click="clearAll">Execute Everyone</button>
    </div>
    <ul>
      <li v-for="(bell, index) in paginatedBells" :key="index">
        <hr>
        <div v-if="bell.type === 'event'">
          type: {{ bell.type }} | {{ bell.event_name }}
          <br> group: {{ bell.group_name }}
          <br> <button @click="openGroup(bell)">Open Group</button>
          <button @click="removeBell(bell)">Close</button>
        </div>
        <div v-else-if="bell.type === 'following'">
          type: {{ bell.type }} | {{ bell.first_name }} {{ bell.last_name }} ({{ bell.email }})
          <br> <button @click="acceptFollowRequest(bell)">Accept</button>
          <button @click="rejectFollowRequest(bell)">Reject</button>
        </div>
        <div v-else-if="bell.type === 'invitation'">
          type: {{ bell.type }} | {{ bell.group_name }}
          <br> <button @click="openGroup(bell)">Open Group</button>
          <br> <button @click="acceptInvitation(bell)">Accept</button>
          <button @click="rejectInvitation(bell)">Reject</button>
        </div>
        <div v-else-if="bell.type === 'request'">
          type: {{ bell.type }} | {{ bell.group_name }}
          <button @click="allowJoinRequest(bell)">Accept</button>
          <button @click="rejectJoinRequest(bell)">Reject</button>
        </div>
      </li>
    </ul>
    <div v-if="paginatedBells.length > 0">
      <button @click="previousPage" :disabled="currentPage === 1">Previous Page</button>
      <button @click="nextPage" :disabled="currentPage === totalPages">Next Page</button>
      <button @click="clearAll">Execute Everyone</button>
    </div>
  </div>

</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue';
import { useWebSocketStore } from '@/store/websocket';
import router from '@/router';
import { useGroupStore } from '@/store/group';
import { Bell, TargetProfileRequest, WSMessageType } from '@/api/types';
import { useUUIDStore } from '@/store/uuid';
import { useProfileStore } from '@/store/profile';

const wss = useWebSocketStore()
const bells = computed(() => wss.bellsList);

const currentPage = ref(1);
const itemsPerPage = 3;

const paginatedBells = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage;
  const end = start + itemsPerPage;
  return bells.value.slice(start, end);
});

const totalPages = ref(Math.ceil(bells.value.length / itemsPerPage));

function updateTotalPages() {
  const totalItems = bells.value.length;
  const lastPage = Math.ceil(totalItems / itemsPerPage);
  if (currentPage.value > lastPage) {
    currentPage.value = lastPage;
  }
  totalPages.value = lastPage;
  console.log("pages ", totalPages.value, " paginatedBells.value.length ", paginatedBells.value.length, "bells.value.length ", bells.value.length);
}

const groupStore = useGroupStore();
function openGroup(bell: Bell) {
  // code to open group
  groupStore.setGroupId(bell.group_id);
  router.push({ name: 'group' });
}

function removeBell(bell: Bell) {
  bells.value.splice(index, 1);
  bellStore.setBells(bells.value);

  // Update totalPages when removing the last event from the last page
  updateTotalPages();

}

function acceptFollowRequest(bell: Bell) {
  // code to accept follow request
}

function rejectFollowRequest(bell: Bell) {
  wss.sendMessage({
    type: WSMessageType.FOLLOW_REQUEST_REJECT,
    data: {
      user_uuid: storeUUID.getUUID,
      target_email: profileStore.getTargetUserEmail,
    } as TargetProfileRequest,
  })
  // bells.value.splice(index, 1);
  // bellStore.setBells(bells.value);
  updateTotalPages();
}

function acceptInvitation(bell: Bell) {
  // code to accept invitation
}

function rejectInvitation(bell: Bell) {
  // bells.value.splice(index, 1);
  // bellStore.setBells(bells.value);
  updateTotalPages();
}

function allowJoinRequest(bell: Bell) {
  // code to allow join request
}

function rejectJoinRequest(bell: Bell) {
  // bells.value.splice(email, 1);
  // bellStore.setBells(bells.value);
  updateTotalPages();
}

function previousPage() {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
}

function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
}

function clearAll() {
  // bells.value = [];
  // bellStore.setBells(bells.value);
}

const storeUUID = useUUIDStore();
const profileStore = useProfileStore();
function updateBells() {
  // todo: add x4 cases for each type of bell
  wss.sendMessage({
    type: WSMessageType.FOLLOW_REQUESTS_LIST,
    data: {
      user_uuid: storeUUID.getUUID,
      target_email: profileStore.getTargetUserEmail,
    } as TargetProfileRequest,
  })

}

onMounted(() => {
  updateBells();
  // bells.value = createDummyData();
  // bellStore.setBells(bells.value);
  updateTotalPages();
});

// todo: refactor/comment later. Dummy data section, should be replaced with real data from the backend

function createDummyData(): Bell[] {
  const dummyData: Bell[] = [];

  for (let i = 0; i < 7; i++) {
    const randomTypeIndex = Math.floor(Math.random() * 4);
    const randomType: string = ['event', 'following', 'invitation', 'request'][randomTypeIndex];

    let bell: Bell;

    if (randomType === 'event') {
      const groupId = generateRandomId();
      bell = {
        type: 'event',
        message: `New event created in Group ${groupId}`,
        groupId,
        userId: -1,
      };
    } else if (randomType === 'following') {
      const userId = generateRandomId();
      bell = {
        type: 'following',
        message: `User ${userId} sent you a follow request`,
        groupId: -1,
        userId,
      };
    } else if (randomType === 'invitation') {
      const groupId = generateRandomId();
      bell = {
        type: 'invitation',
        message: `You've been invited to join Group ${groupId}`,
        groupId,
        userId: -1,
      };
    } else {
      const groupId = -1;
      const userId = -1;
      bell = {
        type: 'request',
        message: `User ${userId} requested to join Group ${groupId}`,
        groupId,
        userId,
      };
    }

    dummyData.push(bell);
  }

  return dummyData;
}

function generateRandomInteger(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function generateRandomId(): number {
  return generateRandomInteger(1, 1000);
}



</script>