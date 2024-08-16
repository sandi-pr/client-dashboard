import { writable } from 'svelte/store';

export const call = writable(null)
export const callReady = writable(false)

export function setCall(cr3Call) {
  call.set(cr3Call)
  callReady.set(true)
}