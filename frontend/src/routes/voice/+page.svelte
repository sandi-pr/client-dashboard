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
    BotCard,
    check,
    isBakrie,
    setVoiceLibraryIds,
    voiceLibraryId,
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

  let libraryId, botId, botIdStrapi;
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
  let addCollectionModal = $state(false);
  let addDocumentsModal = $state(false);
  let editModal = $state(false);
  let deleteModal = $state(false);
  let dropdownOpen = $state([false]);
  let acceptedFileTypes = $state(".txt, .pdf");
  let selectedTypeCollection = $state("Documents (TXT, PDF)");
  let timer;
  let fileNotSupport = $state(false);
  let attributeNotUnique = $state(false);
  let nameIsAvailable = $state(true);
  let gettingBotData = $state(false);
  let dropdownOpenBot = $state([]);

  let bots = $state([
    {
      botsName: "Voice’s Assistant",
      botDescription: "",
      botsID: 1,
      urlBot: "",
      deployStatus: "Ready",
      deployedPlatform: "",
      bot_profile: "",
    }
  ]);

  async function getActionItems() {
    actionItems = [
      {
        onclick: () => {
          addCollectionModal = true;
          attributeNotUnique = false;
          tmpCollectionName = "";
          tmpDescription = "";
        },
        // href: await $call.getPage("client-dashboard", "document-collection-create"),
        icon: "penToolPlus",
        actionName: "Add New Voice Record",
        desc: "desc Add New Voice Record",
      },
      {
        onclick: async () => {
          let urlBot = await $call.getPage("hexabot", "home-page") + "?id=" + botId
          window.open(urlBot, "_blank");
        },
        // href: await $call.getPage("hexabot", "home-page") + "?id=" + botId,
        icon: "messageChatSquare",
        actionName: "Chat with Voice Assistant",
        desc: "desc Chat with Voice Assistant",
      },
    ];
  }

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    for (let i = 0; i < dropdownOpenBot.length; i++) {
      dropdownOpenBot[i] = false;
    }
  };

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
    goto(await $call.getPage("client-dashboard", "voice-collections") + "?libraryId=" + libraryId);
  }

  async function checkLibrariesVoice() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "find-create-library-by-name", {
        URLQuery: {},
        URLParams: {},
        Params: {
          library_name: "Voice Analyzer",
          description: "Voice Analyzer Collection",
        },
        Headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.status !== 200) {
        console.log("status not 200", response);
        return;
      }

      console.log("created voice library", response);

      localStorage.setItem("libraryIdVoice", response.result.id);
      setVoiceLibraryIds(response.result.id);
    }
  }

  async function startReading(botId) {
    let token = localStorage.getItem("token");

    if (token) {
      try {
        const response = await $call.call("client-dashboard", "add-document-by-bot-id", {
          URLQuery: {},
          URLParams: {
            id: botId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerToast("error");
          console.log(response);
          return;
        }

        // setTimeout(getBots, 3000);
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function createCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "create-library-collection", {
        URLQuery: {},
        URLParams: {},
        Params: {
          collection_name: tmpCollectionName,
          description: tmpDescription,
          parent_id: parseInt(libraryId),
          bot_id: botIdStrapi,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        if (response.status === 400) {
          console.error("Invalid parameters:", response);
          triggerToast($_("invalidParameters"), true, "error");
        } else if (response.status === 401) {
          console.error("Unauthorized access:", response);
          triggerToast($_("unauthorized"), true, "error");
        } else if (response.status === 500) {
          console.error("Server error:", response);
          triggerToast($_("serverError"), true, "error");
        } else {
          console.error("Unexpected status code:", response);
          triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
        }
        return;
      }

      tmpCollectionName = "";
      tmpDescription = "";

      triggerToast($_("successCreateCollection"), true);

      addCollectionModal = false;
      getCollections();
    }
  }

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      collections = [];
      documents = [];
      let response = await $call.call("client-dashboard", "get-library-collections", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: 1,
          page_size: 8,
          parent_id: parseInt(libraryId),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        gettingDatas = false;
        if (response.status === 400) {
          console.error("Invalid parameters:", response);
          triggerToast($_("invalidParameters"), true, "error");
        } else if (response.status === 401) {
          console.error("Unauthorized access:", response);
          triggerToast($_("unauthorized"), true, "error");
        } else if (response.status === 500) {
          console.error("Server error:", response);
          triggerToast($_("serverError"), true, "error");
        } else {
          console.error("Unexpected status code:", response);
          triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
        }
        return;
      }

      if (response.result.data.length > 0) {
        collections = collections.concat(response.result.data);
      }

      let link = await $call.getPage("client-dashboard", "voice-documents");

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

        return {
          id: collection.id,
          collectionName: collection.attributes.name,
          description: collection.attributes?.meta?.description || "",
          type: collection.attributes.type,
          href: link + "?collectionId=" + collection.id + "&" + `libraryId=${libraryId}`,
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
        let response = await $call.call("client-dashboard", "update-library-collection", {
          URLQuery: {},
          URLParams: {
            id: tmpCollectionId.toString(),
          },
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
          if (response.status === 400) {
            console.error("Invalid parameters:", response);
            triggerToast($_("invalidParameters"), true, "error");
          } else if (response.status === 401) {
            console.error("Unauthorized access:", response);
            triggerToast($_("unauthorized"), true, "error");
          } else if (response.status === 500) {
            console.error("Server error:", response);
            triggerToast($_("serverError"), true, "error");
          } else {
            console.error("Unexpected status code:", response);
            triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
          }
          return;
        }

        tmpCollectionName = "";
        tmpDescription = "";

        editModal = false;
        triggerToast($_("successUpdateCollection"), true);
        getCollections();
      } catch (error) {
        console.error("Error during collection editing:", error);
        triggerToast($_("failedEditCollection"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
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

  const checkLibraryName = async (name) => {
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
    const status = await checkLibraryName(name.trim());

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
    // await checkLibrariesVoice();
    libraryId = localStorage.getItem("libraryId" + "voice");
    botId = localStorage.getItem("botId" + "voice");
    botIdStrapi = localStorage.getItem("botIdStrapi" + "voice");
    bots[0].urlBot = (await $call.getPage("hexabot", "home-page")) + "?id=" + botId;
    bots[0].botsID = botId;

    console.log("bots.urlBot", bots[0].urlBot);
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
        text={$_("voiceAnalyzer")}
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
      <!-- <div class="flex flex-col gap-4">
        <SectionHeader title={$_("reports")} />
        <ChartCard {chartAction} chartList={[dataChartCollecting]} />
      </div> -->

      <div class="flex flex-col gap-4">
        <div
          class="flex items-center justify-between border-b border-secondary-light dark:border-secondary-dark pb-4"
        >
          <Typography preset="Text lg/Semibold" color="primary" text={$_("voiceCollection")} />
          {#if latestDocumentCollection.length > 0}
            <Button on:click={() => viewAll()} color="light">{$_("seeAll")}</Button>
          {/if}
        </div>

        <div>
          {#if latestDocumentCollection.length > 0}
            <div
              class={`${$showAI || $showNotif ? "grid grid-cols-1 xl:grid-cols-3" : "grid grid-cols-2 lg:grid-cols-4"} gap-4`}
            >
              {#each latestDocumentCollection as collection}
              <LibraryCard
                collectionName={collection.collectionName}
                functAddDocuments={() => {
                  tmpCollectionId = collection.id;
                  // tmpHref = collection.href;
                  // categories = collection.categories;
                  // selectTypeCollection("Documents");
                  // categories = collection.categories;
                  // otherCategory = "";
                  // addOtherCategory = false;
                  addDocumentsModal = true;
                }}
                functEdit={() => {
                  tmpCollectionId = collection.id;
                  tmpCollectionName = collection.collectionName;
                  tmpOldName = document.collectionName;
                  tmpDescription = collection.description;
                  nameIsAvailable = true;
                  editModal = true;
                }}
                functDelete={() => {
                  tmpCollectionId = collection.id;
                  tmpCollectionName = collection.collectionName;
                  deleteModal = true;
                }}
                funcDetailCollection={() => {
                  goto(collection.href);
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
      <div class="flex flex-col gap-4">
        <div
          class="flex items-center justify-between border-b border-secondary-light dark:border-secondary-dark pb-4"
        >
          <Typography preset="Text lg/Semibold" color="primary" text={$_("voiceBotsList")} />
        </div>
        {#if gettingBotData}
          <div class="flex flex-col gap-4 w-full">
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          </div>
        {:else if bots.length === 0}
          <StateCard illustration={IllustrationEmpty} description={$_("noDataFound")} />
        {:else}
          <div class={`grid ${$showAI || $showNotif ? "grid-cols-1 xl:grid-cols-3" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4"} gap-4`}>
            {#each bots as bot, index (bot.botsID)}
              <BotCard {bot} {index} {dropdownOpen} {startReading} config={{ openWebchat: true, startReading: true }} {handleDropdownFocusLoss} />
            {/each}
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<Modal title={$_("createNewCollection")} bind:open={addCollectionModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("collectionName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={handleInput}
        placeholder={$_("whatIsYourCollectionTitle")}
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
      class="w-full">{$_("addNewCollection")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editCollectionName")} bind:open={editModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("collectionName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={handleInput}
        placeholder={$_("whatIsYourCollectionTitle")}
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
      "{tmpCollectionName}"?
    </h3>
    <Button onclick={deleteDocumentCollection} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>

<!-- <LoadingContainer /> -->
