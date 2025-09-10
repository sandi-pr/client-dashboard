<script>
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { Button, Modal, Input, Helper, Search } from "flowbite-svelte";

  import {
    call,
    check,
    getIcons,
    checkLibraryName,
    getCollectionAPIConnector,
    createCollectionAPIConnector,
    updateCollectionAPIConnector,
    deleteCollectionAPIConnector,
  } from "$lib";
  import { Typography, ToastContainer, StateCard, LibraryCard } from "$lib";

  const FilePlus02Icon = getIcons("filePlus02");
  const IllustrationEmpty = getIcons("illustrationEmpty");
  const IllustrationProgress = getIcons("illustrationProgress");

  const isFrontend = process.env.NODE_ENV === "development";

  let editModal = $state(false);
  let deleteModal = $state(false);
  let gettingDatas = $state(false);
  let nameIsAvailable = $state(true);
  let addCollectionModal = $state(false);

  let libraryId = $state(0);
  let tmpParentId = $state(0);

  let tmpOldName = $state("");
  let searchTerm = $state("");
  let tmpDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpCollectionName = $state("");

  let collectionListDatas = $state([]);

  const handleCreateCollection = async () => {
    gettingDatas = true;
    addCollectionModal = false;

    const result = await createCollectionAPIConnector($_, $call, tmpCollectionName, tmpDescription);

    if (result) {
      handleGetCollections();
    } else {
      gettingDatas = false;
    }
  };

  const handleGetCollections = async () => {
    gettingDatas = true;
    searchTerm = searchTerm.trim();
    collectionListDatas = await getCollectionAPIConnector($_, $call, searchTerm);
    gettingDatas = false;
  };

  const handleUpdateCollection = async () => {
    gettingDatas = true;
    editModal = false;

    const result = await updateCollectionAPIConnector(
      $_,
      $call,
      tmpParentId ? parseInt(tmpParentId) : 0,
      tmpCollectionId,
      tmpCollectionName,
      tmpDescription,
    );

    if (result) {
      handleGetCollections();
    } else {
      gettingDatas = false;
    }
  };

  const handleDeleteCollection = async () => {
    deleteModal = false;
    gettingDatas = true;

    const result = await deleteCollectionAPIConnector($_, $call, tmpCollectionId);

    if (result) {
      handleGetCollections();
    } else {
      gettingDatas = false;
    }
  };

  const handleInput = async (event, parentId, type) => {
    const name = event.target.value;
    const status = await checkLibraryName($call, name.trim(), parentId, type);

    if (tmpOldName === name) {
      nameIsAvailable = true;
    } else if (status === false) {
      nameIsAvailable = false;
    } else if (status === true) {
      nameIsAvailable = true;
    }
  };

  onMount(async () => {
    if (check("documentAi.documentLibrary:view") === false) {
      window.history.back();
      return;
    }
    await handleGetCollections();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-start justify-between mb-4 md:mb-8">
          <div class="flex flex-col gap-1.5">
            <div class="flex gap-4">
              <Typography
                preset="Display sm/Semibold"
                color="primary"
                text={$_("apiCollectionList")}
                propClass="fluid-xl"
              />
            </div>
          </div>

          <Button
            outline
            class="w-fit gap-1.5 ms-14"
            disabled={!check("documentAi.documentLibrary.collection:create")}
            on:click={() => {
              addCollectionModal = true;
              nameIsAvailable = true;
              tmpCollectionName = "";
              tmpDescription = "";
            }}
          >
            <FilePlus02Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
            <div class="hidden md:block">{$_("addNewCollection")}</div>
          </Button>
        </div>

        <div>
          <Search
            class="bg-primary-light dark:bg-primary-dark mb-2 md:mb-4 rounded-lg"
            placeholder={$_("search")}
            bind:value={searchTerm}
            on:input={handleGetCollections}
          />
        </div>

        {#if collectionListDatas.length > 0 && !gettingDatas}
          <div class="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mt-4">
            {#each collectionListDatas as collection}
              <LibraryCard
                collectionName={collection.collectionName}
                usedFor="libraries_collections"
                functEdit={() => {
                  tmpCollectionId = collection.id;
                  tmpCollectionName = collection.collectionName;
                  tmpOldName = document.collectionName;
                  tmpDescription = collection.description;
                  tmpParentId = collection.parentId;
                  nameIsAvailable = true;
                  editModal = true;
                }}
                functDelete={() => {
                  tmpCollectionId = collection.id;
                  tmpCollectionName = collection.collectionName;
                  deleteModal = true;
                }}
                funcDetailCollection={() => {
                  goto(collection.href);
                }}
                {libraryId}
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
  <ToastContainer />
</div>

<Modal title={$_("createNewCollection")} bind:open={addCollectionModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("collectionName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={(event) => handleInput(event, libraryId, "collection")}
        placeholder={$_("whatIsYourCollectionTitle")}
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
      onclick={() => handleCreateCollection()}
      disabled={!tmpCollectionName || !nameIsAvailable}
      class="w-full">{$_("addNewCollection")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editCollectionName")} bind:open={editModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("collectionName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={(event) => handleInput(event, libraryId, "collection")}
        placeholder={$_("whatIsYourCollectionTitle")}
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
      onclick={handleUpdateCollection}
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
      {$_("collection")}
      "{tmpCollectionName}"?
    </h3>
    <Button onclick={handleDeleteCollection} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>
