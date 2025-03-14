<script>
  import { _ } from "svelte-i18n";
  import MultiSelect from "svelte-multiselect";
  import { Typography, UploadFile, InputRw, FeaturedIcon, getIcons, UploadFolder } from "$lib";
  import { Modal, Button, Alert, Helper } from "flowbite-svelte";
  import { InfoCircleSolid } from "flowbite-svelte-icons";
  import TabDetail from "$lib/components/TabDetail.svelte";

  let {
    addDocumentsModal = $bindable(false),
    nameIsAvailable = $bindable(true),
    selectedTypeCollection,
    files = $bindable([]),
    alltags,
    tags,
    isProcess,
    handleFileUploadWrapper,
    fileNotSupport,
    hasDocOrDocx,
    tmpDocumentTitle,
    tmpDocumentDescription,
    tmpDocumentLink,
    handleInput,
    buttonDisabled = $bindable(true),
    addDocumentByCollection,
    actionDelete,
    isFolder,
    groupedFiles = $bindable(null),
  } = $props();

  let tabDatas = $state({
    selected: "Upload",
    data: [
      {
        name: "Upload",
      },
      {
        name: "Link",
      },
    ],
  });

  let stepLink = $state("input");

</script>

<Modal bind:open={addDocumentsModal} size="sm">
  <svelte:fragment slot="header">
    <div class="flex items-center gap-4">
      <FeaturedIcon classesSvg={getIcons("folderDownload")} />
      <div>
        <Typography preset="Text lg/Semibold" text={$_("uploadCollection")} />
        <Typography preset="Text sm/Regular" text={$_("descUploadCollection")} />
      </div>
    </div>
  </svelte:fragment>
  <div class="flex flex-col gap-4">
    <div>
      <TabDetail
        {tabDatas}
        on:tabChange={(e) => (tabDatas.selected = e.detail)}
        spaceTop="mt-0 !mx-0"
      />
    </div>
    <div class="flex flex-col gap-1.5 relative">
      <Typography
        preset="Text sm/medium"
        color="secondary"
        text={$_("documentType")}
      />
      <InputRw onEdit={false} value={selectedTypeCollection} />
    </div>
    {#if tabDatas.selected === "Upload"}
      {#if isFolder}
        <UploadFolder
          title={$_("uploadDocuments")}
          handleFileUploadWrapper={() => handleFileUploadWrapper(isFolder=true)}
          files = {groupedFiles}
          bind:isProcess
          multiple={true}
          actionButton={true}
          actionView={true}
          {actionDelete}
        />
      {/if}
      {#if fileNotSupport}
        <Alert color="red">
          <InfoCircleSolid slot="icon" class="w-5 h-5" />
          {$_("fileNotSupport")}
        </Alert>
      {:else if hasDocOrDocx}
        <Alert color="yellow">
          <InfoCircleSolid slot="icon" class="w-5 h-5" />
          {$_("docOrDocxWillConvertToPdf")}
        </Alert>
      {/if}
    {/if}
    {#if tabDatas.selected === "Link"}
      <div class="flex flex-col gap-1.5">
        <Typography
          preset="Text sm/medium"
          color="secondary"
          text={$_("importLink")}
        />
        <InputRw
          onEdit={true}
          bind:value={tmpDocumentLink}
          placeholder={$_("descImportLink")}
        />
      </div>
      {#if stepLink === "notFound"}
        <Helper color="red">
          <span class="font-medium"
            >{$_("File not found. Make sure the link is correct")}</span
          >
        </Helper>
      {:else if stepLink === "checked"}
        <UploadFile
          title={$_("uploadedDocument")}
          handleFileUploadWrapper={() => handleFileUploadWrapper()}
          file={files[0]}
          {files}
          bind:isProcess
          multiple={false}
          actionView={true}
          actionButton={true}
          {actionDelete}
        />
      {/if}
    {/if}
  </div>
  <svelte:fragment slot="footer">
    {#if tabDatas.selected === "Link" && stepLink === "input"}
      <Button
        onclick={() => (stepLink = "checked")}
        disabled={!tmpDocumentLink}
        class="w-full">{$_("checkLink")}</Button
      >
    {:else}
      <Button
        onclick={addDocumentByCollection}
        disabled={buttonDisabled}
        class="w-full">{$_("addNewDocument")}</Button
      >
    {/if}
  </svelte:fragment>
</Modal>
