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
  import FolderTypeIcon from "$lib/icon/folder-type.svg";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";
  import JPEGSolidIcon from "$lib/icon/JPEGSolid.svg";
  import JPGSolidIcon from "$lib/icon/JPGSolid.svg";
  import PNGSolidIcon from "$lib/icon/PNGSolid.svg";
  import CSVIcon from "$lib/icon/csv.svg";
  import TrashIcon from "$lib/icon/trash-03.svelte";
  import EyeIcon from "$lib/icon/eye.svelte";
  import UploadFileIcon from "$lib/icon/upload-01.svelte";
  import GoogleDriveIcon from "$lib/icon/google-drive.svelte";
  import PdfjsViewer from "$lib/components/PdfjsViewer.svelte";
  import FolderSolid from "$lib/icon/folder-solid.svelte";
  import MP3SolidIcon from "$lib/icon/MP3Solid.svg";
  import ChevronUp from "$lib/icon/chevron-up.svelte";
  import ChevronDown from "$lib/icon/chevron-down.svelte";
  import { onMount } from "svelte";

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
    folderMode = false,
    handleDragOver = () => {},
    handleDrop = () => {},
  } = $props();

  let textarea = $state();
  let modalAllDocuments = $state(false);
  let modalAllFolderDocuments = $state(false);
  let previewModal = $state(false);
  let selectedFile = $state(null);
  let expanded = $state({});

  function actionDeleteFile(fileToDelete) {
    const index = files.findIndex((file) => file === fileToDelete);
    if (index !== -1) {
      files.splice(index, 1);
    }
    files = [...files];
  }

  function actionViewFile(fileToView) {
    console.log("cek file", fileToView);
    // return;
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

  function transformData(input, parentPath = "") {
    let result = {};

    input.forEach((folder) => {
      const currentPath = parentPath ? `${parentPath}/${folder.name}` : folder.name;
      result[currentPath] = {};

      folder.files.forEach((file) => {
        if (file.files) {
          // Jika ada nested folder, panggil rekursif
          Object.assign(result, transformData([file], currentPath));
        } else {
          // Jika itu file, masukkan ke dalam folder saat ini
          result[currentPath][file.name] = file;
        }
      });
    });

    return result;
  }

  function getFileIcon(filename) {
    if (filename.includes(".pdf")) return PDFSolidIcon;
    if (filename.includes(".docx")) return DOCXSolidIcon;
    if (filename.includes(".doc")) return DOCSolidIcon;
    if (filename.includes(".jpg")) return JPGSolidIcon;
    if (filename.includes(".jpeg")) return JPEGSolidIcon;
    if (filename.includes(".png")) return PNGSolidIcon;
    if (filename.includes(".mp3")) return MP3SolidIcon;
    if (filename.includes(".csv")) return CSVIcon;
    return PDFSolidIcon; // Default icon
  }

  // Toggle Expand/Collapse Folder
  function toggleFolder(path) {
    expanded[path] = !expanded[path];
  }

  $effect(() => {
    if (files) {
      Object.keys(files).forEach((pathFolder) => {
        if (pathFolder.includes("/")) {
          expanded[pathFolder] = false; // Default collapsed
        }
      });
    }
  });

  $effect(() => {
    if (textarea) {
      autosize();
    }
  });
</script>

<div class="relative flex flex-col gap-1.5">
  <Typography preset="Text sm/Medium" text={title} color="secondary" />
  {#if file === null || (isProcess && files.length < 1)}
    <div
      class={`flex gap-3 ${isPreview ? (largePreview ? "flex-col md:flex-row" : "flex-row") : "flex-col"}`}
      ondragover={handleDragOver}
      ondrop={handleDrop}
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
        class="grow rounded-xl border gap-3 py-4 flex flex-col justify-start items-center {disabled
          ? 'cursor-not-allowed focus:ring-0 border-disabled-subtle-light dark:border-disabled-subtle-dark hover:border-disabled-subtle-light hover:dark:border-disabled-subtle-dark focus:border-disabled-subtle-light focus:dark:border-disabled-subtle-dark bg-disabled-subtle-light dark:bg-disabled-subtle-dark hover:bg-disabled-subtle-light dark:hover:bg-disabled-subtle-dark text-disabled-light dark:text-disabled-dark'
          : 'border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark focus:ring-0'}"
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
        </div>
      </Button>
    </div>
  {:else if files.length > 0}
    <!-- {console.log("cek file", file)} -->
    <div
      class={`flex ${isFlex ? "flex-row items-center justify-between" : "flex-col"} gap-3 px-6 py-4 rounded-xl border border-primary-light dark:border-primary-dark`}
    >
      {#if files.length > 1 || folderMode}
        <div class="flex flex-col gap-2 grow items-center justify-center">
          {#if !folderMode}
            <div class="relative">
              <img src={FolderTypeIcon} alt="" class="w-16 h-16" />
              <Indicator border size="xl" color="purple" placement="center-right">
                <span class="text-white text-xs font-bold">{files.length}</span>
              </Indicator>
            </div>
          {:else}
            <div class="w-full flex justify-between items-center">
              <div class="flex gap-4 items-center p-4">
                <GoogleDriveIcon
                  class="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
                />
                <Typography preset="Text sm/Medium" text={files[0].name} color="primary" />
              </div>
              <div>
                <Typography
                  preset="Text sm/Semibold"
                  text={files[0].files.length + " files"}
                  color="tertiary"
                />
              </div>
            </div>
          {/if}
          {#if actionButton}
            <div class="flex w-full items-center justify-center gap-1">
              {#if !folderMode}
                <Button
                  class="w-max h-max gap-1.5 p-2"
                  color="light"
                  on:click={handleFileUploadWrapper}
                >
                  <UploadFileIcon
                    svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
                  />
                  {$_("uploadMore")}
                </Button>
              {/if}
              <Button
                class="{folderMode ? 'w-full' : 'w-max'} h-max p-2"
                color="light"
                on:click={() => {
                  if (folderMode) {
                    files = transformData(files);
                    modalAllFolderDocuments = true;
                  } else {
                    modalAllDocuments = true;
                  }
                }}
              >
                {$_("viewAll")}
              </Button>
              <Button class="w-max h-max p-2" color="light" on:click={actionDelete}>
                <TrashIcon svgClass="w-5 h-5 text-hc-error-600 dark:text-hc-error-500" />
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
          {:else if file.name.includes(".csv") || file.name.includes(".CSV")}
            <img src={CSVIcon} alt="" />
          {:else if file.name.includes(".mp3") || file.name.includes(".MP3")}
            <img src={MP3SolidIcon} alt="" />
          {:else}
            <img src={GoogleDriveIcon} alt="" />
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
        {#if actionButton}
          <div class="flex items-end justify-end gap-3">
            {#if multiple}
              <Button class="grow h-max p-2" color="light" on:click={handleFileUploadWrapper}
                >{$_("uploadMore")}</Button
              >
              {#if actionView}
                <Button class="w-max h-max p-2" color="light" on:click={() => actionViewFile(file)}>
                  <EyeIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
                </Button>
              {/if}
            {:else if actionView}
              <Button
                class="grow h-max gap-1 p-2"
                color="light"
                on:click={() => actionViewFile(file)}
              >
                <EyeIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
                <Typography preset="Text sm/Regular" text={$_("preview")} color="secondary" />
              </Button>
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

<Modal title={$_("selectedFileToUpload")} bind:open={modalAllDocuments} size="sm" outsideclose>
  {#if files.length > 1}
    {#each files as file}
      <div class="flex items-center justify-between">
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
        <div class="flex items-center justify-end gap-1">
          <Button
            class="w-max h-max p-2 !border-0"
            color="light"
            on:click={() => actionViewFile(file)}
          >
            <EyeIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
          </Button>
          <Button
            class="w-max h-max p-2 !border-0"
            color="light"
            on:click={() => actionDeleteFile(file)}
          >
            <TrashIcon svgClass="w-5 h-5 text-hc-error-600 dark:text-hc-error-500" />
          </Button>
        </div>
      </div>
    {/each}
  {/if}
</Modal>

<Modal
  title={$_("selectedFileToUpload")}
  bind:open={modalAllFolderDocuments}
  size="sm"
  outsideclose
  bodyClass="p-4 md:p-5 space-y-4 flex-1 max-h-[60vh] overflow-y-auto overscroll-contain"
>
  {#each Object.entries(files) as [pathFolder, folder]}
    {#if pathFolder.includes("/")}
      <!-- Tampilkan FOLDER terlebih dahulu -->
      <button
        class="flex items-center w-full justify-between cursor-pointer rounded hover:bg-gray-100 dark:hover:bg-gray-800"
        onclick={() => toggleFolder(pathFolder)}
      >
        <div class="flex items-center gap-3">
          <FolderSolid svgClass="text-fg-quaternary-light dark:text-fg-quaternary-dark" />
          <Typography
            preset="Text sm/Semibold"
            propClass="truncate"
            text={pathFolder.split("/").slice(-1)[0]}
            color="secondary"
          />
        </div>
        {#if expanded[pathFolder]}
          <ChevronUp svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
        {:else}
          <ChevronDown svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
        {/if}
        <!-- <span class="text-gray-500">{expanded[pathFolder] ? "▼" : "▶"}</span> -->
      </button>

      <!-- Jika folder di-expand, tampilkan isinya -->
      {#if expanded[pathFolder]}
        {#each Object.entries(folder) as [_, file]}
          <div class="flex items-center justify-start gap-3 ml-6">
            <img src={getFileIcon(file.name)} alt="" />
            <div class="flex flex-col max-w-72 justify-center">
              <Typography
                preset="Text sm/Semibold"
                propClass="truncate"
                text={file.name}
                color="secondary"
              />
              <Typography preset="Text sm/Regular" text={file.size} color="tertiary" />
            </div>
          </div>
        {/each}
      {/if}
    {/if}
  {/each}

  <!-- Sekarang tampilkan file yang bukan di dalam folder -->
  {#each Object.entries(files) as [pathFolder, folder]}
    {#if !pathFolder.includes("/")}
      {#each Object.entries(folder) as [_, file]}
        <div class="flex items-center justify-start gap-3">
          <img src={getFileIcon(file.name)} alt="" />
          <div class="flex flex-col max-w-72 justify-center">
            <Typography
              preset="Text sm/Semibold"
              propClass="truncate"
              text={file.name}
              color="secondary"
            />
            <Typography preset="Text sm/Regular" text={file.size} color="tertiary" />
          </div>
        </div>
      {/each}
    {/if}
  {/each}
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
      {:else if selectedFile.ext.includes(".mp3") || selectedFile.ext.includes(".wav") || selectedFile.ext.includes(".ogg")}
        <audio controls class="w-full max-w-md mx-auto">
          <source
            src={`data:${selectedFile.type};base64,${selectedFile.base64}`}
            type="audio/mpeg"
          />
          Your browser does not support the audio element.
        </audio>
      {/if}
    </div>
  </div>
</Modal>

<!-- <style>
  audio {
    border: none; /* Remove border from audio element */
  }
  audio::-webkit-media-controls-panel {
    background-color: #1f2937; /* Example color */
    border: none; /* Remove border from controls panel */
  }
  audio::-webkit-media-controls-play-button,
  audio::-webkit-media-controls-volume-slider,
  audio::-webkit-media-controls-current-time-display,
  audio::-webkit-media-controls-time-remaining-display {
    color: #f9fafb; /* Example color */
    border: none; /* Remove border from individual controls */
  }
</style> -->
