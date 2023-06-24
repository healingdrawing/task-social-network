<template>
  <h1>Profile:</h1>

  <div>
    <button @click="handleFollowing()">
      {{ isVisitorNotFollowerAndDidNotRequested ? 'Request To Follow' : 'Unfollow' }}
    </button>
  </div>

  <div v-if="isProfilePublicOrVisitorFollower">
    <!-- todo: remove later . show id for dev needs-->
    <p>Id: {{ profileStore.getTargetUserId }}</p>
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
      @click="piniaManageData(post)">
        {{ post.title }}
      </router-link>
    </div>
    <!-- ( :to="{ name: 'post' }" ) also can be ( :to="'/post'" ) -->
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted } from 'vue';
import { usePostStore } from '@/store/post';
import { useProfileStore } from '@/store/profile';

// if true/false, then show follow/unfollow text on button
const isVisitorNotFollowerAndDidNotRequested = ref(true);

watch(isVisitorNotFollowerAndDidNotRequested, (newValue, oldValue) => {
  alert(`isVisitorNotFollowerAndDidNotRequested: ${newValue}`);
  // handleFollowing(newValue);
});

function handleFollowing() {
  // Call your method here
  isVisitorNotFollowerAndDidNotRequested.value = !isVisitorNotFollowerAndDidNotRequested.value;
}

// if true then show all profile information on screen
const isProfilePublicOrVisitorFollower = ref(true);

function getImgUrl(imageNameWithExtension: string) {
  return require(`../assets/${imageNameWithExtension}`)
}

const email = ref('john.doe@example.com');
const firstName = ref('John');
const lastName = ref('Doe');
const dob = ref('01/01/1990');
const avatar = ref('logo.png'); //todo: have to be at least logo.png in intial moment or fail with updateProfile() inside onMounted(). Outside onMounted() it works fine, with any initial string value
const nickname = ref('johndoe');
const aboutMe = ref('Lorem ipsum dolor sit amet, consectetur adipiscing elit.');

//todo: remove/refactor later, to get data from backend
interface Profile {
  email: string;
  firstName: string;
  lastName: string;
  dob: string;
  avatar: string;
  nickname: string;
  aboutMe: string;
}

function fillDummyProfileData(userId: number): Profile {
  return {
    email: 'dummy@mail.com',
    firstName: 'Dummy',
    lastName: 'Doe',
    dob: '01/01/1990',
    avatar: 'logo.png',
    nickname: 'dummydoe',
    aboutMe: 'Very interesting text.',
  };
}

const profileStore = useProfileStore();
/** Function to update the profile data using dummy data at the moment*/
function updateProfile() {
  //todo: get data from backend based on user() id from profileStore
  const profile = fillDummyProfileData(profileStore.getTargetUserId);

  email.value = profile.email;
  firstName.value = profile.firstName;
  lastName.value = profile.lastName;
  dob.value = profile.dob;
  avatar.value = profile.avatar;//todo: placeholder.it must be image(uploaded or anonymous placeholder), so or from assets, like it is now, or from public folder(requires another code) * /
  nickname.value = profile.nickname;
  aboutMe.value = profile.aboutMe;
}

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
function updatePostsList() {
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


const postStore = usePostStore()
const piniaManageData = (post: Post) => {
  postStore.setPostId(post.id)
}

onMounted(() => {
  updateProfile();
  updateFollowingList();
  updateFollowersList();
  updatePostsList();
});

</script>