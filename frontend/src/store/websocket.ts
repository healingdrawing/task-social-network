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

interface Message<T> {
  type: MessageType;
  content: T;
}

interface CommentSubmit {
  type: MessageType.COMMENT_SUBMIT
  content: MessageUnionType
}
interface CommentsList {
  type: MessageType.COMMENTS_LIST
  content: MessageUnionType
}

interface ChatMessagesList {
  type: MessageType.CHAT_MESSAGES_LIST
  content: MessageUnionType
}
interface ChatUsersList {
  type: MessageType.CHAT_USERS_LIST
  content: MessageUnionType
}
interface ChatMessageSubmit {
  type: MessageType.CHAT_MESSAGE_SUBMIT
  content: MessageUnionType
}
interface ChatTyping {
  type: MessageType.CHAT_TYPING
  content: MessageUnionType
}

interface FollowRequestReject {
  type: MessageType.FOLLOW_REQUEST_REJECT
  content: MessageUnionType
}
interface FollowRequestAccept {
  type: MessageType.FOLLOW_REQUEST_ACCEPT
  content: MessageUnionType
}
interface FollowRequestsList {
  type: MessageType.FOLLOW_REQUESTS_LIST
  content: MessageUnionType
}

interface PostsList {
  type: MessageType.POSTS_LIST
  content: MessageUnionType
}
interface PostSubmit {
  type: MessageType.POST_SUBMIT
  content: MessageUnionType
}

interface GroupsList {
  type: MessageType.GROUPS_LIST
  content: MessageUnionType
}
interface GroupSubmit {
  type: MessageType.GROUP_SUBMIT
  content: MessageUnionType
}
interface GroupPostSubmit {
  type: MessageType.GROUP_POST_SUBMIT
  content: MessageUnionType
}
interface GroupPostsList {
  type: MessageType.GROUP_POSTS_LIST
  content: MessageUnionType
}
interface GroupPostCommentSubmit {
  type: MessageType.GROUP_POST_COMMENT_SUBMIT
  content: MessageUnionType
}
interface GroupPostCommentsList {
  type: MessageType.GROUP_POST_COMMENTS_LIST
  content: MessageUnionType
}
interface GroupJoin {
  type: MessageType.GROUP_JOIN
  content: MessageUnionType
}
interface GroupLeave {
  type: MessageType.GROUP_LEAVE
  content: MessageUnionType
}
interface GroupInvite {
  type: MessageType.GROUP_INVITE
  content: MessageUnionType
}
interface GroupInvited {
  type: MessageType.GROUP_INVITED
  content: MessageUnionType
}
interface GroupInviteAccept {
  type: MessageType.GROUP_INVITE_ACCEPT
  content: MessageUnionType
}
interface GroupInviteReject {
  type: MessageType.GROUP_INVITE_REJECT
  content: MessageUnionType
}
interface GroupRequestsList {
  type: MessageType.GROUP_REQUESTS_LIST
  content: MessageUnionType
}
interface GroupRequestAccept {
  type: MessageType.GROUP_REQUEST_ACCEPT
  content: MessageUnionType
}
interface GroupRequestReject {
  type: MessageType.GROUP_REQUEST_REJECT
  content: MessageUnionType
}

interface GroupEventSubmit {
  type: MessageType.GROUP_EVENT_SUBMIT
  content: MessageUnionType
}
interface GroupEventsList {
  type: MessageType.GROUP_EVENTS_LIST
  content: MessageUnionType
}
interface GroupEventParticipantsList {
  type: MessageType.GROUP_EVENT_PARTICIPANTS_LIST
  content: MessageUnionType
}
interface GroupEventAttend {
  type: MessageType.GROUP_EVENT_ATTEND
  content: MessageUnionType
}
interface GroupEventNotAttend {
  type: MessageType.GROUP_EVENT_NOT_ATTEND
  content: MessageUnionType
}

interface UserCheck {
  type: MessageType.USER_CHECK
  content: MessageUnionType
}
interface UserFollowingList {
  type: MessageType.USER_FOLLOWING_LIST
  content: MessageUnionType
}
interface UserFollowersList {
  type: MessageType.USER_FOLLOWERS_LIST
  content: MessageUnionType
}
interface UserFollow {
  type: MessageType.USER_FOLLOW
  content: MessageUnionType
}
interface UserLogin {
  type: MessageType.USER_LOGIN
  content: MessageUnionType
}
interface UserLogout {
  type: MessageType.USER_LOGOUT
  content: MessageUnionType
}
interface UserPostsList {
  type: MessageType.USER_POSTS_LIST
  content: MessageUnionType
}
interface UserPrivacy {
  type: MessageType.USER_PRIVACY
  content: MessageUnionType
}
interface UserProfile {
  type: MessageType.USER_PROFILE
  content: MessageUnionType
}
interface UserRegister {
  type: MessageType.USER_REGISTER
  content: MessageUnionType
}
interface UserUnfollow {
  type: MessageType.USER_UNFOLLOW
  content: MessageUnionType
}


type CommentSubmitType = Message<CommentSubmit>
type CommentsListType = Message<CommentsList>

type ChatMessagesListType = Message<ChatMessagesList>
type ChatUsersListType = Message<ChatUsersList>
type ChatMessageSubmitType = Message<ChatMessageSubmit>
type ChatTypingType = Message<ChatTyping>

type FollowRequestRejectType = Message<FollowRequestReject>
type FollowRequestAcceptType = Message<FollowRequestAccept>
type FollowRequestsListType = Message<FollowRequestsList>

type PostsListType = Message<PostsList>
type PostSubmitType = Message<PostSubmit>

type GroupsListType = Message<GroupsList>
type GroupSubmitType = Message<GroupSubmit>
type GroupPostSubmitType = Message<GroupPostSubmit>
type GroupPostsListType = Message<GroupPostsList>
type GroupPostCommentSubmitType = Message<GroupPostCommentSubmit>
type GroupPostCommentsListType = Message<GroupPostCommentsList>
type GroupJoinType = Message<GroupJoin>
type GroupLeaveType = Message<GroupLeave>
type GroupInviteType = Message<GroupInvite>
type GroupInvitedType = Message<GroupInvited>
type GroupInviteAcceptType = Message<GroupInviteAccept>
type GroupInviteRejectType = Message<GroupInviteReject>
type GroupRequestsListType = Message<GroupRequestsList>
type GroupRequestAcceptType = Message<GroupRequestAccept>
type GroupRequestRejectType = Message<GroupRequestReject>

type GroupEventSubmitType = Message<GroupEventSubmit>
type GroupEventsListType = Message<GroupEventsList>
type GroupEventParticipantsGetType = Message<GroupEventParticipantsList>
type GroupEventAttendType = Message<GroupEventAttend>
type GroupEventNotAttendType = Message<GroupEventNotAttend>

type UserCheckType = Message<UserCheck>
type UserFollowingListType = Message<UserFollowingList>
type UserFollowersListType = Message<UserFollowersList>
type UserFollowType = Message<UserFollow>
type UserLoginType = Message<UserLogin>
type UserLogoutType = Message<UserLogout>
type UserPostsListType = Message<UserPostsList>
type UserPrivacyType = Message<UserPrivacy>
type UserProfileType = Message<UserProfile>
type UserRegisterType = Message<UserRegister>
type UserUnfollowType = Message<UserUnfollow>

type MessageUnionType =
  CommentSubmit |
  CommentsList |

  ChatMessagesList |
  ChatUsersList |
  ChatMessageSubmit |
  ChatTyping |

  FollowRequestReject |
  FollowRequestAccept |
  FollowRequestsList |

  PostsList |
  PostSubmit |

  GroupsList |
  GroupSubmit |
  GroupPostSubmit |
  GroupPostsList |
  GroupPostCommentSubmit |
  GroupPostCommentsList |
  GroupJoin |
  GroupLeave |
  GroupInvite |
  GroupInvited |
  GroupInviteAccept |
  GroupInviteReject |
  GroupRequestsList |
  GroupRequestAccept |
  GroupRequestReject |

  GroupEventSubmit |
  GroupEventsList |
  GroupEventParticipantsList |
  GroupEventAttend |
  GroupEventNotAttend |

  UserCheck |
  UserFollowingList |
  UserFollowersList |
  UserFollow |
  UserLogin |
  UserLogout |
  UserPostsList |
  UserPrivacy |
  UserProfile |
  UserRegister |
  UserUnfollow

export const useWebSocketStore = defineStore({
  id: 'websocket',
  state: () => ({
    socket: null as WebSocket | null,
    messages: [] as MessageUnionType[],
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
        const message = JSON.parse(event.data) as MessageUnionType;
        this.messages.push(message);
      };
    },
    sendMessage(message: MessageUnionType) {
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

    commentsList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.COMMENTS_LIST) },
    chatMessagesList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.CHAT_MESSAGES_LIST) },
    chatUsersList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.CHAT_USERS_LIST) },
    followRequestsList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.FOLLOW_REQUESTS_LIST) },
    postsList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.POSTS_LIST) },
    groupsList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.GROUPS_LIST) },
    groupPostsList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.GROUP_POSTS_LIST) },
    groupPostCommentsList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.GROUP_POST_COMMENTS_LIST) },
    groupRequestsList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.GROUP_REQUESTS_LIST) },
    groupEventsList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.GROUP_EVENTS_LIST) },
    groupEventParticipantsList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.GROUP_EVENT_PARTICIPANTS_LIST) },
    userFollowingList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.USER_FOLLOWING_LIST) },
    userFollowersList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.USER_FOLLOWERS_LIST) },
    userPostsList(): Message<MessageUnionType>[] { return this.messages.filter((message) => message.type === MessageType.USER_POSTS_LIST) },

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
