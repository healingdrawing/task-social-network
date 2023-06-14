<template>
  <div>
    <h1>Post:</h1>
    <h2>id: {{ postId }}</h2>
    <h2>author: {{ postAuthor }}</h2>
    <h2>title: {{ postTitle }}</h2>
    <h2>content: {{ postContent }}</h2>
  </div>
  <div>
    <h2>Comments:</h2>
    <!-- add new comment using text area -->
    <div>
      <form @submit.prevent="addComment">
        <label for="commentContent">Comment Content:</label>
        <textarea id="commentContent" v-model="commentContent" required></textarea>
        <br>
        <button type="submit">Submit</button>
        <!-- todo: add image or gif to comment required in task. Perhaps, to prevent posting "anacondas" and "caves" photos, the images can be limited from allowed lists of images, but generally it sounds like they expect any image upload, which is unsafe, like in avatar too -->
      </form>
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
    </div>
  </div>
  
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { usePostStore } from '@/store/pinia';
import { useProfileStore } from '@/store/pinia';

const postStore = usePostStore();

// todo: plan to use this view(PostView.vue) like target view from both "ProfileView.vue" and "PostView.vue", when user click on post title, it will be redirected to "PostView.vue" and this postId will be used to get post from backend, and show it on "PostView.vue". So "userId" must be provided some way(pinia storage or cookie/session not sure) too, to create comment for this post

// todo: refactor to get post from backend , and update postTitle, and postContent(include some way images, which is part of task)
function updateFullPost() {
  postId.value = postStore.getPostId;
  // todo: get post from backend, using postId, and update x3 data bottom
  postAuthor.value = "must be collected from backend later inside updateFullPost()";
  postTitle.value = "must be collected from backend later inside updateFullPost()";
  postContent.value = "must be collected from backend later inside updateFullPost()";
  console.log('postId inside updatePost PostView.vue ', postId);
}

// const postId = computed(() => profileStore.getPostId); // reactivity syntax
/** managed by pinia storage, to keep links in browser as possible short, just "/post" */
const postId = ref(-1); //no need reactivity in this case
const postAuthor = ref(''); //todo: full name of author
const postTitle = ref('');
const postContent = ref(''); //todo: not sure it can be just string, because the images can be part of post content

interface Comment {
  id: number; // comment id, unique, autoincrement, primary key, all comments must be stored one table in database
  authorId: number; //todo: need to implement clickable link to user profile
  authorFullName: string; //todo: need to implement clickable link to user profile
  content: string;
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
  };
  commentsList.value.unshift(comment);
  commentContent.value = '';

}

const profileStore = useProfileStore();
function piniaManageData(comment: Comment) {
  profileStore.setTargetUserId(comment.authorId);
}

onMounted(() => {
  updateFullPost();
});

</script>
