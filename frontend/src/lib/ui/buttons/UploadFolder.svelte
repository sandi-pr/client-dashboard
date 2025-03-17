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
  import { Button, Avatar, Indicator, Modal } from "flowbite-svelte";
  import { Typography } from "$lib";
  import UploadIcon from "$lib/icon/upload-cloud-02.svelte";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";
  import JPEGSolidIcon from "$lib/icon/JPEGSolid.svg";
  import JPGSolidIcon from "$lib/icon/JPGSolid.svg";
  import PNGSolidIcon from "$lib/icon/PNGSolid.svg";
  import MP3SolidIcon from "$lib/icon/MP3Solid.svg";
  import CSVSolidIcon from "$lib/icon/CSVSolid.svg";
  import CSVIcon from "$lib/icon/csv.svg";
  import TrashIcon from "$lib/icon/trash-03.svelte";
  import EyeIcon from "$lib/icon/eye.svelte";
  import PdfjsViewer from "$lib/components/PdfjsViewer.svelte";
  import FolderSolidIcon from "$lib/icon/folder-solid.svelte";
  import FolderSolid from "$lib/icon/folder-solid.svelte";

  let {
    title,
    handleFileUploadWrapper,
    multiple = false,
    status = $bindable(""),
    file = null,
    files = [],
    tmpFiles = [],
    isProcess = $bindable(false),
    actionDelete,
    actionButton = false,
    actionView = null,
    isFlex = false,
    isPreview = false,
    largePreview = false,
    srcAvatar = null,
    disabled = false,
  } = $props();

  let textarea = $state();
  let modalAllDocuments = $state(false);
  let previewModal = $state(false);
  let selectedFile = $state(null);
  let folderName = $state("");

  const textFileType = (type) => {
    if (type === "document") {
      return "PDF (Max. 15MB)";
    } else {
      return "";
    }
  };

  function actionDeleteFile(fileToDelete) {
    const index = tmpFiles.findIndex((file) => file === fileToDelete);
    if (index !== -1) {
      tmpFiles.splice(index, 1);
    }
    tmpFiles = [...tmpFiles];
  }

  function actionViewFile(fileToView) {
    previewModal = true;
    selectedFile = fileToView;
  }

  // Function to adjust the height of the textarea
  function autosize() {
    textarea.style.height = "auto"; // Reset height
    textarea.style.height = `${textarea.scrollHeight}px`; // Set height to scroll height

    textarea.addEventListener("input", (e) => {
      e.target.style.height = "auto";
      e.target.style.height = `${e.target.scrollHeight}px`;
    });
  }

  function getFiles () {
    var tmpFiles = []

    if (!files) {
      return tmpFiles
    }

    for (const folder in files) {
      const parts = folder.split("/");
      folderName = parts[0];
      let _folder = files[folder]
      for (const f in _folder) {
        tmpFiles.push(_folder[f])
      }
    }

    file = tmpFiles[0]
    return tmpFiles
  }

  // function actionDelete () {
  //   files = {}
  //   tmpFiles = []
  //   file = null
  // }

  $effect(() => {
    if (textarea) {
      autosize();
    }
    tmpFiles = getFiles()
    if (tmpFiles.length > 0) {
      file = tmpFiles[0]
    }
  });

</script>
<div class="relative flex flex-col gap-1.5">
  <Typography preset="Text sm/Medium" text={title} color="secondary" />
  {#if file === null || tmpFiles.length < 1}
    <div
      class={`flex gap-3 ${isPreview ? (largePreview ? "flex-col md:flex-row" : "flex-row") : "flex-col"}`}
    >
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
          <Typography
            preset="Text sm/Semibold"
            text={$_("clickToUpload")}
            color="brand-secondary"
          />
          <!-- <Typography preset="Text xs/Regular" text={textFileType(type)} color="tertiary" /> -->
        </div>
      </Button>
    </div>
  {:else if tmpFiles.length > 0}
    <div
      class={`flex ${isFlex ? "flex-row items-center justify-between" : "flex-col"} gap-3 px-6 py-4 rounded-xl border border-primary-light dark:border-primary-dark`}
    >
      {#if tmpFiles}
        <div class="flex flex-col gap-2 grow items-center justify-center">
          <div class="flex items-center justify-between w-full relative">
            <div class="flex items-center gap-3">
              <FolderSolid svgClass="text-fg-quaternary-light dark:text-fg-quaternary-dark" />
              <Typography preset="Text md/Medium" text={folderName} color="primary" />
            </div>
            <Typography preset="Text sm/Semibold" text={`${tmpFiles.length} ${$_("files")}`} color="primary" />
          </div>
          {#if actionButton}
            <div class="flex w-full items-center justify-center gap-1">
              <Button
                class="w-full h-max p-2"
                color="light"
                on:click={() => (modalAllDocuments = true)}
              >
                {$_("previewCollection")}
              </Button>
              <Button class="w-max h-max p-2" color="light" on:click={actionDelete}>
                <TrashIcon svgClass="w-5 h-5 text-hc-error-600 dark:text-hc-error-500" />
              </Button>
            </div>
          {/if}
        </div>
      {:else if !isProcess}
        {#if actionButton}
          <div class="flex items-end justify-end gap-3">
            {#if multiple}
              <Button class="grow h-max p-2" color="light" on:click={handleFileUploadWrapper}>{$_("uploadMore")}</Button>
              {#if actionView}
                <Button class="w-max h-max p-2" color="light" on:click={() => actionViewFile(file)}>
                  <EyeIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
                </Button>
              {/if}
            {:else}
              {#if actionView}
                <Button class="grow h-max gap-1 p-2" color="light" on:click={() => actionViewFile(file)}>
                  <EyeIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
                  <Typography preset="Text sm/Regular" text={$_("preview")} color="secondary" />
                </Button>
              {/if}
            {/if}
            <Button class="w-max h-max p-2" color="light" on:click={actionDelete}>
              <TrashIcon svgClass="w-5 h-5 text-hc-error-600 dark:text-hc-error-500" />
            </Button>
          </div>
        {/if}
      {:else}
        <div class="grow"></div>
      {/if}
    </div>
  {/if}
</div>

<Modal title={$_("selectedFileToUpload")} bind:open={modalAllDocuments} size="xs" outsideclose>
  {#if tmpFiles.length > 0}
    {#each Object.entries(files) as [pathFolder, folder]}
      {#if pathFolder.split('/').length > 1}
        <div class="flex items-center justify-start gap-3">
          <FolderSolid svgClass="text-fg-quaternary-light dark:text-fg-quaternary-dark" />
          <div class="flex flex-col max-w-72 justify-center">
            <Typography
              preset="Text sm/Semibold"
              propClass="truncate"
              text={pathFolder.split('/')[1].trim()}
              color="secondary"
            />
          </div>
        </div>
      {/if}
      {#each Object.entries(folder) as [_, file]}
        <div class="flex items-center justify-between">
          <div class={`flex items-center justify-start gap-3 $  ${pathFolder.split('/').length > 1 ? 'ml-3' : ''}`}>
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
            {:else if file.name.includes(".mp3") || file.name.includes(".MP3")}
              <img src={MP3SolidIcon} alt="" />
            {:else if file.name.includes(".csv") || file.name.includes(".CSV")}
              <img src={CSVIcon} alt="" />
            {/if}
            <div class="flex flex-col max-w-72 justify-center">
              <Typography
                preset="Text sm/Semibold"
                propClass="truncate"
                text={file.name}
                color="secondary"
              />
              <Typography preset="Text sm/Regular" text={file.size} color="tertiary" />
              {#if status !== ""}
                <Typography preset="Text sm/Regular" text={`status: ${status}`} color="secondary" />
              {/if}
            </div>
          </div>
        </div>
      {/each}
    {/each}
  {/if}
</Modal>

<Modal title={$_("preview")} bind:open={previewModal} size="lg" outsideclose>
  <div class="flex">
    <div
      class="flex w-full justify-center items-center py-6 rounded-lg bg-secondary-light dark:bg-secondary-dark"
    >
      {#if selectedFile.ext.includes(".png") || selectedFile.ext.includes(".jpg") || selectedFile.ext.includes(".jpeg")}
        <img
          src={`data:${selectedFile.type};base64,${selectedFile.base64}`}
          style="max-width: 80%; height: auto;"
          alt="PNG Preview"
        />
      {:else if selectedFile.ext.includes(".pdf")}
        <PdfjsViewer pdfFile={selectedFile.base64} />
      {/if}
    </div>
  </div>
</Modal>
