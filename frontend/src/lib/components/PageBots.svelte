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
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";
  import { Avatar, Button, Modal, Indicator } from "flowbite-svelte";
  import { ChevronDownOutline, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import {
    Badge,
    Toggle,
    Divider,
    InputRw,
    Dropdown,
    StateCard,
    TabDetail,
    TextareaRw,
    Typography,
    UploadFile,
    InputSlider,
    SelectedCard,
    DataTableListCollection,
  } from "$lib";
  import {
    call,
    getIcons,
    triggerToast,
    triggerLoading,
    findIconFromModelName,
    getCollectionAPIConnector,
  } from "$lib";

  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";

  // const PlusIcon = getIcons("plus");
  const CopyIcon = getIcons("copy03");
  // const Trash01Icon = getIcons("trash01");
  // const BrandIcon = getIcons("brand");
  // const Send03Icon = getIcons("send03");
  // const ThumbsDownIcon = getIcons("thumbsDown");
  // const RefreshCw05Icon = getIcons("refreshCw05");
  const VerifiedTickIcon = getIcons("verifiedTick");
  const ApiConnectorSupportIcon = getIcons("apiConnectorSupport");

  const isFrontend = process.env.NODE_ENV === "development";

  const params = new URLSearchParams(window.location.search);
  const edit = params.get("edit");
  const botId = params.get("botId");

  let { pageAction } = $props();

  let formattedDate = "Created on 2021-08-12 14:00 PM";

  let onEdit = $state(false);
  let search = $state(false);
  let needAuth = $state(false);
  let checkAll = $state(false);
  let modalAPI = $state(false);
  let modalDelete = $state(false);
  let onMountState = $state(true);
  let modalKnowledge = $state(false);
  let modalThemeColor = $state(false);
  let apiSourceEnable = $state(false);
  let modalOverlayColor = $state(false);
  let isSelectedLibrary = $state(false);
  let dropdownBgAlignOpen = $state(false);
  let gettingDatasKnowledge = $state(false);
  let screenWidth = $state(0);
  let selectedplatforms = $state(0);
  let botName = $state("");
  let searchTerm = $state("");
  let description = $state("");
  let tmpThemeColor = $state("");
  let tmpOverlayColor = $state("");
  let additionalPrompt = $state("");
  let selectedCollectionId = $state("");
  let tmpSelectedCollectionId = $state("");
  let fileImage = $state(null);
  let collectionId = $state(null);
  let fileImageName = $state(null);
  let fileImageResult = $state(null);
  let oldSelectedCollectionId = $state(null);
  let datas = $state([]);
  let allData = $state([]);
  let platforms = $state([]);
  let collections = $state([]);
  let paginationDocs = $state([]);
  let apiCollections = $state([]);
  let recommendations = $state([]);
  let initCollections = $state([]);
  let datasCheckedApi = $state([]);
  let collectionsChecked = $state([]);
  let datasCheckedKnowledge = $state([]);
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
  let bgAlign = $state([
    { name: "Center", value: "center" },
    { name: "Left", value: "left" },
    { name: "Right", value: "right" },
  ]);
  let users = $state({});
  let configsApi = $state({});
  let searchCategory = $state({});
  let tmpDatasChecked = $state({});
  let configsKnowledge = $state({});
  let configsCategories = $state({});
  let searchTermCategory = $state({});
  let datasCategoriesChecked = $state({});
  let botFile = $state({ name: "" });
  let userFile = $state({ name: "" });
  let logoFile = $state({ name: "" });
  let backgroundFile = $state({ name: "" });
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });
  let paginationCollection = $state({ page: 1, pageSize: 20, pageCount: 0, total: 0 });
  let paginationApi = $state({ page: 1, pageSize: 1000, pageCount: 0, total: 0 });
  let tabDatas = $state({
    selected: "Information",
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

  const getLibraries = async () => {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let page = 1;
        initCollections = [];
        while (true) {
          let response = await $call.call("client-dashboard", "get-librarys", {
            URLQuery: {},
            URLParams: {},
            Params: {
              page: page,
              page_size: 25,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (
            response === null ||
            response === undefined ||
            (response?.status && response.status !== 200)
          ) {
            errorResponseTrigger(error);
            return;
          }

          if (response.result.data.length > 0) {
            initCollections = initCollections.concat(response.result.data);
            page++;
          } else {
            break;
          }
        }

        // console.log("libraries", initCollections);

        collections = initCollections.map((collection) => {
          return {
            value: collection.id,
            id: collection.id,
            name: collection.attributes.name,
            description: collection.attributes.meta?.description,
          };
        });
      } catch (error) {
        errorResponseTrigger(error);
      }
    }
  };

  const getCollections = async (id) => {
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        if (!search) {
          gettingDatasKnowledge = true;
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

        if (
          response === null ||
          response === undefined ||
          (response?.status && response.status !== 200)
        ) {
          errorResponseTrigger(response);
          gettingDatasKnowledge = false;
          return;
        }

        // console.log("collections", response);

        if (response.result.data.length > 0) {
          pagination = response.result.meta.pagination;
          configsKnowledge = {
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
                  handleCheckAllCollection(false, "knowledge");
                },
              },
              {
                id: "select_all",
                name: "Select All",
                type: "brand-inverse",
                onclick: () => {
                  handleCheckAllCollection(true, "knowledge");
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
          // console.log("all data with subcollections", allData);

          if (pageAction === "create") {
            if (checkAll) {
              datasCheckedKnowledge = mapAllDataChecked(allData, true);
            } else {
              if (tmpDatasChecked[pagination.page] === undefined) {
                datasCheckedKnowledge = mapAllDataChecked(allData, false);
              } else {
                datasCheckedKnowledge = tmpDatasChecked[pagination.page];
              }
            }
          } else if (pageAction === "detail") {
            datasCheckedKnowledge = mapAllDataChecked(allData, false);

            if (checkAll) {
              datasCheckedKnowledge = mapAllDataChecked(allData, true);
            } else {
              if (tmpDatasChecked[pagination.page] === undefined) {
                if (selectedCollectionId === tmpSelectedCollectionId) {
                  tmpSelectedCollectionId = null;
                  datasCheckedKnowledge = datasCheckedKnowledge.map((data) => {
                    if (collectionsChecked.includes(data.id)) {
                      data.checked = true;
                    }
                    data.sub_collection = data.sub_collection.map((sub) => {
                      if (collectionsChecked.includes(sub.id)) {
                        sub.checked = true;
                      }
                      return sub;
                    });
                    return data;
                  });
                } else {
                  datasCheckedKnowledge = allData.map((data) => ({
                    id: data.id,
                    checked: false,
                  }));
                }
              } else {
                datasCheckedKnowledge = tmpDatasChecked[pagination.page];
              }
            }
          }
        } else {
          pagination = response.result.meta.pagination;
          allData = [];
          datasCheckedKnowledge = [];
          configsKnowledge = {
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
                  handleCheckAllCollection(false, "knowledge");
                },
              },
              {
                id: "select_all",
                name: "Select All",
                type: "brand-inverse",
                onclick: () => {
                  handleCheckAllCollection(true, "knowledge");
                },
              },
            ],
          };
        }
      } catch (error) {
        errorResponseTrigger(error);
      } finally {
        gettingDatasKnowledge = false;
        search = false;
      }
    }
  };

  const getSubCollections = async (id) => {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        gettingDatasKnowledge = true;
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
            gettingDatasKnowledge = false;
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
        // console.log("subCollectionList", subCollectionList);
        gettingDatasKnowledge = false;
        return subCollectionList;
      } catch (error) {
        gettingDatasKnowledge = false;
        errorResponseTrigger(error);
        return;
      }
    }
  };

  const getModels = async () => {
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

        // console.log("models", response);

        if (
          response === null ||
          response === undefined ||
          (response?.status && response.status !== 200)
        ) {
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
            return;
          } else {
            triggerLoading(false);
            errorResponseTrigger(response);
            return;
          }
        }

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
            recommended: model.name === "Llama3.1" || model.name === "Llama3.2" ? true : false,
            apiSupported: model.name === "Llama3.2" ? true : false,
          };
        });

        platforms = await Promise.all(promise);
        // console.log("platforms", platforms);
      } catch (error) {
        errorResponseTrigger(error);
      }
    }
  };

  const getDocument = async (id) => {
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

      if (
        response === null ||
        response === undefined ||
        (response?.status && response.status !== 200)
      ) {
        triggerLoading(false);
        errorResponseTrigger(response);
        gettingDatasKnowledge = false;
        return;
      }

      // console.log("ayam", response);

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

      configsCategories[id].pagination = response.result.meta.pagination;
    }
  };

  const getDocuments = async () => {
    let token = localStorage.getItem("token");
    for (let i = 0; i < collectionId.length; i++) {
      // console.log("collectionId", collectionId[i]);
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

        if (
          response === null ||
          response === undefined ||
          (response?.status && response.status !== 200)
        ) {
          triggerLoading(false);
          errorResponseTrigger(response);
          gettingDatasKnowledge = false;
          return;
        }

        // console.log("ayam", response);

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

          // Remove any existing object with the same key
          const key = collectionId[i];
          datas = [
            ...datas.filter((obj) => !Object.prototype.hasOwnProperty.call(obj, key)),
            { [key]: await Promise.all(promises) },
          ];
        } else {
          // Remove any existing object with the same key
          const key = collectionId[i];
          datas = [
            ...datas.filter((obj) => !Object.prototype.hasOwnProperty.call(obj, key)),
            { [key]: [] },
          ];
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
          gettingDatasKnowledge = false;
          return;
        }

        // console.log("cekkkkk", response2);

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
    gettingDatasKnowledge = false;
  };

  const handleGetAPICollections = async () => {
    apiCollections = await getCollectionAPIConnector($_, $call, searchTerm, paginationApi);
    if (apiCollections.length > 0) {
      apiCollections = apiCollections.map((collection) => {
        return {
          id: collection.id,
          collectionName: collection.collectionName,
        };
      });

      datasCheckedApi = apiCollections.map((data) => ({
        id: data.id,
        checked: false,
      }));

      configsApi = {
        search: true,
        searchHandler: handleSearch,
        pagination: paginationApi,
        space: "p-4",
        actionInLeft: true,
        actions: [
          {
            id: "deselect_all",
            name: "Deselect All",
            onclick: () => {
              handleCheckAllCollection(false, "api");
            },
          },
          {
            id: "select_all",
            name: "Select All",
            type: "brand-inverse",
            onclick: () => {
              handleCheckAllCollection(true, "api");
            },
          },
        ],
      };
    }
    // console.log("apiCollections", apiCollections);
  };

  const getBot = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        !onMountState && triggerLoading(true);
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

        if (
          response === null ||
          response === undefined ||
          (response?.status && response.status !== 200)
        ) {
          triggerLoading(false);
          errorResponseTrigger(response);
          goto(isFrontend ? "/bots/list" : await $call.getPage("client-dashboard", "bots-list"));
          return;
        }

        const bot = response.result;
        const attributes = bot.attributes;

        let urlBot = (await $call.getPage("hexabot", "home-page")) + "?id=" + attributes.bot_id;

        collectionsChecked = attributes.collection_id.split(",").map(Number);
        collectionId = collectionsChecked;
        tmpSelectedCollectionId = attributes?.library_id;
        selectedCollectionId = attributes?.library_id;

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
          layout:
            attributes.layout !== null
              ? layouts.find((layout) => layout.name === attributes.layout).id
              : 1,
          background_align:
            attributes.background_align !== null ? attributes.background_align : "center",
          embedCode: `<iframe style=”width: 100%; height: 100%; border: 0px; border-radius: 10px; display: block;" src="${urlBot}"></iframe>`,
          needAuth: attributes.need_auth,
          documentOptions: attributes.manage_document,
          apiSourceEnable: attributes?.function_call || false,
        };

        additionalPrompt = attributes.additional_prompt;
        selectedplatforms = platforms.find(
          (platform) => platform.model === attributes.ai_engine,
        ).id;

        // console.log("selectedplatforms", selectedplatforms);

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

        datasCheckedApi = apiCollections.map((data) => ({
          id: data.id,
          checked: attributes.collection_api_ids
            ? attributes.collection_api_ids.split(",").map(Number).includes(data.id)
            : false,
        }));
      } catch (error) {
        errorResponseTrigger(error);
      } finally {
        !onMountState && triggerLoading(false);
      }
    }
  };

  const createBot = async () => {
    let token = localStorage.getItem("token");

    if (token !== null) {
      triggerLoading(true);

      if (platforms?.find((platform) => platform.id === selectedplatforms)?.model !== "openai") {
        let allCollection = "";

        // console.log("tmpDatasChecked", tmpDatasChecked);

        let pageNow = pagination.page;
        if (tmpDatasChecked[pageNow] === undefined) {
          tmpDatasChecked = {
            ...tmpDatasChecked,
            [pageNow]: datasCheckedKnowledge,
          };
        } else {
          tmpDatasChecked[pageNow] = datasCheckedKnowledge;
        }

        Object.values(tmpDatasChecked).forEach((pageData) => {
          if (!pageData) {
            triggerLoading(false);
            return;
          }

          // console.log("pageData", pageData);

          const collections = pageData
            .flatMap((item) => {
              // Ambil ID dari sub_collection jika checked = true
              const subCollectionIds =
                item.sub_collection?.filter((sub) => sub.checked).map((sub) => sub.id) || [];

              // Tambahkan item.id jika checked = true
              return item.checked ? [item.id, ...subCollectionIds] : subCollectionIds;
            })
            .join(","); // Gabungkan ID dengan koma

          // console.log("collections check", collections);

          if (collections) {
            allCollection += collections + ",";
          }
        });

        allCollection = allCollection.replace(/,$/, "");

        // console.log("All Collection IDs (checked = true):", allCollection);
        // return;

        const apiCollectionSend = apiCollections
          .filter((collection) =>
            datasCheckedApi.some((data) => data.id === collection.id && data.checked),
          )
          .map((collection) => collection.id)
          .join(",");

        // console.log("apiCollectionSend", apiCollectionSend);
        // return;

        let dataSend = {
          bot_name: botName,
          description: description,
          ai_engine: platforms.find((platform) => platform.id === selectedplatforms).model,
          library_id: selectedCollectionId,
          collection_id: checkAll ? "all" : allCollection,
          additional_prompt: additionalPrompt,
          need_auth: needAuth,
          function_call: apiSourceEnable,
          collection_api_ids: apiCollectionSend,
        };

        if (fileImage !== "" && fileImage !== null && fileImage !== undefined) {
          dataSend.bot_logo = JSON.stringify({ name: fileImageName, base64Data: fileImage });
        }

        // console.log("dataSend", dataSend);

        try {
          const response = await $call.call("client-dashboard", "create-bot", {
            URLQuery: {},
            URLParams: {},
            Params: dataSend,
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          // console.log(response);

          if (
            response === null ||
            response === undefined ||
            (response?.status && response.status !== 200)
          ) {
            triggerLoading(false);
            errorResponseTrigger(response);
            return;
          }

          goto(isFrontend ? "/bots" : await $call.getPage("client-dashboard", "bots"));
          triggerToast("created");
        } catch (error) {
          errorResponseTrigger(error);
        } finally {
          triggerLoading(false);
        }
      } else {
        const response = await $call.call("client-dashboard", "find-create-library-llm", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (
          response === null ||
          response === undefined ||
          (response?.status && response.status !== 200)
        ) {
          triggerLoading(false);
          errorResponseTrigger(response);
          return;
        }

        let tmpCollectionId = response.result.collection_id;
        let libraryId = response.result.library_id;

        let dataSend = {
          bot_name: botName,
          description: description,
          ai_engine: platforms.find((platform) => platform.id === selectedplatforms).model,
          library_id: parseInt(libraryId),
          collection_id: tmpCollectionId.toString(),
          additional_prompt: additionalPrompt,
          need_auth: needAuth,
        };

        if (fileImage !== "" && fileImage !== null && fileImage !== undefined) {
          dataSend.bot_logo = JSON.stringify({ name: fileImageName, base64Data: fileImage });
        }

        // console.log("dataSend", dataSend);

        try {
          const response = await $call.call("client-dashboard", "create-bot", {
            URLQuery: {},
            URLParams: {},
            Params: dataSend,
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          // console.log(response);

          if (
            response === null ||
            response === undefined ||
            (response?.status && response.status !== 200)
          ) {
            triggerLoading(false);
            errorResponseTrigger(response);
            return;
          }

          goto(isFrontend ? "/bots" : await $call.getPage("client-dashboard", "bots"));
          triggerToast("created");
        } catch (error) {
          errorResponseTrigger(error);
        } finally {
          triggerLoading(false);
        }
      }
    }
  };

  const updateBot = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      triggerLoading(true);
      let allCollection = "";

      let pageNow = pagination.page;
      if (tmpDatasChecked[pageNow] === undefined) {
        tmpDatasChecked = {
          ...tmpDatasChecked,
          [pageNow]: datasCheckedKnowledge,
        };
      } else {
        tmpDatasChecked[pageNow] = datasCheckedKnowledge;
      }

      // Iterasi key dalam tmpDatasChecked
      Object.values(tmpDatasChecked).forEach((pageData) => {
        if (!pageData) {
          triggerLoading(false);
          return;
        } // Skip jika pageData undefined/null

        // console.log("pageData", pageData);

        const collections = pageData
          .flatMap((item) => {
            // Ambil ID dari sub_collection jika checked = true
            const subCollectionIds =
              item.sub_collection?.filter((sub) => sub.checked).map((sub) => sub.id) || [];

            // Tambahkan item.id jika checked = true
            return item.checked ? [item.id, ...subCollectionIds] : subCollectionIds;
          })
          .join(","); // Gabungkan ID dengan koma

        // console.log("collections check", collections);

        if (collections) {
          allCollection += collections + ",";
        }
      });

      // Hapus koma terakhir jika ada
      allCollection = allCollection.replace(/,$/, "");

      const apiCollectionSend = apiCollections
        .filter((collection) =>
          datasCheckedApi.some((data) => data.id === collection.id && data.checked),
        )
        .map((collection) => collection.id)
        .join(",");

      try {
        let response = await $call.call("client-dashboard", "update-bot", {
          URLQuery: {},
          URLParams: {
            id: botId,
          },
          Params: {
            additional_prompt: additionalPrompt,
            ai_engine: platforms.find((platform) => platform.id === selectedplatforms).model,
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
            library_id: selectedCollectionId,
            manage_document: users.documentOptions,
            function_call: users.apiSourceEnable,
            collection_api_ids: apiCollectionSend,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (
          response === null ||
          response === undefined ||
          (response?.status && response.status !== 200)
        ) {
          triggerLoading(false);
          errorResponseTrigger(response);
          return;
        }

        triggerToast("updated");
      } catch (error) {
        errorResponseTrigger(error);
      } finally {
        getBot();
      }
    }
  };

  const deleteBot = async (id) => {
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

        if (
          response === null ||
          response === undefined ||
          (response?.status && response.status !== 200)
        ) {
          triggerLoading(false);
          errorResponseTrigger(response);
          return;
        }

        triggerToast("deleted");
        goto(isFrontend ? "/bots/list" : await $call.getPage("client-dashboard", "bots-list"));
      } catch (error) {
        errorResponseTrigger(error);
      } finally {
        triggerLoading(false);
      }
    }
  };

  const handleCheck = async (id, isSubCollection = false, parentId = null, type = "knowledge") => {
    if (type === "knowledge") {
      if (!Array.isArray(datasCheckedKnowledge)) {
        datasCheckedKnowledge = [];
      }

      datasCheckedKnowledge = datasCheckedKnowledge.map((data) => {
        // Jika parentId tidak ada, kita memproses level utama
        if (!isSubCollection && data.id === id) {
          if (data.checked) {
            checkAll = false;
          }
          return {
            ...data,
            checked: !data.checked,
          };
        }

        // console.log("checkAll", checkAll);

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
          [pageNow]: datasCheckedKnowledge,
        };
      } else {
        tmpDatasChecked[pageNow] = datasCheckedKnowledge;
      }
      // console.log("datasCheckedKnowledge setelah pembaruan:", datasCheckedKnowledge);
    } else {
      if (!Array.isArray(datasCheckedApi)) {
        datasCheckedApi = [];
      }

      datasCheckedApi = datasCheckedApi.map((data) => {
        if (data.id === id) {
          if (data.checked) {
            checkAll = false;
          }
          return {
            ...data,
            checked: !data.checked,
          };
        }
        return data;
      });

      // console.log("datasCheckedApi setelah pembaruan:", datasCheckedApi);
    }
  };

  const handleCheckAll = async (allChecked, type = "knowledge") => {
    if (type === "knowledge") {
      datasCheckedKnowledge = datasCheckedKnowledge.map((data) => ({
        id: data.id,
        checked: allChecked ? false : true,
        sub_collection: data.sub_collection.map((sub) => ({
          id: sub.id,
          checked: allChecked ? false : true,
        })),
      }));
    } else if (type === "api") {
      datasCheckedApi = datasCheckedApi.map((data) => ({
        id: data.id,
        checked: allChecked ? false : true,
      }));
    }

    // console.log("Updated datasCheckedKnowledge (handleCheckAll):", datasCheckedKnowledge);
  };

  const handleCheckAllCollection = async (allChecked, type = "knowledge") => {
    if (type === "knowledge") {
      checkAll = allChecked;

      datasCheckedKnowledge = datasCheckedKnowledge.map((data) => ({
        id: data.id,
        checked: allChecked,
        sub_collection: data.sub_collection.map((sub) => ({
          id: sub.id,
          checked: allChecked,
        })),
      }));

      if (!allChecked) {
        // Jika allChecked adalah false, update tmpDatasChecked
        Object.keys(tmpDatasChecked).forEach((page) => {
          tmpDatasChecked[page] = tmpDatasChecked[page].map((data) => ({
            id: data.id,
            checked: false,
            sub_collection: data.sub_collection.map((sub) => ({
              id: sub.id,
              checked: false,
            })),
          }));
        });
      }
    } else if (type === "api") {
      checkAll = allChecked;

      datasCheckedApi = datasCheckedApi.map((data) => ({
        id: data.id,
        checked: allChecked,
      }));

      if (!allChecked) {
        // Jika allChecked adalah false, update tmpDatasChecked
        Object.keys(tmpDatasChecked).forEach((page) => {
          tmpDatasChecked[page] = tmpDatasChecked[page].map((data) => ({
            id: data.id,
            checked: false,
          }));
        });
      }
    }

    // console.log("Updated datasCheckedKnowledge (handleCheckAllCollection):", datasCheckedKnowledge);
  };

  const handleCheckCategories = async (category, id) => {
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
  };

  const handleCheckAllCategories = async (category, allChecked) => {
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
  };

  const next = async () => {
    if (pagination.page >= pagination.pageCount) return;

    let pageNow = pagination.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasCheckedKnowledge,
      };
    } else {
      tmpDatasChecked[pageNow] = datasCheckedKnowledge;
    }

    pagination.page = pagination.page + 1;
    // console.log("cek tmpDatasChecked", tmpDatasChecked);
    getCollections(selectedCollectionId);
  };

  const previous = async () => {
    if (pagination.page <= 1) return;

    let pageNow = pagination.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasCheckedKnowledge,
      };
    } else {
      tmpDatasChecked[pageNow] = datasCheckedKnowledge;
    }

    pagination.page = pagination.page - 1;
    getCollections(selectedCollectionId);
  };

  const nextByCategory = async (category) => {
    if (
      configsCategories[category].pagination.page >=
      configsCategories[category].pagination.pageCount
    )
      return;
    configsCategories[category].pagination.page++;
    getDocument(category);
  };

  const previousByCategory = async (category) => {
    if (configsCategories[category].pagination.page <= 1) return;
    configsCategories[category].pagination.page--;
    getDocument(category);
  };

  const handleSearch = async (value, type = "knowledge") => {
    if (type === "knowledge") {
      search = true;
      pagination.page = 1;
      searchTerm = value;

      await getCollections(selectedCollectionId);
    } else {
      paginationApi.page = 1;
      searchTerm = value;

      await handleGetAPICollections();
    }
  };

  const handleSearchCategory = async (value, category) => {
    searchCategory[category] = true;
    configsCategories[category].pagination.page = 1;
    searchTermCategory = { ...searchTermCategory, [category]: value };

    await getDocument(category);
  };

  const readFileAsBase64 = async (file) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        const base64Data = reader.result.split(",")[1];
        resolve(base64Data);
      };
      reader.onerror = reject;

      reader.readAsDataURL(file);
    });
  };

  const areAllUnchecked = () => {
    // Gabungkan data dari halaman sekarang dengan data yang tersimpan di tmpDatasChecked
    const allData = [...datasCheckedKnowledge, ...Object.values(tmpDatasChecked).flat()];

    // Cek apakah semua data tidak ada yang checked
    return allData.every((item) => item.checked === false);
  };

  const handleInputFile = (usage) => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = "image/jpeg, image/png";

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file && (file.type === "image/jpeg" || file.type === "image/png")) {
        if (pageAction === "create" && usage === "bot") {
          fileImageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          fileImage = await readFileAsBase64(file);
          fileImageName = file.name;
        } else if (usage === "bot") {
          botFile.imageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          botFile.base64Data = await readFileAsBase64(file);
          botFile.name = file.name;
        } else if (usage === "user") {
          userFile.imageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          userFile.base64Data = await readFileAsBase64(file);
          userFile.name = file.name;
        } else if (usage === "logo") {
          logoFile.imageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          logoFile.base64Data = await readFileAsBase64(file);
          logoFile.name = file.name;
        } else if (usage === "background") {
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

  const mapAllDataChecked = (allData, checked = false) => {
    return allData.map((data) => ({
      id: data.id,
      checked,
      sub_collection: data.sub_collections
        ? data.sub_collections.map((sub) => ({
            id: sub.id,
            checked,
          }))
        : [],
    }));
  };

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

  const getFile = (base64, ext, name) => {
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
  };

  const base64toBlob = (base64Data, fileType) => {
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
  };

  const getBrandColorName = (selectedId, value) => {
    const selectedColor = themeColors.find((color) => color.id === selectedId);
    return selectedColor
      ? value === "name"
        ? selectedColor.name
        : selectedColor.value
      : "Unknown Color";
  };

  const copyToClipboard = (text) => {
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
  };

  const errorResponseTrigger = (error) => {
    triggerToast("error");
    console.log("error", error);
  };

  onMount(() => {
    const unsubscribe = page.subscribe(async () => {
      try {
        triggerLoading(true);
        if (edit === "true" || pageAction === "create") onEdit = true;

        await getLibraries();
        await getModels();
        await handleGetAPICollections();

        if (!botId && pageAction === "detail") {
          triggerToast("error");
          throw new Error("No botId provided");
        } else if (pageAction === "detail") {
          await getBot();
        }
      } catch (error) {
        console.log("error", error);
        if (window.history.length > 1) {
          goto(-1);
        } else {
          goto(isFrontend ? "/bots/list" : await $call.getPage("client-dashboard", "bots-list"));
        }
      } finally {
        onMountState = false;
        triggerLoading(false);
      }
    });
    return () => unsubscribe();
  });

  $effect(() => {
    if (selectedCollectionId) {
      if (
        typeof selectedCollectionId === "number" &&
        selectedCollectionId !== oldSelectedCollectionId
      ) {
        oldSelectedCollectionId = selectedCollectionId;
        tmpDatasChecked = {};
        getCollections(selectedCollectionId);
      }
    }
  });
</script>

<svelte:window bind:innerWidth={screenWidth} />

<div class="m-4 pb-12">
  <div class="sticky z-[1] top-0 bg-primary-light dark:bg-primary-dark">
    <div class="flex flex-col md:flex-row gap-4 justify-between">
      <div class="flex items-center gap-4 ms-4">
        <Avatar
          class="size-10 md:size-14 border-[1px] border-hc-graylight-200"
          src={pageAction === "create"
            ? fileImageName
              ? fileImageResult
              : ""
            : botFile.name
              ? botFile.imageResult
              : ""}
          stacked
        />
        <div>
          <div class="flex items-center gap-2">
            {#if pageAction === "create" && botName !== ""}
              <Typography preset="Display sm/Semibold" text={botName} />
            {:else if pageAction === "detail"}
              <Typography preset="Display sm/Semibold" text={users.botName} />
            {:else}
              <Typography preset="Display sm/Semibold" text={$_("createNewBots")} />
            {/if}
            {#if pageAction !== "create"}
              <Badge type={classStatusType(users.deployedPlatform)} text={users.deployedPlatform} />
            {/if}
          </div>
          {#if pageAction !== "create"}
            <Typography
              preset="Text md/Medium"
              text={formattedDate}
              color="quarterary"
              propClass="fluid-sm md:fluid-base"
            />
          {/if}
        </div>
      </div>
      <div class="flex items-center gap-2">
        {#if pageAction === "create"}
          <Button
            on:click={async () => {
              goto(isFrontend ? "/bots" : await $call.getPage("client-dashboard", "bots"));
            }}
            color="light"
            class="w-full h-min"
          >
            {$_("cancel")}
          </Button>
          <Button
            disabled={botName === "" || areAllUnchecked()}
            on:click={() => {
              createBot();
            }}
            class="w-full h-min bg-hc-brand-600 hover:bg-hc-brand-700 dark:bg-hc-brand-600 dark:hover:bg-hc-brand-700"
          >
            {$_("create")}
          </Button>
        {:else if onEdit}
          <Button
            onclick={() => {
              onEdit = false;
              tabDatas.selected = "Information";
              const url = new URL(window.location.href);
              if (url.searchParams.get("edit")) {
                url.searchParams.delete("edit");
                window.history.replaceState({}, "", url);
              }
              getBot();
            }}
            color="light"
            class="w-full px-3 focus-within:!ring-1"
          >
            {$_("cancel")}
          </Button>
          <Button
            on:click={() => {
              onEdit = false;
              tabDatas.selected = "Information";
              const url = new URL(window.location.href);
              if (url.searchParams.get("edit")) {
                url.searchParams.delete("edit");
                window.history.replaceState({}, "", url);
              }
              updateBot();
            }}
            class="w-full px-3 focus-within:!ring-1"
          >
            {$_("save")}
          </Button>
        {:else}
          <Button
            onclick={async function () {
              goto(isFrontend ? "/bots" : await $call.getPage("client-dashboard", "bots"));
            }}
            color="light"
            class="w-full px-3 focus-within:!ring-1"
          >
            {$_("back")}
          </Button>
          <Button
            onclick={() => (modalDelete = true)}
            outline
            color="red"
            class="w-full px-3 focus-within:!ring-1"
          >
            {$_("delete")}
          </Button>
          <Button
            onclick={() => {
              onEdit = true;
            }}
            class="w-full px-3 focus-within:!ring-1"
          >
            {$_("edit")}
          </Button>
        {/if}
      </div>
    </div>
    <!-- Submenu -->
    {#if pageAction !== "create"}
      <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />
    {/if}
  </div>

  <!-- Content -->
  {#if tabDatas.selected === "Information"}
    <div class="mt-6 space-y-6">
      <div class="md:flex">
        <Typography
          preset="Text md/Semibold"
          text={$_("botsInformation")}
          propClass="w-full md:w-1/4"
        />
        <div class="w-full md:w-3/4 mt-2 md:mt-0 grid grid-cols-1 md:grid-cols-4 gap-3">
          <div class="flex flex-col gap-1.5 col-span-4">
            <Typography preset="Text sm/Medium" text={$_("botsName")} color="secondary" />
            {#if pageAction === "create"}
              <InputRw bind:value={botName} {onEdit} />
            {:else}
              <InputRw bind:value={users.botName} {onEdit} />
            {/if}
          </div>
          <div class="flex flex-col gap-1.5 col-span-4">
            <Typography
              preset="Text sm/Medium"
              text={$_("descriptionOptional")}
              color="secondary"
            />
            {#if pageAction === "create"}
              <TextareaRw bind:value={description} {onEdit} />
            {:else}
              <TextareaRw bind:value={users.description} {onEdit} />
            {/if}
          </div>
          {#if pageAction === "create"}
            <div class="col-span-4">
              <UploadFile
                title={$_("botsProfilePicture")}
                handleFileUploadWrapper={() => handleInputFile("bot")}
                srcAvatar={fileImageName ? fileImageResult : ""}
              />
            </div>
          {/if}
        </div>
      </div>
      <Divider />
      <div class="md:flex">
        <Typography
          preset="Text md/Semibold"
          text={$_("aiEngines")}
          color="secondary"
          propClass="w-full md:w-1/4"
        />
        <div class="w-full md:w-3/4 mt-2 md:mt-0">
          <SelectedCard
            selectionData={platforms}
            bind:selectedData={selectedplatforms}
            numGrid={3}
            isRadio
            disabled={!onEdit ||
              platforms?.find((platform) => platform.id === users?.platform)?.model === "openai"}
          />
          <div class="flex flex-col gap-1 mt-2">
            <div class="flex items-center gap-2">
              <VerifiedTickIcon class="w-4 h-4" />
              <Typography
                preset="Text sm/Regular"
                text={$_("recommendedModels")}
                color="secondary"
              />
            </div>
            <div class="flex items-center gap-2">
              <ApiConnectorSupportIcon class="w-4 h-4" />
              <Typography
                preset="Text sm/Regular"
                text={$_("supportForApiConnector")}
                color="secondary"
              />
            </div>
          </div>
        </div>
      </div>
      {#if platforms?.find((platform) => platform.id === selectedplatforms)?.model !== "openai"}
        <Divider />
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text={$_("knowledgeSource")}
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4 mt-2 md:mt-0">
            <div class="flex items-center gap-4">
              <Button
                color="light"
                class="w-max rounded-lg border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1.5 flex justify-between items-center py-2.5 px-3"
                on:click={() => {
                  modalKnowledge = true;
                }}
                disabled={!onEdit}
              >
                <Typography
                  preset="Text md/Semibold"
                  text={selectedCollectionId ? $_("configured") : $_("configure")}
                  color={selectedCollectionId ? "fg-success-primary" : "secondary"}
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
          text={$_("apiSources")}
          color="secondary"
          propClass="w-full md:w-1/4"
        />
        <div class="w-full md:w-3/4 mt-2 md:mt-0">
          <div class="flex flex-col justify-center gap-4">
            {#if pageAction === "create"}
              {#if selectedplatforms === 1}
                <Toggle {onEdit} bind:checked={apiSourceEnable} text="enable" />
              {:else}
                <Toggle onEdit={false} bind:checked={apiSourceEnable} text="enable" />
              {/if}
            {:else if Object.keys(users).length > 0}
              {#if selectedplatforms === 1}
                <Toggle {onEdit} bind:checked={users.apiSourceEnable} text="enable" />
              {:else}
                <Toggle onEdit={false} bind:checked={users.apiSourceEnable} text="enable" />
              {/if}
            {/if}
            <Button
              color="light"
              class="w-max rounded-lg border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1.5 flex justify-between items-center py-2.5 px-3"
              on:click={() => {
                modalAPI = true;
              }}
              disabled={!onEdit ||
                !(pageAction === "create" ? apiSourceEnable : users.apiSourceEnable)}
            >
              <Typography
                preset="Text md/Semibold"
                text={!datasCheckedApi.map((item) => item.checked).every((checked) => !checked)
                  ? $_("configured")
                  : $_("configure")}
                color={!datasCheckedApi.map((item) => item.checked).every((checked) => !checked)
                  ? "fg-success-primary"
                  : "secondary"}
              />
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
          text={$_("additionalPrompts")}
          color="secondary"
          propClass="w-full md:w-1/4"
        />
        <div class="w-full md:w-3/4 mt-2 md:mt-0 flex flex-col gap-1.5">
          <Typography preset="Text sm/Medium" text={$_("botsRole")} color="secondary" />
          <TextareaRw bind:value={additionalPrompt} {onEdit} rows={8} />
        </div>
      </div>
      <Divider />
      <div class="md:flex">
        <Typography
          preset="Text md/Semibold"
          text={$_("additionalSettings")}
          color="secondary"
          propClass="w-full md:w-1/4"
        />
        <div class="w-full md:w-3/4 mt-2 md:mt-0 grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="col-span-4">
            {#if pageAction === "create"}
              <Toggle {onEdit} bind:checked={needAuth} text="requiresLogin" />
            {:else if pageAction === "detail" && users?.needAuth !== undefined}
              <Toggle {onEdit} bind:checked={users.needAuth} text="requiresLogin" />
            {/if}
            <Typography
              preset="Text sm/Regular"
              text={$_("requiresLoginDesc")}
              color="tertiary"
              propClass="mt-1"
            />
          </div>
          {#if pageAction === "detail" && users?.documentOptions !== undefined}
            <div class="col-span-4">
              <Toggle {onEdit} bind:checked={users.documentOptions} text="documentOptions" />
            </div>
          {/if}
        </div>
      </div>
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
        <div class="w-full md:w-3/4 mt-2 md:mt-0 grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="flex flex-col gap-1.5 col-span-4">
            <Typography preset="Text sm/Medium" text={$_("welcomeMessage")} color="secondary" />
            <TextareaRw bind:value={users.welcomeMessage} {onEdit} />
          </div>
        </div>
      </div>
      <!-- <Divider />
      <div class="flex flex-col md:flex-row gap-1.5">
        <div class="flex flex-col gap-2 w-full md:w-1/4">
          <Typography
            preset="Text md/Semibold"
            text={$_("promptRecommendation")}
            color="secondary"
          />
          <Typography
            preset="Text md/Regular"
            text={$_("promptRecommendationDesc")}
            color="tertiary"
          />
        </div>
        <div class="flex flex-col gap-3 w-full md:w-3/4">
          {#each recommendations as _unneeded, index}
            <div class="flex flex-col gap-2">
              <div class="flex justify-between items-center">
                <Typography
                  preset="Text sm/Medium"
                  text={$_("promptRecommendation") + ` #${index + 1}`}
                  color="secondary"
                />
                <Button
                  color="none"
                  disabled={!onEdit}
                  class="p-0 mr-8 focus:ring-0 outline-none"
                  on:click={() => {
                    if (onEdit) {
                      recommendations = recommendations.filter((_, i) => i !== index);
                    }
                  }}
                >
                  <Trash01Icon
                    svgClass="w-5 h-5 text-error-primary-light dark:text-error-primary-dark"
                  />
                </Button>
              </div>
              <TextareaRw bind:value={recommendations[index]} {onEdit} rows={3} />
            </div>
          {/each}
          <Button
            class="flex gap-2 w-fit focus:ring-0"
            on:click={() => {
              if (onEdit) {
                recommendations = [...recommendations, ""];
              }
            }}
            disabled={!onEdit || recommendations.length >= 3}
          >
            <PlusIcon svgClass="w-5 h-5 text-primary-light dark:text-primary-dark" />
            <Typography preset="Text md/Regular" text={$_("addPromptRecommendation")} />
          </Button>
          {#if recommendations.length <= 0}
            <Typography
              preset="Text md/Regular"
              text={$_("promptRecommendationNotFound")}
              color="secondary"
            />
          {/if}
        </div>
      </div> -->
      <Divider />
      <div class="md:flex">
        <Typography
          preset="Text md/Semibold"
          text={$_("botsProfilePicture")}
          color="secondary"
          propClass="w-full md:w-1/4"
        />
        <div class="w-full md:w-3/4">
          <UploadFile
            title={$_("botsProfilePicture")}
            handleFileUploadWrapper={() => handleInputFile("bot")}
            isPreview
            srcAvatar={botFile.name ? botFile.imageResult : ""}
            disabled={!onEdit}
          />
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
        <div class="w-full md:w-3/4 mt-2 md:mt-0 grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="flex flex-col gap-1.5 col-span-1">
            <Typography preset="Text sm/Medium" text={$_("primaryBrandColor")} color="secondary" />
            <Button
              on:click={() => {
                if (onEdit) {
                  tmpThemeColor = users.primaryBrandColor;
                  modalThemeColor = true;
                }
              }}
              class="flex justify-between border {!onEdit
                ? 'cursor-not-allowed focus:ring-0 border-disabled-subtle-light dark:border-disabled-subtle-dark hover:border-disabled-subtle-light hover:dark:border-disabled-subtle-dark focus:border-disabled-subtle-light focus:dark:border-disabled-subtle-dark bg-disabled-subtle-light dark:bg-disabled-subtle-dark hover:bg-disabled-subtle-light dark:hover:bg-disabled-subtle-dark text-disabled-light dark:text-disabled-dark'
                : 'border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark focus:ring-0'}"
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
        <div class="w-full md:w-3/4 mt-2 md:mt-0">
          {#if users?.layout !== undefined}
            <SelectedCard
              isRadio
              widthFlex="w-full"
              spacing="!p-4"
              numGrid={2}
              disabled={!onEdit}
              selectionData={layouts}
              bind:selectedData={users.layout}
            />
          {/if}
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
          <div class="w-full md:w-3/4 mt-2 md:mt-0 grid grid-cols-1 md:grid-cols-4 gap-4">
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
                class={`flex justify-between ${!onEdit ? "cursor-default focus-within:ring-0 hover:bg-primary-light hover:dark:bg-secondary-dark" : ""} border-primary-light dark:border-primary-dark bg-primary-light dark:bg-secondary-dark`}
                color="light"
                on:click={!onEdit
                  ? () => {}
                  : () => {
                      tmpOverlayColor = users.overlayColor;
                      modalOverlayColor = true;
                    }}
                disabled={!onEdit}
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
              <Typography preset="Text sm/Medium" text={$_("backgroundAlign")} color="secondary" />
              <Dropdown
                dropdownOpen={dropdownBgAlignOpen}
                bind:selectedItem={users.background_align}
                items={bgAlign}
                isOverlay={true}
                disabled={!onEdit}
              />
            </div>
            <!-- <div class="col-span-4">
              <div>
                <Typography
                  preset="Text sm/Medium"
                  text={$_("overlayTransparency") + ": " + users.backgroundOpacity + "%"}
                  color="secondary"
                />
                <InputSlider bind:value={users.backgroundOpacity} />
              </div>
            </div> -->
          </div>
        </div>
      {/if}
      <!-- Preview section -->
      <!-- <Divider />
      <div class="md:flex">
        <Typography
          preset="Text md/Semibold"
          text={$_("preview")}
          color="secondary"
          propClass="w-full md:w-1/4"
        />
        <div class="w-full md:w-3/4 mt-2 md:mt-0 grid grid-cols-1 md:grid-cols-4 gap-4">
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
                    <BrandIcon svgClass="w-5 h-5 text-gray-700 dark:text-gray-400" />
                  {/if}
                </div>
                {#if users.logoName && !users.logoOnly}
                  <Typography preset="Text xl/Semibold" text={users.logoName} />
                {/if}
              </div>
              <div class="flex flex-col items-center gap-4 px-4 pb-4">
                <div class="flex w-full md:w-1/2 justify-end">
                  <Button outline color="alternative" class="mt-4 w-fit">
                    <RefreshCw05Icon svgClass="w-4 h-4 hover:text-white" />
                    &nbsp;
                    {$_("clearConversations")}
                  </Button>
                </div>
                <div class="flex flex-col w-full md:w-1/2 gap-4">
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
                        <ThumbsDownIcon svgClass="w-5 h-5 text-gray-700 dark:text-gray-400" />
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
                      <ThumbsDownIcon svgClass="w-5 h-5 text-gray-700 dark:text-gray-400" />
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
                      <ThumbsDownIcon svgClass="w-5 h-5 text-gray-700 dark:text-gray-400" />
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
                      <Send03Icon svgClass="w-5 h-5 text-white" />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div> -->
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
                  (isFrontend
                    ? "documents/libraries/collections/documents" +
                      `?collectionId=${collectionId[index]}`
                    : await $call.getPage("client-dashboard", "document-library-documents")) +
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
                bind:datasCheckedKnowledge={datasCategoriesChecked[collectionId[index]]}
                bind:configsKnowledge={configsCategories[collectionId[index]]}
                bind:gettingDatasKnowledge
                next={nextByCategory}
                previous={previousByCategory}
                handleCheckAll={handleCheckAllCategories}
                handleCheck={handleCheckCategories}
              />
            {:else if gettingDatasKnowledge}
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

        <div class="w-full md:w-3/4 mt-2 md:mt-0 grid grid-cols-1 md:grid-cols-4 gap-4">
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

<Modal bind:open={modalKnowledge} outsideclose>
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
        bind:selectedData={selectedCollectionId}
        numGrid={2}
        isRadio
      />
    </div>
  {:else if allData.length > 0 || searchTerm !== ""}
    {#if !gettingDatasKnowledge}
      <DataTableListCollection
        bind:datas={allData}
        bind:datasChecked={datasCheckedKnowledge}
        bind:configs={configsKnowledge}
        bind:gettingDatas={gettingDatasKnowledge}
        {next}
        {previous}
        {handleCheckAll}
        {handleCheck}
      />
    {:else}
      <StateCard illustration={IllustrationProgress} description={$_("fetchingDataPleaseWait")} />
    {/if}
  {:else if gettingDatasKnowledge}
    <StateCard illustration={IllustrationProgress} description={$_("fetchingDataPleaseWait")} />
  {:else}
    <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
  {/if}
  <svelte:fragment slot="footer">
    {#if !isSelectedLibrary}
      <Button
        on:click={() => {
          isSelectedLibrary = true;
          searchTerm = "";
        }}
        class="w-full"
        disabled={!selectedCollectionId}>{$_("next")}</Button
      >
    {:else}
      <div class="flex gap-4 w-full">
        <button
          onclick={() => {
            isSelectedLibrary = false;
            searchTerm = "";
          }}
          class="px-5">{$_("goBack")}</button
        >
        <Button
          on:click={() => {
            modalKnowledge = false;
            searchTerm = "";
          }}
          class="grow"
          disabled={areAllUnchecked()}>{$_("finish")}</Button
        >
      </div>
    {/if}
  </svelte:fragment>
</Modal>

<Modal bind:open={modalAPI} outsideclose>
  <svelte:fragment slot="header">
    <div class="flex gap-4 me-6 grow w-full">
      <div class="grow">
        <Typography
          preset="Text md/Bold"
          text={$_("chooseApiSources")}
          color={!isSelectedLibrary ? "brand-secondary" : "secondary"}
        />
        <Typography
          preset="Text sm/Regular"
          text={$_("apiSourceAsThePrimarySource")}
          color={!isSelectedLibrary ? "brand-tertiary" : "tertiary"}
        />
      </div>
    </div>
  </svelte:fragment>
  {#if apiCollections.length > 0 || searchTerm !== ""}
    <DataTableListCollection
      bind:datas={apiCollections}
      bind:datasChecked={datasCheckedApi}
      bind:configs={configsApi}
      isApi={true}
      {next}
      {previous}
      {handleCheckAll}
      {handleCheck}
    />
  {:else}
    <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
  {/if}
  <svelte:fragment slot="footer">
    <div class="flex gap-4 w-full">
      <button
        onclick={() => {
          modalAPI = false;
          searchTerm = "";
        }}
        class="px-5"
      >
        {$_("cancel")}
      </button>
      <Button
        on:click={() => {
          modalAPI = false;
          searchTerm = "";
        }}
        class="grow"
        disabled={datasCheckedApi.map((item) => item.checked).every((checked) => !checked)}
      >
        {$_("finish")}
      </Button>
    </div>
  </svelte:fragment>
</Modal>
