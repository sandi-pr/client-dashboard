<script>
  import { onMount } from "svelte";
  import {
    Table,
    Modal,
    Button,
    TableHead,
    TableBody,
    TableBodyRow,
    TableHeadCell,
    TableBodyCell,
  } from "flowbite-svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import MultiSelect from "svelte-multiselect";
  import { _, locale } from "svelte-i18n";
  import { id as localeId, enUS as localeEn } from "date-fns/locale";
  import { formatDistanceToNow } from "date-fns";

  import {
    InputRw,
    Divider,
    TextareaRw,
    Typography,
    Pagination,
    FeaturedIcon,
    ToastContainer,
  } from "$lib"; // Library for components
  import {
    call,
    getUser,
    getIcons,
    getLabels,
    triggerToast,
    truncateText,
    handleInputName,
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
  const StarsIcon = getIcons("stars");

  let {
    data = $bindable(),
    tags = $bindable(),
    CSVMetadata = $bindable(),
    CSVData = $bindable(),
    onEdit,
    base64,
    fileUrl,
    fileReady,
    libraryId,
    documentId,
    defaultData,
    collectionId,
    pageDataType,
    attributeNotUnique,
    updateData = () => {},
    deleteData = () => {},
    CSVNext = () => {},
    CSVPrevious = () => {},
  } = $props();

  let innerWidth = $state(0);
  let uploadedBy = $state("");
  let documentTitle = $state("");
  let dropdownAction = $state(false);
  let modalDelete = $state(false);
  let modalDetail = $state(false);
  let modalDetectedText = $state(false);
  let nameIsAvailable = $state(true);
  let alltags = $state([]);

  const handleDropdownFocusLoss = () => {
    setTimeout(() => {
      dropdownAction = false;
    }, 200);
  };

  const getTROCRPerDataType = (dataType, infoType) => {
    const dataTypeMapping = {
      Handwriting: {
        info: $_("handwritingInformation"),
        text: $_("detectedText"),
      },
      "Document Images": {
        info: $_("documentInformation"),
        text: $_("documentText"),
      },
      "Document Audio": {
        info: $_("audioInformation"),
        text: $_("detectedAudio"),
      },
    };

    return dataTypeMapping[dataType]?.[infoType] || null;
  };

  const getDataTypeTextPerExt = (ext) => {
    const extMapping = {
      image: {
        extensions: [".jpg", ".jpeg", ".png", ".webp"],
        text: $_("image") + " (JPG, JPEG, PNG, WEBP)",
      },
      audio: {
        extensions: [".mp3", ".wav", ".ogg"],
        text: $_("audio") + " (MP3, WAV, OGG)",
      },
      pdf: {
        extensions: [".pdf"],
        text: $_("handwriting") + " (PDF)",
      },
      csv: {
        extensions: [".csv"],
        text: $_("data") + " (CSV)",
      },
    };

    for (const key in extMapping) {
      if (extMapping[key].extensions.includes(ext)) {
        return extMapping[key].text;
      }
    }

    return $_("unknownFileType");
  };

  onMount(async () => {
    const resizeHandler = () => {
      innerWidth = globalThis.innerWidth;
    };
    window.addEventListener("resize", resizeHandler);
    resizeHandler();

    alltags = await getLabels($call);

    return () => {
      window.removeEventListener("resize", resizeHandler);
    };
  });

  $effect(async () => {
    if (Object.keys(data).length > 0) {
      uploadedBy = await getUser($call, data);
    }
  });

  $effect(() => {
    if (innerWidth < 1024) {
      documentTitle = truncateText(defaultData.document_title || "", 2);
    } else {
      documentTitle = defaultData.document_title;
    }
  });
</script>

<div class="flex">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div class="flex flex-col grow gap-1">
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
          <Typography
            preset="Text sm/Medium"
            color="tertiary"
            text={fileReady
              ? $_("createOn") +
                " " +
                $_(
                  formatDistanceToNow(data.createdAt, {
                    addSuffix: true,
                    locale: $locale === "id" ? localeId : localeEn,
                  }),
                )
              : ""}
          />
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
              <Typography preset="Text sm/Semibold" color="error-primary" text={$_("delete")} />
            </button>
          </div>
        </div>
      </div>

      <!-- Content -->
      <div class="flex flex-col gap-4 pt-4 md:pt-6 md:gap-6">
        {#if fileReady}
          <div class="flex flex-col md:flex-row gap-6">
            <div class="w-60">
              <Typography
                preset="Text md/Semibold"
                color="secondary"
                text={[".mp3", ".wav", ".ogg"].includes(data.ext) ? $_("listener") : $_("viewer")}
              />
            </div>
            {#if [".mp3", ".wav", ".ogg"].includes(data.ext)}
              {#if fileUrl}
                <AudioPlayer src={fileUrl} />
              {/if}
            {:else if [".csv"].includes(data.ext)}
              <div
                class="bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-700 rounded-xl pt-4"
              >
                <Table class="w-full">
                  <TableHead>
                    {#each Object.keys(CSVData[0]) as key}
                      {#if key === "id"}
                        <!-- empty -->
                      {:else}
                        <TableHeadCell>{key.replace(/_/g, " ")}</TableHeadCell>
                      {/if}
                    {/each}
                  </TableHead>
                  <TableBody tableBodyClass="divide-y">
                    {#each CSVData as data}
                      <TableBodyRow class="border-0 !border-gray-200 dark:!border-gray-700">
                        {#each Object.entries(data) as [key, value]}
                          {#if key === "id"}
                            <!-- empty -->
                          {:else}
                            <TableBodyCell tdClass="px-6 py-4 font-medium">{value}</TableBodyCell>
                          {/if}
                        {/each}
                      </TableBodyRow>
                    {/each}
                  </TableBody>
                </Table>
                <Pagination
                  pagination={CSVMetadata}
                  previousPage={CSVPrevious}
                  nextPage={CSVNext}
                />
              </div>
            {:else}
              <div
                class={`flex flex-col relative ${data.document_type === "Document Audio" ? "items-start ps-8" : "items-center"} grow py-4 md:py-12 rounded-lg bg-secondary-light dark:bg-secondary-dark`}
              >
                {#if data.ext === ".pdf"}
                  {#if fileReady}
                    <PdfjsViewer pdfFile={base64} />
                  {/if}
                {:else if data.ext === ".jpg" || data.ext === ".png" || data.ext === ".jpeg" || data.ext === ".webp"}
                  <div class="flex justify-center">
                    <img class="w-full h-screen" src={`${fileUrl}`} alt="imageViewer" />
                  </div>
                {:else}
                  {#each contentFile.split("\n") as paragraphs, i}
                    <p class="dark:text-white">{paragraphs}</p>
                  {/each}
                {/if}
              </div>
            {/if}
          </div>
          {#if ["Handwriting", "Document Images", "Document Audio"].includes(data.document_type)}
            <Divider />
            <div class="flex flex-col md:flex-row gap-4 md:gap-6">
              <div class="flex flex-col w-60 gap-2">
                <Typography
                  preset="Text md/Semibold"
                  color="secondary"
                  text={getTROCRPerDataType(data.document_type, "info")}
                />
                <button
                  onclick={() => {
                    modalDetectedText = true;
                    onEdit = false;
                  }}
                >
                  <Typography
                    preset="Text sm/Semibold"
                    color="fg-brand-primary"
                    text={$_("wrongDetectedTextCorrectHere")}
                  />
                </button>
              </div>
              <div class="flex flex-col grow gap-4 md:gap-6">
                <div class="grid grid-cols-1">
                  <div class="flex flex-col gap-2">
                    <Typography
                      preset="Text sm/Medium"
                      text={getTROCRPerDataType(data.document_type, "text")}
                    />
                    <TextareaRw {onEdit} bind:value={data.detectedText} />
                  </div>
                </div>
              </div>
            </div>
          {/if}
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
      <Typography preset="Text lg/Semibold" text={$_("generalInformation")} />
    </div>
  </svelte:fragment>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col gap-1">
      <div class="flex flex-col gap-2">
        <Typography preset="Text sm/Medium" color="secondary" text={`${$_("documentTitle")}`} />
        <InputRw
          {onEdit}
          oninput={async (evt) => {
            const res = await handleInputName(
              $call,
              evt,
              defaultData.document_title,
              collectionId,
              pageDataType,
            );

            if (res !== undefined) {
              nameIsAvailable = res;
            }
          }}
          bind:value={data.document_title}
        />
      </div>
      {#if !nameIsAvailable}
        <Typography preset="Text xs/Thin" color="error-primary" text={$_("nameAlreadyTaken")} />
      {/if}
    </div>
    <div class="flex flex-col gap-2">
      <Typography preset="Text sm/Medium" color="secondary" text={`${$_("description")}`} />
      <TextareaRw {onEdit} bind:value={data.description} />
    </div>
    <div class="flex flex-col gap-2">
      <Typography preset="Text sm/Medium" color="secondary" text={`${$_("documentType")}`} />
      <InputRw onEdit={false} value={getDataTypeTextPerExt(data.ext)} />
    </div>
    <div class="flex flex-col gap-2">
      <Typography preset="Text sm/Medium" color="secondary" text={`${$_("labels")}`} />
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
      <Typography preset="Text sm/Medium" color="secondary" text={$_("uploadBy")} />
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
          <Typography preset="Text md/Semibold" color="fg-secondary" text={$_("cancel")} />
        </Button>
        <Button
          class="w-full gap-2"
          disabled={!nameIsAvailable || attributeNotUnique}
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
          <Typography preset="Text md/Semibold" color="primary-onbrand" text={$_("save")} />
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
            preset="Text md/Semibold"
            color="fg-secondary"
            text={$_("editInformations")}
          />
        </Button>
      {/if}
    </div>
  </svelte:fragment>
</Modal>

<Modal
  bind:open={modalDetectedText}
  onclose={() => {
    data = { ...defaultData };
  }}
  outsideclose="true"
>
  <svelte:fragment slot="header">
    <div class="flex items-center gap-2">
      <FeaturedIcon classesSvg={StarsIcon} />
      <Typography preset="Text lg/Semibold" text={$_("AIHandwritingInformation")} />
    </div>
  </svelte:fragment>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col gap-2">
      <Typography preset="Text sm/Medium" color="secondary" text={`${$_("detectedText")}`} />
      <TextareaRw {onEdit} bind:value={data.detectedText} />
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
          }}
        >
          <Typography preset="Text md/Semibold" color="fg-secondary" text={$_("cancel")} />
        </Button>
        <Button
          class="w-full gap-2"
          onclick={() => {
            updateData();
            onEdit = false;
            modalDetectedText = false;
          }}
        >
          <Save01Icon
            svgClass="w-5 h-5 text-primary-onbrand-light dark:text-primary-onbrand-dark"
          />
          <Typography preset="Text md/Semibold" color="primary-onbrand" text={$_("save")} />
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
            preset="Text md/Semibold"
            color="fg-secondary"
            text={$_("editDetectedText")}
          />
        </Button>
      {/if}
    </div>
  </svelte:fragment>
</Modal>
