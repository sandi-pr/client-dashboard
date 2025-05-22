<script>
  import { Search } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import TabBar from "$lib/components/TabBar.svelte";
  import {
  call,
  Typography,
  toCamelCase,
  tabDatas,
  cpuData, 
  realTimeDatasCPU,
  memoryData,
  realTimeDatasMemory,
  diskData,
  realTimeDatasDisk,
  gpuData,
  realTimeDatasGPU} from "$lib";
  import Chart from "$lib/components/Chart.svelte";
  import { onDestroy, onMount } from "svelte";

  const isFrontend = process.env.NODE_ENV === "development";

  // let tabDatas = $state({
  //   selected: "CPU",
  //   data: [
  //     {
  //       name: "CPU",
  //     },
  //     {
  //       name: "Memory",
  //     },
  //     {
  //       name: "Disk",
  //     },
  //     {
  //       name: "GPU",
  //     },
  //   ],
  // });

  let titleList = $state([
    {
      name: "CPU",
      value: "",
    },
    {
      name: "Memory",
      value: "",
    },
    {
      name: "Disk",
      value: "",
    },
    {
      name: "GPU",
      value: "",
    },
  ]);

  let titleListDisk = $state([]);

  // let realTimeDatasCPU = [
  //   {
  //     name: "Utilization",
  //     value: "68%",
  //   },
  //   {
  //     name: "Speed",
  //     value: "2,4 GHz",
  //   },
  //   {
  //     name: "Process",
  //     value: "216",
  //   },
  //   {
  //     name: "Thread",
  //     value: "2443",
  //   },
  //   {
  //     name: "Handles",
  //     value: "98078",
  //   },
  //   {
  //     name: "Up Time",
  //     value: "5:14:58:42",
  //   },
  // ];

  // let realTimeDatasMemory = [
  //   {
  //     name: "In Use",
  //     value: "8,3 GB",
  //   },
  //   {
  //     name: "Available",
  //     value: "7,1 GB",
  //   },
  //   {
  //     name: "Committed",
  //     value: "12,4/18,0 GB",
  //   },
  //   {
  //     name: "Cached",
  //     value: "7,0GB",
  //   },
  //   {
  //     name: "Paged Pool",
  //     value: "949 MB",
  //   },
  //   {
  //     name: "Non-Paged Pool",
  //     value: "396MB",
  //   },
  // ];

  // let realTimeDatasDisk = [
  //   {
  //     name: "Active Time",
  //     value: "1%",
  //   },
  //   {
  //     name: "Average Response Time",
  //     value: "0,2 ms",
  //   },
  //   {
  //     name: "Read Speed",
  //     value: "0,0 MB/s",
  //   },
  //   {
  //     name: "Write Speed",
  //     value: "0,0 MB/s",
  //   },
  // ];

  // let realTimeDatasGPU = [
  //   {
  //     name: "Utilization",
  //     value: "68%",
  //   },
  //   {
  //     name: "GPU Memory",
  //     value: "0,8/8,2 GB",
  //   },
  //   {
  //     name: "Dedicated GPU Memory",
  //     value: "451/512 MB",
  //   },
  //   {
  //     name: "Shared GPU Memory",
  //     value: "0,3/7,7 GB",
  //   },
  //   {
  //     name: "GPU Temperature",
  //     value: "80 °C",
  //   },
  // ];

  let realTimeDatas = $state([]);

  let generalDatasCPU = $state([
    {
      name: "Base speed",
      value: "",
    },
    {
      name: "Sockets",
      value: "",
    },
    {
      name: "Cores",
      value: "",
    },
    {
      name: "Logical Processors",
      value: "",
    },
    {
      name: "Virtualization",
      value: "",
    },
    {
      name: "Hyper-V Support",
      value: "",
    },
    {
      name: "L1 Cache",
      value: "",
    },
    {
      name: "L2 Cache",
      value: "",
    },
  ]);

  let generalDatasMemory = $state([
    {
      name: "Speed",
      value: "",
    },
    {
      name: "Slots Used",
      value: "",
    },
    {
      name: "Form Factor",
      value: "",
    },
    {
      name: "Hardware Reserved",
      value: "",
    },
  ]);

  let generalDatasDisk = $state([]);

  let generalDatasGPU = $state([
    {
      name: "Driver Version",
      value: "496.76",
    },
    // {
    //   name: "Driver Date",
    //   value: "2021-11-09",
    // },
    // {
    //   name: "DirectX Version",
    //   value: "DirectX 12 Ultimate",
    // },
    {
      name: "Physical Allocation",
      value: "PCI Slot 1",
    },
    {
      name: "Hardware Reserved Memory",
      value: "512 MB",
    },
  ]);

  let generalDatas = $state([]);

  let optionColor = $derived(
    $tabDatas.selected === "CPU"
      ? "#7F56D9"
      : $tabDatas.selected === "Memory"
        ? "#007AFF"
        : $tabDatas.selected === "Disk"
          ? "#B54708"
          : $tabDatas.selected === "GPU"
            ? "#067647"
            : "#7F56D9",
  );

  let optionData = $derived(
    $tabDatas.selected === "CPU"
      ? $cpuData
      : $tabDatas.selected === "Memory"
        ? $memoryData
        : $tabDatas.selected === "Disk"
          ? [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
          : $tabDatas.selected === "GPU"
            ? $gpuData
            : [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
  );

  let optionsDisks = $derived([]);

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
      animations: {
        enabled: false,
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
        name: $_(toCamelCase("Utilization")),
        data: optionData,
        color: optionColor,
      },
    ],
    xaxis: {
      title: {
        text: $_(toCamelCase("Real Time")) + " (m)",
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
        text: $_(toCamelCase("Utilization")) + " (%)",
      },
    },
  });

  let ws;
  let shouldReconnect = true;

  function connectWebSocket() {
    if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
      console.log("[Usage] Closing existing WebSocket connection...");
      shouldReconnect = false; 
      ws.close();
    }

    var base_url =
      (location.protocol == "https:" ? "wss:" : "ws:") +
      "//" +
      location.hostname +
      (location.port == "" ? "" : ":" + location.port);

    if (isFrontend) {
      ws = new WebSocket(
        `wss://tank.hexacode.ai/client-dashboard/usage/${$tabDatas.selected}/${encodeURIComponent(localStorage.getItem("token"))}`
      );
    } else {
      ws = new WebSocket(
        `${base_url}/client-dashboard/usage/${$tabDatas.selected}/${encodeURIComponent(localStorage.getItem("token"))}`
      );
    }

    ws.onopen = function (event) {
      console.log(`[Usage] [${$tabDatas.selected}] WebSocket connected, event`);
      shouldReconnect = true; 
    };

    ws.onerror = function (error) {
      console.error(`[Usage] [${$tabDatas.selected}] WebSocket error occurred:, error`);
    };

    ws.onmessage = function (event) {
      try {
        const data = JSON.parse(event.data);
        console.log("[Usage] Received WebSocket message:", data);

        switch (data.type) {
          case "cpu_realtime":
            updateCpuData(data);
            break;
          case "memory_realtime":
            updateMemoryData(data);
            break;
          case "disk_realtime":
            updateDiskData(data);
            break;
          case "gpu_realtime":
            updateGpuData(data);
            // console.log("[GPU] Received WebSocket message:", data);
            break;
          case "error":
            console.error("WebSocket error:", data.message);
            break;
          default:
            console.warn("Unknown message type:", data.type);
        }
      } catch (e) {
        console.error("Failed to parse WebSocket message", e);
      }
    };

    ws.onclose = () => {
      console.log("[Usage] WebSocket disconnected.");
      // if (shouldReconnect) {
      //   setTimeout(() => {
      //     connectWebSocket();
      //   }, 3000);
      // } else {
      //   console.log("[Usage] Manual disconnect, skipping reconnect.");
      // }
    };
  }

  function updateCpuData(data) {
    cpuData.update(current => {
      const newValue = Math.round(data.utilization);
      const updated = [...current, newValue];
      if (updated.length > 10) updated.shift();
      return updated;
    });

    realTimeDatasCPU.set([
      {
        name: "Utilization",
        value: `${data.utilization.toFixed(0)}%`,
      },
      {
        name: "Speed",
        value: `${(data.speed / 1000).toFixed(1)} GHz`,
      },
      {
        name: "Process",
        value: `${data.process ?? 0}`,
      },
      {
        name: "Thread",
        value: `${data.thread ?? 0}`,
      },
      // {
      //   name: "Handles",
      //   value: ${data.handles ?? 0},
      // },
      {
        name: "Up Time",
        value: `${data.uptimeReadable ?? "-"}`,
      },
    ]);
  }

  function formatBytes(bytes) {
    const gb = 1024 * 1024 * 1024;
    const mb = 1024 * 1024;
    if (bytes >= gb) return `${(bytes / gb).toFixed(1)} GB`;
    if (bytes >= mb) return `${(bytes / mb).toFixed(0)} MB`;
    return `${bytes} B`;
  }

  function updateMemoryData(data) {
    memoryData.update(current => {
      const percentageUsed = Math.round((data.used / data.total) * 100);
      const updated = [...current, percentageUsed];
      if (updated.length > 10) updated.shift();
      return updated;
    });

    realTimeDatasMemory.set([
      {
        name: "In Use",
        value: formatBytes(data.used),
      },
      {
        name: "Available",
        value: formatBytes(data.available),
      },
      {
        name: "Committed",
        value: `${formatBytes(data.committed)}/${formatBytes(data.total)}`,
      },
      {
        name: "Cached",
        value: formatBytes(data.cached),
      },
      {
        name: "Paged Pool",
        value: formatBytes(data.paged_pool),
      },
      {
        name: "Non-Paged Pool",
        value: formatBytes(data.non_paged_pool),
      },
    ]);
  }

  function updateDiskData(data) {
    diskData.update(current => {
      const updated = current.map(d => {
        const match = data.list.find(item => item.device === d.device);
        if (match) {
          const utilizationValue = parseFloat(match.utilization?.replace("%", "") || "0");
          const roundedValue = Math.round(utilizationValue);

          const newHistory = [...d.history, roundedValue];
          if (newHistory.length > 10) newHistory.shift();

          return { ...d, history: newHistory };
        }
        return d;
      });

      data.list.forEach(item => {
        const exists = current.some(d => d.device === item.device);
        if (!exists) {
          const utilizationValue = parseFloat(item.utilization?.replace("%", "") || "0");
          const roundedValue = Math.round(utilizationValue);

          updated.push({
            device: item.device,
            history: new Array(9).fill(0).concat(roundedValue),
          });
        }
      });

      return updated;
    });

    for (var i = 0; i < data.list.length; i++) {
      optionsDisks[i].series[0].data = $diskData[i].history;
    }

    console.log("Disk Data:", optionsDisks);

    realTimeDatasDisk.update(current => {
      const updated = [...current];

      for (const item of data.list) {
        const index = updated.findIndex(d => d.device === item.device);

        const newData = {
          device: item.device,
          data: [
            { name: "Active Time", value: item.utilization || "0%" },
            { name: "Average Response Time", value: item.avg_response || "0 ms" },
            { name: "Read Speed", value: item.read_speed || "0 MB/s" },
            { name: "Write Speed", value: item.write_speed || "0 MB/s" },
          ]
        };

        if (index !== -1) {
          updated[index] = newData;
        } else {
          updated.push(newData);
        }
      }

      return updated;
    });
  }

  function updateGpuData(data) {
    gpuData.update(current => {
      const updated = [...current, parseFloat(data.utilization)];
      if (updated.length > 10) updated.shift();
      return updated;
    });

    realTimeDatasGPU.set([
      {
        name: "Utilization",
        value: data.utilization || "N/A",
      },
      {
        name: "GPU Memory",
        value: data.memory || "N/A",
      },
      {
        name: "Dedicated GPU Memory",
        value: data.dedicatedMemory || "N/A",
      },
      {
        name: "GPU Temperature",
        value: data.temperature || "N/A",
      },
    ]);
  }


  async function getGeneralIformationCPU() {
    let token = localStorage.getItem("token");
    if (!token) return;

    const response = await $call.call("client-dashboard", "get-cpu-info", {
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

    const data = response.result;
    console.log("getGeneralIformationCPU", data);

    titleList[0].value = `CPU ${data.model ?? "-"}`;

    generalDatasCPU = [
      {
        name: "Base speed",
        value: data.base_speed ? `${(data.base_speed / 1000).toFixed(2)} GHz` : "-",
      },
      {
        name: "Sockets",
        value: "",
      },
      {
        name: "Cores",
        value: `${data.physical_cores ?? "-"}`,
      },
      {
        name: "Logical Processors",
        value: `${data.logical_cores ?? "-"}`,
      },
      {
        name: "Virtualization",
        value: data.supports_virtualization ? "Enabled" : "Disabled",
      },
      {
        name: "Hyper-V Support",
        value: data.supports_hyperthreading ? "Yes" : "No",
      },
      {
        name: "L1 Cache",
        value: `${(data.l1_data_cache + data.l1_instruction_cache) || 0} KB`,
      },
      {
        name: "L2 Cache",
        value: `${data.l2_cache ?? 0} KB`,
      },
    ];
  }

  async function getGeneralIformationMemory() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "get-memory-info", {
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

      const data = response.result;
      console.log("Memory Information:", data);

      generalDatasMemory = [
        {
          name: "Speed",
          value: data.speed || "-",
        },
        {
          name: "Slots Used",
          value: `${data.slots_used || 0} of ?`,
        },
        {
          name: "Form Factor",
          value: data.form_factor || "-",
        },
        {
          name: "Hardware Reserved",
          value: data.hardware_reserved ?? "-",
        },
      ];
    }
  }

  async function getGeneralIformationDisk() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "get-disk-info", {
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

      const data = response.result;
      console.log("Disk Information:", data);

      for (let i = 0; i < data.list.length; i++) {
        const item = data.list[i];

        if (!generalDatasDisk[i]) {
          generalDatasDisk.push([
            { name: "Capacity", value: item.capacity },
            { name: "Formatted Capacity", value: item.formatted_capacity },
            { name: "System Disk", value: item.system_disk },
            { name: "Page File", value: item.page_file },
            { name: "Type", value: item.type },
          ]);
        } else {
          generalDatasDisk[i] = [
            { name: "Capacity", value: item.capacity },
            { name: "Formatted Capacity", value: item.formatted_capacity },
            { name: "System Disk", value: item.system_disk },
            { name: "Page File", value: item.page_file },
            { name: "Type", value: item.type },
          ];
        }

        titleListDisk.push(item.model);

        var newOptions = options;

        newOptions.series[0].data = [0,0,0,0,0,0,0,0,0,0];

        optionsDisks.push(newOptions);
      }
    }
  }

  async function getGeneralIformationGPU() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "get-gpu-info", {
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

      const data = response.result;
      console.log("GPU Information:", data);

      titleList[3].value = `${data.model ?? "-"}`;

      generalDatasGPU = [
      {
        name: "Driver Version",
        value: data.driver_version || "N/A",
      },
      // {
      //   name: "Driver Date",
      //   value: data.driver_date || "N/A", // kalau tidak ada di response, kasih default
      // },
      // {
      //   name: "DirectX Version",
      //   value: "DirectX 12 Ultimate", // hardcoded karena tidak ada di response
      // },
      {
        name: "Physical Allocation",
        value: `${data.physical_allocation?.replace("00000000:", "") || "N/A"}`,
      },
      {
        name: "Hardware Reserved Memory",
        value: data.hardware_reserved || "N/A",
      },
    ];
    }
  }

  function disconnectWebSocket() {
    if (ws) {
      ws.close();
    }
  }

  onMount(async () => {
    connectWebSocket();
    getGeneralIformationCPU();
    getGeneralIformationMemory();
    getGeneralIformationDisk();
    getGeneralIformationGPU();
  });

  onDestroy(() => {
    disconnectWebSocket();
  });

  $effect(() => {
    generalDatas =
      $tabDatas.selected === "CPU"
        ? generalDatasCPU
        : $tabDatas.selected === "Memory"
          ? generalDatasMemory
          : $tabDatas.selected === "Disk"
            ? generalDatasDisk
            : $tabDatas.selected === "GPU"
              ? generalDatasGPU
              : generalDatasCPU;

    realTimeDatas =
      $tabDatas.selected === "CPU"
        ? $realTimeDatasCPU
        : $tabDatas.selected === "Memory"
          ? $realTimeDatasMemory
          : $tabDatas.selected === "Disk"
            ? $realTimeDatasDisk
            : $tabDatas.selected === "GPU"
              ? $realTimeDatasGPU
              : realTimeDatasCPU;
  });
</script>

<div class="m-4 mb-32 md:m-4">
  <div class="flex justify-between items-center">
    <Typography preset="Display sm/Semibold" text="Usage" />
    <div class="hidden md:block">
      <Search class="bg-primary-light dark:bg-primary-dark" />
    </div>
  </div>

  <TabBar tabDatas={$tabDatas} on:tabChange={(e) => {
      $tabDatas.selected = e.detail;
      // if ($tabDatas.selected !== "Disk") {
        realTimeDatas = [];
        generalDatas = [];
      // }
      connectWebSocket();
    }} />

  {#if $tabDatas.selected}
    {#if $tabDatas.selected !== "Disk"}
      <Typography
        preset="Text sm/Semibold"
        text={titleList.find((item) => item.name === $tabDatas.selected).value}
        propClass="mt-8 mb-4"
      />
      <div class="md:max-w-[75%]">
        <Chart dataChart={options} />
      </div>
      <Typography
        preset="Text xs/Medium"
        text={$_(toCamelCase("Real Time Information"))}
        propClass="mt-8 mb-4"
      />
      <div class="flex flex-wrap gap-6">
        {#if realTimeDatas.length === 0}
          <div class="flex flex-col">
            <Typography
              preset="Text sm/Medium"
              text={$_(toCamelCase("No data available"))}
              color="quarterary"
            />
            <Typography preset="Text lg/Semibold" text={$_(toCamelCase("Try again later"))} />
          </div>
        {:else}
          {#each realTimeDatas as data}
            <div class="flex flex-col">
              <Typography
                preset="Text sm/Medium"
                text={$_(toCamelCase(data.name))}
                color="quarterary"
              />
              <Typography preset="Text lg/Semibold" text={data.value} />
            </div>
          {/each}
        {/if}
      </div>
      <Typography
        preset="Text xs/Medium"
        text={$_(toCamelCase("General Information"))}
        propClass="mt-8 mb-4"
      />
      <div class="md:max-w-[50%]">
        {#if generalDatas.length === 0}
          <div class="flex flex-col">
            <Typography
              preset="Text sm/Medium"
              text={$_(toCamelCase("No data available"))}
              color="quarterary"
            />
            <Typography preset="Text lg/Semibold" text={$_(toCamelCase("Try again later"))} />
          </div>
        {:else}
          {#each generalDatas as data}
            <div class="w-full flex justify-between">
              <Typography
                preset="Text sm/Medium"
                text={$_(toCamelCase(data.name))}
                color="quarterary"
              />
              <Typography preset="Text sm/Semibold" text={data.value} color="secondary" />
            </div>
          {/each}
        {/if}
      </div>
    {:else}
      {#if realTimeDatas.length > 0}
        {#each optionsDisks as disk, i}
          <Typography
            preset="Text sm/Semibold"
            text={titleListDisk[i]}
            propClass="mt-8 mb-4"
          />
          <div class="md:max-w-[75%]">
              <Chart dataChart={disk} />
          </div>
          <Typography
            preset="Text xs/Medium"
            text={$_(toCamelCase("Real Time Information"))}
            propClass="mt-8 mb-4"
          />
          <div class="flex flex-wrap gap-6">
            {#if realTimeDatas[i].data.length === 0}
              <div class="flex flex-col">
                <Typography
                  preset="Text sm/Medium"
                  text={$_(toCamelCase("No data available"))}
                  color="quarterary"
                />
                <Typography preset="Text lg/Semibold" text={$_(toCamelCase("Try again later"))} />
              </div>
            {:else}
              {#each realTimeDatas[i].data as data}
                <div class="flex flex-col">
                  <Typography
                    preset="Text sm/Medium"
                    text={$_(toCamelCase(data.name))}
                    color="quarterary"
                  />
                  <Typography preset="Text lg/Semibold" text={data.value} />
                </div>
              {/each}
            {/if}
          </div>
          <Typography
            preset="Text xs/Medium"
            text={$_(toCamelCase("General Information"))}
            propClass="mt-8 mb-4"
          />
          <div class="md:max-w-[50%]">
            {#if generalDatasDisk[i].length === 0}
              <div class="flex flex-col">
                <Typography
                  preset="Text sm/Medium"
                  text={$_(toCamelCase("No data available"))}
                  color="quarterary"
                />
                <Typography preset="Text lg/Semibold" text={$_(toCamelCase("Try again later"))} />
              </div>
            {:else}
              {#each generalDatasDisk[i] as data}
                <div class="w-full flex justify-between">
                  <Typography
                    preset="Text sm/Medium"
                    text={$_(toCamelCase(data.name))}
                    color="quarterary"
                  />
                  <Typography preset="Text sm/Semibold" text={data.value} color="secondary" />
                </div>
              {/each}
            {/if}
          </div>
        {/each}
      {/if}
    {/if}
  {:else}
    <div>Hello World</div>
  {/if}
</div>
