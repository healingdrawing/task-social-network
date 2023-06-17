<template>
  user id transfered using pinia: {{ profileStore.userId }}
  <br>
  group id transfered using pinia: {{ groupStore.groupId }}
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
    <h1>{{ group.title }}</h1>
    <p>{{ group.description }}</p>
    <div v-if="isMember">
      <button @click="groupChat">Open Group Chat</button>
      <button @click="groupInvite">Invite User</button>
      <button @click="groupPosts">Group Posts</button>
      <div>
        <h2>Create New Event</h2>
        <form @submit.prevent="createEvent">
          <label for="title">Title:</label>
          <input type="text" id="title" v-model="event.title" required>
          <br>
          <label for="datetime">Date and Time:</label>
          <input type="datetime-local" id="datetime" v-model="event.datetime" required>
          <br>
          <label for="description">Description:</label>
          <textarea id="description" v-model="event.description"></textarea>
          <br>
          <label>Going to Event:</label>
          <input type="radio" id="going" value="1" v-model="event.going">
          <label for="going">Yes</label>
          <input type="radio" id="not-going" value="0" v-model="event.going">
          <label for="not-going">No</label>
          <br>
          <button type="submit">Create Event</button>
        </form>
      </div>
      <div>
        <h2>List of Group Events</h2>
        <ul>
          <li v-for="event in events" :key="event.id">
            <h3>{{ event.title }}</h3>
            <p>{{ event.datetime }}</p>
            <p>{{ event.description }}</p>
            <div>
              <input type="radio" id="going-yes" name="going-{{ event.id }}" value="1" v-bind:checked="event.going === '1'" v-on:change="goingYes(event, '1')">
              <label for="going-yes">Yes</label>
              <input type="radio" id="going-no" name="going-{{ event.id }}" value="0" v-bind:checked="event.going === '0'" v-on:change="goingNo(event, '0')">
              <label for="going-no">No</label>
            </div>
          </li>
        </ul>
      </div>
    </div>
    <div v-else>
      <button @click="joinGroup">Join Group</button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import router from '@/router/index';
import { useProfileStore } from '@/store/profile';
import { useGroupStore } from '@/store/group';
import { useChatStore } from '@/store/chat';
const profileStore = useProfileStore();
const groupStore = useGroupStore();
const chatStore = useChatStore();

interface Group {
  title: string;
  description: string;
}

interface Event {
  id: number;
  title: string;
  datetime: string;
  description: string;
  going: string;
}

const group = ref<Group>({ title: '', description: '' });
const isMember = ref(false);
const event = ref<Event>({ id: 1, title: '', datetime: '', description: '', going: '1' });
const events = ref<Event[]>([]);

//todo: implement createEvent() function
const createEvent = async () => {
  event.value.id = events.value.length + 1
  events.value.unshift(event.value)
  event.value = { id: -1, title: '', datetime: '', description: '', going: '1' }
}

const goingYes = (event: Event, value: string) => {
  event.going = value;
};

const goingNo = (event: Event, value: string) => {
  event.going = value;
};

//todo: implement joinGroup() function
const joinGroup = () => { console.log('join group') }

// open ChatView.vue
const groupChat = () => { router.push({ name: 'chat' }) }

// open GroupInviteView.vue
const groupInvite = () => { router.push({ name: 'groupInvite' }) }

// open PostsView.vue . // todo: Group posts only. Filtered by group id on backend inside onMount() of PostsView.vue. The groupId is already set in groupStore in GroupsView.vue piniaManageData(), before visit that view
const groupPosts = () => { router.push({ name: 'posts' }) }

onMounted(() => {
  group.value.title = 'Dummy Group Title'; //todo: replace with actual group title
  group.value.description = 'Dummy Group Description'; //todo: replace with actual group description

  //todo: get chat id for the group from backend using groupStore.getGroupId
  const chatId = 77; // replace with actual chat ID
  //then set it using pinia, to use it in ChatView.vue, to collect needed data from backend
  chatStore.setChatId(chatId);

  isMember.value = true; //todo: replace with actual membership check
});
</script>