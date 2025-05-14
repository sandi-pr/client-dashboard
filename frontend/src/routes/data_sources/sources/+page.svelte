<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { Button, Modal, Input, Alert } from "flowbite-svelte";
  import { format } from "date-fns";

  import {
    Typography,
    ToastContainer,
    StateCard,
    call,
    InputRw,
    triggerToast,
    triggerLoading,
    UploadFile,
    check,
    Divider,
    SelectedCard,
  } from "$lib";
  import DataTableCsv from "$lib/components/DataTableCsv.svelte";
  import { SearchOutline, ExclamationCircleSolid, InfoCircleSolid } from "flowbite-svelte-icons";
  import Trash01Icon from "$lib/icon/trash-01.svelte";
  import FilePlus02Icon from "$lib/icon/file-plus-02.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import Copy01Icon from "$lib/icon/copy-01.svelte";

  const isFrontend = process.env.NODE_ENV === "development";

  let subCollectionId;
  let collectionId;
  let libraryId;
  let checkAllState = $state(false);
  let addDocumentsModal = $state(false);
  let popupModal = $state(false);
  let updateModal = $state(false);
  let deleteModal = $state(false);
  let dropdownOpen = $state(false);
  let fileNotSupport = $state(false);
  let buttonDisabled = $state(false);
  let gettingDatas = $state(false);
  let gettingDocs = $state(false);
  let buttonDisabledCheckAction = $state(false);
  let pagination = $state({ page: 1, pageSize: 20, pageCount: 1, total: 3 });
  let paginationCollection = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpCollectionName = $state("");
  let tmpDocumentType = $state("Documents");
  let acceptedFileTypes = $state(".csv");
  let search = $state("");
  let tags = $state([]);
  let newTag = $state("");
  let files = $state([]);
  let checkboxDelete = $state(false);
  let step = $state(1);
  let linkCsv = $state("");
  let tmpId = $state(null);
  let tabDatas = $state({
    selected: "Collections",
    data: [
      {
        name: "Collections",
      },
    ],
  });

  let source = $state([
    {
      id: 1,
      name: "Upload CSV File",
    },
    {
      id: 2,
      name: "Upload with link",
    },
  ]);

  let selectedSource = $state(1);

  let collectionDatas = $state({
    title: "Data Sources Name",
    description: "",
    createdAt: "",
  });
  let collectionListDatas = $state([]);
  let csvListDatas = $state([
    // {
    //   id: 1,
    //   filename: 'user_1.csv',
    //   filesize: '200 KB',
    //   fileext: 'csv',
    //   checked: false,
    //   createdAt: 'Jan 4, 2022',
    // },
    // {
    //   id: 2,
    //   filename: 'user_2.csv',
    //   filesize: '200 KB',
    //   fileext: 'csv',
    //   checked: false,
    //   createdAt: 'Jan 4, 2022',
    // },
    // {
    //   id: 3,
    //   filename: 'user_3.csv',
    //   filesize: '200 KB',
    //   fileext: 'csv',
    //   checked: false,
    //   createdAt: 'Jan 4, 2022',
    // },
  ]);
  let childCollectionListDatas = $state([]);
  let childCollectionList = $state([]);
  let collectionList = $state([]);
  let tmpSelectedCollection = $state(null);
  let previewModal = $state(false);
  let selectedFile = $state(null);
  let dataSend = $state(null);
  let textarea = $state();
  let nameIsAvailable = $state(true);
  let isProcess = $state(false);

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
            parent_id: parseInt(subCollectionId),
            search: search,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

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

        // console.log("documets", response.result.data);

        pagination = response.result.meta.pagination;

        let link = isFrontend
          ? "/data_sources/collections/documents/sub_documents/detail"
          : await $call.getPage("client-dashboard", "document-library-sub-documents-detail");

        let promise = response.result.data.map(async (collection) => {
          const params = new URLSearchParams(window.location.search);

          params.delete("documentId");

          let uploadBy = "";

          try {
            let response = await $call.call("client-dashboard", "get-user", {
              URLQuery: {},
              URLParams: { id: collection.attributes.meta.upload_by.toString() },
              Params: {},
              Headers: {
                Authorization: "Bearer " + token,
              },
            });

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

            uploadBy = response.result.user.username;
          } catch (error) {
            console.error("Error during user get:", error);
            gettingDocs = false;
            triggerToast($_("failedGetUser"), true, "error");
          }

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
            upload_by: uploadBy,
            filesize: collection.attributes.meta.file_size,
            createdAt: format(new Date(collection.attributes.createdAt), "MMM d, yyyy"),
            href: link + "?documentId=" + collection.id + "&" + queryString,
            checked: false,
          };
        });

        collectionListDatas = await Promise.all(promise);
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

  const checkSingle = (index) => {
    csvListDatas = csvListDatas.map((data, i) => {
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
    for (const data of csvListDatas) {
      if (data.checked === false) {
        findCheckAll = false;
        break;
      }
    }

    if (findCheckAll) {
      csvListDatas = csvListDatas.map((data) => {
        return {
          ...data,
          checked: false,
        };
      });
    } else {
      csvListDatas = csvListDatas.map((data) => {
        return {
          ...data,
          checked: true,
        };
      });
    }
  };

  const nextPage = (evt) => {
    console.log("next page");
  };

  const prevPage = (evt) => {
    console.log("prev page");
  };

  const handleSearch = (evt) => {
    console.log("clicked");
    getDocuments();
  };

  async function addDocumentByCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        addDocumentsModal = false;
        triggerLoading(true);

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
              console.log("trocr-predict", response3);
            } else {
              response3 = await $call.call("hexa-ai-rag", "audio-predict", {
                URLQuery: {},
                URLParams: {},
                Params: {
                  audio_file: files[i].base64,
                  ext: files.name.split(".").pop(),
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
                document_name: files.length > 1 ? files[i].document_title : tmpDocumentTitle,
                description: tmpDocumentDescription,
                parent_id: parseInt(subCollectionId),
                documents: JSON.stringify([files[i]]),
                document_type: tmpDocumentType,
                category: tags.join(","),
                detected_text: response3.result.text,
                result_trocr: JSON.stringify(response3.result),
                edited: false,
              };

              triggerLoading(false);
            }
          }

          if (tmpDocumentType === "Handwriting" || tmpDocumentType === "Document Images") {
            if (files.length > 1) {
              await sendDocument();
            } else {
              previewModal = true;
              selectedFile = files[0];
            }
          } else {
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
              parent_id: parseInt(subCollectionId),
              documents: JSON.stringify([files[i]]),
              document_type: tmpDocumentType,
              category: tags.join(","),
            };
            await sendDocument();
          }
        }
        files = [];

        tmpDocumentTitle = "";
        tmpDocumentDescription = "";
        tags = [];
        attributeNotUnique = false;
        triggerLoading(false);

        if (
          tmpDocumentType !== "Handwriting" &&
          tmpDocumentType !== "Document Images" &&
          tmpDocumentType !== "Document Audio"
        ) {
          dataSend = null;
          previewModal = false;
          addDocumentsModal = false;
          getDocuments();
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

  function addDataSource() {
    addDocumentsModal = false;
    csvListDatas = [
      ...csvListDatas,
      {
        id: csvListDatas.length + 1,
        filename: files[0].name,
        filesize: files[0].size,
        fileext: files[0].ext,
        checked: false,
        createdAt: format(new Date(), "MMM d, yyyy"),
      },
    ];

    localStorage.setItem("csvListDatas", JSON.stringify(csvListDatas));

    files = [];

    triggerToast($_("succesfullyUploaded"), true);
  }

  function deleteSource() {
    if (checkboxDelete === false) {
      csvListDatas = csvListDatas.filter((item) => item.id !== tmpId);
    } else {
      csvListDatas = csvListDatas.filter((item) => !item.checked);
    }

    localStorage.setItem("csvListDatas", JSON.stringify(csvListDatas));

    tmpId = null;

    triggerToast($_("succesfullyDeleted"), true);
  }

  async function sendDocument(single = false) {
    let token = localStorage.getItem("token");

    if (token === null) {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
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
      getDocuments();
      triggerToast($_("succesfullyUploaded"), true);
    }
    return;
  }

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

      // console.log(fileTypeRegex.test(file.name));
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

          if (files.length === 1) {
            const status = await checkLibraryName(tmpDocumentTitle.trim());
            if (status === true) {
              nameIsAvailable = true;
            } else if (status === false) {
              nameIsAvailable = false;
            }
          }
          files = [...files];
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

        triggerToast($_("successDeleteDocument"), true);
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
      triggerToast($_("successDeleteDocument"), true);
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
        link = isFrontend
          ? "/data_sources/collections/documents/sub_documents"
          : await $call.getPage("client-dashboard", "document-library-sub-documents");
      } else {
        link = isFrontend
          ? "/data_sources/collections/documents"
          : await $call.getPage("client-dashboard", "document-library-documents");
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
        href: link + "?collectionId=" + collection.id + "&" + queryString,
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

        getCollections();
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

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
  };

  onMount(() => {
    csvListDatas = JSON.parse(localStorage.getItem("csvListDatas") || "[]");
  });

  $effect(() => {
    checkAllState = csvListDatas.every((data) => data.checked === true);
  });

  $effect(() => {
    buttonDisabledCheckAction = !csvListDatas.some((data) => data.checked === true);
  });

  $effect(() => {
    buttonDisabled = files.length !== 0 || linkCsv !== "" ? false : true;
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between mb-3 md:mb-8">
          <div class="flex flex-col gap-2">
            <div class="flex gap-4">
              <Typography
                preset="Display sm/Semibold"
                color="primary"
                text={collectionDatas.title}
                propClass="fluid-xl"
              />
              <button
                onclick={() => {
                  navigator.clipboard
                    .writeText(collectionDatas.title)
                    .then(() => {
                      console.log("Text copied to clipboard");
                      triggerToast("copyToClipboard");
                    })
                    .catch((err) => {
                      console.error("Failed to copy: ", err);
                    });
                }}
                class="p-2 border-2 rounded-lg border-secondary-light dark:border-secondary-dark"
              >
                <Copy01Icon
                  svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
                />
              </button>
            </div>
          </div>

          <div class="flex relative gap-3">
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
                class="p-2 md:px-3 md:py-2 gap-2 text-fg-secondary-light dark:text-fg-secondary-dark border-secondary-light dark:border-secondary-dark hover:bg-secondary-light hover:dark:bg-secondary-dark hover:text-fg-secondary-light dark:hover:text-fg-secondary-dark"
                outline
              >
                <Trash01Icon svgClass="w-5 h-5" />
                <div class="hidden md:block">{$_("delete")}</div>
              </Button>
            {/if}
            <Button
              disabled={!check("documentAi.documentLibrary.document:create")}
              outline
              on:click={() => {
                addDocumentsModal = true;
                step = 1;
              }}
              class="p-2 md:px-3 md:py-2 gap-2 whitespace-nowrap"
            >
              <FilePlus02Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("addDocument")}</div>
            </Button>
          </div>
        </div>

        <Divider />

        {#if csvListDatas.length > 0}
          <DataTableCsv
            collectionListDatas={csvListDatas}
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
        {:else if gettingDocs}
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

<Modal title={$_("addNewCsv")} bind:open={addDocumentsModal} size="xs">
  <div class="flex flex-col gap-4">
    {#if step === 0}
      <SelectedCard selectionData={source} bind:selectedData={selectedSource} isFlex isRadio />
    {:else if step === 1}
      {#if selectedSource === 1}
        <UploadFile
          title={`${$_("uploadDocuments")}*`}
          handleFileUploadWrapper={() => handleFileUploadWrapper()}
          file={files[0]}
          {files}
          bind:isProcess
          multiple={false}
          isFlex
          actionButton={true}
          actionDelete={() => (files = [])}
        />
        {#if fileNotSupport}
          <Alert color="red">
            <InfoCircleSolid slot="icon" class="w-5 h-5" />
            {$_("fileNotSupport")}
          </Alert>
        {/if}
      {:else if selectedSource === 2}
        <div class="flex flex-col gap-1.5 col-span-4">
          <Typography preset="Text sm/Medium" text="Import Link" color="secondary" />
          <InputRw bind:value={linkCsv} onEdit />
        </div>
      {/if}
    {/if}
  </div>
  <svelte:fragment slot="footer">
    {#if step === 0}
      <Button
        on:click={() => {
          step = 2;
        }}
        class="w-full"
        disabled={!selectedSource}>{$_("next")}</Button
      >
    {:else if step === 1}
      <div class="flex gap-4 w-full">
        <!-- <button
          onclick={() => {
            step = 1;
            linkCsv = "";
            files = [];
          }}
          class="px-5">{$_("back")}</button
        > -->
        <Button onclick={addDataSource} disabled={buttonDisabled} class="w-full"
          >{$_("upload")}</Button
        >
      </div>
    {/if}
  </svelte:fragment>
</Modal>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteSource} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <p>{tmpCollectionName}</p>
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteDocumentCollection} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>
