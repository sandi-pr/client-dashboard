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
  import { P, Card, Avatar, Heading, A, Button } from "flowbite-svelte";
  import { StateCard } from "$lib";
  import { _ } from "svelte-i18n";
  import UserPlusIcon from "$lib/icon/user-plus-02.svelte";
  import PackagePlusIcon from "$lib/icon/package-plus.svelte";
  import PackageSearchIcon from "$lib/icon/package-search.svelte";
  import Chart from "$lib/components/Chart.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";

  const call = new CR3Call(globalThis.location.origin);

  let { showAI = false } = $props();

  let latestPersons = $state([]);
  let configs = $state({});
  let documentsReady = $state(false);
  let personsIsEmpty = $state(true);
  // let documents = [];
  // let documentsNeed = [];
  // let documentsComplete = [];
  // let allUrls = [];
  // let chart = true;
  let monthlyPersonsCounts = [3, 6, 1, 8, 2, 7, 4, 9, 5, 0, 6, 3];
  let monthlyActivePersonsCounts = [9, 2, 7, 4, 0, 6, 3, 5, 8, 1, 7, 2];
  let monthlyInactivePersonsCounts = [0, 6, 3, 5, 8, 1, 7, 2, 9, 4, 0, 6];
  let dataChartAll = $state({});
  let dataChartNeed = $state({});
  let dataChartComplete = $state({});
  let chartReady = $state(false);
  let actionItems = $state([]);
  let chartAction = $state([]);

  async function getActionItems() {
    actionItems = [
      {
        href: await call.getPage("client-dashboard", "document-management-create"),
        icon: UserPlusIcon,
        actionName: "Add new document",
        desc: "desc Add new document",
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
      {
        name: "Users",
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

    dataChartAll = {
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

    dataChartNeed = {
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

    dataChartComplete = {
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

  async function getChartAction() {
    chartAction = [
      {
        href: await call.getPage("client-dashboard", "doc-list"),
        dataName: "Total Documents",
        value: "10",
        change: "3",
        changeType: "increase",
      },
      {
        href: await call.getPage("client-dashboard", "doc-list"),
        dataName: "General Documents",
        value: "7",
        change: "3",
        changeType: "decrease",
      },
      {
        href: await call.getPage("client-dashboard", "doc-list"),
        dataName: "Contract Documents",
        value: "3",
        change: "3",
        changeType: "increase",
      },
    ];
  }

  let latestDocuments = $state([
    {
      id: 1,
      upload_date: "01/12/2024 04:07",
      document_name: "Work Contract - Hexacode 1",
      type: "Contract",
      actions: {
        objects: [
          {
            href: "document_management/detail",
            name: $_("view"),
          },
        ],
      },
    },
  ]);

  let latestProductServices = $state([]);
  let latestUsers = $state([]);

  async function viewAll() {
    // let href = await call.getPage("client-dashboard", "doc-list");
    let href = "document_management/list";
    globalThis.location.href = href;
  }

  onMount(async () => {
    getChartAction();
    generateChartData();
    getActionItems();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Heading tag="h2" class="font-jakarta text-primary-light dark:text-primary-dark">
          {$_("documentManagementOverview")}
        </Heading>

        <!-- Action Items -->
        <div>
          <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
            <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
              {$_("action")}
            </Heading>
          </div>
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 w-full items-center justify-center">
            {#each actionItems as item}
              <Card
                href={item.href}
                horizontal="true"
                class="cursor-pointer items-center w-full !p-4 min-h-28 bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark border-brand-light dark:border-brand-dark shadow-sm hover:no-underline"
                size="none"
              >
                <Avatar
                  size="sm"
                  border
                  class="!ring-0 bg-brand-secondary-light dark:bg-brand-secondary-dark !p-2 !h-10 !w-10"
                  ><svelte:component
                    this={item.icon}
                    svgClass="w-5 h-5 text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark"
                  /></Avatar
                >
                <div class="ml-4 gap[2px]">
                  <P
                    class="font-jakarta text-base font-semibold leading-normal text-primary-light dark:text-primary-dark"
                    >{$_(toCamelCase(item.actionName))}</P
                  >
                  <P
                    class="font-jakarta text-sm font-normal leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_(toCamelCase(item.desc))}</P
                  >
                </div>
              </Card>
            {/each}
          </div>
        </div>

        <!-- Reports Monitoring -->
        <div>
          <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
            <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
              {$_("reports")}
            </Heading>
          </div>
          <div
            class={`grid ${showAI ? "grid-cols-1 xl:grid-cols-2" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-3"} gap-4 items-center justify-center`}
          >
            {#each chartAction as { href, dataName, value, change, changeType }, index ({ href, dataName, value, change, changeType })}
              <Card
                horizontal="true"
                class="grid grid-cols-12 items-center !px-4 !py-0 min-h-28 justify-between bg-primary-light dark:bg-primary-dark border-secondary-light dark:border-secondary-dark shadow-sm"
                size="none"
              >
                <div class="space-y-2 col-span-6 sm:col-span-7 md:col-span-5">
                  <div class="flex items-center gap-2">
                    <P
                      size="base"
                      weight="semibold"
                      color="text-secondary-light dark:text-secondary-dark"
                      class="font-jakarta">{$_(toCamelCase(dataName))}</P
                    >
                    <P
                      size="sm"
                      weight="bold"
                      color={`${
                        changeType === "decrease"
                          ? "text-error-primary-light dark:text-error-primary-dark"
                          : "text-success-primary-light dark:text-success-primary-dark"
                      }`}
                      class="font-jakarta leading-tight"
                      >{`${changeType === "decrease" ? "-" : "+"}${change}`}</P
                    >
                  </div>
                  <P
                    size="xl"
                    weight="semibold"
                    color="text-gray-900"
                    class="font-jakarta leading-[30px] text-primary-light dark:text-primary-dark"
                    >{value}</P
                  >
                  <A
                    {href}
                    class="text-sm font-semibold text-fg-brand-primary-light dark:text-fg-brand-primary-dark font-jakarta leading-tight"
                    >{$_("seeAll")}</A
                  >
                </div>
                <div class="col-span-6 sm:col-span-5 md:col-span-7">
                  {#if chartReady}
                    {#if index === 0}
                      <Chart dataChart={dataChartAll} />
                    {:else if index === 1}
                      <Chart dataChart={dataChartNeed} />
                    {:else}
                      <Chart dataChart={dataChartComplete} />
                    {/if}
                  {:else}
                    <div></div>
                  {/if}
                </div>
              </Card>
            {/each}
          </div>
        </div>

        <div class="flex flex-col">
          <div class="flex justify-between">
            <Typography preset="Text lg/Semibold" color="primary" text={$_("documentLatest")} />
            <div></div>
            <Button on:click={() => viewAll()} color="light">{$_("seeAll")}</Button>
          </div>

          <div>
            {#if latestDocuments.length > 0}
              <DataTable bind:datas={latestDocuments} bind:configs />
            {:else}
              <div class="pt-6">
                <StateCard
                  illustration={IllustrationEmpty}
                  description={$_("youDontHaveAnyDataYet")}
                />
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
