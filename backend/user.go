package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
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
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Dob         string `json:"dob"`
	Avatar      string `json:"avatar"`
	avatarBytes []byte `sqlite3:"avatar"`
	Nickname    string `json:"nickname"`
	AboutMe     string `json:"about_me"`
	Public      bool   `json:"public"`
	Privacy     string `sqlite3:"privacy"`
}

type loginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UUIDData struct {
	UUID string `json:"UUID"`
}

func isFollowing(myID int, ID int) (bool, error) {
	rows, err := statements["doesSecondFollowFirst"].Query(ID, myID)
	if err != nil {
		log.Println("doesSecondFollowFirst query failed", err.Error())
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
		log.Println("dob is empty")
		jsonResponse(w, http.StatusUnprocessableEntity, "Invalid date of birth")
		return
	}

	// check the avatar validity
	if data.Avatar != "" {
		// cut prefix "data:image/jpeg;base64," or "data:image/png;base64,"
		// use comma as delimiter
		imageData, err := extractImageData(data.Avatar)
		if err != nil {
			log.Println("=FAIL extractImageData:", err.Error())
			jsonResponse(w, http.StatusUnprocessableEntity, err.Error()) //error here is handmade
			return
		}
		avatarData, err := base64.StdEncoding.DecodeString(imageData)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusUnprocessableEntity, "Invalid avatar")
			return
		}
		if !isImage(avatarData) {
			log.Println("avatar is not a valid image")
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

	data.Privacy = map[bool]string{true: "public", false: "private"}[data.Public]

	onlyEnglishRegex := regexp.MustCompile(`^[a-zA-Z0-9]{2,15}$`)

	if len([]rune(strings.TrimSpace(data.Nickname))) > 15 {
		data.Nickname =
			string([]rune(strings.TrimSpace(data.Nickname))[:15])
	}

	if len(data.FirstName) < 1 {
		message := `Invalid first name length: ` + data.FirstName + `
		First name must be between 1 and 32 characters long`
		log.Println(message)
		jsonResponse(w, http.StatusUnprocessableEntity, message)
		return
	}
	if len([]rune(strings.TrimSpace(data.FirstName))) > 32 {
		data.FirstName = string([]rune(strings.TrimSpace(data.FirstName))[:32])
	}

	if len(data.LastName) < 1 {
		message := `Invalid last name length: ` + data.LastName + `
		Last name must be between 1 and 32 characters long`
		log.Println(message)
		jsonResponse(w, http.StatusUnprocessableEntity, message)
		return
	}
	if len([]rune(strings.TrimSpace(data.LastName))) > 32 {
		data.LastName = string([]rune(strings.TrimSpace(data.LastName))[:32])
	}

	emailRegex := regexp.MustCompile(`^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,4})+$`)

	if !emailRegex.MatchString(strings.ToLower((data.Email))) {
		message := `Invalid email: ` + data.Email + `
		Email must be a valid email address`
		log.Println(message)
		jsonResponse(w, http.StatusUnprocessableEntity, message)
		return
	}

	if len(data.Password) < 6 || len(data.Password) > 15 {
		message := `Invalid password length: ` + data.Password + `
		Password must be between 6 and 15 characters long`
		log.Println(message)
		jsonResponse(w, http.StatusUnprocessableEntity, message)
		return
	}

	if !onlyEnglishRegex.MatchString(data.Password) {
		message := `Invalid password: ` + data.Password + `
Password must only contain english characters and numbers`
		log.Println(message)
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
		log.Println(err.Error())
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

	if rows.Next() {
		err = rows.Scan(&email, &hash)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, "scan credentials failed")
			return
		}
		rows.Close()
	}

	if email == "" || hash == "" {
		log.Println("Invalid credentials. Email or password is incorrect")
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

	http.SetCookie(w, &http.Cookie{
		Name:    "user_uuid",
		Value:   UUID,
		Expires: time.Now().Add(24 * time.Hour),
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "user_email",
		Value:   email,
		Expires: time.Now().Add(24 * time.Hour),
	})

	http.SetCookie(w, &http.Cookie{
		Value: UUID,
	})

	jsonResponseObj, _ := json.Marshal(map[string]string{
		"UUID":  UUID,
		"email": email, // todo: perhaps remove later
	})
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj)<-UUID,email failed")
		return
	}
}

func userLogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer recovery(w)

	cookie, err := r.Cookie("user_uuid")
	if err != nil || cookie.Value == "" || cookie == nil {
		log.Println(err.Error())
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
	ID, err := get_user_id_by_email(email)
	if err != nil {
		return "", err
	}
	_, err = statements["addSession"].Exec(UUID, ID)
	if err != nil {
		return "", err
	}
	return UUID, nil
}

// todo : still not refactored to ws. Because not used at the moment
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
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj)<-Exists:true failed")
		return
	}

}
