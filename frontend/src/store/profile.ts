import { defineStore } from 'pinia';

interface ProfileState {
  user_id: number;
  user_email: string;
  target_user_email: string; // used when the user visit other user's profile
}

export const useProfileStore = defineStore({
  id: 'profile',
  state: (): ProfileState => ({
    // Define your state properties here
    user_id: localStorage.getItem("user_id") !== null ? parseInt(localStorage.getItem("user_id")!) : -1
    ,// the order of these properties is not matter. it is object
    user_email: localStorage.getItem("user_email") || "",
    target_user_email: localStorage.getItem("target_user_email") || "",
  }),
  getters: {
    // Define your getters here
    getUserId(): number { return this.user_id; },
    getUserEmail(): string { return this.user_email; },
    getTargetUserEmail(): string { return this.target_user_email; },
  },
  actions: {
    // Define your actions here
    setUserId(user_id: number) {
      this.user_id = user_id;
      localStorage.setItem("user_id", user_id.toString());
    },
    setUserEmail(user_email: string) {
      this.user_email = user_email;
      localStorage.setItem("user_email", user_email);
    },
    setTargetUserEmail(target_user_email: string) {
      this.target_user_email = target_user_email;
      localStorage.setItem("target_user_email", target_user_email);
    },
  },
});
