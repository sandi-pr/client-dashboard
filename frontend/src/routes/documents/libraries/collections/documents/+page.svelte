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
    LibraryCard,
    isBakrie,
    ModalObject,
    AddDocModal,
    ModalFilePreview,
  } from "$lib";
  import DataTableLibraryCollection from "$lib/components/DataTableLibraryCollection.svelte";
  import DataTableLibraryCollectionVertical from "$lib/components/DataTableLibraryCollectionVertical.svelte";
  import {
    SearchOutline,
    ExclamationCircleSolid,
    ArrowUpOutline,
  } from "flowbite-svelte-icons";
  import Trash01Icon from "$lib/icon/trash-01.svelte";
  import MoveIcon from "$lib/icon/move.svelte";
  import FilePlus02Icon from "$lib/icon/file-plus-02.svelte";
  import FolderPlusIcon from "$lib/icon/folder-plus.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import ModalFolder from "$lib/ui/modal/ModalFolder.svelte";
  
  let collectionId;
  let libraryId = $state(null);
  let inputRef;
  let checkAllState = $state(false);
  let addDocumentsModal = $state(false);
  let addFolderModal = $state(false);
  let addCollectionModal = $state(false);
  let addCategoriesModal = $state(false);
  let modalObjectOpen = $state(false);
  let popupModal = $state(false);
  let updateModal = $state(false);
  let editModal = $state(false);
  let deleteModal = $state(false);
  let deleteCategoriesModal = $state(false);
  let dropdownOpen = $state(false);
  let dropdownMore = $state(false);
  let editCollectionModal = $state(false);
  let dropdownUpload = $state(false);
  let dropdownUploadEmpty = $state(false);
  let fileNotSupport = $state(false);
  let buttonDisabled = $state(true);
  let gettingDatas = $state(false);
  let gettingDocs = $state(false);
  let buttonDisabledCheckAction = $state(false);
  let pagination = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });
  let paginationCollection = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpCollectionName = $state("");
  let tmpOldName = $state("");
  let tmpDescription = $state("");
  let tmpDocumentType = $state("Documents");
  let selectedTypeCollection = $state("Documents (TXT, PDF, DOC, DOCX)");
  let acceptedFileTypes = $state(".txt, .pdf, .doc, .docx");
  let search = $state("");
  let alltags = $state([]);
  let tags = $state([]);
  let newTag = $state("");
  let files = $state([]);
  let groupedFiles = $state(null);
  let attributeNotUnique = $state(false);
  let checkboxDelete = $state(false);
  let checkboxUpdate = $state(false);
  let tmpId = $state(null);
  let isMobileView = $state(false);
  let hasDocOrDocx = $state(false);

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
  let div = $state();
  let nameIsAvailable = $state(true);
  let isProcess = $state(false);
  let selected = $state([]);
  let dummyTags = $state(["general", "news", "article"]);
  let folderMode = $state(false);

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

        let link = await $call.getPage("client-dashboard", "document-library-detail");

        let promise = response.result.data.map(async (collection) => {
          const params = new URLSearchParams(window.location.search);

          params.delete("documentId");

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

        console.log(response);

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
    console.log("next page");
  };

  const prevPage = (evt) => {
    if (pagination.page > 1) pagination.page = pagination.page - 1;
    getDocuments();
    console.log("prev page");
  };

  const handleSearch = (evt) => {
    console.log("clicked");
    getDocuments();
  };

  let audioUrl = $state(null);

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

        if (
          tmpDocumentType !== "Handwriting" &&
          tmpDocumentType !== "Document Images" &&
          tmpDocumentType !== "Document Audio"
        ) {
          dataSend = null;
          previewModal = false;
          addDocumentsModal = false;
          addFolderModal = false;
          getDocuments();
          getCollections(true);
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

  //#FUNCTION UNTUK TRIGGERTOAST DARI RESPONSE
  //TODO: BISA DIJADIKAN UTIL
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

  async function groupFilesByPath(files) {
    const structure = {};
    console.log("masuk ke func group files")

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
    getCollections(true);
    triggerLoading(false);
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
    dropdownUploadEmpty = false;
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
        goto(await $call.getPage("client-dashboard", "document-library"));
        return;
      }
      getCollections(true, false);
      getCollections(false, false);
      getDocuments();
      getLabels();
    });

    checkMobileView();
    window.addEventListener("resize", checkMobileView);

    return () => {
      unsubscribe();
      window.removeEventListener("resize", checkMobileView);
    };
  });

  $effect(() => {
    checkAllState = collectionListDatas.every((data) => data.checked === true);
    if (files.length > 0) {
      hasDocOrDocx = files.some((file) => file.ext === ".doc" || file.ext === ".docx");
      buttonDisabledCheckAction = !collectionListDatas.some((data) => data.checked === true);
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
              color="primary"
              outline
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
                  disabled={!check("documentAi.documentLibrary.document:create")}
                  onclick={() => {
                    addFolderModal = true;
                    groupedFiles = null;
                  }}
                  ><FolderPlusIcon
                    svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark"
                  />{$_("uploadCollection")}</button
                >
                <button
                  class="flex p-2 w-full gap-1.5"
                  disabled={!check("documentAi.documentLibrary.document:create")}
                  onclick={() => {
                    addFolderModal = true;
                    groupedFiles = null;
                  }}
                  ><FolderPlusIcon
                    svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark"
                  />{$_("uploadCollection")}</button
                >
                <button
                  class="flex p-2 w-full gap-1.5"
                  disabled={!check("documentAi.documentLibrary.document:create")}
                  onclick={() => {
                    addDocumentsModal = true;
                    folderMode = true;
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
                    dropdownUpload = false;
                    nameIsAvailable = true;
                  }}
                  ><FilePlus02Icon
                    svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark"
                  />{$_("addDocument")}</button
                >
              </ul>
            </div>
            {#if !checkTitle(collectionDatas.title)}
              {#if check("documentAi.documentLibrary.collection:edit") || check("documentAi.documentLibrary.collection:delete")}
                <div onfocusout={handleDropdownFocusLoss} class=" flex items-center relative">
                  <Button
                    on:click={() => (dropdownMore = !dropdownMore)}
                    color="light"
                    class="p-2.5"
                  >
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
                          }}
                        >
                          {$_("delete")}
                        </button>
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
            {#each childCollectionList as collection}
              <LibraryCard
                collectionName={collection.collectionName}
                usedFor="libraries_collections"
                functEdit={() => {
                  tmpCollectionId = collection.id;
                  tmpCollectionName = collection.collectionName;
                  tmpOldName = collection.collectionName;
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
                {libraryId}
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
              <FilePlus02Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("addNew")}</div>
            </Button>
            <div
              class={`z-50 ${dropdownUploadEmpty ? "absolute top-[190%] left-1/2 transform -translate-x-1/2 -translate-y-1/2" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}
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
                  <FilePlus02Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                  {$_("addDocument")}
                </button>
              </ul>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
  <ToastContainer />
</div>

<!-- <AddDocModal
  bind:addDocumentsModal
  bind:nameIsAvailable
  {selectedTypeCollection}
  bind:files
  {alltags}
  {tags}
  {isProcess}
  {handleFileUploadWrapper}
  {fileNotSupport}
  {hasDocOrDocx}
  {tmpDocumentTitle}
  {tmpDocumentDescription}
  {handleInput}
  bind:buttonDisabled
  {addDocumentByCollection}
  actionDelete={() => (files = [])}
/> -->

<ModalFolder
  bind:addDocumentsModal={addFolderModal}
  bind:nameIsAvailable
  {selectedTypeCollection}
  bind:files
  {alltags}
  {tags}
  {isProcess}
  {handleFileUploadWrapper}
  {fileNotSupport}
  {hasDocOrDocx}
  {tmpDocumentTitle}
  {tmpDocumentDescription}
  {handleInput}
  bind:buttonDisabled
  {addDocumentByCollection}
  actionDelete={() => (files = [])}
/>

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
  {hasDocOrDocx}
  {folderMode}
  maxDepth={1}
/>

<Modal title={$_("createNewCategories")} bind:open={addCategoriesModal}>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("categories")}`} />
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

<ModalFilePreview
  bind:previewModal
  bind:dataSend
  cancelOnClick={() => {
    previewModal = false;
    addDocumentsModal = true;
  }}
  okOnClick={sendDocument(true)}
  {selectedFile}
  {audioUrl}
/>

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
