import { defineStore } from 'pinia';

export const useProfileStore = defineStore({
  id: 'profile',
  state: () => ({
    // Define your state properties here
    postId: -1
  }),
  getters: {
    // Define your getters here
    getPostId(): number {
      return this.postId;
    }
  },
  actions: {
    // Define your actions here
    setPostId(postId: number) {
      this.postId = postId;
    }
  },
});
