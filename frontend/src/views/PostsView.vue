<template>
  <div>
    <h1>Create Post:</h1>
    <form @submit.prevent="addPost">
      <label for="postTitle">Post Title:</label>
      <input type="text" id="postTitle" v-model="postTitle" required>
      <br>
      <label for="postTags">Post Tags:</label>
      <input type="text" id="postTags" v-model="postTags"> <!-- todo: the "required" field removed, because in "real-time-forum" backend the empty tag was implemented as decoration using randomly colored circles from emoji, and i like it generally. There is no strict requirements about tags in the task. And no post filtering required to implement in the task, the tags can be not clickable. So it is ok at the moment. But check and keep in mind this -->
      <br>
      <label for="postContent">Post Content:</label>
      <textarea id="postContent" v-model="postContent" required></textarea>
      
      <br>
      <label for="postPrivacy">Post Privacy:</label>
      <br>
      <input type="radio" id="public" name="postPrivacy" value="public" v-model="postPrivacy">
      <label for="public">Public - for all users</label>
      <br>
      <input type="radio" id="private" name="postPrivacy" value="private" v-model="postPrivacy">
      <label for="private">Private - for all followers</label>
      <br>
      <input type="radio" id="almostPrivate" name="postPrivacy" value="almostPrivate" v-model="postPrivacy">
      <label for="almostPrivate">Almost Private - for selected followers</label>
      <br>
      <select v-if="postPrivacy === 'almostPrivate'" multiple v-model="selectedFollowers">
        <option v-for="follower in followers" :key="follower.email" :value="follower.email">{{ follower.full_name }} ({{ follower.email }})</option>
      </select>
      
      <div>
        <label for="picture">Picture:</label>
        <input type="file" id="picture" accept="image/jpeg, image/png, image/gif" @change="handlePictureChange">
        <div class="optional">(optional)</div>
      </div>

      <br>
      <button type="submit">Submit</button>
      <!-- todo: add image or gif to post required in task. Perhaps, to prevent posting "anacondas" and "caves" photos, the images can be limited from allowed lists of images, but generally it sounds like they expect any image upload, which is unsafe, like in picture too -->
    </form>
    <div v-if="pictureStore.pictureError">{{ pictureStore.pictureError }}</div>
  </div>
  <div>
    <h2>Posts:</h2>
    <!-- add posts list , already created -->
    <div v-for="post in postsList"
      :key="post.id">
      <hr>
      <router-link
      :to="{ name: 'post' }"
      @click="piniaManageData(post)">
        <p>Post Author id: {{ post.author_id }}</p>
        <h3>Post Author: {{ post.author_full_name }}</h3>
        <h3>Post Author email:{{ post.author_email }}</h3>
        <p>Post id: {{ post.id }}</p>
        <p>Post title: {{ post.title }}</p>
        <p>Post tags: {{ post.tags }}</p>
        <p>Post content: {{ post.content }}</p>
        <p>Post privacy: {{ post.privacy }}</p><!-- todo: no need to display -->
        <p>Post picture: {{ post.picture }}</p>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Ref, onMounted, ref } from 'vue';
import { usePostStore } from '@/store/pinia';
import { usePictureStore } from '@/store/pinia';
import { useWebSocketStore } from '@/store/websocket';
import { WSMessage, WSMessageType, PostSubmit, Post } from '@/api/types';

interface Follower {
  full_name: string;
  email: string;
}

//todo: remove/refactor later, dummy data, must be collected from backend
function getPosts() {
  const posts: Post[] = [
    { id: 1, author_id: 11, author_full_name: 'John Doe 11', author_email: "email11@mail.com", title: "Dummy post title", tags: "dummy, post, 111", content: 'Dummy Post content text.', privacy: 'public' },
    { id: 2, author_id: 22, author_full_name: 'Jane Doe 22', author_email: "email22@mail.com", title: "Dummy post title", tags: "dummy, post, 222", content: 'Dummy Post content text.', privacy: 'private' },
  ];
  return posts;
}

const webSocketStore = useWebSocketStore();
const postsList = webSocketStore.postsList;
// const postsList = ref(getPosts());

const postTitle = ref('');
const postTags = ref('');
const postContent = ref('');
const postPrivacy = ref('public');
const selectedFollowers = ref<string[]>([]);
const followers = ref<Follower[]>([]);
// const picture: Ref<Blob | null> = ref(null); //todo: chat gpt solution, to fix null value case, because field is optional

const pictureStore = usePictureStore();
function handlePictureChange(event: Event) {
  pictureStore.handlePictureUpload(event);
  // picture.value = (event.target as HTMLInputElement).files?.[0] ?? null;
}

//todo: refactor to send post to backend
async function addPost() {
  //todo: first send to backend, and if success, then add to postsList on page, so all the data finally must be from backend. The "postTitle.value" etc must be used in request to backend, to create new post. The authorId is the current user id, f.e. managed using pinia storage, or cookies, not sure yet.

  const picture = await pictureStore.getBase64forJson
  const postSubmit: PostSubmit = {
    user_uuid: "123", //todo: get from pinia storage, or cookies, or backend

    title: postTitle.value,
    tags: postTags.value, //todo: comma separated tags, but for dummy case just string on screen
    content: postContent.value,
    privacy: postPrivacy.value,
    picture: picture,
  };

  if (postPrivacy.value === 'almostPrivate') {
    postSubmit.followers = selectedFollowers.value;
  }

  const message: WSMessage = {
    type: WSMessageType.POST_SUBMIT,
    data: postSubmit,
  };
  webSocketStore.sendMessage(message);

  // postsList.value.unshift(post);

  postTitle.value = '';
  postTags.value = '';
  postContent.value = '';
  postPrivacy.value = 'public';
  selectedFollowers.value = [];
  // picture = "null";

}

//todo: remove/refactor later, dummy data, must be collected from backend
function updateFollowersList() {
  followers.value = [
    { full_name: 'John Doe 11', email: 'John_Doe@mail.com' },
    { full_name: 'Jane Doe 22', email: 'Jane_Doe@mail.com' },
    { full_name: 'Sir Flex 33', email: 'Sir_Flex@mail.com' },
  ];
}

const postStore = usePostStore();
function piniaManageData(post: Post) {
  postStore.setPostId(post.id);
}

onMounted(() => {
  updateFollowersList();
});

</script>