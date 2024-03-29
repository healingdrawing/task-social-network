import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
// import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'login',
    component: LoginView
  },
  {
    path: '/signup',
    name: 'signup',
    /*todo: comment bottom created by vue create myappname command. lazy loading*/
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "signup" */ '../views/SignupView.vue')
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import(/* webpackChunkName: "profile" */ '../views/ProfileView.vue')
  },
  {
    path: '/target',
    name: 'target',
    component: () => import(/* webpackChunkName: "target" */ '../views/TargetView.vue')
  },
  {
    path: '/post',
    name: 'post',
    component: () => import(/* webpackChunkName: "post" */ '../views/PostView.vue')
  },
  {
    path: '/posts',
    name: 'posts',
    component: () => import(/* webpackChunkName: "posts" */ '../views/PostsView.vue')
  },
  {
    path: '/group_post',
    name: 'group_post',
    component: () => import(/* webpackChunkName: "group_post" */ '../views/GroupPostView.vue')
  },
  {
    path: '/group_posts',
    name: 'group_posts',
    component: () => import(/* webpackChunkName: "group_posts" */ '../views/GroupPostsView.vue')
  },
  {
    path: '/group',
    name: 'group',
    component: () => import(/* webpackChunkName: "group" */ '../views/GroupView.vue')
  },
  {
    path: '/group_invite',
    name: 'group_invite',
    component: () => import(/* webpackChunkName: "group_invite" */ '../views/GroupInviteView.vue')
  },
  {
    path: '/groups',
    name: 'groups',
    component: () => import(/* webpackChunkName: "groups" */ '../views/GroupsView.vue')
  },
  {
    path: '/groups_all',
    name: 'groups_all',
    component: () => import(/* webpackChunkName: "groups_all" */ '../views/GroupsAllView.vue')
  },
  {
    path: '/group_chat',
    name: 'group_chat',
    component: () => import(/* webpackChunkName: "group_chat" */ '../views/GroupChatView.vue')
  },
  {
    path: '/chat',
    name: 'chat',
    component: () => import(/* webpackChunkName: "chat" */ '../views/ChatView.vue')
  },
  {
    path: '/chats',
    name: 'chats',
    component: () => import(/* webpackChunkName: "chats" */ '../views/ChatsView.vue')
  },
  {
    path: '/bell',
    name: 'bell',
    component: () => import(/* webpackChunkName: "bell" */ '../views/BellView.vue')
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior(/*to, from, savedPosition*/) { //todo: BE CAREFUL. scrollBehavior(to, from, savedPosition) { - raise warning(NOT DEADLY). To mute this shit of vscode, incoming parameters was removed/commented. This can produce weird behavior if you will try to use router with scrollBehavior later (according to perplexity message).
    return { top: 0 };//scroll to top in time of navigation between routes
  },
})

export default router
