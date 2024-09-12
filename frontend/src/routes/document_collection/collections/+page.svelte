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
  import { _ } from "svelte-i18n";
  import { CR3Call } from "hexacomponents-svelte";
  import {
    StateCard,
    Typography,
    CollectionCard,
    triggerToast,
    ToastContainer,
    UploadFile,
  } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { format } from "date-fns";
  import { Button, Modal, Input, P, Indicator, Tooltip } from "flowbite-svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import FileTypeIcon from "$lib/icon/file-type-icon.svelte";
  import { ExclamationCircleSolid, ChevronDownOutline } from "flowbite-svelte-icons";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import {
    triggerSelectCollection,
    triggerShowAI,
    triggerSelectDropdown,
    triggerReload,
    triggerLoading,
    showAI,
  } from "$lib";
  import LoadingContainer from "$lib/components/LoadingContainer.svelte";

  import { goto } from "$app/navigation";

  import { call } from "$lib";

  let datas = $state([
    // {
    //   id: "1",
    //   collectionName: "Bakrie’s Handwriting",
    //   totalDocuments: 2,
    //   type: "",
    //   status: "ready",
    // },
    // {
    //   id: "2",
    //   collectionName: "Study Scriptures of Qur’an Tafsir",
    //   totalDocuments: 2,
    //   type: "",
    //   status: "reading",
    // },
    // {
    //   id: "3",
    //   collectionName: "Laws and Regulations",
    //   totalDocuments: 0,
    //   type: "",
    //   status: "collecting data",
    // },
  ]);

  let gettingDatas = $state(false);

  let collections = $state([]);
  let documents = $state([]);
  let files = $state([]);
  let tmpCollectionName = $state("");
  let tmpDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpDocumentType = $state("");
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let addCollectionsModal = $state(false);
  let addDocumentsModal = $state(false);
  let editModal = $state(false);
  let deleteModal = $state(false);
  let dropdownOpen = $state(false);

  let acceptedFileTypes = $state(".txt, .pdf, .doc, .docx");

  let selectedTypeCollection = $state("Documents (TXT, PDF, DOC, DOCX)");

  let statusRAG = $state("");
  let timer;

  async function createCollection() {
    console.log(files);
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "create-document-collection", {
        URLQuery: {},
        URLParams: {},
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
        triggerToast("error");
        return;
      }

      console.log(response);

      // let response2 = await $call.call("hexa-ai-rag", "add-document-by-collection-id", {
      //   URLQuery: {},
      //   URLParams: {},
      //   Params: {
      //     collection_id: response.result.data.id.toString(),
      //   },
      //   Headers: {
      //     Authorization: "Bearer " + token,
      //   },
      // });

      // if (response2.status !== 200) {
      //   console.log(response2);
      //   triggerToast("error");
      //   return;
      // }

      // console.log(response2);

      tmpCollectionName = "";
      tmpDescription = "";

      triggerToast("success");

      addCollectionsModal = false;
      getCollections();
      // goto(await $call.getPage("client-dashboard", "document-collection"));
    }
  }

  async function getCollections() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      let page = 1;
      collections = [];
      documents = [];
      while (true) {
        let response = await $call.call("client-dashboard", "get-document-collections", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 25,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          gettingDatas = false;
          return;
        }

        if (response.result.data.length > 0) {
          collections = collections.concat(response.result.data);
          page++;
        } else {
          break;
        }
      }

      for (const collection of collections) {
        let page = 1;
        while (true) {
          let response3 = await $call.call("client-dashboard", "get-documents-by-collectionid", {
            URLQuery: {},
            URLParams: {},
            Params: {
              collection_id: collection.id.toString(),
              page: page,
              page_size: 25,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response3.status !== 200) {
            console.log(response3);
            gettingDatas = false;
            return;
          }

          if (response3.result.data.length > 0) {
            documents = documents.concat(response3.result.data);
            page++;
          } else {
            break;
          }
        }
      }

      let collectionList = await $call.getPage(
        "client-dashboard",
        "document-collection-list-document",
      );

      datas = collections.map((collection) => {
        const filteredDocuments = documents.filter(
          (document) => document.attributes.collection_id === collection.id,
        );

        // Hitung total dokumen
        const totalDocuments = filteredDocuments.length;

        // Hitung jumlah dokumen dengan status null
        const nullDocuments = filteredDocuments.filter(
          (document) => document.attributes.status === null,
        ).length;

        // Hitung jumlah dokumen dengan status angka
        const numberDocuments = filteredDocuments.filter(
          (document) => typeof document.attributes.status === "number",
        ).length;

        // Hitung total dokumen yang Ready
        const readyDocuments = filteredDocuments.filter(
          (document) => document.attributes.status === "Ready",
        ).length;

        // Tentukan status koleksi berdasarkan kondisi yang diberikan
        let collectionStatus;
        if (nullDocuments === totalDocuments) {
          collectionStatus = null;
        } else if (nullDocuments > 0 || numberDocuments > 0) {
          collectionStatus = "reading";
        } else if (numberDocuments === totalDocuments && readyDocuments === totalDocuments) {
          collectionStatus = "ready";
        }
        return {
          collectionId: collection.id,
          activationTime: format(collection.attributes.createdAt, "dd-MM-yyyy HH:mm"),
          collectionName: collection.attributes.collection_name,
          description: collection.attributes.description,
          totalDocuments: totalDocuments,
          status: collectionStatus,
          href: `${collectionList}?collectionId=${collection.id}`,
          webchat: collection.attributes.webchat_tmp === true ? true : false,
        };
      });

      gettingDatas = false;

      // timer = setTimeout(getCollections, 2000);
    }
  }

  async function editCollectionName() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "update-document-collection", {
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
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        for (const document of documents) {
          if (document.attributes.collection_id === tmpCollectionId) {
            let response = await $call.call("client-dashboard", "delete-document", {
              URLQuery: {},
              URLParams: {
                id: document.id.toString(),
              },
              Params: {},
              Headers: {
                Authorization: "Bearer " + token,
              },
            });
            if (response.status !== 200) {
              console.log(response);
              triggerToast("error");
              return;
            }
          }
        }
        let response2 = await $call.call("client-dashboard", "delete-document-collection", {
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

        console.log("ayam", response2);

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

  async function addDocumentByCollection() {
    files[0].document_title = tmpDocumentTitle;
    files[0].description = tmpDocumentDescription;

    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        triggerLoading(true);
        let response = await $call.call("client-dashboard", "create-document", {
          URLQuery: {},
          URLParams: {},
          Params: {
            documents: JSON.stringify(files),
            collection_id: tmpCollectionId.toString(),
            document_type: tmpDocumentType,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          triggerToast("error");
          triggerLoading(false);
          return;
        }

        console.log(response);

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

            if (response3.status !== 200) {
              triggerToast("error");
              triggerLoading(false);
              console.log(response);
              return;
            }

            console.log("trocr-predict", response3);

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

              if (response4.status !== 200) {
                triggerToast("error");
                triggerLoading(false);
                console.log(response4);
                return;
              }

              console.log("ayam", response4);
            }
          }
        }
      } catch (error) {
        triggerToast("error");
        triggerLoading(false);
        console.log(error);
      }

      files = [];

      tmpCollectionId = "";
      getCollections();
      triggerLoading(false);
      addDocumentsModal = false;

      triggerToast("created");
    }
  }

  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  const handleFileUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = acceptedFileTypes;
    inputElement.multiple = false;

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file) {
        const base64Data = await readFileAsBase64(file);
        files.push({
          name: file.name,
          document_title: tmpDocumentTitle,
          description: tmpDocumentDescription,
          size: formatBytes(file.size),
          base64: base64Data,
        });
      } else {
        alert("Invalid file. Please select a valid file.");
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

  const selectTypeCollection = (type) => {
    tmpDocumentType = type;

    if (type === "Documents") {
      selectedTypeCollection = "Documents (TXT, PDF, DOC, DOCX)";
      acceptedFileTypes = ".txt, .pdf, .doc, .docx";
    } else if (type === "Images") {
      selectedTypeCollection = "Images (JPG, JPEG, PNG, WEBP, IMG)";
      acceptedFileTypes = ".jpg, .jpeg, .png, .webp, .img";
    } else if (type === "Handwriting") {
      selectedTypeCollection = "Handwriting (JPG, JPEG, PNG)";
      acceptedFileTypes = ".jpg, .jpeg, .png";
      // acceptedFileTypes = ".pdf, .jpg, .jpeg, .png, .webp, .img";
    }

    files = [];

    dropdownOpen = false;
  };

  const classStatusType = (type) => {
    if (type === "Ready") {
      return "bg-hc-success-50 border-hc-success-200 text-hc-success-700 dark:bg-hc-success-900 dark:border-hc-success-700 dark:text-hc-success-300";
    } else if (type === "Not Ready") {
      return "bg-hc-error-50 border-hc-error-200 text-hc-error-700 dark:bg-hc-error-900 dark:border-hc-error-700 dark:text-hc-error-300";
    } else if (type === "On Progress") {
      return "bg-hc-warning-50 border-hc-warning-200 text-hc-warning-700 dark:bg-hc-warning-900 dark:border-hc-warning-700 dark:text-hc-warning-300";
    }
  };

  // async function next() {
  //   getUsers(pagination.page + 1);
  // }

  // async function previous() {
  //   getUsers(pagination.page - 1);
  // }

  async function getStatusRag() {
    let token = sessionStorage.getItem("token");
    try {
      let rag = await $call.call("hexa-ai-rag", "get-rag-status", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      let status = rag.result.rag_status;
      if (statusRAG.includes("!")) {
        let status2 = status.split("!")[0].trim();
        if (
          status2.includes("(NOVERB)") ||
          status2.includes("(INOVERB)") ||
          status2.includes("(!INOVERB)")
        ) {
          let status3 = status2
            .replace("(NOVERB)", "")
            .replace("(INOVERB)", "")
            .replace("(!INOVERB)", "")
            .replace("(!NOVERB)", "");
          statusRAG = status3;
        } else {
          statusRAG = status2;
        }
      } else {
        statusRAG = status;
      }
      console.log(rag); // Mengganti statusRAG dengan rag

      if (statusRAG === "Ready") {
        return;
      }
    } catch (error) {
      console.error("Error:", error);
    }
    timer = setTimeout(getStatusRag, 3000);
  }

  async function selectCollection(collection) {
    console.log("ayam collections");
    triggerShowAI(true);
    localStorage.setItem("showAI", "true");
    triggerSelectDropdown("collections");
    localStorage.setItem("selectedDropdown", "collections");
    triggerReload(true);
    triggerSelectCollection(collection);
    localStorage.setItem("selectedCollection", JSON.stringify(collection));
    return;
  }

  async function UpdateCollectionWebchat() {
    let token = sessionStorage.getItem("token");
    try {
      let res = await $call.call("client-dashboard", "update-collection-webchat", {
        URLQuery: {},
        URLParams: {},
        Params: {
          collection_id: tmpCollectionId,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (res.status !== 200) {
        console.log(res);
        triggerToast("error");
        return;
      }

      getCollections();
      triggerToast("success");
      return;
    } catch (error) {
      console.error("Error:", error);
    }
  }

  let buttonDisabled = $state(false);

  $effect(() => {
    buttonDisabled = files.length === 0 || tmpDocumentTitle === "";
  });

  onMount(async () => {
    // await getStatusRag();
    await getCollections();
  });

  onDestroy(() => {
    if (timer) {
      clearTimeout(timer);
    }
  });

  // async function addDocument() {
  //   goto(await $call.getPage("client-dashboard", "document-collection-create"));
  //   return;
  // }

  async function getDocumentCollectionDetail(id) {
    return (
      `${await $call.getPage("client-dashboard", "document-collection-list-document")}` +
      "?collectionId=" +
      id
    );
  }

  async function askToAI(id) {
    return (
      `${await $call.getPage("client-dashboard", "document-collection-chat")}` +
      "?collectionId=" +
      id
    );
  }
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text={$_("collectionList")} />
          <Button on:click={() => (addCollectionsModal = true)} icon={PlusIcon}
            ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" />
            {$_("addNew")}</Button
          >
        </div>

        <!-- Bots List Data -->
        <div>
          <!-- <Search
            class="bg-white dark:bg-hc-graydark-950 mb-6 rounded-lg"
            placeholder={$_("searchBotName")}
            bind:value={searchTerm}
          /> -->
          {#if datas.length > 0}
            <div
              class={`${$showAI ? "grid grid-cols-1 xl:grid-cols-3" : "grid grid-cols-2 lg:grid-cols-4"} gap-4`}
            >
              {#each datas as document}
                <CollectionCard
                  collectionName={document.collectionName}
                  totalDocuments={document.totalDocuments}
                  status={document.status}
                  webchat={document.webchat}
                  functButton={() => {
                    selectCollection({
                      id: document.collectionId,
                      name: document.collectionName,
                      description: document.description,
                      totalDocuments: document.totalDocuments,
                    });
                  }}
                  functAddDocuments={() => {
                    tmpCollectionId = document.collectionId;
                    selectTypeCollection("Documents");
                    addDocumentsModal = true;
                  }}
                  functEdit={() => {
                    tmpCollectionId = document.collectionId;
                    tmpCollectionName = document.collectionName;
                    tmpDescription = document.description;
                    editModal = true;
                  }}
                  functDelete={() => {
                    tmpCollectionId = document.collectionId;
                    deleteModal = true;
                  }}
                  funcSelectCollectionWebchat={() => {
                    tmpCollectionId = document.collectionId;
                    UpdateCollectionWebchat();
                  }}
                  funcDetailCollection={() => {
                    goto(document.href);
                  }}
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
  <ToastContainer />
</div>

<Modal title={$_("createNewCollection")} bind:open={addCollectionsModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("collectionName")}*`} />
      <Input bind:value={tmpCollectionName} placeholder={$_("whatIsYourCollectionTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={createCollection} disabled={!tmpCollectionName} class="w-full"
      >{$_("addNewCollection")}</Button
    >
  </svelte:fragment>
</Modal>

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
            class="flex p-2"
            onclick={() => {
              selectTypeCollection("Documents");
            }}>Documents (TXT, PDF, DOC, DOCX)</button
          >
          <button
            class="flex p-2"
            onclick={() => {
              selectTypeCollection("Images");
            }}>Images (JPG, JPEG, PNG, WEBP, IMG)</button
          >
          <button
            class="flex p-2"
            onclick={() => {
              selectTypeCollection("Handwriting");
            }}>Handwriting (JPG, JPEG, PNG)</button
          >
        </ul>
      </div>
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("documentTitle")}*`} />
      <Input bind:value={tmpDocumentTitle} placeholder={$_("enterDocumentTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Input bind:value={tmpDocumentDescription} placeholder={$_("enterDocumentDescription")} />
    </div>
    <UploadFile
      title={`${$_("uploadDocuments")}*`}
      handleFileUploadWrapper={() => handleFileUploadWrapper()}
      file={files[0]}
      actionButton={true}
      actionDelete={() => (files = [])}
      isFlex
    />
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={addDocumentByCollection} disabled={buttonDisabled} class="w-full"
      >{$_("addNewDocument")}</Button
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
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={editCollectionName} class="me-2">{$_("save")}</Button>
  </svelte:fragment>
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteDocumentCollection} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>

<!-- <LoadingContainer /> -->
