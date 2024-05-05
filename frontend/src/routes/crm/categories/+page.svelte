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
  import { P, Card, Avatar, Heading, A } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import UserPlusIcon from "$lib/icon/user-plus-02.svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import Chart from "$lib/components/Chart.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import SectionHeader from "$lib/ui/header/SectionHeader.svelte";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { Button } from "flowbite-svelte";

  const call = new CR3Call(globalThis.location.origin);

  let personsLatest = [];
  let documentsReady = $state(false);
  let actionItems = $state([]);

  let pagination = { page: 1 };

  let datas = $state([]);

  let configs = $state({});

  function getStatus(is_active) {
    if (is_active === "Yes") {
      return {
        type: "success",
        name: "Active",
      };
    } else {
      return {
        type: "danger",
        name: "Inactive",
      };
    }
  }

  async function getCategories() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-categories", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: pagination.page,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status === 200) {
        configs = {
          search: true,
          pagination: response.result.meta.pagination,
        };
        let promises = response.result.data.map(async (category) => {
          return {
            id: category.id,
            category_name: category.attributes.category_name,
            description: category.attributes.description,
            status: getStatus(category.attributes.is_active),
            actions: {
              objects: [
                {
                  href:
                    `${await call.getPage("client-dashboard", "crm-category-detail")}` +
                    "?catergoryId=" +
                    category.id,
                  name: "Edit",
                },
              ],
            },
          };
        });

        datas = await Promise.all(promises);
        documentsReady = true;
      }
    }
  }

  async function addCategory() {
    await call.goToPage("client-dashboard", "crm-category-create");
    return;
  }

  onMount(async () => {
    // getActionItems();
    getCategories();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text="Manage Categories" />
          <Button onclick={addCategory} icon={PlusIcon}
            >Add new <img class="w-4 h-4" src={PlusIcon} alt="plus-icon" /></Button
          >
        </div>

        <!-- Person Table Data -->
        <div>
          {#if documentsReady}
            <DataTable bind:datas bind:configs />
          {:else}
            <DataTable {datas} {configs} />
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>
