<script>
  import { Input, Button, Modal } from "flowbite-svelte";
  import { Typography, StateCard, ModalObject, call, check, triggerToastResponse, triggerToast, triggerLoading } from "$lib";
  import { format } from "date-fns";
  import { onMount } from "svelte";
  import { SearchOutline, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { _ } from "svelte-i18n";
  
  import PdfViewer from "svelte-pdf";
  import DataTable from "$lib/components/DataTable.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import PenToolPlusIcon from "$lib/icon/pen-tool-plus.svelte";
  import AlertCircleIcon from "$lib/icon/alert-circle.svg";
  import DataFlow from "$lib/icon/dataflow-03.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  
  let datas = $state([]);
  let files = $state([]);
  let configs = $state({});
  let dataSend = $state({});
  let paginationHandwriting = $state({
    page: 1,
    pageCount: 1,
    pageSize: 5,
    total: 10,
  });
  let gettingDatas = $state(true);
  let isProcess = $state(false);
  let addDocumentsModal = $state(false);
  let fileNotSupport = $state(false);
  let previewModal = $state(false);
  let popupModal = $state(false);
  let disableButtonHandwriting = $state(false);
  let selectedFile = $state(null);
  let statusRequest = $state(null);
  let tmpId = $state(null);
  let tmpName = $state(null);
  let searchTerm = $state("");
  let tmpDocumentType = $state("Handwriting");
  let selectedTypeCollection = $state("Handwriting (JPG, JPEG, PNG)");
  let acceptedFileTypes = $state(".jpg, .jpeg, .png");

  const getDatasetsHandwriting = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      if (searchTerm !== "") {
        paginationHandwriting.page = 1;
      }

      const response = await $call.call("client-dashboard", "get-datasets-handwriting", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: paginationHandwriting.page,
          page_size: paginationHandwriting.pageSize,
          search: searchTerm,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerToastResponse(response);
        return;
      }

      if (response.result.data.length < 1) {
        disableButtonHandwriting = true;
      } else {
        disableButtonHandwriting = false;
      }

      console.log("response getDatasetsHandwriting" + paginationHandwriting.page, response);

      paginationHandwriting = response.result.meta.pagination;

      configs = {
        pagination: response.result.meta.pagination,
      };

      let promise = response.result.data.map(async (data) => {
        let link = await $call.getPage("client-dashboard", "datasets-detail");

        return {
          id: data.id,
          fileName: data.attributes.name + "." + data.attributes.path.split(".").pop(),
          // fileext: data.attributes.path.split(".").pop(),
          filesize: data.attributes.file_size,
          createdAt: format(new Date(data.attributes.createdAt), "MMM d, yyyy"),
          actions: {
            objects: [
              {
                href: link + "?datasetId=" + data.id,
                icon: "eye",
                name: "View",
              },
              {
                onclick: () => {
                  tmpId = data.id;
                  tmpName = data.attributes.name;
                  popupModal = true;
                },
                icon: "trash",
                name: "Delete",
              },
            ],
          },
        };
      });

      datas = await Promise.all(promise);

      gettingDatas = false;
    }
  };

  const handleSearch = (evt) => {
    console.log("clicked");
    getDatasetsHandwriting();
  };

  const handleFileUploadWrapper = () => {

    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = acceptedFileTypes;
    inputElement.multiple = true;
    inputElement.replaceWith(inputElement.cloneNode(true));

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
            size: formatBytes(file.size),
            base64: base64Data,
            ext: "." + file.name.split(".").pop(),
            type: tmpDocumentType,
          });

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

  async function addDatasetHandwriting() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        addDocumentsModal = false;
        triggerLoading(true);

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
            let originalName = files[i].name.split(".").slice(0, -1).join(".");
            let newName = originalName;

            files[i].name = newName + files[i].ext;

            dataSend = {
              name: originalName,
              documents: JSON.stringify([files[i]]),
              detected_text: response3.result.text,
              result_trocr: JSON.stringify(response3.result),
            };

            if (files.length < 2) {
              triggerLoading(false);
            }
          }

          if (files.length > 1) {
            await sendDocument();
          } else {
            previewModal = true;
            selectedFile = files[0];
          }
        }

        addDocumentsModal = false;

        if (files.length > 1) {
          getDatasetsHandwriting();
          files = [];
          dataSend = null;
          previewModal = false;
          triggerLoading(false);
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

  async function sendDocument() {
    let token = localStorage.getItem("token");

    if (token === null) {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }

    triggerLoading(true);
    let response = await $call.call("client-dashboard", "create-dataset-handwriting", {
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

    if (files.length === 1) {
      getDatasetsHandwriting();
      triggerLoading(false);
      files = [];
      dataSend = null;
      previewModal = false;
      triggerLoading(false);
      triggerToast($_("succesfullyUploaded"), true);
    }

    return;
  }

  async function deleteDataset() {
    popupModal = false;
    triggerLoading(true);
    let token = localStorage.getItem("token");

    if (token === null) {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }

    try {
      const response = await $call.call("client-dashboard", "delete-dataset-handwriting", {
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
      getDatasetsHandwriting();
    } catch (error) {
      console.error("Error during document deletion:", error);
      triggerToast($_("failedDeleteDocument"), true, "error");
    }
  }

  async function processDatasetHandwriting() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        const response = await $call.call("client-dashboard", "request-fine-tuning-handwriting", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          // statusRequestHandwriting = "REQUEST";
          triggerToast("created");
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function getStatusRequestHandwriting() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        const response = await $call.call(
          "client-dashboard",
          "get-status-fine-tuning-handwriting",
          {
            URLQuery: {},
            URLParams: {},
            Params: {},
            Headers: {
              Authorization: "Bearer " + token,
            },
          },
        );
        if (response.status === 200) {
          statusRequest = response.result.status;
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  // svelte-ignore non_reactive_update
  async function nextPage() {
    if (paginationHandwriting.page >= paginationHandwriting.pageCount) return;
    paginationHandwriting.page++;

    console.log("cek next paginationHandwriting", paginationHandwriting.page);
    await getDatasetsHandwriting();
  };

  // svelte-ignore non_reactive_update
  async function prevPage() {
    if (paginationHandwriting.page <= 1) return;
    paginationHandwriting.page--;
    await getDatasetsHandwriting();
  };

  onMount(() => {
    getDatasetsHandwriting();
    getStatusRequestHandwriting();
  })
</script>

<div class="relative h-full flex flex-col gap-8 px-4 pt-4 pb-24">
  <div class="flex items-center justify-between">
    <Typography preset="Display sm/Semibold" text={$_("handwritingDatasetHistory")} color="primary" propClass="grow" />
    <div class="max-w-64">
      <Input
        id="search"
        type="search"
        placeholder={$_("search")}
        bind:value={searchTerm}
        on:input={handleSearch}
      >
        <SearchOutline
          slot="left"
          class="w-5 h-5 text-quinary-light dark:text-quinary-dark"
        />
      </Input>
    </div>
    <Button
      on:click={() => {
        addDocumentsModal = true;
      }}
      class="p-2 md:px-3 ms-3 md:py-2.5 gap-2 whitespace-nowrap"
    >
      <PenToolPlusIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
      <div class="hidden md:block">{$_("addNewDataset")}</div>
    </Button>
  </div>
  <!-- Handwriting Finetune -->
  <div class="flex flex-col gap-4">
    {#if datas.length > 0}
      <DataTable
        bind:datas={datas}
        bind:configs
        bind:next={nextPage}
        bind:previous={prevPage}
      />
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

  <!-- <div class="fixed bottom-4"> -->
    <Button
      disabled={statusRequest === "process" ||
        !check("aiManagement.datasets:process") ||
        disableButtonHandwriting}
      on:click={() => {
        processDatasetHandwriting();
      }}
      outline
      class="!flex hover:bg-primary-light dark:hover:bg-primary-dark sticky bottom-0 w-full h-12 rounded-2xl"
    >
      <DataFlow
        svgClass="w-5 h-5 text-utility-brand-700-light dark:text-utility-brand-700-dark"
      />
      &nbsp; &nbsp;
      <Typography preset="Text md/Semibold" text={$_("process")} color="brand-secondary" />
    </Button>
  <!-- </div> -->
</div>

<ModalObject
  bind:modalObjectOpen={addDocumentsModal}
  {handleFileUploadWrapper}
  bind:isProcess
  bind:files
  addDocumentByCollection={addDatasetHandwriting}
  usedFor={"Handwriting"}
  maxDepth={1}
/>

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
      class="flex w-1/2 justify-center items-center py-6 rounded-lg bg-secondary-light dark:bg-secondary-dark"
    >
      {#if files[0].ext.includes(".png") || files[0].ext.includes(".jpg") || files[0].ext.includes(".jpeg")}
        <img
          src={`data:${"Handwriting"};base64,${files[0].base64}`}
          style="max-width: 80%; height: auto;"
          alt="PNG Preview"
        />
      {:else if files[0].ext.includes(".pdf")}
        <PdfViewer
          showBorder={false}
          showButtons={false}
          scale={1}
          data={atob(files[0].base64)}
        />
      {/if}
    </div>
    <!-- detected handwriting -->
    {#if files[0].type === "Handwriting"}
      <div class="flex w-full flex-col gap-3 px-4 pb-1">
        <Typography preset="Text md/Semibold" text="Handwriting Information" color="secondary" />
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/Medium" text="Detected Text" color="secondary" />
          <textarea
            bind:value={dataSend.detected_text}
            id="chat"
            rows="4"
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
          files = [];
        }}
        class="w-full">{$_("cancel")}</Button
      >
      <Button color="primary" on:click={sendDocument} class="rounded-lg w-full">{$_("yesImSure")}</Button>
    </div>
  </svelte:fragment>
</Modal>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
      {$_("thisDocument")}
      "{tmpName}"
      <!-- {$_("thisDocument")} -->
    </h3>
    <Button onclick={deleteDataset} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>