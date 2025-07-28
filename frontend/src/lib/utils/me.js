import { writable } from "svelte/store";

export const me = writable({ username: "", email: "" });
