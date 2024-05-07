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

  const call = new CR3Call(globalThis.location.origin);

  let personsLatest = [];
  let documentsReady = $state(false);
  let pagination = $state({ page: 1, pageSize: 25, pageCount: 1, total: 25 });
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
        name: "Active",
      };
    }

    if (is_active === "No") {
      return {
        type: "danger",
        name: "Inactive",
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
                    product.id,
                  name: "View",
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
    await call.goToPage("client-dashboard", "crm-product-service-create");
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
          <Typography preset="Display sm/Semibold" color="primary" text="Product & Service" />
          <Button onclick={addProduct} icon={PlusIcon}
            >Add new <img class="w-4 h-4" src={PlusIcon} alt="plus-icon" /></Button
          >
        </div>

        <!-- Person Table Data -->
        <div>
          {#if documentsReady}
            <DataTable bind:datas bind:configs bind:next bind:previous />
          {:else}
            <DataTable {datas} {configs} />
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>
