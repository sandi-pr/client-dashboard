<script>
  import { onMount } from "svelte";
  import { Button, Modal } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { call, triggerToast, triggerLoading, getIcons } from "$lib";
  import { Typography, PdfjsViewer, Divider, DataTableAIReview, ReportCard } from "$lib";

  const isFrontend = process.env.NODE_ENV === "development";

  let projectId = $state(null);
  let selectedFile = $state(null);
  let projectDocumentId = $state(null);

  let isCSV = $state(false);
  let sofModal = $state(false);
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
  let actionItems = $state([]);
  let datasExtraction = $state([]);

  let referenceDocumentNames = $state({});
  let CSVResultSummary = $state({
    valid: 0,
    invalid: 0,
    total: 0,
    // need_review: 0,
    // resolved: 0,
    // unresolved: 0,
  });

  const normalizeResult = (raw) => {
    const result = [];
    const summary = [];
    let total = 0;
    let invalid = 0;

    for (const [cif, forms] of Object.entries(raw)) {
      for (const [form, content] of Object.entries(forms)) {

        // Hanya form yang punya row_data
        if (!form.startsWith("F") || !content?.row_data || content.row_data.length === 0) continue;

        total++;

        const firstRow = content.row_data[0];

        const logItem = {
          no: 1,
          row_data: firstRow,
          reference_document: {
            id: null,
            name: "No Reference",
            prompt: "",
            result: ""
          },
          ai_extraction: "",
          status_icon: "⚠️",
          status_label: "Needs Review",
          ai_suggestion:
            forms?.finding?.completeness?.[form] ||
            forms?.finding?.comparison?.[form]?.reason ||
            "Review required"
        };

        // ============= BUILD RESULT ENTRY UI =============
        result.push({
          title: `CIF ${cif} - ${form}`,
          generated: {
            results_query: [],
            tabular_result: { logs: [logItem] }
          }
        });

        invalid++;

        // ============= BUILD SUMMARY SECTION =============
        if (forms.finding) {
          if (forms.finding.comparison?.[form]) {
            summary.push({
              row: `CIF ${cif} - ${form}`,
              context: "Comparison",
              summary: forms.finding.comparison[form].reason
            });
          }

          if (forms.finding.accuracy?.[form]) {
            const firstKey = Object.keys(forms.finding.accuracy[form])[0];
            summary.push({
              row: `CIF ${cif} - ${form}`,
              context: "Accuracy",
              summary: forms.finding.accuracy[form][firstKey]
            });
          }

          if (forms.finding.completeness?.[form]) {
            summary.push({
              row: `CIF ${cif} - ${form}`,
              context: "Completeness",
              summary: forms.finding.completeness[form]
            });
          }
        }
      }
    }

    return {
      result,
      summary_of_findings: summary.map((s, i) => ({ id: i + 1, ...s })),
      csv_summary: { valid: total - invalid, invalid, total }
    };
  };


  const getProjectResult = async () => {
    try {
      triggerLoading(true);

      let token = localStorage.getItem("token");
      if (!token) throw new Error("No token found");

      const response = await $call.call("client-dashboard", "get-contract-audit-result", {
        URLQuery: {},
        URLParams: { id: projectId },
        Params: {},
        Headers: { Authorization: "Bearer " + token }
      });

      if (response.status !== 200) throw new Error("Failed to fetch");

      const raw = response.result.attributes.result;

      // 🔥 NORMALIZER RUN HERE
      const normalized = normalizeResult(raw);

      // 🔥 Set ke UI state
      projectName = response.result.attributes.project_name || "Audit Result";
      generatedOn = response.result.attributes.createdAt;

      datasSOF = normalized.summary_of_findings;
      CSVResultSummary = normalized.csv_summary;

      // Build table UI format
      datas = buildDatasFromTabularLogs(normalized.result);
      console.log("==================", datas);
      allDatas = datas;

      getActionItems();
      triggerLoading(false);

    } catch (err) {
      console.error(err);
      triggerToast("error");
      triggerLoading(false);
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

  const getActionItems = () => {
    actionItems = [
      {
        icon: "checkCircle",
        iconColor: "fg-success-primary",
        actionName: CSVResultSummary.valid.toString(),
        desc: "compliantData",
        subtext: "20%",
        changeType: "increase",
      },
      {
        icon: "alertTriangle",
        iconColor: "fg-warning-primary",
        actionName: CSVResultSummary.invalid.toString(),
        desc: "nonCompliantData",
        subtext: "50%",
        changeType: "decrease",
      },
      {
        icon: "infoCircle",
        iconColor: "fg-secondary",
        actionName: CSVResultSummary.total.toString(),
        desc: "totalData",
        subtext: "40%",
        changeType: "increase",
      },
    ];
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

  // const downloadAsCSV = async () => {
  //   if (!isCSV) {
  //     let csvContent = "data:text/csv;charset=utf-8,";

  //     // 1. Header row: all names
  //     const headerRow = datasExtraction
  //       .map((data) => `"${data.name.replace(/"/g, '""')}"`)
  //       .join(",");
  //     csvContent += headerRow + "\n";

  //     // 2. Data row: all results
  //     const dataRow = datasExtraction
  //       .map((data) => `"${data.result.replace(/"/g, '""')}"`)
  //       .join(",");
  //     csvContent += dataRow + "\n";

  //     const encodedUri = encodeURI(csvContent);
  //     const link = document.createElement("a");
  //     link.setAttribute("href", encodedUri);
  //     link.setAttribute("download", projectName + ".csv");
  //     document.body.appendChild(link); // Required for FF

  //     link.click();
  //     document.body.removeChild(link);
  //   } else {
  //     let csvContent = "data:text/csv;charset=utf-8,";

  //     // 1. Header row
  //     const headerRow = [
  //       "No",
  //       "Section Name",
  //       "Section Prompt",
  //       "Section Result",
  //       "Reference Document Name",
  //       "Reference Document Prompt",
  //       "Reference Document Result",
  //       "AI Extraction",
  //       "AI Suggestion",
  //     ].join(",");
  //     csvContent += headerRow + "\n";

  //     // 2. Data rows
  //     allDatas.forEach((data) => {
  //       const dataRow = [
  //         data.no,
  //         `"${data.section.name.replace(/"/g, '""')}"`,
  //         `"${data.section.prompt.replace(/"/g, '""')}"`,
  //         `"${data.section.result.replace(/"/g, '""')}"`,
  //         `"${data.reference_documents.name.replace(/"/g, '""')}"`,
  //         `"${data.reference_documents.prompt.replace(/"/g, '""')}"`,
  //         `"${data.reference_documents.result.replace(/"/g, '""')}"`,
  //         `"${data.ai_extraction.replace(/"/g, '""')}"`,
  //         `"${data.ai_suggestion.replace(/"/g, '""')}"`,
  //       ].join(",");
  //       csvContent += dataRow + "\n";
  //     });

  //     // console.log("csvContent", csvContent);

  //     const encodedUri = encodeURI(csvContent);
  //     const link = document.createElement("a");
  //     link.setAttribute("href", encodedUri);
  //     link.setAttribute("download", projectName + ".csv");
  //     document.body.appendChild(link); // Required for FF

  //     link.click();
  //     document.body.removeChild(link);
  //   }
  // };

  onMount(async () => {
    projectId = $page.url.searchParams.get("projectId");
    if (projectId !== null) {
      await getProjectResult();
    } else {
      triggerToast("error");
      goto(
        isFrontend
          ? "/document_reviewer/audit_scan"
          : await $call.getPage("client-dashboard", "contract-audit-scan"),
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
          console.log("TODO: download");
        }}
      >
        <Typography preset="Text md/Semibold" text={$_("download")} color="fg-white" />
      </Button>
      <Button
        class="h-min"
        on:click={async () => {
          const gotoUrl = isFrontend
            ? "/document_reviewer/audit_scan/detail"
            : await $call.getPage("client-dashboard", "contract-audit-scan-detail");
          goto(gotoUrl + `?projectId=${projectId}`);
        }}
      >
        <Typography preset="Text md/Semibold" text={$_("backToProjectList")} color="fg-white" />
      </Button>
      <!-- <Button class="h-min" color="light" outline on:click={downloadAsCSV}>
        <Typography preset="Text md/Semibold" text={$_("downloadAsCSV")} color="fg-secondary" />
      </Button> -->
    </div>
  </div>
  <div class="flex gap-3">
    <div class="flex flex-col grow gap-6">
      <div class="flex flex-col gap-6 w-full">
        {#if actionItems.length > 0}
          <div>
            {#if actionItems.some((item) => item.disabled !== true)}
              <div
                class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4"
              >
                <Typography preset="Text lg/Semibold" color="primary" text={$_("overallSummary")} />
              </div>
            {/if}
            <ReportCard {actionItems} />
          </div>
        {/if}
        {#if datasSOF.length > 0}
          <Divider />
          <div class="flex justify-between">
            <Typography preset="Text md/Semibold" text={$_("summaryOfFindings")} color="primary" />
            <Button
              outline
              color="light"
              on:click={() => {
                sofModal = true;
              }}
            >
              <Typography preset="Text md/Semibold" text={$_("viewAll")} color="fg-secondary" />
            </Button>
          </div>
          <DataTableAIReview bind:datas={datasSOF} />
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

<Modal title={$_("editResult")} bind:open={sofModal} size="lg">
  <svelte:fragment slot="header">
    <div class="flex flex-col gap-2">
      <Typography preset="Text lg/Semibold" text={$_("aiSuggestion")} color="primary" />
      <Typography preset="Text sm/Regular" text={$_("descSofModal")} color="tertiary" />
    </div>
  </svelte:fragment>
  {#if datasSOF.length > 0}
    <ul class="flex flex-col gap-2 list-disc pl-5">
      {#each datasSOF as data}
        <li>
          <Typography preset="Text sm/Regular" text={data.summary} color="primary" />
        </li>
      {/each}
    </ul>
  {/if}
</Modal>
