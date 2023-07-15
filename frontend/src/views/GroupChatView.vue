<template>
  <div class="chat-view">
    <h1>Group: {{ group.name }}</h1>
    <h1>Group Chat with id: {{ group.id }}</h1>
    <div>
      <textarea v-model="messageText" @keydown.enter.ctrl.prevent="sendMessage"></textarea>
      <br>
      <button @click="sendMessage">Send</button>
    </div>
    <div class="messages" ref="messages"><!--  :scrollTop="messages?.scrollHeight"> commented, because now the chat history placed under the message box -->
      <div class="message" v-for="message in messages_list" :key="message.created_at">
        <hr>
        <h6>
          <router-link
            :to="{ name: 'target' }"
            @click="piniaManageData(message)">            
            {{ message.first_name }} {{ message.last_name }} ({{ message.email }})
          </router-link>
        </h6>
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
import { ref, onMounted, computed } from 'vue';
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

const messages = ref<HTMLElement | null>(null);
/* 
function scrollToBottom() {
  if (messages?.value) {
    messages.value.scrollTop = messages.value.scrollHeight;
  }
}
 */

function piniaManageData(message: GroupChatMessage) {
  profileStore.setTargetUserEmail(message.email);
}

onMounted(() => {
  messages.value = document.querySelector('.messages'); // bind to div by class name

  // scrollToBottom(); // commented, because now the chat history placed under the message box
});
</script>