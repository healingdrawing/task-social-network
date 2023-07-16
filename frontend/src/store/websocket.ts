import { defineStore } from 'pinia';
import { WSMessage, WSMessageType, Post, GroupPost, Comment, UserProfile, UserForList, UserForChatList, UserVisitorStatus as Visitor, Bell, BellType, Group, Event, PrivateChatMessage, GroupChatMessage } from '@/api/types';
import router from '@/router/index';

const websockets: (WebSocket | null)[] = [];

export const useWebSocketStore = defineStore({
  id: 'websocket',
  state: () => ({
    socket: null as WebSocket | null,
    messages: [] as WSMessage[],
    private_chat_user_id: 0,
    group_chat_id: 0,
  }),
  actions: {
    /**for internal usage of send message for private chat */
    set_private_chat_user_id(id: number) {
      this.private_chat_user_id = id;
    },
    /**for internal usage of send message for group chat */
    set_group_chat_id(id: number) {
      this.group_chat_id = id;
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
        console.log(`Received message: ${event.data}`);
        const message = JSON.parse(event.data) as WSMessage;
        this.clearMessagesWhenNewMessageArrives(message)
        this.messages.unshift(message);
      };
    },
    sendMessage(message: WSMessage) {
      // this.clearMessagesBeforeSendMessage(message);
      const messageString = JSON.stringify(message);
      console.log(`\n=Sending message.\ntype: ${message.type}\ndata(json string): ${messageString}`);
      this.socket?.send(messageString);
    },
    disconnect() {
      console.log('WebSocket disconnecting');
      this.socket?.close();
      this.socket = null;
      console.log('socket', this.socket);
    },

    //todo: at the moment not used, may be implement clear chats history button later
    clearChatMessages() {
      this.messages = this.messages.filter((message) =>
        message.type !== WSMessageType.PRIVATE_CHAT_MESSAGE
        && message.type !== WSMessageType.GROUP_CHAT_MESSAGE
      );
    },

    /**clearMessagesWhenNewMessageArrives removes all the messages of type = message.Type, before unshift new message, to prevent duplication of messages in getters ( -> screen/view) */
    clearMessagesWhenNewMessageArrives(new_message: WSMessage) {
      console.log('==================\n=clearMessagesWhenNewMessageArrives===', new_message.type, '\n==================');

      this.messages = this.messages.filter((message) => message.type !== WSMessageType.SUCCESS_RESPONSE);
      this.messages = this.messages.filter((message) => message.type !== WSMessageType.ERROR_RESPONSE);


      switch (new_message.type) {
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

    /**clearMessagesBeforeSendMessage removes all the messages of type = message.Type, before fetch fresh, to prevent duplication of messages in getters ( -> screen/view) */
    clearMessagesBeforeSendMessage(message: WSMessage) {
      console.log('==================\n=clearMessagesBeforeSendMessage===', message.type, '\n==================');

      this.messages = this.messages.filter((message) => message.type !== WSMessageType.SUCCESS_RESPONSE);

      switch (message.type) {
        //todo: NOPE. FORGET ABOUT IT!
        // perhaps refactor to replace response to list for all possible cases, and then, in case of success, function can be oneline without switch

        case WSMessageType.FOLLOW_REQUEST_ACCEPT:
        case WSMessageType.FOLLOW_REQUEST_REJECT:
        case WSMessageType.FOLLOW_REQUESTS_LIST:
          console.log('case clearMessages==============================', message.type);
          // this.messages = this.messages.filter((message) => message.type !== WSMessageType.FOLLOW_REQUEST_RESPONSE);
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.FOLLOW_REQUESTS_LIST);
          break;

        case WSMessageType.GROUP_INVITE_ACCEPT:
        case WSMessageType.GROUP_INVITE_REJECT:
        case WSMessageType.GROUP_INVITES_LIST:
          console.log('case clearMessages==============================', message.type);
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.GROUP_INVITES_LIST);
          break;

        case WSMessageType.GROUP_REQUEST_ACCEPT:
        case WSMessageType.GROUP_REQUEST_REJECT:
        case WSMessageType.GROUP_REQUESTS_LIST:
          console.log('case clearMessages==============================', message.type);
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.GROUP_REQUESTS_LIST);
          break;

        case WSMessageType.GROUP_EVENT_SUBMIT:
        case WSMessageType.GROUP_EVENT_GOING:
        case WSMessageType.GROUP_EVENT_NOT_GOING:
        case WSMessageType.GROUP_EVENTS_LIST:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.GROUP_EVENTS_LIST);
          break;

        case WSMessageType.USER_GROUPS_FRESH_EVENTS_LIST: //todo: for BellView.vue
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.USER_GROUPS_FRESH_EVENTS_LIST);
          break;

        case WSMessageType.POST_SUBMIT:
        case WSMessageType.POSTS_LIST:
        case WSMessageType.USER_POSTS_LIST:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.POSTS_LIST);
          break;

        case WSMessageType.COMMENT_SUBMIT:
        case WSMessageType.COMMENTS_LIST:
        case WSMessageType.GROUP_POST_COMMENT_SUBMIT:
        case WSMessageType.GROUP_POST_COMMENTS_LIST:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.COMMENTS_LIST);
          break;

        case WSMessageType.GROUP_POST_SUBMIT:
        case WSMessageType.GROUP_POSTS_LIST:
        case WSMessageType.USER_GROUP_POSTS_LIST:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.GROUP_POSTS_LIST);
          break;

        case WSMessageType.GROUP_SUBMIT:
        case WSMessageType.GROUPS_LIST:
        case WSMessageType.GROUPS_ALL_LIST:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.GROUPS_LIST);
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

        case WSMessageType.USER_GROUP_VISITOR_STATUS:
          this.messages = this.messages.filter((message) => message.type !== WSMessageType.USER_GROUP_VISITOR_STATUS);
          break;

        default:
          console.log('SKIP clearMessagesBeforeSendMessage default============', message.type);
      }
    },

    /**remove all the messages from the store */
    facepalm() {
      console.log('===BEFORE facepalm===');
      this.$state.messages = [];
      console.log('===facepalm===');
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
      // websockets.length = 0; // it is const so "= []" raises error
      this.facepalm();
      router.push({ name: 'login' });
      console.log('========\n=======before reset and after push ======\n======');
      // this.$reset();
    },

    /** todo: later check/remove. Try to prevent artefacts between routing. Looks like it works, but facepalm generally */
    clearOnBeforeRouteLeave(path: string) {
      if (path == "/posts") {
        this.messages = this.messages.filter((message) => message.type !== WSMessageType.POSTS_LIST);
        this.messages = this.messages.filter((message) => message.type !== WSMessageType.POST_SUBMIT);
        this.messages = this.messages.filter((message) => message.type !== WSMessageType.POST_RESPONSE);
      }
      console.log('clearBeforeRouteChange ==============================');
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
    groupVisitor(): Visitor {
      // filter in exact visitor status response messages
      const visitor_messages = this.messages.filter((message) => message.type === WSMessageType.USER_GROUP_VISITOR_STATUS);
      const visitor = visitor_messages.map((message) => message.data as Visitor)[0];

      console.log('groupVisitor.status========getter========', visitor);
      return visitor;

    },
    userProfile(): UserProfile {
      const profile_messages = this.messages.filter((message) => message.type === WSMessageType.USER_PROFILE);
      const profile = profile_messages.map((message) => message.data as UserProfile)[0];
      return profile;
    },
    userFollowingList(): UserForList[] {
      const followingListMessages = this.messages.filter((message) => message.type === WSMessageType.USER_FOLLOWING_LIST && message.data !== null && message.data !== undefined)
      const followingList = followingListMessages.map((message) => message.data as UserForList).flat()
      return followingList
    },
    userFollowersList(): UserForList[] {
      const followersListMessages = this.messages.filter((message) => message.type === WSMessageType.USER_FOLLOWERS_LIST && message.data !== null && message.data !== undefined)
      const followersList = followersListMessages.map((message) => message.data as UserForList).flat()
      return followersList
    },
    /**all the posts able to see by user. Excludes group posts(separated view)*/
    postsList(): Post[] {
      const fresh_posts_messages = this.messages.filter((message) => message.type === WSMessageType.POST_RESPONSE);
      const fresh_posts = fresh_posts_messages.map((message) => message.data as Post);
      // todo: probably the function code above is artefact. Check and refactor if so
      const history_posts_messages_list = this.messages.filter((message) => message.type === WSMessageType.POSTS_LIST && message.data !== null);
      const history_posts = history_posts_messages_list.map((message) =>
        (message.data as Post[]).map((post) => post)
      ).flat();

      const posts = [...fresh_posts, ...history_posts];
      return posts;
    },
    commentsList(): Comment[] {
      const comments_messages_list = this.messages.filter((message) => message.type === WSMessageType.COMMENTS_LIST && message.data !== null);
      const comments = comments_messages_list.map((message) =>
        (message.data as Comment[]).map((comment) => comment)
      ).flat();
      return [...comments];
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
      return [...groups]
    },

    /**all the groups to discover. The type is same as above, but request type is GROUPS_ALL_LIST */
    groupsAllList(): Group[] {
      const groups_messages_list = this.messages.filter((message) => message.type === WSMessageType.GROUPS_LIST && message.data !== null);
      const groups = groups_messages_list.map((message) =>
        (message.data as Group[]).map((group) => group)
      ).flat();
      return [...groups];
    },

    groupEventsList(): Event[] {
      const group_events_messages_list = this.messages.filter((message) => message.type === WSMessageType.GROUP_EVENTS_LIST && message.data !== null);
      const group_events = group_events_messages_list.map((message) =>
        (message.data as Event[]).map((group_event) => group_event)
      ).flat();
      return [...group_events];
    },

    group_chat_messages_list(): GroupChatMessage[] {
      const group_chat_messages = this.messages.filter((message) => message.type === WSMessageType.GROUP_CHAT_MESSAGE && message.data !== null);
      const chat_messages = group_chat_messages.map((message) =>
        (message.data as GroupChatMessage)).filter((message) => message.group_id === this.group_chat_id)
      return [...chat_messages];
    },

    private_chat_messages_list(): PrivateChatMessage[] {
      const private_chat_messages = this.messages.filter((message) => message.type === WSMessageType.PRIVATE_CHAT_MESSAGE && message.data !== null);

      const chat_messages = private_chat_messages.map((message) =>
        (message.data as PrivateChatMessage))

      const messages = chat_messages.filter((message) =>
        message.user_id === this.private_chat_user_id || message.target_user_id === this.private_chat_user_id
      )

      return [...messages];
    },

    private_chat_users_list(): UserForChatList[] {
      const private_chat_users = this.messages.filter((message) => message.type === WSMessageType.PRIVATE_CHAT_USERS_LIST && message.data !== null);
      const users = private_chat_users.map((message) =>
        (message.data as UserForChatList)).flat()
      return [...users];
    },

    followRequestsBellList(): Bell[] {
      const fresh_follow_requests_messages = this.messages.filter((message) => message.type === WSMessageType.FOLLOW_REQUEST_RESPONSE)
      const fresh_follow_requests = fresh_follow_requests_messages.map((message) => message.data as Bell)

      // todo: not clear place, tried not null check, and for another lists too
      const history_follow_requests_messages_list = this.messages.filter((message) => message.type === WSMessageType.FOLLOW_REQUESTS_LIST && message.data !== null)
      const history_follow_requests = history_follow_requests_messages_list.map((message) =>
        (message.data as Bell[]).map((bell) => bell)
      ).flat()

      const follow_requests = [...fresh_follow_requests, ...history_follow_requests]

      // prepare for display, fill the empty fields
      follow_requests.forEach((bell) => {
        bell.type = BellType.FOLLOWING
      })

      console.log('pinia \n follow_requests========== ', follow_requests.length,
        '\n fresh_follow_requests_messages========== ', fresh_follow_requests_messages.length,
        '\n fresh_follow_requests========== ', fresh_follow_requests.length,
        '\n history_follow_requests_messages_list========== ', history_follow_requests_messages_list.length,
        '\n history_follow_requests========== ', history_follow_requests.length);
      return follow_requests
    },

    groupInvitesBellList(): Bell[] {
      // todo: not clear place, tried not null check, and for another lists too
      const group_invites_messages_list = this.messages.filter((message) => message.type === WSMessageType.GROUP_INVITES_LIST && message.data !== null)
      const group_invites = group_invites_messages_list.map((message) =>
        (message.data as Bell[]).map((bell) => bell)
      ).flat()

      // prepare for display, fill the empty fields
      group_invites.forEach((bell) => {
        bell.type = BellType.INVITATION
      })

      console.log('pinia \n group_invites========== ', group_invites.length,
        '\n group_invites_messages_list========== ', group_invites_messages_list.length);
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

      console.log('pinia \n group_requests========== ', group_requests.length,
        '\n group_requests_messages_list========== ', group_requests_messages_list.length);
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

      console.log('pinia \n group_requests========== ', group_fresh_events.length,
        '\n group_requests_messages_list========== ', group_fresh_events_messages_list.length);
      return group_fresh_events
    },

    bellsList(): Bell[] {
      //TODO: add other bells x4 summary
      return [...this.followRequestsBellList, ...this.groupInvitesBellList, ...this.groupRequestsBellList, ...this.groupEventsBellList]
    },

  },

});
