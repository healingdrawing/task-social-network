interface ResizedImage {
  blob: Blob;
  width: number;
  height: number;
}

/**prepare picture data for send to backend in string form*/
export function blobToBase64(blob: Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(blob);
    reader.onloadend = () => {
      const base64String = reader.result?.toString().split(',')[1];
      if (base64String) {
        resolve(base64String);
      } else {
        reject(new Error('Failed to convert Blob to Base64-encoded string'));
      }
    };
    reader.onerror = () => {
      reject(new Error('Failed to read Blob'));
    };
  });
}

// todo: remove later. not used at the moment. not tested.
export function resizeImage(file: File, maxSizeKb: number, maxWidthPx: number, maxHeightPx: number): Promise<ResizedImage> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.onload = () => {
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');
      const MAX_SIZE = maxSizeKb * 1024; // 500kb
      let width = img.width;
      let height = img.height;

      if (width > height) {
        if (width > maxWidthPx) {
          height *= maxWidthPx / width;
          width = maxWidthPx;
        }
      } else {
        if (height > maxHeightPx) {
          width *= maxHeightPx / height;
          height = maxHeightPx;
        }
      }

      canvas.width = width;
      canvas.height = height;
      ctx?.drawImage(img, 0, 0, width, height);
      canvas.toBlob((blob) => {
        if (blob) {
          if (blob.size > MAX_SIZE) {
            reject(new Error('Image size must be less than 500kb'));
          } else {
            const resizedImage: ResizedImage = {
              blob: blob,
              width: width,
              height: height
            };
            resolve(resizedImage);
          }
        }
      }, 'image/jpeg', 0.7);
    };
    img.onerror = () => {
      reject(new Error('Failed to load image'));
    };
    img.src = URL.createObjectURL(file);
  });
}