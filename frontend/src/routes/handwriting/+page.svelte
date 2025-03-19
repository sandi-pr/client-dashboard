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
  import { Modal, Button, Input, Helper } from "flowbite-svelte";
  import {
    StateCard,
    Typography,
    SectionHeader,
    LibraryCard,
    triggerToast,
    ActionCard,
    BotCard,
    check,
    isBakrie,
    ModalObject,
    triggerLoading,
    ModalFilePreview,
  } from "$lib";
  import { _ } from "svelte-i18n";
  import {
    ExclamationCircleSolid,
  } from "flowbite-svelte-icons";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { goto } from "$app/navigation";
  import {
    triggerSelectCollection,
    triggerShowAI,
    triggerSelectDropdown,
    triggerReload,
    showAI,
    showNotif,
  } from "$lib";
  
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
  let tmpOldName = $state("");
  let tmpDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpHref = $state("");
  let categories = $state([]);
  let otherCategory = $state("");
  let addOtherCategory = $state(false);
  let tmpCategories = $state("");
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let addLibraryModal = $state(false);
  let addCollectionModal = $state(false);
  let addDocumentsModal = $state(false);
  let editModal = $state(false);
  let deleteModal = $state(false);
  let dropdownOpen = $state([false]);
  let tmpDocumentType = $state("Handwriting");
  let selectedTypeCollection = $state("Handwriting (JPG, JPEG, PNG)");
  let acceptedFileTypes = $state(".jpg, .jpeg, .png");
  let timer;
  let fileNotSupport = $state(false);
  let attributeNotUnique = $state(false);
  let nameIsAvailable = $state(true);
  let gettingBotData = $state(false);
  let isProcess = $state(false);
  let dropdownOpenBot = $state([]);
  let alltags = $state([]);
  let tags = $state([]);
  let folderMode = $state(false);
  let audioUrl = $state(null);
  let previewModal = $state(false);
  let dataSend = $state({ detected_text: "" });
  let selectedFile = $state(null);
  let groupedFiles = $state(null);

  let libraryId, botId, botIdStrapi, collectionId;

  let bots = $state([
    {
      botsName: "Handwriting’s Assistant",
      botDescription: "",
      botsID: 12345,
      urlBot: "#",
      deployStatus: "Ready",
      deployedPlatform: "",
      bot_profile: "",
    }
  ]);

  async function getActionItems() {
    actionItems = [
      {
        onclick: () => {
          addDocumentsModal = true;
          attributeNotUnique = false;
          tmpCollectionName = "";
          tmpDescription = "";
        },
        // href: await $call.getPage("client-dashboard", "document-collection-create"),
        icon: "penToolPlus",
        actionName: "Add New Handwriting to Read",
        desc: "desc Add New Handwriting to Read",
      },
      {
        onclick: async () => {
          let urlBot = await $call.getPage("hexabot", "home-page") + "?id=" + botId
          window.open(urlBot, "_blank");
        },
        icon: "messageChatSquare",
        actionName: "Chat with Handwriting Assistant",
        desc: "desc Chat with Handwriting Assistant",
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
  let gettingStatusBot = $state(false);

  async function viewAll() {
    goto(await $call.getPage("client-dashboard", "handwriting-collections") + "?libraryId=" + libraryId);
  }

  async function startReading(botId) {
    let token = localStorage.getItem("token");

    if (token) {
      try {
        // gettingStatusBot = true;

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

      let link = await $call.getPage("client-dashboard", "handwriting-documents");

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

    console.log("tmpCollectionId", tmpCollectionId);

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

  function triggerToastResponse (response) {
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
      return false;
    }
    return true
  }

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

  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  async function groupFilesByPath(files) {
    const structure = {};

    for (const file of files) {
      const parts = file.webkitRelativePath.split("/");
      const fileName = parts.pop();
      const base64Data = await readFileAsBase64(file);

      // 0, 1 Sub Folder
      let folderKey = parts.slice(0, 1).join("/");
      if (!folderKey) folderKey = "root";

      if (!structure[folderKey]) {
        structure[folderKey] = {};
      }

      const fileObject = {
        id: Object.keys(structure[folderKey]).length + 1,
        name: file.name,
        path: file.webkitRelativePath || file.name,
        document_title: file.name,
        description: tmpDocumentDescription,
        size: formatBytes(file.size),
        ext: "." + file.name.split(".").pop(),
        base64: base64Data,
        type: tmpDocumentType
      };

      structure[folderKey][fileName] = fileObject;
    }

    console.log('structure::::::::::::')
    console.log(structure)
    return structure

  }

  const handleFileUploadWrapper = (isFolder = false) => {

    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.webkitdirectory = isFolder;
    inputElement.accept = acceptedFileTypes;
    inputElement.multiple = true;
    inputElement.replaceWith(inputElement.cloneNode(true));

    inputElement.addEventListener("change", async (event) => {  
      if (isFolder) {
        groupedFiles = await groupFilesByPath(event.target.files);
        return;
      }

      const filesSelect = event.target.files;

      const fileTypeRegex = new RegExp(
        acceptedFileTypes
          .split(",")
          .map((type) => type.trim().replace(".", "\\."))
          .join("|") + "$",
        "i",
      );

      // console.log(filesSelect.length);
      isProcess = true;

      for (const file of filesSelect) {
        if (file && fileTypeRegex.test(file.name)) {
          fileNotSupport = false;
          const base64Data = await readFileAsBase64(file);

          files.push({
            id: files.length < 1 ? 1 : files.length + 1,
            name: file.name,
            document_title: tmpDocumentTitle,
            description: tmpDocumentDescription,
            size: formatBytes(file.size),
            base64: base64Data,
            ext: "." + file.name.split(".").pop(),
            type: tmpDocumentType,
          });
          tmpDocumentTitle = file.name.split(".").slice(0, -1).join(".");

          // console.log("cek files", files[0]);
          const collectionId = localStorage.getItem("collectionIdhandwriting");

          if (files.length === 1) {
            const status = await checkLibraryName(tmpDocumentTitle.trim(), collectionId, "document");
            if (status === true) {
              nameIsAvailable = true;
            } else if (status === false) {
              nameIsAvailable = false;
            }
          }

          files = [...files];
          // console.log(JSON.stringify(files));
          // console.log("total file", files.length);
        } else {
          fileNotSupport = true;
        }
      }

      isProcess = false;
    });
    inputElement.click();
  };

  async function sendFolder() {
    triggerLoading(true);
    var currentFolder = []

    // Buat Collection by Name Folder Dan Document
    for (var folder in groupedFiles) {
      var parent_id;
      var sub_folder;
      var main_folder;

      if (folder.split('/').length > 1) {
        main_folder = folder.split('/')[0].trim()
        sub_folder = folder.split('/')[1].trim()

        if (currentFolder.length) {
          let finded = currentFolder.find((f) => f[main_folder] == main_folder)
          if (finded) {
            parent_id = finded.id
          }
        }
      }

      dataSend = {
        collection_name: folder,
        description: '',
        parent_id: parseInt(collectionId),
      };

      if (sub_folder && parent_id) {
        dataSend["collection_name"] = sub_folder
        dataSend["parent_id"] = parseInt(parent_id)
      }
    }

    let token = localStorage.getItem("token");
    let res = await $call.call(
        "client-dashboard", "create-library-collection", {
            URLQuery: {},
            URLParams: {},
            Params: dataSend,
            Headers: {
              Authorization: "Bearer " + token,
            },
        });

    let statusRes = triggerToastResponse(res)
    if (!statusRes) {
      triggerLoading(false);
      return
    }

    currentFolder.push(res.result.data)
    files = groupedFiles[folder]

    if (!parent_id && !sub_folder) {
      parent_id = res.result.data.id
    }

    for (const doc in files) {
      if (sub_folder && parent_id) {
        files[doc].id = parent_id
      }

      dataSend = {
        document_name: files[doc].name,
        description: files[doc].description,
        parent_id: parseInt(res.result.data.id),
        document_type: files[doc].type,
        documents: JSON.stringify(Array(files[doc])),
        category: "",
      };

      let resdoc = await $call.call("client-dashboard", "create-library-document", {
        URLQuery: {},
        URLParams: {},
        Params: dataSend,
        Headers: {
          Authorization: "Bearer " + token
        },
      });

      let statusResDoc = triggerToastResponse(resdoc)
      if (!statusResDoc) {
        triggerLoading(false);
        return
      }
    }

    addFolderModal = false;
    getCollections();
    triggerLoading(false);
  }

  async function addDocumentByCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {

      if(groupedFiles != null) {
        sendFolder();
        return
      }

      try {
        addDocumentsModal = false;
        triggerLoading(true);
        console.log("trigger loading");

        for (let i = 0; i < files.length; i++) {
          if (
            tmpDocumentType === "Handwriting" ||
            tmpDocumentType === "Document Images" ||
            tmpDocumentType === "Document Audio"
          ) {
            let response3;
            if (tmpDocumentType === "Handwriting" || tmpDocumentType === "Document Images") {
              response3 = await $call.call("hexa-ai-rag", "trocr-predict", {
                URLQuery: {},
                URLParams: {},
                Params: {
                  image: files[i].base64,
                  model:
                    tmpDocumentType === "Document Images"
                      ? "trocr:base"
                      : sessionStorage.getItem("model_trocr"),
                },
                Headers: {
                  Authorization: "Bearer " + token,
                },
              });
              // console.log("trocr-predict", response3);
              // response3 = {
              //   status: 200,
              //   result: {
              //     slices: [],
              //     text: "Belajian Tikin Bahagia\nHari-hari gajian selalu karena udah makan dalam\nnikin\nnakan.\nHal-\nbak aja uang ters-\nn kelasan dari pada menyat mengak dalam lari\nmelam lari bikin kita.\nmelu kita,\nmeli kita aja pada uang bisa,\npada bisa pada bahwa jaju kali bisa uang mengikin kekali bah turu kosia.\nTer.\nTakarta,\nTerimas AHY berkeni\nTerinti kali yang takarta.",
              //   },
              // };
            } else {
              response3 = await $call.call("hexa-ai-rag", "audio-predict", {
                URLQuery: {},
                URLParams: {},
                Params: {
                  audio_file: files[i].base64,
                  ext: files[i].name.split(".").pop(),
                },
                Headers: {
                  Authorization: "Bearer " + token,
                },
              });
              console.log("audio-predict", response3);
            }

            if (response3.status !== 200) {
              addDocumentsModal = true;
              triggerLoading(false);
              if (response3.status === 400) {
                console.error("Invalid parameters:", response3);
                triggerToast($_("invalidParameters"), true, "error");
              } else if (response3.status === 401) {
                console.error("Unauthorized access:", response3);
                triggerToast($_("unauthorized"), true, "error");
              } else if (response3.status === 500) {
                console.error("Server error:", response3);
                triggerToast($_("serverError"), true, "error");
              } else {
                console.error("Unexpected status code:", response3);
                triggerToast(`${$_("unexpectedError")}: ${response3.status}`, true, "error");
              }
              return;
            }

            if (response3?.result?.text !== null) {
              let index = 1;
              let originalName = files[i].name.split(".").slice(0, -1).join(".");
              let newName = originalName;

              for (let i = 0; i < index; i++) {
                let statusAvailable = await checkLibraryName(newName.trim(), collectionId, "document");

                if (statusAvailable === true) {
                  break;
                } else {
                  newName = `${originalName}(${index})`;
                  index++;
                }
              }

              files[i].name = newName + files[i].ext;

              dataSend = {
                document_name: files.length > 1 ? newName : tmpDocumentTitle,
                description: tmpDocumentDescription,
                parent_id: parseInt(collectionId),
                documents: JSON.stringify([files[i]]),
                document_type: tmpDocumentType,
                category: tags.join(","),
                detected_text: response3.result.text,
                result_trocr: JSON.stringify(response3.result),
                edited: false,
              };

              // convert base64 to url
              if (tmpDocumentType === "Document Audio") {
                let file = base64toBlob(files[i].base64, "audio/mpeg");
                audioUrl = URL.createObjectURL(new Blob([file], { type: "audio/mp3" }));
              }
              triggerLoading(false);
            }
          }

          if (
            tmpDocumentType === "Handwriting" ||
            tmpDocumentType === "Document Images" ||
            tmpDocumentType === "Document Audio"
          ) {
            console.log("ayam");
            if (files.length > 1) {
              await sendDocument();
            } else {
              console.log("bebek");
              previewModal = true;
              selectedFile = files[0];
            }
          } else {
            console.log("cacing");

            let index = 1;
            let originalName = files[i].name.split(".").slice(0, -1).join(".");
            let newName = originalName;

            for (let i = 0; i < index; i++) {
              let statusAvailable = await checkLibraryName(newName.trim(), collectionId, "document");

              if (statusAvailable === true) {
                break;
              } else {
                newName = `${originalName}(${index})`;
                index++;
              }
            }

            files[i].name = newName + files[i].ext;

            dataSend = {
              document_name: files.length > 1 ? newName : tmpDocumentTitle,
              // document_name: newName,
              description: tmpDocumentDescription,
              parent_id: parseInt(collectionId),
              documents: JSON.stringify([files[i]]),
              document_type: tmpDocumentType,
              category: tags.join(","),
            };

            // console.log("cek data send", dataSend);
            await sendDocument();
          }
        }

        files = [];

        tmpDocumentTitle = "";
        tmpDocumentDescription = "";
        tags = [];
        attributeNotUnique = false;
        triggerLoading(false);
        triggerToast($_("succesfullyUploaded"), true);

        if (
          tmpDocumentType !== "Handwriting" &&
          tmpDocumentType !== "Document Images" &&
          tmpDocumentType !== "Document Audio"
        ) {
          dataSend = null;
          previewModal = false;
          addDocumentsModal = false;
          getCollections();
          triggerToast($_("succesfullyUploaded"), true);
        }
      } catch (error) {
        addDocumentsModal = true;
        console.error("Error during document upload:", error);
        triggerToast($_("failedUploadDocument"), true, "error");
        triggerLoading(false);
        return;
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  async function sendDocument(single = false) {
    let token = localStorage.getItem("token");

    if (token === null) {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
    triggerLoading(true);
    let response = await $call.call("client-dashboard", "create-library-document", {
      URLQuery: {},
      URLParams: {},
      Params: dataSend,
      Headers: {
        Authorization: "Bearer " + token,
      },
    });

    if (response.status !== 200) {
      addDocumentsModal = true;
      triggerLoading(false);
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

    if (single) {
      dataSend = null;
      previewModal = false;
      addDocumentsModal = false;
      triggerLoading(false);
      triggerToast($_("succesfullyUploaded"), true);
    }
    return;
  }

  async function getLabels() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      let response = await $call.call("client-dashboard", "get-label", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });
      if (response.status !== 200) {
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

      alltags = response.result.labels;
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
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

  const handleInput = async (event, parentId = collectionId, type = "document") => {
    const name = event.target.value;
    const status = await checkLibraryName(name.trim(), parentId, type);

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
    libraryId = localStorage.getItem("libraryId" + "handwriting");
    collectionId = localStorage.getItem("collectionId" + "handwriting");
    botId = localStorage.getItem("botId" + "handwriting");
    botIdStrapi = localStorage.getItem("botIdStrapi" + "handwriting");
    bots[0].urlBot = (await $call.getPage("hexabot", "home-page")) + "?id=" + botId;
    bots[0].botsID = botIdStrapi;
    getActionItems();
    getCountDocumentsPerMonth();
    getCollections();
    getLabels();
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
        text={$_("handwritingReader")}
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
          <Typography preset="Text lg/Semibold" color="primary" text={$_("handwritingCollection")} />
          {#if latestDocumentCollection.length > 0}
            <Button on:click={() => viewAll()} color="light">{$_("seeAll")}</Button>
          {/if}
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
                  usedFor="libraries_collections"
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
                    tmpCollectionId = document.id;
                    tmpHref = document.href;
                    // categories = document.categories;
                    // selectTypeCollection("Documents");
                    // categories = document.categories;
                    otherCategory = "";
                    addOtherCategory = false;
                    addDocumentsModal = true;
                  }}
                  functEdit={() => {
                    tmpCollectionId = document.id;
                    tmpCollectionName = document.collectionName;
                    tmpOldName = document.collectionName;
                    tmpDescription = document.description;
                    attributeNotUnique = false;
                    nameIsAvailable = true;
                    editModal = true;
                  }}
                  functDelete={() => {
                    tmpCollectionId = document.id;
                    tmpCollectionName = document.collectionName;
                    deleteModal = true;
                  }}
                  funcDetailCollection={() => {
                    goto(document.href);
                  }}
                  {libraryId}
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
          <Typography preset="Text lg/Semibold" color="primary" text={$_("handwritingBotsList")} />
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

<ModalObject
  bind:modalObjectOpen={addDocumentsModal}
  bind:nameIsAvailable
  bind:tmpDocumentTitle
  {handleFileUploadWrapper}
  bind:isProcess
  bind:files
  {alltags}
  bind:tags
  {addDocumentByCollection}
  {handleInput}
  {fileNotSupport}
  {folderMode}
  usedFor={"Handwriting"}
  maxDepth={1}
/>

<ModalFilePreview
  bind:previewModal
  bind:dataSend
  cancelOnClick={() => {
    previewModal = false;
    addDocumentsModal = true;
  }}
  okOnClick={sendDocument}
  {selectedFile}
  {audioUrl}
/>

<Modal title={$_("createNewCollection")} bind:open={addCollectionModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("collectionName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={(event) => handleInput(event, libraryId, 'collection')}
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
        on:input={(event) => handleInput(event, libraryId, 'collection')}
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

<!-- <Modal title={$_("editLibraryName")} bind:open={editModal} size="xs">
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
</Modal> -->

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
