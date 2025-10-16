<script>
  import { _ } from "svelte-i18n";
  import { Button, Modal } from "flowbite-svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { call, getIcons } from "$lib";
  import { Typography, Divider, DataTable, InputRw, TextareaRw } from "$lib";
  import { onMount } from "svelte";

  const Trash03Icon = getIcons("trash03");
  const BrushCleaningIcon = getIcons("brushCleaning");

  let checkAllState = $state(false);
  let deleteModal = $state(false);
  let isAllLogs = $state(false);
  let isViewAll = $state(true);
  let pagination = $state({
    page: 1,
    pageSize: 10,
    pageCount: 0,
    total: 0,
  });
  let datas = $state([]);
  let data = $state({
    id: 1,
    user: "",
    botsName: "",
    chatName: "",
    message: "",
    datetime: "",
    checkbox: {
      checked: false,
    },
  });

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
      const response = await $call.call("hexa-ai-rag", "get-log-messages", {
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

      // console.log("response", response);
      if (response?.status !== 200) {
        const error = new Error("Failed to fetch chat logs");
        error.status = response?.status || response?.result?.status || 500;
        error.response = response;
        throw error;
      }

      datas = response.result.data.map((item) => ({
        id: item.id,
        user: item.user || "Assistant",
        botsName: item.bot_name || "Support Bot",
        chatName: item.chat_name || "",
        message: abbreviateMessage(item.message) || "",
        datetime: new Date(item.datetime).toLocaleString(),
        checkbox: {
          checked: false,
        },
        actions: {
          objects: [
            {
              icon: $_("eye"),
              onclick: async () => {
                await getChatLog(item.id);
                isViewAll = false;
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

  const getChatLog = async (id) => {
    let token = localStorage.getItem("token");
    if (!token) return;
    try {
      const response = await $call.call("hexa-ai-rag", "get-log-message", {
        URLQuery: {},
        URLParams: { id: id.toString() },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });
      // console.log("response", response);

      if (response?.status !== 200) {
        const error = new Error("Failed to fetch chat log");
        error.status = response?.status || response?.result?.status || 500;
        error.response = response;
        throw error;
      }

      data = {
        id: response.result.id,
        user: response.result.user || "Assistant",
        botsName: response.result.bot_name || "Support Bot",
        chatName: response.result.chat_name || "",
        message: response.result.message || "",
        datetime: new Date(response.result.datetime).toLocaleString(),
        checkbox: {
          checked: false,
        },
      };
    } catch (error) {
      console.log("error", error);
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

  onMount(async () => {
    await getChatLogs();
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
            {#if isViewAll}
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
                <p class="font-semibold">{$_("deleteSelectedLogs")}</p>
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
                <p class="font-semibold">{$_("cleanAllLogs")}</p>
              </Button>
            {:else}
              <Button
                outline
                color="alternative"
                class="min-w-max gap-1.5"
                on:click={() => {
                  isViewAll = true;
                  data = {
                    id: 1,
                    user: "",
                    botsName: "",
                    chatName: "",
                    message: "",
                    datetime: "",
                    checkbox: {
                      checked: false,
                    },
                  };
                }}
              >
                <p class="font-semibold">{$_("back")}</p>
              </Button>
              <Button
                outline
                color="red"
                class="min-w-max gap-1.5"
                on:click={() => {
                  const foundIndex = datas.findIndex((item) => item.id === data.id);
                  if (foundIndex !== -1) {
                    datas[foundIndex].checkbox.checked = !datas[foundIndex].checkbox.checked;
                  }
                  isAllLogs = false;
                  deleteModal = true;
                }}
              >
                <BrushCleaningIcon svgClass="w-5 h-5" />
                <p class="font-semibold">{$_("deleteLog")}</p>
              </Button>
            {/if}
          </div>
        </div>
        <Divider />
      </div>
      {#if isViewAll}
        {#if datas.length > 0}
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
        {/if}
      {:else}
        <div class="mt-8">
          <div class="w-full grid grid-cols-1 md:grid-cols-4">
            <div class="col-span-1">
              <Typography preset="Text md/Semibold" text={$_("generalInformation")} />
            </div>
            <div class="col-span-3 grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="col-span-2 flex flex-col gap-2">
                <Typography preset="Text sm/Medium" text={$_("botsName")} />
                <InputRw bind:value={data.botsName} onEdit={false} />
              </div>
              <div class="col-span-2 flex flex-col gap-2">
                <Typography preset="Text sm/Medium" text={$_("chatName")} />
                <InputRw bind:value={data.chatName} onEdit={false} />
              </div>
              <div class="col-span-1 flex flex-col gap-2">
                <Typography preset="Text sm/Medium" text={$_("userType")} />
                <InputRw bind:value={data.user} onEdit={false} />
              </div>
              <div class="col-span-1 flex flex-col gap-2">
                <Typography preset="Text sm/Medium" text={$_("datetime")} />
                <InputRw bind:value={data.datetime} onEdit={false} />
              </div>
            </div>
          </div>
          <Divider className="my-6" />
          <div class="w-full grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="col-span-1">
              <Typography preset="Text md/Semibold" text={$_("content")} />
            </div>
            <div class="col-span-3 flex flex-col gap-2">
              <Typography preset="Text sm/Medium" text={$_("message")} />
              <TextareaRw bind:value={data.message} rows={10} onEdit={false} />
            </div>
          </div>
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
