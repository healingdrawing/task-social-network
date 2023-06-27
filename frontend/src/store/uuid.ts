import { defineStore } from 'pinia';

interface UUIDState {
  UUID: string;
}

export const useUUIDStore = defineStore({
  id: 'uuid',
  state: (): UUIDState => ({
    // Define your state properties here
    UUID: "-1",// the order of these properties is not matter. it is object field
  }),
  getters: {
    // Define your getters here
    getUUID(): string { return this.UUID; },
  },
  actions: {
    // Define your actions here
    setUUID(UUID: string) { this.UUID = UUID; },
  },
});
