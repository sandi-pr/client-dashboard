<script>
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { P, Button, Modal, Helper } from "flowbite-svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import MultiSelect from "svelte-multiselect";
  import { _, locale } from "svelte-i18n";
  import { id as localeId, enUS as localeEn } from "date-fns/locale";
  import { formatDistanceToNow } from "date-fns";

  import { InputRw, TextareaRw, Typography, ToastContainer, Divider, FeaturedIcon } from "$lib"; // Library for components
  import {
    call,
    check,
    getUser,
    getIcons,
    getLabels,
    triggerToast,
    triggerLoading,
    downloadFileBlob,
    downloadFileBase64,
  } from "$lib"; // Library for functions
  import PdfjsViewer from "$lib/components/PdfjsViewer.svelte";
  import AudioPlayer from "$lib/components/AudioPlayer.svelte";

  const Copy01Icon = getIcons("copy01");
  const DotsVerticalIcon = getIcons("dotsVertical");
  const Download02Icon = getIcons("download02");
  const InfoCircleIcon = getIcons("infoCircle");
  const Edit05Icon = getIcons("edit05");
  const Trash01Icon = getIcons("trash01");
  const Save01Icon = getIcons("save01");

  const isFrontend = process.env.NODE_ENV === "development";

  let {
    data = $bindable(),
    tags = $bindable(),
    documentTitle = $bindable(),
    onEdit,
    base64,
    fileUrl,
    pdfReady,
    fileReady,
    libraryId,
    documentId,
    defaultData,
    collectionId,
    pageDataType,
    nameIsAvailable,
    attributeNotUnique,
    updateData = () => {},
    deleteData = () => {},
    handleInputName = () => {},
  } = $props();

  let uploadedBy = $state("");
  let dropdownAction = $state(false);
  let alltags = $state([]);
  let modalDelete = $state(false);
  let modalDetail = $state(true);

  const handleDropdownFocusLoss = () => {
    setTimeout(() => {
      dropdownAction = false;
    }, 200);
  };

  onMount(async () => {
    alltags = await getLabels($call);
  });

  $effect(async () => {
    if (Object.keys(data).length > 0) {
      uploadedBy = await getUser($call, data);
    }
  });
</script>

<div class="flex">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div class="flex flex-col grow">
          <div class="flex gap-2">
            <Typography preset="Display sm/Semibold" text={documentTitle} />
            <button
              onclick={() => {
                navigator.clipboard
                  .writeText(documentTitle)
                  .then(() => {
                    triggerToast("copyToClipboard");
                  })
                  .catch((err) => {
                    console.error("Failed to copy: ", err);
                  });
              }}
              class="p-2 border-2 rounded-lg border-secondary-light dark:border-secondary-dark"
            >
              <Copy01Icon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
            </button>
          </div>
          <P
            class="fluid-sm md:fluid-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
          >
            {#if fileReady}
              {#if $locale === "id"}
                {$_("createOn")}
                {formatDistanceToNow(data.createdAt, { addSuffix: true, locale: localeId })}
              {:else}
                {$_("createOn")}
                {formatDistanceToNow(data.createdAt, { addSuffix: true, locale: localeEn })}
              {/if}
            {/if}
          </P>
        </div>
        <div onfocusout={handleDropdownFocusLoss} class="flex items-center relative">
          <button
            class="p-2 rounded-lg bg-secondary-light dark:bg-secondary-dark border border-secondary-light dark:border-secondary-dark"
            onclick={() => {
              dropdownAction = !dropdownAction;
            }}
          >
            <DotsVerticalIcon
              svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
            />
          </button>
          <div
            class="flex flex-col z-10 rounded-lg shadow w-max bg-primary-light dark:bg-primary-dark p-2 {dropdownAction
              ? 'absolute top-10 right-0'
              : 'hidden'}"
          >
            <button
              class="flex items-center gap-2 p-2 rounded-lg hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
              onclick={() => {
                if (pageDataType === "document") {
                  downloadFileBase64(base64, data.document_title, data.ext);
                } else {
                  downloadFileBlob(fileUrl, data.document_title);
                }
                dropdownAction = false;
              }}
            >
              <Download02Icon
                svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
              />
              <Typography preset="Text sm/Semibold" color="fg-tertiary" text={$_("download")} />
            </button>
            <button
              class="flex items-center gap-2 p-2 rounded-lg hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
              onclick={() => {
                onEdit = false;
                modalDetail = true;
                dropdownAction = false;
              }}
            >
              <InfoCircleIcon
                svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
              />
              <Typography preset="Text sm/Semibold" color="fg-tertiary" text={$_("details")} />
            </button>
            <button
              class="flex items-center gap-2 p-2 rounded-lg hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
              onclick={() => {
                onEdit = true;
                modalDetail = true;
                dropdownAction = false;
              }}
            >
              <Edit05Icon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
              <Typography preset="Text sm/Semibold" color="fg-tertiary" text={$_("edit")} />
            </button>
            <button
              class="flex items-center gap-2 p-2 rounded-lg hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
              onclick={() => {
                modalDelete = true;
                dropdownAction = false;
              }}
            >
              <Trash01Icon
                svgClass="w-5 h-5 text-error-primary-light dark:text-error-primary-dark"
              />
              <Typography preset="Text sm/semibold" color="error-primary" text={$_("delete")} />
            </button>
          </div>
        </div>
      </div>

      <!-- Content -->
      <div class="flex flex-col gap-4 pt-4 md:pt-6 md:gap-6">
        {#if fileReady}
          {#if data.document_type === "Handwriting" || data.document_type === "Document Images" || data.document_type === "Document Audio"}
            <div class="flex flex-col md:flex-row gap-4 md:gap-6">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {#if data.document_type === "Handwriting"}
                    {$_("handwritingInformation")}
                  {:else if data.document_type === "Document Images"}
                    {$_("documentInformation")}
                  {:else if data.document_type === "Document Audio"}
                    {$_("audioInformation")}
                  {/if}
                </P>
              </div>
              <div class="flex flex-col grow gap-4 md:gap-6">
                <div class="grid grid-cols-1">
                  <div class="flex flex-col gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {#if data.document_type === "Handwriting"}
                        {$_("handwritingText")}
                      {:else if data.document_type === "Document Images"}
                        {$_("documentText")}
                      {:else if data.document_type === "Document Audio"}
                        {$_("audioText")}
                      {/if}</P
                    >
                    <TextareaRw {onEdit} bind:value={data.detectedText} />
                  </div>
                </div>
              </div>
            </div>
            <Divider />
          {/if}
          <div class="flex flex-col md:flex-row gap-6">
            <div class="w-60">
              <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                {#if data.ext === ".mp3" || data.ext === ".wav" || data.ext === ".ogg"}
                  {$_("audioPreview")}
                {:else}
                  {$_("viewer")}
                {/if}
              </P>
            </div>
            <div
              class={`flex flex-col relative ${data.document_type === "Document Audio" ? "items-start ps-8" : "items-center"} grow py-4 md:py-12 rounded-lg bg-secondary-light dark:bg-secondary-dark`}
            >
              {#if data.ext === ".pdf"}
                {#if pdfReady}
                  <PdfjsViewer pdfFile={base64} />
                {/if}
              {:else if data.ext === ".jpg" || data.ext === ".png" || data.ext === ".jpeg" || data.ext === ".webp"}
                <div class="flex justify-center">
                  <img class="w-full h-screen" src={`${fileUrl}`} alt="imageViewer" />
                </div>
              {:else if data.ext === ".mp3" || data.ext === ".wav" || data.ext === ".ogg"}
                {#if fileUrl}
                  <AudioPlayer src={fileUrl} />
                {/if}
              {:else}
                {#each contentFile.split("\n") as paragraphs, i}
                  <p class="dark:text-white">{paragraphs}</p>
                {/each}
              {/if}
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
  <ToastContainer />
</div>

<Modal bind:open={modalDelete} autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteData} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>

<Modal
  bind:open={modalDetail}
  onclose={() => {
    data = { ...defaultData };
    attributeNotUnique = false;
    nameIsAvailable = true;
  }}
  outsideclose="true"
>
  <svelte:fragment slot="header">
    <div class="flex items-center gap-2">
      <FeaturedIcon classesSvg={InfoCircleIcon} />
      <Typography preset="Text lg/semibold" text={$_("generalInformation")} />
    </div>
  </svelte:fragment>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col gap-2">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("documentTitle")}*`} />
      <InputRw
        {onEdit}
        oninput={(evt) => {
          handleInputName(evt, defaultData.document_title);
        }}
        bind:value={data.document_title}
      />
      {#if !nameIsAvailable}
        <Helper color="red">
          <span class="font-medium">{$_("nameAlreadyTaken")}</span>
        </Helper>
      {/if}
    </div>
    <div class="flex flex-col gap-2">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}`} />
      <TextareaRw {onEdit} bind:value={data.description} />
    </div>
    <div class="flex flex-col gap-2">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("labels")}`} />
      <MultiSelect
        --sms-min-height="40px"
        options={alltags || []}
        disabled={!onEdit}
        allowUserOptions="append"
        bind:selected={tags}
        autoScroll={true}
        noMatchingOptionsMsg=""
        createOptionMsg={null}
        outerDivClass="rounded-lg {onEdit
          ? 'bg-primary-light dark:bg-primary-dark border-primary-light dark:border-primary-dark text-primary-light dark:text-primary-dark'
          : 'bg-disabled-subtle-light dark:bg-disabled-subtle-dark border-disabled-subtle-light dark:border-disabled-subtle-dark text-disabled-light dark:text-disabled-dark'}"
        liOptionClass="bg-primary-light dark:bg-secondary-dark hover:bg-secondary-light hover:dark:bg-primary-dark text-primary-light dark:text-primary-dark"
        liSelectedClass={onEdit
          ? "bg-tertiary-light dark:bg-tertiary-dark"
          : "bg-primary-light dark:bg-primary-dark"}
      />
    </div>
    <div class="flex flex-col gap-2">
      <Typography preset="Text sm/medium" color="secondary" text={$_("uploadBy")} />
      <InputRw onEdit={false} value={uploadedBy} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <div class="w-full flex items-center gap-2">
      {#if onEdit}
        <Button
          color="alternative"
          class="w-full"
          onclick={() => {
            onEdit = false;
            data = { ...defaultData };
            attributeNotUnique = false;
            nameIsAvailable = true;
          }}
        >
          <Typography preset="Text md/semibold" color="fg-secondary" text={$_("cancel")} />
        </Button>
        <Button
          class="w-full gap-2"
          onclick={() => {
            updateData();
            onEdit = false;
            attributeNotUnique = false;
            nameIsAvailable = true;
            modalDetail = false;
          }}
        >
          <Save01Icon
            svgClass="w-5 h-5 text-primary-onbrand-light dark:text-primary-onbrand-dark"
          />
          <Typography preset="Text md/semibold" color="primary-onbrand" text={$_("save")} />
        </Button>
      {:else}
        <Button
          color="alternative"
          class="w-full gap-2"
          onclick={() => {
            onEdit = true;
          }}
        >
          <Edit05Icon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
          <Typography
            preset="Text md/semibold"
            color="fg-secondary"
            text={$_("editInformations")}
          />
        </Button>
      {/if}
    </div>
  </svelte:fragment>
</Modal>
