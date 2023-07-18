# task-social-network

After run project, open the browser on `http://localhost:3000/`  

## audit run to show containers  

To show two docker containers inside docker desktop or terminal, one for backend, and one for frontend (docker desktop must be started before), inside root folder of repository where `run.sh` placed, using terminal execute command:
- `./run.sh`  
to show containers in terminal:  
- `docker ps`  

## dev run  

### backend  

- `cd backend`
- `go run .`

### frontend  

- `cd frontend`
- `npm run serve`

## At the moment the recomendations are:  

- DO NOT USE `zero` branch, which is default, it can be used later as final destination  
- use your own branches , and merge them into `dev` branch, using pull requests  
- start branches from "wip/" prefix, f.e.: `wip/the-name-of-branch` or from the "futures/" prefix, f.e.: `futures/my-new-success`  
- before start work session always refresh you local branch from remote, to not forget something or do not recreate the same twice  
- after work session completed, even if the code not very polished, you should commit your changes into your branch to do not loose the local progress on the other computer(when you will open repo from other workplace), even if code is not very polished at the moment. It is your branch and you can use it as you want  

## task description and audit questions, on github  

https://github.com/01-edu/public/tree/master/subjects/social-network  

# According to task requirements, at the moment, there is no any restrictions for connection type, and level of privacy of the project.  

So using **websockets** for all interaction between client and server can be the best solution.

And websocket connection established only once, in time of user login success. Registration process before is required. After registration the client redirects user to login page, which is default page.

In task the require from us to use sessions and cookies, perhaps it can be limited by only login logout process, other interaction can be done by websockets, and in our case **pinia** storage, which is **Vue** js framework package.

Also database **migrations** are required on **golang based backend** side, and some allowed packages are provided for these needs.  
To implement some migrations more soft way it looks like not bad approach to **overwrite old tables structure from real-time-forum backend, and add to "Post/post" table new field(column) groupId(group_id), and overwrite the old records in table with value 0 or -1(depends on sqlite autoincrement minimum value) for the old posts not binded to group, using new migration**.
In the "ProfileView.vue" both type of posts (user created, and user created inside groups with membership) must be collected in one section, maybe on backend side, and sorted by date, descending order, the most fresh first.

New tables must be implemented for group functionality. Need discuss it.
F.e. it can be:
- "group" table: id, title, description, creator_id, created_at
- wip
