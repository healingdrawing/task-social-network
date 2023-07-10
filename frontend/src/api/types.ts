// message type enum
export enum WSMessageType {
  ERROR_RESPONSE = "error_response", //  errorHandler // not sure this needed
  SUCCESS_RESPONSE = "success_response", //  successHandler // not sure this needed

  COMMENT_SUBMIT = "comment_submit", //  commentNewHandler
  COMMENT_RESPONSE = "comment_response", //  NEW 
  COMMENTS_LIST = "comments_list", //  commentsGetHandler

  CHAT_USERS_LIST = "chat_users_list", //  chatUsersHandler
  CHAT_MESSAGE_SUBMIT = "chat_message_submit", //  chatNewHandler

  FOLLOW_REQUEST_REJECT = "follow_request_reject", //  rejectFollowerHandler
  FOLLOW_REQUEST_ACCEPT = "follow_request_accept", //  approveFollowerHandler
  FOLLOW_REQUESTS_LIST = "follow_requests_list", //  followRequestListHandler
  FOLLOW_REQUEST_RESPONSE = "follow_request_response", //  NEW , potential use when new follow request raised

  POST_SUBMIT = "post_submit", //  postNewHandler
  POST_RESPONSE = "post_response",
  POSTS_LIST = "posts_list", //  postsGetHandler
  // ANY_PROFILE_VIEW_POSTS_LIST = "any_profile_view_posts_list", // NEW

  GROUPS_LIST = "groups_list", //  NEW to get all groups where user is member
  GROUPS_ALL_LIST = "groups_all_list", //  NEW to discover all groups
  GROUP_SUBMIT = "group_submit", //  groupNewHandler
  GROUP_POST_SUBMIT = "group_post_submit", //  groupPostNewHandler
  GROUP_POSTS_LIST = "group_posts_list", //  groupPostsGetHandler
  GROUP_POST_COMMENT_SUBMIT = "group_post_comment_submit", //  groupCommentNewHandler
  GROUP_POST_COMMENTS_LIST = "group_post_comments_list", //  groupCommentsGetHandler
  GROUP_JOIN = "group_join", //  groupJoinHandler

  GROUP_INVITES_SUBMIT = "group_invites_submit", //  groupInviteHandler
  GROUP_INVITES_LIST = "group_invites_list", //todo: attention, FOR USER
  GROUP_INVITE_ACCEPT = "group_invite_accept", //  groupInviteAcceptHandler
  GROUP_INVITE_REJECT = "group_invite_reject", //  groupInviteRejectHandler

  GROUP_REQUEST_SUBMIT = "group_request_submit", //  new
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

  USER_VISITOR_STATUS = "user_visitor_status", //  NEW
  USER_GROUP_VISITOR_STATUS = "user_group_visitor_status", //  NEW
}

export enum SuccessContent {
  FOLLOWER_WAS_ADDED = "Follower was added",
  FOLLOW_REQUEST_WAS_ADDED = "Request to become a follower was added",
}

export enum VisitorStatus {
  OWNER = "owner",
  FOLLOWER = "follower",
  REQUESTER = "requester",
  MEMBER = "member", // injection for groupVisitor() case
  VISITOR = "visitor",
}

export interface WSMessage {
  type: WSMessageType;
  data: object;
}


// used in NavBar.vue signup.ts login.ts
export interface ErrorResponse {
  message: string;
}

export interface SuccessResponse {
  message: string;
}

export interface Comment {
  content: string;
  picture: string;
  created_at: string;
  email: string;
  first_name: string;
  last_name: string;
}

export interface CommentSubmit {
  user_uuid: string;
  post_id: number; // -ugly golang cant casting number in json using interface.(int). It returns 0 and ok=false, In same time print of map shows number > 0. So data must be reconverted using float64 intermediate value. facepalm
  content: string;
  picture?: string; // jpeg,png,gif
}

export interface CommentsListRequest {
  user_uuid: string;
  post_id: number;
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

export interface PostsListRequest {
  user_uuid: string;
}

export interface Group {
  id: number;
  name: string;
  description: string;
  created_at: string;
  email: string;
  first_name: string;
  last_name: string;
}

export interface GroupSubmit {
  user_uuid: string;
  name: string;
  description: string;
  invited_emails: string;// space separated, to send invites in group creation process
}

export interface GroupsListRequest {
  user_uuid: string;
}

export interface TargetProfileRequest {
  user_uuid: string;
  target_email: string;
}

export interface ChangePrivacyRequest {
  user_uuid: string;
  make_public: string;
}

export interface UserProfile {
  email: string;
  first_name: string;
  last_name: string;
  dob: string;
  avatar: string;
  nickname: string;
  about_me: string;
  public: boolean;
}

export interface UserForList {
  email: string;
  first_name: string;
  last_name: string;
}

export interface GroupVisitorStatusRequest {
  user_uuid: string;
  group_id: number;
}

/**used for two cases visitor() and groupVisitor() */
export interface UserVisitorStatus {
  status: VisitorStatus;
}

// used for two cases Accept and Reject
export interface GroupRequestActionSubmit {
  user_uuid: string;
  group_id: number;
  requester_email: string;
}

export enum BellType {
  EVENT = "event",
  FOLLOWING = "following",
  INVITATION = "invitation",
  REQUEST = "request",
}

export type BellStatus =
  | 'visible'
  | 'hidden';

export interface Bell {
  type: BellType;
  status: BellStatus; // to remove from list, or to hide

  group_id: number; // invitation to group , and request to join group, and event
  group_name: string; // invitation to group , and request to join group , and event

  event_id: number; // event
  event_name: string; // event

  email: string; // invitation to group , and request to join group, following
  first_name: string; // invitation to group , and request to join group, following
  last_name: string; // invitation to group , and request to join group, following
  // todo: tables not ready even
}

export interface BellState {
  bells: Bell[];
}