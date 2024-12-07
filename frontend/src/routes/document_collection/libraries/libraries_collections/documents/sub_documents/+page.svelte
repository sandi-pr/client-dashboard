<script>
  import { onMount } from "svelte";
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
    LibraryCard,
  } from "$lib";
  import DataTableLibraryCollection from "$lib/components/DataTableLibraryCollection.svelte";
  import {
    SearchOutline,
    ExclamationCircleSolid,
    ChevronDownOutline,
    InfoCircleSolid,
  } from "flowbite-svelte-icons";
  import Trash01Icon from "$lib/icon/trash-01.svelte";
  import MoveIcon from "$lib/icon/move.svelte";
  import FilePlus02Icon from "$lib/icon/file-plus-02.svelte";
  import FolderPlusIcon from "$lib/icon/folder-plus.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import XCloseIcon from "$lib/icon/x-close.svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import AlertCircleIcon from "$lib/icon/alert-circle.svg";

  let subCollectionId;
  let collectionId;
  let libraryId;
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
  let fileNotSupport = $state(false);
  let buttonDisabled = $state(false);
  let buttonDisabledDelete = $state(false);
  let pagination = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });
  let paginationCollection = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpCollectionName = $state("");
  let tmpDescription = $state("");
  let tmpDocumentType = $state("Documents");
  let selectedTypeCollection = $state("Documents (TXT, PDF)");
  let acceptedFileTypes = $state(".txt, .pdf");
  let search = $state("");
  let tags = $state([]);
  let newTag = $state("");
  let files = $state([]);
  let attributeNotUnique = $state(false);
  let checkboxDelete = $state(false);
  let checkboxUpdate = $state(false);
  let tmpId = $state(null);
  let tabDatas = $state({
    selected: "Collections",
    data: [
      {
        name: "Collections",
      },
    ],
  });

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
  let dataSend = $state(null);
  let textarea = $state();

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
        description: response.result.attributes.meta.description,
        createdAt: format(new Date(response.result.attributes.createdAt), "dd-MM-yyyy"),
      };
    }
  }

  async function getDocuments() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      getCollection();

      if (search !== "") {
        pagination.page = 1;
      }

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
        return;
      }

      // console.log("documets", response.result.data);

      pagination = response.result.meta.pagination;

      let link = await $call.getPage("client-dashboard", "document-library-sub-documents-detail");

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
            return;
          }

          uploadBy = response.result.user.username;
        } catch (error) {
          console.log(error);
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
    }
  }

  async function createCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "create-library-collection", {
        URLQuery: {},
        URLParams: {},
        Params: {
          collection_name: tmpCollectionName,
          description: tmpDescription,
          parent_id: parseInt(collectionId),
          // bakrie
          // bot_id: localStorage.getItem("botIdStrapi").toString(),
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

      tmpCollectionName = "";
      tmpDescription = "";

      triggerToast("success");

      addCollectionModal = false;
      getDocuments();
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
              triggerToast("error");
              triggerLoading(false);
              console.log(response3);
              return;
            }

            if (response3?.result?.text !== null) {
              dataSend = {
                document_name: tmpDocumentTitle,
                description: tmpDocumentDescription,
                parent_id: parseInt(subCollectionId),
                documents: JSON.stringify(files),
                document_type: tmpDocumentType,
                category: tags.join(","),
                detected_text: response3.result.text,
                result_trocr: JSON.stringify(response3.result),
                edited: false,
              };

              triggerLoading(false);
            }
          }

          if (dataSend === null) {
            dataSend = {
              document_name: tmpDocumentTitle,
              description: tmpDocumentDescription,
              parent_id: parseInt(subCollectionId),
              documents: JSON.stringify(files),
              document_type: tmpDocumentType,
              category: tags.join(","),
            };
          }

          if (tmpDocumentType === "Handwriting" || tmpDocumentType === "Document Images") {
            previewModal = true;
            selectedFile = files[0];
          } else {
            sendDocument();
          }
        }
      } catch (error) {
        addDocumentsModal = true;
        console.log(error);
        triggerToast("error");
        triggerLoading(false);
        return;
      }
    }
  }

  async function sendDocument() {
    let token = localStorage.getItem("token");

    if (token === null) {
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
      if (response && response.errors) {
        const errors = response.errors.flatMap((error) =>
          error.errors.flatMap((subError) => subError.message),
        );

        if (errors.includes("This attribute must be unique")) {
          attributeNotUnique = true;
        } else {
          console.log("Pesan error :", errors);
          triggerToast("error");
        }
      }
      addDocumentsModal = true;
      triggerLoading(false);
      return;
    }

    files = [];
    dataSend = null;

    tmpDocumentTitle = "";
    tmpDocumentDescription = "";
    tags = [];
    attributeNotUnique = false;

    previewModal = false;
    addDocumentsModal = false;
    getDocuments();
    triggerLoading(false);

    triggerToast("created");
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
          ext: "." + file.name.split(".").pop(),
          type: tmpDocumentType,
        });
        tmpDocumentTitle = file.name.split(".").slice(0, -1).join(".");
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

      try {
        const resp = await $call.call("client-dashboard", "delete-library-document", {
          URLQuery: {},
          URLParams: {
            id: id.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log(resp);

        if (resp.status !== 200) {
          triggerLoading(false);
          popupModal = true;
          triggerToast("error");
          return;
        }

        triggerLoading(false);
        triggerToast("deleted");
        getDocuments();
      } catch (e) {
        console.log(e);
      }
    }
  }

  async function deleteDocument() {
    popupModal = false;
    triggerLoading(true);
    let token = localStorage.getItem("token");

    try {
      const resp = await $call.call("client-dashboard", "delete-library-document", {
        URLQuery: {},
        URLParams: {
          id: tmpId.toString(),
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      console.log(resp);

      if (resp.status !== 200) {
        triggerLoading(false);
        popupModal = true;
        triggerToast("error");
        return;
      }

      tmpId = null;

      triggerLoading(false);
      triggerToast("deleted");
      getDocuments();
    } catch (e) {
      console.log(e);
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
          triggerLoading(false);
          triggerToast("error");
          console.log(response);
          return;
        }

        triggerLoading(false);
        console.log(response);
        triggerToast("updated");
        getDocuments();
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
        console.log("cek response");
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
        description: collection.attributes.meta.description,
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
          tabDatas.data = [
            ...tabDatas.data,
            {
              name: "Sub Collections",
            },
          ];
        } else {
          tabDatas.data = tabDatas.data.filter((item) => item.name !== "Sub Collections");
        }
      } else {
        console.log("collections", collectionList);
        collectionList = await Promise.all(promise);
      }
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
          return;
        }

        tmpCollectionName = "";
        tmpDescription = "";

        editModal = false;
        triggerToast("updated");
        getCollections();
      } catch (error) {
        console.log(error);
        triggerToast("error");
      }
    } else {
      console.log("Token not found");
    }
  }

  async function deleteDocumentCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response2 = await $call.call("client-dashboard", "delete-library-collection", {
          URLQuery: {},
          URLParams: {
            id: tmpCollectionId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          console.log(response2);
          triggerToast("error");
          return;
        }

        triggerToast("deleted");

        getCollections();
        deleteModal = false;
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    } else {
      console.log("Token not found");
    }
  }

  const selectCollection = (item, subcollection = false) => {
    tmpSelectedCollection = item;
    if (!subcollection) {
      getCollections(false, true);
    }

    console.log("selected collection", tmpSelectedCollection);
  };

  onMount(() => {
    if (check("documentAi.documentLibrary:view") === false) {
      window.history.back();
      return;
    }
    const unsubscribe = page.subscribe(async ($page) => {
      subCollectionId = $page.url.searchParams.get("subCollectionId");
      collectionId = $page.url.searchParams.get("collectionId");
      libraryId = $page.url.searchParams.get("libraryId");
      if (collectionId === null) {
        goto(await $call.getPage("client-dashboard", "document-library"));
        return;
      }
      getCollections(true);
      getCollections();
      getDocuments();
    });
    return () => unsubscribe();
  });

  $effect(() => {
    checkAllState = collectionListDatas.every((data) => data.checked === true);
  });

  $effect(() => {
    buttonDisabledDelete = !collectionListDatas.some((data) => data.checked === true);
  });

  $effect(() => {
    buttonDisabled = files.length === 0 || tmpDocumentTitle === "" || tmpDocumentDescription === "";
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
            </div>
          </div>

          <div class="flex gap-4">
            {#if collectionListDatas.length > 0}
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
              <Button
                on:click={() => {
                  popupModal = true;
                  checkboxDelete = true;
                  // id
                }}
                disabled={buttonDisabledDelete}
                class="p-2 md:px-3 md:py-2 gap-2 text-fg-secondary-light dark:text-fg-secondary-dark border-secondary-light dark:border-secondary-dark hover:bg-secondary-light hover:dark:bg-secondary-dark hover:text-fg-secondary-light dark:hover:text-fg-secondary-dark"
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
                disabled={buttonDisabledDelete}
                class="p-2 md:px-3 md:py-2 gap-2 text-fg-secondary-light dark:text-fg-secondary-dark border-secondary-light dark:border-secondary-dark hover:bg-secondary-light hover:dark:bg-secondary-dark hover:text-fg-secondary-light dark:hover:text-fg-secondary-dark"
                outline
              >
                <MoveIcon svgClass="w-5 h-5" />
                <div class="hidden md:block">{$_("move")}</div>
              </Button>
            {/if}
            <Button
              disabled={!check("documentAi.documentLibrary.document:create") || $onProgress}
              outline
              on:click={() => {
                addDocumentsModal = true;
                tags = [];
                tmpDocumentTitle = "";
                tmpDocumentDescription = "";
                attributeNotUnique = false;
              }}
              class="p-2 md:px-3 md:py-2 gap-2 whitespace-nowrap"
            >
              <FilePlus02Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("addDocument")}</div>
            </Button>
          </div>
        </div>

        <Divider />

        {#if collectionListDatas.length > 0}
          <DataTableLibraryCollection
            {collectionListDatas}
            metadata={pagination}
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
          <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
        {/if}
      </div>
    </div>
  </div>
  <ToastContainer />
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
            }}>Document Images (JPG, JPEG, PNG, WEBP, IMG)</button
          >
          <button
            class="flex p-2 w-full"
            onclick={() => {
              selectTypeCollection("Handwriting");
            }}>Handwriting (JPG, JPEG, PNG)</button
          >
          <!-- <button
            class="flex p-2 w-full"
            onclick={() => {
              selectTypeCollection("Document Audio");
            }}>Document Audio (MP3, WAV, OGG)</button
          > -->
          <!-- <button
            class="flex p-2 w-full"
            onclick={() => {
              selectTypeCollection("CSV");
            }}>CSV</button
          > -->
        </ul>
      </div>
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
    {#if files.length > 0}
      <div class="flex flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={`${$_("documentTitle")}*`} />
        {#if !attributeNotUnique}
          <InputRw
            onEdit={true}
            bind:value={tmpDocumentTitle}
            placeholder={$_("enterDocumentTitle")}
          />
        {:else}
          <Input
            id="error"
            color="red"
            bind:value={tmpDocumentTitle}
            placeholder={$_("enterDocumentTitle")}
          />
          <Helper color="red">
            <span class="font-medium">{$_("nameAlreadyExist")}</span>
          </Helper>
        {/if}
        <!-- <InputRw
          onEdit={true}
          bind:value={tmpDocumentTitle}
          placeholder={$_("enterDocumentTitle")}
        /> -->
      </div>
      <div class="flex flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}*`} />
        <InputRw
          onEdit={true}
          bind:value={tmpDocumentDescription}
          placeholder={$_("enterDocumentDescription")}
        />
      </div>
    {/if}
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("labels")}`} />
      <div class="flex">
        <button
          class="flex w-full flex-wrap gap-2 items-center border border-primary-light dark:border-primary-dark rounded-lg py-2.5 px-2 bg-primary-light dark:bg-secondary-dark"
          onclick={() => {
            focusInput(inputRef);
          }}
        >
          {#each tags as tag}
            <span
              class="flex items-center space-x-0.5 bg-utility-brand-100-light dark:bg-utility-brand-100-dark text-brand-secondary-light dark:text-brand-secondary-dark rounded-full px-2"
            >
              <span class="text-sm pb-0.5">{tag}</span>
              <!-- svelte-ignore node_invalid_placement_ssr -->
              <button
                onclick={() => {
                  removeTag(tag);
                }}
                class="text-brand-primary-light dark:text-brand-primary-dark"
                ><XCloseIcon
                  svgClass="w-3 h-3 text-utility-brand-400-light dark:text-utility-brand-400-dark"
                /></button
              >
            </span>
          {/each}

          <input
            bind:this={inputRef}
            type="text"
            bind:value={newTag}
            placeholder={`${tags.length < 1 ? $_("enterNewLabel") : ""}`}
            onkeydown={handleKeyDown}
            class={`border-none flex-grow !pl-0 !py-0.5 text-sm bg-primary-light dark:bg-secondary-dark  focus:ring-transparent ${tags.length > 0 ? "max-w-min" : "w-full"}`}
          />
        </button>
        <button
          onclick={addTag}
          class="bg-secondary-light dark:bg-secondary-dark text-primary-light dark:text-primary-dark px-4 py-1 rounded"
        >
          Add
        </button>
      </div>
      <!-- <MultiSelect items={categories} bind:value={tmpCategories}></MultiSelect> -->
      <!-- <Select
        class="mt-2"
        items={categories}
        bind:value={tmpCategories}
        placeholder={$_("selectCategory")}
      />
      {#if categories.length === 0}
        <Alert color="yellow">
          <InfoCircleSolid slot="icon" class="w-5 h-5" />
          {$_("categoryIsEmptyAdd")}
          <button
            onclick={() => {
              addDocumentsModal = false;
              addCategoriesModal = true;
            }}
            class="font-semibold underline hover:text-yellow-800 dark:hover:text-yellow-900"
            >{$_("inHere")}</button
          >
        </Alert>
      {/if} -->
    </div>
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
    <Button onclick={deleteCategories} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>

<Modal title={$_("createNewCollection")} bind:open={addCollectionModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("collectionName")}*`} />
      <Input bind:value={tmpCollectionName} placeholder={$_("whatIsYourCollectionTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}*`} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={createCollection}
      disabled={!tmpCollectionName || !tmpDescription}
      class="w-full">{$_("addNewCollection")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editCollectionName")} bind:open={editModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("collectionName")}*`} />
      <Input bind:value={tmpCollectionName} placeholder={$_("whatIsYourCollectionTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}*`} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={editCollectionName}
      disabled={!tmpCollectionName && !tmpDescription}
      class="me-2 w-full">{$_("save")}</Button
    >
  </svelte:fragment>
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
      <Button color="primary" on:click={sendDocument} class="rounded-lg w-full">I’m sure</Button>
    </div>
  </svelte:fragment>
</Modal>
