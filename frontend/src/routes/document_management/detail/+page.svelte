<script>
  import {
    P,
    Span,
    Heading,
    Input,
    Textarea,
    Button,
    ButtonGroup,
    Dropdown,
    DropdownItem,
    Avatar,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    Modal,
    Select,
  } from "flowbite-svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import PDFViewer from "svelte-pdf";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import DownloadIcon from "$lib/icon/file-download-03.svg";
  import Divider from "../../../lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { add, format } from "date-fns";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import pdf from "$lib/template.pdf";

  let onEdit = $state(false);

  let tabDatas = $state({
    selected: "Documents",
    data: [
      {
        name: "Documents",
      },
      {
        name: "Information",
      },
    ],
  });

  let pdfReady = $state(true);

  let data = $state({
    document_title: "",
    document_upload: null,
    description: "",
    document_type: "",
    person_client: "",
    contract_start: "Jan 6, 2024",
    contract_end: "Jan 7, 2024",
  });

  let categories = $state([
    {
      value: "general",
      name: "General",
    },
    {
      value: "contract",
      name: "Contract",
    },
  ]);

  let personClient = $state([
    {
      value: "1",
      name: "Azizi",
    },
    {
      value: "2",
      name: "Habibi",
    },
  ]);

  onMount(() => {
    //
  });

</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex items-center pb-8">
        <div class="flex flex-col grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
              >Document Name <Button
                outline
                color="red"
                class="hover:!bg-transparent !border-hc-graylight-200 dark:!border-hc-graydark-600 focus-within:!ring-1 !p-2"
              >
                <CopyIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
              </Button></Heading
            >
          </div>
          <P
            class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
            >{$_("createOn")}{` uploadDate`}</P
          >
        </div>
      </div>

      <!-- Tab Bar -->
      <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

      <!-- Content -->
      <div class="">
        {#if tabDatas.selected === "Documents"}
          <div class="flex justify-between items-center pt-6 pb-4">
            <P class="font-semibold text-hc-graylight-700 dark:text-hc-graydark-300"
              >{$_("documentName")}</P
            >
            <div class="flex items-center gap-2">
              <Button type="submit" class="w-max px-3 py-[10px] gap-1 focus-within:!ring-1">
                <img alt="refresh" src={DownloadIcon} class="w-5 h-5" />
                {$_("downdloadDocument")}
              </Button>
            </div>
          </div>
          <div class="flex justify-center w-full py-10 bg-gray-50 rounded-lg dark:bg-gray-800">
            {#if pdfReady}
              <PDFViewer showBorder={false} showButtons={false} scale={1.2} url={pdf} />
            {/if}
          </div>
        {:else if tabDatas.selected === "Information"}
          <div class="flex flex-col gap-4 pt-6">
            <!-- general information -->
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {$_("generalInformation")}
                </P>
              </div>
              <div class="flex flex-col grow gap-6">
                <div class="flex gap-5">
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("documentName")}
                    </P>
                    <Input
                      readonly={!onEdit}
                      bind:value={data.document_title}
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                </div>
                <div class="grid grid-cols-1">
                  <div class="flex flex-col gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >{$_("description")}</P
                    >
                    <Textarea
                      readonly={!onEdit}
                      bind:value={data.document_description}
                      rows="4"
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                </div>
              </div>
            </div>
            <Divider />
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {$_("documentType")}*
                </P>
              </div>
              <div class="flex flex-col grow gap-6">
                <div class="flex flex-col gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("documentType")}*</P
                  >
                  <Select class="bg-primary-light dark:bg-secondary-dark" items={categories} bind:value={data.document_type} />
                </div>
              </div>
            </div>
            <Divider />
            <!-- general information -->
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {$_("contractInformation")}
                </P>
              </div>
              <div class="flex grow flex-col gap-6">
                <div class="flex flex-col grow gap-6">
                  <div class="flex flex-col gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >{$_("personsOrClient")}*</P
                    >
                    <Select class="bg-primary-light dark:bg-secondary-dark" items={personClient} bind:value={data.person_client} />
                  </div>
                </div>
                <div class="flex justify-center items-center gap-5">
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("contractStartPeriod")}
                    </P>
                    <Input
                      readonly={!onEdit}
                      bind:value={data.contract_start}
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                  <div class="mt-6">
                    <P>-</P>
                  </div>
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("contractEndPeriod")}
                    </P>
                    <Input
                      readonly={!onEdit}
                      bind:value={data.contract_end}
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
