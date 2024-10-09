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
  import { Button } from "flowbite-svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import { goto } from "$app/navigation";
  import { StateCard } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";

  import { call } from "$lib";

  let paginationOrders = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });
  let paginationInvoices = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });

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

  let gettingDatas = $state(false);

  let dataOrders = $state([]);

  let dataInvoices = $state([]);

  let configsOrders = $state({ search: true });

  let configsInvoices = $state({ search: true });

  let orderReady = $state(false);

  let invoiceReady = $state(false);

  async function getOrder() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      try {
        let response = await $call.call("client-dashboard", "get-orders", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: paginationOrders.page,
            page_size: paginationOrders.pageSize,
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
            let promise = response.result.data.map(async (data) => {
              return {
                id: data.id,
                customer_name: data.attributes.person.data?.attributes.fullname,
                order: data.attributes.product.data?.attributes.product_name,
                total_amount: data.attributes.total_amount,
                order_status:
                  data.attributes.order_status[data.attributes.order_status.length - 1].status,
                actions: {
                  objects: [
                    {
                      href:
                        `${await $call.getPage("client-dashboard", "crm-order-invoice-order")}` +
                        "?orderId=" +
                        data.id,
                      name: "View",
                    },
                  ],
                },
              };
            });

            dataOrders = await Promise.all(promise);
          }
          gettingDatas = false;
        }
      } catch (error) {
        console.log(error);
        gettingDatas = false;
      }
    }
  }

  async function getInvoice() {
    let token = localStorage.getItem("token");

    if (token != null) {
      gettingDatas = true;
      let response = await $call.call("client-dashboard", "get-invoices", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: paginationInvoices.page,
          page_size: paginationInvoices.pageSize,
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
          let promise = response.result.data.map(async (data) => {
            return {
              id: data.id,
              customer_name: data.attributes.person.data.attributes.fullname,
              invoice_date: data.attributes.invoice_date,
              due_date: data.attributes.due_date,
              total_amount: data.attributes.order_booking.data?.attributes.total_amount,
              payment_status:
                data.attributes.payment_status[data.attributes.payment_status.length - 1].status,
              actions: {
                objects: [
                  {
                    href:
                      `${await $call.getPage("client-dashboard", "crm-order-invoice-invoice")}` +
                      "?invoiceId=" +
                      data.id,
                    name: "View",
                  },
                ],
              },
            };
          });

          dataInvoices = await Promise.all(promise);

          gettingDatas = false;
        }
      }
    }
  }

  async function addNew() {
    goto(await $call.getPage("client-dashboard", "crm-order-invoice-create"));
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
        <div class="flex items-center justify-between pb-2">
          <Typography
            preset="Display sm/Semibold"
            color="primary"
            text={$_(toCamelCase(tabDatas.selected))}
          />

          <Button onclick={addNew} icon={PlusIcon}
            ><img class="me-1 w-4 h-4" src={PlusIcon} alt="plus-icon" /> {$_("addNewOrder")}</Button
          >
        </div>

        <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

        <!-- Table Data -->
        {#if tabDatas.selected === "Orders"}
          {#if dataOrders.length > 0}
            <DataTable bind:datas={dataOrders} bind:configs={configsOrders} />
          {:else if gettingDatas}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else}
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          {/if}
        {:else if tabDatas.selected === "Invoices"}
          {#if dataInvoices.length > 0}
            <DataTable bind:datas={dataInvoices} bind:configs={configsInvoices} />
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
