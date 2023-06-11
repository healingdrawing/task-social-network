<template>
  <div class="nav-bar">
    <div class="nav-bar__logo">
      <img src="../assets/logo.png" alt="Vue logo" />
    </div>
    <div class="nav-bar__links">
      <router-link to="/profile">Profile</router-link> |
      <router-link to="/posts">Posts</router-link> |
      <router-link to="/chats">Chats</router-link> |
      <router-link to="/groups">Groups</router-link> |
      <router-link to="/">Logout</router-link>
    </div>
    <router-view/>
  </div>
</template>