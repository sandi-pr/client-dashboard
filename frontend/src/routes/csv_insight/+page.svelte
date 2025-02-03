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
  import { Modal, Button, Input, Card, Avatar } from "flowbite-svelte";
  import {
    StateCard,
    Typography,
    triggerToast,
  } from "$lib";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import LightbulbIcon from "$lib/icon/lightbulb-04.svg";
  import Star04Icon from "$lib/icon/star-04.svelte";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
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
  let tmpCollectionName = $state("");
  let tmpDescription = $state("");
  let tmpCollectionId = $state("");
  let tmpDocumentTitle = $state("");
  let addLibraryModal = $state(false);
  let editModal = $state(false);
  let deleteModal = $state(false);
  let timer;

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

  let dropdownOpen = $state([]);

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    for (let i = 0; i < dropdownOpen.length; i++) {
      dropdownOpen[i] = false;
    }
  };

  onMount(async () => {
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
      text={$_("csvInsights")}
      propClass="fluid-xl"
    />

    <div class="flex flex-col gap-4">
      <div
        class="flex items-center justify-between border-b border-secondary-light dark:border-secondary-dark pb-4"
      >
        <Typography preset="Text lg/Semibold" color="primary" text={$_("csvInsightsBotsList")} />
        <!-- <Button on:click={() => viewAll()} color="light">{$_("seeAll")}</Button> -->
      </div>

      <div>
        {#if latestDataSources.length > 0}
          <div
            class={`${$showAI || $showNotif ? "grid grid-cols-1 xl:grid-cols-3" : "grid grid-cols-2 lg:grid-cols-4"} gap-4`}
          >
            {#each latestDataSources as bot, index}
              <Card padding="none" class="relative shadow-none grow w-full">
                <!-- {#if bot.deployStatus === "Reading"}
                  <div class="absolute top-2 right-2">
                    <Badge type="warning" text={bot.deployStatus} />
                  </div>
                {:else if bot.deployStatus === "Ready"}
                  <div class="absolute top-2 right-2">
                    <Badge type="success" text={bot.deployStatus} />
                  </div>
                {/if} -->
                <div
                  class="bg-active-light dark:bg-active-dark border-b border-secondary-light dark:border-secondary-dark px-4 py-6"
                >
                  <div class="flex justify-between items-center">
                    <div class="flex items-center gap-3">
                      <Avatar src={LightbulbIcon} class="p-2 bg-fg-brand-primary-light dark:bg-fg-brand-primary-dark" />
                      <div class="flex flex-col">
                        <Typography
                          preset="Text lg/Semibold"
                          color="primary"
                          text={bot.collectionName}
                        />
                      </div>
                    </div>
                    <div
                      onfocusout={handleDropdownFocusLoss}
                      class=" flex items-center relative"
                    >
                      <Button
                        on:click={() => (dropdownOpen[index] = !dropdownOpen[index])}
                        class="focus-within:ring-0 !p-0"
                        color="none"
                        size="xs"
                        ><DotsVerticalIcon
                          class="text-fg-quinary-light dark:text-fg-quinary-dark"
                        /></Button
                      >
                      <div
                        class={`z-10 ${dropdownOpen[index] ? "absolute top-0 right-4" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
                      >
                        <ul class="py-2 grow text-sm">
                          <!-- <button
                            class="flex py-2 px-4 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                            onclick={goToEditPage(bot.botsID)}>{$_("edit")}</button
                          > -->
                          <button
                            class="flex py-2 px-4 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                            onclick={() => {
                              popupModal = true;
                              tmpId = bot.botsID;
                              dropdownOpen[index] = !dropdownOpen[index];
                            }}
                            >{$_("delete")}</button
                          >
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="flex gap-4 px-4 pt-4 pb-6">
                  <div class="grow gap-3">
                    <Button
                      on:click={async () => {
                        goto(bot.href);
                      }}
                      outline
                      class="!flex w-full gap-1.5 hover:bg-primary-light dark:hover:bg-primary-dark rounded-lg"
                    >
                      <Typography
                        preset="Text sm/Semibold"
                        text={$_("askForInsight")}
                        color="brand-secondary"
                      />
                      <Star04Icon svgClass="w-4 h-4 text-brand-secondary-light dark:text-brand-secondary-dark" />
                    </Button>
                  </div>
                  <button
                    onclick={() => {
                      navigator.clipboard
                        .writeText(bot.urlBot)
                        .then(() => {
                          console.log("Text copied to clipboard");
                          triggerToast("copyToClipboard");
                        })
                        .catch((err) => {
                          console.error("Failed to copy: ", err);
                        });
                    }}
                  >
                    <Copy01Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                  </button>
                </div>
              </Card>
              <!-- <DataSourcesCard
                collectionName={document.collectionName}
                totalDocuments={document.totalDocuments}
                status={document.status}
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
              /> -->
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
