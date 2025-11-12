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
  import Markdown from "svelte-exmarkdown";
  import { gfmPlugin } from "svelte-exmarkdown/gfm";
  import { _ } from "svelte-i18n";
  import {
    P,
    Table,
    Modal,
    Button,
    TableBody,
    TableHead,
    TableBodyRow,
    TableBodyCell,
    TableHeadCell,
  } from "flowbite-svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";

  import { toCamelCase, getIcons } from "$lib";
  import { FeaturedIcon, Badge, Typography } from "$lib";
  import { pagesBuilder } from "../utils/pagesBuilder";
  import DatatablePagination from "./DatatablePagination.svelte";

  const JpgIcon = getIcons("JpgIcon");
  const InfoCircleIcon = getIcons("infoCircle");
  const Stars01Icon = getIcons("stars01");
  const PDFOutlineIcon = getIcons("PDFOutline");
  const CSVSolidIcon = getIcons("CSVSolid");

  const plugins = [gfmPlugin()];

  export let datas;
  export let configs;
  export let next = () => {};
  export let previous = () => {};
  export let goToPage = () => {};

  let sectionRefModal = false;
  let reviewModal = false;
  let suggestionModal = false;
  let dropdownInput = true;
  let dropdownOutput = true;
  let sectionRefData = {};
  let reviewData = "";
  let suggestionData = "";
  let tmpData = {};
  let tmpRefDoc = {};

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

  const getObjectDataTypeClass = (type, color) => {
    if (type === "pill") {
      return classPillType(color);
    }
  };
</script>

<div>
  <div
    class="[&>div]:sm:rounded-lg [&>div]:border-2 overx [&>div]:border-primary-light [&>div]:dark:border-primary-dark"
  >
    <Table>
      <TableHead class="bg-hc-graylight-50 dark:bg-hc-graydark-900">
        {#if datas.length > 0}
          {#each Object.keys(datas[0]) as key}
            {#if key !== "actions" && key !== "id"}
              {#if key === "ai_suggestion"}
                <TableHeadCell
                  class="bg-brand-solid-light dark:bg-brand-solid-dark text-primary-onbrand-light dark:text-primary-onbrand-dark"
                  >{splitCamelCase($_(toCamelCase(key)))}</TableHeadCell
                >
              {:else}
                <TableHeadCell>{splitCamelCase($_(toCamelCase(key)))}</TableHeadCell>
              {/if}
            {/if}
          {/each}
          {#if datas.some((data) => data.actions && Object.keys(data.actions).length > 0)}
            <TableHeadCell>Actions</TableHeadCell>
          {/if}
        {/if}
      </TableHead>
      <TableBody>
        {#each configs?.search ? filteredItems : datas as data, index (data.id)}
          <TableBodyRow
            class="w-full dark:bg-hc-graydark-950 border-hc-graylight-200 dark:border-hc-graydark-800"
          >
            {#each Object.keys(data) as key}
              {#if key !== "actions" && key !== "id"}
                {#if typeof data[key] === "string" || typeof data[key] === "number"}
                  {#if key === "ai_suggestion"}
                    <TableBodyCell
                      on:click={() => {
                        suggestionModal = true;
                        suggestionData = data[key];
                      }}
                      class="w-1/4 hover:cursor-pointer bg-brand-primary-light dark:bg-brand-primary-dark"
                    >
                      <P
                        size="sm"
                        class="underline line-clamp-3 text-brand-primary-light dark:text-brand-primary-dark"
                      >
                        <Markdown md={data[key]} {plugins} />
                      </P>
                    </TableBodyCell>
                  {:else if key === "ai_review" || key === "ai_extraction"}
                    <TableBodyCell
                      on:click={() => {
                        reviewModal = true;
                        reviewData = data[key];
                      }}
                      class="w-1/4 hover:cursor-pointer"
                    >
                      <P
                        size="sm"
                        class="text-tertiary-light dark:text-tertiary-dark underline line-clamp-2"
                      >
                        <Markdown md={data[key]} {plugins} />
                      </P>
                    </TableBodyCell>
                  {:else}
                    <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400">
                      {data[key]}
                    </TableBodyCell>
                  {/if}
                {:else if typeof data[key] === "object" && key === "section"}
                  <TableBodyCell
                    on:click={() => {
                      tmpData = data;
                      tmpRefDoc = data[key];
                      sectionRefModal = true;
                      sectionRefData = data[key];
                    }}
                    tdClass="hover:cursor-pointer underline"
                  >
                    <Typography preset="Text sm/Medium" text={data[key].name} />
                  </TableBodyCell>
                {:else if typeof data[key] === "object" && key === "reference_documents"}
                  <TableBodyCell
                    on:click={() => {
                      tmpData = data;
                      tmpRefDoc = data[key];
                      sectionRefModal = true;
                      sectionRefData = data[key];
                    }}
                    tdClass="w-1/6 hover:cursor-pointer"
                  >
                    <div class="flex items-center justify-start gap-3">
                      {#if data[key].name.includes(".pdf") || data[key].name.includes(".PDF")}
                        <PDFOutlineIcon svgClass="w-32 h-10" />
                      {:else if data[key].name.includes(".csv") || data[key].name.includes(".CSV")}
                        <CSVSolidIcon svgClass="w-10 h-10" />
                      {/if}
                      <div class="line-clamp-2">
                        <P size="sm">{data[key].name.replace(/\.[^/.]+$/, "")}</P>
                      </div>
                    </div>
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
                        <JpgIcon svgClass="w-6 h-6" />
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
            {#if data.actions && Object.keys(data.actions).length > 0}
              <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400">
                <div class="flex justify-end gap-8">
                  {#each Object.entries(data.actions) as [actionKey, actionValue] (actionKey)}
                    {#if actionKey === "status" && actionValue === "Done"}
                      <div
                        class={`py-[2px] px-2 border-[1px] rounded-xl text-center w-min ${getObjectDataTypeClass("pill", data.actions.result === "Changed" ? "success" : "error")}`}
                      >
                        {data.actions.result}
                      </div>
                    {:else if actionKey === "status" && actionValue === "Processing"}
                      {#if data.actions.reject}
                        <button
                          on:click={data.actions.reject}
                          class="font-semibold text-tertiary-light dark:text-tertiary-dark"
                          >Reject</button
                        >
                      {/if}
                      {#if data.actions.change}
                        <button
                          on:click={data.actions.change}
                          class="font-semibold text-brand-tertiary-alt-light dark:text-brand-tertiary-alt-dark"
                          >Change</button
                        >
                      {/if}
                    {/if}
                  {/each}
                </div>
              </TableBodyCell>
            {/if}
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>
</div>

<DatatablePagination {configs} {goToPage} {next} {previous} />

<Modal
  bind:open={sectionRefModal}
  dialogClass="fixed top-0 start-0 end-0 h-modal md:inset-0 md:h-full z-50 w-full p-8 md:p-12 flex"
  bodyClass="bg-secondary-light dark:bg-secondary-dark rounded-lg"
  size="lg"
  outsideclose
>
  <div class="flex flex-col gap-4 p-8 md:p-12 bg-primary-light dark:bg-primary-dark">
    <div class="flex flex-col gap-2">
      <Typography
        preset="Text lg/Semibold"
        text={`No.${tmpData.no} - ${$_("section")} "${tmpData.section.name}"`}
      />
      <div class="flex items-center justify-start gap-1.5">
        {#if tmpRefDoc.name.includes(".pdf") || tmpRefDoc.name.includes(".PDF")}
          <img src={PDFOutlineIcon} alt="" class="w-6 h-6" />
        {/if}
        <Typography
          preset="Text sm/Regular"
          text={tmpRefDoc.name.replace(/\.[^/.]+$/, "")}
          color="secondary"
        />
      </div>
    </div>
    <Button
      on:click={() => (dropdownInput = !dropdownInput)}
      class="flex justify-between border-0 bg-secondary-light dark:bg-secondary-dark dark:hover:bg-secondary-hover-dark"
      color="light"
    >
      <div class="flex items-center gap-2">
        <P class="text-xs font-semibold text-utility-gray-500-light dark:text-utility-gray-500-dark"
          >{$_("PROMPT")}</P
        >
        <Badge text="Input" type="warning" />
      </div>
      <ChevronDownOutline
        class="w-8 h-8 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
      />
    </Button>
    {#if dropdownInput}
      <div class="flex flex-col gap-1">
        <Markdown md={sectionRefData.prompt} {plugins}></Markdown>
      </div>
    {/if}
    <Button
      on:click={() => (dropdownOutput = !dropdownOutput)}
      class="flex justify-between border-0 bg-secondary-light dark:bg-secondary-dark dark:hover:bg-secondary-hover-dark"
      color="light"
    >
      <div class="flex items-center gap-2">
        <P class="text-xs font-semibold text-utility-gray-500-light dark:text-utility-gray-500-dark"
          >{$_("AI READ")}</P
        >
        <Badge text="Output" />
      </div>
      <ChevronDownOutline
        class="w-8 h-8 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
      />
    </Button>
    {#if dropdownOutput}
      <div class="flex flex-col gap-1">
        <Markdown md={sectionRefData.result} {plugins}></Markdown>
      </div>
    {/if}
  </div>
</Modal>

<Modal bind:open={reviewModal} autoclose outsideclose>
  <div class="flex flex-col gap-4">
    <FeaturedIcon classesSvg={InfoCircleIcon} />
    <div class="flex flex-col gap-1">
      <Typography preset="Text lg/Semibold" text={$_("aiReview")} />
      <Markdown md={reviewData} {plugins}></Markdown>
    </div>
  </div>
</Modal>

<Modal bind:open={suggestionModal} autoclose outsideclose>
  <div class="flex flex-col gap-4">
    <FeaturedIcon classesSvg={Stars01Icon} />
    <div class="flex flex-col gap-1">
      <Typography preset="Text lg/Semibold" text={$_("aiSuggestion")} />
      <Markdown md={suggestionData} {plugins}></Markdown>
    </div>
  </div>
</Modal>
