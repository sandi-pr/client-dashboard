<script>
  import { Avatar, Button } from "flowbite-svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import EN from "../icon/GB.svg";
  import ID from "../icon/ID.svg";
  let dropdownOpen = false;
  let locale = "en";
  export let changeLanguage;
  const handleChangeLanguage = (evt, lang) => {
    evt.preventDefault();
    changeLanguage(lang);
    locale = lang;
    dropdownOpen = false;
  };
</script>

<div>
  <Button color="light" on:click={() => (dropdownOpen = !dropdownOpen)} id="dropdownDefaultButton">
    {#if locale === "id"}
      <Avatar src={ID} alt="id" size="xs" class="mr-2" />
    {:else}
      <Avatar src={EN} alt="us" size="xs" class="mr-2" />
    {/if}
    <ChevronDownOutline class="w-3 h-3 ms-2 dark:text-white" />
  </Button>

  <!-- Dropdown menu -->
  <div
    id="dropdown"
    class={`z-10 ${dropdownOpen ? "fixed" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}
  >
    <ul
      class="py-2 text-sm text-gray-700 dark:text-gray-200"
      aria-labelledby="dropdownDefaultButton"
    >
      <button
        class="flex p-2"
        onclick={(evt) => {
          handleChangeLanguage(evt, "id");
        }}
      >
        <Avatar src={ID} alt="id" size="xs" class="mr-2" /> Indonesia
      </button>
      <button
        class="flex p-2"
        onclick={(evt) => {
          handleChangeLanguage(evt, "en");
        }}
      >
        <Avatar src={EN} alt="us" size="xs" class="mr-2" /> English
      </button>
    </ul>
  </div>
</div>
