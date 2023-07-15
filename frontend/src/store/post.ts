import { GroupPost, Post } from '@/api/types';
import { defineStore } from 'pinia';

interface PostState {
  post: Post;
  group_post: GroupPost;
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
    group_post: {
      group_id: -1,
      group_name: '',
      group_description: '',
      id: -1,
      title: '',
      categories: '',
      content: '',
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
    getGroupPost(): GroupPost { return this.group_post; },
  },
  actions: {
    // Define your actions here
    setPost(post: Post) { this.post = post; },
    setGroupPost(group_post: GroupPost) { this.group_post = group_post; },
  },
});
