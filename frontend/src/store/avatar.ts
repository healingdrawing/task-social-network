import { defineStore } from "pinia";

interface AvatarState {
  avatarError: string;
  avatarBlob: Blob | null;
}

export const useAvatarStore = defineStore({
  id: "avatar",
  state: (): AvatarState => ({
    avatarError: "",
    avatarBlob: null,
  }),
  getters: {
    getAvatarError(): string {
      return this.avatarError;
    }
  },
  actions: {
    setAvatarError(avatarError: string) {
      this.avatarError = avatarError;
    },
    handleAvatarUpload(event: Event) {
      const file = (event.target as HTMLInputElement).files?.[0];
      if (!file) return;

      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => {
        const img = new Image();
        img.src = reader.result as string;
        img.onload = () => {
          const canvas = document.createElement('canvas');
          const ctx = canvas.getContext('2d');
          const MAX_WIDTH = 500;
          const MAX_HEIGHT = 500;
          let width = img.width;
          let height = img.height;

          if (width > height) {
            if (width > MAX_WIDTH) {
              alert('store/avatar.ts width > height && width > MAX_WIDTH');
              height *= MAX_WIDTH / width;
              width = MAX_WIDTH;
            }
          } else {
            if (height > MAX_HEIGHT) {
              alert('store/avatar.ts height > MAX_HEIGHT');
              width *= MAX_HEIGHT / height;
              height = MAX_HEIGHT;
            }
          }

          canvas.width = width;
          canvas.height = height;

          ctx?.drawImage(img, 0, 0, width, height);

          canvas.toBlob((blob) => {
            if (!blob) return;

            if (blob.size > 500000) {
              this.avatarError = 'The image must be less than 500 KB.';
            } else if (!['image/jpeg', 'image/png'].includes(blob.type)) {
              this.avatarError = 'The image must be in JPEG or PNG format.';
            } else if (width !== img.width || height !== img.height) {
              this.avatarError = 'The image must be less than or equal to 500 pixels in width and height.';
            } else {
              this.avatarBlob = blob;
              this.avatarError = '';
            }
          }, 'image/jpeg', 0.8);
        };
      };
    },
  }
});