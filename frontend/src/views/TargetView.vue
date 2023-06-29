<template>
  <h1>Profile:</h1>

  <div>
    <button @click="handleFollowing()">
      {{ isVisitorNotFollowerAndNotRequester ? 'Request To Follow' : 'Unfollow' }}
    </button>
  </div>

  <div v-if="isProfilePublicOrVisitorFollower">
    <!-- todo: remove later . show id for dev needs-->
    <p>target_email: {{ profileStore.getTargetUserEmail }}</p>
    <!-- add user information -->
    <div v-if="profile">
      <p>Email: {{ profile.email }}</p>
      <p>First Name: {{ profile.first_name }}</p>
      <p>Last Name: {{ profile.last_name }}</p>
      <p>Date of Birth: {{ profile.dob }}</p>
      <p>Nickname: {{ profile.nickname }}</p>
      <p>About Me: {{ profile.about_me }}</p>
    </div>
    <!-- separately add avatar, perhaps it should be on the right half of screen -->
    <div v-if="profile">
      <!-- <p>Avatar: <img :src="getImgUrl(profile.avatar)" alt="fail again"></p> -->
      <p>Avatar: {{ profile.avatar }} </p>
    </div>
    <!-- add following list. The other users followed by the user -->
    <h2>Following:</h2>
    <div class="user-list" style="height: 100px; overflow-y: scroll;">
      <div v-for="user in followingList" :key="user.email">{{ `${user.first_name} ${user.last_name} (${user.email})` }}</div>
    </div>
    <!-- add followers list. The other users following the user -->
    <h2>Followers:</h2>
    <div class="user-list" style="height: 100px; overflow-y: scroll;">
      <div v-for="user in followersList" :key="user.email">{{ `${user.first_name} ${user.last_name} (${user.email})` }}</div>
    </div>
    <!-- add user posts list. The posts created by the user -->
    <h2>Posts:</h2>
    <div v-for="post in postsList"
      :key="post.id">
      <router-link
      :to="{ name: 'post' }"
      @click="piniaManageData(post)">
        {{ post.title }}
      </router-link>
    </div>
    <!-- ( :to="{ name: 'post' }" ) also can be ( :to="'/post'" ) -->
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted, onBeforeMount, computed, reactive } from 'vue';
import { useWebSocketStore } from '@/store/websocket';
import { useUUIDStore } from '@/store/uuid';
import { usePostStore } from '@/store/post';
import { useProfileStore } from '@/store/profile';
import { WSMessageType, TargetProfileRequest, UserProfile } from '@/api/types';

// if true/false, then show follow/unfollow text on button
const isVisitorNotFollowerAndNotRequester = ref(true);

watch(isVisitorNotFollowerAndNotRequester, (newValue, oldValue) => {
  alert(`isVisitorNotFollowerAndDidNotRequested: ${newValue}`);
  // handleFollowing(newValue);
});

function handleFollowing() {
  // Call your method here
  isVisitorNotFollowerAndNotRequester.value = !isVisitorNotFollowerAndNotRequester.value;
}

// if true then show all profile information on screen
const isProfilePublicOrVisitorFollower = ref(true);

function getImgUrl(imageNameWithExtension: string) {
  return require(`../assets/${imageNameWithExtension}`)
}

const webSocketStore = useWebSocketStore();
const storeUUID = useUUIDStore();
const profileStore = useProfileStore();


const profile = computed(() => webSocketStore.userProfile);
/** updateProfile updates the profile data from server*/
function updateProfile() {
  webSocketStore.sendMessage({
    type: WSMessageType.USER_PROFILE,
    data: {
      user_uuid: storeUUID.getUUID,
      target_email: profileStore.getTargetUserEmail,
    } as TargetProfileRequest,
  })
}
const followingList = computed(() => webSocketStore.userFollowingList);
/** updateFollowingList updates the following list from server*/
function updateFollowingList() {
  webSocketStore.sendMessage({
    type: WSMessageType.USER_FOLLOWING_LIST,
    data: {
      user_uuid: storeUUID.getUUID,
      target_email: profileStore.getTargetUserEmail,
    } as TargetProfileRequest,
  })
}

const followersList = computed(() => webSocketStore.userFollowersList);

// todo: dummy data, remove/refactor later
function updateFollowersList() {
  webSocketStore.sendMessage({
    type: WSMessageType.USER_FOLLOWERS_LIST,
    data: {
      user_uuid: storeUUID.getUUID,
      target_email: profileStore.getTargetUserEmail,
    } as TargetProfileRequest,
  })
}

// user posts section
interface Post {
  id: number;
  title: string;
}

const postsList = ref<Post[]>([]);

// todo: dummy data, remove/refactor later
function updatePostsList() {
  // Code to get the user posts goes here
  const posts: Post[] = [
    { id: 1, title: 'Dummy Post 1 Title' },
    { id: 2, title: 'Dummy Post 2 Title' },
    { id: 3, title: 'Dummy Post 3 Title' },
    { id: 4, title: 'Dummy Post 4 Title' },
    { id: 5, title: 'Dummy Post 5 Title' },
    { id: 6, title: 'Dummy Post 6 Title' },
    { id: 7, title: 'Dummy Post 7 Title' },
    { id: 8, title: 'Dummy Post 8 Title' },
    { id: 9, title: 'Dummy Post 9 Title' },
    { id: 10, title: 'Dummy Post 10 Title' },
  ];
  postsList.value = posts;
  console.log('Posts list updated');
}


const postStore = usePostStore()
const piniaManageData = (post: Post) => {
  postStore.setPostId(post.id)
}

onMounted(() => {
  console.log('profile.value', profile)
  updateProfile();
  updateFollowingList();
  updateFollowersList();
  updatePostsList();
});



</script>