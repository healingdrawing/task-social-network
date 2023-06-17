import { defineStore } from 'pinia';

interface ProfileState {
  userId: number;
  targetUserId: number; // used when the user visit other user's profile
}

export const useProfileStore = defineStore({
  id: 'profile',
  state: ():ProfileState => ({
    // Define your state properties here
    userId: -1,// the order of these properties is not matter. it is object
    targetUserId: -1,
  }),
  getters: {
    // Define your getters here
    getUserId(): number { return this.userId; },
    getTargetUserId(): number { return this.targetUserId; },
  },
  actions: {
    // Define your actions here
    setUserId(userId: number) { this.userId = userId; },
    setTargetUserId(targetUserId: number) { this.targetUserId = targetUserId; },
  },
});
