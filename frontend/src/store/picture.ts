import { defineStore } from "pinia";
import { blobToBase64 } from "@/api/tools";

interface PictureState {
  pictureError: string;
  pictureBlob: Blob | null;
  pictureBase64String: string;
}

export const usePictureStore = defineStore({
  id: "picture",
  state: (): PictureState => ({
    pictureError: "",
    pictureBlob: null,
    pictureBase64String: "",
  }),
  getters: {
    getPictureError(): string {
      return this.pictureError;
    },
    getPictureBlob(): Blob | null {
      return this.pictureBlob;
    },
    // async getBase64forJson(state) {
    //   if (!state.pictureBlob) {
    //     return "null";
    //   }
    //   const base64String = await blobToBase64(state.pictureBlob);
    //   return base64String;
    // },
    getPictureBase64String(): string {
      return this.pictureBase64String;
    },
  },
  actions: {
    setPictureError(pictureError: string) {
      this.pictureError = pictureError;
    },
    handlePictureUpload(event: Event) {
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
              alert('store/picture.ts width > height && width > MAX_WIDTH');
              height *= MAX_WIDTH / width;
              width = MAX_WIDTH;
            }
          } else {
            if (height > MAX_HEIGHT) {
              alert('store/picture.ts height > MAX_HEIGHT');
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
              this.pictureError = 'The image must be less than 500 KB.';
            } else if (!['image/jpeg', 'image/png'].includes(blob.type)) {
              this.pictureError = 'The image must be in JPEG or PNG format.';
            } else if (width !== img.width || height !== img.height) {
              this.pictureError = 'The image must be less than or equal to 500 pixels in width and height.';
            } else {
              this.pictureBlob = blob;
              this.pictureError = '';
              this.pictureBase64String = reader.result?.toString() || '';
            }
          }, 'image/jpeg', 0.8);
        };
      };
    },
  }
});