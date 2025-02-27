import { writable } from 'svelte/store';

// Define writable stores for each collection ID
export const voiceLibraryId = writable(null);

// Function to update values in the stores
export function setVoiceLibraryIds(id) {
  voiceLibraryId.set(id);
}
