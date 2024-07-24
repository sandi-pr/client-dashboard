<script>
  import { onMount } from "svelte";
  import { Button, Modal, P } from "flowbite-svelte";
  import { Typography } from "$lib";
  import PdfViewer from "svelte-pdf";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import FileDownloadIcon from "$lib/icon/file-download-03.svg";
  import { ChevronDownOutline, RefreshOutline } from "flowbite-svelte-icons";
  import templatePdf from "$lib/template.pdf";
  import DataTableAiReview from "$lib/components/DataTableAIReview.svelte";
  // import { Svroller, Svrollbar } from 'svrollbar';

  let configs = $state({});
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });
  let dropdownToggle = $state([false, false, false]);
  let pdfData = $state([null, null, null, null]);
  let datasGTC = $state([
    {
      id: 1,
      no: 1,
      parameter: "Parameter 1",
      referenceDocumentGTC: "12345",
      sellerGTC: "54321",
      recommendation: "No Action Required",
    },
    {
      id: 2,
      no: 2,
      parameter: "Parameter 2",
      referenceDocumentGTC: "pertamina",
      sellerGTC: "pertamina",
      recommendation: "No Action Required",
    },
    {
      id: 3,
      no: 3,
      parameter: "Parameter 3",
      referenceDocumentGTC: "01-Aug-2021",
      sellerGTC: "15-Aug-2021",
      recommendation: "Adjust to match 01-Aug-2021",
      actions: {
        status: "Done",
        result: "Changed",
      },
    },
    {
      id: 4,
      no: 4,
      parameter: "Parameter 4",
      referenceDocumentGTC: "Net 30 Days",
      sellerGTC: "Net 60 Days",
      recommendation: "Adjust to match Net 30 Days",
      actions: {
        status: "Done",
        result: "Rejected",
      },
    },
    {
      id: 5,
      no: 5,
      parameter: "Parameter 5",
      referenceDocumentGTC: "100.000 Barrels",
      sellerGTC: "200.000 Barrels",
      recommendation: "Verify reason for change; update if acceptable",
      actions: {
        status: "Processing",
        reject: () => handleReject(5),
        change: () => handleChange(5),
      },
    },
  ]);
  let datasDealDone = $state([
    {
      id: 1,
      no: 1,
      parameter: "Parameter 1",
      referenceDocumentGTC: "12345",
      sellerGTC: "54321",
      recommendation: "No Action Required",
    },
    {
      id: 2,
      no: 2,
      parameter: "Parameter 2",
      referenceDocumentGTC: "pertamina",
      sellerGTC: "pertamina",
      recommendation: "No Action Required",
    },
    {
      id: 3,
      no: 3,
      parameter: "Parameter 3",
      referenceDocumentGTC: "01-Aug-2021",
      sellerGTC: "15-Aug-2021",
      recommendation: "Adjust to match 01-Aug-2021",
      actions: {
        status: "Done",
        result: "Changed",
      },
    },
    {
      id: 4,
      no: 4,
      parameter: "Parameter 4",
      referenceDocumentGTC: "Net 30 Days",
      sellerGTC: "Net 60 Days",
      recommendation: "Adjust to match Net 30 Days",
      actions: {
        status: "Done",
        result: "Rejected",
      },
    },
    {
      id: 5,
      no: 5,
      parameter: "Parameter 5",
      referenceDocumentGTC: "100.000 Barrels",
      sellerGTC: "200.000 Barrels",
      recommendation: "Verify reason for change; update if acceptable",
      actions: {
        status: "Processing",
        reject: () => handleReject(5),
        change: () => handleChange(5),
      },
    },
  ]);

  async function getDocumentsGTC(page) {
    datasGTC = [
      {
        id: 1,
        no: 1,
        parameter: "Parameter 1",
        referenceDocumentGTC: "12345",
        sellerGTC: "54321",
        recommendation: "No Action Required",
      },
      {
        id: 2,
        no: 2,
        parameter: "Parameter 2",
        referenceDocumentGTC: "pertamina",
        sellerGTC: "pertamina",
        recommendation: "No Action Required",
      },
      {
        id: 3,
        no: 3,
        parameter: "Parameter 3",
        referenceDocumentGTC: "01-Aug-2021",
        sellerGTC: "15-Aug-2021",
        recommendation: "Adjust to match 01-Aug-2021",
        actions: {
          status: "Done",
          result: "Changed",
        },
      },
      {
        id: 4,
        no: 4,
        parameter: "Parameter 4",
        referenceDocumentGTC: "Net 30 Days",
        sellerGTC: "Net 60 Days",
        recommendation: "Adjust to match Net 30 Days",
        actions: {
          status: "Done",
          result: "Rejected",
        },
      },
      {
        id: 5,
        no: 5,
        parameter: "Parameter 5",
        referenceDocumentGTC: "100.000 Barrels",
        sellerGTC: "200.000 Barrels",
        recommendation: "Verify reason for change; update if acceptable",
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

  async function getDocumentsDealDone(page) {
    datasDealDone = [
      {
        id: 1,
        no: 1,
        parameter: "Parameter 1",
        referenceDocumentGTC: "12345",
        sellerGTC: "54321",
        recommendation: "No Action Required",
      },
      {
        id: 2,
        no: 2,
        parameter: "Parameter 2",
        referenceDocumentGTC: "pertamina",
        sellerGTC: "pertamina",
        recommendation: "No Action Required",
      },
      {
        id: 3,
        no: 3,
        parameter: "Parameter 3",
        referenceDocumentGTC: "01-Aug-2021",
        sellerGTC: "15-Aug-2021",
        recommendation: "Adjust to match 01-Aug-2021",
        actions: {
          status: "Done",
          result: "Changed",
        },
      },
      {
        id: 4,
        no: 4,
        parameter: "Parameter 4",
        referenceDocumentGTC: "Net 30 Days",
        sellerGTC: "Net 60 Days",
        recommendation: "Adjust to match Net 30 Days",
        actions: {
          status: "Done",
          result: "Rejected",
        },
      },
      {
        id: 5,
        no: 5,
        parameter: "Parameter 5",
        referenceDocumentGTC: "100.000 Barrels",
        sellerGTC: "200.000 Barrels",
        recommendation: "Verify reason for change; update if acceptable",
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
    getDocumentsGTC(pagination.page + 1);
  }

  async function previousGTC() {
    getDocumentsGTC(pagination.page - 1);
  }

  async function nextDealDone() {
    getDocumentsDealDone(pagination.page + 1);
  }

  async function previousDealDone() {
    getDocumentsDealDone(pagination.page - 1);
  }

  onMount(async () => {
    await getDocumentsGTC(1);
    await getDocumentsDealDone(1);

    const response = await fetch(templatePdf);
    const blob = await response.blob();
    const reader = new FileReader();

    reader.onloadend = () => {
      // Split the data URL at the comma and take the second part
      pdfData[0] = reader.result.split(",")[1];
      pdfData[1] = reader.result.split(",")[1];
      pdfData[2] = reader.result.split(",")[1];
      pdfData[3] = reader.result.split(",")[1];
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
          <div class="grid grid-cols-2 gap-4 mb-8">
            <div class="col-span-1">
              <P
                class="font-bold py-4 text-brand-tertiary-alt-light dark:text-brand-tertiary-alt-dark"
                >{$_("referenceDocumentGTC")}</P
              >
              <!-- <embed
                src={templatePdf + "#toolbar=0"}
                type="application/pdf"
                width="100%"
                height="100%"
              /> -->
              {#if pdfData[0]}
                <div class="[&>*]:!m-0 [&>*]:!p-0 [&>*>div>*>*>canvas]:mx-auto">
                  <PdfViewer
                    showBorder={false}
                    showButtons={["navigation", "zoom", "rotate", "autoflip", "pageInfo"]}
                    scale={1}
                    data={atob(pdfData[0])}
                  />
                </div>
              {:else}
                <P>No Document</P>
              {/if}
            </div>
            <div class="col-span-1">
              <P class="font-bold py-4 text-primary-light dark:text-primary-dark"
                >{$_("sellerGTC")}</P
              >
              <!-- <embed
                src={templatePdf + "#toolbar=0"}
                type="application/pdf"
                width="100%"
                height="100%"
              /> -->
              {#if pdfData[1]}
                <div class="[&>*]:!m-0 [&>*]:!p-0 [&>*>div>*>*>canvas]:mx-auto">
                  <PdfViewer
                    showBorder={false}
                    showButtons={["navigation", "zoom", "rotate", "autoflip", "pageInfo"]}
                    scale={1}
                    data={atob(pdfData[1])}
                  />
                </div>
              {:else}
                <P>No Document</P>
              {/if}
            </div>
            <div class="col-span-1">
              <P
                class="font-bold py-4 text-brand-tertiary-alt-light dark:text-brand-tertiary-alt-dark"
                >{$_("referenceDocumentDealDone")}</P
              >
              <!-- <embed
                src={templatePdf + "#toolbar=0"}
                type="application/pdf"
                width="100%"
                height="100%"
              /> -->
              {#if pdfData[2]}
                <div class="[&>*]:!m-0 [&>*]:!p-0 [&>*>div>*>*>canvas]:mx-auto">
                  <PdfViewer
                    showBorder={false}
                    showButtons={["navigation", "zoom", "rotate", "autoflip", "pageInfo"]}
                    scale={1}
                    data={atob(pdfData[2])}
                  />
                </div>
              {:else}
                <P>No Document</P>
              {/if}
            </div>
            <div class="col-span-1">
              <P class="font-bold py-4 text-primary-light dark:text-primary-dark"
                >{$_("documentDealDone")}</P
              >
              <!-- <embed
                src={templatePdf + "#toolbar=0"}
                type="application/pdf"
                width="100%"
                height="100%"
              /> -->
              {#if pdfData[3]}
                <div class="[&>*]:!m-0 [&>*]:!p-0 [&>*>div>*>*>canvas]:mx-auto">
                  <PdfViewer
                    showBorder={false}
                    showButtons={["navigation", "zoom", "rotate", "autoflip", "pageInfo"]}
                    scale={1}
                    data={atob(pdfData[3])}
                  />
                </div>
              {:else}
                <P>No Document</P>
              {/if}
            </div>
          </div>
          <div class="mb-8">
            <Divider />
          </div>
        {/if}

        <!-- first table section -->
        <Button
          on:click={() => (dropdownToggle[1] = !dropdownToggle[1])}
          class="flex justify-between border-primary-light dark:border-primary-dark bg-secondary-light dark:bg-secondary-dark"
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
          <DataTableAiReview
            bind:datas={datasGTC}
            bind:configs
            bind:next={nextGTC}
            bind:previous={previousGTC}
          />
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
          <DataTableAiReview
            bind:datas={datasDealDone}
            bind:configs
            bind:next={nextDealDone}
            bind:previous={previousDealDone}
          />
        {/if}
      </div>
    </div>
  </div>
</div>
