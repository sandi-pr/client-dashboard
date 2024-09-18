<script>
  import { Button, Card, Indicator, Dropdown, DropdownItem } from "flowbite-svelte";
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
    webchat = false,
    functButton = () => {},
    functEdit = () => {},
    functAddDocuments = () => {},
    functDelete = () => {},
    funcSelectCollectionWebchat = () => {},
    dropdownOpen = false,
    funcDetailCollection = () => {},
  } = $props();

  const getTextStatus = (status) => {
    if (status === "Ready") {
      return "ready";
    } else if (status === "Reading") {
      return "aiIsReading";
    } else {
      return "noDocumnetAvailable";
    }
  };

  const getTextButton = (status) => {
    if (status === "Ready") {
      return "askAI";
    } else if (status === "Reading") {
      return "reading";
    } else {
      return "startReading";
    }
  };

  const getIconButton = (status) => {
    if (status === "Ready") {
      return MessageChatCircle;
    } else {
      return Stars03Icon;
    }
  };

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
  };
</script>

<Card class="gap-2.5 justify-between !p-0">
  <div
    class="flex justify-between px-4 pb-4 pt-3 border-b border-secondary-light dark:border-secondary-dark"
  >
    <button onclick={funcDetailCollection}>
      <Typography
        preset="Text md/Semibold"
        text={webchat ? collectionName + " (Webchat)" : collectionName}
        propClass="hover:underline line-clamp-1 md:line-clamp-2 fluid-sm"
      />
    </button>
    <div onfocusout={handleDropdownFocusLoss} class=" flex items-center relative">
      <Button
        on:click={() => (dropdownOpen = !dropdownOpen)}
        class="focus-within:ring-0 !p-0"
        color="none"
        size="xs"
        ><DotsVerticalIcon class="text-fg-quinary-light dark:text-fg-quinary-dark" /></Button
      >
      <div
        class={`z-10 ${dropdownOpen ? "absolute top-10" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
      >
        <ul class="py-2 text-sm text-gray-700 dark:text-gray-200">
          <button
            class="flex p-2"
            onclick={() => {
              functAddDocuments(), (dropdownOpen = !dropdownOpen);
            }}>{$_("addDocuments")}</button
          >
          <button
            class="flex p-2"
            onclick={() => {
              functEdit(), (dropdownOpen = !dropdownOpen);
            }}>{$_("editCollectionName")}</button
          >
          <button
            class="flex p-2"
            onclick={() => {
              functDelete(), (dropdownOpen = !dropdownOpen);
            }}>{$_("delete")}</button
          >
          <button
            class="flex p-2"
            onclick={() => {
              funcSelectCollectionWebchat(), (dropdownOpen = !dropdownOpen);
            }}>{$_("selectCollectionWebchat")}</button
          >
        </ul>
      </div>
    </div>
  </div>
  <Card
    on:click={funcDetailCollection}
    class="flex hover:cursor-pointer flex-col items-center border-none shadow-none gap-2.5 px-6 pb-2 md:pb-6 pt-2"
  >
    <div class="relative">
      <FileTypeIcon svgClass="text-utility-brand-600-light dark:text-utility-brand-600-dark" />
      {#if totalDocuments > 0}
        <Indicator
          border
          placement="center-right"
          class="!w-8 !h-8 text-white bg-utility-brand-600-light dark:bg-utility-brand-600-dark text-xs font-semibold p-4">+{totalDocuments}</Indicator
        >
      {/if}
    </div>
    <Badge
      text={$_(getTextStatus(status))}
      type={`${status === "Ready" ? "success" : status === "Reading" ? "warning" : "info"}`}
    />
  </Card>
  <div class="flex justify-center border-t border-secondary-light dark:border-secondary-dark">
    <Button
      disabled={status === undefined || status === null || status === "Reading" || status === "null"}
      on:click={() => {
        functButton();
      }}
      color="light"
      class="flex w-full gap-1 md:gap-1.5 border-none fluid-xs md:fluid-sm"
    >
      <svelte:component
        this={getIconButton(status)}
        svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
      />
      {$_(getTextButton(status))}
    </Button>
  </div>
</Card>
