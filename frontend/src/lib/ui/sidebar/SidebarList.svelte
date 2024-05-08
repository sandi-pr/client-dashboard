<!--
  CR3 Client Dashboard
  Package: gitlab.hexacode.org/public-repo/js/hexacomponents-svelte
  Maintainer:
  * Sandi Permana <sandy@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>

  Copyright (C) 2024 Hexacode Teknologi Indonesia
  All Rights Reserved.
-->

<script lang="ts">
  import { page } from "$app/stores";
  import { LogoBrand, SidebarNavItem, Typography, getIcons } from "$lib";
  import { goto } from "$app/navigation";
  import {
    Sidebar,
    SidebarGroup,
    SidebarItem,
    SidebarWrapper,
    Search,
    Button,
    Avatar,
    CloseButton,
    SidebarDropdownItem,
    SidebarDropdownWrapper,
  } from "flowbite-svelte";
  import { ArrowRightToBracketOutline } from "flowbite-svelte-icons";
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
    funcButton,
    logout,
  }: {
    getModule: any;
    toCamelCase: any;
    name: string;
    email: string;
    brand: string;
    isHidden: boolean;
    funcButton?: any;
    logout: () => void;
  } = $props();

  let site = {
    name: brand,
    href: "/",
    img: LogoIcon,
  };

  let sidebarItems = $state([
    { name: "home", href: "/", icon: "home", children: [] },
  ]);

  const sidebarBuilder = async () => {
    let moduleClient = await getModule();
    const sidebarHalf = moduleClient.map(
      (item: { name: any; href: any; icon: any; children: any }) => {
        return {
          name: item.name,
          href: item.href,
          icon: item.icon,
          children: item.children,
        };
      },
    );

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
                classesActive="flex w-full gap-3 items-center px-11 py-2 hover:no-underline bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700", 
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
      <SidebarGroup>
        <!-- <SidebarNavItem href="/admin_setting" {activeUrl} icon="setting" typographyData={{ preset: "Text md/Semibold", text: "Super Admin Setting", color: "secondary" }} /> -->
        <SidebarNavItem
          href="/settings"
          {activeUrl}
          icon="setting"
          typographyData={{ preset: "Text md/Semibold", text: $_("settings"), color: "secondary" }}
        />
      </SidebarGroup>
      <SidebarGroup border>
        <div class="flex justify-between">
          <div class="flex w-full h-full space-x-4 items-center pl-2 cursor-pointer">
            <Avatar />
            <div class="grow">
              <Typography preset="Text md/Semibold" text={name} color="secondary" />
              <Typography preset="Text sm/Regular" text={email} color="tertiary" />
            </div>
          </div>
          <Button
            on:click={handleLogout}
            color="light"
            class="focus-within:ring-0 border-0 p-2 bg-transparent dark:bg-transparent"
          >
            <ArrowRightToBracketOutline
              class="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark transition duration-75 group-hover:text-gray-900 dark:group-hover:text-white"
            />
          </Button>
        </div>
      </SidebarGroup>
    </div>
  </SidebarWrapper>
</Sidebar>
