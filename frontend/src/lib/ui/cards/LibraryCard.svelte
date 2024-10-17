<script>
  import { Button, Card, Indicator, Dropdown, DropdownItem } from "flowbite-svelte";
  import { Typography, Badge } from "$lib";
  import { _ } from "svelte-i18n";
  import FileTypeIcon from "$lib/icon/file-type-icon.svelte";
  import Stars03Icon from "$lib/icon/stars-03.svelte";
  import MessageChatCircle from "$lib/icon/message-chat-circle.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import BookClosedIcon from "$lib/icon/book-closed.svelte";
  import DotIcon from "$lib/icon/dot.svelte";
  import PlusIcon from "$lib/icon/plus.svelte";
  import Folder from "$lib/icon/folder.svelte";

  let {
    collectionName,
    totalDocuments = 0,
    status,
    usedFor = "libraries",
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
      return "chatReady";
    } else if (status === "Reading") {
      return "aiIsReading";
    } else if (status === "Empty") {
      return "noDocumnetAvailable";
    } else {
      return "startReading";
    }
  };

  const getTextButton = (status) => {
    if (status === "Ready") {
      return "askAI";
    } else if (status === "Reading") {
      return "reading";
    } else if (status === "Empty") {
      return "addDocuments";
    } else {
      return "startReading";
    }
  };

  const getIconButton = (status) => {
    if (status === "Ready") {
      return MessageChatCircle;
    } else if (status === "Reading") {
      return DotIcon;
    } else if (status === "Empty") {
      return "";
    } else {
      return Stars03Icon;
    }
  };

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
  };
</script>

<Card
  class={`gap-2.5 justify-between !p-0 ${usedFor === "libraries_collections" ? "max-w-full" : ""}`}
>
  <div
    class="flex justify-between px-4 pb-4 pt-3 border-b border-secondary-light dark:border-secondary-dark"
  >
    <button onclick={funcDetailCollection} class="flex items-center gap-3 grow">
      {#if usedFor === "libraries"}
        <BookClosedIcon svgClass="text-fg-quaternary-light dark:text-fg-quaternary-dark" />
        <div class="flex flex-col gap-1">
          <Typography
            preset="Text md/Semibold"
            text={collectionName}
            propClass="line-clamp-1 md:line-clamp-2 fluid-sm"
          />
          <Badge
            icon={getIconButton(status)}
            iconClass="w-3 h-3 text-utility-success-500-light dark:text-utility-success-500-dark"
            text={$_(getTextStatus(status))}
            type={`${status === "Ready" ? "success" : status === "Reading" ? "warning" : status === "Empty" ? "info" : "brand"}`}
          />
        </div>
      {:else if usedFor === "libraries_collections"}
        <Folder svgClass="w-5 h-5 text-quarterary-light dark:text-quarterary-dark" />
        <Typography preset="Text sm/Semibold" color="secondary" text={collectionName} />
      {:else}
        <BookClosedIcon svgClass="text-fg-quaternary-light dark:text-fg-quaternary-dark" />
        <div class="flex flex-col gap-1">
          <Typography
            preset="Text md/Semibold"
            text={collectionName}
            propClass="line-clamp-1 md:line-clamp-2 fluid-sm"
          />
          <Badge
            icon={getIconButton(status)}
            iconClass="w-3 h-3 text-utility-success-500-light dark:text-utility-success-500-dark"
            text={$_(getTextStatus(status))}
            type={`${status === "Ready" ? "success" : status === "Reading" ? "warning" : status === "Empty" ? "info" : "brand"}`}
          />
        </div>
      {/if}
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
        class={`z-10 ${dropdownOpen ? "absolute top-10 right-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
      >
        <ul class="py-2 grow text-sm">
          <button
            class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
            onclick={() => {
              functEdit(), (dropdownOpen = !dropdownOpen);
            }}>{$_("editLibraryName")}</button
          >
          <button
            class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
            onclick={() => {
              functDelete(), (dropdownOpen = !dropdownOpen);
            }}>{$_("delete")}</button
          >
        </ul>
      </div>
    </div>
  </div>
</Card>
