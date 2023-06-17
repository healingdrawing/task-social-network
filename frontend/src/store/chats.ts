import { defineStore } from 'pinia';

export interface Message {
  message: string;
  groupId: number;
  userId: number;
}

interface ChatsState {
  chatMessages: Map<number, string[]>;
}

export const useChatsStore = defineStore({
  id: 'chats',
  state: (): ChatsState => ({
    // Define your state properties here
    // todo: remove this dummy data later
    chatMessages: new Map<number, string[]>([
      [1, ["apple", "banana", "orange"]],
      [2, ["cat", "dog"]],
    ]),
    // chatMessages: new Map<number, string[]>(), // todo: use this in production
  }),
  getters: {
    // Define your getters here
    // todo: weird... cant refactor this function without arrow function notation, it always shows some shit. Need sleep session
    getNewMessages: (state) => (chatId: number): string[] => { return state.chatMessages.get(chatId) || [] },
    // this method use to highlight NavBar.vue section "Chats" , if at least one fresh(not displayed on screen once) is present
    // in messages storage , which is map of chats by chatId with array of messages for each chat.
    // Plan to fill this storage using new messages on background, and once window of chat opened,
    // messages from this storage for opened chat must be removed.
    hasNewMessages: (state) => [...state.chatMessages.values()].some(array => array.length > 0)
  },
  actions: {
    // Define your actions here
    setNewMessages(chatId: number, messages: string[]) { this.chatMessages.set(chatId, messages) },
  },
});
