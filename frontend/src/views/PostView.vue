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
          <img :src="`data:image/jpeg;base64,${post.picture}`" alt="picture" />
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
      :key="comment.id">
      <hr>
      <p>Comment Author id: {{ comment.authorId }}</p>
      <router-link
      :to="{ name: 'target' }"
      @click="piniaManageData(comment)">
      <h3>Comment Author: {{ comment.authorFullName }}</h3>
      </router-link>
      <p>Comment id: {{ comment.id }}</p>
      <p>Comment content: {{ comment.content }}</p>
      <p>Comment picture: {{ comment.picture }}</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, Ref, computed } from 'vue';
import { usePostStore } from '@/store/pinia';
import { useProfileStore } from '@/store/pinia';
import { usePictureStore } from '@/store/pinia';

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

  // todo: send message through websocket to refresh comments list
}

interface Comment {
  id: number; // comment id, unique, autoincrement, primary key, all comments must be stored one table in database
  authorId: number; //todo: need to implement clickable link to user profile
  authorFullName: string; //todo: need to implement clickable link to user profile
  content: string;
  picture?: Blob | null;
}

//todo: remove/refactor later, dummy data, must be collected from backend
function getComments() {
  const comments: Comment[] = [
    { id: 1, authorId: 11, authorFullName: 'John Doe 11', content: 'Dummy comment.', },
    { id: 2, authorId: 22, authorFullName: 'Jane Doe 22', content: 'Dummy comment.', },
  ];
  return comments;
}
const commentsList = ref(getComments());

const commentContent = ref('');
const commentAuthorFullName = ref('');

// todo: refactor to add comment to backend
function addComment() {
  // todo: perpahs here call the backend to add comment, using user id, post id, and comment content, and get back the comment id, and author full name(to not provide it using pinia storage, which not looks too good).
  commentAuthorFullName.value = 'dummy author full name'
  const comment: Comment = {
    id: 3,
    authorId: 33,
    authorFullName: commentAuthorFullName.value,
    content: commentContent.value,
    picture: picture.value,
  };
  commentsList.value.unshift(comment);

  commentContent.value = '';
  picture.value = null;

}


function piniaManageData(comment: Comment) {
  profileStore.setTargetUserEmail("comment.authorEmail");
}

onMounted(() => {
  updatePostComments();
});

</script>
