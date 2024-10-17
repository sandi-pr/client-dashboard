<script>
  import { createEventDispatcher } from "svelte";
  import { _ } from "svelte-i18n";
  import { toCamelCase } from "$lib";
  export let tabDatas;
  const dispatch = createEventDispatcher();
  function capitalizeFirstLetter(string) {
    return string
      .split(" ")
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(" ");
  }
  function handleTabClick(name) {
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
      class={`p-2 px-4 font-semibold text-quaternary-light dark:text-white hover:bg-tertiary-light hover:dark:bg-brand-solid-dark ${tabDatas.selected === submenu.name && "border-b-2 border-brand-solid-light !text-brand-tertiary-light focus:bg-transparent dark:focus:bg-transparent dark:!text-brand-tertiary-light hover:bg-secondary-light hover:dark:bg-primary-light"}`}
      >{capitalizeFirstLetter($_(toCamelCase(submenu.name)))}</button
    >
  {/each}
</div>

<div class="mx-4 sm:mx-0 border-t-[1px] border-hc-graylight-300 dark:border-hc-graydark-700"></div>
