package main

func registerHandlers() {
	// Websocket
	CustomRouter.HandleFunc("/ws", wsConnection)

	// API
	CustomRouter.HandleFunc("/api/comment/submit", commentNewHandler)
	CustomRouter.HandleFunc("/api/comments/get", commentsGetHandler)

	CustomRouter.HandleFunc("/api/chat/getmessages", chatMessagesHandler)
	CustomRouter.HandleFunc("/api/chat/getusers", chatUsersHandler)
	CustomRouter.HandleFunc("/api/chat/newmessage", chatNewHandler)
	CustomRouter.HandleFunc("/api/chat/typing", chatTypingHandler)

	CustomRouter.HandleFunc("/api/followrequest/reject", rejectFollowerHandler)
	CustomRouter.HandleFunc("/api/followrequest/accept", approveFollowerHandler)
	CustomRouter.HandleFunc("/api/followrequestlist", followRequestListHandler)

	CustomRouter.HandleFunc("/api/posts/get", postsGetHandler)
	CustomRouter.HandleFunc("/api/post/submit", postNewHandler)

	CustomRouter.HandleFunc("/api/group/submit", groupNewHandler)
	CustomRouter.HandleFunc("/api/group/post/submit", groupPostNewHandler)
	CustomRouter.HandleFunc("/api/group/posts/get", groupPostsGetHandler)
	CustomRouter.HandleFunc("/api/group/comment/submit", groupCommentNewHandler)
	CustomRouter.HandleFunc("/api/group/comments/get", groupCommentsGetHandler)
	CustomRouter.HandleFunc("/api/group/join", groupJoinHandler)
	CustomRouter.HandleFunc("/api/group/leave", groupLeaveHandler) // TODO: not part of audit, so untested
	CustomRouter.HandleFunc("/api/group/invite", groupInviteHandler)
	CustomRouter.HandleFunc("/api/group/invited", groupInvitedHandler)
	CustomRouter.HandleFunc("/api/group/invite/accept", groupInviteAcceptHandler)
	CustomRouter.HandleFunc("/api/group/invite/reject", groupInviteRejectHandler)
	CustomRouter.HandleFunc("/api/group/requests", groupRequestsHandler)
	CustomRouter.HandleFunc("/api/group/request/accept", groupRequestAcceptHandler)
	CustomRouter.HandleFunc("/api/group/request/reject", groupRequestRejectHandler)

	CustomRouter.HandleFunc("/api/event/submit", eventNewHandler)
	CustomRouter.HandleFunc("/api/events/get", eventsGetHandler)
	CustomRouter.HandleFunc("/api/event/participants/get", eventParticipantsGetHandler)
	CustomRouter.HandleFunc("/api/event/attend", eventAttendHandler)
	CustomRouter.HandleFunc("/api/event/notattend", eventNotAttendHandler)

	CustomRouter.HandleFunc("/api/user/check", sessionCheckHandler)
	CustomRouter.HandleFunc("/api/user/following", followingHandler)
	CustomRouter.HandleFunc("/api/user/followers", followersHandler)
	CustomRouter.HandleFunc("/api/user/follow", followHandler)
	CustomRouter.HandleFunc("/api/user/login", userLoginHandler)
	CustomRouter.HandleFunc("/api/user/logout", userLogoutHandler)
	CustomRouter.HandleFunc("/api/user/posts", userPostsHandler)
	CustomRouter.HandleFunc("/api/user/privacy", changePrivacyHandler)
	CustomRouter.HandleFunc("/api/user/profile", userProfileHandler)
	CustomRouter.HandleFunc("/api/user/register", userRegisterHandler)
	CustomRouter.HandleFunc("/api/user/unfollow", unfollowHandler)
}
