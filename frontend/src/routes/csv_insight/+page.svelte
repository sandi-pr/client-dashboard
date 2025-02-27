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
    call,
    BotCard,
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
  
  let bots = $state([]);

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

  async function getBots() {
    let token = localStorage.getItem("token");

    if (token) {
      gettingDatas = true;
      try {
        // const response = await $call.call("client-dashboard", "get-bots", {
        const response = await $call.call("client-dashboard", "get-bots", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: 1,
            page_size: 10,
            search: "",
            module_type: "insights",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response === null) {
          gettingDatas = false;
          console.log("response null", response);
          return;
        }

        if (response.status !== 200) {
          gettingDatas = false;
          console.log("status not 200", response);
          return;
        }

        let tmp = response.result.data.map(async (bot) => {
          let fileImageResult = "";
          if (bot.attributes.bot_profile.data !== null) {
            let ext = bot.attributes.bot_profile.data.attributes.ext;

            let fileType;
            if (ext === ".png") {
              fileType = "image/png";
            } else if (ext === ".jpg") {
              fileType = "image/jpg";
            } else if (ext === ".jpeg") {
              fileType = "image/jpeg";
            } else if (ext === ".svg") {
              fileType = "image/svg+xml";
            } else {
              fileType = "image/png";
            }

            if (bot.bot_profile !== "") {
              let fileImage = bot.bot_profile;
              let blob = base64toBlob(fileImage, fileType);
              fileImageResult = URL.createObjectURL(blob);
            }
          }

          let urlBot = await $call.getPage("client-dashboard", "csv-insight-ask-csv");

          return {
            botsName: bot.attributes.bot_name,
            botDescription: bot.attributes.description,
            botsID: bot.id,
            // clientID: bot.attributes.client_id,
            activationTime: bot.attributes.createdAt,
            urlBot: urlBot,
            timestamp: bot.attributes.updatedAt,
            deployStatus: bot.attributes.status,
            deployedPlatform: bot.attributes.ai_engine,
            bot_profile: fileImageResult,
          };
        });

        bots = await Promise.all(tmp);

        dropdownOpen = [];
        for (let i = 0; i < bots.length; i++) {
          dropdownOpen.push(false);
        }

        gettingDatas = false;
        console.log("cek bot insight", bots);
      } catch (error) {
        gettingDatas = false;
        console.log(error);
      }
    } else {
      console.log("token not found");
    }
  }

  onMount(async () => {
    // latestDataSources = JSON.parse(localStorage.getItem("latestDataSources")) || [];
    getBots();
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
        {#if bots.length > 0}
          <div class={`grid ${$showAI || $showNotif ? "grid-cols-1 xl:grid-cols-3" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4"} gap-4`}>
            {#each bots as bot, index (bot.botsID)}
              <BotCard {bot} {index} {dropdownOpen} config={{ insights: true }} {handleDropdownFocusLoss} />
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
