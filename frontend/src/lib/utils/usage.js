import { writable, derived } from "svelte/store";

export const cpuData = writable([0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
export const memoryData = writable([0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
export const diskData = writable([]);
export const gpuData = writable([0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
export const tabDatas = writable({ 
    selected: "CPU" , 
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

export const realTimeDatasCPU = writable([
    { name: "Utilization", value: "0%" },
    { name: "Speed", value: "0 GHz" },
    { name: "Process", value: "0" },
    { name: "Thread", value: "0" },
    // { name: "Handles", value: "0" },
    { name: "Up Time", value: "0s" },
]);

export const realTimeDatasMemory = writable([
    {
        name: "In Use",
        value: "0 GB",
    },
    {
        name: "Available",
        value: "0 GB",
    },
    {
        name: "Committed",
        value: "0 GB",
    },
    {
        name: "Cached",
        value: "0 GB",
    },
    {
        name: "Paged Pool",
        value: "0 MB",
    },
    {
        name: "Non-Paged Pool",
        value: "0 MB",
    },
]);

export const realTimeDatasDisk = writable([]);

export const realTimeDatasGPU = writable([
  {
    name: "Utilization",
    value: "",
  },
  {
    name: "GPU Memory",
    value: "",
  },
  {
    name: "Dedicated GPU Memory",
    value: "",
  },
  // {
  //   name: "Shared GPU Memory",
  //   value: "0,3/7,7 GB",
  // },
  {
    name: "GPU Temperature",
    value: "",
  },
]);

export const optionData = derived(tabDatas, ($tabDatas, set) => {
  if ($tabDatas.selected === "CPU") {
    cpuData.subscribe(set);
  } else if ($tabDatas.selected === "Memory") {
    memoryData.subscribe(set);
  } else if ($tabDatas.selected === "Disk") {
    diskData.subscribe(set);
  } else if ($tabDatas.selected === "GPU") {
    gpuData.subscribe(set);
  } else {
    set([0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
  }
});