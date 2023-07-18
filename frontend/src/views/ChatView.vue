<template>
  <div class="chat-view">
    <h3>
      Rap Battle Versus
      <br> {{ chat.first_name }} {{ chat.last_name }} ({{ chat.email }})
    </h3>
    <h1> Let's get it on! </h1>
    <div>
      <textarea v-model="messageText" @keydown.enter.ctrl.prevent="sendMessage" required></textarea>
      <br>
      <button @click="sendMessage">Send</button>
    </div>
    <div class="messages">
      <div class="message" v-for="message in messages_list" :key="message.created_at">
        <hr>
        <router-link
          :to="{ name: 'target' }"
          @click="piniaManageData(message)">            
          {{ message.first_name }} {{ message.last_name }} ({{ message.email }})
        </router-link>
        <br>
        {{ message.content }}
      </div>
    </div>
  </div>
</template>

<style>
.message {
  white-space: pre-wrap;
  overflow: auto;
  /* add a vertical scrollbar when content exceeds height */
}
</style>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue';
import { useWebSocketStore } from '@/store/websocket';
import { useUUIDStore } from '@/store/uuid';
import { useChatStore } from '@/store/chat';
import { useProfileStore } from '@/store/profile';
import { PrivateChatMessage } from '@/api/types';

const wss = useWebSocketStore();
const UUIDStore = useUUIDStore();
const chatStore = useChatStore();
const profileStore = useProfileStore();

const messages_list = computed(() => wss.private_chat_messages_list);
const messageText = ref('');
const chat = chatStore.get_target_user;

function sendMessage() {
  if (messageText.value.trim() === '') {
    return;
  }

  wss.send_private_chat_message(
    messageText.value,
    UUIDStore.getUUID,
  );

  messageText.value = '';
}

function piniaManageData(message: PrivateChatMessage) {
  profileStore.setTargetUserEmail(message.email);
}

onMounted(() => {
  wss.set_group_chat_id(0)
  wss.set_private_chat_user_id(chat.user_id)
});
</script>