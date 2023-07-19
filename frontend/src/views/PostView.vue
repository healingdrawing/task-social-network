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
    <!-- add new comment using text area -->
    <hr>
    <div>
      <form @submit.prevent="addComment">
        <label for="commentContent"> Create Comment </label>
        <br> <textarea id="commentContent" v-model="commentContent" required></textarea>
        <div>
          <label for="picture"> with picture(optional): </label>
          <br> <input type="file" id="picture" accept="image/jpeg, image/png, image/gif" @change="handlePictureChange">
        </div>
        <br>
        <button type="submit">Submit</button>
      </form>
      <div v-if="pictureStore.pictureError">{{ pictureStore.pictureError }}</div>
    </div>
    <!-- add comments list , already created -->
    <h2>Comments:</h2>
    <div v-for="comment in commentsList"
      :key="comment.created_at">
      <hr>
      <p>Comment content: {{ comment.content }}</p>
      <div v-if="comment.picture !== ''">
        <p>Comment picture: 
          <br> <img :src="`data:image/jpeg;base64,${comment.picture}`" alt="picture" />
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

const picture: Ref<Blob | null> = ref(null);
const pictureStore = usePictureStore();
function handlePictureChange(event: Event) {
  pictureStore.handlePictureUpload(event);
  picture.value = (event.target as HTMLInputElement).files?.[0] ?? null;
}

const postStore = usePostStore();
const post = computed(() => postStore.getPost);

function updatePostComments() {
  wss.sendMessage({
    type: WSMessageType.COMMENTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      post_id: post.value.id,
    } as CommentsListRequest,
  });
}

const commentsList = computed(() => wss.commentsList);

const commentContent = ref('');

function addComment() {
  const commentSubmit: CommentSubmit = {
    user_uuid: UUIDStore.getUUID,
    post_id: post.value.id,
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

onMounted(async () => {
  wss.refresh_websocket()
  await wss.waitForConnection();
  updatePostComments();
});

</script>
