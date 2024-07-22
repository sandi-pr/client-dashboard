<script>
  import { onMount } from "svelte";
  import { Button, Modal, P } from "flowbite-svelte";
  import { FeaturedIcon, Typography, UploadFile } from "$lib";
  import PdfViewer from "svelte-pdf";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import FilePlusIcon from "$lib/icon/file-plus-02.svelte";
  import FileDownloadIcon from "$lib/icon/file-download-03.svg";
  import { ChevronDownOutline, RefreshOutline } from "flowbite-svelte-icons";
  import templatePdf from "$lib/template.pdf";
  console.log(templatePdf);

  let dropdownToggle = $state([true, false, false]);
  let pdfData = $state(null);
  let documentViewModal = $state(false);
  let selectedIndustry = $state("");
  let files = $state([]);
  let selectedFile = $state(null);
  let requiredFileTypes = $state(["documentGtc", "tradeConfirmation", "dealDoneDocument"]);
  let roles = $state(["buyer", "seller"]);

  let industries = $state([
    { name: "Oil & Gas", value: "oil-gas" },
    { name: "Other", value: "other" },
  ]);

  // const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
  //   if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
  //   dropdownToggle = false;
  // };

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

  const handleFileUploadWrapper = (type, role) => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = ".pdf, .doc, .docx";

    inputElement.addEventListener("change", async (event) => {
      const uploadedFile = event.target.files[0];

      if (uploadedFile) {
        const base64Data = await readFileAsBase64(uploadedFile);
        files.push({
          type: type,
          role: role,
          name: uploadedFile.name,
          size: formatBytes(uploadedFile.size),
          base64: base64Data,
        });
      } else {
        alert("Invalid file. Please select a valid file.");
      }
    });

    inputElement.click();
  };

  function getFileByType(type, role) {
    const index = files.findIndex((file) => file.type === type && file.role === role);
    return index !== -1 ? files[index] : [];
  }

  function deleteFileByType(type, role) {
    const index = files.findIndex((file) => file.type === type && file.role === role);
    if (index !== -1) {
      files.splice(index, 1);
    }
  }

  function areRequiredFilesPresent(files) {
    return roles.every((role) =>
      requiredFileTypes.every((type) =>
        files.some((file) => file.type === type && file.role === role)
      )
    );
  }

  let allRequiredFilesPresent = $state(false);

  $effect(() => {
    allRequiredFilesPresent = areRequiredFilesPresent(files);
  });

  onMount(async () => {
    const response = await fetch(templatePdf);
    const blob = await response.blob();
    const reader = new FileReader();

    reader.onloadend = () => {
      pdfData = reader.result;
    };

    reader.readAsDataURL(blob);
  });
</script>

<div class="flex flex-col px-4 pt-6 gap-6 pb-12">
  <!-- section header -->
  <div class="flex justify-between">
    <div class="flex items-center gap-3">
      <div class="flex flex-col">
        <P class="text-3xl text-primary-light dark:text-primary-dark font-semibold"
          >{$_("aiReview")}</P
        >
        <Typography preset="Text sm/Regular" text={"Nama Project"} color="tertiary" />
      </div>
    </div>
    <div class="flex gap-3">
      <Button outline on:click={() => {}} color="light">
        <RefreshOutline class="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark" />
        &nbsp; Refresh
      </Button>
      <Button on:click={() => {}}>
        <img src={FileDownloadIcon} alt="file download icon" />
        &nbsp;
        {$_("generateConfirmation")}</Button
      >
    </div>
  </div>
  <Divider />
  <div class="flex gap-3">
    <div class="flex flex-col grow gap-6">
      <!-- <div onfocusout={handleDropdownFocusLoss} class="flex flex-col gap-1 w-full"> -->
      <div class="flex flex-col gap-1 w-full">
        <!-- document section -->
        <Button
          on:click={() => (dropdownToggle[0] = !dropdownToggle[0])}
          class="flex justify-between border-primary-light dark:border-primary-dark bg-secondary-light dark:bg-secondary-dark"
          color="light"
        >
          <P class="text-lg font-semibold text-tertiary-light dark:text-tertiary-dark"
            >{$_("documentViewerEditor")}</P
          >
          <ChevronDownOutline
            class="w-8 h-8 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
          />
        </Button>
        {#if dropdownToggle[0]}
          <div class="grid grid-cols-2 gap-4 mb-32">
            <div class="col-span-1 h-[50rem]">
              <P
                class="font-bold py-4 text-brand-tertiary-alt-light dark:text-brand-tertiary-alt-dark"
                >{$_("referenceDocumentGTC")}</P
              >
              <embed
                src={templatePdf + "#toolbar=0"}
                type="application/pdf"
                width="100%"
                height="100%"
              />
            </div>
            <div class="col-span-1 h-[50rem]">
              <P class="font-bold py-4 text-primary-light dark:text-primary-dark"
                >{$_("sellerGTC")}</P
              >
              <embed
                src={templatePdf + "#toolbar=0"}
                type="application/pdf"
                width="100%"
                height="100%"
              />
            </div>
            <div class="col-span-1 h-[50rem]">
              <P
                class="font-bold pt-16 pb-4 text-brand-tertiary-alt-light dark:text-brand-tertiary-alt-dark"
                >{$_("referenceDocumentDealDone")}</P
              >
              <embed
                src={templatePdf + "#toolbar=0"}
                type="application/pdf"
                width="100%"
                height="100%"
              />
            </div>
            <div class="col-span-1 h-[50rem]">
              <P class="font-bold pt-16 pb-4 text-primary-light dark:text-primary-dark"
                >{$_("documentDealDone")}</P
              >
              <embed
                src={templatePdf + "#toolbar=0"}
                type="application/pdf"
                width="100%"
                height="100%"
              />
            </div>
            <!-- {#if pdfData}
              <PdfViewer
                showBorder={false}
                showButtons={["navigation", "zoom", "rotate", "autoflip", "pageInfo"]}
                scale={1}
                data={pdfData}
              />
            {/if} -->
          </div>
        {/if}

        <Divider />

        <!-- first table section -->
        <Button
          on:click={() => (dropdownToggle[1] = !dropdownToggle[1])}
          class="flex justify-between border-primary-light dark:border-primary-dark bg-secondary-light dark:bg-secondary-dark mt-6"
          color="light"
        >
          <P class="text-lg font-semibold text-tertiary-light dark:text-tertiary-dark"
            >{$_("summaryOfFindingsGtcDocument")}</P
          >
          <ChevronDownOutline
            class="w-8 h-8 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
          />
        </Button>
        {#if dropdownToggle[1]}
          <ul
            style:visibility={dropdownToggle[1] ? "visible" : "hidden"}
            class="py-2 flex flex-col items-center border border-secondary-light dark:border-secondary-dark rounded-lg text-sm text-primary-light dark:text-primary-dark"
            aria-labelledby="dropdownDefaultButton"
          >
            hello
          </ul>
        {/if}

        <!-- second table section -->
        <Button
          on:click={() => (dropdownToggle[2] = !dropdownToggle[2])}
          class="flex justify-between border-primary-light dark:border-primary-dark bg-secondary-light dark:bg-secondary-dark"
          color="light"
        >
          <P class="text-lg font-semibold text-tertiary-light dark:text-tertiary-dark"
            >{$_("summaryOfFindingsDocumentDealDone")}</P
          >
          <ChevronDownOutline
            class="w-8 h-8 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
          />
        </Button>
        {#if dropdownToggle[2]}
          <ul
            style:visibility={dropdownToggle[2] ? "visible" : "hidden"}
            class="py-2 flex flex-col items-center border border-secondary-light dark:border-secondary-dark rounded-lg text-sm text-primary-light dark:text-primary-dark"
            aria-labelledby="dropdownDefaultButton"
          >
            hello
          </ul>
        {/if}
      </div>
    </div>
  </div>
</div>

<Modal title={$_("documentViewer")} bind:open={documentViewModal} autoclose outsideclose>
  {#if selectedFile && (selectedFile.name.includes(".pdf") || selectedFile.name.includes(".PDF"))}
    <div class="flex justify-center">
      <PdfViewer
        showBorder={false}
        showButtons={["navigation", "zoom", "rotate", "autoflip", "pageInfo"]}
        scale={1}
        data={atob(selectedFile.base64)}
      />
    </div>
  {/if}
</Modal>
