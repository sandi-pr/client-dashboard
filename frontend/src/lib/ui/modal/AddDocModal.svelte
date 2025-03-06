<script>
  import { _ } from "svelte-i18n";
  import MultiSelect from "svelte-multiselect";
  import { Typography, UploadFile, InputRw } from "$lib";
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

<Modal title={$_("addDocuments")} bind:open={addDocumentsModal} size="xs">
  <svelte:fragment slot="header">
    <Typography preset="Text md/Semibold" color="secondary" text={$_("addDocuments")} />
  </svelte:fragment>
  <div class="flex flex-col gap-4">
    <div>
      <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} spaceTop="mt-0 !mx-0" />
    </div>
    <div
      class="flex flex-col gap-1.5 relative"
    >
      <Typography
        preset="Text sm/medium"
        color="secondary"
        text={$_("documentType")}
      />
      <InputRw
        onEdit={false}
        value={selectedTypeCollection}
      />
    </div>
    {#if tabDatas.selected === "Upload"}
      <UploadFile
        title={$_("uploadDocuments")}
        handleFileUploadWrapper={() => handleFileUploadWrapper()}
        file={files[0]}
        {files}
        bind:isProcess
        multiple={true}
        actionButton={true}
        actionView={true}
        {actionDelete}
      />
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
          <span class="font-medium">{$_("File not found. Make sure the link is correct")}</span>
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
    {#if tabDatas.selected === "Upload" || (tabDatas.selected === "Link" && stepLink === "checked")}
      <div class="flex flex-col gap-1.5">
        <Typography
          preset="Text sm/medium"
          color="secondary"
          text={`${$_("labels")}`}
        />
        <MultiSelect
          --sms-min-height="40px"
          options={alltags}
          allowUserOptions="append"
          bind:selected={tags}
          autoScroll={true}
          noMatchingOptionsMsg=""
          createOptionMsg={null}
          ulSelectedClass="text-primary-light dark:text-primary-dark"
          liSelectedClass="border border-primary-light dark:border-primary-dark"
          ulOptionsClass="bg-secondary-light dark:bg-tertiary-dark"
          outerDivClass="bg-primary-light dark:bg-secondary-dark border rounded-lg border-gray-300 dark:border-gray-600 focus:none focus:ring-0 focus:border-gray-300 dark:focus:border-gray-600"
        />
      </div>
      {#if files.length === 1 && !isProcess}
        <div class="flex flex-col gap-1.5">
          <Typography
            preset="Text sm/medium"
            color="secondary"
            text={`${$_("documentTitle")}*`}
          />
          <InputRw
            onEdit={true}
            bind:value={tmpDocumentTitle}
            placeholder={$_("enterDocumentTitle")}
            oninput={handleInput}
          />
          {#if !nameIsAvailable}
            <Helper color="red">
              <span class="font-medium">{$_("nameAlreadyTaken")}</span>
            </Helper>
          {/if}
        </div>
        <div class="flex flex-col gap-1.5">
          <Typography
            preset="Text sm/medium"
            color="secondary"
            text={`${$_("description")}`}
          />
          <InputRw
            onEdit={true}
            bind:value={tmpDocumentDescription}
            placeholder={$_("enterDocumentDescription")}
          />
        </div>
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
