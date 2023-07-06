<template>
  <div>
    <h1>Create Post:</h1>

    <div><hr><button type="button" @click="crap" title="remove in production">Fill Debug / remove later</button><hr></div> <!-- todo: remove later -->

    <form @submit.prevent="addPost">
      <label for="postTitle">Post Title:</label>
      <input type="text" id="postTitle" v-model="postTitle" required>
      <br>
      <label for="postTags">Post Tags:</label>
      <input type="text" id="postTags" v-model="postTags">
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
import { usePostStore } from '@/store/pinia';
import { useProfileStore } from '@/store/pinia';
import { usePictureStore } from '@/store/pinia';
import { useWebSocketStore } from '@/store/websocket';
import { WSMessage, WSMessageType, PostSubmit, Post, PostsListRequest, TargetProfileRequest } from '@/api/types';
import { onBeforeRouteLeave, useRoute } from 'vue-router';

interface Follower {
  full_name: string;
  email: string;
}

const wss = useWebSocketStore();
const postsList = computed(() => wss.postsList); // ref and reactive failed to work here, so computed used. Straight way put webSocketStore.postsList to template works too,

const postTitle = ref('');
const postTags = ref('');
const postContent = ref('');
const postPrivacy = ref('public');
const selectedFollowers = ref<string[]>([]);
const followers = ref<Follower[]>([]);
const picture: Ref<Blob | null> = ref(null); //todo: chat gpt solution, to fix null value case, because field is optional

const pictureStore = usePictureStore();
function handlePictureChange(event: Event) {
  pictureStore.handlePictureUpload(event);
  picture.value = (event.target as HTMLInputElement).files?.[0] ?? null;
}

const storeUUID = useUUIDStore();
async function addPost() {
  const postSubmit: PostSubmit = {
    user_uuid: storeUUID.getUUID,

    title: postTitle.value,
    categories: postTags.value,
    content: postContent.value,
    privacy: postPrivacy.value,
    able_to_see: selectedFollowers.value.join(' '), //list of emails, separated by space
    picture: pictureStore.getPictureBase64String,
  };

  // if (postPrivacy.value === 'almostPrivate') {
  //   postSubmit.able_to_see = selectedFollowers.value.join(' ');
  // }

  const message: WSMessage = {
    type: WSMessageType.POST_SUBMIT,
    data: postSubmit,
  };
  wss.sendMessage(message);

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
function piniaManageDataPost(post: Post) {
  postStore.setPostId(post.id);
}

const profileStore = useProfileStore();
function piniaManageDataProfile(email: string) {
  profileStore.setTargetUserEmail(email);
}

// send request to get old posts list, used inside onMounted
function updatePostsList() {
  console.log('=======FIRED======= updatePostsList');

  wss.sendMessage({
    type: WSMessageType.POSTS_LIST,
    data: {
      user_uuid: storeUUID.getUUID,
    } as PostsListRequest,
  });
}

onMounted(() => {
  updatePostsList();
  updateFollowersList();
});

const crap = () => {
  postTitle.value = 'Dummy post title';
  postTags.value = 'dummy, post, 111, test';
  postContent.value = 'Dummy Post content text.';
  postPrivacy.value = 'public';
}

watch(postsList, (newVal, oldVal) => {
  console.log('Posts list:', newVal);
});

onBeforeRouteLeave((to, from, next) => {
  wss.facepalm();
  // wss.clearOnBeforeRouteLeave(from.path) // todo: clear the messages. ugly but works
  console.log('onBeforeRouteLeave fired');
  next();
});

</script>