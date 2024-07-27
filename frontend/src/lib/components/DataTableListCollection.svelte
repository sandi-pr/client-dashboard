<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Dean Ramhan <dean@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import {
    P,
    Input,
    Avatar,
    Pagination,
    Button,
    Heading,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
  } from "flowbite-svelte";
  import * as Icon from "flowbite-svelte-icons";
  import { pagesBuilder } from "../utils/pagesBuilder";
  import { _ } from "svelte-i18n";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import JpgIcon from "$lib/icon/jpg-icon.svelte";
  import PngIcon from "$lib/icon/png-icon.svelte";
  import PdfIcon from "$lib/icon/pdf-icon.svelte";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";
  import JPEGSolidIcon from "$lib/icon/JPEGSolid.svg";
  import JPGSolidIcon from "$lib/icon/JPGSolid.svg";
  import PNGSolidIcon from "$lib/icon/PNGSolid.svg";
  // import { toCamelCase } from "../utils/toCamelCase.js";

  export let datas;
  export let datasChecked;
  export let configs;
  export let next = () => {};
  export let previous = () => {};
  export let handleCheckAll = () => {};
  export let handleCheck = () => {};
  // export let divider = false;

  let allChecked = false;

  let pages;
  if (configs?.pagination && Object.keys(configs.pagination).length > 0) {
    pages = pagesBuilder(configs.pagination);
  }

  const splitCamelCase = (value) => {
    return value.replace(/([a-z0-9])([A-Z])/g, "$1 $2").replace(/_/g, " ");
  };

  const classPillType = (type) => {
    if (type === "success") {
      return "bg-utility-success-50-light dark:bg-utility-success-50-dark border-utility-success-200-light dark:border-utility-success-200-dark text-utility-success-700-light dark:text-utility-success-700-dark";
    } else if (type === "error") {
      return "bg-utility-error-50-light dark:bg-utility-error-50-dark border-utility-error-200-light dark:border-utility-error-200-dark text-utility-error-700-light dark:text-utility-error-700-dark";
    } else if (type === "warning") {
      return "bg-utility-warning-50-light dark:bg-utility-warning-50-dark border-utility-warning-200-light dark:border-utility-warning-200-dark text-utility-warning-700-light dark:text-utility-warning-700-dark";
    } else if (type === "info") {
      return "bg-utility-gray-50-light dark:bg-utility-gray-50-dark border-utility-gray-200-light dark:border-utility-gray-200-dark text-utility-gray-700-light dark:text-utility-gray-700-dark";
    } else if (type === "brand") {
      return "bg-utility-brand-50-light dark:bg-utility-brand-50-dark border-utility-brand-200-light dark:border-utility-brand-200-dark text-utility-brand-700-light dark:text-utility-brand-700-dark";
    }
  };

  const classActionType = (type) => {
    if (type === "success") {
      return "bg-hc-success-50 border-hc-success-200 text-hc-success-700";
    } else if (type === "warning") {
      return "bg-hc-warning-50 border-hc-warning-200 text-hc-warning-700";
    } else if (type === "info") {
      return "bg-hc-graylight-50 dark:bg-hc-graydark-900 hover:bg-hc-graylight-100 hover:dark:bg-hc-graydark-800 focus-within:ring-0 hover:no-underline";
    } else if (type === "brand") {
      return "text-white bg-hc-brand-600 dark:bg-hc-brand-600 hover:bg-hc-brand-700 dark:hover:bg-hc-brand-700 focus-within:ring-0 hover:no-underline";
    } else if (type === "brand-inverse") {
      return "text-hc-brand-600 dark:text-hc-brand-400 bg-white dark:bg-hc-graydark-900 hover:bg-hc-graylight-100 hover:dark:bg-hc-graydark-800 focus-within:ring-0 hover:no-underline";
    } else if (type === "error") {
      return "text-hc-error-600 dark:text-hc-error-400 bg-white dark:bg-hc-graydark-900 hover:bg-hc-graylight-100 hover:dark:bg-hc-graydark-800 focus-within:ring-0 hover:no-underline";
    }
  };

  const withSeperator = (configs) => {
    if (configs?.seperator && configs.seperator === true) {
      return "pb-4 border-b border-hc-graylight-200 dark:border-hc-graydark-800";
    }
  };

  let searchTerm = "";

  // $: filteredItems = datas.filter((data) => {
  //   if (searchTerm === "") return data;
  //   if (data?.name) {
  //     return data.name.toLowerCase().indexOf(searchTerm.toLowerCase()) !== -1;
  //   }
  // });

  $: filteredItems = datas.filter((data) => {
    if (searchTerm === "") return true;
    if (data) {
      return Object.values(data).some((value) => {
        if (typeof value === "string") {
          return value.toLowerCase().includes(searchTerm.toLowerCase());
        }
        return false;
      });
    }
  });

  const handleSearch = (event) => {
    searchTerm = event.target.value;
  };

  const getCheckedStatus = (id) => {
    const item = datasChecked.find((data) => data.id === id);
    if (item) {
      allChecked = datasChecked.every((data) => data.checked);
      return item.checked;
    }
    return false;
  };

  const getObjectDataTypeClass = (type, color) => {
    if (type === "pill") {
      return classPillType(color);
    }
  };

  // onMount(() => {
  // Object.keys(datas[0]).forEach((key) => {
  //   console.log(key);
  // });
  // });
</script>

<div>
  <div
    class={`flex items-center justify-between p-4 mt-8 border-t-2 border-x-2 border-primary-light dark:border-primary-dark rounded-t-lg ${withSeperator(configs)}`}
  >
    {#if configs?.headers}
      <Heading tag="h6" class="font-jakarta w-fit text-hc-graylight-900 dark:text-hc-graydark-50"
        >{configs.headers}</Heading
      >
    {/if}
    <div class="flex items-center gap-3">
      <Input
        id="search"
        type="search"
        class="bg-white dark:bg-hc-graydark-950 py-3"
        placeholder={$_("search")}
        on:change={handleSearch}
      >
        <Icon.SearchOutline slot="left" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
      </Input>
      {#if configs?.actions && Object.keys(configs.actions).length > 0}
        {#each configs.actions as action (action.id)}
          <Button
            onclick={action.onclick}
            color="light"
            class={`${classActionType(action.type)} min-w-max px-2 py-3`}
          >
            {#if action?.icon}
              {#if action.icon === "plus"}
                <Icon.PlusOutline class="w-5 h-5" />&nbsp;
              {/if}
              {#if action.icon === "trash"}
                <Icon.TrashBinOutline class="w-5 h-5" />&nbsp;
              {/if}
            {/if}
            {action.name}
          </Button>
        {/each}
      {/if}
    </div>
  </div>
  <div
    class="[&>div]:sm:rounded-none [&>div]:sm:rounded-b-lg [&>div]:border-x-2 [&>div]:border-b-2 [&>div]:border-primary-light [&>div]:dark:border-primary-dark"
  >
    <Table class="font-jakarta dark:border-white">
      <TableHead class="bg-hc-graylight-50 dark:bg-hc-graydark-900">
        <TableHeadCell class="!p-4">
          <input
            type="checkbox"
            id="checkAll"
            name="checkAll"
            checked={allChecked}
            on:click={() => {
              handleCheckAll(allChecked);
            }}
            class="rounded h-5 w-5 dark:bg-primary-dark text-white focus:outline-none checked:bg-hc-brand-500 checked:dark:bg-hc-brand-500 focus:ring-2 focus:ring-hc-brand-500"
          />
        </TableHeadCell>
        {#if datas.length > 0}
          {#each Object.keys(datas[0]) as key}
            {#if key !== "actions" && key !== "id"}
              <TableHeadCell>{splitCamelCase($_(toCamelCase(key)))}</TableHeadCell>
            {/if}
          {/each}
          {#if datas[0]?.actions && Object.keys(datas[0].actions).length > 0}
            <TableHeadCell></TableHeadCell>
          {/if}
        {/if}
      </TableHead>
      <TableBody>
        {#each configs?.search ? filteredItems : datas as data, index (data.id)}
          <TableBodyRow
            class="dark:bg-hc-graydark-950 border-hc-graylight-200 dark:border-hc-graydark-800"
          >
            <TableBodyCell class="!p-4">
              <input
                type="checkbox"
                id={data.id}
                name={data.id}
                checked={getCheckedStatus(data.id)}
                on:click={() => {
                  handleCheck(data.id);
                }}
                class="rounded h-5 w-5 dark:bg-primary-dark text-white focus:outline-none checked:bg-hc-brand-500 checked:dark:bg-hc-brand-500 focus:ring-2 focus:ring-hc-brand-500"
              />
            </TableBodyCell>
            {#each Object.keys(datas[0]) as key}
              {#if key !== "actions" && key !== "id"}
                {#if typeof data[key] === "string" || typeof data[key] === "number"}
                  <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400">
                    {data[key]}
                  </TableBodyCell>
                {:else if typeof data[key] === "object"}
                  {#if data[key]?.type !== "custom"}
                    <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400">
                      <div
                        class={`py-[2px] px-2 border-[1px] rounded-xl text-center w-min ${getObjectDataTypeClass(data[key]?.type, data[key]?.color)}`}
                      >
                        {#if data[key]?.name}
                          {data[key].name}
                        {/if}
                      </div>
                    </TableBodyCell>
                  {:else}
                    <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400">
                      <div class="flex gap-2">
                        {#if data[key]?.name?.includes(".jpg") || data[key]?.name?.includes(".JPG")}
                          <img src={JPGSolidIcon} alt="" />
                        {:else if data[key]?.name?.includes(".jpeg") || data[key]?.name?.includes(".JPEG")}
                          <img src={JPEGSolidIcon} alt="" />
                        {:else if data[key]?.name?.includes(".png") || data[key]?.name?.includes(".PNG")}
                          <img src={PNGSolidIcon} alt="" />
                        {:else if data[key]?.name?.includes(".pdf") || data[key]?.name?.includes(".PDF")}
                          <img src={PDFSolidIcon} alt="" />
                        {:else if data[key]?.name?.includes(".doc") || data[key]?.name?.includes(".DOC")}
                          <img src={DOCSolidIcon} alt="" />
                        {:else if data[key]?.name?.includes(".docx") || data[key]?.name?.includes(".DOCX")}
                          <img src={DOCXSolidIcon} alt="" />
                        {/if}
                        <!-- {#if data[key]?.type?.name?.includes(".jpg") || data[key]?.type?.name?.includes(".JPG") || data[key]?.type?.name?.includes(".jpeg") || data[key]?.type?.name?.includes(".JPEG")}
                          <JpgIcon />
                        {:else if data[key]?.type?.name?.includes(".png") || data[key]?.type?.name?.includes(".PNG")}
                          <PngIcon />
                        {:else if data[key]?.type?.name?.includes(".pdf") || data[key]?.type?.name?.includes(".PDF")}
                          <PdfIcon />
                        {/if} -->
                        <!-- <JpgIcon /> -->
                        <div class="flex flex-col">
                          <P>{data[key].name}</P>
                          <P size="sm" class="text-hc-graylight-500">{data[key].size}</P>
                        </div>
                      </div>
                    </TableBodyCell>
                  {/if}
                {/if}
              {/if}
            {/each}
            {#if datas[0]?.actions && Object.keys(datas[0].actions).length > 0}
              {#each Object.keys(datas[0]) as key}
                {#if key === "actions"}
                  <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400">
                    <div class="flex justify-end gap-8">
                      {#each datas[index].actions.objects as action (action.name)}
                        {#if action.href}
                          <a
                            class="hover:no-underline text-hc-brand-600 font-semibold"
                            href={action.href}>{$_(toCamelCase(action.name))}</a
                          >
                        {:else}
                          <button
                            class="hover:no-underline text-hc-brand-600 font-semibold"
                            on:click={() => action.function(datas[index].id)}
                          >
                            <!-- {$_(toCamelCase(action.name))} -->
                            <Icon.DotsVerticalOutline /></button
                          >
                        {/if}
                      {/each}
                    </div>
                  </TableBodyCell>
                {/if}
              {/each}
            {/if}
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>
</div>

{#if configs?.pagination && Object.keys(configs.pagination).length > 0}
  <div class="flex justify-between">
    <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
      {"showing"}
      <span
        >{configs.pagination.page === 1
          ? 1
          : (configs.pagination.page - 1) * configs.pagination.pageSize + 1}</span
      >
      {"to"}
      <span
        >{configs.pagination.page * configs.pagination.pageSize >= configs.pagination.total
          ? configs.pagination.total
          : configs.pagination.page * configs.pagination.pageSize}</span
      >
      {"of"}
      <span>{configs.pagination.total}</span>
      {"entries"}
    </div>

    <div class="m-4">
      <Pagination
        bind:pages
        on:next={next}
        on:previous={previous}
        activeClass="bg-hc-graylight-100 dark:bg-hc-graydark-700 dark:border-[#374151] dark:text-hc-graydark-300"
      >
        <span slot="prev">{"previous"}</span>
        <span slot="next">{"next"}</span>
      </Pagination>
    </div>
  </div>
{/if}
