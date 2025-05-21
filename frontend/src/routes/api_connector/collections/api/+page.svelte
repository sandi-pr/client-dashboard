<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { Button, Modal, Input, P, Helper, Card } from "flowbite-svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { format } from "date-fns";

  import {
    Typography,
    ToastContainer,
    StateCard,
    call,
    triggerToast,
    triggerLoading,
    check,
    Divider,
    bakrieLibraryId,
    isBakrie,
    ModalObject,
    ModalFolder,
    InputRw,
    getIcons,
  } from "$lib";
  import { ExclamationCircleSolid, ArrowUpOutline } from "flowbite-svelte-icons";
  import AudioPlayer from "$lib/components/AudioPlayer.svelte";
  import Trash01Icon from "$lib/icon/trash-01.svelte";
  import MoveIcon from "$lib/icon/move.svelte";
  import PenToolPlusIcon from "$lib/icon/pen-tool-plus.svelte";
  import FolderPlusIcon from "$lib/icon/folder-plus.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import AlertCircleIcon from "$lib/icon/alert-circle.svg";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import TabBar from "$lib/components/TabBar.svelte";

  const isFrontend = process.env.NODE_ENV === "development";

  let collectionId = $state(null);
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
  let dropdownMore = $state(false);
  let editCollectionModal = $state(false);
  let dropdownUploadEmpty = $state(false);
  let dropdownAuthorizationType = $state(false);
  let fileNotSupport = $state(false);
  let buttonDisabled = $state(false);
  let gettingDatas = $state(false);
  let gettingDocs = $state(false);
  let buttonDisabledCheckAction = $state(true);
  let pagination = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });
  let paginationCollection = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpCollectionName = $state("");
  let tmpOldName = $state("");
  let tmpDescription = $state("");
  let tmpDocumentType = $state("Handwriting");
  let selectedTypeCollection = $state("Handwriting (JPG, JPEG, PNG)");
  let acceptedFileTypes = $state(".jpg, .jpeg, .png");
  let search = $state("");
  let alltags = $state([]);
  let tags = $state([]);
  let newTag = $state("");
  let files = $state([]);
  let checkboxDelete = $state(false);
  let tmpId = $state(null);
  let isMobileView = $state(false);
  let hasDocOrDocx = $state(false);
  let tabDatas = $state({
    selected: "My API",
    data: [
      {
        name: "My API",
      },
      {
        name: "Authorization",
      },
      {
        name: "Information",
      },
    ],
  });

  let authorizationTypeList = $state({
    selected: "No Auth",
    data: [
      { name: "No Auth" },
      { name: "Inherit Auth from Parent" },
      { name: "Basic Auth" },
      { name: "Bearer Token" },
    ],
  });

  let tmpSelectedAuthorizationType = $state("No Auth");
  let authBasicAuth = $state({
    username: "",
    password: "",
  });
  let authBearerToken = $state({
    token: "",
  });

  let groupedFiles = $state(null);
  let addFolderModal = $state(false);

  let collectionDatas = $state({
    title: "",
    description: "",
    createdAt: "",
  });
  let collectionListDatas = $state([]);
  let childCollectionListDatas = $state([]);
  let childCollectionList = $state([]);
  let collectionList = $state([]);
  let tmpSelectedCollection = $state(null);
  let previewModal = $state(false);
  let selectedFile = $state(null);
  let dataSend = $state({ detected_text: "" });
  let textarea = $state();
  let nameIsAvailable = $state(true);
  let isProcess = $state(false);
  let folderMode = $state(false);
  let apisData = $state([
    {
      id: 1,
      // checkbox: {
      //   checked: false,
      // },
      name: "GetUserAPI",
      method: {
        type: "success",
        name: "GET",
      },
      apiLink: "https://hexacode.ai/api/character?name",
      dateUploaded: "Jan 4, 2022",
      actions: {
        objects: [
          {
            href: "#",
            name: "View",
          },
        ],
      },
    },
  ]);

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

      // console.log("collection", response);

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

        // console.log("documets", response.result.data);

        pagination = response.result.meta.pagination;

        let link = isFrontend
          ? "/handwriting/collections/documents/detail"
          : await $call.getPage("client-dashboard", "handwriting-detail");

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

        // console.log("collectionListDatas", collectionListDatas);
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

        // console.log(response);

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
    // console.log("cek single", collectionListDatas);
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

  const nextPage = (evt) => {
    if (pagination.pageCount > pagination.page) pagination.page = pagination.page + 1;
    getDocuments();
    // console.log("next page");
  };

  const prevPage = (evt) => {
    if (pagination.page > 1) pagination.page = pagination.page - 1;
    getDocuments();
    // console.log("prev page");
  };

  let audioUrl = $state(null);

  async function addDocumentByCollection(folderFiles = null, parentId = null) {
    if (folderFiles != null) {
      files = folderFiles;
    }

    // console.log("cekkkkk files", files);
    let token = localStorage.getItem("token");

    if (token !== null) {
      if (groupedFiles != null) {
        sendFolder();
        return;
      }

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
              // console.log("audio-predict", response3);
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
                let statusAvailable = await checkLibraryName(
                  newName.trim(),
                  collectionId,
                  "document",
                );

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

              // console.log("cekk dataSend", dataSend);

              // convert base64 to url
              if (tmpDocumentType === "Document Audio") {
                let file = base64toBlob(files[i].base64, "audio/mpeg");
                audioUrl = URL.createObjectURL(new Blob([file], { type: "audio/mp3" }));
              }
              // triggerLoading(false);
            }
          }

          if (
            tmpDocumentType === "Handwriting" ||
            tmpDocumentType === "Document Images" ||
            tmpDocumentType === "Document Audio"
          ) {
            // console.log("ayam");
            if (files.length > 1) {
              await sendDocument();
            } else {
              // console.log("bebek");
              previewModal = true;
              triggerLoading(false);
              selectedFile = files[0];
              return;
            }
          } else {
            // console.log("cacing");

            let index = 1;
            let originalName = files[i].name.split(".").slice(0, -1).join(".");
            let newName = originalName;

            for (let i = 0; i < index; i++) {
              let statusAvailable = await checkLibraryName(
                newName.trim(),
                collectionId,
                "document",
              );

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
              parent_id: !parentId ? parseInt(collectionId) : parseInt(parentId),
              documents: JSON.stringify([files[i]]),
              document_type: tmpDocumentType,
              category: tags.join(","),
            };

            // console.log("cekkkk dataSend", dataSend);

            // console.log("cek data send", dataSend);
            await sendDocument();
          }
        }

        files = [];

        tmpDocumentTitle = "";
        tmpDocumentDescription = "";
        tags = [];
        attributeNotUnique = false;
        dataSend = null;
        triggerLoading(false);
        getDocuments();
        triggerToast($_("succesfullyUploaded"), true);

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
      getDocuments();
      triggerLoading(false);
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

  $effect(() => {
    if (groupedFiles != null) {
      buttonDisabled = false;
    } else {
      buttonDisabled = true;
    }
  });

  const handleFileUploadWrapper = (isFolder = false) => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.webkitdirectory = isFolder;
    inputElement.accept = acceptedFileTypes;
    inputElement.multiple = true;

    inputElement.addEventListener("change", async (event) => {
      if (isFolder) {
        // Buat Collection by Name Folder Dan Document
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

          if (files.length === 1) {
            const status = await checkLibraryName(
              tmpDocumentTitle.trim(),
              collectionId,
              "document",
            );
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

  //#FUNCTION UNTUK TRIGGERTOAST DARI RESPONSE
  function triggerToastResponse(response) {
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
    return true;
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
        // console.log(response);
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
        link = isFrontend
          ? "/handwriting/collections/documents/sub_documents"
          : await $call.getPage("client-dashboard", "handwriting-sub-documents");
      } else {
        link = isFrontend
          ? "/handwriting/collections/documents"
          : await $call.getPage("client-dashboard", "handwriting-documents");
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
        // console.log("ischild", childCollectionList);
      } else if (subCollections) {
        childCollectionListDatas = await Promise.all(promise);
        // console.log("subcollection", childCollectionListDatas);
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
        // console.log("collections", collectionList);
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
        let link = isFrontend
          ? "/handwriting"
          : await $call.getPage("client-dashboard", "handwriting");
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
    if (parseInt(libraryId) === parseInt($bakrieLibraryId)) {
      return (
        title === "Documents" ||
        title === "Images" ||
        title === "Document Images" ||
        title === "Handwriting"
      );
    } else if (
      parseInt(libraryId) === parseInt(localStorage.getItem("libraryId" + "handwriting"))
    ) {
      return title === "My Handwriting";
    } else if (parseInt(libraryId) === parseInt(localStorage.getItem("libraryId" + "voice"))) {
      return title === "My Voice";
    } else {
      return false;
    }
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

      // console.log("cek resp", resp);

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

  const selectCollection = (item, subcollection = false) => {
    tmpSelectedCollection = item;
    if (!subcollection) {
      getCollections(false, true);
    }

    // console.log("selected collection", tmpSelectedCollection);
  };

  const checkMobileView = () => {
    isMobileView = window.innerWidth <= 768;
  };

  const scrollToTop = () => {
    // console.log("Scrolling to top using scrollIntoView");

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
    dropdownUploadEmpty = false;
    dropdownMore = false;
    dropdownAuthorizationType = false;
  };

  async function groupFilesByPath(files) {
    const structure = {};

    const fileTypeRegex = new RegExp(
      acceptedFileTypes
        .split(",")
        .map((type) => type.trim().replace(".", "\\."))
        .join("|") + "$",
      "i",
    );

    for (const file of files) {
      if (file && fileTypeRegex.test(file.name)) {
        const parts = file.webkitRelativePath.split("/");
        if (parts.length > 2) {
          continue;
        }
        const fileName = parts.pop();
        const base64Data = await readFileAsBase64(file);

        let folderKey = parts.slice(0, 1).join("/");
        if (!folderKey) folderKey = "root";

        if (parts.length > 2) {
          continue;
        }

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
          type: tmpDocumentType,
        };

        structure[folderKey][fileName] = fileObject;
      }
    }

    return structure;
  }

  const checkNewName = async (name, parentId, type) => {
    let index = 1;
    let originalName = name;
    let newName = originalName;

    for (let i = 0; i < index; i++) {
      let statusAvailable = await checkLibraryName(newName.trim(), parentId, type);

      if (statusAvailable === true) {
        break;
      } else {
        newName = `${originalName}(${index})`;
        index++;
      }
    }

    return newName;
  };

  async function sendFolder() {
    triggerLoading(true);
    var currentFolder = [];

    for (var folder in groupedFiles) {
      var parent_id;
      var sub_folder;
      var main_folder;

      if (folder.split("/").length > 1) {
        main_folder = folder.split("/")[0].trim();
        sub_folder = folder.split("/")[1].trim();

        if (currentFolder.length) {
          let finded = currentFolder.find((f) => f[main_folder] == main_folder);
          if (finded) {
            parent_id = finded.id;
          }
        }
      }

      let newName = await checkNewName(folder, collectionId, "collection");

      dataSend = {
        collection_name: newName,
        description: "",
        parent_id: parseInt(collectionId),
      };

      if (sub_folder && parent_id) {
        dataSend["collection_name"] = sub_folder;
        dataSend["parent_id"] = parseInt(parent_id);
      }

      let token = localStorage.getItem("token");
      let res = await $call.call("client-dashboard", "create-library-collection", {
        URLQuery: {},
        URLParams: {},
        Params: dataSend,
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      let statusRes = triggerToastResponse(res);
      if (!statusRes) {
        triggerLoading(false);
        return;
      }

      currentFolder.push(res.result.data);
      files = groupedFiles[folder];

      if (!parent_id && !sub_folder) {
        parent_id = res.result.data.id;
      }

      for (const doc in files) {
        if (sub_folder && parent_id) {
          files[doc].id = parent_id;
        }

        let response3;
        response3 = await $call.call("hexa-ai-rag", "trocr-predict", {
          URLQuery: {},
          URLParams: {},
          Params: {
            image: files[doc].base64,
            model:
              tmpDocumentType === "Document Images"
                ? "trocr:base"
                : sessionStorage.getItem("model_trocr"),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        let statusDetectedDoc = triggerToastResponse(response3);
        if (!statusDetectedDoc) {
          // triggerLoading(false);
          return;
        }

        if (response3?.result?.text !== null) {
          let index = 1;
          let originalName = files[doc].name.split(".").slice(0, -1).join(".");
          let newName = originalName;

          for (let i = 0; i < index; i++) {
            let statusAvailable = await checkLibraryName(
              newName.trim(),
              res.result.data.id,
              "document",
            );

            if (statusAvailable === true) {
              break;
            } else {
              newName = `${originalName}(${index})`;
              index++;
            }
          }

          files[doc].name = newName + files[doc].ext;

          dataSend = {
            document_name: newName,
            description: files[doc].description,
            parent_id: parseInt(res.result.data.id),
            document_type: files[doc].type,
            documents: JSON.stringify(Array(files[doc])),
            category: "",
            detected_text: response3.result.text,
            result_trocr: JSON.stringify(response3.result),
            edited: false,
          };
          // triggerLoading(false);
        }

        let resdoc = await $call.call("client-dashboard", "create-library-document", {
          URLQuery: {},
          URLParams: {},
          Params: dataSend,
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        let statusResDoc = triggerToastResponse(resdoc);
        if (!statusResDoc) {
          // triggerLoading(false);
          // addFolderModal = false;
          return;
        }
      }
    }

    groupedFiles = null;
    addDocumentsModal = false;
    getCollections(true);
    triggerLoading(false);

    triggerToast($_("succesfullyUploaded"), true);
  }

  onMount(async() => {
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
        goto(await $call.getPage("client-dashboard", "handwriting"));
        return;
      }
      getCollections(true, false);
      getCollections(false, false);
      getDocuments();
      getLabels();
    });

    const hrefDetailApi = isFrontend ? "/api_connector/collections/api/detail" : await $call.getPage("client-dashboard", "api-connector-collections-api-detail");

    apisData.forEach((api) => {
      api.actions.objects[0].href = hrefDetailApi + "?apiId=" + api.id + "&collectionId=" + collectionId + "&libraryId=" + libraryId;
    })

    checkMobileView();
    window.addEventListener("resize", checkMobileView);

    return () => {
      unsubscribe();
      window.removeEventListener("resize", checkMobileView);
    };
  });

  $effect(() => {
    checkAllState = collectionListDatas.every((data) => data.checked === true);
    buttonDisabledCheckAction = !collectionListDatas.some((data) => data.checked === true);
    if (files.length > 0) {
      hasDocOrDocx = files.some((file) => file.ext === ".doc" || file.ext === ".docx");
      buttonDisabled =
        files.length === 0 || files.length === 1
          ? tmpDocumentTitle === "" || !nameIsAvailable
          : false;
    } else if (groupedFiles != null) {
      buttonDisabled = false;
    } else {
      buttonDisabled = true;
    }
  });
</script>

<div class="top-element"></div>
<!-- Add this element at the top of the page -->
<div class="flex dark:bg-primary-dark overflow-auto md:overflow-y-hidden md:hover:overflow-y-auto">
  {#if isMobileView}
    <button
      class="fixed bottom-5 right-5 z-40 bg-brand-primary-light dark:bg-brand-primary-dark text-white rounded-full w-12 h-12 flex items-center justify-center shadow-lg"
      onclick={scrollToTop}
    >
      <ArrowUpOutline class="w-6 h-6" />
    </button>
  {/if}
  <div class="w-full pb-10 min-h-[30rem]">
    <div class="flex flex-col m-4">
      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-2 mb-4">
          <div class="flex flex-col gap-2">
            <div class="flex gap-4">
              <Typography
                preset="Display sm/Semibold"
                color="primary"
                text={collectionDatas.title}
                propClass="fluid-xl line-clamp-2"
              />
              <button
                onclick={() => {
                  navigator.clipboard
                    .writeText(collectionDatas.title)
                    .then(() => {
                      // console.log("Text copied to clipboard");
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
            <P
              class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
              >{collectionDatas.description}</P
            >
            <P
              class="text-sm font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
              >{$_("createOn")}{` ${collectionDatas.createdAt}`}</P
            >
          </div>

          <div class="flex relative gap-3" onfocusout={handleDropdownFocusLoss}>
            {#if !buttonDisabledCheckAction}
              <Button
                on:click={() => {
                  popupModal = true;
                  checkboxDelete = true;
                  // id
                }}
                disabled={!check("documentAi.documentLibrary.document:delete")}
                class="p-2 md:px-3 md:py-2 gap-2 text-error-primary-light dark:text-error-primary-dark hover:text-error-primary-light dark:hover:text-error-primary-dark border-secondary-light dark:border-secondary-dark hover:bg-secondary-light hover:dark:bg-secondary-dark"
                outline
              >
                <Trash01Icon svgClass="w-5 h-5" />
                <div class="hidden md:block">{$_("delete")}</div>
              </Button>
              <Button
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
              </Button>
            {/if}
            <Button
              outline
              class="min-w-max gap-1.5"
              disabled={!check("documentAi.documentLibrary.document:create")}
              on:click={async () => {
                let link = isFrontend
                  ? "/api_connector/collections/api/create" + `?collectionId=${collectionId}&libraryId=${libraryId}`
                  : await $call.getPage("client-dashboard", "api-connector-collections-api-create") + `?collectionId=${collectionId}&libraryId=${libraryId}`;
                goto(link);
              }}
            >
              <PenToolPlusIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("addNewApi")}</div>
            </Button>
            {#if check("documentAi.documentLibrary.collection:delete")}
              <div onfocusout={handleDropdownFocusLoss} class=" flex items-center relative">
                <Button on:click={() => (dropdownMore = !dropdownMore)} color="light" class="p-2.5">
                  <DotsVerticalIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                </Button>
                <div
                  class={`z-10 ${dropdownMore ? "absolute top-10 right-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
                >
                  <ul class="py-2 grow text-sm">
                    <!-- {#if check("documentAi.documentLibrary.collection:edit")}
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
                    {/if} -->
                    {#if check("documentAi.documentLibrary.collection:delete")}
                      <button
                        class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                        onclick={() => {
                          dropdownMore = !dropdownMore;
                          // tmpCollectionId = collectionId;
                          // tmpCollectionName = collectionDatas.title;
                          // deleteModal = true;
                        }}
                      >
                        {$_("delete")}
                      </button>
                    {/if}
                  </ul>
                </div>
              </div>
            {/if}
          </div>
        </div>

        <TabBar {tabDatas} />
        <Divider class="mb-4" />

        {#if tabDatas.selected === "My API"}
          {#if apisData.length > 0}
            <DataTable
              bind:datas={apisData}
              next={nextPage}
              previous={prevPage}
              {checkAll}
              {checkSingle}
            />
          {:else if gettingDatas || gettingDocs}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else}
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
            <div
              class="relative w-full justify-center flex mb-40"
              onfocusout={handleDropdownFocusLoss}
            >
              <Button
                color="primary"
                class="w-fit gap-1.5 mt-4"
                disabled={!check("documentAi.documentLibrary.document:create")}
                on:click={() => (dropdownUploadEmpty = !dropdownUploadEmpty)}
              >
                <PenToolPlusIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                <div class="hidden md:block">{$_("addNew")}</div>
              </Button>
              <div
                class={`z-50 ${dropdownUploadEmpty ? "absolute top-32 left-1/2 transform -translate-x-1/2 -translate-y-1/2" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}
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
                      dropdownUploadEmpty = false;
                      nameIsAvailable = true;
                    }}
                  >
                    <FolderPlusIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                    {$_("addSubCollection")}
                  </button>
                  <button
                    class="flex p-2 w-full gap-1.5"
                    disabled={!check("documentAi.documentLibrary.document:create")}
                    onclick={() => {
                      addDocumentsModal = true;
                      folderMode = true;
                      files = [];
                      groupedFiles = null;
                      attributeNotUnique = false;
                      dropdownUploadEmpty = false;
                      nameIsAvailable = true;
                    }}
                  >
                    <FolderPlusIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                    {$_("addFolder")}
                  </button>
                  <button
                    class="flex p-2 w-full gap-1.5"
                    onclick={() => {
                      addDocumentsModal = true;
                      folderMode = false;
                      tags = [];
                      files = [];
                      tmpDocumentTitle = "";
                      tmpDocumentDescription = "";
                      attributeNotUnique = false;
                      dropdownUploadEmpty = false;
                      nameIsAvailable = true;
                    }}
                  >
                    <PenToolPlusIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                    {$_("addNewHandwriting")}
                  </button>
                </ul>
              </div>
            </div>
          {/if}
        {:else if tabDatas.selected === "Authorization"}
          <div class="grid grid-cols-12 mt-4 gap-4">
            <Typography
              preset="Text md/Semibold"
              color="secondary"
              text={$_("authType")}
              propClass="col-span-3"
            />
            <div
              onfocusout={handleDropdownFocusLoss}
              class=" flex flex-col gap-2 relative col-span-9"
            >
              <Typography
                preset="Text sm/Medium"
                color="secondary"
                text={$_("authorizationType")}
              />
              <Button
                on:click={() => (dropdownAuthorizationType = !dropdownAuthorizationType)}
                class="flex justify-between bg-primary-light dark:bg-primary-dark border border-primary-light dark:border-primary-dark hover:bg-primary-light dark:hover:bg-primary-dark hover:border-primary-light dark:hover:border-primary-dark focus:ring-0"
              >
                <Typography preset="Text md/Reguler" text={authorizationTypeList.selected} />
                <svelte:component
                  this={getIcons("chevronDown")}
                  svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark"
                />
              </Button>
              <Typography
                preset="Text sm/Reguler"
                color="tertiary"
                text={$_("theAuthorizationHeaderWillBeAutomaticallyGeneratedWhenYouSendTheRequest")}
              />
              <div
                class={`z-10 ${dropdownAuthorizationType ? "absolute top-[5rem] left-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
              >
                <ul class="py-2 grow text-sm">
                  {#each authorizationTypeList.data as auth}
                    <button
                      class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                      onclick={() => {
                        dropdownAuthorizationType = !dropdownAuthorizationType;
                        authorizationTypeList.selected = auth.name;
                        authBasicAuth.username = "";
                        authBasicAuth.password = "";
                        authBearerToken.token = "";
                      }}
                    >
                      {$_(auth.name)}
                    </button>
                  {/each}
                </ul>
              </div>
            </div>

            {#if authorizationTypeList.selected === "Basic Auth"}
              <Divider className="col-span-12 my-4" />
              <Typography
                preset="Text md/Semibold"
                color="secondary"
                text={$_("token")}
                propClass="col-span-3"
              />
              <div class="col-span-9 flex flex-col gap-2">
                <Typography preset="Text sm/Medium" color="secondary" text={$_("username")} />
                <InputRw bind:value={authBasicAuth.username} />
              </div>
              <div class="col-span-3"></div>
              <div class="col-span-9 flex flex-col gap-2">
                <Typography preset="Text sm/Medium" color="secondary" text={$_("password")} />
                <InputRw bind:value={authBasicAuth.password} type="password" />
              </div>
            {:else if authorizationTypeList.selected === "Bearer Token"}
              <Divider className="col-span-12 my-4" />
              <Typography
                preset="Text md/Semibold"
                color="secondary"
                text={$_("token")}
                propClass="col-span-3"
              />
              <div class="col-span-9 flex flex-col gap-2">
                <Typography preset="Text sm/Medium" color="secondary" text={$_("token")} />
                <InputRw bind:value={authBearerToken.token} />
              </div>
            {/if}
            <div class="col-span-12 flex justify-end">
              <Button
                class="min-w-max gap-1.5"
                disabled={authorizationTypeList.selected === "No Auth" ||
                authorizationTypeList.selected === "Inherit Auth from Parent"
                  ? authorizationTypeList.selected === tmpSelectedAuthorizationType
                  : authorizationTypeList.selected === "Basic Auth"
                    ? authBasicAuth.username === "" && authBasicAuth.password === ""
                    : authorizationTypeList.selected === "Bearer Token"
                      ? authBearerToken.token === ""
                      : false}
                on:click={async () => {
                  console.log("clicked");
                }}
              >
                <svelte:component
                  this={getIcons("save01")}
                  svgClass="w-5 h-5 text-fg-primary-light dark:text-fg-primary-dark"
                />
                <div class="hidden md:block">{$_("saveChanges")}</div>
              </Button>
            </div>
          </div>
        {:else if tabDatas.selected === "Information"}
          <div class="grid grid-cols-12 mt-4 gap-4">
            <Typography
              preset="Text md/Semibold"
              color="secondary"
              text={$_("generalInformation")}
              propClass="col-span-3"
            />
            <div class="flex flex-col col-span-9 gap-2">
              <Typography preset="Text sm/Medium" color="secondary" text={$_("collectionTitle")} />
              <InputRw value={collectionDatas.title} />
            </div>
            <div class="col-span-3"></div>
            <div class="flex flex-col col-span-9 gap-2">
              <Typography preset="Text sm/Medium" color="secondary" text={$_("description")} />
              <InputRw value={collectionDatas.description} />
            </div>
          </div>
        {:else}
          Hello
        {/if}
      </div>
    </div>
  </div>
  <ToastContainer />
</div>

<ModalFolder
  bind:addDocumentsModal={addFolderModal}
  bind:nameIsAvailable
  bind:files
  bind:buttonDisabled
  bind:groupedFiles
  {selectedTypeCollection}
  {alltags}
  {tags}
  {isProcess}
  {handleFileUploadWrapper}
  {fileNotSupport}
  {hasDocOrDocx}
  {tmpDocumentTitle}
  {tmpDocumentDescription}
  {handleInput}
  {addDocumentByCollection}
  isFolder={true}
  actionDelete={() => (files = [])}
/>

<ModalObject
  bind:modalObjectOpen={addDocumentsModal}
  bind:nameIsAvailable
  bind:tmpDocumentTitle
  bind:groupedFiles
  bind:isProcess
  bind:files
  bind:tags
  {handleFileUploadWrapper}
  {alltags}
  {addDocumentByCollection}
  {handleInput}
  {fileNotSupport}
  {hasDocOrDocx}
  {folderMode}
  maxDepth={1}
  {createCollection}
  {collectionId}
  {checkLibraryName}
  {sendDocument}
  usedFor="Handwriting"
/>

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
        on:input={(event) => handleInput(event, collectionId, "collection")}
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

<Modal title={$_("editSubCollectionName")} bind:open={editModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("subCollectionName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={(event) => handleInput(event, collectionId, "collection")}
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
        <Typography preset="Text md/Semibold" text="Audio Information" color="secondary" />
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
        on:input={(event) => handleInput(event, libraryId, "collection")}
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
