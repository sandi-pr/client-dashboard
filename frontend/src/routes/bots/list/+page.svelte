<script>
  import { Button, Toast, Search, Modal } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import {
    Typography,
    call,
    triggerToast,
    StateCard,
    newNotification,
    notifications,
    BotCard,
    showAI,
    showNotif,
    Pagination,
    onProgress,
    triggerLoading,
  } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { onMount } from "svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";

  let searchTerm = $state("");
  let serachPagination = $state(false);
  let gettingDatas = $state(true);
  let bots = $state([]);
  let tmpId = $state(null);
  let page = $state(1);
  let pagination = $state({
    total: 0,
    page: 0,
    pageSize: 0,
  });

  let popupModal = $state(false);
  let gettingStatusBot = $state(false);

  function copyToClipboard(text) {
    navigator.clipboard.writeText(text);
    trigger();
  }

  let open = false;
  let counter = 6;

  function trigger() {
    open = true;
    counter = 6;
    timeout();
  }

  function timeout() {
    if (--counter > 0) return setTimeout(timeout, 1000);
    open = false;
  }

  async function getBots() {
    let token = localStorage.getItem("token");

    if (searchTerm !== "" && !serachPagination) {
      pagination.page = 1;
    }

    if (token) {
      gettingDatas = true;
      try {
        const response = await $call.call("client-dashboard", "get-bots", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 18,
            search: searchTerm,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response === null) {
          gettingDatas = false;
          console.log(response);
          return;
        }

        if (response.status !== 200) {
          gettingDatas = false;
          console.log(response);
          return;
        }

        pagination = response.result.meta.pagination;

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

          let urlBot =
            (await $call.getPage("hexabot", "home-page")) + "?id=" + bot.attributes.bot_id;

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
        // console.log(bots);
      } catch (error) {
        gettingDatas = false;
        console.log(error);
      }
    } else {
      console.log("token not found");
    }
  }

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

  function nextPage() {
    if (page < pagination.pageCount) {
      if (searchTerm !== "") {
        serachPagination = true;
      } else {
        serachPagination = false;
      }
      page = page + 1;
      getBots();
    }
  }

  function previousPage() {
    if (page > 1) {
      if (searchTerm !== "") {
        serachPagination = true;
      } else {
        serachPagination = false;
      }
      page = page - 1;
      getBots();
    }
  }

  let dropdownOpen = $state([]);

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    for (let i = 0; i < dropdownOpen.length; i++) {
      dropdownOpen[i] = false;
    }
  };

  async function deleteBot(id) {
    let token = localStorage.getItem("token");

    if (token) {
      triggerLoading(true);
      try {
        const response = await $call.call("client-dashboard", "delete-bot", {
          URLQuery: {},
          URLParams: {
            id: id.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          console.log(response);
          return;
        }

        tmpId = null;
        popupModal = false;
        triggerLoading(false);
        triggerToast("deleted");
        getBots();
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function startReading(botId) {
    let token = localStorage.getItem("token");

    if (token) {
      try {
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
          triggerToast("error");
          console.log(response);
          return;
        }

        // // await updateBot(botId, "Reading");

        // // bots = bots.map((bot) => {
        // //   if (bot.botsID === botId) {
        // //     bot.deployStatus = "Reading";
        // //   }
        // //   return bot;
        // // });

        //demo
        // const response = await $call.call("hexa-ai-rag", "add-document-by-bot-id", {
        //   URLQuery: {},
        //   URLParams: {
        //     id: botId.toString(),
        //   },
        //   Params: {},
        //   Headers: {
        //     Authorization: "Bearer " + token,
        //   },
        // });

        // if (response.status !== 200) {
        //   triggerToast("error");
        //   console.log(response);
        //   return;
        // }

        // await updateBot(botId, "Reading");

        // bots = bots.map((bot) => {
        //   if (bot.botsID === botId) {
        //     bot.deployStatus = "Reading";
        //   }
        //   return bot;
        // });

        // triggerToast("success");
        setTimeout(getBots, 3000);
      } catch (error) {
        console.log(error);
      }
    }
  }

  onMount(() => {
    getBots();
  });

  $effect(() => {
    if ($onProgress) {
      console.log("notification", $notifications);
      if ($notifications.length > 0 && $notifications[0].attributes.content.includes("100")) {
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
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Typography
          preset="Display sm/Semibold"
          text={$_("recentDeployedBots")}
          color="secondary"
        />

        <!-- Bots List Data -->
        <div>
          <Search
            class="bg-primary-light dark:bg-primary-dark mb-2 md:mb-8 rounded-lg"
            placeholder={$_("search")}
            bind:value={searchTerm}
            on:input={getBots}
          />

          {#if gettingDatas}
            <div class="flex flex-col gap-4 w-full">
              <StateCard
                illustration={IllustrationProgress}
                description={$_("fetchingDataPleaseWait")}
              />
            </div>
          {:else if bots.length === 0}
            <StateCard
              illustration={IllustrationEmpty}
              description={searchTerm !== "" ? $_("noDataFound") : $_("youDontHaveAnyDataYet")}
            />
          {:else}
            <div
              class={`grid ${$showAI || $showNotif ? "grid-cols-1 xl:grid-cols-3" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4"} gap-4`}
            >
              {#each bots as bot, index (bot.botsID)}
                <BotCard
                  {bot}
                  {index}
                  {dropdownOpen}
                  bind:popupModal
                  bind:tmpId
                  config={{ startReading: true, detail: true, detailView: true, configure: true }}
                  {handleDropdownFocusLoss}
                  {startReading}
                />
              {/each}
            </div>
            <Pagination {pagination} {previousPage} {nextPage} />
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>

<div class="flex justify-center">
  <Toast bind:open>{$_("copiedToClipboard")}</Toast>
</div>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={() => deleteBot(tmpId)} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>
