import { writable, derived } from "svelte/store";

export const cpuData = writable([0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
export const memoryData = writable([0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
export const diskData = writable([]);
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

export const optionData = derived(tabDatas, ($tabDatas, set) => {
  if ($tabDatas.selected === "CPU") {
    cpuData.subscribe(set);
  } else if ($tabDatas.selected === "Memory") {
    memoryData.subscribe(set);
  } else if ($tabDatas.selected === "Disk") {
    diskData.subscribe(set);
  } else if ($tabDatas.selected === "GPU") {
    set([50, 20, 50, 20, 50, 75, 10, 30, 60]);
  } else {
    set([25, 10, 80, 50, 95, 100, 75, 45, 30, 10]);
  }
});