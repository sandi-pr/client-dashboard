<script>
  import { Chart, Search } from "flowbite-svelte";
  import TabBar from "$lib/components/TabBar.svelte";
  import { Typography } from "$lib";

  let tabDatas = $state({
    selected: "CPU",
    data: [
      {
        name: "CPU",
      },
      {
        name: "Memory",
      },
      {
        name: "Disk",
      },
      {
        name: "GPU",
      },
    ],
  });

  let realTimeDatasCPU = [
    {
      name: "Utilization",
      value: "68%",
    },
    {
      name: "Speed",
      value: "2,4 GHz",
    },
    {
      name: "Process",
      value: "216",
    },
    {
      name: "Thread",
      value: "2443",
    },
    {
      name: "Handles",
      value: "98078",
    },
    {
      name: "Up Time",
      value: "5:14:58:42",
    },
  ];

  let realTimeDatasMemory = [
    {
      name: "In Use",
      value: "8,3 GB",
    },
    {
      name: "Available",
      value: "7,1 GB",
    },
    {
      name: "Committed",
      value: "12,4/18,0 GB",
    },
    {
      name: "Cached",
      value: "7,0GB",
    },
    {
      name: "Paged Pool",
      value: "949 MB",
    },
    {
      name: "Non-Paged Pool",
      value: "396MB",
    },
  ];

  let realTimeDatasDisk = [
    {
      name: "Active Time",
      value: "1%",
    },
    {
      name: "Average Response Time",
      value: "0,2 ms",
    },
    {
      name: "Read Speed",
      value: "0,0 MB/s",
    },
    {
      name: "Write Speed",
      value: "0,0 MB/s",
    },
  ];

  let realTimeDatasGPU = [
    {
      name: "Utilization",
      value: "68%",
    },
    {
      name: "GPU Memory",
      value: "0,8/8,2 GB",
    },
    {
      name: "Dedicated GPU Memory",
      value: "451/512 MB",
    },
    {
      name: "Shared GPU Memory",
      value: "0,3/7,7 GB",
    },
    {
      name: "GPU Temperature",
      value: "80 °C",
    },
  ];

  let realTimeDatas = $state([]);

  let generalDatasCPU = [
    {
      name: "Base speed",
      value: "2,50 GHz",
    },
    {
      name: "Sockets",
      value: "1",
    },
    {
      name: "Cores",
      value: "4",
    },
    {
      name: "Logical Processors",
      value: "4",
    },
    {
      name: "Virtualization",
      value: "Disabled",
    },
    {
      name: "Hyper-V Support",
      value: "Yes",
    },
    {
      name: "L1 Cache",
      value: "512 KB",
    },
    {
      name: "L2 Cache",
      value: "4.0 MB",
    },
  ];

  let generalDatasMemory = [
    {
      name: "Speed",
      value: "2400 MHz",
    },
    {
      name: "Slots Used",
      value: "2 of 4",
    },
    {
      name: "Form Factor",
      value: "DIMM",
    },
    {
      name: "Hardware Reserved",
      value: "555 MB",
    },
  ];

  let generalDatasDisk = [
    {
      name: "Capacity",
      value: "477 GB",
    },
    {
      name: "Formatted Capacity",
      value: "477 GB",
    },
    {
      name: "System Disk",
      value: "Yes",
    },
    {
      name: "Page File",
      value: "Yes",
    },
    {
      name: "Type",
      value: "SSD",
    },
  ];

  let generalDatasGPU = [
    {
      name: "Driver Version",
      value: "496.76",
    },
    {
      name: "Driver Date",
      value: "2021-11-09",
    },
    {
      name: "DirectX Version",
      value: "DirectX 12 Ultimate",
    },
    {
      name: "Physical Allocation",
      value: "PCI Slot 1",
    },
    {
      name: "Hardware Reserved Memory",
      value: "512 MB",
    },
  ];

  let generalDatas = $state([]);

  let optionColor = $derived(
    tabDatas.selected === "CPU"
      ? "#7F56D9"
      : tabDatas.selected === "Memory"
        ? "#007AFF"
        : tabDatas.selected === "Disk"
          ? "#B54708"
          : tabDatas.selected === "GPU"
            ? "#067647"
            : "#7F56D9",
  );

  let optionData = $derived(
    tabDatas.selected === "CPU"
      ? [25, 10, 80, 50, 95, 100, 75, 45, 30, 10]
      : tabDatas.selected === "Memory"
        ? [10, 20, 30, 40, 50, 60, 70, 80, 90]
        : tabDatas.selected === "Disk"
          ? [20, 25, 10, 90, 65, 60, 70, 100, 30]
          : tabDatas.selected === "GPU"
            ? [50, 20, 50, 20, 50, 75, 10, 30, 60]
            : [25, 10, 80, 50, 95, 100, 75, 45, 30, 10],
  );

  let options = $derived({
    chart: {
      height: "400px",
      maxWidth: "100%",
      type: "area",
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
        shade: "#1C64F2",
        gradientToColors: ["#1C64F2"],
      },
    },
    dataLabels: {
      enabled: false,
    },
    stroke: {
      width: 4,
    },
    grid: {
      show: true,
      strokeDashArray: 4,
      padding: {
        left: 2,
        right: 2,
        top: 0,
      },
    },
    series: [
      {
        name: "Utilization",
        data: optionData,
        color: optionColor,
      },
    ],
    xaxis: {
      title: {
        text: "Real Time (m)",
      },
      categories: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10"],
      labels: {
        show: true,
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
        text: "Utilization (%)",
      },
    },
  });

  $effect(() => {
    generalDatas =
      tabDatas.selected === "CPU"
        ? generalDatasCPU
        : tabDatas.selected === "Memory"
          ? generalDatasMemory
          : tabDatas.selected === "Disk"
            ? generalDatasDisk
            : tabDatas.selected === "GPU"
              ? generalDatasGPU
              : generalDatasCPU;

    realTimeDatas =
      tabDatas.selected === "CPU"
        ? realTimeDatasCPU
        : tabDatas.selected === "Memory"
          ? realTimeDatasMemory
          : tabDatas.selected === "Disk"
            ? realTimeDatasDisk
            : tabDatas.selected === "GPU"
              ? realTimeDatasGPU
              : realTimeDatasCPU;
  });
</script>

<div class="m-4 mb-32 md:m-4">
  <div class="flex justify-between items-center">
    <Typography preset="Display sm/Semibold" text="Usage" />
    <div>
      <Search class="bg-primary-light dark:bg-primary-dark hidden md:block" />
    </div>
  </div>

  <TabBar {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

  {#if tabDatas.selected}
    <Typography
      preset="Text sm/Semibold"
      text="CPU Intel Core i7-12700K @ 3.60GHz"
      propClass="mt-8 mb-4"
    />
    <div class="md:max-w-[75%]">
      <Chart {options} />
    </div>
    <Typography preset="Text xs/Medium" text="Real Time Information" propClass="mt-8 mb-4" />
    <div class="flex flex-wrap gap-6">
      {#if realTimeDatas.length === 0}
        <div class="flex flex-col">
          <Typography preset="Text sm/Medium" text="No data available" color="quarterary" />
          <Typography preset="Text lg/Semibold" text="Try again later" />
        </div>
      {:else}
        {#each realTimeDatas as data}
          <div class="flex flex-col">
            <Typography preset="Text sm/Medium" text={data.name} color="quarterary" />
            <Typography preset="Text lg/Semibold" text={data.value} />
          </div>
        {/each}
      {/if}
    </div>
    <Typography preset="Text xs/Medium" text="General Information" propClass="mt-8 mb-4" />
    <div class="md:max-w-[50%]">
      {#if generalDatas.length === 0}
        <div class="flex flex-col">
          <Typography preset="Text sm/Medium" text="No data available" color="quarterary" />
          <Typography preset="Text lg/Semibold" text="Try again later" />
        </div>
      {:else}
        {#each generalDatas as data}
          <div class="w-full flex justify-between">
            <Typography preset="Text sm/Medium" text={data.name} color="quarterary" />
            <Typography preset="Text sm/Semibold" text={data.value} color="secondary" />
          </div>
        {/each}
      {/if}
    </div>
  {:else}
    <div>Hello World</div>
  {/if}
</div>
