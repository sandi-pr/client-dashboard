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
  import { Typography, showAI } from "$lib";
  import { Drawer, Button, CloseButton, P } from "flowbite-svelte";
  import SidebarList from "./SidebarList.svelte";
  import LogoIcon from "$lib/icon/Logomark.svg";
  import MenuIcon from "$lib/icon/menu-02.svelte";
  import ArrowLeftIcon from "$lib/icon/arrow-left.svelte";
  import { sineIn } from "svelte/easing";
  import { page } from "$app/stores";
  import { onMount } from "svelte";
  let hidden2 = $state(true);
  let transitionParams = {
    x: -320,
    duration: 200,
    easing: sineIn,
  };

  let {
    getModule,
    toCamelCase,
    name = $bindable("HexaCode.AI"),
    email = $bindable(""),
    brand = $bindable(""),
    logoUrl = $bindable(""),
    innerWidth = 0,
    logout,
  }: {
    getModule: any;
    toCamelCase: any;
    name: string;
    email: string;
    brand: string;
    logoUrl: string;
    innerWidth: number;
    logout: () => void;
  } = $props();

  let url = $state("");
  let currentUrl = $state($page.url.pathname);

  $effect(() => {
    // console.log(currentUrl);
    if (currentUrl !== $page.url.pathname) {
      currentUrl = $page.url.pathname;
      // console.log("current url", currentUrl);
    }
  });
  onMount(() => {
    url = globalThis.location.pathname;
    // console.log("initial url", currentUrl);
  });
</script>

<!-- Sidebar for small screen -->
<Drawer
  transitionType="fly"
  {transitionParams}
  bind:hidden={hidden2}
  class="!p-0"
  divClass="overflow-y-hidden hover:overflow-y-auto !z-50 p-4 bg-white dark:bg-gray-800"
>
  <!-- <div class="fixed z-10 top-3 left-[280px]"> -->
  <!-- <div class="absolute z-50 top-3 start-[280px]">
    <CloseButton on:click={() => (hidden2 = true)} class=" dark:text-white" />
  </div> -->
  <SidebarList
    {getModule}
    {toCamelCase}
    {name}
    {email}
    {brand}
    {logoUrl}
    isHidden={false}
    funcButton={() => (hidden2 = true)}
    {logout}
    {innerWidth}
  />
</Drawer>

<!-- Sidebar for large screen -->
<SidebarList
  {getModule}
  {toCamelCase}
  {name}
  {email}
  {brand}
  {logoUrl}
  isHidden={innerWidth < 1024}
  funcButton={() => (hidden2 = true)}
  {logout}
  {innerWidth}
/>

<!-- Appbar for Drawer -->
{#if innerWidth < 1024 && !$showAI}
  <div
    class="lg:hidden fixed top-0 w-full z-40 px-2 py-1 md:py-3 justify-between items-center flex"
  >
    <!-- {#if currentUrl.startsWith('/client-dashboard/') && currentUrl.split('/').length > 3}
      <div class="p-2 rounded-lg justify-center items-center flex">
        <Button on:click={() => history.back()} color="light" class="p-2">
          <ArrowLeftIcon svgClass="w-6 h-6 text-fg-secondary-light dark:text-fg-secondary-dark" />
        </Button>
      </div>
    {/if} -->
    <div class="justify-between items-center gap-2 pl-2 flex">
      <div class="w-8 h-8 justify-center items-center flex">
        <img alt="logo" src={logoUrl !== "" ? logoUrl : LogoIcon} class="w-8 h-8" />
      </div>
      <Typography preset="Text lg/Bold" text={brand} color="primary" />
    </div>
    <div class="p-2 rounded-lg justify-center items-center flex">
      <Button on:click={() => (hidden2 = false)} color="light" class="p-2">
        <MenuIcon svgClass="w-6 h-6 text-fg-secondary-light dark:text-fg-secondary-dark" />
      </Button>
    </div>
  </div>
{/if}

<style>
  /* Hide scrollbar for Chrome, Safari and Opera */
  .no-scrollbar::-webkit-scrollbar {
    display: none;
  }

  /* Hide scrollbar for IE, Edge and Firefox */
  .no-scrollbar {
    -ms-overflow-style: none; /* IE and Edge */
    scrollbar-width: none; /* Firefox */
  }
</style>
