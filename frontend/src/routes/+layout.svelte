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
  import { CR3Call } from "hexacomponents-svelte";
  import { Pane, Splitpanes } from "svelte-splitpanes";
  import { Navbar, Sidebar, Typography, TabBar, Card } from "$lib";
  import { onMount } from "svelte";
  import { locale, waitLocale, isLoading, _ } from "svelte-i18n";
  import "$lib/service/i18n";
  import SideDrawer from "$lib/components/SideDrawer.svelte";
  import LoadingSpinner from "$lib/components/LoadingSpinner.svelte";

  let innerWidth = $state(0);
  let { showAI = false } = $props();
  let isShowAI = $state(false);
  let activities = $state([]);
  // let isActiveUrl = $state($page.url.pathname);
  let activeUrl = $state($page.url.pathname);
  let appUrl = "/"; // change to application url to make the breadcrumb not display the application name
  let userData = $state({ username: "", email: "", brand: "" });
  let isFrontend = $state(false); // for frontend change to true

  let tabDatas = $state({
    selected: "Latest Activity",
    data: [
      {
        name: "Latest Activity",
      },
      {
        name: "Hexa AI",
      },
    ],
  });

  const call = isFrontend ? "" : new CR3Call(globalThis.location.origin);

  const hideLayoutList = [
    "/forgot-password",
    "/workflow/update",
    "/workflow/create",
    "/hexabot/webchat",
  ];

  const hideLayout = () => {
    return !hideLayoutList.some((path) => {
      try {
        if (path.endsWith("/*")) {
          return globalThis.location.pathname.startsWith(path.slice(0, -1));
        } else {
          return globalThis.location.pathname === path;
        }
      } catch (error) {
        console.log(error);
        if (path.endsWith("/*")) {
          return window.location.pathname.startsWith(path.slice(0, -1));
        } else {
          return window.location.pathname === path;
        }
      }
    });
  };

  const changeShowAI = () => {
    showAI = !showAI;
    isShowAI = !isShowAI;
    localStorage.setItem("showAI", showAI);
  };

  const getModules = async () => {
    return [
      {
        icon: "users",
        name: "crmOverview",
        desc: "descCrm",
        status: "coming soon",
        children: [
          {
            href: isFrontend ? "/crm" : await call.getPage("client-dashboard", "crm-overview"),
            name: "overview",
          },
          {
            href: isFrontend ? "/crm/person" : await call.getPage("client-dashboard", "crm-person"),
            name: "persons",
          },
          {
            href: isFrontend
              ? "/crm/product_service"
              : await call.getPage("client-dashboard", "crm-product-service"),
            name: "product&Services",
          },
          {
            href: isFrontend
              ? "/crm/categories"
              : await call.getPage("client-dashboard", "crm-categories"),
            name: "categories",
          },
          {
            href: isFrontend
              ? "/crm/calendar"
              : await call.getPage("client-dashboard", "crm-calendar"),
            name: "calendar",
          },
        ],
      },
      {
        icon: "fileShield",
        name: "documentManagement",
        children: [
          {
            href: isFrontend
              ? "/document_management"
              : await call.getPage("client-dashboard", "document-management-overview"),
            name: "overview",
          },
        ],
      },
      {
        icon: "userSquare",
        name: "HR",
        children: [
          {
            href: isFrontend
              ? "/hr"
              : await call.getPage("client-dashboard", "hr-overview"),
            name: "overview",
          },
        ],
      },
    ];
  };

  const logout = async () => {
    sessionStorage.clear();
    localStorage.clear();

    try {
      if (!isFrontend) {
        await call.goToPage("client-login", "home-page");
      }
    } catch (error) {
      console.log(error);
    }
  };

  const reset_database = async () => {
    try {
      await call.call("client-dashboard", "reset", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + sessionStorage.getItem("token"),
        },
      });

      await logout();
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

  async function CheckSetup() {
    try {
      var resp_setup = await call.call("client-login", "is-need-setup", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {},
      });
      console.log(resp_setup);
      if (resp_setup.status === 200) {
        if (resp_setup.result.is_need_setup === true) {
          if (sessionStorage.getItem("token") !== null) sessionStorage.removeItem("token");
          if (sessionStorage.getItem("user") !== null) sessionStorage.removeItem("user");
          await call.goToPage("hexapod-setup", "home-page");
          return;
        }
      }
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
            sessionStorage.setItem("user", JSON.stringify(resp.result));
            userData = resp.result;
          }
        }
      } else {
        await call.goToPage("client-login", "home-page");
        return;
      }
    } catch (error) {
      // console.log(error);
      userData = { username: "Hexa Staff", email: "hexacode@hexacode.org", brand: "hexacode" };
    }
  }

  onMount(() => {
    CheckSetup();
    if (localStorage.getItem("lang") !== null) {
      locale.set(localStorage.getItem("lang"));
    }

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
  if (typeof localStorage !== "undefined") {
    showAI = localStorage.getItem("showAI") === "true" ? true : false;
  } else {
    console.log("localstorage is undefined");
  }

  $effect(() => {
    if (activeUrl !== $page.url.pathname) {
      activeUrl = $page.url.pathname;
    }
  });
</script>

{#if $isLoading}
  <!-- <LoadingSpinner /> -->
{:else if winInnerWidth > 0}
  {#if hideLayout()}
    <div class="lg:flex">
      <Sidebar
        getModule={getModules}
        bind:name={userData.username}
        bind:email={userData.email}
        bind:brand={userData.brand}
        bind:isFrontend
        innerWidth={winInnerWidth}
        {logout}
      />
      <div
        class="w-full min-h-screen flex lg:pl-80 border-l border-secondary-light dark:border-secondary-dark"
      >
        <!-- Navbar -->
        <Splitpanes class="[&>*]:!bg-primary-light [&>*]:dark:!bg-primary-dark [&>*]:!border-none">
          <Pane minSize={innerWidth > 1024 ? "70" : "0"}>
            <div class="transition-all duration-300">
              <Navbar {activeUrl} {showAI} {changeShowAI} {appUrl} {reset_database} />
              <slot />
            </div>
          </Pane>
          <!-- Hexabot Sidebar -->
          {#if showAI}
            <SideDrawer innerWidth={winInnerWidth} {changeShowAI} />
          {/if}
        </Splitpanes>
      </div>
    </div>
  {:else}
    <slot />
  {/if}
{/if}
