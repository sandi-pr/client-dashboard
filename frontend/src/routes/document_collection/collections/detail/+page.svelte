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
  import { P, Heading, Input, Textarea, Button, Select, Modal } from "flowbite-svelte";
  import {
    ToastContainer,
    triggerToast,
    InputRw,
    TextareaRw,
    Typography,
    triggerLoading,
  } from "$lib";
  import { CheckCircleSolid, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import DownloadIcon from "$lib/icon/file-download-03.svg";
  import Divider from "$lib/components/Divider.svelte";
  import { onMount, onDestroy } from "svelte";
  import { _ } from "svelte-i18n";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { format, formatDistanceToNow } from "date-fns";
  import { goto } from "$app/navigation";
  import { createEventDispatcher } from "svelte";
  import PdfViewer from "svelte-pdf";
  import { generateUniqueId } from "$lib";
  import { page } from "$app/stores";
  import LoadingContainer from "$lib/components/LoadingContainer.svelte";

  let documentId, collectionId;

  const call = new CR3Call(globalThis.location.origin);

  // const params = new URLSearchParams(window.location.search);
  // let documentId = params.get("documentId");
  // let collectionId = params.get("collectionId");
  let documentIdOld = $state("");
  let collectionIdOld = $state("");
  const dispatch = createEventDispatcher();
  let timer;

  let onEdit = $state(false);

  let popupModal = $state(false);

  let contentFile = $state(null);
  let base64 = $state(null);
  let fileUrl = $state("null");
  let pdfReady = $state(false);

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

  let fileReady = $state(false);

  let data = $state({});

  let categories = $state([
    {
      value: "GENERAL",
      name: "General",
    },
    {
      value: "CONTRACT",
      name: "Contract",
    },
  ]);

  let personClient = $state([]);

  async function getDocument() {
    triggerLoading(true);
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await call.call("client-dashboard", "get-document", {
          URLQuery: {},
          URLParams: {
            id: documentId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          // triggerToast("error");
          console.log(response);
          return;
        }

        if (response.result.document.data.length > 0) {
          let documentData = response.result.document.data[0];

          base64 = response.result.base64;
          pdfReady = true;

          contentFile = atob(base64);

          let fileType = "";
          if (data.ext === ".pdf") {
            fileType = "application/pdf";
          } else if (data.ext === ".docx") {
            fileType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
          } else if (data.ext === ".doc") {
            fileType = "application/msword";
          } else {
            fileType = "text/plain";
          }

          const blob = base64toBlob(base64, fileType);
          fileUrl = URL.createObjectURL(blob);
          console.log(fileUrl);

          // let response2;
          // if (documentData.attributes.document_type === "Handwriting") {
          //   if (
          //     documentData.attributes.detected_text === null ||
          //     documentData.attributes.result_trocr === null
          //   ) {
          //     response2 = await call.call("hexa-ai-rag", "trocr-predict", {
          //       URLQuery: {},
          //       URLParams: {},
          //       Params: {
          //         image: base64,
          //       },
          //       Headers: {
          //         Authorization: "Bearer " + token,
          //       },
          //     });

          //     if (response2.status !== 200) {
          //       // triggerToast("error");
          //       console.log(response);
          //       return;
          //     }

          //     console.log("trocr-predict", response2);

          //     if (response2?.result?.text !== null) {
          //       const response3 = await call.call("client-dashboard", "update-document", {
          //         URLQuery: {},
          //         URLParams: {
          //           id: documentId.toString(),
          //         },
          //         Params: {
          //           detected_text: response2.result.text,
          //           result_trocr: JSON.stringify(response2.result),
          //           edited: false,
          //         },
          //         Headers: {
          //           Authorization: "Bearer " + token,
          //         },
          //       });

          //       if (response3.status !== 200) {
          //         // triggerToast("error");
          //         console.log(response);
          //         return;
          //       }

          //       console.log("ayam", response3);
          //     }
          //   }
          // }

          data = {
            id: documentData.id,
            createdAt: formatDistanceToNow(documentData.attributes.createdAt, { addSuffix: true }),
            document_title: documentData.attributes.document_title,
            description: documentData.attributes.description,
            document_type: documentData.attributes.document_type,
            ext: documentData.attributes.document.data.attributes.ext,
            detectedText: documentData.attributes.detected_text
              ? documentData.attributes.detected_text
              : "",
            slices: documentData.attributes.result_trocr
              ? documentData.attributes.result_trocr.slices
              : "",
            result_trocr: documentData.attributes.result_trocr,
            // description: documentData.attributes.description,
            // document_type: documentData.attributes.document_type,
            // person_client: documentData.attributes.person.data?.id,
            // contract_start:
            //   documentData.attributes.contract_start !== null
            //     ? format(
            //         new Date(documentData.attributes.contract_start),
            //         "yyyy-MM-dd",
            //       )
            //     : "",
            // contract_end:
            //   documentData.attributes.contract_end !== null
            //     ? format(
            //         new Date(documentData.attributes.contract_end),
            //         "yyyy-MM-dd",
            //       )
            //     : "",
          };

          fileReady = true;
          triggerLoading(false);
        }
      } catch (error) {
        // triggerToast("error");
        console.log(error);
      }
    } else {
      // triggerToast("error");
    }
  }

  async function deleteDocument() {
    let token = sessionStorage.getItem("token");
    if (documentId !== "" && documentId !== null && documentId !== undefined) {
      try {
        let response = await call.call("client-dashboard", "delete-document", {
          URLQuery: {},
          URLParams: {
            id: documentId.toString(),
          },
          Headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status !== 200) {
          triggerToast("error");
          console.log(response);
          return;
        }

        if (response.status === 200) {
          let response2 = await call.call("hexa-ai-rag", "delete-document-by-id", {
            URLQuery: {},
            URLParams: {},
            Params: {
              document_id: documentId.toString(),
            },
            Headers: {
              Authorization: `Bearer ${token}`,
            },
          });

          if (response2.status !== 200) {
            console.log(response2);
            return;
          }

          console.log(response2);

          triggerToast("deleted");
          // setTimeout(async () => {
          goto(
            (await call.getPage("client-dashboard", "document-collection-list-document")) +
              "?collectionId=" +
              collectionId,
          );
          // }, 3000);
          return;
        } else {
          triggerToast("error");
          console.log(response);
          return;
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function getPersons() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let page = 1;
      while (true) {
        let response = await call.call("client-dashboard", "get-persons", {
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
          return;
        }

        if (response.result.data.length > 0) {
          let promise = response.result.data.map((person) => {
            return {
              value: person.id,
              name:
                person.attributes.fullname !== null
                  ? person.attributes.fullname
                  : `<${person.attributes.email}>`,
            };
          });

          let data = await Promise.all(promise);
          personClient = personClient.concat(data);
          page++;
        } else {
          break;
        }
      }
    }
  }

  function downloadFile(content, fileName, ext) {
    let fileType = "";
    if (ext === ".pdf") {
      fileType = "application/pdf";
    } else if (ext === ".docx") {
      fileType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
    } else if (ext === ".doc") {
      fileType = "application/msword";
    } else {
      fileType = "text/plain";
    }
    const a = document.createElement("a");
    a.href = URL.createObjectURL(base64toBlob(base64, fileType));
    a.download = fileName;
    a.click();
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

  function getDocumentId() {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get("documentId");
  }

  function getCollectionId() {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get("collectionId");
  }

  // $effect(async () => {
  //   documentId = getDocumentId();
  //   collectionId = getCollectionId();
  //   if (documentId && documentId !== documentIdOld) {
  //     documentIdOld = documentId;
  //     collectionIdOld = collectionId;
  //     getDocument();
  //   }
  //   if (collectionId && collectionId !== collectionIdOld) {
  //     collectionIdOld = collectionId;
  //     documentIdOld = documentId;
  //     getDocument();
  //   }
  // });

  // async function editDetectedText() {
  //   onEdit = false;
  //   let token = sessionStorage.getItem("token");
  //   if (token !== null) {
  //     let response = await call.call("client-dashboard", "update-document", {
  //       URLQuery: {},
  //       URLParams: {
  //         id: documentId.toString(),
  //       },
  //       Params: {
  //         detected_text: data.detectedText,
  //         result_trocr: JSON.stringify(data.result_trocr),
  //         edited: true,
  //       },
  //       Headers: {
  //         Authorization: "Bearer " + token,
  //       },
  //     });
  //     if (response.status !== 200) {
  //       console.log(response);
  //       return;
  //     }

  //     console.log(response);
  //     triggerToast("updated");
  //     getDocument();
  //   }
  // }

  async function updateDocument() {
    onEdit = false;
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response;
      if (data.document_type === "Handwriting") {
        response = await call.call("client-dashboard", "update-document-handwriting", {
          URLQuery: {},
          URLParams: {
            id: documentId.toString(),
          },
          Params: {
            detected_text: data.detectedText,
            result_trocr: JSON.stringify(data.result_trocr),
            edited: true,
            document_title: data.document_title,
            description: data.description,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          console.log(response);
          return;
        }

        console.log(response);
        triggerToast("updated");
        getDocument();
      } else {
        response = await call.call("client-dashboard", "update-document", {
          URLQuery: {},
          URLParams: {
            id: documentId.toString(),
          },
          Params: {
            document_title: data.document_title,
            description: data.description,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          console.log(response);
          return;
        }

        console.log(response);
        triggerToast("updated");
        getDocument();
      }
    }
  }

  onMount(() => {
    const unsubscribe = page.subscribe(($page) => {
      documentId = $page.url.searchParams.get("documentId");
      collectionId = $page.url.searchParams.get("collectionId");
      base64 = null;
      pdfReady = false;
      getDocument();
    });
    return () => unsubscribe();
  });

  $effect(async () => {
    const unsubscribe = page.subscribe(($page) => {
      documentId = $page.url.searchParams.get("documentId");
      collectionId = $page.url.searchParams.get("collectionId");
      base64 = null;
      pdfReady = false;
      getDocument();
    });
    onDestroy(() => unsubscribe());
  });

  // onDestroy(() => {
  // if (timer) {
  //   clearTimeout(timer);
  // }
  // });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex items-center justify-between pb-8">
        <div class="flex flex-col grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
              >{data.document_title}<Button
                outline
                color="red"
                class="hover:!bg-transparent !border-hc-graylight-200 dark:!border-hc-graydark-600 focus-within:!ring-1 !p-2"
              >
                <CopyIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
              </Button></Heading
            >
          </div>
          <P
            class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
            >{$_("createOn")}{` ${data.createdAt}`}</P
          >
        </div>

        <div class="flex gap-2">
          {#if onEdit}
            <Button
              onclick={() => {
                onEdit = false;
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button
              onclick={() => {
                updateDocument();
              }}
              class="px-3 focus-within:!ring-1">{$_("save")}</Button
            >
          {:else}
            <Button
              onclick={async function () {
                goto(
                  (await call.getPage("client-dashboard", "document-collection-list-document")) +
                    `?collectionId=${collectionId}`,
                );
                return;
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("back")}</Button
            >
            <Button
              onclick={() => (popupModal = true)}
              outline
              color="red"
              class="px-3 focus-within:!ring-1">{$_("delete")}</Button
            >
            <!-- {#if data.document_type !== "Handwriting"} -->
            <Button
              onclick={() => {
                onEdit = true;
              }}
              class="px-3 focus-within:!ring-1">{$_("edit")}</Button
            >
            <!-- {/if} -->
          {/if}
        </div>
      </div>

      <!-- Tab Bar -->
      <!-- <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} /> -->

      <!-- Content -->
      <div class="flex flex-col gap-4 pt-6">
        {#if fileReady}
          <!-- general information -->
          <div class="flex gap-6">
            <div class="w-60">
              <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                {$_("generalInformation")}
              </P>
            </div>
            <div
              class="flex flex-col grow gap-6 rounded-lg bg-secondary-light dark:bg-secondary-dark"
            >
              <div class="flex flex-col gap-1.5">
                <Typography
                  preset="Text sm/medium"
                  color="secondary"
                  text={`${$_("documentTitle")}*`}
                />
                <InputRw {onEdit} bind:value={data.document_title} />
              </div>
              <div class="flex flex-col gap-1.5">
                <Typography
                  preset="Text sm/medium"
                  color="secondary"
                  text={`${$_("description")}`}
                />
                <InputRw {onEdit} bind:value={data.description} />
              </div>
            </div>
          </div>
          <Divider />
          {#if data.document_type === "Handwriting"}
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {$_("handwritingInformation")}
                </P>
                <!-- <div class="pt-2">
                  {#if onEdit}
                    <div class="flex gap-2">
                      <button
                        onclick={() => (onEdit = false)}
                        class="text-start text-sm text-tertiary-light dark:text-tertiary-dark"
                        >{$_("cancel")}</button
                      >
                      <button
                        onclick={editDetectedText}
                        class="text-start text-sm text-brand-tertiary-light dark:text-brand-tertiary-dark"
                        >{$_("confirm")}</button
                      >
                    </div>
                  {:else}
                    <Typography
                      preset="Text sm/Regular"
                      color="quarterary"
                      text={$_("wrongDetectedText")}
                    />
                    <button
                      onclick={() => (onEdit = true)}
                      class="text-start text-sm text-brand-tertiary-light dark:text-brand-tertiary-dark"
                      >{$_("correctHere")}</button
                    >
                  {/if}
                </div> -->
              </div>
              <div class="flex flex-col grow gap-6">
                <div class="grid grid-cols-1">
                  <div class="flex flex-col gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >{$_("Handwriting text")}</P
                    >
                    <TextareaRw {onEdit} bind:value={data.detectedText} />
                    <!-- <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >{$_("Result slice image")}</P
                    >
                    {#each data.slices as slice}
                      <div class="relative">
                        <img src={slice.image} alt="Dynamic Image" class="w-full h-auto" />
                        <P
                          class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                          >{slice.text}</P
                        >
                      </div>
                    {/each} -->
                  </div>
                </div>
              </div>
            </div>
            <Divider />
          {/if}
          <div class="flex gap-6">
            <div class="w-60">
              <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                {$_("viewer")}
              </P>
            </div>
            <div
              class="flex flex-col grow gap-6 rounded-lg bg-secondary-light dark:bg-secondary-dark"
            >
              {#if data.ext === ".pdf"}
                <!-- <iframe class="w-full h-screen" src={`${fileUrl}`}></iframe> -->
                <div class="flex justify-center">
                  {#if pdfReady}
                    <PdfViewer
                      showBorder={false}
                      showButtons={["navigation", "zoom", "rotate", "autoflip", "pageInfo"]}
                      scale={1}
                      data={atob(base64)}
                    />
                  {/if}
                </div>
              {:else if data.ext === ".jpg" || data.ext === ".png" || data.ext === ".jpeg" || data.ext === ".webp"}
                <div class="flex justify-center">
                  <img class="w-full h-screen" src={`${fileUrl}`} alt="imageViewer" />
                </div>
              {:else}
                {#each contentFile.split("\n") as paragraphs, i}
                  <p class="dark:text-white">{paragraphs}</p>
                {/each}
              {/if}
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
  <ToastContainer />
</div>

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

<!-- <LoadingContainer /> -->
