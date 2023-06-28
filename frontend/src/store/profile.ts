import { defineStore } from 'pinia';

interface ProfileState {
  userId: number;
  targetUserEmail: string; // used when the user visit other user's profile
}

export const useProfileStore = defineStore({
  id: 'profile',
  state: (): ProfileState => ({
    // Define your state properties here
    userId: -1,// the order of these properties is not matter. it is object
    targetUserEmail: "",
  }),
  getters: {
    // Define your getters here
    getUserId(): number { return this.userId; },
    getTargetUserEmail(): string { return this.targetUserEmail; },
  },
  actions: {
    // Define your actions here
    setUserId(userId: number) { this.userId = userId; },
    setTargetUserEmail(targetUserEmail: string) { this.targetUserEmail = targetUserEmail; },
  },
});
