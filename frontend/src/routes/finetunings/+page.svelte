<script>
  import { Typography, SectionHeader, ActionCard, StateCard, call, triggerToast, triggerToastResponse, setSelectedConversation } from "$lib";
  import { format } from "date-fns";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  
  import DataTable from "$lib/components/DataTable.svelte";
  import DataTableLibraryCollection from "$lib/components/DataTableLibraryCollection.svelte";
  import DataTableLibraryCollectionBakrie from "$lib/components/DataTableLibraryCollectionBakrie.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { goto } from "$app/navigation";

  let actionItems = $state([]);
  let datas = $state([]);
  let datasHandwriting = $state([]);
  let datasConversation = $state([]);
  let sidebarData = $state({});
  let configs = $state({});
  let pagination = $state({
    page: 1,
    pageCount: 1,
    pageSize: 10,
    total: 10,
  });
  let paginationHandwriting = $state({
    page: 1,
    pageCount: 1,
    pageSize: 5,
    total: 10,
  });
  let tmpId = $state(null);
  let popupModal = $state(false);
  let gettingDatas = $state(true);

  // let datasConversation = $state([
  //   {
  //     id: 1,
  //     conversatin_name: "NewChat 21.02.2025",
  //     actions: {
  //       objects: [
  //         {
  //           href: "",
  //           name: "View",
  //         },
  //       ],
  //     },
  //   }
  // ]);

  function getActionItems() {
    actionItems = [
      {
        onclick: () => {
          // TODO
        },
        icon: "penToolPlus",
        actionName: "Finetune for Handwriting",
        desc: "desc Finetune for Handwriting",
      },
      {
        onclick: () => {
          // TODO
        },
        icon: "messageChatSquare",
        actionName: "Finetune for Conversation",
        desc: "desc Finetune for Conversation",
      },
    ];
  }

  const getDatasetsHandwriting = async (search) => {
    let token = localStorage.getItem("token");
    if (token) {
      if (search !== "") {
        paginationHandwriting.page = 1;
      }

      const response = await $call.call("client-dashboard", "get-datasets-handwriting", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: paginationHandwriting.page,
          page_size: paginationHandwriting.pageSize,
          search: search,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerToastResponse(response);
        return;
      }

      console.log("response getDatasetsHandwriting", response);

      let promise = response.result.data.map(async (data) => {
        let link = await $call.getPage("client-dashboard", "datasets-detail");

        return {
          id: data.id,
          fileName: data.attributes.name + "." + data.attributes.path.split(".").pop(),
          // fileext: data.attributes.path.split(".").pop(),
          filesize: data.attributes.file_size,
          createdAt: format(new Date(data.attributes.createdAt), "MMM d, yyyy"),
          actions: {
            objects: [
              {
                href: link + "?datasetId=" + data.id,
                icon: "eye",
                name: "View",
              },
            ],
          },
        };
      });

      datas = await Promise.all(promise);

      gettingDatas = false;
    }
  };

  const getListConversations = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "get-conversations", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: 1,
          page_size: 10,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerToast("error");
        return;
      }

      let promise = response.result.data.map(async (conversation) => {
        let updateAt = format(new Date(conversation.attributes.updatedAt), "yyyy-MM-dd");
        let title = conversation.attributes.title;
        let id = conversation.attributes.conversation_id;

        return {
          id: id,
          conversation_name: title,
          actions: {
            objects: [
              {
                onclick: () => {
                  setSelectedConversation(conversation.id);
                  goto("/finetunings/conversation");
                },
                icon: "eye",
                name: "View",
              },
            ],
          },
        };
      });

      datasConversation = await Promise.all(promise);

      console.log("datasConversation", datasConversation);

      gettingDatas = false;
    }
  };

  let next = async () => {
    // getCategories(pagination.page + 1);
  };

  let previous = async () => {
    // getCategories(pagination.page - 1);
  };

  function handleDeleteDataset(id) {
    tmpId = id;
    popupModal = true;
  }

  onMount(() => {
    getActionItems();
    getDatasetsHandwriting();
    getListConversations();
  })
</script>

<div class="flex flex-col gap-8 px-4 pt-4 pb-24">
  <Typography preset="Display sm/Semibold" text={$_("finetuning")} color="primary" />
  <div class="flex flex-col gap-4">
    <SectionHeader title={$_("quickActions")} />
    <ActionCard {actionItems} />
  </div>
  <!-- Conversation Finetune -->
  <div class="flex flex-col gap-4">
    <SectionHeader title={$_("conversationFinetuneHistory")} hrefAction="/finetunings/conversation" />
    {#if datas.length > 0}
      <DataTable bind:datas={datasConversation} bind:configs />
    {:else if gettingDatas}
      <StateCard
        illustration={IllustrationProgress}
        description={$_("fetchingDataPleaseWait")}
      />
    {:else}
      <StateCard
        illustration={IllustrationEmpty}
        description={$_("youDontHaveAnyDataYet")}
      />
    {/if}
  </div>
  <!-- Handwriting Finetune -->
  <div class="flex flex-col gap-4">
    <SectionHeader title={$_("handwritingFinetuneHistory")} hrefAction="/finetunings/handwriting" />
    {#if datas.length > 0}
      <DataTable
        bind:datas={datas}
        bind:configs
      />
    {:else if gettingDatas}
      <StateCard
        illustration={IllustrationProgress}
        description={$_("fetchingDataPleaseWait")}
      />
    {:else}
      <StateCard
        illustration={IllustrationEmpty}
        description={$_("youDontHaveAnyDataYet")}
      />
    {/if}
  </div>
</div>