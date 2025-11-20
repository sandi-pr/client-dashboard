<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer:
  - Dean Ramhan <dean@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->
<script>
  import { onMount } from "svelte";
  import { Heading } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { call, check } from "$lib";
  import {
    Typography,
    DataTable,
    getIcons,
    ReportCard,
    ColumnChart,
    getProjects,
    StateCard,
  } from "$lib";

  const IllustrationEmpty = getIcons("illustrationEmpty");
  const IllustrationProgress = getIcons("illustrationProgress");

  const isFrontend = process.env.NODE_ENV === "development";
  let gettingDatas = $state(false);
  let gettingDataChart = $state(true);
  let actionItems = $state([]);
  let datasTable = $state([]);
  let datasConfigs = $state({});
  let chartCategories = $state([]);

  let chartSeries = $state([
    { name: "Compliant Data", color: "#6941C6", data: [] },
    { name: "Non-compliant Data", color: "#9E77ED", data: [] },
    { name: "Total Data", color: "#EAECF0", data: [] },
  ]);

  const getActionItems = async (summary) => {
    actionItems = [
      {
        href: isFrontend
          ? "/document_reviewer/report/summary" + "?type=compliant&projectDocumentId=" + 82
          : (await $call.getPage("client-dashboard", "report-summary")) +
            "?type=compliant&projectDocumentId=" +
            82,
        icon: "checkCircle",
        iconColor: "fg-success-primary",

        // gunakan data API
        actionName: summary.valid.toLocaleString(),
        desc: "compliantData",

        // persentase juga dari API
        subtext: summary.percentage_valid + "%",

        changeType: "increase",
      },
      {
        href: isFrontend
          ? "/document_reviewer/report/summary" + "?type=nonCompliant&projectDocumentId=" + 82
          : (await $call.getPage("client-dashboard", "report-summary")) +
            "?type=nonCompliant&projectDocumentId=" +
            82,
        icon: "alertTriangle",
        iconColor: "fg-warning-primary",

        actionName: summary.invalid.toLocaleString(),
        desc: "nonCompliantData",
        subtext: summary.percentage_invalid + "%",

        changeType: "decrease",
      },
      {
        href: isFrontend
          ? "/document_reviewer/report/summary" + "?type=total&projectDocumentId=" + 82
          : (await $call.getPage("client-dashboard", "report-summary")) +
            "?type=total&projectDocumentId=" +
            82,
        icon: "infoCircle",
        iconColor: "fg-secondary",

        actionName: summary.total_people.toLocaleString(),
        desc: "totalData",
        subtext: summary.percentage_total_people + "%",

        changeType: "increase",
      },
    ];
  };

  const handleGetProjects = async () => {
    gettingDatas = true;
    let token = localStorage.getItem("token");

    const response = await $call.call("client-dashboard", "get-project-report", {
      URLQuery: {},
      URLParams: {},
      Params: {
        page: 1,
        page_size: 25,
      },
      Headers: {
        Authorization: "Bearer " + token,
      },
    });

    // console.log("check get-project-report", response);

    if (response === null) {
      gettingDatas = false;
      return;
    }

    if (response.status !== 200) {
      gettingDatas = false;
      console.log(response);
      return;
    }

    let projects = response.result.data;

    if (projects.length > 0) {
      datasConfigs = {
        headers: "projectReportLists",
        actions: [
          {
            name: "viewAll",
            href: isFrontend
              ? "/document_reviewer/report/project"
              : await $call.getPage("client-dashboard", "contract-report-project-list"),
            type: "info",
          },
        ],
        seperator: true,
      };

      let url = isFrontend
        ? "/document_reviewer/report/project/detail"
        : await $call.getPage("client-dashboard", "contract-report-project-detail");

      datasTable = projects.map((project) => {
        const { document_rule_id, ...rest } = {
          id: project.id,
          project_name: project.attributes.project_name,
          actions: {
            objects: [
              {
                href: url + `?projectDocumentId=82&projectId=${project.id}`,
                name: "View",
                type: "primary",
              },
            ],
          },
        };
        return {
          ...rest,
        };
      });

      if (datasTable.length > 0) {
        for (let i = 0; i < datasTable.length; i++) {
          const responseDocument = await $call.call(
            "client-dashboard",
            "get-document-contract-project",
            {
              URLQuery: {},
              URLParams: {},
              Params: {
                project_id: datasTable[i].id.toString(),
              },
              Headers: {
                Authorization: "Bearer " + token,
              },
            },
          );

          // console.log("responseDocument", responseDocument);

          if (responseDocument.status !== 200) {
            let error = new Error("Failed to fetch document for project id " + datasTable);
            error.response = responseDocument;
            throw error;
          }

          const responseFiles = await getStatusGenerated(responseDocument.result.files);

          const projectDocumentIdStrings = responseFiles.result.data.map((data) =>
            data.id.toString(),
          );

          let url = isFrontend
            ? "/document_reviewer/report/project/detail"
            : await $call.getPage("client-dashboard", "contract-report-project-detail");

          datasTable[i].actions = {
            objects: [
              check("aiManagement.contractreviewer.project:view") && {
                href:
                  url +
                  `?projectDocumentId=${projectDocumentIdStrings.join(",")}&projectId=${datasTable[i].id}`,
                name: "View",
                type: "primary",
              },
            ].filter(Boolean),
          };
        }
      }
    } else {
      datasTable = [];
    }

    gettingDatas = false;
  };

  const getStatusGenerated = async (files) => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        for (let i = 0; i < files.length; i++) {
          if (
            files[i].status_generated !== "START" &&
            files[i].status_generated !== "FINISH" &&
            files[i].status_generated !== "GENERATING"
          )
            continue;
          else {
            const response = await $call.call(
              "client-dashboard",
              "get-contract-project-result-with-document-id",
              {
                URLQuery: {},
                URLParams: {},
                Params: {
                  project_document_id: files[i].id.toString(),
                },
                Headers: {
                  Authorization: "Bearer " + token,
                },
              },
            );

            // console.log("response", response);

            if (response.status !== 200) {
              triggerLoading(false);
              triggerToast("error");
              console.log(response);
              break;
            }

            return response;
          }
        }
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        console.log(error);
      }
    }
  };

  function rotateMonths(monthNames) {
    const now = new Date();
    const currentMonth = now.getMonth();
    const before = monthNames.slice(currentMonth + 1);
    const after = monthNames.slice(0, currentMonth + 1);

    return [...before, ...after];
  }

  function formatChartData(chartData) {
    const monthNames = [
      "Jan",
      "Feb",
      "Mar",
      "Apr",
      "May",
      "Jun",
      "Jul",
      "Aug",
      "Sep",
      "Oct",
      "Nov",
      "Dec",
    ];

    const orderedMonths = rotateMonths(monthNames);

    const compliantData = [];
    const nonCompliantData = [];
    const totalData = [];

    const dataMap = new Map();
    chartData.forEach((item) => {
      const monthIndex = parseInt(item.month.split("-")[1]) - 1;
      const label = monthNames[monthIndex];
      dataMap.set(label, {
        valid: item.valid ?? 0,
        invalid: item.invalid ?? 0,
        total_people: item.total_people ?? 0,
      });
    });

    orderedMonths.forEach((label) => {
      const d = dataMap.get(label) || { valid: 0, invalid: 0, total_people: 0 };

      compliantData.push({ x: label, y: d.valid });
      nonCompliantData.push({ x: label, y: d.invalid });
      totalData.push({ x: label, y: d.total_people });
    });

    return { compliantData, nonCompliantData, totalData, orderedMonths };
  }

  const getReport = async () => {
    let token = localStorage.getItem("token");
    if (!token) return;

    try {
      const response = await $call.call("client-dashboard", "get-all-report-document-reviewer", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      // console.log("check get-all-report-document-reviewer", response);

      if (response.result && response.result.chart && response.result.chart.length > 0) {
        const { compliantData, nonCompliantData, totalData, orderedMonths } = formatChartData(
          response.result.chart,
        );

        chartSeries = [
          { ...chartSeries[0], data: compliantData },
          { ...chartSeries[1], data: nonCompliantData },
          { ...chartSeries[2], data: totalData },
        ];

        chartCategories = orderedMonths;

        // console.log("check orderedMonths", orderedMonths);
        // console.log("check chartSeries", chartSeries);
      }

      gettingDataChart = false;
      return response.result.summary;
    } catch (error) {
      console.log(error);
    }
  };

  onMount(async () => {
    await handleGetProjects();
    const summary = await getReport();
    await getActionItems(summary);
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-32">
    <div class="flex flex-col m-4">
      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Typography preset="Display sm/Semibold" text={$_("auditReporting")} propClass="fluid-xl" />

        {#if gettingDataChart}
          <ColumnChart />
        {:else}
          <ColumnChart {chartSeries} categories={chartCategories} />
        {/if}

        <!-- Action Items -->
        <div>
          {#if actionItems.some((item) => item.disabled !== true)}
            <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
              <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
                {$_("recentProjectSummary")}
              </Heading>
            </div>
          {/if}
          <ReportCard {actionItems} />
        </div>

        {#if gettingDatas}
          <StateCard
            illustration={IllustrationProgress}
            description={$_("fetchingDataPleaseWait")}
          />
        {/if}

        {#if datasTable.length > 0}
          <DataTable bind:datas={datasTable} bind:configs={datasConfigs} />
        {/if}

        {#if !gettingDatas && datasTable.length === 0}
          <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
        {/if}
      </div>
    </div>
  </div>
</div>
