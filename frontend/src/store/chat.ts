import { UserForChatList } from '@/api/types';
import { defineStore } from 'pinia';

interface ChatState {
  target_user: UserForChatList;
}

export const useChatStore = defineStore({
  id: 'chat',
  state: (): ChatState => ({
    // Define your state properties here
    target_user: {
      user_id: 0,
      email: '',
      first_name: '',
      last_name: '',
    }
  }),
  getters: {
    // Define your getters here
    get_target_user(): UserForChatList { return this.target_user; },
  },
  actions: {
    // Define your actions here
    set_target_user(user: UserForChatList) { this.target_user = user; },
  },
});
