<template>
  <div class="chat-view">
    <h1>Group: {{ group.name }}</h1>
    <h1>Group Chat with id: {{ group.id }}</h1>
    <div>
      <textarea v-model="messageText" @keydown.enter.ctrl.prevent="sendMessage" required></textarea>
      <br>
      <button @click="sendMessage">Send</button>
    </div>
    <div class="messages">
      <div class="message single_div_box " v-for="message in messages_list" :key="message.created_at">
        <br>
        <router-link
          :to="{ name: 'target' }"
          @click="piniaManageData(message)">
          <div class="router_link_box">
            {{ message.first_name }} {{ message.last_name }} ({{ message.email }})
          </div>
        </router-link>
        <br> <br> <br>
        {{ message.content }}
      </div>
    </div>
  </div>
</template>

<style>
.message {
  white-space: pre-wrap;
}
</style>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue';
import { useWebSocketStore } from '@/store/pinia';
import { useUUIDStore } from '@/store/pinia';
import { useGroupStore } from '@/store/pinia';
import { useProfileStore } from '@/store/pinia';
import { GroupChatMessage } from '@/api/types';

const wss = useWebSocketStore();
const UUIDStore = useUUIDStore();
const groupStore = useGroupStore();
const profileStore = useProfileStore();

const group = groupStore.getGroup

const messages_list = computed(() => wss.group_chat_messages_list);

const messageText = ref('');

function sendMessage() {
  if (messageText.value.trim() === '') {
    messageText.value = '';
    return;
  }

  wss.send_group_chat_message(
    messageText.value,
    UUIDStore.getUUID,
  )

  messageText.value = '';

  // scrollToBottom(); // commented, because now the chat history placed under the message box
}

function piniaManageData(message: GroupChatMessage) {
  profileStore.setTargetUserEmail(message.email);
}

onMounted(async () => {
  wss.refresh_websocket()
  await wss.waitForConnection();
});
</script>