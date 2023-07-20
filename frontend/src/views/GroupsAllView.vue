<template>
  <div>
    <br>
    <router-link to="/groups">
      <div class="router_link_box">
        Back to Groups
      </div>
    </router-link>
    <div v-for="group in groups_list" :key="group.id">
      <div class="single_div_box">
        <br>
        <h3> Group name: </h3> <p> {{ group.name }} </p>
        <h3> Group description: </h3> <p> {{ group.description }} </p>
        <h3> Group created: </h3> <p> {{ group.created_at }} </p>
        <router-link
          :Title = "group.first_name + '\n' + group.last_name + '\n' + group.email"
          :to="{ name: 'target' }" @click="piniaManageDataProfile(group.email)">
          <div class="router_link_box">
            visit creator profile
          </div>
        </router-link>
        <br>
        <router-link :to="{ name: 'group' }" @click="piniaManageDataGroup(group)">
          <div class="router_link_box">
            visit group {{ group.id }}
          </div>
        </router-link>
      </div>
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


onMounted(async () => {
  wss.refresh_websocket()
  await wss.waitForConnection();
  updateGroupsList()
});

</script>