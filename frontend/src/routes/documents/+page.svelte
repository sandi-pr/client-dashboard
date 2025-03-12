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
  import { Modal, Button, Input, P, Indicator, Alert, Select, Helper } from "flowbite-svelte";
  import {
    StateCard,
    Typography,
    SectionHeader,
    Card,
    ChartCard,
    CollectionCard,
    LibraryCard,
    triggerToast,
    UploadFile,
    ActionCard,
    check,
    isBakrie,
  } from "$lib";
  import { _ } from "svelte-i18n";
  import PlusIcon from "$lib/icon/plus.svg";
  import FileTypeIcon from "$lib/icon/file-type-icon.svelte";
  import {
    ExclamationCircleSolid,
    ChevronDownOutline,
    InfoCircleSolid,
  } from "flowbite-svelte-icons";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { CR3Call } from "hexacomponents-svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { format } from "date-fns";
  import { goto } from "$app/navigation";
  import {
    triggerSelectCollection,
    triggerShowAI,
    triggerSelectDropdown,
    triggerReload,
    triggerLoading,
    showAI,
    showNotif,
  } from "$lib";
  import LoadingContainer from "$lib/components/LoadingContainer.svelte";

  import { call } from "$lib";
  import { on } from "svelte/events";

  // let { showAI = true } = $props();
  // let showAI = $state(false);

  let documents = $state([]);
  let collections = $state([]);
  let files = $state([]);
  let monthlyCollectionsCounts = $state([]);
  let monthlyCollectionReadyCounts = $state([]);
  let monthlyCollectingDataCounts = $state([]);
  let tmpMonthlyCollectionsCounts = $state([]);
  let tmpMonthlyCollectionReadyCounts = $state([]);
  let tmpMonthlyCollectingDataCounts = $state([]);
  let dataChartAllDocs = $state({});
  let dataChartCollections = $state({});
  let dataChartCollecting = $state({});
  let chartReady = $state(false);
  let actionItems = $state([]);
  let chartAction = $state([]);
  let tmpCollectionName = $state("");
  let tmpOldName = $state("");
  let tmpDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpHref = $state("");
  let categories = $state([]);
  let otherCategory = $state("");
  let addOtherCategory = $state(false);
  let tmpCategories = $state("");
  let tmpDocumentType = $state("");
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let addLibraryModal = $state(false);
  let addDocumentsModal = $state(false);
  let editModal = $state(false);
  let deleteModal = $state(false);
  let dropdownOpen = $state(false);
  let acceptedFileTypes = $state(".txt, .pdf");
  let selectedTypeCollection = $state("Documents (TXT, PDF)");
  let timer;
  let fileNotSupport = $state(false);
  let attributeNotUnique = $state(false);
  let nameIsAvailable = $state(true);

  async function getActionItems() {
    actionItems = [
      {
        onclick: () => {
          addLibraryModal = true;
          attributeNotUnique = false;
          tmpCollectionName = "";
          tmpDescription = "";
        },
        // href: await $call.getPage("client-dashboard", "document-collection-create"),
        icon: "filePlus",
        actionName: "Add new library collection",
        desc: "desc Add new library collection",
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

    dataChartAllDocs = {
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
          name: "Total library",
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

    dataChartCollections = {
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
          name: "Documents collection",
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
          name: "Total document",
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

  async function getChartAction(totalCollections, dataCollecting) {
    chartAction = [
      // {
      //   href: await $call.getPage("client-dashboard", "document-library"),
      //   dataName: "Total Library",
      //   value: totalCollections.total,
      //   change: totalCollections.change,
      //   changeType: totalCollections.changeType,
      // },
      {
        href: await $call.getPage("client-dashboard", "document-collection"),
        dataName: "Total Document",
        value: dataCollecting.total,
        change: dataCollecting.change,
        changeType: dataCollecting.changeType,
      },
    ];
  }

  let latestDocuments = $state([]);

  let latestDocumentCollection = $state([]);

  let gettingDatas = $state(false);

  async function viewAll() {
    goto(await $call.getPage("client-dashboard", "document-library"));
  }

  async function createCollection() {
    console.log(files);
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "create-library", {
        URLQuery: {},
        URLParams: {},
        Params: {
          library_name: tmpCollectionName,
          description: tmpDescription,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        if (response && response.errors) {
          const errors = response.errors.flatMap((error) =>
            error.errors.flatMap((subError) => subError.message),
          );

          if (errors.includes("This attribute must be unique")) {
            attributeNotUnique = true;
          } else {
            console.log("Pesan error :", errors);
            triggerToast("error");
          }
        }
        return;
      }

      console.log(response);

      // let response2 = await $call.call("hexa-ai-rag", "add-document-by-collection-id", {
      //   URLQuery: {},
      //   URLParams: {},
      //   Params: {
      //     collection_id: response.result.data.id.toString(),
      //   },
      //   Headers: {
      //     Authorization: "Bearer " + token,
      //   },
      // });

      // if (response2.status !== 200) {
      //   console.log(response2);
      //   triggerToast("error");
      //   return;
      // }

      // console.log(response2);

      tmpCollectionName = "";
      tmpDescription = "";
      attributeNotUnique = false;

      triggerToast("created");

      addLibraryModal = false;
      getCollections();
      getCountDocumentsPerMonth();
      // goto(await $call.getPage("client-dashboard", "document-collection"));
    }
  }

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      collections = [];
      documents = [];
      let response = await $call.call("client-dashboard", "get-librarys", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: 1,
          page_size: 4,
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
        collections = collections.concat(response.result.data);
      }

      // let page2 = 1;
      // while (true) {
      //   let response2 = await $call.call("client-dashboard", "get-documents", {
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

      let collectionList = await $call.getPage("client-dashboard", "document-library-collections");

      let promise = collections.map(async (collection) => {
        console.log("collection", collection);
        let response4 = await $call.call("client-dashboard", "get-status-library", {
          URLQuery: {},
          URLParams: {
            id: collection.id.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response4.status !== 200) {
          console.log(response4);
          gettingDatas = false;
          return;
        }

        console.log("status", response4);
        let collectionStatus = response4.result.status;

        return {
          collectionId: collection.id,
          activationTime: format(collection.attributes.createdAt, "dd-MM-yyyy HH:mm"),
          collectionName: collection.attributes.name,
          description: collection.attributes?.meta?.description || "",
          // totalDocuments: totalDocuments,
          // categories: categoriesCollection,
          status: collectionStatus,
          href: `${collectionList}?libraryId=${collection.id}`,
          // webchat: collection.attributes.webchat_tmp === true ? true : false,
        };
      });

      latestDocumentCollection = await Promise.all(promise);

      gettingDatas = false;
    }
  }

  async function editCollectionName() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "update-library", {
          URLQuery: {},
          URLParams: {
            id: tmpCollectionId.toString(),
          },
          Params: {
            library_name: tmpCollectionName,
            description: tmpDescription,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          if (response && response.errors) {
            const errors = response.errors.flatMap((error) =>
              error.errors.flatMap((subError) => subError.message),
            );

            if (errors.includes("This attribute must be unique")) {
              attributeNotUnique = true;
            } else {
              console.log("Pesan error :", errors);
              triggerToast("error");
            }
          }
          return;
        }

        tmpCollectionName = "";
        tmpDescription = "";
        attributeNotUnique = false;

        editModal = false;
        triggerToast("updated");
        getCollections();
      } catch (error) {
        console.log(error);
        triggerToast("error");
      }
    } else {
      console.log("Token not found");
    }
  }

  async function deleteDocumentCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response2 = await $call.call("client-dashboard", "delete-library", {
          URLQuery: {},
          URLParams: {
            id: tmpCollectionId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          console.log(response2);
          triggerToast("error");
          return;
        }

        console.log("ayam", response2);

        tmpCollectionId = "";
        triggerToast("deleted");

        getCollections();
        deleteModal = false;
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    } else {
      console.log("Token not found");
    }
  }

  async function getCountDocumentsPerMonth() {
    tmpMonthlyCollectionsCounts = [];

    const today = new Date();
    let earlyDateCurrentMonth = new Date(today.getFullYear(), today.getMonth(), 1);
    let lastDateCurrentMonth = new Date(today.getFullYear(), today.getMonth() + 1, 1);
    let token = localStorage.getItem("token");

    if (token !== null) {
      for (let i = 0; i < 12; i++) {
        let response = await $call.call("client-dashboard", "get-count-libraries-per-month", {
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

        tmpMonthlyCollectionsCounts.push(response.result.library.meta.pagination.total);
        tmpMonthlyCollectingDataCounts.push(response.result.document.meta.pagination.total);

        earlyDateCurrentMonth = new Date(
          earlyDateCurrentMonth.getFullYear(),
          earlyDateCurrentMonth.getMonth() - 1,
          1,
        );
        lastDateCurrentMonth = new Date(
          lastDateCurrentMonth.getFullYear(),
          lastDateCurrentMonth.getMonth() - 1,
          1,
        );
      }

      tmpMonthlyCollectionsCounts.reverse();
      tmpMonthlyCollectingDataCounts.reverse();

      monthlyCollectionsCounts = tmpMonthlyCollectionsCounts;
      monthlyCollectingDataCounts = tmpMonthlyCollectingDataCounts;

      let totalCollections = groupingDataChart(monthlyCollectionsCounts);
      let totalCollectingData = groupingDataChart(monthlyCollectingDataCounts);

      console.log("totalCollectingData", totalCollectingData);
      console.log("totalCollections", monthlyCollectingDataCounts);

      getChartAction(totalCollections, totalCollectingData);
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

    let change = data[data.length - 1];
    // let change = data[data.length - 1] - data[data.length - 2];

    let changeType = change > 0 ? "increase" : change < 0 ? "decrease" : "-";

    return {
      total: total,
      change: change,
      changeType: changeType,
    };
  }

  async function selectCollection(collection) {
    console.log("ayam overview");
    triggerShowAI(true);
    localStorage.setItem("showAI", "true");
    triggerSelectDropdown("collections");
    localStorage.setItem("selectedDropdown", "collections");
    triggerReload(true);
    triggerSelectCollection(collection);
    localStorage.setItem("selectedCollection", JSON.stringify(collection));
    return;
  }

  const checkLibraryName = async (name, parentId, type) => {
    let token = localStorage.getItem("token");
    if (!name || !token) {
      return null;
    }

    try {
      let resp = await $call.call("client-dashboard", "is-library-name-available", {
        URLQuery: {},
        URLParams: {},
        Params: {
          name: name,
          parent_id: parseInt(parentId),
          type: type,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      console.log("cek resp", resp);

      if (resp === null) {
        throw new Error("Server sibuk");
      }

      if (resp.status === 200) {
        return resp.result.available; // Mengembalikan status (true: tersedia, false: tidak tersedia)
      } else {
        console.error("Error memeriksa name:", resp);
        return null;
      }
    } catch (error) {
      console.error("Terjadi kesalahan:", error);
      return null;
    }
  };

  const handleInput = async (event) => {
    const name = event.target.value;
    const status = await checkLibraryName(name.trim(), 0, "library");

    if (tmpOldName === name) {
      nameIsAvailable = true;
    } else if (status === false) {
      nameIsAvailable = false;
    } else if (status === true) {
      nameIsAvailable = true;
    }
  };

  onMount(async () => {
    if (
      $isBakrie &&
      typeof window !== "undefined" &&
      window.document !== null &&
      window.localStorage.getItem("libraryIdWebchatBakrie")
    ) {
      const libraryIdWebchatBakrie = window.localStorage.getItem("libraryIdWebchatBakrie");
      const redirectUrl = `documents/libraries/collections/?libraryId=${libraryIdWebchatBakrie}`;
      goto(redirectUrl);
    }
    if (check("documentAi.documentLibrary:view") === false) {
      window.history.back();
      return;
    }
    getActionItems();
    getCountDocumentsPerMonth();
    getCollections();
    // generateChartData();
  });

  onDestroy(() => {
    if (timer) {
      clearTimeout(timer);
    }
  });

  let buttonDisabled = $state(false);

  $effect(() => {
    buttonDisabled = files.length === 0 || tmpDocumentTitle === "";
  });
</script>

{#if !$isBakrie}
  <div class="flex dark:bg-primary-dark">
    <div class="flex flex-col w-full pb-20 m-4 gap-3 md:gap-8">
      <!-- User Information -->
      <Typography
        preset="Display sm/Semibold"
        text={$_("documentLibraryOverview")}
        propClass="fluid-xl"
      />

      {#if check("documentAi.documentLibrary:create")}
        <!-- Action Items -->
        <div class="flex flex-col gap-4">
          <SectionHeader title={$_("action")} />
          <ActionCard {actionItems} />
        </div>
      {/if}

      <!-- Reports Monitoring -->
      <div class="flex flex-col gap-4">
        <SectionHeader title={$_("reports")} />
        <ChartCard {chartAction} chartList={[dataChartCollecting]} />
      </div>

      <div class="flex flex-col gap-4">
        <div
          class="flex items-center justify-between border-b border-secondary-light dark:border-secondary-dark pb-4"
        >
          <Typography preset="Text lg/Semibold" color="primary" text={$_("documentLibraryList")} />
          <Button on:click={() => viewAll()} color="light">{$_("seeAll")}</Button>
        </div>

        <div>
          {#if latestDocumentCollection.length > 0}
            <div
              class={`${$showAI || $showNotif ? "grid grid-cols-1 xl:grid-cols-3" : "grid grid-cols-2 lg:grid-cols-4"} gap-4`}
            >
              {#each latestDocumentCollection as document}
                <LibraryCard
                  collectionName={document.collectionName}
                  totalDocuments={document.totalDocuments}
                  status={document.status}
                  functButton={() => {
                    // selectCollection({
                    //   id: document.collectionId,
                    //   name: document.collectionName,
                    //   description: document.description,
                    //   totalDocuments: document.totalDocuments,
                    // });
                  }}
                  functAddDocuments={() => {
                    tmpCollectionId = document.collectionId;
                    tmpHref = document.href;
                    // categories = document.categories;
                    // selectTypeCollection("Documents");
                    // categories = document.categories;
                    otherCategory = "";
                    addOtherCategory = false;
                    addDocumentsModal = true;
                  }}
                  functEdit={() => {
                    tmpCollectionId = document.collectionId;
                    tmpCollectionName = document.collectionName;
                    tmpOldName = document.collectionName;
                    tmpDescription = document.description;
                    attributeNotUnique = false;
                    nameIsAvailable = true;
                    editModal = true;
                  }}
                  functDelete={() => {
                    tmpCollectionId = document.collectionId;
                    tmpCollectionName = document.collectionName;
                    deleteModal = true;
                  }}
                  funcDetailCollection={() => {
                    goto(document.href);
                  }}
                />
              {/each}
            </div>
          {:else if gettingDatas}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else}
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          {/if}
        </div>
      </div>
    </div>
  </div>
{/if}

<Modal title={$_("createNewLibrary")} bind:open={addLibraryModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("libraryName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={handleInput}
        placeholder={$_("whatIsYourLibraryTitle")}
      />
      {#if !nameIsAvailable}
        <Helper color="red">
          <span class="font-medium">{$_("nameAlreadyTaken")}</span>
        </Helper>
      {/if}
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}`} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={createCollection}
      disabled={!tmpCollectionName || !nameIsAvailable}
      class="w-full">{$_("addNewLibrary")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editLibraryName")} bind:open={editModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("libraryName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={handleInput}
        placeholder={$_("whatIsYourLibraryTitle")}
      />
      {#if !nameIsAvailable}
        <Helper color="red">
          <span class="font-medium">{$_("nameAlreadyTaken")}</span>
        </Helper>
      {/if}
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={editCollectionName}
      disabled={!tmpCollectionName || !nameIsAvailable}
      class="me-2 w-full">{$_("save")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
      {$_("library")}
      {tmpCollectionName}?
    </h3>
    <Button onclick={deleteDocumentCollection} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>

<!-- <LoadingContainer /> -->
