<script>
  import {
    Button,
    Card,
    Indicator,
    Dropdown,
    DropdownItem,
  } from "flowbite-svelte";
  import { Typography, Badge } from "$lib";
  import { _ } from "svelte-i18n";
  import FileTypeIcon from "$lib/icon/file-type-icon.svelte";
  import Stars03Icon from "$lib/icon/stars-03.svelte";
  import MessageChatCircle from "$lib/icon/message-chat-circle.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";

  let {
    collectionName,
    totalDocuments = 0,
    status,
    functButton = () => {},
    functEdit = () => {},
    functAddDocuments = () => {},
    functDelete = () => {},
    dropdownOpen = false,
    funcDetailCollection = () => {},
  } = $props();

  const getTextStatus = (status) => {
    if (status === "ready") {
      return "ready";
    } else if (status === "reading") {
      return "aiIsReading";
    } else {
      return "collectingData";
    }
  };

  const getTextButton = (status) => {
    if (status === "ready") {
      return "askAI";
    } else if (status === "reading") {
      return "reading";
    } else {
      return "startReading";
    }
  };

  const getIconButton = (status) => {
    if (status === "ready") {
      return MessageChatCircle;
    } else {
      return Stars03Icon;
    }
  };
</script>

<Card class="gap-2.5 justify-between !p-0">
  <div
    class="flex justify-between px-4 pb-4 pt-3 border-b border-secondary-light dark:border-secondary-dark"
  >
    <button onclick={funcDetailCollection}>
      <Typography preset="Text md/Semibold" text={collectionName} propClass="hover:underline" />
    </button>
    <div class=" flex items-center relative">
      <Button on:click={() => (dropdownOpen = !dropdownOpen)} class="focus-within:ring-0 !p-0" color="none" size="xs"
        ><DotsVerticalIcon
          class="text-fg-quinary-light dark:text-fg-quinary-dark"
        /></Button
      >
      <div class={`z-10 ${dropdownOpen ? "absolute top-10" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}>
        <ul class="py-2 text-sm text-gray-700 dark:text-gray-200">
          <button class="flex p-2" onclick={() => {functAddDocuments(), dropdownOpen = !dropdownOpen}}>{$_("addDocuments")}</button>
          <button class="flex p-2" onclick={() => {functEdit(), dropdownOpen = !dropdownOpen}}>{$_("editCollectionName")}</button>
          <button class="flex p-2" onclick={() => {functDelete(), dropdownOpen = !dropdownOpen}}>{$_("delete")}</button>
        </ul>
      </div>
    </div>
  </div>
  <Card on:click={funcDetailCollection} class="flex hover:cursor-pointer flex-col items-center border-none shadow-none gap-2.5 px-6 pb-6 pt-2">
    <div class="relative">
      <FileTypeIcon />
      {#if totalDocuments > 0}
        <Indicator
          color="purple"
          border
          placement="center-right"
          class="!w-8 !h-8 text-white text-xs font-semibold p-4"
          >+{totalDocuments}</Indicator
        >
      {/if}
    </div>
    <Badge
      text={$_(getTextStatus(status))}
      type={`${status === "ready" ? "success" : status === "reading" ? "warning" : "info"}`}
    />
  </Card>
  <div class="flex justify-center border-t border-secondary-light dark:border-secondary-dark">
    <Button
      on:click={functButton}
      disabled={status === "reading"}
      color="light"
      class="flex w-full gap-1.5 border-none"
    >
      <svelte:component
        this={getIconButton(status)}
        svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
      />
      {$_(getTextButton(status))}
    </Button>
  </div>
</Card>
