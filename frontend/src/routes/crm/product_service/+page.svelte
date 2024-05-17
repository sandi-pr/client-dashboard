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
  import { StateCard } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { goto } from "$app/navigation";

  const call = new CR3Call(globalThis.location.origin);

  let personsLatest = [];
  let documentsReady = $state(false);
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });
  let actionItems = $state([]);

  let datas = $state([]);

  let configs = $state({});

  // function getCategory(category) {
  //   return [
  //     {
  //       type: "brand",
  //       name: category.data.attributes.category_name,
  //     },
  //   ];
  // }

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
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-products", {
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
                {
                  href:
                    `${await call.getPage("client-dashboard", "crm-product-service-detail")}` +
                    "?productId=" +
                    product.id +
                    "&edit=" +
                    false,
                  name: "View",
                },
                {
                  href:
                    `${await call.getPage("client-dashboard", "crm-product-service-detail")}` +
                    "?productId=" +
                    product.id +
                    "&edit=" +
                    true,
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

  async function addProduct() {
    goto(await call.getPage("client-dashboard", "crm-product-service-create"));
    return;
  }

  async function next() {
    getProducts(pagination.page + 1);
  }

  async function previous() {
    getProducts(pagination.page - 1);
  }

  onMount(async () => {
    getProducts(pagination.page);
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text={$_("product&Service")} />
          <Button onclick={addProduct} icon={PlusIcon}
            ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" /> {$_("addNew")}</Button
          >
        </div>

        <!-- Person Table Data -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs bind:next bind:previous />
        {:else}
          <div class="pt-8">
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
