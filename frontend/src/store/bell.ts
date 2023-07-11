import { defineStore } from 'pinia';
import { Bell, BellState } from '@/api/types';

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
