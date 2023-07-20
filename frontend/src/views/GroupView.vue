<template>
  <div>
    <h6 title="group id">{{ group.id }}</h6>
    <h1 title="group name/title">{{ group.name }}</h1>
    <p title="group description">{{ group.description }}</p>
    <!--
      link to group creator link is not required and not implemented,
      because it is headache when jump into group from invite,
      because you need provide giant pile of rare needed info into <Bell>.
      Also sql must be refactored, only for this not required case. So, no.
      Commented to keep the same style for all cases of visiting group.
      With other cases, visiting group not from invite,
      the uncommented code section bottom works correct, because filled naturally.
    -->
    <!-- <p v-if="group." title="group creator">
      <router-link :to="{ name: 'target' }" @click="piniaManageDataProfile(group.email)">
        {{ group.first_name }} {{ group.last_name }} ({{ group.email }})
      </router-link>
    </p> -->
    <div v-if="group_visitor">
      <div v-if="group_visitor.status === VisitorStatus.MEMBER">
        <button @click= "groupChat"> Open Group Chat </button>
        <button @click="groupInvite"> Invite Followers </button>
        <button @click="groupPosts"> Group Posts </button>
        <div>
          <h2> Create Event </h2>
          <form @submit.prevent="createEvent">
            <label for="title"> Title: </label>
            <br> <input type="text" id="title" v-model="event.title" required>
            <br>
            <label for="datetime"> Date and Time: </label>
            <br> <input type="datetime-local" id="datetime" v-model="event.date" required>
            <br>
            <label for="description"> Description: </label>
            <br> <textarea id="description" v-model="event.description" required> </textarea>
            <br>
            <label>Going to Event: </label>
            <input type="radio" id="going" value="going" v-model="event.decision">
            <label for="going"> going </label>
            <input type="radio" id="not going" value="not going" v-model="event.decision">
            <label for="not going"> not going </label>
            <br>
            <button type= "submit"> Create Event </button>
          </form>
        </div>
        <h2>List of Group Events: </h2>
        <div v-for="event in events_list" :key="event.id">
          <div class="single_div_box">
            <br>
            <h3> Event title: </h3> <p> {{ event.title }} </p>
            <h3> Event description: </h3> <p> {{ event.description }} </p>
            <h3> Event date: </h3> <p>{{ event.date }}</p>
            <h3> Event decision: </h3>
            <div v-if="event.decision === 'waiting'">
              <button @click="going_yes(event)" >going</button>
              <button @click="going_no(event)" >not going</button>
            </div>
            <div v-else>
              {{ event.decision }}
            </div>
          </div>
        </div>
      </div>
      <div v-else-if="group_visitor.status === VisitorStatus.REQUESTER">
        <p>Request to join group is pending.</p>
      </div>
      <div v-else>
        <button @click="joinGroup">Join Group</button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue';
import router from '@/router/index';
import { useWebSocketStore } from '@/store/pinia';
import { useUUIDStore } from '@/store/pinia';
import { useGroupStore } from '@/store/group';
import { GroupVisitorStatusRequest, VisitorStatus, WSMessageType, GroupEventSubmit, WSMessage, Event, GroupEventAction, GroupEventsListRequest } from '@/api/types';
const wss = useWebSocketStore();
const UUIDStore = useUUIDStore();
const groupStore = useGroupStore();

const group_visitor = computed(() => wss.groupVisitor)
function updateGroupVisitor() {
  wss.sendMessage({
    type: WSMessageType.USER_GROUP_VISITOR_STATUS,
    data: {
      user_uuid: UUIDStore.getUUID,
      group_id: groupStore.getGroup.id,
    } as GroupVisitorStatusRequest,
  })
}

const group = computed(() => groupStore.getGroup)
const event = ref<Event>({ id: 1, title: '', date: '', description: '', decision: 'going' });
const events_list = computed(() => wss.groupEventsList)

function updateGroupEventsList() {
  wss.sendMessage({
    type: WSMessageType.GROUP_EVENTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      group_id: groupStore.getGroup.id,
    } as GroupEventsListRequest,
  })
}

//todo: implement createEvent() function
const createEvent = async () => {

  const group_event_submit: GroupEventSubmit = {
    user_uuid: UUIDStore.getUUID,
    group_id: groupStore.getGroup.id,
    title: event.value.title,
    date: event.value.date, // to do not refactor all, will be just date
    description: event.value.description,
    decision: event.value.decision,
  }

  const message: WSMessage = {
    type: WSMessageType.GROUP_EVENT_SUBMIT,
    data: group_event_submit,
  };
  wss.sendMessage(message);

  event.value = { id: -1, title: '', date: '', description: '', decision: 'going' }
}

const going_yes = (event: Event) => {

  const group_event_going = {
    user_uuid: UUIDStore.getUUID,
    event_id: event.id,
    decision: 'going',
    group_id: groupStore.getGroup.id, // to collect list of events after decision
  } as GroupEventAction;

  const message: WSMessage = {
    type: WSMessageType.GROUP_EVENT_GOING,
    data: group_event_going,
  };
  wss.sendMessage(message);
};

const going_no = (event: Event) => {

  const group_event_not_going = {
    user_uuid: UUIDStore.getUUID,
    event_id: event.id,
    decision: 'not going',
    group_id: groupStore.getGroup.id, // to collect list of events after decision
  } as GroupEventAction;

  const message: WSMessage = {
    type: WSMessageType.GROUP_EVENT_NOT_GOING,
    data: group_event_not_going,
  };
  wss.sendMessage(message);
};

const joinGroup = () => {
  wss.sendMessage({
    type: WSMessageType.GROUP_REQUEST_SUBMIT,
    data: {
      user_uuid: UUIDStore.getUUID,
      group_id: groupStore.getGroup.id,
    } as GroupVisitorStatusRequest, // the fields are same , so reuse it
  })
}

// open ChatView.vue
const groupChat = () => { router.push({ name: 'group_chat' }) }

// open GroupInviteView.vue
const groupInvite = () => { router.push({ name: 'group_invite' }) }

// open GroupPostsView.vue . // todo: saparated table for group posts used
const groupPosts = () => { router.push({ name: 'group_posts' }) }

onMounted(async () => {
  wss.refresh_websocket()
  await wss.waitForConnection();

  updateGroupVisitor();
  updateGroupEventsList();

  wss.set_group_chat_id(group.value.id)
  wss.set_private_chat_user_id(0)
});
</script>