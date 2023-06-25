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
    async signup(userData: User) {
      try {
        if (!userData.avatar) {
          userData.avatar = ''
        }
        console.log("request signup userData: ", userData);
        console.log("request signup userData: ", JSON.stringify(userData));
        const response = await fetch('http://127.0.0.1:8000/api/user/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(userData),
        });
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const resp = await response;
        console.log(resp);
        const data = await resp.json();

        console.log(data);
        this.data = data;
      } catch (error) {
        const errorResponse = JSON.parse(error as string) as ErrorResponse;
        this.error = errorResponse.message;
      }
    }
  },
});