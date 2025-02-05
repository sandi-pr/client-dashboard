<script>
  import {
    P,
    Card,
    Button,
    Heading,
    Avatar,
    Toast,
    Pagination,
    PaginationItem,
    Modal,
  } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import ArrowUpRightIcon from "$lib/icon/arrow-up-right.svelte";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  import {
    Typography,
    call,
    triggerToast,
    StateCard,
    newNotification,
    notifications,
    Badge,
    showAI, 
    showNotif,
  } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { check } from "$lib";

  let searchTerm = $state("");
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

    if (token) {
      try {
        const response = await $call.call("client-dashboard", "get-bots", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 18,
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
        console.log(bots);
      } catch (error) {
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
      page = page + 1;
      getBots();
    }
  }

  function previousPage() {
    if (page > 1) {
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
          triggerToast("error");
          console.log(response);
          return;
        }

        tmpId = null;
        popupModal = false;
        triggerToast("deleted");
        getBots();
      } catch (error) {
        console.log(error);
      }
    }
  }

  onMount(() => {
    getBots();
  });

  $effect(() => {
    if ($newNotification) {
      if ($notifications[0].attributes.type === "hexabot") {
        getBots();
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
        <Typography preset="Display sm/Semibold" text={$_("deployedBots")} color="secondary" />

        <!-- Bots List Data -->
        <div>
          <!-- <Search
            class="bg-white dark:bg-hc-graydark-950 mb-6 rounded-lg"
            placeholder={$_("searchBotName")}
            bind:value={searchTerm}
          /> -->

          {#if gettingDatas}
            <div class="flex flex-col gap-4 w-full">
              <StateCard
                illustration={IllustrationProgress}
                description={$_("fetchingDataPleaseWait")}
              />
            </div>
          {:else if bots.length === 0}
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          {:else}
            <div class={`grid ${$showAI || $showNotif ? "grid-cols-1 xl:grid-cols-3" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4"} gap-4`}>
              {#each bots as bot, index (bot.botsID)}
                <!-- {#if bot.deployStatus === "Done"} -->
                <Card padding="none" class="relative shadow-none grow w-full">
                  {#if bot.deployStatus === "Reading"}
                    <div class="absolute top-2 right-2">
                      <Badge type="warning" text={bot.deployStatus} />
                    </div>
                  {:else if bot.deployStatus === "Ready"}
                    <div class="absolute top-2 right-2">
                      <Badge type="success" text={bot.deployStatus} />
                    </div>
                  {/if}
                  <div
                    class="bg-active-light dark:bg-active-dark border-b border-secondary-light dark:border-secondary-dark px-4 py-6"
                  >
                    <div class="flex justify-between items-center">
                      <div class="flex items-center gap-3">
                        <Avatar src={bot.bot_profile ? bot.bot_profile : ""} />
                        <div class="flex flex-col">
                          <Typography
                            preset="Text lg/Semibold"
                            color="primary"
                            text={bot.botsName}
                          />
                          <Typography
                            preset="Text sm/Regular"
                            color="secondary"
                            text={bot.botDescription}
                          />
                        </div>
                      </div>
                      {#if check("aiManagement.hexabots:delete")}
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
                            class={`z-10 ${dropdownOpen[index] ? "absolute top-10" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
                          >
                            <ul class="py-2 text-sm text-gray-700 dark:text-gray-200">
                              <button
                                class="flex p-2"
                                onclick={() => {
                                  popupModal = true;
                                  tmpId = bot.botsID;
                                  dropdownOpen[index] = !dropdownOpen[index];
                                }}>{$_("delete")}</button
                              >
                            </ul>
                          </div>
                        </div>
                      {/if}
                    </div>
                  </div>
                  <div class="flex gap-3 px-4 pt-4 pb-6">
                    <div class="grid grow grid-cols-2 gap-3">
                      <Button
                        disabled={!check("aiManagement.hexabots:view")}
                        on:click={async () => {
                          goto(
                            (await $call.getPage("client-dashboard", "bots-detail")) +
                              "?botId=" +
                              bot.botsID,
                          );
                        }}
                        outline
                        class="!flex grow hover:bg-primary-light dark:hover:bg-primary-dark rounded-lg"
                      >
                        <Typography
                          preset="Text sm/Semibold"
                          text={$_("view")}
                          color="brand-secondary"
                        />
                      </Button>
                      <Button
                        on:click={() => {
                          window.open(bot.urlBot, "_blank");
                        }}
                        outline
                        class="!flex grow hover:bg-primary-light dark:hover:bg-primary-dark rounded-lg border-secondary-light dark:border-secondary-dark"
                      >
                        <Typography
                          preset="Text sm/Semibold"
                          text={$_("webchat")}
                          color="secondary"
                        />
                        <ArrowUpRightIcon
                          svgClass="w-5 h-5 text-secondary-light dark:text-secondary-dark"
                        />
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
              {/each}
            </div>
            <div class="flex justify-between">
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

              <div class="m-4">
                <Pagination
                  on:next={nextPage}
                  on:previous={previousPage}
                  activeClass="bg-hc-graylight-100 dark:bg-hc-graydark-700 dark:border-[#374151] dark:text-hc-graydark-300"
                >
                  <span slot="prev">{$_("previous")}</span>
                  <span slot="next">{$_("next")}</span>
                </Pagination>
              </div>
            </div>
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
