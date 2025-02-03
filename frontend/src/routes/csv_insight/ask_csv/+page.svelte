<script>
  import { Divider, Typography, StateCard } from '$lib';
  import { _ } from 'svelte-i18n';  
  import Markdown from "svelte-exmarkdown";
  import { gfmPlugin } from "svelte-exmarkdown/gfm";
  import UserIcon from "$lib/icon/user-01.svg";
  import LightbulbIcon from "$lib/icon/lightbulb-04.svg";
  import File06Icon from '$lib/icon/file-06.svelte';
  import SendIcon from "$lib/icon/send-01.svg";
  import ChevronDownIcon from "$lib/icon/chevron-down.svelte";
  import DataTableListCollection from "$lib/components/DataTableListCollection.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { Button, Modal, Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell } from 'flowbite-svelte';
  import { onMount } from 'svelte';

  const plugins = [gfmPlugin()];

  let messageInput = $state("");
  let textarea = $state(null);
  let modalOpen = $state(false);
  let allData = $state([
    {
      id: 1,
      file_name: 'quarterly_revenue.csv',
    },
    {
      id: 2,
      file_name: 'ad_campaign_performance.csv',
    }
  ]);
  let datasChecked = $state([
    {
      id: 1,
      checked: false,
    },
    {
      id: 2,
      checked: false,
    }
  ]);
  let tmpDatasChecked = $state({});
  let configs = $state({hideSearch: true});
  let searchTerm = $state("");
  let gettingDatas = $state(false);
  let innerWidth = $state(0);
  let winInnerWidth = $derived(innerWidth);
  
  let pagination = $state({ page: 1, pageSize: 3, pageCount: 0, total: 0 });

  const autosize = () => {
    textarea.addEventListener("input", (e) => {
      setTimeout(() => {
        textarea.style.height = "auto";
        let scHeight = e.target.scrollHeight + 2; // Tambahkan sedikit untuk padding
        textarea.style.height = `${scHeight}px`;
      }, 0);
    });
  };

  const resetTextarea = () => {
    textarea.style.height = "auto";
  };

  async function sendMessageInput() {
    messageInput = "";
    resetTextarea();
    console.log("messageInput", messageInput);
  }

  const handleKeyDown = (event) => {
    if (winInnerWidth > 768) {
      if (event.key === "Enter") {
        if (!event.shiftKey) {
          event.preventDefault();
          if (
            messageInput.trim() !== ""
          ) {
            onclickSend();
          }
        }
      }
    }
  };

  

  const onclickSend = () => {
    if (messageInput.trim() !== "") {
      // sendMessageInput();
    }
  };

  async function next() {
    if (pagination.page >= pagination.pageCount) return;

    let pageNow = pagination.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasChecked,
      };
    } else {
      tmpDatasChecked[pageNow] = datasChecked;
    }

    pagination.page = pagination.page + 1;
    // TODO: create function getCsv();
  }

  async function previous() {
    if (pagination.page <= 1) return;

    let pageNow = pagination.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasChecked,
      };
    } else {
      tmpDatasChecked[pageNow] = datasChecked;
    }

    pagination.page = pagination.page - 1;
    // TODO: create function getCsv();
  }

  async function handleCheckAll(allChecked) {
    datasChecked = datasChecked.map((data) => ({
      id: data.id,
      checked: allChecked ? false : true,
    }));

    // console.log("Updated datasChecked (handleCheckAll):", datasChecked);
  }

  async function handleCheck(id) {
    if (!Array.isArray(datasChecked)) {
      datasChecked = [];
    }

    datasChecked = datasChecked.map((data) => {
      if (data.id === id) {
        return { id: data.id, checked: !data.checked };
      }
      return data;
    });

    let pageNow = pagination.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasChecked,
      };
    } else {
      tmpDatasChecked[pageNow] = datasChecked;
    }
    // console.log("datasChecked setelah pembaruan:", datasChecked);
  }

  function areAllUnchecked() {
    // Gabungkan data dari halaman sekarang dengan data yang tersimpan di tmpDatasChecked
    const allData = [
      ...datasChecked,
      ...Object.values(tmpDatasChecked).flat(),
    ];

    // Cek apakah semua data tidak ada yang checked
    return allData.every((item) => item.checked === false);
  }

  onMount(() => {
    if (textarea) {
      autosize();
    }
    const resizeHandler = () => {
      innerWidth = globalThis.innerWidth;
    };
    globalThis.addEventListener("resize", resizeHandler);
    resizeHandler(); // call once initially to set the correct value
    return () => {
      globalThis.removeEventListener("resize", resizeHandler);
    };
  });

</script>

<div class="flex items-center justify-between mb-3 md:mb-8 pt-5 px-4">
  <div class="flex flex-col gap-2">
    <div class="flex gap-4">
      <Typography
        preset="Display sm/Semibold"
        color="primary"
        text="Customer Feedback Data"
        propClass="fluid-xl"
      />
    </div>
  </div>

  <div class="flex relative gap-3">
    <Button
      color="light"
      class="min-w-max gap-1.5 px-3.5"
      on:click={() => {modalOpen = true;}}
    >
      <File06Icon svgClass="w-5 h-5 text-fg-brand-primary-light dark:text-fg-brand-primary-dark" />
      <div class="font-semibold text-fg-brand-primary-light dark:text-fg-brand-primary-dark">{$_("changeDataSource")}</div>
      <ChevronDownIcon svgClass="w-5 h-5 text-fg-brand-primary-light dark:text-fg-brand-primary-dark" />
    </Button>
  </div>
</div>

<Divider />

<div class="flex flex-col md:flex-row w-full pb-12">
  <div class="w-full md:w-1/2 pt-8 px-4 flex flex-col gap-6">
    <Typography preset="Text md/Semibold" text={$_("providedTable")} />
    <Table>
      <TableHead>
        <TableHeadCell>Product name</TableHeadCell>
        <TableHeadCell>Color</TableHeadCell>
        <TableHeadCell>Category</TableHeadCell>
        <TableHeadCell>Price</TableHeadCell>
      </TableHead>
      <TableBody tableBodyClass="divide-y">
        <TableBodyRow>
          <TableBodyCell>Apple MacBook Pro 17"</TableBodyCell>
          <TableBodyCell>Sliver</TableBodyCell>
          <TableBodyCell>Laptop</TableBodyCell>
          <TableBodyCell>$2999</TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell>Microsoft Surface Pro</TableBodyCell>
          <TableBodyCell>White</TableBodyCell>
          <TableBodyCell>Laptop PC</TableBodyCell>
          <TableBodyCell>$1999</TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell>Magic Mouse 2</TableBodyCell>
          <TableBodyCell>Black</TableBodyCell>
          <TableBodyCell>Accessories</TableBodyCell>
          <TableBodyCell>$99</TableBodyCell>
        </TableBodyRow>
      </TableBody>
    </Table>
  </div>
  <div class="w-full md:w-1/2 pt-8 px-4 flex flex-col gap-4">
    <Typography preset="Text md/Semibold" text={$_("csvInsight")} />
    <div class="flex gap-4 p-4">
      <div
        class="w-8 h-7 md:w-12 md:h-11 flex justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
      >
        <div>
          <img src={UserIcon} alt="User" class="w-4 h-4 md:w-6 md:h-6" />
        </div>
      </div>
      <div class="w-full">
        <div
          class="prose max-w-none font-medium fluid-base md:fluid-base text-tertiary-light dark:text-primary-dark"
        >
          <Markdown md={"Hi! I’m your smart CSV assistant. Drop a file, ask me anything, and let’s unlock the secrets hidden in your data!"} {plugins} />
        </div>
      </div>
    </div>

    <div class="flex gap-4 p-4">
      <div
        class="w-8 h-7 md:w-12 md:h-11 flex justify-center items-center bg-fg-brand-primary-light dark:bg-fg-brand-primary-dark rounded-full"
      >
        <div>
          <img src={LightbulbIcon} alt="User" class="w-4 h-4 md:w-6 md:h-6" />
        </div>
      </div>
      <div class="w-full">
        <div
          class="prose max-w-none font-medium fluid-base md:fluid-base text-tertiary-light dark:text-primary-dark"
        >
          <Markdown md={"Lorem ipsum dolor sit amet, consectetur adipiscing elit. In elit tellus, ullamcorper ut placerat sit amet, volutpat id erat. Duis elit justo, congue a tincidunt consequat, varius ut nibh. Ut eget mauris malesuada, dictum nulla sollicitudin, vulputate dolor. Integer aliquam viverra libero, eget tincidunt risus pharetra ut. Aenean lobortis, diam at tempor tempus, diam arcu fermentum mauris, id rutrum augue ante et tortor. Donec suscipit mi ornare nunc facilisis, et placerat ligula aliquam. Nam a ante faucibus, feugiat nibh eget, laoreet magna."} {plugins} />
        </div>
      </div>
    </div>

    <form
      onsubmit={messageInput !== "" ? sendMessageInput : () => {}}
      class="w-full p-6"
    >
      <div class="flex items-center">
        <div class="w-full flex gap-4">
          <!-- svelte-ignore event_directive_deprecated -->
          <textarea
            bind:this={textarea}
            bind:value={messageInput}
            onkeydown={handleKeyDown}
            id="chat"
            rows="1"
            class="block resize-none p-2.5 w-full fluid-sm text-brand-primary-light dark:text-primary-dark bg-white dark:bg-primary-dark rounded-xl border border-utility-brand-500-light dark:border-utility-brand-500-dark focus:border-utility-brand-500-light dark:focus:border-utility-brand-500-dark focus:ring-utility-brand-400-light dark:focus:ring-utility-brand-400-light dark:placeholder-gray-400"
            placeholder="Send a message..."
          ></textarea>
          <Button
            disabled={!messageInput}
            onclick={onclickSend}
            class="rounded-full h-11 w-11 !p-3"
          >
            <img src={SendIcon} alt="SendIcon" class="w-5 h-5" />
          </Button>
        </div>
      </div>
    </form>
  </div>
</div>

<Modal bind:open={modalOpen} outsideclose>
  <svelte:fragment slot="header">
    <div
      class={`grow pt-2`}
    >
      <Typography
        preset="Text md/Bold"
        text={$_("changeDataSource")}
        color="primary"
      />
      <Typography
        preset="Text sm/Regular"
        text={$_("changeYourDataSourcesFiles")}
        color="tertiary"
      />
    </div>
  </svelte:fragment>
  {#if allData.length > 0 || searchTerm !== ""}
    {#if !gettingDatas}
      <DataTableListCollection
        bind:datas={allData}
        bind:datasChecked
        bind:configs
        bind:gettingDatas
        bind:next
        bind:previous
        bind:handleCheckAll
        bind:handleCheck
      />
    {:else}
      <StateCard illustration={IllustrationProgress} description={$_("fetchingDataPleaseWait")} />
    {/if}
  {:else if gettingDatas}
    <StateCard illustration={IllustrationProgress} description={$_("fetchingDataPleaseWait")} />
  {:else}
    <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
  {/if}
  <svelte:fragment slot="footer">
    <div class="flex gap-4 w-full">
      <button
        onclick={() => {
          modalOpen = false;
        }}
        class="px-5">{$_("cancel")}</button
      >
      <Button
        on:click={() => {
          // TODO: send data to the API
        }}
        class="grow"
        disabled={areAllUnchecked()}>{$_("finish")}</Button
      >
    </div>
  </svelte:fragment>
</Modal>