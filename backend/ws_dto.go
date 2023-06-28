package main

type WS_RESPONSE_MESSAGE_DTO struct {
	Type WSMT        `json:"type"`
	Data interface{} `json:"data"`
}

type WS_ERROR_RESPONSE_DTO struct {
	Content string `json:"content"`
}

type WS_SUCCESS_RESPONSE_DTO struct {
	Content string `json:"content"`
}

// perhaps refactor bottom code

// type WS_COMMENT_SUBMIT_DTO struct {
// 	Type string `json:"type"`
// 	Data struct {
// 		User_uuid  string `json:"user_uuid"`
// 		Content    string `json:"content"`
// 		Picture    string `json:"picture"`
// 		Created_at string `json:"created_at"`
// 	} `json:"data"`
// }
// type WS_COMMENTS_LIST_DTO struct {
// 	Type string `json:"type"`
// 	Data []struct {
// 		Email      string `json:"email"`
// 		Full_name  string `json:"full_name"`
// 		Content    string `json:"content"`
// 		Picture    string `json:"picture"`
// 		Created_at string `json:"created_at"`
// 	} `json:"data"`
// }

// type WS_CHAT_USERS_LIST_DTO struct {
// 	Type string `json:"type"`
// 	Data []struct {
// 		Email     string `json:"email"`
// 		Full_name string `json:"full_name"`
// 	} `json:"data"`
// }

// type WS_CHAT_MESSAGE_SUBMIT_DTO struct {
// 	Type string `json:"type"`
// 	Data struct {
// 		User_uuid    string `json:"user_uuid"`
// 		To_email     string `json:"to_email"`
// 		To_full_name string `json:"to_full_name"`
// 		Content      string `json:"content"`
// 	} `json:"data"`
// }
// type WS_CHAT_MESSAGE_BROADCAST_DTO struct {
// 	Type string `json:"type"`
// 	Data struct {
// 		From_email     string `json:"from_email"`
// 		From_full_name string `json:"from_full_name"`
// 		Content        string `json:"content"`
// 	} `json:"data"`
// }

// type WS_FOLLOW_REQUEST_SUBMIT_DTO struct { // also used for accept and reject
// 	Type string `json:"type"`
// 	Data struct {
// 		From_email     string `json:"from_email"`
// 		From_full_name string `json:"from_full_name"`
// 		To_email       string `json:"to_email"`
// 		To_full_name   string `json:"to_full_name"`
// 	} `json:"data"`
// }

// type WS_FOLLOW_REQUESTS_LIST_DTO struct {
// 	Type string `json:"type"`
// 	Data []struct {
// 		From_email     string `json:"from_email"`
// 		From_full_name string `json:"from_full_name"`
// 		To_email       string `json:"to_email"`
// 		To_full_name   string `json:"to_full_name"`
// 	} `json:"data"`
// }

// // type WS_POST_SUBMIT_DTO struct {
// // 	Type string `json:"type"`
// // 	Data struct {
// // 		User_uuid   int    `json:"user_uuid"`
// // 		Title       string `json:"title"`
// // 		Content     string `json:"content"`
// // 		Categories  string `json:"categories"`
// // 		Picture     string `json:"picture"`
// // 		Privacy     string `json:"privacy"`
// // 		Able_to_see string `json:"able_to_see"` // only list of emails
// // 	} `json:"data"`
// // }

// // type WS_POSTS_LIST_DTO struct {
// // 	Type string `json:"type"`
// // 	Data []struct {
// // 		Id               int    `json:"id"`
// // 		Author_full_name string `json:"author_full_name"`
// // 		Author_email     string `json:"author_email"`
// // 		Title            string `json:"title"`
// // 		Content          string `json:"content"`
// // 		Categories       string `json:"categories"`
// // 		Picture          string `json:"picture"`
// // 		Privacy          string `json:"privacy"`
// // 	} `json:"data"`
// // }
