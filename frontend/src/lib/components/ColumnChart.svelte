<script lang="ts">
  import type { ApexOptions } from "apexcharts";
  import { Chart } from "@flowbite-svelte-plugins/chart";
  import { Indicator } from "flowbite-svelte";
  import { UsersGroupOutline, ArrowUpOutline } from "flowbite-svelte-icons";
  import { _ } from "svelte-i18n";
  import { Typography } from "$lib";
  import { onMount } from "svelte";

  export let chartSeries = [];
  export let categories = [];

  const options: ApexOptions = {
    colors: ["#6941C6", "#9E77ED", "#EAECF0"],
    series: chartSeries,
    chart: {
      type: "bar",
      height: "320px",
      fontFamily: "Inter, sans-serif",
      toolbar: {
        show: false,
      },
    },
    plotOptions: {
      bar: {
        horizontal: false,
        columnWidth: "70%",
        borderRadiusApplication: "end",
        borderRadius: 8,
      },
    },
    tooltip: {
      shared: true,
      intersect: false,
      style: {
        fontFamily: "Inter, sans-serif",
      },
    },
    states: {
      hover: {
        filter: {
          type: "darken",
        },
      },
    },
    stroke: {
      show: true,
      width: 0,
      colors: ["transparent"],
    },
    grid: {
      show: false,
      strokeDashArray: 4,
      padding: {
        left: 2,
        right: 2,
        top: -14,
      },
    },
    dataLabels: {
      enabled: false,
    },
    legend: {
      show: false,
    },
    xaxis: {
      categories: categories,
      floating: false,
      labels: {
        show: true,
        style: {
          fontFamily: "Inter, sans-serif",
          cssClass: "text-xs font-normal fill-gray-500 dark:fill-gray-400",
        },
      },
      axisBorder: {
        show: false,
      },
      axisTicks: {
        show: false,
      },
    },
    yaxis: {
      show: true,
      title: {
        text: "Debitur",
        style: {
          fontSize: "12px",
          fontWeight: 500,
          color: "#475467",
        },
      },
    },
    fill: {
      opacity: 1,
    },
  };

  // onMount(() => {
  //   console.log("check chartSeries", chartSeries);
  // })
</script>

<div class="">
  <Typography preset="Text lg/Semibold" text={$_("auditPerformanceTrends")} />
  <div class="flex justify-end gap-3">
    <div class="flex items-center">
      <Indicator
        size="sm"
        color="none"
        class="me-1.5 bg-utility-brand-700-light dark:bg-utility-brand-700-dark"
      />
      <Typography preset="Text sm/Regular" text={$_("compliantData")} color="tertiary" />
    </div>
    <div class="flex items-center">
      <Indicator
        size="sm"
        color="none"
        class="me-1.5 bg-utility-brand-500-light dark:bg-utility-brand-500-dark"
      />
      <Typography preset="Text sm/Regular" text={$_("nonCompliantData")} color="tertiary" />
    </div>
    <div class="flex items-center">
      <Indicator
        size="sm"
        color="none"
        class="me-1.5 bg-utility-gray-200-light dark:bg-utility-gray-200-dark"
      />
      <Typography preset="Text sm/Regular" text={$_("totalData")} color="tertiary" />
    </div>
  </div>
  <Chart {options} />
</div>
