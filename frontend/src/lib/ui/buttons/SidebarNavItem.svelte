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
  import { goto } from "$app/navigation";
  import { Typography, getIcons, Badge } from "$lib";
  import { page } from "$app/stores";
  import { _ } from "svelte-i18n";
  import { sidebarLargeClosed } from "$lib/utils/sidebarLargeClosed";
  import { onMount } from "svelte";

  let {
    href = "/",
    activeUrl = "/",
    commingSoon = false,
    classes = "flex w-full gap-3 items-center justify-between px-3 py-2 hover:no-underline rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700",
    classesActive = "flex w-full gap-3 items-center px-3 py-2 hover:no-underline bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700",
    icon = "",
    classesSvg = "w-5 h-5 flex text-hc-graylight-500 dark:text-hc-graydark-400",
    typographyData = {},
    collapse = false,
  } = $props();

  function extractPath(url) {
    if (!url) return;
    const match = url.match(/^(?:https?:\/\/)?[^\/]+(\/[^?#]*)/);
    return match ? match[1] : url;
  }

  let activeUrl2 = $derived($page.url.pathname);

  onMount(() => {
    console.log(`cek comming soon ${icon}`, commingSoon);
  })
</script>

<button
  class={!$sidebarLargeClosed
    ? activeUrl2 === extractPath(href)
      ? classesActive
      : classes
    : "flex w-full py-2 m-0 !my-4 hover:no-underline rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 justify-center"}
  disabled={commingSoon}
  onclick={() => {
    // console.log(href);
     goto(href);
  }}
>
  <div class="flex gap-3">
    <div class="w-max">
      {#if icon}
        <svelte:component this={getIcons(icon)} svgClass={classesSvg} />
      {/if}
    </div>
    <Typography {...typographyData} />
  </div>
  {#if commingSoon}
    <Badge text={$_("soon")} type="brand" />
  {/if}
  <!-- <Typography preset={typographyData.preset} text={typographyData.text} color={typographyData.color} /> -->
</button>

<!--
@component
## Config
@prop label: string; // Label of the item.
@prop href: string; // URL of the item.
@prop activeUrl: string; // Current active URL.
@prop classes: string; // Class of the item.
@prop classesActive: string; // Class of the item when active.
@prop icon: string; // Icon of the item.
@prop classesSvg: string; // Class of the icon.
@prop typographyData: object; // Data of Typography object.
-->
