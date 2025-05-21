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
  import { P, Heading, Button, Modal, Helper } from "flowbite-svelte";
  import {
    ToastContainer,
    triggerToast,
    InputRw,
    TextareaRw,
    Typography,
    triggerLoading,
    check,
    truncateText,
    PageFileViewer,
  } from "$lib";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { _, locale } from "svelte-i18n";
  import { id as localeId, enUS as localeEn } from "date-fns/locale";
  import { formatDistanceToNow } from "date-fns";
  import { goto } from "$app/navigation";
  import PdfjsViewer from "$lib/components/PdfjsViewer.svelte";
  import { page } from "$app/stores";
  import AudioPlayer from "$lib/components/AudioPlayer.svelte";
  import MultiSelect from "svelte-multiselect";

  const isFrontend = process.env.NODE_ENV === "development";

  let documentId, collectionId;

  import { call } from "$lib";

  let onEdit = $state(false);
  let popupModal = $state(false);
  let pdfReady = $state(false);
  let attributeNotUnique = $state(false);
  let nameIsAvailable = $state(true);
  let fileReady = $state(false);
  let fileUrl = $state("null");
  let uploadBy = $state("");
  let data = $state({});
  let tmpData = $state({});
  let alltags = $state([]);
  let tags = $state([]);
  let contentFile = $state(null);
  let base64 = $state(null);
  let innerWidth = $state(0);

  const base64toBlob = (base64Data, fileType) => {
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
  };

  const getLabels = async () => {
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
  };

  // const truncateText = (text, maxWords, maxFirstWordLength = 14) => {
  //   const words = text.split(" ");

  //   // Check if the first word is too long, or if the second word is long
  //   if (words[0] && words[0].length > maxFirstWordLength) {
  //     return words[0].slice(0, maxFirstWordLength) + "...";
  //   } else if (words[1] && words[1].length > maxFirstWordLength) {
  //     return words[0] + " " + words[1].slice(0, maxFirstWordLength) + "...";
  //   }

  //   // Limit to maxWords and append "..." if necessary
  //   if (words.length > maxWords) {
  //     return words.slice(0, maxWords).join(" ") + " ...";
  //   }

  //   return words.join(" ");
  // };

  const downloadFile = (content, fileName, ext) => {
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
  };

  const getDocument = async () => {
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

        if (response.status !== 200) {
          console.log(response);
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
          goto(await $call.getPage("client-dashboard", "document-library"));
          return;
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

        data = {
          id: response.result.id,
          createdAt: documentData.createdAt,
          document_title: documentData.name,
          description: documentData.meta?.description || "",
          category: documentData.meta.category,
          document_type: documentData.meta.document_type,
          upload_by: documentData.meta.upload_by,
          ext: response.result.ext,
          // document_type: "Document Audio",
          // ext: ".mp3",
          detectedText: documentData.meta.detected_text ? documentData.meta.detected_text : "",
          slices: documentData.meta.result_trocr ? documentData.meta.result_trocr.slices : "",
          result_trocr: documentData.meta.result_trocr,
        };

        tmpData = {
          id: response.result.id,
          createdAt: documentData.createdAt,
          document_title: documentData.name,
          description: documentData.meta?.description || "",
          category: documentData.meta.category,
          document_type: documentData.meta.document_type,
          upload_by: documentData.meta.upload_by,
          ext: response.result.ext,
          // document_type: "Document Audio",
          // ext: ".mp3",
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
        console.error("Error during document get:", error);
        triggerToast($_("failedGetDocument"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  };

  const deleteDocument = async () => {
    let token = localStorage.getItem("token");
    if (token === null) {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
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

        triggerToast($_("successDeleteDocument"), true);
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
        console.error("Error during document deletion:", error);
        triggerToast($_("failedDeleteDocument"), true, "error");
      }
    }
  };

  const updateDocument = async () => {
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
          console.log(response);
          triggerLoading(false);
          onEdit = true;
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

        // console.log(response);
        attributeNotUnique = false;
        triggerToast("successUpdateDocument", true);
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
          console.log(response);
          triggerLoading(false);
          onEdit = true;
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

        triggerLoading(false);
        // console.log(response);
        triggerToast("successUpdateDocument", true);
        getDocument();
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
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

      // console.log("cek resp", resp);

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

  const handleInputName = async (event) => {
    const name = event.target.value;
    const status = await checkLibraryName(name.trim());

    if (status === true) {
      nameIsAvailable = true;
    } else if (status === false) {
      nameIsAvailable = false;
    }
  };

  onMount(() => {
    if (check("documentAi.documentLibrary:view") === false) {
      window.history.back();
      return;
    }
    const unsubscribe = page.subscribe(async ($page) => {
      documentId = $page.url.searchParams.get("documentId");
      if (documentId === null || documentId === undefined) {
        triggerToast("error");
        goto(
          isFrontend ? "/documents" : await $call.getPage("client-dashboard", "document-library"),
        );
        return;
      }
      base64 = null;
      pdfReady = false;
      await getDocument();
      getLabels();
    });
    const resizeHandler = () => {
      innerWidth = globalThis.innerWidth;
    };
    globalThis.addEventListener("resize", resizeHandler);
    resizeHandler();
    return () => {
      unsubscribe();
      globalThis.removeEventListener("resize", resizeHandler);
    };
  });

  $effect(() => {
    if (innerWidth < 1024) {
      tmpData.document_title = truncateText(tmpData.document_title || "", 2, 2);
    } else {
      tmpData.document_title = data.document_title;
    }
  });
</script>

<!-- <div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex items-center justify-between">
        <div class="flex flex-col grow">
          <div class="gap-2">
            <Heading
              tag="h3"
              class="fluid-xl md:fluid-xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {tmpData.document_title}
              &nbsp;
              <button
                onclick={() => {
                  navigator.clipboard
                    .writeText(tmpData.document_title)
                    .then(() => {
                      console.log("Text copied to clipboard");
                      triggerToast("copyToClipboard");
                    })
                    .catch((err) => {
                      console.error("Failed to copy: ", err);
                    });
                }}
                class="p-2 border-2 rounded-lg border-secondary-light dark:border-secondary-dark"
              >
                <Copy01Icon
                  svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
                />
              </button>
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
                attributeNotUnique = false;
                nameIsAvailable = true;
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button
              onclick={() => {
                updateDocument();
              }}
              disabled={data.document_title === "" || !nameIsAvailable}
              class="px-3 focus-within:!ring-1">{$_("save")}</Button
            >
          {:else}
            <Button
              onclick={async function () {
                triggerLoading(false);
                let collectionId = $page.url.searchParams.get("collectionId");
                let libraryId = $page.url.searchParams.get("libraryId");
                goto(
                  (isFrontend
                    ? "/documents/libraries/collections/documents"
                    : await $call.getPage("client-dashboard", "document-library-documents")) +
                    `?collectionId=${collectionId}&libraryId=${libraryId}`,
                );
                return;
              }}
              color="light"
              class="hidden md:block px-3 focus-within:!ring-1">{$_("back")}</Button
            >
            <Button
              disabled={!check("documentAi.documentLibrary.document:delete")}
              onclick={() => (popupModal = true)}
              outline
              color="red"
              class="px-3 focus-within:!ring-1">{$_("delete")}</Button
            >
            <Button
              disabled={!check("documentAi.documentLibrary.document:update")}
              onclick={() => {
                onEdit = true;
              }}
              class="px-3 focus-within:!ring-1">{$_("edit")}</Button
            >
          {/if}
        </div>
      </div>

      <div class="flex flex-col gap-4 pt-4 md:pt-6 md:gap-6">
        {#if fileReady}
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
                <InputRw {onEdit} oninput={handleInputName} bind:value={data.document_title} />
                {#if !nameIsAvailable}
                  <Helper color="red">
                    <span class="font-medium">{$_("nameAlreadyTaken")}</span>
                  </Helper>
                {/if}
              </div>
              <div class="flex flex-col gap-1.5">
                <Typography
                  preset="Text sm/medium"
                  color="secondary"
                  text={`${$_("description")}`}
                />
                <InputRw {onEdit} bind:value={data.description} />
              </div>
              <div class="flex flex-col gap-1.5">
                <Typography preset="Text sm/medium" color="secondary" text={`${$_("labels")}`} />
                <MultiSelect
                  --sms-min-height="40px"
                  options={alltags || []}
                  disabled={!onEdit}
                  allowUserOptions="append"
                  bind:selected={tags}
                  autoScroll={true}
                  noMatchingOptionsMsg=""
                  createOptionMsg={null}
                  outerDivClass="bg-primary-light dark:bg-secondary-dark border-primary-light dark:border-primary-dark text-primary-light dark:text-primary-dark rounded-lg"
                  liOptionClass="bg-primary-light dark:bg-secondary-dark hover:bg-secondary-light hover:dark:bg-primary-dark text-primary-light dark:text-primary-dark"
                  liSelectedClass="bg-tertiary-light dark:bg-primary-dark"
                />
              </div>
              <div class="flex flex-col gap-1.5">
                <Typography preset="Text sm/medium" color="secondary" text={$_("uploadBy")} />
                <InputRw onEdit={false} bind:value={uploadBy} />
              </div>
            </div>
          </div>
          <Divider />
          {#if data.document_type === "Handwriting" || data.document_type === "Document Images" || data.document_type === "Document Audio"}
            <div class="flex flex-col md:flex-row gap-4 md:gap-6">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {#if data.document_type === "Handwriting"}
                    {$_("handwritingInformation")}
                  {:else if data.document_type === "Document Images"}
                    {$_("documentInformation")}
                  {:else if data.document_type === "Document Audio"}
                    {$_("audioInformation")}
                  {/if}
                </P>
              </div>
              <div class="flex flex-col grow gap-4 md:gap-6">
                <div class="grid grid-cols-1">
                  <div class="flex flex-col gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {#if data.document_type === "Handwriting"}
                        {$_("handwritingText")}
                      {:else if data.document_type === "Document Images"}
                        {$_("documentText")}
                      {:else if data.document_type === "Document Audio"}
                        {$_("audioText")}
                      {/if}</P
                    >
                    <TextareaRw {onEdit} bind:value={data.detectedText} />
                  </div>
                </div>
              </div>
            </div>
            <Divider />
          {/if}
          <div class="flex flex-col md:flex-row gap-6">
            <div class="w-60">
              <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                {#if data.ext === ".mp3" || data.ext === ".wav" || data.ext === ".ogg"}
                  {$_("audioPreview")}
                {:else}
                  {$_("viewer")}
                {/if}
              </P>
            </div>
            <div
              class={`flex flex-col relative ${data.document_type === "Document Audio" ? "items-start ps-8" : "items-center"} grow py-4 md:py-12 rounded-lg bg-secondary-light dark:bg-secondary-dark`}
            >
              {#if data.ext === ".pdf"}
                {#if pdfReady}
                  <PdfjsViewer pdfFile={base64} />
                {/if}
              {:else if data.ext === ".jpg" || data.ext === ".png" || data.ext === ".jpeg" || data.ext === ".webp"}
                <div class="flex justify-center">
                  <img class="w-full h-screen" src={`${fileUrl}`} alt="imageViewer" />
                </div>
              {:else if data.ext === ".mp3" || data.ext === ".wav" || data.ext === ".ogg"}
                {#if fileUrl}
                  <AudioPlayer src={fileUrl} />
                {/if}
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
</Modal> -->

<PageFileViewer
  {data}
  {tags}
  {onEdit}
  {base64}
  {tmpData}
  {alltags}
  {fileUrl}
  {pdfReady}
  {uploadBy}
  {fileReady}
  {popupModal}
  {updateDocument}
  {deleteDocument}
  {nameIsAvailable}
  {handleInputName}
  {attributeNotUnique}
/>
