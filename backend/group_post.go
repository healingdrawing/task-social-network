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
	defer recovery(w)
	userID, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	var data GroupPostRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnprocessableEntity, "bad request")
		return
	}
	data.Categories = sanitizeCategories(data.Categories)

	// blob the picture
	postPicture := []byte{}
	if data.Picture != "" {
		avatarData, err := base64.StdEncoding.DecodeString(data.Picture)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusUnprocessableEntity, "Invalid avatar")
			return
		}
		if !isImage(avatarData) {
			log.Println(err.Error())
			jsonResponse(w, http.StatusUnsupportedMediaType, "avatar is not a valid image")
			return
		}
		postPicture = avatarData
	}

	var result sql.Result
	result, err = statements["addGroupPost"].Exec(userID, data.Title, data.Categories, data.Content, postPicture, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " addGroupPost query failed")
		return
	}

	// get id of new group post to make group_post_membership
	postID, err := result.LastInsertId()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "LastInsertId failed")
		return
	}

	// add group_post_membership
	groupID := data.GroupID
	_, err = statements["addGroupPostMembership"].Exec(groupID, postID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "addGroupPostMembership failed")
		return
	}

	jsonResponse(w, http.StatusOK, "Post created")
	// todo: do the websocket sending group_posts with GroupPost type
	// rows, err := statements["getGroupPosts"].Query(userID, groupID)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	w.WriteHeader(500)
	// 	jsonResponseObj, _ := json.Marshal(map[string]string{
	// 		"message": "internal server error, getGroupPosts failed",
	// 	})
	// 	w.Write(jsonResponseObj)
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
	defer recovery(w)

	// check session from cookie
	_, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	// get the posts of the user
	rows, err := statements["getGroupPosts"].Query()
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "getGroupPosts query failed")
		return
	}

	// create a slice of posts
	var posts []PostDTOelement

	// iterate over the rows and append the posts to the slice
	for rows.Next() {
		var post PostDTOelement
		var firstName, lastName string
		var pictureBlob []byte
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Categories, &firstName, &lastName, &post.CreatorEmail, &post.CreatedAt, &pictureBlob)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "posts scan failed")
			return
		}
		post.CreatorFullName = firstName + " " + lastName
		post.Picture = base64.StdEncoding.EncodeToString(pictureBlob)
		posts = append(posts, post)
	}

	// close the rows
	rows.Close()

	// add the posts to the map
	postsMap := map[string][]PostDTOelement{
		"posts": posts,
	}

	// marshal the map into json
	jsonResponseObj, err := json.Marshal(postsMap)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "json marshal failed")
		return
	}

	// write the response
	w.WriteHeader(200)
	w.Write(jsonResponseObj)
}
