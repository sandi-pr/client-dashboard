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
  import { check, onProgress, call } from "$lib";

  import ArrowLeftIcon from "$lib/icon/arrow-left.svelte";
  import ArrowRightIcon from "$lib/icon/arrow-right.svelte";
  import EyeIcon from "$lib/icon/eye.svelte";
  import Star04Icon from "$lib/icon/star-04.svelte";
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

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownActionOpen = dropdownActionOpen.map((item) => false);
  };

  onMount(() => {
    for (const data of collectionListDatas) {
      dropdownActionOpen.push(false);
    }
  });
</script>

<div
  on:focusout={handleDropdownFocusLoss}
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
      <TableHeadCell>{$_("action")}</TableHeadCell>
      <TableHeadCell>{$_("fileSize")}</TableHeadCell>
      <TableHeadCell>{$_("uploadDetail")}</TableHeadCell>
      <TableHeadCell></TableHeadCell>
    </TableHead>
    <TableBody tableBodyClass="divide-y">
      {#each collectionListDatas as collectionListData, index (collectionListData.id)}
        <TableBodyRow
          on:click={() => {
            checkSingle(index);
          }}
        >
          <TableBodyCell class="!p-4">
            {#if collectionListData.checked}
              <Checkbox
                checked={collectionListData.checked}
                on:click={() => {
                  // checkSingle(index);
                }}
              />
            {:else}
              <Checkbox
                checked={collectionListData.checked}
                on:click={() => {
                  // checkSingle(index);
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
            <Button
              disabled={!check("documentAi.documentLibrary.document:view")}
              on:click={async () => {
                let link = await $call.getPage("client-dashboard", "insights");
                goto(link);
              }}
              color="white"
              class="!p-2 focus-within:!ring-1 text-brand-tertiary-light dark:text-brand-tertiary-dark"
              size="lg"
            >
              {$_("askMyCsv")}&nbsp;
              <Star04Icon svgClass="w-5 h-5" />
            </Button>
          </TableBodyCell>
          <TableBodyCell>{collectionListData.filesize}</TableBodyCell>
          <TableBodyCell>
            <div class="flex items-center gap-1 text-disabled-light dark:text-disabled-dark">
              <Avatar size="xs" />
              {collectionListData.createdAt}
            </div>
          </TableBodyCell>
          <TableBodyCell>
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
                  class={`z-10 ${dropdownActionOpen[index] ? "absolute -top-4 right-20" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
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
          </TableBodyCell>
        </TableBodyRow>
      {/each}
    </TableBody>
  </Table>
  <div class="flex justify-between !border-0">
    {#if pagination.total > 0}
      <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
        Showing
        <span>{pagination.page === 1 ? 1 : (pagination.page - 1) * pagination.pageSize + 1}</span>
        -
        <!-- {$_("to")} -->
        <!-- <span
          >{pagination.page * pagination.pageSize >= pagination.total
            ? pagination.total
            : pagination.page * pagination.pageSize}</span
        > -->
        <span>{collectionListDatas.length}</span>
        of
        <span>{collectionListDatas.length}</span>
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
  {#if pagination.pageCount > 1}
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
        disabled={pagination.page === pagination.pageCount}
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
</div>
