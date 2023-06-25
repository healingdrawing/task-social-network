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
        <option v-for="follower in followers" :key="follower.id" :value="follower.id">{{ follower.name }}</option>
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
        <p>Post Author id: {{ post.authorId }}</p>
        <h3>Post Author: {{ post.authorFullName }}</h3>
        <p>Post id: {{ post.id }}</p>
        <p>Post title: {{ post.title }}</p>
        <p>Post tags: {{ post.tags }}</p>
        <p>Post content: {{ post.content }}</p>
        <p>Post privacy: {{ post.privacy }}</p><!-- todo: no need to display -->
        <p>Post followers: {{ post.followers }}</p> <!-- todo: no need to display it of course, it is used on backend side, before return post as visible or not -->
        <p>Post picture: {{ post.picture }}</p> <!-- todo: no need to display it of course, it is used on backend side, before return post as visible or not -->
      </router-link>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Ref, onMounted, ref } from 'vue';
import { usePostStore } from '@/store/pinia';
import { usePictureStore } from '@/store/pinia';
import { useWebSocketStore } from '@/store/websocket';
import { Message, MessageType, Post } from '@/api/types';

interface Follower {
  id: number;
  name: string;
}

//todo: remove/refactor later, dummy data, must be collected from backend
function getPosts() {
  const posts: Post[] = [
    { id: 1, authorId: 11, authorFullName: 'John Doe 11', title: "Dummy post title", tags: "dummy, post, 111", content: 'Dummy Post content text.', privacy: 'public' },
    { id: 2, authorId: 22, authorFullName: 'Jane Doe 22', title: "Dummy post title", tags: "dummy, post, 222", content: 'Dummy Post content text.', privacy: 'private' },
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
const selectedFollowers = ref([]);
const followers = ref<Follower[]>([]);
const picture: Ref<Blob | null> = ref(null); //todo: chat gpt solution, to fix null value case, because field is optional

const pictureStore = usePictureStore();
function handlePictureChange(event: Event) {
  pictureStore.handlePictureUpload(event);
  picture.value = (event.target as HTMLInputElement).files?.[0] ?? null;
}

//todo: refactor to send post to backend
function addPost() {
  //todo: first send to backend, and if success, then add to postsList on page, so all the data finally must be from backend. The "postTitle.value" etc must be used in request to backend, to create new post. The authorId is the current user id, f.e. managed using pinia storage, or cookies, not sure yet.
  const post: Post = {
    id: -1, //todo: must be collected from backend
    authorId: -1, //todo: must be the current user id , returned from backend
    authorFullName: "must be the current user full name",
    title: postTitle.value,
    tags: postTags.value, //todo: comma separated tags, but for dummy case just string on screen
    content: postContent.value,
    privacy: postPrivacy.value,
    picture: picture.value,
  };

  if (postPrivacy.value === 'almostPrivate') {
    post.followers = selectedFollowers.value;
  }

  const message: Message = {
    messageType: MessageType.POST_SUBMIT,
    content: JSON.stringify(post),
  };
  webSocketStore.sendMessage(message);

  // postsList.value.unshift(post);

  postTitle.value = '';
  postTags.value = '';
  postContent.value = '';
  postPrivacy.value = 'public';
  selectedFollowers.value = [];
  picture.value = null;

}

//todo: remove/refactor later, dummy data, must be collected from backend
function updateFollowersList() {
  followers.value = [
    { id: 11, name: 'John Doe 11' },
    { id: 22, name: 'Jane Doe 22' },
    { id: 33, name: 'Sir Flex 33' },
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