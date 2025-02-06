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

  function truncateText(text, maxWords, maxFirstWordLength = 20) {
    const words = text.split(" ");

    // Check if the first word is too long, or if the second word is long
    if (words[0] && words[0].length > maxFirstWordLength) {
      return words[0].slice(0, maxFirstWordLength) + "...";
    } else if (words[1] && words[1].length > maxFirstWordLength) {
      return words[0] + " " + words[1].slice(0, maxFirstWordLength) + "...";
    }

    // Limit to maxWords and append "..." if necessary
    if (words.length > maxWords) {
      return words.slice(0, maxWords).join(" ") + " ...";
    }

    return words.join(" ");
  }

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    // dropdownActionOpen = false;
  };

  onMount(() => {
    for (const data of collectionListDatas) {
      dropdownActionOpen.push(false);
    }
  });
</script>

<div
  class="mt-4 rounded-lg [&>div]:rounded-t-lg bg-gray-50 dark:bg-gray-700 [&>div]:border [&>div]:border-primary-light [&>div]:dark:border-primary-dark"
>
  <Table hoverable={true}>
    <TableBody tableBodyClass="divide-y">
      {#each collectionListDatas as collectionListData, index (collectionListData.id)}
        <TableBodyRow>
          <TableBodyCell class="!p-4">{$_("fileName")}</TableBodyCell>
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
              <P class="line-clamp-2 w-full truncate fluid-sm font-semibold">{truncateText(collectionListData.filename, 2)}</P>
            </div>
          </TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell>{$_("label")}</TableBodyCell>
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
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell>{$_("type")}</TableBodyCell>
          <TableBodyCell>
            <div class={pillStyle()}>
              {collectionListData.type}
            </div>
          </TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell>{$_("fileSize")}</TableBodyCell>
          <TableBodyCell>{collectionListData.filesize}</TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell>{$_("dateUpload")}</TableBodyCell>
          <TableBodyCell>{collectionListData.createdAt}</TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell>{$_("actions")}</TableBodyCell>
          <TableBodyCell class="flex justify-between mr-4">
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
            <!-- {#if check("documentAi.documentLibrary.document:delete")}
              <div on:focusout={handleDropdownFocusLoss} class=" flex items-center relative">
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
                  class={`z-10 ${dropdownActionOpen[index] ? "absolute top-10 right-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
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
            {/if} -->
            <Button
              on:click={() => {
                functDeleteCategories(collectionListData.id);
              }}
              disabled={!check("documentAi.documentLibrary.document:delete") || $onProgress}
              class="focus-within:ring-0 !p-0 text-brand-tertiary-light dark:text-brand-tertiary-dark"
              color="none"
              size="xs"
            >
              <Trash01Icon svgClass="w-5 h-5" />
              &nbsp;&nbsp;
              <P class="text-brand-tertiary-light dark:text-brand-tertiary-dark">{$_("delete")}</P>
            </Button>
            <div class="flex items-center ml-1">
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
              <P class="text-brand-tertiary-light hidden dark:text-brand-tertiary-dark">{$_("select")}</P>
            </div>
          </TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell></TableBodyCell>
          <TableBodyCell></TableBodyCell>
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
