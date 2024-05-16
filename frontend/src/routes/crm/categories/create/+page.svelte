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
  import { P, Heading, Input, Button, Select, Textarea, Toggle, Toast } from "flowbite-svelte";
  import { CheckCircleSolid, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { fly } from "svelte/transition";
  import { _ } from "svelte-i18n";
  import { CR3Call } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  let onEdit = $state(true);

  let buttonDisable = $state(true);

  let toastCreated = $state(false);
  let toastError = $state(false);
  let counter = $state(6);

  let data = $state({
    category_name: "",
    description: "",
    is_active: true,
  });

  async function createCategory() {
    if (data.category_name === "" || data.description === "") return;

    let token = sessionStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await call.call("client-dashboard", "create-category", {
          URLQuery: {},
          URLParams: {},
          Params: {
            category_name: data.category_name,
            description: data.description,
            is_active: data.is_active ? "Yes" : "No",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        console.log("kelinci", response);
        if (response.status === 200) {
          triggerToast("created");
          await call.goToPage("client-dashboard", "crm-categories");
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  function triggerToast(target) {
    if (target === "created") {
      toastCreated = true;
      counter = 6;
      timeout();
    } else if (target === "error") {
      toastError = true;
      counter = 6;
      timeout();
    } else {
      console.log("target not found");
    }
  }

  function timeout() {
    if (--counter > 0) return setTimeout(timeout, 1000);
    open = false;
    toastError = false;
  }

  $effect(() => {
    buttonDisable =
      data.category_name.length < 1 ||
      data.description.length < 1;
  });

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
              {$_("createNewCategoriesProduct")}
            </Heading>
          </div>
        </div>
        <div class="flex gap-2">
          {#if onEdit}
            <Button
              onclick={async function () {
                await call.goToPage("client-dashboard", "crm-categories");
                return;
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button disabled={buttonDisable} onclick={createCategory} class="px-3 focus-within:!ring-1"
              >{$_("create")}</Button
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
                  {$_("categoryName")}*
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
                  >{$_("description")}*</P
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
  <div class="fixed bottom-5 end-5">
    <Toast color="green" transition={fly} params={{ x: 200 }} bind:open={toastCreated}>
      <svelte:fragment slot="icon">
        <CheckCircleSolid class="w-5 h-5" />
      </svelte:fragment>
      {$_("successfullyCreated")}
    </Toast>
    <Toast
      class="!max-w-sm"
      color="red"
      transition={fly}
      params={{ x: 200 }}
      bind:open={toastError}
    >
      <svelte:fragment slot="icon">
        <ExclamationCircleSolid class="w-5 h-5" />
      </svelte:fragment>
      {$_("sometingWentWrongPleaseTryAgain")}
    </Toast>
  </div>
</div>
