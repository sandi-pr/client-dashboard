<script>
  import { onMount, tick } from "svelte";
  import { _ } from "svelte-i18n";
  import { Button, Modal, Input, Alert, Select, P, Helper, Card } from "flowbite-svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { format, parse } from "date-fns";

  import {
    Typography,
    ToastContainer,
    StateCard,
    call,
    InputRw,
    TextareaRw,
    triggerToast,
    triggerLoading,
    UploadFile,
    check,
    onProgress,
    Divider,
    SectionHeader,
    bakrieLibraryId,
    LibraryCard,
    isBakrie,
  } from "$lib";
  import DataTableLibraryCollection from "$lib/components/DataTableLibraryCollection.svelte";
  import DataSourcesCard from "$lib/ui/cards/DataSourcesCard.svelte";
  import DataTableCsv from "$lib/components/DataTableCsv.svelte";
  import DataTableLibraryCollectionVertical from "$lib/components/DataTableLibraryCollectionVertical.svelte";
  import {
    SearchOutline,
    ExclamationCircleSolid,
    ChevronDownOutline,
    InfoCircleSolid,
    ArrowUpOutline,
  } from "flowbite-svelte-icons";
  import AudioPlayer from "$lib/components/AudioPlayer.svelte";
  import Trash01Icon from "$lib/icon/trash-01.svelte";
  import MoveIcon from "$lib/icon/move.svelte";
  import FilePlus02Icon from "$lib/icon/file-plus-02.svelte";
  import FolderPlusIcon from "$lib/icon/folder-plus.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import XCloseIcon from "$lib/icon/x-close.svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import AlertCircleIcon from "$lib/icon/alert-circle.svg";
  import MultiSelect from "svelte-multiselect";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import { get } from "svelte/store";

  let collectionId;
  let libraryId = $state(null);
  let inputRef;
  let checkAllState = $state(false);
  let addDocumentsModal = $state(false);
  let addCollectionModal = $state(false);
  let addCategoriesModal = $state(false);
  let popupModal = $state(false);
  let updateModal = $state(false);
  let editModal = $state(false);
  let deleteModal = $state(false);
  let deleteCategoriesModal = $state(false);
  let dropdownOpen = $state(false);
  let dropdownMore = $state(false);
  let editCollectionModal = $state(false);
  let dropdownUpload = $state(false);
  let fileNotSupport = $state(false);
  let buttonDisabled = $state(false);
  let gettingDatas = $state(false);
  let gettingDocs = $state(false);
  let buttonDisabledCheckAction = $state(false);
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });
  let paginationCollection = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpCollectionName = $state("");
  let tmpOldName = $state("");
  let tmpDescription = $state("");
  let tmpDocumentType = $state("Data");
  let selectedTypeCollection = $state("Data (CSV)");
  let acceptedFileTypes = $state(".csv");
  let search = $state("");
  let alltags = $state([]);
  let tags = $state([]);
  let newTag = $state("");
  let files = $state([]);
  let attributeNotUnique = $state(false);
  let checkboxDelete = $state(false);
  let checkboxUpdate = $state(false);
  let tmpId = $state(null);
  let isMobileView = $state(false);
  let hasDocOrDocx = $state(false);
  let tabDatas = $state({
    selected: "Collections",
    data: [
      {
        name: "Collections",
      },
    ],
  });

  let csvListDatas = $state([
    {
      id: 1,
      filename: 'user_1.csv',
      filesize: '200 KB',
      fileext: 'csv',
      type: 'CSV',
      checked: false,
      createdAt: 'Jan 4, 2022',
    },
    {
      id: 2,
      filename: 'user_2.csv',
      filesize: '200 KB',
      fileext: 'csv',
      type: 'CSV',
      checked: false,
      createdAt: 'Jan 4, 2022',
    },
    {
      id: 3,
      filename: 'user_3.csv',
      filesize: '200 KB',
      fileext: 'csv',
      type: 'CSV',
      checked: false,
      createdAt: 'Jan 4, 2022',
    },
  ]);

  let csvListData = $state([]);

  let collectionDatas = $state({
    title: "Sample Data",
    description: "",
    createdAt: "",
  });
  let collectionListDatas = $state([]);
  let childCollectionListDatas = $state([]);
  let childCollectionList = $state([]);
  let subDataSources = $state([]);
  let collectionList = $state([]);
  let tmpSelectedCollection = $state(null);
  let previewModal = $state(false);
  let selectedFile = $state(null);
  let dataSend = $state({ detected_text: "" });
  let textarea = $state();
  let div = $state();
  let nameIsAvailable = $state(true);
  let isProcess = $state(false);
  let selected = $state([]);
  let dummyTags = $state(["general", "news", "article"]);

  async function getCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "get-library-collection", {
        URLQuery: {},
        URLParams: {
          id: collectionId,
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      console.log("collection", response);

      collectionDatas = {
        title: response.result.attributes.name,
        description: response.result.attributes?.meta?.description || "",
        createdAt: format(new Date(response.result.attributes.createdAt), "dd-MM-yyyy"),
      };
    }
  }

  async function getDocuments() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDocs = true;
      getCollection();

      if (search !== "") {
        pagination.page = 1;
      }

      try {
        let response = await $call.call("client-dashboard", "get-library-documents", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: pagination.page,
            page_size: pagination.pageSize,
            parent_id: parseInt(collectionId),
            search: search,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        collectionListDatas = [];

        if (response.status !== 200) {
          console.log(response);
          gettingDocs = false;
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

        pagination = response.result.meta.pagination;

        let link = await $call.getPage("client-dashboard", "data-sources-detail");

        let promise = response.result.data.map(async (collection) => {
          const params = new URLSearchParams(window.location.search);

          params.delete("documentId");

          // let uploadBy = "";

          // try {
          //   let response = await $call.call("client-dashboard", "get-user", {
          //     URLQuery: {},
          //     URLParams: { id: collection.attributes.meta.upload_by.toString() },
          //     Params: {},
          //     Headers: {
          //       Authorization: "Bearer " + token,
          //     },
          //   });

          //   if (response.status !== 200) {
          //     console.log(response);
          //     gettingDocs = false;
          //     if (response.status === 400) {
          //       console.error("Invalid parameters:", response);
          //       triggerToast($_("invalidParameters"), true, "error");
          //     } else if (response.status === 401) {
          //       console.error("Unauthorized access:", response);
          //       triggerToast($_("unauthorized"), true, "error");
          //     } else if (response.status === 500) {
          //       console.error("Server error:", response);
          //       triggerToast($_("serverError"), true, "error");
          //     } else {
          //       console.error("Unexpected status code:", response);
          //       triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
          //     }
          //     return;
          //   }

          //   uploadBy = response.result.user.username;
          // } catch (error) {
          //   console.error("Error during user get:", error);
          //   gettingDocs = false;
          //   triggerToast($_("failedGetUser"), true, "error");
          // }

          const queryParams = {};
          for (const [key, value] of params.entries()) {
            queryParams[key] = value;
          }
          const queryString = new URLSearchParams(queryParams).toString();
          return {
            id: collection.id,
            filename: collection.attributes.name,
            fileext: collection.attributes.path.split(".").pop(),
            label:
              collection.attributes.meta.category === undefined
                ? []
                : collection.attributes.meta.category.split(","),
            type: collection.attributes.meta.document_type,
            upload_by: collection.attributes.meta.author,
            filesize: collection.attributes.meta.file_size,
            createdAt: format(new Date(collection.attributes.createdAt), "MMM d, yyyy"),
            href: link + "?documentId=" + collection.id + "&" + queryString,
            checked: false,
          };
        });

        collectionListDatas = await Promise.all(promise);

        console.log("collectionListDatas", collectionListDatas);
      } catch (error) {
        console.error("Error during documents get:", error);
        gettingDocs = false;
        triggerToast($_("failedGetDocuments"), true, "error");
      }
      gettingDocs = false;
    } else {
      triggerToast($_("unauthorized"), true, "error");
      gettingDocs = false;
      return;
    }
  }

  async function createCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let dataSend = {};
        if ($isBakrie) {
          dataSend = {
            collection_name: tmpCollectionName,
            description: tmpDescription,
            parent_id: parseInt(collectionId),
            // bakrie
            bot_id: localStorage.getItem("botIdStrapi").toString(),
          };
        } else {
          dataSend = {
            collection_name: tmpCollectionName,
            description: tmpDescription,
            parent_id: parseInt(collectionId),
          };
        }

        let response = await $call.call("client-dashboard", "create-library-collection", {
          URLQuery: {},
          URLParams: {},
          Params: dataSend,
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

        console.log("sub collection created", response);

        tmpCollectionName = "";
        tmpDescription = "";

        triggerToast($_("successCreateSubCollection"), true);

        addCollectionModal = false;
        getCollections(true);
      } catch (error) {
        console.error("Error during sub collection creation:", error);
        triggerToast($_("failedCreateSubCollection"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  function createDataSources() {
    subDataSources = [
      ...subDataSources,
      {
        collectionId: subDataSources.length + 1,
        collectionName: tmpCollectionName,
      },
    ];

    localStorage.setItem("subDataSources", JSON.stringify(subDataSources));

    tmpCollectionName = "";
    addCollectionModal = false;
  }

  let editIndex = $state(0);

  function editDataSourcesName() {
    subDataSources[editIndex].collectionName = tmpCollectionName;
    localStorage.setItem("subDataSources", JSON.stringify(subDataSources));
    editModal = false;

    triggerToast($_("succesfullyUpdated"), true);
  }

  function deleteDataSources(id) {
    subDataSources = subDataSources.filter(
      (item) => item.collectionId !== id,
    );
    localStorage.setItem("subDataSources", JSON.stringify(subDataSources));
    deleteModal = false;

    triggerToast($_("succesfullyDeleted"), true);
  }

  function deleteSource() {
    if (checkboxDelete === false) {
      csvListDatas = csvListDatas.filter((item) => item.id !== tmpId);
      console.log("csvListDatas", csvListDatas);
    } else {
      csvListDatas = csvListDatas.filter((item) => !item.checked);
    }

    localStorage.setItem("csvListDatas", JSON.stringify(csvListDatas));

    tmpId = null;

    triggerToast($_("succesfullyDeleted"), true);
  }

  function addDataSource() {
    addDocumentsModal = false;
    csvListDatas = [
      ...csvListDatas,
      ...files.map((file) => ({
        id: csvListDatas.length + 1 + files.indexOf(file),
        filename: file.name,
        filesize: file.size,
        fileext: file.ext,
        type: 'CSV',
        checked: false,
        createdAt: format(new Date(), "MMM d, yyyy"),
      })),
    ]

    localStorage.setItem("csvListDatas", JSON.stringify(csvListDatas));

    files = [];

    triggerToast($_("succesfullyUploaded"), true);
  }

  const checkSingle = (index) => {
    collectionListDatas = collectionListDatas.map((data, i) => {
      if (i === index) {
        return {
          ...data,
          checked: !data.checked,
        };
      }
      return data;
    });
  };

  const checkAll = () => {
    let findCheckAll = true;
    for (const data of collectionListDatas) {
      if (data.checked === false) {
        findCheckAll = false;
        break;
      }
    }

    if (findCheckAll) {
      collectionListDatas = collectionListDatas.map((data) => {
        return {
          ...data,
          checked: false,
        };
      });
    } else {
      collectionListDatas = collectionListDatas.map((data) => {
        return {
          ...data,
          checked: true,
        };
      });
    }
  };

  // const checkSingle = (index) => {
  //   csvListDatas = csvListDatas.map((data, i) => {
  //     if (i === index) {
  //       return {
  //         ...data,
  //         checked: !data.checked,
  //       };
  //     }
  //     return data;
  //   });
  // };

  // const checkAll = () => {
  //   let findCheckAll = true;
  //   for (const data of csvListDatas) {
  //     if (data.checked === false) {
  //       findCheckAll = false;
  //       break;
  //     }
  //   }

  //   if (findCheckAll) {
  //     csvListDatas = csvListDatas.map((data) => {
  //       return {
  //         ...data,
  //         checked: false,
  //       };
  //     });
  //   } else {
  //     csvListDatas = csvListDatas.map((data) => {
  //       return {
  //         ...data,
  //         checked: true,
  //       };
  //     });
  //   }
  // };

  // const nextPage = (evt) => {
  //   if (pagination.pageCount > pagination.page) pagination.page = pagination.page + 1;
  //   getDocuments();
  //   console.log("next page");
  // };

  // const prevPage = (evt) => {
  //   if (pagination.page > 1) pagination.page = pagination.page - 1;
  //   getDocuments();
  //   console.log("prev page");
  // };

  function nextPage() {
    if (pagination.page < pagination.pageCount) {
      pagination.page += 1;
    }
  }

  function prevPage() {
    if (pagination.page > 1) {
      pagination.page -= 1;
    }
  }

  const handleSearch = (evt) => {
    console.log("clicked");
    getDocuments();
  };

  let audioUrl = $state(null);

  async function addDocumentByCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        addDocumentsModal = false;
        triggerLoading(true);

        for (let i = 0; i < files.length; i++) {

          let index = 1;
          let originalName = files[i].name.split(".").slice(0, -1).join(".");
          let newName = originalName;

          for (let i = 0; i < index; i++) {
            let statusAvailable = await checkLibraryName(newName.trim());

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
          };

          // console.log("data send", dataSend);
          await sendDocument();
        }

        files = [];
        dataSend = null;

        tmpDocumentTitle = "";
        tmpDocumentDescription = "";
        tags = [];
        attributeNotUnique = false;
        getDocuments();
        triggerLoading(false);
        triggerToast($_("succesfullyUploaded"), true);
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
    }
    return;
  }

  const selectTypeCollection = (type) => {
    tmpDocumentType = type;

    if (type === "Documents") {
      selectedTypeCollection = "Documents (TXT, PDF, DOC, DOCX)";
      acceptedFileTypes = ".txt, .pdf, .doc, .docx";
    } else if (type === "Images") {
      selectedTypeCollection = "Images (JPG, JPEG, PNG, WEBP, IMG)";
      acceptedFileTypes = ".jpg, .jpeg, .png, .webp, .img";
    } else if (type === "Document Images") {
      selectedTypeCollection = "Document Images (JPG, JPEG, PNG, WEBP, IMG)";
      acceptedFileTypes = ".jpg, .jpeg, .png, .webp, .img";
    } else if (type === "Handwriting") {
      selectedTypeCollection = "Handwriting (PDF, JPG, JPEG, PNG, WEBP, IMG)";
      acceptedFileTypes = ".jpg, .jpeg, .png";
    } else if (type === "Document Audio") {
      selectedTypeCollection = "Document Audio (MP3, WAV, OGG)";
      acceptedFileTypes = ".mp3, .wav, .ogg";
    } else if (type === "CSV") {
      selectedTypeCollection = "CSV";
      acceptedFileTypes = ".csv";
    }

    files = [];

    dropdownOpen = false;
  };

  // Function to adjust the height of the textarea
  function autosize() {
    textarea.style.height = "auto"; // Reset height
    textarea.style.height = `${textarea.scrollHeight}px`; // Set height to scroll height

    textarea.addEventListener("input", (e) => {
      e.target.style.height = "auto";
      e.target.style.height = `${e.target.scrollHeight}px`;
    });
  }

  $effect(() => {
    if (textarea) {
      autosize();
    }
  });

  const handleFileUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = acceptedFileTypes;
    inputElement.multiple = true;

    inputElement.addEventListener("change", async (event) => {
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

          if (files.length === 1) {
            const status = await checkLibraryName(tmpDocumentTitle.trim());
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

  async function deleteDocuments() {
    let datasDelete = collectionListDatas.filter((data) => data.checked === true);
    popupModal = false;
    triggerLoading(true);

    for (const data of datasDelete) {
      const id = data.id;
      let token = localStorage.getItem("token");

      if (token === null) {
        triggerToast($_("unauthorized"), true, "error");
        return;
      }

      try {
        const response = await $call.call("client-dashboard", "delete-library-document", {
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
          console.log(response);
          triggerLoading(false);
          popupModal = true;
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

        triggerToast("successDeleteDocument", true);
      } catch (error) {
        console.error("Error during document deletion:", error);
        triggerToast($_("failedDeleteDocument"), true, "error");
      }
    }

    getDocuments();
    triggerLoading(false);
  }

  async function deleteDocument() {
    popupModal = false;
    triggerLoading(true);
    let token = localStorage.getItem("token");

    if (token === null) {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }

    try {
      const response = await $call.call("client-dashboard", "delete-library-document", {
        URLQuery: {},
        URLParams: {
          id: tmpId.toString(),
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        triggerLoading(false);
        popupModal = true;
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

      tmpId = null;

      triggerLoading(false);
      triggerToast("successDeleteDocument", true);
      getDocuments();
    } catch (error) {
      console.error("Error during document deletion:", error);
      triggerToast($_("failedDeleteDocument"), true, "error");
    }
  }

  async function updateDocument() {
    let datasUpdate = collectionListDatas.filter((data) => data.checked === true);
    updateModal = false;
    triggerLoading(true);

    for (const data of datasUpdate) {
      const id = data.id;
      let token = localStorage.getItem("token");
      if (token !== null) {
        let response;
        response = await $call.call("client-dashboard", "update-library-document-parent", {
          URLQuery: {},
          URLParams: {},
          Params: {
            document_id: parseInt(id),
            parent_id: parseInt(tmpSelectedCollection.id),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          console.log(response);
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

        triggerLoading(false);
        console.log(response);
        triggerToast("successUpdateDocument", true);
        getDocuments();
      } else {
        triggerToast($_("unauthorized"), true, "error");
        return;
      }
    }
  }

  async function getCollections(isChild = false, subCollections = false) {
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
            parent_id: parseInt(
              isChild ? collectionId : subCollections ? tmpSelectedCollection.id : libraryId,
            ),
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

      let link;

      if (isChild) {
        link = await $call.getPage("client-dashboard", "document-library-sub-documents");
      } else {
        link = await $call.getPage("client-dashboard", "document-library-documents");
      }

      const params = new URLSearchParams(window.location.search);
      params.delete("collectionId");

      const queryParams = {};
      for (const [key, value] of params.entries()) {
        queryParams[key] = value;
      }
      const queryString = new URLSearchParams(queryParams).toString();

      const promise = aggregatedData.map(async (collection) => ({
        id: collection.id,
        collectionName: collection.attributes.name,
        description: collection.attributes?.meta?.description || "",
        type: collection.attributes.type,
        href:
          link +
          "?subCollectionId=" +
          collection.id +
          "&" +
          "collectionId=" +
          collectionId +
          "&" +
          queryString,
      }));

      if (isChild) {
        childCollectionList = await Promise.all(promise);
        console.log("ischild", childCollectionList);
      } else if (subCollections) {
        childCollectionListDatas = await Promise.all(promise);
        console.log("subcollection", childCollectionListDatas);
        if (childCollectionListDatas.length > 0) {
          if (tabDatas.data.length < 2) {
            tabDatas.data = [
              ...tabDatas.data,
              {
                name: "Sub Collections",
              },
            ];
          }
        } else {
          tabDatas.data = tabDatas.data.filter((item) => item.name !== "Sub Collections");
        }
      } else {
        console.log("collections", collectionList);
        collectionList = await Promise.all(promise);
      }
      gettingDatas = false;
    } else {
      triggerToast($_("unauthorized"), true, "error");
      gettingDatas = false;
      return;
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
        getCollection();
        getCollections(true);
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
        let response = await $call.call("client-dashboard", "delete-library-collection", {
          URLQuery: {},
          URLParams: {
            id: tmpCollectionId.toString(),
          },
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

        triggerToast($_("successDeleteCollection"), true);

        getCollections(true);
        deleteModal = false;
      } catch (error) {
        console.error("Error during collection deletion:", error);
        triggerToast($_("failedDeleteCollection"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  async function deleteCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "delete-library-collection", {
          URLQuery: {},
          URLParams: {
            id: tmpCollectionId.toString(),
          },
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

        triggerToast($_("successDeleteCollection"), true);

        deleteModal = false;
        const libraryIdWebchatBakrie = window.localStorage.getItem("libraryIdWebchatBakrie");
        let link = await $call.getPage("client-dashboard", "document-library");
        goto(`${link}?libraryId=${libraryIdWebchatBakrie}`);
      } catch (error) {
        console.error("Error during collection deletion:", error);
        triggerToast($_("failedDeleteCollection"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  function checkTitle(title) {
    console.log(libraryId, $bakrieLibraryId);
    if (parseInt(libraryId) === parseInt($bakrieLibraryId)) {
      console.log("check title", true);
      return (
        title === "Documents" ||
        title === "Images" ||
        title === "Document Images" ||
        title === "Handwriting"
      );
    } else {
      console.log("check title", false);
      return false;
    }
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

  const selectCollection = (item, subcollection = false) => {
    tmpSelectedCollection = item;
    if (!subcollection) {
      getCollections(false, true);
    }

    console.log("selected collection", tmpSelectedCollection);
  };

  const checkMobileView = () => {
    isMobileView = window.innerWidth <= 768;
  };

  const scrollToTop = () => {
    console.log("Scrolling to top using scrollIntoView");

    // Use an existing element at the top of the page
    const topElement = document.querySelector(".top-element");
    if (topElement) {
      topElement.scrollIntoView({ behavior: "smooth" });
    } else {
      console.error("Top element not found");
    }
  };

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
    dropdownUpload = false;
    dropdownMore = false;
  };

  onMount(() => {
    if (check("documentAi.documentLibrary:view") === false) {
      window.history.back();
      return;
    }
    const unsubscribe = page.subscribe(async ($page) => {
      collectionId = $page.url.searchParams.get("collectionId");
      libraryId = $page.url.searchParams.get("libraryId");
      if (
        collectionId === null ||
        collectionId === undefined ||
        libraryId === null ||
        libraryId === undefined
      ) {
        goto(await $call.getPage("client-dashboard", "data-sources"));
        return;
      }
      getCollections(true, false);
      getCollections(false, false);
      getDocuments();
      getLabels();
    });

    // subDataSources = JSON.parse(localStorage.getItem("subDataSources")) || [];
    // csvListDatas = JSON.parse(localStorage.getItem("csvListDatas")) || [];
    
    // collectionDatas.title = $page.url.searchParams.get("collectionName");

    checkMobileView();
    window.addEventListener("resize", checkMobileView);

    return () => {
      unsubscribe();
      window.removeEventListener("resize", checkMobileView);
    };
  });

  $effect(() => {
    hasDocOrDocx = files.some(file => file.ext === '.doc' || file.ext === '.docx');
  })

  $effect(() => {
    checkAllState = collectionListDatas.every((data) => data.checked === true);
  });

  $effect(() => {
    buttonDisabledCheckAction = !collectionListDatas.some((data) => data.checked === true);
  });

  // $effect(() => {
  //   pagination.total = csvListDatas.length;
  //   pagination.pageCount = Math.ceil(pagination.total / pagination.pageSize);

  //   let start = (pagination.page - 1) * pagination.pageSize;
  //   let end = start + pagination.pageSize;

  //   csvListData = csvListDatas.slice(start, end);
  // });

  // $effect(() => {
  //   checkAllState = csvListDatas.every((data) => data.checked === true);
  // });

  // $effect(() => {
  //   buttonDisabledCheckAction = !csvListDatas.some((data) => data.checked === true);
  // });

  $effect(() => {
    buttonDisabled =
      files.length === 0 || files.length === 1
        ? tmpDocumentTitle === "" || !nameIsAvailable
        : false;
  });
</script>

<div class="top-element"></div>
<!-- Add this element at the top of the page -->
<div
  bind:this={div}
  class="flex dark:bg-primary-dark overflow-auto md:overflow-y-hidden md:hover:overflow-y-auto"
>
  {#if isMobileView}
    <button
      class="fixed bottom-5 right-5 z-40 bg-brand-primary-light dark:bg-brand-primary-dark text-white rounded-full w-12 h-12 flex items-center justify-center shadow-lg"
      onclick={scrollToTop}
    >
      <ArrowUpOutline class="w-6 h-6" />
    </button>
  {/if}
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-2 mb-3 md:mb-8">
          <div class="flex flex-col gap-2">
            <div class="flex gap-4">
              <Typography
                preset="Display sm/Semibold"
                color="primary"
                text={collectionDatas.title}
                propClass="fluid-xl line-clamp-2"
              />
            </div>
          </div>

          <div class="flex relative gap-3" onfocusout={handleDropdownFocusLoss}>
            <Input
              id="search"
              type="search"
              class="bg-white dark:bg-hc-graydark-950"
              placeholder={$_("search")}
              bind:value={search}
              on:input={handleSearch}
            >
              <SearchOutline
                slot="left"
                class="w-5 h-5 text-quinary-light dark:text-quinary-dark"
              />
            </Input>
            {#if !buttonDisabledCheckAction}
              <Button
                on:click={() => {
                  popupModal = true;
                  checkboxDelete = true;
                  // id
                }}
                disabled={!check("documentAi.documentLibrary.document:delete")}
                class="p-2 md:px-3 md:py-2 gap-2 text-error-primary-light dark:text-error-primary-dark border-button-secondary-light dark:border-button-secondary-dark hover:bg-secondary-light hover:dark:bg-secondary-dark hover:text-fg-secondary-light dark:hover:text-fg-secondary-dark"
                outline
              >
                <Trash01Icon svgClass="w-5 h-5 text-error-primary-light dark:text-error-primary-dark" />
                <div class="hidden md:block">{$_("delete")}</div>
              </Button>
              <!-- <Button
                on:click={() => {
                  updateModal = true;
                  checkboxUpdate = true;
                  paginationCollection.page = 1;
                  // id
                }}
                disabled={!check("documentAi.documentLibrary.document:update")}
                class="p-2 md:px-3 md:py-2 gap-2 text-fg-secondary-light dark:text-fg-secondary-dark border-secondary-light dark:border-secondary-dark hover:bg-secondary-light hover:dark:bg-secondary-dark hover:text-fg-secondary-light dark:hover:text-fg-secondary-dark"
                outline
              >
                <MoveIcon svgClass="w-5 h-5" />
                <div class="hidden md:block">{$_("move")}</div>
              </Button> -->
            {/if}
            <Button
              color="light"
              class="min-w-max gap-1.5"
              disabled={!check("documentAi.documentLibrary.document:create")}
              on:click={() => (dropdownUpload = !dropdownUpload)}
            >
              <FilePlus02Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("addNew")}</div>
            </Button>
            <div
              class={`z-50 ${dropdownUpload ? "absolute right-0 top-12" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}
            >
              <ul class="py-2 text-sm text-gray-700 dark:text-gray-200">
                <button
                  class="flex p-2 w-full gap-1.5"
                  disabled={!check("documentAi.documentLibrary.document:create")}
                  onclick={() => {
                    addCollectionModal = true;
                    tmpCollectionName = "";
                    tmpDescription = "";
                    attributeNotUnique = false;
                    dropdownUpload = false;
                    nameIsAvailable = true;
                  }}
                  ><FolderPlusIcon
                    svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark"
                  />{$_("addSubCollection")}</button
                >
                <button
                  class="flex p-2 w-full gap-1.5"
                  onclick={() => {
                    addDocumentsModal = true;
                    tags = [];
                    files = [];
                    tmpDocumentTitle = "";
                    tmpDocumentDescription = "";
                    attributeNotUnique = false;
                    dropdownUpload = false;
                    nameIsAvailable = true;
                  }}
                  ><FilePlus02Icon
                    svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark"
                  />{$_("addCsv")}</button
                >
              </ul>
            </div>
            {#if !checkTitle(collectionDatas.title)}
              {#if check("documentAi.documentLibrary.collection:edit") || check("documentAi.documentLibrary.collection:delete")}
                <div onfocusout={handleDropdownFocusLoss} class=" flex items-center relative">
                  <Button on:click={() => (dropdownMore = !dropdownMore)} color="light" class="p-2.5">
                    <DotsVerticalIcon
                      svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark"
                    />
                  </Button>
                  <div
                    class={`z-10 ${dropdownMore ? "absolute top-10 right-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
                  >
                    <ul class="py-2 grow text-sm">
                      {#if check("documentAi.documentLibrary.collection:edit")}
                        <button
                          class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                          onclick={() => {
                            dropdownMore = !dropdownMore;
                            tmpCollectionId = collectionId;
                            tmpCollectionName = collectionDatas.title;
                            tmpOldName = collectionDatas.title;
                            tmpDescription = collectionDatas.description;
                            nameIsAvailable = true;
                            editCollectionModal = true;
                          }}>{$_("editCollectionName")}</button
                        >
                      {/if}
                      {#if check("documentAi.documentLibrary.collection:delete")}
                        <button
                          class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                          onclick={() => {
                            dropdownMore = !dropdownMore;
                            tmpCollectionId = collectionId;
                            tmpCollectionName = collectionDatas.title;
                            deleteModal = true;
                          }}>{$_("delete")}</button
                        >
                      {/if}
                    </ul>
                  </div>
                </div>
              {/if}
            {/if}
          </div>
        </div>

        {#if childCollectionList.length > 0 || collectionListDatas.length > 0}
          <div class="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 my-4">
            {#each childCollectionList as document, index}
              <DataSourcesCard
                collectionName={document.collectionName}
                totalDocuments={document.totalDocuments}
                status={document.status}
                functEdit={() => {
                  tmpCollectionId = document.id;
                  tmpCollectionName = document.collectionName;
                  tmpDescription = document.description;
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
              />
            {/each}
          </div>
          <Divider />
          {#if collectionListDatas.length > 0}
            {#if isMobileView}
              <DataTableLibraryCollectionVertical
                {collectionListDatas}
                {pagination}
                {checkAllState}
                {checkSingle}
                {checkAll}
                next={nextPage}
                previous={prevPage}
                functDeleteCategories={(id) => {
                  tmpId = id;
                  popupModal = true;
                  checkboxDelete = false;
                }}
              />
            {:else}
              <DataTableLibraryCollection
                {collectionListDatas}
                {pagination}
                {checkAllState}
                {checkSingle}
                {checkAll}
                next={nextPage}
                previous={prevPage}
                functDeleteCategories={(id) => {
                  tmpId = id;
                  popupModal = true;
                  checkboxDelete = false;
                }}
              />
            {/if}
          {/if}
        {:else if gettingDatas || gettingDocs}
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
  <ToastContainer />
</div>

<Modal title={$_("addCsv")} bind:open={addDocumentsModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5 relative" onfocusout={handleDropdownFocusLoss}>
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("documentType")}*`} />
      <Button
        color="light"
        class="flex justify-between"
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
            }}>Documents (TXT, PDF, DOC, DOCX)</button
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
            }}>Document Images (JPG, JPEG, PNG, WEBP, IMG)</button
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
              selectTypeCollection("Document Audio");
            }}>Document Audio (MP3, WAV, OGG)</button
          >
          <!-- <button
            class="flex p-2 w-full"
            onclick={() => {
              selectTypeCollection("CSV");
            }}>CSV</button
          > -->
        </ul>
      </div>
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("labels")}`} />
      <MultiSelect
        --sms-min-height="40px"
        options={alltags}
        allowUserOptions="append"
        bind:selected={tags}
        autoScroll={true}
        noMatchingOptionsMsg=""
        createOptionMsg={null}
      />
    </div>
    <UploadFile
      title={`${$_("uploadDocuments")}*`}
      handleFileUploadWrapper={() => handleFileUploadWrapper()}
      file={files[0]}
      {files}
      bind:isProcess
      multiple={true}
      actionButton={true}
      actionDelete={() => (files = [])}
    />
    {#if fileNotSupport}
      <Alert color="red">
        <InfoCircleSolid slot="icon" class="w-5 h-5" />
        {$_("fileNotSupport")}
      </Alert>
    {:else if hasDocOrDocx}
      <Alert color="yellow">
        <InfoCircleSolid slot="icon" class="w-5 h-5" />
        {$_("docOrDocxWillConvertToPdf")}
      </Alert>
    {/if}
    <!-- <label class="inline-flex items-center cursor-pointer">
      <input type="checkbox" value="" class="hidden peer" bind:checked={bulk_upload} />
      <div
        class="relative w-11 h-6 bg-tertiary-light dark:bg-tertiary-dark peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-utility-brand-200-light dark:peer-focus:ring-utility-brand-200-dark rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-brand-solid-light"
      ></div>
      <span class="ms-3 text-sm font-medium text-gray-900 dark:text-gray-300"
        >{$_("bulkUpload")}</span
      >
    </label> -->
    {#if files.length === 1 && !isProcess}
      <div class="flex flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={`${$_("documentTitle")}*`} />
        <InputRw
          onEdit={true}
          bind:value={tmpDocumentTitle}
          placeholder={$_("enterDocumentTitle")}
          oninput={handleInput}
        />
        {#if !nameIsAvailable}
          <Helper color="red">
            <span class="font-medium">{$_("nameAlreadyTaken")}</span>
          </Helper>
        {/if}
      </div>
      <div class="flex flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}`} />
        <InputRw
          onEdit={true}
          bind:value={tmpDocumentDescription}
          placeholder={$_("enterDocumentDescription")}
        />
      </div>
    {/if}
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={addDocumentByCollection} disabled={buttonDisabled} class="w-full"
      >{$_("add")}</Button
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
      {checkboxDelete ? $_("theseDocuments") : $_("thisDocument")}
    </h3>
    <Button onclick={checkboxDelete ? deleteDocuments : deleteDocument} color="red" class="me-2"
      >{$_("yesImSure")}</Button
    >
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>

<Modal bind:open={updateModal}>
  <svelte:fragment slot="header">
    <div class="flex flex-col">
      <Typography preset="Text lg/Semibold" text={$_("moveDocument")} color="primary" />
      <Typography
        preset="Text sm/Regular"
        text={$_("whereDoYouWantToMoveThisDocument")}
        color="tertiary"
      />
    </div>
  </svelte:fragment>

  <div class="flex flex-col">
    <div class="bg-white dark:bg-gray-800 sticky -top-5">
      <TabDetail spaceTop="mt-0" {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />
    </div>
    {#if tabDatas.selected === "Collections"}
      <div class="flex flex-col justify-center mt-3 gap-3">
        <div class="flex flex-col gap-6">
          {#each collectionList as item, index (item)}
            <Card
              class="{tmpSelectedCollection === item
                ? 'border-brand-solid-light dark:border-brand-solid-dark ring-brand-solid-light dark:ring-brand-solid-dark'
                : ''}bg-hc-graylight-50 hover:bg-hc-brand-50 cursor-pointer gap-3"
              size="xl"
              on:click={() => selectCollection(item)}
            >
              <div class="flex flex-col gap-0.5">
                <P
                  weight="semibold"
                  class={`${tmpSelectedCollection === item ? "text-brand-primary-light dark:text-brand-primary-dark" : "text-primary-light dark:text-primary-dark"}`}
                  >{item.collectionName}</P
                >
                <P
                  size="xs"
                  class={`${tmpSelectedCollection === item ? "text-brand-secondary-light dark:text-brand-secondary-dark" : "text-secondary-light dark:text-secondary-dark"}`}
                  >{item.description}</P
                >
              </div>
            </Card>
          {/each}
        </div>
      </div>
    {:else if tabDatas.selected === "Sub Collections"}
      <div class="flex flex-col justify-center mt-3 gap-3">
        <div class="flex flex-col gap-6">
          {#each childCollectionListDatas as item, index (item)}
            <Card
              class="{tmpSelectedCollection === item
                ? 'border-brand-solid-light dark:border-brand-solid-dark ring-brand-solid-light dark:ring-brand-solid-dark'
                : ''}bg-hc-graylight-50 hover:bg-hc-brand-50 cursor-pointer gap-3"
              size="xl"
              on:click={() => selectCollection(item, true)}
            >
              <div class="flex flex-col gap-0.5">
                <P
                  weight="semibold"
                  class={`${tmpSelectedCollection === item ? "text-brand-primary-light dark:text-brand-primary-dark" : "text-primary-light dark:text-primary-dark"}`}
                  >{item.collectionName}</P
                >
                <P
                  size="xs"
                  class={`${tmpSelectedCollection === item ? "text-brand-secondary-light dark:text-brand-secondary-dark" : "text-secondary-light dark:text-secondary-dark"}`}
                  >{item.description}</P
                >
              </div>
            </Card>
          {/each}
        </div>
      </div>
    {/if}
  </div>
  <svelte:fragment slot="footer">
    <Button on:click={updateDocument} class="w-full" disabled={!tmpSelectedCollection}
      >{$_("move")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal bind:open={deleteCategoriesModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDeleteCategories")}
    </h3>
    <Button onclick={() => {}} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>

<Modal title={$_("createNewSubCollection")} bind:open={addCollectionModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("subCollectionName")}*`} />
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
      class="w-full">{$_("addNewSubCollection")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editSubCollectionName")} bind:open={editModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("subCollectionName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={handleInput}
        placeholder={$_("whatIsYourSubCollectionTitle")}
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
      {$_("collection")}
      "{tmpCollectionName}"?
    </h3>
    <Button
      onclick={tmpCollectionId === collectionId ? deleteCollection : deleteDocumentCollection}
      color="red"
      class="me-2">{$_("yesImSure")}</Button
    >
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>

<Modal bind:open={previewModal} size="lg">
  <svelte:fragment slot="header">
    <div class="flex justify-start items-center gap-4 py-4 px-6">
      <div class="w-8 h-7 p-2 bg-[#feefc6] rounded-full flex justify-center items-center">
        <img src={AlertCircleIcon} alt="AlertCircleIcon" class="w-5 h-5" />
      </div>
      <div class="flex flex-col">
        <Typography preset="Text lg/Semibold" text="Double-check" color="primary" />
        <Typography
          preset="Text sm/Regular"
          text="Are you sure this is the correct document that you want to upload?"
          color="tertiary"
        />
      </div>
    </div>
  </svelte:fragment>
  <div class="flex">
    <div
      class="flex w-full justify-center items-center py-6 rounded-lg bg-secondary-light dark:bg-secondary-dark"
    >
      {#if selectedFile.ext.includes(".png") || selectedFile.ext.includes(".jpg") || selectedFile.ext.includes(".jpeg")}
        <img
          src={`data:${selectedFile.type};base64,${selectedFile.base64}`}
          style="max-width: 80%; height: auto;"
          alt="PNG Preview"
        />
      {/if}
      {#if audioUrl}
        <AudioPlayer src={audioUrl} />
      {/if}
    </div>
    <!-- detected handwriting -->
    {#if selectedFile.type === "Handwriting" || selectedFile.type === "Document Images"}
      <div class="flex w-full flex-col gap-3 px-4 pb-1">
        <Typography
          preset="Text md/Semibold"
          text={selectedFile.type === "Handwriting"
            ? "Handwriting Information"
            : "Document Information"}
          color="secondary"
        />
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/Medium" text="Detected Text" color="secondary" />
          <textarea
            bind:this={textarea}
            bind:value={dataSend.detected_text}
            id="chat"
            rows="1"
            class="block p-2.5 w-full text-sm text-brand-primary-light dark:text-primary-dark bg-white dark:bg-primary-dark rounded-lg border border-utility-brand-500-light dark:border-utility-brand-500-dark focus:border-utility-brand-500-light dark:focus:utility-brand-500-dark focus:ring-utility-brand-400-light dark:focus:ring-utility-brand-400-light dark:placeholder-gray-400"
          ></textarea>
        </div>
      </div>
    {/if}
    {#if selectedFile.type === "Document Audio"}
      <div class="flex w-full flex-col gap-3 px-4 pb-1">
        <Typography
          preset="Text md/Semibold"
          text="Audio Information"
          color="secondary"
        />
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/Medium" text="Detected Text" color="secondary" />
          <textarea
            bind:this={textarea}
            bind:value={dataSend.detected_text}
            id="chat"
            rows="1"
            class="block p-2.5 w-full text-sm text-brand-primary-light dark:text-primary-dark bg-white dark:bg-primary-dark rounded-lg border border-utility-brand-500-light dark:border-utility-brand-500-dark focus:border-utility-brand-500-light dark:focus:utility-brand-500-dark focus:ring-utility-brand-400-light dark:focus:ring-utility-brand-400-light dark:placeholder-gray-400"
          ></textarea>
        </div>
      </div>
    {/if}
  </div>
  <svelte:fragment slot="footer">
    <div class="flex w-full gap-2 justify-center items-center py-3 px-4">
      <Button
        color="light"
        on:click={() => {
          previewModal = false;
          addDocumentsModal = true;
        }}
        class="w-full">Cancel</Button
      >
      <Button color="primary" on:click={sendDocument(true)} class="rounded-lg w-full"
        >I’m sure</Button
      >
    </div>
  </svelte:fragment>
</Modal>

<Modal title={$_("editCollectionName")} bind:open={editCollectionModal} size="xs" autoclose>
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
      onclick={editCollectionName}
      disabled={!tmpCollectionName || !nameIsAvailable}
      class="me-2 w-full">{$_("save")}</Button
    >
  </svelte:fragment>
</Modal>

<style>
  html,
  body {
    height: 100%;
    margin: 0;
    padding: 0;
    overflow: auto;
  }

  .content {
    min-height: 200vh; /* Ensure content is long enough to require scrolling */
  }
</style>
