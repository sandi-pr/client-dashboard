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

  let documentsReady = $state(false);
  let pagination = { page: 1, pageSize: 25, pageCount: 1, total: 25 };
  let actionItems = $state([]);

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

  let dataOrders = $state([
    {
      id: 1,
      booking_date: "05/05/2024 20:49",
      customer_name: "Christy Weber",
      order: "Solid Hardwood Flooring",
      total_amount: "$2766",
      order_status: "Calendars Booked",
      actions: {
        objects: [
          {
            href: "/crm/order_invoice/order",
            name: "Detail",
          },
        ],
      },
    }
  ]);

  let dataInvoices = $state([
    {
      id: 1,
      customer_name: "Christy Weber",
      invoice_date: "05/05/2024 20:49",
      due_date: "05/05/2024 20:49",
      total_amount: "$2766",
      payment_status: "Unpaid",
      actions: {
        objects: [
          {
            href: "/crm/order_invoice/invoice",
            name: "Detail",
          },
        ],
      },
    }
  ]);

  let configs = $state({ search: true });

  onMount(async () => {

  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <Typography preset="Display sm/Semibold" color="primary" text={$_(toCamelCase(tabDatas.selected))} />

        <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

        <!-- Table Data -->
        {#if tabDatas.selected === "Orders"}
          <div>
            <DataTable bind:datas={dataOrders} bind:configs />
          </div>
        {:else if tabDatas.selected === "Invoices"}
          <div>
            <DataTable bind:datas={dataInvoices} bind:configs />
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
