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
  import { locale, waitLocale, isLoading } from "svelte-i18n";
  import "$lib/service/i18n";
  import SideDrawer from "$lib/components/SideDrawer.svelte";

  let innerWidth = $state(0);
  let showAI = $state(false);
  let activities = $state([]);
  let activeUrl = $derived($page.url.pathname);

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
            icon: "users",
            name: "Overview",
            desc: "descCrmOverview",
            status: "coming soon",
          },
          {
            href: await call.getPage("client-dashboard", "crm-person"), // for frontend comment this
            icon: "users",
            name: "Person",
            desc: "descCrmPerson",
            status: "coming soon",
          },
          {
            href: await call.getPage("client-dashboard", "crm-calendar"), // for frontend comment this
            icon: "users",
            name: "Calendar",
            desc: "descCrmCalendar",
            status: "coming soon",
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
    CheckSetup();
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
        class="w-full min-h-screen flex lg:ml-80 border-l border-secondary-light dark:border-secondary-dark"
      >
        <!-- Navbar -->
        <Splitpanes class="[&>*]:!bg-primary-light [&>*]:dark:!bg-primary-dark [&>*]:!border-none">
          <Pane>
            <div class="transition-all duration-300">
              <Navbar {activeUrl} {showAI} {changeShowAI} />
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
