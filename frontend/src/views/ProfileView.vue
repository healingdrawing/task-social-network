<template>
  <h1>Profile</h1>
  <!-- add checkbox to make profile public -->
  <div>
    <label>
      <input type="checkbox" v-model="isPublic" />
      is public
    </label>
  </div>
  <!-- add user information -->
  <div>
    <p>Email: {{ email }}</p>
    <p>First Name: {{ firstName }}</p>
    <p>Last Name: {{ lastName }}</p>
    <p>Date of Birth: {{ dob }}</p>
    <p>Nickname: {{ nickname }}</p>
    <p>About Me: {{ aboutMe }}</p>
  </div>
  <!-- separately add avatar, perhaps it should be on the right half of screen -->
  <div>
    <p>Avatar: <img :src="getImgUrl(avatar)" alt="fail again"></p>
  </div>
  <!-- add following list. The other users followed by the user -->
  <h2>Following:</h2>
  <div class="user-list" style="height: 100px; overflow-y: scroll;">
    <div v-for="user in followingList" :key="user.id">{{ user.name }}</div>
  </div>
  <!-- add followers list. The other users following the user -->
  <h2>Followers:</h2>
  <div class="user-list" style="height: 100px; overflow-y: scroll;">
    <div v-for="user in followersList" :key="user.id">{{ user.name }}</div>
  </div>
  <!-- add user posts list. The posts created by the user -->
  <h2>Posts:</h2>
  <div v-for="post in postsList"
    :key="post.id">
    <router-link
    :to="{ name: 'post' }"
    @click="updateStorageForPostView(post)">
      {{ post.title }}
    </router-link>
  </div>
  <!-- ( :to="{ name: 'post' }" ) also can be ( :to="'/post'" ) -->
</template>

<script lang="ts" setup>
import { ref, watch, onMounted } from 'vue';
import { useProfileStore } from '@/store/profile';

const isPublic = ref(false);

watch(isPublic, (newValue, oldValue) => {
  handleCheckboxChange(newValue);
});

function handleCheckboxChange(value: boolean) {
  // Call your method here
  alert(value)
}

function getImgUrl(imageNameWithExtension: string) {
  return require(`../assets/${imageNameWithExtension}`)
}

const email = ref('john.doe@example.com');
const firstName = ref('John');
const lastName = ref('Doe');
const dob = ref('01/01/1990');
const avatar = ref('facepalm, need black magic to load image, wtf');
const nickname = ref('johndoe');
const aboutMe = ref('Lorem ipsum dolor sit amet, consectetur adipiscing elit.');

//todo: remove/refactor later
/** Function to update the profile data using dummy data at the moment*/
function updateProfile() {
  email.value = 'jane.doe@example.com';
  firstName.value = 'Jane';
  lastName.value = 'Doe';
  dob.value = '02/02/1995';
  avatar.value = 'logo.png'; /* todo: fix this later!!!. the worst part of placeholder. it must be image(uploaded or anonymous placeholder), so or from assets, like it is now, or from public folder(requires another code) */
  nickname.value = 'janedoe';
  aboutMe.value = 'Very interesting text.';
}

updateProfile();

// following and followers section
interface User {
  id: number;
  name: string;
}

const followingList = ref<User[]>([]);

// todo: dummy data, remove/refactor later
function updateFollowingList() {
  // Code to get the user list goes here
  const users: User[] = [
    { id: 1, name: 'John' },
    { id: 2, name: 'Jane' },
    { id: 3, name: 'Bob' },
    { id: 4, name: 'Alice' },
    { id: 5, name: 'Mike' },
    { id: 6, name: 'Sara' },
    { id: 7, name: 'Tom' },
    { id: 8, name: 'Kate' },
    { id: 9, name: 'David' },
    { id: 10, name: 'Emily' },
  ];
  followingList.value = users;
  console.log('Following list updated');
}

const followersList = ref<User[]>([]);

// todo: dummy data, remove/refactor later
function updateFollowersList() {
  // Code to get the user list goes here
  const users: User[] = [
    { id: 1, name: 'John follower' },
    { id: 2, name: 'Jane follower' },
    { id: 3, name: 'Bob follower' },
    { id: 4, name: 'Alice follower' },
    { id: 5, name: 'Mike follower' },
    { id: 6, name: 'Sara follower' },
    { id: 7, name: 'Tom follower' },
    { id: 8, name: 'Kate follower' },
    { id: 9, name: 'David follower' },
    { id: 10, name: 'Emily follower' },
  ];
  followersList.value = users;
  console.log('Followers list updated');
}

// user posts section
interface Post {
  id: number;
  title: string;
}

const postsList = ref<Post[]>([]);

// todo: dummy data, remove/refactor later
function updatePostList() {
  // Code to get the user posts goes here
  const posts: Post[] = [
    { id: 1, title: 'Dummy Post 1 Title' },
    { id: 2, title: 'Dummy Post 2 Title' },
    { id: 3, title: 'Dummy Post 3 Title' },
    { id: 4, title: 'Dummy Post 4 Title' },
    { id: 5, title: 'Dummy Post 5 Title' },
    { id: 6, title: 'Dummy Post 6 Title' },
    { id: 7, title: 'Dummy Post 7 Title' },
    { id: 8, title: 'Dummy Post 8 Title' },
    { id: 9, title: 'Dummy Post 9 Title' },
    { id: 10, title: 'Dummy Post 10 Title' },
  ];
  postsList.value = posts;
  console.log('Posts list updated');
}


const profileStore = useProfileStore()
const updateStorageForPostView = (post: Post) => {
  profileStore.setPostId(post.id)
  // router.push('/post') //todo: remove later, no need because of router-link
}

onMounted(() => {
  updateFollowingList();
  updateFollowersList();
  updatePostList();
});

</script>