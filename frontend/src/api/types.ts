// message type enum
export enum WSMessageType {
  ERROR_RESPONSE = "error_response",
  SUCCESS_RESPONSE = "success_response",

  COMMENT_SUBMIT = "comment_submit",
  COMMENT_RESPONSE = "comment_response",
  COMMENTS_LIST = "comments_list",

  PRIVATE_CHAT_USERS_LIST = "private_chat_users_list",
  PRIVATE_CHAT_MESSAGE = "private_chat_message",
  GROUP_CHAT_MESSAGE = "group_chat_message",

  FOLLOW_REQUEST_REJECT = "follow_request_reject",
  FOLLOW_REQUEST_ACCEPT = "follow_request_accept",
  FOLLOW_REQUESTS_LIST = "follow_requests_list",
  FOLLOW_REQUEST_RESPONSE = "follow_request_response",

  POST_SUBMIT = "post_submit",
  POST_RESPONSE = "post_response",
  POSTS_LIST = "posts_list",

  GROUPS_LIST = "groups_list",
  GROUPS_ALL_LIST = "groups_all_list",
  GROUP_SUBMIT = "group_submit",
  GROUP_POST_SUBMIT = "group_post_submit",
  GROUP_POSTS_LIST = "group_posts_list",
  GROUP_POST_COMMENT_SUBMIT = "group_post_comment_submit",
  GROUP_POST_COMMENTS_LIST = "group_post_comments_list",
  GROUP_JOIN = "group_join",

  GROUP_INVITES_SUBMIT = "group_invites_submit",
  GROUP_INVITES_LIST = "group_invites_list", // FOR SINGLE USER
  GROUP_INVITE_ACCEPT = "group_invite_accept",
  GROUP_INVITE_REJECT = "group_invite_reject",

  GROUP_REQUEST_SUBMIT = "group_request_submit",
  GROUP_REQUESTS_LIST = "group_requests_list",
  GROUP_REQUEST_ACCEPT = "group_request_accept",
  GROUP_REQUEST_REJECT = "group_request_reject",

  GROUP_EVENT_SUBMIT = "group_event_submit",
  GROUP_EVENTS_LIST = "group_events_list",
  GROUP_EVENT_PARTICIPANTS_LIST = "group_event_participants_list",
  GROUP_EVENT_GOING = "group_event_going",
  GROUP_EVENT_NOT_GOING = "group_event_not_going",

  USER_CHECK = "user_check",
  USER_FOLLOWING_LIST = "user_following_list",
  USER_FOLLOWERS_LIST = "user_followers_list",
  USER_FOLLOW = "user_follow",
  USER_LOGIN = "user_login",
  USER_LOGOUT = "user_logout",
  USER_POSTS_LIST = "user_posts_list",
  USER_GROUP_POSTS_LIST = "user_group_posts_list",
  USER_GROUPS_FRESH_EVENTS_LIST = "user_groups_fresh_events_list",
  USER_PRIVACY = "user_privacy",
  USER_PROFILE = "user_profile",
  USER_REGISTER = "user_register",
  USER_UNFOLLOW = "user_unfollow",

  USER_VISITOR_STATUS = "user_visitor_status",
  USER_GROUP_VISITOR_STATUS = "user_group_visitor_status",
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

export interface PrivateChatMessage {
  content: string;
  email: string;
  first_name: string;
  last_name: string;
  created_at: string; // use inside v-for like id, because no id in chat message
  target_user_id: number;
}

export interface GroupChatMessage {
  content: string;
  email: string;
  first_name: string;
  last_name: string;
  created_at: string; // use inside v-for like id, because no id in chat message
  group_id: number; // which is group_chat_id in exact
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


export interface GroupPost {
  group_id: number;
  group_name: string;
  group_description: string;
  id: number; // group post id, unique, autoincrement, primary key, all group posts must be stored one table in database
  title: string;
  categories: string;
  content: string;
  picture?: string;
  created_at: string;
  email: string;
  first_name: string;
  last_name: string;
}

export interface GroupPostSubmit {
  user_uuid: string;
  group_id: number;
  title: string;
  categories: string;
  content: string;
  picture?: string; // jpeg,png,gif
}

export interface GroupPostsListRequest {
  user_uuid: string;
  group_id: number;
}


export interface GroupEventSubmit {
  user_uuid: string;
  group_id: number;
  title: string;
  description: string;
  date: string;
  decision: string;
}

export interface GroupEventsListRequest {
  user_uuid: string;
  group_id: number;
}

export interface Event {
  id: number;
  title: string;
  date: string;
  description: string;
  decision: string;
}

export interface GroupEventAction {
  user_uuid: string;
  event_id: number;
  decision: string;
  group_id: number;
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

// To send group invitations by group member to followers
export interface GroupInvitesSubmit {
  user_uuid: string;
  group_id: number;
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
  group_description: string; // to fill the groupStore before routing to group page

  event_id: number; // event
  event_title: string; // event
  event_description: string; // event
  event_date: string; // event //todo: not sure it needed, at the moment not used

  email: string; // invitation to group , and request to join group, following
  first_name: string; // invitation to group , and request to join group, following
  last_name: string; // invitation to group , and request to join group, following
}

export interface BellState {
  bells: Bell[];
}

export interface BellRequest {
  user_uuid: string;
}