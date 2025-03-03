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
  import { Modal, Button, Input, Helper } from "flowbite-svelte";
  import {
    StateCard,
    Typography,
    SectionHeader,
    ActionCard,
    check,
  } from "$lib";
  import DataSourcesCard from "$lib/ui/cards/DataSourcesCard.svelte";
  import { _ } from "svelte-i18n";
  import {
    ExclamationCircleSolid,
  } from "flowbite-svelte-icons";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { goto } from "$app/navigation";
  import {
    showAI,
    showNotif,
    call,
    datasourceLibraryId,
    setDatasourceLibraryIds,
    triggerToast,
  } from "$lib";

  let libraryId;
  let collections = $state([]);
  let latestDocumentCollection = $state([]);
  let files = $state([]);
  let actionItems = $state([]);
  let tmpOldName = $state("");
  let tmpCollectionName = $state("");
  let tmpDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpDocumentTitle = $state("");
  let addLibraryModal = $state(false);
  let editModal = $state(false);
  let deleteModal = $state(false);
  let timer;
  let nameIsAvailable = $state(true);

  async function getActionItems() {
    actionItems = [
      {
        onclick: () => {
          addLibraryModal = true;
          tmpCollectionName = "";
          tmpDescription = "";
        },
        // href: await $call.getPage("client-dashboard", "document-collection-create"),
        icon: "filePlus",
        actionName: "Add data source",
        desc: "desc Add data source",
      },
    ];
  }

  let gettingDatas = $state(false);

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      collections = [];
      let response = await $call.call("client-dashboard", "get-library-collections", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: 1,
          page_size: 8,
          parent_id: parseInt(libraryId),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
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
        return;
      }

      if (response.result.data.length > 0) {
        collections = collections.concat(response.result.data);
      }

      let link = await $call.getPage("client-dashboard", "data-sources-documents");

      let promise = collections.map(async (collection) => {
        console.log("collection", collection);
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

        return {
          id: collection.id,
          collectionName: collection.attributes.name,
          description: collection.attributes?.meta?.description || "",
          type: collection.attributes.type,
          href: link + "?collectionId=" + collection.id + "&" + `libraryId=${libraryId}`,
        };
      });

      latestDocumentCollection = await Promise.all(promise);

      console.log("latestDocumentCollection", latestDocumentCollection);

      gettingDatas = false;
    }
  }

  const handleInput = async (event) => {
    const name = event.target.value;
    const status = await checkLibraryName(name.trim());

    if (tmpOldName === name) {
      nameIsAvailable = true;
    } else if (status === false) {
      nameIsAvailable = false;
    } else if (status === true) {
      nameIsAvailable = true;
    }
  };

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

  async function createCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "create-library-collection", {
        URLQuery: {},
        URLParams: {},
        Params: {
          collection_name: tmpCollectionName,
          description: tmpDescription,
          parent_id: parseInt(libraryId),
        },
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

      tmpCollectionName = "";
      tmpDescription = "";

      triggerToast($_("successCreateCollection"), true);

      addLibraryModal = false;
      getCollections();
    }
  }

  async function editCollectionName() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "update-library-collection", {
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

        tmpCollectionName = "";
        tmpDescription = "";

        editModal = false;
        triggerToast($_("successUpdateCollection"), true);
        getCollections();
      } catch (error) {
        console.error("Error during collection editing:", error);
        triggerToast($_("failedEditCollection"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  async function deleteDocumentCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "delete-library-collection", {
          URLQuery: {},
          URLParams: {
            id: tmpCollectionId.toString(),
          },
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

        triggerToast($_("successDeleteCollection"), true);

        getCollections();
        deleteModal = false;
      } catch (error) {
        console.error("Error during collection deletion:", error);
        triggerToast($_("failedDeleteCollection"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  async function viewAll() {
    let link = await $call.getPage("client-dashboard", "data-sources-collections");
    goto(link + `?libraryId=${libraryId}`);
  }

  onMount(async () => {
    getActionItems();

    // let idDatasource = localStorage.getItem("libraryIdDatasource");
    // await setDatasourceLibraryIds(idDatasource);

    libraryId = localStorage.getItem("libraryId" + "insights");

    getCollections();
    // latestDataSources = JSON.parse(localStorage.getItem("latestDataSources")) || [];
  });

  onDestroy(() => {
    if (timer) {
      clearTimeout(timer);
    }
  });

  let buttonDisabled = $state(false);

  $effect(() => {
    buttonDisabled = files.length === 0 || tmpDocumentTitle === "";
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="flex flex-col w-full pb-20 m-4 gap-3 md:gap-8">
    <!-- User Information -->
    <Typography
      preset="Display sm/Semibold"
      text={$_("dataSources")}
      propClass="fluid-xl"
    />

    {#if check("documentAi.documentLibrary:create")}
      <!-- Action Items -->
      <div class="flex flex-col gap-4">
        <SectionHeader title={$_("action")} />
        <ActionCard {actionItems} />
      </div>
    {/if}

    <div class="flex flex-col gap-4">
      <div
        class="flex items-center justify-between border-b border-secondary-light dark:border-secondary-dark pb-4"
      >
        <Typography preset="Text lg/Semibold" color="primary" text={$_("dataSourcesList")} />
        <Button on:click={() => viewAll()} color="light">{$_("seeAll")}</Button>
      </div>

      <div>
        {#if latestDocumentCollection.length > 0}
          <div
            class={`${$showAI || $showNotif ? "grid grid-cols-1 xl:grid-cols-3" : "grid grid-cols-2 lg:grid-cols-4"} gap-4`}
          >
            {#each latestDocumentCollection as document, index}
              <DataSourcesCard
                collectionName={document.collectionName}
                totalDocuments={document.totalDocuments}
                status={document.status}
                linkIsFrontend={'/data_sources/collections/documents?collectionName=' + document.collectionName}
                functEdit={() => {
                  tmpCollectionId = document.collectionId;
                  tmpOldName = document.collectionName;
                  tmpCollectionName = document.collectionName;
                  tmpDescription = document.description;
                  nameIsAvailable = true;
                  editModal = true;
                }}
                functDelete={() => {
                  tmpCollectionId = document.collectionId;
                  tmpCollectionName = document.collectionName;
                  deleteModal = true;
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

<Modal title={$_("createNewDataSources")} bind:open={addLibraryModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("dataSourcesName")}*`} />
      <Input bind:value={tmpCollectionName} on:input={handleInput} placeholder={$_("whatIsYourDataSourcesTitle")} />
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
      onclick={createCollection}
      disabled={!tmpCollectionName || !nameIsAvailable}
      class="w-full">{$_("addNewDataSources")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editDataSourcesName")} bind:open={editModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("dataSourcesName")}*`} />
      <Input bind:value={tmpCollectionName} on:input={handleInput} placeholder={$_("whatIsYourDataSourcesTitle")} />
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
      "{tmpCollectionName}"?
    </h3>
    <Button onclick={deleteDocumentCollection} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>
