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
  import { page } from "$app/stores";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { Navbar, Sidebar } from "$lib";
  import { onMount } from "svelte";
  import { locale, waitLocale, isLoading } from "svelte-i18n";
  import "$lib/service/i18n";
  import checkIsLoggedIn from "$lib/utils/checkIsLoggedIn";

  let innerWidth = $state(0);
  let activeUrl = $derived($page.url.pathname);
  // const call = new CR3Call(globalThis.location.origin);

  const hideLayoutList = ["/forgot-password", "/workflow/update", "/workflow/create"];

  const hideLayout = () => {
    return !hideLayoutList.some((path) => {
      if (path.endsWith("/*")) {
        return globalThis.location.pathname.startsWith(path.slice(0, -1));
      } else {
        return globalThis.location.pathname === path;
      }
    });
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

  const logout = async () => {
    sessionStorage.removeItem("user");
    sessionStorage.removeItem("token");
    try {
      await call.goToPage("client-login", "home-page");
    } catch (error) {
      console.error(error);
    }
  };

  export const load = async () => {
    if (browser) {
      locale.set(globalThis.navigator.language);
    }
    await waitLocale();
  };

  onMount(() => {
    const resizeHandler = () => {
      innerWidth = globalThis.innerWidth;
    };
    globalThis.addEventListener("resize", resizeHandler);
    resizeHandler(); // call once initially to set the correct value
    return () => {
      globalThis.removeEventListener("resize", resizeHandler);
    };
  });

  let winInnerWidth = $derived(innerWidth);
  let userData = $state({ username: "", email: "", brand: "" });

  $effect(async () => {
    var token = sessionStorage.getItem("token");
    if (token !== null) {
      var resp = await call.call("client-login", "is-token-valid", {
        URLQuery: {},
        URLParams: {},
        Params: {
          token: token,
        },
        Headers: {},
      });

      if (resp.status === 200) {
        if (resp.result.confirmed === true && resp.result.blocked === false) {
          sessionStorage.setItem("user", resp.result);
          userData = resp.result;
        }
      }
    }
  });
</script>

{#if $isLoading}
  <div>Loading...</div>
{:else if winInnerWidth > 0}
  {#if hideLayout()}
    <div class="lg:flex">
      <Sidebar
        getModule={getModules}
        bind:name={userData.username}
        bind:email={userData.email}
        bind:brand={userData.brand}
        innerWidth={winInnerWidth}
        {logout}
      />
      <div
        class="w-full h-screen lg:ml-80 border-l border-secondary-light dark:border-secondary-dark"
      >
        <Navbar {activeUrl} />
        <slot />
      </div>
    </div>
  {:else}
    <slot />
  {/if}
{/if}
