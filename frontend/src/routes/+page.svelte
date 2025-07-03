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
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";
  import { Alert, CloseButton, Button as ButtonFlowbite } from "flowbite-svelte";
  import { Card, Icon, Divider, StateCard, ChartCard, Typography, CardHome } from "$lib";
  import { call, showAI, showNotif, isOpenSidebarStore } from "$lib";
  import IllustrationFirst from "$lib/icon/illustration_first_time.svelte";

  const isFrontend = process.env.NODE_ENV === "development";

  let userData = $state({ username: "", email: "", brand: "" });

  let isFirstTime = $state(true);
  let chartActionAi = $state([]);
  let dataChartAllAi = $state([]);
  let hiddenContent = $state([]);
  let cardsArrays = $state([]);
  let isFirstTimeTimeoutId;

  const cardsAssistantActions = [
    {
      icon: "var95",
      title: "assistant",
      description: "descAssistant",
      page: isFrontend ? "/assistant" : "assistant",
    },
  ];

  const cardsDocumentActions = [
    {
      icon: "var91",
      title: "documents",
      description: "descDocumentLibrary",
      page: isFrontend ? "/documents" : "document-collection-overview",
    },
    // {
    //   icon: "graduationHat02",
    //   title: "eLearning",
    //   description: "descELearning",
    //   page: isFrontend ? "/e_learning" : "e-learning",
    // },
    // {
    //   icon: "data",
    //   title: "apiConnector",
    //   description: "descApiConnector",
    //   page: isFrontend ? "api_connector" : "api-connector",
    // },
  ];

  const cardsAnalyticsActions = [
    {
      icon: "var86",
      title: "dataSources",
      description: "descDataSources",
      page: isFrontend ? "/data_sources" : "data-sources",
    },
    {
      // BRI ADJUSTMENT
      icon: "var87",
      title: "insights",
      description: "descInsights",
      page: isFrontend ? "data_sources" : "csv-insight",
    },
  ];

  const cardsHandwritingActions = [
    {
      icon: "var88",
      title: "handwriting",
      description: "descHandwriting",
      page: isFrontend ? "/handwriting" : "handwriting",
    },
  ];

  const cardsVoiceActions = [
    {
      icon: "var92",
      title: "tryModels",
      description: "descTryModels",
      page: isFrontend ? "/llms" : "llms",
    },
  ];

  const cardsLlmsActions = [
    {
      icon: "var94",
      title: "tryModels",
      description: "descTryModels",
      page: isFrontend ? "/llms" : "llms",
    },
  ];

  const cardsConfigurationActions = [
    {
      icon: "var99",
      title: "createBots",
      description: "descCreateBots",
      page: isFrontend ? "/bots" : "bots",
    },
    {
      icon: "var76",
      title: "finetunings",
      description: "descFinetunings",
      page: isFrontend ? "/finetunings" : "finetunings",
    },
  ];

  function setIsFirstTime() {
    isFirstTime = false;
    localStorage.setItem("isFirstTime", isFirstTime);
  }

  function getIsFirstTime() {
    if (typeof localStorage !== "undefined") {
      isFirstTime = localStorage.getItem("isFirstTime") === "false" ? false : true;
      // console.log("isFirstTime", isFirstTime);

      if (isFirstTime) {
        isFirstTimeTimeoutId = setTimeout(() => {
          isFirstTime = false;
          // localStorage.setItem("isFirstTime", isFirstTime);
        }, 6000);
      }
    } else {
      console.log("localstorage is undefined");
    }
  }

  async function CheckSetup() {
    try {
      // var resp_setup = await $call.call("client-login", "is-need-setup", {
      //   URLQuery: {},
      //   URLParams: {},
      //   Params: {},
      //   Headers: {},
      // });
      // console.log(resp_setup);
      // if (resp_setup.status === 200) {
      //   if (resp_setup.result.is_need_setup === true) {
      //     if (localStorage.getItem("token") !== null) localStorage.removeItem("token");
      //     if (localStorage.getItem("user") !== null) localStorage.removeItem("user");
      //     await $call.goToPage("hexapod-setup", "home-page");
      //     return;
      //   }
      // }
      var token = localStorage.getItem("token");
      if (token !== null) {
        var resp = await $call.call("client-login", "is-token-valid", {
          URLQuery: {},
          URLParams: {},
          Params: {
            token: token,
          },
          Headers: {},
        });

        if (resp.status === 200) {
          if (resp.result.confirmed === true && resp.result.blocked === false) {
            localStorage.setItem("user", JSON.stringify(resp.result));
            userData = resp.result;
          }
        }
      } else {
        await $call.goToPage("client-login", "home-page");
        return;
      }
    } catch (error) {
      // console.log(error);
      userData = { username: "Hexa Staff", email: "hexacode@hexacode.org", brand: "hexacode" };
    }
  }

  const navigateToPage = async (page, param = "") => {
    if (page === "bakrie-webchat") {
      goto(isFrontend ? "/" : await $call.getPage("bakrie", "home-page"));
    } else {
      goto((isFrontend ? page : await $call.getPage("client-dashboard", page)) + param);
    }
  };

  const cardsNotEmpty = (cardsArrays) => {
    // console.log("cardsArrays", cardsArrays);
    let counter = 0;
    for (let cardArray of cardsArrays) {
      if (cardArray.array.length > 0) {
        counter++;
      }
    }
    return counter;
  };

  onMount(async () => {
    // triggerLoading(true);
    // setTimeout(async () => {
    //   triggerLoading(false);
    // }, 3000);
    // await getCompanyType();
    // getActionItems();
    // getCountPerModulePerMonth();
    cardsArrays = [
      { name: "assistant", array: cardsAssistantActions },
      { name: "documents", array: cardsDocumentActions },
      { name: "analytics", array: cardsAnalyticsActions },
      { name: "handwriting", array: cardsHandwritingActions },
      { name: "voice", array: cardsVoiceActions },
      { name: "llms", array: cardsLlmsActions },
      { name: "configurations", array: cardsConfigurationActions },
    ];

    let totalCardsSection = cardsNotEmpty(cardsArrays);
    hiddenContent = Array.from({ length: totalCardsSection }, () => true);
    getIsFirstTime();
    CheckSetup();
  });
</script>

<div class="flex flex-col gap-8 px-4 md:py-6 mb-20">
  {#if isFirstTime}
    <Alert dismissable on:close={setIsFirstTime} class="relative flex justify-center">
      <CloseButton
        slot="close-button"
        let:close
        class="absolute top-2 right-2 ms-1.5 dark:hover:bg-gray-700"
        color="primary"
        on:click={close}
      />
      <StateCard illustration={IllustrationFirst} />
      <div class="flex flex-col justify-center items-center gap-1 pt-2">
        <Typography
          preset="Display sm/Semibold"
          text={$_("welcometoHexapod") + sessionStorage.getItem("sidebar_title") || "HexaCode.AI"}
          color="primary"
        />
      </div>
    </Alert>
  {:else}
    <div class="flex flex-col gap-1">
      <Typography
        preset="Display sm/Semibold"
        text={`${$_("welcometoHexapod")} ${sessionStorage.getItem("sidebar_title") || "HexaCode.AI"}`}
        color="primary"
      />
      <Typography text="{userData.username} - {userData.email}" color="tertiary" />
    </div>
  {/if}

  {#each cardsArrays as cardArray, index}
    {#if cardArray.array.length > 0}
      <div class="flex flex-col gap-4">
        <ButtonFlowbite
          class="p-0 pr-4 justify-between border-none focus:outline-none focus:ring-0"
          color="none"
          outline
          on:click={() => (hiddenContent[index] = !hiddenContent[index])}
        >
          <Typography text={$_(cardArray.name)} preset="Text lg/Semibold" />
          {#if hiddenContent[index]}
            <Icon name="chevronUp" />
          {:else}
            <Icon name="chevronDown" />
          {/if}
        </ButtonFlowbite>
        <Divider />
        {#if hiddenContent[index]}
          {#if cardArray.name === "aiAnalyzer"}
            {#if chartActionAi.length > 0}
              <ChartCard isLibrary chartAction={chartActionAi} chartList={dataChartAllAi} />
            {/if}
            <div
              class={`grid ${$showAI || $showNotif ? "grid-cols-1 xl:grid-cols-2" : "grid-cols-1 md:grid-cols-2"} gap-4`}
            >
              {#each cardArray.array as card}
                {#if !chartActionAi.find((item) => item.title === card.title)}
                  <Card
                    icon={card.icon}
                    title={$_(card.title)}
                    description={$_(card.description)}
                    onclick={() => navigateToPage(card.page, card.param)}
                    comingSoon={card?.comingSoon}
                  />
                {/if}
              {/each}
            </div>
          {:else}
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {#each cardArray.array as card}
                <CardHome
                  icon={card.icon}
                  title={$_(card.title)}
                  description={$_(card.description)}
                  onclick={() => {
                    // const isOpen = JSON.parse(sessionStorage.getItem("isOpen")) || {};
                    // if (cardArray.name in isOpen) {
                    //   isOpen[cardArray.name] = true;
                    //   sessionStorage.setItem("isOpen", JSON.stringify(isOpen));
                    // }
                    // // console.log(cardArray.name);
                    // // console.log(isOpen);
                    isOpenSidebarStore.update((isOpen) => {
                      if (cardArray.name in isOpen) {
                        isOpen[cardArray.name] = true;
                      }
                      return isOpen;
                    });
                    navigateToPage(card.page);
                  }}
                  comingSoon={card?.comingSoon}
                />
              {/each}
            </div>
          {/if}
        {/if}
      </div>
    {/if}
  {/each}
</div>
