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
    <div v-for="group in groups_list" :key="group.id">
      <hr>
      <router-link :to="{ name: 'group' }" @click="piniaManageDataGroup(group)">
        group id: {{ group.id }}
        <br> group name: {{ group.name }}
        <br> group description: {{ group.description }}
        <br> group created: {{ group.created_at }}
      </router-link>
      <router-link :to="{ name: 'target' }" @click="piniaManageDataProfile(group.email)">
        <br> group creator: {{ group.first_name }} {{ group.last_name }} ({{ group.email }})
      </router-link>
    </div>
  </div>

</template>

<script lang="ts" setup>
import { computed, onMounted } from 'vue';
import { useWebSocketStore } from '@/store/pinia';
import { useUUIDStore } from '@/store/pinia';
import { useProfileStore } from '@/store/pinia';
import { useGroupStore } from '@/store/group';
import { Group, TargetProfileRequest, WSMessageType } from '@/api/types';

const wss = useWebSocketStore();
const UUIDStore = useUUIDStore();
const profileStore = useProfileStore();

const groups_list = computed(() => wss.groupsList);
function updateGroupsList() {
  wss.sendMessage({
    type: WSMessageType.GROUPS_ALL_LIST,
    data: {
      user_uuid: UUIDStore.getUUID,
      target_email: profileStore.getUserEmail,
    } as TargetProfileRequest,
  })
}

const groupStore = useGroupStore();
function piniaManageDataGroup(group: Group) {
  groupStore.setGroup(group)
}

const piniaManageDataProfile = (email: string) => {
  profileStore.setTargetUserEmail(email);
};


onMounted(() => {
  updateGroupsList()
});

</script>