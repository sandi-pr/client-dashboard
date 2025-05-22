<script>
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { P, Button, Modal, Helper } from "flowbite-svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import MultiSelect from "svelte-multiselect";
  import { _, locale } from "svelte-i18n";
  import { id as localeId, enUS as localeEn } from "date-fns/locale";
  import { formatDistanceToNow } from "date-fns";

  import { InputRw, TextareaRw, Typography, ToastContainer, Divider } from "$lib"; // Library for components
  import {
    call,
    check,
    getIcons,
    triggerToast,
    downloadFileBase64,
    downloadFileBlob,
    triggerLoading,
  } from "$lib"; // Library for functions
  import PdfjsViewer from "$lib/components/PdfjsViewer.svelte";
  import AudioPlayer from "$lib/components/AudioPlayer.svelte";

  const Copy01Icon = getIcons("copy01");

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
    popupModal,
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
  let alltags = $state([]);

  const getLabels = async () => {
    let token = localStorage.getItem("token");
    if (token !== null) {
      let response = await $call.call("client-dashboard", "get-label", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });
      if (response.status !== 200) {
        if (response.status === 400) {
          console.error("Invalid parameters:", response);
          triggerToast($_("invalidParameters"), true, "error");
        } else if (response.status === 401) {
          console.error("Unauthorized access:", response);
          triggerToast($_("unauthorized"), true, "error");
        } else if (response.status === 500) {
          console.error("Server error:", response);
          triggerToast($_("serverError"), true, "error");
        } else {
          console.error("Unexpected status code:", response);
          triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
        }
        return;
      }

      alltags = response.result.labels;
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  };

  const getUser = async (userData) => {
    let token = localStorage.getItem("token");
    try {
      let response = await $call.call("client-dashboard", "get-user", {
        URLQuery: {},
        URLParams: { id: userData.upload_by.toString() },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        if (response.status === 400) {
          console.error("Invalid parameters:", response);
          triggerToast($_("invalidParameters"), true, "error");
        } else if (response.status === 401) {
          console.error("Unauthorized access:", response);
          triggerToast($_("unauthorized"), true, "error");
        } else if (response.status === 500) {
          console.error("Server error:", response);
          triggerToast($_("serverError"), true, "error");
        } else {
          console.error("Unexpected status code:", response);
          triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
        }
        return;
      }

      uploadedBy = response?.result?.user?.username || "Unknown User";
    } catch (error) {
      console.error("Error during user get:", error);
      triggerToast($_("failedGetUser"), true, "error");
    }
  };

  onMount(async () => {
    await getLabels();
  });

  $effect(async () => {
    if (Object.keys(data).length > 0) {
      await getUser(data);
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
        <div class="flex gap-2">
          {#if onEdit}
            <Button
              onclick={() => {
                onEdit = false;
                data = { ...defaultData };
                attributeNotUnique = false;
                nameIsAvailable = true;
              }}
              color="light"
              class="px-3 focus-within:!ring-1"
            >
              {$_("cancel")}
            </Button>
            <Button
              onclick={() => {
                updateData();
                onEdit = false;
                attributeNotUnique = false;
                nameIsAvailable = true;
              }}
              disabled={data.document_title === "" || !nameIsAvailable}
              class="px-3 focus-within:!ring-1"
            >
              {$_("save")}
            </Button>
          {:else}
            <Button
              onclick={async function () {
                triggerLoading(false);
                goto(
                  (isFrontend
                    ? "/documents/libraries/collections/documents"
                    : await $call.getPage("client-dashboard", "document-library-documents")) +
                    `?collectionId=${collectionId}&libraryId=${libraryId}`,
                );
                return;
              }}
              color="light"
              class="hidden md:block px-3 focus-within:!ring-1"
            >
              {$_("back")}
            </Button>
            <Button
              disabled={!check("documentAi.documentLibrary.document:delete")}
              onclick={() => (popupModal = true)}
              outline
              color="red"
              class="px-3 focus-within:!ring-1"
            >
              {$_("delete")}
            </Button>
            <Button
              disabled={!check("documentAi.documentLibrary.document:update")}
              onclick={() => {
                onEdit = true;
              }}
              class="px-3 focus-within:!ring-1"
            >
              {$_("edit")}
            </Button>
            <Button
              onclick={() => {
                if (pageDataType === "document") {
                  downloadFileBase64(base64, data.document_title, data.ext);
                } else {
                  downloadFileBlob(fileUrl, data.document_title);
                }
              }}
              class="px-3 focus-within:!ring-1"
            >
              {$_("download")}
            </Button>
          {/if}
        </div>
      </div>

      <!-- Content -->
      <div class="flex flex-col gap-4 pt-4 md:pt-6 md:gap-6">
        {#if fileReady}
          <!-- general information -->
          <div class="flex flex-col md:flex-row gap-3 md:gap-6">
            <div class="w-60">
              <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                {$_("generalInformation")}
              </P>
            </div>
            <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
              <div class="flex flex-col gap-1.5">
                <Typography
                  preset="Text sm/medium"
                  color="secondary"
                  text={`${$_("documentTitle")}*`}
                />
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
              <div class="flex flex-col gap-1.5">
                <Typography
                  preset="Text sm/medium"
                  color="secondary"
                  text={`${$_("description")}`}
                />
                <InputRw {onEdit} bind:value={data.description} />
              </div>
              <div class="flex flex-col gap-1.5">
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
                  outerDivClass="bg-primary-light dark:bg-secondary-dark border-primary-light dark:border-primary-dark text-primary-light dark:text-primary-dark rounded-lg"
                  liOptionClass="bg-primary-light dark:bg-secondary-dark hover:bg-secondary-light hover:dark:bg-primary-dark text-primary-light dark:text-primary-dark"
                  liSelectedClass="bg-tertiary-light dark:bg-primary-dark"
                />
              </div>
              <div class="flex flex-col gap-1.5">
                <Typography preset="Text sm/medium" color="secondary" text={$_("uploadBy")} />
                <InputRw onEdit={false} value={uploadedBy} />
              </div>
            </div>
          </div>
          <Divider />
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

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteData} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>
