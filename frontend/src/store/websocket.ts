import { defineStore } from 'pinia';
import { WSMessage, WSMessageType, Post, UserProfile, UserForList, UserVisitorStatus as Visitor, SuccessResponse, SuccessContent, VisitorStatus, Bell } from '@/api/types';


export const useWebSocketStore = defineStore({
  id: 'websocket',
  state: () => ({
    socket: null as WebSocket | null,
    messages: [] as WSMessage[],
  }),
  actions: {
    connect(uuid: string) {
      this.socket = new WebSocket(`ws://localhost:8080/ws?uuid=${uuid}`);

      this.socket.onopen = () => {
        console.log('WebSocket connected');
      };

      this.socket.onclose = () => {
        console.log('WebSocket disconnected');
      };

      this.socket.onmessage = (event) => {
        console.log(`Received message: ${event.data}`);
        const message = JSON.parse(event.data) as WSMessage;
        this.messages.unshift(message);
      };
    },
    sendMessage(message: WSMessage) {
      this.clearMessages(message);
      const messageString = JSON.stringify(message);
      console.log(`Sending message json string: ${messageString}`);
      this.socket?.send(messageString);
    },
    disconnect() {
      this.socket?.close();
      this.socket = null;
    },
    /**clearMessages removes all the messages of type = message.Type, before fetch fresh, to prevent duplication of messages in getters ( -> screen/view) */
    clearMessages(message: WSMessage) {
      console.log('clearMessages==============================', message.type);

      this.messages = this.messages.filter((message) => message.type !== WSMessageType.SUCCESS_RESPONSE);

      switch (message.type) {
        //todo: add x4 cases for each type of bell
        case WSMessageType.FOLLOW_REQUESTS_LIST:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.FOLLOW_REQUEST_RESPONSE);
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.FOLLOW_REQUESTS_LIST);
          break;

        case WSMessageType.USER_POSTS_LIST:
        case WSMessageType.POSTS_LIST:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.POST_RESPONSE);
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.POSTS_LIST);
          break;
        case WSMessageType.USER_PROFILE:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.USER_PROFILE);
          break;
        case WSMessageType.USER_FOLLOWING_LIST:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.USER_FOLLOWING_LIST);
          break;
        case WSMessageType.USER_FOLLOWERS_LIST:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.USER_FOLLOWERS_LIST);
          break;
        case WSMessageType.USER_FOLLOW:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.USER_FOLLOW);
          break;
        case WSMessageType.USER_UNFOLLOW:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.USER_UNFOLLOW);
          break;
        case WSMessageType.USER_VISITOR_STATUS:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.USER_VISITOR_STATUS);
          break;

        default:
          console.log('SKIP clearMessages default============', message.type);
      }
    },
  },
  getters: {
    isConnected(): boolean {
      return this.socket !== null && this.socket.readyState === WebSocket.OPEN;
    },
    visitor(): Visitor {
      // filter in exact visitor status response messages
      const visitor_messages = this.messages.filter((message) => message.type === WSMessageType.USER_VISITOR_STATUS);
      const visitor = visitor_messages.map((message) => message.data as Visitor)[0];

      console.log('visitor.status========getter========', visitor);
      return visitor;

    },
    userProfile(): UserProfile {
      const profile_messages = this.messages.filter((message) => message.type === WSMessageType.USER_PROFILE);
      const profile = profile_messages.map((message) => message.data as UserProfile)[0];
      return profile;
    },
    userFollowingList(): UserForList[] {
      const followingListMessages = this.messages.filter((message) => message.type === WSMessageType.USER_FOLLOWING_LIST && message.data !== null && message.data !== undefined)
      const followingList = followingListMessages.map((message) => message.data as UserForList)
      return followingList
    },
    userFollowersList(): UserForList[] {
      const followersListMessages = this.messages.filter((message) => message.type === WSMessageType.USER_FOLLOWERS_LIST && message.data !== null && message.data !== undefined)
      const followersList = followersListMessages.map((message) => message.data as UserForList)
      return followersList
    },
    /**all the posts able to see by user. Excludes group posts(separated view)*/
    postsList(): Post[] {
      const fresh_posts_messages = this.messages.filter((message) => message.type === WSMessageType.POST_RESPONSE);
      const fresh_posts = fresh_posts_messages.map((message) => message.data as Post);

      const history_posts_messages_list = this.messages.filter((message) => message.type === WSMessageType.POSTS_LIST);
      const history_posts = history_posts_messages_list.map((message) =>
        (message.data as Post[]).map((post) => post)
      ).flat();

      const posts = [...fresh_posts, ...history_posts];
      return posts;
    },
    followRequestsList(): Bell[] {
      const fresh_follow_requests_messages = this.messages.filter((message) => message.type === WSMessageType.FOLLOW_REQUEST_RESPONSE)
      const fresh_follow_requests = fresh_follow_requests_messages.map((message) => message.data as Bell)

      // todo: not clear place, tried not null check, and for another lists too
      const history_follow_requests_messages_list = this.messages.filter((message) => message.type === WSMessageType.FOLLOW_REQUESTS_LIST && message.data !== null)
      const history_follow_requests = history_follow_requests_messages_list.map((message) =>
        (message.data as Bell[]).map((bell) => bell)
      ).flat()

      const follow_requests = [...fresh_follow_requests, ...history_follow_requests]
      return follow_requests
    },
    bellsList(): Bell[] {
      return [...this.followRequestsList] //TODO: add other bells x4 summary
    },

    commentsList(): WSMessage[] { return this.messages.filter((message) => message.type === WSMessageType.COMMENTS_LIST) },
    chatUsersList(): WSMessage[] { return this.messages.filter((message) => message.type === WSMessageType.CHAT_USERS_LIST) },
    groupsList(): WSMessage[] { return this.messages.filter((message) => message.type === WSMessageType.GROUPS_LIST) },
    groupPostsList(): WSMessage[] { return this.messages.filter((message) => message.type === WSMessageType.GROUP_POSTS_LIST) },
    groupPostCommentsList(): WSMessage[] { return this.messages.filter((message) => message.type === WSMessageType.GROUP_POST_COMMENTS_LIST) },
    groupRequestsList(): WSMessage[] { return this.messages.filter((message) => message.type === WSMessageType.GROUP_REQUESTS_LIST) },
    groupEventsList(): WSMessage[] { return this.messages.filter((message) => message.type === WSMessageType.GROUP_EVENTS_LIST) },
    groupEventParticipantsList(): WSMessage[] { return this.messages.filter((message) => message.type === WSMessageType.GROUP_EVENT_PARTICIPANTS_LIST) },


    // userPostsList(): WSMessage[] { return this.messages.filter((message) => message.type === WSMessageType.USER_POSTS_LIST) },

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
