import { Group } from '@/api/types';
import { defineStore } from 'pinia';

interface GroupState {
  group: Group;
  groupId: number;
}

export const useGroupStore = defineStore({
  id: 'group',
  state: (): GroupState => ({
    group: {
      id: -1,
      name: '',
      description: '',
      created_at: '',
      email: '',
      first_name: '',
      last_name: '',
    },
    // Define your state properties here
    groupId: -1,
  }),
  getters: {
    // Define your getters here
    getGroup(): Group { return this.group; },
    getGroupId(): number { return this.groupId; },
  },
  actions: {
    // Define your actions here
    setGroup(group: Group) { this.group = group; },
    setGroupId(postId: number) { this.groupId = postId; },
  },
});
