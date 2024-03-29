import { defineStore } from 'pinia';
import { WSMessage, WSMessageType, Post, GroupPost, Comment, UserProfile, UserForList, UserForChatList, UserVisitorStatus as Visitor, Bell, BellType, Group, Event, PrivateChatMessage, GroupChatMessage } from '@/api/types';
import router from '@/router/index';

const websockets: (WebSocket | null)[] = [];

export const useWebSocketStore = defineStore({
  id: 'websocket',
  state: () => ({
    socket: null as WebSocket | null,
    messages: [] as WSMessage[],
    private_chat_user_id: sessionStorage.getItem("private_chat_user_id") !== null ? parseInt(sessionStorage.getItem("private_chat_user_id")!) : -1,
    group_chat_id: sessionStorage.getItem("group_chat_id") !== null ? parseInt(sessionStorage.getItem("group_chat_id")!) : -1,
  }),
  actions: {
    /**for internal usage of send message for private chat */
    set_private_chat_user_id(id: number) {
      this.private_chat_user_id = id;
      sessionStorage.setItem("private_chat_user_id", id.toString());
    },
    /**for internal usage of send message for group chat */
    set_group_chat_id(id: number) {
      this.group_chat_id = id;
      sessionStorage.setItem("group_chat_id", id.toString());
    },

    send_private_chat_message(message: string, uuid: string) {
      const wsMessage: WSMessage = {
        type: WSMessageType.PRIVATE_CHAT_MESSAGE,
        data: {
          user_uuid: uuid,
          target_user_id: this.private_chat_user_id,
          content: message,
        },
      };
      this.sendMessage(wsMessage);
    },

    send_group_chat_message(message: string, uuid: string) {
      const wsMessage: WSMessage = {
        type: WSMessageType.GROUP_CHAT_MESSAGE,
        data: {
          user_uuid: uuid,
          group_id: this.group_chat_id,
          content: message,
        },
      };
      this.sendMessage(wsMessage);
    },

    waitForConnection() {
      return new Promise<void>((resolve) => {
        const interval = setInterval(() => {
          if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            console.log('=== Connection is established ===')
            clearInterval(interval);
            resolve();
          }
        }, 100);
      });
    },

    refresh_websocket() {
      console.log("=== refreshing websocket ===");
      if (this.socket) {
        return;
      }
      this.killThemAll();
      const uuid = sessionStorage.getItem('UUID');
      console.log('=== uuid inside refresh_websocket:\n' + uuid);
      if (uuid) {
        this.connect(uuid);
      } else {
        alert('uuid is null');
        router.push('/login');
      }
    },
    connect(uuid: string) {
      this.socket = new WebSocket(`ws://localhost:8080/ws?uuid=${uuid}`);

      websockets.unshift(this.socket);

      this.socket.onopen = () => {
        console.log('WebSocket connected');
      };

      this.socket.onclose = () => {
        console.log('WebSocket disconnected');
      };

      this.socket.onmessage = (event) => {
        console.log(`= Received message =`);
        // console.log(`Received message: ${event.data}`); // todo: remove debug, giant output for picture, etc
        const message = JSON.parse(event.data) as WSMessage;
        this.clearMessagesWhenNewMessageArrives(message)
        this.messages.unshift(message);
      };
    },
    sendMessage(message: WSMessage) {
      // this.clearMessagesBeforeSendMessage(message);
      const message_string = JSON.stringify(message);
      console.log(`\n=Sending message.\ntype: ${message.type}\ndata(json string): ${message_string}`);
      this.socket?.send(message_string);
    },
    disconnect() {
      console.log('WebSocket disconnecting');
      this.socket?.close();
      this.socket = null;
      console.log('socket', this.socket);
      this.killThemAll();
    },

    /**clearMessagesWhenNewMessageArrives removes all the messages of type = message.Type, before unshift new message, to prevent duplication of messages in getters ( -> screen/view) */
    clearMessagesWhenNewMessageArrives(new_message: WSMessage) {
      console.log(
        '====================================\n', '=clearMessagesWhenNewMessageArrives=\n', new_message.type, '\n',
        '====================================');

      switch (new_message.type) {
        case WSMessageType.ERROR_RESPONSE:
        case WSMessageType.INFO_RESPONSE:
        case WSMessageType.SUCCESS_RESPONSE:
        case WSMessageType.FOLLOW_REQUESTS_LIST:
        case WSMessageType.GROUP_INVITES_LIST:
        case WSMessageType.GROUP_REQUESTS_LIST:
        case WSMessageType.GROUP_EVENTS_LIST:
        case WSMessageType.USER_GROUPS_FRESH_EVENTS_LIST: //todo: for BellView.vue
        case WSMessageType.POSTS_LIST:
        case WSMessageType.COMMENTS_LIST:
        case WSMessageType.GROUP_POSTS_LIST:
        case WSMessageType.GROUPS_LIST:
        case WSMessageType.USER_PROFILE:
        case WSMessageType.USER_FOLLOWING_LIST:
        case WSMessageType.USER_FOLLOWERS_LIST:
        case WSMessageType.USER_VISITOR_STATUS:
        case WSMessageType.USER_GROUP_VISITOR_STATUS:
        case WSMessageType.PRIVATE_CHAT_USERS_LIST:
          this.messages = this.messages.filter((message) => message.type !== new_message.type);
          break;

        default:
          console.log('SKIP clearMessagesWhenNewMessageArrives default============', new_message.type);
      }
    },

    /**remove all the messages from the store */
    facepalm() {
      console.log('= BEFORE facepalm =');
      this.$state.messages = [];
      console.log('= facepalm =');
    },
    /**
     * close all the websockets and empty the global array const websockets: (WebSocket | null)[] = [];
     * Attempt to prevent creation of artefacts for emergency case, hardreload etc*/
    killThemAll() {
      console.log('===killThemAll===');
      console.log('websockets', websockets.length);
      websockets.forEach((websocket) => {
        console.log('============================')
        console.log('forEach websocket disconnecting');
        console.log('forEach websocket', websocket);
        websocket?.close(1000, 'killThemAll');
        websocket = null;
        console.log('forEach websocket', websocket);
      });
      websockets.length = 0; // it is const so "= []" raises error.
      this.facepalm();
    },

  },
  getters: {
    isConnected(): boolean {
      return this.socket !== null && this.socket.readyState === WebSocket.OPEN;
    },
    visitor(): Visitor {
      // filter in exact user visitor status response messages
      const visitor_messages = this.messages.filter((message) => message.type === WSMessageType.USER_VISITOR_STATUS);
      const visitor = visitor_messages.map((message) => message.data as Visitor)[0];

      return visitor;
    },
    groupVisitor(): Visitor {
      // filter in exact group visitor status response messages
      const visitor_messages = this.messages.filter((message) => message.type === WSMessageType.USER_GROUP_VISITOR_STATUS);
      const visitor = visitor_messages.map((message) => message.data as Visitor)[0];

      return visitor;
    },
    userProfile(): UserProfile {
      const profile_messages = this.messages.filter((message) => message.type === WSMessageType.USER_PROFILE);
      const profile = profile_messages.map((message) => message.data as UserProfile)[0];
      return profile;
    },
    userFollowingList(): UserForList[] {
      const following_list_messages = this.messages.filter((message) => message.type === WSMessageType.USER_FOLLOWING_LIST && message.data !== null && message.data !== undefined)
      const following_list = following_list_messages.map((message) => message.data as UserForList).flat()

      return following_list
    },
    userFollowersList(): UserForList[] {
      const followers_list_messages = this.messages.filter((message) => message.type === WSMessageType.USER_FOLLOWERS_LIST && message.data !== null && message.data !== undefined)
      const followers_list = followers_list_messages.map((message) => message.data as UserForList).flat()

      return followers_list
    },
    /**all the posts able to see by user. Excludes group posts(separated view)*/
    postsList(): Post[] {
      const posts_messages_list = this.messages.filter((message) => message.type === WSMessageType.POSTS_LIST && message.data !== null);
      const posts = posts_messages_list.map((message) =>
        (message.data as Post[]).map((post) => post)
      ).flat();

      return posts;
    },

    /** commentsList used for both posts and group_posts,
     *  because no difference in structure,
     * but requests have different type,
     * so data received from different db tables, using different queries */
    commentsList(): Comment[] {
      const comments_messages_list = this.messages.filter((message) => message.type === WSMessageType.COMMENTS_LIST && message.data !== null);
      const comments = comments_messages_list.map((message) =>
        (message.data as Comment[]).map((comment) => comment)
      ).flat();

      return comments;
    },

    groupPostsList(): GroupPost[] {
      const group_posts_messages_list = this.messages.filter((message) => message.type === WSMessageType.GROUP_POSTS_LIST && message.data !== null);
      const group_posts = group_posts_messages_list.map((message) =>
        (message.data as GroupPost[]).map((post) => post)
      ).flat();

      return group_posts;
    },

    /**all the groups where user is member*/
    groupsList(): Group[] {
      const groups_messages_list = this.messages.filter((message) => message.type === WSMessageType.GROUPS_LIST && message.data !== null);
      const groups = groups_messages_list.map((message) =>
        (message.data as Group[]).map((group) => group)
      ).flat();

      return groups
    },

    /**all the groups to discover. The type is same as above, but request type is GROUPS_ALL_LIST. So it is just for better visual */
    groupsAllList(): Group[] {
      const groups_messages_list = this.messages.filter((message) => message.type === WSMessageType.GROUPS_LIST && message.data !== null);
      const groups = groups_messages_list.map((message) =>
        (message.data as Group[]).map((group) => group)
      ).flat();

      return groups;
    },

    groupEventsList(): Event[] {
      const group_events_messages_list = this.messages.filter((message) => message.type === WSMessageType.GROUP_EVENTS_LIST && message.data !== null);
      const group_events = group_events_messages_list.map((message) =>
        (message.data as Event[]).map((group_event) => group_event)
      ).flat();

      return group_events;
    },

    group_chat_messages_list(): GroupChatMessage[] {
      const group_chat_messages = this.messages.filter((message) => message.type === WSMessageType.GROUP_CHAT_MESSAGE && message.data !== null);
      const chat_messages = group_chat_messages.map((message) =>
        (message.data as GroupChatMessage)).filter((message) => message.group_id === this.group_chat_id)

      return chat_messages;
    },

    private_chat_messages_list(): PrivateChatMessage[] {
      const private_chat_messages = this.messages.filter((message) => message.type === WSMessageType.PRIVATE_CHAT_MESSAGE && message.data !== null);

      const chat_messages = private_chat_messages.map((message) =>
        (message.data as PrivateChatMessage))

      const messages = chat_messages.filter((message) =>
        message.user_id === this.private_chat_user_id || message.target_user_id === this.private_chat_user_id
      )

      return messages;
    },

    private_chat_users_list(): UserForChatList[] {
      const private_chat_users = this.messages.filter((message) => message.type === WSMessageType.PRIVATE_CHAT_USERS_LIST && message.data !== null);
      const users = private_chat_users.map((message) =>
        (message.data as UserForChatList)).flat()

      return users;
    },

    followRequestsBellList(): Bell[] {
      const follow_requests_messages_list = this.messages.filter((message) => message.type === WSMessageType.FOLLOW_REQUESTS_LIST && message.data !== null)
      const follow_requests = follow_requests_messages_list.map((message) =>
        (message.data as Bell[]).map((bell) => bell)
      ).flat()

      // prepare for display, fill the empty fields
      follow_requests.forEach((bell) => {
        bell.type = BellType.FOLLOWING
      })

      return follow_requests
    },

    groupInvitesBellList(): Bell[] {
      const group_invites_messages_list = this.messages.filter((message) => message.type === WSMessageType.GROUP_INVITES_LIST && message.data !== null)
      const group_invites = group_invites_messages_list.map((message) =>
        (message.data as Bell[]).map((bell) => bell)
      ).flat()

      // prepare for display, fill the empty fields
      group_invites.forEach((bell) => {
        bell.type = BellType.INVITATION
      })

      return group_invites
    },

    groupRequestsBellList(): Bell[] {
      const group_requests_messages_list = this.messages.filter((message) => message.type === WSMessageType.GROUP_REQUESTS_LIST && message.data !== null)
      const group_requests = group_requests_messages_list.map((message) =>
        (message.data as Bell[]).map((bell) => bell)
      ).flat()

      // prepare for display, fill the empty fields
      group_requests.forEach((bell) => {
        bell.type = BellType.REQUEST
      })

      return group_requests
    },

    groupEventsBellList(): Bell[] {
      const group_fresh_events_messages_list = this.messages.filter((message) => message.type === WSMessageType.USER_GROUPS_FRESH_EVENTS_LIST && message.data !== null)
      const group_fresh_events = group_fresh_events_messages_list.map((message) =>
        (message.data as Bell[]).map((bell) => bell)
      ).flat()

      // prepare for display, fill the empty fields
      group_fresh_events.forEach((bell) => {
        bell.type = BellType.EVENT
      })

      return group_fresh_events
    },

    bellsList(): Bell[] {
      return [...this.followRequestsBellList, ...this.groupInvitesBellList, ...this.groupRequestsBellList, ...this.groupEventsBellList]
    },

  },

});
