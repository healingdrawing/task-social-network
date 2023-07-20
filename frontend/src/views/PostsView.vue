<template>
  <div>
    <h1>Create Post</h1>

    <div><button type="button" @click="crap" title="remove in production">Fill Debug / remove later</button></div> <!-- todo: remove later -->

    <form @submit.prevent="addPost">
      <label for="postTitle">Post Title:</label>
      <br> <input type="text" id="postTitle" v-model="postTitle" required>
      <br>
      <label for="postTags">Post Tags:</label>
      <br> <input title="comma separated" type="text" id="postTags" v-model="postTags">
      <br>
      <label for="postContent">Post Content:</label>
      <br> <textarea id="postContent" v-model="postContent" required></textarea>
      
      <br>
      <label for="postPrivacy">Post Privacy:</label>
      <br>
      <input type="radio" id="public" name="postPrivacy" value="public" v-model="postPrivacy">
      <label for="public">Public - for all users</label>
      <br>
      <input type="radio" id="private" name="postPrivacy" value="private" v-model="postPrivacy">
      <label for="private">Private - for all followers</label>
      <br>
      <input type="radio" id="almost_private" name="postPrivacy" value="almost private" v-model="postPrivacy">
      <label for="almost_private">Almost Private - for selected followers</label>
      <br>
      <select v-if="postPrivacy === 'almost private'" multiple v-model="selectedFollowers" class="users_list_with_scroll">
        <option v-for="follower in followersList" :key="follower.email" :value="follower.email">{{ follower.first_name }} {{ follower.last_name }} ({{ follower.email }})</option>
      </select>
      
      <div>
        <br>
        <label for="picture" class="label_file_upload">
          with picture(optional):
          <input type="file" id="picture" accept="image/jpeg, image/png, image/gif" @change="handlePictureChange">
        </label>
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
  </div>
</template>

<script lang="ts" setup>
import { Ref, computed, onMounted, ref } from 'vue';
import { useUUIDStore } from '@/store/pinia';
import { usePostStore } from '@/store/pinia';
import { useProfileStore } from '@/store/pinia';
import { usePictureStore } from '@/store/pinia';
import { useWebSocketStore } from '@/store/websocket';
import { WSMessage, WSMessageType, PostSubmit, Post, PostsListRequest, TargetProfileRequest } from '@/api/types';

const wss = useWebSocketStore();
const postsList = computed(() => wss.postsList); // ref and reactive failed to work here, so computed used. Straight way put wss.postsList to template works too,


const postTitle = ref('');
const postTags = ref('');
const postContent = ref('');
const postPrivacy = ref('public');
const selectedFollowers = ref<string[]>([]);
const picture: Ref<Blob | null> = ref(null);

const pictureStore = usePictureStore();
function handlePictureChange(event: Event) {
  pictureStore.handlePictureUpload(event);
  picture.value = (event.target as HTMLInputElement).files?.[0] ?? null;
}

const UUIDStore = useUUIDStore();
async function addPost() {
  const postSubmit: PostSubmit = {
    user_uuid: UUIDStore.getUUID,

    title: postTitle.value,
    categories: postTags.value,
    content: postContent.value,
    privacy: postPrivacy.value,
    able_to_see: selectedFollowers.value.join(' '), //list of emails, separated by space
    picture: pictureStore.getPictureBase64String,
  };

  const message: WSMessage = {
    type: WSMessageType.POST_SUBMIT,
    data: postSubmit,
  };
  wss.sendMessage(message);

  postTitle.value = '';
  postTags.value = '';
  postContent.value = '';
  postPrivacy.value = 'public';
  selectedFollowers.value = [];
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

const followersList = computed(() => wss.userFollowersList);
function updateFollowersList() {
  wss.sendMessage({
    type: WSMessageType.USER_FOLLOWERS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      target_email: profileStore.getUserEmail,
    } as TargetProfileRequest,
  })
}

// send request to get old posts list, used inside onMounted
function updatePostsList() {
  wss.sendMessage({
    type: WSMessageType.POSTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
    } as PostsListRequest,
  });
}

onMounted(async () => {
  wss.refresh_websocket()
  await wss.waitForConnection();
  updatePostsList();
  updateFollowersList();
});

const crap = () => {
  postTitle.value = 'Dummy post title';
  postTags.value = 'dummy, post, 111, test';
  postContent.value = 'Dummy Post content text.';
  postPrivacy.value = 'public';
}

</script>