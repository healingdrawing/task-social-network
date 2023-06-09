<template>
  <nav>
    <router-link to="/">Login</router-link> |
    <router-link to="/signup">Signup</router-link>
  </nav>
  <router-view :key="$route.fullPath" /> <!--TODO:  not sure it needed, perplexity said this is to only show one view at one time, but without "key" it show the same -->
</template>