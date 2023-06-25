export interface User {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  dob: string;
  avatar: Blob | null | string;
  nickname: string;
  aboutMe: string;
  public: boolean;
}
