<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Dean Ramhan <dean@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<!--
  dataExample = {
    selected: "information",
    data: [
      {
        name: "information",
      },
      {
        name: "recipients",
      }
    ]
  }
-->

<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { _ } from "svelte-i18n";
  import { toCamelCase } from "../utils/toCamelCase.js";
  export let tabDatas: any;

  const dispatch = createEventDispatcher();

  function capitalizeFirstLetter(string: string) {
    return string
      .split(" ")
      .map((word: string) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(" ");
  }

  function handleTabClick(name: string) {
    tabDatas.selected = name;
    dispatch("tabChange", name);
  }
</script>

<div
  class="flex whitespace-nowrap mt-2 mx-4 sm:px-0 sm:mx-0 max-w-full overflow-x-auto overflow-hidden"
>
  {#each tabDatas.data as submenu (submenu)}
    <button
      name={submenu.name}
      on:click={() => handleTabClick(submenu.name)}
      class={`p-2 px-4 font-semibold text-hc-graylight-500 dark:text-white hover:bg-hc-graylight-100 hover:dark:bg-hc-brand-600 ${
        tabDatas.selected === submenu.name &&
        "border-b-2 border-hc-brand-600 !text-hc-brand-600 focus:bg-transparent dark:focus:bg-transparent dark:!text-hc-brand-600 hover:bg-hc-brand-100 hover:dark:bg-white"
      }`}>{capitalizeFirstLetter($_(toCamelCase(submenu.name)))}</button
    >
  {/each}
</div>

<div class="mx-4 sm:mx-0 border-t-[1px] border-hc-graylight-300 dark:border-hc-graydark-700"></div>
