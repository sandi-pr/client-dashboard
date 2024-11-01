<!--
  CR3 Client Dashboard
  Package: gitlab.hexacode.org/public-repo/js/hexacomponents-svelte
  Maintainer:
  * Sandi Permana <sandy@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>

  Copyright (C) 2024 Hexacode Teknologi Indonesia
  All Rights Reserved.
-->

<script>
  import { page } from "$app/stores";
  import { LogoBrand, SidebarNavItem, Typography, getIcons, LanguageSelector, toCamelCase } from "$lib";
  import { goto } from "$app/navigation";
  import { CR3Call } from "hexacomponents-svelte";
  import {
    Search,
    Button,
    Avatar,
    Sidebar,
    SidebarGroup,
    SidebarWrapper,
    SidebarDropdownWrapper,
    SidebarDropdownItem,
    DarkMode,
  } from "flowbite-svelte";
  import {
    CogOutline,
    UsersOutline,
    ChevronUpOutline,
    ChevronDownOutline,
    ArrowRightToBracketOutline,
  } from "flowbite-svelte-icons";
  import LogoIcon from "$lib/icon/Logomark.svg";
  import { sidebarLargeClosed, toggleSidebarLargeClosed } from "$lib/utils/sidebarLargeClosed";
  import SidebarCloseIcon from "$lib/icon/sidebar-close.svelte";
  import SidebarOpenIcon from "$lib/icon/sidebar-open.svelte";
  import { onMount } from "svelte";
  import Me from "$lib/utils/me";
  import { _ } from "svelte-i18n";

  let {
    getModule,
    name = $bindable(""),
    email = $bindable(""),
    brand = $bindable(""),
    logoUrl = $bindable(""),
    isHidden = true,
    isFrontend = $bindable(false),
    funcButton,
    logout,
    innerWidth = 0,
  } = $props();

  let sidebarItems = $state([]);
  let dropdownUser = $state(false);
  let dropdownManage = $state(false);
  let isOpen = $state({});
  let isOpenOld = $state({});

  let hrefUser = $state("");
  let hrefJobLevels = $state("");

  import { call } from "$lib";

  const sidebarBuilder = async () => {
    let moduleClient = await getModule();
    const sidebarHalf = moduleClient.map((item) => {
      return {
        name: item.name,
        href: item.href,
        icon: item.icon,
        children: item.children,
      };
    });

    // sidebarHalf.push({ name: "Digital Identity", href: "/blockchain", icon: "blockchain" });

    return sidebarHalf;
  };

  let activeUrl = $derived($page.url.pathname);

  const handleLogout = () => {
    logout();
  };

  const changeSidbar = () => {
    toggleSidebarLargeClosed(!$sidebarLargeClosed);
    console.log($sidebarLargeClosed);
    localStorage.setItem("sidebarLargeClosed", $sidebarLargeClosed ? "true" : "false");
  };

  onMount(async () => {
    const sidebarAddition = await sidebarBuilder();

    hrefUser = await $call.getPage("client-dashboard", "users");
    hrefJobLevels = await $call.getPage("client-dashboard", "job-management");

    let home = await $call.getPage("client-dashboard", "home-page");
    sidebarItems = [
      { name: "home", href: home, icon: "home", children: [] },
      ...sidebarAddition,
    ];

    // Inisialisasi isOpen dari sessionStorage atau buat baru jika tidak ada
    let isOpenSession = JSON.parse(sessionStorage.getItem("isOpen"));
    if (isOpenSession !== null) {
      isOpen = isOpenSession;
    } else {
      for (const item of sidebarItems) {
        isOpen[item.name] = false;
      }
      sessionStorage.setItem("isOpen", JSON.stringify(isOpen));
    }
  });

  $effect(() => {
    if (JSON.stringify(isOpen) !== JSON.stringify(isOpenOld)) {
      isOpenOld = { ...isOpen };
      sessionStorage.setItem("isOpen", JSON.stringify(isOpen));
    }
  });
</script>

<Sidebar
  {activeUrl}
  class={` ${isHidden ? "hidden h-screen" : "fixed h-full"} !overflow-x-visible z-50 border-r border-secondary-light dark:border-secondary-dark lg:fixed lg:inset-y-0 lg:z-10 lg:flex ${$sidebarLargeClosed ? "w-20" : "w-80"} transition-all duration-300 ease-in-out`}
>
  <SidebarWrapper
    class="h-full lg:pt-6 lg:px-4 flex flex-col grow justify-between bg-white dark:bg-primary-dark"
    divClass="overflow-y-hidden hover:overflow-y-auto no-scrollbar py-4 px-3 bg-gray-50 rounded dark:bg-gray-800"
  >
    <SidebarGroup>
      <div
        class={`relative flex items-start ${!$sidebarLargeClosed ? "justify-between" : "justify-center"}`}
      >
        {#if isHidden}
          <div></div>
        {:else if !$sidebarLargeClosed}
          <LogoBrand
            brandName={brand !== "" ? brand : "Hexapod"}
            href="/"
            brandLogo={logoUrl !== "" ? logoUrl : LogoIcon}
            imgClass="h-8"
          />
          {#if !(innerWidth < 1024)}
            <button
              onclick={() => {
                changeSidbar();
              }}
            >
              <SidebarCloseIcon
                svgClass="w-8 h-8 cursor-pointer text-utility-brand-500-alt-light dark:text-utility-brand-500-alt-dark"
              />
            </button>
          {/if}
        {:else}
          <LogoBrand
            brandName={""}
            href="/"
            brandLogo={logoUrl !== "" ? logoUrl : LogoIcon}
            imgClass="h-8"
          />
          {#if !(innerWidth < 1024)}
            <button
              onclick={() => {
                changeSidbar();
              }}
              class="fixed left-16 z-50 overflow-visible"
            >
              <SidebarOpenIcon
                svgClass="w-8 h-8 cursor-pointer text-utility-brand-500-alt-light dark:text-utility-brand-500-alt-dark"
              />
            </button>
          {/if}
        {/if}
      </div>
      {#if !$sidebarLargeClosed}
        <div class="md:hidden flex gap-2 pb-1">
          <div class="grow w-full">
            <Search placeholder={$_("search")} class="bg-white dark:bg-primary-dark !py-2" />
          </div>
          <DarkMode
            btnClass="text-gray-500 h-max dark:text-gray-400 bg-white dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none rounded-lg p-2.5 border border-gray-300 dark:border-gray-600 dark:hover:border-gray-600 focus-within:ring-4 focus-within:outline-none focus-within:ring-gray-200 dark:focus-within:ring-gray-700"
          />
          <LanguageSelector isMobileScreen />
        </div>
        <div class="hidden md:flex pb-4">
          <Search placeholder={$_("search")} class="bg-white dark:bg-primary-dark" />
        </div>
        {#each sidebarItems as item}
          {#if item.children.length > 0}
            <SidebarDropdownWrapper
              bind:isOpen={isOpen[item.name]}
              ulClass="!space-y-0"
              spanClass="flex-1 text-left whitespace-nowrap text-secondary-light dark:text-secondary-dark text-base font-semibold leading-normal"
              btnClass="flex w-full gap-3 items-center px-3 py-2 hover:no-underline rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
              label={$_(item.name)}
            >
              <svelte:fragment slot="icon">
                <svelte:component
                  this={getIcons(item.icon)}
                  svgClass="w-5 h-5 text-hc-graylight-500 dark:text-hc-graydark-400"
                />
              </svelte:fragment>
              {#each item.children as child}
                <SidebarNavItem
                  rel="external"
                  href={child.href}
                  {activeUrl}
                  icon={child.icon}
                  classes="flex w-full gap-3 items-center pl-11 py-2 hover:no-underline rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
                  classesActive="flex w-full gap-3 items-center pl-11 py-2 hover:no-underline bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
                  typographyData={{
                    preset: "Text md/Regular",
                    text: $_(toCamelCase(child.name)),
                    color: "secondary",
                  }}
                  collapse={false}
                />
              {/each}
            </SidebarDropdownWrapper>
          {:else}
            <SidebarNavItem
              href={item.href}
              {activeUrl}
              icon={item.icon}
              typographyData={{
                preset: "Text md/Semibold",
                text: $_(item.name),
                color: "secondary",
              }}
              collapse={false}
            />
          {/if}
        {/each}
      {:else}
        <Search placeholder={$_("")} class="bg-white dark:bg-primary-dark mb-6 w-0 px-0" />
        {#each sidebarItems as item}
          <SidebarNavItem
            href={item.children.length > 0 ? item.children[0].href : item.href}
            {activeUrl}
            icon={item.icon}
            typographyData={{
              preset: "Text md/Semibold",
              color: "secondary",
            }}
            collapse={true}
          />
        {/each}
      {/if}
    </SidebarGroup>
    <div>
      <!-- <SidebarGroup>
        <SidebarNavItem
          href="/settings"
          {activeUrl}
          icon="setting"
          typographyData={{ preset: "Text md/Semibold", text: $_("settings"), color: "secondary" }}
        />
      </SidebarGroup> -->
      <div class="relative">
        <!-- Dropdown menu user -->
        <div
          class={`z-10 ${dropdownUser ? "fixed" : "hidden"} bottom-16 p-2 mb-2 bg-white dark:bg-gray-700 divide-y divide-gray-100 rounded-lg border border-secondary-light dark:border-secondary-dark shadow dark:shadow-lg`}
        >
          <ul aria-labelledby="dropdownDefaultButton">
            <div class="flex pb-3 justify-between rounded-lg border-b">
              <div class="flex w-full h-full space-x-4 items-center cursor-pointer">
                <Avatar />
                <div class="grow">
                  <Typography preset="Text md/Semibold" text={name} color="secondary" />
                  <Typography preset="Text sm/Regular" text={email} color="tertiary" />
                </div>
              </div>
              <Button
                on:click={handleLogout}
                color="light"
                class="focus-within:ring-0 border-0 p-2 bg-transparent dark:bg-transparent hover:bg-secondary-light dark:hover:bg-secondary-dark"
              >
                <ArrowRightToBracketOutline
                  class="size-7 text-fg-quaternary-light dark:text-fg-quaternary-dark transition duration-75 group-hover:text-gray-900 dark:group-hover:text-white"
                />
              </Button>
            </div>
            <SidebarDropdownWrapper
              spanClass="flex-1 ms-0 text-left whitespace-nowrap"
              label={$_("manage")}
            >
              <svelte:fragment slot="icon">
                <UsersOutline class="size-6 mr-4 text-secondary-light dark:text-secondary-dark" />
              </svelte:fragment>
              <SidebarDropdownItem href={hrefUser} class="!ps-14" label={$_("users")} />
              <!-- <SidebarDropdownItem href={hrefJobLevels} class="!ps-14" label={$_("job")} /> -->
            </SidebarDropdownWrapper>
            <button
              onclick={async () => {
                let url = await $call.getPage("client-dashboard", "settings");
                goto(url);
                // else goto("/settings");
              }}
              class="flex p-2 items-center w-full hover:bg-secondary-light dark:hover:bg-secondary-dark rounded-lg"
            >
              <CogOutline class="size-6 mr-4 text-secondary-light dark:text-secondary-dark" />
              <Typography
                preset="Text md/Regular"
                propClass="text-base font-normal"
                text={$_("settings")}
                color="primary"
              />
            </button>
          </ul>
        </div>
        <SidebarGroup>
          <button
            type="button"
            class="flex w-full justify-between rounded-lg hover:bg-secondary-light dark:hover:bg-secondary-dark"
            onclick={() => {
              dropdownUser = !dropdownUser;
            }}
          >
            <div
              class={`flex w-full h-full ${$sidebarLargeClosed ? "" : "gap-4"} items-center cursor-pointer`}
            >
              <Avatar />
              {#if !$sidebarLargeClosed}
                <div class="grow">
                  <Typography preset="Text md/Semibold" text={name} color="secondary" />
                  <Typography preset="Text sm/Regular" text={email} color="tertiary" />
                </div>
                <Button
                  color="light"
                  class="focus-within:ring-0 border-0 p-2 bg-transparent dark:bg-transparent hover:bg-primary-light dark:hover:bg-primary-dark"
                >
                  {#if dropdownUser}
                    <ChevronUpOutline
                      class="size-7 text-fg-quaternary-light dark:text-fg-quaternary-dark transition duration-75 group-hover:text-gray-900 dark:group-hover:text-white"
                    />
                  {:else}
                    <ChevronDownOutline
                      class="size-7 text-fg-quaternary-light dark:text-fg-quaternary-dark transition duration-75 group-hover:text-gray-900 dark:group-hover:text-white"
                    />
                  {/if}
                </Button>
              {:else}
                <div></div>
              {/if}
            </div>
          </button>
        </SidebarGroup>
      </div>
    </div>
  </SidebarWrapper>
</Sidebar>

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
