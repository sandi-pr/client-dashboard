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
  import Divider from "./Divider.svelte";
  import { _ } from "svelte-i18n";
  import { toCamelCase } from "../utils/toCamelCase.js";
  export let tabDatas: any;
  export let spaceTop = "mt-2";

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

<div class="px-0 mx-0 max-w-full overflow-x-auto overflow-hidden">
  <div class={`flex whitespace-nowrap ${spaceTop}`}>
    {#each tabDatas.data as submenu (submenu)}
      <button
        name={submenu.name}
        on:click={() => handleTabClick(submenu.name)}
        class={`p-2 px-4 font-semibold text-quarterary-light dark:text-quarterary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark ${
          tabDatas.selected === submenu.name &&
          "border-b-2 border-fg-brand-primary-light dark:border-fg-brand-primary-dark !text-brand-secondary-light dark:!text-brand-secondary-dark"
        }`}>{capitalizeFirstLetter($_(toCamelCase(submenu.name)))}</button
      >
    {/each}
  </div>
  <Divider />
</div>
