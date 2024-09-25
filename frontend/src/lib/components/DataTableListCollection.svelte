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
  import Dropdown from "$lib/ui/buttons/Dropdown.svelte";
  import * as Icon from "flowbite-svelte-icons";
  import { pagesBuilder } from "../utils/pagesBuilder";
  import { _, locale } from "svelte-i18n";
  import { formatDistanceToNow } from "date-fns";
  import { id as localeId, enUS as localeEn } from "date-fns/locale";
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
  import ArrowLeftIcon from "$lib/icon/arrow-left.svelte";
  import ArrowRightIcon from "$lib/icon/arrow-right.svelte";
  import FIlePlusIcon from "$lib/icon/file-plus-02.svelte";
  import TrashIcon from "$lib/icon/trash-01.svelte";
  import ChevronDownIcon from "$lib/icon/chevron-down.svelte";
  import ChevronUpIcon from "$lib/icon/chevron-up.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  // import { toCamelCase } from "../utils/toCamelCase.js";

  export let datas;
  export let datasChecked;
  export let gettingDatas;
  export let configs;
  export let next;
  export let previous;
  export let handleCheckAll;
  export let handleCheck;
  export let functDeleteCategories;

  let dropdownToggle = true;

  let dropdownActionOpen = false;
  let selectedAction = "";
  let action = [
    {
      value: "delete",
      name: "Delete",
    },
  ];
  // export let divider = false;

  let allChecked = false;

  let pages;
  if (configs?.pagination && Object.keys(configs.pagination).length > 0) {
    pages = pagesBuilder(configs.pagination);
  }

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownActionOpen = false;
  };

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
      return "text-utility-brand-700-light dark:text-utility-gray-700-dark bg-primary-alt-light dark:bg-primary-alt-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus-within:ring-0 hover:no-underline";
    } else if (type === "error") {
      return "text-utility-error-700-light dark:text-utility-error-700-dark bg-primary-alt-light dark:bg-primary-alt-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark border-utility-error-300-light dark:border-utility-error-200-dark focus-within:ring-0 hover:no-underline";
    }
  };

  const withSeperator = (configs) => {
    if (configs?.seperator && configs.seperator === true) {
      return "pb-4 border-b border-primary-light dark:border-primary-dark";
    }
  };

  let searchTerm = "";

  // $: filteredItems = datas.filter((data) => {
  //   if (searchTerm === "") return data;
  //   if (data?.name) {
  //     return data.name.toLowerCase().indexOf(searchTerm.toLowerCase()) !== -1;
  //   }
  // });

  // $: filteredItems = datas.filter((data) => {
  //   if (searchTerm === "") return true;
  //   if (data) {
  //     return Object.values(data).some((value) => {
  //       if (typeof value === "string") {
  //         return value.toLowerCase().includes(searchTerm.toLowerCase());
  //       }
  //       return false;
  //     });
  //   }
  // });

  const handleSearch = (event) => {
    console.log(event.target.value);
    if (configs?.isCategories) {
      configs?.searchHandler(event.target.value, configs?.headers);
    } else {
      configs?.searchHandler(event.target.value);
    }
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
    class={`flex items-center justify-between p-4 mt-6 ${configs?.isCategories ? "mb-4 rounded-lg bg-secondary-light dark:bg-secondary-dark" : "border-t border-x border-primary-light dark:border-primary-dark rounded-t-lg"} ${withSeperator(configs)}`}
  >
    {#if configs?.headers}
      <Heading tag="h6" class="font-jakarta w-fit text-hc-graylight-900 dark:text-hc-graydark-50"
        >{configs?.isCategories ? configs?.headers : $_(toCamelCase(configs.headers))}</Heading
      >
    {/if}
    <div class="flex items-center gap-2">
      <div class="hidden lg:block">
        <Input
          id="search"
          type="search"
          class="bg-white dark:bg-hc-graydark-950 py-3"
          placeholder={$_("search")}
          on:change={handleSearch}
        >
          <Icon.SearchOutline slot="left" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
        </Input>
      </div>
      {#if configs?.actions && Object.keys(configs.actions).length > 0}
        {#each configs.actions as action (action.id)}
          {#if action?.dropdown && action.dropdown === "delete"}
            <div onfocusout={handleDropdownFocusLoss} class=" flex items-center relative">
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
                  <button
                    class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                    onclick={() => {
                      functDeleteCategories(configs?.headers),
                        (dropdownActionOpen = !dropdownActionOpen);
                    }}>{$_("delete")}</button
                  >
                </ul>
              </div>
            </div>
          {:else}
            <Button
              onclick={action?.id === "expand"
                ? () => (dropdownToggle = !dropdownToggle)
                : action.onclick}
              color="light"
              disabled={action?.disabled}
              class={`${classActionType(action.type)} min-w-max p-2 md:p-2.5 gap-1.5`}
            >
              {#if action?.icon}
                {#if action.icon === "plus"}
                  <FIlePlusIcon
                    svgClass="w-5 h-5 text-utility-brand-700-light dark:text-utility-gray-700-dark"
                  />
                {/if}
                {#if action.icon === "trash"}
                  <TrashIcon
                    svgClass="w-5 h-5 text-utility-error-700-light dark:text-utility-error-700-dark"
                  />
                {/if}
                {#if action.icon === "chevron"}
                  {#if dropdownToggle === true}
                    <ChevronDownIcon
                      svgClass="w-5 h-5 text-utility-gray-700-light dark:text-utility-gray-700-dark"
                    />
                  {:else}
                    <ChevronUpIcon
                      svgClass="w-5 h-5 text-utility-gray-700-light dark:text-utility-gray-700-dark"
                    />
                  {/if}
                {/if}
              {/if}
              {#if action?.name}
                <div class="hidden lg:block">
                  {$_(toCamelCase(action.name))}
                </div>
              {/if}
            </Button>
          {/if}
        {/each}
      {/if}
    </div>
  </div>
  {#if dropdownToggle === true}
    <div
      class={`${configs?.isCategories ? "[&>div]:sm:rounded-t-lg" : ""} [&>div]:border [&>div]:border-primary-light [&>div]:dark:border-primary-dark`}
    >
      <Table class="font-jakarta dark:border-white">
        <TableHead class="bg-hc-graylight-50 dark:bg-hc-graydark-900">
          {#if datasChecked && datasChecked.length > 0}
            <TableHeadCell class="!p-4">
              <input
                type="checkbox"
                id="checkAll"
                name="checkAll"
                checked={allChecked}
                onclick={() => {
                  if (configs?.isCategories) {
                    handleCheckAll(configs?.headers, allChecked);
                  } else {
                    handleCheckAll(allChecked);
                  }
                }}
                class="rounded h-5 w-5 dark:bg-primary-dark text-white focus:outline-none checked:bg-hc-brand-500 checked:dark:bg-hc-brand-500 focus:ring-2 focus:ring-hc-brand-500"
              />
            </TableHeadCell>
          {/if}
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
          {#each datas as data, index (data.id)}
            <TableBodyRow
              class="dark:bg-hc-graydark-950 border-hc-graylight-200 dark:border-hc-graydark-800"
            >
              {#if datasChecked}
                <TableBodyCell class="!px-4">
                  <input
                    type="checkbox"
                    id={data.id}
                    name={data.id}
                    checked={getCheckedStatus(data.id)}
                    onclick={() => {
                      if (configs?.isCategories) {
                        handleCheck(configs?.headers, data.id);
                      } else {
                        handleCheck(data.id);
                      }
                    }}
                    class="rounded h-5 w-5 dark:bg-primary-dark text-white focus:outline-none checked:bg-hc-brand-500 checked:dark:bg-hc-brand-500 focus:ring-2 focus:ring-hc-brand-500"
                  />
                </TableBodyCell>
              {/if}
              {#each Object.keys(datas[0]) as key}
                {#if key !== "actions" && key !== "id"}
                  {#if typeof data[key] === "string" || typeof data[key] === "number"}
                    <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400">
                      {#if key === "upload_date"}
                        {#if $locale === "id"}
                          {$_("createOn")}
                          {formatDistanceToNow(data[key], { addSuffix: true, locale: localeId })}
                        {:else}
                          {$_("createOn")}
                          {formatDistanceToNow(data[key], { addSuffix: true, locale: localeEn })}
                        {/if}
                      {:else}
                        {data[key]}
                      {/if}
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
                      <TableBodyCell
                        class="text-hc-graylight-600 dark:text-hc-graydark-400 px-3 py-2 md:px-6 md:py-4 pe-8"
                      >
                        <div class="flex gap-2 items-center">
                          {#if data[key]?.filename?.includes(".jpg") || data[key]?.filename?.includes(".JPG")}
                            <img src={JPGSolidIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
                          {:else if data[key]?.filename?.includes(".jpeg") || data[key]?.filename?.includes(".JPEG")}
                            <img src={JPEGSolidIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
                          {:else if data[key]?.filename?.includes(".png") || data[key]?.filename?.includes(".PNG")}
                            <img src={PNGSolidIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
                          {:else if data[key]?.filename?.includes(".pdf") || data[key]?.filename?.includes(".PDF")}
                            <img src={PDFSolidIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
                          {:else if data[key]?.filename?.includes(".doc") || data[key]?.filename?.includes(".DOC")}
                            <img src={DOCSolidIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
                          {:else if data[key]?.filename?.includes(".docx") || data[key]?.filename?.includes(".DOCX")}
                            <img src={DOCXSolidIcon} alt="" class="w-8 h-8 md:w-10 md:h-10" />
                          {/if}
                          <!-- {#if data[key]?.type?.name?.includes(".jpg") || data[key]?.type?.name?.includes(".JPG") || data[key]?.type?.name?.includes(".jpeg") || data[key]?.type?.name?.includes(".JPEG")}
                            <JpgIcon />
                          {:else if data[key]?.type?.name?.includes(".png") || data[key]?.type?.name?.includes(".PNG")}
                            <PngIcon />
                          {:else if data[key]?.type?.name?.includes(".pdf") || data[key]?.type?.name?.includes(".PDF")}
                            <PdfIcon />
                          {/if} -->
                          <!-- <JpgIcon /> -->
                          <div class="flex flex-col min-w-40 md:min-w-64">
                            <P class="line-clamp-2 fluid-sm">{data[key].name}</P>
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
                              onclick={() => action.function(datas[index].id)}
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
      {#if configs?.pagination && Object.keys(configs.pagination).length > 0 && gettingDatas === false}
        <div class="flex justify-between !border-t-0 !rounded-t-none rounded-b-lg">
          <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
            {$_("showing")}
            <span
              >{configs.pagination.page === 1
                ? 1
                : (configs.pagination.page - 1) * configs.pagination.pageSize + 1}</span
            >
            -
            <!-- {$_("to")} -->
            <span
              >{configs.pagination.page * configs.pagination.pageSize >= configs.pagination.total
                ? configs.pagination.total
                : configs.pagination.page * configs.pagination.pageSize}</span
            >
            {$_("of")}
            <span>{configs.pagination.total}</span>
            <!-- {$_("entries")} -->
          </div>
          <div class="flex">
            {#if configs.pagination.total > configs.pagination.pageSize}
              <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
                {$_("page")}
                <span>{configs.pagination.page}</span>
                {$_("of")}
                <span>{configs.pagination.pageCount}</span>
              </div>
            {/if}
            {#if configs.pagination.total > configs.pagination.pageSize}
              <div class="flex justify-end">
                <Button
                  disabled={configs.pagination.page === 1}
                  on:click={configs?.isCategories ? previous(configs?.headers) : previous}
                  color="light"
                  class="fluid-sm gap-1 md:gap-2 rounded-l-lg border-0 rounded-r-none px-3 py-2 md:px-4 md:py-2.5"
                >
                  <ArrowLeftIcon
                    svgClass="text-fg-secondary-light dark:text-fg-secondary-dark w-4 h-4 md:w-5 md:h-5"
                  />
                  {$_("previous")}
                </Button>
                <Button
                  disabled={configs.pagination.page * configs.pagination.pageSize >=
                    configs.pagination.total}
                  on:click={configs?.isCategories ? next(configs?.headers) : next}
                  color="light"
                  class="fluid-sm gap-1 md:gap-2 rounded-r-lg border-0 rounded-l-none px-3 py-2 md:px-4 md:py-2.5"
                >
                  {$_("next")}
                  <ArrowRightIcon
                    svgClass="text-fg-secondary-light dark:text-fg-secondary-dark w-4 h-4 md:w-5 md:h-5"
                  />
                </Button>
              </div>
            {/if}
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>
