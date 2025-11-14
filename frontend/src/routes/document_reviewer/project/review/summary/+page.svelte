<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { call, triggerToast, triggerLoading } from "$lib";
  import { ContractReviewDetailView } from "$lib";

  const isFrontend = process.env.NODE_ENV === "development";

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

  const goBackPage = async () => {
    const gotoUrl = isFrontend
      ? "/document_reviewer/project/review"
      : await $call.getPage("client-dashboard", "contract-project-review");
    goto(gotoUrl + `?projectDocumentId=${projectDocumentId}`);
  };

  onMount(async () => {
    projectDocumentId = $page.url.searchParams.get("projectDocumentId");
    if (projectDocumentId !== null) {
      await getProjectResult();
    } else {
      triggerToast("error");
      goto(
        isFrontend
          ? "/document_reviewer/project"
          : await $call.getPage("client-dashboard", "contract-project-list"),
      );
    }
  });
</script>

<ContractReviewDetailView
  title={$_("summaryOfFindings")}
  {isCSV}
  bind:datas
  bind:configs
  {next}
  {previous}
  {goToPage}
  {goBackPage}
/>
