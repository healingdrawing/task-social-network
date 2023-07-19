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
      user_id: localStorage.getItem("chat_target_user_id") !== null ? parseInt(localStorage.getItem("chat_target_user_id")!) : -1,
      email: localStorage.getItem("chat_target_user_email") || "",
      first_name: localStorage.getItem("chat_target_user_first_name") || "",
      last_name: localStorage.getItem("chat_target_user_last_name") || "",
    }
  }),
  getters: {
    // Define your getters here
    get_target_user(): UserForChatList { return this.target_user; },
  },
  actions: {
    // Define your actions here
    set_target_user(user: UserForChatList) {
      this.target_user = user;
      localStorage.setItem("chat_target_user_id", user.user_id.toString());
      localStorage.setItem("chat_target_user_email", user.email);
      localStorage.setItem("chat_target_user_first_name", user.first_name);
      localStorage.setItem("chat_target_user_last_name", user.last_name);
    },
  },
});
