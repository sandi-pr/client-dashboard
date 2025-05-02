<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer:
  - Dean Ramhan <dean@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
  - Taqie Fadlillah <taqief@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->
<script>
  import { _ } from "svelte-i18n";
  import { onMount } from "svelte";
  import { format } from "date-fns";
  import { Button } from "flowbite-svelte";
  import { goto } from "$app/navigation";

  import DataTable from "$lib/components/DataTable.svelte";
  import { Typography, triggerToast, StateCard } from "$lib";

  import DataFlow from "$lib/icon/dataflow-03.svelte";
  import MessagePlusSquareIcon from "$lib/icon/message-plus-square.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";

  import { call, check } from "$lib";

  let gettingDatas = $state(false);
  let disableButtonConversation = $state(true);
  let statusRequest = $state(null);

  let page = $state(1);
  let pageCount = $state(1);

  let searchTerm = $state("");

  let datas = $state([]);
  let allDatas = $state([]);

  let configs = $state({});

  const handleSearch = async (value) => {
    searchTerm = value;
    page = 1;

    console.log(searchTerm);

    await getListConversations();
  };

  const getListConversations = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      gettingDatas = true;
      try {
        const response = await $call.call("client-dashboard", "get-conversations", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 5,
            search: searchTerm,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerToast("error");
          gettingDatas = false;
          return;
        }

        console.log("response", response);

        let conversations = [];

        conversations = response.result.data;

        pageCount = response.result.meta.pagination.pageCount;
        configs = {
          pagination: response.result.data.length > 0 ? response.result.meta.pagination : null,
          search: true,
          searchHandler: handleSearch,
        };

        if (response.result.data.length < 1) {
          disableButtonConversation = true;
        } else {
          disableButtonConversation = false;
        }

        let link = await $call.getPage("client-dashboard", "finetunings-conversation-chat");
        let promise = conversations.map((conversation) => {
          return {
            id: conversation.attributes.conversation_id,
            conversationName: conversation.attributes.title,
            updateAt: format(new Date(conversation.attributes.updatedAt), "MMM d, yyyy"),
            actions: {
              objects: [
                {
                  href: link + "?id=" + conversation.attributes.conversation_id,
                  icon: "eye",
                  name: "View",
                },
              ],
            },
          };
        });
        allDatas = await Promise.all(promise);
        datas = allDatas.slice(0, 10);
        gettingDatas = false;
        // }
      } catch (error) {
        console.log("error", error);
        gettingDatas = false;
      }
    }
  };

  async function getStatusRequest() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        const response = await $call.call("client-dashboard", "get-status-request-fine-tuning", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status === 200) {
          if (response.result.data.length > 0) {
            statusRequest = response.result.data[0].attributes.status;
          }
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function processDataset() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        const response = await $call.call("client-dashboard", "request-fine-tuning", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          statusRequest = "REQUEST";
          triggerToast("created");
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  const nextPage = async () => {
    if (page >= pageCount) return;
    page++;
    getListConversations();
  };

  const prevPage = async () => {
    console.log(page);
    if (page <= 1) return;
    page--;
    getListConversations();
  };

  let hrefConversationChat = $state("");

  async function getHrefAction(href) {
    let link = "";
    if (href === "conversation-chat") {
      link = await $call.getPage("client-dashboard", "finetunings-conversation-chat");
      return link;
    }
  }

  onMount(async () => {
    if (check("aiManagement.datasets:view") === false) {
      window.history.back();
      return;
    }
    hrefConversationChat = await getHrefAction("conversation-chat");
    getHrefAction("conversation-chat");
    getListConversations();
    getStatusRequest();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 gap-8">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography
            preset="Display sm/Semibold"
            text={$_("conversationFinetuningHistory")}
            propClass="fluid-lg md:fluid-xl"
          />
          <Button
            outline
            on:click={() => {
              goto(hrefConversationChat);
            }}
          >
            <MessagePlusSquareIcon svgClass="w-5 h-5 md:me-1" />
            <div class="hidden md:block">{$_("addNewFinetuneConversation")}</div>
          </Button>
        </div>

        <!-- table latest -->
        {#if datas.length > 0 || searchTerm !== ""}
          <DataTable bind:datas bind:configs next={nextPage} previous={prevPage} />
        {:else if gettingDatas}
          <StateCard
            illustration={IllustrationProgress}
            description={$_("fetchingDataPleaseWait")}
          />
        {:else}
          <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
        {/if}
      </div>

      {#if datas.length > 0}
        <Button
          disabled={statusRequest === "process" ||
            !check("aiManagement.datasets:process") ||
            disableButtonConversation}
          on:click={() => {
            processDataset();
          }}
          outline
          class="!flex hover:bg-primary-light dark:hover:bg-primary-dark sticky bottom-0 w-full h-12 rounded-2xl"
        >
          <DataFlow
            svgClass="w-5 h-5 text-utility-brand-700-light dark:text-utility-brand-700-dark"
          />
          &nbsp; &nbsp;
          <Typography preset="Text md/Semibold" text={$_("process")} color="brand-secondary" />
        </Button>
      {/if}
    </div>
  </div>
</div>
