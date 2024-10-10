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
  import { CR3Call } from "hexacomponents-svelte";
  import { Card, Typography, SectionHeader, triggerLoading } from "$lib";
  import { StateCard } from "$lib";
  import IllustrationFirst from "$lib/icon/illustration_first_time.svelte";
  import { Alert, CloseButton, P } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { goto } from "$app/navigation";
  import { onMount, tick } from "svelte";
  import { writable } from "svelte/store";
  import LoadingContainer from "../lib/components/LoadingContainer.svelte";
  import { call, showAI } from "$lib";

  let userData = $state({ username: "", email: "", brand: "" });

  let isFirstTime = $state(true);

  const cardsBusiness = [
    // {
    //   icon: "user",
    //   title: "person",
    //   description: "descPerson",
    //   page: "crm-person",
    // },
    // {
    //   icon: "faceHappy",
    //   title: "HR",
    //   description: "descHr",
    //   page: "hr-overview",
    // },
    // {
    //   icon: "package",
    //   title: "product&Services",
    //   description: "descProduct&Services",
    //   page: "crm-product-service",
    // },
    // {
    //   icon: "receipt",
    //   title: "order&Invoice",
    //   description: "descOrder&Invoice",
    //   page: "crm-order-invoice",
    // },
    // {
    //   icon: "calendar",
    //   title: "calendar",
    //   description: "descCalendar",
    //   page: "crm-calendar",
    // },
  ];

  const cardsAi = [
    {
      icon: "messageChat",
      title: "hexabot",
      description: "descHexabot",
      page: "hexabot",
    },
    {
      icon: "archive",
      title: "documentCollection",
      description: "descDocumentCollection",
      page: "document-collection-overview",
    },
    // {
    //   icon: "fileCheck",
    //   title: "contractReviewer",
    //   description: "descContractReviewer",
    //   page: "contract-reviewer-overview",
    // },
    {
      icon: "data",
      title: "datasets",
      description: "descDatasets",
      page: "datasets",
    },
  ];

  function setFristTime() {
    isFirstTime = false;
    localStorage.setItem("isFirstTime", isFirstTime);
  }

  function getFristTime() {
    if (typeof localStorage !== "undefined") {
      isFirstTime = localStorage.getItem("isFirstTime") === "true" ? true : false;
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

  const navigateToPage = async (page) => {
    if (page === "bakrie-webchat") {
      goto(await $call.getPage("bakrie", "home-page"));
    } else {
      goto(await $call.getPage("client-dashboard", page));
    }
  };

  onMount(() => {
    // triggerLoading(true);
    // setTimeout(async () => {
    //   triggerLoading(false);
    // }, 3000);
    getFristTime();
    CheckSetup();
  });
</script>

<div class="flex flex-col gap-8 px-4 md:py-6">
  {#if isFirstTime}
    <Alert dismissable on:close={setFristTime} class="relative flex justify-center">
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
          text={$_("welcometoHexapod") + sessionStorage.getItem("sidebar_title")}
          color="primary"
        />
        <!-- <Typography
        text="{userData?.data?.username || 'Hexastaff'} - {userData?.data?.email ||
          'staff@hexacode.org'}"
        color="tertiary"
      /> -->
      </div>
    </Alert>
  {:else}
    <div class="flex flex-col gap-1">
      <Typography
        preset="Display sm/Semibold"
        text={`${sessionStorage.getItem("sidebar_title")} ${$_("library")}`}
        color="primary"
      />
      <Typography text="{userData.username} - {userData.email}" color="tertiary" />
    </div>
  {/if}
  {#if cardsBusiness.length > 0}
    <div class="flex flex-col gap-4">
      <SectionHeader title={$_("businessManagements")} />
      <div
        class={`grid ${$showAI ? "grid-cols-1 xl:grid-cols-2" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-4"} gap-4`}
      >
        {#each cardsBusiness as card}
          <Card
            icon={card.icon}
            title={$_(card.title)}
            description={$_(card.description)}
            onclick={() => navigateToPage(card.page)}
          />
        {/each}
      </div>
    </div>
  {/if}
  {#if cardsAi.length > 0}
    <div class="flex flex-col gap-4">
      <SectionHeader title={$_("aiManagements")} />
      <div
        class={`grid ${$showAI ? "grid-cols-1 xl:grid-cols-2" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-4"} gap-4`}
      >
        {#each cardsAi as card}
          <Card
            icon={card.icon}
            title={$_(card.title)}
            description={$_(card.description)}
            onclick={() => navigateToPage(card.page)}
          />
        {/each}
      </div>
    </div>
  {/if}
</div>

<!-- <LoadingContainer /> -->
