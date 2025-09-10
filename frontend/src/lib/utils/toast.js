import { writable } from "svelte/store";

export const toastCreated = writable(false);
export const toastUpdated = writable(false);
export const toastDeleted = writable(false);
export const toastError = writable(false);
export const toastNotAvailable = writable(false);
export const toastDayNotAvailable = writable(false);
export const toastTimeNotAvailable = writable(false);
export const toastCopyToClipboard = writable(false);
export let counter = 6;

export const customToasts = writable([]);

function triggerToast(target, isCustom = false, type = "success") {
  // console.log(target, isCustom, type);

  if (isCustom) {
    // console.log("custom");
    const toastId = Date.now(); // Menggunakan timestamp sebagai ID unik

    customToasts.update((toasts) => {
      return [...toasts, { id: toastId, target, isOpen: true, type }];
    });

    timeout(toastId, true); // Mengirimkan ID untuk timeout
    return;
  }

  switch (target) {
    case "created":
      toastCreated.set(true);
      break;
    case "updated":
      toastUpdated.set(true);
      break;
    case "deleted":
      toastDeleted.set(true);
      break;
    case "error":
      toastError.set(true);
      break;
    case "notAvailable":
      toastNotAvailable.set(true);
      break;
    case "dayNotAvailable":
      toastDayNotAvailable.set(true);
      break;
    case "timeNotAvailable":
      toastTimeNotAvailable.set(true);
      break;
    case "copyToClipboard":
      toastCopyToClipboard.set(true);
      break;
    default:
      console.log("target not found");
      return;
  }
  counter = 6;
  timeout(target);
}

function timeout(toastId, isCustom = false) {
  setTimeout(() => {
    if (isCustom) {
      customToasts.update((toasts) => {
        return toasts.filter((toast) => toast.id !== toastId);
      });
    } else {
      switch (toastId) {
        case "created":
          toastCreated.set(false);
          break;
        case "updated":
          toastUpdated.set(false);
          break;
        case "deleted":
          toastDeleted.set(false);
          break;
        case "error":
          toastError.set(false);
          break;
        case "notAvailable":
          toastNotAvailable.set(false);
          break;
        case "dayNotAvailable":
          toastDayNotAvailable.set(false);
          break;
        case "timeNotAvailable":
          toastTimeNotAvailable.set(false);
          break;
        case "copyToClipboard":
          toastCopyToClipboard.set(false);
          break;
        default:
          console.log("target not found");
          return;
      }
    }
  }, 5000); // 5 detik waktu untuk timeout, bisa disesuaikan sesuai kebutuhan
}

export function triggerToastResponse(response) {
  if (response.status !== 200) {
    console.log(response);
    if (response.status === 400) {
      console.error("Invalid parameters:", response);
      triggerToast($_("invalidParameters"), true, "error");
    } else if (response.status === 401) {
      console.error("Unauthorized access:", response);
      triggerToast($_("unauthorized"), true, "error");
    } else if (response.status === 500) {
      console.error("Server error:", response);
      triggerToast($_("serverError"), true, "error");
    } else {
      console.error("Unexpected status code:", response);
      triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
    }
    return false;
  }
  return true;
}

export default triggerToast;
