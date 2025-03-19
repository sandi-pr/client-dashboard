<script>
  import { _ } from "svelte-i18n";
  import { Button, Modal, Input, P, Helper } from "flowbite-svelte";

  import { call } from "$lib";
  import {
    Typography,
    ToastContainer,
    StateCard,
    triggerToast,
    InputRw,
    TextareaRw,
    LibraryCard,
    check,
    isBakrie,
    ModalFolder,
    ModalObject,
    triggerLoading,
  } from "$lib";
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";

  import { format } from "date-fns";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  import FolderPlusIcon from "$lib/icon/folder-plus.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import FilePlus02Icon from "$lib/icon/file-plus-02.svelte";
  import TabBar from "$lib/components/TabBar.svelte";

  let addCollectionModal = $state(false);
  let editModal = $state(false);
  let editLibraryModal = $state(false);
  let deleteModal = $state(false);
  let deleteLibraryModal = $state(false);
  let tmpCollectionName = $state("");
  let tmpOldName = $state("");
  let tmpDescription = $state("");
  let tmpLibraryName = $state("");
  let tmpDescLibrary = $state("");
  let tmpCollectionId = $state("");
  let libraryId = $state("");
  let botIdStrapi = $state(null);
  let dropdownOpen = $state(false);
  let nameIsAvailable = $state(true);
  let gettingDatas = $state(false);
  
  //Ketbutuhan Untuk Modal DropDown
  let tmpDocumentType = $state("Handwriting");
  let selectedTypeCollection = $state("Handwriting (JPG, JPEG, PNG)");
  let acceptedFileTypes = $state(".jpg, .jpeg, .png");
  let files = $state([]);
  let dataSend = $state({ detected_text: "" });
  let groupedFiles = $state({});
  let dropdownUpload = $state(false);
  let addFolderModal = $state(false);
  let previewModal = $state(false);
  let alltags = $state([]);
  let tags = $state([]);
  let fileNotSupport = $state(false);
  let hasDocOrDocx = $state(false);
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let buttonDisabled = $state(true);
  let isProcess = $state(false);
  let folderMode = $state(false);

  let pagination = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });

  let collectionDatas = $state({
    title: "",
    description: "",
    createdAt: "",
  });

  let tabDatas = $state({
    selected: "My Collections",
    data: [
      {
        name: "My Collections",
      },
      {
        name: "Information",
      },
    ],
  });
  let collectionListDatas = $state([]);
  let informationDatas = {
    id: "1",
    collectionName: "Bakrie’s Handwriting",
    totalDocuments: 2,
    type: "",
    status: "ready",
  };

  async function getLibrary() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "get-library", {
        URLQuery: {},
        URLParams: {
          id: libraryId,
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      console.log("library", response);

      collectionDatas = {
        title: response.result.attributes.name,
        description: response.result.attributes?.meta?.description || "",
        createdAt: format(new Date(response.result.attributes.createdAt), "dd-MM-yyyy"),
      };
    }
  }

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      try {
        getLibrary();

        let response = await $call.call("client-dashboard", "get-library-collections", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: pagination.page,
            page_size: pagination.pageSize,
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

        // id: "4",
        // collectionName: "Study Scriptures of Qur’an Tafsir",
        // totalDocuments: 2,
        // type: "",
        // status: "reading",

        let link = await $call.getPage("client-dashboard", "handwriting-documents");

        let promise = response.result.data.map(async (collection) => {
          const params = new URLSearchParams(window.location.search);

          params.delete("collectionId");

          const queryParams = {};
          for (const [key, value] of params.entries()) {
            queryParams[key] = value;
          }
          const queryString = new URLSearchParams(queryParams).toString();
          return {
            id: collection.id,
            collectionName: collection.attributes.name,
            description: collection.attributes?.meta?.description || "",
            type: collection.attributes.type,
            href: link + "?collectionId=" + collection.id + "&" + queryString,
          };
        });

        collectionListDatas = await Promise.all(promise);
        gettingDatas = false;
      } catch (error) {
        console.error("Error during collection get:", error);
        gettingDatas = false;
        triggerToast($_("failedGetCollection"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
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
            parent_id: parseInt(libraryId),
            // bakrie
            bot_id: localStorage.getItem("botIdStrapi").toString(),
          };
        } else {
          dataSend = {
            collection_name: tmpCollectionName,
            description: tmpDescription,
            parent_id: parseInt(libraryId),
            bot_id: botIdStrapi,
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
        addCollectionModal = false;

        triggerToast($_("successCreateCollection"), true);
        getCollections();
        // goto(await $call.getPage("client-dashboard", "document-library"));
      } catch (error) {
        console.error("Error during collection creation:", error);
        triggerToast($_("failedCreateCollection"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  async function editLibraryName() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "update-library", {
          URLQuery: {},
          URLParams: {
            id: libraryId.toString(),
          },
          Params: {
            library_name: tmpLibraryName,
            description: tmpDescLibrary,
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

        tmpLibraryName = "";
        tmpDescLibrary = "";

        editLibraryModal = false;
        triggerToast($_("successUpdateLibrary"), true);
        getCollections();
      } catch (error) {
        console.error("Error during library editing:", error);
        triggerToast($_("failedEditLibrary"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
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

  async function deleteLibrary() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        // for (const document of documents) {
        //   if (document.attributes.collection_id === tmpCollectionId) {
        //     let response = await $call.call("client-dashboard", "delete-document", {
        //       URLQuery: {},
        //       URLParams: {
        //         id: document.id.toString(),
        //       },
        //       Params: {},
        //       Headers: {
        //         Authorization: "Bearer " + token,
        //       },
        //     });
        //     if (response.status !== 200) {
        //       console.log(response);
        //       triggerToast("error");
        //       return;
        //     }
        //   }
        // }
        let response = await $call.call("client-dashboard", "delete-library", {
          URLQuery: {},
          URLParams: {
            id: libraryId.toString(),
          },
          Params: {},
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

        console.log("ayam", response);

        triggerToast($_("successDeleteLibrary"), true);

        getCollections();
        deleteLibraryModal = false;
      } catch (error) {
        console.error("Error during library deletion:", error);
        triggerToast($_("failedDeleteLibrary"), true, "error");
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

        console.log("ayam", response);

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

  const handleInput = async (event, parentId, type) => {
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
  }

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
    dropdownUpload = false;
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
        const fileName = parts.pop();
        const base64Data = await readFileAsBase64(file);

        let folderKey = parts.slice(0, 2).join("/");
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
    }

    return structure

  }

  async function sendFolder () {
    triggerLoading(true);
    var currentFolder = []

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

      let newName = await checkNewName(folder, libraryId, "collection");

      dataSend = {
        collection_name: newName,
        description: '',
        parent_id: parseInt(libraryId),
      };

      if (sub_folder && parent_id) {
        dataSend["collection_name"] = sub_folder
        dataSend["parent_id"] = parseInt(parent_id)
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
          return
        }

        if (response3?.result?.text !== null) {

          let newName = await checkNewName(files[doc].name.split(".").slice(0, -1).join("."), res.result.data.id, "document");

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
            Authorization: "Bearer " + token
          },
        });

        let statusResDoc = triggerToastResponse(resdoc)
        if (!statusResDoc) {
          // triggerLoading(false);
          // addFolderModal = false;
          return
        }
      }
    }

    groupedFiles = null;
    addFolderModal = false;
    getCollections();
    triggerLoading(false);

    triggerToast($_("succesfullyUploaded"), true);

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
          const collectionId = localStorage.getItem("collectionIdvoice");

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


  //#FUNCTION UNTUK TRIGGERTOAST DARI RESPONSE
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

            let statusDetectedDoc = triggerToastResponse(response3);
            if (!statusDetectedDoc) {
              triggerLoading(false);
              return
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
                parent_id: parseInt(libraryId),
                documents: JSON.stringify([files[i]]),
                document_type: tmpDocumentType,
                category: tags.join(","),
                detected_text: response3.result.text,
                result_trocr: JSON.stringify(response3.result),
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
              parent_id: parseInt(libraryId),
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
        triggerLoading(false);

        if (
          tmpDocumentType !== "Handwriting" &&
          tmpDocumentType !== "Document Images" &&
          tmpDocumentType !== "Document Audio"
        ) { 
          dataSend = null;
          previewModal = false;
          addDocumentsModal = false;  
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

  onMount(() => {
    if (check("documentAi.documentLibrary:view") === false) {
      window.history.back();
      return;
    }

    botIdStrapi = localStorage.getItem("botIdStrapi" + "handwriting");

    const unsubscribe = page.subscribe(async ($page) => {
      libraryId = $page.url.searchParams.get("libraryId");
      if (libraryId === null) {
        // goto(await $call.getPage("client-dashboard", "document-library"));
        return;
      }
      getCollections();
    });
    return () => unsubscribe();
  });

  $effect(() => {
    if (groupedFiles != null) {
      buttonDisabled = false;
    } else {
      buttonDisabled = true;
    }
  });

  // $effect(() => {
  //   const unsubscribe = page.subscribe(async ($page) => {
  //     libraryId = $page.url.searchParams.get("libraryId");
  //     if (libraryId === null) {
  //       goto(await $call.getPage("client-dashboard", "document-library"));
  //       return;
  //     }
  //     getCollections();
  //   });
  //   onDestroy(() => unsubscribe());
  // });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-start justify-between mb-3 md:mb-8">
          <div class="flex flex-col gap-1.5">
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
              >
                <Copy01Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
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
            <Button
              color="primary"
              class="w-fit gap-1.5 ms-14"
              disabled={!check("documentAi.documentLibrary.collection:create")}
              on:click={() => (dropdownUpload = !dropdownUpload)}
            >
              <FilePlus02Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("addNew")}</div>
            </Button>

            <div
              class={`z-50 ${dropdownUpload ? "absolute left-0 top-12" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}
            >
              <ul class="py-2 text-sm text-gray-700 dark:text-gray-200">
                <button
                  class="flex p-2 w-full gap-1.5"
                  disabled={!check("documentAi.documentLibrary.document:create")}
                  onclick={() => {
                    addCollectionModal = true;
                    nameIsAvailable = true;
                    tmpCollectionName = "";
                    tmpDescription = "";
                  }}
                >
                  <FolderPlusIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                  {$_("addCollection")}
                </button>

                <button
                  class="flex p-2 w-full gap-1.5"
                  disabled={!check("documentAi.documentLibrary.document:create")}
                  onclick={() => {
                    folderMode = true;
                    dropdownUpload = false;
                    addFolderModal = true;
                    groupedFiles = null;
                  }}
                  ><FolderPlusIcon
                    svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark"
                  />{$_("addFolder")}
                </button>
              </ul>
            </div>
          </div>
        </div>

        <TabBar {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

        {#if tabDatas.selected === "My Collections" && collectionListDatas.length > 0}
          <div class="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mt-4">
            {#each collectionListDatas as collection}
              <LibraryCard
                collectionName={collection.collectionName}
                usedFor="libraries_collections"
                functButton={() => {
                  // selectCollection({
                  //   id: collection.id,
                  //   name: collection.collectionName,
                  //   description: collection.description,
                  //   totalDocuments: collection.totalDocuments,
                  // });
                }}
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
                {libraryId}
              />
            {/each}
          </div>
        {:else if tabDatas.selected === "Information" && Object.keys(informationDatas).length > 0}
          <div class="grid grid-cols-12 mt-4 gap-2">
            <div class="col-span-4">
              <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                {$_("generalInformation")}
              </P>
            </div>
            <div class="col-span-8">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >{$_("collectionTitle")}</P
              >
            </div>
            <div class="col-span-4"></div>
            <div class="col-span-8">
              <InputRw onEdit={false} bind:value={collectionDatas.title} />
            </div>
            <div class="col-span-4"></div>
            <div class="col-span-8">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >{$_("description")}</P
              >
            </div>
            <div class="col-span-4"></div>
            <div class="col-span-8">
              <TextareaRw onEdit={false} bind:value={collectionDatas.description} />
            </div>
          </div>
        {:else if gettingDatas}
          <StateCard
            illustration={IllustrationProgress}
            description={$_("fetchingDataPleaseWait")}
          />
        {:else}
          <StateCard
            illustration={IllustrationEmpty}
            description={$_("youDontHaveAnyDataYet")}
          />
        {/if}
      </div>
    </div>
  </div>
  <ToastContainer />
</div>

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

<Modal title={$_("editLibraryName")} bind:open={editLibraryModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("libraryName")}*`} />
      <Input
        bind:value={tmpLibraryName}
        on:input={(event) => handleInput(event, 0, 'collection')}
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
      <Input bind:value={tmpDescLibrary} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={editLibraryName}
      disabled={!tmpLibraryName || !nameIsAvailable}
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
    <Button onclick={deleteDocumentCollection} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>

<Modal bind:open={deleteLibraryModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
      {$_("library")}
      {collectionDatas.title}?
    </h3>
    <Button onclick={deleteLibrary} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteLibraryModal = false)}
      >{$_("noCancel")}</Button
    >
  </div>
</Modal>

<!-- <ModalFolder
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

  isFolder = {true}
  actionDelete={() => (files = [])}
/> -->

<ModalObject
  bind:modalObjectOpen={addFolderModal}
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
  maxDepth={2}
  {createCollection}
  collectionId={libraryId}
  {checkLibraryName}
  usedFor="Handwriting"
/>