import { Post } from '@/api/types';
import { defineStore } from 'pinia';

interface PostState {
  post: Post;
}

export const usePostStore = defineStore({
  id: 'post',
  state: (): PostState => ({
    // Define your state properties here
    post: {
      id: -1,
      title: '',
      categories: '',
      content: '',
      privacy: '',
      picture: '',
      created_at: '',
      email: '',
      first_name: '',
      last_name: '',
    },
  }),
  getters: {
    // Define your getters here
    getPost(): Post { return this.post; },
  },
  actions: {
    // Define your actions here
    setPost(post: Post) { this.post = post; },
  },
});
