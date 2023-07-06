<template>
  <router-link to="/bell" @click="wss.facepalm()" >Express Royal Will</router-link>
  <h1>Profile:</h1>
  <!-- todo: add button to open BellView.vue this button should be highlighted in case of still present a new, not marked by user as read already, notifications -->
  <!-- add checkbox to make profile public -->
  <div>
    <label v-if="profile">
      <input type="checkbox" v-model="profile.public" />
      is public
    </label>
  </div>
  <!-- add user information -->
  <div v-if="profile">
    <p>Email: {{ profile.email }}</p>
    <p>First Name: {{ profile.first_name }}</p>
    <p>Last Name: {{ profile.last_name }}</p>
    <p>Date of Birth: {{ profile.dob }}</p>
    <p>Nickname: {{ profile.nickname }}</p>
    <p>About Me: {{ profile.about_me }}</p>
    <p>Public: {{ profile.public }}</p>
  </div>
  <!-- separately add avatar, perhaps it should be on the right half of screen -->
  <div v-if="profile">
    <p>Avatar: {{ profile.avatar }}</p>
  </div>
  <!-- add following list. The other users followed by the user -->
  <h2>Following:</h2>
    <div v-if="followingList.length > 0" class="user-list" style="height: 100px; overflow-y: scroll;">{{ followingList.length }} <br> {{ followingList }}
      <div v-for="user in followingList" :key="user.email">{{ `${user.first_name} ${user.last_name} (${user.email})` }}</div>
    </div>
    <div v-else>No following</div>

    <!-- add followers list. The other users following the user -->
    <h2>Followers:</h2>
    <div v-if="followersList.length > 0" class="user-list" style="height: 100px; overflow-y: scroll;"> {{ followersList.length }} <br> {{ followersList }}
      <div v-for="user in followersList" :key="user.email">{{ `${user.first_name} ${user.last_name} (${user.email})` }}</div>
    </div>
    <div v-else>No followers</div>

    <!-- add user posts list. The posts created by the user -->
    <h2>Posts:</h2>
    <div v-for="post in postsList"
      :key="post.id">
      <hr>
      <router-link
        :to="{ name: 'post' }"
        @click="piniaManageDataPost(post)">
        <p>Post id: {{ post.id }}</p>
        <p>Post title: {{ post.title }}</p>
        <p>Post tags: {{ post.categories }}</p>
        <p>Post content: {{ post.content }}</p>
        <p>Post privacy: {{ post.privacy }}</p><!-- todo: no need to display -->
        <p>Post picture: {{ post.picture }}</p>
        <p>Post created: {{ post.created_at }}</p>
      </router-link>
      <router-link
      :to="{ name: 'target' }"
      @click="piniaManageDataProfile(post.email)">
        <h3>
          Author: {{ post.first_name }}
          {{ post.last_name }} 
          ({{ post.email }})
        </h3>
      </router-link>
  </div>
  <!-- ( :to="{ name: 'post' }" ) also can be ( :to="'/post'" ) -->
</template>

<script lang="ts" setup>
import { ref, watch, onMounted, computed } from 'vue';
import { useWebSocketStore } from '@/store/websocket';
import { useUUIDStore } from '@/store/uuid';
import { usePostStore } from '@/store/post';
import { useProfileStore } from '@/store/profile';
import { WSMessageType, ChangePrivacyRequest, TargetProfileRequest, Post } from '@/api/types';

const wss = useWebSocketStore();

const isPublic = ref(false);

watch(isPublic, (newValue, oldValue) => {
  handleCheckboxChange(newValue);
});

const storeUUID = useUUIDStore();
function handleCheckboxChange(value: boolean) {
  wss.sendMessage({
    type: WSMessageType.USER_PRIVACY,
    data: {
      user_uuid: storeUUID.getUUID,
      make_public: value,
    } as ChangePrivacyRequest,
  })
  // Call your method here
  alert(value + ' . Checkbox changed. ProfileView.vue');
}

function getImgUrl(imageNameWithExtension: string) {
  return require(`../assets/${imageNameWithExtension}`)
}

const profileStore = useProfileStore();
function piniaManageDataProfile(email: string) {
  profileStore.setTargetUserEmail(email);
}

const profile = computed(() => wss.userProfile);
/** updateProfile updates the profile data from server*/
function updateProfile() {
  wss.sendMessage({
    type: WSMessageType.USER_PROFILE,
    data: {
      user_uuid: storeUUID.getUUID,
      target_email: profileStore.getUserEmail,
    } as TargetProfileRequest,
  })
}

const followingList = computed(() => wss.userFollowingList);
/** updateFollowingList updates the following list from server*/
function updateFollowingList() {
  wss.sendMessage({
    type: WSMessageType.USER_FOLLOWING_LIST,
    data: {
      user_uuid: storeUUID.getUUID,
      target_email: profileStore.getUserEmail,
    } as TargetProfileRequest,
  })
}

const followersList = computed(() => wss.userFollowersList);
/** updateFollowersList updates the followers list from server*/
function updateFollowersList() {
  wss.sendMessage({
    type: WSMessageType.USER_FOLLOWERS_LIST,
    data: {
      user_uuid: storeUUID.getUUID,
      target_email: profileStore.getUserEmail,
    } as TargetProfileRequest,
  })
}

const postsList = computed(() => wss.postsList);

/** updatePostsList updates the posts list from the server, able to see for the visitor */
function updatePostsList() {
  wss.sendMessage({
    type: WSMessageType.USER_POSTS_LIST, // todo: do not forget filter by able to see
    data: {
      user_uuid: storeUUID.getUUID,
      target_email: profileStore.getTargetUserEmail,
    } as TargetProfileRequest,
  })
  console.log('Posts list updated');
}


const postStore = usePostStore()
const piniaManageDataPost = (post: Post) => {
  postStore.setPostId(post.id)
}

onMounted(() => {
  updateProfile();
  updateFollowingList();
  updateFollowersList();
  updatePostsList();
});

</script>