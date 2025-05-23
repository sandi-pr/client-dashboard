import { checkLibraryName } from "./checkLibraryName";

let debounceTimeout;

export const handleInputName = (call, event, defaultName, collectionId, pageDataType) => {
  return new Promise((resolve) => {
    const name = event.target.value.trim();

    // Reset the debounce timer
    clearTimeout(debounceTimeout);

    // If the name matches the default name, resolve immediately
    if (name === defaultName) {
      resolve(true);
      return;
    }

    // Set a debounce delay of 0.5 seconds before calling checkLibraryName
    debounceTimeout = setTimeout(async () => {
      const status = await checkLibraryName(call, name, collectionId, pageDataType);
      resolve(!!status); // Resolve with the result of checkLibraryName
    }, 500);
  });
};
