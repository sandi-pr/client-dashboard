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
  import { CloseButton, Button, P } from "flowbite-svelte";
  import HexaAi from "$lib/icon/hexa-ai.svelte";
  import HexaAIChat from "$lib/components/HexaAI.svelte";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";

  import { CR3Call } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

  let { innerWidth, changeShowAI } = $props();

  let activities = $state([]);
  let statusRAG = $state("progress");

  let imageFileName = $state("");
  let imageUrl = $state(null);
  let imageFile = $state(null);

  let tabDatas = $state({
    selected: "Hexa AI",
    data: [
      {
        name: "Hexa AI",
      },
      // {
      //   name: "Latest Activity",
      // },
    ],
  });

  const changeTabbar = (tabName) => {
    tabDatas.selected = tabName;
  };

  async function getStatusRag() {
    let token = sessionStorage.getItem("token");
    const callStatusRag = async () => {
      try {
        let rag = await call.call("hexa-ai", "get-rag-status", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        let status = rag.result.rag_status;
        if (statusRAG.includes("!")) {
          let status2 = status.split("!")[0].trim();
          if (status2.includes("(NOVERB)") || status2.includes("(INOVERB)") || status2.includes("(!INOVERB)")) {
            let status3 = status2.replace("(NOVERB)", "").replace("(INOVERB)", "").replace("(!INOVERB)", "").replace("(!NOVERB)", "");
            statusRAG = status3;
          } else {
            statusRAG = status2;
          }
        } else {
          statusRAG = status;
        }
        console.log(rag); // Mengganti statusRAG dengan rag
      } catch (error) {
        console.error("Error:", error);
      }
      setTimeout(callStatusRag, 2000);
    };
    callStatusRag();
  }

  const handleFileUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = ".jpg, .jpeg, .png";

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file) {
        imageUrl = URL.createObjectURL(file);
        imageFile = file;

        imageFileName = file.name;
      } else {
        alert("Invalid file. Please select a valid file.");
      }
    });

    inputElement.click();
  };

  onMount(() => {
    // getStatusRag();
  });
</script>

<Pane
  size={innerWidth > 1024 ? 30 : 100}
  minSize={innerWidth > 1024 ? 30 : 100}
>
  <div
    class="h-screen transition-all duration-300 border-l-2 border-gray-200 dark:border-gray-800"
  >
    <div class="sticky top-0 z-10">
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
            <div class="flex grow justify-between">
              <div class="flex items-center">
                <div
                  class="w-10 h-10 p-2.5 mr-3 bg-brand-secondary-light dark:bg-brand-secondary-dark rounded-full justify-center items-center"
                >
                  <svelte:component
                    this={HexaAi}
                    svgClass="text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark"
                  />
                </div>
                <div>
                  <Typography
                    preset="Text xl/Semibold"
                    text="Hexa AI"
                    color="primary"
                  />
                  <!-- <Typography
                    preset="Text sm/Regular"
                    text={`Status RAG: ${statusRAG}`}
                    color="tertiary"
                  /> -->
                  <!-- <Typography
                    preset="Text sm/Regular"
                    text="Assistant"
                    color="tertiary"
                  /> -->
                </div>
              </div>
              <Button
                on:click={() => {
                  handleFileUploadWrapper();
                }}
                color="light"
                class={`w-max mt-2 ms-4 rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1 flex flex-col justify-start`}
              >
                <P>{$_("newOcr")}</P>
              </Button>
            </div>
          {/if}
          {#if innerWidth < 1024}
            <CloseButton on:click={changeShowAI} class="mb-4 dark:text-white" />
          {/if}
        </div>
      </div>
      <!-- <TabBar {tabDatas} {changeTabbar} /> -->
    </div>
    <div class="flex w-full flex-col h-screen">
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
        <HexaAIChat {imageUrl} {imageFileName} {imageFile} />
      {/if}
    </div>
  </div>
</Pane>
