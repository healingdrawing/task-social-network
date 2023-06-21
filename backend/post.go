package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strings"
)

type PostRequest struct {
	Title      string `json:"title"`
	Categories string `json:"categories"`
	Content    string `json:"content"`
	Privacy    string `json:"privacy"`
	CreatedAt  string `json:"created_at"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	ID         int    `json:"ID"`
	UserID     int    `json:"user_id"`
	Title      string `json:"title"`
	Categories string `json:"categories"`
	Content    string `json:"content"`
}

type PostDTOelement struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	Categories      string `json:"categories"`
	CreatorFullName string `json:"creatorFullName"`
	CreatorEmail    string `json:"creatorEmail"`
	CreatedAt       string `json:"createdAt"`
}

func postNewHandler(w http.ResponseWriter, r *http.Request) {
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
	var data PostRequest
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
	// get user id from cookie
	cookie, err := r.Cookie("user_uuid")
	if err != nil || cookie.Value == "" || cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "desired cookie not present",
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
			"message": "You are not logged in, or your cookie is invalid",
		})
		w.Write(jsonResponse)
		return
	}
	_, err = statements["addPost"].Exec(userID, data.Title, data.Categories, data.Content, data.Privacy, data.CreatedAt)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, addPost query failed",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "Post created",
	})
	w.Write(jsonResponse)
	rows, err := statements["getPosts"].Query(userID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, getPosts query failed",
		})
		w.Write(jsonResponse)
		return
	}
	var post Post
	rows.Next()
	rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Categories, &post.Content)
	rows.Close()
	sendPost(post)
}

func postGetHandler(w http.ResponseWriter, r *http.Request) {
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
	rows, err := statements["getPosts"].Query()
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	var posts Posts
	for rows.Next() {
		var post Post
		rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Categories, &post.Content)
		posts.Posts = append(posts.Posts, post)
	}
	rows.Close()
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(posts)
	w.Write(jsonResponse)
}

func sanitizeCategories(data string) string {
	newdata := strings.Split(data, ",")
	returndata := ""
	for i, w := range newdata {
		w = strings.TrimSpace(w)
		if w != "" && i > 0 {
			returndata += (", " + w)
		} else {
			returndata += w
		}
		if returndata == "" {
			returndata = generateRandomEmojiSequence()
		}
	}
	return returndata
}

func generateRandomEmojiSequence() string {
	rounds := []string{"🔴", "🟠", "🟡", "🟢", "🔵", "🟣", "🟤", "⚫", "⚪"}
	// Shuffle the rounds using Fisher-Yates algorithm
	for i := len(rounds) - 1; i > 0; i-- {
		bi := big.NewInt(3)
		bj, err := rand.Int(rand.Reader, bi)
		if err != nil {
			log.Fatal(err)
		}
		// convert big.Int to int
		j := int(bj.Int64())
		rounds[i], rounds[j] = rounds[j], rounds[i]
	}

	// Join the shuffled rounds into a single string
	mixedRounds := strings.Join(rounds, " ")

	return mixedRounds
}

func userPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error. we could not get your posts",
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
	userID, err := getIDbyUUID(cookie.Value)
	if err != nil {
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized. You are not logged in",
		})
		w.Write(jsonResponse)
		return
	}

	// get the posts of the user
	rows, err := statements["getPosts"].Query(userID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error. getPosts query failed",
		})
		w.Write(jsonResponse)
		return
	}

	log.Println("starting to iterate over the rows")

	// create a slice of posts
	var posts []PostDTOelement

	// iterate over the rows and append the posts to the slice
	for rows.Next() {
		var post PostDTOelement
		var firstName, lastName string
		rows.Scan(&post.ID, &post.Title, &post.Content, &post.Categories, &firstName, &lastName, &post.CreatorEmail, &post.CreatedAt)
		post.CreatorFullName = firstName + " " + lastName
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
