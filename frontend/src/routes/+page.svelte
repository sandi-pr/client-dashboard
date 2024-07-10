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
  import { Card, Typography, SectionHeader } from "$lib";
  import { StateCard } from "$lib";
  import IllustrationFirst from "$lib/icon/illustration_first_time.svelte";
  import { Alert, CloseButton, P } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  const call = new CR3Call(globalThis.location.origin);
  let userData = $state({ username: "", email: "", brand: "" });

  let isFirstTime = $state(true);

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
      var resp_setup = await call.call("client-login", "is-need-setup", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {},
      });
      console.log(resp_setup);
      if (resp_setup.status === 200) {
        if (resp_setup.result.is_need_setup === true) {
          if (sessionStorage.getItem("token") !== null) sessionStorage.removeItem("token");
          if (sessionStorage.getItem("user") !== null) sessionStorage.removeItem("user");
          await call.goToPage("hexapod-setup", "home-page");
          return;
        }
      }
      var token = sessionStorage.getItem("token");
      if (token !== null) {
        var resp = await call.call("client-login", "is-token-valid", {
          URLQuery: {},
          URLParams: {},
          Params: {
            token: token,
          },
          Headers: {},
        });

        if (resp.status === 200) {
          if (resp.result.confirmed === true && resp.result.blocked === false) {
            sessionStorage.setItem("user", JSON.stringify(resp.result));
            userData = resp.result;
          }
        }
      } else {
        await call.goToPage("client-login", "home-page");
        return;
      }
    } catch (error) {
      // console.log(error);
      userData = { username: "Hexa Staff", email: "hexacode@hexacode.org", brand: "hexacode" };
    }
  }

  onMount(() => {
    getFristTime();
    CheckSetup();
  })

</script>

{#if !isFirstTime}
  <div class="flex flex-col gap-1 px-4 pt-6">
    <Typography preset="Display sm/Semibold" text={$_("welcometoHexapod")} color="primary" />
    <Typography
      text="{userData.username} - {userData.email}"
      color="tertiary"
    />
  </div>
{/if}
<div class="flex flex-col gap-4 px-4 pt-8">
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
      <Typography preset="Display sm/Semibold" text={$_("welcometoHexapod")} color="primary" />
      <!-- <Typography
        text="{userData?.data?.username || 'Hexastaff'} - {userData?.data?.email ||
          'staff@hexacode.org'}"
        color="tertiary"
      /> -->
    </div>
  </Alert>
  {/if}
  <SectionHeader title={$_("availableModules")} />
  <Card
    icon="users"
    title="{$_('crmOverview')}"
    description={$_("descCrm")}
    onclick={async () => {
      goto(await call.getPage("client-dashboard", "crm-overview"));
      return;
    }}
  />
  <Card
    icon="fileShield"
    title="{$_('documentCollection')}"
    description={$_("documentCollection")}
    onclick={async () => {
      goto(await call.getPage("client-dashboard", "document-collection"));
      return;
    }}
  />
  <Card
    icon="userSquare"
    title="{$_('HR')}"
    description={$_("descHr")}
    onclick={async () => {
      goto(await call.getPage("client-dashboard", "hr-overview"));
      return;
    }}
  />
</div>
