<template>
  <h1>Profile:</h1>
  <!-- todo: add button to open BellView.vue this button should be highlighted in case of still present a new, not marked by user as read already, notifications -->
  <!-- add checkbox to make profile public -->
  <div>
    <label v-if="profile">
      <input type="checkbox" v-model="isPublic" />
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
  </div>
  <!-- separately add avatar, perhaps it should be on the right half of screen -->
  <div v-if="profile && profile.avatar !== ''">
    <p>Avatar: 
      <br> <img :src="`data:image/jpeg;base64,${profile.avatar}`" alt="avatar" />
    </p>
  </div>
  <!-- add following list. The other users followed by the user -->
  <h2>Following:</h2>
  <div v-if="followingList.length > 0" class="users_list_with_scroll">
    <!-- {{ followingList.length }} <br> {{ followingList }} -->
    <div v-for="user in followingList" :key="user.email">{{ `${user.first_name} ${user.last_name} (${user.email})` }}</div>
  </div>
  <div v-else>No following</div>

  <!-- add followers list. The other users following the user -->
  <h2>Followers:</h2>
  <div v-if="followersList.length > 0" class="users_list_with_scroll">
    <!-- {{ followersList.length }} <br> {{ followersList }} -->
    <div v-for="user in followersList" :key="user.email">{{ `${user.first_name} ${user.last_name} (${user.email})` }}</div>
  </div>
  <div v-else>No followers</div>

  <!-- add user posts list. The posts created by the user -->
  <h2>Posts:</h2>
  <div v-for="post in postsList"
      :key="post.id">
      <div class="single_div_box">
        <br>
        <h3> Post title: </h3> <p> {{ post.title }} </p>
        <h3> Post tags: </h3> <p> {{ post.categories }} </p>
        <h3> Post content: </h3> <p> {{ post.content }} </p>
        <h3> Post privacy: </h3> <p> {{ post.privacy }} </p>
        <h3> Post created: </h3> <p> {{ post.created_at }} </p>
        <div v-if="post.picture !== ''">
          <h3> Post picture: </h3>
          <p>
            <img :src="`data:image/jpeg;base64,${post.picture}`" alt="picture" />
          </p>
        </div>
        <router-link
        :Title="post.first_name + '\n' + post.last_name + '\n' + post.email"
        :to="{ name: 'target' }"
        @click="piniaManageDataProfile(post.email)">
          <div class="router_link_box">
            visit author profile
          </div>
        </router-link>
        <br>
        <router-link
        :to="{ name: 'post' }"
        @click="piniaManageDataPost(post)">
          <div class="router_link_box">
            comment post {{ post.id }}
          </div>
        </router-link>
      </div>
    </div>
  <!-- ( :to="{ name: 'post' }" ) also can be ( :to="'/post'" ) -->

  <!-- add user group posts list. The group posts created by the user in time of group membership -->
  <h2>Group Posts:</h2>
  <div v-for="group_post in groupPostsList"
    :key="group_post.id">

    <div class="single_div_box">

      <h3> Group Post id: </h3> <p> {{ group_post.id }} </p>
      <h3> Group Post title: </h3> <p> {{ group_post.title }} </p>
      <h3> Group Post tags: </h3> <p> {{ group_post.categories }} </p>
      <h3> Group Post content: </h3> <p> {{ group_post.content }} </p>
      <h3> Group Post created: </h3> <p> {{ group_post.created_at }} </p>
      <div v-if="group_post.picture !== ''">
        <h3> Group Post picture: </h3>
        <p>
          <br> <img :src="`data:image/jpeg;base64,${group_post.picture}`" alt="picture" />
        </p>
      </div>
      
      <router-link
      :Title="group_post.first_name + '\n' + group_post.last_name + '\n' + group_post.email"
      :to="{ name: 'target' }"
      @click="piniaManageDataProfile(group_post.email)">
        <div class="router_link_box">
          visit author profile
        </div>
      </router-link>
      <br>
      <router-link
      :Title="group_post.group_name"
      :to="{ name: 'group' }"
      @click="piniaManageDataGroupPost(group_post)">
        <div class="router_link_box">
          visit group {{ group_post.group_id }}
        </div>
      </router-link>

    </div>
    
  </div>

</template>

<script lang="ts" setup>
import { ref, watch, onMounted, computed } from 'vue';
import { useWebSocketStore } from '@/store/websocket';
import { useUUIDStore } from '@/store/uuid';
import { useProfileStore } from '@/store/profile';
import { usePostStore } from '@/store/post';
import { useGroupStore } from '@/store/group';
import { WSMessageType, ChangePrivacyRequest, TargetProfileRequest, Post, GroupPost, GroupPostsListRequest, Group } from '@/api/types';

const wss = useWebSocketStore();

const isPublic = ref(false);

watch(isPublic, (newValue) => {
  handleCheckboxChange(newValue);
});

const UUIDStore = useUUIDStore();
function handleCheckboxChange(value: boolean) {
  console.log('= handleCheckboxChange', value);
  // todo: undefined happens in time of move to BellView.vue, perhaps because of wss.facepalm() cleaning + reactivity
  if (value !== undefined) {
    wss.sendMessage({
      type: WSMessageType.USER_PRIVACY,
      data: {
        user_uuid: UUIDStore.getUUID,
        make_public: value.toString(),
      } as ChangePrivacyRequest,
    })
  }
}

const profileStore = useProfileStore();
function piniaManageDataProfile(email: string) {
  profileStore.setTargetUserEmail(email);
}

const profile = computed(() => wss.userProfile);

watch(() => profile.value?.public, (newPublic) => {
  isPublic.value = newPublic;
});

/** updateProfile updates the profile data from server*/
function updateProfile() {
  wss.sendMessage({
    type: WSMessageType.USER_PROFILE,
    data: {
      user_uuid: UUIDStore.getUUID,
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
      user_uuid: UUIDStore.getUUID,
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
      user_uuid: UUIDStore.getUUID,
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
      user_uuid: UUIDStore.getUUID,
      target_email: profileStore.getUserEmail,
    } as TargetProfileRequest,
  })
}

const postStore = usePostStore()
const piniaManageDataPost = (post: Post) => {
  postStore.setPost(post)
}

const groupPostsList = computed(() => wss.groupPostsList);
// send request to get all group posts list, created by user in time of membering groups
function updateUserGroupPostsList() {
  wss.sendMessage({
    type: WSMessageType.USER_GROUP_POSTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      group_id: -1, // no needed, because it is all group posts collector for user
    } as GroupPostsListRequest,
  });
}

const groupStore = useGroupStore()
const piniaManageDataGroupPost = (group_post: GroupPost) => {
  groupStore.setGroup(
    {
      id: group_post.group_id,
      name: group_post.group_name,
      description: group_post.group_description,
    } as Group
  )
}

onMounted(async () => {
  wss.refresh_websocket()
  await wss.waitForConnection();
  updateProfile();
  updateFollowingList();
  updateFollowersList();
  updatePostsList();
  updateUserGroupPostsList();
});

</script>