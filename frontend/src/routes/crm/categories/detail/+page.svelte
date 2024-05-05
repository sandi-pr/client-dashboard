<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: 
  - Azzis Arswendo <azzis@hexacode.org>
  - Dean Ramhan <dean@hexacode.org>
  - Taqie Fadlillah <taqief@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->
<script>
  import { P, Heading, Input, Button, Select, Textarea, Toggle } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  const catergoryId = params.get("catergoryId");
  let onEdit = $state(false);

  let data = $state({});

  function formatDate(inputDate) {
    const date = new Date(inputDate);

    const formattedDate = `${(date.getMonth() + 1).toString().padStart(2, "0")}/${date
      .getDate()
      .toString()
      .padStart(2, "0")}/${date.getFullYear()}, ${date
      .getHours()
      .toString()
      .padStart(2, "0")}:${date
      .getMinutes()
      .toString()
      .padStart(2, "0")} ${date.getHours() >= 12 ? "PM" : "AM"}`;

    return formattedDate;
  }

  async function getCategory() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-category", {
        URLQuery: {},
        URLParams: {
          id: catergoryId,
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });
      console.log("kelinci", response);
      if (response.status === 200) {
        data = {
          createdAt: formatDate(response.result.data[0].attributes.createdAt),
          category_name: response.result.data[0].attributes.category_name,
          description: response.result.data[0].attributes.description,
          is_active: response.result.data[0].attributes.is_active,
        };
      }
    }
  }

  async function updateCategory() {
    if (data.category_name === "" || data.description === "") return;
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await call.call("client-dashboard", "update-category", {
          URLQuery: {},
          URLParams: {
            id: catergoryId,
          },
          Params: {
            category_name: data.category_name,
            description: data.description,
            is_active: data.is_active ? "Yes" : "No",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status === 200) {
          await call.goToPage("client-dashboard", "crm-categories");
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function deleteCategory() {
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "delete-category", {
        URLQuery: {},
        URLParams: {
          id: catergoryId,
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });
      if (response.status === 200) {
        await call.goToPage("client-dashboard", "crm-categories");
      }
    } catch (error) {
      console.log(error);
    }
  }

  onMount(async () => {
    getCategory();
  });

  // let dropdownOpen = false;
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- Header -->
      <div class="flex items-center pb-8">
        <div class="flex flex-col gap-2 grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {data.category_name}
              <Button
                outline
                color="red"
                class="hover:!bg-transparent !border-hc-graylight-200 dark:!border-hc-graydark-600 focus-within:!ring-1 !p-2"
              >
                <CopyIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
              </Button>
            </Heading>
          </div>
          <P
            class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
          >
            {$_("createOn")}{` ${data.createdAt}`}
          </P>
        </div>
        <div class="flex gap-2">
          {#if onEdit}
            <Button
              onclick={() => {
                onEdit = false;
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button onclick={updateCategory} class="px-3 focus-within:!ring-1">{$_("save")}</Button>
          {:else}
            <Button onclick={deleteCategory} color="light" class="px-3 focus-within:!ring-1"
              >{$_("delete")}</Button
            >
            <Button
              onclick={() => {
                onEdit = true;
              }}
              class="px-3 focus-within:!ring-1">{$_("edit")}</Button
            >
          {/if}
        </div>
      </div>
      <div class="flex flex-col gap-4 pt-8">
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
                  {$_("categoryName")}
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.category_name}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2 justify-center w-max px-4">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("active")}</P
                >
                <Toggle bind:checked={data.is_active}></Toggle>
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
                  bind:value={data.description}
                  rows="4"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
