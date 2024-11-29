<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: Azzis Arswendo <azzis@hexacode.org>

Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->

<script>
  import {
    Typography,
    TabBar,
    call,
    notifications,
    notificationsUnread,
    notificationsRead,
    paginationAll,
    paginationRead,
    paginationUnread,
    notificationsNewCount,
  } from "$lib";
  import InfoCircleIcon from "$lib/icon/info-circle.svelte";
  import { Pane } from "svelte-splitpanes";
  import { CloseButton, Button, Toast, Progressbar } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { onMount } from "svelte";
  import { get } from "svelte/store";
  import { goto } from "$app/navigation";

  let { innerWidth, changeShowNotif } = $props();

  let tabDatas = $state({
    selected: "All",
    data: [
      {
        name: "All",
      },
      {
        name: "Unread",
      },
      {
        name: "Read",
      },
    ],
  });

  const changeTabbar = (tabName) => {
    tabDatas.selected = tabName;
  };

  async function getNotifications() {
    let token = localStorage.getItem("token");

    try {
      const response = await $call.call("client-dashboard", "get-notifications", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: $paginationAll.page,
          page_size: $paginationAll.page_size,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      if (response.result.data.length > 0) {
        notifications.update((data) => [response.result.data, ...data]);
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function getNotificationsUnread() {
    let token = localStorage.getItem("token");

    try {
      const response = await $call.call("client-dashboard", "get-notifications-unread", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: $paginationUnread.page,
          page_size: $paginationUnread.page_size,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      if (response.result.data.length > 0) {
        notificationsUnread.update((data) => [response.result.data, ...data]);
        notificationsNewCount.update((data) => data + response.result.meta.pagination.total);
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function getNotificationsRead() {
    let token = localStorage.getItem("token");

    try {
      const response = await $call.call("client-dashboard", "get-notifications-read", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: $paginationRead.page,
          page_size: $paginationRead.page_size,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      if (response.result.data.length > 0) {
        notificationsRead.update((data) => [response.result.data, ...data]);
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function markAsRead(id) {
    let token = localStorage.getItem("token");

    try {
      const response = await $call.call("client-dashboard", "read-notification", {
        URLQuery: {},
        URLParams: {
          id: id.toString(),
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      notificationsUnread.set([]);
      paginationUnread.set(1);
      notificationsRead.set([]);
      paginationRead.set(1);
      notificationsUnread.set([]);
      paginationUnread.set(0);

      getNotifications();
      getNotificationsRead();
      getNotificationsUnread();
    } catch (error) {
      console.log(error);
    }
  }

  onMount(() => {
    // getNotifications();
    // getNotificationsUnread();
    // getNotificationsRead();
  });
</script>

<Pane size={innerWidth > 1024 ? 40 : 100} minSize={innerWidth > 1024 ? 40 : 100}>
  <div class="h-screen transition-all duration-300 border-l-2 border-gray-200 dark:border-gray-800">
    <div class="sticky top-0 z-10">
      <div class="flex items-center">
        <div
          id="drawer-label"
          class="w-full relative inline-flex items-center mx-4 mt-4 text-base font-semibold text-gray-500 dark:text-gray-400"
        >
          <div class="flex flex-col grow">
            <div class="flex grow justify-between pb-4">
              <Typography preset="Text xl/Semibold" text="Notifications" color="primary" />
              <CloseButton
                on:click={changeShowNotif}
                class="mb-4 absolute right-0 top-1 dark:text-white"
              />
            </div>
          </div>
        </div>
      </div>
      <TabBar {tabDatas} {changeTabbar} />
    </div>
    <div class="flex w-full flex-col h-screen">
      {#if tabDatas.selected === "All"}
        <div class="p-4">
          {#if $notifications.length === 0}
            <Typography
              preset="Text sm/Medium"
              text="No Notification"
              color="secondary"
              class="w-full text-center"
            />
          {:else}
            {#each $notifications as data}
              <Toast
                align={false}
                divClass="w-full p-4 text-gray-500 bg-white shadow dark:text-gray-400 dark:bg-gray-800 gap-3"
              >
                <svelte:fragment slot="icon">
                  <InfoCircleIcon svgClass="w-6 h-6" />
                </svelte:fragment>

                <span class="font-semibold text-gray-900 dark:text-white"
                  >{data.attributes.title}</span
                >
                <div class="mt-1 gap-2">
                  <div class="mb-2 text-sm font-normal">
                    {data.attributes.content}
                  </div>
                  <Progressbar progress="50" size="h-4" labelInside />
                  <div class="flex gap-3 pt-2">
                    {#if !data.attributes.read}
                      <button
                        onclick={async () => {
                          await markAsRead(data.id);
                        }}
                        class="text-sm font-semibold text-fg-button-tertiary-light dark:text-fg-button-tertiary-dark"
                        >Mark as read</button
                      >
                    {/if}
                    {#if data.attributes.type === "hexabot" && data.attributes.meta && "bot_id" in data.attributes.meta}
                      <button
                        onclick={async () => {
                          goto(
                            (await $call.getPage("client-dashboard", "hexabot-detail")) +
                              "?botId=" +
                              data.attributes.meta.bot_id,
                          );
                        }}
                        class="text-sm font-semibold text-fg-button-tertiary-color-light dark:text-fg-button-tertiary-color-dark"
                        >View detail</button
                      >
                    {/if}
                  </div>
                </div>
              </Toast>
            {/each}
          {/if}
        </div>
      {:else if tabDatas.selected === "Unread"}
        <div class="p-4">
          {#if $notificationsUnread.length === 0}
            <Typography
              preset="Text sm/Medium"
              text="No Notification"
              color="secondary"
              class="w-full text-center"
            />
          {:else}
            {#each $notificationsUnread as data}
              <Toast
                align={false}
                divClass="w-full p-4 text-gray-500 bg-white shadow dark:text-gray-400 dark:bg-gray-800 gap-3"
              >
                <svelte:fragment slot="icon">
                  <InfoCircleIcon svgClass="w-6 h-6" />
                </svelte:fragment>

                <span class="font-semibold text-gray-900 dark:text-white"
                  >{data.attributes.title}</span
                >
                <div class="mt-1">
                  <div class="mb-2 text-sm font-normal">
                    {data.attributes.content}
                  </div>
                  <div class="flex gap-3">
                    <button
                      class="text-sm font-semibold text-fg-button-tertiary-light dark:text-fg-button-tertiary-dark"
                      >Mark as read</button
                    >
                    {#if data.attributes.type === "hexabot" && data.attributes.meta && "bot_id" in data.attributes.meta}
                      <button
                        onclick={async () => {
                          goto(
                            (await $call.getPage("client-dashboard", "hexabot-detail")) +
                              "?botId=" +
                              data.attributes.meta.bot_id,
                          );
                        }}
                        class="text-sm font-semibold text-fg-button-tertiary-color-light dark:text-fg-button-tertiary-color-dark"
                        >View detail</button
                      >
                    {/if}
                  </div>
                </div>
              </Toast>
            {/each}
          {/if}
        </div>
      {:else if tabDatas.selected === "Read"}
        <div class="p-4">
          {#if $notificationsRead.length === 0}
            <Typography
              preset="Text sm/Medium"
              text="No Notification"
              color="secondary"
              class="w-full text-center"
            />
          {:else}
            {#each $notificationsRead as data}
              <Toast
                align={false}
                divClass="w-full p-4 text-gray-500 bg-white shadow dark:text-gray-400 dark:bg-gray-800 gap-3"
              >
                <svelte:fragment slot="icon">
                  <InfoCircleIcon svgClass="w-6 h-6" />
                </svelte:fragment>

                <span class="font-semibold text-gray-900 dark:text-white"
                  >{data.attributes.title}</span
                >
                <div class="mt-1">
                  <div class="mb-2 text-sm font-normal">
                    {data.attributes.content}
                  </div>
                  <div class="flex gap-3">
                    <button
                      onclick={async () => {
                        await markAsRead(data.id);
                      }}
                      class="text-sm font-semibold text-fg-button-tertiary-light dark:text-fg-button-tertiary-dark"
                      >Mark as read</button
                    >
                    {#if data.attributes.type === "hexabot" && data.attributes.meta && "bot_id" in data.attributes.meta}
                      <button
                        onclick={async () => {
                          goto(
                            (await $call.getPage("client-dashboard", "hexabot-detail")) +
                              "?botId=" +
                              data.attributes.meta.bot_id,
                          );
                        }}
                        class="text-sm font-semibold text-fg-button-tertiary-color-light dark:text-fg-button-tertiary-color-dark"
                        >View detail</button
                      >
                    {/if}
                  </div>
                </div>
              </Toast>
            {/each}
          {/if}
        </div>
      {/if}
    </div>
  </div>
</Pane>
