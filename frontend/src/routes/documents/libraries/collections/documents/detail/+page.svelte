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
  import {
    call,
    check,
    triggerToast,
    truncateText,
    base64toBlob,
    triggerLoading,
    PageFileViewer,
    checkLibraryName,
  } from "$lib";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  const isFrontend = process.env.NODE_ENV === "development";
  const pageDataType = "document";

  let debounceTimeout;

  let onEdit = $state(false);
  let popupModal = $state(false);
  let pdfReady = $state(false);
  let attributeNotUnique = $state(false);
  let nameIsAvailable = $state(true);
  let fileReady = $state(false);
  let documentId = $state(null);
  let collectionId = $state(null);
  let libraryId = $state(null);
  let fileUrl = $state("null");
  let uploadBy = $state("");
  let documentTitle = $state("");
  let data = $state({});
  let defaultData = $state({});
  let tags = $state([]);
  let contentFile = $state(null);
  let base64 = $state(null);
  let innerWidth = $state(0);

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

        defaultData = { ...data };

        tags =
          documentData.meta.category === undefined ? [] : documentData.meta.category.split(",");

        console.log("tags get", tags);

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
        triggerToast($_("successUpdateDocument"), true);
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
        triggerToast($_("successUpdateDocument"), true);
        getDocument();
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  };

  const handleInputName = async (event, defaultName) => {
    const name = event.target.value.trim();

    // Reset the debounce timer
    clearTimeout(debounceTimeout);

    // If the name matches the default name, set availability and return early
    if (name === defaultName) {
      nameIsAvailable = true;
      return;
    }

    // Set a debounce delay of 0.5 seconds before calling checkLibraryName
    debounceTimeout = setTimeout(async () => {
      const status = await checkLibraryName($call, name, collectionId, pageDataType);

      // Update the availability status based on the result
      nameIsAvailable = !!status;
    }, 500);
  };

  onMount(() => {
    if (check("documentAi.documentLibrary:view") === false) {
      window.history.back();
      return;
    }

    const unsubscribe = page.subscribe(async ($page) => {
      documentId = $page.url.searchParams.get("documentId");
      collectionId = $page.url.searchParams.get("collectionId");
      libraryId = $page.url.searchParams.get("libraryId");
      if (documentId === null || documentId === undefined) {
        triggerToast($_("error"));
        goto(
          isFrontend ? "/documents" : await $call.getPage("client-dashboard", "document-library"),
        );
        return;
      }
      base64 = null;
      pdfReady = false;
      await getDocument();
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
      documentTitle = truncateText(defaultData.document_title || "", 2, 6);
    } else {
      documentTitle = defaultData.document_title;
    }
  });
</script>

<PageFileViewer
  bind:data
  bind:tags
  bind:documentTitle
  {onEdit}
  {base64}
  {fileUrl}
  {pdfReady}
  {uploadBy}
  {fileReady}
  {libraryId}
  {popupModal}
  {documentId}
  {defaultData}
  {collectionId}
  {pageDataType}
  {nameIsAvailable}
  {attributeNotUnique}
  updateData={updateDocument}
  deleteData={deleteDocument}
  {handleInputName}
/>
