export interface User {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  dob: string;
  avatar: Blob | null;
  nickname: string;
  aboutMe: string;
  public: boolean;
}
