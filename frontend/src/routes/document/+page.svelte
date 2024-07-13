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
  import { onMount, onDestroy } from "svelte";
  import { Modal, Button, Input } from "flowbite-svelte";
  import {
    StateCard,
    Typography,
    SectionHeader,
    Card,
    ChartCard,
    CollectionCard,
    triggerToast,
  } from "$lib";
  import { _ } from "svelte-i18n";
  import FolderPlusIcon from "$lib/icon/folder-plus.svelte";
  import FilePlusIcon from "$lib/icon/file-plus-02.svelte";
  import Chart from "$lib/components/Chart.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { CR3Call } from "hexacomponents-svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { format } from "date-fns";
  import { goto } from "$app/navigation";

  const call = new CR3Call(globalThis.location.origin);

  let { showAI = false } = $props();

  let documents = $state([]);
  let collections = $state([]);
  let configs = $state({});
  let documentsReady = $state(false);
  let personsIsEmpty = $state(true);
  // let documentsNeed = [];
  // let documentsComplete = [];
  // let allUrls = [];
  // let chart = true;
  let monthlyCollectionsCounts = $state([]);
  let monthlyCollectionReadyCounts = $state([]);
  let monthlyCollectingDataCounts = $state([]);
  let dataChartAll = $state({});
  let dataChartReady = $state({});
  let dataChartCollecting = $state({});
  let chartReady = $state(false);
  let actionItems = $state([]);
  let chartAction = $state([]);
  let tmpCollectionName = $state("");
  let tmpDescription = $state("");
  let addCollectionsModal = $state(false);
  let timer;

  async function getActionItems() {
    actionItems = [
      {
        href: await call.getPage("client-dashboard", "document-collection-create"),
        icon: "filePlus",
        actionName: "Add new document collection",
        desc: "desc Add new document collection",
      },
    ];
  }

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
          data: monthlyCollectionsCounts,
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

    dataChartReady = {
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
          data: monthlyCollectionReadyCounts,
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

    dataChartCollecting = {
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
          data: monthlyCollectingDataCounts,
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

  async function getChartAction(totalCollections, readyCollection, dataCollecting) {
    chartAction = [
      {
        href: await call.getPage("client-dashboard", "document-collection"),
        dataName: "Total Documents",
        value: totalCollections.total,
        change: totalCollections.change,
        changeType: totalCollections.changeType,
      },
      {
        href: await call.getPage("client-dashboard", "document-collection"),
        dataName: "Document Collection",
        value: readyCollection.total,
        change: readyCollection.change,
        changeType: readyCollection.changeType,
      },
      {
        href: await call.getPage("client-dashboard", "document-collection"),
        dataName: "Collecting Data",
        value: dataCollecting.total,
        change: dataCollecting.change,
        changeType: dataCollecting.changeType,
      },
    ];
  }

  let latestDocuments = $state([]);

  let latestDocumentCollection = $state([]);

  async function viewAll() {
    goto(await call.getPage("client-dashboard", "document-collection"));
  }

  async function createCollection() {
    console.log(files);
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      let response = await call.call("client-dashboard", "create-document-collection", {
        URLQuery: {},
        URLParams: {},
        Params: {
          collection_name: tmpCollectionName,
          description: tmpDescription,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        triggerToast("error");
        return;
      }

      console.log(response);

      let response2 = await call.call("hexa-ai-rag", "add-document-by-collection-id", {
        URLQuery: {},
        URLParams: {},
        Params: {
          collection_id: response.result.data.id.toString(),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response2.status !== 200) {
        console.log(response2);
        triggerToast("error");
        return;
      }

      console.log(response2);

      tmpCollectionName = "";
      tmpDescription = "";

      triggerToast("success");

      addCollectionsModal = false;
      getCollections();
      // goto(await call.getPage("client-dashboard", "document-collection"));
    }
  }

  async function getCollections() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      collections = [];
      documents = [];
      let response = await call.call("client-dashboard", "get-document-collections", {
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
        return;
      }

      if (response.result.data.length > 0) {
        collections = collections.concat(response.result.data);
      }

      // let page2 = 1;
      // while (true) {
      //   let response2 = await call.call("client-dashboard", "get-documents", {
      //     URLQuery: {},
      //     URLParams: {},
      //     Params: {
      //       page: page2,
      //       page_size: 25,
      //     },
      //     Headers: {
      //       Authorization: "Bearer " + token,
      //     },
      //   });

      //   if (response2.status !== 200) {
      //     console.log(response2);
      //     return;
      //   }

      //   if (response2.result.data.length > 0) {
      //     documents = documents.concat(response2.result.data);
      //     page2++;
      //   } else {
      //     break;
      //   }
      // }

      for (const collection of collections) {
        let page = 1;
        while (true) {
          let response3 = await call.call("client-dashboard", "get-documents-by-collectionid", {
            URLQuery: {},
            URLParams: {},
            Params: {
              collection_id: collection.id.toString(),
              page: page,
              page_size: 25,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response3.status !== 200) {
            console.log(response3);
            return;
          }

          if (response3.result.data.length > 0) {
            documents = documents.concat(response3.result.data);
            page++;
          } else {
            break;
          }
        }
      }

      let collectionList = await call.getPage("client-dashboard", "document-collection-list");

      latestDocumentCollection = collections.map((collection) => {
        let totalDocuments = documents.filter(
          (document) => document.attributes.collection_id === collection.id,
        ).length;
        return {
          collectionId: collection.id,
          activationTime: format(collection.attributes.createdAt, "dd-MM-yyyy HH:mm"),
          collectionName: collection.attributes.collection_name,
          totalDocuments: totalDocuments,
          href: `${collectionList}?collectionId=${collection.id}`,
        };
      });
    }
  }

  async function getCountDocumentsPerMonth() {
    const today = new Date();
    let earlyDateCurrentMonth = new Date(today.getFullYear(), today.getMonth(), 1);
    let lastDateCurrentMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0);
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      for (let i = 0; i < 12; i++) {
        let response = await call.call("client-dashboard", "get-count-documents-per-month", {
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

        monthlyCollectionsCounts.push(response.result.document.meta.pagination.total);
        // TODO: adjust response to collection ready
        monthlyCollectionReadyCounts.push(response.result.general.meta.pagination.total);
        // TODO: adjust response to collecting data
        monthlyCollectingDataCounts.push(response.result.contract.meta.pagination.total);

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

      monthlyCollectionsCounts.reverse();
      monthlyCollectionReadyCounts.reverse();
      monthlyCollectingDataCounts.reverse();

      let totalCollections = groupingDataChart(monthlyCollectionsCounts);
      let readyCollection = groupingDataChart(monthlyCollectionReadyCounts);
      let dataCollecting = groupingDataChart(monthlyCollectingDataCounts);

      getChartAction(totalCollections, readyCollection, dataCollecting);
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
    getChartAction(
      { total: 0, change: 0, changeType: "-" },
      { total: 0, change: 0, changeType: "-" },
      { total: 0, change: 0, changeType: "-" },
    );
    getCollections();
    getCountDocumentsPerMonth();
    generateChartData();
    getActionItems();
  });

  onDestroy(() => {
    if (timer) {
      clearTimeout(timer);
    }
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Typography preset="Display sm/Semibold" text={$_("documentCollectionOverview")} />

        <!-- Action Items -->
        <div class="flex flex-col gap-4">
          <SectionHeader title={$_("action")} />
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 w-full items-center justify-center">
            {#each actionItems as item}
              <Card
                onclick={() => (addCollectionsModal = true)}
                icon={item.icon}
                title={$_(toCamelCase(item.actionName))}
                description={$_(toCamelCase(item.desc))}
              />
            {/each}
          </div>
        </div>

        <!-- Reports Monitoring -->
        <div class="flex flex-col gap-4">
          <SectionHeader title={$_("reports")} />
          <div
            class={`grid ${showAI ? "grid-cols-1 xl:grid-cols-2" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-3"} gap-4 items-center justify-center`}
          >
            {#each chartAction as { href, dataName, value, change, changeType }, index}
              <ChartCard
                {href}
                dataName={$_(toCamelCase(dataName))}
                {value}
                {change}
                {changeType}
                {index}
                chartList={[dataChartAll, dataChartReady, dataChartCollecting]}
              />
            {/each}
          </div>
        </div>

        <div class="flex flex-col gap-4">
          <div
            class="flex justify-between border-b border-secondary-light dark:border-secondary-dark pb-4"
          >
            <Typography preset="Text lg/Semibold" color="primary" text={$_("documentLatest")} />
            <div></div>
            <Button on:click={() => viewAll()} color="light">{$_("seeAll")}</Button>
          </div>

          <div>
            {#if latestDocumentCollection.length > 0}
              <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
                {#each latestDocumentCollection as document}
                  <CollectionCard
                    collectionName={document.collectionName}
                    totalDocuments={document.totalDocuments}
                    status={document.status}
                  />
                {/each}
              </div>
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

<Modal title={$_("createNewCollection")} bind:open={addCollectionsModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography
        preset="Text sm/medium"
        color="seaddCollectionsModal
      addCollectionsModalcondary"
        text={`${$_("collectionName")}*`}
      />
      <Input bind:value={tmpCollectionName} placeholder={$_("whatIsYourCollectionTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={createCollection} disabled={!tmpCollectionName} class="w-full"
      >{$_("addNewCollection")}</Button
    >
  </svelte:fragment>
</Modal>
