import { defineStore } from 'pinia';

interface PostState {
  postId: number;
}

export const usePostStore = defineStore({
  id: 'post',
  state: (): PostState => ({
    // Define your state properties here
    postId: -1,
  }),
  getters: {
    // Define your getters here
    getPostId(): number { return this.postId; },
  },
  actions: {
    // Define your actions here
    setPostId(postId: number) { this.postId = postId; },
  },
});
