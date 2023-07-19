<template>
  <nav>
    <router-link to="/">
      <div class="router_link_box">
        Login
      </div>
    </router-link> |
    <router-link to="/signup">
      <div class="router_link_box">
        Signup
      </div>
    </router-link>
  </nav>
  <router-view :key="$route.fullPath" /> <!--TODO:  not sure it needed, perplexity said this is to only show one view at one time, but without "key" it show the same -->
</template>