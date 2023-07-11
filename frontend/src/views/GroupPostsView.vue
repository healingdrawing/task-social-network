<template>
  <div>
    <router-link to="/group">Back to Group</router-link>
    <h1>Create Group Post:</h1>

    <div><hr><button type="button" @click="crap" title="remove in production">Fill Debug / remove later</button><hr></div> <!-- todo: remove later -->

    <form @submit.prevent="addGroupPost">
      <label for="postTitle">Post Title:</label>
      <input type="text" id="postTitle" v-model="postTitle" required>
      <br>
      <label for="postTags">Post Tags:</label>
      <input type="text" id="postTags" v-model="postTags">
      <br>
      <label for="postContent">Post Content:</label>
      <textarea id="postContent" v-model="postContent" required></textarea>
      
      <div>
        <label for="picture">Picture:</label>
        <input type="file" id="picture" accept="image/jpeg, image/png, image/gif" @change="handlePictureChange">
        <div class="optional">(optional)</div>
      </div>

      <br>
      <button type="submit">Submit</button>
    </form>
    <div v-if="pictureStore.pictureError">{{ pictureStore.pictureError }}</div>
  </div>
  <div>
    <h2>Group Posts:</h2>
    <!-- add posts list , already created -->
    <div v-for="post in group_posts_list"
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
  </div>
</template>

<script lang="ts" setup>
import { Ref, computed, onMounted, ref, watch } from 'vue';
import { useUUIDStore } from '@/store/pinia';
import { useWebSocketStore } from '@/store/websocket';
import { useGroupStore } from '@/store/pinia';
import { usePostStore } from '@/store/pinia';
import { useProfileStore } from '@/store/pinia';
import { usePictureStore } from '@/store/pinia';
import { WSMessage, WSMessageType, Post, GroupPostSubmit, GroupPostsListRequest } from '@/api/types';
import { onBeforeRouteLeave } from 'vue-router';

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
function piniaManageDataPost(post: Post) {
  postStore.setPost(post);
}

const profileStore = useProfileStore();
function piniaManageDataProfile(email: string) {
  profileStore.setTargetUserEmail(email);
}

// send request to get old posts list, used inside onMounted
function updateGroupPostsList() {
  console.log('=======FIRED======= updateGroupPostsList');

  wss.sendMessage({
    type: WSMessageType.GROUP_POSTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      group_id: groupStore.getGroup.id,
    } as GroupPostsListRequest,
  });
}

onMounted(() => {
  updateGroupPostsList();
});

const crap = () => {
  postTitle.value = 'Dummy post title';
  postTags.value = 'dummy, post, 111, test';
  postContent.value = 'Dummy Post content text.';
}

watch(group_posts_list, (newVal) => {
  console.log('Posts list:', newVal);
});

onBeforeRouteLeave((to, from, next) => {
  wss.facepalm();
  // wss.clearOnBeforeRouteLeave(from.path) // todo: clear the messages. ugly but works
  console.log('onBeforeRouteLeave fired');
  next();
});

</script>