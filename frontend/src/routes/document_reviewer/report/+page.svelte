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
  let actionItems = $state([]);
  let datasTable = $state([]);
  let datasConfigs = $state({});

  let chartSeries = [
    { name: "Compliant Data", color: "#6941C6", data: [] },
    { name: "Non-compliant Data", color: "#9E77ED", data: [] },
    { name: "Total Data", color: "#EAECF0", data: [] },
  ];

  const getActionItems = async () => {
    actionItems = [
      {
        href: isFrontend
          ? "/document_reviewer/report/summary" + "?type=compliant&projectDocumentId=" + 82
          : (await $call.getPage("client-dashboard", "report-summary")) +
            "?type=compliant&projectDocumentId=" +
            82,
        icon: "checkCircle",
        iconColor: "fg-success-primary",
        actionName: "1,273",
        desc: "compliantData",
        subtext: "20%",
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
        actionName: "6,256",
        desc: "nonCompliantData",
        subtext: "50%",
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
        actionName: "11,890",
        desc: "totalData",
        subtext: "40%",
        changeType: "increase",
      },
    ];
  };

  const handleGetProjects = async () => {
    gettingDatas = true;

    const projects = await getProjects($call, { page: 1, pageSize: 5 }, false);
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
        ? "/document_reviewer/report/detail"
        : await $call.getPage("client-dashboard", "contract-project-detail");

      const promise = projects.map((project) => {
        const { document_rule_id, ...rest } = {
          id: project.id,
          project_name: project.attributes.project_name,
          actions: {
            objects: [
              {
                href: url + `?projectId=${project.id}`,
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

      datasTable = await Promise.all(promise);
    } else {
      datasTable = [];
    }

    gettingDatas = false;
  };

  function formatChartData(chartData) {
    const compliantData = [];
    const nonCompliantData = [];
    const totalData = [];
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

    const dataMap = new Map();
    chartData.forEach((item) => {
      const monthIndex = parseInt(item.month.split("-")[1]) - 1;
      const monthLabel = monthNames[monthIndex];

      dataMap.set(monthLabel, {
        valid: item.valid,
        invalid: item.invalid,
        total_people: item.total_people,
      });
    });

    monthNames.forEach((monthLabel) => {
      const data = dataMap.get(monthLabel) || { valid: 0, invalid: 0, total_people: 0 };

      compliantData.push({ x: monthLabel, y: data.valid });
      nonCompliantData.push({ x: monthLabel, y: data.invalid });
      totalData.push({ x: monthLabel, y: data.total_people });
    });

    return { compliantData, nonCompliantData, totalData };
  }

  const getReport = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        // ... (panggilan API Anda) ...
        const response = await $call.call("client-dashboard", "get-all-report-document-reviewer", {
          // ... Params dan Headers Anda ...
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        console.log("check get-all-report-document-reviewer", response);

        // --- PROSES DATA UNTUK CHART ---
        if (response.result && response.result.chart && response.result.chart.length > 0) {
          const { compliantData, nonCompliantData, totalData } = formatChartData(
            response.result.chart,
          );

          // 4. Perbarui variabel reaktif chartSeries
          // Mengganti seluruh array untuk memicu reaktivitas
          chartSeries = [
            { ...chartSeries[0], data: compliantData },
            { ...chartSeries[1], data: nonCompliantData },
            { ...chartSeries[2], data: totalData },
          ];
        }
        // --- AKHIR PROSES DATA ---
      } catch (error) {
        console.log(error);
      }
    }
  };

  onMount(async () => {
    await getReport();
    await handleGetProjects();
    await getActionItems();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-32">
    <div class="flex flex-col m-4">
      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Typography preset="Display sm/Semibold" text={$_("auditReporting")} propClass="fluid-xl" />

        <ColumnChart />

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
