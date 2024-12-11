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
  import { Button, Avatar, Indicator } from "flowbite-svelte";
  import { Typography } from "$lib";
  import UploadIcon from "$lib/icon/upload-cloud-02.svelte";
  import FolderTypeIcon from "$lib/icon/folder-type.svg";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";
  import JPEGSolidIcon from "$lib/icon/JPEGSolid.svg";
  import JPGSolidIcon from "$lib/icon/JPGSolid.svg";
  import PNGSolidIcon from "$lib/icon/PNGSolid.svg";
  import TrashIcon from "$lib/icon/trash-03.svelte";
    import { is } from "date-fns/locale";

  let {
    title,
    handleFileUploadWrapper,
    multiple = false,
    type = "document",
    status = $bindable(""),
    file = null,
    files = [],
    isProcess = $bindable(false),
    actionButton = false,
    actionView = null,
    actionDelete = () => {},
    isFlex = false,
    isPreview = false,
    largePreview = false,
    srcAvatar = null,
    disabled = false,
  } = $props();

  const textFileType = (type) => {
    if (type === "document") {
      return "PDF (Max. 15MB)";
    } else {
      return "";
    }
  };
</script>

<div class="relative flex flex-col gap-1.5">
  <Typography preset="Text sm/Medium" text={title} color="secondary" />
  {#if file === null || isProcess && files.length < 1}
    <div class={`flex gap-3 ${isPreview ? largePreview ? "flex-col md:flex-row" : "flex-row" : "flex-col"}`}>
      {#if isPreview}
        {#if largePreview}
          {#if srcAvatar}
            <img src={srcAvatar} class="h-32 rounded-xl object-cover" alt="preview" />
          {/if}
        {:else}
          <Avatar class="size-12  object-cover" src={srcAvatar} />
        {/if}
      {/if}
      <Button
        on:click={disabled ? () => {} : handleFileUploadWrapper}
        color="light"
        class={`grow rounded-xl ${disabled ? "cursor-default focus-within:ring-0 hover:bg-primary-light hover:dark:bg-secondary-dark" : ""} border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-secondary-dark gap-3 py-4 flex flex-col justify-start items-center`}
      >
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
    </div>
  {:else if files.length > 0}
    {console.log("cek file", file)}
    <div
      class={`flex ${isFlex ? "flex-row items-center justify-between" : "flex-col"} gap-3 px-6 py-4 rounded-xl border border-primary-light dark:border-primary-dark`}
    >
      {#if files.length > 1}
        <div class="flex flex-col gap-2 grow items-center justify-center">
          <div class="relative">
            <img src={FolderTypeIcon} alt="" class="w-16 h-16" />
            <Indicator border size="xl" color="purple" placement="center-right">
              <span class="text-white text-xs font-bold">{files.length}</span>
            </Indicator>
          </div>
          {#if actionButton}
            <div class="flex w-full gap-3">
              {#if actionView}
                <Button class="w-full" color="light" on:click={actionView}>
                  <Typography preset="Text sm/Semibold" text={$_("viewDocument")} color="secondary" />
                </Button>
              {/if}
              <Button class="grow h-max p-2" color="light" on:click={handleFileUploadWrapper}>{$_("uploadMore")}</Button>
              <Button class="w-max h-max p-2" color="light" on:click={actionDelete}>
                <TrashIcon class="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
              </Button>
            </div>
          {/if}
        </div>
      {:else if !isProcess}
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
          <div class="flex flex-col justify-center">
            <Typography preset="Text sm/Semibold" text={file.name} color="secondary" />
            <Typography preset="Text sm/Regular" text={file.size} color="tertiary" />
            {#if status !== ""}
              <Typography preset="Text sm/Regular" text={`status: ${status}`} color="secondary" />
            {/if}
          </div>
        </div>
        {#if actionButton}
          <div class="flex items-end justify-end gap-3">
            {#if actionView}
              <Button class="w-full" color="light" on:click={actionView}>
                <Typography preset="Text sm/Semibold" text={$_("viewDocument")} color="secondary" />
              </Button>
            {/if}
            {#if multiple}
              <Button class="grow h-max p-2" color="light" on:click={handleFileUploadWrapper}>{$_("uploadMore")}</Button>
            {/if}
            <Button class="w-max h-max p-2" color="light" on:click={actionDelete}>
              <TrashIcon class="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
            </Button>
          </div>
        {/if}
      {:else}
      <div class="grow"></div>
      {/if}
    </div>
  {/if}
</div>
