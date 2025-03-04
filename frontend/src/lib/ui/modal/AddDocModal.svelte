<script>
  import { _ } from "svelte-i18n";
  import MultiSelect from "svelte-multiselect";
  import { Typography, UploadFile, InputRw } from "$lib";
  import { Modal, Button, Alert } from "flowbite-svelte";
  import { InfoCircleSolid } from "flowbite-svelte-icons";

  let {
    addDocumentsModal = $bindable(false),
    selectedTypeCollection,
    files,
    alltags,
    tags,
    isProcess,
    handleFileUploadWrapper,
    fileNotSupport,
    hasDocOrDocx,
    tmpDocumentTitle,
    tmpDocumentDescription,
    handleInput,
    buttonDisabled,
    addDocumentByCollection,
  } = $props();

</script>

<Modal title={$_("addDocuments")} bind:open={addDocumentsModal} size="xs">
  <svelte:fragment slot="header">
    <Typography preset="Text md/Semibold" color="secondary" text={$_("addDocuments")} />
  </svelte:fragment>
  <div class="flex flex-col gap-4">
    <div
      class="flex flex-col gap-1.5 relative"
    >
      <Typography
        preset="Text sm/medium"
        color="secondary"
        text={`${$_("documentType")}*`}
      />
      <InputRw
        onEdit={false}
        value={selectedTypeCollection}
      />
    </div>
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
    <UploadFile
      title={`${$_("uploadDocuments")}*`}
      handleFileUploadWrapper={() => handleFileUploadWrapper()}
      file={files[0]}
      {files}
      bind:isProcess
      multiple={true}
      actionButton={true}
      actionDelete={() => (files = [])}
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
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={addDocumentByCollection}
      disabled={buttonDisabled}
      class="w-full">{$_("addNewDocument")}</Button
    >
  </svelte:fragment>
</Modal>
