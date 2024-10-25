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
  } from "flowbite-svelte";

  import { _ } from "svelte-i18n";
  import { goto } from "$app/navigation";

  import EyeIcon from "$lib/icon/eye.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";
  import JPEGSolidIcon from "$lib/icon/JPEGSolid.svg";
  import JPGSolidIcon from "$lib/icon/JPGSolid.svg";
  import PNGSolidIcon from "$lib/icon/PNGSolid.svg";

  export let collectionListDatas;
  export let metadata;
  export let checkSingle;
  export let checkAll;
  export let checkAllState;
  export let next;
  export let previous;
  export let functDeleteCategories;

  let dropdownActionOpen = false;

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

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownActionOpen = false;
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
      <TableHeadCell>File name</TableHeadCell>
      <TableHeadCell>Label</TableHeadCell>
      <TableHeadCell>Type</TableHeadCell>
      <TableHeadCell>File size</TableHeadCell>
      <TableHeadCell>Date Upload</TableHeadCell>
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
            <div class="flex gap-2 items-center">
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
              {/if}
              <P class="line-clamp-2 fluid-sm font-semibold">{collectionListData.filename}</P>
            </div>
          </TableBodyCell>
          <TableBodyCell>
            <div class="flex gap-1">
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
            </div>
          </TableBodyCell>
          <TableBodyCell
            ><div class={pillStyle()}>
              {collectionListData.type}
            </div></TableBodyCell
          >
          <TableBodyCell>{collectionListData.filesize}</TableBodyCell>
          <TableBodyCell>{collectionListData.createdAt}</TableBodyCell>
          <TableBodyCell>
            <Button
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
            <div on:focusout={handleDropdownFocusLoss} class=" flex items-center relative">
              <Button
                on:click={() => (dropdownActionOpen = !dropdownActionOpen)}
                class="focus-within:ring-0 !p-0"
                color="none"
                size="xs"
                ><DotsVerticalIcon
                  class="text-fg-quinary-light dark:text-fg-quinary-dark"
                /></Button
              >
              <div
                class={`z-10 ${dropdownActionOpen ? "absolute top-10 right-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
              >
                <ul class="py-2 grow text-sm">
                  <!-- <button
                    class="flex p-2 px-6 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                    onclick={() => {isEdit = true}}>{$_("edit")}</button
                  > -->
                  <button
                    class="flex p-2 px-6 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                    on:click={() => {
                      functDeleteCategories(),
                        (dropdownActionOpen = !dropdownActionOpen);
                    }}>{$_("delete")}</button
                  >
                </ul>
              </div>
            </div>
          </TableBodyCell>
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
