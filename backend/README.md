# DTO (Data Transfer Object) for  

1. frontend(TypeScript)  
2. JSON(Transfer)  
3. backend(Golang)  

- `ResponseError` used for error handling from backend to frontend.  
- Avatar image not managed properly in dummy code(it is just link to assets at the moment)
- `POST` requests should be used, because lot of `FUNCTIONALITY AVAILABLE ONLY FOR LOGGED IN USER`.  
Except `LoginView.vue` and `SignupView.vue`.
- `Èmail` is used as `ID` for user identification, because email is unique and used for login.
- `Nickname` is optional, so not used to identify user.

---

## LoginView.vue requests and responses

request is denoted by symbol ![request][request]
response is denoted by symbol ![response][response]

![request][request] Login

1. TypeScript

```typescript
interface Login {
  email: string;
  password: string;
}
```

2. JSON

```json
{
  "email": "string",
  "password": "string"
}
```

3. Golang

```go
type Login struct {
  Email    string `json:"email"`
  Password string `json:"password"`
}
```

![response][response]

  - - `FAIL` case

3. Golang

```go
type ResponseError struct {
  Type string `json:"type"`
  ErrorText string `json:"errorText"`
}
```

2. JSON

```json
{
  "type": "error",
  "errorText": "string"
}
```

1. TypeScript

```typescript
interface ResponseError {
  type: string;
  errorText: string;
}
```

  - - `SUCCESS` case

### To show profile data, redirect to ProfileView.vue

---

## ProfileView.vue  

### To show profile data, inside `onBeforeRouterEnter()` requests to backend, to fetch data before rendering page  

- ![request][request] (to get profile data) endpoint: `/api/user/profile`  

- ![response][response]

3. Golang

```go
type Profile struct {
  Email     string `json:"email"`
  FirstName string `json:"firstName"`
  LastName  string `json:"lastName"`
  Dob       string `json:"dob"`
  Avatar    string `json:"avatar"`
  Nickname  string `json:"nickname"`
  AboutMe   string `json:"aboutMe"`
  Public    bool   `json:"public"`
}
```

2. JSON

```json
{
  "email": "string",
  "firstName": "string",
  "lastName": "string",
  "dob": "string",
  "avatar": "string",
  "nickname": "string",
  "aboutMe": "string",
  "public": false
}
```

1. TypeScript

```typescript
interface Profile {
  email: string;
  firstName: string;
  lastName: string;
  dob: string;
  avatar: string;
  nickname: string;
  aboutMe: string;
  public: boolean;
}
```

- ![request][request] (following users list) endpoint: `/api/user/following`

- ![response][response] from backend

3. Golang

```go
type User struct {
  Email    string `json:"email"`
  FullName string `json:"fullName"`
}
type UserList struct {
  Users []User `json:"users"`
}
```

2. JSON

```json
{
  "users": [
    {
      "email": "string",
      "fullName": "string"
    }
  ]
}
```

1. TypeScript

```typescript
interface User {
  email: string;
  fullName: string;
}
interface UserList {
  users: User[];
}
```

- ### `REQUEST` (followers users list) endpoint: `/api/user/followers`

- ### `RESPONSE` from backend

SAME STRUCTURE AS FOR `/api/user/following`

- ### `REQUEST` (user's posts list) endpoint: `/api/user/posts`

- ### `RESPONSE` from backend

3. Golang

```go
type Post struct {
  Id        int    `json:"id"`
  Title     string `json:"title"`
  Tags      string `json:"tags"`
  CreatedAt string `json:"createdAt"`
  CreatorFullName string `json:"creatorFullName"`
  CreatorEmail string `json:"creatorEmail"`
}
type PostList struct {
  Posts []Post `json:"posts"`
}
```

2. JSON

```json
{
  "posts": [
    {
      "id": 0,
      "title": "string",
      "tags": "string",
      "createdAt": "string",
      "creatorFullName": "string",
      "creatorEmail": "string"
    }
  ]
}
```

1. TypeScript

```typescript
interface Post {
  id: number;
  title: string;
  tags: string;
  createdAt: string;
  creatorFullName: string;
  creatorEmail: string;
}
interface PostList {
  posts: Post[];
}
```

## User action requests and responses  

- ### `REQUEST` (change profile privacy) engpoint: `/api/user/privacy`

- ### `RESPONSE` from backend
    - - `SUCCESS` case  

3. Golang

```go
type Privacy struct {
  Public bool `json:"public"`
}
```

2. JSON

```json
{
  "public": false
}
```

1. TypeScript

```typescript
interface Privacy {
  public: boolean;
}
```

---

## TargetView.vue  

### To show target user profile data, inside `onBeforeRouterEnter()` requests to backend, to fetch data before rendering page

- ### `REQUEST` (target user profile check request status) endpoint: `/api/user/profile/request/{email}`

- ### `RESPONSE` from backend

    - - logged in user is `NOT FOLLOWER` of target user, and `PROFILE IS PRIVATE` and `REQUEST WAS NOT MADE` case
Show the `Request To Follow` button.

3. Golang

```go
type IsVisitorNotFollowerAndDidNotRequested struct {
  IsVisitorNotFollowerAndDidNotRequested bool `json:"isVisitorNotFollowerAndDidNotRequested"`
}
```

2. JSON

```json
{
  "isVisitorNotFollowerAndDidNotRequested": true
}
```

1. TypeScript

```typescript
interface IsVisitorNotFollowerAndDidNotRequested {
  isVisitorNotFollowerAndDidNotRequested: boolean;
}
```

isProfilePublicOrVisitorFollower

- SUCCESS case

SAME STRUCTURE AS FOR `/api/user/profile`


![request][request] (target user following users list) endpoint: `/api/user/following/{email}`  
  - logged in user is `NOT FOLLOWER` of target user, and `PROFILE IS PRIVATE` case  
Redirect to `ProfileView.vue`. This is the case impossible using normal navigation.
  - `SUCCESS` case
SAME STRUCTURE AS FOR `/api/user/following`
Also next endpoints responses structured the same way as for logged in user profile.

![request][request] (target user followers users list) endpoint: `/api/user/followers/{email}`  

- `REQUEST` (target user posts list) endpoint: `/api/user/posts/{email}`  

## User action follow requests and responses

![request][request] (follow target user) endpoint: `/api/user/follow/{email}`

![response][response] `SUCCESS` case :

3. Golang

```go
type IsVisitorNotFollowerAndDidNotRequested struct {
  IsVisitorNotFollowerAndDidNotRequested bool `json:"isVisitorNotFollowerAndDidNotRequested"`
}
```

2. JSON

```json
{
  "isVisitorNotFollowerAndDidNotRequested": false
}
```

1. TypeScript

```typescript
interface IsVisitorNotFollowerAndDidNotRequested {
  isVisitorNotFollowerAndDidNotRequested: boolean;
}
```

---

## SignupView.vue  

- (signup new user) endpoint: `/api/signup`

![request][request]

1. TypeScript

```typescript
interface User {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  dob: string;
  avatar: File | null;
  nickname: string;
  aboutMe: string;
  public: boolean;
}
```

2. JSON

```json
{
  "email": "string",
  "password": "string",
  "firstName": "string",
  "lastName": "string",
  "dob": "string",
  "avatar": "string",
  "nickname": "string",
  "aboutMe": "string",
  "public": false
}
```

3. Golang

```go
type User struct {
  Email    string `json:"email"`
  Password string `json:"password"`
  FirstName string `json:"firstName"`
  LastName string `json:"lastName"`
  Dob string `json:"dob"`
  Avatar string `json:"avatar"`
  Nickname string `json:"nickname"`
  AboutMe string `json:"aboutMe"`
  Public bool `json:"public"`
}
```

![response][response]

- `SUCCESS` case:
Redirect to `LoginView.vue`. Perhaps will be better to sign in user automatically and redirect to `ProfileView.vue`.

---

# TODO: sergei see this and make them alright, make them go at the right place int he readme. Below this line is rought work

// Incoming JSON DTO for group creation over handler groupCreateHandler

```json
{
  "UUID": 1,
  "name": "group name",
  "description": "group description",
  "privacy": "public"
}
```

// Incoming JSON DTO for group joining over handler groupJoinHandler

```json
{
  "group_id": 1,
  "member_id": 1
}
```

[request]: https://i.postimg.cc/C5c0SmRN/requestww.jpg "Request"
[response]: https://i.postimg.cc/RhVRwFqd/responsew.jpg "Response"
