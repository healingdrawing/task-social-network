import { defineStore } from 'pinia';
import { User } from '../api/types';

interface SignupLoginResponse {
  UUID: string;
  email: string;
}

interface ErrorResponse {
  message: string;
}

export const useSignupLoginStore = defineStore('signupLogin', {
  state: () => ({
    data: {} as SignupLoginResponse,
    error: '',
  }),
  getters: {
    getError(): string {
      return this.error;
    },
    getData(): SignupLoginResponse {
      return this.data;
    }
  },
  actions: {
    async storeFetchData(userData: User) {
      console.log("stage 0")
      try {
        if (!userData.avatar) {
          userData.avatar = ''
        }
        console.log("request signup userData: ", userData);
        const bodyJson = JSON.stringify(userData);
        console.log("request signup userData: ", bodyJson);

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
        console.log("stage 1")
        const resp = await response;
        if (resp.status !== 200) {
          throw new Error('Network response was not ok');
        }
        console.log("stage 2")
        console.log(resp);
        const data = await resp.json();
        if (data.error) {
          throw new Error(data.error as string + "problem with json parsing of response");
        }
        console.log("stage 3")

        console.log(data);
        this.data = data;
      } catch (error) {
        const errorResponse = error as ErrorResponse;
        this.error = errorResponse.message;
      } finally {
        console.log("stage 4")
      }
    }
  },
});