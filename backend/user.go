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
	defer recovery(w)

	cookie, err := r.Cookie("user_uuid")
	if err != nil || cookie.Value == "" || cookie == nil {
		jsonResponse(w, http.StatusUnauthorized, "You are not logged in")
		return
	}

	uuid := cookie.Value

	ID, err := getIDbyUUID(uuid)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "You are not logged in")
		return
	}

	incomingData := map[string]any{}
	// decode the request body into the DTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&incomingData)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	wantPublic := incomingData["public"].(bool)

	privacyvalue := map[bool]string{true: "public", false: "private"}[wantPublic]

	_, err = statements["updateUserPrivacy"].Exec(privacyvalue, ID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "updateUserPrivacy query failed")
		return
	}

	jsonResponse(w, http.StatusOK, "Privacy updated")
}

func userProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)

	// get uuid from the cookie
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnauthorized, "cookie not found")
		return
	}
	myuuid := cookie.Value

	// get the ID of the user that is currently logged in
	myID, err := getIDbyUUID(myuuid)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "You are not logged in")
		return
	}

	// get the email from the json request
	var incomingData map[string]any
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&incomingData)
	if err != nil && err.Error() != "EOF" {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "Bad request")
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
			jsonResponse(w, http.StatusUnauthorized, "getEmailByID query failed, user not found")
			return
		}
		defer rows.Close()
		rows.Next()
		err = rows.Scan(&incomingEmail)
		if err != nil {
			jsonResponse(w, http.StatusUnauthorized, "getEmailByID query failed, incomingEmail scan failed")
			return
		}
		rows.Close()
	}

	log.Println(incomingEmail)

	// get the ID of the user that we want to see
	ID, err := getIDbyEmail(incomingEmail)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnauthorized, "User not found, getIDbyEmail failed")
		return
	}

	var profile ProfileData
	rows, err := statements["getUserProfile"].Query(ID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnauthorized, "getUserProfile query failed, user not found")
		return
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&profile.Email, &profile.FirstName, &profile.LastName, &profile.Dob,
		&profile.avatarBytes, &profile.Nickname, &profile.AboutMe, &profile.Privacy)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "getUserProfile -> profile scan failed")
		return
	}
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
		jsonResponseObj, _ := json.Marshal(profileDTO)
		_, err = w.Write(jsonResponseObj)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj) <- profileDTO failed")
		}
		return
	}

	// if the profile is private, check if I am following the user
	// if I am not following the user, return 401
	if profile.Privacy == "private" {
		following, err := isFollowing(myID, ID)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, "isFollowing failed")
			return
		}
		if !following {
			jsonResponse(w, http.StatusUnauthorized, "User is private, you cannot check their profile without following them")
			return
		}
	}

	w.WriteHeader(200)
	jsonResponseObj, err := json.Marshal(profileDTO)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "json.Marshal(profileDTO) failed")
		return
	}
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj) <- profileDTO failed")
	}

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
	defer recovery(w)

	var data signupData
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}

	// if dob is empty, return error
	if data.Dob == "" {
		jsonResponse(w, http.StatusUnprocessableEntity, "Invalid date of birth")
		return
	}

	// check the avatar validity
	if data.Avatar != "" {
		avatarData, err := base64.StdEncoding.DecodeString(data.Avatar)
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
		data.avatarBytes = avatarData
	}

	if data.Avatar == "" {
		rn := randomNum(0, 5)
		defaultAvatar, err := os.Open("./assets/images/profile/defaults/" + strconv.Itoa(rn) + ".jpeg")
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusFailedDependency, "failed to load default avatar")
			return
		}
		defer defaultAvatar.Close()
		data.avatarBytes, err = ioutil.ReadAll(defaultAvatar)
		defaultAvatar.Close()
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusFailedDependency, "failed to load default avatar")
			return
		}
		defaultAvatar.Close()
	}

	// if data.Public {
	// 	data.Privacy = "public"
	// } else {
	// 	data.Privacy = "private"
	// }

	data.Privacy = map[bool]string{true: "public", false: "private"}[data.Public]

	onlyEnglishRegex := regexp.MustCompile(`^[a-zA-Z0-9]{2,15}$`)

	if data.Nickname != "" {
		if !onlyEnglishRegex.MatchString(data.Nickname) {
			message := `Invalid nickname: ` + data.Nickname + `
			Nickname must only contain english letters and numbers.
			Nickname must be between 2 and 15 characters long.`
			jsonResponse(w, http.StatusUnprocessableEntity, message)
			return
		}
	}

	if len(data.FirstName) < 1 || len(data.FirstName) > 32 {
		message := `Invalid first name length: ` + data.FirstName + `
		First name must be between 1 and 32 characters long`
		jsonResponse(w, http.StatusUnprocessableEntity, message)
		return
	}

	if len(data.LastName) < 1 || len(data.LastName) > 32 {
		message := `Invalid last name length: ` + data.LastName + `
		Last name must be between 1 and 32 characters long`
		jsonResponse(w, http.StatusUnprocessableEntity, message)
		return
	}

	emailRegex := regexp.MustCompile(`^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,4})+$`)

	if !emailRegex.MatchString(strings.ToLower((data.Email))) {
		message := `Invalid email: ` + data.Email + `
		Email must be a valid email address`
		jsonResponse(w, http.StatusUnprocessableEntity, message)
		return
	}

	if len(data.Password) < 6 || len(data.Password) > 15 {
		message := `Invalid password length: ` + data.Password + `
		Password must be between 6 and 15 characters long`
		jsonResponse(w, http.StatusUnprocessableEntity, message)
		return
	}

	if !onlyEnglishRegex.MatchString(data.Password) {
		message := `Invalid password: ` + data.Password + `
Password must only contain english characters and numbers`
		jsonResponse(w, http.StatusUnprocessableEntity, message)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnprocessableEntity, "dependency failure: could not hash password")
		return
	}
	_, err = statements["addUser"].Exec(data.Email, string(hash), data.FirstName, data.LastName, data.Dob, data.avatarBytes, data.Nickname, data.AboutMe, data.Privacy)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			log.Println(err.Error())
			jsonResponse(w, http.StatusConflict, "This email is already taken")
			return
		}

		log.Println(err.Error())
		jsonResponse(w, http.StatusUnprocessableEntity, "database entry for adding user failed")
		return
	}
	UUID, err := createSession(data.Email)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, " create session failed")
		return
	}
	w.WriteHeader(200)
	jsonResponseObj, _ := json.Marshal(map[string]string{
		"UUID":  UUID,
		"email": data.Email,
	})
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj)<-UUID,email failed")
		return
	}

}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	var data loginData
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnprocessableEntity, "Bad request. The JSON body is not as expected")
		return
	}

	var email, hash string
	rows, err := statements["getUserCredentials"].Query(data.Email)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&email, &hash)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "scan credentials failed")
		return
	}
	rows.Close()

	if email == "" || hash == "" {
		jsonResponse(w, http.StatusUnauthorized, "Invalid credentials. Email or password is incorrect")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(data.Password))
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnauthorized, "Invalid credentials. Forgot password?")
		return
	}

	UUID, err := createSession(email)

	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "create session failed")
		return
	}
	w.WriteHeader(200)
	jsonResponseObj, _ := json.Marshal(map[string]string{
		"UUID":  UUID,
		"email": email,
	})
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj)<-UUID,email failed")
		return
	}
}

func sessionCheckHandler(w http.ResponseWriter, r *http.Request) {
	defer recovery(w)
	var data UUIDData
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}
	rows, err := statements["getSession"].Query(data.UUID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}
	defer rows.Close()
	if !rows.Next() {
		w.WriteHeader(200)
		jsonResponseObj, err := json.Marshal(map[string]bool{
			"Exists": false,
		})
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "json.Marshal(map[string]bool{\"Exists\": false}) failed")
			return
		}
		_, err = w.Write(jsonResponseObj)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj)<-Exists:false failed")
		}
		return
	}
	rows.Close()
	w.WriteHeader(200)
	jsonResponseObj, _ := json.Marshal(map[string]bool{
		"Exists": true,
	})
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj)<-Exists:true failed")
		return
	}

}

func userLogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer recovery(w)

	cookie, err := r.Cookie("user_uuid")
	if err != nil || cookie.Value == "" || cookie == nil {
		jsonResponse(w, http.StatusOK, "You are not logged in")
		return
	}

	uuid := cookie.Value

	_, err = statements["removeSession"].Exec(uuid)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "removeSession query failed")
		return
	}
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
	jsonResponse(w, http.StatusOK, "Session deleted")
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
	rows, err := statements["getUserIDByEmail"].Query(email)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&ID)
	if err != nil {
		return 0, err
	}
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
	err = rows.Scan(&ID)
	if err != nil {
		return 0, err
	}
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
	err = rows.Scan(&email)
	if err != nil {
		return "", err
	}
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

// todo: CHECK IT! , it is refactored to prevent warning
func randomNum(min, max int) int {
	rng := rand.New(rand.NewSource(time.Now().Unix()))
	rng.Seed(time.Now().Unix())
	return rng.Intn(max-min) + min
	// rand.Seed(time.Now().Unix())
	// return rand.Intn(max-min) + min
}
