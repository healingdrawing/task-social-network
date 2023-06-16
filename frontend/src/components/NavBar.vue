<template>
  <div class="nav-bar">
    <div class="nav-bar__logo">
      <img src="../assets/logo.png" alt="Vue logo" />
    </div>
    <div class="nav-bar__links">
      <!-- todo: add implementation to mark/hightlight "Profile" and "Chats" -->
      <router-link :class="{ highlighted: hasNewBells }" to="/profile">Profile</router-link> |
      <router-link to="/posts" @click="piniaManageData()">Posts</router-link> |
      <router-link :class="{ highlighted: true }" to="/chats">Chats</router-link> |
      <router-link to="/groups">Groups</router-link> |
      <router-link to="/">Logout</router-link>
    </div>
    <router-view/>
  </div>
</template>

<style scoped>
.highlighted {
  background-color: gold;
}
</style>

<script lang="ts" setup>
import { computed } from 'vue';
import { useBellStore } from '@/store/bell';
import { useGroupStore } from '@/store/group';

// when "posts" click happens, reset group id to -1 or 0, to prevent backend filtering of the posts to not show group only posts, but show all
const groupStore = useGroupStore();
function piniaManageData() {
  alert('piniaManageData posts NavBar.vue');
  groupStore.setGroupId(-1); //todo: implement on backend. Now posts will be not filtered by group id. But only filtered by date, fresh first
}

const bellStore = useBellStore();
const hasNewBells = computed(() => bellStore.bells.length > 0);
</script>
