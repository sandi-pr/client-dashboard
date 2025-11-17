<script>
  import { onMount } from "svelte";
  import { Button } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { call, triggerToast, triggerLoading } from "$lib";
  import { Typography, Divider, DataTableAIReview } from "$lib";

  const isFrontend = process.env.NODE_ENV === "development";

  let summaryType = $state("");
  let projectDocumentId = $state(null);

  let isCSV = $state(true);

  let datas = $state([]);
  let allDatas = $state([]);

  let configs = $state({
    pagination: {
      page: 1,
      pageSize: 10,
      pageCount: 1,
      total: 10,
    },
  });

  const getProjectResult = async () => {
    try {
      triggerLoading(true);

      let token = localStorage.getItem("token");
      if (!token) {
        throw new Error("No token found");
      }

      if (datas.length > 0) {
        datas = [];
      }

      let response = await $call.call("client-dashboard", "get-contract-project-result", {
        URLQuery: {},
        URLParams: {
          id: projectDocumentId,
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        let error = new Error("Failed to fetch project result");
        error.response = response;
        throw error;
      }

      // console.log("response", response);

      allDatas = buildDatasFromSOF(response.result.data[0].attributes.summary_of_findings);
      datas = allDatas.slice(0, 10);
      configs.pagination.total = allDatas.length;
      configs.pagination.pageCount = Math.ceil(allDatas.length / configs.pagination.pageSize);

      triggerLoading(false);
    } catch (error) {
      console.log(error);
      triggerLoading(false);
      triggerToast("error");
    }
  };

  const buildDatasFromSOF = (resultArray) => {
    return resultArray.map((item, index) => {
      return {
        id: index + 1,
        no: item.no,
        row: item.row,
        context: item.context,
        summary: item.summary,
      };
    });
  };

  const formatRowData = (row) => {
    if (row == null) return "";
    if (typeof row === "string") return row;
    if (Array.isArray(row)) {
      return row
        .map((item) => (typeof item === "object" ? formatRowData(item) : String(item)))
        .join("  \n");
    }
    if (typeof row === "object") {
      return Object.entries(row)
        .map(([k, v]) => `${k}: ${formatRowData(v)}`)
        .join("  \n");
    }

    return String(row);
  };

  const next = () => {
    if (configs.pagination.page < configs.pagination.pageCount) {
      configs.pagination.page += 1;
      const start = (configs.pagination.page - 1) * configs.pagination.pageSize;
      const end = start + configs.pagination.pageSize;
      datas = allDatas.slice(start, end);
    }
  };

  const previous = () => {
    if (configs.pagination.page > 1) {
      configs.pagination.page -= 1;
      const start = (configs.pagination.page - 1) * configs.pagination.pageSize;
      const end = start + configs.pagination.pageSize;
      datas = allDatas.slice(start, end);
    }
  };

  const goToPage = (pageNumber) => {
    if (pageNumber >= 1 && pageNumber <= configs.pagination.pageCount) {
      configs.pagination.page = pageNumber;
      const start = (configs.pagination.page - 1) * configs.pagination.pageSize;
      const end = start + configs.pagination.pageSize;
      datas = allDatas.slice(start, end);
    }
  };

  onMount(async () => {
    summaryType = $page.url.searchParams.get("type");
    projectDocumentId = $page.url.searchParams.get("projectDocumentId");
    if (projectDocumentId !== null) {
      await getProjectResult();
    } else {
      triggerToast("error");
      goto(
        isFrontend
          ? "/document_reviewer/report/project"
          : await $call.getPage("client-dashboard", "contract-report-project-list"),
      );
    }
  });
</script>

<div class="flex flex-col px-4 pt-6 gap-6 pb-12">
  <div class="flex justify-between">
    <div class="flex flex-col gap-1">
      <Typography preset="Display sm/Semibold" text={summaryType + " Data"} color="primary" />
    </div>
    <Button
      class="h-min"
      color="light"
      outline
      on:click={async () => {
        // const gotoUrl = isFrontend
        //   ? "/document_reviewer/report/project"
        //   : await $call.getPage("client-dashboard", "contract-report-project-list");
        // goto(gotoUrl);

        window.history.back();
      }}
    >
      <Typography preset="Text md/Semibold" text={$_("goBack")} color="primary" />
    </Button>
  </div>
  <Divider />
  <div class="flex gap-3">
    <div class="flex flex-col grow gap-6">
      <div class="flex flex-col gap-6 w-full">
        {#if isCSV && datas.length > 0}
          <DataTableAIReview bind:datas bind:configs {next} {previous} {goToPage} />
        {/if}
      </div>
    </div>
  </div>
</div>
