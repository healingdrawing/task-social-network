# task-social-network

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