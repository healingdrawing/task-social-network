<template>
  <br>
  <router-link to="/group">
    <div class="router_link_box">
      Back to Group
    </div>
  </router-link>
  <br>
  <div>
    <h1>Create Group Post</h1>

    <div><button type="button" @click="crap" title="remove in production">Fill Debug / remove later</button></div> <!-- todo: remove later -->

    <form @submit.prevent="addGroupPost">
      <label for="postTitle">Group Post title:</label>
      <br> <input type="text" id="postTitle" v-model="postTitle" required>
      <br>
      <label for="postTags">Group Post tags:</label>
      <br> <input title="comma separated" type="text" id="postTags" v-model="postTags">
      <br>
      <label for="postContent">Group Post content:</label>
      <br> <textarea id="postContent" v-model="postContent" required></textarea>
      
      <div>
        <label for="picture" class="label_file_upload">
          with picture(optional):
          <input type="file" id="picture" accept="image/jpeg, image/png, image/gif" @change="handlePictureChange">
        </label>
      </div>

      <button type="submit">Submit</button>
    </form>
    <div v-if="pictureStore.pictureError">{{ pictureStore.pictureError }}</div>
  </div>
  <div>
    <h2>Group Posts:</h2>
    <!-- add posts list , already created -->
    <div v-for="group_post in group_posts_list"
      :key="group_post.id">
      <div class="single_div_box">
        <br>
        <h3> Group Post title: </h3> <p> {{ group_post.title }}</p>
        <h3> Group Post tags: </h3> <p> {{ group_post.categories }}</p>
        <h3> Group Post content: </h3> <p> {{ group_post.content }}</p>
        <h3> Group Post created: </h3> <p> {{ group_post.created_at }}</p>
        <div v-if="group_post.picture !== ''">
          <h3> Group Post picture: </h3>
          <p>
            <img :src="`data:image/jpeg;base64,${group_post.picture}`" alt="picture" />
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
          :to="{ name: 'group_post' }"
          @click="piniaManageDataGroupPost(group_post)">
          <div class="router_link_box">
            comment group post {{ group_post.id }}
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Ref, computed, onMounted, ref } from 'vue';
import { useUUIDStore } from '@/store/pinia';
import { useWebSocketStore } from '@/store/websocket';
import { useGroupStore } from '@/store/pinia';
import { usePostStore } from '@/store/pinia';
import { useProfileStore } from '@/store/pinia';
import { usePictureStore } from '@/store/pinia';
import { WSMessage, WSMessageType, GroupPost, GroupPostSubmit, GroupPostsListRequest } from '@/api/types';

const wss = useWebSocketStore();
const group_posts_list = computed(() => wss.groupPostsList); // ref and reactive failed to work here, so computed used. Straight way put webSocketStore.groupPostsList to template works too,


const postTitle = ref('');
const postTags = ref('');
const postContent = ref('');
const picture: Ref<Blob | null> = ref(null);

const pictureStore = usePictureStore();
function handlePictureChange(event: Event) {
  pictureStore.handlePictureUpload(event);
  picture.value = (event.target as HTMLInputElement).files?.[0] ?? null;
}

const UUIDStore = useUUIDStore();
const groupStore = useGroupStore();
async function addGroupPost() {
  const group_post_submit: GroupPostSubmit = {
    user_uuid: UUIDStore.getUUID,
    group_id: groupStore.getGroup.id,
    title: postTitle.value,
    categories: postTags.value,
    content: postContent.value,
    picture: pictureStore.getPictureBase64String,
  };

  const message: WSMessage = {
    type: WSMessageType.GROUP_POST_SUBMIT,
    data: group_post_submit,
  };
  wss.sendMessage(message);

  postTitle.value = '';
  postTags.value = '';
  postContent.value = '';
  picture.value = null;
  (document.getElementById("picture") as HTMLInputElement).value = "";
  pictureStore.resetPicture()
}

const postStore = usePostStore();
function piniaManageDataGroupPost(group_post: GroupPost) {
  postStore.setGroupPost(group_post);
}

const profileStore = useProfileStore();
function piniaManageDataProfile(email: string) {
  profileStore.setTargetUserEmail(email);
}

// send request to get old group posts list, used inside onMounted
function updateGroupPostsList() {
  wss.sendMessage({
    type: WSMessageType.GROUP_POSTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      group_id: groupStore.getGroup.id,
    } as GroupPostsListRequest,
  });
}

onMounted(async () => {
  wss.refresh_websocket()
  await wss.waitForConnection();
  updateGroupPostsList();
});

const crap = () => {
  postTitle.value = 'Dummy post title';
  postTags.value = 'dummy, post, 111, test';
  postContent.value = 'Dummy Post content text.';
}

</script>