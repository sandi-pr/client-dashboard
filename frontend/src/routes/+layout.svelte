<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: Azzis Arswendo <azzis@hexacode.org>

Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->

<script>
  import "../app.css";
  import { browser } from "$app/environment";
  import { Sidebar, Navbar } from "hexacomponents-svelte";
  import { onMount, onDestroy } from "svelte";
  // import { getModules } from "$lib/utils/getModule";

  let innerWidth = $state(0);

  // if (typeof window !== "undefined") {
  //   innerHeight = window.innerHeight;
  //   innerWidth = window.innerWidth;
  //   console.log(window.innerHeight, window.innerWidth);
  // }

  onMount(() => {
    const resizeHandler = () => {
      innerWidth = window.innerWidth;
    };
    window.addEventListener("resize", resizeHandler);
    resizeHandler(); // call once initially to set the correct value
    return () => {
      window.removeEventListener("resize", resizeHandler);
    };
  });

  let winInnerHeight = $derived(innerHeight);
  let winInnerWidth = $derived(innerWidth);
  // import "$lib/service/i18n.js";
  // import { locale, waitLocale, isLoading } from "svelte-i18n";

  // export const load = async () => {
  //   if (browser) {
  //     locale.set(window.navigator.language);
  //   }
  //   await waitLocale();
  // };
  const getModules = async () => {
    return [
      {
        href: "/person",
        icon: "users",
        name: "CRM",
        desc: "descCrm",
        status: "coming soon",
      },
      {
        href: "/hexa-bots",
        icon: "messageChat",
        name: "Hexa Bots",
        desc: "descHexaBots",
        status: "coming soon",
      },
    ];
  };
</script>

<!-- {#if $isLoading}
    <p>Loading...</p>
  {:else} -->
{#if winInnerWidth > 0}
  <div class="lg:flex">
    <Sidebar
      getModule={getModules}
      name="Hexastaff"
      email="staff@hexacode.org"
      brand="Crossroads"
      innerWidth={winInnerWidth}
    />
    <div class="w-full h-screen lg:ml-80">
      <Navbar />
      <slot />
    </div>
  </div>
{/if}
<!-- {/if} -->
