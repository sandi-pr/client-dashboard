<script>
  import {
    P,
    Table,
    Button,
    Checkbox,
    TableBody,
    TableHead,
    Pagination,
    TableBodyRow,
    TableBodyCell,
    TableHeadCell,
    Avatar,
  } from "flowbite-svelte";

  import { _ } from "svelte-i18n";
  import { goto } from "$app/navigation";
  import { check, onProgress } from "$lib";

  import ArrowLeftIcon from "$lib/icon/arrow-left.svelte";
  import ArrowRightIcon from "$lib/icon/arrow-right.svelte";
  import EyeIcon from "$lib/icon/eye.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";
  import JPEGSolidIcon from "$lib/icon/JPEGSolid.svg";
  import JPGSolidIcon from "$lib/icon/JPGSolid.svg";
  import PNGSolidIcon from "$lib/icon/PNGSolid.svg";
  import CSVIcon from "$lib/icon/csv.svg";
  import Trash01Icon from "$lib/icon/trash-01.svelte";
  import { onMount } from "svelte";

  export let collectionListDatas;
  export let pagination;
  export let checkSingle;
  export let checkAll;
  export let checkAllState;
  export let next;
  export let previous;
  export let functDeleteCategories;

  let dropdownActionOpen = [];

  const pillStyle = (style) => {
    if (style === "brand") {
      return "py-[2px] px-2 border-[1px] rounded-xl text-center w-min bg-utility-brand-50-light dark:bg-utility-brand-50-dark border-utility-brand-200-light dark:border-utility-brand-200-dark text-utility-brand-700-light dark:text-utility-brand-700-dark";
    } else {
      return "py-[2px] px-2 border-[1px] rounded-xl text-center w-min bg-utility-brand-50-light dark:bg-utility-brand-50-dark border-utility-brand-200-light dark:border-utility-brand-200-dark text-utility-brand-700-light dark:text-utility-brand-700-dark";
    }
  };

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownActionOpen = dropdownActionOpen.map((item) => false);
  };

  // const filesizePostprocessing = (filesize) => {
  //   let size = filesize;
  //   let unit = "B";
  //   if (size >= 1024) {
  //     size = size / 1024;
  //     unit = "KB";
  //   }
  //   if (size >= 1024) {
  //     size = size / 1024;
  //     unit = "MB";
  //   }
  //   if (size >= 1024) {
  //     size = size / 1024;
  //     unit = "GB";
  //   }
  //   return `${size.toFixed(2)} ${unit}`;
  // };

  // console.log(collectionListDatas);

  onMount(() => {
    for (const data of collectionListDatas) {
      dropdownActionOpen.push(false);
    }
  });
</script>

<div
  class="mt-4 rounded-lg [&>div]:sm:rounded-t-lg bg-gray-50 dark:bg-gray-700 [&>div]:border [&>div]:border-primary-light [&>div]:dark:border-primary-dark"
>
  <Table hoverable={true}>
    <TableHead>
      <TableHeadCell class="!p-4">
        {#if checkAllState}
          <Checkbox checked on:click={checkAll} />
        {:else}
          <Checkbox on:click={checkAll} />
        {/if}
      </TableHeadCell>
      <TableHeadCell>{$_("fileName")}</TableHeadCell>
      <TableHeadCell>{$_("label")}</TableHeadCell>
      <TableHeadCell>{$_("type")}</TableHeadCell>
      <TableHeadCell>{$_("fileSize")}</TableHeadCell>
      <TableHeadCell>{$_("uploadDetail")}</TableHeadCell>
      <TableHeadCell></TableHeadCell>
      <TableHeadCell></TableHeadCell>
    </TableHead>
    <TableBody tableBodyClass="divide-y">
      {#each collectionListDatas as collectionListData, index (collectionListData.id)}
        <TableBodyRow>
          <TableBodyCell class="!p-4">
            {#if collectionListData.checked}
              <Checkbox
                checked={collectionListData.checked}
                on:click={() => {
                  checkSingle(index);
                }}
              />
            {:else}
              <Checkbox
                checked={collectionListData.checked}
                on:click={() => {
                  checkSingle(index);
                }}
              />
            {/if}
          </TableBodyCell>
          <TableBodyCell>
            <div class="flex gap-2 items-center min-w-64">
              {#if collectionListData.fileext === "jpg"}
                <img src={JPGSolidIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
              {:else if collectionListData.fileext === "jpeg"}
                <img src={JPEGSolidIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
              {:else if collectionListData.fileext === "png"}
                <img src={PNGSolidIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
              {:else if collectionListData.fileext === "pdf"}
                <img src={PDFSolidIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
              {:else if collectionListData.fileext === "doc"}
                <img src={DOCSolidIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
              {:else if collectionListData.fileext === "docx"}
                <img src={DOCXSolidIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
              {:else if collectionListData.fileext === "csv" || collectionListData.fileext === ".csv"}
                <img src={CSVIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
              {/if}
              <P class="line-clamp-2 fluid-sm font-semibold">{collectionListData.filename}</P>
            </div>
          </TableBodyCell>
          <TableBodyCell>
            <div class="flex flex-wrap gap-1 min-w-32 max-w-32">
              {#if collectionListData.label.length === 0}
                <div class="ps-3">-</div>
              {:else}
                {#each collectionListData.label as label, index (label)}
                  {#if index <= 2}
                    <div class={pillStyle()}>
                      {label}
                    </div>
                  {:else if index === 3}
                    <div class={pillStyle()}>
                      +{collectionListData.label.length - 3}
                    </div>
                  {/if}
                {/each}
              {/if}
            </div>
          </TableBodyCell>
          <TableBodyCell
            ><div class={pillStyle()}>
              {collectionListData.type}
            </div></TableBodyCell
          >
          <TableBodyCell>{collectionListData.filesize}</TableBodyCell>
          <TableBodyCell>
            <div class="flex items-center gap-1 text-disabled-light dark:text-disabled-dark">
              <Avatar size="xs" />
              {collectionListData.upload_by + " - " + collectionListData.createdAt}
            </div>
          </TableBodyCell>
          <TableBodyCell>
            <Button
              disabled={!check("documentAi.documentLibrary.document:view")}
              on:click={() => {
                goto(collectionListData.href);
              }}
              color="white"
              class="!p-2 focus-within:!ring-1 text-brand-tertiary-light dark:text-brand-tertiary-dark"
              size="lg"
            >
              <EyeIcon svgClass="w-5 h-5" />
              &nbsp;{$_("view")}
            </Button>
          </TableBodyCell>
          <TableBodyCell>
            <!-- <Button
              class="p-2 md:px-3 md:py-2 bg-transparent !text-black"
              color="white"
              on:click={() => {
                // addLibraryModal = true;
                // tmpCollectionName = "";
                // tmpDescription = "";
              }}
            >
              <DotsVerticalIcon svgClass="w-5 h-5 !text-quinary-light !dark:text-quinary-dark" />
            </Button> -->

            {#if check("documentAi.documentLibrary.document:delete")}
              <div class=" flex items-center relative">
                <Button
                  on:click={() => (dropdownActionOpen[index] = !dropdownActionOpen[index])}
                  class="focus-within:ring-0 !p-0"
                  color="none"
                  size="xs"
                  ><DotsVerticalIcon
                    class="text-fg-quinary-light dark:text-fg-quinary-dark"
                  /></Button
                >
                <div
                  class={`z-10 ${dropdownActionOpen[index] ? "absolute -top-4 right-12" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
                >
                  <ul class="py-2 grow text-sm">
                    <button
                      class="flex p-2 px-6 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                      on:click={() => {
                        dropdownActionOpen[index] = !dropdownActionOpen[index];
                        functDeleteCategories(collectionListData.id);
                      }}>{$_("delete")}</button
                    >
                  </ul>
                </div>
              </div>
            {/if}
            <!-- <Button
              on:click={() => {
                functDeleteCategories(collectionListData.id);
              }}
              disabled={!check("documentAi.documentLibrary.document:delete")}
              class="focus-within:ring-0 !p-0 text-brand-tertiary-light dark:text-brand-tertiary-dark"
              color="none"
              size="xs"
            >
              <Trash01Icon svgClass="w-5 h-5" />
              &nbsp;&nbsp;
              <P class="text-brand-tertiary-light dark:text-brand-tertiary-dark">{$_("delete")}</P>
            </Button> -->
          </TableBodyCell>
        </TableBodyRow>
      {/each}
    </TableBody>
  </Table>
  <div class="flex justify-between !border-0">
    {#if pagination.total > pagination.pageSize}
      <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
        Showing
        <span>{pagination.page === 1 ? 1 : (pagination.page - 1) * pagination.pageSize + 1}</span>
        -
        <!-- {$_("to")} -->
        <span
          >{pagination.page * pagination.pageSize >= pagination.total
            ? pagination.total
            : pagination.page * pagination.pageSize}</span
        >
        of
        <span>{pagination.total}</span>
        <!-- {$_("entries")} -->
      </div>

      <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
        page
        <span>{pagination.page}</span>
        of
        <span>{pagination.pageCount}</span>
      </div>
    {/if}
  </div>
  {#if pagination.total > pagination.pageSize}
    <div class="flex justify-end !border-0">
      <Button
        disabled={pagination.page === 1}
        on:click={previous}
        color="light"
        class="fluid-sm gap-1 md:gap-2 rounded-l-lg rounded-r-none px-3 py-2 md:px-4 md:py-2.5"
      >
        <ArrowLeftIcon
          svgClass="text-fg-secondary-light dark:text-fg-secondary-dark w-4 h-4 md:w-5 md:h-5"
        />
        previous
      </Button>
      <Button
        disabled={pagination.page * pagination.pageSize >= pagination.total}
        on:click={next}
        color="light"
        class="fluid-sm gap-1 md:gap-2 rounded-r-lg rounded-l-none px-3 py-2 md:px-4 md:py-2.5"
      >
        next
        <ArrowRightIcon
          svgClass="text-fg-secondary-light dark:text-fg-secondary-dark w-4 h-4 md:w-5 md:h-5"
        />
      </Button>
    </div>
  {/if}
  <!-- <div class="flex justify-between items-center !rounded-none !rounded-b-lg !border-t-0">
    <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
      {$_("showing")}
      <span>{metadata.page === 1 ? 1 : (metadata.page - 1) * metadata.pageSize + 1}</span>
      {$_("to")}
      <span
        >{metadata.page * metadata.pageSize >= metadata.total
          ? metadata.total
          : metadata.page * metadata.pageSize}</span
      >
      {$_("of")}
      <span>{metadata.total}</span>
      {$_("entries")}
    </div>

    <div class="m-4">
      <Pagination
        on:next={next}
        on:previous={previous}
        {pages}
        activeClass="bg-hc-graylight-100 dark:bg-hc-graydark-700 dark:border-[#374151] dark:text-hc-graydark-300"
      >
        <span slot="prev">{$_("previous")}</span>
        <span slot="next">{$_("next")}</span>
      </Pagination>
    </div>
  </div> -->
</div>
