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
  import { Typography } from "$lib";
  import { Drawer, Button, CloseButton, P } from "flowbite-svelte";
  import SidebarList from "./SidebarList.svelte";
  import LogoIcon from "$lib/icon/Logomark.svg";
  import MenuIcon from "$lib/icon/menu-02.svelte";
  import { sineIn } from "svelte/easing";
  let hidden2 = $state(true);
  let transitionParams = {
    x: -320,
    duration: 200,
    easing: sineIn
  };

  let { getModule, toCamelCase, name = $bindable("Hexapod"), email = $bindable(""), brand = $bindable(""), innerWidth = 0, logout }: { getModule: any; toCamelCase: any; name: string; email: string; brand: string; innerWidth: number; logout: () => void } = $props();
</script>

<!-- Sidebar for small screen -->
<Drawer transitionType="fly" {transitionParams} bind:hidden={hidden2} class="!p-0">
  <!-- <div class="fixed z-10 top-3 left-[280px]"> -->
  <!-- <div class="absolute z-50 top-3 start-[280px]">
    <CloseButton on:click={() => (hidden2 = true)} class=" dark:text-white" />
  </div> -->
  <SidebarList {getModule} {toCamelCase} {name} {email} {brand} isHidden={false} funcButton={() => (hidden2 = true)} {logout} />
</Drawer>

<!-- Sidebar for large screen -->
<SidebarList {getModule} {toCamelCase} {name} {email} {brand} isHidden={innerWidth < 1024} funcButton={() => (hidden2 = true)} {logout} />

<!-- Appbar for Drawer -->
{#if innerWidth < 1024}
  <div class="lg:hidden pl-4 pr-2 py-3 justify-between items-center flex">
    <div class="justify-between items-center gap-2 flex">
      <div class="w-8 h-8 justify-center items-center flex">
        <img alt="logo" src={LogoIcon} class="w-8 h-8" />
      </div>
      <Typography preset="Text lg/Bold" text={brand} color="primary" />
    </div>
    <div class="p-2 rounded-lg justify-center items-center flex">
      <Button on:click={() => (hidden2 = false)} color="light">
        <MenuIcon svgClass="w-6 h-6 text-fg-secondary-light dark:text-fg-secondary-dark" />
      </Button>
    </div>
  </div>
{/if}
