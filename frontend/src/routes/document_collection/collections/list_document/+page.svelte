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
  import DataTableListCollection from "$lib/components/DataTableListCollection.svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { StateCard, InputRw, TextareaRw, triggerToast } from "$lib";
  import TabBar from "$lib/components/TabBar.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { format, formatDistanceToNow } from "date-fns";
  import { Button, P, Modal, Tooltip, Indicator, Input } from "flowbite-svelte";
  import FileTypeIcon from "$lib/icon/file-type-icon.svelte";
  import { ExclamationCircleSolid, ChevronDownOutline } from "flowbite-svelte-icons";
  import PlusIcon from "$lib/icon/plus.svg";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import {
    triggerSelectCollection,
    triggerShowAI,
    triggerSelectDropdown,
    triggerLoading,
  } from "$lib";

  import { goto } from "$app/navigation";

  const params = new URLSearchParams(window.location.search);
  let collectionId = params.get("collectionId");

  const call = new CR3Call(globalThis.location.origin);

  let documentsReady = $state(false);

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });

  let datas = $state([]);

  let onEdit = $state(false);

  let popupModal = $state(false);

  let datasChecked = $state([]);

  let collection = $state({});

  let configs = $state({});

  let tabDatas = $state({
    selected: "Documents",
    data: [
      {
        name: "Documents",
      },
      {
        name: "Information",
      },
    ],
  });

  let files = $state([]);
  let tmpDocumentType = $state("");
  let selectedTypeCollection = $state("Documents (TXT, PDF, DOC, DOCX)");
  let acceptedFileTypes = $state(".txt, .pdf, .doc, .docx");
  let addDocumentsModal = $state(false);
  let dropdownOpen = $state(false);
  let newModelName = $state("");
  let trainingModelModal = $state(false);

  let timer;

  async function getDocuments(page) {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-documents-by-collectionid", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: page,
          page_size: pagination.pageSize,
          collection_id: collectionId,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      console.log("ayam", response);

      // let response = {
      //   result: {
      //     data: [
      //       {
      //         id: 1,
      //         attributes: {
      //           createdAt: new Date(),
      //           document_title: "Document 1",
      //           document_type: "GENERAL",
      //           status: "PENDING",
      //         },
      //       },
      //       {
      //         id: 2,
      //         attributes: {
      //           createdAt: new Date(),
      //           document_title: "Document 2",
      //           document_type: "CONTRACT",
      //           status: "APPROVED",
      //         },
      //       },
      //     ],
      //     meta: {
      //       pagination: {
      //         page: 1,
      //         pageSize: 10,
      //         pageCount: 1,
      //         total: 2,
      //       },
      //     },
      //   },
      // };

      if (response.result.data.length > 0) {
        configs = {
          search: true,
          headers: "File Uploaded",
          pagination: response.result.meta.pagination,
          actions: [
            {
              id: "delete",
              name: "Delete",
              icon: "trash",
              type: "error",
              onclick: () => {
                if (datasChecked.find((data) => data.checked)) {
                  popupModal = true;
                }
              },
            },
            {
              id: "add_new",
              name: "Add New",
              icon: "plus",
              type: "brand-inverse",
              onclick: () => {
                selectTypeCollection("Documents");
                addDocumentsModal = true;
              },
            },
          ],
        };
        let promises = response.result.data.map(async (document) => {
          let link = await call.getPage("client-dashboard", "document-collection-detail");
          return {
            id: document.id,
            file_name: {
              name: document.attributes.document_title,
              type: "custom",
              size:
                document.attributes.file_size === null ? "Unknown" : document.attributes.file_size,
            },
            type: {
              name: document.attributes.document_type,
              type: "pill",
              color: "brand",
            },
            file_size:
              document.attributes.file_size === null ? "Unknown" : document.attributes.file_size,
            upload_date: formatDistanceToNow(document.attributes.createdAt, { addSuffix: true }),
            uploaded_by:
              document.attributes.upload_by.data === null
                ? "Unknown"
                : document.attributes.upload_by.data.attributes.fullname,
            actions: {
              objects: [
                {
                  // href: "/document/document_collection/detail?documentId=1&collectionId=1",
                  href: `${link}` + "?documentId=" + document.id + "&collectionId=" + collectionId,
                  name: "View",
                },
                // {
                //   // href: "localhost:5173/client-dashboard/document-collection-detail?documentId=1&collectionId=1",
                //   // `${await call.getPage("client-dashboard", "document-collection-detail")}` +
                //   // "?documentId=" +
                //   // document.id +
                //   // "&collectionId=" +
                //   // collectionId,
                //   function: handleTitikTiga,
                //   name: "Titik tiga",
                // },
              ],
            },
          };
        });

        datas = await Promise.all(promises);
        datasChecked = datas.map((data) => {
          return { id: data.id, checked: false };
        });
      }
    }

    documentsReady = true;
  }

  async function getCollection() {
    // search collection using params collectionId here
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await call.call("client-dashboard", "get-collection-by-id", {
          URLQuery: {},
          URLParams: {
            id: collectionId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          console.log(response);
          return;
        }

        collection = {
          id: response.result.data.id,
          name: response.result.data.attributes.collection_name,
          createdAt: formatDistanceToNow(response.result.data.attributes.createdAt, {
            addSuffix: true,
          }),
          description: response.result.data.attributes.description,
        };
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function addDocumentByCollection() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        addDocumentsModal = false;
        triggerLoading(true);
        let response = await call.call("client-dashboard", "create-document", {
          URLQuery: {},
          URLParams: {},
          Params: {
            documents: JSON.stringify(files),
            collection_id: collectionId,
            document_type: tmpDocumentType,
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

        let response2 = await call.call("hexa-ai-rag", "add-document-by-collection-id", {
          URLQuery: {},
          URLParams: {},
          Params: {
            collection_id: collectionId.toString(),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          addDocumentsModal = true;
          console.log(response2);
          triggerToast("error");
          triggerLoading(false);
          return;
        }

        console.log(response2);

        if (tmpDocumentType === "Handwriting") {
          for (let i = 0; i < files.length; i++) {
            let response3 = await call.call("hexa-ai-rag", "trocr-predict", {
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
              triggerToast("error");
              triggerLoading(false);
              console.log(response3);
              return;
            }

            console.log("trocr-predict", response3);

            if (response3?.result?.text !== null) {
              const response4 = await call.call("client-dashboard", "update-document", {
                URLQuery: {},
                URLParams: {
                  id: response.result.data[i].id.toString(),
                },
                Params: {
                  detected_text: response3.result.text,
                  result_trocr: JSON.stringify(response3.result),
                  edited: false,
                },
                Headers: {
                  Authorization: "Bearer " + token,
                },
              });

              if (response3.status !== 200) {
                addDocumentsModal = true;
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
        addDocumentsModal = true;
        console.log(error);
        triggerToast("error");
        triggerLoading(false);
        return;
      }

      files = [];

      addDocumentsModal = false;
      getDocuments(pagination.page);
      triggerLoading(false);

      triggerToast("created");
    }
  }

  async function deleteDocument() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      let ids = datasChecked.filter((data) => data.checked).map((data) => data.id);
      try {
        for (const id of ids) {
          let response = await call.call("client-dashboard", "delete-document", {
            URLQuery: {},
            URLParams: {
              id: id.toString(),
            },
            Headers: {
              Authorization: `Bearer ${token}`,
            },
          });

          if (response.status !== 200) {
            console.log(response);
            return;
          }

          if (response.status === 200) {
            let response2 = await call.call("hexa-ai-rag", "delete-document-by-id", {
              URLQuery: {},
              URLParams: {},
              Params: {
                document_id: id.toString(),
              },
              Headers: {
                Authorization: `Bearer ${token}`,
              },
            });

            if (response2.status !== 200) {
              console.log(response2);
              return;
            }
          } else {
            triggerToast("error");
            console.log(response);
            return;
          }
        }

        triggerToast("deleted");
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
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
    inputElement.multiple = true;

    inputElement.addEventListener("change", async (event) => {
      const uploadedFiles = event.target.files;

      if (uploadedFiles.length > 0) {
        for (const file of uploadedFiles) {
          const base64Data = await readFileAsBase64(file);
          files.push({
            name: file.name,
            size: formatBytes(file.size),
            base64: base64Data,
          });
        }
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
      selectedTypeCollection = "Handwriting (PDF, JPG, JPEG, PNG, WEBP, IMG)";
      acceptedFileTypes = ".jpg, .jpeg, .png";
    }

    files = [];

    dropdownOpen = false;
  };

  let buttonDisabled = $state(false);

  $effect(() => {
    buttonDisabled = files.length === 0;
  });

  const changeTabbar = (name) => {
    tabDatas.selected = name;
  };

  async function next() {
    getDocuments(pagination.page + 1);
  }

  async function previous() {
    getDocuments(pagination.page - 1);
  }

  async function handleCheckAll(allChecked) {
    datasChecked = datasChecked.map((data) => {
      return { id: data.id, checked: allChecked ? false : true };
    });
  }

  async function handleCheck(id) {
    let datasCheckedCopy = JSON.parse(JSON.stringify(datasChecked));
    datasCheckedCopy = datasCheckedCopy.map((data) => {
      if (data.id === id) {
        return { id: data.id, checked: !data.checked };
      }
      return data;
    });
    datasChecked = datasCheckedCopy;
  }

  async function handleTitikTiga(id) {
    console.log("Titik tiga ", id);
  }

  async function selectCollection(collection) {
    triggerShowAI(true);
    triggerSelectDropdown("collections");
    triggerSelectCollection(collection);
  }

  onMount(async () => {
    getDocuments(pagination.page);
    collectionId ? getCollection() : null;
  });

  onDestroy(() => {
    if (timer) clearTimeout(timer);
  });

  async function startTraining() {
    if (collectionId) {
      try {
        trainingModelModal = false;
        triggerLoading(true);
        let token = sessionStorage.getItem("token");
        const response = await call.call("hexa-ai-rag", "trocr-train", {
          URLQuery: {},
          URLParams: {},
          Params: {
            new_model_name: newModelName,
            collection_id: collectionId,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log(response);
        if (response === 200) {
          sessionStorage.setItem("model_trocr", `trocr:${newModelName}`);
          triggerToast("created");
          triggerLoading(false);
        } else {
          trainingModelModal = true;
          triggerToast("error");
        }
      } catch (error) {
        trainingModelModal = true;
        console.log(error);
      }
    }
  }
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <div>
            <Typography
              preset="Display sm/Semibold"
              color="primary"
              text={collection.name ? collection.name : $_("loading")}
            />
            <P
              class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
            >
              {$_("createOn")}{` ${collection.createdAt}`}
            </P>
          </div>
          <div>
            <div class="flex gap-2">
              {#if collectionId === "29"}
                <Button
                  disabled={datas.length === 0}
                  onclick={() => {
                    trainingModelModal = true;
                  }}
                  icon={PlusIcon}
                  ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" />
                  {$_("trainingModel")}</Button
                >
              {/if}
              <Button
                disabled={datas.length === 0}
                onclick={() => {
                  selectCollection({
                    id: collectionId,
                    name: collection.name,
                    description: collection.description,
                    totalDocuments: datas.length,
                  });
                }}
                icon={PlusIcon}
                ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" />
                {$_("startReading")}</Button
              >
            </div>
          </div>
        </div>

        <!-- Submenu -->
        <!-- <TabBar {tabDatas} {changeTabbar} /> -->
        <TabBar {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

        {#if tabDatas.selected === "Documents"}
          <!-- Person Table Data -->
          {#if datas.length > 0}
            <DataTableListCollection
              bind:datas
              bind:datasChecked
              bind:configs
              bind:next
              bind:previous
              bind:handleCheckAll
              bind:handleCheck
            />
          {:else}
            <div class="pt-8">
              <StateCard
                illustration={IllustrationEmpty}
                description={$_("youDontHaveAnyDataYet")}
                actionText={$_("addNewDocument")}
                funcButton={() => (addDocumentsModal = true)}
              />
            </div>
          {/if}
        {:else if tabDatas.selected === "Information"}
          <div class="flex flex-col gap-4 pt-6">
            <!-- general information -->
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {$_("generalInformation")}
                </P>
              </div>
              <div class="flex flex-col grow gap-6">
                <div class="flex gap-5">
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("collectionName")}
                    </P>
                    <InputRw {onEdit} bind:value={collection.name} />
                  </div>
                </div>
                <div class="grid grid-cols-1">
                  <div class="flex flex-col gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >{$_("description")}</P
                    >
                    <TextareaRw {onEdit} bind:value={collection.description} />
                  </div>
                </div>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
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
      <P class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
        >{$_("uploadDocuments")}</P
      >
      {#if files.length === 0}
        <Button
          onclick={() => {
            handleFileUploadWrapper();
          }}
          color="light"
          class={`w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-3 py-4 flex flex-col justify-start items-center`}
        >
          <div
            class="rounded-xl shadow border border-primary-light dark:border-primary-dark justify-center items-center inline-flex pointer-events-none"
          >
            <img src={upload} alt="mail.svg" class="p-4 dark:invert pointer-events-none" />
          </div>
          <div class="flex space-x-1">
            <P class="text-hc-brand-700 pointer-events-none">{$_("clickToUpload")}</P>
            <P class="pointer-events-none">
              <!-- {$_("orDragAndDrop")} -->
            </P>
          </div>
          <!-- <P class="pointer-events-none">PDF, DOC, DOCX (Max. 15MB)</P> -->
        </Button>
      {:else}
        <div class="flex items-center justify-center relative">
          <div class="relative w-min">
            <FileTypeIcon />
            <Indicator
              color="purple"
              border
              placement="center-right"
              class="!w-8 !h-8 text-white text-xs font-semibold p-4">+{files.length}</Indicator
            >
          </div>
          <div class="absolute bottom-0 right-0">
            <Button on:click={handleFileUploadWrapper} pill={true} id="placement-1" class="!p-2"
              ><img class="w-4 h-4" src={PlusIcon} alt="plus-icon" /></Button
            >
            <Tooltip type="light" triggeredBy="#placement-1" placement="left">
              {$_("addMoreFiles")}
            </Tooltip>
          </div>
        </div>
      {/if}
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={addDocumentByCollection} disabled={buttonDisabled} class="w-full"
      >{$_("addNewDocument")}</Button
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
    <Button onclick={deleteDocument} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>

<Modal title={$_("trainingModel")} bind:open={trainingModelModal} size="xs">
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
</Modal>
