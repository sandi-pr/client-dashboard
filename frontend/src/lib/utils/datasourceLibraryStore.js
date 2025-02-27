import { writable } from 'svelte/store';

// Define writable stores for each collection ID
export const datasourceLibraryId = writable(null);

// Function to update values in the stores
export function setDatasourceLibraryIds(id) {
  datasourceLibraryId.set(id);
}
