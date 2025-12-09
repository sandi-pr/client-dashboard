<script>
  import { onMount } from "svelte";
  import { Button, Modal } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { call, triggerToast, triggerLoading, getIcons } from "$lib";
  import { Typography, PdfjsViewer, Divider, DataTableAIReview, ReportCard } from "$lib";
  import { Document, Packer, Paragraph, TextRun } from "docx";

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

  let allDatas = $state([]);
  let datasSOF = $state([]);
  let tmpResult = $state([]);
  let actionItems = $state([]);

  let CSVResultSummary = $state({
    valid: 0,
    invalid: 0,
    total: 0,
    // need_review: 0,
    // resolved: 0,
    // unresolved: 0,
  });
  let rawResult = $state({});

  // const normalizeResult = (raw) => {
  //   const result = [];
  //   const summary = [];
  //   let total = 0;
  //   let invalid = 0;

  //   for (const [cif, forms] of Object.entries(raw)) {
  //     for (const [form, content] of Object.entries(forms)) {
  //       if (!form.startsWith("F") || !content?.row_data || content.row_data.length === 0) continue;

  //       total++;

  //       const firstRow = content.row_data[0];

  //       const logItem = {
  //         no: 1,
  //         row_data: firstRow,
  //         reference_document: {
  //           id: null,
  //           name: "No Reference",
  //           prompt: "",
  //           result: "",
  //         },
  //         ai_extraction: "",
  //         status_icon: "⚠️",
  //         status_label: "Needs Review",
  //         ai_suggestion:
  //           forms?.finding?.completeness?.[form] ||
  //           forms?.finding?.comparison?.[form]?.reason ||
  //           "Review required",
  //       };

  //       result.push({
  //         title: `CIF ${cif} - ${form}`,
  //         generated: {
  //           results_query: [],
  //           tabular_result: { logs: [logItem] },
  //         },
  //       });

  //       invalid++;

  //       if (forms.finding) {
  //         if (forms.finding.comparison?.[form]) {
  //           summary.push({
  //             row: `CIF ${cif} - ${form}`,
  //             context: "Comparison",
  //             summary: forms.finding.comparison[form].reason,
  //           });
  //         }

  //         if (forms.finding.accuracy?.[form]) {
  //           const firstKey = Object.keys(forms.finding.accuracy[form])[0];
  //           summary.push({
  //             row: `CIF ${cif} - ${form}`,
  //             context: "Accuracy",
  //             summary: forms.finding.accuracy[form][firstKey],
  //           });
  //         }

  //         if (forms.finding.completeness?.[form]) {
  //           summary.push({
  //             row: `CIF ${cif} - ${form}`,
  //             context: "Completeness",
  //             summary: forms.finding.completeness[form],
  //           });
  //         }
  //       }
  //     }
  //   }

  //   return {
  //     result,
  //     summary_of_findings: summary.map((s, i) => ({ id: i + 1, ...s })),
  //     csv_summary: { valid: total - invalid, invalid, total },
  //   };
  // };

  const getProjectResult = async () => {
    try {
      triggerLoading(true);

      let token = localStorage.getItem("token");
      if (!token) throw new Error("No token found");

      const response = await $call.call("client-dashboard", "get-contract-audit-result", {
        URLQuery: {},
        URLParams: { id: projectId },
        Params: {},
        Headers: { Authorization: "Bearer " + token },
      });

      // console.log("response", response);
      // console.log(
      //   "here",
      //   Object.entries(response.result.result["009597"])
      //     .map((form) => {
      //       if (form[1]?.row_data && form[1].row_data.length === 0) {
      //         return form[0];
      //       } else {
      //         return null;
      //       }
      //     })
      //     .filter((f) => f !== null)
      //     .join(", "),
      // );

      if (response.status !== 200) throw new Error("Failed to fetch");

      projectName = response.result.project_name || "Audit Result";
      generatedOn = response.result.generated_on || new Date().toISOString();
      CSVResultSummary = response.result.csv_summary;
      rawResult = response.result.result;
      datasSOF = response.result.summary_of_findings;
      allDatas = datasSOF;

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

  // Function to generate and download the DOCX file

  const generateDocx = async () => {
    const sofBlocks = [];

    // make datasSOF loopable instead of hardcoding index 0
    if (Array.isArray(datasSOF) && datasSOF.length > 0) {
      datasSOF.forEach((sof, idx) => {
        // Seperator
        sofBlocks.push(
          new Paragraph({
            children: [
              new TextRun({
                text: "================================",
                break: 1,
              }),
            ],
          }),
        );

        // CIF header
        sofBlocks.push(
          new Paragraph({
            children: [
              new TextRun({
                text: `💠 CIF: ${sof?.cif ?? ""}`,
                bold: true,
                break: 1,
              }),
            ],
          }),
        );

        // unreported form
        sofBlocks.push(
          new Paragraph({
            children: [
              new TextRun({
                text: "📍 Form belum dilaporkan: ",
                bold: true,
                break: 1,
              }),
              new TextRun({
                text: Object.entries(rawResult[sof?.cif ?? ""])
                  .map((form) => {
                    if (form[1]?.row_data && form[1].row_data.length === 0) {
                      return form[0];
                    } else {
                      return null;
                    }
                  })
                  .filter((f) => f !== null)
                  .join(", "),
              }),
            ],
          }),
        );

        // section title (keeps previous wording)
        sofBlocks.push(
          new Paragraph({
            children: [
              new TextRun({
                text: "📍 Temuan Kolektibilitas:",
                bold: true,
                break: 1,
              }),
            ],
          }),
        );

        // parsed markdown summary
        sofBlocks.push(...markdownToDocxBlocks(sof?.summary ?? ""));

        // add a blank paragraph between entries for readability
        if (idx < datasSOF.length - 1) {
          sofBlocks.push(
            new Paragraph({
              children: [new TextRun({ text: "" })],
            }),
          );
        }
      });
    } else {
      // fallback when no datasSOF
      sofBlocks.push(
        new Paragraph({
          children: [new TextRun({ text: datasSOF?.summary ?? "" })],
        }),
      );
    }

    const doc = new Document({
      sections: [
        {
          properties: {},
          children: [
            new Paragraph({
              children: [
                new TextRun({
                  text: "LAPORAN TEMUAN AUDIT DATA SLIK",
                }),
              ],
            }),
            new Paragraph({
              children: [
                new TextRun({
                  text: "Project: Audit Result",
                }),
              ],
            }),
            new Paragraph({
              children: [
                new TextRun({
                  text: "Generated: " + new Date(generatedOn).toLocaleString(),
                }),
              ],
            }),
            ...sofBlocks,
          ],
        },
      ],
    });

    const mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document";

    const blob = await Packer.toBlob(doc);

    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `${projectName}.docx`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const createInlineRuns = (text) => {
    const runs = [];
    const inlineRe = /(\*\*(.+?)\*\*|\*(.+?)\*|`(.+?)`)/g;
    let lastIndex = 0;
    let m;
    while ((m = inlineRe.exec(text)) !== null) {
      if (m.index > lastIndex) {
        runs.push(new TextRun({ text: text.slice(lastIndex, m.index) }));
      }
      if (m[2]) {
        runs.push(new TextRun({ text: m[2], bold: true }));
      } else if (m[3]) {
        runs.push(new TextRun({ text: m[3], italics: true }));
      } else if (m[4]) {
        runs.push(new TextRun({ text: m[4], font: "Courier New" }));
      }
      lastIndex = m.index + m[0].length;
    }
    if (lastIndex < text.length) {
      runs.push(new TextRun({ text: text.slice(lastIndex) }));
    }
    return runs;
  };

  const markdownToDocxBlocks = (md) => {
    if (!md) return [new Paragraph({ children: [new TextRun({ text: "" })] })];

    const lines = md.replace(/\r\n/g, "\n").split("\n");
    const blocks = [];
    let inCode = false;
    let codeBuf = [];

    const pushCodeBlock = () => {
      if (codeBuf.length === 0) return;
      blocks.push(
        new Paragraph({
          children: [new TextRun({ text: codeBuf.join("\n"), font: "Courier New" })],
        }),
      );
      codeBuf = [];
    };

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];

      if (line.trim().startsWith("```")) {
        inCode = !inCode;
        if (!inCode) pushCodeBlock();
        continue;
      }

      if (inCode) {
        codeBuf.push(line);
        continue;
      }

      const h = line.match(/^(#{1,6})\s+(.*)$/);
      if (h) {
        const level = h[1].length;
        const text = h[2];
        const sizeMap = { 1: 28, 2: 24, 3: 20, 4: 18, 5: 16, 6: 14 };
        blocks.push(
          new Paragraph({
            children: [new TextRun({ text, bold: true, size: sizeMap[level] || 20 })],
          }),
        );
        continue;
      }

      const li = line.match(/^\s*[-*]\s+(.*)$/);
      if (li) {
        const children = createInlineRuns("• " + li[1]);
        blocks.push(new Paragraph({ children }));
        continue;
      }

      if (line.trim() === "") {
        blocks.push(new Paragraph({ children: [new TextRun({ text: "" })] }));
        continue;
      }

      // normal paragraph with inline formatting
      blocks.push(new Paragraph({ children: createInlineRuns(line) }));
    }

    // flush any leftover code (if no closing ``` was present)
    if (inCode && codeBuf.length) pushCodeBlock();

    return blocks;
  };

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
          await generateDocx();
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
