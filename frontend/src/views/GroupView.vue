<template>
  user id transfered using pinia: {{ profileStore.getUserId }}
  <!-- todo: implement view "GroupView.vue" -->
  <pre style="text-align: left;">
    one group view gap.
    Check membership inside onMounted hook using call to backend.
    view must include:
    -------------------------------------------------
    - if NO MEMBERSHIP, then:
    - - button to make request to join group
    - - title + description of the group
    -------------------------------------------------
    - if MEMBERSHIP, then:
    - - title + description of the group
    - - open group chat button -> opens "ChatView.vue"
    - - invite user button -> opens "GroupInviteView.vue"
    - - group posts button -> opens "PostsView.vue" , but with posts only from this group
    
    - - create new even section
    - - - event title - input text
    - - - event date+time - input date + input time(or some combined Vue widget)
    - - - event description - textarea
    - - - radio button with two variants 0 - not going to event, 1 - going to event
    - - - - (for creator too, inside event creation section,
    and default value for creator is going to event)
    - - - submit button to create event

    - - list of group events section
    - - - for each user, every event has:
     title, date+time, description, going to event button.
     To prevent spamming , only once choice are allowed.
     After that, going button is disabled/replaced by text/or restyled to green etc.
     So user can only once choose to go to event.
  </pre>

  <div>
    <h6 title="group id">{{ group.id }}</h6>
    <h1 title="group name/title">{{ group.name }}</h1>
    <p title="group description">{{ group.description }}</p>
    <!--
      link to group creator link is not required and not implemented,
      because it is headache when jump into group from invite,
      because you need provide giant pile of rare needed info into <Bell>.
      Also sql must be refactored, only for this not required case. So, no.
      Commened to keep the same style for all cases of visiting group.
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
        <button @click="groupInvite"> Invite User </button>
        <button @click="groupPosts"> Group Posts </button>
        <div>
          <h2>Create New Event </h2>
          <form @submit.prevent="createEvent">
            <label for="title"> Title: </label>
            <input type="text" id="title" v-model="event.title" required>
            <br>
            <label for="datetime"> Date and Time: </label>
            <input type="datetime-local" id="datetime" v-model="event.datetime" required>
            <br>
            <label for="description"> Description: </label>
            <textarea id="description" v-model="event.description" required> </textarea>
            <br>
            <label>Going to Event: </label>
            <input type="radio" id="going" value="1" v-model="event.going">
            <label for="going"> Yes </label>
            <input type="radio" id="not-going" value="0" v-model="event.going">
            <label for="not-going"> No </label>
            <br>
            <button type= "submit"> Create Event </button>
          </form>
        </div>
        <h2>List of Group Events </h2>
        <ul>
          <li v-for="event in events" :key="event.id">
            <h3>{{
              event.title }}</h3>
            <p>{{ event.datetime }}</p>
            <p>{{ event.description }}</p>
            <div>
              <input type="radio" id="going-yes" name="going-{{ event.id }}" value="1" v-bind:checked="event.going === '1'" v-on:change="going_yes(event, '1')">
              <label for="going-yes">Yes</label>
              <input type="radio" id="going-no" name="going-{{ event.id }}" value="0" v-bind:checked="event.going === '0'" v-on:change="going_no(event, '0')">
              <label for="going-no">No</label>
            </div>
          </li>
        </ul>
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
import { useProfileStore } from '@/store/profile';
import { useGroupStore } from '@/store/group';
import { useChatStore } from '@/store/chat';
import { GroupVisitorStatusRequest, VisitorStatus, WSMessageType } from '@/api/types';
const wss = useWebSocketStore();
const UUIDStore = useUUIDStore();
const profileStore = useProfileStore();
const groupStore = useGroupStore();
const chatStore = useChatStore();

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

//todo: artefact to provide commented functionality in template section of this file
// const piniaManageDataProfile = (email: string) => {
//   profileStore.setTargetUserEmail(email);
// };

//dummy code


interface Event {
  id: number;
  title: string;
  datetime: string;
  description: string;
  going: string;
}

const group = computed(() => groupStore.getGroup)
const event = ref<Event>({ id: 1, title: '', datetime: '', description: '', going: '1' });
const events = ref<Event[]>([]);

//todo: implement createEvent() function
const createEvent = async () => {
  event.value.id = events.value.length + 1
  events.value.unshift(event.value)
  event.value = { id: -1, title: '', datetime: '', description: '', going: '1' }
}

const going_yes = (event: Event, value: string) => {
  event.going = value;
};

const going_no = (event: Event, value: string) => {
  event.going = value;
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
const groupChat = () => { router.push({ name: 'chat' }) }

// open GroupInviteView.vue
const groupInvite = () => { router.push({ name: 'group_invite' }) }

// open GroupPostsView.vue . // todo: saparated table for group posts used
const groupPosts = () => { router.push({ name: 'group_posts' }) }

onMounted(() => {

  updateGroupVisitor();

  //todo: get chat id for the group from backend using groupStore.getGroupId
  const chatId = 77; // replace with actual chat ID
  //then set it using pinia, to use it in ChatView.vue, to collect needed data from backend
  chatStore.setChatId(chatId);
});
</script>