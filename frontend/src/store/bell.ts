import { defineStore } from 'pinia';

export type BellType =
  | 'event'
  | 'following'
  | 'invitation'
  | 'request';

export interface Bell {
  type: BellType;
  message: string;
  groupId: number;
  userId: number;
}

interface BellState {
  bells: Bell[];
}

export const useBellStore = defineStore({
  id: 'bell',
  state: (): BellState => ({
    // Define your state properties here
    bells: [],
  }),
  getters: {
    // Define your getters here
    getBells(): Bell[] { return this.bells; },
  },
  actions: {
    // Define your actions here
    setBells(bells: Bell[]) { this.bells = bells; },
  },
});
