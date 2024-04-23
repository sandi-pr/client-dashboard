<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Dean Ramhan <dean@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import {
    Heading,
    Avatar,
    Button,
    DarkMode,
    // Dropdown, DropdownItem
  } from "flowbite-svelte";
  // import { ChevronDownOutline } from "flowbite-svelte-icons";
  import US from "../icon/US.svg";
  import ID from "../icon/ID.svg";
  import logo from "../icon/Logomark.svg";
  import { _, locale } from "svelte-i18n";

  // let dropdownOpen = false;

  const changeLanguage = (lang) => {
    locale.set(lang);
    // dropdownOpen = false;
  };
</script>

<div class="flex justify-between items-center w-full">
  <div class="flex items-center">
    <img src={logo} alt="logomark" />
    <Heading tag="h4" class="ml-2">Crossroads</Heading>
  </div>

  <div class="flex gap-4 items-center">
    <div>
      <Button color="light">
        {#if $locale === "id"}
          <Avatar src={ID} alt="id" size="xs" class="mr-2" />
        {:else}
          <Avatar src={US} alt="us" size="xs" class="mr-2" />
        {/if}
        <div class="hidden md:block">
          {#if $locale === "id"}
            Indonesia
          {:else}
            English (US)
          {/if}
        </div>
        <!-- <ChevronDownOutline class="w-3 h-3 ms-2 dark:text-white" /> -->
      </Button>
      <!-- <Dropdown
        bind:open={dropdownOpen}
        classContainer="border rounded-lg shadow-sm"
        class="w-48 overflow-y-auto py-1"
      >
        <DropdownItem
          on:click={() => changeLanguage("en")}
          class="flex items-center text-base font-semibold gap-2"
        >
          <Avatar src={US} alt="us" size="xs" class="mr-2" /> English (US)
        </DropdownItem>
        <DropdownItem
          on:click={() => changeLanguage("id")}
          class="flex items-center text-base font-semibold gap-2"
        >
          <Avatar src={ID} alt="us" size="xs" class="mr-2" /> Indonesia
        </DropdownItem>
      </Dropdown> -->
    </div>
    <DarkMode size="lg" />
  </div>
</div>
