<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { Button, Modal, Input, Alert, Select, P } from "flowbite-svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { format } from "date-fns";

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
  } from "$lib";
  import DataTableLibraryCollection from "$lib/components/DataTableLibraryCollection.svelte";
  import {
    SearchOutline,
    ExclamationCircleSolid,
    ChevronDownOutline,
    InfoCircleSolid,
  } from "flowbite-svelte-icons";
  import Trash01Icon from "$lib/icon/trash-01.svelte";
  import FilePlus02Icon from "$lib/icon/file-plus-02.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import XCloseIcon from "$lib/icon/x-close.svelte";

  let collectionId;
  let inputRef;
  let checkAllState = $state(false);
  let addDocumentsModal = $state(false);
  let addCategoriesModal = $state(false);
  let popupModal = $state(false);
  let deleteCategoriesModal = $state(false);
  let dropdownOpen = $state(false);
  let fileNotSupport = $state(false);
  let buttonDisabled = $state(false);
  let buttonDisabledDelete = $state(false);
  let pagination = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let tmpDocumentType = $state("Documents");
  let selectedTypeCollection = $state("Documents (TXT, PDF)");
  let acceptedFileTypes = $state(".txt, .pdf");
  let tags = $state([]);
  let newTag = $state("");
  let files = $state([]);

  let collectionDatas = $state({
    title: "",
    description: "",
    createdAt: "",
  });
  let collectionListDatas = $state([]);

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

      let response = await $call.call("client-dashboard", "get-library-documents", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: pagination.page,
          page_size: pagination.pageSize,
          parent_id: parseInt(collectionId),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      console.log("documets", response.result.data);

      pagination = response.result.meta.pagination;

      let link = await $call.getPage("client-dashboard", "document-library-detail");

      let promise = response.result.data.map(async (collection) => {
        return {
          id: collection.id,
          filename: collection.attributes.name,
          fileext: collection.attributes.path.split(".").pop(),
          label:
            collection.attributes.meta.category === undefined
              ? []
              : collection.attributes.meta.category.split(","),
          type: collection.attributes.meta.document_type,
          filesize: collection.attributes.meta.file_size,
          createdAt: format(new Date(collection.attributes.createdAt), "yyyy-mm-dd"),
          href: link + "?documentId=" + collection.id,
          checked: false,
        };
      });

      collectionListDatas = await Promise.all(promise);
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
  };

  async function addDocumentByCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        addDocumentsModal = false;
        triggerLoading(true);
        let response = await $call.call("client-dashboard", "create-library-document", {
          URLQuery: {},
          URLParams: {},
          Params: {
            document_name: tmpDocumentTitle,
            description: tmpDocumentDescription,
            parent_id: parseInt(collectionId),
            documents: JSON.stringify(files),
            document_type: tmpDocumentType,
            category: tags.join(","),
          },
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

        if (tmpDocumentType === "Handwriting") {
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

            if (response3?.result?.text !== null) {
              const response4 = await $call.call(
                "client-dashboard",
                "update-document-handwriting",
                {
                  URLQuery: {},
                  URLParams: {
                    id: response.result.data[i].id.toString(),
                  },
                  Params: {
                    detected_text: response3.result.text,
                    result_trocr: JSON.stringify(response3.result),
                    edited: false,
                    document_title: tmpDocumentTitle,
                    description: tmpDocumentDescription,
                  },
                  Headers: {
                    Authorization: "Bearer " + token,
                  },
                },
              );

              if (response3.status !== 200) {
                addDocumentsModal = true;
                triggerToast("error");
                triggerLoading(false);
                console.log(response4);
                return;
              }

              console.log("ayam", response4);
              triggerLoading(false);
            }
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
      tags = [];

      addDocumentsModal = false;
      getDocuments();
      triggerLoading(false);

      triggerToast("created");
    }
  }

  const selectTypeCollection = (type) => {
    tmpDocumentType = type;

    if (type === "Documents") {
      selectedTypeCollection = "Documents (TXT, PDF)";
      acceptedFileTypes = ".txt, .pdf";
    } else if (type === "Images") {
      selectedTypeCollection = "Images (JPG, JPEG, PNG, WEBP, IMG)";
      acceptedFileTypes = ".jpg, .jpeg, .png, .webp, .img";
    } else if (type === "Handwriting") {
      selectedTypeCollection = "Handwriting (PDF, JPG, JPEG, PNG, WEBP, IMG)";
      acceptedFileTypes = ".jpg, .jpeg, .png";
    } else if (type === "Document Audio") {
      selectedTypeCollection = "Document Audio (MP3)";
      acceptedFileTypes = ".mp3";
    } else if (type === "CSV") {
      selectedTypeCollection = "CSV";
      acceptedFileTypes = ".csv";
    }

    files = [];

    dropdownOpen = false;
  };

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
        tmpDocumentTitle = file.name.split('.').slice(0, -1).join('.');
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

  onMount(() => {
    const unsubscribe = page.subscribe(async ($page) => {
      collectionId = $page.url.searchParams.get("collectionId");
      if (collectionId === null) {
        goto(await $call.getPage("client-dashboard", "document-library"));
        return;
      }
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
            <Input
              id="search"
              type="search"
              class="bg-white dark:bg-hc-graydark-950"
              placeholder={$_("search")}
              on:change={handleSearch}
            >
              <SearchOutline
                slot="left"
                class="w-5 h-5 text-quinary-light dark:text-quinary-dark"
              />
            </Input>
            <Button
              on:click={() => {
                popupModal = true;
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
                addDocumentsModal = true;
                tags = [];
                tmpDocumentTitle = "";
                tmpDocumentDescription = "";
              }}
              class="p-2 md:px-3 md:py-2 gap-2 whitespace-nowrap"
            >
              <FilePlus02Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("addNew")}</div>
            </Button>
          </div>
        </div>

        <div
          class="mx-4 sm:mx-0 border-t-[1px] border-hc-graylight-300 dark:border-hc-graydark-700"
        ></div>

        {#if collectionListDatas.length > 0}
          <DataTableLibraryCollection
            {collectionListDatas}
            metadata={pagination}
            {checkAllState}
            {checkSingle}
            {checkAll}
            next={nextPage}
            previous={prevPage}
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
              selectTypeCollection("Document Audio");
            }}>Document Audio (MP3)</button
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
        <InputRw onEdit={true} bind:value={tmpDocumentTitle} placeholder={$_("enterDocumentTitle")} />
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
    <Button onclick={deleteDocuments} color="red" class="me-2">{$_("yesImSure")}</Button>
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

<!-- <Modal title={$_("trainingModel")} bind:open={trainingModelModal} size="xs">
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
</Modal> -->
