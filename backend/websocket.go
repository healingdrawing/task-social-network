package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	clients  = &sync.Map{}
)

type wsInput struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type wsPost struct {
	Type string `json:"type"`
	Data Post   `json:"post"`
}

type wsComment struct {
	Type   string  `json:"type"`
	PostID int     `json:"postID"`
	Data   Comment `json:"comment"`
}

type wsStatus struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Online   bool   `json:"online"`
}

type wsMessage struct {
	Type    string  `json:"type"`
	Message Message `json:"message"`
}

type wsTyping struct {
	Type         string `json:"type"`
	UsernameFrom string `json:"usernameFrom"`
	Typing       bool   `json:"typing"`
}

type wsError struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}

type wsGroupPosts struct {
	Type  string              `json:"type"`
	Posts []PostDTOoutElement `json:"posts"`
}

type wsGroupPost struct {
	Type string `json:"type"`
	Post Post   `json:"post"`
}

type wsGroupPostComment struct {
	Type    string       `json:"type"`
	PostID  int          `json:"postID"`
	Comment GroupComment `json:"group_comment"`
}

type wsGroupPostComments struct {
	Type          string         `json:"type"`
	PostID        int            `json:"postID"`
	GroupComments []GroupComment `json:"group_comments"`
}

type wsNotification struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func wsConnection(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	uuid := strings.TrimSpace(r.URL.Query().Get("uuid"))
	log.Println("wsConnection uuid: ", uuid) //todo: delete debug
	if uuid == "" {
		log.Println("====================================")
		log.Println("uuid is empty")
		log.Println("====================================")
		return
	}
	reader(ws, uuid)
}

func reader(conn *websocket.Conn, uuid string) {
	clients.Store(conn, uuid)
	defer clients.Delete(conn)
	defer conn.Close()
	for {
		messageType, incoming, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			log.Println("=== error in reader , before delete and close ws ===")
			return
		}

		log.Println("=================\nread message:",
			"\nincoming as string:", string(incoming),
			"\nmessageType: ", messageType) //todo: delete debug

		if messageType == websocket.TextMessage {
			log.Println("Text message received")
			var data wsInput
			if err := json.Unmarshal(incoming, &data); err != nil {
				log.Println(err)
				return
			}

			log.Println("data after unmarshalling: ", data) //todo: delete debug

			switch data.Type {
			case string(WS_GROUP_SUBMIT):
				wsGroupSubmitHandler(conn, data.Data)
			case string(WS_GROUPS_LIST):
				wsGroupsListHandler(conn, data.Data)
			case string(WS_GROUPS_ALL_LIST):
				wsGroupsAllListHandler(conn, data.Data)

			case string(WS_GROUP_POST_SUBMIT):
				wsGroupPostSubmitHandler(conn, data.Data)
			case string(WS_GROUP_POSTS_LIST):
				wsGroupPostsListHandler(conn, data.Data)
			case string(WS_USER_GROUP_POSTS_LIST):
				wsUserGroupPostsListHandler(conn, data.Data)

			case string(WS_POST_SUBMIT):
				wsPostSubmitHandler(conn, data.Data)
			case string(WS_POSTS_LIST):
				wsPostsListHandler(conn, data.Data)
			case string(WS_USER_POSTS_LIST):
				wsUserPostsListHandler(conn, data.Data)

			case string(WS_COMMENT_SUBMIT):
				wsCommentSubmitHandler(conn, data.Data)
			case string(WS_COMMENTS_LIST):
				wsCommentsListHandler(conn, data.Data)

			case string(WS_USER_PROFILE):
				wsUserProfileHandler(conn, data.Data)
			case string(WS_USER_PRIVACY):
				wsChangePrivacyHandler(conn, data.Data)

			case string(WS_USER_FOLLOWING_LIST):
				wsFollowingListHandler(conn, data.Data)
			case string(WS_USER_FOLLOWERS_LIST):
				wsFollowersListHandler(conn, data.Data)
			case string(WS_USER_FOLLOW):
				wsFollowHandler(conn, data.Data)
			case string(WS_USER_UNFOLLOW):
				wsUnfollowHandler(conn, data.Data)

			case string(WS_FOLLOW_REQUESTS_LIST):
				wsFollowRequestsListHandler(conn, data.Data)
			case string(WS_FOLLOW_REQUEST_ACCEPT):
				wsAcceptFollowerHandler(conn, data.Data)
			case string(WS_FOLLOW_REQUEST_REJECT):
				wsRejectFollowerHandler(conn, data.Data)

			case string(WS_GROUP_REQUEST_SUBMIT): // for frontend button GroupView.vue
				wsGroupRequestSubmitHandler(conn, data.Data)
			case string(WS_GROUP_REQUEST_ACCEPT):
				wsGroupRequestAcceptHandler(conn, data.Data)
			case string(WS_GROUP_REQUEST_REJECT):
				wsGroupRequestRejectHandler(conn, data.Data)
			case string(WS_GROUP_REQUESTS_LIST):
				wsGroupRequestsListHandler(conn, data.Data)

			case string(WS_GROUP_INVITES_SUBMIT): // string of emails space separated
				wsGroupInvitesSubmitHandler(conn, data.Data)
			case string(WS_GROUP_INVITE_ACCEPT):
				wsGroupInviteAcceptHandler(conn, data.Data)
			case string(WS_GROUP_INVITE_REJECT):
				wsGroupInviteRejectHandler(conn, data.Data)
			case string(WS_GROUP_INVITES_LIST):
				wsGroupInvitesListHandler(conn, data.Data)

			case string(WS_GROUP_EVENT_SUBMIT):
				wsGroupEventSubmitHandler(conn, data.Data)
			case string(WS_GROUP_EVENTS_LIST):
				wsGroupEventsListHandler(conn, data.Data)
			case string(WS_GROUP_EVENT_GOING):
				wsGroupEventGoingHandler(conn, data.Data)
			case string(WS_GROUP_EVENT_NOT_GOING):
				wsGroupEventNotGoingHandler(conn, data.Data)

			case string(WS_USER_VISITOR_STATUS):
				wsUserVisitorStatusHandler(conn, data.Data)
			case string(WS_USER_GROUP_VISITOR_STATUS):
				wsUserGroupVisitorStatusHandler(conn, data.Data)

			case "login":
				clients.Store(conn, data.Data["username"])
				sendStatus(data.Data["username"].(string), true)
				defer sendStatus(data.Data["username"].(string), false)
			case "logout":
				conn.Close()
				clients.Delete(conn)
				sendStatus(data.Data["username"].(string), false)
			default:
				log.Println("Unknown type: ", data.Type)
			}
		}
		if messageType == websocket.BinaryMessage {
			log.Println("Binary message received")
		}
	}
}

//todo: there is tiny chance it can be simplified. Check after all wsSend... is done

func wsSendError(msg WS_ERROR_RESPONSE_DTO) {
	outputMessage, err := wsCreateResponseMessage(WS_ERROR_RESPONSE, msg)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func wsSendSuccess(msg WS_SUCCESS_RESPONSE_DTO) {
	outputMessage, err := wsCreateResponseMessage(WS_SUCCESS_RESPONSE, msg)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

// func wsSendPost(post WS_POST_RESPONSE_DTO) {

// 	outputMessage, err := wsCreateResponseMessage(WS_POST_RESPONSE, post)

// 	if err != nil {
// 		log.Println(err)
// 	}
// 	clients.Range(func(key, value interface{}) bool {
// 		if c, ok := key.(*websocket.Conn); ok {
// 			err = c.WriteMessage(websocket.TextMessage, outputMessage)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 		}
// 		return true
// 	})
// }

func wsSendPostsList(posts_list WS_POSTS_LIST_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_POSTS_LIST, posts_list)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func wsSendGroupPostsList(group_posts_list WS_GROUP_POSTS_LIST_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_GROUP_POSTS_LIST, group_posts_list)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func wsSendGroupsList(groups_list WS_GROUPS_LIST_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_GROUPS_LIST, groups_list)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

// func wsSendComment(comment WS_COMMENT_RESPONSE_DTO) {

// 	outputMessage, err := wsCreateResponseMessage(WS_COMMENT_RESPONSE, comment)

// 	if err != nil {
// 		log.Println(err)
// 	}
// 	clients.Range(func(key, value interface{}) bool {
// 		if c, ok := key.(*websocket.Conn); ok {
// 			err = c.WriteMessage(websocket.TextMessage, outputMessage)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 		}
// 		return true
// 	})
// }

func wsSendCommentsList(comments_list WS_COMMENTS_LIST_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_COMMENTS_LIST, comments_list)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func wsSendUserProfile(profile WS_USER_PROFILE_RESPONSE_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_USER_PROFILE, profile)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func wsSendFollowingList(following_list WS_FOLLOWING_LIST_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_USER_FOLLOWING_LIST, following_list)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func wsSendFollowersList(following_list WS_FOLLOWERS_LIST_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_USER_FOLLOWERS_LIST, following_list)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func wsSendFollowRequestsList(follow_requests_list WS_FOLLOW_REQUESTS_LIST_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_FOLLOW_REQUESTS_LIST, follow_requests_list)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func wsSendUserVisitorStatus(user_visitor_status WS_USER_VISITOR_STATUS_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_USER_VISITOR_STATUS, user_visitor_status)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func wsSendUserGroupVisitorStatus(user_group_visitor_status WS_USER_VISITOR_STATUS_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_USER_GROUP_VISITOR_STATUS, user_group_visitor_status)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func wsSendGroupInvitesList(invites_list WS_GROUP_INVITES_LIST_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_GROUP_INVITES_LIST, invites_list)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func wsSendGroupRequestsList(requests_list WS_GROUP_REQUESTS_LIST_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_GROUP_REQUESTS_LIST, requests_list)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func wsSendGroupEventsList(events_list WS_GROUP_EVENTS_LIST_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_GROUP_EVENTS_LIST, events_list)

	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, outputMessage)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

////////////////////////////
// old code
////////////////////////////

func sendPost(post Post) {
	data := wsPost{"post", post}
	output, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, output)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func sendComment(postID int, comment Comment) {
	data := wsComment{"comment", postID, comment}
	output, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			err = c.WriteMessage(websocket.TextMessage, output)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func sendStatus(username string, online bool) {
	data := wsStatus{"status", username, online}
	output, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if value.(string) != "" {
			err = key.(*websocket.Conn).WriteMessage(websocket.TextMessage, output) // todo: CHECK IT! err was added, not sure it is correct
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func sendMessage(message Message) {
	data := wsMessage{"message", message}
	output, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if value.(string) == message.UsernameFrom || value.(string) == message.UsernameTo {
			err = key.(*websocket.Conn).WriteMessage(websocket.TextMessage, output) // todo: CHECK IT! err was added, not sure it is correct
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func sendTyping(typing Typing) {
	data := wsTyping{"typing", typing.UsernameFrom, typing.Typing}
	output, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if value.(string) == typing.UsernameTo {
			err = key.(*websocket.Conn).WriteMessage(websocket.TextMessage, output) // todo: CHECK IT! err was added, not sure it is correct
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
}
