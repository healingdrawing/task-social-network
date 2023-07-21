<template>
  <br>
  <router-link to="/group_posts">
    <div class="router_link_box">
      Back to Group Posts
    </div>  
  </router-link>
  <div>
    <h1>Group Post</h1>
    <div class="single_div_box">
      <br>
      <h3>Group Post title: </h3> <p> {{ group_post.title }} </p>
      <h3>Group Post tags: </h3> <p> {{ group_post.categories }} </p>
      <h3>Group Post content: </h3> <p> {{ group_post.content }} </p>
      <h3>Group Post created: </h3> <p> {{ group_post.created_at }} </p>
      <div v-if="group_post.picture !== ''">
        <h3>Group Post picture: </h3>
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
        <br>
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
const group_post = computed(() => postStore.getGroupPost);

function updatePostComments() {
  wss.sendMessage({
    type: WSMessageType.GROUP_POST_COMMENTS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      post_id: group_post.value.id,
    } as CommentsListRequest,
  });
}

const commentsList = computed(() => wss.commentsList);

const commentContent = ref('');

function addComment() {
  const group_post_comment_submit: CommentSubmit = {
    user_uuid: UUIDStore.getUUID,
    post_id: group_post.value.id,
    content: commentContent.value,
    picture: pictureStore.getPictureBase64String,
  };

  const message: WSMessage = {
    type: WSMessageType.GROUP_POST_COMMENT_SUBMIT,
    data: group_post_comment_submit,
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
