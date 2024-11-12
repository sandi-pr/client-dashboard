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
  import { _, locale } from "svelte-i18n";
  import DataTableListCollection from "$lib/components/DataTableListCollection.svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { StateCard, InputRw, TextareaRw, triggerToast, UploadFile, Typography } from "$lib";
  import TabBar from "$lib/components/TabBar.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { format, formatDistanceToNow } from "date-fns";
  import { id as localeId, enUS as localeEn } from "date-fns/locale";
  import XCloseIcon from "$lib/icon/x-close.svelte";
  import {
    Button,
    P,
    Modal,
    Tooltip,
    Indicator,
    Input,
    Alert,
    Select,
    MultiSelect,
  } from "flowbite-svelte";
  import FileTypeIcon from "$lib/icon/file-type-icon.svelte";
  import {
    ExclamationCircleSolid,
    ChevronDownOutline,
    InfoCircleSolid,
  } from "flowbite-svelte-icons";
  import PlusIcon from "$lib/icon/plus.svg";
  import MessageChatCircle from "$lib/icon/message-chat-circle.svelte";
  import FolderPlusIcon from "$lib/icon/folder-plus.svelte";
  import Stars03Icon from "$lib/icon/stars-03.svelte";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import { goto } from "$app/navigation";
  import {
    triggerSelectCollection,
    triggerShowAI,
    triggerSelectDropdown,
    triggerLoading,
    triggerReload,
    check,
  } from "$lib";

  const params = new URLSearchParams(window.location.search);
  let collectionId = params.get("collectionId");

  import { call } from "$lib";

  let documentsReady = $state(false);

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });

  let allData = $state([]);
  let datas = $state([]);

  let status = $state(null);

  let statusRag = $state(null);

  let gettingDatas = $state(true);

  let onEdit = $state(false);

  let popupModal = $state(false);

  let deleteCategory = $state(false);

  let deleteCategoriesModal = $state(false);

  let datasChecked = $state([]);
  let selectCategories = $state(null);

  let datasCategoriesChecked = $state({});

  let collection = $state({});

  let configs = $state({});

  let configsCategories = $state({});

  // let dropdownToggle = $state([true, false]);

  let tabDatas = $state({
    selected: "Lists",
    data: [
      {
        name: "Lists",
      },
      // {
      //   name: "By Categories",
      // },
      {
        name: "Information",
      },
    ],
  });

  let files = $state([]);
  let tmpDocumentType = $state("Documents");
  let tmpDocumentTitle = $state("");
  let categories = $state([]);
  let tmpCategories = $state([""]);
  let tmpSelectedCategories = $state(["ayam"]);
  let otherCategory = $state("");
  let tmpDocumentDescription = $state("");
  let selectedTypeCollection = $state("Documents (TXT, PDF)");
  let acceptedFileTypes = $state(".txt, .pdf");
  let addDocumentsModal = $state(false);
  let addCategoriesModal = $state(false);
  let dropdownOpen = $state(false);
  let newModelName = $state("");
  let trainingModelModal = $state(false);
  let fileNotSupport = $state(false);
  let addOtherCategory = $state(false);
  let search = $state(false);

  let timer;

  let page = $state(1);

  let searchTerm = $state("");
  let searchTermCategory = $state({});

  const handleSearch = async (value) => {
    search = true;
    page = 1;
    searchTerm = value;

    await getDocuments();
  };

  const handleSearchCategory = async (value, category) => {
    configsCategories[category].pagination.page = 1;
    searchTermCategory = { ...searchTermCategory, [category]: value };

    await getDocumentsByCategory(category);
  };

  const dummyData = {
    id: "1",
    file_name: {
      name: "dummy document",
      filename: "dummy.pdf",
      type: "custom",
      size: "2MB",
    },
    category: {
      name: "dummy category",
      type: "pill",
      color: "info",
    },
    type: {
      name: "dummy type",
      type: "pill",
      color: "brand",
    },
    file_size: "2MB",
    upload_date: "2024-01-01",
    actions: {
      objects: [
        {
          href: "/document-collection/collections/detail/document?documentId=1&collectionId=1",
          name: "View",
        },
      ],
    },
  };

  async function getDocuments() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      if (!search) {
        gettingDatas = true;
      }
      let response = await $call.call("client-dashboard", "get-documents-by-collectionid", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: page,
          page_size: pagination.pageSize,
          collection_id: collectionId,
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

      if (response.result?.status_rag) {
        statusRag = response.result.status_rag;
      }

      if (response.result.data.length > 0) {
        pagination = response.result.meta.pagination;
        configs = {
          search: true,
          searchHandler: handleSearch,
          headers: "File Uploaded",
          pagination: response.result.meta.pagination,
          isCategories: false,
          actions: [
            check("documentAi.documentLibrary.document:delete") && {
              id: "delete",
              icon: "trash",
              type: "error",
              disabled: datasChecked.find((data) => data.checked) ? false : true,
              onclick: () => {
                if (datasChecked.find((data) => data.checked)) {
                  deleteCategory = false;
                  popupModal = true;
                }
              },
            },
            check("documentAi.documentLibrary.document:create") && {
              id: "add_new",
              name: "Add New",
              icon: "plus",
              type: "brand-inverse",
              onclick: () => {
                selectTypeCollection("Documents");
                addDocumentsModal = true;
                addOtherCategory = false;
              },
            },
          ].filter(Boolean),
        };

        let promises = response.result.data.map(async (document) => {
          let link = await $call.getPage("client-dashboard", "document-collection-detail");
          return {
            id: document.id,
            file_name: {
              name: document.attributes.document_title,
              filename: document.attributes.document.data.attributes.name,
              type: "custom",
              size:
                document.attributes.file_size === null ? "Unknown" : document.attributes.file_size,
            },
            // category: {
            //   name: document.attributes.category,
            //   type: "pill",
            //   color: "info",
            // },
            type: {
              name: document.attributes.document_type,
              type: "pill",
              color: "brand",
            },
            file_size:
              document.attributes.file_size === null ? "Unknown" : document.attributes.file_size,
            upload_date: document.attributes.createdAt,

            // uploaded_by:
            //   document.attributes.upload_by.data === null
            //     ? "Unknown"
            //     : document.attributes.upload_by.data.attributes.fullname,
            actions: {
              objects: [
                check("documentAi.documentLibrary.document:view") && {
                  // href: "/document/document_collection/detail?documentId=1&collectionId=1",
                  href: `${link}` + "?documentId=" + document.id + "&collectionId=" + collectionId,
                  name: "View",
                },
                // {
                //   // href: "localhost:5173/client-dashboard/document-collection-detail?documentId=1&collectionId=1",
                //   // `${await $call.getPage("client-dashboard", "document-library-detail")}` +
                //   // "?documentId=" +
                //   // document.id +
                //   // "&collectionId=" +
                //   // collectionId,
                //   function: handleTitikTiga,
                //   name: "Titik tiga",
                // },
              ].filter(Boolean),
            },
          };
        });

        allData = await Promise.all(promises);
        console.log("all data", allData);
        datasChecked = allData.map((data) => {
          return { id: data.id, checked: false };
        });
      } else {
        allData = [];
        datasChecked = [];
        configs = {
          search: true,
          searchHandler: handleSearch,
          headers: "File Uploaded",
          pagination: response.result.meta.pagination,
          actions: [
            {
              id: "delete",
              name: "Delete",
              icon: "trash",
              type: "error",
              disabled: datasChecked.find((data) => data.checked) ? false : true,
              onclick: () => {
                if (datasChecked.find((data) => data.checked)) {
                  popupModal = true;
                }
              },
            },
            {
              id: "add_new",
              name: "Add New",
              icon: "plus",
              type: "brand-inverse",
              onclick: () => {
                selectTypeCollection("Documents");
                addDocumentsModal = true;
                addOtherCategory = false;
              },
            },
          ],
        };
      }

      if (!search) {
        // await getAllDocumentsByCategory();
      }

      gettingDatas = false;
      search = false;
    }
  }

  async function getAllDocumentsByCategory() {
    let token = localStorage.getItem("token");
    const response2 = await $call.call("client-dashboard", "get-document-by-categories", {
      URLQuery: {},
      URLParams: {},
      Params: {
        collection_id: collectionId.toString(),
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

    console.log("ayam", response2);

    let datasTmp = response2.result.data;

    categories = datasTmp.map((document) => {
      return {
        value: document.category,
        name: document.category,
      };
    });

    datas = [];

    for (let i = 0; i < datasTmp.length; i++) {
      if (!datasCategoriesChecked[datasTmp[i].category]) {
        datasCategoriesChecked[datasTmp[i].category] = [];
      }

      if (!searchTermCategory[datasTmp[i].category]) {
        searchTermCategory[datasTmp[i].category] = "";
      }

      let promises = datasTmp[i].data.map(async (document) => {
        datasCategoriesChecked[datasTmp[i].category].push({
          id: document.id,
          checked: false,
        });

        let link = await $call.getPage("client-dashboard", "document-collection-detail");
        return {
          id: document.id,
          file_name: {
            name: document.attributes.document_title,
            filename: document.attributes.document.data.attributes.name,
            type: "custom",
            size:
              document.attributes.file_size === null ? "Unknown" : document.attributes.file_size,
          },
          category: {
            name: document.attributes.category,
            type: "pill",
            color: "info",
          },
          type: {
            name: document.attributes.document_type,
            type: "pill",
            color: "brand",
          },
          file_size:
            document.attributes.file_size === null ? "Unknown" : document.attributes.file_size,
          upload_date: document.attributes.createdAt,

          // uploaded_by:
          //   document.attributes.upload_by.data === null
          //     ? "Unknown"
          //     : document.attributes.upload_by.data.attributes.fullname,
          actions: {
            objects: [
              {
                // href: "/document/document_collection/detail?documentId=1&collectionId=1",
                href: `${link}` + "?documentId=" + document.id + "&collectionId=" + collectionId,
                name: "View",
              },
              // {
              //   // href: "localhost:5173/client-dashboard/document-collection-detail?documentId=1&collectionId=1",
              //   // `${await $call.getPage("client-dashboard", "document-library-detail")}` +
              //   // "?documentId=" +
              //   // document.id +
              //   // "&collectionId=" +
              //   // collectionId,
              //   function: handleTitikTiga,
              //   name: "Titik tiga",
              // },
            ],
          },
        };
      });

      let category = datasTmp[i].category;

      datas.push({
        category: category,
        data: await Promise.all(promises),
      });

      configsCategories = {
        ...configsCategories,
        [category]: {
          search: true,
          searchHandler: handleSearchCategory,
          headers: category,
          pagination: datasTmp[i].meta.pagination,
          isCategories: true,
          actions: [
            {
              id: "delete",
              icon: "trash",
              type: "error",
              disabled: datasCategoriesChecked[category].find((data) => data.checked)
                ? false
                : true,
              onclick: () => {
                if (datasCategoriesChecked[category].find((data) => data.checked)) {
                  deleteCategory = true;
                  selectCategories = category;
                  popupModal = true;
                }
              },
            },
            {
              id: "add_new",
              // name: "Add New",
              icon: "plus",
              type: "brand-inverse",
              onclick: () => {
                selectTypeCollection("Documents");
                addDocumentsModal = true;
                addOtherCategory = false;
                tmpCategories = [];
                for (let i = 0; i < categories.length; i++) {
                  if (categories[i]?.other === true) {
                    categories = categories.slice(0, i).concat(categories.slice(i + 1));
                    break;
                  }
                }
                tmpCategories = category;
              },
            },
            {
              id: "ai",
              icon: "ai",
              name: "ai",
              type: "brand",
              disabled: datasTmp[i].data.length === 0,
              onclick: () => {
                selectCollection({
                  id: collectionId,
                  name: collection.name,
                  description: collection.description,
                  totalDocuments: datas.length,
                });
              },
            },
            {
              id: "expand",
              icon: "chevron",
              type: "info",
            },
            {
              id: "dots-vertical",
              icon: "dots-vertical",
              type: "info",
              dropdown: "more",
            },
          ],
        },
      };

      // configsCategories.push({
      //   // search: true,
      //   // searchHandler: handleSearch,
      //   headers: category,
      //   pagination: datasTmp[i].meta.pagination,
      //   isCategories: true,
      //   actions: [
      //     {
      //       id: "delete",
      //       icon: "trash",
      //       type: "error",
      //       disabled: datasCategoriesChecked[category].find((data) => data.checked)
      //         ? false
      //         : true,
      //       onclick: () => {
      //         if (datasCategoriesChecked[category].find((data) => data.checked)) {
      //           deleteCategory = true;
      //           selectCategories = category;
      //           popupModal = true;
      //         }
      //       },
      //     },
      //     {
      //       id: "add_new",
      //       name: "Add New",
      //       icon: "plus",
      //       type: "brand-inverse",
      //       onclick: () => {
      //         selectTypeCollection("Documents");
      //         addDocumentsModal = true;
      //         tmpCategories = [];
      //         tmpCategories = category;
      //       },
      //     },
      //     {
      //       id: "expand",
      //       icon: "chevron",
      //       type: "info",
      //     },
      //     {
      //       id: "dots-vertical",
      //       icon: "dots-vertical",
      //       type: "info",
      //       dropdown: "delete",
      //     },
      //   ],
      // });
    }
  }

  async function getDocumentsByCategory(category) {
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "get-document-by-category", {
          URLQuery: {},
          URLParams: {},
          Params: {
            category: category,
            collection_id: collectionId.toString(),
            page: configsCategories[category].pagination.page,
            page_size: configsCategories[category].pagination.pageSize,
            search: searchTermCategory[category],
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

        for (let i = 0; i < datas.length; i++) {
          if (datas[i].category === category) {
            datasCategoriesChecked[datas[i].category] = [];
            let promises = response.result.data.map(async (document) => {
              datasCategoriesChecked[category].push({
                id: document.id,
                checked: false,
              });

              let link = await $call.getPage("client-dashboard", "document-collection-detail");
              return {
                id: document.id,
                file_name: {
                  name: document.attributes.document_title,
                  filename: document.attributes.document.data.attributes.name,
                  type: "custom",
                  size:
                    document.attributes.file_size === null
                      ? "Unknown"
                      : document.attributes.file_size,
                },
                category: {
                  name: document.attributes.category,
                  type: "pill",
                  color: "info",
                },
                type: {
                  name: document.attributes.document_type,
                  type: "pill",
                  color: "brand",
                },
                file_size:
                  document.attributes.file_size === null
                    ? "Unknown"
                    : document.attributes.file_size,
                upload_date: document.attributes.createdAt,

                actions: {
                  objects: [
                    {
                      // href: "/document/document_collection/detail?documentId=1&collectionId=1",
                      href:
                        `${link}` + "?documentId=" + document.id + "&collectionId=" + collectionId,
                      name: "View",
                    },
                  ],
                },
              };
            });

            datas[i].data = await Promise.all(promises);

            break;
          }
        }

        configsCategories[category].pagination = response.result.meta.pagination;
        // gettingDatas = false;
      } catch (error) {
        console.log(error);
        triggerToast("error");
      }
    }
  }

  $effect(() => {
    if (datasChecked) {
      if (configs?.actions) {
        configs.actions.find((config) => {
          if (config.id === "delete") {
            config.disabled = datasChecked.find((data) => data.checked) ? false : true;
          }
        });
      }
    }
  });

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
        gettingDatas = false;
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function addDocumentByCollection() {
    files[0].document_title = tmpDocumentTitle;
    files[0].description = tmpDocumentDescription;

    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        addDocumentsModal = false;
        triggerLoading(true);

        let dataSend = null;

        if (tmpDocumentType === "Handwriting" || tmpDocumentType === "Document Images") {
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

            console.log("trocr-predict", response3);

            if (response3.status !== 200) {
              addDocumentsModal = true;
              triggerToast("error");
              triggerLoading(false);
              console.log(response3);
              return;
            }

            dataSend = {
              documents: JSON.stringify(files),
              collection_id: collectionId,
              document_type: tmpDocumentType,
              category: tmpCategories,
              detected_text: response3.result.text,
              result_trocr: JSON.stringify(response3.result),
              edited: false,
            };

            let response = await $call.call("client-dashboard", "create-document", {
              URLQuery: {},
              URLParams: {},
              Params: dataSend,
              Headers: {
                Authorization: "Bearer " + token,
              },
            });

            if (response.status !== 200) {
              addDocumentsModal = true;
              console.log(response);
              triggerToast("error");
              triggerLoading(false);
              return;
            }
          }
        } else {
          if (dataSend === null) {
            dataSend = {
              documents: JSON.stringify(files),
              collection_id: collectionId,
              document_type: tmpDocumentType,
              category: tmpCategories,
            };
          }

          let response = await $call.call("client-dashboard", "create-document", {
            URLQuery: {},
            URLParams: {},
            Params: dataSend,
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response.status !== 200) {
            addDocumentsModal = true;
            console.log(response);
            triggerToast("error");
            triggerLoading(false);
            return;
          }
        }
      } catch (error) {
        addDocumentsModal = true;
        console.log(error);
        triggerToast("error");
        triggerLoading(false);
        return;
      }

      files = [];

      tmpDocumentTitle = "";
      tmpDocumentDescription = "";
      tmpCategories = [];

      addDocumentsModal = false;
      getDocuments();
      triggerLoading(false);

      triggerToast("created");
    }
  }

  function functDeleteCategories(categories) {
    deleteCategoriesModal = true;
    selectCategories = categories;
  }

  async function deleteCategories() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      popupModal = false;
      triggerLoading(true);
      try {
        const response = await $call.call("client-dashboard", "delete-document-by-categories", {
          URLQuery: {},
          URLParams: {},
          Params: {
            category: selectCategories,
            collection_id: collectionId.toString(),
          },
          Headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          console.log(response);
          return;
        }

        await getCollection();
        await getDocuments();
        triggerToast("deleted");
        triggerLoading(false);
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function deleteDocument() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      popupModal = false;
      triggerLoading(true);
      let ids = [];
      if (deleteCategory === true) {
        ids = datasCategoriesChecked[selectCategories]
          .filter((data) => data.checked)
          .map((data) => data.id);
      } else {
        ids = datasChecked.filter((data) => data.checked).map((data) => data.id);
      }
      try {
        for (const id of ids) {
          let response = await $call.call("client-dashboard", "delete-document", {
            URLQuery: {},
            URLParams: {
              id: id.toString(),
            },
            Headers: {
              Authorization: `Bearer ${token}`,
            },
          });

          if (response.status !== 200) {
            triggerLoading(false);
            triggerToast("error");
            popupModal = true;
            console.log(response);
            return;
          }
        }

        await getCollection();
        await getDocuments();
        triggerToast("deleted");
        triggerLoading(false);
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
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
      selectedTypeCollection = "Documents (TXT, PDF)";
      acceptedFileTypes = ".txt, .pdf";
    } else if (type === "Images") {
      selectedTypeCollection = "Images (JPG, JPEG, PNG, WEBP, IMG)";
      acceptedFileTypes = ".jpg, .jpeg, .png, .webp, .img";
    } else if (type === "Document Images") {
      selectedTypeCollection = "Document Images (JPG, JPEG, PNG)";
      acceptedFileTypes = ".jpg, .jpeg, .png";
    } else if (type === "Handwriting") {
      selectedTypeCollection = "Handwriting (JPG, JPEG, PNG)";
      acceptedFileTypes = ".jpg, .jpeg, .png";
    } else if (type === "CSV") {
      selectedTypeCollection = "CSV";
      acceptedFileTypes = ".csv";
    }

    files = [];

    dropdownOpen = false;
  };

  let buttonDisabled = $state(false);

  $effect(() => {
    buttonDisabled = files.length === 0 || tmpDocumentTitle === "";
  });

  const changeTabbar = (name) => {
    tabDatas.selected = name;
  };

  async function next() {
    if (page >= pagination.pageCount) return;
    page++;
    getDocuments();
  }

  async function previous() {
    if (page <= 1) return;
    page--;
    getDocuments();
  }

  async function nextByCategory(category) {
    if (
      configsCategories[category].pagination.page >=
      configsCategories[category].pagination.pageCount
    )
      return;
    configsCategories[category].pagination.page++;
    getDocumentsByCategory(category);
  }

  async function previousByCategory(category) {
    if (configsCategories[category].pagination.page <= 1) return;
    configsCategories[category].pagination.page--;
    getDocumentsByCategory(category);
  }

  async function handleCheckAll(allChecked) {
    datasChecked = datasChecked.map((data) => {
      return { id: data.id, checked: allChecked ? false : true };
    });
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

  async function handleCheck(id) {
    let datasCheckedCopy = datasChecked;
    datasCheckedCopy = datasCheckedCopy.map((data) => {
      if (data.id === id) {
        return { id: data.id, checked: !data.checked };
      }
      return data;
    });
    datasChecked = datasCheckedCopy;
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

  async function handleTitikTiga(id) {
    console.log("Titik tiga ", id);
  }

  async function selectCollection(collection) {
    console.log("ayam list");
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
    getDocuments(pagination.page);
    collectionId ? getCollection() : null;
  });

  onDestroy(() => {
    if (timer) clearTimeout(timer);
  });

  async function startTraining() {
    if (collectionId) {
      try {
        trainingModelModal = false;
        triggerLoading(true);
        let token = localStorage.getItem("token");
        const response = await $call.call("hexa-ai-rag", "trocr-train", {
          URLQuery: {},
          URLParams: {},
          Params: {
            new_model_name: `trocr:${newModelName}`,
            collection_id: collectionId,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log(response);
        if (response.status === 200) {
          const response2 = await $call.call("client-dashboard", "update-model-trocr", {
            URLQuery: {},
            URLParams: {},
            Params: {
              model_trocr: `trocr:${newModelName}`,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response2.status !== 200) {
            console.log(response2);
            return;
          }
          sessionStorage.setItem("model_trocr", `trocr:${newModelName}`);
          triggerToast("created");
          triggerLoading(false);
        } else {
          trainingModelModal = true;
          triggerLoading(false);
          triggerToast("error");
        }
      } catch (error) {
        trainingModelModal = true;
        triggerLoading(false);
        console.log(error);
      }
    }
  }

  async function addCategories() {
    let tmp = categories.map((category) => {
      return category.name;
    });

    for (let i = 0; i < tmp.length; i++) {
      if (tags.includes(tmp[i])) {
        tags = tags.filter((tag) => {
          return tag !== tmp[i];
        });
      }
    }

    let tmp2 = [...tags, ...tmp];

    addCategoriesModal = false;
    triggerLoading(true);
    let token = localStorage.getItem("token");
    try {
      const response = await $call.call("client-dashboard", "update-categories", {
        URLQuery: {},
        URLParams: {},
        Params: {
          collection_id: collectionId.toString(),
          categories: JSON.stringify(tmp2),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerLoading(false);
        addCategoriesModal = true;
        console.log(response);
        triggerToast("error");
        return;
      }

      // await getAllDocumentsByCategory();
      triggerLoading(false);
      triggerToast("created");
      tags = [];
      addCategoriesModal = false;
    } catch (error) {
      console.log(error);
      triggerToast("error");
    }
  }

  let tags = $state([]); // Daftar tag yang sudah dimasukkan pengguna
  let newTag = $state(""); // Menyimpan input dari pengguna
  let inputRef;

  // Fungsi untuk menambahkan tag
  function addTag() {
    if (newTag.trim() && !tags.includes(newTag)) {
      tags = [...tags, newTag];
      newTag = ""; // Reset input setelah menambahkan tag
    }
  }

  // Fungsi untuk menghapus tag
  function removeTag(tagToRemove) {
    tags = tags.filter((tag) => tag !== tagToRemove);
  }

  // Fungsi untuk menambah tag saat menekan Enter
  function handleKeyDown(event) {
    if (event.key === "Enter") {
      addTag();
    }
  }

  function focusInput(input) {
    input.focus();
  }
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <div>
            <Typography
              preset="Display sm/Semibold"
              color="primary"
              text={collection.name ? collection.name : $_("loading")}
              propClass="fluid-xl"
            />
            <P
              class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
            >
              {#if !gettingDatas}
                {#if $locale === "id"}
                  {$_("createOn")}
                  {formatDistanceToNow(collection.createdAt, { addSuffix: true, locale: localeId })}
                {:else}
                  {$_("createOn")}
                  {formatDistanceToNow(collection.createdAt, { addSuffix: true, locale: localeEn })}
                {/if}
              {/if}
            </P>
          </div>
          <div>
            <div class="flex gap-2">
              <!-- {#if collectionId === "30"}
                <Button
                  disabled={datas.length === 0}
                  onclick={() => {
                    trainingModelModal = true;
                  }}
                  icon={PlusIcon}
                  ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" />
                  {$_("trainingModel")}</Button
                >
              {/if} -->
              <!-- <Button
                onclick={() => {
                  addCategoriesModal = true;
                }}
                color="light"
                class="gap-2"
                ><FolderPlusIcon
                  svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
                />{$_("addNewCategories")}</Button
              > -->
              <Button
                disabled={allData.length === 0 || statusRag === "Reading"}
                onclick={() => {
                  selectCollection({
                    id: collectionId,
                    name: collection.name,
                    description: collection.description,
                    totalDocuments: datas.length,
                  });
                }}
                class="gap-2"
              >
                {#if statusRag === "Ready"}
                  <MessageChatCircle svgClass="text-white w-5 h-5" />
                {:else}
                  <Stars03Icon svgClass="text-white w-5 h-5" />
                {/if}
                <div class="hidden lg:block">
                  {#if statusRag === null}
                    {$_("startReading")}
                  {:else if statusRag === "Reading"}
                    {$_("reading")}
                  {:else if statusRag === "Ready"}
                    {$_("askToAI")}
                  {/if}
                </div>
              </Button>
            </div>
          </div>
        </div>

        <!-- Submenu -->
        <!-- <TabBar {tabDatas} {changeTabbar} /> -->
        <TabBar {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

        {#if tabDatas.selected === "Lists"}
          <!-- List Table Data -->
          {#if allData.length > 0 || searchTerm !== ""}
            {#if !gettingDatas}
              <DataTableListCollection
                bind:datas={allData}
                bind:datasChecked
                bind:configs
                bind:gettingDatas
                bind:next
                bind:previous
                bind:handleCheckAll
                bind:handleCheck
              />
            {:else}
              <StateCard
                illustration={IllustrationProgress}
                description={$_("fetchingDataPleaseWait")}
              />
            {/if}
          {:else if gettingDatas}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else if check("documentAi.documentLibrary.document:create")}
            <StateCard
              illustration={IllustrationEmpty}
              description={$_("youDontHaveAnyDataYet")}
              actionText={$_("addNewDocument")}
              funcButton={() => {
                addDocumentsModal = true;
                addOtherCategory = false;
              }}
            />
          {:else}
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          {/if}
        {:else if tabDatas.selected === "By Categories"}
          {#if datas.length > 0 || searchTerm !== ""}
            {#if !gettingDatas}
              {#each datas as data, index}
                <DataTableListCollection
                  bind:datas={data.data}
                  bind:datasChecked={datasCategoriesChecked[data.category]}
                  bind:configs={configsCategories[data.category]}
                  bind:gettingDatas
                  bind:next={nextByCategory}
                  bind:previous={previousByCategory}
                  bind:handleCheckAll={handleCheckAllCategories}
                  bind:handleCheck={handleCheckCategories}
                  bind:functDeleteCategories
                  bind:statusRag
                />
              {/each}
            {:else}
              <StateCard
                illustration={IllustrationProgress}
                description={$_("fetchingDataPleaseWait")}
              />
            {/if}
          {:else if gettingDatas}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else if tabDatas.selected === "Lists"}
            <StateCard
              illustration={IllustrationEmpty}
              description={$_("youDontHaveAnyDataYet")}
              actionText={$_("addNewDocument")}
              funcButton={() => {
                addDocumentsModal = true;
                addOtherCategory = false;
              }}
            />
          {:else if tabDatas.selected === "By Categories"}
            <StateCard
              illustration={IllustrationEmpty}
              description={$_("youDontHaveAnyDataYet")}
              actionText={$_("addNewCategories")}
              funcButton={() => {
                addCategoriesModal = true;
              }}
            />
          {/if}
        {:else if tabDatas.selected === "Information"}
          <div class="flex flex-col gap-4 pt-4 md:pt-6">
            <!-- general information -->
            <div class="flex flex-col md:flex-row gap-3 md:gap-6">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {$_("generalInformation")}
                </P>
              </div>
              <div class="flex flex-col grow gap-2 md:gap-6">
                <div class="flex gap-5">
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("collectionName")}
                    </P>
                    <InputRw {onEdit} bind:value={collection.name} />
                  </div>
                </div>
                <div class="grid grid-cols-1">
                  <div class="flex flex-col gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >{$_("description")}</P
                    >
                    <TextareaRw {onEdit} bind:value={collection.description} />
                  </div>
                </div>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

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
            class="flex p-2 w-full"
            onclick={() => {
              selectTypeCollection("Documents");
            }}>Documents (TXT, PDF)</button
          >
          <button
            class="flex p-2 w-full"
            onclick={() => {
              selectTypeCollection("Images");
            }}>Images (JPG, JPEG, PNG, WEBP, IMG)</button
          >
          <button
            class="flex p-2 w-full"
            onclick={() => {
              selectTypeCollection("Document Images");
            }}>Document Images (JPG, JPEG, PNG)</button
          >
          <button
            class="flex p-2 w-full"
            onclick={() => {
              selectTypeCollection("Handwriting");
            }}>Handwriting (JPG, JPEG, PNG)</button
          >
          <button
            class="flex p-2 w-full"
            onclick={() => {
              selectTypeCollection("CSV");
            }}>CSV</button
          >
        </ul>
      </div>
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("documentTitle")}*`} />
      <InputRw onEdit={true} bind:value={tmpDocumentTitle} placeholder={$_("enterDocumentTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <InputRw
        onEdit={true}
        bind:value={tmpDocumentDescription}
        placeholder={$_("enterDocumentDescription")}
      />
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

<Modal title={$_("createNewCategories")} bind:open={addCategoriesModal}>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("categories")}`} />
      <!-- <MultiSelect items={categories} bind:value={tmpCategories}></MultiSelect> -->
      <!-- <Select
        class="mt-2"
        items={tmpCategories}
        placeholder={$_("selectCategory")}
      /> -->
      <!-- buat class disable -->
      <div class="flex">
        <button
          class="flex w-full flex-wrap gap-2 items-center border rounded-lg p-2"
          onclick={() => {
            focusInput(inputRef);
          }}
        >
          {#each tags as tag}
            <span
              class="flex items-center space-x-1 bg-utility-brand-100-light dark:bg-utility-brand-100-dark text-brand-secondary-light dark:text-brand-secondary-dark rounded-full px-2 py-1"
            >
              <span>{tag}</span>
              <button
                onclick={() => {
                  removeTag(tag);
                }}
                class="text-brand-primary-light dark:text-brand-primary-dark">x</button
              >
            </span>
          {/each}

          <input
            bind:this={inputRef}
            type="text"
            bind:value={newTag}
            placeholder={`${tags.length < 1 ? $_("enterNewCategory") : ""}`}
            onkeydown={handleKeyDown}
            class={`border-none flex-grow focus:ring-transparent ${tags.length > 0 ? "max-w-min" : "w-full"}`}
          />
        </button>
        <button
          onclick={addTag}
          class="bg-secondary-light dark:bg-secondary-dark text-primary-light dark:text-primary-dark px-4 py-1 rounded"
        >
          Add
        </button>
      </div>
      <!-- <MultiSelect
        class="no-pointer-events"
        items={tmpCategories}
        bind:value={tmpSelectedCategories}
      />
      <Input
        placeholder={$_("enterNewCategory")}
        bind:value={otherCategory}
        onkeydown={(e) => {
          if (
            e.key === "Enter" &&
            otherCategory !== "" &&
            !tmpCategories.find((category) => category.value === otherCategory)
          ) {
            // Tambahkan otherCategory ke tmpCategories
            tmpCategories = [...tmpCategories, { value: otherCategory, name: otherCategory }];
            tmpSelectedCategories = [...tmpSelectedCategories, otherCategory];
            otherCategory = "";
          }
        }}
      /> -->
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={() => addCategories()} disabled={tags.length === 0} class="w-full"
      >{$_("addNewCategories")}</Button
    >
  </svelte:fragment>
</Modal>

<!-- <LoadingContainer /> -->
<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteDocument} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>

<Modal bind:open={deleteCategoriesModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDeleteCategories")}
    </h3>
    <Button onclick={deleteCategories} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>

<Modal title={$_("trainingModel")} bind:open={trainingModelModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("newModelName")}*`} />
      <Input bind:value={newModelName} placeholder={$_("whatIsYourNewModelName")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={startTraining} disabled={newModelName === ""} class="w-full"
      >{$_("startTraining")}</Button
    >
  </svelte:fragment>
</Modal>
