<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: 
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->
<script>
  import { _ } from "svelte-i18n";
  import { Button } from "flowbite-svelte";
  import { Typography, Dropdown } from "$lib";
  import UploadIcon from "$lib/icon/upload-cloud-02.svelte";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";
  import JPEGSolidIcon from "$lib/icon/JPEGSolid.svg";
  import JPGSolidIcon from "$lib/icon/JPGSolid.svg";
  import PNGSolidIcon from "$lib/icon/PNGSolid.svg";
  import EyeIcon from "$lib/icon/eye.svelte";
  import TrashIcon from "$lib/icon/trash-03.svelte";
  import StarsIcon from "$lib/icon/stars-01.svelte";

  let {
    title,
    handleFileUploadWrapper,
    multiple = false,
    type = "document",
    file = [],
    rules = [],
    selectedItem = $bindable(""),
    placeholderDropdown = "selectItem",
    status = "",
    dropdownOpen = false,
    actionButton = false,
    actionView = null,
    actionDelete = () => {},
    actionGenerate = () => {},
  } = $props();

  const textFileType = (type) => {
    if (type === "document") {
      return "PDF, DOC, DOCX (Max. 15MB)";
    } else {
      return "";
    }
  };
</script>

<div class="relative flex flex-col gap-1.5">
  {#if file.length === 0}
    <Button
      on:click={handleFileUploadWrapper}
      color="light"
      class={`w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-3 py-4 flex flex-col justify-start items-center`}
    >
      <Typography
        preset="Text md/Semibold"
        text={title}
        color="primary"
        propClass="w-full text-start"
      />
      <div
        class="rounded-xl p-2.5 shadow border border-primary-light dark:border-primary-dark justify-center items-center"
      >
        <UploadIcon class="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
      </div>
      <div class="flex flex-col items-center gap-0.5">
        <Typography preset="Text sm/Semibold" text={$_("clickToUpload")} color="brand-secondary" />
        <Typography preset="Text xs/Regular" text={textFileType(type)} color="tertiary" />
      </div>
    </Button>
  {:else}
    <div
      class="flex flex-col gap-3 px-6 py-4 rounded-xl border border-primary-light dark:border-primary-dark"
    >
      <Typography
        preset="Text md/Semibold"
        text={title}
        color="primary"
        propClass="w-full text-start"
      />
      <div class="flex items-center justify-start gap-3">
        {#if file.name.includes(".pdf") || file.name.includes(".PDF")}
          <img src={PDFSolidIcon} alt="" />
        {:else if file.name.includes(".docx") || file.name.includes(".DOCX")}
          <img src={DOCXSolidIcon} alt="" />
        {:else if file.name.includes(".doc") || file.name.includes(".DOC")}
          <img src={DOCSolidIcon} alt="" />
        {:else if file.name.includes(".jpg") || file.name.includes(".JPG")}
          <img src={JPGSolidIcon} alt="" />
        {:else if file.name.includes(".jpeg") || file.name.includes(".JPEG")}
          <img src={JPEGSolidIcon} alt="" />
        {:else if file.name.includes(".png") || file.name.includes(".PNG")}
          <img src={PNGSolidIcon} alt="" />
        {/if}
        <div class="flex flex-col justify-center w-full pe-12">
          <Typography
            preset="Text sm/Semibold"
            text={file.name}
            color="secondary"
            propClass="truncate"
          />
          <Typography preset="Text sm/Regular" text={file.size} color="tertiary" />
          <Typography preset="Text sm/Regular" text={"status : " + status} color="tertiary" />
        </div>
      </div>
      {#if actionButton}
        <div class="flex gap-3">
          <!-- <Dropdown
            dropdownOpen={dropdownOpen}
            bind:selectedItem={selectedItem}
            items={rules}
            placeholder={$_(placeholderDropdown)}
            isOverlay={true}
          /> -->
          {#if actionView}
            <Button class="p-2 w-full gap-1.5" color="light" on:click={actionView}>
              <EyeIcon class="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
              <Typography preset="Text sm/Semibold" text={$_("view")} color="secondary" />
            </Button>
          {/if}
          <Button class="p-2 w-full gap-1.5" color="light" on:click={actionDelete}>
            <TrashIcon class="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
            <Typography preset="Text sm/Semibold" text={$_("remove")} color="secondary" />
          </Button>
        </div>
      {/if}
      <Button
        color="light"
        on:click={actionGenerate}
        class="gap-1.5 border-brand-light dark:border-primary-dark"
      >
        <Typography preset="Text sm/Semibold" text={"AI-Generate"} color="brand-secondary" />
        <StarsIcon svgClass="w-5 h-5 text-brand-secondary-light dark:text-brand-secondary-dark" />
      </Button>
    </div>
  {/if}
</div>
