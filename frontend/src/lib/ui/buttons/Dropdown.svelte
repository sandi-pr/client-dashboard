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
  import { twMerge } from "tailwind-merge";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import { _ } from "svelte-i18n";
  import Divider from "../../components/Divider.svelte";

  let {
    dropdownOpen,
    selectedItem = $bindable(""),
    items = [],
    width = "w-full",
    placeholder = "select item",
    isOverlay = false,
    disabled = false,
    onClick = null,
    widget = null,
    propClass = "",
  } = $props();

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
  };

  const buttonClass = () => {
    const defaultClass = `flex px-4 py-2 justify-between border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark focus:ring-0 ${
      disabled
        ? "focus:ring-0 cursor-not-allowed bg-secondary-light dark:bg-secondary-dark hover:bg-secondary-light hover:dark:bg-secondary-dark hover:border-primary-light hover:dark:border-primary-dark"
        : ""
    }`;

    if (propClass) {
      return twMerge(defaultClass, propClass);
    } else {
      return defaultClass;
    }
  };
</script>

<div onfocusout={handleDropdownFocusLoss} class={`flex flex-col relative gap-1 ${width}`}>
  <Button
    on:click={() => (dropdownOpen = !dropdownOpen)}
    class={buttonClass()}
    color="light"
    {disabled}
  >
    <div class="flex items-center gap-2">
      {#if widget}
        <svelte:component
          this={widget}
          color={items.find((item) => item.value === selectedItem)?.value}
        />
      {/if}
      {selectedItem === "" ? placeholder : items.find((item) => item.value === selectedItem)?.name}
    </div>
    <ChevronDownOutline class="w-7 h-7 text-primary-light dark:text-primary-dark" />
  </Button>
  {#if dropdownOpen && items.length > 0 && !disabled}
    {#if isOverlay}
      <div
        id="dropdown"
        class={`z-10 ${width} ${dropdownOpen ? "absolute top-14" : "hidden"} bg-white dark:bg-gray-700 rounded-lg shadow`}
        style="max-height: 200px; overflow-y: auto;"
      >
        <ul
          class="text-sm text-gray-700 dark:text-gray-200"
          aria-labelledby="dropdownDefaultButton"
        >
          {#each items as item, index}
            {#if index !== 0}
              <Divider />
            {/if}
            <button
              class="flex w-full text-left py-2 px-4 hover:bg-gray-100 dark:hover:bg-gray-600"
              onclick={() => {
                dropdownOpen = false;
                selectedItem = item.value;
                if (onClick !== null) onClick(item.value);
              }}
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
        style="max-height: 200px; overflow-y: auto;"
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
