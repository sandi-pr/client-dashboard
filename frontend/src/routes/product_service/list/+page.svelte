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
  import { _ } from "svelte-i18n";
  import PlusIcon from "$lib/icon/plus.svg";
  import DataTable from "$lib/components/DataTable.svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { Button } from "flowbite-svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import { StateCard, check } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { goto } from "$app/navigation";

  import { call } from "$lib";

  let personsLatest = [];
  let documentsReady = $state(false);
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });
  let actionItems = $state([]);

  let datas = $state([]);
  let productServiceData = $state([]);
  let categoriesData = $state([]);

  let gettingDatas = $state(false);

  let configs = $state({});
  let configsCategories = $state({});

  // function getCategory(category) {
  //   return [
  //     {
  //       type: "brand",
  //       name: category.data.attributes.category_name,
  //     },
  //   ];
  // }

  let tabDatas = $state({
    selected: "Product & Services",
    data: [
      {
        name: "Product & Services",
      },
      {
        name: "Categories",
      },
    ],
  });

  function getStatus(is_active) {
    if (is_active === "Yes") {
      return {
        type: "success",
        name: $_("active"),
      };
    }

    if (is_active === "No") {
      return {
        type: "danger",
        name: $_("inActive"),
      };
    }

    return;
  }

  async function getProducts(page) {
    console.log("page", page);
    console.log("pagination", pagination.page);
    if (page > pagination.pageCount) return;
    let token = localStorage.getItem("token");
    if (token !== null) {
      gettingDatas = true;
      let response = await $call.call("client-dashboard", "get-products", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: page,
          page_size: pagination.pageSize,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        gettingDatas = false;
        return;
      }

      if (response.result.data.length > 0) {
        pagination = response.result.meta.pagination;
        configs = {
          search: true,
          pagination: response.result.meta.pagination,
        };
        console.log("ayam", response.result.data);
        let promises = response.result.data.map(async (product, index) => {
          return {
            id: product.id,
            product_name: product.attributes.product_name,
            price: "$" + product.attributes.price + "/hr",
            // category: getCategory(product.category),
            status: getStatus(product.attributes.is_active),
            actions: {
              objects: [
                check("businessManagement.products_services:view") && {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-product-service-detail")}` +
                    "?productId=" +
                    product.id +
                    "&edit=" +
                    false,
                  name: "View",
                },
                check("businessManagement.products_services:edit") && {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-product-service-detail")}` +
                    "?productId=" +
                    product.id +
                    "&edit=" +
                    true,
                  name: "Edit",
                },
              ].filter(Boolean),
            },
          };
        });

        productServiceData = await Promise.all(promises);
      }

      gettingDatas = false;
    }
  }

  async function getCategories(page) {
    if (page > pagination.pageCount) return;
    let token = localStorage.getItem("token");
    if (token !== null) {
      gettingDatas = true;
      let response = await $call.call("client-dashboard", "get-categories", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: page,
          page_size: pagination.pageSize,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status === 200) {
        pagination = response.result.meta.pagination;
        configsCategories = {
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
                check("businessManagement.products_services:view") && {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-category-detail")}` +
                    "?catergoryId=" +
                    category.id +
                    "&edit=" +
                    false,
                  name: "View",
                },
                check("businessManagement.products_services:edit") && {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-category-detail")}` +
                    "?catergoryId=" +
                    category.id +
                    "&edit=" +
                    true,
                  name: "Edit",
                },
              ].filter(Boolean),
            },
          };
        });

        categoriesData = await Promise.all(promises);
      }

      gettingDatas = false;
    }
  }

  async function addProduct() {
    goto(await $call.getPage("client-dashboard", "crm-product-service-create"));
    return;
  }

  async function addCategory() {
    goto(await $call.getPage("client-dashboard", "crm-category-create"));
    return;
  }

  async function next() {
    getProducts(pagination.page + 1);
  }

  async function previous() {
    getProducts(pagination.page - 1);
  }

  async function nextCategory() {
    getCategories(pagination.page + 1);
  }

  async function previousCategory() {
    getCategories(pagination.page - 1);
  }

  onMount(async () => {
    if (check("businessManagement.products_services:view") === false) {
      window.history.back();
      return;
    }
    getProducts(pagination.page);
    getCategories(pagination.page);
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography
            preset="Display sm/Semibold"
            color="primary"
            text={tabDatas.selected === "Product & Services"
              ? $_("product&Service")
              : $_("manageCategories")}
          />
          {#if tabDatas.selected === "Product & Services"}
            <Button
              disabled={!check("businessManagement.products_services:create")}
              onclick={addProduct}
              icon={PlusIcon}
              ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" /> {$_("addNew")}</Button
            >
          {:else if tabDatas.selected === "Categories"}
            <Button
              disabled={!check("businessManagement.products_services:create")}
              onclick={addCategory}
              icon={PlusIcon}
              ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" />
              {$_("addNewCategory")}</Button
            >
          {/if}
        </div>

        <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

        {#if tabDatas.selected === "Product & Services"}
          <!-- Produc service Table Data -->
          {#if productServiceData.length > 0}
            <DataTable bind:datas={productServiceData} bind:configs bind:next bind:previous />
          {:else if gettingDatas}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else}
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          {/if}
        {:else if tabDatas.selected === "Categories"}
          {#if categoriesData.length > 0}
            <DataTable
              bind:datas={categoriesData}
              bind:configs={configsCategories}
              bind:next={nextCategory}
              bind:previous={previousCategory}
            />
          {:else if gettingDatas}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else}
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          {/if}
        {/if}
      </div>
    </div>
  </div>
</div>
