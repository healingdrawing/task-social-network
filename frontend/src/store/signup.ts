import { defineStore } from 'pinia';
import { ErrorResponse } from '@/api/types';

// interface SignupLoginResponse {
//   UUID: string;
//   email: string;
// }


export interface SignupSubmit {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  dob: string;
  avatar: Blob | null | string;
  nickname: string;
  about_me: string;
  public: boolean;
}
export interface SignupResponse {
  UUID: string;
  email: string;
}

export const useSignupStore = defineStore('signup', {
  state: () => ({
    data: {} as SignupResponse,
    error: '',
  }),
  getters: {
    getError(): string {
      return this.error;
    },
    getData(): SignupResponse {
      return this.data;
    }
  },
  actions: {
    async fetchData(userData: SignupSubmit) {
      console.log("= signup =")
      try {
        if (!userData.avatar) {
          userData.avatar = ''
        }
        const bodyJson = JSON.stringify(userData);

        const response = await fetch('http://localhost:8080/api/user/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Origin': 'http://localhost:8080'
          },
          body: bodyJson,
          mode: 'cors',
          credentials: 'omit' // <-- modified option
        });
        const resp = await response;
        if (resp.status !== 200) {
          throw new Error('Network response was not ok');
        }
        console.log(resp);

        const data = await resp.json();
        if (data.error) {
          throw new Error(data.error as string + "problem with json parsing of response");
        }

        console.log(data);
        this.data = data;

      } catch (error) {
        const errorResponse = error as ErrorResponse;
        this.error = errorResponse.message;
      } finally {
        console.log("= signup = 'finally' fired")
      }
    }
  },
});