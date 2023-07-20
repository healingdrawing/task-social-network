<template>
  <div>
    <h1>Post</h1>
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
      </div>
  </div>
  <div>
    <!-- add new comment using text area -->
    <div>
      <form @submit.prevent="addComment">
        <label for="commentContent"> Create Comment </label>
        <br> <textarea id="commentContent" v-model="commentContent" required></textarea>
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
    <!-- add comments list , already created -->
    <h2>Comments:</h2>
    <div v-for="comment in commentsList"
      :key="comment.created_at">
      <br v-if="comment.content.trim() !== '' || comment.picture !== ''">
      <div v-if="comment.content.trim() !== '' || comment.picture !== ''" class="single_div_box">
        <div v-if="comment.content.trim() !== ''">
          <h3> Comment content: </h3> <p> {{ comment.content }}</p>
        </div>
        <div v-if="comment.picture !== ''">
          <h3> Comment picture: </h3>
          <p>
            <br> <img :src="`data:image/jpeg;base64,${comment.picture}`" alt="picture" />
          </p>
        </div>
        <router-link
        :Title="comment.first_name + '\n' + comment.last_name + '\n' + comment.email"
        :to="{ name: 'target' }"
        @click="piniaManageDataProfile(comment.email)">
          <div class="router_link_box">
            visit author profile
          </div>
        </router-link>
      </div>
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
