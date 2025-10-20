// src/lib/stores/window.js
import { readable } from "svelte/store";

// A simple readable store for the window dimensions
export const windowSize = readable(
  {
    width: typeof window !== "undefined" ? window.innerWidth : 0,
    height: typeof window !== "undefined" ? window.innerHeight : 0,
  },
  function start(set) {
    // Only set up the listener if we're in the browser
    if (typeof window === "undefined") {
      return;
    }

    function handleResize() {
      set({
        width: window.innerWidth,
        height: window.innerHeight,
      });
    }

    // Set initial value
    handleResize();

    // Add the event listener for resizing
    window.addEventListener("resize", handleResize);

    // Return the stop function to clean up when no one is subscribed
    return function stop() {
      window.removeEventListener("resize", handleResize);
    };
  },
);
