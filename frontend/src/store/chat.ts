import { defineStore } from 'pinia';

interface ChatState {
  target_user_id: number;
}

export const useChatStore = defineStore({
  id: 'chat',
  state: (): ChatState => ({
    // Define your state properties here
    target_user_id: 0,
  }),
  getters: {
    // Define your getters here
    get_target_user_id(): number { return this.target_user_id; },
  },
  actions: {
    // Define your actions here
    set_target_user_id(user_id: number) { this.target_user_id = user_id; },
  },
});
