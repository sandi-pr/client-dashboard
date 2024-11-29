import { writable } from 'svelte/store';

// Define writable stores for each collection ID
export const bakrieLibraryId = writable(null);
export const documentsCollectionId = writable(null);
export const imagesCollectionId = writable(null);
export const documentImagesCollectionId = writable(null);
export const handwritingCollectionId = writable(null);
export const audioCollectionId = writable(null);

// Function to update values in the stores
export function setLibraryIds(id) {
  bakrieLibraryId.set(id);
}

export function setCollectionIds(ids) {
  const {
    documentsCollectionIdValue,
    imagesCollectionIdValue,
    documentImagesCollectionIdValue,
    handwritingCollectionIdValue,
    audioCollectionIdValue,
  } = ids;

  documentsCollectionId.set(documentsCollectionIdValue);
  imagesCollectionId.set(imagesCollectionIdValue);
  documentImagesCollectionId.set(documentImagesCollectionIdValue);
  handwritingCollectionId.set(handwritingCollectionIdValue);
  audioCollectionId.set(audioCollectionIdValue);
}
