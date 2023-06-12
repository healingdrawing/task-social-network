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
      </form>
    </div>
    <!-- add comments list , already created -->
    <div v-for="comment in commentsList"
      :key="comment.id">
      <h3>Comment Author: {{ comment.authorFullName }}</h3>
      <p>Comment id: {{ comment.id }}</p>
      <p>Comment content: {{ comment.content }}</p>
    </div>
  </div>
  
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { useProfileStore } from '@/store/pinia';

const profileStore = useProfileStore();

// todo: plan to use this view(PostView.vue) like target view from both "ProfileView.vue" and "PostView.vue", when user click on post title, it will be redirected to "PostView.vue" and this postId will be used to get post from backend, and show it on "PostView.vue". So "userId" must be provided some way(pinia storage or cookie/session not sure) too, to create comment for this post

// todo: refactor to get post from backend , and update postTitle, and postContent(include some way images, which is part of task)
function updateFullPost(postId: number) {
  // todo: get post from backend, using postId, and update x3 data bottom
  postAuthor.value = "must be collected from backend later inside updateFullPost()";
  postTitle.value = "must be collected from backend later inside updateFullPost()";
  postContent.value = "must be collected from backend later inside updateFullPost()";
  console.log('postId inside updatePost PostView.vue ', postId);
}

// const postId = computed(() => profileStore.getPostId); // reactivity syntax
/** managed by pinia storage, to keep links in browser as possible short, just "/post" */
const postId = profileStore.getPostId; //no need reactivity in this case
const postAuthor = ref(''); //todo: full name of author
const postTitle = ref('');
const postContent = ref(''); //todo: not sure it can be just string, because the images can be part of post content

interface Comment {
  id: number;
  authorFullName: string; //todo: need to implement clickable link to user profile
  content: string;
}

//todo: remove/refactor later, dummy data, must be collected from backend
function getComments() {
  const comments: Comment[] = [
    { id: 1, authorFullName: 'John Doe', content: 'Dummy comment.', },
    { id: 2, authorFullName: 'Jane Doe', content: 'Dummy comment.', },
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
    authorFullName: commentAuthorFullName.value,
    content: commentContent.value,
  };
  commentsList.value.unshift(comment);
  commentContent.value = '';

}

onMounted(() => {
  updateFullPost(profileStore.getPostId); //no needed at this point, just test of pinia storage implementation
});

</script>
