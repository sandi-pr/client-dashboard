<!--
  CR3 Client Dashboard
  Package: gitlab.hexacode.org/public-repo/js/hexacomponents-svelte
  Maintainer:
  * Sandi Permana <sandy@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>

  Copyright (C) 2024 Hexacode Teknologi Indonesia
  All Rights Reserved.
-->

<script lang="ts">
  // import { page } from "$app/stores";
  import { sineIn } from "svelte/easing";
  import { Breadcrumb, Button, LanguageSelector, Typography, TabBar, Card } from "$lib";
  // import TabDetail from "$lib/components/TabDetail.svelte";
  // import AIBots from "$lib/components/AIBots.svelte";
  import { DarkMode, Drawer, CloseButton, P } from "flowbite-svelte";
  import BellIcon from "$lib/icon/bell-01.svg";
  import StarsIcon from "$lib/icon/stars-01.svg";
  import DataflowIcon from "$lib/icon/dataflow-solid.svelte";
  import HexaAi from "$lib/icon/hexa-ai.svelte";

  let drawerIsOpen = $state(true);

  let transitionParamsRight = {
    x: 320,
    duration: 200,
    easing: sineIn
  };

  let tabDatas = $state({
    selected: "Latest Activity",
    data: [
      {
        name: "Latest Activity"
      },
      {
        name: "Hexa AI"
      }
    ]
  });

  let { activities = [], activeUrl = "/", changeLanguage = () => {} } = $props();

  const numberMonthToTextMonth = (month: any) => {
    const monthNames = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
    return monthNames[month - 1];
  };
</script>

<div>
  <div class="flex justify-between px-4 lg:pt-8">
    <Breadcrumb {activeUrl} />
    <div class="flex space-x-4">
      <LanguageSelector {changeLanguage} />
      <DarkMode btnClass="text-gray-500 dark:text-gray-400 bg-white dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none rounded-lg p-2.5 border border-gray-300 dark:border-gray-600 dark:hover:border-gray-600 focus-within:ring-4 focus-within:outline-none focus-within:ring-gray-200 dark:focus-within:ring-gray-700" />
      <Button color="light" icon={BellIcon} classesSvg="bell" spaceY="mt-0" />
      <Button
        color="primary"
        icon={StarsIcon}
        spaceY="mt-0"
        funcButton={() => {
          drawerIsOpen = !drawerIsOpen;
          console.log(drawerIsOpen);
        }} />
    </div>
  </div>
</div>

<Drawer placement="right" transitionType="fly" transitionParams={transitionParamsRight} bind:hidden={drawerIsOpen} id="sidebar-dashboard" width="w-full lg:w-96" class="h-screen bg-primary-light dark:bg-secondary-dark">
  <div class="flex items-center">
    <div id="drawer-label" class="inline-flex items-center mb-4 text-base font-semibold text-gray-500 dark:text-gray-400">
      {#if tabDatas.selected === "Latest Activity"}
        <div class="w-10 h-10 p-2.5 mr-3 bg-brand-secondary-light dark:bg-brand-secondary-dark rounded-full justify-center items-center">
          <svelte:component this={DataflowIcon} svgClass="text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark" />
        </div>
        &nbsp;
        <div>
          <Typography preset="Text xl/Semibold" text={"latestActivity"} color="primary" />
          <Typography preset="Text sm/Regular" text={"yourLatestActivity"} color="tertiary" />
        </div>
      {:else if tabDatas.selected === "Hexa AI"}
        <!-- <img src={HexaAi} alt="Icon Hexa AI" class="bg-primary-100 p-2 rounded-full" /> -->
        <div class="w-10 h-10 p-2.5 mr-3 bg-brand-secondary-light dark:bg-brand-secondary-dark rounded-full justify-center items-center">
          <svelte:component this={HexaAi} svgClass="text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark" />
        </div>
        &nbsp;
        <div>
          <Typography preset="Text xl/Semibold" text="Hexa AI" color="primary" />
          <Typography preset="Text sm/Regular" text={"yourCrossroadsPersonalAssistant"} color="tertiary" />
        </div>
      {/if}
    </div>
    <CloseButton
      on:click={() => {
        drawerIsOpen = true;
        console.log(drawerIsOpen);
      }}
      class="mb-4 dark:text-white" />
  </div>

  <!-- <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} /> -->
  <TabBar {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

  {#if tabDatas.selected === "Latest Activity"}
    {#if activities.length > 0}
      {#each activities as activityYear (activityYear)}
        {#each Object.keys(activityYear) as year (year)}
          <div class="flex flex-col gap-6 mt-4">
            <!-- <Typography preset="Text md/Medium" text={year} color="primary" /> -->
            {#each Object.keys(activityYear[year]) as month (month)}
              <div class="flex flex-col gap-3">
                <Typography preset="Text sm/Medium" text="{numberMonthToTextMonth(month)} {year}" color="tertiary" />
                {#each activityYear[year][month] as activity (activity)}
                  <Card type="ledger" title={activity.activity} description={activity.address} subtext={activity.date} />
                {/each}
              </div>
            {/each}
          </div>
        {/each}
      {/each}
    {:else}
      <Typography preset="Text sm/Medium" text="No Activity" color="primary" />
    {/if}
  {:else if tabDatas.selected === "Hexa AI"}
    <div class="w-full h-full mt-4">
      <!-- <AIBots /> -->
    </div>
  {:else}
    ""
  {/if}
</Drawer>
