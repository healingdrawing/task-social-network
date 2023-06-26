import { defineStore } from 'pinia';
import { Message, MessageType, Post } from '@/api/types';


export const useWebSocketStore = defineStore({
  id: 'websocket',
  state: () => ({
    socket: null as WebSocket | null,
    messages: [] as Message[],
  }),
  actions: {
    connect() {
      this.socket = new WebSocket('ws://localhost:8080/ws');

      this.socket.onopen = () => {
        console.log('WebSocket connected');
      };

      this.socket.onclose = () => {
        console.log('WebSocket disconnected');
      };

      this.socket.onmessage = (event) => {
        console.log(`Received message: ${event.data}`);
        const message = JSON.parse(event.data) as Message;
        this.messages.push(message);
      };
    },
    sendMessage(message: Message) {
      const messageString = JSON.stringify(message);
      console.log(`Sending message json string: ${messageString}`);
      this.socket?.send(messageString);
    },
    disconnect() {
      this.socket?.close();
      this.socket = null;
    },
  },
  getters: {
    isConnected(): boolean {
      return this.socket !== null && this.socket.readyState === WebSocket.OPEN;
    },

    commentsList(): Message[] { return this.messages.filter((message) => message.type === MessageType.COMMENTS_LIST) },
    chatMessagesList(): Message[] { return this.messages.filter((message) => message.type === MessageType.CHAT_MESSAGES_LIST) },
    chatUsersList(): Message[] { return this.messages.filter((message) => message.type === MessageType.CHAT_USERS_LIST) },
    followRequestsList(): Message[] { return this.messages.filter((message) => message.type === MessageType.FOLLOW_REQUESTS_LIST) },
    postsList(): Post[] {
      const postsMessages = this.messages.filter((message) => message.type === MessageType.POSTS_LIST);
      const posts = postsMessages.map((message) => message.data as Post);
      return posts;
    },
    groupsList(): Message[] { return this.messages.filter((message) => message.type === MessageType.GROUPS_LIST) },
    groupPostsList(): Message[] { return this.messages.filter((message) => message.type === MessageType.GROUP_POSTS_LIST) },
    groupPostCommentsList(): Message[] { return this.messages.filter((message) => message.type === MessageType.GROUP_POST_COMMENTS_LIST) },
    groupRequestsList(): Message[] { return this.messages.filter((message) => message.type === MessageType.GROUP_REQUESTS_LIST) },
    groupEventsList(): Message[] { return this.messages.filter((message) => message.type === MessageType.GROUP_EVENTS_LIST) },
    groupEventParticipantsList(): Message[] { return this.messages.filter((message) => message.type === MessageType.GROUP_EVENT_PARTICIPANTS_LIST) },
    userFollowingList(): Message[] { return this.messages.filter((message) => message.type === MessageType.USER_FOLLOWING_LIST) },
    userFollowersList(): Message[] { return this.messages.filter((message) => message.type === MessageType.USER_FOLLOWERS_LIST) },
    userPostsList(): Message[] { return this.messages.filter((message) => message.type === MessageType.USER_POSTS_LIST) },

    // chatMessages(): Message[] {
    //   return this.messages.filter((message) => message.type === 'chat');
    // },
    // postMessages(): Message[] {
    //   return this.messages.filter((message) => message.type === 'post');
    // },
    // notificationMessages(): Message[] {
    //   return this.messages.filter((message) => message.type === 'notification');
    // },
  },
});

// Usage in a component
// <script lang="ts" setup >
// import { useWebSocketStore } from '@/stores/websocket';

// const socketStore = useWebSocketStore();

// socketStore.connect();

// // Send a chat message
// const chatMessage = { type: 'chat', text: 'Hello, world!', sender: 'Alice' };
// socketStore.sendMessage(chatMessage);

// // Send a post message
// const postMessage = { type: 'post', text: 'Check out my new blog post!', sender: 'Bob' };
// socketStore.sendMessage(postMessage);

// // Send a notification message
// const notificationMessage = { type: 'notification', text: 'You have a new message!', sender: 'Charlie' };
// socketStore.sendMessage(notificationMessage);

// // Disconnect from the WebSocket
// socketStore.disconnect();
// </script>
