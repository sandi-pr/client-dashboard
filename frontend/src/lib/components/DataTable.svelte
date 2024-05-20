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
  import RefreshIcon from "../icon/refresh.svelte";
  import EyeIcon from "../icon/eye.svelte";
  import CheckIcon from "../icon/check.svg";
  import XIcon from "../icon/x.svg";
  import CheckWhiteIcon from "../icon/check-white.svg";
  import ID from "../icon/ID.svg";
  import { _ } from "svelte-i18n";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  // import { toCamelCase } from "../utils/toCamelCase.js";

  export let datas;
  export let configs;
  export let next = () => {};
  export let previous = () => {};
  // export let divider = false;

  let pages;
  if (configs?.pagination && Object.keys(configs.pagination).length > 0) {
    pages = pagesBuilder(configs.pagination);
  }

  const splitCamelCase = (value) => {
    return value.replace(/([a-z0-9])([A-Z])/g, "$1 $2").replace(/_/g, " ");
  };
  

  const normalContent = (key) => {
    return (
      key !== "id" &&
      key !== "actions" &&
      key !== "status" &&
      key !== "parentUnit" &&
      key !== "unit" &&
      key !== "role" &&
      key !== "country" &&
      key !== "avatar" &&
      key !== "parentUnit" &&
      key !== "unit" &&
      key !== "role" &&
      key !== "member" &&
      key !== "recipients" &&
      key !== "category" &&
      key !== "authority" &&
      key !== "attendance_status"
    );
  };

  const pillContent = (key) => {
    return (
      key === "status" ||
      key === "parentUnit" ||
      key === "unit" ||
      key === "role" ||
      key === "category" ||
      key === "authority" ||
      key === "attendance_status"
    );
  };

  const classStatusType = (type) => {
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
    } else if (type === "error") {
      return "bg-hc-error-50 border-hc-error-200 text-hc-error-700";
    } else if (type === "warning") {
      return "bg-hc-warning-50 border-hc-warning-200 text-hc-warning-700";
    } else if (type === "info") {
      return "bg-hc-graylight-50 dark:bg-hc-graydark-900 hover:bg-hc-graylight-100 hover:dark:bg-hc-graydark-800 focus-within:ring-0 hover:no-underline";
    } else if (type === "brand") {
      return "text-white bg-hc-brand-600 dark:bg-hc-brand-600 hover:bg-hc-brand-700 dark:hover:bg-hc-brand-700 focus-within:ring-0 hover:no-underline";
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

  // onMount(() => {
  // Object.keys(datas[0]).forEach((key) => {
  //   console.log(key);
  // });
  // });
</script>

<div>
  <div class={`flex items-center justify-between mb-4 ${withSeperator(configs)}`}>
    {#if configs?.headers}
      <Heading tag="h6" class="font-jakarta w-fit text-hc-graylight-900 dark:text-hc-graydark-50"
        >{configs.headers}</Heading
      >
    {/if}
    <div>
      {#if configs?.actions && Object.keys(configs.actions).length > 0}
        {#each configs.actions as action (action.id)}
          <Button href={action.href} color="light" class={`${classActionType(action.type)}`}>
            {#if action?.icon}
              {#if action.icon === "plus"}
                <Icon.PlusOutline class="w-4 h-4" />&nbsp;
              {/if}
            {/if}
            {action.name}
          </Button>
        {/each}
      {/if}
    </div>
    {#if configs?.actionUrl}
      <Button
        href={configs.actionUrl.url}
        color="light"
        class="bg-hc-graylight-50 dark:bg-hc-graydark-900 border-hc-graylight-200 dark:border-hc-graydark-800 hover:bg-hc-graylight-100 hover:dark:bg-hc-graydark-800 focus-within:ring-0"
        >{configs.actionUrl.text}</Button
      >
    {/if}
  </div>

  {#if configs?.search || configs?.documentManagement}
    <div class={`flex ${configs?.search ? "mb-8 mt-6" : ""} gap-x-4`}>
      {#if configs.search}
        <Input
          id="search"
          type="search"
          class="bg-white dark:bg-hc-graydark-950"
          placeholder={$_("search")}
          on:change={handleSearch}
        >
          <Icon.SearchOutline slot="left" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
        </Input>
        {#if configs.documentManagement}
          <div class="flex gap-2">
            <Button
              outline
              color="red"
              class=" w-max hover:!bg-transparent hover:text-hc-error-600 !border-hc-error-300 dark:!border-hc-error-400 focus-within:!ring-1 px-3 py-[10px] gap-1"
            >
              <Icon.TrashBinOutline
                class="w-5 h-5 text-fg-error-secondary-light dark:text-fg-error-secondary-dark"
              />
              {"batchDelete"}
            </Button>
            <Button type="submit" class="w-max px-3 py-[10px] gap-1 focus-within:!ring-1">
              <img alt="refresh" src={CheckWhiteIcon} class="w-5 h-5" />
              {"batchApprove"}
            </Button>
          </div>
        {/if}
      {/if}
    </div>
  {/if}

  <Table shadow class="font-jakarta dark:border-white">
    {#if configs?.doubleHeader}
      <TableHead class="bg-hc-graylight-50 dark:bg-hc-graydark-900">
        <TableHeadCell class="flex items-center">
          <img
            alt="Home Icon"
            src={configs?.iconDoubleHeader}
            class="w-5 h-5 mr-4"
          />{configs.doubleHeader}</TableHeadCell
        >
        {#if datas.length > 0}
          {#each Object.keys(datas[0]) as key, index (key)}
            {#if index < Object.keys(datas[0]).length - 2}
              <TableHeadCell></TableHeadCell>
            {/if}
          {/each}
        {/if}
      </TableHead>
    {/if}
    <TableHead class="bg-hc-graylight-50 dark:bg-hc-graydark-900">
      {#if datas.length > 0}
        {#each Object.keys(datas[0]) as key}
          {#if normalContent(key)}
            <TableHeadCell>{splitCamelCase($_(toCamelCase(key)))}</TableHeadCell>
          {/if}
        {/each}
        {#each Object.keys(datas[0]) as key}
          {#if pillContent(key)}
            <TableHeadCell>{key}</TableHeadCell>
          {/if}
        {/each}
        {#if datas[0]?.avatar && Object.keys(datas[0].avatar).length > 0}
          <TableHeadCell>avatar</TableHeadCell>
        {/if}
        {#if datas[0]?.recipients && Object.keys(datas[0].recipients).length > 0}
          <TableHeadCell>{"recipients"}</TableHeadCell>
        {/if}
        {#if datas[0]?.country}
          <TableHeadCell>{$_("country")}</TableHeadCell>
        {/if}
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
          {#each Object.keys(datas[0]) as key}
            {#if normalContent(key)}
              <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400">
                {data[key]}
              </TableBodyCell>
            {/if}
          {/each}
          {#if datas[0]?.parentUnit || datas[0]?.unit || datas[0]?.role || datas[0]?.category || datas[0]?.authority}
            {#each Object.keys(datas[0]) as key}
              {#if key === "parentUnit" || key === "unit" || key === "role" || key === "category" || key === "authority"}
                <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400">
                  {#if data[key]?.type}
                    <p
                      class={`py-[2px] px-2 border-[1px] rounded-xl text-center w-min ${classStatusType(
                        data[key].type,
                      )}`}
                    >
                      {data[key].name}
                    </p>
                  {:else}
                    <p>
                      {data[key]}
                    </p>
                  {/if}
                </TableBodyCell>
              {/if}
            {/each}
          {/if}
          {#if datas[0]?.status && Object.keys(datas[0].status).length > 0}
            {#each Object.keys(datas[0]) as key}
              {#if key === "status"}
                <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400"
                  ><p
                    class={`py-[2px] px-2 border-[1px] rounded-xl text-center w-min ${classStatusType(
                      data.status.type,
                    )}`}
                  >
                    {data[key].name}
                  </p>
                </TableBodyCell>
              {/if}
            {/each}
          {/if}
          {#if datas[0]?.avatar && Object.keys(datas[0].avatar).length > 0}
            <div class="flex">
              {#each data.avatar.slice(0, 2) as avatar (avatar.id)}
                <Avatar src={avatar.photo} stacked />
              {/each}
              {#if data.avatar.length > 2}
                <Avatar stacked href="/" class="bg-gray-700 text-white hover:bg-gray-600 text-sm"
                  >+{data.avatar.length - 2}</Avatar
                >
              {/if}
            </div>
          {/if}
          {#if datas[0]?.member}
            {#each Object.keys(datas[0]) as key}
              {#if key === "member"}
                <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400">
                  <div class="flex">
                    <p class="mr-2">{`${data.member.length} users`}</p>
                    {#each data.member as member, index (member)}
                      {#if index < 5}
                        {#if member.url === "" || member.url === "#"}
                          <p>
                            <Avatar stacked class="w-6 h-6" />
                          </p>
                        {:else}
                          <p>
                            <Avatar src={member.url} class="w-6 h-6" />
                          </p>
                        {/if}
                      {/if}
                    {/each}
                    {#if data.member.length > 5}
                      <Avatar src="#" class="w-6 h-6" alt={`+${data.member.length - 5}`} />
                    {/if}
                  </div>
                </TableBodyCell>
              {/if}
            {/each}
          {/if}
          {#if datas[0]?.recipients}
            {#each Object.keys(datas[0]) as key}
              {#if key === "recipients"}
                <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400">
                  <div class="flex">
                    {#each data.recipients as recipients, index (recipients)}
                      {#if index < 5}
                        {#if recipients.url === ""}
                          <Avatar
                            stacked
                            class="bg-gray-100 text-slate-600 text-sm"
                            style="width: 40px; height: 40px; overflow: hidden;"
                          />
                        {:else}
                          <Avatar
                            src={recipients.url}
                            alt="img"
                            style="width: 100%; height: 100%; object-fit: cover; border-radius: 50%;"
                          />
                        {/if}
                      {/if}
                    {/each}
                    {#if data.recipients.length > 5}
                      <Avatar stacked alt={`+${data.recipients.length - 5}`}
                        >{`+${data.recipients.length - 5}`}</Avatar
                      >
                    {/if}
                  </div>
                </TableBodyCell>
              {/if}
            {/each}
          {/if}
          {#if datas[0]?.country}
            {#each Object.keys(datas[0]) as key}
              {#if key === "country"}
                <TableBodyCell
                  class="flex items-center text-hc-graylight-600 dark:text-hc-graydark-400"
                  ><p>
                    {#if data.country === "ID"}
                      <img src={ID} alt="ID" class="w-6 h-6" />
                    {/if}
                  </p>
                  <p class="ml-2 text-md">{data.country}</p>
                </TableBodyCell>
              {/if}
            {/each}
          {/if}
          {#if datas[0]?.attendance_status && Object.keys(datas[0].attendance_status).length > 0}
            {#each Object.keys(datas[0]) as key}
              {#if key === "attendance_status"}
                <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400"
                  ><p
                    class={`py-[2px] px-2 border-[1px] rounded-xl text-center w-min ${classStatusType(
                      data.attendance_status.type,
                    )}`}
                  >
                    {data[key].name}
                  </p>
                </TableBodyCell>
              {/if}
            {/each}
          {/if}
          {#if datas[0]?.attendance_status}
            {#each Object.keys(datas[0]) as key}
                {#if key === "attendance_status"}
                  <TableBodyCell
                    class="flex items-center text-hc-graylight-600 dark:text-hc-graydark-400"
                    ><p>
                      {#if data.country === "ID"}
                        <img src={ID} alt="ID" class="w-6 h-6" />
                      {/if}
                    </p>
                    <p class="ml-2 text-md">{data.country}</p>
                  </TableBodyCell>
                {/if}
              {/each}
          {/if}
          {#if datas[0]?.actions && Object.keys(datas[0].actions).length > 0}
            {#each Object.keys(datas[0]) as key}
              {#if key === "actions"}
                <TableBodyCell class="text-hc-graylight-600 dark:text-hc-graydark-400">
                  <div class="flex justify-end gap-8">
                    {#each datas[index].actions.objects as action (action.name)}
                      {#if action?.iconOnly && action.iconOnly === true}
                        {#if configs?.documentManagement}
                          <div class="flex gap-1">
                            {#if action?.actionButton?.approve}
                              <Button
                                outline
                                color="green"
                                class="hover:!bg-transparent !border-hc-success-500 dark:!border-hc-success-400 focus-within:!ring-1 !p-2"
                              >
                                <img alt="refresh" src={CheckIcon} class="w-5 h-5" />
                              </Button>
                            {/if}
                            {#if action?.actionButton?.reject}
                              <Button
                                outline
                                color="red"
                                class="hover:!bg-transparent !border-hc-error-300 dark:!border-hc-error-400 focus-within:!ring-1 !p-2"
                              >
                                <img alt="refresh" src={XIcon} class="w-5 h-5" />
                              </Button>
                            {/if}
                            {#if action?.actionButton?.recall}
                              <Button color="white" class="!p-2 focus-within:!ring-1 " size="lg">
                                <RefreshIcon
                                  class="w-5 h-5 text-fg-tertiary-light dark:text-fg-tertiary-dark"
                                />
                              </Button>
                            {/if}
                            <Button
                              href="/documents/single?documentId={data.id}"
                              color="white"
                              class="!p-2 focus-within:!ring-1"
                              size="lg"
                            >
                              <EyeIcon
                                svgClass="w-5 h-5 text-fg-tertiary-light dark:text-fg-tertiary-dark"
                              />
                            </Button>
                          </div>
                        {:else if action.name === "Delete"}
                          <a href={action.href}>
                            <Icon.TrashBinOutline class="w-5 h-5 text-gray-600" />
                          </a>
                        {:else if action.name === "Edit"}
                          <a href={action.href}>
                            <Icon.PenOutline class="w-5 h-5 text-gray-600" />
                          </a>
                        {:else}
                          <a href={action.href}>
                            <Icon.TrashBinOutline class="w-3 h-3 text-gray-600" />
                          </a>
                        {/if}
                      {:else}
                        <a
                          class="hover:no-underline text-hc-brand-600 font-semibold"
                          href={action.href}>{$_(toCamelCase(action.name))}</a
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
