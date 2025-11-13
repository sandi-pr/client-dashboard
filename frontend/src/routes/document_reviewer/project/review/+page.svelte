<script>
  import { onMount } from "svelte";
  import { Button, Modal } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { call, triggerToast, triggerLoading, getIcons } from "$lib";
  import { Typography, ListMarkdown, PdfjsViewer, Divider, DataTableAIReview } from "$lib";

  const CheckCircleIcon = getIcons("checkCircle");
  const AlertTriangleIcon = getIcons("alertTriangle");
  const InfoCircleIcon = getIcons("infoCircle");

  const isFrontend = process.env.NODE_ENV === "development";

  let selectedFile = $state(null);
  let projectDocumentId = $state(null);

  let isCSV = $state(false);
  let editModal = $state(false);
  let documentViewModal = $state(false);

  let tmpText = $state("");
  let tmpTitle = $state("");
  let projectName = $state("");
  let generatedOn = $state("");

  let datas = $state([]);
  let allDatas = $state([]);
  let datasSOF = $state([]);
  let tmpResult = $state([]);
  let datasExtraction = $state([]);

  let referenceDocumentNames = $state({});
  let CSVResultSummary = $state({
    valid: 0,
    invalid: 0,
    total: 0,
    need_review: 0,
    resolved: 0,
    unresolved: 0,
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

      if (CSVResultSummary.total > 0) {
        CSVResultSummary.valid = 0;
        CSVResultSummary.invalid = 0;
        CSVResultSummary.total = 0;
        CSVResultSummary.need_review = 0;
        CSVResultSummary.resolved = 0;
        CSVResultSummary.unresolved = 0;
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

      let projectId = response.result.data[0].attributes.project_id;

      const response2 = await $call.call("client-dashboard", "get-contract-project", {
        URLQuery: {},
        URLParams: {
          id: projectId.toString(),
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response2.status !== 200) {
        let error = new Error("Failed to fetch project details");
        error.response = response2;
        throw error;
      }

      tmpResult = response.result.data[0].attributes.result;
      generatedOn = response.result.data[0].attributes.createdAt;
      isCSV = tmpResult[0]?.generated?.tabular_result;

      projectName = response2.result.data[0].attributes.project_name;

      if (!isCSV) {
        const referenceDocumentIds = tmpResult
          .flatMap((section) =>
            section.generated.results_query.flatMap((item) => item.reference_document_ids),
          )
          .filter((value, index, self) => self.indexOf(value) === index);

        // console.log("referenceDocumentIds", referenceDocumentIds);

        await buildReferenceDocumentNames(referenceDocumentIds, token);

        datas = buildDatasFromResultsQuery(response.result.data[0].attributes.result);
        // console.log("datas", datas);

        datasExtraction = response.result.data[0].attributes.result.map((item) => {
          return {
            name: item.title,
            result: removeCodeBlocks(item.generated.final_result).replace(/\n/g, "  \n"),
          };
        });
        // console.log("datasExtraction", datasExtraction);
      } else {
        // accumulate CSV summary
        response.result.data[0].attributes.result.forEach((item) => {
          CSVResultSummary.valid += item.generated.tabular_result.summary.valid;
          CSVResultSummary.invalid += item.generated.tabular_result.summary.invalid;
          CSVResultSummary.total += item.generated.tabular_result.summary.total;
          CSVResultSummary.need_review += item.generated.tabular_result.summary.need_review;
          CSVResultSummary.resolved += item.generated.tabular_result.summary.resolved;
          CSVResultSummary.unresolved += item.generated.tabular_result.summary.unresolved;
        });
        // console.log("CSVResultSummary", CSVResultSummary);

        allDatas = buildDatasFromTabularLogs(response.result.data[0].attributes.result);
        datas = allDatas.slice(0, 5);
        // console.log("allDatas", allDatas);

        if (response.result.data[0].attributes.summary_of_findings) {
          datasSOF = buildDatasFromSOF(response.result.data[0].attributes.summary_of_findings);
          // console.log("datasSOF", datasSOF);
          datasSOF = datasSOF.slice(0, 3);
        }
      }

      triggerLoading(false);
    } catch (error) {
      console.log(error);
      triggerLoading(false);
      triggerToast("error");
    }
  };

  const updateResult = async () => {
    editModal = false;
    triggerLoading(true);
    let token = localStorage.getItem("token");
    if (token !== null) {
      let index = tmpResult.findIndex((data) => data.title === tmpTitle);
      if (index !== -1) {
        tmpResult[index].text = tmpText;
      }

      try {
        let response = await $call.call("client-dashboard", "update-contract-project-result", {
          URLQuery: {},
          URLParams: {
            id: projectDocumentId,
          },
          Params: {
            result: JSON.stringify(tmpResult),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          editModal = true;
          triggerToast("error");
          console.log(response);
          return;
        }

        editModal = false;
        await getProjectResult();

        triggerToast("updated");
        triggerLoading(false);
      } catch (error) {
        console.log(error);
      }
    }
  };

  const buildReferenceDocumentNames = async (ids, token) => {
    for (const referenceDocumentId of ids) {
      const response3 = await $call.call("client-dashboard", "get-reference-document-by-id", {
        URLQuery: {},
        URLParams: {
          id: referenceDocumentId.toString(),
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response3.status !== 200) {
        continue;
      }

      if (response3.result.data.length === 0) {
        continue;
      } else {
        referenceDocumentNames[referenceDocumentId] =
          response3.result.data[0].attributes.document_name;
      }
    }
  };

  const buildDatasFromResultsQuery = (resultArray) => {
    let idCounter = 1;
    return resultArray.flatMap((item) =>
      item.generated.results_query.map((generatedItem) => {
        const documentId = generatedItem.reference_document_ids[0];

        return {
          no: idCounter,
          id: idCounter++,
          section: {
            name: item.title,
            prompt: generatedItem.section_prompt,
            result: generatedItem.section_results[0],
          },
          reference_documents: {
            id: documentId,
            name: referenceDocumentNames[documentId] || "Reference Document " + documentId + ".pdf",
            prompt: generatedItem.reference_prompt,
            result: generatedItem.reference_results[0],
          },
          ai_review: removeCodeBlocks(generatedItem.ai_review).replace(/\n/g, "  \n"),
          ai_suggestion: removeCodeBlocks(generatedItem.ai_suggestion).replace(/\n/g, "  \n"),
        };
      }),
    );
  };

  const buildDatasFromTabularLogs = (resultArray) => {
    let datasIndex = 1;
    return resultArray.flatMap((item) =>
      item.generated.tabular_result.logs.map((log) => {
        const documentId = log.reference_document.id;

        return {
          id: datasIndex++,
          no: log.no,
          section: {
            name: item.title,
            prompt: log.section?.prompt,
            result: formatRowData(log.row_data),
          },
          reference_documents: {
            id: documentId,
            name: log.reference_document.name,
            prompt: log.reference_document.prompt,
            result: log.reference_document.result,
          },
          ai_extraction: log.ai_extraction,
          ai_suggestion: log.status_icon + " " + log.status_label + ": " + log.ai_suggestion,
        };
      }),
    );
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

  const removeCodeBlocks = (text) => {
    return text
      .replace(/```/g, "") // Hapus semua backticks triple
      .replace(/“/g, "") // Hapus semua tanda kutip buka
      .replace(/”/g, "") // Hapus semua tanda kutip tutup
      .replace(/\s{3,}/g, "  ") // Gantikan lebih dari tiga spasi dengan dua spasi
      .trim();
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

  const downloadAsCSV = async () => {
    if (!isCSV) {
      let csvContent = "data:text/csv;charset=utf-8,";

      // 1. Header row: all names
      const headerRow = datasExtraction
        .map((data) => `"${data.name.replace(/"/g, '""')}"`)
        .join(",");
      csvContent += headerRow + "\n";

      // 2. Data row: all results
      const dataRow = datasExtraction
        .map((data) => `"${data.result.replace(/"/g, '""')}"`)
        .join(",");
      csvContent += dataRow + "\n";

      const encodedUri = encodeURI(csvContent);
      const link = document.createElement("a");
      link.setAttribute("href", encodedUri);
      link.setAttribute("download", projectName + ".csv");
      document.body.appendChild(link); // Required for FF

      link.click();
      document.body.removeChild(link);
    } else {
      let csvContent = "data:text/csv;charset=utf-8,";

      // 1. Header row
      const headerRow = [
        "No",
        "Section Name",
        "Section Prompt",
        "Section Result",
        "Reference Document Name",
        "Reference Document Prompt",
        "Reference Document Result",
        "AI Extraction",
        "AI Suggestion",
      ].join(",");
      csvContent += headerRow + "\n";

      // 2. Data rows
      allDatas.forEach((data) => {
        const dataRow = [
          data.no,
          `"${data.section.name.replace(/"/g, '""')}"`,
          `"${data.section.prompt.replace(/"/g, '""')}"`,
          `"${data.section.result.replace(/"/g, '""')}"`,
          `"${data.reference_documents.name.replace(/"/g, '""')}"`,
          `"${data.reference_documents.prompt.replace(/"/g, '""')}"`,
          `"${data.reference_documents.result.replace(/"/g, '""')}"`,
          `"${data.ai_extraction.replace(/"/g, '""')}"`,
          `"${data.ai_suggestion.replace(/"/g, '""')}"`,
        ].join(",");
        csvContent += dataRow + "\n";
      });

      // console.log("csvContent", csvContent);

      const encodedUri = encodeURI(csvContent);
      const link = document.createElement("a");
      link.setAttribute("href", encodedUri);
      link.setAttribute("download", projectName + ".csv");
      document.body.appendChild(link); // Required for FF

      link.click();
      document.body.removeChild(link);
    }
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

<div class="flex flex-col px-4 pt-6 gap-6 pb-12">
  <!-- section header -->
  <div class="flex justify-between">
    <div class="flex flex-col gap-1">
      <Typography preset="Display sm/Semibold" text={$_("aiReview")} color="primary" />
      <Typography preset="Text sm/Medium" text={projectName} color="secondary" />
      <Typography
        preset="Text md/Regular"
        text={$_("generatedOn") + ": " + new Date(generatedOn).toLocaleString()}
        color="tertiary"
      />
    </div>
    <div class="flex items-center gap-3">
      <Button
        class="h-min"
        on:click={async () => {
          const gotoUrl = isFrontend
            ? "/document_reviewer/project/detail"
            : await $call.getPage("client-dashboard", "contract-project-detail");
          goto(gotoUrl + `?projectId=${$page.url.searchParams.get("projectId")}`);
        }}
      >
        <Typography preset="Text md/Semibold" text={$_("backToProjectList")} color="primary" />
      </Button>
      <Button class="h-min" color="light" outline on:click={downloadAsCSV}>
        <Typography preset="Text md/Semibold" text={$_("downloadAsCSV")} color="primary" />
      </Button>
    </div>
  </div>
  <Divider />
  <div class="flex gap-3">
    <div class="flex flex-col grow gap-6">
      <div class="flex flex-col gap-6 w-full">
        {#if CSVResultSummary.total > 0}
          <Typography preset="Text md/Semibold" text={$_("overallSummary")} color="primary" />
          <div class="grid grid-cols-3 gap-4">
            {#each [...Array(3).keys()] as item}
              <div
                class="rounded-lg border border-primary-light dark:border-primary-dark p-6 bg-primary-light dark:bg-primary-dark flex flex-col gap-2"
              >
                <div class="flex items-center gap-2">
                  {#if item === 0}
                    <CheckCircleIcon
                      svgClass="w-5 h-5 text-fg-success-primary-light dark:text-fg-success-primary-dark"
                    />
                  {:else if item === 1}
                    <AlertTriangleIcon
                      svgClass="w-5 h-5 text-fg-warning-primary-light dark:text-fg-warning-primary-dark"
                    />
                  {:else}
                    <InfoCircleIcon
                      svgClass="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark"
                    />
                  {/if}
                  <Typography
                    preset="Text md/Regular"
                    text={item === 0
                      ? $_("validData")
                      : item === 1
                        ? $_("invalidData")
                        : $_("totalData")}
                    color="secondary"
                  />
                </div>
                <Typography
                  preset="Display sm/Bold"
                  text={item === 0
                    ? CSVResultSummary.valid
                    : item === 1
                      ? CSVResultSummary.invalid
                      : CSVResultSummary.total}
                  color="primary"
                />
              </div>
            {/each}
          </div>
        {/if}
        {#if datasExtraction.length > 0}
          <Divider />
          <Typography preset="Text md/Semibold" text={$_("extractionResults")} color="primary" />
          <div
            class="flex flex-col justify-center bg-tertiary-light dark:bg-tertiary-dark rounded-lg p-4"
          >
            <ListMarkdown id="content-to-download" bind:datas={datasExtraction} />
          </div>
        {/if}
        {#if datasSOF.length > 0}
          <Divider />
          <div class="flex justify-between">
            <Typography preset="Text md/Semibold" text={$_("summaryOfFindings")} color="primary" />
            <Button
              outline
              color="alternative"
              on:click={async () => {
                const gotoUrl = isFrontend
                  ? "/document_reviewer/project/review/summary"
                  : await $call.getPage("client-dashboard", "contract-project-review-summary");
                goto(gotoUrl + `?projectDocumentId=${projectDocumentId}`);
              }}
            >
              <Typography preset="Text md/Semibold" text={$_("viewAll")} color="primary" />
            </Button>
          </div>
          <DataTableAIReview bind:datas={datasSOF} />
        {/if}
        {#if (Object.keys(referenceDocumentNames).length > 0 || isCSV) && datas.length > 0}
          <Divider />
          <div class="flex justify-between">
            <Typography
              preset="Text md/Semibold"
              text={isCSV ? $_("AIScanningLogs") : $_("reviewerSummary")}
              color="primary"
            />
            {#if isCSV}
              <Button
                outline
                color="alternative"
                on:click={async () => {
                  const gotoUrl = isFrontend
                    ? "/document_reviewer/project/review/logs"
                    : await $call.getPage("client-dashboard", "contract-project-review-logs");
                  goto(gotoUrl + `?projectDocumentId=${projectDocumentId}`);
                }}
              >
                <Typography preset="Text md/Semibold" text={$_("viewAll")} color="primary" />
              </Button>
            {/if}
          </div>
          <DataTableAIReview bind:datas />
        {/if}
      </div>
    </div>
  </div>
</div>

<Modal title={$_("documentViewer")} bind:open={documentViewModal} autoclose outsideclose>
  {#if selectedFile && (selectedFile.name.includes(".pdf") || selectedFile.name.includes(".PDF"))}
    <div class="flex justify-center">
      <PdfjsViewer pdfFile={selectedFile.base64} />
    </div>
  {/if}
</Modal>

<Modal title={$_("editResult")} bind:open={editModal} size="lg">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={tmpTitle} />
      <Input bind:value={tmpText} placeholder={$_("whatIsYourChangeResult")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={() => {
        updateResult();
      }}
      class="w-full"
      disabled={buttonDisabled}>{$_("change")}</Button
    >
  </svelte:fragment>
</Modal>
