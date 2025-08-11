import { writable } from "svelte/store";

export const me = writable({ username: "", email: "" });
export const isDemoTrue = writable(false);
export const isLimitTrue = writable(false);
