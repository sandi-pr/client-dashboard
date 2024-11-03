<script>
  import { Button, Card, Indicator, Dropdown, DropdownItem } from "flowbite-svelte";
  import { Typography, Badge, checkOr } from "$lib";
  import { _ } from "svelte-i18n";
  import FileTypeIcon from "$lib/icon/file-type-icon.svelte";
  import Stars03Icon from "$lib/icon/stars-03.svelte";
  import MessageChatCircle from "$lib/icon/message-chat-circle.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import PlusIcon from "$lib/icon/plus.svelte";

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
    } else if (status === "Empty") {
      return "noDocumnetAvailable";
    } else {
      return "collectingData";
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
    } else if (status === "Empty") {
      return PlusIcon;
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
    <button disabled={!checkOr(["Admin", "Manager", "User"])} onclick={funcDetailCollection}>
      <Typography
        preset="Text md/Semibold"
        text={webchat ? collectionName + " (Webchat)" : collectionName}
        propClass="hover:underline line-clamp-1 md:line-clamp-2 fluid-sm"
      />
    </button>
    {#if checkOr(["Admin", "Manager"])}
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
                functAddDocuments(), (dropdownOpen = !dropdownOpen);
              }}>{$_("addDocuments")}</button
            >
            <button
              class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
              onclick={() => {
                functEdit(), (dropdownOpen = !dropdownOpen);
              }}>{$_("editCollectionName")}</button
            >
            {#if checkOr(["Admin"])}
              <button
                class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                onclick={() => {
                  functDelete(), (dropdownOpen = !dropdownOpen);
                }}>{$_("delete")}</button
              >
            {/if}
            <!-- <button
            disabled={totalDocuments === 0}
            class={`flex p-2 w-full ${totalDocuments === 0 ? 'text-disabled-light dark:text-disabled-dark' : 'text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark'}`}
            onclick={() => {
              funcSelectCollectionWebchat(), (dropdownOpen = !dropdownOpen);
            }}
            ><div class={totalDocuments === 0 ? "text-disabled-light dark:text-disabled-dark" : ""}>
              {$_("selectCollectionWebchat")}
            </div></button
          > -->
          </ul>
        </div>
      </div>
    {/if}
  </div>

  {#if checkOr(["Admin", "Manager", "User"])}
    <Card
      on:click={funcDetailCollection}
      class="flex hover:cursor-pointer flex-col items-center border-none shadow-none gap-2.5 px-6 pb-2 md:pb-6 pt-2"
    >
      <div class="relative">
        <FileTypeIcon
          svgClass="text-utility-brand-600-light dark:text-utility-brand-600-dark w-16 h-16 md:w-20 md:h-20"
        />
        {#if totalDocuments > 0}
          <Indicator
            border
            placement="center-right"
            class="!w-6 !h-6 !md:w-8 !md:h-8 text-white bg-utility-brand-600-light dark:bg-utility-brand-600-dark fluid-xs font-semibold p-2 md:p-4"
            >+{totalDocuments}</Indicator
          >
        {/if}
      </div>
      <Badge
        text={$_(getTextStatus(status))}
        type={`${status === "Ready" ? "success" : status === "Reading" ? "warning" : "info"}`}
      />
    </Card>
  {:else}
    <Card
      class="flex hover:cursor-pointer flex-col items-center border-none shadow-none gap-2.5 px-6 pb-2 md:pb-6 pt-2"
    >
      <div class="relative">
        <FileTypeIcon
          svgClass="text-utility-brand-600-light dark:text-utility-brand-600-dark w-16 h-16 md:w-20 md:h-20"
        />
        {#if totalDocuments > 0}
          <Indicator
            border
            placement="center-right"
            class="!w-6 !h-6 !md:w-8 !md:h-8 text-white bg-utility-brand-600-light dark:bg-utility-brand-600-dark fluid-xs font-semibold p-2 md:p-4"
            >+{totalDocuments}</Indicator
          >
        {/if}
      </div>
      <Badge
        text={$_(getTextStatus(status))}
        type={`${status === "Ready" ? "success" : status === "Reading" ? "warning" : "info"}`}
      />
    </Card>
  {/if}
  <div class="flex justify-center border-t border-secondary-light dark:border-secondary-dark">
    <Button
      disabled={status === "Reading" || status === "Empty"
        ? !checkOr(["Admin", "Manager"])
        : !checkOr(["Admin", "Manager", "User"])}
      on:click={() => {
        if (status === "Empty") {
          functAddDocuments();
        } else {
          functButton();
        }
      }}
      color="light"
      class="flex w-full px-2 py-2 md:px-4 md:py-2.5 gap-1 md:gap-1.5 border-none fluid-xs md:fluid-sm"
    >
      <svelte:component
        this={getIconButton(status)}
        svgClass="w-4 h-4 md:w-5 md:h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
      />
      {$_(getTextButton(status))}
    </Button>
  </div>
</Card>
