<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Dean Ramhan <dean@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->
<script>
  import { P, Avatar, Button, Indicator, Modal } from "flowbite-svelte";
  import DataTableListCollection from "$lib/components/DataTableListCollection.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import { ChevronDownOutline, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import { _ } from "svelte-i18n";
  import { goto } from "$app/navigation";
  import {
    Typography,
    Badge,
    InputRw,
    TextareaRw,
    SelectedCard,
    UploadFile,
    Dropdown,
    triggerToast,
    triggerLoading,
    StateCard,
    Toggle,
    call,
  } from "$lib";
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import { findIconFromModelName, getIcons } from "$lib";
  // import Highlight from "svelte-highlight";
  // import typescript from "svelte-highlight/languages/typescript";
  // import github from "svelte-highlight/styles/github";

  // const API_URL = process.env.NODE_ENV === "development" ? PUBLIC_API_URL_DEV : PUBLIC_API_URL_PROD;
  // const STRAPI_URL =
  //   process.env.NODE_ENV === "development"
  //     ? PUBLIC_STRAPI_MOTHERSHIP_DEV
  //     : PUBLIC_STRAPI_MOTHERSHIP_PROD;

  const isFrontend = process.env.NODE_ENV === "development";

  const params = new URLSearchParams(window.location.search);
  const edit = params.get("edit");

  let onEdit = $state(false);

  let ready = $state(false);
  let botId;
  let botFile = $state({ name: "" });
  let userFile = $state({ name: "" });
  let logoFile = $state({ name: "" });
  let backgroundFile = $state({
    name: "",
  });
  let modalThemeColor = $state(false);
  let modalOverlayColor = $state(false);
  let dropdownBgAlignOpen = $state(false);
  let modalDelete = $state(false);
  let background_align = $state("center");
  let bgAlign = $state([
    { name: "Center", value: "center" },
    { name: "Left", value: "left" },
    { name: "Right", value: "right" },
  ]);

  let users = $state({});

  let collectionId = $state(null);
  let collection = $state({});
  let datas = $state([]);
  let gettingDatas = $state(false);
  let datasChecked = $state([]);
  let tmpDatasChecked = $state([]);
  let configs = $state({});
  let configsChecked = $state({});
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });
  let paginationCollection = $state({ page: 1, pageSize: 20, pageCount: 0, total: 0 });
  let paginationDocs = $state([]);
  let searchTerm = $state("");
  let dropdownOpen = $state(false);
  let modalOpen = $state(false);
  let collections = $state([]);
  let selectedCollection = $state("");
  let tmpSelectedCollection = $state("");
  let oldSelectedCollection = $state("");
  let initCollections = $state([]);
  let isSelectedLibrary = $state(false);
  let allData = $state([]);
  let search = $state(false);
  let searchCategory = $state({});
  let pages = $state(1);
  let disabledButton = $state(false);
  let collectionsChecked = $state([]);
  let checkAll = $state(false);
  let selectCategories = $state(null);
  let datasCategoriesChecked = $state({});
  let configsCategories = $state({});
  let searchTermCategory = $state({});

  let tmpThemeColor = $state("");
  let tmpOverlayColor = $state("");

  let tabDatas = $state({
    selected: "Customization",
    data: [
      {
        name: "Information",
      },
      {
        name: "Customization",
      },
      {
        name: "Data",
      },
      {
        name: "Widget Embed",
      },
    ],
  });

  let platforms = $state([]);

  let layouts = $state([
    {
      id: 1,
      name: "Modern Minimalism",
    },
    {
      id: 2,
      name: "Modern with Background",
    },
  ]);

  let themeColors = $state([
    {
      id: 1,
      name: "purple",
      value: "purple",
    },
    {
      id: 2,
      name: "blue",
      value: "blue",
    },
    {
      id: 3,
      name: "green",
      value: "green",
    },
  ]);

  let dropdownUsecaseOpen = $state(false);
  let usecase = $state([
    { name: "My Bots (Default)", value: "myBots" },
    { name: "Hexassistant", value: "hexassistant" },
    { name: "Coding Assistant", value: "codingAssistant" },
  ]);
  let usecaseSelected = $state("myBots");

  function getBrandColorName(selectedId, value) {
    const selectedColor = themeColors.find((color) => color.id === selectedId);
    return selectedColor
      ? value === "name"
        ? selectedColor.name
        : selectedColor.value
      : "Unknown Color";
  }

  let timeoutId;
  // old
  // const handleSearch = async (value) => {
  //   pagination.page = 1;
  //   searchTerm = value;

  //   await getDocuments(1);
  // };
  const handleSearch = async (value) => {
    search = true;
    pagination.page = 1;
    searchTerm = value;

    await getCollectionSub(selectedCollection);
  };

  let conversationDatas = [
    {
      id: 1,
      timestamp: "03/19/2024 22:57",
      customerName: "Elias Hoppe IV",
      conversationId: "4c64986428240cn948929284",
      actions: {
        objects: [
          {
            name: "View",
            href: "/conversation_logs",
          },
        ],
      },
    },
    {
      id: 2,
      timestamp: "03/20/2024 09:15",
      customerName: "Amelia Jones",
      conversationId: "8d74683756ds33l90928f",
      actions: {
        objects: [
          {
            name: "View",
            href: "/conversation_logs",
          },
        ],
      },
    },
    {
      id: 3,
      timestamp: "03/21/2024 13:45",
      customerName: "Mikhail Petrov",
      conversationId: "5f739dks76402dk23kj4",
      actions: {
        objects: [
          {
            name: "View",
            href: "/conversation_logs",
          },
        ],
      },
    },
  ];

  let conversationConfigs = {
    search: handleSearch,
    pagination: {
      page: 1,
      pageSize: 5,
      pageCount: 1,
      total: 5,
    },
  };

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let pageNow = 1;
      initCollections = [];
      while (true) {
        let response = await $call.call("client-dashboard", "get-librarys", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: pageNow,
            page_size: 25,
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
          initCollections = initCollections.concat(response.result.data);
          pageNow++;
        } else {
          break;
        }
      }

      console.log("libraries", initCollections);

      collections = initCollections.map((collection) => {
        return {
          value: collection.id,
          id: collection.id,
          name: collection.attributes.name,
          description: collection.attributes.meta?.description,
        };
      });
    }
  }

  async function getCollectionSub(id) {
    let token = localStorage.getItem("token");
    if (token !== null) {
      if (!search) {
        gettingDatas = true;
      }
      let response = await $call.call("client-dashboard", "get-library-collections", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: pagination.page,
          page_size: pagination.pageSize,
          parent_id: id,
          search: searchTerm,
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

      console.log("collections", response);

      if (response.result.data.length > 0) {
        pagination = response.result.meta.pagination;
        configsChecked = {
          search: true,
          searchHandler: handleSearch,
          pagination: response.result.meta.pagination,
          space: "p-4",
          actionInLeft: true,
          actions: [
            {
              id: "deselect_all",
              name: "Deselect All",
              onclick: () => {
                handleCheckAllCollection(false);
              },
            },
            {
              id: "select_all",
              name: "Select All",
              type: "brand-inverse",
              onclick: () => {
                handleCheckAllCollection(true);
              },
            },
          ],
        };
        let promises = response.result.data.map(async (document) => {
          const sub_collections = await getSubCollections(document.id);
          return {
            id: document.id,
            collection_name: document.attributes.name,
            sub_collections,
          };
        });
        allData = await Promise.all(promises);
        console.log("all data", allData);
        datasChecked = allData.map((data) => {
          return {
            id: data.id,
            checked: false,
            sub_collection: data.sub_collections.map((sub) => ({
              id: sub.id,
              checked: false,
            })),
          };
        });

        if (checkAll) {
          datasChecked = allData.map((data) => {
            return {
              id: data.id,
              checked: true,
              sub_collection: data.sub_collections.map((sub) => ({
                id: sub.id,
                checked: true,
              })),
            };
          });
        } else {
          if (tmpDatasChecked[pagination.page] === undefined) {
            if (selectedCollection === tmpSelectedCollection) {
              tmpSelectedCollection = null;
              datasChecked = datasChecked.map((data) => {
                if (collectionsChecked.includes(data.id)) {
                  data.checked = true;
                }

                // Perbarui status checked untuk setiap sub_collection
                data.sub_collection = data.sub_collection.map((sub) => {
                  if (collectionsChecked.includes(sub.id)) {
                    sub.checked = true;
                  }
                  return sub;
                });

                return data;
              });
            } else {
              datasChecked = allData.map((data) => {
                return { id: data.id, checked: false };
              });
            }
          } else {
            datasChecked = tmpDatasChecked[pagination.page];
          }
        }
      } else {
        pagination = response.result.meta.pagination;
        allData = [];
        datasChecked = [];
        configsChecked = {
          search: true,
          searchHandler: handleSearch,
          pagination: response.result.meta.pagination,
          space: "p-4",
          actionInLeft: true,
          actions: [
            {
              id: "deselect_all",
              name: "Deselect All",
              onclick: () => {
                handleCheckAllCollection(false);
              },
            },
            {
              id: "select_all",
              name: "Select All",
              type: "brand-inverse",
              onclick: () => {
                handleCheckAllCollection(true);
              },
            },
          ],
          // actions: [
          //   {
          //     id: "delete",
          //     name: "Delete",
          //     icon: "trash",
          //     type: "error",
          //     disabled: datasChecked.find((data) => data.checked) ? false : true,
          //     onclick: () => {
          //       if (datasChecked.find((data) => data.checked)) {
          //         popupModal = true;
          //       }
          //     },
          //   },
          //   {
          //     id: "add_new",
          //     name: "Add New",
          //     icon: "plus",
          //     type: "brand-inverse",
          //     onclick: () => {
          //       selectTypeCollection("Documents");
          //       addDocumentsModal = true;
          //       addOtherCategory = false;
          //     },
          //   },
          // ],
        };
      }

      gettingDatas = false;
      search = false;
    }
  }

  async function getSubCollections(id) {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      const aggregatedData = [];

      paginationCollection.page = 1;

      do {
        let response = await $call.call("client-dashboard", "get-library-collections", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: paginationCollection.page,
            page_size: paginationCollection.pageSize,
            parent_id: parseInt(id),
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

        aggregatedData.push(...response.result.data);

        paginationCollection = response.result.meta.pagination;

        paginationCollection.page++;
      } while (paginationCollection.page <= paginationCollection.pageCount);

      const promise = aggregatedData.map(async (collection) => ({
        id: collection.id,
        collection_name: collection.attributes.name,
      }));

      const subCollectionList = await Promise.all(promise);
      console.log("subCollectionList", subCollectionList);
      gettingDatas = false;
      return subCollectionList;
    } else {
      triggerToast($_("unauthorized"), true, "error");
      gettingDatas = false;
      return;
    }
  }

  function areAllUnchecked() {
    // Gabungkan data dari halaman sekarang dengan data yang tersimpan di tmpDatasChecked
    const allData = [...datasChecked, ...Object.values(tmpDatasChecked).flat()];

    // Cek apakah semua data tidak ada yang checked
    return allData.every((item) => item.checked === false);
  }

  const classStatusType = (type) => {
    if (type === "Done") {
      return "success";
    } else if (type === "Not Ready") {
      return "error";
    } else if (type === "On Progress") {
      return "warning";
    } else if (type === "WhatsApp") {
      return "success";
    } else if (type === "Web") {
      return "sky";
    }
  };

  const handleInputFile = (as) => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = "image/jpeg, image/png";

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file && (file.type === "image/jpeg" || file.type === "image/png")) {
        if (as === "bot") {
          botFile.imageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          botFile.base64Data = await readFileAsBase64(file);
          botFile.name = file.name;
        } else if (as === "user") {
          userFile.imageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          userFile.base64Data = await readFileAsBase64(file);
          userFile.name = file.name;
        } else if (as === "logo") {
          logoFile.imageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          logoFile.base64Data = await readFileAsBase64(file);
          logoFile.name = file.name;
        } else if (as === "background") {
          backgroundFile.imageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          backgroundFile.base64Data = await readFileAsBase64(file);
          backgroundFile.name = file.name;
        }
      } else {
        alert("Invalid file. Please select a valid .jpg or .png file.");
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

  function copyToClipboard(text) {
    navigator.clipboard
      .writeText(text)
      .then(() => {
        console.log("Text copied to clipboard");
        // Opsional: Tambahkan notifikasi atau toast untuk menunjukkan teks telah disalin
        triggerToast("copyToClipboard");
      })
      .catch((err) => {
        console.error("Failed to copy: ", err);
      });
  }

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
  };

  async function next() {
    if (pagination.page >= pagination.pageCount) return;

    let pageNow = pagination.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasChecked,
      };
    } else {
      tmpDatasChecked[pageNow] = datasChecked;
    }

    pagination.page = pagination.page + 1;
    // console.log(pagination.page);
    // getDocuments();
    getCollectionSub(selectedCollection);
  }

  async function previous() {
    if (pagination.page <= 1) return;

    let pageNow = pagination.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasChecked,
      };
    } else {
      tmpDatasChecked[pageNow] = datasChecked;
    }

    pagination.page = pagination.page - 1;
    // getDocuments();
    getCollectionSub(selectedCollection);
  }

  async function handleCheckAll(allChecked) {
    datasChecked = datasChecked.map((data) => ({
      id: data.id,
      checked: allChecked ? false : true,
      sub_collection: data.sub_collection.map((sub) => ({
        id: sub.id,
        checked: allChecked ? false : true,
      })),
    }));

    // console.log("Updated datasChecked (handleCheckAll):", datasChecked);
  }

  async function handleCheckAllCollection(allChecked) {
    checkAll = allChecked;

    datasChecked = datasChecked.map((data) => ({
      id: data.id,
      checked: allChecked,
      sub_collection: data.sub_collection.map((sub) => ({
        id: sub.id,
        checked: allChecked,
      })),
    }));

    // console.log("Updated datasChecked (handleCheckAllCollection):", datasChecked);
  }

  async function handleCheck(id, isSubCollection = false, parentId = null) {
    if (!Array.isArray(datasChecked)) {
      datasChecked = [];
    }

    datasChecked = datasChecked.map((data) => {
      // Jika parentId tidak ada, kita memproses level utama
      if (!isSubCollection && data.id === id) {
        if (!data.checked) {
          checkAll = false;
        }
        return {
          ...data,
          checked: !data.checked,
        };
      }

      // Jika parentId cocok, kita memproses sub_collection
      if (isSubCollection && data.id === parentId) {
        const updatedSubCollection = data.sub_collection.map((sub) => {
          if (sub.id === id) {
            return { ...sub, checked: !sub.checked };
          }
          return sub;
        });

        return {
          ...data,
          sub_collection: updatedSubCollection,
        };
      }

      return data;
    });

    let pageNow = pagination.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasChecked,
      };
    } else {
      tmpDatasChecked[pageNow] = datasChecked;
    }
    // console.log("datasChecked setelah pembaruan:", datasChecked);
  }

  let formattedDate = "Created on 2021-08-12 14:00 PM";

  let code = `<iframe style=”width: 100%; height: 100%; border: 0px; border-radius: 10px; display: block;" src="http://webchat.hexacode.org/19ca7b92-4bfd-4c20-b6ab-048207e66258"></iframe>`;

  async function getModels() {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        let response = await $call.call("hexa-ai-rag", "get-models", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("models", response);

        if (response === null) {
          if (isFrontend) {
            platforms = [
              {
                id: 1,
                name: "Llama3.1",
                model: "llama3.1:8b",
                description: "Llama3.1 Develop by Meta",
                create_at: "2024-07-23",
              },
            ];
          } else {
            triggerLoading(false);
            triggerToast("error");
          }

          return;
        }

        if (response.status === 200) {
          // BRI ADJUSTMENT
          response.result.models = response.result.models.filter(
            (model) => !model.name.includes("Qurani"),
          );

          let i = 0;
          let promise = response.result.models.map((model) => {
            return {
              id: i++,
              name: model.name,
              model: model.model,
              description: model.description,
              create_at: model.create_at,
              icon: findIconFromModelName(model.name),
              recommended: model.name === "Llama3.1" ? true : false,
            };
          });

          platforms = await Promise.all(promise);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function getCollection() {
    // search collection using params collectionId here
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "get-collection-by-id", {
          URLQuery: {},
          URLParams: {
            id: collectionId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          console.log(response);
          triggerToast("error");
          goto(await $call.getPage("client-dashboard", "document-collection"));
          return;
        }

        collection = {
          id: response.result.data.id,
          name: response.result.data.attributes.collection_name,
          createdAt: response.result.data.attributes.createdAt,
          description: response.result.data.attributes.description,
        };
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function getDocuments() {
    let token = localStorage.getItem("token");
    for (let i = 0; i < collectionId.length; i++) {
      console.log("collectionId", collectionId[i]);
      if (paginationDocs[collectionId[i]] === undefined) {
        paginationDocs[collectionId[i]] = {
          page: 1,
          pageSize: 10,
          pageCount: 0,
          total: 0,
        };
      }
      if (token !== null) {
        let response = await $call.call("client-dashboard", "get-library-documents", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: paginationDocs[collectionId[i]].page,
            page_size: paginationDocs[collectionId[i]].pageSize,
            parent_id: parseInt(collectionId[i]),
            search: searchTerm,
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

        console.log("ayam", response);

        if (response.result.data.length > 0) {
          let promises = response.result.data.map(async (document) => {
            let link = await $call.getPage("client-dashboard", "document-library-detail");
            return {
              id: document.id,
              file_name: {
                name: document.attributes.name,
                filename:
                  document.attributes.name + "." + document.attributes.path.split(".").pop(),
                type: "custom",
                size:
                  document.attributes.meta.file_size === null
                    ? "Unknown"
                    : document.attributes.meta.file_size,
              },
              type: {
                name: document.attributes.meta.document_type,
                type: "pill",
                color: "brand",
              },
              upload_date: document.attributes.createdAt,
              actions: {
                objects: [
                  {
                    href: `${link}` + "?documentId=" + document.id,
                    name: "View",
                  },
                ],
              },
            };
          });

          datas = [...datas, { [collectionId[i]]: await Promise.all(promises) }];
        } else {
          datas = [...datas, { [collectionId[i]]: [] }];
        }

        searchCategory = {
          ...searchCategory,
          [collectionId[i]]: false,
        };

        datasCategoriesChecked = {
          ...datasCategoriesChecked,
          [collectionId[i]]: datas.map((data) => {
            return { id: data.id, checked: false };
          }),
        };

        searchTermCategory = {
          ...searchTermCategory,
          [collectionId[i]]: "",
        };

        let response2 = await $call.call("client-dashboard", "get-library-collection", {
          URLQuery: {},
          URLParams: {
            id: collectionId[i].toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          console.log(response2);
          gettingDatas = false;
          return;
        }

        console.log("cekkkkk", response2);

        configsCategories = {
          ...configsCategories,
          [collectionId[i]]: {
            search: true,
            searchHandler: handleSearchCategory,
            headers: "Files",
            name: response2.result.attributes.name,
            pagination: response.result.meta.pagination,
            // isCategories: true,
            collectionId: collectionId[i],
            actions: [
              {
                id: "expand",
                icon: "chevron",
                type: "info",
              },
            ],
          },
        };
      }
    }
    gettingDatas = false;
  }

  async function getDocument(id) {
    let token = localStorage.getItem("token");
    if (token !== null) {
      let response = await $call.call("client-dashboard", "get-library-documents", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: paginationDocs[id].page,
          page_size: paginationDocs[id].pageSize,
          parent_id: parseInt(id),
          search: searchTermCategory[id],
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

      console.log("ayam", response);

      let promises = response.result.data.map(async (document) => {
        let link = await $call.getPage("client-dashboard", "document-library-detail");
        return {
          id: document.id,
          file_name: {
            name: document.attributes.name,
            filename: document.attributes.name + "." + document.attributes.path.split(".").pop(),
            type: "custom",
            size:
              document.attributes.meta.file_size === null
                ? "Unknown"
                : document.attributes.meta.file_size,
          },
          type: {
            name: document.attributes.meta.document_type,
            type: "pill",
            color: "brand",
          },
          upload_date: document.attributes.createdAt,
          actions: {
            objects: [
              {
                href: `${link}` + "?documentId=" + document.id,
                name: "View",
              },
            ],
          },
        };
      });

      datas = await Promise.all(
        datas.map(async (data) => {
          if (parseInt(Object.keys(data)[0]) === id) {
            return { [id]: await Promise.all(promises) };
          } else {
            return data;
          }
        }),
      );

      // datasChecked = datas.map((data) => {
      //   return { id: data.id, checked: false };
      // });

      configsCategories[id].pagination = response.result.meta.pagination;
    }
  }

  async function getBot() {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        let response = await $call.call("client-dashboard", "get-bot", {
          URLQuery: {},
          URLParams: {
            id: botId,
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          triggerLoading(false);
          triggerToast("error");
          goto(await $call.getPage("client-dashboard", "bots-list"));
          return;
        }

        const bot = response.result;
        const attributes = bot.attributes;

        let urlBot = (await $call.getPage("hexabot", "home-page")) + "?id=" + attributes.bot_id;

        // collectionId = attributes.collection_id;
        collectionsChecked = attributes.collection_id.split(",").map(Number);
        collectionId = collectionsChecked;
        console.log("dihh", attributes.collection_id);
        tmpSelectedCollection = attributes?.library_id;
        selectedCollection = attributes?.library_id;
        // datasChecked = collections.map((data) => {
        //   return { id: data, checked: true };
        // });

        // await getCollection();
        await getDocuments();

        if (attributes.ai_engine === "openai") {
          // hanya tampilkan platform yang modelnya openai
          platforms = platforms.filter((platform) => platform.model === "openai");
          // hapus tabDatas yang namanya Data
          tabDatas.data = tabDatas.data.filter((tab) => tab.name !== "Data");
        }

        users = {
          botName: attributes.bot_name,
          logoName: attributes.logo_name,
          logoOnly: attributes.logo_only,
          deployedPlatform: attributes?.deployed_platform || "Web",
          description: attributes.description,
          additionalPrompts: attributes.additional_prompt,
          welcomeMessage: attributes.greeting_text,
          username: attributes?.username || "User",
          primaryBrandColor:
            attributes.color_theme !== null
              ? themeColors.find((color) => color.value === attributes.color_theme).id
              : 1,
          overlayColor:
            attributes.overlay_color !== null
              ? themeColors.find((color) => color.value === attributes.overlay_color).id
              : 1,
          platform: platforms.find((platform) => platform.model === attributes.ai_engine).id,
          layout:
            attributes.layout !== null
              ? layouts.find((layout) => layout.name === attributes.layout).id
              : 1,
          background_align:
            attributes.background_align !== null ? attributes.background_align : "center",
          embedCode: `<iframe style=”width: 100%; height: 100%; border: 0px; border-radius: 10px; display: block;" src="${urlBot}"></iframe>`,
          needAuth: attributes.need_auth,
          documentOptions: attributes.manage_document,
        };

        if (attributes.background.data !== null) {
          backgroundFile = getFile(
            bot.background,
            attributes.background.data.attributes.ext,
            attributes.background.data.attributes.name,
          );
        }

        if (attributes.bot_logo.data !== null) {
          logoFile = getFile(
            bot.bot_logo,
            attributes.bot_logo.data.attributes.ext,
            attributes.bot_logo.data.attributes.name,
          );
        }

        if (attributes.user_profile.data !== null) {
          userFile = getFile(
            bot.user_profile,
            attributes.user_profile.data.attributes.ext,
            attributes.user_profile.data.attributes.name,
          );
        }

        if (attributes.bot_profile.data !== null) {
          botFile = getFile(
            bot.bot_profile,
            attributes.bot_profile.data.attributes.ext,
            attributes.bot_profile.data.attributes.name,
          );
        }

        ready = true;
        triggerLoading(false);
      } catch (error) {
        console.log(error);
      }
    }
  }

  function getFile(base64, ext, name) {
    let fileType;
    if (ext === ".png") {
      fileType = "image/png";
    } else if (ext === ".jpg") {
      fileType = "image/jpg";
    } else if (ext === ".jpeg") {
      fileType = "image/jpeg";
    } else if (ext === ".svg") {
      fileType = "image/svg+xml";
    } else {
      fileType = "image/png";
    }

    let blob = base64toBlob(base64, fileType);
    let fileImageResult = URL.createObjectURL(blob);
    return {
      imageResult: fileImageResult,
      name: name,
      base64: base64,
    };
  }

  function base64toBlob(base64Data, fileType) {
    const byteCharacters = atob(base64Data);
    const byteArrays = [];

    for (let offset = 0; offset < byteCharacters.length; offset += 512) {
      const slice = byteCharacters.slice(offset, offset + 512);

      const byteNumbers = new Array(slice.length);
      for (let i = 0; i < slice.length; i++) {
        byteNumbers[i] = slice.charCodeAt(i);
      }

      const byteArray = new Uint8Array(byteNumbers);
      byteArrays.push(byteArray);
    }

    const blob = new Blob(byteArrays, { type: fileType });
    return blob;
  }

  async function updateBot() {
    let token = localStorage.getItem("token");
    if (token) {
      triggerLoading(true);
      let allCollection = "";

      let pageNow = pagination.page;
      if (tmpDatasChecked[pageNow] === undefined) {
        tmpDatasChecked = {
          ...tmpDatasChecked,
          [pageNow]: datasChecked,
        };
      } else {
        tmpDatasChecked[pageNow] = datasChecked;
      }

      // Iterasi key dalam tmpDatasChecked
      Object.values(tmpDatasChecked).forEach((pageData) => {
        if (!pageData) {
          triggerLoading(false);
          return;
        } // Skip jika pageData undefined/null

        console.log("pageData", pageData);

        const collections = pageData
          .flatMap((item) => {
            // Ambil ID dari sub_collection jika checked = true
            const subCollectionIds =
              item.sub_collection?.filter((sub) => sub.checked).map((sub) => sub.id) || [];

            // Tambahkan item.id jika checked = true
            return item.checked ? [item.id, ...subCollectionIds] : subCollectionIds;
          })
          .join(","); // Gabungkan ID dengan koma

        console.log("collections check", collections);

        if (collections) {
          allCollection += collections + ",";
        }
      });

      // Hapus koma terakhir jika ada
      allCollection = allCollection.replace(/,$/, "");
      try {
        let response = await $call.call("client-dashboard", "update-bot", {
          URLQuery: {},
          URLParams: {
            id: botId,
          },
          Params: {
            additional_prompt: users.additionalPrompts,
            ai_engine: platforms.find((platform) => platform.id === users.platform).model,
            background: JSON.stringify(backgroundFile),
            background_align: users.background_align,
            bot_logo: JSON.stringify(logoFile),
            bot_name: users.botName,
            bot_profile: JSON.stringify(botFile),
            color_theme: themeColors.find((color) => color.id === users.primaryBrandColor).value,
            description: users.description,
            greeting_text: users.welcomeMessage,
            layout: layouts.find((layout) => layout.id === users.layout).name,
            logo_name: users.logoName,
            logo_only: users.logoOnly,
            overlay_color: themeColors.find((color) => color.id === users.overlayColor).value,
            user_profile: JSON.stringify(userFile),
            deployed_platform: users.deployedPlatform,
            username: users.username,
            need_auth: users.needAuth,
            collection_id: checkAll ? "all" : allCollection,
            library_id: selectedCollection,
            manage_document: users.documentOptions,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("update bot", response);
        onEdit = false;
        triggerLoading(false);
        triggerToast("updated");
      } catch (error) {
        triggerLoading(false);
        console.log(error);
      }
    }
  }

  async function deleteBot(id) {
    let token = localStorage.getItem("token");

    if (token) {
      triggerLoading(true);
      try {
        const response = await $call.call("client-dashboard", "delete-bot", {
          URLQuery: {},
          URLParams: {
            id: id.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          console.log(response);
          return;
        }

        // popupModal = false;
        triggerLoading(false);
        triggerToast("deleted");
        goto(await $call.getPage("client-dashboard", "bots-list"));
      } catch (error) {
        triggerLoading(false);
        console.log(error);
      }
    }
  }

  async function getDocumentsByCategory() {}

  const handleSearchCategory = async (value, category) => {
    searchCategory[category] = true;
    configsCategories[category].pagination.page = 1;
    searchTermCategory = { ...searchTermCategory, [category]: value };

    await getDocument(category);
  };

  async function nextByCategory(category) {
    if (
      configsCategories[category].pagination.page >=
      configsCategories[category].pagination.pageCount
    )
      return;
    configsCategories[category].pagination.page++;
    getDocument(category);
  }

  async function previousByCategory(category) {
    if (configsCategories[category].pagination.page <= 1) return;
    configsCategories[category].pagination.page--;
    getDocument(category);
  }

  async function handleCheckAllCategories(category, allChecked) {
    configsCategories[category].actions.find((config) => {
      if (config.id === "delete") {
        config.disabled = allChecked
          ? datasCategoriesChecked[category].find((data) => data.checked)
          : false;
      }
    });
    datasCategoriesChecked[category] = datasCategoriesChecked[category].map((data) => {
      return { id: data.id, checked: allChecked ? false : true };
    });
  }

  async function handleCheckCategories(category, id) {
    let datasCategoriesCheckedCopy = datasCategoriesChecked[category];
    datasCategoriesCheckedCopy = datasCategoriesCheckedCopy.map((data) => {
      if (data.id === id) {
        if (!data.checked) {
          if (
            configsCategories[category]?.actions &&
            configsCategories[category].headers === category
          ) {
            configsCategories[category].actions.find((config) => {
              if (config.id === "delete") {
                config.disabled = datasCategoriesChecked[category].find((data) => data.checked);
              }
            });
          }
        } else {
          if (
            configsCategories[category]?.actions &&
            configsCategories[category].headers === category
          ) {
            configsCategories[category].actions.find((config) => {
              if (config.id === "delete") {
                config.disabled = true;
              }
            });
          }
        }
        return { id: data.id, checked: !data.checked };
      }
      return data;
    });
    datasCategoriesChecked[category] = datasCategoriesCheckedCopy;
  }

  onMount(() => {
    if (edit === "true") onEdit = true;
    triggerLoading(true);
    getCollections();
    // getCollectionSub();

    const unsubscribe = page.subscribe(async ($page) => {
      botId = $page.url.searchParams.get("botId");
      if (!botId) {
        triggerToast("error");
        triggerLoading(false);
        goto(await $call.getPage("client-dashboard", "bots-list"));
        return;
      }
      await getModels();
      getBot();

      triggerLoading(false);
    });
    return () => unsubscribe();
  });

  $effect(() => {
    if (selectedCollection) {
      if (typeof selectedCollection === "number" && selectedCollection !== oldSelectedCollection) {
        oldSelectedCollection = selectedCollection;
        tmpDatasChecked = [];
        getCollectionSub(selectedCollection);
      }
    }
  });

  // $: disabledButton = users.botName === "" || users.username === "";

  $effect(() => {
    if (disabledButton) {
      disabledButton = users.botName === "" || users.username === "";
    }
  });
</script>

<!-- <svelte:head>
  {@html github}
</svelte:head> -->

<!-- <svelte:window bind:innerWidth={screenWidth} /> -->

<!-- <LoadingSpinner {isLoading} /> -->
{#if ready}
  <div class="m-4 pb-12">
    <div class="sticky top-0 bg-primary-light dark:bg-primary-dark">
      <div class="flex justify-between items-center">
        <div class="flex items-center ms-4">
          <Avatar
            class="size-10 md:size-14 border-[1px] border-hc-graylight-200"
            src={botFile.name ? botFile.imageResult : ""}
            stacked
          />
          <div class="ml-2 md:ml-4">
            <div class="flex items-center gap-2">
              <Typography
                preset="Display sm/Semibold"
                text={users.botName}
                color="primary"
                propClass="fluid-xl"
              />
              <Badge type={classStatusType(users.deployedPlatform)} text={users.deployedPlatform} />
            </div>
            <Typography
              preset="Text md/Medium"
              text={formattedDate}
              color="quarterary"
              propClass="fluid-sm md:fluid-base"
            />
          </div>
        </div>
        <div class="flex gap-2">
          {#if onEdit}
            <Button
              onclick={() => {
                onEdit = false;
                getBot();
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button on:click={() => updateBot()} class="px-3 focus-within:!ring-1"
              >{$_("save")}</Button
            >
          {:else}
            <div class="hidden md:flex gap-2">
              <Button
                onclick={async function () {
                  if (isFrontend) {
                    goto("/bots");
                  } else {
                    goto(await $call.getPage("client-dashboard", "bots"));
                  }
                }}
                color="light"
                class="px-3 focus-within:!ring-1">{$_("back")}</Button
              >
              <Button
                onclick={() => (modalDelete = true)}
                outline
                color="red"
                class="px-3 focus-within:!ring-1">{$_("delete")}</Button
              >
              <Button
                onclick={() => {
                  onEdit = true;
                }}
                class="px-3 focus-within:!ring-1">{$_("edit")}</Button
              >
            </div>
            <div onfocusout={handleDropdownFocusLoss} class="flex md:hidden items-center relative">
              <Button
                on:click={() => (dropdownOpen = !dropdownOpen)}
                class="focus-within:ring-0 !p-0"
                color="none"
                size="xs"
                ><DotsVerticalIcon
                  class="text-fg-quinary-light dark:text-fg-quinary-dark"
                /></Button
              >
              <div
                class={`z-10 ${dropdownOpen ? "absolute top-0 right-6" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
              >
                <ul class="py-2 grow text-sm">
                  <button
                    class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                    onclick={() => {
                      (onEdit = true), (dropdownOpen = !dropdownOpen);
                    }}>{$_("edit")}</button
                  >
                  <button
                    class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                    onclick={() => {
                      (modalDelete = true), (dropdownOpen = !dropdownOpen);
                    }}>{$_("delete")}</button
                  >
                </ul>
              </div>
            </div>
          {/if}
        </div>
      </div>

      <!-- Submenu -->
      <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />
    </div>

    <!-- Content -->
    {#if tabDatas.selected === "Information"}
      <div class="mt-6 space-y-6">
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text={$_("botsInformation")}
            color="secondary"
            propClass="w-full md:w-1/4"
          />

          <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-3">
            <div class="flex flex-col gap-1.5 col-span-4">
              <Typography preset="Text sm/Medium" text={$_("botsName")} color="secondary" />
              <InputRw bind:value={users.botName} {onEdit} />
            </div>
            <div class="flex flex-col gap-1.5 col-span-4">
              <Typography
                preset="Text sm/Medium"
                text={$_("descriptionOptional")}
                color="secondary"
              />
              <TextareaRw bind:value={users.description} {onEdit} />
            </div>
          </div>
        </div>
        <Divider />
        <div class="flex flex-col md:flex-row gap-1.5">
          <Typography
            preset="Text md/Semibold"
            text={$_("aiEngines")}
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4">
            <SelectedCard
              selectionData={platforms}
              bind:selectedData={users.platform}
              numGrid={3}
              isRadio
              disabled={!onEdit ||
                platforms?.find((platform) => platform.id === users?.platform)?.model === "openai"}
            />
          </div>
        </div>
        <!-- <Divider />
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text={$_("useCase")}
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-auto">
            <Dropdown
              dropdownOpen={dropdownUsecaseOpen}
              bind:selectedItem={usecaseSelected}
              items={usecase}
              isOverlay={true}
            />
          </div>
        </div> -->
        {#if platforms?.find((platform) => platform.id === users?.platform)?.model !== "openai"}
          <Divider />
          <div class="md:flex">
            <Typography
              preset="Text md/Semibold"
              text={$_("knowledgeSource")}
              color="secondary"
              propClass="w-full md:w-1/4"
            />
            <div class="w-full md:w-3/4">
              <div class="flex items-center gap-4">
                <Button
                  on:click={() => (modalOpen = true)}
                  disabled={!onEdit}
                  color="light"
                  class={`w-max rounded-lg border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1.5 flex justify-between items-center py-2.5 px-3`}
                >
                  <Typography
                    preset="Text md/Semibold"
                    text={selectedCollection
                      ? collections.find((item) => item.value === selectedCollection)?.name
                      : $_("configure")}
                    color="secondary"
                  />
                  <ChevronDownOutline
                    class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
                  />
                </Button>
              </div>
            </div>
          </div>
        {/if}
        <Divider />
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text={$_("additionalPrompts")}
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4 flex flex-col gap-1.5">
            <Typography preset="Text sm/Medium" text={$_("botsRole")} color="secondary" />
            <TextareaRw bind:value={users.additionalPrompts} {onEdit} rows={8} />
          </div>
        </div>
        <Divider />
        <!-- <div class="flex justify-end gap-3">
          <Button on:click={() => {}} color="light" class="px-3 focus-within:!ring-1"
            >{$_("cancel")}</Button
          >
          <Button
            disabled={disabledButton}
            on:click={() => {
              updateBot();
            }}
            class="px-3 focus-within:!ring-1">{$_("save")}</Button
          >
        </div> -->
      </div>
    {:else if tabDatas.selected === "Customization"}
      <div class="mt-6 space-y-6">
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text={$_("general")}
            color="secondary"
            propClass="w-full md:w-1/4"
          />

          <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="flex flex-col gap-1.5 col-span-4">
              <Typography preset="Text sm/Medium" text={$_("welcomeMessage")} color="secondary" />
              <TextareaRw bind:value={users.welcomeMessage} {onEdit} />
            </div>
            <div class="flex gap-4 col-span-4">
              <div class="flex flex-col gap-1.5 grow">
                <Typography preset="Text sm/Medium" text={$_("botsName")} color="secondary" />
                <InputRw bind:value={users.botName} {onEdit} />
              </div>
              <div class="flex flex-col gap-1.5 grow">
                <Typography preset="Text sm/Medium" text={$_("userName")} color="secondary" />
                <InputRw bind:value={users.username} {onEdit} />
              </div>
            </div>
            <div class="col-span-4 md:col-span-2">
              <UploadFile
                title={$_("botsProfilePicture")}
                handleFileUploadWrapper={() => handleInputFile("bot")}
                isPreview
                srcAvatar={botFile.name ? botFile.imageResult : ""}
                disabled={!onEdit}
              />
            </div>
            <div class="col-span-4 md:col-span-2">
              <UploadFile
                title={$_("userprofilePicture")}
                handleFileUploadWrapper={() => handleInputFile("user")}
                isPreview
                srcAvatar={userFile.name ? userFile.imageResult : ""}
                disabled={!onEdit}
              />
            </div>
          </div>
        </div>
        <Divider />
        <div class="flex flex-col md:flex-row gap-1.5">
          <Typography
            preset="Text md/Semibold"
            text={$_("logo")}
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-2 md:gap-4">
            <div class="col-span-4">
              <Toggle {onEdit} bind:checked={users.logoOnly} text="logoOnly" />
            </div>
            <div class="col-span-4 md:col-span-2">
              <UploadFile
                title={$_("logo")}
                handleFileUploadWrapper={() => handleInputFile("logo")}
                isPreview
                srcAvatar={logoFile.name ? logoFile.imageResult : ""}
                disabled={!onEdit}
              />
            </div>
            <div class="flex flex-col gap-1.5 col-span-4 md:col-span-2">
              <Typography preset="Text sm/Medium" text={$_("logoName")} color="secondary" />
              <InputRw bind:value={users.logoName} {onEdit} />
            </div>
          </div>
        </div>
        <Divider />
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text={$_("themeColors")}
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="flex flex-col gap-1.5 col-span-2">
              <Typography
                preset="Text sm/Medium"
                text={$_("primaryBrandColor")}
                color="secondary"
              />
              <Button
                on:click={!onEdit
                  ? () => {}
                  : () => {
                      tmpThemeColor = users.primaryBrandColor;
                      modalThemeColor = true;
                    }}
                class={`flex justify-between ${!onEdit ? "cursor-default focus-within:ring-0 hover:bg-primary-light hover:dark:bg-secondary-dark" : ""} border-primary-light dark:border-primary-dark bg-primary-light dark:bg-secondary-dark`}
                color="light"
              >
                <div class="flex items-center gap-2">
                  <Indicator color={getBrandColorName(users.primaryBrandColor)} />
                  {getBrandColorName(users.primaryBrandColor, "name")}
                </div>
                <ChevronDownOutline
                  class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
                />
              </Button>
            </div>
          </div>
        </div>
        <Divider />
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text={$_("layout")}
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4">
            <SelectedCard
              selectionData={layouts}
              bind:selectedData={users.layout}
              isFlex
              isRadio
              widthFlex="w-full"
              spacing="!p-4"
              disabled={!onEdit}
            />
          </div>
        </div>
        {#if users.layout === 2}
          <Divider />
          <div class="md:flex">
            <Typography
              preset="Text md/Semibold"
              text={$_("background")}
              color="secondary"
              propClass="w-full md:w-1/4"
            />
            <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
              <div class="col-span-4">
                <UploadFile
                  title={$_("backgroundPicture")}
                  handleFileUploadWrapper={() => handleInputFile("background")}
                  isPreview
                  largePreview
                  srcAvatar={backgroundFile.name ? backgroundFile.imageResult : ""}
                  disabled={!onEdit}
                />
              </div>
              <div class="flex flex-col gap-1.5 col-span-2">
                <Typography preset="Text sm/Medium" text={$_("overlayColor")} color="secondary" />
                <Button
                  on:click={!onEdit
                    ? () => {}
                    : () => {
                        tmpOverlayColor = users.overlayColor;
                        modalOverlayColor = true;
                      }}
                  class={`flex justify-between ${!onEdit ? "cursor-default focus-within:ring-0 hover:bg-primary-light hover:dark:bg-secondary-dark" : ""} border-primary-light dark:border-primary-dark bg-primary-light dark:bg-secondary-dark`}
                  color="light"
                >
                  <div class="flex items-center gap-2">
                    <Indicator color={getBrandColorName(users.overlayColor)} />
                    {getBrandColorName(users.overlayColor, "name")}
                  </div>
                  <ChevronDownOutline
                    class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
                  />
                </Button>
              </div>
              <div class="flex flex-col gap-1.5 col-span-2">
                <Typography
                  preset="Text sm/Medium"
                  text={$_("backgroundAlign")}
                  color="secondary"
                />
                <Dropdown
                  dropdownOpen={dropdownBgAlignOpen}
                  bind:selectedItem={users.background_align}
                  items={bgAlign}
                  isOverlay={true}
                  disabled={!onEdit}
                />
              </div>
            </div>
          </div>
        {/if}
        <Divider />
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text={$_("additionalSettings")}
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="col-span-4">
              <Toggle {onEdit} bind:checked={users.needAuth} text="requiresLogin" />
            </div>
            <div class="col-span-4">
              <Toggle {onEdit} bind:checked={users.documentOptions} text="documentOptions" />
            </div>
          </div>
        </div>
        <Divider />
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text={$_("preview")}
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="col-span-4">
              <div
                class="flex flex-col w-full border border-secondary-light dark:border-secondary-dark rounded-lg"
              >
                <div
                  class="flex gap-4 {backgroundFile.name
                    ? 'flex-col justify-center items-center mb-4'
                    : 'items-center px-4 pt-4'}"
                >
                  {#if backgroundFile.name}
                    <div class="relative overflow-hidden max-h-[20rem]">
                      <img
                        src={backgroundFile.imageResult}
                        alt="Background"
                        class="w-full h-full object-cover object-center rounded-lg"
                      />
                    </div>
                  {/if}
                  <div
                    class={backgroundFile.name
                      ? "mt-[-2.5rem] z-10 bg-primary-light dark:bg-primary-dark rounded-full p-2"
                      : ""}
                  >
                    {#if logoFile.name}
                      <Avatar src={logoFile.imageResult} />
                    {:else}
                      <svelte:component
                        this={getIcons("brand")}
                        svgClass="w-5 h-5 text-gray-700 dark:text-gray-400"
                      />
                    {/if}
                  </div>
                  {#if users.logoName && !users.logoOnly}
                    <Typography preset="Text xl/Semibold" text={users.logoName} />
                  {/if}
                </div>
                <div class="flex flex-col items-center gap-4 px-4 pb-4">
                  <div class="flex w-1/2 justify-end">
                    <Button outline color="alternative" class="mt-4 w-fit ">
                      <svelte:component
                        this={getIcons("refreshCw05")}
                        svgClass="w-4 h-4 hover:text-white"
                      />
                      &nbsp;
                      {$_("clearConversations")}
                    </Button>
                  </div>
                  <div class="flex flex-col w-1/2 gap-4">
                    {#if users?.welcomeMessage}
                      <div
                        class="w-full px-10 py-6 border rounded-lg border-secondary-light dark:border-secondary-dark bg-secondary-light dark:bg-secondary-dark"
                      >
                        <div class="flex items-center justify-between">
                          <div class="flex items-center gap-4">
                            <Avatar src={botFile.name ? botFile.imageResult : ""} />
                            <div class="flex gap-1 items-center">
                              <Typography
                                preset="Text md/Semibold"
                                text={users?.botName ? users.botName : "Bot"}
                              />
                              <Typography
                                preset="Text sm/regular"
                                color="quarterary"
                                text={$_("24/04/2024 03:15 PM")}
                              />
                            </div>
                          </div>
                          <svelte:component
                            this={getIcons("thumbsDown")}
                            svgClass="w-5 h-5 text-gray-700 dark:text-gray-400"
                          />
                        </div>
                        <div class="mt-4">
                          <Typography
                            preset="Text md/Regular"
                            color="tertiary"
                            text={users.welcomeMessage}
                          />
                        </div>
                      </div>
                    {/if}
                    <div
                      class="w-full px-10 py-6 border rounded-lg bg-secondary-light dark:bg-secondary-dark {users.primaryBrandColor ===
                      1
                        ? 'border-brand-light'
                        : users.primaryBrandColor === 2
                          ? 'border-blue-500'
                          : users.primaryBrandColor === 3
                            ? 'border-green-500'
                            : 'border-brand-light'}"
                    >
                      <div class="flex items-center justify-between">
                        <div class="flex items-center gap-4">
                          <Avatar src={userFile.name ? userFile.imageResult : ""} />
                          <div class="flex gap-1 items-center">
                            <Typography
                              preset="Text md/Semibold"
                              text={users?.username ? users.username : "User"}
                            />
                            <Typography
                              preset="Text sm/regular"
                              color="quarterary"
                              text={$_("24/04/2024 03:15 PM")}
                            />
                          </div>
                        </div>
                        <svelte:component
                          this={getIcons("thumbsDown")}
                          svgClass="w-5 h-5 text-gray-700 dark:text-gray-400"
                        />
                      </div>
                      <div class="mt-4">
                        <Typography
                          preset="Text md/Regular"
                          color="tertiary"
                          text={$_(
                            "Hello! Can you explain the meaning of Surah Al-Fatihah, verse 1?",
                          )}
                        />
                      </div>
                    </div>
                    <div
                      class="w-full px-10 py-6 border rounded-lg border-secondary-light dark:border-secondary-dark bg-secondary-light dark:bg-secondary-dark"
                    >
                      <div class="flex items-center justify-between">
                        <div class="flex items-center gap-4">
                          <Avatar src={botFile.name ? botFile.imageResult : ""} />
                          <div class="flex gap-1 items-center">
                            <Typography
                              preset="Text md/Semibold"
                              text={users?.botName ? users.botName : "Bot"}
                            />
                            <Typography
                              preset="Text sm/regular"
                              color="quarterary"
                              text={$_("24/04/2024 03:15 PM")}
                            />
                          </div>
                        </div>
                        <svelte:component
                          this={getIcons("thumbsDown")}
                          svgClass="w-5 h-5 text-gray-700 dark:text-gray-400"
                        />
                      </div>
                      <div class="mt-4">
                        <Typography
                          preset="Text md/Regular"
                          color="tertiary"
                          text={$_(
                            "Certainly! Surah Al-Fatihah, verse 1, is known as the Opener of the Book. It begins with ‘Al-Hamdu lillahi Rabbil-`Alamin,’ which means ‘All praise is due to Allah, the Lord of existence.’ This verse emphasizes praising Allah and recognizing His sovereignty. It’s also called Umm Al-Kitab (the Mother of the Book) because it contains the essence of the entire Quran. Would you like more details?",
                          )}
                        />
                      </div>
                    </div>
                    <div class="w-full flex items-center gap-4">
                      <div
                        class="rounded-2xl border-2 p-4 grow {users.primaryBrandColor === 1
                          ? 'border-brand-light'
                          : users.primaryBrandColor === 2
                            ? 'border-blue-500'
                            : users.primaryBrandColor === 3
                              ? 'border-green-500'
                              : 'border-brand-light'}"
                      >
                        <Typography preset="Text md/regular" text={$_("Send a message...")} />
                      </div>
                      <div
                        class="flex justify-center items-center p-3 rounded-full {users.primaryBrandColor ===
                        1
                          ? 'bg-brand-primary-dark'
                          : users.primaryBrandColor === 2
                            ? 'bg-blue-500'
                            : users.primaryBrandColor === 3
                              ? 'bg-green-500'
                              : 'bg-brand-primary-dark'}"
                      >
                        <svelte:component this={getIcons("send03")} svgClass="w-5 h-5 text-white" />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    {:else if tabDatas.selected === "Data"}
      {#each datas as data, index}
        {#if data[collectionId[index]].length > 0}
          <div class="md:flex">
            <div class="w-full md:w-1/4 mt-6">
              <Typography
                preset="Text md/Semibold"
                text={configsCategories[collectionId[index]].name}
                color="secondary"
                propClass=""
              />
              <button
                onclick={async function () {
                  goto(
                    (await $call.getPage("client-dashboard", "document-library-documents")) +
                      `?collectionId=${collectionId[index]}`,
                  );
                  return;
                }}
                class="text-sm text-brand-tertiary-light dark:text-brand-tertiary-dark"
              >
                {$_("goToDocumentCollection")}
              </button>
            </div>
            <div class="w-full md:w-3/4">
              {#if data[collectionId[index]].length > 0 || searchCategory[collectionId[index]]}
                <DataTableListCollection
                  bind:datas={data[collectionId[index]]}
                  bind:datasChecked={datasCategoriesChecked[collectionId[index]]}
                  bind:configs={configsCategories[collectionId[index]]}
                  bind:gettingDatas
                  bind:next={nextByCategory}
                  bind:previous={previousByCategory}
                  bind:handleCheckAll={handleCheckAllCategories}
                  bind:handleCheck={handleCheckCategories}
                />
              {:else if gettingDatas}
                <StateCard
                  illustration={IllustrationProgress}
                  description={$_("fetchingDataPleaseWait")}
                />
              {/if}
            </div>
          </div>
        {/if}
      {/each}
    {:else if tabDatas.selected === "Widget Embed"}
      <div class="mt-6 space-y-6">
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text={$_("embedCodes")}
            color="secondary"
            propClass="w-full md:w-1/4"
          />

          <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
            <Button
              outline
              on:click={() => copyToClipboard(users.embedCode)}
              color="purple"
              class=" w-max col-span-4 hover:!bg-transparent !text-sm !font-semibold !text-utility-brand-600-light dark:!text-utility-brand-600-dark !border-brand-light dark:!border-brand-dark focus-within:!ring-1 px-3 py-[10px] gap-1"
            >
              <CopyIcon
                svgClass="w-5 h-5 text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
              />
              {$_("copyToClipboard")}
            </Button>
            <div class="col-span-4">
              <TextareaRw bind:value={users.embedCode} onEdit={false} />
            </div>
          </div>
        </div>
      </div>
    {/if}
  </div>

  <Modal title={$_("chooseYourPreferedColor")} bind:open={modalThemeColor} autoclose outsideclose>
    <SelectedCard
      selectionData={themeColors}
      bind:selectedData={tmpThemeColor}
      numGrid="2"
      isRadio
      widthFlex="w-full"
      spacing="!p-4"
      widgetIndicator={true}
    />
    <svelte:fragment slot="footer">
      <Button
        on:click={() => (users.primaryBrandColor = tmpThemeColor)}
        class="px-3 w-full focus-within:!ring-1">{$_("save")}</Button
      >
    </svelte:fragment>
  </Modal>

  <Modal title="Terms of Service" bind:open={modalOverlayColor} autoclose outsideclose>
    <SelectedCard
      selectionData={themeColors}
      bind:selectedData={tmpOverlayColor}
      numGrid="2"
      isRadio
      widthFlex="w-full"
      spacing="!p-4"
      widgetIndicator={true}
    />
    <svelte:fragment slot="footer">
      <Button
        on:click={() => (users.overlayColor = tmpOverlayColor)}
        class="px-3 w-full focus-within:!ring-1">{"Save"}</Button
      >
    </svelte:fragment>
  </Modal>

  <Modal bind:open={modalDelete} size="xs" autoclose>
    <div class="text-center">
      <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
      <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
        {$_("areYouSureYouWantToDelete")}
      </h3>
      <Button onclick={() => deleteBot(botId)} color="red" class="me-2">{$_("yesImSure")}</Button>
      <Button color="alternative">{$_("noCancel")}</Button>
    </div>
  </Modal>

  <Modal bind:open={modalOpen} outsideclose>
    <svelte:fragment slot="header">
      <div class="flex gap-4 me-6 grow w-full">
        <div
          class="grow pt-2 border-t-4 border-fg-brand-primary-alt-light dark:border-fg-brand-primary-alt-dark"
        >
          <Typography
            preset="Text md/Bold"
            text={$_("chooseDocumentsLibrary")}
            color={!isSelectedLibrary ? "brand-secondary" : "secondary"}
          />
          <Typography
            preset="Text sm/Regular"
            text={$_("documentLibraryAsThePrimarySource")}
            color={!isSelectedLibrary ? "brand-tertiary" : "tertiary"}
          />
        </div>

        <div
          class={`grow pt-2 border-t-4 ${isSelectedLibrary ? "border-fg-brand-primary-alt-light dark:border-fg-brand-primary-alt-dark" : "border-utility-gray-200-light dark:border-utility-gray-200-dark"}`}
        >
          <Typography
            preset="Text md/Bold"
            text={$_("chooseWhichCollections")}
            color={isSelectedLibrary ? "brand-secondary" : "secondary"}
          />
          <Typography
            preset="Text sm/Regular"
            text={$_("chooseTheOneYouWantToInclude")}
            color={isSelectedLibrary ? "brand-tertiary" : "tertiary"}
          />
        </div>
      </div>
    </svelte:fragment>
    {#if !isSelectedLibrary}
      <div class="flex flex-col justify-center gap-4">
        <Typography preset="Text sm/Medium" text={$_("latestDocumentLibrary")} color="quaterary" />

        <SelectedCard
          selectionData={collections}
          bind:selectedData={selectedCollection}
          numGrid={2}
          isRadio
        />
      </div>
    {:else if allData.length > 0 || searchTerm !== ""}
      {#if !gettingDatas}
        <DataTableListCollection
          bind:datas={allData}
          bind:datasChecked
          bind:configs={configsChecked}
          bind:gettingDatas
          bind:next
          bind:previous
          bind:handleCheckAll
          bind:handleCheck
        />
      {:else}
        <StateCard illustration={IllustrationProgress} description={$_("fetchingDataPleaseWait")} />
      {/if}
    {:else if gettingDatas}
      <StateCard illustration={IllustrationProgress} description={$_("fetchingDataPleaseWait")} />
    {:else}
      <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
    {/if}
    <svelte:fragment slot="footer">
      {#if !isSelectedLibrary}
        <Button
          on:click={() => {
            isSelectedLibrary = true;
          }}
          class="w-full"
          disabled={!selectedCollection}>{$_("next")}</Button
        >
      {:else}
        <div class="flex gap-4 w-full">
          <button
            onclick={() => {
              isSelectedLibrary = false;
            }}
            class="px-5">{$_("goBack")}</button
          >
          <Button
            on:click={() => {
              modalOpen = false;
            }}
            class="grow"
            disabled={areAllUnchecked()}>{$_("finish")}</Button
          >
        </div>
      {/if}
    </svelte:fragment>
  </Modal>
{/if}
