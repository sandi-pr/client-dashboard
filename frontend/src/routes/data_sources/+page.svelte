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
  import { Modal, Button, Input } from "flowbite-svelte";
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
  } from "$lib";

  let files = $state([]);
  let actionItems = $state([]);
  let tmpCollectionName = $state("");
  let tmpDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpDocumentTitle = $state("");
  let addLibraryModal = $state(false);
  let editModal = $state(false);
  let deleteModal = $state(false);
  let timer;
  let attributeNotUnique = $state(false);
  let nameIsAvailable = $state(true);

  async function getActionItems() {
    actionItems = [
      {
        onclick: () => {
          addLibraryModal = true;
          attributeNotUnique = false;
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

  let latestDataSources = $state([]);

  let gettingDatas = $state(false);

  function createDataSources() {
    latestDataSources = [
      ...latestDataSources,
      {
        collectionId: latestDataSources.length + 1,
        collectionName: tmpCollectionName,
      },
    ];

    localStorage.setItem("latestDataSources", JSON.stringify(latestDataSources));

    tmpCollectionName = "";
    addLibraryModal = false;
  }

  let editIndex = $state(0);

  function editDataSourcesName() {
    latestDataSources[editIndex].collectionName = tmpCollectionName;
    localStorage.setItem("latestDataSources", JSON.stringify(latestDataSources));
    editModal = false;
  }

  function deleteDataSources(id) {
    latestDataSources = latestDataSources.filter(
      (item) => item.collectionId !== id,
    );
    localStorage.setItem("latestDataSources", JSON.stringify(latestDataSources));
    deleteModal = false;
  }

  onMount(async () => {
    getActionItems();
    latestDataSources = JSON.parse(localStorage.getItem("latestDataSources")) || [];
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
        <!-- <Button on:click={() => viewAll()} color="light">{$_("seeAll")}</Button> -->
      </div>

      <div>
        {#if latestDataSources.length > 0}
          <div
            class={`${$showAI || $showNotif ? "grid grid-cols-1 xl:grid-cols-3" : "grid grid-cols-2 lg:grid-cols-4"} gap-4`}
          >
            {#each latestDataSources as document, index}
              <DataSourcesCard
                collectionName={document.collectionName}
                totalDocuments={document.totalDocuments}
                status={document.status}
                linkIsFrontend={'/data_sources/collections/documents?collectionName=' + document.collectionName}
                functEdit={() => {
                  tmpCollectionId = document.collectionId;
                  tmpCollectionName = document.collectionName;
                  tmpDescription = document.description;
                  attributeNotUnique = false;
                  nameIsAvailable = true;
                  editIndex = index;
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
      <Input bind:value={tmpCollectionName} placeholder={$_("whatIsYourDataSourcesTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}`} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={createDataSources}
      disabled={!tmpCollectionName}
      class="w-full">{$_("addNewDataSources")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editDataSourcesName")} bind:open={editModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("dataSourcesName")}*`} />
      <Input bind:value={tmpCollectionName} placeholder={$_("whatIsYourDataSourcesTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={editDataSourcesName}
      disabled={!tmpCollectionName}
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
    <Button onclick={() => deleteDataSources(tmpCollectionId)} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>
