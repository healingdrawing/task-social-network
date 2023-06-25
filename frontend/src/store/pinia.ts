import { createPinia } from 'pinia';

export const pinia = createPinia();

// Export the Pinia app instance for use in your main entry file
// export const usePinia = () => pinia;

export { useProfileStore } from '@/store/profile';
export { useAvatarStore } from '@/store/avatar';
export { usePictureStore } from '@/store/picture';
export { usePostStore } from '@/store/post';
export { useChatStore } from '@/store/chat';
export { useChatsStore } from '@/store/chats';
export { useGroupStore } from '@/store/group';
export { useBellStore } from '@/store/bell';
export { useSignupLoginStore } from '@/store/signupLogin';
