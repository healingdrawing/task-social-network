import { defineStore } from 'pinia';

interface GroupState {
  groupId: number;
}

export const useGroupStore = defineStore({
  id: 'group',
  state: ():GroupState => ({
    // Define your state properties here
    groupId: -1, 
  }),
  getters: {
    // Define your getters here
    getGroupId(): number { return this.groupId; },
  },
  actions: {
    // Define your actions here
    setGroupId(postId: number) { this.groupId = postId; },
  },
});
