package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type GroupPostRequest struct {
	Title      string `json:"title"`
	Categories string `json:"categories"`
	Content    string `json:"content"`
	Picture    string `json:"picture"`
	GroupID    int    `json:"group_id"`
}

type GroupPosts struct {
	GroupPosts []GroupPost `json:"group_posts"`
}

type GroupPost struct {
	ID         int    `json:"ID"`
	UserID     int    `json:"user_id"`
	Title      string `json:"title"`
	Categories string `json:"categories"`
	Content    string `json:"content"`
}

type GroupPostDTOelement struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	Categories      string `json:"categories"`
	CreatorFullName string `json:"creatorFullName"`
	CreatorEmail    string `json:"creatorEmail"`
	CreatedAt       string `json:"createdAt"`
}

// # groupPostNewHandler create a new group post
//
// @r.param {group_id int, title string, categories string, content string, picture string}
func groupPostNewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()
	var data GroupPostRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request",
		})
		w.Write(jsonResponse)
		return
	}
	cookie, err := r.Cookie("user_uuid")
	if err != nil || cookie.Value == "" || cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "You are not logged in, named cookie not present",
		})
		w.Write(jsonResponse)
		return
	}
	uuid := cookie.Value
	data.Categories = sanitizeCategories(data.Categories)
	userID, err := getIDbyUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "You are not logged in",
		})
		w.Write(jsonResponse)
		return
	}

	// blob the picture
	postPicture := []byte{}
	if data.Picture != "" {
		avatarData, err := base64.StdEncoding.DecodeString(data.Picture)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnprocessableEntity)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "Invalid avatar",
			})
			w.Write(jsonResponse)
			return
		}
		if !isImage(avatarData) {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnsupportedMediaType)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "avatar is not a valid image",
			})
			w.Write(jsonResponse)
			return
		}
		postPicture = avatarData
	}

	var result sql.Result
	result, err = statements["addGroupPost"].Exec(userID, data.Title, data.Categories, data.Content, postPicture, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, addGroupPost failed",
		})
		w.Write(jsonResponse)
		return
	}

	// get id of new group post to make group_post_membership
	postID, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, LastInsertId failed",
		})
		w.Write(jsonResponse)
		return
	}

	// add group_post_membership
	groupID := data.GroupID
	_, err = statements["addGroupPostMembership"].Exec(groupID, postID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, addGroupPostMembership failed",
		})
		w.Write(jsonResponse)
		return
	}

	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "Post created",
	})
	w.Write(jsonResponse)
	// todo: do the websocket sending group_posts with GroupPost type
	// rows, err := statements["getGroupPosts"].Query(userID, groupID)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	w.WriteHeader(500)
	// 	jsonResponse, _ := json.Marshal(map[string]string{
	// 		"message": "internal server error, getGroupPosts failed",
	// 	})
	// 	w.Write(jsonResponse)
	// 	return
	// }
	// todo: make post picture into encoded string and send it

	// todo: 99% it must be remastered too, but 1% still exists :D . Check it precisely
	// var post Post
	// rows.Next()
	// rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Categories, &post.Content)
	// rows.Close()
	// sendPost(post)
}

// # groupPostGetHandler get all group posts
//
// @r.Params {group_id int}
func groupPostsGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()

	// get the uuid of the current user from the cookies
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request. Bad cookie, not tasty",
		})
		w.Write(jsonResponse)
		return
	}

	// get the user id from the uuid
	_, err = getIDbyUUID(cookie.Value)
	if err != nil {
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized. You are not logged in",
		})
		w.Write(jsonResponse)
		return
	}

	// get the posts of the user
	rows, err := statements["getGroupPosts"].Query()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error. getGroupPosts query failed",
		})
		w.Write(jsonResponse)
		return
	}

	// create a slice of posts
	var posts []PostDTOelement

	// iterate over the rows and append the posts to the slice
	for rows.Next() {
		var post PostDTOelement
		var firstName, lastName string
		var pictureBlob []byte
		rows.Scan(&post.ID, &post.Title, &post.Content, &post.Categories, &firstName, &lastName, &post.CreatorEmail, &post.CreatedAt, &pictureBlob)
		post.CreatorFullName = firstName + " " + lastName
		post.Picture = base64.StdEncoding.EncodeToString(pictureBlob)
		posts = append(posts, post)
	}

	// close the rows
	rows.Close()

	// create a map to store the posts
	var postsMap map[string][]PostDTOelement

	// add the posts to the map
	postsMap = map[string][]PostDTOelement{
		"posts": posts,
	}

	// marshal the map into json
	jsonResponse, err := json.Marshal(postsMap)
	if err != nil {
		w.WriteHeader(500)
		// todo: the message bottom looks too strange, for the "userPostsHandler" function
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error. we could not register you at this time",
		})
		w.Write(jsonResponse)
		return
	}

	// write the response
	w.WriteHeader(200)
	w.Write(jsonResponse)

}
