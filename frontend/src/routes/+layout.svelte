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
  import { Navbar, Sidebar, Typography, TabBar, Card } from "$lib";
  import { CloseButton } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { locale, waitLocale, isLoading } from "svelte-i18n";
  import DataflowIcon from "$lib/icon/dataflow-solid.svelte";
  import HexaAi from "$lib/icon/hexa-ai.svelte";
  import "$lib/service/i18n";
  import checkIsLoggedIn from "$lib/utils/checkIsLoggedIn";

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

  const changeLanguage = (lang) => {
    locale.set(lang);
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
        href: await call.getPage("client-dashboard", "crm-overview"), // for frontend comment this
        href: "crm/person",
        icon: "users",
        name: "CRM",
        desc: "descCrm",
        status: "coming soon",
      },
      {
        href: await call.getPage("client-dashboard", "hexabot-overview"), // for frontend comment this
        href: "hexabot",
        icon: "messageChat",
        name: "Hexa Bot",
        desc: "descHexaBots",
        status: "coming soon",
      },
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

  // for frontend comment below
  $effect(async () => {
    var resp_setup = await call.call("client-login", "is-need-setup", {
      URLQuery: {},
      URLParams: {},
      Params: {},
      Headers: {},
    });
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
        class="w-full h-screen flex lg:ml-80 border-l border-secondary-light dark:border-secondary-dark"
      >
        <!-- Navbar -->
        <div class={`transition-all duration-300 ${showAI ? "w-[80%]" : "w-full"}`}>
          <Navbar {activeUrl} {changeLanguage} {showAI} {changeShowAI} />
          <slot />
        </div>

        <!-- Hexabot Sidebar -->
        {#if showAI}
          <div
            class="w-[20%] h-screen transition-all duration-300 border-l-2 border-gray-200 dark:border-gray-800"
          >
            <div class="flex items-center">
              <div
                id="drawer-label"
                class="inline-flex items-center m-4 text-base font-semibold text-gray-500 dark:text-gray-400"
              >
                {#if tabDatas.selected === "Latest Activity"}
                  <div
                    class="w-10 h-10 p-2.5 mr-3 bg-brand-secondary-light dark:bg-brand-secondary-dark rounded-full justify-center items-center"
                  >
                    <svelte:component
                      this={DataflowIcon}
                      svgClass="text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark"
                    />
                  </div>
                  &nbsp;
                  <div>
                    <Typography preset="Text xl/Semibold" text={"latestActivity"} color="primary" />
                    <Typography
                      preset="Text sm/Regular"
                      text={"yourLatestActivity"}
                      color="tertiary"
                    />
                  </div>
                {:else if tabDatas.selected === "Hexa AI"}
                  <div
                    class="w-10 h-10 p-2.5 mr-3 bg-brand-secondary-light dark:bg-brand-secondary-dark rounded-full justify-center items-center"
                  >
                    <svelte:component
                      this={HexaAi}
                      svgClass="text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark"
                    />
                  </div>
                  &nbsp;
                  <div>
                    <Typography preset="Text xl/Semibold" text="Hexa AI" color="primary" />
                    <Typography
                      preset="Text sm/Regular"
                      text={"yourCrossroadsPersonalAssistant"}
                      color="tertiary"
                    />
                  </div>
                {/if}
              </div>
              <CloseButton
                on:click={() => {
                  showAI = false;
                }}
                class="mb-4 dark:text-white"
              />
            </div>

            <TabBar {tabDatas} {changeTabbar} />

            {#if tabDatas.selected === "Latest Activity"}
              <div class="p-4">
                {#if activities.length > 0}
                  {#each activities as activityYear (activityYear)}
                    {#each Object.keys(activityYear) as year (year)}
                      <div class="flex flex-col gap-6">
                        {#each Object.keys(activityYear[year]) as month (month)}
                          <div class="flex flex-col gap-3">
                            <Typography
                              preset="Text sm/Medium"
                              text="{numberMonthToTextMonth(month)} {year}"
                              color="tertiary"
                            />
                            {#each activityYear[year][month] as activity (activity)}
                              <Card
                                type="ledger"
                                title={activity.activity}
                                description={activity.address}
                                subtext={activity.date}
                              />
                            {/each}
                          </div>
                        {/each}
                      </div>
                    {/each}
                  {/each}
                {:else}
                  <Typography preset="Text sm/Medium" text="No Activity" color="primary" />
                {/if}
              </div>
            {:else if tabDatas.selected === "Hexa AI"}
              <div class="w-full h-full m-4">
                <!-- <AIBots /> -->
              </div>
            {:else}
              ""
            {/if}
          </div>
        {/if}
      </div>
    </div>
  {:else}
    <slot />
  {/if}
{/if}
