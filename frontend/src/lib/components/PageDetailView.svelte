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
  import { call, check, triggerToast, base64toBlob, triggerLoading, PageFileViewer } from "$lib";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  const isFrontend = process.env.NODE_ENV === "development";
  let { pageDataType = "document" } = $props();

  let onEdit = $state(false);
  let attributeNotUnique = $state(false);
  let fileReady = $state(false);
  let fileUrl = $state("null");
  let uploadBy = $state("");
  let data = $state({});
  let defaultData = $state({});
  let dataTables = $state([]);
  let dataTablesSelected = $state("");
  let CSVMetadata = $state({ page: 1, pageSize: 10, pageCount: 2, total: 13 });
  let datatablesMetadata = $state({});
  let tags = $state([]);
  let CSVData = $state([]);
  let documentId = $state(null);
  let collectionId = $state(null);
  let libraryId = $state(null);
  let contentFile = $state(null);
  let base64 = $state(null);

  const getData = async (state = "all") => {
    let token = localStorage.getItem("token");
    if (token !== null) {
      if (state === "all") {
        triggerLoading(true);
      }
      try {
        const urlAction =
          pageDataType === "document" ? "get-library-document" : "get-library-document-csv";
        const urlParam =
          pageDataType === "document"
            ? {}
            : {
                page: CSVMetadata.page,
                page_size: CSVMetadata.pageSize,
                table_name: dataTablesSelected,
              };
        let response = await $call.call("client-dashboard", urlAction, {
          URLQuery: {},
          URLParams: {
            id: documentId.toString(),
          },
          Params: {
            ...urlParam,
          },
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

        // console.log("response", response);

        let documentData = response.result.data;
        dataTables =
          dataTables.length > 0
            ? dataTables
            : response.result?.tables?.map((table, index) => {
                return {
                  id: index + 1,
                  name: table,
                  description: table,
                };
              }) || [];

        dataTablesSelected =
          dataTablesSelected !== ""
            ? dataTablesSelected
            : dataTables.length > 0
              ? dataTables[0].name
              : "";

        datatablesMetadata =
          Object.keys(datatablesMetadata).length > 0
            ? datatablesMetadata
            : {
                page: 1,
                pageSize: 10,
                pageCount: dataTables.length > 0 ? Math.ceil(dataTables.length / 10) : 1,
                total: dataTables.length,
              };

        if (pageDataType !== "Data") {
          base64 = response.result.base64;

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
        }

        data = {
          id: response.result.id,
          createdAt: documentData.createdAt,
          document_title: documentData.name,
          description: documentData.meta?.description || "",
          category: documentData.meta.category,
          document_type: documentData.meta.document_type,
          upload_by: documentData.meta.upload_by,
          ext: response.result.ext,
          detectedText: documentData.meta.detected_text ? documentData.meta.detected_text : "",
          slices: documentData.meta.result_trocr ? documentData.meta.result_trocr.slices : "",
          result_trocr: documentData.meta.result_trocr,
        };

        defaultData = { ...data };

        if (pageDataType === "Data") {
          CSVMetadata.pageCount = response.result.result.total_pages;
          CSVMetadata.total = response.result.result.total_items;

          CSVData = response.result.result.data;
        }

        tags =
          documentData.meta.category === undefined ? [] : documentData.meta.category.split(",");

        fileReady = true;
      } catch (error) {
        console.error("Error during document get:", error);
        triggerToast($_("failedGetDocument"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
    }
    triggerLoading(false);
  };

  const deleteData = async () => {
    let token = localStorage.getItem("token");
    if (token !== null) {
      triggerLoading(true);
      if (documentId !== "" && documentId !== null && documentId !== undefined) {
        try {
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

          triggerToast($_("successDeleteDocument"), true);
          goto(
            (await $call.getPage("client-dashboard", "document-library-documents")) +
              "?collectionId=" +
              collectionId,
          );
          return;
        } catch (error) {
          console.error("Error during document deletion:", error);
          triggerToast($_("failedDeleteDocument"), true, "error");
        }
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
    }
    triggerLoading(false);
  };

  const updateData = async () => {
    let token = localStorage.getItem("token");
    if (token !== null) {
      triggerLoading(true);
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
        getData();
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

        triggerToast($_("successUpdateDocument"), true);
        getData();
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
    }
    triggerLoading(false);
  };

  const tableNextPage = (table) => {
    if (table.page >= table.pageCount) return;
    table.page += 1;
    getData("partial");
  };

  const tablePreviousPage = (table) => {
    if (table.page > 1) {
      table.page -= 1;
      getData("partial");
    }
  };

  // Use the generic table pagination for CSV and Datatables
  const CSVNext = () => tableNextPage(CSVMetadata);
  const CSVPrevious = () => tablePreviousPage(CSVMetadata);
  const DatatableNext = () => tableNextPage(datatablesMetadata);
  const DatatablePrevious = () => tablePreviousPage(datatablesMetadata);

  onMount(() => {
    if (check("documentAi.documentLibrary.document:view") === false) {
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
      await getData();
    });

    return () => {
      unsubscribe();
    };
  });
</script>

<PageFileViewer
  bind:data
  bind:tags
  bind:CSVMetadata
  bind:CSVData
  bind:datatablesMetadata
  bind:dataTables
  bind:dataTablesSelected
  {onEdit}
  {base64}
  {fileUrl}
  {uploadBy}
  {fileReady}
  {libraryId}
  {documentId}
  {defaultData}
  {collectionId}
  {pageDataType}
  {attributeNotUnique}
  {getData}
  {updateData}
  {deleteData}
  {CSVNext}
  {CSVPrevious}
  {DatatableNext}
  {DatatablePrevious}
/>
