<script>
  import {
    Button,
    Input,
    Textarea,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    Pagination,
  } from "flowbite-svelte";
  // import Datatable from "$lib/components/Datatable.svelte";
  import { _ } from "svelte-i18n";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";

  import { Typography, Divider, StateCard } from "$lib";
  import { onMount } from "svelte";
  
  import { page } from "$app/stores";

  let canEdit = $state(false);
  let dummyPageData = $state({
    title: "quarterly_revenue.csv",
    content: "Page Content",
  });
  let dummyDocumentData = $state({
    title: "quarterly_revenue.csv",
    description: "This is a document of revenue per quarter",
  });
  let dummyCsvDatas = $state([
    {
      id: 1,
      customer_id: 1,
      name: "John Doe",
      rating: 5,
      date: "2021-08-01",
    },
    {
      id: 2,
      customer_id: 2,
      name: "Jane Doe",
      rating: 4,
      date: "2021-08-02",
    },
    {
      id: 3,
      customer_id: 3,
      name: "John Smith",
      rating: 3,
      date: "2021-08-03",
    },
    {
      id: 4,
      customer_id: 4,
      name: "Jane Smith",
      rating: 2,
      date: "2021-08-04",
    },
    {
      id: 5,
      customer_id: 5,
      name: "John Doe",
      rating: 1,
      date: "2021-08-05",
    },
    {
      id: 6,
      customer_id: 6,
      name: "Jane Doe",
      rating: 5,
      date: "2021-08-06",
    },
    {
      id: 7,
      customer_id: 7,
      name: "John Smith",
      rating: 4,
      date: "2021-08-07",
    },
    {
      id: 8,
      customer_id: 8,
      name: "Jane Smith",
      rating: 3,
      date: "2021-08-08",
    },
    {
      id: 9,
      customer_id: 9,
      name: "John Doe",
      rating: 2,
      date: "2021-08-09",
    },
    {
      id: 10,
      customer_id: 10,
      name: "Jane Doe",
      rating: 1,
      date: "2021-08-10",
    },
    {
      id: 11,
      customer_id: 11,
      name: "John Smith",
      rating: 5,
      date: "2021-08-11",
    },
    {
      id: 12,
      customer_id: 12,
      name: "Jane Smith",
      rating: 4,
      date: "2021-08-12",
    },
    {
      id: 13,
      customer_id: 13,
      name: "John Doe",
      rating: 3,
      date: "2021-08-13",
    },
  ]);
  let currentArrayData = $state([]);
  // let dummyCsvConfigs = $state({
  //   pagination: {
  //     page: 1,
  //     pageCount: 1,
  //     pageSize: 10,
  //     total: 13,
  //   },
  // });

  let dummyCsvMeta = $state({ page: 1, pageSize: 10, pageCount: 2, total: 13 });

  const underscoreToSpace = (text) => {
    return text.replace(/_/g, " ");
  };

  async function next() {
    // if (paginationModal.page >= paginationModal.pageCount) return;
    // let pageNow = paginationModal.page;
    // if (tmpDatasChecked[pageNow] === undefined) {
    //   tmpDatasChecked = {
    //     ...tmpDatasChecked,
    //     [pageNow]: datasChecked,
    //   };
    // } else {
    //   tmpDatasChecked[pageNow] = datasChecked;
    // }
    // paginationModal.page = paginationModal.page + 1;
    // TODO: create function getCsv();
  }

  async function previous() {
    // if (paginationModal.page <= 1) return;
    // let pageNow = paginationModal.page;
    // if (tmpDatasChecked[pageNow] === undefined) {
    //   tmpDatasChecked = {
    //     ...tmpDatasChecked,
    //     [pageNow]: datasChecked,
    //   };
    // } else {
    //   tmpDatasChecked[pageNow] = datasChecked;
    // }
    // paginationModal.page = paginationModal.page - 1;
    // TODO: create function getCsv();
  }

  onMount(async () => {
    if (dummyCsvMeta) {
      currentArrayData = dummyCsvDatas.slice(
        (dummyCsvMeta.page - 1) * dummyCsvMeta.pageSize,
        dummyCsvMeta.page * dummyCsvMeta.pageSize,
      );
    }

    dummyPageData.title = $page.url.searchParams.get("docName");
  });
</script>

<div class="m-4">
  <div class="flex justify-between items-center">
    <Typography text={dummyPageData.title} preset="Display sm/Semibold" />
    <div class="flex gap-2">
      {#if canEdit}
        <Button color="alternative" on:click={() => (canEdit = false)}>
          <Typography text="Cancel" preset="Text sm/Semibold" color="tertiary" />
        </Button>
        <Button color="primary">
          <Typography text="Save" preset="Text sm/Semibold" color="primary-onbrand" />
        </Button>
      {:else}
        <Button color="alternative">
          <Typography text="Delete" preset="Text sm/Semibold" color="error-primary" />
        </Button>
        <Button color="alternative" on:click={() => (canEdit = true)}>
          <Typography
            text="edit"
            preset="Text sm/Semibold"
            propClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
          />
        </Button>
      {/if}
    </div>
  </div>
  <Divider className="my-4" />
  <!-- general information -->
  <div class="flex gap-6">
    <div class="w-60">
      <Typography preset="Text md/Semibold" color="secondary" text={$_("generalInformation")} />
    </div>
    <div class="flex flex-col grow gap-6">
      <div class="flex gap-5">
        <div class="flex flex-col grow gap-2">
          <Typography preset="Text sm/Semibold" color="secondary" text={$_("documentTitle")} />
          <Input
            readonly={!canEdit}
            bind:value={dummyPageData.title}
            class="bg-primary-light dark:bg-secondary-dark focus:none focus:ring-0 focus:border-gray-300 dark:focus:border-gray-600"
          />
        </div>
      </div>
      <div class="grid grid-cols-1">
        <div class="flex flex-col gap-2">
          <Typography preset="Text sm/Semibold" color="secondary" text={$_("description")} />
          <Textarea
            readonly={!canEdit}
            bind:value={dummyDocumentData.description}
            rows="4"
            class="bg-primary-light dark:bg-secondary-dark focus:none focus:ring-0 border-gray-300 dark:border-600 focus:!border-gray-300 dark:focus:!border-gray-600"
          />
        </div>
      </div>
    </div>
  </div>
  <Divider className="my-4" />
  <!-- Table Viewer -->
  <div class="flex gap-6">
    <div class="w-60">
      <Typography preset="Text md/Semibold" color="secondary" text={$_("generalInformation")} />
    </div>
    <div class="flex flex-col grow gap-6">
      <!-- <Datatable datas={dummyCsvDatas} configs={dummyCsvConfigs} /> -->
      {#if currentArrayData.length === 0}
        <StateCard illustration={IllustrationEmpty} description={$_("noDataAvailable")} />
      {:else}
        <div
          class="bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-700 rounded-xl pt-4"
        >
          <Table class="w-full">
            <TableHead>
              {#each Object.keys(currentArrayData[0]) as key}
                {#if key === "id"}
                  <!-- empty -->
                {:else}
                  <TableHeadCell>{underscoreToSpace(key)}</TableHeadCell>
                {/if}
              {/each}
            </TableHead>
            <TableBody tableBodyClass="divide-y">
              {#each currentArrayData as data}
                <TableBodyRow class="border-0 !border-gray-200 dark:!border-gray-700">
                  {#each Object.entries(data) as [key, value]}
                    {#if key === "id"}
                      <!-- empty -->
                    {:else}
                      <TableBodyCell>{value}</TableBodyCell>
                    {/if}
                  {/each}
                </TableBodyRow>
              {/each}
            </TableBody>
          </Table>
          <div class="flex justify-between items-center p-3 rounded-b-xl">
            <Typography
              preset="Text sm/Regular"
              text={`Showing ${dummyCsvMeta.page * dummyCsvMeta.pageSize - dummyCsvMeta.pageSize + 1} to ${
                dummyCsvMeta.page * dummyCsvMeta.pageSize
              } of ${dummyCsvMeta.total} entries`}
            />
            <Pagination
              on:next={next}
              on:previous={previous}
              activeClass="bg-hc-graylight-100 dark:bg-hc-graydark-700 dark:border-[#374151] dark:text-hc-graydark-300"
            >
              <span slot="prev">{$_("previous")}</span>
              <span slot="next">{$_("next")}</span>
            </Pagination>
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>
