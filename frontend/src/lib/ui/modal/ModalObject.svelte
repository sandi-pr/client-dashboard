<script>
  import { Modal, Button, Helper, Alert } from "flowbite-svelte";
  import {
    Typography,
    FeaturedIcon,
    getIcons,
    InputRw as Input,
    CardFeature as Card,
    UploadFile,
    UploadFolder,
    triggerLoading,
  } from "$lib";
  import { getBgColors } from "$lib";
  import MultiSelect from "svelte-multiselect";
  import TabBar from "$lib/components/TabBar.svelte";
  import { InfoCircleSolid } from "flowbite-svelte-icons";
  import { _ } from "svelte-i18n";
  import { onMount } from "svelte";

  let {
    modalObjectOpen = $bindable(false),
    handleFileUploadWrapper,
    files = $bindable([]),
    groupedFiles = $bindable(null),
    isProcess = $bindable(false),
    addDocumentByCollection,
    alltags = [],
    tags = $bindable([]),
    handleInput,
    nameIsAvailable = $bindable(true),
    tmpDocumentTitle = $bindable(""),
    usedFor = "Documents",
    fileNotSupport = false,
    hasDocOrDocx = false,
    folderMode = false,
    maxDepth = 3,
    createCollection = () => {},
    collectionId = 1,
    checkLibraryName = () => {},
    limitation = [],
    allowedFileTypesDragAndDrop = [],
  } = $props();

  let tabBarSelected = $state("Upload");
  let linkImportLink = $state("");
  let linkChecked = $state(false);
  let urlCheckerRunning = $state(false);
  let firstLayer = $state([]);
  let secondLayer = $state([]);
  let isProcessing = $state(false);
  let droppedItemType = $state("");
  let delayTimeout;

  let mainContentUpload = [
    {
      tabBar: {
        tabDatas: {
          selected: "Upload",
          data: [
            {
              name: "Upload",
            },
            {
              name: "Link",
            },
          ],
        },
      },
    },
    {
      text: {
        text: $_("documentType"),
        color: "secondary",
        preset: "Text sm/Medium",
      },
    },
    {
      input: {
        type: "text",
        onEdit: false,
        value: "Documents (TXT, PDF, DOC, DOCX)",
      },
    },
    {
      uploadFile: {
        title: $_("uploadDocuments"),
      },
    },
    {
      text: {
        text: "Label",
        color: "primary",
        preset: "Text sm/Medium",
      },
    },
    {
      multiSelect: {
        // selected: tags,
        // options: alltags,
        allowUserOptions: "append",
        autoScroll: true,
        noMatchingOptionsMsg: "",
        createOptionMsg: null,
        ulSelectedClass: "text-primary-light dark:text-primary-dark",
        liSelectedClass: "border border-primary-light dark:border-primary-dark",
        ulOptionsClass: "bg-secondary-light dark:bg-tertiary-dark",
        outerDivClass:
          "bg-primary-light dark:bg-secondary-dark border rounded-lg border-gray-300 dark:border-gray-600 focus:none focus:ring-0 focus:border-gray-300 dark:focus:border-gray-600 py-2",
      },
    },
    {
      metaFile: {
        content: [
          {
            text: {
              text: `${$_("documentTitle")}*`,
              color: "secondary",
              preset: "Text sm/Medium",
            },
          },
          {
            input: {
              id: "docTitle",
              type: "text",
              value: "",
              onEdit: true,
              oninput: handleInput,
            },
          },
          // {
          //   text: {
          //     text: $_("description"),
          //     color: "secondary",
          //     preset: "Text sm/Medium",
          //   },
          // },
          // {
          //   input: {
          //     id: "docDescription",
          //     type: "text",
          //     value: "",
          //     onEdit: true,
          //   },
          // },
        ],
      },
    },
  ];

  let mainContentLink = [
    {
      tabBar: {
        tabDatas: {
          selected: "Link",
          data: [
            {
              name: "Upload",
            },
            {
              name: "Link",
            },
          ],
        },
      },
    },
    {
      text: {
        text: $_("documentType"),
        color: "primary",
        preset: "Text md/Semibold",
      },
    },
    {
      input: {
        type: "text",
        onEdit: false,
        value: "Documents (TXT, PDF, DOC, DOCX)",
      },
    },
    {
      text: {
        text: $_("importLink"),
        color: "primary",
        preset: "Text sm/Medium",
      },
    },
    {
      input: {
        id: "link",
        type: "text",
        onEdit: true,
        placeholder: "Paste image link (e.g. https://plus.unsplash.com/...)",
        oninput: (e) => {
          linkImportLink = e.target.value;
          configModal.footer[0].button.buttonDisabled = e.target.value ? false : true;
        },
      },
    },
    {
      card: {
        cardData: {
          // iconData: { icon: "filePlus" },
          contentCenter: [
            [
              {
                icon: {
                  icon: "infoCircle",
                  color: "error",
                },
                text: {
                  text: "",
                  color: "fg-error-primary",
                  preset: "Text sm/Regular",
                },
                justify: "justify-start",
              },
            ],
          ],
        },
        propClass: getBgColors("error-primary") + " border-0 rounded-xl",
      },
    },
    {
      text: {
        text: $_("labels"),
        color: "primary",
        preset: "Text sm/Medium",
      },
    },
    {
      multiSelect: {
        // selected: tags,
        // options: alltags,
        allowUserOptions: "append",
        autoScroll: true,
        noMatchingOptionsMsg: "",
        createOptionMsg: null,
        ulSelectedClass: "text-primary-light dark:text-primary-dark",
        liSelectedClass: "border border-primary-light dark:border-primary-dark",
        ulOptionsClass: "bg-secondary-light dark:bg-tertiary-dark",
        outerDivClass:
          "bg-primary-light dark:bg-secondary-dark border rounded-lg border-gray-300 dark:border-gray-600 focus:none focus:ring-0 focus:border-gray-300 dark:focus:border-gray-600 py-2",
      },
    },
    {
      metaFile: {
        content: [
          {
            text: {
              text: `${$_("documentTitle")}*`,
              color: "secondary",
              preset: "Text sm/Medium",
            },
          },
          {
            input: {
              id: "docTitle",
              type: "text",
              value: "file",
              onEdit: true,
              oninput: handleInput,
            },
          },
          // {
          //   text: {
          //     text: $_("description"),
          //     color: "secondary",
          //     preset: "Text sm/Medium",
          //   },
          // },
          // {
          //   input: {
          //     id: "docDescription",
          //     type: "text",
          //     value: "",
          //     onEdit: true,
          //   },
          // },
        ],
      },
    },
  ];

  let footerUpload = [
    {
      button: {
        text: "addNewDocument",
        onClick: addDocumentByCollection,
        buttonDisabled: true,
      },
    },
  ];

  const createCollectionByLinkGD = async () => {
    const processFolder = async (folder, parentId) => {
      try {
        const createdCollectionId = await createCollection({
          collection_name: folder.name,
          description: "",
          parent_id: parseInt(parentId),
        });

        let randomString = Math.random().toString(36).substring(7);

        for (let i = 0; i < folder.files.length; i++) {
          if (folder.files[i].files) {
            await processFolder(folder.files[i], createdCollectionId);
          } else {
            folder.files[i].name =
              folder.files[i].name.split(".")[0] + "-" + randomString + folder.files[i].ext;
          }
        }

        const filesToUpload = folder.files.filter((file) => !file.files);
        if (filesToUpload.length > 0) {
          await addDocumentByCollection(filesToUpload, createdCollectionId);
        }
      } catch (error) {
        console.error(`Error processing folder ${folder.name}:`, error);
        // Handle the error appropriately, e.g., show an alert or notification
      }
    };

    try {
      for (const folder of files) {
        await processFolder(folder, collectionId);
      }
    } catch (error) {
      console.error("Error creating collections by link:", error);
      // Handle the error appropriately, e.g., show an alert or notification
    }
  };

  let footerUploadCollection = [
    {
      button: {
        text: "addNewCollection",
        onClick: createCollectionByLinkGD,
        buttonDisabled: true,
      },
    },
  ];

  let footerLink = [
    {
      button: {
        text: "checkLink",
        onClick: async () => {
          // console.log(linkImportLink);
          let urlValid = await urlChecker(linkImportLink, limitation, maxDepth);
          if (urlValid) {
            linkChecked = true;
            configModal.footer = !folderMode ? footerUpload : footerUploadCollection;
            configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text = "";
            configModal.mainContent[3].text.text = $_("uploadDocuments");
          } else {
            linkChecked = false;
            linkImportLink = "";
          }
          urlCheckerRunning = false;
        },
        buttonDisabled: true,
      },
    },
  ];

  let configModal = $state({
    header: {
      iconLeft: {
        icon: "filePlus",
      },
      contentCenter: [
        {
          text: {
            text: "Add Document",
            color: "primary",
            preset: "Text lg/Semibold",
          },
        },
        {
          text: {
            text: "Upload & Organize Your Files Effortlessly",
            color: "tertiary",
            preset: "Text sm/Regular",
          },
        },
      ],
    },
    mainContent: mainContentUpload,
    footer: footerUpload,
  });

  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
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

  const urlChecker = async (url, limitation = [], maxDepth = 3) => {
    urlCheckerRunning = true;
    let errorObject;
    try {
      let response;
      const isGoogleDrive = url.includes("drive.google.com");
      const fileId = isGoogleDrive ? url.match(/[-\w]{25,}/) : null;

      if (isGoogleDrive && !fileId) {
        throw new Error("Invalid Google Drive URL");
      }

      if (isGoogleDrive && url.includes("/drive/folders/")) {
        await getFolderContents(fileId, {}, 1, maxDepth);
        console.log("files", files);
        return true;
      } else {
        // Handle Google Drive file or other URL
        const fetchUrl = isGoogleDrive
          ? `https://www.googleapis.com/drive/v3/files/${fileId[0]}?alt=media&key=AIzaSyDi62RQsNNjA-_IfGyaDfWc7L11dVm75YI`
          : url;
        response = await fetch(fetchUrl);

        if (!response.ok) {
          errorObject = await getErrorObject(response);
          throw new Error("Network response was not ok");
        }

        // Fetch file metadata
        const fileMetadataUrl = `https://www.googleapis.com/drive/v3/files/${fileId[0]}?key=AIzaSyDi62RQsNNjA-_IfGyaDfWc7L11dVm75YI`;
        const fileMetadataResponse = await fetch(fileMetadataUrl);
        if (!fileMetadataResponse.ok) {
          throw new Error("Failed to fetch file metadata");
        }
        const fileMetadata = await fileMetadataResponse.json();

        const contentType = response.headers.get("content-type");
        if (limitation.length === 0 || limitation.some((type) => contentType.includes(type))) {
          await processFile(response, contentType, url, fileMetadata);
          return true;
        } else {
          setErrorText("fileNotSupported");
          throw new Error("URL does not contain a valid resource");
        }
      }
    } catch (error) {
      handleError(error, errorObject);
      return false;
    }
  };

  const getFolderContents = async (fileId, fileMetadata = {}, depth = 1, maxDepth = 3) => {
    if (depth > maxDepth) return;

    // Fetch folder metadata (optimized)
    let folderMetadata;
    if (!Object.keys(fileMetadata).length > 0) {
      const folderMetadataUrl = `https://www.googleapis.com/drive/v3/files/${fileId[0]}?fields=id,name,mimeType&key=AIzaSyDi62RQsNNjA-_IfGyaDfWc7L11dVm75YI`;
      const folderMetadataResponse = await fetch(folderMetadataUrl);
      if (!folderMetadataResponse.ok) {
        throw new Error("Failed to fetch folder metadata");
      }
      folderMetadata = await folderMetadataResponse.json();
    } else {
      folderMetadata = fileMetadata;
    }

    // Handle Google Drive folder
    const folderUrl = `https://www.googleapis.com/drive/v3/files?q='${fileId}'+in+parents&fields=files(id,name,mimeType)&key=AIzaSyDi62RQsNNjA-_IfGyaDfWc7L11dVm75YI`;
    const folderResponse = await fetch(folderUrl);
    if (!folderResponse.ok) {
      throw new Error("Failed to fetch folder contents");
    }
    const folderData = await folderResponse.json();

    await insertLayer(folderMetadata, folderData, depth);

    if (depth >= maxDepth) {
      await processLayer();
      return;
    }

    // console.log("folderData", folderData);

    let counterFolder = 0;
    if (folderData?.files?.length > 0) {
      folderData.files.forEach(async (file) => {
        if (file.mimeType === "application/vnd.google-apps.folder") {
          counterFolder++;
          await getFolderContents([file.id], file, depth + 1, maxDepth);
        }
      });
    }

    if (counterFolder === 0) {
      await processLayer();
    }
  };

  const processLayer = async () => {
    for (let i = 0; i <= maxDepth; i++) {
      if (i === 1) {
        firstLayer = await removeFolderObject(firstLayer);
        firstLayer = await hydrateContent(firstLayer);
      } else if (i === 2) {
        secondLayer = await removeFolderObject(secondLayer);
        secondLayer = await hydrateContent(secondLayer);
      }
    }

    files = firstLayer.map((folder) => {
      return {
        ...folder,
        files: [...folder.files, ...secondLayer],
      };
    });

    console.log("cek files processLayer", files);
  };

  const insertLayer = async (folderMetadata, folderData, depth) => {
    if (depth === 1) {
      firstLayer.push({
        folderId: folderMetadata.id,
        folderName: folderMetadata.name,
        files: folderData.files,
      });
    } else if (depth === 2) {
      secondLayer.push({
        folderId: folderMetadata.id,
        folderName: folderMetadata.name,
        files: folderData.files,
      });
    }
  };

  const removeFolderObject = (layer) => {
    for (let i = 0; i < layer.length; i++) {
      layer[i].files = layer[i].files.filter((file) => {
        return file.mimeType !== "application/vnd.google-apps.folder";
      });
    }
    return layer;
  };

  const hydrateContent = async (layer) => {
    const hydratedLayer = await Promise.all(
      layer.map(async (folder, index) => {
        const filesPromise = await Promise.all(
          folder.files.map(async (file) => {
            if (file?.mimeType) {
              const fileUrl = `https://www.googleapis.com/drive/v3/files/${file.id}?alt=media&key=AIzaSyDi62RQsNNjA-_IfGyaDfWc7L11dVm75YI`;
              const response = await fetch(fileUrl);

              if (!response.ok) {
                throw new Error("Failed to fetch file content");
              }

              const blob = await response.blob();
              const base64data = await readFileAsBase64(blob);
              const fileType = response.headers.get("content-type").split("/")[1];
              const fileName = `${file.name.slice(0, 10)}.${fileType}`;
              const fileData = {
                id: files.length < 1 ? 1 : files.length + 1,
                name: file.name || fileName,
                size: formatBytes(blob.size),
                base64: base64data,
                ext: `.${fileType}`,
                type: getFileType(response.headers.get("content-type")),
              };
              return fileData;
            }
          }),
        );

        return {
          id: folder.files[index]?.mimetype ? index + 1 : folder.files[index].id,
          name: folder.folderName,
          files: filesPromise,
        };
      }),
    );

    return hydratedLayer;
  };

  const processFile = async (response, contentType, url, fileMetadata) => {
    const blob = await response.blob();
    const base64data = await readFileAsBase64(blob);
    const fileType = contentType.split("/")[1];
    const fileName = `${url.split("/")[3].slice(0, 10)}.${fileType}`;
    const fileData = {
      id: files.length < 1 ? 1 : files.length + 1,
      name: fileMetadata.name || fileName,
      size: formatBytes(blob.size),
      base64: base64data,
      ext: `.${fileType}`,
      type: getFileType(contentType),
    };
    files.push(fileData);
  };

  const getFileType = (contentType) => {
    if (contentType.includes("image")) return "image";
    if (contentType.includes("text/plain")) return "text";
    if (contentType.includes("application/pdf")) return "pdf";
    if (
      contentType.includes("application/msword") ||
      contentType.includes(
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      )
    )
      return "document";
    return "unknown";
  };

  const getErrorObject = async (response) => {
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
      return await response.json();
    } else {
      return { error: { message: "Non-JSON response received" } };
    }
  };

  const setErrorText = (key) => {
    configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text = $_(key);
  };

  const handleError = (error, errorObject) => {
    if (errorObject && errorObject.error) {
      if (errorObject.error.code === 403) {
        setErrorText("accessRestrictedCheckSharingSettings");
      } else if (errorObject.error.code === 404) {
        setErrorText("fileNotFoundMakeSureTheLinkIsCorrectAndAccessible");
      }
    } else if (error instanceof TypeError && error.message === "Failed to fetch") {
      setErrorText("accessRestrictedCheckSharingSettings");
    } else {
      setErrorText("failedToFetchResource");
    }
    console.error("Error fetching the URL:", error);
  };

  // Drag and Drop functionality
  function handleDragOver(event) {
    event.preventDefault();
    event.stopPropagation();
  }

  function handleDrop(event) {
    event.preventDefault();
    event.stopPropagation();
    isProcessing = true;
    const items = event.dataTransfer.items;
    let hasFile = false;
    let hasFolder = false;
    let fileCount = 0;
    let folderCount = 0;

    for (let i = 0; i < items.length; i++) {
      const item = items[i].webkitGetAsEntry();
      if (item) {
        if (item.isFile) {
          hasFile = true;
          fileCount++;
        } else if (item.isDirectory) {
          hasFolder = true;
          folderCount++;
        }
      }
    }

    if (folderCount > 1 || (hasFile && hasFolder)) {
      alert("Please drop either a single file, multiple files, or a single folder.");
      return;
    }

    for (let i = 0; i < items.length; i++) {
      const item = items[i].webkitGetAsEntry();
      if (item) {
        traverseFileTree(item);
      }
    }

    // console.log(files);

    if (hasFile) {
      droppedItemType = fileCount > 1 ? "multiple files" : "file";
    } else if (hasFolder) {
      droppedItemType = "folder";
    }
  }

  function traverseFileTree(item, path = "", parentFolder = null) {
    if (item.isFile) {
      item.file(async (file) => {
        if (
          !allowedFileTypesDragAndDrop.includes(file.type) &&
          allowedFileTypesDragAndDrop.length !== 0
        ) {
          console.log(`File type ${file.type} is not allowed.`);
          return;
        }
        let base64 = await toBase64(file);
        // remove the prefix (pdf, image, anything) if it exists
        if (base64.startsWith("data:application/pdf;base64,")) {
          let prefix = [
            "data:application/pdf;base64,",
            "data:image/jpeg;base64,",
            "data:image/png;base64,",
            "data:image/jpg;base64,",
            "data:audio/mpeg;base64,",
            "data:audio/wav;base64,",
            "data:audio/ogg;base64,",
            "data:text/csv;base64,",
          ];
          for (let i = 0; i < prefix.length; i++) {
            base64 = base64.replace(prefix[i], "");
          }
        }

        let nameIsAvailable = await checkLibraryName(
          file.name.split(".")[0],
          collectionId,
          "document",
        );
        let newFileName = "";
        if (!nameIsAvailable) {
          let randomString = Math.random().toString(36).substring(7);
          newFileName =
            file.name.split(".")[0] + "-" + randomString + "." + file.name.split(".")[1];
        }
        tmpDocumentTitle = newFileName || file.name.split(".")[0];

        const fileData = {
          id: files.length < 1 ? 1 : files.length + 1,
          name: newFileName || file.name,
          documentTitle: file.name.split(".")[0],
          description: "",
          size: formatBytes(file.size),
          base64: base64,
          ext: "." + file.type.split("/")[file.type.split("/").length - 1],
          type: getFileType(file.type),
        };
        if (parentFolder) {
          parentFolder.files = [...parentFolder.files, fileData];
          // console.log("parentFolder", parentFolder);
        } else {
          files = [...files, fileData];
        }
        // console.log("File added:", fileData);
      });
    } else if (item.isDirectory) {
      const existingFolder = files.find((folder) => folder.name === item.name);
      const folderData = existingFolder || {
        id: uuidv4(),
        name: item.name,
        files: [],
        parentId: parentFolder ? parentFolder.id : null,
      };
      if (!existingFolder) {
        files = [...files, folderData];
      }
      const dirReader = item.createReader();
      dirReader.readEntries((entries) => {
        // console.log("Folder entries:", entries);
        for (let i = 0; i < entries.length; i++) {
          traverseFileTree(entries[i], path + item.name + "/", folderData);
        }
      });
    }
  }

  function toBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result);
      reader.onerror = (error) => reject(error);
    });
  }

  function processFolders(folders) {
    if (folders.length <= 1) {
      isProcessing = false;
      return;
    }

    const folderMap = new Map();
    folders.forEach((folder) => {
      folderMap.set(folder.id, folder);
    });

    folders.forEach((folder) => {
      if (folder.parentId) {
        const parentFolder = folderMap.get(folder.parentId);
        if (parentFolder) {
          parentFolder.files = [...parentFolder.files, folder];
        }
      }
    });

    folders = folders.filter((folder) => !folder.parentId);
  }

  onMount(async () => {
    if (usedFor === "Documents") {
      configModal.header.contentCenter[0].text.text = $_("addDocuments");
    } else if (usedFor === "Handwriting") {
      configModal.header.contentCenter[0].text.text = $_("addHandwriting");
      mainContentUpload[2].input.value = "Handwriting (JPG, JPEG, PNG)";
      mainContentLink[2].input.value = "Handwriting (JPG, JPEG, PNG)";
      limitation = ["image/jpeg", "image/png", "image/jpg"];
    } else if (usedFor === "Audio") {
      configModal.header.contentCenter[0].text.text = $_("addAudio");
      mainContentUpload[2].input.value = "Audio (MP3, WAV, OGG)";
      mainContentLink[2].input.value = "Audio (MP3, WAV, OGG)";
      limitation = ["audio/mpeg", "audio/wav", "audio/ogg"];
    } else if (usedFor === "Data") {
      configModal.header.contentCenter[0].text.text = $_("addCsv");
      mainContentUpload[2].input.value = "Data (CSV)";
      mainContentLink[2].input.value = "Data (CSV)";
      limitation = ["text/csv"];
    }
  });

  $effect(() => {
    if (groupedFiles) {
      configModal.footer[0].button.buttonDisabled = groupedFiles === null;
    }
  });

  $effect(() => {
    if (files.length > 1) {
      configModal.footer[0].button.buttonDisabled = false;
      console.log("masuk files.length > 1");
    } else if (files.length === 1 && tabBarSelected === "Upload") {
      configModal.footer[0].button.buttonDisabled = tmpDocumentTitle === "" || !nameIsAvailable;
      console.log("masuk Upload files.length === 1");
    } else if (files.length === 1 && tabBarSelected === "Link") {
      configModal.footer[0].button.buttonDisabled = false;
      console.log("masuk Link files.length === 1");
    }
  });

  $effect(() => {
    if (urlCheckerRunning) {
      configModal.footer[0].button.buttonDisabled = true;
    } else {
      configModal.footer[0].button.buttonDisabled = linkImportLink ? false : true;
    }
  });

  $effect(() => {
    if (files.length > 0 && droppedItemType === "folder") {
      if (delayTimeout) {
        clearTimeout(delayTimeout);
      }
      delayTimeout = setTimeout(() => {
        processFolders(files);
        files = files[0];
        console.log("Folders", files);
      }, 1000);
    } else if (droppedItemType === "file") {
      console.log(files);
    }
  });
</script>

<Modal
  bind:open={modalObjectOpen}
  bodyClass="p-4 md:p-5 !space-y-4 flex-1 overflow-y-auto overscroll-contain"
>
  <svelte:fragment slot="header">
    <div class="flex items-center gap-4">
      {#if configModal?.header}
        {#if configModal.header?.iconLeft}
          <FeaturedIcon classesSvg={getIcons(configModal.header.iconLeft?.icon)} />
        {/if}
        <div>
          {#each configModal.header?.contentCenter as config}
            <Typography {...config.text} />
          {/each}
        </div>
      {/if}
    </div>
  </svelte:fragment>

  {#if configModal?.mainContent}
    {#each configModal?.mainContent as content}
      {#if content?.tabBar}
        <TabBar
          tabDatas={content.tabBar?.tabDatas}
          on:tabChange={(e) => {
            tabBarSelected = e.detail;
            files = [];
            configModal.mainContent =
              tabBarSelected === "Upload" ? mainContentUpload : mainContentLink;
            configModal.footer =
              tabBarSelected === "Upload" ? footerUpload : linkChecked ? footerUpload : footerLink;
            linkChecked = false;
            linkImportLink = "";
          }}
        />
      {/if}
      {#if content?.text}
        <Typography {...content.text} />
      {/if}
      {#if content?.input && content.input.id === "link" && linkChecked}
        <UploadFile
          title=""
          handleFileUploadWrapper={() => {}}
          file={files[0]}
          {files}
          bind:isProcess
          multiple={folderMode}
          actionButton={true}
          actionDelete={() => {
            files = [];
            linkChecked = false;
            linkImportLink = "";
            configModal.footer = footerLink;
          }}
          actionView={true}
          {folderMode}
          {handleDragOver}
          {handleDrop}
        />
      {:else if content?.input}
        <Input {...content.input} />
      {/if}
      {#if content?.uploadFile}
        {#if folderMode}
          <UploadFolder
            title={content.uploadFile?.title}
            handleFileUploadWrapper={() => handleFileUploadWrapper(true)}
            files={groupedFiles}
            bind:isProcess
            multiple={true}
            actionButton={true}
            actionView={true}
            actionDelete={() => {
              groupedFiles = null;
              console.log("cek groupedFiles", groupedFiles);
              console.log(
                "cek configModal.footer[0].button.buttonDisabled",
                configModal.footer[0].button.buttonDisabled,
              );
            }}
          />
        {:else}
          <UploadFile
            title={content.uploadFile?.title}
            handleFileUploadWrapper={() => handleFileUploadWrapper()}
            file={files[0]}
            {files}
            bind:isProcess
            multiple={true}
            actionButton={true}
            actionView={true}
            actionDelete={() => {
              files = [];
            }}
            {handleDragOver}
            {handleDrop}
          />
        {/if}
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
      {/if}
      {#if content?.multiSelect}
        <MultiSelect {...content.multiSelect} bind:selected={tags} options={alltags} />
      {/if}
      {#if content?.card && content.card?.cardData.contentCenter[0][0].text.text}
        <Card {...content.card} />
      {/if}
      {#if files.length === 1 && !folderMode}
        {#if content?.metaFile}
          {#each content.metaFile.content as content}
            {#if content?.text}
              <Typography {...content.text} />
            {/if}
            {#if content?.input}
              <Input {...content.input} bind:value={tmpDocumentTitle} />
              {#if !nameIsAvailable}
                <Helper color="red">
                  <span class="font-medium">{$_("nameAlreadyTaken")}</span>
                </Helper>
              {/if}
            {/if}
          {/each}
        {/if}
      {/if}
    {/each}
  {/if}

  <svelte:fragment slot="footer">
    {#if configModal?.footer}
      {#each configModal?.footer as item}
        {#if item?.button}
          <Button
            onclick={item.button?.onClick}
            disabled={item.button?.buttonDisabled || false}
            class="w-full">{$_(item.button?.text)}</Button
          >
        {/if}
      {/each}
    {/if}
  </svelte:fragment>
</Modal>
