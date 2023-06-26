package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

const (
	WS_ERROR          = "error"
	WS_COMMENT_SUBMIT = "comment_submit"
	WS_COMMENTS_LIST  = "comments_list"

	WS_CHAT_USERS_LIST     = "chat_users_list"
	WS_CHAT_MESSAGE_SUBMIT = "chat_message_submit"

	WS_FOLLOW_REQUEST_REJECT         = "follow_request_reject"
	WS_FOLLOW_REQUEST_ACCEPT         = "follow_request_accept"
	WS_FOLLOW_REQUESTS_LIST          = "follow_requests_list"
	WS_POSTS_LIST                    = "posts_list"
	WS_POST_SUBMIT                   = "post_submit"
	WS_GROUPS_LIST                   = "groups_list"
	WS_GROUP_SUBMIT                  = "group_submit"
	WS_GROUP_POST_SUBMIT             = "group_post_submit"
	WS_GROUP_POSTS_LIST              = "group_posts_list"
	WS_GROUP_POST_COMMENT_SUBMIT     = "group_post_comment_submit"
	WS_GROUP_POST_COMMENTS_LIST      = "group_post_comments_list"
	WS_GROUP_JOIN                    = "group_join"
	WS_GROUP_LEAVE                   = "group_leave"
	WS_GROUP_INVITE                  = "group_invite"
	WS_GROUP_INVITED                 = "group_invited"
	WS_GROUP_INVITE_ACCEPT           = "group_invite_accept"
	WS_GROUP_INVITE_REJECT           = "group_invite_reject"
	WS_GROUP_REQUESTS_LIST           = "group_requests_list"
	WS_GROUP_REQUEST_ACCEPT          = "group_request_accept"
	WS_GROUP_REQUEST_REJECT          = "group_request_reject"
	WS_GROUP_EVENT_SUBMIT            = "group_event_submit"
	WS_GROUP_EVENTS_LIST             = "group_events_list"
	WS_GROUP_EVENT_PARTICIPANTS_LIST = "group_event_participants_list"
	WS_GROUP_EVENT_ATTEND            = "group_event_attend"
	WS_GROUP_EVENT_NOT_ATTEND        = "group_event_not_attend"
	WS_USER_CHECK                    = "user_check"
	WS_USER_FOLLOWING_LIST           = "user_following_list"
	WS_USER_FOLLOWERS_LIST           = "user_followers_list"
	WS_USER_FOLLOW                   = "user_follow"
	WS_USER_LOGIN                    = "user_login"
	WS_USER_LOGOUT                   = "user_logout"
	WS_USER_POSTS_LIST               = "user_posts_list"
	WS_USER_PRIVACY                  = "user_privacy"
	WS_USER_PROFILE                  = "user_profile"
	WS_USER_REGISTER                 = "user_register"
	WS_USER_UNFOLLOW                 = "user_unfollow"
)

// # jsonResponse marshals and forwards json response writing to http.ResponseWriter
//
// @params {w http.ResponseWriter, statusCode int, data any}
// @sideEffect {jsonResponse -> w}
func wsJsonMarshal(data any) []byte {

	jsonMarshaledObj := []byte{}
	// if data type is string
	if message, ok := data.(string); ok {
		jsonMarshaledObj, _ = json.Marshal(map[string]string{
			"message": ": " + message,
		})
	}
	// if data type is int
	if message, ok := data.(int); ok {
		jsonMarshaledObj, _ = json.Marshal(map[string]int{
			"message": message,
		})
	}
	// if data type is bool
	if message, ok := data.(bool); ok {
		jsonMarshaledObj, _ = json.Marshal(map[string]bool{
			"message": message,
		})
	}
	// if data type is slice
	if _, ok := data.([]any); ok {
		jsonMarshaledObj, _ = json.Marshal(map[string][]any{
			"data": data.([]any),
		})
	}
	// if data type is object
	if _, ok := data.(map[string]any); ok {
		jsonMarshaledObj, _ = json.Marshal(map[string]any{
			"data": data.(map[string]any),
		})
	}

	return jsonMarshaledObj
}

// # recovery is a utility function to recover from panic and send a json err response over http
//
// @sideEffect {log, debug}
//
// - for further debugging uncomment {print stack trace}
func wsRecover() {
	if r := recover(); r != nil {
		fmt.Println("=====================================")
		stackTrace := debug.Stack()
		lines := strings.Split(string(stackTrace), "\n")
		relevantPanicLines := []string{}
		for _, line := range lines {
			if strings.Contains(line, "backend/") {
				relevantPanicLines = append(relevantPanicLines, line)
			}
		}
		if len(relevantPanicLines) > 1 {
			for i, line := range relevantPanicLines {
				if strings.Contains(line, "utils.go") {
					relevantPanicLines = append(relevantPanicLines[:i], relevantPanicLines[i+1:]...)
				}
			}
		}
		relevantPanicLine := strings.Join(relevantPanicLines, "\n")
		log.Println(relevantPanicLines)
		wsJsonMarshal(relevantPanicLine)
		fmt.Println("=====================================")
		// to print the full stack trace
		log.Println(string(stackTrace))
	}
}

type WSRecover struct {
	messageType string
	data        string
}
