<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: 
  - Azzis Arswendo <azzis@hexacode.org>
  - Dean Ramhan <dean@hexacode.org>
  - Taqie Fadlillah <taqief@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->

<script>
  import { CR3Call } from "hexacomponents-svelte";
  import {
    Card,
    Typography,
    SectionHeader,
    triggerLoading,
    ChartCard,
    ActionCard,
    industryData,
  } from "$lib";
  import { StateCard } from "$lib";
  import IllustrationFirst from "$lib/icon/illustration_first_time.svelte";
  import { Alert, CloseButton, P } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { goto } from "$app/navigation";
  import { onMount, tick } from "svelte";
  import { writable } from "svelte/store";
  import LoadingContainer from "../lib/components/LoadingContainer.svelte";
  import { call, showAI, showNotif } from "$lib";

  let userData = $state({ username: "", email: "", brand: "" });

  let isFirstTime = $state(true);

  let chartActionBusiness = $state([]);
  let chartActionAi = $state([]);
  let chartReady = $state(false);

  let monthlyEmployeeCounts = $state([]);
  let monthlyPersonsCounts = $state([]);
  let monthlyProductCounts = $state([]);
  let monthlyOrderCounts = $state([]);
  let monthlyScheduleCounts = $state([]);

  let monthlyBotCounts = $state([]);
  let monthlyCollectionsCounts = $state([]);
  let monthlyDatasetCounts = $state([]);

  let dataChartAllBusiness = $state([]);
  let dataChartAllAi = $state([]);
  let selectedIndustryType = $state(null);

  const cardsBusiness = [
    {
      icon: "user",
      title: "person",
      description: "descPerson",
      page: "crm-person",
    },
    {
      icon: "faceHappy",
      title: "HR",
      description: "descHr",
      page: "hr-overview",
    },
    {
      icon: "package",
      title: "product&Services",
      description: "descProduct&Services",
      page: "crm-product-service",
    },
    {
      icon: "receipt",
      title: "order&Invoice",
      description: "descOrder&Invoice",
      page: "crm-order-invoice",
    },
    {
      icon: "calendar",
      title: "calendar",
      description: "descCalendar",
      page: "crm-calendar",
    },
  ];

  const cardsAi = [
    {
      icon: "messageChat",
      title: "hexabot",
      description: "descHexabot",
      page: "hexabot",
    },
    {
      icon: "archive",
      title: "documentCollection",
      description: "descDocumentCollection",
      page: "document-collection-overview",
    },
    {
      icon: "fileCheck",
      title: "contractReviewer",
      description: "descContractReviewer",
      page: "contract-reviewer-overview",
    },
    {
      icon: "data",
      title: "datasets",
      description: "descDatasets",
      page: "datasets",
    },
  ];

  let docAiCard = $state([]);

  async function getIndustryPoints(industryId) {
    // Find the industry object by id
    if (typeof industryId === "string") {
      industryId = parseInt(industryId);
    }
    const industry = industryData.find((item) => item.id === industryId);
    if (!industry) return [];

    let industryLibrary = industry.points.map((point) => ({
      actionName: point.title,
      desc: point.desc,
    }));

    // return industry.points.map((point) => ({
    //   href: "",
    //   icon: "bookClosed",
    //   actionName: point.title,
    //   desc: `desc ${point.title}`,
    // }));

    let initIndustryLibrary = industry.points.map((point) => ({
      href: "",
      icon: "bookClosed",
      actionName: point.title,
      desc: `desc ${point.title}`,
    }));

    for (const point of industryLibrary) {
      let token = localStorage.getItem("token");
      if (token) {
        const response = await $call.call("client-dashboard", "find-create-library-by-name", {
          URLQuery: {},
          URLParams: {},
          Params: {
            library_name: point.actionName,
            description: point.desc,
          },
          Headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        let link =
          (await $call.getPage("client-dashboard", "document-library-collections")) +
          "?libraryId=" +
          response.result.id;

        initIndustryLibrary.find((item) => item.actionName === point.actionName).href = link;
      }
    }

    return initIndustryLibrary;
  }

  async function getActionItems() {
    docAiCard = [
      {
        href: await $call.getPage("client-dashboard", "document-collection-overview"),
        icon: "bookOpen",
        actionName: "Document library",
        desc: "desc Document library",
      },
      // ...(await getIndustryPoints(sessionStorage.getItem("industry_id"))),
    ];
  }

  async function getCompanyType() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "get-company-type", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      selectedIndustryType = parseInt(response.result.data.attributes.company_industry);
    }
  }

  async function getChartAction(total, module) {
    let type = "business";

    let tmp = cardsBusiness.find((card) => card.title === module);
    if (tmp === undefined) {
      tmp = cardsAi.find((card) => card.title === module);
      type = "ai";
    }

    if (tmp === undefined) {
      console.log(module);
      return;
    }

    let tmpName = "";
    if (module === "person") {
      tmpName = "Total Person";
    } else if (module === "HR") {
      tmpName = "Total Employee";
    } else if (module === "product&Services") {
      tmpName = "Total Product";
    } else if (module === "order&Invoice") {
      tmpName = "Total Order";
    } else if (module === "calendar") {
      tmpName = "Total Schedule";
    } else if (module === "hexabot") {
      tmpName = "Total Bot";
    } else if (module === "documentCollection") {
      tmpName = "Total Collection";
    } else if (module === "datasets") {
      tmpName = "Total Dataset";
    }

    if (type === "business") {
      chartActionBusiness = [
        ...chartActionBusiness,
        {
          icon: tmp.icon,
          title: tmp.title,
          description: tmp.description,
          href: await $call.getPage("client-dashboard", tmp.page),
          dataName: tmpName,
          value: total.total,
          change: total.change,
          changeType: total.changeType,
        },
      ];
    } else if (type === "ai") {
      chartActionAi = [
        ...chartActionAi,
        {
          icon: tmp.icon,
          title: tmp.title,
          description: tmp.description,
          href: await $call.getPage("client-dashboard", tmp.page),
          dataName: tmpName,
          value: total.total,
          change: total.change,
          changeType: total.changeType,
        },
      ];
    }

    generateChartData(module, type);
  }

  function generateChartData(module, type) {
    const currentDate = new Date();
    const twelveMonthsBefore = new Date();
    twelveMonthsBefore.setMonth(currentDate.getMonth() - 12);

    const formatDate = (date) => {
      const options = { month: "short", year: "numeric" };
      return new Intl.DateTimeFormat("en-US", options).format(date);
    };

    let tmpName = "";
    let tmpData = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];
    if (module === "person") {
      tmpName = "Total Person";
      tmpData = monthlyPersonsCounts;
    } else if (module === "HR") {
      tmpName = "Total Employee";
      tmpData = monthlyEmployeeCounts;
    } else if (module === "product&Services") {
      tmpName = "Total Product";
      tmpData = monthlyProductCounts;
    } else if (module === "order&Invoice") {
      tmpName = "Total Order";
      tmpData = monthlyOrderCounts;
    } else if (module === "calendar") {
      tmpName = "Total Schedule";
      tmpData = monthlyScheduleCounts;
    } else if (module === "hexabot") {
      tmpName = "Total Bot";
      tmpData = monthlyBotCounts;
    } else if (module === "documentCollection") {
      tmpName = "Total Collection";
      tmpData = monthlyCollectionsCounts;
    } else if (module === "datasets") {
      tmpName = "Total Dataset";
      tmpData = monthlyDatasetCounts;
    }

    if (type === "business") {
      dataChartAllBusiness = [
        ...dataChartAllBusiness,
        {
          chart: {
            type: "area",
            height: "150px",
            fontFamily: "Inter, sans-serif",
            dropShadow: {
              enabled: false,
            },
            toolbar: {
              show: false,
            },
          },
          tooltip: {
            enabled: true,
            x: {
              show: false,
            },
          },
          fill: {
            type: "gradient",
            gradient: {
              opacityFrom: 0.55,
              opacityTo: 0,
              shade: "#17B169",
              gradientToColors: ["#17B169"],
            },
          },
          dataLabels: {
            enabled: false,
          },
          stroke: {
            width: 2,
          },
          grid: {
            show: false,
            strokeDashArray: 4,
            padding: {
              left: 2,
              right: 2,
              top: 0,
            },
          },
          series: [
            {
              name: tmpName,
              data: tmpData,
              color: "#17B169",
            },
          ],
          xaxis: {
            categories: Array.from({ length: 12 }, (_, i) => {
              const date = new Date(twelveMonthsBefore);
              date.setMonth(twelveMonthsBefore.getMonth() + i + 1);
              return formatDate(date);
            }),
            labels: {
              show: false,
            },
            axisBorder: {
              show: false,
            },
            axisTicks: {
              show: false,
            },
          },
          yaxis: {
            show: false,
          },
        },
      ];
    } else if (type === "ai") {
      dataChartAllAi = [
        ...dataChartAllAi,
        {
          chart: {
            type: "area",
            height: "150px",
            fontFamily: "Inter, sans-serif",
            dropShadow: {
              enabled: false,
            },
            toolbar: {
              show: false,
            },
          },
          tooltip: {
            enabled: true,
            x: {
              show: false,
            },
          },
          fill: {
            type: "gradient",
            gradient: {
              opacityFrom: 0.55,
              opacityTo: 0,
              shade: "#17B169",
              gradientToColors: ["#17B169"],
            },
          },
          dataLabels: {
            enabled: false,
          },
          stroke: {
            width: 2,
          },
          grid: {
            show: false,
            strokeDashArray: 4,
            padding: {
              left: 2,
              right: 2,
              top: 0,
            },
          },
          series: [
            {
              name: tmpName,
              data: tmpData,
              color: "#17B169",
            },
          ],
          xaxis: {
            categories: Array.from({ length: 12 }, (_, i) => {
              const date = new Date(twelveMonthsBefore);
              date.setMonth(twelveMonthsBefore.getMonth() + i + 1);
              return formatDate(date);
            }),
            labels: {
              show: false,
            },
            axisBorder: {
              show: false,
            },
            axisTicks: {
              show: false,
            },
          },
          yaxis: {
            show: false,
          },
        },
      ];
    }
  }

  async function getCountPerModulePerMonth() {
    const today = new Date();
    let earlyDateCurrentMonth = new Date(today.getFullYear(), today.getMonth(), 1);
    let lastDateCurrentMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0);
    let token = localStorage.getItem("token");

    let tmpMonthlyPersonsCounts = [];
    let tmpMonthlyEmployeeCounts = [];
    let tmpMonthlyProductCounts = [];
    let tmpMonthlyOrderCounts = [];
    let tmpMonthlyScheduleCounts = [];

    let tmpMonthlyBotCounts = [];
    let tmpMonthlyCollectionsCounts = [];
    let tmpMonthlyDatasetCounts = [];

    if (token !== null) {
      for (let i = 0; i < 12; i++) {
        let response = await $call.call("client-dashboard", "get-count-persons-per-month", {
          URLQuery: {},
          URLParams: {},
          Params: {
            earlyDate: earlyDateCurrentMonth.toISOString().split("T")[0],
            lastDate: lastDateCurrentMonth.toISOString().split("T")[0],
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          console.log(response.status);
          console.log(response);
          return;
        }

        tmpMonthlyPersonsCounts.push(response.result.persons.meta.pagination.total);

        let response2 = await $call.call("client-dashboard", "get-count-employee-per-month", {
          URLQuery: {},
          URLParams: {},
          Params: {
            earlyDate: earlyDateCurrentMonth.toISOString().split("T")[0],
            lastDate: lastDateCurrentMonth.toISOString().split("T")[0],
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response2.status !== 200) {
          console.log(response2.status);
          console.log(response2);
          return;
        }

        tmpMonthlyEmployeeCounts.push(response2.result.employee.meta.pagination.total);

        let response3 = await $call.call("client-dashboard", "get-count-product-per-mount", {
          URLQuery: {},
          URLParams: {},
          Params: {
            earlyDate: earlyDateCurrentMonth.toISOString().split("T")[0],
            lastDate: lastDateCurrentMonth.toISOString().split("T")[0],
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response3.status !== 200) {
          console.log(response3.status);
          console.log(response3);
          return;
        }

        tmpMonthlyProductCounts.push(response3.result.meta.pagination.total);

        let response4 = await $call.call("client-dashboard", "get-count-order-per-mount", {
          URLQuery: {},
          URLParams: {},
          Params: {
            earlyDate: earlyDateCurrentMonth.toISOString().split("T")[0],
            lastDate: lastDateCurrentMonth.toISOString().split("T")[0],
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response4.status !== 200) {
          console.log(response4.status);
          console.log(response4);
          return;
        }

        tmpMonthlyOrderCounts.push(response4.result.meta.pagination.total);

        let response5 = await $call.call("client-dashboard", "get-count-schedule-per-mount", {
          URLQuery: {},
          URLParams: {},
          Params: {
            earlyDate: earlyDateCurrentMonth.toISOString().split("T")[0],
            lastDate: lastDateCurrentMonth.toISOString().split("T")[0],
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response5.status !== 200) {
          console.log(response5.status);
          console.log(response5);
          return;
        }

        tmpMonthlyScheduleCounts.push(response5.result.meta.pagination.total);

        let response6 = await $call.call("client-dashboard", "get-count-bot-per-mount", {
          URLQuery: {},
          URLParams: {},
          Params: {
            earlyDate: earlyDateCurrentMonth.toISOString().split("T")[0],
            lastDate: lastDateCurrentMonth.toISOString().split("T")[0],
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response6.status !== 200) {
          console.log(response6.status);
          console.log(response6);
          return;
        }

        tmpMonthlyBotCounts.push(response6.result.meta.pagination.total);

        let response7 = await $call.call("client-dashboard", "get-count-documents-per-month", {
          URLQuery: {},
          URLParams: {},
          Params: {
            earlyDate: earlyDateCurrentMonth.toISOString().split("T")[0],
            lastDate: lastDateCurrentMonth.toISOString().split("T")[0],
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response7.status !== 200) {
          console.log(response7.status);
          console.log(response7);
          return;
        }

        tmpMonthlyCollectionsCounts.push(response7.result.document.meta.pagination.total);

        let response8 = await $call.call("client-dashboard", "get-count-dataset-per-mount", {
          URLQuery: {},
          URLParams: {},
          Params: {
            earlyDate: earlyDateCurrentMonth.toISOString().split("T")[0],
            lastDate: lastDateCurrentMonth.toISOString().split("T")[0],
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response8.status !== 200) {
          console.log(response8.status);
          console.log(response8);
          return;
        }

        tmpMonthlyDatasetCounts.push(response8.result.meta.pagination.total);

        earlyDateCurrentMonth = new Date(
          earlyDateCurrentMonth.getFullYear(),
          earlyDateCurrentMonth.getMonth() - 1,
          1,
        );
        lastDateCurrentMonth = new Date(
          lastDateCurrentMonth.getFullYear(),
          lastDateCurrentMonth.getMonth(),
          0,
        );
      }

      monthlyPersonsCounts = tmpMonthlyPersonsCounts.reverse();
      monthlyEmployeeCounts = tmpMonthlyEmployeeCounts.reverse();
      monthlyProductCounts = tmpMonthlyProductCounts.reverse();
      monthlyOrderCounts = tmpMonthlyOrderCounts.reverse();
      monthlyScheduleCounts = tmpMonthlyScheduleCounts.reverse();

      monthlyBotCounts = tmpMonthlyBotCounts.reverse();
      monthlyCollectionsCounts = tmpMonthlyCollectionsCounts.reverse();
      monthlyDatasetCounts = tmpMonthlyDatasetCounts.reverse();

      let totalPersons = groupingDataChart(monthlyPersonsCounts);
      let totalEmployees = groupingDataChart(monthlyEmployeeCounts);
      let totalProducts = groupingDataChart(monthlyProductCounts);
      let totalOrders = groupingDataChart(monthlyOrderCounts);
      let totalSchedules = groupingDataChart(monthlyScheduleCounts);

      let totalBots = groupingDataChart(monthlyBotCounts);
      let totalCollections = groupingDataChart(monthlyCollectionsCounts);
      let totalDataset = groupingDataChart(monthlyDatasetCounts);

      await getChartAction(totalPersons, "person");
      await getChartAction(totalEmployees, "HR");
      await getChartAction(totalProducts, "product&Services");
      await getChartAction(totalOrders, "order&Invoice");
      await getChartAction(totalSchedules, "calendar");

      await getChartAction(totalBots, "hexabot");
      await getChartAction(totalCollections, "documentCollection");
      await getChartAction(totalDataset, "datasets");

      chartReady = true;
      // generateChartData();
    } else {
      console.log("no token");
    }
  }

  function groupingDataChart(data) {
    let total = 0;
    for (let i = 0; i < data.length; i++) {
      total += data[i];
    }

    let change = data[data.length - 1] - data[data.length - 2];

    let changeType = change > 0 ? "increase" : change < 0 ? "decrease" : "-";

    return {
      total: total,
      change: change,
      changeType: changeType,
    };
  }

  function setFristTime() {
    isFirstTime = false;
    localStorage.setItem("isFirstTime", isFirstTime);
  }

  function getFristTime() {
    if (typeof localStorage !== "undefined") {
      isFirstTime = localStorage.getItem("isFirstTime") === "false" ? false : true;
      console.log("isFirstTime", isFirstTime);
    } else {
      console.log("localstorage is undefined");
    }
  }

  async function CheckSetup() {
    try {
      // var resp_setup = await $call.call("client-login", "is-need-setup", {
      //   URLQuery: {},
      //   URLParams: {},
      //   Params: {},
      //   Headers: {},
      // });
      // console.log(resp_setup);
      // if (resp_setup.status === 200) {
      //   if (resp_setup.result.is_need_setup === true) {
      //     if (localStorage.getItem("token") !== null) localStorage.removeItem("token");
      //     if (localStorage.getItem("user") !== null) localStorage.removeItem("user");
      //     await $call.goToPage("hexapod-setup", "home-page");
      //     return;
      //   }
      // }
      var token = localStorage.getItem("token");
      if (token !== null) {
        var resp = await $call.call("client-login", "is-token-valid", {
          URLQuery: {},
          URLParams: {},
          Params: {
            token: token,
          },
          Headers: {},
        });

        if (resp.status === 200) {
          if (resp.result.confirmed === true && resp.result.blocked === false) {
            localStorage.setItem("user", JSON.stringify(resp.result));
            userData = resp.result;
          }
        }
      } else {
        await $call.goToPage("client-login", "home-page");
        return;
      }
    } catch (error) {
      // console.log(error);
      userData = { username: "Hexa Staff", email: "hexacode@hexacode.org", brand: "hexacode" };
    }
  }

  const navigateToPage = async (page) => {
    if (page === "bakrie-webchat") {
      goto(await $call.getPage("bakrie", "home-page"));
    } else {
      goto(await $call.getPage("client-dashboard", page));
    }
  };

  onMount(async () => {
    // triggerLoading(true);
    // setTimeout(async () => {
    //   triggerLoading(false);
    // }, 3000);
    // await getCompanyType();
    getActionItems();
    getCountPerModulePerMonth();
    getFristTime();
    CheckSetup();
  });
</script>

<div class="flex flex-col gap-8 px-4 md:py-6">
  {#if isFirstTime}
    <Alert dismissable on:close={setFristTime} class="relative flex justify-center">
      <CloseButton
        slot="close-button"
        let:close
        class="absolute top-2 right-2 ms-1.5 dark:hover:bg-gray-700"
        color="primary"
        on:click={close}
      />
      <StateCard illustration={IllustrationFirst} />
      <div class="flex flex-col justify-center items-center gap-1 pt-2">
        <Typography
          preset="Display sm/Semibold"
          text={$_("welcometoHexapod") + sessionStorage.getItem("sidebar_title") || "Hexapod"}
          color="primary"
        />
        <!-- <Typography
        text="{userData?.data?.username || 'Hexastaff'} - {userData?.data?.email ||
          'staff@hexacode.org'}"
        color="tertiary"
      /> -->
      </div>
    </Alert>
  {:else}
    <div class="flex flex-col gap-1">
      <Typography
        preset="Display sm/Semibold"
        text={`${$_("welcometoHexapod")} ${sessionStorage.getItem("sidebar_title") || "Hexapod"}`}
        color="primary"
      />
      <Typography text="{userData.username} - {userData.email}" color="tertiary" />
    </div>
  {/if}
  <!-- <div class="flex flex-col gap-4">
    <SectionHeader title={$_("documentAi")} />
    <ActionCard actionItems={docAiCard} numGrid="2" />
  </div>
  {#if cardsBusiness.length > 0}
    <div class="flex flex-col gap-4">
      <SectionHeader title={$_("businessManagements")} />
      {#if chartActionBusiness.length > 0}
        <ChartCard isLibrary chartAction={chartActionBusiness} chartList={dataChartAllBusiness} />
      {/if}
      <div
        class={`grid ${$showAI ? "grid-cols-1 xl:grid-cols-2" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-4"} gap-4`}
      >
        {#each cardsBusiness as card}
          {#if !chartActionBusiness.find((item) => item.title === card.title)}
            <Card
              icon={card.icon}
              title={$_(card.title)}
              description={$_(card.description)}
              onclick={() => navigateToPage(card.page)}
            />
          {/if}
        {/each}
      </div>
    </div>
  {/if} -->
  {#if cardsAi.length > 0}
    <div class="flex flex-col gap-4">
      <SectionHeader title={$_("aiManagements")} />
      {#if chartActionAi.length > 0}
        <ChartCard isLibrary chartAction={chartActionAi} chartList={dataChartAllAi} />
      {/if}
      <div
        class={`grid ${$showAI || $showNotif ? "grid-cols-1 xl:grid-cols-2" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-4"} gap-4`}
      >
        {#each cardsAi as card}
          {#if !chartActionAi.find((item) => item.title === card.title)}
            <Card
              icon={card.icon}
              title={$_(card.title)}
              description={$_(card.description)}
              onclick={() => navigateToPage(card.page)}
            />
          {/if}
        {/each}
      </div>
    </div>
  {/if}
</div>

<!-- <LoadingContainer /> -->
