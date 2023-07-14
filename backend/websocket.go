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
	reader(uuid, ws)
}

func reader(uuid string, conn *websocket.Conn) {
	clients.Store(uuid, conn)
	defer clients.Delete(uuid)
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

			case string(WS_GROUP_POST_COMMENT_SUBMIT):
				wsGroupPostCommentSubmitHandler(conn, data.Data)
			case string(WS_GROUP_POST_COMMENTS_LIST):
				wsGroupPostCommentsListHandler(conn, data.Data)

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
			case string(WS_USER_GROUPS_FRESH_EVENTS_LIST):
				wsUserGroupsFreshEventsListHandler(conn, data.Data)
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

// send message to connections by uuids provided
func wsSend(message_type WSMT, message interface{}, uuids []string) {
	outputMessage, err := wsCreateResponseMessage(message_type, message)

	if err != nil {
		log.Println(err)
	}

	for _, uuid := range uuids {
		if conn, ok := clients.Load(uuid); ok {
			if c, ok := conn.(*websocket.Conn); ok {
				err = c.WriteMessage(websocket.TextMessage, outputMessage)
				if err != nil {
					log.Println(err)
				}
			} else {
				log.Println("wsSend: clients.Load(uuid) is not *websocket.Conn")
			}
		} else {
			log.Println("wsSend: client not found . clients.Load(uuid) failed")
		}
	}
}

// //////////////////////////
// fragments of old code. remove later if full cleaning will be executed
// //////////////////////////
// todo: remove code bottom only in case of full cleaning from old http implementation only. Can be not safe remove this part only

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
