<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: Azzis Arswendo <azzis@hexacode.org>

Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->

<script>
  import { Typography, TabBar, Icon } from "$lib";
  import { Pane } from "svelte-splitpanes";
  import DataflowIcon from "$lib/icon/dataflow-solid.svelte";
  import {
    CloseButton,
    Button,
    P,
    Dropdown,
    DropdownItem,
    Modal,
    Input,
    Card,
    Badge,
  } from "flowbite-svelte";
  import { ChevronRightOutline, SearchOutline } from "flowbite-svelte-icons";
  import SettingsIcon from "$lib/icon/settings-02.svelte";
  import RefreshCwIcon from "$lib/icon/refresh-cw-01.svelte";
  import HexaAi from "$lib/icon/hexa-ai.svelte";
  import HexaAIChat from "$lib/components/HexaAI.svelte";
  import ArrowUpRightIcon from "$lib/icon/arrow-up-right.svelte";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";

  import { CR3Call } from "hexacomponents-svelte";
  import CollectionsAi from "./CollectionsAI.svelte";
  import Ocr from "./OCR.svelte";
  import {
    triggerSelectCollection,
    selectedCollection,
    selectedDropdown,
    triggerReload,
    triggerSelectDropdown,
  } from "$lib";

  import { call } from "$lib";

  let { innerWidth, changeShowAI } = $props();

  let activities = $state([]);
  let statusRAG = $state("progress");

  let imageFileName = $state("");
  let imageUrl = $state(null);
  let imageFile = $state(null);

  let dropdownOpen = $state(false);
  let modalOpen = $state(false);
  // let selectedCollection = $state(null);
  let tmpSelectedCollection = $state(null);
  let disableButton = $state(false);
  let disableConfigButton = $state(false);

  let initCollections = $state([]);
  let documents = $state([]);
  let collections = $state([]);
  let nameCollection = $state(null);
  let reloadCollection = $state(true);

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

  let dropdownOpenSource = $state([{ id: 1, bool: false }]);

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    // make selection of buttons doesn't work
    // dropdownOpenSource = dropdownOpenSource.map((item) => {
    //   return { id: item.id, bool: false };
    // });
  };

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let page = 1;
      initCollections = [];
      while (true) {
        let response = await $call.call("client-dashboard", "get-document-collections", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 25,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        if (response.result.data.length > 0) {
          initCollections = initCollections.concat(response.result.data);
          page++;
        } else {
          break;
        }

        let page2 = 1;
        documents = [];
        while (true) {
          let response2 = await $call.call("client-dashboard", "get-documents", {
            URLQuery: {},
            URLParams: {},
            Params: {
              page: page2,
              page_size: 25,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response2.status !== 200) {
            console.log(response2);
            return;
          }

          if (response2.result.data.length > 0) {
            documents = documents.concat(response2.result.data);
            page2++;
          } else {
            break;
          }
        }

        collections = initCollections.map((collection) => {
          let totalDocuments = documents.filter(
            (document) => document.attributes.collection_id === collection.id,
          ).length;

          console.log("=====================", collection.attributes);

          // let documentCategory = [
          //   {
          //     id: 1,
          //     name: "Category 1",
          //   },
          //   {
          //     id: 2,
          //     name: "Category 2",
          //   },
          // ];

          if (
            collection.attributes.categories !== null &&
            collection.attributes.categories.length > 0
          ) {
            dropdownOpenSource ==
              collection.attributes.categories.map((item) => {
                return { id: item.id, bool: false };
              });
          }

          return {
            id: collection.id,
            name: collection.attributes.collection_name,
            description: collection.attributes.description,
            totalDocuments: totalDocuments,
            category: collection.attributes.categories || [],
            selectedCategory: { id: 1, name: "Category 1" } || null,
          };
        });
      }
    }
  }

  async function getStatusRag() {
    let token = localStorage.getItem("token");
    const callStatusRag = async () => {
      try {
        let rag = await $call.call("hexa-ai", "get-rag-status", {
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
          if (
            status2.includes("(NOVERB)") ||
            status2.includes("(INOVERB)") ||
            status2.includes("(!INOVERB)")
          ) {
            let status3 = status2
              .replace("(NOVERB)", "")
              .replace("(INOVERB)", "")
              .replace("(!INOVERB)", "")
              .replace("(!NOVERB)", "");
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

  const selectDropdown = (item) => {
    dropdownOpen = false;
    if (item === "general") {
      nameCollection = null;
    }

    if (item === "collections") {
      if ($selectedCollection !== null) {
        tmpSelectedCollection = $selectedCollection;
        callColection();
      }
    }

    triggerSelectDropdown(item);
    localStorage.setItem("selectedDropdown", item);
  };

  const selectCollection = (item) => {
    tmpSelectedCollection = item;
  };

  async function callColection() {
    let token = localStorage.getItem("token");

    let response2 = await $call.call("hexa-ai-rag", "add-document-by-collection-id", {
      URLQuery: {},
      URLParams: {},
      Params: {
        collection_id: tmpSelectedCollection.id.toString(),
      },
      Headers: {
        Authorization: "Bearer " + token,
      },
    });

    if (response2.status !== 200) {
      console.log(response2);
      return;
    }

    triggerReload(true);
    triggerSelectCollection(tmpSelectedCollection);
    localStorage.setItem("selectedCollection", JSON.stringify(tmpSelectedCollection));
    console.log(response2);
  }

  function clearChat() {
    if ($selectedDropdown === "collections") {
      sessionStorage.removeItem("conversationIdCollection-" + $selectedCollection.id);
    } else {
      sessionStorage.removeItem("conversationId");
    }
    window.location.reload();
  }

  $effect(() => {
    disableButton = tmpSelectedCollection === null ? true : false;
    disableConfigButton =
      $selectedDropdown === "general" ? true : $selectedDropdown === "ocr" ? true : false;
  });

  $effect(() => {
    if ($selectedDropdown === "collections") {
      if ($selectedCollection !== null) {
        nameCollection = $selectedCollection.name;
        if (reloadCollection) {
          reloadCollection = false;
          tmpSelectedCollection = $selectedCollection;
          callColection();
        }
      } else {
        nameCollection = null;
      }
    }
  });

  onMount(() => {
    // getStatusRag();
    getCollections();
  });
</script>

<Pane size={innerWidth > 1024 ? 40 : 100} minSize={innerWidth > 1024 ? 40 : 100}>
  <div class="h-screen transition-all duration-300 border-l-2 border-gray-200 dark:border-gray-800">
    <div class="sticky top-0 z-10">
      <div class="flex items-center">
        <div
          id="drawer-label"
          class="w-full relative inline-flex items-center mx-4 mt-4 text-base font-semibold text-gray-500 dark:text-gray-400"
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
              <Typography preset="Text xl/Semibold" text="Latest Activity" color="primary" />
              <Typography preset="Text sm/Regular" text="Your latest activity" color="tertiary" />
            </div>
          {:else if tabDatas.selected === "Hexa AI"}
            <div class="flex flex-col grow">
              <div
                class="flex grow justify-between pb-4 border-b border-primary-light dark:border-primary-dark"
              >
                <div class="flex w-full items-center justify-between">
                  <div
                    class="w-10 h-10 p-2.5 mr-3 bg-brand-secondary-light dark:bg-brand-secondary-dark rounded-full justify-center items-center"
                  >
                    <svelte:component
                      this={HexaAi}
                      svgClass="text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark"
                    />
                  </div>
                  <div>
                    <Typography preset="Text xl/Semibold" text="Hexa AI" color="primary" />
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
                  <div class="flex gap-2">
                    {#if !disableConfigButton}
                      <Button
                        disabled={disableConfigButton}
                        on:click={() => (modalOpen = true)}
                        color="light"
                        class={`w-max mt-4 rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-2 flex justify-start`}
                      >
                        <!-- <P>{$_("configure")}</P> -->
                        <SettingsIcon
                          svgClass="w-4 h-4 text-fg-secondary-light dark:text-fg-secondary-dark"
                        />
                      </Button>
                    {/if}
                    <Button
                      on:click={clearChat}
                      color="light"
                      class={`w-max mt-4 !px-3 rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-2 flex justify-start`}
                    >
                      <RefreshCwIcon
                        svgClass="w-4 h-4 text-fg-secondary-light dark:text-fg-secondary-dark"
                      />
                    </Button>
                  </div>
                </div>
              </div>
              <div class="flex justify-end items-center">
                <!-- <Button
                  on:click={() => (dropdownOpen = !dropdownOpen)}
                  color="light"
                  class={`w-max-fit mr-3 mt-4 !pr-3 rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1 flex justify-start`}
                >
                  <P class="line-clamp-1"
                    >{selectedDropdown === "general"
                      ? $_($selectedDropdown)
                      : nameCollection !== null
                        ? nameCollection
                        : $_($selectedDropdown)}</P
                  >
                  <ChevronRightOutline
                    class="text-fg-secondary-light dark:text-fg-secondary-dark"
                  />
                </Button> -->
                <!-- <Typography preset="Text md/Medium" text={$_("general")} color="tertiary" /> -->
                <!-- Dropdown menu -->
                <!-- <Dropdown placement="right" containerClass="border rounded-xl">
                  <DropdownItem>
                    <button onclick={() => selectDropdown("general")}>
                      {$_("general")}
                    </button>
                  </DropdownItem>
                  <DropdownItem>
                    <button onclick={() => selectDropdown("collections")}>
                      {$_("collections")}
                    </button>
                  </DropdownItem>
                </Dropdown> -->
                <!-- <div class="flex gap-2">
                  {#if !disableConfigButton}
                    <Button
                      disabled={disableConfigButton}
                      on:click={() => (modalOpen = true)}
                      color="light"
                      class={`w-max mt-4 rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-2 flex justify-start`}
                    >
                      <SettingsIcon
                        svgClass="w-4 h-4 text-fg-secondary-light dark:text-fg-secondary-dark"
                      />
                    </Button>
                  {/if}
                  <Button
                    on:click={clearChat}
                    color="light"
                    class={`w-max mt-4 !px-3 rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-2 flex justify-start`}
                  >
                    <RefreshCwIcon
                      svgClass="w-4 h-4 text-fg-secondary-light dark:text-fg-secondary-dark"
                    />
                  </Button>
                </div> -->
              </div>
            </div>
          {/if}
          {#if innerWidth < 1024}
            <CloseButton
              on:click={changeShowAI}
              class="mb-4 absolute right-0 top-1 dark:text-white"
            />
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
            <Typography preset="Text sm/Medium" text="No Activity" color="primary" />
          {/if}
        </div>
      {:else if tabDatas.selected === "Hexa AI"}
        {#if $selectedDropdown === "general"}
          <HexaAIChat />
        {:else if $selectedDropdown === "collections"}
          <CollectionsAi />
          <!-- {:else if selectedDropdown === "ocr"}
        <Ocr /> -->
        {/if}
      {/if}
    </div>
  </div>
  <Modal title={$_("chooseSource")} bind:open={modalOpen} autoclose outsideclose>
    <div class="flex flex-col justify-center gap-4">
      <Input
        id="collection"
        size="lg"
        color="base"
        class="focus:border-hc-brand-600 focus:ring-hc-brand-600"
        type="text"
        placeholder={$_("search")}
      >
        <SearchOutline slot="left" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
      </Input>
      <div class="flex flex-col justify-center mt-3 gap-3">
        <div class="flex text-start">
          <P class="text-tertiary-light dark:text-tertiary-dark text-base font-normal"
            >{$_("recentCollections")}</P
          >
        </div>

        <div class="flex flex-col gap-6">
          {#each collections as item, index (item)}
            <Card
              class="{tmpSelectedCollection === item
                ? 'border-brand-solid-light dark:border-brand-solid-dark ring-brand-solid-light dark:ring-brand-solid-dark'
                : ''}bg-hc-graylight-50 {item.totalDocuments > 0
                ? 'hover:bg-hc-brand-50 cursor-pointer'
                : 'opacity-50'} gap-3"
              size="xl"
              on:click={item.totalDocuments > 0 ? () => selectCollection(item) : null}
            >
              <div class="flex flex-col gap-0.5">
                <P
                  weight="semibold"
                  class={`${tmpSelectedCollection === item ? "text-brand-primary-light dark:text-brand-primary-dark" : "text-primary-light dark:text-primary-dark"}`}
                  >{item.name}</P
                >
                <P
                  size="xs"
                  class={`${tmpSelectedCollection === item ? "text-brand-secondary-light dark:text-brand-secondary-dark" : "text-secondary-light dark:text-secondary-dark"}`}
                  >{item.description}</P
                >
              </div>
              <div class="flex justify-between">
                <Badge
                  rounded
                  class="bg-utility-brand-50-light dark:bg-utility-brand-50-dark text-utility-brand-700-light dark:text-utility-brand-700-dark border border-utility-brand-200-light dark:border-utility-brand-200-dark"
                  >{item.totalDocuments} {$_("documents")}</Badge
                >
              </div>
              <!-- {#if item.category.length > 0}
                <div class="flex gap-4 items-center">
                  <P size="sm" class="text-secondary-light dark:text-secondary-dark"
                    >{$_("category")}</P
                  >
                  <div onfocusout={handleDropdownFocusLoss} class="relative">
                    <Button
                      outline
                      on:click={(event) => {
                        event.stopPropagation();
                        dropdownOpenSource = dropdownOpenSource.map((item2) => {
                          if (item2.id === item.category[index].id) {
                            return { id: item2.id, bool: !item2.bool };
                          } else {
                            return { id: item2.id, bool: false };
                          }
                        });
                      }}
                      class={`py-2 w-full justify-between gap-1 px-2.5 md:px-3 hover:bg-primary-light dark:hover:bg-primary-dark rounded-lg border-secondary-light dark:border-secondary-dark`}
                    >
                      <Typography
                        preset="Text sm/Semibold"
                        text={item.selectedCategory.name}
                        color="secondary"
                      />
                      <Icon name="chevronDown" />
                    </Button>

                    {#if item.category[index]?.id === undefined}
                      {console.log(item.category)}
                    {/if}
                    <div
                      id={"dropdown-" + item.category[index].id}
                      class={`z-10 bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700 ${
                        dropdownOpenSource.find((item2) => item2.id === item.category[index].id)
                          .bool
                          ? "fixed"
                          : "hidden"
                      }`}
                    >
                      <ul
                        class="py-2 text-sm text-gray-700 dark:text-gray-200"
                        aria-labelledby="dropdownDefaultButton"
                      >
                        {#each item.category as category (category)}
                          <button
                            class="flex p-2"
                            onclick={(event) => {
                              console.log("test before");
                              event.stopPropagation();
                              console.log("test");
                              item.selectedCategory = category;
                              dropdownOpenSource = dropdownOpenSource.map((item2) => {
                                return { id: item2.id, bool: false };
                              });
                            }}
                          >
                            <P>{category.name}</P>
                          </button>
                        {/each}
                      </ul>
                    </div>
                  </div>
                  <Button
                    on:click={(event) => {
                      event.stopPropagation();
                      window.open("", "_blank");
                    }}
                    outline
                    class="!flex hover:bg-primary-light dark:hover:bg-primary-dark rounded-lg border-secondary-light dark:border-secondary-dark"
                  >
                    <Typography preset="Text sm/Semibold" text={$_("detail")} color="secondary" />
                    <ArrowUpRightIcon
                      svgClass="w-5 h-5 text-secondary-light dark:text-secondary-dark"
                    />
                  </Button>
                </div>
              {/if} -->
            </Card>
          {/each}
        </div>
      </div>
    </div>
    <svelte:fragment slot="footer">
      <Button on:click={callColection} class="w-full" disabled={disableButton}
        >{$_("select")}</Button
      >
    </svelte:fragment>
  </Modal>
</Pane>
