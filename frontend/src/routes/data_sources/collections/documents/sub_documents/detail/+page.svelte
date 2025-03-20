<script>
  import {
    Button,
    Input,
    Textarea,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    Modal,
    Helper,
  } from "flowbite-svelte";
  // import Datatable from "$lib/components/Datatable.svelte";
  import { _ } from "svelte-i18n";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { Typography, Divider, StateCard, call, triggerLoading, triggerToast, check, Pagination } from "$lib";
  import { onMount, onDestroy  } from "svelte";
  import MultiSelect from 'svelte-multiselect';
  
  import { page } from "$app/stores";
  // import { lib } from "crypto-js/core";

  let documentId, collectionId, libraryId;

  let deleteModal = $state(false);
  let canEdit = $state(false);
  let alltags = $state([]);
  let tags = $state([]);
  let gettingDatas = $state(false);
  let nameIsAvailable = $state(true);

  let data = $state({});
  let tmpData = $state({});

  let csvMeta = $state({ page: 1, pageSize: 10, pageCount: 2, total: 13 });

  const underscoreToSpace = (text) => {
    return text.replace(/_/g, " ");
  };
  
  let currentArrayData = $state([]);

  async function getDocument() {
    triggerLoading(true);
    gettingDatas = true;
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "get-library-document-csv", {
          URLQuery: {},
          URLParams: {
            id: documentId.toString(),
          },
          Params: {
            page: csvMeta.page,
            page_size: csvMeta.pageSize,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("document csv response", response);

        if (response.status !== 200) {
          console.log(response);
          triggerLoading(false);
          gettingDatas = false;
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

        // csvMeta = { page: 1, pageSize: 10, pageCount: 2, total: 13 }

        csvMeta.pageCount = response.result.result.total_pages;
        csvMeta.total = response.result.result.total_items;

        currentArrayData = response.result.result.data;

        console.log("currentArrayData", currentArrayData);

        tags =
          documentData.meta.category === undefined ? [] : documentData.meta.category.split(",");

        console.log("data csv", data);

        // fileReady = true;
        triggerLoading(false);
        gettingDatas = false;
      } catch (error) {
        triggerLoading(false);
        gettingDatas = false;
        console.error("Error during document get:", error);
        triggerToast($_("failedGetDocument"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  async function getCsvData() {
    gettingDatas = true;
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "get-library-document-csv", {
          URLQuery: {},
          URLParams: {
            id: documentId.toString(),
          },
          Params: {
            page: csvMeta.page,
            page_size: csvMeta.pageSize,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("document csv response", response);

        if (response.status !== 200) {
          console.log(response);
          triggerLoading(false);
          gettingDatas = false;
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

        csvMeta.pageCount = response.result.result.total_pages;
        csvMeta.total = response.result.result.total_items;

        currentArrayData = response.result.result.data;

        console.log("currentArrayData", currentArrayData);

        gettingDatas = false;
      } catch (error) {
        gettingDatas = false;
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

  async function updateDocument() {
    triggerLoading(true);
    canEdit = false;
    let token = localStorage.getItem("token");
    if (token !== null) {
      let response;
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
        canEdit = true;
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
      triggerToast($_("successUpdateDocument"), true);
      getDocument();
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
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
          (await $call.getPage("client-dashboard", "data-sources-sub-documents")) +
            "?collectionId=" +
            collectionId +
            "&libraryId=" +
            libraryId,
        );
        // }, 3000);
        return;
      } catch (error) {
        console.error("Error during document deletion:", error);
        triggerToast($_("failedDeleteDocument"), true, "error");
      }
    }
  }

  async function next() {
    if (csvMeta.page >= csvMeta.pageCount) return;
    csvMeta.page += 1;
    // TODO: create function getCsv();
    getCsvData();
  }

  async function previous() {
    if (csvMeta.page > 1) {
      csvMeta.page -= 1;
      getCsvData();
    }
  }

  const handleInput = async (event) => {
    const name = event.target.value;
    const status = await checkLibraryName(name.trim(), collectionId, "document");

    if (status === true) {
      nameIsAvailable = true;
    } else if (status === false) {
      nameIsAvailable = false;
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

  onMount(() => {
    if (check("documentAi.documentLibrary:view") === false) {
      window.history.back();
      return;
    }
    const unsubscribe = page.subscribe(async ($page) => {
      documentId = $page.url.searchParams.get("documentId");
      libraryId = $page.url.searchParams.get("libraryId");
      if (documentId === null || documentId === undefined) {
        triggerToast("error");
        goto(await $call.getPage("client-dashboard", "data-sources"));
        return;
      }
      await getDocument();
      getLabels();
    });
    return () => unsubscribe();
  });
</script>

<div class="m-4">
  <div class="flex justify-between items-center">
    <Typography text={tmpData.document_title} preset="Display sm/Semibold" />
    <div class="flex gap-2">
      {#if canEdit}
        <Button color="alternative" on:click={() => {canEdit = false; data.document_title = tmpData.document_title}}>
          <Typography text="Cancel" preset="Text sm/Semibold" color="tertiary" />
        </Button>
        <Button color="primary" disabled={data.document_title === "" || !nameIsAvailable} on:click={() => updateDocument()}>
          <Typography text="Save" preset="Text sm/Semibold" color="primary-onbrand" />
        </Button>
      {:else}
        <Button color="alternative" on:click={() => (deleteModal = true)}>
          <Typography text="Delete" preset="Text sm/Semibold" color="error-primary" />
        </Button>
        <Button color="alternative" on:click={() => (canEdit = true)}>
          <Typography
            text="edit"
            preset="Text sm/Semibold"
            propClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
          />
        </Button>
      {/if}
    </div>
  </div>
  <Divider className="my-4" />
  <!-- general information -->
  <div class="flex gap-6">
    <div class="w-60">
      <Typography preset="Text md/Semibold" color="secondary" text={$_("generalInformation")} />
    </div>
    <div class="flex flex-col grow gap-6">
      <div class="flex gap-5">
        <div class="flex flex-col grow gap-2">
          <Typography preset="Text sm/Semibold" color="secondary" text={$_("documentTitle")} />
          <Input
            readonly={!canEdit}
            oninput={handleInput}
            bind:value={data.document_title}
            class="bg-primary-light dark:bg-secondary-dark focus:none focus:ring-0 focus:border-gray-300 dark:focus:border-gray-600"
          />
          {#if !nameIsAvailable}
            <Helper color="red">
              <span class="font-medium">{$_("nameAlreadyTaken")}</span>
            </Helper>
          {/if}
        </div>
      </div>
      <div class="grid grid-cols-1">
        <div class="flex flex-col gap-2">
          <Typography preset="Text sm/Semibold" color="secondary" text={$_("description")} />
          <Textarea
            readonly={!canEdit}
            bind:value={data.description}
            rows="4"
            class="bg-primary-light dark:bg-secondary-dark focus:none focus:ring-0 border-gray-300 dark:border-600 focus:!border-gray-300 dark:focus:!border-gray-600"
          />
        </div>
      </div>
      <div class="flex flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={`${$_("labels")}`} />
        <MultiSelect
          --sms-min-height="40px"
          options={alltags || []} 
          disabled={!canEdit}
          allowUserOptions="append"
          bind:selected={tags}
          autoScroll={true}
          noMatchingOptionsMsg=""
          ulSelectedClass="text-primary-light dark:text-primary-dark"
          liSelectedClass="border border-primary-light dark:border-primary-dark"
          outerDivClass="bg-primary-light dark:bg-secondary-dark border border-gray-300 dark:border-gray-600 focus:none focus:ring-0 focus:border-gray-300 dark:focus:border-gray-600"
          createOptionMsg={null}
        ><span slot="disabled-icon"></span></MultiSelect>
      </div>
    </div>
  </div>
  <Divider className="my-4" />
  <!-- Table Viewer -->
  <div class="flex gap-6">
    <div class="w-60">
      <Typography preset="Text md/Semibold" color="secondary" text={$_("generalInformation")} />
    </div>
    <div class="flex flex-col grow gap-6">
      <!-- <Datatable datas={dummyCsvDatas} configs={dummyCsvConfigs} /> -->
      {#if currentArrayData.length === 0}
        <StateCard illustration={IllustrationEmpty} description={$_("noDataAvailable")} />
      {:else if gettingDatas}
        <StateCard illustration={IllustrationProgress} description={$_("fetchingDataPleaseWait")} />
      {:else}
        <div
          class="bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-700 rounded-xl pt-4"
        >
          <Table class="w-full">
            <TableHead>
              {#each Object.keys(currentArrayData[0]) as key}
                {#if key === "id"}
                  <!-- empty -->
                {:else}
                  <TableHeadCell>{underscoreToSpace(key)}</TableHeadCell>
                {/if}
              {/each}
            </TableHead>
            <TableBody tableBodyClass="divide-y">
              {#each currentArrayData as data}
                <TableBodyRow class="border-0 !border-gray-200 dark:!border-gray-700">
                  {#each Object.entries(data) as [key, value]}
                    {#if key === "id"}
                      <!-- empty -->
                    {:else}
                      <TableBodyCell tdClass="px-6 py-4 font-medium">{value}</TableBodyCell>
                    {/if}
                  {/each}
                </TableBodyRow>
              {/each}
            </TableBody>
          </Table>
          <Pagination pagination={csvMeta} previousPage={previous} nextPage={next} />
        </div>
      {/if}
    </div>
  </div>
</div>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteDocument} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>
