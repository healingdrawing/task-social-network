import { defineStore } from 'pinia';

interface ChatState {
  chatId: number;
}

export const useChatStore = defineStore({
  id: 'chat',
  state: ():ChatState => ({
    // Define your state properties here
    chatId: -1, 
  }),
  getters: {
    // Define your getters here
    getChatId(): number { return this.chatId; },
  },
  actions: {
    // Define your actions here
    setChatId(postId: number) { this.chatId = postId; },
  },
});
