<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: Azzis Arswendo <azzis@hexacode.org>

Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->

<script>
  import "../app.css";
  import "@fontsource/plus-jakarta-sans";
  import { Sidebar, Navbar } from "hexacomponents-svelte";
  import { onMount } from "svelte";
  import { locale, waitLocale, isLoading } from "svelte-i18n";
  import "$lib/service/i18n";
  import checkIsLoggedIn from "$lib/utils/checkIsLoggedIn";

  let innerWidth = $state(0);

  const hideLayoutList = ["/login", "/register", "/forgot-password", "/auth"];

  const hideLayout = () => {
    return !hideLayoutList.includes(window.location.pathname);
  };

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

  const logout = () => {
    sessionStorage.removeItem("user");
    window.location.replace("/login");
  };

  export const load = async () => {
    if (browser) {
      locale.set(window.navigator.language);
    }
    await waitLocale();
  };

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

  let winInnerWidth = $derived(innerWidth);
  let userData = $state({});

  $effect(() => {
    if (hideLayout() && checkIsLoggedIn() === null) {
      window.location.replace("/login");
    } else if (!hideLayout() && checkIsLoggedIn() !== null) {
      window.location.replace("/");
    }
    userData = checkIsLoggedIn();
  });
</script>

{#if $isLoading}
  <div>Loading...</div>
{:else if winInnerWidth > 0}
  {#if hideLayout()}
    <div class="lg:flex">
      <Sidebar
        getModule={getModules}
        name={userData?.data?.username || "Hexastaff"}
        email={userData?.data?.email || "staff@hexacode.org"}
        brand={userData?.data?.brand || "Crossroads"}
        innerWidth={winInnerWidth}
        {logout}
      />
      <div
        class="w-full h-screen lg:ml-80 pt-8 border-l border-secondary-light dark:border-secondary-dark"
      >
        <Navbar />
        <slot />
      </div>
    </div>
  {:else}
    <slot />
  {/if}
{/if}
