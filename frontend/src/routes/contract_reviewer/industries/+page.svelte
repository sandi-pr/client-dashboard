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
  import { onMount } from "svelte";
  import { triggerToast } from "$lib";
  import { _ } from "svelte-i18n";
  import PlusIcon from "$lib/icon/plus.svg";
  import DataTable from "$lib/components/DataTable.svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { Button, Modal, Input, Textarea } from "flowbite-svelte";
  import { StateCard } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";

  const call = new CR3Call(globalThis.location.origin);

  let ready = $state(false);

  let industryModal = $state(false);

  let deleteModal = $state(false);

  let update = $state(false);

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });

  let datas = $state([]);

  let configs = $state({});

  let tmpIndustryId = $state("");
  let tmpIndustryName = $state("");
  let tmpIndustryDescription = $state("");

  async function createIndustry() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "create-industry", {
        URLQuery: {},
        URLParams: {},
        Params: {
          industry_name: tmpIndustryName,
          industry_description: tmpIndustryDescription,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      industryModal = false;
      triggerToast("created");
      getIndustries();
    }
  }

  async function getIndustries() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-industries", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: pagination.page,
          page_size: pagination.pageSize,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      if (response.result.data.length > 0) {
        configs = {
          pagination: response.result.meta.pagination,
        };
        let promise = response.result.data.map((industry) => {
          return {
            id: industry.id,
            name: industry.attributes.industry_name,
            description: industry.attributes.industry_description,
            actions: {
              objects: [
                {
                  onclick: () => {
                    tmpIndustryId = industry.id;
                    tmpIndustryName = industry.attributes.industry_name;
                    tmpIndustryDescription = industry.attributes.industry_description;
                    update = true;
                    industryModal = true;
                  },
                  name: "Edit",
                },
                {
                  onclick: () => {
                    tmpIndustryId = industry.id;
                    deleteModal = true;
                  },
                  name: "Delete",
                },
              ],
            },
          };
        });

        datas = await Promise.all(promise);
        ready = true;
      } else {
        datas = [];
      }
    }
  }

  async function deleteIndustry() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "delete-industry", {
        URLQuery: {},
        URLParams: {
          id: tmpIndustryId.toString(),
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      tmpIndustryId = "";
      triggerToast("deleted");
      getIndustries();
    }
  }

  async function updateIndustry() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      let response = await call.call("client-dashboard", "update-industry", {
        URLQuery: {},
        URLParams: {
          id: tmpIndustryId.toString(),
        },
        Params: {
          industry_name: tmpIndustryName,
          industry_description: tmpIndustryDescription,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      tmpIndustryId = "";
      tmpIndustryName = "";
      tmpIndustryDescription = "";

      update = false;
      industryModal = false;

      triggerToast("updated");
      getIndustries();
    }
  }

  async function nextPage() {
    if (pagination.page < pagination.pageCount) return;
    pagination.page = pagination.page + 1;
    getIndustries();
  }

  async function previousPage() {
    if (pagination.page > 1) return;
    pagination.page = pagination.page - 1;
    getIndustries();
  }

  onMount(() => {
    getIndustries();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text={$_("manageIndustries")} />
          <Button
            onclick={() => {
              industryModal = true;
            }}
            icon={PlusIcon}
            ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" /> {$_("addNew")}</Button
          >
        </div>

        <!-- Person Table Data -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs bind:next={nextPage} bind:previous={previousPage} />
        {:else}
          <div class="pt-8">
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<Modal title={$_("createNewIndustry")} bind:open={industryModal} size="lg">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("industryName")}*`} />
      <Input
        bind:value={tmpIndustryName}
        placeholder={$_("whatIsYourIndustryName")}
        class="bg-primary-light dark:bg-secondary-dark"
      />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Textarea
        bind:value={tmpIndustryDescription}
        placeholder={$_("placholderDescription")}
        class="bg-primary-light dark:bg-secondary-dark"
      />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={update ? updateIndustry : createIndustry}
      disabled={!tmpIndustryName}
      class="w-full">{update ? $_("updateIndustry") : $_("addNewIndustry")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteIndustry} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button
      color="alternative"
      onclick={() => {
        deleteModal = false;
        tmpIndustryId = "";
      }}>{$_("noCancel")}</Button
    >
  </div>
</Modal>
