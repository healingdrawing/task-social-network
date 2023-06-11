import { createPinia } from 'pinia';

export const pinia = createPinia();

// Export the Pinia app instance for use in your main entry file
// export const usePinia = () => pinia;

export { useProfileStore } from '@/store/profile';
