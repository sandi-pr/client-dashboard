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
  export let tabDatas: { selected: string; data: { name: string }[] };
  export let changeTabbar: (tabName: string) => void;

  function capitalizeFirstLetter(string: string) {
    return string
      .split(" ")
      .map((word: string) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(" ");
  }
</script>

<div
  class="flex whitespace-nowrap mt-2 mx-4 sm:px-0 sm:mx-0 max-w-full overflow-x-auto overflow-hidden"
>
  {#each tabDatas.data as submenu (submenu)}
    <button
      name={submenu.name}
      on:click={() => changeTabbar(submenu.name)}
      class={`p-2 px-4 font-semibold text-quaternary-light dark:text-white hover:bg-tertiary-light hover:dark:bg-brand-solid-dark ${tabDatas.selected === submenu.name && "border-b-2 border-brand-solid-light !text-brand-tertiary-light focus:bg-transparent dark:focus:bg-transparent dark:!text-brand-tertiary-light hover:bg-secondary-light hover:dark:bg-primary-light"}`}
      >{capitalizeFirstLetter(submenu.name)}</button
    >
  {/each}
</div>

<div class="mx-4 sm:mx-0 border-t-[1px] border-hc-graylight-300 dark:border-hc-graydark-700"></div>

<!--
@component
## Config
@prop tabDatas: object; // Configuration for the tab bar.
@prop changeTabbar: function; // Function to change the tab bar.
-->
