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
      id: sessionStorage.getItem("post_id") !== null ? parseInt(sessionStorage.getItem("post_id")!) : -1,
      title: sessionStorage.getItem("post_title") || "",
      categories: sessionStorage.getItem("post_categories") || "",
      content: sessionStorage.getItem("post_content") || "",
      privacy: sessionStorage.getItem("post_privacy") || "",
      picture: sessionStorage.getItem("post_picture") || "",
      created_at: sessionStorage.getItem("post_created_at") || "",
      email: sessionStorage.getItem("post_email") || "",
      first_name: sessionStorage.getItem("post_first_name") || "",
      last_name: sessionStorage.getItem("post_last_name") || "",
    },
    group_post: {
      group_id: sessionStorage.getItem("group_post_group_id") !== null ? parseInt(sessionStorage.getItem("group_post_group_id")!) : -1,
      group_name: sessionStorage.getItem("group_post_group_name") || "",
      group_description: sessionStorage.getItem("group_post_group_description") || "",
      id: sessionStorage.getItem("group_post_id") !== null ? parseInt(sessionStorage.getItem("group_post_id")!) : -1,
      title: sessionStorage.getItem("group_post_title") || "",
      categories: sessionStorage.getItem("group_post_categories") || "",
      content: sessionStorage.getItem("group_post_content") || "",
      picture: sessionStorage.getItem("group_post_picture") || "",
      created_at: sessionStorage.getItem("group_post_created_at") || "",
      email: sessionStorage.getItem("group_post_email") || "",
      first_name: sessionStorage.getItem("group_post_first_name") || "",
      last_name: sessionStorage.getItem("group_post_last_name") || "",
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
      sessionStorage.setItem("post_id", post.id.toString());
      sessionStorage.setItem("post_title", post.title);
      sessionStorage.setItem("post_categories", post.categories);
      sessionStorage.setItem("post_content", post.content);
      sessionStorage.setItem("post_privacy", post.privacy);
      sessionStorage.setItem("post_picture", post.picture || "");
      sessionStorage.setItem("post_created_at", post.created_at);
      sessionStorage.setItem("post_email", post.email);
      sessionStorage.setItem("post_first_name", post.first_name);
      sessionStorage.setItem("post_last_name", post.last_name);
    },
    setGroupPost(group_post: GroupPost) {
      this.group_post = group_post;
      sessionStorage.setItem("group_post_group_id", group_post.group_id.toString());
      sessionStorage.setItem("group_post_group_name", group_post.group_name);
      sessionStorage.setItem("group_post_group_description", group_post.group_description);
      sessionStorage.setItem("group_post_id", group_post.id.toString());
      sessionStorage.setItem("group_post_title", group_post.title);
      sessionStorage.setItem("group_post_categories", group_post.categories);
      sessionStorage.setItem("group_post_content", group_post.content);
      sessionStorage.setItem("group_post_picture", group_post.picture || "");
      sessionStorage.setItem("group_post_created_at", group_post.created_at);
      sessionStorage.setItem("group_post_email", group_post.email);
      sessionStorage.setItem("group_post_first_name", group_post.first_name);
      sessionStorage.setItem("group_post_last_name", group_post.last_name);
    },
  },
});
