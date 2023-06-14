<template>
  <!-- todo: implement view "GroupsAllView.vue" -->
  <pre style="text-align: left;">
    browse all groups view gap
    Includes:
    - router-link -> "back to groups" - opens "GroupsView.vue"
    - list of all groups. Perhaps separated by x20 groups per page.
    Every item in list is clickable and opens "GroupView.vue".
    Every item in list includes group title, group members count.
  </pre>

  <div>
    <router-link to="/groups">Back to Groups</router-link>
    <div>
      <button @click="previousPage">Previous Page</button>
      <button @click="nextPage">Next Page</button>
    </div>
    <div v-for="group in displayedGroups" :key="group.id">
      <router-link
      :to="{ name: 'group' }"
      @click="piniaManageData(group)"
      >
        <h2>{{ group.title }}</h2>
      </router-link>
      <p>Members: {{ group.membersCount }}</p>
      <p>Creator: {{ group.creatorFullName }}</p>
    </div>
    <div>
      <button @click="previousPage">Previous Page</button>
      <button @click="nextPage">Next Page</button>
    </div>
  </div>

</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { useGroupStore } from '@/store/group';

interface Group {
  id: number;
  title: string;
  membersCount: number;
  creatorFullName: string;
}

const groups = ref<Group[]>([
  { id: 1, title: 'Group 1', membersCount: 10, creatorFullName: 'John Doe' },
  { id: 2, title: 'Group 2', membersCount: 20, creatorFullName: 'Jane Smith' },
  { id: 3, title: 'Group 3', membersCount: 30, creatorFullName: 'Bob Johnson' },
  { id: 4, title: 'Group 4', membersCount: 40, creatorFullName: 'Alice Brown' },
  { id: 5, title: 'Group 5', membersCount: 50, creatorFullName: 'Tom Wilson' },
  { id: 6, title: 'Group 6', membersCount: 60, creatorFullName: 'Sara Lee' },
  { id: 7, title: 'Group 7', membersCount: 70, creatorFullName: 'Mike Davis' },
  { id: 8, title: 'Group 8', membersCount: 80, creatorFullName: 'Lisa Green' },
  { id: 9, title: 'Group 9', membersCount: 90, creatorFullName: 'David Lee' },
  { id: 10, title: 'Group 10', membersCount: 100, creatorFullName: 'Amy Chen' },
]);

const itemsPerPage = 3;
const currentPage = ref(1);

function previousPage() {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
  updateDisplayedGroups()
}

function nextPage() {
  if (currentPage.value < Math.ceil(groups.value.length / itemsPerPage)) {
    currentPage.value++;
  }
  updateDisplayedGroups()
}

function getGroups() {
  const startIndex = (currentPage.value - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  return groups.value.slice(startIndex, endIndex);
}

const groupStore = useGroupStore();
function piniaManageData(group: Group) {
  groupStore.setGroupId(group.id)
}

function updateDisplayedGroups() {
  displayedGroups.value = getGroups()
}

const displayedGroups = ref<Group[]>([])
onMounted(() => {
  updateDisplayedGroups()
});

</script>