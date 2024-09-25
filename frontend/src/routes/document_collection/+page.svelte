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
  import { Modal, Button, Input, P, Indicator, Alert, Select } from "flowbite-svelte";
  import {
    StateCard,
    Typography,
    SectionHeader,
    Card,
    ChartCard,
    CollectionCard,
    triggerToast,
    UploadFile,
    ActionCard,
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
  } from "$lib";
  import LoadingContainer from "$lib/components/LoadingContainer.svelte";

  import { call } from "$lib";

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
  let tmpDescription = $state("");
  let tmpCollectionId = $state("");
  let categories = $state([]);
  let otherCategory = $state("");
  let addOtherCategory = $state(false);
  let tmpCategories = $state("");
  let tmpDocumentType = $state("");
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let addCollectionsModal = $state(false);
  let addDocumentsModal = $state(false);
  let editModal = $state(false);
  let deleteModal = $state(false);
  let dropdownOpen = $state(false);
  let acceptedFileTypes = $state(".txt, .pdf, .doc, .docx");
  let selectedTypeCollection = $state("Documents (TXT, PDF, DOC, DOCX)");
  let timer;
  let fileNotSupport = $state(false);

  async function getActionItems() {
    actionItems = [
      {
        href: await $call.getPage("client-dashboard", "document-collection-create"),
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
          name: "Total documents",
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
          name: "Collecting data",
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
        href: await $call.getPage("client-dashboard", "document-collection"),
        dataName: "Total Documents",
        value: totalCollections.total,
        change: totalCollections.change,
        changeType: totalCollections.changeType,
      },
      {
        href: await $call.getPage("client-dashboard", "document-collection"),
        dataName: "Document Collection",
        value: readyCollection.total,
        change: readyCollection.change,
        changeType: readyCollection.changeType,
      },
      // {
      //   href: await $call.getPage("client-dashboard", "document-collection"),
      //   dataName: "Collecting Data",
      //   value: dataCollecting.total,
      //   change: dataCollecting.change,
      //   changeType: dataCollecting.changeType,
      // },
    ];
  }

  let latestDocuments = $state([]);

  let latestDocumentCollection = $state([]);

  let gettingDatas = $state(false);

  async function viewAll() {
    goto(await $call.getPage("client-dashboard", "document-collection"));
  }

  async function createCollection() {
    console.log(files);
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "create-document-collection", {
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

      let response2 = await $call.call("hexa-ai-rag", "add-document-by-collection-id", {
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

      triggerToast("created");

      addCollectionsModal = false;
      getCollections();
      // goto(await $call.getPage("client-dashboard", "document-collection"));
    }
  }

  async function getCollections() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      collections = [];
      documents = [];
      let response = await $call.call("client-dashboard", "get-document-collections", {
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

      for (const collection of collections) {
        let page = 1;
        while (true) {
          let response3 = await $call.call("client-dashboard", "get-documents-by-collectionid", {
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
            gettingDatas = false;
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

      let collectionList = await $call.getPage(
        "client-dashboard",
        "document-collection-list-document",
      );

      latestDocumentCollection = collections.map(async (collection) => {
        let categories = await getCategories(collection.id);
        const filteredDocuments = documents.filter(
          (document) => document.attributes.collection_id === collection.id,
        );

        // Hitung total dokumen
        const totalDocuments = filteredDocuments.length;

        // Hitung jumlah dokumen dengan status null
        const nullDocuments = filteredDocuments.filter(
          (document) => document.attributes.status === null,
        ).length;

        // Hitung jumlah dokumen dengan status angka
        const numberDocuments = filteredDocuments.filter((document) =>
          document.attributes.status?.includes("Reading"),
        ).length;

        // Hitung total dokumen yang Ready
        const readyDocuments = filteredDocuments.filter(
          (document) => document.attributes.status === "Ready",
        ).length;

        // Tentukan status koleksi berdasarkan kondisi yang diberikan
        let collectionStatus;
        if (nullDocuments === totalDocuments) {
          if (totalDocuments === 0) {
            // console.log("empty", collection.id);
            collectionStatus = "Empty";
          } else {
            // console.log("null", collection.id);
            collectionStatus = null;
          }
        } else if (nullDocuments > 0 || numberDocuments > 0) {
          // console.log("read", collection.id);
          collectionStatus = "Reading";
        } else if (readyDocuments === totalDocuments) {
          // console.log("ready", collection.id);
          collectionStatus = "Ready";
        }
        return {
          collectionId: collection.id,
          activationTime: format(collection.attributes.createdAt, "dd-MM-yyyy HH:mm"),
          collectionName: collection.attributes.collection_name,
          description: collection.attributes.description,
          totalDocuments: totalDocuments,
          categories: categories,
          status: collectionStatus,
          href: `${collectionList}?collectionId=${collection.id}`,
          webchat: collection.attributes.webchat_tmp === true ? true : false,
        };
      });

      console.log("cek status", latestDocumentCollection);

      gettingDatas = false;
    }
  }

  async function editCollectionName() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "update-document-collection", {
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
          return;
        }

        tmpCollectionName = "";
        tmpDescription = "";

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
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        for (const document of documents) {
          if (document.attributes.collection_id === tmpCollectionId) {
            let response = await $call.call("client-dashboard", "delete-document", {
              URLQuery: {},
              URLParams: {
                id: document.id.toString(),
              },
              Params: {},
              Headers: {
                Authorization: "Bearer " + token,
              },
            });
            if (response.status !== 200) {
              console.log(response);
              triggerToast("error");
              return;
            }
          }
        }
        let response2 = await $call.call("client-dashboard", "delete-document-collection", {
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

  async function addDocumentByCollection() {
    files[0].document_title = tmpDocumentTitle;
    files[0].description = tmpDocumentDescription;

    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        triggerLoading(true);
        let response = await $call.call("client-dashboard", "create-document", {
          URLQuery: {},
          URLParams: {},
          Params: {
            documents: JSON.stringify(files),
            collection_id: tmpCollectionId.toString(),
            document_type: tmpDocumentType,
            category: tmpCategories,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          triggerToast("error");
          triggerLoading(false);
          return;
        }

        console.log(response);

        if (tmpDocumentType === "Handwriting") {
          for (let i = 0; i < files.length; i++) {
            let response3 = await $call.call("hexa-ai-rag", "trocr-predict", {
              URLQuery: {},
              URLParams: {},
              Params: {
                image: files[i].base64,
                model: sessionStorage.getItem("model_trocr"),
              },
              Headers: {
                Authorization: "Bearer " + token,
              },
            });

            if (response3.status !== 200) {
              triggerToast("error");
              triggerLoading(false);
              console.log(response3);
              return;
            }

            console.log("trocr-predict", response3);

            if (response3?.result?.text !== null) {
              const response4 = await $call.call(
                "client-dashboard",
                "update-document-handwriting",
                {
                  URLQuery: {},
                  URLParams: {
                    id: response.result.data[i].id.toString(),
                  },
                  Params: {
                    detected_text: response3.result.text,
                    result_trocr: JSON.stringify(response3.result),
                    edited: false,
                    document_title: tmpDocumentTitle,
                    description: tmpDocumentDescription,
                  },
                  Headers: {
                    Authorization: "Bearer " + token,
                  },
                },
              );

              if (response4.status !== 200) {
                triggerToast("error");
                triggerLoading(false);
                console.log(response4);
                return;
              }

              console.log("ayam", response4);
            }
          }
        }
      } catch (error) {
        triggerToast("error");
        triggerLoading(false);
        console.log(error);
      }

      tmpCollectionId = "";
      getCollections();
      triggerLoading(false);
      addDocumentsModal = false;

      triggerToast("created");
    }
  }

  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  const handleFileUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = acceptedFileTypes;
    inputElement.multiple = false;

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      const fileTypeRegex = new RegExp(
        acceptedFileTypes
          .split(",")
          .map((type) => type.trim().replace(".", "\\."))
          .join("|") + "$",
        "i",
      );

      console.log(fileTypeRegex.test(file.name));

      if (file && fileTypeRegex.test(file.name)) {
        fileNotSupport = false;
        const base64Data = await readFileAsBase64(file);
        files.push({
          name: file.name,
          document_title: tmpDocumentTitle,
          description: tmpDocumentDescription,
          size: formatBytes(file.size),
          base64: base64Data,
        });
      } else {
        fileNotSupport = true;
      }
    });

    inputElement.click();
  };

  async function readFileAsBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        const base64Data = reader.result.split(",")[1];
        resolve(base64Data);
      };
      reader.onerror = reject;

      reader.readAsDataURL(file);
    });
  }

  const selectTypeCollection = (type) => {
    tmpDocumentType = type;

    if (type === "Documents") {
      selectedTypeCollection = "Documents (TXT, PDF, DOC, DOCX)";
      acceptedFileTypes = ".txt, .pdf, .doc, .docx";
    } else if (type === "Images") {
      selectedTypeCollection = "Images (JPG, JPEG, PNG, WEBP, IMG)";
      acceptedFileTypes = ".jpg, .jpeg, .png, .webp, .img";
    } else if (type === "Handwriting") {
      selectedTypeCollection = "Handwriting (JPG, JPEG, PNG)";
      acceptedFileTypes = ".jpg, .jpeg, .png";
      // acceptedFileTypes = ".pdf, .jpg, .jpeg, .png, .webp, .img";
    }

    files = [];

    dropdownOpen = false;
  };

  async function getCountDocumentsPerMonth() {
    const today = new Date();
    let earlyDateCurrentMonth = new Date(today.getFullYear(), today.getMonth(), 1);
    let lastDateCurrentMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0);
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      for (let i = 0; i < 12; i++) {
        let response = await $call.call("client-dashboard", "get-count-documents-per-month", {
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

        tmpMonthlyCollectionsCounts.push(response.result.document.meta.pagination.total);
        // TODO: adjust response to collection ready
        tmpMonthlyCollectionReadyCounts.push(response.result.general.meta.pagination.total);
        // TODO: adjust response to collecting data
        tmpMonthlyCollectingDataCounts.push(response.result.contract.meta.pagination.total);

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

      tmpMonthlyCollectionsCounts.reverse();
      tmpMonthlyCollectionReadyCounts.reverse();
      tmpMonthlyCollectingDataCounts.reverse();

      monthlyCollectionsCounts = tmpMonthlyCollectionsCounts;
      monthlyCollectionReadyCounts = tmpMonthlyCollectionReadyCounts;
      monthlyCollectingDataCounts = tmpMonthlyCollectingDataCounts;

      let totalCollections = groupingDataChart(monthlyCollectionsCounts);
      let readyCollection = groupingDataChart(monthlyCollectionReadyCounts);
      let dataCollecting = groupingDataChart(monthlyCollectingDataCounts);

      getChartAction(totalCollections, readyCollection, dataCollecting);
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

  async function getCategories(id) {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "get-categories-by-collection-id", {
        URLQuery: {},
        URLParams: {},
        Params: {
          collection_id: id.toString(),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(id, response);
        return;
      }

      return response.result.categories.map((category) => {
        return {
          name: category,
          value: category,
        };
      });
    }
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

  onMount(async () => {
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

<div class="flex dark:bg-primary-dark">
  <div class="flex flex-col w-full pb-10 m-4 gap-3 md:gap-8">
    <!-- User Information -->
    <Typography
      preset="Display sm/Semibold"
      text={$_("documentCollectionOverview")}
      propClass="fluid-xl"
    />

    <!-- Action Items -->
    <div class="flex flex-col gap-4">
      <SectionHeader title={$_("action")} />
      <ActionCard {actionItems} onclick={() => (addCollectionsModal = true)} />
    </div>

    <!-- Reports Monitoring -->
    <div class="flex flex-col gap-4">
      <SectionHeader title={$_("reports")} />
      <ChartCard {chartAction} chartList={[dataChartAllDocs, dataChartCollections]} />
    </div>

    <div class="flex flex-col gap-4">
      <div
        class="flex items-center justify-between border-b border-secondary-light dark:border-secondary-dark pb-4"
      >
        <Typography preset="Text lg/Semibold" color="primary" text={$_("documentLatest")} />
        <div></div>
        <Button on:click={() => viewAll()} color="light">{$_("seeAll")}</Button>
      </div>

      <div>
        {#if latestDocumentCollection.length > 0}
          <div
            class={`${$showAI ? "grid grid-cols-1 xl:grid-cols-3" : "grid grid-cols-2 lg:grid-cols-4"} gap-4`}
          >
            {#each latestDocumentCollection as document}
              <CollectionCard
                collectionName={document.collectionName}
                totalDocuments={document.totalDocuments}
                status={document.status}
                functButton={() => {
                  selectCollection({
                    id: document.collectionId,
                    name: document.collectionName,
                    description: document.description,
                    totalDocuments: document.totalDocuments,
                  });
                }}
                functAddDocuments={() => {
                  tmpCollectionId = document.collectionId;
                  categories = document.categories;
                  selectTypeCollection("Documents");
                  categories = document.categories;
                  addDocumentsModal = true;
                }}
                functEdit={() => {
                  tmpCollectionId = document.collectionId;
                  tmpCollectionName = document.collectionName;
                  tmpDescription = document.description;
                  editModal = true;
                }}
                functDelete={() => {
                  tmpCollectionId = document.collectionId;
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

<Modal title={$_("addDocuments")} bind:open={addDocumentsModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5 relative">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("documentType")}*`} />
      <Button
        color="light"
        class="flex justify-between"
        on:click={() => (dropdownOpen = !dropdownOpen)}
      >
        {`${selectedTypeCollection ? selectedTypeCollection : $_("selectDocumentType")}`}
        <ChevronDownOutline class="w-6 h-6 text-fg-quaternary-light dark:text-quarterary-dark" />
      </Button>
      <div
        class={`z-50 ${dropdownOpen ? "absolute top-20" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}
      >
        <ul class="py-2 text-sm text-gray-700 dark:text-gray-200">
          <button
            class="flex p-2"
            onclick={() => {
              selectTypeCollection("Documents");
            }}>Documents (TXT, PDF, DOC, DOCX)</button
          >
          <button
            class="flex p-2"
            onclick={() => {
              selectTypeCollection("Images");
            }}>Images (JPG, JPEG, PNG, WEBP, IMG)</button
          >
          <button
            class="flex p-2"
            onclick={() => {
              selectTypeCollection("Handwriting");
            }}>Handwriting (JPG, JPEG, PNG)</button
          >
        </ul>
      </div>
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("categories")}`} />
      <!-- <MultiSelect items={categories} bind:value={tmpCategories}></MultiSelect> -->
      <Select class="mt-2" items={categories} bind:value={tmpCategories} placeholder={$_("selectCategory")} />
      {#if !addOtherCategory}
        <p>
          {$_("addOtherCategory")}
          <button onclick={() => {addOtherCategory = true;}} class="ms-1 text-brand-secondary-light dark:text-brand-secondary-dark hover:underline" >
            {$_("clickHere")}
          </button>
        </p>
      {/if}
      {#if addOtherCategory}
        <Input
          placeholder={$_("enterOtherCategory")}
          bind:value={otherCategory}
          onkeydown={(e) => {
            if (
              e.key === "Enter" &&
              otherCategory !== "" &&
              !tmpCategories.includes(otherCategory)
            ) {
              for (let i = 0; i < categories.length; i++) {
                if (categories[i]?.other === true) {
                  categories.splice(i, 1);
                  break;
                }
              }
              categories = [
                ...categories,
                { value: otherCategory, name: otherCategory, other: true },
              ];
              tmpCategories = otherCategory;
              otherCategory = "";
            }
          }}
        />
      {/if}
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("documentTitle")}*`} />
      <Input bind:value={tmpDocumentTitle} placeholder={$_("enterDocumentTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Input bind:value={tmpDocumentDescription} placeholder={$_("enterDocumentDescription")} />
    </div>
    <UploadFile
      title={`${$_("uploadDocuments")}*`}
      handleFileUploadWrapper={() => handleFileUploadWrapper()}
      file={files[0]}
      actionButton={true}
      actionDelete={() => (files = [])}
      isFlex
    />
    {#if fileNotSupport}
      <Alert color="red">
        <InfoCircleSolid slot="icon" class="w-5 h-5" />
        {$_("fileNotSupport")}
      </Alert>
    {/if}
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={addDocumentByCollection} disabled={buttonDisabled} class="w-full"
      >{$_("addNewDocument")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editCollectionName")} bind:open={editModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("collectionName")}*`} />
      <Input bind:value={tmpCollectionName} placeholder={$_("whatIsYourCollectionTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={editCollectionName} class="me-2">{$_("save")}</Button>
  </svelte:fragment>
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteDocumentCollection} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>

<!-- <LoadingContainer /> -->
