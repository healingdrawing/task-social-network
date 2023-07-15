package main

func registerHandlers() {
	// Websocket
	CustomRouter.HandleFunc("/ws", wsConnection)

	// API
	CustomRouter.HandleFunc("/api/user/check", sessionCheckHandler)
	CustomRouter.HandleFunc("/api/user/login", userLoginHandler)
	CustomRouter.HandleFunc("/api/user/logout", userLogoutHandler)
	CustomRouter.HandleFunc("/api/user/register", userRegisterHandler)
}
