/*
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: 
  - Taqie Fadlillah <taqief@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
*/

import { writable } from 'svelte/store';

export const toastCreated = writable(false);
export const toastUpdated = writable(false);
export const toastDeleted = writable(false);
export const toastError = writable(false);
export const toastNotAvailable = writable(false);
export const toastDayNotAvailable = writable(false);
export const toastTimeNotAvailable = writable(false);
export const toastCopyToClipboard = writable(false);
export let counter = 6;

export const customToasts = writable({});

function triggerToast(target, isCustom = false, type = 'success') {

  console.log(target, isCustom, type);

  if (isCustom) {
    customToasts.update((toasts) => {
      toasts[target] = { isOpen: true, type };
      return toasts;
    });
    timeout(target, true);
    return;
  }

  switch (target) {
    case 'created':
      toastCreated.set(true);
      break;
    case 'updated':
      toastUpdated.set(true);
      break;
    case 'deleted':
      toastDeleted.set(true);
      break;
    case 'error':
      toastError.set(true);
      break;
    case 'notAvailable':
      toastNotAvailable.set(true);
      break;
    case 'dayNotAvailable':
      toastDayNotAvailable.set(true);
      break;
    case 'timeNotAvailable':
      toastTimeNotAvailable.set(true);
      break;
    case 'copyToClipboard':
      toastCopyToClipboard.set(true);
      break;
    default:
      console.log('target not found');
      return;
  }
  counter = 6;
  timeout(target);
}

function timeout(target, isCustom = false) {
  if (--counter > 0) return setTimeout(() => timeout(target, isCustom), 1000);

  if (isCustom) {
    // Deactivate custom toast
    customToasts.update((toasts) => {
      delete toasts[target];
      return toasts;
    });
    return;
  }

  switch (target) {
    case 'created':
      toastCreated.set(false);
      break;
    case 'updated':
      toastUpdated.set(false);
      break;
    case 'deleted':
      toastDeleted.set(false);
      break;
    case 'error':
      toastError.set(false);
      break;
    case 'notAvailable':
      toastNotAvailable.set(false);
      break;
    case 'dayNotAvailable':
      toastDayNotAvailable.set(false);
      break;
    case 'timeNotAvailable':
      toastTimeNotAvailable.set(false);
      break;
    case 'copyToClipboard':
      toastCopyToClipboard.set(false);
      break;
    default:
      console.log('target not found');
      return;
  }
}

export default triggerToast;
