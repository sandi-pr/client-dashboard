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

  import EyeIcon from "$lib/icon/eye.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";
  import JPEGSolidIcon from "$lib/icon/JPEGSolid.svg";
  import JPGSolidIcon from "$lib/icon/JPGSolid.svg";
  import PNGSolidIcon from "$lib/icon/PNGSolid.svg";
  import Trash01Icon from "$lib/icon/trash-01.svelte";
  import { onMount } from "svelte";

  export let collectionListDatas;
  export let metadata;
  export let checkSingle;
  export let checkAll;
  export let checkAllState;
  export let next;
  export let previous;
  export let functDeleteCategories;
  export let datasChecked;
  export let handleCheckAll;
  export let handleCheck;
  let allChecked = false;

  let dropdownActionOpen = [];

  const pillStyle = (style) => {
    if (style === "brand") {
      return "py-[2px] px-2 border-[1px] rounded-xl text-center w-min bg-utility-brand-50-light dark:bg-utility-brand-50-dark border-utility-brand-200-light dark:border-utility-brand-200-dark text-utility-brand-700-light dark:text-utility-brand-700-dark";
    } else {
      return "py-[2px] px-2 border-[1px] rounded-xl text-center w-min bg-utility-brand-50-light dark:bg-utility-brand-50-dark border-utility-brand-200-light dark:border-utility-brand-200-dark text-utility-brand-700-light dark:text-utility-brand-700-dark";
    }
  };

  const pagesGenerator = () => {
    let length = metadata.total / metadata.pageSize;
    length = Math.ceil(length);
    let pages = [];
    if (length > 6) {
      for (let i = 0; i < 3; i++) {
        pages.push({
          name: `${i + 1}`,
        });
      }
      pages.push({
        name: "...",
      });
      // push the last 3 pages
      for (let i = length - 3; i < length; i++) {
        pages.push({
          name: `${i + 1}`,
        });
      }
    } else {
      for (let i = 0; i < length; i++) {
        pages.push({
          name: `${i + 1}`,
        });
      }
    }
    return pages;
  };

  let pages = pagesGenerator();

  const getCheckedStatus = (id) => {
    const item = datasChecked.find((data) => data.id === id);
    if (item) {
      allChecked = datasChecked.every((data) => data.checked);
      return item.checked;
    }
    return false;
  };

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
      <TableHeadCell>{$_("collectionName")}</TableHeadCell>
    </TableHead>
    <TableBody tableBodyClass="divide-y">
      {#each collectionListDatas as collectionListData, index (collectionListData.id)}
        <TableBodyRow on:click={() => {checkSingle(index)}}>
          <TableBodyCell class="!p-4">
            <input
              type="checkbox"
              id={collectionListDatas.id}
              name={collectionListDatas.id}
              checked={getCheckedStatus(collectionListDatas.id)}
              onclick={() => {
                handleCheck(collectionListDatas.id);
              }}
              class="rounded h-5 w-5 dark:bg-primary-dark text-white focus:outline-none checked:bg-hc-brand-500 checked:dark:bg-hc-brand-500 focus:ring-2 focus:ring-hc-brand-500"
            />
          </TableBodyCell>
          <TableBodyCell>{collectionListData.collection_name}</TableBodyCell>
        </TableBodyRow>
      {/each}
    </TableBody>
  </Table>
  <div class="flex justify-between items-center !rounded-none !rounded-b-lg !border-t-0">
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
        <!--  -->
        <span slot="next">{$_("next")}</span>
      </Pagination>
    </div>
  </div>
</div>
