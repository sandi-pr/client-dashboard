<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: Azzis Arswendo <azzis@hexacode.org>

Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->

<script>
  import { Typography, TabBar, Card } from "$lib";
  import { Pane } from "svelte-splitpanes";
  import DataflowIcon from "$lib/icon/dataflow-solid.svelte";
  import { CloseButton } from "flowbite-svelte";
  import HexaAi from "$lib/icon/hexa-ai.svelte";
  import HexaAIChat from "$lib/components/HexaAI.svelte";
  import { onMount } from "svelte";

  let { innerWidth, changeShowAI } = $props();

  let activities = $state([]);
  let statusRAG = $state("progress");

  let tabDatas = $state({
    selected: "Hexa AI",
    data: [
      {
        name: "Hexa AI",
      },
      {
        name: "Latest Activity",
      },
    ],
  });

  const changeTabbar = (tabName) => {
    tabDatas.selected = tabName;
  };

  async function getStatusRag() {
    let token = sessionStorage.getItem("token");
    const callStatusRag = async () => {
      try {
        let rag = await call.call("hexa-ai", "get-status-rag", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        statusRAG = rag.result;
        console.log(rag); // Mengganti statusRAG dengan rag
      } catch (error) {
        console.error("Error:", error);
      }
      setTimeout(callStatusRag, 2000);
    };
    callStatusRag();
  }

  onMount(() => {
    getStatusRag();
  });
</script>

<Pane
  size={innerWidth > 1024 ? "30" : "100"}
  minSize={innerWidth > 1024 ? "30" : "100"}
  maxSize={innerWidth > 1024 ? "50" : "100"}
>
  <div
    class="h-screen fixed transition-all duration-300 border-l-2 border-gray-200 dark:border-gray-800"
  >
    <div class="absolute w-screen top-0 z-10">
      <div class="flex items-center">
        <div
          id="drawer-label"
          class="w-full inline-flex items-center mx-4 mt-4 text-base font-semibold text-gray-500 dark:text-gray-400"
        >
          {#if tabDatas.selected === "Latest Activity"}
            <div
              class="w-10 h-10 p-2.5 mr-3 bg-brand-secondary-light dark:bg-brand-secondary-dark rounded-full justify-center items-center"
            >
              <svelte:component
                this={DataflowIcon}
                svgClass="text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark"
              />
            </div>
            &nbsp;
            <div>
              <Typography
                preset="Text xl/Semibold"
                text="Latest Activity"
                color="primary"
              />
              <Typography
                preset="Text sm/Regular"
                text="Your latest activity"
                color="tertiary"
              />
            </div>
          {:else if tabDatas.selected === "Hexa AI"}
            <div class="flex w-100  justify-between">
              <div class="flex">
                <div
                  class="w-10 h-10 p-2.5 mr-3 bg-brand-secondary-light dark:bg-brand-secondary-dark rounded-full justify-center items-center"
                >
                  <svelte:component
                    this={HexaAi}
                    svgClass="text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark"
                  />
                </div>
                &nbsp;
                <div>
                  <Typography
                    preset="Text xl/Semibold"
                    text="Hexa AI"
                    color="primary"
                  />
                  <Typography
                    preset="Text sm/Regular"
                    text="Assistant"
                    color="tertiary"
                  />
                </div>
              </div>
              <Typography
                preset="Text sm/Regular"
                text={`Status RAG: ${statusRAG}`}
                color="tertiary"
              />
            </div>
          {/if}
          {#if innerWidth < 1024}
            <CloseButton on:click={changeShowAI} class="mb-4 dark:text-white" />
          {/if}
        </div>
      </div>
      <TabBar {tabDatas} {changeTabbar} />
    </div>
    <div class="flex flex-col h-screen pt-[118px]">
      {#if tabDatas.selected === "Latest Activity"}
        <div class="p-4">
          {#if activities.length > 0}
            {#each activities as activityYear (activityYear)}
              {#each Object.keys(activityYear) as year (year)}
                <div class="flex flex-col gap-6">
                  {#each Object.keys(activityYear[year]) as month (month)}
                    <div class="flex flex-col gap-3">
                      <Typography
                        preset="Text sm/Medium"
                        text="{numberMonthToTextMonth(month)} {year}"
                        color="tertiary"
                      />
                      {#each activityYear[year][month] as activity (activity)}
                        <Card
                          type="ledger"
                          title={activity.activity}
                          description={activity.address}
                          subtext={activity.date}
                        />
                      {/each}
                    </div>
                  {/each}
                </div>
              {/each}
            {/each}
          {:else}
            <Typography
              preset="Text sm/Medium"
              text="No Activity"
              color="primary"
            />
          {/if}
        </div>
      {:else if tabDatas.selected === "Hexa AI"}
        <HexaAIChat />
      {/if}
    </div>
  </div>
</Pane>
