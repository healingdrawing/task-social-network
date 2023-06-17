<template>
  <div>
    <h1>Chats</h1>
    <ul class="chats-list">
      <li v-for="oponent in oponents" :key="oponent.userId">
        <router-link :to="{ name: 'chat' }" @click="piniaManageData(oponent.chatId)">
          <span v-if="oponent.online">😀</span>
          <span v-else>😴</span>
          {{ oponent.fullName }} ({{ oponent.email }})
        </router-link>
      </li>
    </ul>
  </div>
</template>

<style>
.chats-list {
  height: 300px;
  /* set the height of the list */
  overflow-y: auto;
  /* add a vertical scrollbar when content exceeds height */
}
</style>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useChatStore } from '@/store/chat';

const chatStore = useChatStore()

interface Oponent {
  userId: number
  email: string
  fullName: string
  online: boolean
  chatId: number
}

const oponents = ref<Oponent[]>([])

/** save data in pinia storage, to use inside "ChatView.vue", to do not have deal with params */
function piniaManageData(chatId: number) {
  chatStore.setChatId(chatId)
}

function updateOponents() {
  //todo: get oponents from backend, and create new private chat if not exists
  //dummy data
  oponents.value = [
    { userId: 1, email: 'dummy1@mail.com', fullName: 'Dummy User1', online: true, chatId: 1, },
    { userId: 2, email: 'dummy2@mail.com', fullName: 'Dummy User2', online: false, chatId: 2, },
    { userId: 3, email: 'dummy3@mail.com', fullName: 'Dummy User3', online: true, chatId: 3, },
  ]
}

onMounted(() => {
  updateOponents()
})
</script>
