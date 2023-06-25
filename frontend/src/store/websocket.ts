import { defineStore } from 'pinia';

// message type enum
export enum MessageType {
  COMMENT_SUBMIT = "comment_submit", //  commentNewHandler)
  COMMENTS_LIST = "comments_list", //  commentsGetHandler)

  CHAT_MESSAGES_LIST = "chat_messages_list", //  chatMessagesHandler)
  CHAT_USERS_LIST = "chat_users_list", //  chatUsersHandler)
  CHAT_MESSAGE_SUBMIT = "chat_message_submit", //  chatNewHandler)
  CHAT_TYPING = "chat_typing", //  chatTypingHandler)

  FOLLOW_REQUEST_REJECT = "follow_request_reject", //  rejectFollowerHandler)
  FOLLOW_REQUEST_ACCEPT = "follow_request_accept", //  approveFollowerHandler)
  FOLLOW_REQUESTS_LIST = "follow_requests_list", //  followRequestListHandler)

  POSTS_LIST = "posts_list", //  postsGetHandler)
  POST_SUBMIT = "post_submit", //  postNewHandler)

  GROUPS_LIST = "groups_list", //  !!!groupsGetHandler) //todo: NOT IMPLEMENTED ON OLD BACKEND
  GROUP_SUBMIT = "group_submit", //  groupNewHandler)
  GROUP_POST_SUBMIT = "group_post_submit", //  groupPostNewHandler)
  GROUP_POSTS_LIST = "group_posts_list", //  groupPostsGetHandler)
  GROUP_POST_COMMENT_SUBMIT = "group_post_comment_submit", //  groupCommentNewHandler)
  GROUP_POST_COMMENTS_LIST = "group_post_comments_list", //  groupCommentsGetHandler)
  GROUP_JOIN = "group_join", //  groupJoinHandler)
  GROUP_LEAVE = "group_leave", //  groupLeaveHandler), // TODO: not part of audit, so untested
  GROUP_INVITE = "group_invite", //  groupInviteHandler)
  GROUP_INVITED = "group_invited", //  groupInvitedHandler)
  GROUP_INVITE_ACCEPT = "group_invite_accept", //  groupInviteAcceptHandler)
  GROUP_INVITE_REJECT = "group_invite_reject", //  groupInviteRejectHandler)
  GROUP_REQUESTS_LIST = "group_requests_list", //  groupRequestsHandler)
  GROUP_REQUEST_ACCEPT = "group_request_accept", //  groupRequestAcceptHandler)
  GROUP_REQUEST_REJECT = "group_request_reject", //  groupRequestRejectHandler)

  GROUP_EVENT_SUBMIT = "group_event_submit", //  eventNewHandler)
  GROUP_EVENTS_LIST = "group_events_list", //  eventsGetHandler)
  GROUP_EVENT_PARTICIPANTS_LIST = "group_event_participants_list", //  eventParticipantsGetHandler)
  GROUP_EVENT_ATTEND = "group_event_attend", //  eventAttendHandler)
  GROUP_EVENT_NOT_ATTEND = "group_event_not_attend", //  eventNotAttendHandler)

  USER_CHECK = "user_check", //  sessionCheckHandler)
  USER_FOLLOWING_LIST = "user_following_list", //  followingHandler)
  USER_FOLLOWERS_LIST = "user_followers_list", //  followersHandler)
  USER_FOLLOW = "user_follow", //  followHandler)
  USER_LOGIN = "user_login", //  userLoginHandler)
  USER_LOGOUT = "user_logout", //  userLogoutHandler)
  USER_POSTS_LIST = "user_posts_list", //  userPostsHandler)
  USER_PRIVACY = "user_privacy", //  changePrivacyHandler)
  USER_PROFILE = "user_profile", //  userProfileHandler)
  USER_REGISTER = "user_register", //  userRegisterHandler)
  USER_UNFOLLOW = "user_unfollow", //  unfollowHandler)
}

interface Message {
  messageType: MessageType;
  content: string;
}



export const useWebSocketStore = defineStore({
  id: 'websocket',
  state: () => ({
    socket: null as WebSocket | null,
    messages: [] as Message[],
  }),
  actions: {
    connect() {
      this.socket = new WebSocket('ws://localhost:8080');

      this.socket.onopen = () => {
        console.log('WebSocket connected');
      };

      this.socket.onclose = () => {
        console.log('WebSocket disconnected');
      };

      this.socket.onmessage = (event) => {
        const message = JSON.parse(event.data) as Message;
        this.messages.push(message);
      };
    },
    sendMessage(message: Message) {
      const messageString = JSON.stringify(message);
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

    commentsList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.COMMENTS_LIST) },
    chatMessagesList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.CHAT_MESSAGES_LIST) },
    chatUsersList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.CHAT_USERS_LIST) },
    followRequestsList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.FOLLOW_REQUESTS_LIST) },
    postsList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.POSTS_LIST) },
    groupsList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.GROUPS_LIST) },
    groupPostsList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.GROUP_POSTS_LIST) },
    groupPostCommentsList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.GROUP_POST_COMMENTS_LIST) },
    groupRequestsList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.GROUP_REQUESTS_LIST) },
    groupEventsList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.GROUP_EVENTS_LIST) },
    groupEventParticipantsList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.GROUP_EVENT_PARTICIPANTS_LIST) },
    userFollowingList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.USER_FOLLOWING_LIST) },
    userFollowersList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.USER_FOLLOWERS_LIST) },
    userPostsList(): Message[] { return this.messages.filter((message) => message.messageType === MessageType.USER_POSTS_LIST) },

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
