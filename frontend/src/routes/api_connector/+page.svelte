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
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";
  import { Modal, Button, Input, Helper } from "flowbite-svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import {
    call,
    check,
    showAI,
    getIcons,
    isBakrie,
    showNotif,
    onProgress,
    errorHandling,
    notifications,
    triggerLoading,
    checkLibraryName,
    getCollectionAPIConnector,
    createCollectionAPIConnector,
    updateCollectionAPIConnector,
    deleteCollectionAPIConnector,
  } from "$lib";
  import { StateCard, ActionCard, Typography, LibraryCard, SectionHeader } from "$lib";

  const IllustrationEmpty = getIcons("illustrationEmpty");
  const IllustrationProgress = getIcons("illustrationProgress");

  const isFrontend = process.env.NODE_ENV === "development";

  let botId, botIdStrapi, collectionId;
  let timer;

  let editModal = $state(false);
  let deleteModal = $state(false);
  let gettingDatas = $state(false);
  let nameIsAvailable = $state(true);
  let buttonDisabled = $state(false);
  let gettingStatusBot = $state(false);
  let attributeNotUnique = $state(false);
  let addCollectionModal = $state(false);

  let libraryId = $state(0);
  let tmpParentId = $state(0);

  let tmpOldName = $state("");
  let tmpDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpDocumentTitle = $state("");
  let tmpCollectionName = $state("");

  let files = $state([]);
  let actionItems = $state([]);
  let latestApiCollection = $state([]);
  let bots = $state([
    {
      botsName: "Handwriting’s Assistant",
      botDescription: "",
      botsID: 12345,
      urlBot: "#",
      deployStatus: "Inactive",
      deployedPlatform: "",
      bot_profile: "",
    },
  ]);

  const getActionItems = async () => {
    actionItems = [
      {
        onclick: () => {
          addCollectionModal = true;
          attributeNotUnique = false;
          tmpCollectionName = "";
          tmpDescription = "";
        },
        icon: "dataflow02",
        actionName: "addNewApiCollection",
        desc: "descAddNewApiCollection",
      },
    ];
  };

  const getBot = async () => {
    try {
      let token = localStorage.getItem("token");
      if (token) {
        let response = await $call.call("client-dashboard", "get-bot", {
          URLQuery: {},
          URLParams: {
            id: botIdStrapi.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          const error = new Error("Failed to fetch bot data");
          error.response = response;
          throw error;
        }

        const bot = response.result;
        bots[0].deployStatus = bot.status;
      }
    } catch (error) {
      errorHandling($_, error, "get");
      goto(await $call.getPage("client-dashboard", "hexabot-list"));
    } finally {
      triggerLoading(false);
    }
  };

  const startReading = async (botId) => {
    try {
      triggerLoading(true);
      let token = localStorage.getItem("token");

      if (token) {
        gettingStatusBot = true;

        const response = await $call.call("client-dashboard", "add-document-by-bot-id", {
          URLQuery: {},
          URLParams: {
            id: botId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          const error = new Error("Failed to start reading process");
          error.response = response;
          throw error;
        }

        setTimeout(getBot, 3000);
      }
    } catch (error) {
      errorHandling($_, error, "get");
      triggerLoading(false);
    }
  };

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
    latestApiCollection = await getCollectionAPIConnector($_, $call);
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

  const handleInput = async (event, parentId = collectionId, type = "document") => {
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
    if (
      $isBakrie &&
      typeof window !== "undefined" &&
      window.document !== null &&
      window.localStorage.getItem("libraryIdWebchatBakrie")
    ) {
      const libraryIdWebchatBakrie = window.localStorage.getItem("libraryIdWebchatBakrie");
      const redirectUrl = `documents/libraries/collections/?libraryId=${libraryIdWebchatBakrie}`;
      goto(redirectUrl);
    }

    if (check("documentAi.documentLibrary:view") === false) {
      window.history.back();
      return;
    }

    await getActionItems();
    await handleGetCollections();
  });

  onDestroy(() => {
    if (timer) {
      clearTimeout(timer);
    }
  });

  $effect(() => {
    if ($onProgress) {
      // console.log("notification", $notifications);
      if (
        $notifications.length > 0 &&
        $notifications[0].attributes.meta?.progress.includes("100")
      ) {
        if (!gettingStatusBot) {
          startReading($notifications[0].attributes.meta.bot_id);
          console.log("start reading", $notifications[0].attributes.meta.bot_id);
        } else {
          if (timerGetStatusBot) {
            clearTimeout(timerGetStatusBot);
          }
          timerGetStatusBot = setTimeout(() => {
            startReading($notifications[0].attributes.meta.bot_id);
            console.log("start reading", $notifications[0].attributes.meta.bot_id);
          }, 10000);
        }
      }
    }
  });

  $effect(() => {
    buttonDisabled = files.length === 0 || tmpDocumentTitle === "";
  });
</script>

{#if !$isBakrie}
  <div class="flex dark:bg-primary-dark">
    <div class="flex flex-col w-full pb-20 m-4 gap-3 md:gap-8">
      <!-- User Information -->
      <Typography
        preset="Display sm/Semibold"
        text={$_("apiConnectorOverview")}
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
          <Typography preset="Text lg/Semibold" color="primary" text={$_("apiCollectionList")} />
          {#if latestApiCollection.length > 0}
            <Button
              on:click={async () => {
                const newLink = isFrontend
                  ? "/api_connector/collections"
                  : await $call.getPage("client-dashboard", "api-connector-collections");
                goto(newLink);
              }}
              color="light"
            >
              {$_("viewAll")}
            </Button>
          {/if}
        </div>

        <div>
          {#if latestApiCollection.length > 0 && !gettingDatas}
            <div
              class={`${$showAI || $showNotif ? "grid grid-cols-1 xl:grid-cols-3" : "grid grid-cols-2 lg:grid-cols-4"} gap-4`}
            >
              {#each latestApiCollection as api}
                <LibraryCard
                  collectionName={api.collectionName}
                  usedFor="libraries_collections"
                  status={api.status}
                  functEdit={() => {
                    tmpCollectionId = api.id;
                    tmpCollectionName = api.collectionName;
                    tmpOldName = api.collectionName;
                    tmpDescription = api.description;
                    tmpParentId = api.parent_id;
                    attributeNotUnique = false;
                    nameIsAvailable = true;
                    editModal = true;
                  }}
                  functDelete={() => {
                    tmpCollectionId = api.id;
                    tmpCollectionName = api.collectionName;
                    deleteModal = true;
                  }}
                  funcDetailCollection={() => {
                    goto(api.href);
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
  </div>
{/if}

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
      onclick={handleCreateCollection}
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
      {$_("library")}
      "{tmpCollectionName}"?
    </h3>
    <Button onclick={handleDeleteCollection} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>
