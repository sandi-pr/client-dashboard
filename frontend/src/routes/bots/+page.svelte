<script>
  import { onMount } from "svelte";
  import {
    P,
    Modal,
    Button,
    Avatar,
    Heading,
    // Dropdown,
    // DropdownItem,
  } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import {
    Typography,
    SectionHeader,
    ActionCard,
    ChartCard,
    call,
    triggerToast,
    StateCard,
    check,
    Badge,
    notifications,
    onProgress,
    showAI,
    showNotif,
    BotCard,
    triggerLoading,
  } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import Stars03Icon from "$lib/icon/stars-03.svelte";
  import { goto } from "$app/navigation";
  import { format } from "date-fns";
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

  let bots = $state([]);

  let tmpId = $state(null);

  let popupModal = $state(false);

  let deployStatusModal = $state(false);

  let gettingDatas = $state(true);

  async function getActionItems() {
    actionItems = [
      {
        disabled: !check("aiManagement.hexabots:create"),
        href: await $call.getPage("client-dashboard", "bots-create"),
        icon: "messagePlus",
        actionName: "Create New Bots",
        desc: "desc Create New Bots",
      },
      {
        disabled: !check("aiManagement.hexabots:deploy"),
        onclick: () => {
          deployStatusModal = true;
        },
        icon: "eye",
        actionName: "View bot deployment status",
        desc: "desc View bot deployment status",
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
          name: "New bot",
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
    let lastDateCurrentMonth = new Date(today.getFullYear(), today.getMonth() + 1, 1);
    let token = localStorage.getItem("token");

    if (token !== null) {
      for (let i = 0; i < 12; i++) {
        let response = await $call.call("client-dashboard", "get-count-bot-per-mount", {
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

        tmpMonthlyCollectionReadyCounts.push(response.result.meta.pagination.total);

        earlyDateCurrentMonth = new Date(
          earlyDateCurrentMonth.getFullYear(),
          earlyDateCurrentMonth.getMonth() - 1,
          1,
        );
        lastDateCurrentMonth = new Date(
          lastDateCurrentMonth.getFullYear(),
          lastDateCurrentMonth.getMonth() - 1,
          1,
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

    let change = data[data.length - 1];
    // let change = data[data.length - 1] - data[data.length - 2];

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
        href: await $call.getPage("client-dashboard", "bots-list"),
        dataName: "Total Bot",
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

  async function goToEditPage(botId) {
    goto(
      (await $call.getPage("client-dashboard", "bots-detail")) + "?botId=" + botId + "&edit=true",
    );
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

  async function updateBot(botId, status) {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        let response = await $call.call("client-dashboard", "update-bot-status", {
          URLQuery: {},
          URLParams: {
            id: botId.toString(),
          },
          Params: {
            status: status,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          triggerToast("error");
          return;
        }

        console.log("update bot", response);
      } catch (error) {
        console.log(error);
      }
    }
  }

  let filteredBots = $state([]);

  $effect(() => {
    filteredBots = bots.filter(
      (bot) => bot.deployStatus === "Inactive" || bot.deployStatus === "Reading",
    );
  });

  onMount(() => {
    getActionItems();
    getCountDocumentsPerMonth();
    getBots();
  });

  let gettingStatusBot = $state(false);
  let timerGetStatusBot;

  $effect(() => {
    if ($onProgress) {
      console.log("notification", $notifications);
      if ($notifications.length > 0 && $notifications[0].attributes.meta?.progress.includes("100")) {
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
          text={$_("hexabotsDeployment")}
          propClass="fluid-xl"
        />

        <!-- Action Items -->
        <div class="flex flex-col gap-4">
          {#if actionItems.some((item) => item.disabled !== true)}
            <SectionHeader title={$_("action")} />
          {/if}
          <ActionCard {actionItems} />
        </div>

        <!-- Reports Monitoring -->
        <!-- <div class="flex flex-col gap-4">
          <SectionHeader title={$_("reports")} />
          <ChartCard {chartAction} chartList={[dataChartCollections]} />
        </div> -->

        <!-- Document Table Data -->
        <div>
          <div
            class={`flex items-center justify-between pb-4 mb-4 border-b border-hc-graylight-200 dark:border-hc-graydark-800`}
          >
            <Heading
              tag="h6"
              class="font-jakarta w-fit text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {$_("recentDeployedBots")}
            </Heading>
            <Button
              on:click={async () => {
                goto(await $call.getPage("client-dashboard", "bots-list"));
              }}
              color="light"
              class="bg-hc-graylight-50 dark:bg-hc-graydark-900 border-hc-graylight-200 dark:border-hc-graydark-800 hover:bg-hc-graylight-100 hover:dark:bg-hc-graydark-800 focus-within:ring-0"
            >
              {$_("viewAll")}
            </Button>
          </div>

          {#if gettingDatas}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else if bots.length === 0}
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
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
                  config={{ startReading: true, delete: true, detail: true, detailView: true, configure: true }}
                  {handleDropdownFocusLoss}
                  {startReading}
                />
              {/each}
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>

<Modal bind:open={deployStatusModal} outsideclose>
  <svelte:fragment slot="header">
    <div class="">
      <Heading tag="h5">{$_("botDeploymentStatus")}</Heading>
      <P>{$_("descBotDeploymentStatus")}</P>
    </div>
  </svelte:fragment>
  {#if gettingDatas}
    <StateCard illustration={IllustrationProgress} description={$_("fetchingDataPleaseWait")} />
  {:else if bots.length !== 0}
    <div
      class={`${filteredBots.length > 0 ? "py-4 border rounded-lg space-y-4 border-secondary-light dark:border-secondary-dark" : ""} `}
    >
      {#if filteredBots.length > 0}
        {#each filteredBots as bot, index}
          <div class="flex px-4 justify-between items-center">
            <div class="flex items-center gap-3">
              <Avatar src={bot.bot_profile ? bot.bot_profile : ""} />
              <div class="flex flex-col">
                <div class="flex items-center gap-1">
                  <Typography preset="Text lg/Semibold" color="primary" text={bot.botsName} />
                  <Typography
                    preset="Text xs/Regular"
                    color="quarterary"
                    text={format(new Date(bot.timestamp), "HH:mm, dd/MM/yyyy")}
                  />
                </div>
                <Typography preset="Text sm/Regular" color="secondary" text={bot.botDescription} />
              </div>
            </div>
            {#if bot.deployStatus === "Inactive"}
              <Button
                on:click={() => {
                  startReading(bot.botsID);
                }}
                class="p-2 md:px-3 md:py-2 gap-2"
              >
                <Stars03Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                <div class="hidden md:block">{$_("startReading")}</div>
              </Button>
            {:else if bot.deployStatus === "Reading"}
              <Badge type="warning" text={bot.deployStatus} />
            {/if}
          </div>
          {#if index !== filteredBots.length - 1}
            <Divider />
          {/if}
        {/each}
      {:else}
        <StateCard
          illustration={IllustrationEmpty}
          description={$_("YouDontHaveAnyBotsDeploymentOngoing")}
        />
      {/if}
    </div>
  {:else}
    <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
  {/if}
  <svelte:fragment slot="footer">
    <Button
      on:click={() => {
        deployStatusModal = false;
      }}
      class="w-full"
      color="light">{$_("close")}</Button
    >
  </svelte:fragment>
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
