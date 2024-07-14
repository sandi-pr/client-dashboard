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

export const toastLoading = writable(false);

function triggerLoading(loading) {
  if (loading) {
    toastLoading.set(true);
  } else {
    toastLoading.set(false);
  }
}

export default triggerLoading;
