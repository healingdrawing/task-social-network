package main

import "net/http"

func registerHandlers() {
	// Websocket
	http.HandleFunc("/ws", wsConnection)

	// API
	http.HandleFunc("/api/comment/submit", commentNewHandler)
	http.HandleFunc("/api/comments/get", commentsGetHandler)

	http.HandleFunc("/api/chat/getmessages", chatMessagesHandler)
	http.HandleFunc("/api/chat/getusers", chatUsersHandler)
	http.HandleFunc("/api/chat/newmessage", chatNewHandler)
	http.HandleFunc("/api/chat/typing", chatTypingHandler)

	http.HandleFunc("/api/followrequest/reject", rejectFollowerHandler)
	http.HandleFunc("/api/followrequest/accept", approveFollowerHandler)
	http.HandleFunc("/api/followrequestlist", followRequestListHandler)

	http.HandleFunc("/api/posts/get", postsGetHandler)
	http.HandleFunc("/api/post/submit", postNewHandler)

	http.HandleFunc("/api/group/submit", groupNewHandler)
	http.HandleFunc("/api/group/post/submit", groupPostNewHandler)
	http.HandleFunc("/api/group/posts/get", groupPostsGetHandler)
	http.HandleFunc("/api/group/comment/submit", groupCommentNewHandler)
	http.HandleFunc("/api/group/comments/get", groupCommentsGetHandler)
	http.HandleFunc("/api/group/join", groupJoinHandler)
	http.HandleFunc("/api/group/leave", groupLeaveHandler) // TODO: not part of audit, so untested
	http.HandleFunc("/api/group/invite", groupInviteHandler)
	http.HandleFunc("/api/group/invited", groupInvitedHandler)
	http.HandleFunc("/api/group/invite/accept", groupInviteAcceptHandler)
	http.HandleFunc("/api/group/invite/reject", groupInviteRejectHandler)
	http.HandleFunc("/api/group/requests", groupRequestsHandler)
	http.HandleFunc("/api/group/request/accept", groupRequestAcceptHandler)
	http.HandleFunc("/api/group/request/reject", groupRequestRejectHandler)

	http.HandleFunc("/api/event/submit", eventNewHandler)
	http.HandleFunc("/api/events/get", eventsGetHandler)
	http.HandleFunc("/api/event/participants/get", eventParticipantsGetHandler)
	http.HandleFunc("/api/event/attend", eventAttendHandler)
	http.HandleFunc("/api/event/notattend", eventNotAttendHandler)

	http.HandleFunc("/api/user/check", sessionCheckHandler)
	http.HandleFunc("/api/user/following", followingHandler)
	http.HandleFunc("/api/user/followers", followersHandler)
	http.HandleFunc("/api/user/follow", followHandler)
	http.HandleFunc("/api/user/login", userLoginHandler)
	http.HandleFunc("/api/user/logout", userLogoutHandler)
	http.HandleFunc("/api/user/posts", userPostsHandler)
	http.HandleFunc("/api/user/privacy", changePrivacyHandler)
	http.HandleFunc("/api/user/profile", userProfileHandler)
	http.HandleFunc("/api/user/register", userRegisterHandler)
	http.HandleFunc("/api/user/unfollow", unfollowHandler)
}
