<script lang="ts">
  import type { ApexOptions } from "apexcharts";
  import { Chart } from "@flowbite-svelte-plugins/chart";
  import { Indicator } from "flowbite-svelte";
  import { UsersGroupOutline, ArrowUpOutline } from "flowbite-svelte-icons";
  import { _ } from "svelte-i18n";
  import { Typography } from "$lib";

  const options: ApexOptions = {
    colors: ["#6941C6", "#9E77ED", "#EAECF0"],
    series: [
      {
        name: "Compliant Data",
        color: "#6941C6",
        data: [
          { x: "Jan", y: 300 },
          { x: "Feb", y: 370 },
          { x: "Mar", y: 230 },
          { x: "Apr", y: 320 },
          { x: "May", y: 280 },
          { x: "Jun", y: 350 },
          { x: "Jul", y: 320 },
          { x: "Aug", y: 300 },
          { x: "Sep", y: 290 },
          { x: "Oct", y: 330 },
          { x: "Nov", y: 370 },
          { x: "Dec", y: 280 }
        ]
      },
      {
        name: "Non-compliant Data",
        color: "#9E77ED",
        data: [
          { x: "Jan", y: 250 },
          { x: "Feb", y: 350 },
          { x: "Mar", y: 170 },
          { x: "Apr", y: 270 },
          { x: "May", y: 120 },
          { x: "Jun", y: 310 },
          { x: "Jul", y: 250 },
          { x: "Aug", y: 280 },
          { x: "Sep", y: 250 },
          { x: "Oct", y: 280 },
          { x: "Nov", y: 350 },
          { x: "Dec", y: 220 }
        ]
      },
      {
        name: "Total Data",
        color: "#EAECF0",
        // Total Data = Compliant Data + Non-compliant Data
        data: [
          { x: "Jan", y: 550 }, // 300 + 250
          { x: "Feb", y: 720 }, // 370 + 350
          { x: "Mar", y: 400 }, // 230 + 170
          { x: "Apr", y: 590 }, // 320 + 270
          { x: "May", y: 400 }, // 280 + 120
          { x: "Jun", y: 660 }, // 350 + 310
          { x: "Jul", y: 570 }, // 320 + 250
          { x: "Aug", y: 580 }, // 300 + 280
          { x: "Sep", y: 540 }, // 290 + 250
          { x: "Oct", y: 610 }, // 330 + 280
          { x: "Nov", y: 720 }, // 370 + 350
          { x: "Dec", y: 500 }  // 280 + 220
        ]
      }
    ],
    chart: {
      type: "bar",
      height: "320px",
      fontFamily: "Inter, sans-serif",
      toolbar: {
        show: false
      }
    },
    plotOptions: {
      bar: {
        horizontal: false,
        columnWidth: "70%",
        borderRadiusApplication: "end",
        borderRadius: 8
      }
    },
    tooltip: {
      shared: true,
      intersect: false,
      style: {
        fontFamily: "Inter, sans-serif"
      }
    },
    states: {
      hover: {
        filter: {
          type: "darken"
        }
      }
    },
    stroke: {
      show: true,
      width: 0,
      colors: ["transparent"]
    },
    grid: {
      show: false,
      strokeDashArray: 4,
      padding: {
        left: 2,
        right: 2,
        top: -14
      }
    },
    dataLabels: {
      enabled: false
    },
    legend: {
      show: false
    },
    xaxis: {
      categories: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
      floating: false,
      labels: {
        show: true,
        style: {
          fontFamily: "Inter, sans-serif",
          cssClass: "text-xs font-normal fill-gray-500 dark:fill-gray-400"
        }
      },
      axisBorder: {
        show: false
      },
      axisTicks: {
        show: false
      }
    },
    yaxis: {
      show: true,
      title: {
        text: "Debitur",
        style: {
          fontSize: "12px",
          fontWeight: 500,
          color: '#475467'
        }
      },
    },
    fill: {
      opacity: 1
    }
  };
</script>

<div class="">
  <Typography
    preset="Text lg/Semibold"
    text={$_("auditPerformanceTrends")}
  />
  <div class="flex justify-end gap-3">
    <div class="flex items-center">
      <Indicator size="sm" color="none" class="me-1.5 bg-utility-brand-700-light dark:bg-utility-brand-700-dark" />
      <Typography
        preset="Text sm/Regular"
        text={$_("compliantData")}
        color="tertiary"
      />
    </div>
    <div class="flex items-center">
      <Indicator size="sm" color="none" class="me-1.5 bg-utility-brand-500-light dark:bg-utility-brand-500-dark" />
      <Typography
        preset="Text sm/Regular"
        text={$_("nonCompliantData")}
        color="tertiary"
      />
    </div>
    <div class="flex items-center">
      <Indicator size="sm" color="none" class="me-1.5 bg-utility-gray-200-light dark:bg-utility-gray-200-dark" />
      <Typography
        preset="Text sm/Regular"
        text={$_("totalData")}
        color="tertiary"
      />
    </div>
  </div>
  <Chart {options} />
</div>