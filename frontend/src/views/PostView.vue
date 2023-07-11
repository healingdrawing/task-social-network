<template>
  <div>
    <h1>Post:</h1>
    <div>
      <p>Post id: {{ post.id }}</p>
      <p>Post title: {{ post.title }}</p>
      <p>Post tags: {{ post.categories }}</p>
      <p>Post content: {{ post.content }}</p>
      <p>Post privacy: {{ post.privacy }}</p><!-- todo: no need to display -->
      <p>Post created: {{ post.created_at }}</p>
      <div v-if="post.picture !== ''">
        <p>Post picture: 
          <br> <img :src="`data:image/jpeg;base64,${post.picture}`" alt="picture" />
        </p>
      </div>
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
  <div>
    <h2>Comments:</h2>
    <!-- add new comment using text area -->
    <div>
      <form @submit.prevent="addComment">
        <label for="commentContent">Comment Content:</label>
        <textarea id="commentContent" v-model="commentContent" required></textarea>

        <div>
          <label for="picture">Picture:</label>
          <input type="file" id="picture" accept="image/jpeg, image/png, image/gif" @change="handlePictureChange">
          <div class="optional">(optional)</div>
        </div>

        <br>
        <button type="submit">Submit</button>
        <!-- todo: add image or gif to comment required in task. Perhaps, to prevent posting "anacondas" and "caves" photos, the images can be limited from allowed lists of images, but generally it sounds like they expect any image upload, which is unsafe, like in avatar too -->
      </form>
      <div v-if="pictureStore.pictureError">{{ pictureStore.pictureError }}</div>
    </div>
    <!-- add comments list , already created -->
    <div v-for="comment in commentsList"
      :key="comment.created_at">
      <hr>
      <p>Comment content: {{ comment.content }}</p>
      <div v-if="comment.picture !== ''">
        <p>Comment picture: 
          <img :src="`data:image/jpeg;base64,${comment.picture}`" alt="picture" />
        </p>
      </div>
      <router-link
      :to="{ name: 'target' }"
      @click="piniaManageDataProfile(comment.email)">
      <h6>Comment Author:
        {{ comment.first_name }} 
        {{ comment.last_name }} 
        ({{ comment.email }})
      </h6>
      </router-link>
      
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, Ref, computed } from 'vue';
import { useWebSocketStore } from '@/store/pinia';
import { useUUIDStore } from '@/store/pinia';
import { usePostStore } from '@/store/pinia';
import { useProfileStore } from '@/store/pinia';
import { usePictureStore } from '@/store/pinia';
import { CommentsListRequest, CommentSubmit, WSMessage, WSMessageType } from '@/api/types';

const wss = useWebSocketStore();
const UUIDStore = useUUIDStore();

const profileStore = useProfileStore();
function piniaManageDataProfile(email: string) {
  profileStore.setTargetUserEmail(email);
}

const picture: Ref<Blob | null> = ref(null); //todo: chat gpt solution, to fix null value case, because field is optional
const pictureStore = usePictureStore();
function handlePictureChange(event: Event) {
  pictureStore.handlePictureUpload(event);
  picture.value = (event.target as HTMLInputElement).files?.[0] ?? null;
}

const postStore = usePostStore();
const post = computed(() => postStore.getPost);

// todo: refactor to get comments from backend, using post_id
function updatePostComments() {
  wss.sendMessage({
    type: WSMessageType.COMMENTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      post_id: post.value.id,
    } as CommentsListRequest,
  });
  // todo: send message through websocket to refresh comments list
}

const commentsList = computed(() => wss.commentsList);

const commentContent = ref('');

function addComment() {
  const commentSubmit: CommentSubmit = {
    user_uuid: UUIDStore.getUUID,
    post_id: post.value.id, // idiotic gap, because golang can cast properly only strings. facepalm
    content: commentContent.value,
    picture: pictureStore.getPictureBase64String,
  };

  const message: WSMessage = {
    type: WSMessageType.COMMENT_SUBMIT,
    data: commentSubmit,
  };
  wss.sendMessage(message);

  commentContent.value = '';
  picture.value = null;
  (document.getElementById("picture") as HTMLInputElement).value = "";
  pictureStore.resetPicture()

}

onMounted(() => {
  updatePostComments();
});

</script>
