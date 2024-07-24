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

export const selectedCollection = writable(null);
export const selectedDropdown = writable("general");
export const reload = writable(true);

export function triggerSelectCollection(select) {
  selectedCollection.set(select);
}

export function triggerSelectDropdown(select) {
  selectedDropdown.set(select);
}

export function triggerReload(value) {
  reload.set(value);
}
