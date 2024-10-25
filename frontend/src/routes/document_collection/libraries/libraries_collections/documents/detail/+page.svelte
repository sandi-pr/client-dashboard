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
  import { _, locale } from "svelte-i18n";
  import { id as localeId, enUS as localeEn } from "date-fns/locale";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { format, formatDistanceToNow } from "date-fns";
  import { goto } from "$app/navigation";
  import { createEventDispatcher } from "svelte";
  import PdfViewer from "svelte-pdf";
  import PdfjsViewer from "$lib/components/PdfjsViewer.svelte";
  import { generateUniqueId } from "$lib";
  import { page } from "$app/stores";
  import LoadingContainer from "$lib/components/LoadingContainer.svelte";
  import XCloseIcon from "$lib/icon/x-close.svelte";

  let documentId, collectionId;

  import { call } from "$lib";

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
  let tmpData = $state({});

  // let categoriesTab = $state([
  //   {
  //     value: "GENERAL",
  //     name: "General",
  //   },
  //   {
  //     value: "CONTRACT",
  //     name: "Contract",
  //   },
  // ]);

  let personClient = $state([]);

  async function getDocument() {
    triggerLoading(true);
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "get-library-document", {
          URLQuery: {},
          URLParams: {
            id: documentId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("document", response);

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          console.log(response);
          goto(await $call.getPage("client-dashboard", "document-library"));
        }

        let documentData = response.result.data;

        collectionId = response.result.data.parent_id;

        base64 = response.result.base64;
        pdfReady = true;

        contentFile = atob(base64);

        let fileType = "";
        if (response.result.ext === ".pdf") {
          fileType = "application/pdf";
        } else if (response.result.ext === ".docx") {
          fileType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
        } else if (response.result.ext === ".doc") {
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
        //     response2 = await $call.call("hexa-ai-rag", "trocr-predict", {
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
        //       const response3 = await $call.call("client-dashboard", "update-document", {
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

        // categories = await getCategories(collectionId);

        data = {
          id: response.result.id,
          createdAt: documentData.createdAt,
          document_title: documentData.name,
          description: documentData.meta.description,
          category: documentData.meta.category,
          document_type: documentData.meta.document_type,
          ext: response.result.ext,
          detectedText: documentData.meta.detected_text ? documentData.meta.detected_text : "",
          slices: documentData.meta.result_trocr ? documentData.meta.result_trocr.slices : "",
          result_trocr: documentData.meta.result_trocr,
        };

        tmpData = {
          id: response.result.id,
          createdAt: documentData.createdAt,
          document_title: documentData.name,
          description: documentData.meta.description,
          category: documentData.meta.category,
          document_type: documentData.meta.document_type,
          ext: response.result.ext,
          detectedText: documentData.meta.detected_text ? documentData.meta.detected_text : "",
          slices: documentData.meta.result_trocr ? documentData.meta.result_trocr.slices : "",
          result_trocr: documentData.meta.result_trocr,
        };

        tags =
          documentData.meta.category === undefined ? [] : documentData.meta.category.split(",");

        fileReady = true;
        triggerLoading(false);
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        console.log(error);
      }
    } else {
      // triggerToast("error");
    }
  }

  async function getCategories(id) {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "get-categories-by-collection-id", {
        URLQuery: {},
        URLParams: {},
        Params: {
          collection_id: id.toString(),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(id, response);
        return;
      }

      return response.result.categories.map((category) => {
        return {
          name: category,
          value: category,
        };
      });
    }
  }

  async function deleteDocument() {
    let token = localStorage.getItem("token");
    if (documentId !== "" && documentId !== null && documentId !== undefined) {
      try {
        // let response2 = await $call.call("hexa-ai-rag", "delete-document-by-id", {
        //   URLQuery: {},
        //   URLParams: {},
        //   Params: {
        //     document_id: documentId.toString(),
        //   },
        //   Headers: {
        //     Authorization: `Bearer ${token}`,
        //   },
        // });

        // if (response2.status !== 200) {
        //   console.log(response2);
        //   return;
        // }

        // console.log(response2);

        // if (response2.status === 200) {
        let response = await $call.call("client-dashboard", "delete-library-document", {
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

        triggerToast("deleted");
        // setTimeout(async () => {
        goto(
          (await $call.getPage("client-dashboard", "document-library-documents")) +
            "?collectionId=" +
            collectionId,
        );
        // }, 3000);
        return;
        // } else {
        //   triggerToast("error");
        //   console.log(response);
        //   return;
        // }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function getPersons() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      let page = 1;
      while (true) {
        let response = await $call.call("client-dashboard", "get-persons", {
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
  //   let token = localStorage.getItem("token");
  //   if (token !== null) {
  //     let response = await $call.call("client-dashboard", "update-document", {
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
    triggerLoading(true);
    onEdit = false;
    let token = localStorage.getItem("token");
    if (token !== null) {
      let response;
      if (data.document_type === "Handwriting") {
        response = await $call.call("client-dashboard", "update-library-document", {
          URLQuery: {},
          URLParams: {
            id: documentId.toString(),
          },
          Params: {
            detected_text: data.detectedText,
            document_name: data.document_title,
            description: data.description,
            category: tags.join(","),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          triggerLoading(false);
          onEdit = true;
          console.log(response);
          return;
        }

        console.log(response);
        triggerToast("updated");
        triggerLoading(false);
        getDocument();
      } else {
        response = await $call.call("client-dashboard", "update-library-document", {
          URLQuery: {},
          URLParams: {
            id: documentId.toString(),
          },
          Params: {
            document_name: data.document_title,
            description: data.description,
            category: tags.join(","),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          triggerLoading(false);
          onEdit = true;
          console.log(response);
          return;
        }

        triggerLoading(false);
        console.log(response);
        triggerToast("updated");
        getDocument();
      }
    }
  }

  let tags = $state([]); // Daftar tag yang sudah dimasukkan pengguna
  let newTag = $state(""); // Menyimpan input dari pengguna
  let inputRef;

  // Fungsi untuk menambahkan tag
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

  onMount(() => {
    const unsubscribe = page.subscribe(($page) => {
      documentId = $page.url.searchParams.get("documentId");
      base64 = null;
      pdfReady = false;
      getDocument();
    });
    return () => unsubscribe();
  });

  // $effect(async () => {
  //   const unsubscribe = page.subscribe(($page) => {
  //     documentId = $page.url.searchParams.get("documentId");
  //     base64 = null;
  //     pdfReady = false;
  //     getDocument();
  //   });
  //   onDestroy(() => unsubscribe());
  // });

  // onDestroy(() => {
  // if (timer) {
  //   clearTimeout(timer);
  // }
  // });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex items-center justify-between">
        <div class="flex flex-col grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="fluid-xl md:fluid-xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
              >{tmpData.document_title}
              <!-- <Button
                outline
                color="red"
                class="hover:!bg-transparent !border-hc-graylight-200 dark:!border-hc-graydark-600 focus-within:!ring-1 !p-2"
              >
              </Button> -->
            </Heading>
          </div>
          <P
            class="fluid-sm md:fluid-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
          >
            {#if fileReady}
              {#if $locale === "id"}
                {$_("createOn")}
                {formatDistanceToNow(data.createdAt, { addSuffix: true, locale: localeId })}
              {:else}
                {$_("createOn")}
                {formatDistanceToNow(data.createdAt, { addSuffix: true, locale: localeEn })}
              {/if}
            {/if}</P
          >
        </div>

        <div class="flex gap-2">
          {#if onEdit}
            <Button
              onclick={() => {
                onEdit = false;
                data = { ...tmpData };
                console.log(tmpData);
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button
              onclick={() => {
                updateDocument();
              }}
              disabled={data.document_title === "" || data.description === ""}
              class="px-3 focus-within:!ring-1">{$_("save")}</Button
            >
          {:else}
            <Button
              onclick={async function () {
                goto(
                  (await $call.getPage("client-dashboard", "document-library-documents")) +
                    `?collectionId=${collectionId}`,
                );
                return;
              }}
              color="light"
              class="hidden md:block px-3 focus-within:!ring-1">{$_("back")}</Button
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
      <div class="flex flex-col gap-4 pt-4 md:pt-6 md:gap-6">
        {#if fileReady}
          <!-- general information -->
          <div class="flex flex-col md:flex-row gap-3 md:gap-6">
            <div class="w-60">
              <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                {$_("generalInformation")}
              </P>
            </div>
            <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
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
                  text={`${$_("description")}*`}
                />
                <InputRw {onEdit} bind:value={data.description} />
              </div>
              <div class="flex flex-col gap-1.5">
                <Typography
                  preset="Text sm/medium"
                  color="secondary"
                  text={`${$_("categories")}`}
                />
                <!-- {#if !onEdit}
                  <InputRw bind:value={data.category} />
                {:else} -->
                <div class="flex">
                  <button
                    class="flex w-full flex-wrap gap-2 items-center border rounded-lg py-2.5 px-2"
                    onclick={onEdit
                      ? () => {
                          focusInput(inputRef);
                        }
                      : () => {}}
                  >
                    {#each tags as tag}
                      <span
                        class="flex items-center space-x-0.5 bg-utility-brand-100-light dark:bg-utility-brand-100-dark text-brand-secondary-light dark:text-brand-secondary-dark rounded-full px-2"
                      >
                        <span class="text-sm pb-0.5">{tag}</span>
                        <!-- svelte-ignore node_invalid_placement_ssr -->

                        {#if onEdit}
                          <button
                            onclick={() => {
                              removeTag(tag);
                            }}
                            class="text-brand-primary-light dark:text-brand-primary-dark"
                            ><XCloseIcon
                              svgClass="w-3 h-3 text-utility-brand-400-light dark:text-utility-brand-400-dark"
                            /></button
                          >
                        {/if}
                      </span>
                    {/each}
                    {#if onEdit}
                      <input
                        bind:this={inputRef}
                        type="text"
                        bind:value={newTag}
                        placeholder={`${tags.length < 1 ? $_("enterNewLabel") : ""}`}
                        onkeydown={handleKeyDown}
                        class={`border-none flex-grow !pl-0 !py-0.5 text-sm focus:ring-transparent ${tags.length > 0 ? "max-w-min" : "w-full"}`}
                      />
                    {/if}
                  </button>
                  {#if onEdit}
                    <button
                      onclick={addTag}
                      class="bg-secondary-light dark:bg-secondary-dark text-primary-light dark:text-primary-dark px-4 py-1 rounded"
                    >
                      Add
                    </button>
                  {/if}
                </div>
                <!-- <MultiSelect items={categories} bind:value={tmpCategories}></MultiSelect> -->
                <!-- <Select
                    disabled={!onEdit}
                    class="mt-2"
                    items={categories}
                    bind:value={data.category}
                    placeholder={$_("selectCategory")}
                  /> -->
                <!-- {#if !addOtherCategory}
                    <p>
                      {$_("addOtherCategory")}
                      <button
                        onclick={() => {
                          addOtherCategory = true;
                        }}
                        class="ms-1 text-brand-secondary-light dark:text-brand-secondary-dark hover:underline"
                      >
                        {$_("clickHere")}
                      </button>
                    </p>
                  {/if}
                  {#if addOtherCategory}
                    <Input
                      placeholder={$_("enterOtherCategory")}
                      bind:value={otherCategory}
                      onkeydown={(e) => {
                        if (
                          e.key === "Enter" &&
                          otherCategory !== "" &&
                          !categories.includes(otherCategory)
                        ) {
                          if (categories.length > 0) {
                            for (let i = 0; i < categories.length; i++) {
                              if (categories[i]?.other === true) {
                                categories = categories.slice(0, i).concat(categories.slice(i + 1));
                                break;
                              }
                            }
                          }

                          categories = [
                            ...categories,
                            { value: otherCategory, name: otherCategory, other: true },
                          ];
                          data.category = otherCategory;
                          otherCategory = "";
                        }
                      }}
                    />
                  {/if} -->
                <!-- {/if} -->
              </div>
            </div>
          </div>
          <Divider />
          {#if data.document_type === "Handwriting"}
            <div class="flex flex-col md:flex-row gap-4 md:gap-6">
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
              <div class="flex flex-col grow gap-4 md:gap-6">
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
          <div class="flex flex-col md:flex-row gap-6">
            <div class="w-60">
              <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                {$_("viewer")}
              </P>
            </div>
            <div
              class="flex flex-col relative items-center grow py-4 md:py-12 rounded-lg bg-secondary-light dark:bg-secondary-dark"
            >
              {#if data.ext === ".pdf"}
                {#if pdfReady}
                  <PdfjsViewer pdfFile={base64} />
                {/if}
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
