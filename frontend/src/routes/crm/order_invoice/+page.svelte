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
  import { Heading, Button } from "flowbite-svelte";
  import { StateCard, Typography, ChartCard, ActionCard, check } from "$lib";
  import { _ } from "svelte-i18n";
  import DataTable from "$lib/components/DataTable.svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { goto } from "$app/navigation";

  import { call } from "$lib";

  let { showAI = false } = $props();

  let gettingDatas = $state(false);
  let latestPersons = $state([]);
  let configs = $state({});
  let documentsReady = $state(false);
  let personsIsEmpty = $state(true);
  // let documents = [];
  // let documentsNeed = [];
  // let documentsComplete = [];
  // let allUrls = [];
  // let chart = true;
  let monthlyPersonsCounts = $state([]);
  let monthlyActivePersonsCounts = $state([]);
  let monthlyInactivePersonsCounts = $state([]);
  let dataChartTotalPersons = $state({});
  let dataChartActivePersons = $state({});
  let dataChartInactivePersons = $state({});
  let chartReady = $state(false);
  let actionItems = $state([]);
  let chartAction = $state([]);
  let productServiceData = $state([]);
  let categoriesData = $state([]);
  let configsCategories = $state({});

  let paginationOrders = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });
  let paginationInvoices = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });

  let dataOrders = $state([]);
  let dataInvoices = $state([]);
  let configsOrders = $state({ search: true });
  let configsInvoices = $state({ search: true });

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });

  async function getActionItems() {
    actionItems = [
      {
        disabled: !check("businessManagement.order_invoice:create"),
        href: await $call.getPage("client-dashboard", "crm-order-invoice-create"),
        icon: "packagePlus",
        actionName: "Create new Order Invoice",
        desc: "desc Create new Order Invoice",
      },
    ];
  }

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

  function generateChartData() {
    const currentDate = new Date();
    const twelveMonthsBefore = new Date();
    twelveMonthsBefore.setMonth(currentDate.getMonth() - 12);

    const formatDate = (date) => {
      const options = { month: "short", year: "numeric" };
      return new Intl.DateTimeFormat("en-US", options).format(date);
    };

    dataChartTotalPersons = {
      chart: {
        type: "area",
        fontFamily: "Inter, sans-serif",
        dropShadow: {
          enabled: false,
        },
        toolbar: {
          show: false,
        },
      },
      tooltip: {
        enabled: true,
        x: {
          show: false,
        },
      },
      fill: {
        type: "gradient",
        gradient: {
          opacityFrom: 0.55,
          opacityTo: 0,
          shade: "#17B169",
          gradientToColors: ["#17B169"],
        },
      },
      dataLabels: {
        enabled: false,
      },
      stroke: {
        width: 2,
      },
      grid: {
        show: false,
        strokeDashArray: 4,
        padding: {
          left: 2,
          right: 2,
          top: 0,
        },
      },
      series: [
        {
          name: "New order",
          data: monthlyPersonsCounts,
          color: "#17B169",
        },
      ],
      xaxis: {
        categories: Array.from({ length: 12 }, (_, i) => {
          const date = new Date(twelveMonthsBefore);
          date.setMonth(twelveMonthsBefore.getMonth() + i + 1);
          return formatDate(date);
        }),
        labels: {
          show: false,
        },
        axisBorder: {
          show: false,
        },
        axisTicks: {
          show: false,
        },
      },
      yaxis: {
        show: false,
      },
    };

    chartReady = true;
  }

  async function getChartAction(totalPersons) {
    chartAction = [
      {
        href: await $call.getPage("client-dashboard", "crm-order-invoice"),
        dataName: "Total Order",
        value: totalPersons.total,
        change: totalPersons.change,
        changeType: totalPersons.changeType,
      },
    ];
  }

  async function viewAll() {
    goto(await $call.getPage("client-dashboard", "crm-order-invoice-list"));
  }

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

  async function getCountPersonsPerMonth() {
    const today = new Date();
    let earlyDateCurrentMonth = new Date(today.getFullYear(), today.getMonth(), 1);
    let lastDateCurrentMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0);
    let token = localStorage.getItem("token");

    let tmpMonthlyPersonsCounts = [];

    if (token !== null) {
      for (let i = 0; i < 12; i++) {
        let response = await $call.call("client-dashboard", "get-count-order-per-mount", {
          URLQuery: {},
          URLParams: {},
          Params: {
            earlyDate: earlyDateCurrentMonth.toISOString().split("T")[0],
            lastDate: lastDateCurrentMonth.toISOString().split("T")[0],
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          console.log(response.status);
          console.log(response);
          return;
        }
        // console.log(
        //   "telor",
        //   format(earlyDateCurrentMonth, "yyyy-MM-dd"),
        //   format(lastDateCurrentMonth, "yyyy-MM-dd"),
        // );

        tmpMonthlyPersonsCounts.push(response.result.meta.pagination.total);

        earlyDateCurrentMonth = new Date(
          earlyDateCurrentMonth.getFullYear(),
          earlyDateCurrentMonth.getMonth() - 1,
          1,
        );
        lastDateCurrentMonth = new Date(
          lastDateCurrentMonth.getFullYear(),
          lastDateCurrentMonth.getMonth(),
          0,
        );
      }

      tmpMonthlyPersonsCounts.reverse();

      monthlyPersonsCounts = tmpMonthlyPersonsCounts;

      let totalPersons = groupingDataChart(monthlyPersonsCounts);

      getChartAction(totalPersons);
      generateChartData();
    } else {
      console.log("no token");
    }
  }

  function groupingDataChart(data) {
    let total = 0;
    for (let i = 0; i < data.length; i++) {
      total += data[i];
    }

    let change = data[data.length - 1] - data[data.length - 2];

    let changeType = change > 0 ? "increase" : change < 0 ? "decrease" : "-";

    return {
      total: total,
      change: change,
      changeType: changeType,
    };
  }

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
                    check("businessManagement.order_invoice:view") && {
                      href:
                        `${await $call.getPage("client-dashboard", "crm-order-invoice-order")}` +
                        "?orderId=" +
                        data.id,
                      name: "View",
                    },
                  ].filter(Boolean),
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
                  check("businessManagement.order_invoice:view") && {
                    href:
                      `${await $call.getPage("client-dashboard", "crm-order-invoice-invoice")}` +
                      "?invoiceId=" +
                      data.id,
                    name: "View",
                  },
                ].filter(Boolean),
              },
            };
          });

          dataInvoices = await Promise.all(promise);

          gettingDatas = false;
        }
      }
    }
  }

  onMount(async () => {
    if (check("businessManagement.order_invoice:view") === false) {
      window.history.back();
      return;
    }
    getCountPersonsPerMonth();
    getActionItems();
    getOrder();
    getInvoice();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Typography
          preset="Display sm/Semibold"
          text={$_("orderInvoiceOverview")}
          propClass="fluid-xl"
        />

        <!-- Action Items -->
        <div>
          {#if actionItems.some((item) => item.disabled === false)}
            <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
              <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
                {$_("action")}
              </Heading>
            </div>
          {/if}
          <ActionCard {actionItems} />
        </div>

        <!-- Reports Monitoring -->
        <div>
          <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
            <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
              {$_("reports")}
            </Heading>
          </div>
          <ChartCard {chartAction} chartList={[dataChartTotalPersons, dataChartActivePersons]} />
        </div>

        <div class="flex flex-col">
          <div class="flex justify-between">
            <Typography
              preset="Text lg/Semibold"
              color="primary"
              text={tabDatas.selected === "Orders" ? $_("orders") : $_("invoices")}
            />
            <Button on:click={() => viewAll()} color="light">{$_("seeAll")}</Button>
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
              <StateCard
                illustration={IllustrationEmpty}
                description={$_("youDontHaveAnyDataYet")}
              />
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
              <StateCard
                illustration={IllustrationEmpty}
                description={$_("youDontHaveAnyDataYet")}
              />
            {/if}
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>
