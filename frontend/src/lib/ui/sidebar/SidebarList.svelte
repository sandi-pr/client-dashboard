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
  import { goto } from "$app/navigation";
  import {
    Search,
    Button,
    Avatar,
    Sidebar,
    SidebarGroup,
    SidebarWrapper,
    SidebarDropdownWrapper,
    SidebarDropdownItem,
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
  import Divider from "$lib/components/Divider.svelte";
  import { _ } from "svelte-i18n";
  import { LogoBrand, SidebarNavItem, Typography } from "$lib";
  import {
    me,
    Icon,
    call,
    check,
    getIcons,
    isDemoTrue,
    toCamelCase,
    isOpenSidebarStore,
  } from "$lib";

  const MessageTextSquare01Icon = getIcons("messageTextSquare01");

  const isFrontend = process.env.NODE_ENV === "development";

  let {
    getModule,
    name = $bindable(""),
    email = $bindable(""),
    brand = $bindable(""),
    logoUrl = $bindable(""),
    isHidden = true,
    logout,
    innerWidth = 0,
  } = $props();

  let dropdownUser = $state(false);
  let isInitialized = false; // Flag to ensure `onMount` runs first
  let hrefUser = $state("");
  let hrefJobLevels = $state("");
  let sidebarItems = $state([]);
  let isOpen = $state({});

  let activeUrl = $derived($page.url.pathname);

  const sidebarBuilder = async () => {
    let moduleClient = await getModule();
    const sidebarHalf = moduleClient.map((item) => {
      return {
        name: item.name,
        href: item.href,
        icon: item.icon,
        commingSoon: item.commingSoon || false,
        divider: item.divider || false,
        children: item.children,
      };
    });

    // sidebarHalf.push({ name: "Digital Identity", href: "/blockchain", icon: "blockchain" });

    return sidebarHalf;
  };

  const handleLogout = () => {
    logout();
  };

  const changeSidebar = () => {
    toggleSidebarLargeClosed(!$sidebarLargeClosed);
    // console.log($sidebarLargeClosed);
    localStorage.setItem("sidebarLargeClosed", $sidebarLargeClosed ? "true" : "false");
  };

  const handleDropdownFocusLoss = (event) => {
    if (!event.target.closest(".dropdown-container")) {
      dropdownUser = false;
    }
  };

  onMount(async () => {
    const sidebarAddition = await sidebarBuilder();

    hrefUser = isFrontend ? "/manage_users" : await $call.getPage("client-dashboard", "users");
    hrefJobLevels = await $call.getPage("client-dashboard", "job-management");

    let home = isFrontend ? "/" : await $call.getPage("client-dashboard", "home-page");
    sidebarItems = [{ name: "home", href: home, icon: "home", children: [] }, ...sidebarAddition];

    // Initialize `isOpen` from sessionStorage or create a new object if not present
    let isOpenSession = JSON.parse(sessionStorage.getItem("isOpen"));

    if (isOpenSession !== null) {
      // Ensure all keys from `sidebarItems` exist in `isOpenSession`
      sidebarItems.forEach((item) => {
        if (!(item.name in isOpenSession)) {
          isOpenSession[item.name] = false; // Default to closed if key is missing
        }
      });
      isOpenSidebarStore.set(isOpenSession); // Restore state from sessionStorage
    } else {
      // Create a new `isOpen` object with all keys from `sidebarItems`
      sidebarItems.forEach((item) => {
        isOpen[item.name] = false; // Default all items to closed
      });
      isOpenSidebarStore.set(isOpen); // Set the store
      sessionStorage.setItem("isOpen", JSON.stringify(isOpen)); // Persist the state
    }

    isInitialized = true; // Mark initialization as complete

    document.addEventListener("click", handleDropdownFocusLoss);
    return () => {
      document.removeEventListener("click", handleDropdownFocusLoss);
    };
  });

  $effect(() => {
    // Sync store changes to sessionStorage
    isOpenSidebarStore.subscribe((value) => {
      // Ensure all keys from `sidebarItems` exist in the store value
      sidebarItems.forEach((item) => {
        if (!(item.name in value)) {
          value[item.name] = false; // Default to closed if key is missing
        }
      });
      // Ensure `onMount` has completed before running this effect
      if (isInitialized) {
        sessionStorage.setItem("isOpen", JSON.stringify(value));
      }
    });
  });
</script>

<Sidebar
  {activeUrl}
  class={` ${isHidden ? "hidden h-screen" : "fixed h-full"} !overflow-x-visible z-50 border-r border-secondary-light dark:border-secondary-dark lg:fixed lg:inset-y-0 lg:z-10 lg:flex ${$sidebarLargeClosed ? "w-20" : "w-80"} transition-all duration-300 ease-in-out`}
>
  <SidebarWrapper
    class="h-full lg:pt-6 lg:px-4 pb-24 flex flex-col grow justify-between bg-white dark:bg-primary-dark"
    divClass="overflow-y-hidden hover:overflow-y-auto no-scrollbar py-4 px-3 bg-gray-50 rounded dark:bg-gray-800"
  >
    <SidebarGroup>
      <div class="sticky top-0 flex flex-col bg-primary-light dark:bg-primary-dark">
        <div
          class={`relative flex items-start ${!$sidebarLargeClosed ? "justify-between" : "justify-center"}`}
        >
          {#if isHidden}
            <div></div>
          {:else if !$sidebarLargeClosed}
            <LogoBrand
              brandName={brand !== "" ? brand : "HexaCode.AI"}
              href="/"
              brandLogo={logoUrl !== "" ? logoUrl : LogoIcon}
              imgClass="h-8"
            />
            {#if !(innerWidth < 1024)}
              <button
                onclick={() => {
                  changeSidebar();
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
                  changeSidebar();
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
        <div class="md:hidden flex gap-2 pb-1">
          <div class="grow w-full">
            <Search placeholder={$_("search")} class="bg-white dark:bg-primary-dark !py-2" />
          </div>
        </div>
        <div class="hidden md:flex pb-4">
          <Search
            on:click={() => {
              if ($sidebarLargeClosed) {
                changeSidebar();
              }
            }}
            placeholder={$_("search")}
            class={`bg-white dark:bg-primary-dark ${$sidebarLargeClosed ? "px-0" : ""}`}
          />
        </div>
      </div>
      {#if !$sidebarLargeClosed}
        {#each sidebarItems as item}
          {#if item?.children?.length > 0}
            <SidebarDropdownWrapper
              bind:isOpen={$isOpenSidebarStore[item.name]}
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
                  commingSoon={child.commingSoon}
                  classes="flex w-full gap-3 items-center justify-between pl-11 py-2 hover:no-underline rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
                  classesActive="flex w-full gap-3 items-center justify-between pl-11 py-2 hover:no-underline bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
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
              commingSoon={item.commingSoon}
              typographyData={{
                preset: "Text md/Semibold",
                text: $_(item.name),
                color: "secondary",
              }}
              collapse={false}
            />
          {/if}
          {#if item?.divider}
            <Divider />
          {/if}
        {/each}
      {:else}
        <!-- <Search placeholder={$_("")} class="bg-white dark:bg-primary-dark mb-6 w-0 px-0" /> -->
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
    <div
      class="fixed bottom-0 pb-8 bg-primary-light dark:bg-primary-dark max-w-72 {$sidebarLargeClosed
        ? ''
        : 'w-full'}"
    >
      <div class="relative dropdown-container">
        <div
          class={`z-10 ${dropdownUser ? "fixed" : "hidden"} bottom-20 p-2 mb-2 bg-white dark:bg-gray-700 divide-y divide-gray-100 rounded-lg border border-secondary-light dark:border-secondary-dark shadow dark:shadow-lg w-full max-w-72`}
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
                class="focus-within:ring-0 border-0 p-2 bg-transparent dark:bg-transparent hover:bg-secondary-light dark:hover:bg-gray-700"
              >
                <ArrowRightToBracketOutline
                  class="size-7 text-fg-quaternary-light dark:text-fg-quaternary-dark transition duration-75 group-hover:text-gray-900 dark:group-hover:text-white"
                />
              </Button>
            </div>
            <SidebarDropdownWrapper
              spanClass="flex-1 ms-0 text-left whitespace-nowrap"
              class="hover:bg-secondary-light dark:hover:bg-secondary-dark mt-2"
              label={$_("manage")}
              onclick={(event) => {
                event.stopPropagation(); // Prevent handleDropdownFocusLoss from being triggered
              }}
            >
              <svelte:fragment slot="icon">
                <UsersOutline class="size-6 mr-4 text-secondary-light dark:text-secondary-dark" />
              </svelte:fragment>
              <SidebarDropdownItem
                href={hrefUser}
                onclick={() => {
                  dropdownUser = !dropdownUser;
                }}
                class="!ps-14 hover:bg-secondary-light dark:hover:bg-secondary-dark"
                label={$_("users")}
              />
            </SidebarDropdownWrapper>
            <button
              onclick={async () => {
                dropdownUser = !dropdownUser;
                if ($me.authority.toLowerCase() === "demo") {
                  isDemoTrue.set(true);
                  return;
                }
                let url = isFrontend
                  ? "/preferences"
                  : await $call.getPage("client-dashboard", "preferences");
                goto(url);
              }}
              class="flex p-2 items-center w-full hover:bg-secondary-light dark:hover:bg-secondary-dark rounded-lg"
            >
              <CogOutline class="size-6 mr-4 text-secondary-light dark:text-secondary-dark" />
              <Typography
                preset="Text md/Regular"
                propClass="text-base font-normal"
                text={$_("preferences")}
                color="primary"
              />
            </button>
            <button
              onclick={async () => {
                dropdownUser = !dropdownUser;
                if ($me.authority.toLowerCase() === "demo") {
                  isDemoTrue.set(true);
                  return;
                }
                let url = isFrontend
                  ? "/chat_logs"
                  : await $call.getPage("client-dashboard", "chat-logs");
                goto(url);
              }}
              class="flex p-2 items-center w-full hover:bg-secondary-light dark:hover:bg-secondary-dark rounded-lg"
            >
              <MessageTextSquare01Icon
                svgClass="w-5 h-5 text-secondary-light dark:text-secondary-dark mr-4"
              />
              <Typography
                preset="Text md/Regular"
                propClass="text-base font-normal"
                text={$_("chatLogs")}
                color="primary"
              />
            </button>
            <!-- todo config -->
            <button
              onclick={async () => {
                dropdownUser = !dropdownUser;
                if ($me.authority.toLowerCase() === "demo") {
                  isDemoTrue.set(true);
                  return;
                }
                let url = isFrontend ? "/usage" : await $call.getPage("client-dashboard", "usage");
                if (url !== null) goto(url);
                else {
                  var base_url =
                    location.protocol +
                    "//" +
                    location.hostname +
                    (location.port == "" ? "" : ":" + location.port);

                  goto(`${base_url}/client-dashboard/usage`);
                }
              }}
              class="flex px-3 py-2 items-center w-full hover:bg-secondary-light dark:hover:bg-secondary-dark rounded-lg"
            >
              <Icon name="lineChartUp01" propClass="w-5 h-5" />
              <Typography
                preset="Text md/Regular"
                propClass="text-base font-normal ml-4"
                text={$_("usage")}
                color="primary"
              />
            </button>
          </ul>
        </div>
        <SidebarGroup>
          <button
            type="button"
            class="flex w-full justify-between rounded-lg hover:bg-secondary-light dark:hover:bg-secondary-dark p-2"
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
                <div>
                  {#if dropdownUser}
                    <ChevronUpOutline
                      class="size-7 text-fg-quaternary-light dark:text-fg-quaternary-dark"
                      onclick={(event) => {
                        event.stopPropagation(); // Prevent handleDropdownFocusLoss from being triggered
                        dropdownUser = !dropdownUser;
                      }}
                    />
                  {:else}
                    <ChevronDownOutline
                      class="size-7 text-fg-quaternary-light dark:text-fg-quaternary-dark"
                      onclick={(event) => {
                        event.stopPropagation(); // Prevent handleDropdownFocusLoss from being triggered
                        dropdownUser = !dropdownUser;
                      }}
                    />
                  {/if}
                </div>
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
