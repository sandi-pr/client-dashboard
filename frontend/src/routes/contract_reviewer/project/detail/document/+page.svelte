<script>
  import {
    Typography,
    FeaturedIcon,
    InputRw,
    TextareaRw,
    UploadFile,
    Dropdown,
    triggerLoading,
    triggerToast,
  } from "$lib";
  import { Button, P } from "flowbite-svelte";
  import Divider from "$lib/components/Divider.svelte";
  import FilePlusIcon from "$lib/icon/file-plus-02.svelte";
  import CalendarPlus02 from "$lib/icon/calendar-plus-02.svelte";
  import { DatePicker } from "@svelte-plugins/datepicker";
  import PdfViewer from "svelte-pdf";
  import { format } from "date-fns";
  import { _ } from "svelte-i18n";
  import { onMount, onDestroy } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  const call = new CR3Call(globalThis.location.origin);

  let projectId;
  let projectDocumentId;

  let onEdit = $state(true);
  let onView = $state(false);
  let buttonDisabled = $state(false);
  let isOpen = $state(false);
  let docTypeOpen = $state(false);
  let startDate = $state(null);
  let endDate = $state(null);
  let file = $state(null);
  let dateFormat = $state("d MMM yyyy");
  let formattedStartDate = $state();
  let formattedEndDate = $state("");
  let tmpDocumentType = $state("Documents");
  let tmpSelected = $state("Documents");
  let selectedTypeCollection = $state("Documents");
  let acceptedFileTypes = $state(".txt, .pdf, .doc, .docx");
  let fileUrl = $state("");
  let documentTypes = $state([
    {
      name: "Documents (TXT, PDF, DOC, DOCX)",
      value: "Documents",
    },
    {
      name: "Images (JPG, JPEG, PNG, WEBP, IMG)",
      value: "Images",
    },
    {
      name: "Handwriting (JPG, JPEG, PNG)",
      value: "Handwriting",
    },
  ]);

  let data = $state({
    document_name: "",
    content_date: "",
    description: "",
  });

  const selectTypeCollection = (type) => {
    tmpDocumentType = type;

    if (type === "Documents") {
      selectedTypeCollection = "Documents";
      acceptedFileTypes = ".txt, .pdf, .doc, .docx";
    } else if (type === "Images") {
      selectedTypeCollection = "Images";
      acceptedFileTypes = ".jpg, .jpeg, .png, .webp, .img";
    } else if (type === "Handwriting") {
      selectedTypeCollection = "Handwriting";
      acceptedFileTypes = ".jpg, .jpeg, .png";
    }

    if (tmpSelected !== selectedTypeCollection) {
      file = null;
      tmpSelected = selectedTypeCollection;
    }
  };

  async function readFileAsBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        const base64Data = reader.result.split(",")[1];
        resolve(base64Data);
      };
      reader.onerror = reject;

      reader.readAsDataURL(file);
    });
  }

  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  const handleFileUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = acceptedFileTypes;

    inputElement.addEventListener("change", async (event) => {
      const uploadedFile = event.target.files[0];

      if (uploadedFile) {
        const base64Data = await readFileAsBase64(uploadedFile);
        file = {
          type: data.document_name,
          role: tmpDocumentType,
          name: uploadedFile.name,
          size: formatBytes(uploadedFile.size),
          base64: base64Data,
        };
      } else {
        alert("Invalid file. Please select a valid file.");
      }
    });

    inputElement.click();
  };

  const toggleDatePicker = () => {
    if (onView) return;
    isOpen = !isOpen;
  };

  const formatDate = (dateString) => (dateString && format(new Date(dateString), dateFormat)) || "";

  $effect(() => {
    formattedStartDate = formatDate(startDate);
    formattedEndDate = formatDate(endDate);
  });

  async function createDocument() {
    triggerLoading(true);

    let token = sessionStorage.getItem("token");
    if (token) {
      const response = await call.call("client-dashboard", "create-document-project", {
        URLQuery: {},
        URLParams: {},
        Params: {
          project_id: projectId.toString(),
          content_date: format(formattedStartDate, "yyyy-MM-dd"),
          document: JSON.stringify(file),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerLoading(false);
        triggerToast("error");
        console.log(response);
        return;
      }

      const response2 = await call.call("hexa-ai-contract-review", "process-to-rag", {
        URLQuery: {},
        URLParams: {},
        Params: {
          project_id: projectId.toString(),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response2.status !== 200) {
        triggerLoading(false);
        triggerToast("error");
        console.log(response2);
        return;
      }

      // TODO: save to collection

      triggerToast("created");
      data = {
        document_name: "",
        content_date: "",
        description: "",
      };

      selectedTypeCollection = "Documents";
      acceptedFileTypes = ".txt, .pdf, .doc, .docx";
      file = null;

      setTimeout(async () => {
        triggerLoading(false);
        goto(
          (await call.getPage("client-dashboard", "contract-project-detail")) +
            "?projectId=" +
            projectId,
        );
      }, 2000);
    }
  }

  async function getDocument(id) {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call(
          "client-dashboard",
          "get-contract-project-with-document-id",
          {
            URLQuery: {},
            URLParams: {},
            Params: {
              id: parseInt(id),
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          },
        );

        if (response.status !== 200) {
          triggerToast("error");
          console.log(response);
          return;
        }

        file = {
          base64: response.result.base64,
        };

        data = {
          document_name: response.result.document.data[0].attributes.document_name,
          content_date: response.result.document.data[0].attributes.content_date,
          description: response.result.document.data[0].attributes.description,
          ext: response.result.document.data[0].attributes.document_media_id.data.attributes.ext,
        };

        startDate = new Date(response.result.document.data[0].attributes.content_date);
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  onMount(() => {
    formattedStartDate = formatDate(new Date());
    const unsubscribe = page.subscribe(async ($page) => {
      projectId = $page.url.searchParams.get("projectId");
      projectDocumentId = $page.url.searchParams.get("projectDocumentId");
      if (projectId !== null) {
        onEdit = true;
        onView = false;
      } else if (projectDocumentId !== null) {
        onEdit = false;
        onView = true;
        getDocument(projectDocumentId);
      } else {
        triggerToast("error");
        goto(await call.getPage("client-dashboard", "contract-project-list"));
      }
    });
    onDestroy(() => unsubscribe());
  });
</script>

<div class="flex flex-col px-4 pt-6 gap-6 pb-12">
  <!-- section header -->
  <div class="flex justify-between">
    {#if !onView}
      <div class="flex items-center gap-3">
        <FeaturedIcon classesSvg={FilePlusIcon} />
        <div class="flex flex-col">
          <Typography preset="Text lg/Semibold" text={$_("addNewDocument")} color="primary" />
          <Typography preset="Text sm/Regular" text={$_("descAddNewDocument")} color="tertiary" />
        </div>
      </div>
      <div class="flex gap-3">
        <Button outline on:click={() => {}} color="light">{$_("cancel")}</Button>
        <Button disabled={file === null} on:click={() => {}}>{$_("add")}</Button>
      </div>
    {:else}
      <div class="flex items-center gap-3">
        <Typography preset="Display sm/Semibold" text={data.document_name} color="primary" />
      </div>
    {/if}
  </div>
  <Divider />
  <!-- section form -->
  <div class="flex gap-3 w-full">
    <div class="w-60">
      <Typography preset="Text md/Semibold" text={$_("documentInformation")} color="secondary" />
    </div>
    <div class="flex grow flex-col gap-3">
      <div class="flex gap-4">
        <div class="flex flex-col grow gap-1.5">
          <Typography preset="Text sm/Medium" text={$_("documentName")} color="secondary" />
          <InputRw {onEdit} bind:value={data.document_name} />
        </div>
        <div class="flex grow flex-col gap-1.5">
          <Typography
            preset="Text sm/medium"
            color="secondary"
            text={`${$_("projectDuration")}*`}
          />
          <DatePicker bind:isOpen bind:startDate enableFutureDates>
            <button
              onclick={toggleDatePicker}
              class:open={isOpen}
              class="flex w-full p-2.5 border rounded-lg border-gray-300 dark:border-gray-500 gap-2 justify-center items-center"
            >
              <svelte:component
                this={CalendarPlus02}
                svgClass="w-4 h-4 ms-2 text-secondary-light dark:text-secondary-dark"
              />
              <div class="date text-secondary-light dark:text-secondary-dark">
                {#if startDate}
                  {formattedStartDate}
                {:else}
                  {$_("selectDate")}
                {/if}
              </div>
            </button>
          </DatePicker>
        </div>
      </div>
      <div class="flex flex-col grow gap-1.5">
        <Typography preset="Text sm/Medium" text={$_("description")} color="secondary" />
        <TextareaRw {onEdit} bind:value={data.description} />
      </div>
      {#if !onView}
        <UploadFile
          title={`${$_("documentUpload")}*`}
          handleFileUploadWrapper={() => handleFileUploadWrapper()}
          {file}
          actionButton={true}
          actionDelete={() => (file = null)}
          isFlex
        />
      {/if}
    </div>
  </div>
  <Divider />
  <!-- section form -->
  <div class="flex gap-3">
    <div class="w-60">
      <Typography preset="Text md/Semibold" text={`${$_("documentType")}`} color="secondary" />
    </div>
    <div class="flex grow flex-col gap-3">
      <div class="flex flex-col grow gap-1.5">
        <Typography preset="Text sm/Medium" text={$_("documentType")} color="secondary" />
        <Dropdown
          dropdownOpen={docTypeOpen}
          bind:selectedItem={selectedTypeCollection}
          items={documentTypes}
          placeholder={$_("selectDocumentType")}
          isOverlay={true}
          onClick={selectTypeCollection}
        />
      </div>
    </div>
  </div>
  <Divider />
  {#if !onView}
    <div class="flex justify-end gap-3">
      <Button
        outline
        on:click={async () => {
          goto(
            (await call.getPage("client-dashboard", "contract-project-detail")) +
              "?projectId=" +
              projectId,
          );
        }}
        color="light">{$_("cancel")}</Button
      >
      <Button
        disabled={file === null}
        on:click={() => {
          createDocument();
        }}>{$_("add")}</Button
      >
    </div>
  {/if}
  {#if onView}
    <div class="flex gap-6">
      <div class="w-60">
        <P class="font-semibold text-secondary-light dark:text-secondary-dark">
          {$_("viewer")}
        </P>
      </div>
      <div class="flex flex-col grow gap-6 rounded-lg bg-secondary-light dark:bg-secondary-dark">
        {#if data.ext === ".pdf"}
          <!-- <iframe class="w-full h-screen" src={`${fileUrl}`}></iframe> -->
          <div class="flex justify-center">
            <!-- {#if pdfReady} -->
            <PdfViewer
              showBorder={false}
              showButtons={["navigation", "zoom", "rotate", "autoflip", "pageInfo"]}
              scale={1}
              data={atob(file.base64)}
            />
            <!-- {/if} -->
          </div>
        {:else if data.ext === ".jpg" || data.ext === ".png" || data.ext === ".jpeg" || data.ext === ".webp"}
          <div class="flex justify-center">
            <img class="w-full h-screen" src={`${fileUrl}`} alt="imageViewer" />
          </div>
          <!-- {:else}
          {#each contentFile.split("\n") as paragraphs, i}
            <p class="dark:text-white">{paragraphs}</p>
          {/each} -->
        {/if}
      </div>
    </div>
  {/if}
</div>
