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
  import DataTable from "$lib/components/DataTable.svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";

  const call = new CR3Call(globalThis.location.origin);

  let paginationOrders = $state({ page: 1, pageSize: 25, pageCount: 1, total: 25 });
  let paginationInvoices = $state({ page: 1, pageSize: 25, pageCount: 1, total: 25 });

  let tabDatas = $state({
    selected: "Orders",
    data: [
      {
        name: "Orders",
      },
      {
        name: "Invoices",
      },
    ],
  });

  let dataOrders = $state([]);

  let dataInvoices = $state([]);

  let configsOrders = $state({ search: true });

  let configsInvoices = $state({ search: true });

  let orderReady = $state(false);

  let invoiceReady = $state(false);

  async function getOrder() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await call.call("client-dashboard", "get-orders", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: paginationOrders.page,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("orders", response);

        if (response.status === 200) {
          if (response.result.data.length > 0) {
            configsOrders = {
              search: true,
              pagination: response.result.meta.pagination,
            };
            let promise = response.result.data.map((data) => {
              return {
                id: data.id,
                booking_date: data.attributes.booking_date,
                customer_name: data.attributes.person.data.attributes.fullname,
                order: data.attributes.product.data.attributes.product_name,
                total_amount: data.attributes.total_amount,
                order_status:
                  data.attributes.order_status[data.attributes.order_status.length - 1].status,
              };
            });

            dataOrders = await Promise.all(promise);

            orderReady = true;
          }
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function getInvoice() {
    let token = sessionStorage.getItem("token");

    if (token != null) {
      let response = await call.call("client-dashboard", "get-invoices", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: paginationInvoices.page,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      console.log("invoices", response);

      if (response.status === 200) {
        if (response.result.data.length > 0) {
          configsInvoices = {
            search: true,
            pagination: response.result.meta.pagination,
          };
          let promise = response.result.data.map((data) => {
            return {
              id: data.id,
              customer_name: data.attributes.person.data.attributes.fullname,
              invoice_date: data.attributes.invoice_date,
              due_date: data.attributes.due_date,
              total_amount: data.attributes.total_amount,
              payment_status:
                data.attributes.payment_status[data.attributes.payment_status.length - 1].status,
            };
          });

          dataInvoices = await Promise.all(promise);

          invoiceReady = true;
        }
      }
    }
  }

  onMount(() => {
    getOrder();
    getInvoice();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <Typography
          preset="Display sm/Semibold"
          color="primary"
          text={$_(toCamelCase(tabDatas.selected))}
        />

        <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

        <!-- Table Data -->
        {#if tabDatas.selected === "Orders"}
          {#if orderReady}
            <div>
              <DataTable bind:datas={dataOrders} bind:configs={configsOrders} />
            </div>
          {:else}
            <div>
              <DataTable bind:datas={dataOrders} bind:configs={configsOrders} />
            </div>
          {/if}
        {:else if tabDatas.selected === "Invoices"}
          {#if invoiceReady}
            <div>
              <DataTable bind:datas={dataInvoices} bind:configs={configsInvoices} />
            </div>
          {:else}
            <div>
              <DataTable bind:datas={dataInvoices} bind:configs={configsInvoices} />
            </div>
          {/if}
        {/if}
      </div>
    </div>
  </div>
</div>
