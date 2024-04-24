<script>
  import * as Icon from "flowbite-svelte-icons";
  import { onMount } from "svelte";
  import {
    Avatar,
    Button,
    DarkMode,
    Dropdown,
    Indicator,
    Breadcrumb,
    DropdownItem,
    BreadcrumbItem,
  } from "flowbite-svelte";
  import { _, locale } from "svelte-i18n";
  import { page } from "$app/stores";
  import US from "../icon/US.svg";
  import ID from "../icon/ID.svg";
  import HomeLineIcon from "../icon/home-line.svelte";

  let notifications = [];
  let showNotifications = false;
  let notificationPosition = { x: null, y: null };

  // Function to fetch notifications (dummy data for demonstration)
  function fetchNotifications() {
    notifications = [
      { id: 1, message: "Notification 1" },
      { id: 2, message: "Notification 2" },
      { id: 3, message: "Notification 3" },
    ];
  }

  // Function to toggle visibility of notifications
  function toggleNotifications(event) {
    showNotifications = !showNotifications;
    if (showNotifications && notifications.length === 0) {
      fetchNotifications();
    }
    // Get position of the button
    const rect = event.target.getBoundingClientRect();
    notificationPosition = {
      x: window.innerWidth - rect.right - 20, // Adjust 16 to fit padding/margin
      y: rect.bottom + window.scrollY + 4, // Adjust 4 to fit padding/margin
    };
  }

  const changeLanguage = (lang) => {
    locale.set(lang);
    dropdownOpen = false;
  };

  let dropdownOpen = false;

  // let mnemonicModal = false;
  // let mnemonicPhrases = Array(12).fill("");

  $: activeUrl = $page.url.pathname;
  $: breadcrumbUrl = activeUrl.split("/").filter(Boolean);

  // Close notifications when clicking outside
  onMount(() => {
    fetchNotifications();

    document.addEventListener("click", (event) => {
      const target = event.target;
      const button = document.getElementById("notification-button");
      if (!button?.contains(target)) {
        showNotifications = false;
      }
    });
  });
</script>

<div class="flex justify-between items-center w-full">
  <Breadcrumb aria-label="Default breadcrumb">
    <BreadcrumbItem href="/dashboard" home>
      <svelte:fragment slot="icon">
        <svelte:component
          this={HomeLineIcon}
          svgClass="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark"
        />
      </svelte:fragment>
    </BreadcrumbItem>
    {#if activeUrl !== "/"}
      {#each breadcrumbUrl as breadcrumb, index (breadcrumb)}
        <BreadcrumbItem href={`/${breadcrumbUrl.slice(0, index + 1).join("/")}`}
          >{breadcrumb}</BreadcrumbItem
        >
      {/each}
    {/if}
  </Breadcrumb>

  <div class="flex gap-4 items-center">
    <!-- <Button
      on:click={() => (mnemonicModal = true)}
      color="light"
      class="px-3 bg-white dark:bg-hc-graydark-900 border-hc-graylight-300 dark:border-hc-graydark-800 hover:bg-hc-graylight-50 hover:dark:bg-hc-graydark-800 focus-within:ring-0"
    >
      <span class="flex items-center"
        ><Indicator size="sm" color="red" class="me-1.5" />Mnemonic unverified</span
      >
    </Button> -->
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
        <Icon.ChevronDownSolid class="w-3 h-3 ms-2 dark:text-white" />
      </Button>
      <Dropdown bind:open={dropdownOpen} class="w-48 overflow-y-auto py-1">
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
      </Dropdown>
    </div>
    <DarkMode size="lg" />

    <div class="relative">
      <Button
        on:click={toggleNotifications}
        id="notification-button"
        class="relative size-10 bg-white hover:bg-hc-graylight-200 focus:ring-hc-graylight-300 focus-within:ring-hc-graylight-300 dark:bg-hc-graydark-950 dark:hover:bg-hc-graydark-700 dark:focus:ring-hc-graydark-700 dark:focus-within:ring-hc-graydark-700"
      >
        <Icon.BellOutline
          class="text-fg-secondary-light dark:text-fg-secondary-dark focus:outline-none"
        />
        <span class="sr-only">Notifications</span>
        <Indicator border size="xl" placement="top-right" class="text-xs font-bold bg-purple-600"
          >{notifications.length}</Indicator
        >
      </Button>

      {#if showNotifications}
        <div
          class="absolute z-10 w-[80vw] lg:w-96 bg-white dark:bg-gray-800 shadow-lg rounded-md py-1"
          style="top: {notificationPosition.y}px; right: {notificationPosition.x}px;"
        >
          {#if notifications.length === 0}
            <p class="px-4 py-2 text-gray-600 dark:text-gray-400">No new notifications</p>
          {:else}
            {#each notifications as notification}
              <p class="px-4 py-2 text-gray-800 dark:text-white">{notification.message}</p>
            {/each}
          {/if}
        </div>
      {/if}
    </div>
  </div>
</div>

<!-- <Modal bind:open={mnemonicModal} size="xs" autoclose={false} class="w-full">
  <form class="flex flex-col space-y-4" action="#">
    <div>
      <Heading tag="h4">Blockchain Authentication</Heading>
      <P size="base">To initiate any on-chain activites, you need to fill the fields below by using your mnemonic phrases</P>
    </div>
    <div class="grid grid-cols-3 gap-4 mb-4">
      {#each mnemonicPhrases as phrase, index (index)}
        <Input type="text" name={`p${index + 1}`} placeholder={`Phrase ${index + 1}`} required bind:value={mnemonicPhrases[index]} />
      {/each}
    </div>
    <div class="flex justify-center">
      <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
        Forgot your phrases? <a href="/" class="text-primary-700 hover:underline dark:text-primary-500"> Contact us </a>
      </div>
    </div>
    <div class="flex justify-center gap-4">
      <Button color="light" class="w-full">Cancel</Button>
      <Button type="submit" class="w-full">Authenticate</Button>
    </div>
  </form>
</Modal> -->
