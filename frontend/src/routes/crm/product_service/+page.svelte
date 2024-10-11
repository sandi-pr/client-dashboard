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
  import { StateCard, Typography, ChartCard, ActionCard } from "$lib";
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

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });

  async function getActionItems() {
    actionItems = [
      {
        href: await $call.getPage("client-dashboard", "crm-product-service-create"),
        icon: "packagePlus",
        actionName: "Create new Product & Services",
        desc: "desc Create New Product & Services",
      },
      {
        href: await $call.getPage("client-dashboard", "crm-category-create"),
        icon: "packageSearch",
        actionName: "Add new Category",
        desc: "desc Add new Category",
      },
    ];
  }

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
          name: "New users",
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

    dataChartActivePersons = {
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
          name: "New users",
          data: monthlyActivePersonsCounts,
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

    dataChartInactivePersons = {
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
          name: "New users",
          data: monthlyInactivePersonsCounts,
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

  async function getChartAction(totalPersons, activePersons, inactivePersons) {
    chartAction = [
      {
        href: await $call.getPage("client-dashboard", "crm-person-list"),
        dataName: "Total Persons",
        value: totalPersons.total,
        change: totalPersons.change,
        changeType: totalPersons.changeType,
      },
      {
        href: await $call.getPage("client-dashboard", "crm-person-list"),
        dataName: "Active Persons",
        value: activePersons.total,
        change: activePersons.change,
        changeType: activePersons.changeType,
      },
    ];
  }

  async function viewAll() {
    if (tabDatas.selected === "Product & Services") {
      goto(await $call.getPage("client-dashboard", "crm-product-service-list"));
    } else if (tabDatas.selected === "Categories") {
      goto(await $call.getPage("client-dashboard", "crm-categories"));
    }
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
    let tmpMonthlyActivePersonsCounts = [];
    let tmpMonthlyInactivePersonsCounts = [];

    if (token !== null) {
      for (let i = 0; i < 12; i++) {
        let response = await $call.call("client-dashboard", "get-count-persons-per-month", {
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

        tmpMonthlyPersonsCounts.push(response.result.persons.meta.pagination.total);
        tmpMonthlyActivePersonsCounts.push(response.result.activePersons.meta.pagination.total);
        tmpMonthlyInactivePersonsCounts.push(response.result.inactivePersons.meta.pagination.total);

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
      tmpMonthlyActivePersonsCounts.reverse();
      tmpMonthlyInactivePersonsCounts.reverse();

      monthlyPersonsCounts = tmpMonthlyPersonsCounts;
      monthlyActivePersonsCounts = tmpMonthlyActivePersonsCounts;
      monthlyInactivePersonsCounts = tmpMonthlyInactivePersonsCounts;

      let totalPersons = groupingDataChart(monthlyPersonsCounts);
      let activePersons = groupingDataChart(monthlyActivePersonsCounts);
      let inactivePersons = groupingDataChart(monthlyInactivePersonsCounts);

      getChartAction(totalPersons, activePersons, inactivePersons);
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
                {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-product-service-detail")}` +
                    "?productId=" +
                    product.id +
                    "&edit=" +
                    false,
                  name: "View",
                },
                {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-product-service-detail")}` +
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
                {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-category-detail")}` +
                    "?catergoryId=" +
                    category.id +
                    "&edit=" +
                    false,
                  name: "View",
                },
                {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-category-detail")}` +
                    "?catergoryId=" +
                    category.id +
                    "&edit=" +
                    true,
                  name: "Edit",
                },
              ],
            },
          };
        });

        categoriesData = await Promise.all(promises);
      }

      gettingDatas = false;
    }
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
    getCountPersonsPerMonth();
    getActionItems();
    getProducts(pagination.page);
    getCategories(pagination.page);
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Typography preset="Display sm/Semibold" text={$_("productServiceOverview")} propClass="fluid-xl" />

        <!-- Action Items -->
        <div>
          <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
            <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
              {$_("action")}
            </Heading>
          </div>
          <ActionCard {actionItems} />
        </div>

        <!-- Reports Monitoring -->
        <div>
          <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
            <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
              {$_("reports")}
            </Heading>
          </div>
          <ChartCard
            {chartAction}
            chartList={[dataChartTotalPersons, dataChartActivePersons]}
          />
        </div>

        <div class="flex flex-col">
          <div class="flex justify-between">
            <Typography preset="Text lg/Semibold" color="primary" text={tabDatas.selected === "Product & Services" ? $_("product&Service") : $_("categories")} />
            <Button on:click={() => viewAll()} color="light">{$_("seeAll")}</Button>
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
              <DataTable bind:datas={categoriesData} bind:configs={configsCategories} bind:next={nextCategory} bind:previous={previousCategory} />
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
</div>
