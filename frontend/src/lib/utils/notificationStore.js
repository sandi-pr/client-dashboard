import { writable } from 'svelte/store';

export const notifications = writable([]); // Menyimpan daftar notifikasi
export const notificationsUnread = writable([]); // Menyimpan daftar notifikasi yang belum dibaca
export const notificationsRead = writable([]); // Menyimpan daftar notifikasi yang sudah dibaca
export let socket = null; // WebSocket instance
export let notificationsNewCount = writable(0);
export let newNotification = writable(false);

export const paginationAll = writable({
  page: 1,
  page_size: 10,
  page_count: 0,
  total: 0,
});

export const paginationUnread = writable({
  page: 1,
  page_size: 10,
  page_count: 0,
  total: 0,
});

export const paginationRead = writable({
  page: 1,
  page_size: 10,
  page_count: 0,
  total: 0,
});

export function connectWebSocket(token, isFrontend) {
  console.log("Connecting WebSocket...");

  if (socket && socket.readyState === WebSocket.OPEN) {
    console.log("WebSocket already connected.");
    return;
  }

  var base_url = (location.protocol == 'https:' ? "wss:" : "ws:") + "//" + location.hostname + (location.port == "" ? "" : ":" + location.port);

  if (isFrontend) {
    socket = new WebSocket(
      `ws://192.168.18.14/client-dashboard/notification/${encodeURIComponent(token)}`,
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

      if (message.type === "progress-update") {
        let docId = message.meta?.document_id;

        notifications.update((n) =>
          n.map((notif) =>
            notif.attributes.meta?.document_id === docId
              ? {
                ...notif,
                attributes: {
                  ...notif.attributes,
                  content: message.content,
                  meta: {
                    ...notif.attributes.meta,
                    ...message.meta,
                  },
                },
              }
              : notif
          )
        );
      } else {
        newNotification.update((data) => true);
      }
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