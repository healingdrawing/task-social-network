<template>
  <h1>Target Profile:</h1>

  <div v-if="visitor">
    <button @click="handleFollowing()">{{ button_text }}</button>
  </div>
  <div v-else>what is going on, where is the visitor property? [{{ visitor }}]</div>
 
  <!-- add user information -->
  <div v-if="profile">
    <p>Email: {{ profile.email }}</p>
    <p>First Name: {{ profile.first_name }}</p>
    <p>Last Name: {{ profile.last_name }}</p>
  </div>
  <div v-else>no profile</div>

  <div v-if="profile && (profile.public ||
    visitor && visitor.status == VisitorStatus.FOLLOWER ||
    visitor && visitor.status == VisitorStatus.OWNER)">
    <div v-if="profile">
      <p>Date of Birth: {{ profile.dob }}</p>
      <p>Nickname: {{ profile.nickname }}</p>
      <p>About Me: {{ profile.about_me }}</p>
      <p>{{ profile.public ? "Public" : "Private" }}</p>
    </div>
    <!-- separately add avatar, perhaps it should be on the right half of screen -->
    <div v-if="profile && profile.avatar !== ''">
      <p>Avatar:
        <br> <img :src="`data:image/jpeg;base64,${profile.avatar}`" alt="avatar" />
      </p>
    </div>
    <!-- add following list. The other users followed by the user -->
    <h2>Following:</h2>
    <div v-if="followingList.length > 0" class="user-list" style="height: 100px; overflow-y: scroll;">
      <!-- {{ followingList.length }} <br> {{ followingList }} -->
      <div v-for="user in followingList" :key="user.email">{{ `${user.first_name} ${user.last_name} (${user.email})` }}</div>
    </div>
    <div v-else>No following</div>

    <!-- add followers list. The other users following the user -->
    <h2>Followers:</h2>
    <div v-if="followersList.length > 0" class="user-list" style="height: 100px; overflow-y: scroll;">
      <!-- {{ followersList.length }} <br> {{ followersList }} -->
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
        <p>Post created: {{ post.created_at }}</p>
        <div v-if="post.picture !== ''">
          <p>Post picture: 
            <img :src="`data:image/jpeg;base64,${post.picture}`" alt="picture" />
          </p>
        </div>
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

    <!-- add user group posts list. The group posts created by the user in time of group membership -->
    <h2>Group Posts:</h2>
    <div v-for="group_post in groupPostsList"
      :key="group_post.id">
      <hr>
      <router-link
        :to="{ name: 'group' }"
        @click="piniaManageDataGroupPost(group_post)">
        <p>Group id: {{ group_post.group_id }}</p>
        <p>Group name: {{ group_post.group_name }}</p>
        <p>Group description: {{ group_post.group_description }}</p>
      </router-link>
      <p>Group Post id: {{ group_post.id }}</p>
      <p>Group Post title: {{ group_post.title }}</p>
      <p>Group Post tags: {{ group_post.categories }}</p>
      <p>Group Post content: {{ group_post.content }}</p>
      <p>Group Post created: {{ group_post.created_at }}</p>
      <div v-if="group_post.picture !== ''">
        <p>Group Post picture: 
          <br> <img :src="`data:image/jpeg;base64,${group_post.picture}`" alt="picture" />
        </p>
      </div>
      <h3>
        Author: {{ group_post.first_name }}
        {{ group_post.last_name }} 
        ({{ group_post.email }})
      </h3>
    </div>

  </div>
</template>

<script lang="ts" setup>
import { onMounted, computed } from 'vue';
import { useWebSocketStore } from '@/store/websocket';
import { useUUIDStore } from '@/store/uuid';
import { useProfileStore } from '@/store/profile';
import { usePostStore } from '@/store/post';
import { useGroupStore } from '@/store/group';
import { WSMessageType, TargetProfileRequest, VisitorStatus, Post, GroupPost, GroupPostsListRequest, Group } from '@/api/types';



const wss = useWebSocketStore();
const visitor = computed(() => wss.visitor);

function updateVisitorStatus() {
  wss.sendMessage({
    type: WSMessageType.USER_VISITOR_STATUS,
    data: {
      user_uuid: UUIDStore.getUUID,
      target_email: profileStore.getTargetUserEmail,
    } as TargetProfileRequest,
  })
}

const button_text = computed(() => {
  if (visitor.value.status === 'visitor') {
    return 'Follow'
  } else if (visitor.value.status === 'follower') {
    return 'Unfollow';
  } else if (visitor.value.status === 'requester') {
    return 'Waiting for Decision';
  } else {
    return 'Train Finger Muscles';
  }
});

function handleFollowing() {
  if (visitor.value.status === 'visitor') {
    wss.sendMessage({
      type: WSMessageType.USER_FOLLOW,
      data: {
        user_uuid: UUIDStore.getUUID,
        target_email: profileStore.getTargetUserEmail,
      } as TargetProfileRequest,
    })
  } else if (visitor.value.status === 'follower') {
    wss.sendMessage({
      type: WSMessageType.USER_UNFOLLOW,
      data: {
        user_uuid: UUIDStore.getUUID,
        target_email: profileStore.getTargetUserEmail,
      } as TargetProfileRequest,
    })
  } else if (visitor.value.status === 'requester') {
    alert("Agree!!! it is too long to wait.")
  } else {
    alert("Your prestige is raising!!!")
  }
  updateVisitorStatus();
}

const UUIDStore = useUUIDStore();
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
      user_uuid: UUIDStore.getUUID,
      target_email: profileStore.getTargetUserEmail,
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
      target_email: profileStore.getTargetUserEmail,
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
      target_email: profileStore.getTargetUserEmail,
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
      target_email: profileStore.getTargetUserEmail,
    } as TargetProfileRequest,
  })
}

const postStore = usePostStore()
function piniaManageDataPost(post: Post) {
  postStore.setPost(post);
}

const groupPostsList = computed(() => wss.groupPostsList);
// send request to get all group posts list, created by user in time of membering groups
function updateUserGroupPostsList() {
  console.log('=======FIRED======= updateGroupPostsList');

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


onMounted(() => {
  updateVisitorStatus();
  updateProfile();
  updateFollowingList();
  updateFollowersList();
  updatePostsList();
  updateUserGroupPostsList();
});

</script>