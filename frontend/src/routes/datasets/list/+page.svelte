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
  import { onMount } from "svelte";
  import { Modal, Button, Input, Pagination } from "flowbite-svelte";
  import { Typography, triggerToast, StateCard, triggerLoading } from "$lib";
  import * as Icon from "flowbite-svelte-icons";
  import CalendarPlus02 from "$lib/icon/calendar-plus-02.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { _ } from "svelte-i18n";
  import { CR3Call } from "hexacomponents-svelte";
  import { format } from "date-fns";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import { goto } from "$app/navigation";

  import { call } from "$lib";

  let tmpRulesName = $state("");
  let addRulesModal = $state(false);

  let sidebarData = $state({});

  let configs = $state({});

  let allDatas = $state([]);
  let datas = $state([]);

  let gettingDatas = $state(false);

  let page = $state(1);
  let pageCount = $state(1);
  let searchTerm = $state("");

  const handleSearch = async (value) => {
    searchTerm = value;
    page = 1;

    console.log(searchTerm);

    await getListConversations();
  };

  async function getListConversations() {
    let token = localStorage.getItem("token");
    if (token) {
      gettingDatas = true;
      try {
        const response = await $call.call("client-dashboard", "get-conversations", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 10,
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

        // if (conversations.length > 0) {
        // response.result.data.forEach((conversation) => {
        //   let updateAt = format(new Date(conversation.attributes.updatedAt), "yyyy-MM-dd");
        //   let title = conversation.attributes.title;
        //   let id = conversation.attributes.conversation_id;

        //   if (!sidebarData[updateAt]) {
        //     sidebarData[updateAt] = [];
        //   }

        //   sidebarData[updateAt].push({
        //     title: title,
        //     id: id,
        //   });
        // });

        let link = await $call.getPage("client-dashboard", "datasets");
        let promise = conversations.map((conversation) => {
          return {
            id: conversation.attributes.conversation_id,
            updateAt: format(new Date(conversation.attributes.updatedAt), "yyyy-MM-dd"),
            title: conversation.attributes.title,
            actions: {
              objects: [
                {
                  href: link + "?id=" + conversation.attributes.conversation_id,
                  name: "Select",
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
  }

  async function nextPage() {
    if (page >= pageCount) return;
    page++;
    getListConversations();
  }

  async function prevPage() {
    console.log(page);
    if (page <= 1) return;
    page--;
    getListConversations();
  }

  onMount(() => {
    getListConversations();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex justify-between">
          <Typography preset="Display sm/Semibold" text={$_("allDatasets")} />
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
    </div>
  </div>
</div>
