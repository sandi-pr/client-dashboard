<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: 
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->
<script>
  import { Button } from "flowbite-svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { _ } from "svelte-i18n";

  let {
    dropdownOpen,
    selectedItem = $bindable(""),
    items = [],
    width = "w-full",
    placeholder = "select item",
    isOverlay = false,
    disabled = false,
  } = $props();

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
  };
</script>

<div onfocusout={handleDropdownFocusLoss} class={`flex flex-col relative gap-1 ${width}`}>
  <Button
    on:click={() => (dropdownOpen = !dropdownOpen)}
    class="flex justify-between border-primary-light dark:border-primary-dark"
    color="light"
  >
    {selectedItem === "" ? placeholder : items.find((item) => item.value === selectedItem)?.name}
    <ChevronDownOutline
      class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
    />
  </Button>
  {#if dropdownOpen && items.length > 0 && !disabled}
    {#if isOverlay}
      <div
        id="dropdown"
        class={`z-10 ${dropdownOpen ? "absolute top-12 px-2" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow ${width} dark:bg-gray-700`}
      >
        <ul
          class="py-2 text-sm text-gray-700 dark:text-gray-200"
          aria-labelledby="dropdownDefaultButton"
        >
          {#each items as item}
            <button
              class="flex w-full py-2 px-4"
              onclick={() => ((dropdownOpen = false), (selectedItem = item.value))}
            >
              {item.name}
            </button>
          {/each}
        </ul>
      </div>
    {:else}
      <ul
        style:visibility={dropdownOpen ? "visible" : "hidden"}
        class="py-2 flex flex-col items-center border border-secondary-light dark:border-secondary-dark rounded-lg text-sm text-primary-light dark:text-primary-dark"
        aria-labelledby="dropdownDefaultButton"
      >
        {#each items as item}
          <button
            class="flex w-full py-2 px-4"
            onclick={() => ((dropdownOpen = false), (selectedItem = item.name))}
          >
            {item.name}
          </button>
        {/each}
      </ul>
    {/if}
  {/if}
</div>
