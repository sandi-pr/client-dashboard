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
  import { P, Heading, Input, Button, Select, Textarea, Toggle, Modal } from "flowbite-svelte";
  import { ToastContainer, triggerToast, InputRw, TextareaRw } from "$lib";
  import { CheckCircleSolid, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import {goto} from "$app/navigation";
  import { generateUniqueId } from "$lib";

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  const catergoryId = params.get("catergoryId");
  const edit = params.get("edit");
  let onEdit = $state(false);

  let popupModal = $state(false);

  let data = $state({});

  let toastCreated = $state(false);
  let toastUpdated = $state(false);
  let toastError = $state(false);
  let toastDeleted = $state(false);
  let counter = $state(6);

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
    if (catergoryId === null) {
      triggerToast("error");
      // setTimeout(async () => {
        goto(await call.getPage("client-dashboard", "crm-categories"));
      // }, 3000);
      return;
    };
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
          id: generateUniqueId(response.result.data[0].id, "category"),
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
          triggerToast("updated");
          // setTimeout(async () => {
            goto(await call.getPage("client-dashboard", "crm-categories"));
          // }, 3000);
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
        triggerToast("deleted");
        // setTimeout(async () => {
          goto(await call.getPage("client-dashboard", "crm-categories"));
        // }, 3000);
      }
    } catch (error) {
      console.log(error);
    }
  }

  onMount(async () => {
    if (edit === "true") onEdit = true;
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
            <Button
              onclick={async function () {
                goto(await call.getPage("client-dashboard", "crm-categories"));
                return;
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("back")}</Button
            >
            <Button
              onclick={() => (popupModal = true)}
              outline color="red"
              class="px-3 focus-within:!ring-1">{$_("delete")}</Button
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
                <InputRw {onEdit} bind:value={data.category_name} />
              </div>
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("categoryId")}
                </P>
                <InputRw onEdit={false} bind:value={data.id} />
              </div>
              <div class="flex flex-col gap-2 justify-center w-max px-4">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("active")}</P
                >
                <Toggle disabled={!onEdit} bind:checked={data.is_active}></Toggle>
              </div>
            </div>
            <div class="grid grid-cols-1">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("description")}</P
                >
                <TextareaRw {onEdit} bind:value={data.description} />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <ToastContainer />
</div>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteCategory} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>