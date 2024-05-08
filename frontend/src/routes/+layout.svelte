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

  let innerWidth = $state(0);
  let { showAI = false } = $props();
  let isShowAI = $state(false);
  let activities = $state([]);
  let isActiveUrl = $state($page.url.pathname);
  let activeUrl = $derived(isActiveUrl);
  let appUrl = ""; // change to application url to make the breadcrumb not display the application name

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

  const call = new CR3Call(globalThis.location.origin); // for frontend comment this

  const hideLayoutList = [
    "/forgot-password",
    "/workflow/update",
    "/workflow/create",
    "/hexabot/webchat",
  ];

  const hideLayout = () => {
    return !hideLayoutList.some((path) => {
      if (path.endsWith("/*")) {
        return globalThis.location.pathname.startsWith(path.slice(0, -1)); // for frontend comment this
        // return window.location.pathname.startsWith(path.slice(0, -1)); // for frontend uncomment this
      } else {
        return globalThis.location.pathname === path; // for frontend comment this
        // return window.location.pathname === path; // for frontend uncomment this
      }
    });
  };

  const changeShowAI = () => {
    showAI = !showAI;
    isShowAI = !isShowAI;
    localStorage.setItem("showAI", showAI);
  };

  const changeTabbar = (tabName) => {
    tabDatas.selected = tabName;
  };

  const numberMonthToTextMonth = (month) => {
    const monthNames = [
      "January",
      "February",
      "March",
      "April",
      "May",
      "June",
      "July",
      "August",
      "September",
      "October",
      "November",
      "December",
    ];
    return monthNames[month - 1];
  };

  const getModules = async () => {
    return [
      {
        icon: "users",
        name: "CRM",
        desc: "descCrm",
        status: "coming soon",
        children: [
          {
            href: await call.getPage("client-dashboard", "crm-overview"), // for frontend comment this
            // href: "/crm", // for frontend uncomment this
            name: "overview",
          },
          {
            href: await call.getPage("client-dashboard", "crm-person"), // for frontend comment this
            // href: "/crm/person", // for frontend uncomment this
            name: "persons",
          },
          {
            href: await call.getPage("client-dashboard", "crm-product-service"), // for frontend comment this
            // href: "/crm/product_service", // for frontend uncomment this
            name: "product&Services",
          },
          {
            href: await call.getPage("client-dashboard", "crm-users"), // for frontend comment this
            // href: "/crm/users", // for frontend uncomment this
            name: "users",
          },
          {
            href: await call.getPage("client-dashboard", "crm-categories"), // for frontend comment this
            // href: "/crm/categories", // for frontend uncomment this
            name: "categories",
          },
          {
            href: await call.getPage("client-dashboard", "crm-calendar"), // for frontend comment this
            // href: "/crm/calendar", // for frontend uncomment this
            name: "calendar",
          },
        ],
      },
      {
        icon: "users",
        name: "documentManagement",
        children: [
          {
            href: await call.getPage("client-dashboard", "doc-overview"), // for frontend comment this
            // href: "/document_management", // for frontend uncomment this
            name: "overview",
          },
        ],
      },
      // {
      //   href: await call.getPage("client-dashboard", "hexabot-overview"), // for frontend comment this
      //   name: "Hexa Bot",
      //   desc: "descHexaBots",
      //   status: "coming soon",
      //   children: [],
      // },
    ];
  };

  const logout = async () => {
    if (sessionStorage.getItem("user") !== null) sessionStorage.removeItem("user");
    if (sessionStorage.getItem("token") !== null) sessionStorage.removeItem("token");
    try {
      await call.goToPage("client-login", "home-page"); // for frontend comment this
    } catch (error) {
      console.error(error);
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
      })

      await logout();
    } catch (error) {
      console.error(error);
    }
  }

  export const load = async () => {
    if (browser) {
      locale.set(globalThis.navigator.language);
    }
    await waitLocale();
  };

  async function CheckSetup() {
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
  }

  onMount(() => {
    CheckSetup(); // for frontend comment this
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

  let userData = $state({ username: "", email: "", brand: "" });
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
        class="w-full min-h-screen flex lg:pl-80 border-l border-secondary-light dark:border-secondary-dark"
      >
        <!-- Navbar -->
        <Splitpanes class="[&>*]:!bg-primary-light [&>*]:dark:!bg-primary-dark [&>*]:!border-none">
          <Pane minSize={innerWidth > 1024 ? "70" : "0"}>
            <div class="transition-all duration-300">
              <Navbar {activeUrl} {showAI} {changeShowAI} {appUrl} reset_database={reset_database} />
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
