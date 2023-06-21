package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	uuid "github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

type signupData struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Dob         string `json:"dob"`
	Avatar      string `json:"avatar"`
	avatarBytes []byte `sqlite3:"avatar"`
	Nickname    string `json:"nickname"`
	AboutMe     string `json:"aboutMe"`
	Public      bool   `json:"public"`
	Privacy     string `sqlite3:"privacy"`
}

type ProfileData struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Dob         string `json:"dob"`
	Avatar      string `json:"avatar"`
	avatarBytes []byte `sqlite3:"avatar"`
	Nickname    string `json:"nickname"`
	AboutMe     string `json:"aboutMe"`
	Public      bool   `json:"public"`
	Privacy     string `sqlite3:"privacy"`
}

type ProfileDTOtoFrontend struct {
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Dob       string `json:"dob"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	AboutMe   string `json:"aboutMe"`
	Public    bool   `json:"public"`
}

type loginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UUIDData struct {
	UUID string `json:"UUID"`
}

func changePrivacyHandler(w http.ResponseWriter, r *http.Request) {
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

	cookie, err := r.Cookie("user_uuid")
	if err != nil || cookie.Value == "" || cookie == nil {
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "You are not logged in",
		})
		w.Write(jsonResponse)
		return
	}

	uuid := cookie.Value

	ID, err := getIDbyUUID(uuid)
	if err != nil {
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "You are not logged in",
		})
		w.Write(jsonResponse)
		return
	}

	incomingData := map[string]any{}
	// decode the request body into the DTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&incomingData)
	if err != nil {
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
		return
	}
	wantPublic := incomingData["public"].(bool)

	privacyvalue := ""
	if wantPublic == true {
		privacyvalue = "public"
	} else {
		privacyvalue = "private"
	}

	_, err = statements["updateUserPrivacy"].Exec(privacyvalue, ID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}

	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "Privacy updated",
	})
	w.Write(jsonResponse)
}

func userProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	// get uuid from the cookie
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "User not found",
		})
		w.Write(jsonResponse)
		return
	}
	myuuid := cookie.Value

	// get the ID of the user that is currently logged in
	myID, err := getIDbyUUID(myuuid)

	// get the email from the json request
	var incomingData map[string]any
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&incomingData)
	if err != nil && err.Error() != "EOF" {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request",
		})
		w.Write(jsonResponse)
		return
	}
	incomingEmail := ""
	if incomingData != nil {
		incomingEmail = incomingData["email"].(string)
	}
	if incomingEmail == "" {
		// get the email of the user that is currently logged in
		rows, err := statements["getEmailByID"].Query(myID)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(401)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "User not found",
			})
			w.Write(jsonResponse)
			return
		}
		defer rows.Close()
		rows.Next()
		rows.Scan(&incomingEmail)
		rows.Close()
	}

	log.Println(incomingEmail)

	// get the ID of the user that we want to see
	ID, err := getIDbyEmail(incomingEmail)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "User not found",
		})
		w.Write(jsonResponse)
		return
	}

	var profile ProfileData
	rows, err := statements["getUserProfile"].Query(ID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "User not found",
		})
		w.Write(jsonResponse)
		return
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&profile.Email, &profile.FirstName, &profile.LastName, &profile.Dob,
		&profile.avatarBytes, &profile.Nickname, &profile.AboutMe, &profile.Privacy)
	rows.Close()

	profileDTO := ProfileDTOtoFrontend{}

	profileDTO.Email = profile.Email
	profileDTO.FirstName = profile.FirstName
	profileDTO.LastName = profile.LastName
	profileDTO.Dob = profile.Dob
	profileDTO.Avatar = base64.StdEncoding.EncodeToString(profile.avatarBytes)
	profileDTO.Nickname = profile.Nickname
	profileDTO.AboutMe = profile.AboutMe
	if profile.Privacy == "public" {
		profileDTO.Public = true
	} else {
		profileDTO.Public = false
	}

	if myID == ID {
		w.WriteHeader(200)
		jsonResponse, _ := json.Marshal(profileDTO)
		w.Write(jsonResponse)
		return
	}

	// if the profile is private, check if I am following the user
	// if I am not following the user, return 401
	if profile.Privacy == "private" {
		following, err := isFollowing(myID, ID)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error, doesSecondFollowFirst failed",
			})
			w.Write(jsonResponse)
			return
		}
		if following == false {
			w.WriteHeader(401)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "User is private, you cannot check their profile without following them",
			})
			w.Write(jsonResponse)
			return
		}
	}

	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(profileDTO)
	w.Write(jsonResponse)
}

func isFollowing(myID int, ID int) (bool, error) {
	rows, err := statements["doesSecondFollowFirst"].Query(ID, myID)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	// check length of rows
	// if length is 0, return false
	// else return true
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error. we could not register you at this time",
			})
			w.Write(jsonResponse)
		}
	}()

	var data signupData
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request",
		})
		w.Write(jsonResponse)
		return
	}

	// if dob is empty, return error
	if data.Dob == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Invalid date of birth",
		})
		w.Write(jsonResponse)
		return
	}

	// check the avatar validity
	if data.Avatar != "" {
		avatarData, err := base64.StdEncoding.DecodeString(data.Avatar)
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
		data.avatarBytes = avatarData
	}

	if data.Avatar == "" {
		rn := randomNum(0, 5)
		defaultAvatar, err := os.Open("./assets/images/profile/defaults/" + strconv.Itoa(rn) + ".jpeg")
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusFailedDependency)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "failed to load default avatar",
			})
			w.Write(jsonResponse)
			return
		}
		defer defaultAvatar.Close()
		data.avatarBytes, err = ioutil.ReadAll(defaultAvatar)
		defaultAvatar.Close()
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusFailedDependency)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "default avatar is not a valid image",
			})
			w.Write(jsonResponse)
			return
		}
		defaultAvatar.Close()
	}

	if data.Public == true {
		data.Privacy = "public"
	} else {
		data.Privacy = "private"
	}

	onlyEnglishRegex := regexp.MustCompile(`^[a-zA-Z0-9]{2,15}$`)

	// todo: it looks like not required(but not sure), because of nickname being optional, and does not provide way to login anymore
	if data.Nickname != "" {
		if !onlyEnglishRegex.MatchString(data.Nickname) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": `Invalid nickname: ` + data.Nickname + `
Nickname must only contain english letters and numbers.
Nickname must be between 2 and 15 characters long.`,
			})
			w.Write(jsonResponse)
			return
		}
	}

	if len(data.FirstName) < 1 || len(data.FirstName) > 32 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": `Invalid first name length: ` + data.FirstName + `
First name must be between 1 and 32 characters long`,
		})
		w.Write(jsonResponse)
		return
	}

	if len(data.LastName) < 1 || len(data.LastName) > 32 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": `Invalid last name length: ` + data.LastName + `
Last name must be between 1 and 32 characters long`,
		})
		w.Write(jsonResponse)
		return
	}

	emailRegex := regexp.MustCompile(`^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,4})+$`)

	if !emailRegex.MatchString(strings.ToLower((data.Email))) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": `Invalid email: ` + data.Email + `
Email must be a valid email address`,
		})
		w.Write(jsonResponse)
		return
	}

	if len(data.Password) < 6 || len(data.Password) > 15 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": `Invalid password length: ` + data.Password + `
Password must be between 6 and 15 characters long`,
		})
		w.Write(jsonResponse)
		return
	}

	if !onlyEnglishRegex.MatchString(data.Password) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": `Invalid password: ` + data.Password + `
Password must only contain english characters and numbers`,
		})
		w.Write(jsonResponse)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusFailedDependency)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "dependency failure: could not hash password",
		})
		w.Write(jsonResponse)
		return
	}
	_, err = statements["addUser"].Exec(data.Email, string(hash), data.FirstName, data.LastName, data.Dob, data.avatarBytes, data.Nickname, data.AboutMe, data.Privacy)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			log.Println(err.Error())
			w.WriteHeader(http.StatusConflict)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "This email is already taken",
			})
			w.Write(jsonResponse)
			return
		}

		log.Println(err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "database entry for adding user failed",
		})
		w.Write(jsonResponse)
		return
	}
	UUID, err := createSession(data.Email)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"UUID":  UUID,
		"email": data.Email,
	})
	w.Write(jsonResponse)
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {

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
	var data loginData
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request. The JSON body is not as expected",
		})
		w.Write(jsonResponse)
		return
	}

	var email, hash string
	rows, err := statements["getUserCredentials"].Query(data.Email)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Invalid credentials",
		})
		w.Write(jsonResponse)
		return
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&email, &hash)
	rows.Close()

	if email == "" || hash == "" {
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Invalid credentials. Email or password is incorrect",
		})
		w.Write(jsonResponse)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(data.Password))
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Invalid credentials. Forgot password?",
		})
		w.Write(jsonResponse)
		return
	}

	UUID, err := createSession(email)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"UUID":  UUID,
		"email": email,
	})
	w.Write(jsonResponse)
}

func sessionCheckHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	var data UUIDData
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request",
		})
		w.Write(jsonResponse)
		return
	}
	rows, err := statements["getSession"].Query(data.UUID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	defer rows.Close()
	if !rows.Next() {
		w.WriteHeader(200)
		jsonResponse, _ := json.Marshal(map[string]bool{
			"Exists": false,
		})
		w.Write(jsonResponse)
		return
	}
	rows.Close()
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]bool{
		"Exists": true,
	})
	w.Write(jsonResponse)
}

func userLogoutHandler(w http.ResponseWriter, r *http.Request) {
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

	cookie, err := r.Cookie("user_uuid")
	if err != nil || cookie.Value == "" || cookie == nil {
		w.WriteHeader(200)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "You are not logged in",
		})
		w.Write(jsonResponse)
		return
	}

	uuid := cookie.Value

	_, err = statements["removeSession"].Exec(uuid)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "Session deleted",
	})
	w.Write(jsonResponse)
}

func createSession(email string) (UUID string, err error) {
	random, _ := uuid.NewV4()
	UUID = random.String()
	ID, err := getIDbyEmail(email)
	if err != nil {
		return "", err
	}
	_, err = statements["addSession"].Exec(UUID, ID)
	if err != nil {
		return "", err
	}
	return UUID, nil
}

func getIDbyEmail(email string) (ID int, err error) {
	rows, err := statements["getUserID"].Query(email)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&ID)
	rows.Close()
	return ID, nil
}

// # This function is used to retrieve the ID of the user that is currently logged in
//
// It does so by retrieving the UUID from the request body and then calling getIDbyUUID.
// This function retrieves an ID based on a given UUID by joining the session table with the users table
// and then selecting the ID from the users table
func getIDbyUUID(UUID string) (ID int, err error) {
	rows, err := statements["getIDbyUUID"].Query(UUID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&ID)
	rows.Close()
	return ID, nil
}

func getUserEmailbyID(ID int) (email string, err error) {
	rows, err := statements["getEmailByID"].Query(ID)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&email)
	rows.Close()
	return email, nil
}

func isImage(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return true // JPEG
	}

	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return true // PNG
	}

	if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x38 {
		return true // GIF
	}

	return false
}

func randomNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
