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
  import { P, Heading, Input, Textarea, Button, Select, Modal, Helper } from "flowbite-svelte";
  import {
    ToastContainer,
    triggerToast,
    InputRw,
    TextareaRw,
    Typography,
    triggerLoading,
    check,
  } from "$lib";
  import { CheckCircleSolid, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import DownloadIcon from "$lib/icon/file-download-03.svg";
  import Divider from "$lib/components/Divider.svelte";
  import { onMount, onDestroy } from "svelte";
  import { _, locale } from "svelte-i18n";
  import { id as localeId, enUS as localeEn, ta } from "date-fns/locale";
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
  import AudioPlayer from "$lib/components/AudioPlayer.svelte";
  import ExampleAudio from "$lib/components/strauss.mp3";
  import MultiSelect from 'svelte-multiselect';

  let audioFile = $state(null); // menyimpan file audio sebagai base64
  let audioUrl = $state(null); // menyimpan URL blob untuk pemutaran audio

  let documentId, collectionId;

  import { call } from "$lib";

  let onEdit = $state(false);

  let popupModal = $state(false);

  let contentFile = $state(null);
  let base64 = $state(null);
  let fileUrl = $state("null");
  let pdfReady = $state(false);
  let attributeNotUnique = $state(false);
  let nameIsAvailable = $state(true);

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
  let uploadBy = $state("");
  let alltags = $state([]);

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

        // console.log("document", response);

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
  }

  async function getLabels() {
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
  }

  async function getUser() {
    let token = localStorage.getItem("token");
    try {
      let response = await $call.call("client-dashboard", "get-user", {
        URLQuery: {},
        URLParams: { id: data.upload_by.toString() },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
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

      uploadBy = response.result.user.username;
    } catch (error) {
      console.error("Error during user get:", error);
      triggerToast($_("failedGetUser"), true, "error");
    }
  }

  function truncateText(text, maxWords, maxFirstWordLength = 14) {
    console.log("text", text);
    const words = text.split(" ");

    // Check if the first word is too long, or if the second word is long
    if (words[0] && words[0].length > maxFirstWordLength) {
      return words[0].slice(0, maxFirstWordLength) + "...";
    } else if (words[1] && words[1].length > maxFirstWordLength) {
      return words[0] + " " + words[1].slice(0, maxFirstWordLength) + "...";
    }

    // Limit to maxWords and append "..." if necessary
    if (words.length > maxWords) {
      return words.slice(0, maxWords).join(" ") + " ...";
    }

    return words.join(" ");
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

        triggerToast("successDeleteDocument", true);
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

        console.log(response);
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
        console.log(response);
        triggerToast("successUpdateDocument", true);
        getDocument();
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
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

  // Fungsi untuk mengunggah file audio dan mengonversinya menjadi base64
  async function handleAudioUpload(event) {
    const file = event.target.files[0];
    if (file && (file.type.includes("audio") || file.name.endsWith(".mp3"))) {
      // Baca file sebagai base64
      const base64Audio = await convertToBase64(file);
      audioFile = base64Audio;

      // Membuat URL blob untuk diputar di audio player
      audioUrl = URL.createObjectURL(new Blob([file], { type: "audio/mp3" }));

      console.log("audioUrl: ", audioUrl);
    } else {
      alert("Please upload a valid MP3 file.");
    }
  }

  // Fungsi konversi file ke Base64
  function convertToBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => resolve(reader.result);
      reader.onerror = reject;
      reader.readAsDataURL(file);
    });
  }

  const checkLibraryName = async (name) => {
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
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      console.log("cek resp", resp);

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

  const handleInput = async (event) => {
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
        goto(await $call.getPage("client-dashboard", "document-library"));
        return;
      }
      base64 = null;
      pdfReady = false;
      await getDocument();
      getUser();
      getLabels();
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
              >{globalThis.innerWidth < 768 ? truncateText(tmpData.document_title || "", 2) : tmpData.document_title}
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
                console.log(tmpData);
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
              disabled={!check("documentAi.documentLibrary.document:delete")}
              onclick={() => (popupModal = true)}
              outline
              color="red"
              class="px-3 focus-within:!ring-1">{$_("delete")}</Button
            >
            <!-- {#if data.document_type !== "Handwriting"} -->
            <Button
              disabled={!check("documentAi.documentLibrary.document:update")}
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
                <InputRw {onEdit} oninput={handleInput} bind:value={data.document_title} />
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
                />
              </div>
              <div class="flex flex-col gap-1.5">
                <Typography
                  preset="Text sm/medium"
                  color="secondary"
                  text={$_("uploadBy")}
                />
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
</Modal>

<!-- <LoadingContainer /> -->
