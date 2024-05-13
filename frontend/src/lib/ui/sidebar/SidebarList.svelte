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
  import { LogoBrand, SidebarNavItem, Typography, getIcons } from "$lib";
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
  } from "flowbite-svelte";
  import {
    CogOutline,
    UsersOutline,
    ChevronUpOutline,
    ChevronDownOutline,
    ArrowRightToBracketOutline,
  } from "flowbite-svelte-icons";
  import LogoIcon from "$lib/icon/Logomark.svg";
  import { onMount } from "svelte";
  import Me from "$lib/utils/me";
  import { _ } from "svelte-i18n";

  let {
    getModule,
    toCamelCase,
    name = $bindable(""),
    email = $bindable(""),
    brand = $bindable("Hexapod"),
    isHidden = true,
    isFrontend = $bindable(true),
    funcButton,
    logout,
  } = $props();

  let site = {
    name: brand,
    href: "/",
    img: LogoIcon,
  };

  let sidebarItems = $state([{ name: "home", href: "/", icon: "home", children: [] }]);
  let dropdownUser = $state(false);

  const call = !isFrontend ?? new CR3Call(globalThis.location.origin);

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

  onMount(async () => {
    const sidebarAddition = await sidebarBuilder();
    sidebarItems = [...sidebarItems, ...sidebarAddition];
  });
</script>

<Sidebar
  {activeUrl}
  class={` ${isHidden ? "hidden h-screen" : "fixed h-full"} border-r border-secondary-light dark:border-secondary-dark lg:fixed lg:inset-y-0 lg:z-10 lg:flex w-80`}
>
  <SidebarWrapper
    class="h-full w-full lg:pt-6 lg:px-4 flex flex-col justify-between bg-white dark:bg-primary-dark"
  >
    <SidebarGroup>
      <div class="flex justify-between items-start">
        <LogoBrand brandName="Hexapod" href="/" brandLogo={site.img} />
        {#if isHidden}
          <div></div>
        {:else}
          <!-- <CloseButton on:click={funcButton} class="dark:text-white" /> -->
        {/if}
      </div>
      <Search placeholder={$_("search")} class="bg-white dark:bg-primary-dark mb-6" />
      <!-- TODO Search component belom kelar -->
      {#each sidebarItems as item}
        {#if item.children.length > 0}
          <SidebarDropdownWrapper
            spanClass="flex-1 text-left whitespace-nowrap text-secondary-light dark:text-secondary-dark text-base font-semibold leading-normal"
            btnClass="flex w-full gap-3 items-center px-3 py-2 hover:no-underline rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
            label={$_(item.name)}
          >
            <svelte:fragment slot="icon">
              <svelte:component
                this={getIcons(item.icon)}
                svgClass="w-5 h-5 text-gray-700 dark:text-gray-400"
              />
            </svelte:fragment>
            {#each item.children as child}
              <SidebarNavItem
                href={child.href}
                {activeUrl}
                icon={child.icon}
                classes="flex w-full gap-3 items-center px-11 py-2 hover:no-underline rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
                classesActive="flex w-full gap-3 items-center px-11 py-2 hover:no-underline bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
                ,
                typographyData={{
                  preset: "Text md/Regular",
                  text: $_(child.name),
                  color: "secondary",
                }}
              />
            {/each}
          </SidebarDropdownWrapper>
        {:else}
          <SidebarNavItem
            href={item.href}
            {activeUrl}
            icon={item.icon}
            typographyData={{ preset: "Text md/Semibold", text: $_(item.name), color: "secondary" }}
          />
        {/if}
      {/each}
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
          class={`z-10 ${dropdownUser ? "absolute" : "hidden"} bottom-full w-full p-2 mb-2 bg-white dark:bg-gray-700 divide-y divide-gray-100 rounded-lg border border-secondary-light dark:border-secondary-dark shadow dark:shadow-lg`}
        >
          <ul
            class="text-sm text-gray-700 dark:text-gray-200"
            aria-labelledby="dropdownDefaultButton"
          >
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
            <button
              onclick={async () => {
                window.location = !isFrontend
                  ? await call.getPage("client-dashboard", "users")
                  : "/crm/users";
              }}
              class="flex p-2 items-center w-full hover:bg-secondary-light dark:hover:bg-secondary-dark rounded-lg"
            >
              <UsersOutline class="size-8 mr-4 text-secondary-light dark:text-secondary-dark" />
              <Typography preset="Text lg/Semibold" text="Manage User" color="secondary" />
            </button>
            <button
              class="flex p-2 items-center w-full hover:bg-secondary-light dark:hover:bg-secondary-dark rounded-lg"
            >
              <CogOutline class="size-8 mr-4 text-secondary-light dark:text-secondary-dark" />
              <Typography preset="Text lg/Semibold" text="Settings" color="secondary" />
            </button>
          </ul>
        </div>
        <SidebarGroup>
          <button
            type="button"
            class="flex p-2 justify-between rounded-lg hover:bg-secondary-light dark:hover:bg-secondary-dark"
            onclick={() => {
              dropdownUser = !dropdownUser;
            }}
          >
            <div class="flex w-full h-full space-x-4 items-center cursor-pointer">
              <Avatar />
              <div class="grow">
                <Typography preset="Text md/Semibold" text={name} color="secondary" />
                <Typography preset="Text sm/Regular" text={email} color="tertiary" />
              </div>
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
          </button>
        </SidebarGroup>
      </div>
    </div>
  </SidebarWrapper>
</Sidebar>
