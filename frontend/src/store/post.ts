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
      id: localStorage.getItem("post_id") !== null ? parseInt(localStorage.getItem("post_id")!) : -1,
      title: localStorage.getItem("post_title") || "",
      categories: localStorage.getItem("post_categories") || "",
      content: localStorage.getItem("post_content") || "",
      privacy: localStorage.getItem("post_privacy") || "",
      picture: localStorage.getItem("post_picture") || "",
      created_at: localStorage.getItem("post_created_at") || "",
      email: localStorage.getItem("post_email") || "",
      first_name: localStorage.getItem("post_first_name") || "",
      last_name: localStorage.getItem("post_last_name") || "",
    },
    group_post: {
      group_id: localStorage.getItem("group_post_group_id") !== null ? parseInt(localStorage.getItem("group_post_group_id")!) : -1,
      group_name: localStorage.getItem("group_post_group_name") || "",
      group_description: localStorage.getItem("group_post_group_description") || "",
      id: localStorage.getItem("group_post_id") !== null ? parseInt(localStorage.getItem("group_post_id")!) : -1,
      title: localStorage.getItem("group_post_title") || "",
      categories: localStorage.getItem("group_post_categories") || "",
      content: localStorage.getItem("group_post_content") || "",
      picture: localStorage.getItem("group_post_picture") || "",
      created_at: localStorage.getItem("group_post_created_at") || "",
      email: localStorage.getItem("group_post_email") || "",
      first_name: localStorage.getItem("group_post_first_name") || "",
      last_name: localStorage.getItem("group_post_last_name") || "",
    },
  }),
  getters: {
    // Define your getters here
    getPost(): Post { return this.post; },
    getGroupPost(): GroupPost { return this.group_post; },
  },
  actions: {
    // Define your actions here
    setPost(post: Post) {
      this.post = post;
      localStorage.setItem("post_id", post.id.toString());
      localStorage.setItem("post_title", post.title);
      localStorage.setItem("post_categories", post.categories);
      localStorage.setItem("post_content", post.content);
      localStorage.setItem("post_privacy", post.privacy);
      localStorage.setItem("post_picture", post.picture || "");
      localStorage.setItem("post_created_at", post.created_at);
      localStorage.setItem("post_email", post.email);
      localStorage.setItem("post_first_name", post.first_name);
      localStorage.setItem("post_last_name", post.last_name);
    },
    setGroupPost(group_post: GroupPost) {
      this.group_post = group_post;
      localStorage.setItem("group_post_group_id", group_post.group_id.toString());
      localStorage.setItem("group_post_group_name", group_post.group_name);
      localStorage.setItem("group_post_group_description", group_post.group_description);
      localStorage.setItem("group_post_id", group_post.id.toString());
      localStorage.setItem("group_post_title", group_post.title);
      localStorage.setItem("group_post_categories", group_post.categories);
      localStorage.setItem("group_post_content", group_post.content);
      localStorage.setItem("group_post_picture", group_post.picture || "");
      localStorage.setItem("group_post_created_at", group_post.created_at);
      localStorage.setItem("group_post_email", group_post.email);
      localStorage.setItem("group_post_first_name", group_post.first_name);
      localStorage.setItem("group_post_last_name", group_post.last_name);
    },
  },
});
