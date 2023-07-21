import { Group } from '@/api/types';
import { defineStore } from 'pinia';

interface GroupState {
  group: Group;
}

export const useGroupStore = defineStore({
  id: 'group',
  state: (): GroupState => ({
    group: {
      id: sessionStorage.getItem("group_id") !== null ? parseInt(sessionStorage.getItem("group_id")!) : -1,
      name: sessionStorage.getItem("group_name") || "",
      description: sessionStorage.getItem("group_description") || "",
      created_at: sessionStorage.getItem("group_created_at") || "",
      email: sessionStorage.getItem("group_email") || "",
      first_name: sessionStorage.getItem("group_first_name") || "",
      last_name: sessionStorage.getItem("group_last_name") || "",
    },
  }),
  getters: {
    getGroup(): Group { return this.group; },
  },
  actions: {
    setGroup(group: Group) {
      this.group = group;
      sessionStorage.setItem("group_id", group.id.toString());
      sessionStorage.setItem("group_name", group.name);
      sessionStorage.setItem("group_description", group.description);
      sessionStorage.setItem("group_created_at", group.created_at);
      sessionStorage.setItem("group_email", group.email);
      sessionStorage.setItem("group_first_name", group.first_name);
      sessionStorage.setItem("group_last_name", group.last_name);
    },
  },
});
