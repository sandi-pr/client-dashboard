<!--
  CR3 Client Dashboard
  Package: gitlab.hexacode.org/public-repo/js/hexacomponents-svelte
  Maintainer:
  * Sandi Permana <sandy@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>

  Copyright (C) 2024 Hexacode Teknologi Indonesia
  All Rights Reserved.
-->

<script lang="ts">
  import { Breadcrumb, Button, LanguageSelector } from "$lib";
  import { DarkMode } from "flowbite-svelte";
  import BellIcon from "$lib/icon/bell-01.svg";
  import StarsIcon from "$lib/icon/stars-01.svg";

  let { activeUrl = "/", changeShowAI = () => {} } = $props();
</script>

<div>
  <div class="flex justify-between px-4 lg:pt-8">
    <Breadcrumb {activeUrl} />
    <div class="flex space-x-4">
      <LanguageSelector />
      <DarkMode
        btnClass="text-gray-500 dark:text-gray-400 bg-white dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none rounded-lg p-2.5 border border-gray-300 dark:border-gray-600 dark:hover:border-gray-600 focus-within:ring-4 focus-within:outline-none focus-within:ring-gray-200 dark:focus-within:ring-gray-700"
      />
      <Button color="light" icon={BellIcon} classesSvg="bell" spaceY="mt-0" />
      <Button
        color="primary"
        icon={StarsIcon}
        spaceY="mt-0"
        funcButton={() => {
          changeShowAI();
        }}
      />
    </div>
  </div>
</div>
