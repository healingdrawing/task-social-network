<template>
  <div>
    <h1>Post:</h1>
    <h2>id: {{ postId }}</h2>
    <h2>title: {{ postTitle }}</h2>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { useProfileStore } from '@/store/pinia';

const profileStore = useProfileStore();

function updatePost(postId: number) {
  postTitle.value = "must be collected from backend later inside updatePost()";
  console.log('postId inside updatePost PostView.vue ', postId);
}

// const postId = computed(() => profileStore.getPostId);
const postId = profileStore.getPostId; //no need reactivity in this case
const postTitle = ref('');


onMounted(() => {
  updatePost(profileStore.getPostId); //no needed at this point, just test
});

</script>
