<script>
  import { onMount } from "svelte";
  import { Button, Modal, P } from "flowbite-svelte";
  import { Typography } from "$lib";
  import PdfViewer from "svelte-pdf";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import FileDownloadIcon from "$lib/icon/file-download-03.svelte";
  import DownloadIcon from "$lib/icon/download-01.svelte";
  import { ChevronDownOutline, RefreshOutline } from "flowbite-svelte-icons";
  import templatePdf from "$lib/template.pdf";
  // import PdfViewer from 'svelte-pdf-scrollable';
  import DataTable from "$lib/components/DataTable.svelte";
  import DataTableAiReview from "$lib/components/DataTableAIReview.svelte";
  // import { Svroller, Svrollbar } from 'svrollbar';

  let projectName = $state("Project Name");
  let configs = $state({});
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });
  let dropdownToggle = $state([false, true, false]);
  let documentViewModal = $state(false);
  let selectedFile = $state(null);
  let pdfData = $state([null, null, null, null]);
  let base64Data = $state();
  let fileUrl = $state(null);
  let datas = $state([
    {
      section: "Section 1",
      result: "Result 1",
    }
  ]);
  let dummyDatas = $state([
    {
      id: 1,
      no: 1,
      parameter: "Parameter 1",
      referenceDocument: "12345",
      tradeConfirmation: "54321",
      aiChanges: "No Action Required",
    },
    {
      id: 2,
      no: 2,
      parameter: "Parameter 2",
      referenceDocument: "pertamina",
      tradeConfirmation: "pertamina",
      aiChanges: "No Action Required",
    },
    {
      id: 3,
      no: 3,
      parameter: "Parameter 3",
      referenceDocument: "01-Aug-2021",
      tradeConfirmation: "15-Aug-2021",
      aiChanges: "Adjust to match 01-Aug-2021",
      actions: {
        status: "Done",
        result: "Changed",
      },
    },
    {
      id: 4,
      no: 4,
      parameter: "Parameter 4",
      referenceDocument: "Net 30 Days",
      tradeConfirmation: "Net 60 Days",
      aiChanges: "Adjust to match Net 30 Days",
      actions: {
        status: "Done",
        result: "Rejected",
      },
    },
    {
      id: 5,
      no: 5,
      parameter: "Parameter 5",
      referenceDocument: "100.000 Barrels",
      tradeConfirmation: "200.000 Barrels",
      aiChanges: "Verify reason for change; update if acceptable",
      actions: {
        status: "Processing",
        reject: () => handleReject(5),
        change: () => handleChange(5),
      },
    },
  ]);

  async function getDocument(page) {
    datas = [
      {
        id: 1,
        no: 1,
        parameter: "Parameter 1",
        referenceDocument: "12345",
        tradeConfirmation: "54321",
        aiChanges: "No Action Required",
      },
      {
        id: 2,
        no: 2,
        parameter: "Parameter 2",
        referenceDocument: "pertamina",
        tradeConfirmation: "pertamina",
        aiChanges: "No Action Required",
      },
      {
        id: 3,
        no: 3,
        parameter: "Parameter 3",
        referenceDocument: "01-Aug-2021",
        tradeConfirmation: "15-Aug-2021",
        aiChanges: "Adjust to match 01-Aug-2021",
        actions: {
          status: "Done",
          result: "Changed",
        },
      },
      {
        id: 4,
        no: 4,
        parameter: "Parameter 4",
        referenceDocument: "Net 30 Days",
        tradeConfirmation: "Net 60 Days",
        aiChanges: "Adjust to match Net 30 Days",
        actions: {
          status: "Done",
          result: "Rejected",
        },
      },
      {
        id: 5,
        no: 5,
        parameter: "Parameter 5",
        referenceDocument: "100.000 Barrels",
        tradeConfirmation: "200.000 Barrels",
        aiChanges: "Verify reason for change; update if acceptable",
        actions: {
          status: "Processing",
          reject: () => handleReject(5),
          change: () => handleChange(5),
        },
      },
    ];
    configs = {
      pagination: pagination,
    };
  }

  async function handleReject(no) {
    console.log("reject");
  }

  async function handleChange(no) {
    console.log("change");
  }

  async function nextGTC() {
    getDocument(pagination.page + 1);
  }

  async function previousGTC() {
    getDocument(pagination.page - 1);
  }

  async function nextDealDone() {
    getDocumentsDealDone(pagination.page + 1);
  }

  async function previousDealDone() {
    getDocumentsDealDone(pagination.page - 1);
  }

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

  function base64toBlob(base64Data, fileType) {
    const byteCharacters = atob(base64Data);
    const byteArrays = [];

    for (let offset = 0; offset < byteCharacters.length; offset += 512) {
      const slice = byteCharacters.slice(offset, offset + 512);

      const byteNumbers = new Array(slice.length);
      for (let i = 0; i < slice.length; i++) {
        byteNumbers[i] = slice.charCodeAt(i);
      }

      const byteArray = new Uint8Array(byteNumbers);
      byteArrays.push(byteArray);
    }

    const blob = new Blob(byteArrays, { type: fileType });
    return blob;
  }

  async function createPdfUrl(file) {
    const response = await fetch(templatePdf);
    const blob = await response.blob();
    base64Data = await readFileAsBase64(blob);

    const pdfBlob = base64toBlob(base64Data, blob.type);
    fileUrl = URL.createObjectURL(pdfBlob);
  }

  onMount(async () => {
    // await getDocument(1);
    await createPdfUrl();
  });
</script>

<div class="flex flex-col px-4 pt-6 gap-6 pb-12">
  <!-- section header -->
  <div class="flex justify-between">
    <div class="flex items-center gap-3">
      <div class="flex flex-col gap-1.5">
        <P class="text-3xl text-primary-light dark:text-primary-dark font-semibold"
          >{$_("aiReview")}</P
        >
        <Typography preset="Text sm/Regular" text={projectName} color="tertiary" />
      </div>
    </div>
    <div class="flex items-center gap-3">
      <!-- <Button outline on:click={() => {}} color="light" class="h-max">
        <FileDownloadIcon svgClass="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark" />
        &nbsp;
        {$_("seeReferencedDoc")}
      </Button> -->
      <!-- <Button outline on:click={() => {}} color="light" class="h-max">
        <FileDownloadIcon svgClass="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark" />
        &nbsp;
        {$_("generateNewDocs")}
      </Button>
      <Button on:click={() => {}} class="h-max">
        <DownloadIcon svgClass="w-5 h-5 text-white" />
        &nbsp;
        {$_("download")}</Button
      > -->
    </div>
  </div>
  <Divider />
  <div class="flex gap-3">
    <div class="flex flex-col grow gap-6">
      <!-- <div onfocusout={handleDropdownFocusLoss} class="flex flex-col gap-1 w-full"> -->
      <div class="flex flex-col gap-6 w-full">
        <!-- document section -->
        <!-- <Button
          on:click={() => (dropdownToggle[0] = !dropdownToggle[0])}
          class="flex justify-between border-primary-light dark:border-primary-dark bg-secondary-light dark:bg-secondary-dark"
          color="light"
        >
          <div class="flex items-center gap-3">
            <P class="text-lg font-semibold text-tertiary-light dark:text-tertiary-dark"
              >{$_("generatedDocument")}</P
            >
            <Button outline on:click={() => (documentViewModal = true)} color="light" class="h-max">
              <FileDownloadIcon svgClass="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark" />
              &nbsp;
              {$_("seeReferencedDoc")}
            </Button>
          </div>
          <ChevronDownOutline
            class="w-8 h-8 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
          />
        </Button>
        {#if dropdownToggle[0]}
          <div class="flex justify-center mt-6">
            <div class="bg-tertiary-light dark:bg-tertiary-dark rounded-lg p-4">
              <PdfViewer
                showBorder={false}
                showButtons={["navigation", "zoom", "rotate", "autoflip", "pageInfo"]}
                scale={1}
                data={atob(base64Data)}
              />
            </div>
          </div>
        {/if}

        <Divider /> -->

        <!-- first table section -->
        <!-- <Button
          on:click={() => (dropdownToggle[1] = !dropdownToggle[1])}
          class="flex justify-between border-primary-light dark:border-primary-dark bg-secondary-light dark:bg-secondary-dark"
          color="light"
        >
          <P class="text-lg font-semibold text-tertiary-light dark:text-tertiary-dark"
            >{$_("summaryOfFindings")}</P
          >
          <ChevronDownOutline
            class="w-8 h-8 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
          />
        </Button> -->
        {#if dropdownToggle[1]}
          <!-- <DataTableAiReview
            bind:datas
            bind:configs
            bind:next={nextGTC}
            bind:previous={previousGTC}
          /> -->
          <DataTable bind:datas />
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