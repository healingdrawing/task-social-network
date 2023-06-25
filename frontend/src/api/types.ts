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

export interface Message {
  messageType: MessageType;
  content: string;
}



export interface User {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  dob: string;
  avatar: Blob | null | string;
  nickname: string;
  aboutMe: string;
  public: boolean;
}

export interface Post {
  id: number; // post id, unique, autoincrement, primary key, all posts must be stored one table in database
  authorId: number; //todo: need to implement clickable link to user profile
  authorFullName: string; //todo: need to implement clickable link to user profile
  title: string;
  tags: string;
  content: string;
  privacy: string;
  followers?: number[]; // user ids, to filter posts by privacy, on backend side, before sending to frontend
  picture?: Blob | null; //todo: need to implement image or gif to post required in task. Perhaps, to prevent posting "anacondas" and "caves" photos, the images can be limited from allowed lists of images, but generally it sounds like they expect any image upload, which is unsafe, like in picture too
}