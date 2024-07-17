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

export const showAI = writable(false);

export function triggerShowAI(loading) {
  if (loading) {
    showAI.set(true);
  } else {
    showAI.set(false);
  }
}
