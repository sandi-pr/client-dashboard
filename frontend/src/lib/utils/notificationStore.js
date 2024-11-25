import { writable } from 'svelte/store';

export const notifications = writable([]); // Menyimpan daftar notifikasi
export let socket = null; // WebSocket instance

export function connectWebSocket(token, isFrontend) {
  console.log("Connecting WebSocket...");

  if (socket && socket.readyState === WebSocket.OPEN) {
    console.log("WebSocket already connected.");
    return;
  }

  var base_url = (location.protocol == 'https:' ? "wss:" : "ws:") + "//" + location.hostname + (location.port == "" ? "" : ":" + location.port);

  if (isFrontend) {
    socket = new WebSocket(
      `ws://192.168.1.11/client-dashboard/notification/${encodeURIComponent(token)}`,
    )
  } else {
    socket = new WebSocket(
      `${base_url}/client-dashboard/notification/${encodeURIComponent(token)}`,
    );
  }


  socket.onopen = () => {
    console.log("WebSocket connected.");
  };

  socket.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data);
      console.log("Received notification:", message);
      notifications.update((n) => [message, ...n]); // Tambahkan notifikasi ke store
    } catch (error) {
      console.error("Error parsing WebSocket message:", error);
    }
  };

  socket.onerror = (error) => {
    console.error("WebSocket error:", error);
  };

  socket.onclose = () => {
    console.log("WebSocket disconnected.");
    // Opsi: coba reconnect otomatis jika diperlukan
    setTimeout(() => {
      connectWebSocket(token, isFrontend);
    }, 3000);
  };
}

export function disconnectWebSocket() {
  if (socket) {
    socket.close();
    console.log("WebSocket manually disconnected.");
  }
}