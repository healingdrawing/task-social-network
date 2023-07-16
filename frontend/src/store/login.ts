import { defineStore } from 'pinia';
import { ErrorResponse } from '@/api/types';

export interface LoginSubmit {
  email: string;
  password: string;
}
export interface LoginResponse {
  UUID: string;
  email: string;
}

export const useLoginStore = defineStore('login', {
  state: () => ({
    data: {} as LoginResponse,
    error: '',
  }),
  getters: {
    getError(): string {
      return this.error;
    },
    getData(): LoginResponse {
      return this.data;
    }
  },
  actions: {
    async fetchData(userData: LoginSubmit) {
      console.log("= login =") // todo: remove debug
      try {
        const bodyJson = JSON.stringify(userData);
        const response = await fetch('http://localhost:8080/api/user/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Origin': 'http://localhost:8080'
          },
          body: bodyJson,
          mode: 'cors',
          credentials: 'omit' // <-- modified option
        });
        const data = await response.json();
        if (data.error) {
          throw new Error(data.error as string + "problem with json parsing of response");
        }

        console.log(data);
        this.data = data;

      } catch (error) {
        const errorResponse = error as ErrorResponse;
        this.error = errorResponse.message;
      } finally {
        console.log("= login = 'finally' fired")
      }
    }
  },
});