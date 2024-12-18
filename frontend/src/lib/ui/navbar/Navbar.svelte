<!--
  CR3 Client Dashboard
  Package: gitlab.hexacode.org/public-repo/js/hexacomponents-svelte
  Maintainer:
  * Sandi Permana <sandy@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>

  Copyright (C) 2024 Hexacode Teknologi Indonesia
  All Rights Reserved.
-->

<script>
  import { Breadcrumb, Button, LanguageSelector } from "$lib";
  import RefreshCwIcon from "$lib/icon/refresh.svg";
  import { Modal } from "flowbite-svelte";
  import { DarkMode } from "flowbite-svelte";
  import BellIcon from "$lib/icon/bell-01.svg";
  import StarsIcon from "$lib/icon/stars-01.svg";
  import { ExclamationCircleSolid, ChevronDownOutline } from "flowbite-svelte-icons";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import { _ } from "svelte-i18n";

  let deleteModal = $state(false);
  let dropdownOpen = $state(false);

  let {
    activeUrl = "/",
    changeShowAI = () => {},
    changeShowNotif = () => {},
    appUrl = "/",
    reset_database = () => {},
    authority = $bindable(""),
    hasNewNotif = $bindable(false),
  } = $props();

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
  };
</script>

<div>
  <div class="flex justify-between px-4 pt-1 lg:pt-8">
    <Breadcrumb {activeUrl} {appUrl} />
    <div class="flex items-center gap-2 md:gap-3">
      <div class="hidden md:flex gap-2">
        <LanguageSelector />
        {#if authority === "Admin"}
          <!-- <Button
              color="light"
              icon={RefreshCwIcon}
              text=""
              spaceY="mt-0"
              funcButton={() => (deleteModal = true)}
            /> -->
        {/if}
        <DarkMode
          btnClass="text-gray-500 dark:text-gray-400 bg-white dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none rounded-lg p-2.5 border border-gray-300 dark:border-gray-600 dark:hover:border-gray-600 focus-within:ring-4 focus-within:outline-none focus-within:ring-gray-200 dark:focus-within:ring-gray-700"
        />
      </div>
      <div class="relative">
        <Button
          color="light"
          icon={BellIcon}
          classesSvg="bell"
          spaceY="mt-0"
          funcButton={() => changeShowNotif()}
        />
        {#if hasNewNotif}
          <span class="absolute top-0 right-0 inline-block w-2.5 h-2.5 bg-red-500 rounded-full"
          ></span>
        {/if}
      </div>
      <!-- <Button
        color="primary"
        icon={StarsIcon}
        spaceY="mt-0"
        funcButton={() => {
          changeShowAI();
        }}
      /> -->
    </div>
  </div>
</div>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToResetTheDatabase")}
    </h3>
    <div class="flex justify-center space-x-4">
      <Button
        text={$_("yesImSure")}
        funcButton={() => {
          reset_database();
        }}
        color="red"
        spaceY="mt-0"
        class="me-2"
      ></Button>
      <Button
        text={$_("noCancel")}
        color="alternative"
        funcButton={() => (deleteModal = false)}
        spaceY="mt-0"
      ></Button>
    </div>
  </div>
</Modal>
