<script>
    import {
    Typography,
    SectionHeader,
    ActionCard,
    StateCard,
    call,
    triggerToast,
    triggerLoading,
    triggerToastResponse,
    setSelectedConversation,
    ModalObject,
  } from "$lib";
  import { Button, Modal } from "flowbite-svelte";
  import { format } from "date-fns";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  
  import DataTable from "$lib/components/DataTable.svelte";
  import AlertCircleIcon from "$lib/icon/alert-circle.svg";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { goto } from "$app/navigation";

  let actionItems = $state([]);
  let datas = $state([]);
  let files = $state([]);
  let isProcess = $state(false);
  let fileNotSupport = $state(false);
  let datasConversation = $state([]);
  let configs = $state({});
  let dataSend = $state({});
  let previewModal = $state(false);
  let addDocumentsModal = $state(false);
  let tmpId = $state(null);
  let popupModal = $state(false);
  let tmpName = $state(null);
  let selectedFile = $state(null);
  let paginationHandwriting = $state({
    page: 1,
    pageCount: 1,
    pageSize: 5,
    total: 10,
  });
  let gettingDatas = $state(true);
  let tmpDocumentType = $state("Handwriting");
  let acceptedFileTypes = $state(".jpg, .jpeg, .png");
  let tmpDocumentTitle = $state("");

  // let datasConversation = $state([
  //   {
  //     id: 1,
  //     conversatin_name: "NewChat 21.02.2025",
  //     actions: {
  //       objects: [
  //         {
  //           href: "",
  //           name: "View",
  //         },
  //       ],
  //     },
  //   }
  // ]);

  async function getActionItems() {
    actionItems = [
      {
        onclick: () => {
          addDocumentsModal = true;
        },
        icon: "penToolPlus",
        actionName: "Finetune for Handwriting",
        desc: "desc Finetune for Handwriting",
      },
      {
        href: await $call.getPage("client-dashboard", "finetunings-conversation"),
        icon: "messageChatSquare",
        actionName: "Finetune for Conversation",
        desc: "desc Finetune for Conversation",
      },
    ];
  }

  const getDatasetsHandwriting = async (search) => {
    let token = localStorage.getItem("token");
    if (token) {
      if (search !== "") {
        paginationHandwriting.page = 1;
      }

      const response = await $call.call("client-dashboard", "get-datasets-handwriting", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: paginationHandwriting.page,
          page_size: paginationHandwriting.pageSize,
          search: search,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerToastResponse(response);
        return;
      }

      console.log("response getDatasetsHandwriting", response);

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

  const getListConversations = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "get-conversations", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: 1,
          page_size: 10,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerToast("error");
        return;
      }

      let promise = response.result.data.map(async (conversation) => {
        let title = conversation.attributes.title;
        let id = conversation.attributes.conversation_id;

        return {
          id: id,
          conversation_name: title,
          updateAt: format(new Date(conversation.attributes.updatedAt), "MMM d, yyyy"),
          actions: {
            objects: [
              {
                onclick: () => {
                  setSelectedConversation(conversation.id);
                  goto("/finetunings/conversation");
                },
                icon: "eye",
                name: "View",
              },
            ],
          },
        };
      });

      datasConversation = await Promise.all(promise);

      console.log("datasConversation", datasConversation);

      gettingDatas = false;
    }
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

          tmpDocumentTitle = file.name.split(".").slice(0, -1).join(".");

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
            let originalName = files.length > 1 ? files[i].name.split(".").slice(0, -1).join(".") : tmpDocumentTitle;
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

  async function deleteDatasetHandwriting() {
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

  let hrefHandwriting = $state("");
  let hrefConversation = $state("");

  async function getHrefAction(href) {
    let link = "";
    if (href === "handwriting") {
      link = await $call.getPage("client-dashboard", "finetunings-handwriting");
      return link;
    } else if (href === "conversation") {
      link = await $call.getPage("client-dashboard", "finetunings-conversation");
      return link;
    }
  }

  onMount(async () => {
    hrefHandwriting = await getHrefAction("handwriting");
    hrefConversation = await getHrefAction("conversation");
    getActionItems();
    getDatasetsHandwriting();
    getListConversations();
  })
</script>

<div class="flex flex-col gap-8 px-4 pt-4 pb-24">
  <Typography preset="Display sm/Semibold" text={$_("finetuning")} color="primary" />
  <div class="flex flex-col gap-4">
    <SectionHeader title={$_("quickActions")} />
    <ActionCard {actionItems} />
  </div>
  <!-- Handwriting Finetune -->
  <div class="flex flex-col">
    <SectionHeader title={$_("handwritingFinetuneHistory")} hrefAction={hrefHandwriting} />
    {#if datas.length > 0}
      <DataTable
        bind:datas={datas}
        bind:configs
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
  <!-- Conversation Finetune -->
  <div class="flex flex-col">
    <SectionHeader title={$_("conversationFinetuneHistory")} hrefAction={hrefConversation} />
    {#if datas.length > 0}
      <DataTable bind:datas={datasConversation} bind:configs />
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

<ModalObject
  bind:modalObjectOpen={addDocumentsModal}
  {handleFileUploadWrapper}
  bind:isProcess
  bind:files
  bind:tmpDocumentTitle
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
      <img
        src={`data:${"Handwriting"};base64,${files[0].base64}`}
        style="max-width: 80%; height: auto;"
        alt="PNG Preview"
      />
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
    <Button onclick={deleteDatasetHandwriting} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>