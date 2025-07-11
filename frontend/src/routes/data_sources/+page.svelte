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
    ActionCard,
    check,
    ModalObject,
    triggerLoading,
    ModalFilePreview,
  } from "$lib";
  import DataSourcesCard from "$lib/ui/cards/DataSourcesCard.svelte";
  import { _ } from "svelte-i18n";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { goto } from "$app/navigation";
  import { showAI, showNotif, call, triggerToast } from "$lib";

  const isFrontend = process.env.NODE_ENV === "development";

  let libraryId, collectionId;
  let collections = $state([]);
  let latestDocumentCollection = $state([]);
  let files = $state([]);
  let actionItems = $state([]);
  let tmpOldName = $state("");
  let tmpCollectionName = $state("");
  let tmpDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let tmpDocumentType = $state("Data");
  let selectedTypeCollection = $state("Data (CSV)");
  let acceptedFileTypes = $state(".csv, .db");
  let addDocumentsModal = $state(false);
  let addLibraryModal = $state(false);
  let fileNotSupport = $state(false);
  let editModal = $state(false);
  let deleteModal = $state(false);
  let timer;
  let nameIsAvailable = $state(true);
  let isProcess = $state(false);
  let alltags = $state([]);
  let tags = $state([]);
  let folderMode = $state(false);
  let audioUrl = $state(null);
  let previewModal = $state(false);
  let dataSend = $state({ detected_text: "" });
  let selectedFile = $state(null);
  let groupedFiles = $state(null);

  async function getActionItems() {
    actionItems = [
      {
        onclick: () => {
          addDocumentsModal = true;
          tmpCollectionName = "";
          tmpDescription = "";
        },
        // href: await $call.getPage("client-dashboard", "document-collection-create"),
        icon: "filePlus",
        actionName: "Add data source",
        desc: "desc Add data source",
      },
    ];
  }

  let gettingDatas = $state(false);

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      collections = [];
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

      let link = isFrontend
        ? "/data_sources/collections/documents"
        : await $call.getPage("client-dashboard", "data-sources-documents");

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

        return {
          id: collection.id,
          collectionName: collection.attributes.name,
          description: collection.attributes?.meta?.description || "",
          type: collection.attributes.type,
          href: link + "?collectionId=" + collection.id + "&" + `libraryId=${libraryId}`,
        };
      });

      latestDocumentCollection = await Promise.all(promise);

      console.log("latestDocumentCollection", latestDocumentCollection);

      gettingDatas = false;
    }
  }

  const handleInput = async (event, parentId = collectionId, type = "document") => {
    const name = event.target.value;
    const status = await checkLibraryName(name.trim(), parentId, type);

    console.log("status", status);

    if (tmpOldName === name) {
      nameIsAvailable = true;
    } else if (status === false) {
      nameIsAvailable = false;
    } else if (status === true) {
      nameIsAvailable = true;
    }
  };

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

      addLibraryModal = false;
      getCollections();
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
        type: tmpDocumentType,
      };

      structure[folderKey][fileName] = fileObject;
    }

    console.log("structure::::::::::::");
    console.log(structure);
    return structure;
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

  async function sendFolder() {
    triggerLoading(true);
    var currentFolder = [];

    // Buat Collection by Name Folder Dan Document
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

      dataSend = {
        collection_name: folder,
        description: "",
        parent_id: parseInt(collectionId),
      };

      if (sub_folder && parent_id) {
        dataSend["collection_name"] = sub_folder;
        dataSend["parent_id"] = parseInt(parent_id);
      }
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
          Authorization: "Bearer " + token,
        },
      });

      let statusResDoc = triggerToastResponse(resdoc);
      if (!statusResDoc) {
        triggerLoading(false);
        return;
      }
    }

    addFolderModal = false;
    getCollections();
    triggerLoading(false);
  }

  async function addDocumentByCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      if (groupedFiles != null) {
        sendFolder();
        return;
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

  async function viewAll() {
    let link = isFrontend
      ? "/data_sources/collections"
      : await $call.getPage("client-dashboard", "data-sources-collections");
    goto(link + `?libraryId=${libraryId}`);
  }

  onMount(async () => {
    getActionItems();

    // let idDatasource = localStorage.getItem("libraryIdDatasource");
    // await setDatasourceLibraryIds(idDatasource);

    libraryId = localStorage.getItem("libraryId" + "insights");
    collectionId = localStorage.getItem("collectionId" + "insights");

    getCollections();
    getLabels();
    // latestDataSources = JSON.parse(localStorage.getItem("latestDataSources")) || [];
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
  <div class="flex flex-col w-full pb-20 m-4 gap-3 md:gap-8">
    <!-- User Information -->
    <Typography preset="Display sm/Semibold" text={$_("dataSources")} propClass="fluid-xl" />

    {#if check("documentAi.documentLibrary:create")}
      <!-- Action Items -->
      <div class="flex flex-col gap-4">
        <SectionHeader title={$_("action")} />
        <ActionCard {actionItems} />
      </div>
    {/if}

    <div class="flex flex-col gap-4">
      <div
        class="flex items-center justify-between border-b border-secondary-light dark:border-secondary-dark pb-4"
      >
        <Typography preset="Text lg/Semibold" color="primary" text={$_("dataSourcesList")} />
        <Button on:click={() => viewAll()} color="light">{$_("seeAll")}</Button>
      </div>

      <div>
        {#if latestDocumentCollection.length > 0}
          <div
            class={`${$showAI || $showNotif ? "grid grid-cols-1 xl:grid-cols-3" : "grid grid-cols-2 lg:grid-cols-4"} gap-4`}
          >
            {#each latestDocumentCollection as document, index}
              <DataSourcesCard
                collectionName={document.collectionName}
                totalDocuments={document.totalDocuments}
                status={document.status}
                linkIsFrontend={"/data_sources/collections/documents?collectionName=" +
                  document.collectionName}
                functEdit={() => {
                  tmpCollectionId = document.id;
                  tmpOldName = document.collectionName;
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
  </div>
</div>

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
  usedFor={"Data"}
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

<Modal title={$_("createNewDataSources")} bind:open={addLibraryModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("dataSourcesName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={(event) => handleInput(event, libraryId, "collection")}
        placeholder={$_("whatIsYourDataSourcesTitle")}
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
      class="w-full">{$_("addNewDataSources")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editDataSourcesName")} bind:open={editModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("dataSourcesName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={(event) => handleInput(event, libraryId, "collection")}
        placeholder={$_("whatIsYourDataSourcesTitle")}
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
</Modal>

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
