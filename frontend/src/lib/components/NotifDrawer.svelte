<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: Azzis Arswendo <azzis@hexacode.org>

Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->

<script>
  import { Typography, TabBar, Icon } from "$lib";
  import InfoCircleIcon from "$lib/icon/info-circle.svelte";
  import { Pane } from "svelte-splitpanes";
  import {
    CloseButton,
    Button,
    Toast,
  } from "flowbite-svelte";
  import { _ } from "svelte-i18n";

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
            <div
              class="flex grow justify-between pb-4"
            >
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
          <Toast align={false} divClass='w-full p-4 text-gray-500 bg-white shadow dark:text-gray-400 dark:bg-gray-800 gap-3'>
            <svelte:fragment slot="icon">
              <InfoCircleIcon svgClass="w-6 h-6" />
            </svelte:fragment>
          
            <span class="font-semibold text-gray-900 dark:text-white">Reading [Hexabot’s Name]’s documents</span>
            <div class="mt-1">
              <div class="mb-2 text-sm font-normal">Hexacube is now analyzing the documents uploaded by HexaBot. Sit tight while we process the information!</div>
              <div class="flex gap-3">
                <button class="text-sm font-semibold text-fg-button-tertiary-light dark:text-fg-button-tertiary-dark">Mark as read</button>
                <button class="text-sm font-semibold text-fg-button-tertiary-color-light dark:text-fg-button-tertiary-color-dark">View detail</button>
              </div>
            </div>
          </Toast>
        </div>
      {:else if tabDatas.selected === "Unread"}
        <div class="p-4">
          Unread
        </div>
      {:else if tabDatas.selected === "Read"}
        <div class="p-4">
          Read
        </div>
      {/if}
    </div>
  </div>
</Pane>
