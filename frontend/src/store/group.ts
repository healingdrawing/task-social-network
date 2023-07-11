import { Group } from '@/api/types';
import { defineStore } from 'pinia';

interface GroupState {
  group: Group;
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
  }),
  getters: {
    getGroup(): Group { return this.group; },
  },
  actions: {
    setGroup(group: Group) { this.group = group; },
  },
});
