<template>
  <div>
    <h1>Chats:</h1>
    <div class=" users_list_with_scroll chats_list">
      <div v-for="user in users_list" :key="user.email">
        <br>
        <router-link :to="{ name: 'chat' }" @click="piniaManageData(user)">
          <div class="router_link_box">
            {{ user.first_name }} {{ user.last_name }} ({{ user.email }})
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<style>
.chats_list {
  height: 300px;
  /* set the height of the list */
  overflow-y: auto;
  /* add a vertical scrollbar when content exceeds height */
}
</style>

<script lang="ts" setup>
import { onMounted, computed } from 'vue'
import { useWebSocketStore } from '@/store/pinia'
import { useUUIDStore } from '@/store/pinia'
import { useProfileStore } from '@/store/pinia'
import { useChatStore } from '@/store/chat';
import { TargetProfileRequest, WSMessageType, UserForChatList } from '@/api/types';

const wss = useWebSocketStore()
const UUIDStore = useUUIDStore()
const profileStore = useProfileStore()
const chatStore = useChatStore()

const users_list = computed(() => wss.private_chat_users_list)


/** save data in pinia storage, to use inside "ChatView.vue", to do not have deal with params */
function piniaManageData(user: UserForChatList) {
  chatStore.set_target_user(user)
}

function update_private_chat_users_list() {
  wss.sendMessage({
    type: WSMessageType.PRIVATE_CHAT_USERS_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      target_email: profileStore.getUserEmail,
    } as TargetProfileRequest,
  })
}

onMounted(async () => {
  wss.refresh_websocket()
  await wss.waitForConnection();
  update_private_chat_users_list()
})
</script>
