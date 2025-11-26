<script>
  import { onMount } from "svelte";
  import { format } from "date-fns";
  import { _ } from "svelte-i18n";
  import { Button, Modal } from "flowbite-svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { DatePicker } from "@svelte-plugins/datepicker";
  import { call, getIcons } from "$lib";
  import { Typography, Divider, DataTable, windowSize } from "$lib";

  const Trash03Icon = getIcons("trash03");
  const BrushCleaningIcon = getIcons("brushCleaning");
  const FiltersIcons = getIcons("filtersIcons");
  const ChevronDownIcon = getIcons("chevronDown");
  const CalendarPlus02 = getIcons("calendarPlus02");

  let isAllLogs = $state(false);
  let deleteModal = $state(false);
  let checkAllState = $state(false);
  let isFilterModal = $state(false);
  let filterAssistantDropdown = $state(false);
  let isOpenDate = $state(false);
  let startDate = $state(null);
  let endDate = $state(null);
  let filterAssistantTerm = $state("all");
  let formattedEndDate = $state("");
  let formattedStartDate = $state("");
  let pagination = $state({
    page: 1,
    pageSize: 10,
    pageCount: 0,
    total: 0,
  });
  let datas = $state([]);

  let width = $derived($windowSize.width);

  let isDatasSelected = $derived.by(() => {
    return datas.some((data) => data.checkbox.checked);
  });

  const checkSingle = (index) => {
    datas[index].checkbox.checked = !datas[index].checkbox.checked;
    const allChecked = datas.every((data) => data.checkbox.checked);
    checkAllState = allChecked;
  };

  const checkAll = () => {
    checkAllState = !checkAllState;
    datas = datas.map((data) => ({
      ...data,
      checkbox: { checked: checkAllState },
    }));
  };

  const abbreviateMessage = (fullMessage, maxLength = 50) => {
    // Check if the message is already short enough
    if (fullMessage.length <= maxLength) {
      return fullMessage;
    }

    // Abbreviate the message and append "..."
    const abbreviatedMessage = fullMessage.substring(0, maxLength);
    return abbreviatedMessage + "...";
  };

  const getChatLogs = async () => {
    let token = localStorage.getItem("token");
    if (!token) return;
    try {
      const response = await $call.call("hexa-ai-rag", "get-conversations-log", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: pagination.page,
          page_size: pagination.pageSize,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      console.log("response getChatLogs", response);

      if (response?.status !== 200) {
        const error = new Error("Failed to fetch chat logs");
        error.status = response?.status || response?.result?.status || 500;
        error.response = response;
        throw error;
      }

      datas = response.result.data.map((item) => ({
        id: item.id,
        // user: item.user || "Assistant",
        botsName: item.attributes.bot_name || "",
        chatName:
          item.attributes.title || `Chat ${format(item.attributes.createdAt, "yyyy-MM-dd")}`,
        // message: abbreviateMessage(item.message) || "",
        datetime: new Date(item.attributes.createdAt).toLocaleString(),
        totalMessages: item.attributes.total_messages || 0,
        checkbox: {
          checked: false,
        },
        actions: {
          objects: [
            {
              icon: $_("eye"),
              onclick: async () => {
                // console.log("item.conversation_id", item.attributes.conversation_id);
                // console.log("item.session_id", item.attributes.session_id);
                const url = await $call.getPage("hexabot", "logs");
                window.open(
                  `${url}?conversation_id=${item.attributes.session_id}&bot_id=${item.attributes.conversation_id}`,
                  "_blank",
                );
              },
            },
          ],
        },
      }));

      pagination.pageCount = response.result?.meta?.pagination?.pageCount || 0;
      pagination.total = response.result?.meta?.pagination?.total || 0;
    } catch (error) {
      console.log("error", error);
    }
  };

  const getBot = async (botId) => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        let response = await $call.call("client-dashboard", "get-bot", {
          URLQuery: {},
          URLParams: {
            id: botId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("response getBot", response);

        if (
          response === null ||
          response === undefined ||
          (response?.status && response.status !== 200)
        ) {
          return;
        }

        const bot = response.result;
        const attributes = bot.attributes;
        const botName = attributes.bot_name;

        console.log("botName", botName);

        // return botName;
      } catch (error) {
        // errorResponseTrigger(error);
      }
    }
  };

  const nextPage = async () => {
    if (pagination.page < pagination.pageCount) {
      pagination.page += 1;
      await getChatLogs();
    }
  };

  const prevPage = async () => {
    if (pagination.page > 1) {
      pagination.page -= 1;
      await getChatLogs();
    }
  };

  const goToPage = async (page) => {
    if (page >= 1 && page <= pagination.pageCount) {
      pagination.page = page;
      await getChatLogs();
    }
  };

  const deleteLog = async (logString) => {
    console.log(logString);
    // const token = localStorage.getItem("token");
    // if (!token) return;
    // try {
    //   const response = await $call.call("hexa-ai-rag", "delete-log-message", {
    //     URLQuery: {},
    //     URLParams: { id: logString },
    //     Params: {},
    //     Headers: {
    //       Authorization: "Bearer " + token,
    //     },
    //   });

    //   if (response?.status !== 200) {
    //     const error = new Error("Failed to delete chat log");
    //     error.status = response?.status || response?.result?.status || 500;
    //     error.response = response;
    //     throw error;
    //   }

    //   await getChatLogs();
    // } catch (error) {
    //   console.log("error", error);
    // }
  };

  const handleDeleteLog = async () => {
    const selectedLogs = datas.filter((data) => data.checkbox.checked);
    const log = selectedLogs.map((log) => log.id).join(",");
    await deleteLog(log);
  };

  const handleDeleteAllLogs = async () => {
    await deleteLog("all");
  };

  const toggleDatePicker = () => {
    isOpenDate = !isOpenDate;
  };

  const formatDate = (dateString) =>
    (dateString && format(new Date(dateString), "d MMM yyyy")) || "";

  onMount(async () => {
    await getChatLogs();
  });

  $effect(() => {
    formattedStartDate = formatDate(startDate);
    formattedEndDate = formatDate(endDate);
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex flex-col">
        <div class="flex items-start justify-between mb-3 md:mb-8">
          <div class="flex flex-col gap-1.5">
            <div class="flex gap-4">
              <Typography preset="Display sm/Semibold" text={$_("chatLogs")} />
            </div>
          </div>

          <div class="flex relative gap-3">
            <Button
              outline
              color="red"
              class="min-w-max gap-1.5"
              disabled={isDatasSelected ? false : true}
              on:click={() => {
                isAllLogs = false;
                deleteModal = true;
              }}
            >
              <Trash03Icon svgClass="w-5 h-5" />
              {#if width > 768}
                <p class="font-semibold">{$_("deleteSelectedLogs")}</p>
              {/if}
            </Button>
            <Button
              outline
              color="red"
              class="min-w-max gap-1.5"
              on:click={() => {
                isAllLogs = true;
                deleteModal = true;
              }}
            >
              <BrushCleaningIcon svgClass="w-5 h-5" />
              {#if width > 768}
                <p class="font-semibold">{$_("cleanAllLogs")}</p>
              {/if}
            </Button>
          </div>
        </div>
        <Divider />
      </div>
      {#if datas.length > 0}
        <div class="flex flex-col">
          <Button
            outline
            color="light"
            class="gap-1.5 w-min mt-4"
            on:click={() => {
              isFilterModal = true;
            }}
          >
            <FiltersIcons svgClass="w-5 h-5" />
            {#if width > 768}
              <p class="font-semibold">{$_("filter")}</p>
            {/if}
          </Button>
          <DataTable
            {datas}
            configs={{ pagination: pagination }}
            {checkAllState}
            next={nextPage}
            previous={prevPage}
            {goToPage}
            {checkSingle}
            {checkAll}
          />
        </div>
      {/if}
    </div>
  </div>
</div>

<Modal bind:open={deleteModal} autoclose>
  <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
  {#if isAllLogs}
    <Typography preset="Text sm/Regular" text={$_("deleteAllLogsFunc")} propClass="text-center" />
  {:else}
    <Typography
      preset="Text sm/Regular"
      text={$_("deleteSelectedLogsFunc")}
      propClass="text-center"
    />
  {/if}
  <div class="w-full flex justify-center mt-6 gap-4">
    <Button onclick={() => (deleteModal = false)} color="alternative" class="w-full">
      {$_("noCancel")}
    </Button>
    <Button
      onclick={async () => {
        if (isAllLogs) {
          await handleDeleteAllLogs();
        } else {
          await handleDeleteLog();
        }
        deleteModal = false;
      }}
      color="red"
      class="w-full"
    >
      {$_("yesImSure")}
    </Button>
  </div>
</Modal>

<Modal bind:open={isFilterModal} outsideclose>
  <svelte:fragment slot="header">
    <div class="flex justify-start items-center gap-4">
      <div
        class="p-2 bg-brand-primary-light dark:bg-brand-primary-dark rounded-full flex justify-center items-center"
      >
        <BrushCleaningIcon
          svgClass="w-5 h-5 text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark"
        />
      </div>
      <div class="flex flex-col">
        <Typography preset="Text lg/Semibold" text={$_("filters")} color="primary" />
      </div>
    </div>
  </svelte:fragment>
  <div class="flex flex-col gap-2 mb-4">
    <Typography preset="Text sm/Medium" color="secondary" text={$_("assistantBots")} />

    <div class="relative w-full">
      <Button
        class="w-full flex justify-between"
        color="light"
        onclick={() => (filterAssistantDropdown = !filterAssistantDropdown)}
      >
        <Typography preset="Text md/Regular" text={filterAssistantTerm} />
        <ChevronDownIcon svgClass="w-5 h-5" />
      </Button>

      {#if filterAssistantDropdown}
        <div
          class="absolute left-0 top-full mt-2 w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-primary-dark rounded-md shadow-lg z-[50]"
          style="min-width:12rem;"
        >
          <ul class="max-h-60 overflow-y-auto">
            <li>
              <button
                class="w-full text-left px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700"
                onclick={() => {
                  filterAssistantTerm = "all";
                  filterAssistantDropdown = false;
                }}
              >
                {$_("all")}
              </button>
            </li>
            {#each Array.from(new Set(datas.map((data) => data.botsName))).filter((name) => name) as botName}
              <li>
                <button
                  class="w-full text-left px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700"
                  onclick={() => {
                    filterAssistantTerm = botName;
                    filterAssistantDropdown = false;
                  }}
                >
                  {botName}
                </button>
              </li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
  </div>
  <div class="flex flex-col gap-2 mb-4">
    <Typography preset="Text sm/Medium" color="secondary" text={$_("timePeriod")} />

    <DatePicker
      bind:isOpen={isOpenDate}
      bind:startDate
      bind:endDate
      isRange
      isMultipane
      enableFutureDates
    >
      <button
        onclick={toggleDatePicker}
        class:open={isOpenDate}
        class="flex w-full p-2.5 border rounded-lg border-gray-300 dark:border-gray-500 gap-2 justify-center items-center"
      >
        <CalendarPlus02 svgClass="w-4 h-4 text-gray-700 dark:text-gray-400" />
        <div class="date">
          {#if startDate != endDate}
            {formattedStartDate} - {formattedEndDate}
          {:else if startDate}
            {formattedStartDate}
          {:else}
            {$_("pickADate")}
          {/if}
        </div>
      </button>
    </DatePicker>
  </div>
  <div class="w-full flex justify-center mt-6 gap-4">
    <Button onclick={() => (isFilterModal = false)} color="alternative" class="w-full">
      {$_("noCancel")}
    </Button>
    <Button
      onclick={async () => {
        console.log("belum ada");
      }}
      color="primary"
      class="w-full"
    >
      {$_("applyFilters")}
    </Button>
  </div>
</Modal>
