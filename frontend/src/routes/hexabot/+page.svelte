<script>
  import { onMount } from "svelte";
  import {
    P,
    A,
    Card,
    Modal,
    Button,
    Avatar,
    Heading,
    // Dropdown,
    // DropdownItem,
  } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { DotsVerticalOutline, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import MessagePlus from "$lib/icon/message-plus-square.svelte";
  // import MessageChatSquare from "$lib/icon/message-chat-square.svelte";
  // import FileAttachment_02 from "$lib/icon/file-attachment-02.svelte";
  // import PencilLine from "$lib/icon/pencil-line.svelte";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import Eye from "$lib/icon/eye.svelte";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { Typography, SectionHeader, ActionCard, ChartCard, call, triggerToast, StateCard } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import ArrowUpRightIcon from "$lib/icon/arrow-up-right.svelte";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  import { goto } from "$app/navigation";
  // import { PUBLIC_STRAPI_URL_DEV, PUBLIC_STRAPI_URL_PROD } from "$env/static/public";

  // const isDev = process.env.NODE_ENV === "development";

  // const strapiUrl = isDev ? PUBLIC_STRAPI_URL_DEV : PUBLIC_STRAPI_URL_PROD;

  // export let data;
  let data = {
    admin: false,
  };

  let actionItems = $state([]);

  let chartAction = $state([]);
  let dataChartCollections = $state({});
  let monthlyCollectionReadyCounts = $state([]);
  let tmpMonthlyCollectionReadyCounts = $state([]);

  let chartReady = $state(false);

  let Header;

  let bots = $state([]);

  let tmpId = $state(null);

  let popupModal = $state(false);

  let deployStatusModal = $state(false);

  let gettingDatas = $state(true);

  const classStatusType = (type) => {
    if (type === "Done") {
      return "bg-hc-success-50 border-hc-success-200 text-hc-success-700 dark:bg-hc-success-900 dark:border-hc-success-700 dark:text-hc-success-300";
    } else if (type === "Not Ready") {
      return "bg-hc-error-50 border-hc-error-200 text-hc-error-700 dark:bg-hc-error-900 dark:border-hc-error-700 dark:text-hc-error-300";
    } else if (type === "On Progress") {
      return "bg-hc-warning-50 border-hc-warning-200 text-hc-warning-700 dark:bg-hc-warning-900 dark:border-hc-warning-700 dark:text-hc-warning-300";
    } else if (type === "WhatsApp") {
      return "text-utility-success-700-light dark:text-utility-success-700-dark bg-utility-success-50-light dark:bg-utility-success-50-dark border-utility-success-200-light dark:border-utility-success-200-dark";
    } else if (type === "Web") {
      return "text-sky-700 dark:text-sky-300 bg-sky-50 dark:bg-sky-950 border-sky-200 dark:border-sky-800";
    }
  };

  async function getActionItems() {
    actionItems = [
      {
        href: await $call.getPage("client-dashboard", "hexabot-create"),
        icon: "messagePlus",
        actionName: "Create New Bots",
        desc: "desc Create New Bots",
      },
    ];
  }

  function generateChartData() {
    const currentDate = new Date();
    const twelveMonthsBefore = new Date();
    twelveMonthsBefore.setMonth(currentDate.getMonth() - 12);

    const formatDate = (date) => {
      const options = { month: "short", year: "numeric" };
      return new Intl.DateTimeFormat("en-US", options).format(date);
    };

    dataChartCollections = {
      chart: {
        type: "area",
        fontFamily: "Inter, sans-serif",
        dropShadow: {
          enabled: false,
        },
        toolbar: {
          show: false,
        },
      },
      tooltip: {
        enabled: true,
        x: {
          show: false,
        },
      },
      fill: {
        type: "gradient",
        gradient: {
          opacityFrom: 0.55,
          opacityTo: 0,
          shade: "#17B169",
          gradientToColors: ["#17B169"],
        },
      },
      dataLabels: {
        enabled: false,
      },
      stroke: {
        width: 2,
      },
      grid: {
        show: false,
        strokeDashArray: 4,
        padding: {
          left: 2,
          right: 2,
          top: 0,
        },
      },
      series: [
        {
          name: "Documents collection",
          data: monthlyCollectionReadyCounts,
          color: "#17B169",
        },
      ],
      xaxis: {
        categories: Array.from({ length: 12 }, (_, i) => {
          const date = new Date(twelveMonthsBefore);
          date.setMonth(twelveMonthsBefore.getMonth() + i + 1);
          return formatDate(date);
        }),
        labels: {
          show: false,
        },
        axisBorder: {
          show: false,
        },
        axisTicks: {
          show: false,
        },
      },
      yaxis: {
        show: false,
      },
    };

    chartReady = true;
  }

  async function getCountDocumentsPerMonth() {
    const today = new Date();
    let earlyDateCurrentMonth = new Date(today.getFullYear(), today.getMonth(), 1);
    let lastDateCurrentMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0);
    let token = localStorage.getItem("token");

    if (token !== null) {
      for (let i = 0; i < 12; i++) {
        let response = await $call.call("client-dashboard", "get-count-documents-per-month", {
          URLQuery: {},
          URLParams: {},
          Params: {
            earlyDate: earlyDateCurrentMonth.toISOString().split("T")[0],
            lastDate: lastDateCurrentMonth.toISOString().split("T")[0],
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          console.log(response.status);
          console.log(response);
          return;
        }
        // console.log(
        //   "telor",
        //   format(earlyDateCurrentMonth, "yyyy-MM-dd"),
        //   format(lastDateCurrentMonth, "yyyy-MM-dd"),
        // );

        tmpMonthlyCollectionReadyCounts.push(response.result.general.meta.pagination.total);

        earlyDateCurrentMonth = new Date(
          earlyDateCurrentMonth.getFullYear(),
          earlyDateCurrentMonth.getMonth() - 1,
          1,
        );
        lastDateCurrentMonth = new Date(
          lastDateCurrentMonth.getFullYear(),
          lastDateCurrentMonth.getMonth(),
          0,
        );
      }

      tmpMonthlyCollectionReadyCounts.reverse();

      monthlyCollectionReadyCounts = tmpMonthlyCollectionReadyCounts;

      let readyCollection = groupingDataChart(monthlyCollectionReadyCounts);

      getChartAction(readyCollection);
      generateChartData();
    } else {
      console.log("no token");
    }
  }

  function groupingDataChart(data) {
    let total = 0;
    for (let i = 0; i < data.length; i++) {
      total += data[i];
    }

    let change = data[data.length - 1] - data[data.length - 2];

    let changeType = change > 0 ? "increase" : change < 0 ? "decrease" : "-";

    return {
      total: total,
      change: change,
      changeType: changeType,
    };
  }

  async function getChartAction(readyCollection) {
    chartAction = [
      {
        href: await $call.getPage("client-dashboard", "document-collection"),
        dataName: "Document Collection",
        value: readyCollection.total,
        change: readyCollection.change,
        changeType: readyCollection.changeType,
      },
    ];
  }

  async function getBots() {
    let token = localStorage.getItem("token");

    if (token) {
      try {
        const response = await $call.call("client-dashboard", "get-bots", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: 1,
            page_size: 8,
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
            // deployStatus: bot.attributes.deploy_status,
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

  async function goToEditPage(botId) {
    goto(
      (await $call.getPage("client-dashboard", "hexabot-detail")) +
        "?botId=" +
        botId +
        "&edit=true",
    );
  }

  onMount(() => {
    getActionItems();
    getCountDocumentsPerMonth();
    getBots();
  });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Typography preset="Display sm/Semibold" text={$_("hexaBots")} propClass="fluid-xl" />

        <!-- Action Items -->
        <div class="flex flex-col gap-4">
          <SectionHeader title={$_("action")} />
          <ActionCard {actionItems} />
        </div>

        <!-- Reports Monitoring -->
        <div class="flex flex-col gap-4">
          <SectionHeader title={$_("reports")} />
          <ChartCard {chartAction} chartList={[dataChartCollections]} />
        </div>

        <!-- Document Table Data -->
        <div>
          <div
            class={`flex items-center justify-between pb-4 mb-4 border-b border-hc-graylight-200 dark:border-hc-graydark-800`}
          >
            <Heading
              tag="h6"
              class="font-jakarta w-fit text-hc-graylight-900 dark:text-hc-graydark-50"
              >{$_("recentDeployedBots")}</Heading
            >
            <Button
              on:click={async () => {
                goto(await $call.getPage("client-dashboard", "hexabot-list"));
              }}
              color="light"
              class="bg-hc-graylight-50 dark:bg-hc-graydark-900 border-hc-graylight-200 dark:border-hc-graydark-800 hover:bg-hc-graylight-100 hover:dark:bg-hc-graydark-800 focus-within:ring-0"
              >{$_("viewAll")}</Button
            >
          </div>

          {#if gettingDatas}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else if bots.length === 0}
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          {:else}
            <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4">
              {#each bots as bot, index (bot.botsID)}
                <!-- {#if bot.deployStatus === "Done"} -->
                <Card padding="none" class="shadow-none grow w-full">
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
                      <div onfocusout={handleDropdownFocusLoss} class=" flex items-center relative">
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
                          class={`z-10 ${dropdownOpen[index] ? "absolute top-9 right-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
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
                              }}>{$_("delete")}</button
                            >
                          </ul>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="flex gap-3 px-4 pt-4 pb-6">
                    <div class="grid grow grid-cols-2 gap-3">
                      <Button
                        on:click={async () => {
                          goto(
                            (await $call.getPage("client-dashboard", "hexabot-detail")) +
                              "?botId=" +
                              bot.botsID,
                          );
                        }}
                        outline
                        class="!flex grow hover:bg-primary-light dark:hover:bg-primary-dark rounded-lg"
                      >
                        <Typography
                          preset="Text sm/Semibold"
                          text={$_("detail")}
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
                <!-- <div class="span-col-1">
                  <Card padding="none" size="md">
                    <div
                      class="bg-active-light dark:bg-active-dark border-b border-secondary-light dark:border-secondary-dark px-4 py-3"
                    >
                      <div class="flex justify-between items-center">
                        <div class="flex flex-col">
                          <div class="flex items-center gap-2">
                            <Heading
                              tag="h6"
                              class="!font-semibold text-primary-light dark:text-primary-dark"
                              >{`${bot.botsName}`}</Heading
                            >
                            <div
                              class={`px-2 py-0.5 rounded-full border justify-start items-center inline-flex text-center text-xs font-medium leading-[18px]
                            ${classStatusType(bot.deployedPlatform)}`}
                            >
                              {bot.deployedPlatform}
                            </div>
                          </div>
                          <P size="sm" class="text-secondary-light dark:text-secondary-dark"
                            >{`ID: ${bot.botsID}`}</P
                          >
                        </div>
                        <div class="flex items-center pl-2 gap-2">
                          <A
                            href="hexabot/detail"
                            class="text-fg-brand-primary-light dark:text-fg-brand-primary-dark text-sm font-semibold"
                            >{$_("view")}</A
                          >
                          <div class="cursor-pointer flex items-center">
                            <Button class="focus-within:ring-0 !p-0" color="none" size="xs"
                              ><DotsVerticalOutline /></Button
                            >
                          </div>
                        </div>
                      </div>
                    </div>
                    <div
                      class="pb-4 bg-primary-light dark:bg-primary-dark border-b border-secondary-light dark:border-secondary-dark p-4"
                    >
                      <div class="flex justify-between items-center">
                        <div>
                          <P
                            size="sm"
                            weight="semibold"
                            class="text-primary-light dark:text-primary-dark">{$_("clientId")}</P
                          >
                          <P
                            size="sm"
                            weight="normal"
                            class="text-secondary-light dark:text-secondary-dark line-clamp-1"
                            >{bot.clientID}</P
                          >
                        </div>
                        <div class="flex pl-2">
                          <div class="cursor-pointer">
                            <CopyIcon
                              svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                            />
                          </div>
                        </div>
                      </div>
                    </div>
                    <div
                      class="pb-4 bg-primary-light dark:bg-primary-dark border-b border-secondary-light dark:border-secondary-dark p-4"
                    >
                      <div class="flex justify-between items-center">
                        <div>
                          <P
                            size="sm"
                            weight="semibold"
                            class="text-primary-light dark:text-primary-dark"
                            >{$_("activation")}</P
                          >
                          <P
                            size="sm"
                            weight="normal"
                            class="text-secondary-light dark:text-secondary-dark line-clamp-1"
                            >{bot.activationTime}</P
                          >
                        </div>
                        <div class="flex pl-2">
                          <div class="cursor-pointer">
                            <CopyIcon
                              svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                            />
                          </div>
                        </div>
                      </div>
                    </div>
                    <div
                      class="pb-4 bg-primary-light dark:bg-primary-dark border-b border-secondary-light dark:border-secondary-dark p-4"
                    >
                      <div class="flex justify-between items-center">
                        <div>
                          <P
                            size="sm"
                            weight="semibold"
                            class="text-primary-light dark:text-primary-dark">{$_("urlBot")}</P
                          >
                          <P
                            size="sm"
                            weight="normal"
                            class="text-secondary-light dark:text-secondary-dark line-clamp-1"
                            >{bot.urlBot}</P
                          >
                        </div>
                        <div class="flex pl-2">
                          <div class="cursor-pointer">
                            <CopyIcon
                              svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                            />
                          </div>
                        </div>
                      </div>
                    </div>
                  </Card>
                </div> -->
                <!-- {/if} -->
              {/each}
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>

<Modal bind:open={deployStatusModal} autoclose>
  <div class="space-y-2">
    <Heading tag="h5">{$_("botDeploymentStatus")}</Heading>
    <P>{$_("descBotDeploymentStatus")}</P>
  </div>
  {#each bots as bot}
    {#if bot.deployStatus !== "Done"}
      <Card size="lg">
        <div class="flex justify-between items-center">
          <div>
            <Heading tag="h6">{`${bot.botsName}`}</Heading>
            <P>{`${bot.timestamp}`}</P>
          </div>
          <Button
            class={`p-[2px] border-[1px] h-1/2 rounded-xl text-center ${classStatusType(
              bot.deployStatus,
            )}`}>{bot.deployStatus}</Button
          >
        </div>
      </Card>
    {/if}
  {/each}
  <Button
    on:click={() => {
      deployStatusModal = false;
    }}
    class="w-full"
    color="light">{$_("close")}</Button
  >
</Modal>

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
