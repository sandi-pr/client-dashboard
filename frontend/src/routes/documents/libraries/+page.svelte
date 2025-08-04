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
    me,
    StateCard,
    Typography,
    CollectionCard,
    LibraryCard,
    triggerToast,
    ToastContainer,
    UploadFile,
    check,
  } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { format } from "date-fns";
  import { Button, Modal, Input, Helper, Search } from "flowbite-svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import FileTypeIcon from "$lib/icon/file-type-icon.svelte";
  import {
    ExclamationCircleSolid,
    ChevronDownOutline,
    InfoCircleSolid,
  } from "flowbite-svelte-icons";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import ArrowLeftIcon from "$lib/icon/arrow-left.svelte";
  import ArrowRightIcon from "$lib/icon/arrow-right.svelte";
  import FilePlus02Icon from "$lib/icon/file-plus-02.svelte";

  import {
    triggerSelectCollection,
    triggerShowAI,
    triggerSelectDropdown,
    triggerReload,
    triggerLoading,
    showAI,
    showNotif,
    isBakrie,
  } from "$lib";
  import LoadingContainer from "$lib/components/LoadingContainer.svelte";

  import { goto } from "$app/navigation";

  import { call } from "$lib";

  const isFrontend = process.env.NODE_ENV === "development";

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
  let tmpOldName = $state("");
  let tmpDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpHref = $state("");
  let categories = $state([]);
  let otherCategory = $state("");
  let addOtherCategory = $state(false);
  let tmpCategories = $state("");
  let tmpDocumentType = $state("");
  let tmpDocumentTitle = $state("");
  let tmpDocumentDescription = $state("");
  let addLibraryModal = $state(false);
  let addDocumentsModal = $state(false);
  let editModal = $state(false);
  let deleteModal = $state(false);
  let dropdownOpen = $state(false);
  let page = $state(1);
  let pagination = $state({});
  let fileNotSupport = $state(false);
  let attributeNotUnique = $state(false);
  let nameIsAvailable = $state(true);

  let acceptedFileTypes = $state(".txt, .pdf");

  let selectedTypeCollection = $state("Documents (TXT, PDF)");

  let statusRAG = $state("");
  let timer;

  let searchTerm = $state("");

  async function createCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "create-library", {
        URLQuery: {},
        URLParams: {},
        Params: {
          library_name: tmpCollectionName,
          description: tmpDescription,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        if (response && response.errors) {
          const errors = response.errors.flatMap((error) =>
            error.errors.flatMap((subError) => subError.message),
          );

          if (errors.includes("This attribute must be unique")) {
            attributeNotUnique = true;
          } else {
            console.log("Pesan error :", errors);
            triggerToast("error");
          }
        }
        return;
      }

      console.log(response);

      // RAG Proccessing

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
      attributeNotUnique = false;

      triggerToast("success");

      addLibraryModal = false;
      getCollections();
      // goto(await $call.getPage("client-dashboard", "document-library"));
    }
  }

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      collections = [];
      documents = [];
      let response = await $call.call("client-dashboard", "get-librarys", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: page,
          page_size: 20,
          search: searchTerm,
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

      // console.log("collections", response);

      collections = response.result.data;

      pagination = response.result.meta.pagination;

      let link = isFrontend
        ? "/documents/libraries/collections"
        : await $call.getPage("client-dashboard", "document-library-collections");

      let promise = collections.map(async (collection) => {
        let response4 = await $call.call("client-dashboard", "get-status-library", {
          URLQuery: {},
          URLParams: {
            id: collection.id.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response4.status !== 200) {
          console.log(response4);
          gettingDatas = false;
          return;
        }

        // console.log("status", response4);

        let collectionStatus = response4.result.status;

        return {
          collectionId: collection.id,
          activationTime: format(collection.attributes.createdAt, "dd-MM-yyyy HH:mm"),
          collectionName: collection.attributes.name,
          description: collection.attributes?.meta?.description || "",
          // totalDocuments: totalDocuments,
          // categories: categoriesCollection,
          status: collectionStatus,
          href: `${link}?libraryId=${collection.id}`,
          webchat: collection.attributes.webchat_tmp === true ? true : false,
        };
      });

      datas = await Promise.all(promise);

      if ($me.authority.toLowerCase() === "demo") {
        // filter out collection name that doesn't have "demo" or "Demo" in it
        datas = datas.filter((doc) => doc.collectionName.toLowerCase().includes("demo"));
      }

      gettingDatas = false;

      // timer = setTimeout(getCollections, 2000);
    }
  }

  async function editCollectionName() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "update-library", {
          URLQuery: {},
          URLParams: {
            id: tmpCollectionId.toString(),
          },
          Params: {
            library_name: tmpCollectionName,
            description: tmpDescription,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          if (response && response.errors) {
            const errors = response.errors.flatMap((error) =>
              error.errors.flatMap((subError) => subError.message),
            );

            if (errors.includes("This attribute must be unique")) {
              attributeNotUnique = true;
            } else {
              console.log("Pesan error :", errors);
              triggerToast("error");
            }
          }
          return;
        }

        tmpCollectionName = "";
        tmpDescription = "";
        attributeNotUnique = false;

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
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        // for (const document of documents) {
        //   if (document.attributes.collection_id === tmpCollectionId) {
        //     let response = await $call.call("client-dashboard", "delete-document", {
        //       URLQuery: {},
        //       URLParams: {
        //         id: document.id.toString(),
        //       },
        //       Params: {},
        //       Headers: {
        //         Authorization: "Bearer " + token,
        //       },
        //     });
        //     if (response.status !== 200) {
        //       console.log(response);
        //       triggerToast("error");
        //       return;
        //     }
        //   }
        // }
        let response2 = await $call.call("client-dashboard", "delete-library", {
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
          parent_id: parentId ? parseInt(parentId) : 0,
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

  const handleInput = async (event) => {
    const name = event.target.value;
    const status = await checkLibraryName(name.trim(), 0, "library");

    if (tmpOldName === name) {
      nameIsAvailable = true;
    } else if (status === false) {
      nameIsAvailable = false;
    } else if (status === true) {
      nameIsAvailable = true;
    }
  };

  let buttonDisabled = $state(false);

  $effect(() => {
    buttonDisabled = files.length === 0 || tmpDocumentTitle === "";
  });

  function nextPage() {
    if (page < pagination.pageCount) {
      page = page + 1;
      getCollections();
    }
  }

  function previousPage() {
    if (page > 1) {
      page = page - 1;
      getCollections();
    }
  }

  onMount(async () => {
    // await getStatusRag();
    if (
      $isBakrie &&
      typeof window !== "undefined" &&
      window.document !== null &&
      window.localStorage.getItem("libraryIdWebchatBakrie")
    ) {
      const libraryIdWebchatBakrie = window.localStorage.getItem("libraryIdWebchatBakrie");
      const redirectUrl = `libraries/collections/?libraryId=${libraryIdWebchatBakrie}`;
      goto(redirectUrl);
    }
    if (check("documentAi.documentLibrary:view") === false) {
      window.history.back();
      return;
    }
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
      `${await $call.getPage("client-dashboard", "document-library-documents")}` +
      "?collectionId=" +
      id
    );
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

  async function askToAI(id) {
    return (
      `${await $call.getPage("client-dashboard", "document-collection-chat")}` +
      "?collectionId=" +
      id
    );
  }
</script>

{#if !$isBakrie}
  <div class="flex dark:bg-primary-dark">
    <div class="w-full pb-10">
      <div class="flex flex-col m-4">
        <div class="flex flex-col gap-3 md:gap-8">
          <!-- User Information -->
          <div class="flex items-center justify-between">
            <Typography
              preset="Display sm/Semibold"
              color="primary"
              text={$_("libraryList")}
              propClass="fluid-xl"
            />
            <Button
              outline
              class="min-w-max gap-1.5"
              disabled={!check("documentAi.documentLibrary:create")}
              on:click={() => {
                addLibraryModal = true;
                attributeNotUnique = false;
                nameIsAvailable = true;
                tmpCollectionName = "";
                tmpDescription = "";
              }}
            >
              <FilePlus02Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("addNew")}</div>
            </Button>
          </div>

          <!-- Bots List Data -->
          <div>
            <Search
              class="bg-primary-light dark:bg-primary-dark mb-2 md:mb-8 rounded-lg"
              placeholder={$_("searchLibraryName")}
              bind:value={searchTerm}
              on:input={getCollections}
            />
            {#if datas.length > 0}
              {#if !gettingDatas}
                <div
                  class={`${$showAI || $showNotif ? "grid grid-cols-1 xl:grid-cols-3" : "grid grid-cols-2 lg:grid-cols-4"} items-center gap-4`}
                >
                  {#each datas as document}
                    <LibraryCard
                      collectionName={document.collectionName}
                      status={document.status}
                      functButton={() => {
                        // selectCollection({
                        //   id: document.collectionId,
                        //   name: document.collectionName,
                        //   description: document.description,
                        //   totalDocuments: document.totalDocuments,
                        // });
                      }}
                      functAddDocuments={() => {
                        tmpCollectionId = document.collectionId;
                        tmpHref = document.href;
                        // categories = document.categories;
                        // selectTypeCollection("Documents");
                        addDocumentsModal = true;
                        addOtherCategory = false;
                      }}
                      functEdit={() => {
                        tmpCollectionId = document.collectionId;
                        tmpCollectionName = document.collectionName;
                        tmpOldName = document.collectionName;
                        tmpDescription = document.description;
                        attributeNotUnique = false;
                        nameIsAvailable = true;
                        editModal = true;
                      }}
                      functDelete={() => {
                        tmpCollectionId = document.collectionId;
                        tmpCollectionName = document.collectionName;
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
              {:else}
                <StateCard
                  illustration={IllustrationProgress}
                  description={$_("fetchingDataPleaseWait")}
                />
              {/if}
              {#if !gettingDatas}
                <div class="flex justify-between">
                  <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
                    {$_("showing")}
                    <span
                      >{pagination.page === 1
                        ? 1
                        : (pagination.page - 1) * pagination.pageSize + 1}</span
                    >
                    -
                    <!-- {$_("to")} -->
                    <span
                      >{pagination.page * pagination.pageSize >= pagination.total
                        ? pagination.total
                        : pagination.page * pagination.pageSize}</span
                    >
                    {$_("of")}
                    <span>{pagination.total}</span>
                    <!-- {$_("entries")} -->
                  </div>

                  {#if pagination.total > pagination.pageSize}
                    <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
                      {$_("page")}
                      <span>{pagination.page}</span>
                      {$_("of")}
                      <span>{pagination.pageCount}</span>
                    </div>
                  {/if}
                </div>
                {#if pagination.total > pagination.pageSize}
                  <div class="flex justify-end">
                    <Button
                      disabled={pagination.page === 1}
                      on:click={previousPage}
                      color="light"
                      class="fluid-sm gap-1 md:gap-2 rounded-l-lg rounded-r-none px-3 py-2 md:px-4 md:py-2.5"
                    >
                      <ArrowLeftIcon
                        svgClass="text-fg-secondary-light dark:text-fg-secondary-dark w-4 h-4 md:w-5 md:h-5"
                      />
                      {$_("previous")}
                    </Button>
                    <Button
                      disabled={pagination.page * pagination.pageSize >= pagination.total}
                      on:click={nextPage}
                      color="light"
                      class="fluid-sm gap-1 md:gap-2 rounded-r-lg rounded-l-none px-3 py-2 md:px-4 md:py-2.5"
                    >
                      {$_("next")}
                      <ArrowRightIcon
                        svgClass="text-fg-secondary-light dark:text-fg-secondary-dark w-4 h-4 md:w-5 md:h-5"
                      />
                    </Button>
                  </div>
                {/if}
              {/if}
              <!-- <div class="flex-col md:flex justify-between">
              <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
                {$_("showing")}
                <span
                  >{pagination.page === 1
                    ? 1
                    : (pagination.page - 1) * pagination.pageSize + 1}</span
                >
                {$_("to")}
                <span
                  >{pagination.page * pagination.pageSize >= pagination.total
                    ? pagination.total
                    : pagination.page * pagination.pageSize}</span
                >
                {$_("of")}
                <span>{pagination.total}</span>
                {$_("entries")}
              </div>

              {#if pagination.total > pagination.pageSize}
                <div class="m-4 flex justify-end">
                  <Pagination
                    on:next={nextPage}
                    on:previous={previousPage}
                    activeClass="bg-hc-graylight-100 dark:bg-hc-graydark-700 dark:border-[#374151] dark:text-hc-graydark-300"
                  >
                    <span slot="prev">{$_("previous")}</span>
                    <span slot="next">{$_("next")}</span>
                  </Pagination>
                </div>
              {/if}
            </div> -->
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
      </div>
    </div>
    <ToastContainer />
  </div>
{/if}

<Modal title={$_("createNewLibrary")} bind:open={addLibraryModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("libraryName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={handleInput}
        placeholder={$_("whatIsYourLibraryTitle")}
      />
      {#if !nameIsAvailable}
        <Helper color="red">
          <span class="font-medium">{$_("nameAlreadyTaken")}</span>
        </Helper>
      {/if}
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={createCollection}
      disabled={!tmpCollectionName || !nameIsAvailable}
      class="w-full">{$_("addNewLibrary")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editLibraryName")} bind:open={editModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("libraryName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={handleInput}
        placeholder={$_("whatIsYourLibraryTitle")}
      />
      {#if !nameIsAvailable}
        <Helper color="red">
          <span class="font-medium">{$_("nameAlreadyTaken")}</span>
        </Helper>
      {/if}
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}`} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={editCollectionName}
      disabled={!tmpCollectionName || !nameIsAvailable}
      class="me-2 w-full">{$_("save")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
      {$_("library")}
      {tmpCollectionName}?
    </h3>
    <Button onclick={deleteDocumentCollection} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>

<!-- <LoadingContainer /> -->
