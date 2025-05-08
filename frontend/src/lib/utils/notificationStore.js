import { writable } from "svelte/store";

export const notifications = writable([]); // Menyimpan daftar notifikasi
export const notificationsUnread = writable([]); // Menyimpan daftar notifikasi yang belum dibaca
export const notificationsRead = writable([]); // Menyimpan daftar notifikasi yang sudah dibaca
export let socket = null; // WebSocket instance
export let notificationsNewCount = writable(0);
export let newNotification = writable(false);
export let onProgress = writable(false);

export const paginationAll = writable({
  page: 1,
  page_size: 25,
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

let timeout = null;

export function connectWebSocket(token, isFrontend) {
  console.log("Connecting WebSocket...");

  if (socket && socket.readyState === WebSocket.OPEN) {
    console.log("WebSocket already connected.");
    return;
  }

  let reconnectAttempts = 0; // Track the number of reconnection attempts
  const maxReconnectAttempts = 5; // Maximum allowed reconnection attempts

  const protocol = location.protocol === "https:" ? "wss:" : "ws:";
  const host = location.hostname;
  const port = location.port ? `:${location.port}` : "";
  const base_url = `${protocol}//${host}${port}`;

  if (isFrontend) {
    socket = new WebSocket(
      `ws://tank.hexacode.ai/client-dashboard/notification/${encodeURIComponent(token)}`,
    );
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
      // console.log("Received notification:", message);

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
              : notif,
          ),
        );

        onProgress.set(true);

        if (timeout) {
          clearTimeout(timeout);
        }

        if (message.content.includes("100")) {
          onProgress.set(false);
        }
      } else {
        if (message.type === "progress") {
          onProgress.set(true);
        }
        if (message.title.includes("ready")) {
          onProgress.set(false);
        }

        timeout = setTimeout(() => {
          onProgress.set(false);
        }, 10000);
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

    if (reconnectAttempts < maxReconnectAttempts) {
      reconnectAttempts++;
      console.log(`Reconnection attempt ${reconnectAttempts}...`);
      setTimeout(() => {
        connectWebSocket(token, isFrontend);
      }, 3000);
    } else {
      console.log("Maximum reconnection attempts reached. Stopping further attempts.");
    }
  };
}

export function disconnectWebSocket() {
  if (socket) {
    socket.close();
    console.log("WebSocket manually disconnected.");
  }
}
