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
  import { P, Avatar, Heading, A, Button } from "flowbite-svelte";
  import { StateCard, Typography, Card, ChartCard, ActionCard } from "$lib";
  import { _ } from "svelte-i18n";
  import UserPlusIcon from "$lib/icon/user-plus-02.svelte";
  import PackagePlusIcon from "$lib/icon/package-plus.svelte";
  import PackageSearchIcon from "$lib/icon/package-search.svelte";
  import Chart from "$lib/components/Chart.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { CR3Call } from "hexacomponents-svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { format } from "date-fns";

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

  async function getActionItems() {
    actionItems = [
      {
        href: await $call.getPage("client-dashboard", "crm-person-create"),
        icon: "userPlus",
        actionName: "Create new persons",
        desc: "desc Create New Person",
      },
      {
        href: await $call.getPage("client-dashboard", "crm-product-service-create"),
        icon: "packagePlus",
        actionName: "Create new Product & Services",
        desc: "desc Create New Product & Services",
      },
      {
        href: await $call.getPage("client-dashboard", "crm-categories"),
        icon: "packageSearch",
        actionName: "Manage Categories",
        desc: "desc Manage Categories",
      },
    ];
  }

  let tabDatas = $state({
    selected: "Persons",
    data: [
      {
        name: "Persons",
      },
      {
        name: "Products & Services",
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
        href: await $call.getPage("client-dashboard", "crm-person"),
        dataName: "Total Persons",
        value: totalPersons.total,
        change: totalPersons.change,
        changeType: totalPersons.changeType,
      },
      {
        href: await $call.getPage("client-dashboard", "crm-person"),
        dataName: "Active Persons",
        value: activePersons.total,
        change: activePersons.change,
        changeType: activePersons.changeType,
      },
      {
        href: await $call.getPage("client-dashboard", "crm-person"),
        dataName: "Inactive Persons",
        value: inactivePersons.total,
        change: inactivePersons.change,
        changeType: inactivePersons.changeType,
      },
    ];
  }

  let latestProductServices = $state([]);

  let latestUsers = $state([]);

  async function viewAll(title) {
    if (title === "Persons") {
      let href = await $call.getPage("client-dashboard", "crm-person");
      globalThis.location.href = href;
      return;
    } else if (title === "Products & Services") {
      let href = await $call.getPage("client-dashboard", "crm-product-service");
      globalThis.location.href = href;
      return;
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

  async function getPerson() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      gettingDatas = true;
      let response = await $call.call("client-dashboard", "get-persons", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: 1,
          page_size: 25,
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
        response.result.data.sort(
          (a, b) => new Date(b.attributes.updatedAt) - new Date(a.attributes.updatedAt),
        );
        let latestData = response.result.data.slice(0, 5);

        let promises = latestData.map(async (person) => {
          return {
            id: person.id,
            fullname: person.attributes.fullname,
            email: person.attributes.email,
            phone: person.attributes.phone,
            // address: person.attributes.address,
            // city: person.attributes.city,
            // state: person.attributes.state,
            // zip_code: person.attributes.zip_code,
            country: person.attributes.country,
            // date_of_birth: person.attributes.date_of_birth,
            status: getStatus(person.attributes.is_active),
            actions: {
              objects: [
                {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-person-detail")}` +
                    "?personId=" +
                    person.id +
                    "&edit=false",
                  name: "View",
                },
                {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-person-detail")}` +
                    "?personId=" +
                    person.id +
                    "&edit=true",
                  name: "Edit",
                },
              ],
            },
          };
        });

        latestPersons = await Promise.all(promises);
      }
      
      gettingDatas = false;
    }
  }

  async function getProducts() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      gettingDatas = true;
      let response = await $call.call("client-dashboard", "get-products", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: 1,
          page_size: 5,
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
        let latestData = response.result.data.slice(0, 5);
        let promises = latestData.map(async (product, index) => {
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
                    "&edit=false",
                  name: "View",
                },
                {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-product-service-detail")}` +
                    "?productId=" +
                    product.id +
                    "&edit=true",
                  name: "Edit",
                },
              ],
            },
          };
        });

        latestProductServices = await Promise.all(promises);
        gettingDatas = false;
      }
    }
  }

  async function getCountPersonsPerMonth() {
    const today = new Date();
    let earlyDateCurrentMonth = new Date(today.getFullYear(), today.getMonth(), 1);
    let lastDateCurrentMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0);
    let token = sessionStorage.getItem("token");

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

  onMount(async () => {
    // getChartAction(
    //   { total: 0, change: 0, changeType: "-" },
    //   { total: 0, change: 0, changeType: "-" },
    //   { total: 0, change: 0, changeType: "-" },
    // );
    getCountPersonsPerMonth();
    // generateChartData();
    getActionItems();
    getPerson();
    getProducts();
    // getUsers();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Typography preset="Display sm/Semibold" text={$_("crmOverview")} propClass="fluid-xl" />

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
          <ChartCard {chartAction} chartList={[dataChartTotalPersons, dataChartActivePersons, dataChartInactivePersons,]} />
        </div>

        <div class="flex flex-col">
          <div class="flex justify-between">
            <Typography
              preset="Text lg/Semibold"
              color="primary"
              text={$_(toCamelCase(tabDatas.selected))}
            />
            <div></div>
            <!-- <Button on:click={() => viewAll(tabDatas.selected)} color="light">{$_("seeAll")}</Button> -->
          </div>
          <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

          <div>
            {#if tabDatas.selected === "Persons"}
              <!-- Person Table Data -->
              {#if latestPersons.length > 0}
                <DataTable bind:datas={latestPersons} bind:configs />
              {:else if gettingDatas}
                <div class="pt-4 lg:pt-6">
                  <StateCard
                    illustration={IllustrationProgress}
                    description={$_("fetchingDataPleaseWait")}
                  />
                </div>
              {:else}
                <div class="pt-4 lg:pt-6">
                  <StateCard
                    illustration={IllustrationEmpty}
                    description={$_("youDontHaveAnyDataYet")}
                  />
                </div>
              {/if}
            {:else if tabDatas.selected === "Products & Services"}
              {#if latestProductServices.length > 0}
                <DataTable bind:datas={latestProductServices} bind:configs />
              {:else if gettingDatas}
                <div class="pt-4 lg:pt-6">
                  <StateCard
                    illustration={IllustrationProgress}
                    description={$_("fetchingDataPleaseWait")}
                  />
                </div>
              {:else}
                <div class="pt-4 lg:pt-6">
                  <StateCard
                    illustration={IllustrationEmpty}
                    description={$_("youDontHaveAnyDataYet")}
                  />
                </div>
              {/if}
              <!-- {:else if tabDatas.selected === "Users"}
              {#if latestUsers.length > 0}
                <DataTable bind:datas={latestUsers} bind:configs />
              {:else}
                <div class="pt-6">
                  <StateCard
                    illustration={IllustrationEmpty}
                    description={$_("youDontHaveAnyDataYet")}
                  />
                </div>
              {/if} -->
            {/if}
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
