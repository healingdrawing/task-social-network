# DTO (Data Transfer Object) for  

[request]: https://i.postimg.cc/C5c0SmRN/requestww.jpg "Request"
[response]: https://i.postimg.cc/RhVRwFqd/responsew.jpg "Response"

request is denoted by symbol ![request][request]  
response is denoted by symbol ![response][response]

<hr style="border:2px solid green">

- `ResponseError` used for error handling from backend to frontend.  
- Avatar image not managed properly in dummy code(it is just link to assets at the moment)
- `FUNCTIONALITY AVAILABLE ONLY FOR LOGGED IN USER`.  
Except `LoginView.vue` and `SignupView.vue`.
- `Èmail` is used as `ID` for user identification, because email is unique and used for login.

<hr style="border:4px solid green">

## LoginView.vue

<hr style="border:2px solid green">

![request][request] `/api/user/register`

```json
{
  "email": "string",
  "firstName": "string",
  "lastName": "string",
  "password": "string",
  "dob": "string",
  "avatar": "base64encodeBlobTo - string",
  "nickname": "string",
  "aboutMe": "string",
  "public": boolean
}
```

![response][response]
`SUCCESS`

```json
{
    "UUID": "string",
    "email": "string"
}
```


<hr style="border:2px solid green">

![request][request] `/api/login`

```json
{
  "email": "string",
  "password": "string"
}
```

![response][response]
`FAIL`

```json
{
  "type": "error",
  "errorText": "string"
}
```

`SUCCESS` - to show profile data, redirect to ProfileView.vue

<hr style="border:4px solid green">  

## ProfileView.vue  
To show profile data, inside `onBeforeRouterEnter()` requests to backend, to fetch data before rendering page  

<hr style="border:2px solid green">

![request][request] `/api/user/profile` (to get profile data)

json body only needs email property, because email is used as ID to identify user. If you want to check your own profile, even email property is not needed, because it is taken from token.

```json
{
  "email": "string"
}
```

![response][response]  

```json
{
  "email": "string",
  "firstName": "string",
  "lastName": "string",
  "dob": "string",
  "avatar": "string",
  "nickname": "string",
  "aboutMe": "string",
  "privacy": boolean
}
```

<hr style="border:2px solid green">

![request][request]  `api/followrequestlist` (follow request list, all pending follow requests to the current logged in user)

No need to send any data in request, because user_id is taken from token.

![response][response]  

```json
[
  {
    "email": "string",
    "fullName": "string"
  }
]
```

<hr style="border:2px solid green">

![request][request] `/api/user/follow` (follow this target user):

```json
{
  "email": "string"
}
```

![response][response]
`SUCCESS`  

if the target user was private:

```json
{
  "message": "request sent to follow the user",
}
```

if the target user was public:

```json
{
  "message": "user followed",
}
```

<hr style="border:2px solid green">

![request][request] `/api/user/unfollow` (unfollow this target user):

  ```json
  {
    "email": "string"
  }
  ```

![response][response]
`SUCCESS`  

```json
{
  "message": "user unfollowed",
}
```

<hr style="border:2px solid green">

![request][request] `/api/user/following`  (following users list) :  

> [x] We need to check the following list of other users also, whose profile we visit. So we need to make request to backend to get the list of following users.

- If request has no data, returns list of following users of logged in user.

- If request has JSON email property, then we can get the list of following users of user to whom the email belongs to.

connects to handler `FollowingHandler`

```json
{
  "email": "string",
}
```

![response][response]

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

<hr style="border:2px solid green">

![request][request] `/api/user/followers` (followers users list)

![response][response]  

SAME STRUCTURE AS FOR `/api/user/following`

<hr style="border:2px solid green">

![request][request] `/api/followrequest/reject` (reject folllow request)

```json
{
  "email": "string"
}
```

![response][response]  

```json
{
  "message": "you rejected the follow request",
}
```

<hr style="border:2px solid green">

![request][request] `/api/followrequest/accept` (accept folllow request)

```json
{
  "email": "string"
}
```

![response][response]  

```json
{
  "message": "you accepted the follow request",
}
```


<hr style="border:2px solid green">

![request][request] `/api/post/submit` (add new post)
  
  ```json
  {
    "title": "string",
    "content": "string",
    "categories": "string",
    "picture": "base64encodeBlobTo - string",
    "privacy": "private/public/almost private - string",
    "able_to_see": "string"
  }
  ```

![response][response]
`SUCCESS`  

```json
{
  "message": "Post created"
}
```

<hr style="border:2px solid green">

![request][request] `/api/user/posts` (user's posts list)

![response][response]

```json
{
  "posts": [
    {
      "id": 0,
      "title": "string",
      "content": "string",
      "categories": "string",
      "picture": "base64encodeBlobTo - string",
      "createdAt": "string",
      "creatorFullName": "string",
      "creatorEmail": "string"
    }
  ]
}
```

<hr style="border:2px solid green">

![request][request] "/api/comment/submit" (add new comment)

```json
{
  "postId": "number",
  "content": "string",
  "picture": "base64encodeBlobTo - string"
}
```

![response][response]
`SUCCESS`  

```json
{
  "message": "Comment created"
}
```

## User action requests and responses in ProfileView.vue

![request][request] `/api/user/privacy` (change profile privacy)

![response][response]
`SUCCESS`  

```json
{
  "public": false
}
```

<hr style="border:4px solid green">

![request][request] /api/group/submit (add new group)

```json
{
  "name": "string",
  "description": "string",
  "privacy": "public/private",
  "invited": "string"
}
```

![response][response]
`SUCCESS`  

```json
{
  "message": "Group created"
}
```

<hr style="border:4px solid green">

## TargetView.vue  
To show target user profile data, inside `onBeforeRouterEnter()` requests to backend, to fetch data before rendering page

<hr style="border:2px solid green">

![request][request] `/api/user/profile/request/{email}` (target user profile check request status, to manage button)

![response][response]
`IF` logged in user is NOT FOLLOWER of target user, and the target user PROFILE IS PRIVATE and following to target user REQUEST WAS NOT MADE `THEN` show the `Request To Follow` button.

```json
{
  "isVisitorNotFollowerAndDidNotRequested": true,
}
```

<hr style="border:2px solid green">

![request][request] `api/user/profile/public/{email}` (target user profile check following/public, to hide/show profile info).  
`IF` target profile is public or the visitor is follower `THEN` allow to show the profile section  
![response][response]

```json
{
  "isProfilePublicOrVisitorFollower": true
}
```

<hr style="border:2px solid green">

`IF` profile can be displayed `THEN` make `REQUEST/S` to get data for target profile view.  
`IF` logged in user is `NOT FOLLOWER` of target user, and `PROFILE IS PRIVATE` but requests still happen `THEN` 
Redirect to `ProfileView.vue`. This is the cases impossible using normal navigation(experimentators activity).

<hr style="border:2px solid green">

![request][request] `api/user/profile/{email}`
SAME JSON AS FOR `/api/user/profile`

<hr style="border:2px solid green">

![request][request] `/api/user/following/{email}` (target user following users list)  
![response][response] SAME STRUCTURE AS FOR `/api/user/following`  

Also next two endpoints responses structured the same way as for logged in user own profile.

![request][request] `/api/user/followers/{email}` (target user followers users list)  

![request][request]  `/api/user/posts/{email}` (target user posts list) endpoint:

## User action follow requests and responses

![request][request] `/api/user/follow/{email}` (follow target user)

![response][response]

```json
{
  "isVisitorNotFollowerAndDidNotRequested": false
}
```

<hr style="border:4px solid green">

## SignupView.vue  

<hr style="border:2px solid green">

![request][request] `/api/signup` (signup new user)

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

![response][response]
Redirect to `LoginView.vue`. Perhaps will be better to sign in user automatically and redirect to `ProfileView.vue`.

<hr style="border:2px solid green">

## ERD for database

![erd](https://github.com/healingdrawing/task-social-network/assets/5121817/4b8354f7-c165-4e31-aa63-de0f6fe7a89b)

## To create base64 encoded string from image file, required for testing

I recommend this website

<https://elmah.io/tools/base64-image-encoder>


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
