// message type enum
export enum WSMessageType {
  ERROR = "error", //  errorHandler // not sure this needed

  COMMENT_SUBMIT = "comment_submit", //  commentNewHandler
  COMMENTS_LIST = "comments_list", //  commentsGetHandler

  CHAT_USERS_LIST = "chat_users_list", //  chatUsersHandler
  CHAT_MESSAGE_SUBMIT = "chat_message_submit", //  chatNewHandler

  FOLLOW_REQUEST_REJECT = "follow_request_reject", //  rejectFollowerHandler
  FOLLOW_REQUEST_ACCEPT = "follow_request_accept", //  approveFollowerHandler
  FOLLOW_REQUESTS_LIST = "follow_requests_list", //  followRequestListHandler

  POST_SUBMIT = "post_submit", //  postNewHandler
  POST_RESPONSE = "post_response",
  POSTS_LIST = "posts_list", //  postsGetHandler

  GROUPS_LIST = "groups_list", //  !!!groupsGetHandler //todo: NOT IMPLEMENTED ON OLD BACKEND
  GROUP_SUBMIT = "group_submit", //  groupNewHandler
  GROUP_POST_SUBMIT = "group_post_submit", //  groupPostNewHandler
  GROUP_POSTS_LIST = "group_posts_list", //  groupPostsGetHandler
  GROUP_POST_COMMENT_SUBMIT = "group_post_comment_submit", //  groupCommentNewHandler
  GROUP_POST_COMMENTS_LIST = "group_post_comments_list", //  groupCommentsGetHandler
  GROUP_JOIN = "group_join", //  groupJoinHandler
  GROUP_INVITE = "group_invite", //  groupInviteHandler
  GROUP_INVITES_LIST = "group_invited_list", //todo: attention  groupInvitedHandler
  GROUP_INVITE_ACCEPT = "group_invite_accept", //  groupInviteAcceptHandler
  GROUP_INVITE_REJECT = "group_invite_reject", //  groupInviteRejectHandler
  GROUP_REQUESTS_LIST = "group_requests_list", //  groupRequestsHandler
  GROUP_REQUEST_ACCEPT = "group_request_accept", //  groupRequestAcceptHandler
  GROUP_REQUEST_REJECT = "group_request_reject", //  groupRequestRejectHandler

  GROUP_EVENT_SUBMIT = "group_event_submit", //  eventNewHandler
  GROUP_EVENTS_LIST = "group_events_list", //  eventsGetHandler
  GROUP_EVENT_PARTICIPANTS_LIST = "group_event_participants_list", //  eventParticipantsGetHandler
  GROUP_EVENT_ATTEND = "group_event_attend", //  eventAttendHandler
  GROUP_EVENT_NOT_ATTEND = "group_event_not_attend", //  eventNotAttendHandler

  USER_CHECK = "user_check", //  sessionCheckHandler
  USER_FOLLOWING_LIST = "user_following_list", //  followingHandler
  USER_FOLLOWERS_LIST = "user_followers_list", //  followersHandler
  USER_FOLLOW = "user_follow", //  followHandler
  USER_LOGIN = "user_login", //  userLoginHandler
  USER_LOGOUT = "user_logout", //  userLogoutHandler
  USER_POSTS_LIST = "user_posts_list", //  userPostsHandler
  USER_PRIVACY = "user_privacy", //  changePrivacyHandler
  USER_PROFILE = "user_profile", //  userProfileHandler
  USER_REGISTER = "user_register", //  userRegisterHandler
  USER_UNFOLLOW = "user_unfollow", //  unfollowHandler
}

export interface WSMessage {
  type: WSMessageType;
  data: object;
}


// used in NavBar.vue signup.ts login.ts
export interface ErrorResponse {
  message: string;
}

export interface Post {
  id: number; // post id, unique, autoincrement, primary key, all posts must be stored one table in database
  title: string;
  categories: string;
  content: string;
  privacy: string;
  picture?: string;
  created_at: string;
  email: string;
  first_name: string;
  last_name: string;
}

export interface PostSubmit {
  user_uuid: string;
  title: string;
  categories: string;
  content: string;
  privacy: string;
  able_to_see: string; // user emails, to filter posts by privacy, on backend side, before sending to frontend
  picture?: string; // jpeg,png,gif
}