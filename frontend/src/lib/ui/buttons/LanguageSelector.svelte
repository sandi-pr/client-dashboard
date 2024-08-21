<script lang="ts">
  import { Avatar, Button } from "flowbite-svelte";
  import { Icon } from "$lib";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import EN from "$lib/icon/GB.svg";
  import ID from "$lib/icon/ID.svg";
  import { locale } from "svelte-i18n";
  // let { changeLanguage }: { changeLanguage: (lang: string) => void } = $props();
  let dropdownOpen = $state(false);
  let myLocale = $state("en");

  const changeLanguage = (lang) => {
    locale.set(lang);
    localStorage.setItem("lang", lang);
    dropdownOpen = false;
  };

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return 
    dropdownOpen = false
  }

  // const handleChangeLanguage = (evt: any, lang: string) => {
  //   evt.preventDefault();
  //   changeLanguage(lang);
  //   myLocale = lang;
  //   dropdownOpen = !dropdownOpen;
  // };
</script>

<div onfocusout={handleDropdownFocusLoss}>
  <Button color="light" on:click={() => (dropdownOpen = !dropdownOpen)} id="dropdownDefaultButton" class="py-2 gap-1 px-2.5 md:px-3">
    {#if $locale === "id"}
      <Avatar src={ID} alt="id" size="xs" />
    {:else}
      <Avatar src={EN} alt="us" size="xs" />
    {/if}
    <Icon name="chevronDown" />
    <!-- <ChevronDownOutline class="w-5 h-5 dark:text-white" /> -->
  </Button>

  <!-- Dropdown menu -->
  <div id="dropdown" class={`z-10 ${dropdownOpen ? "fixed" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}>
    <ul class="py-2 text-sm text-gray-700 dark:text-gray-200" aria-labelledby="dropdownDefaultButton">
      <button
        class="flex p-2"
        onclick={() => changeLanguage("id")}>
        <Avatar src={ID} alt="id" size="xs" class="mr-2" /> Indonesia
      </button>
      <button
        class="flex p-2"
        onclick={() => changeLanguage("en")}>
        <Avatar src={EN} alt="us" size="xs" class="mr-2" /> English
      </button>
    </ul>
  </div>
</div>
