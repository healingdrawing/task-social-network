import { Group } from '@/api/types';
import { defineStore } from 'pinia';

interface GroupState {
  group: Group;
}

export const useGroupStore = defineStore({
  id: 'group',
  state: (): GroupState => ({
    group: {
      id: localStorage.getItem("group_id") !== null ? parseInt(localStorage.getItem("group_id")!) : -1,
      name: localStorage.getItem("group_name") || "",
      description: localStorage.getItem("group_description") || "",
      created_at: localStorage.getItem("group_created_at") || "",
      email: localStorage.getItem("group_email") || "",
      first_name: localStorage.getItem("group_first_name") || "",
      last_name: localStorage.getItem("group_last_name") || "",
    },
  }),
  getters: {
    getGroup(): Group { return this.group; },
  },
  actions: {
    setGroup(group: Group) {
      this.group = group;
      localStorage.setItem("group_id", group.id.toString());
      localStorage.setItem("group_name", group.name);
      localStorage.setItem("group_description", group.description);
      localStorage.setItem("group_created_at", group.created_at);
      localStorage.setItem("group_email", group.email);
      localStorage.setItem("group_first_name", group.first_name);
      localStorage.setItem("group_last_name", group.last_name);
    },
  },
});
