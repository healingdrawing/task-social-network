package main

import (
	"encoding/json"
	"log"
	"net/http"
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
	reader(ws)
}

func reader(conn *websocket.Conn) {
	clients.Store(conn, "")
	defer clients.Delete(conn)
	defer conn.Close()
	for {
		messageType, incoming, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("read message: \nincoming bytes: ", incoming,
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
			case string(WS_POST_SUBMIT):
				wsPostSubmitHandler(conn, data.Data)
			case string(WS_POSTS_LIST):
				wsPostsListHandler(conn, data.Data)
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

func wsSendPost(post WS_POST_RESPONSE_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_POST_RESPONSE, post)

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

func wsSendPostsList(postsList WS_POSTS_LIST_DTO) {

	outputMessage, err := wsCreateResponseMessage(WS_POST_RESPONSE, postsList)

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
