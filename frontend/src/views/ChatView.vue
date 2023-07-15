<template>
  <div class="chat-view">
    <h1>Chat: {{ chatId }}</h1>
    <div>
      <textarea v-model="messageText" @keydown.enter.ctrl.prevent="sendMessage"></textarea>
      <br>
      <button @click="sendMessage">Send</button>
    </div>
    <div class="messages" ref="messages"><!--  :scrollTop="messages?.scrollHeight"> commented, because now the chat history placed under the message box -->
      <div class="message" v-for="message in messageList" :key="message.id">
        <!-- todo: implement clickable <router-link> to jump into message sender profile. Use @click="managePiniaData" -->
        <hr>
        <router-link
          :to="{ name: 'target' }"
          @click="piniaManageData(message)">            
          {{ message.fullName }} ({{ message.senderEmail }})
        </router-link>
        <br>
        {{ message.text }}
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
import { ref, onMounted } from 'vue';
import { useChatStore } from '@/store/chat';
import { useProfileStore } from '@/store/profile';

const chatStore = useChatStore();
const profileStore = useProfileStore();

interface Message {
  id: number; // unique id of message in database
  text: string; // here support to upload image/gif is no needed, according to task
  userId: number; // todo: use to implement clickable link to profile, f.e. group chat
  fullName: string;
  senderEmail: string;
}

const messageList = ref<Message[]>([]);
const messageText = ref('');
const chatId = ref(-1);

function sendMessage() {
  if (messageText.value.trim() === '') {
    return;
  }

  const message: Message = {
    id: Date.now(),
    text: messageText.value.trim(),
    userId: 99, // todo: use to implement clickable link to profile, f.e. group chat
    fullName: 'Dummy User 99',
    senderEmail: 'dummy99@mail.com', // because name do not have to be unique
  }

  messageList.value.unshift(message); // add new message to the beginning of the list
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

function piniaManageData(message: Message) {
  profileStore.setTargetUserEmail(message.senderEmail);
}

onMounted(() => {
  messages.value = document.querySelector('.messages'); // bind to div by class name
  chatId.value = chatStore.getChatId;
  // scrollToBottom(); // commented, because now the chat history placed under the message box
});
</script>