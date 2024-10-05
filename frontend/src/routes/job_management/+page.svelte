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
  import { P, Heading, Input, Button, Alert, Modal } from "flowbite-svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import { ToastContainer, triggerToast, StateCard, Typography } from "$lib";
  import { fly } from "svelte/transition";
  import { _ } from "svelte-i18n";
  import { onMount } from "svelte";
  import { dndzone, TRIGGERS } from "svelte-dnd-action";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import DotsGridIcon from "$lib/icon/dots-grid.svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { InfoCircleSolid, ExclamationCircleSolid } from "flowbite-svelte-icons";

  import { call } from "$lib";

  let tmpIdPosition = $state("");
  let tmpJobLevel = $state("");
  let tmpJobPosition = $state("");

  let jobLevels = $state([]);

  let datas = $state([]);
  let configs = $state({});

  let gettingDatas = $state(false);

  let pagination = $state({ page: 1, pageSize: 5, pageCount: 1, total: 0 });

  let addPositionsModal = $state(false);
  let editPositionModal = $state(false);
  let dropdownOpen = $state(false);
  let deleteModal = $state(false);

  let selectedJobLevel = $state("");

  let tabDatas = $state({
    selected: "JobLevel",
    data: [
      {
        name: "JobLevel",
      },
      {
        name: "JobPosition",
      },
    ],
  });

  function addJobLevel() {
    if (tmpJobLevel === "") {
      return;
    }
    const newJobLevel = {
      // dummy id
      id: Date.now(),
      title: tmpJobLevel,
    };
    jobLevels = [...jobLevels, newJobLevel];

    saveJobLevels();

    tmpJobLevel = "";
  }

  function deleteJobLevel(id) {
    jobLevels = jobLevels.filter((jobLevel) => jobLevel.id !== id);
    saveJobLevels();
  }

  function handleDndConsider(e) {
    const {
      items: newItems,
      info: { trigger },
    } = e.detail;
    if (trigger === TRIGGERS.DRAG_STARTED) {
      // Optional: Handle logic when drag starts
    }
    jobLevels = newItems;
    saveJobLevels();
  }

  function handleDndFinalize(e) {
    const { items: newItems } = e.detail;
    jobLevels = newItems;
    saveJobLevels();
  }

  async function saveJobLevels() {
    // if (jobLevels.length === 0) {
    //   triggerToast("error");
    //   return;
    // } else {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "create-job-level", {
          URLQuery: {},
          URLParams: {},
          Params: {
            job_level: JSON.stringify(jobLevels),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log(response);

        if (response.status === 200) {
          // triggerToast("created");
        }
      } catch (error) {
        console.log(error);
        triggerToast("error");
      }
    }
    // }
  }

  async function getJobLevels() {
    let token = localStorage.getItem("token");
    if (token) {
      gettingDatas = true;
      try {
        const response = await $call.call("client-dashboard", "get-job-level", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          jobLevels = response.result.job_level;
        }

        gettingDatas = false;
      } catch (error) {
        triggerToast("error");
        gettingDatas = false;
        console.log(error);
      }
    }
  }

  async function getPositions(page) {
    if (page > pagination.pageCount) return;
    let token = localStorage.getItem("token");
    if (token !== null) {
      gettingDatas = true;
      let response = await $call.call("client-dashboard", "get-job-position", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        triggerToast("error");
        gettingDatas = false;
        return;
      }

      console.log("response", response);

      if (response.result.data !== null) {
        configs = {
          search: true,
        };
        let promises2 = response.result.data.map(async (position) => {
          return {
            id: position.id,
            job_position: position.job_position,
            job_level: position.job_level,
            actions: {
              objects: [
                {
                  onclick: () => {
                    tmpIdPosition = position.id;
                    selectedJobLevel = position.job_level;
                    tmpJobPosition = position.job_position;
                    editPositionModal = true;
                  },
                  name: "Edit",
                },
                {
                  onclick: () => {
                    deleteModal = true;
                    tmpIdPosition = position.id;
                  },
                  name: "Delete",
                },
              ],
            },
          };
        });

        datas = await Promise.all(promises2);
      }

      gettingDatas = false;
    }
  }

  async function createJobPosition() {
    let token = localStorage.getItem("token");
    if (token) {
      datas.push({ id: Date.now(), job_position: tmpJobPosition, job_level: selectedJobLevel });
      try {
        const response = await $call.call("client-dashboard", "save-job-position", {
          URLQuery: {},
          URLParams: {},
          Params: {
            job_position: JSON.stringify(datas),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          tmpJobPosition = "";
          selectedJobLevel = "";
          addPositionsModal = false;
          getPositions();
          triggerToast("created");
          return;
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function editJobPosition() {
    datas = datas.map((data) => {
      if (data.id === tmpIdPosition) {
        return {
          id: data.id,
          job_position: tmpJobPosition,
          job_level: selectedJobLevel,
        };
      } else {
        return data;
      }
    });
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "save-job-position", {
          URLQuery: {},
          URLParams: {},
          Params: {
            job_position: JSON.stringify(datas),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          tmpIdPosition = "";
          tmpJobPosition = "";
          selectedJobLevel = "";
          editPositionModal = false;
          getPositions();
          triggerToast("updated");
          return;
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function deleteJobPosition() {
    datas = datas.filter((data) => data.id !== tmpIdPosition);
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "save-job-position", {
          URLQuery: {},
          URLParams: {},
          Params: {
            job_position: JSON.stringify(datas),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          tmpIdPosition = "";
          tmpJobPosition = "";
          selectedJobLevel = "";
          deleteModal = false;
          getPositions();
          triggerToast("deleted");
          return;
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function next() {
    getPositions(pagination.page + 1);
  }

  async function previous() {
    getPositions(pagination.page - 1);
  }

  let buttonDisable = $state(false);
  let buttonDisabled = $state(false);

  $effect(() => {
    buttonDisable = jobLevels.length < 1;
    buttonDisabled = selectedJobLevel < 1 || tmpJobPosition < 1;
  });

  onMount(async () => {
    getJobLevels();
    getPositions();
  });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- Header -->
      <div class="flex items-center pb-4">
        <div class="flex flex-col gap-2 grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {$_("jobManagement")}
            </Heading>
          </div>
        </div>
        <div class="flex gap-2">
          <!-- <Button
            class="px-3 focus-within:!ring-1"
            disabled={buttonDisable}
            on:click={saveJobLevels}>{$_("save")}</Button
          > -->
          <!-- <Button class="px-3 focus-within:!ring-1">{$_("importCsv")}</Button> -->
        </div>
      </div>

      <!-- Submenu -->
      <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

      {#if tabDatas.selected === "JobLevel"}
        <div class="flex gap-6 pt-8">
          <div class="w-60 space-y-2">
            <div>
              <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                {$_("newJobLevel")}
              </P>
              <P class="text-secondary-light dark:text-secondary-dark">
                {$_("adjust&MoveOrderByDragging")}
              </P>
            </div>
            <Input bind:value={tmpJobLevel} class="bg-primary-light dark:bg-secondary-dark" />
            <Button color="light" on:click={addJobLevel}>{$_("add")}</Button>
          </div>
          <div class="flex flex-col grow gap-6">
            {#if jobLevels.length > 0}
              <div class="bg-secondary-light dark:bg-secondary-dark rounded-lg py-10 pl-4">
                <div
                  use:dndzone={{ items: jobLevels, flipDurationMs: 150 }}
                  onconsider={handleDndConsider}
                  onfinalize={handleDndFinalize}
                  class="flex flex-col gap-4"
                >
                  {#each jobLevels as job (job.id)}
                    <div
                      class="px-4 py-2 bg-white dark:bg-tertiary-dark rounded-lg backdrop-blur-[26px] justify-between items-center flex w-2/5"
                    >
                      <div class="min-w-[172px] self-stretch items-center gap-4 flex">
                        <DotsGridIcon
                          class="w-6 h-6 text-fg-quaternary-light dark:text-fg-quaternary-dark"
                        />
                        <P size="lg" weight="semibold">{job.title}</P>
                      </div>
                      <Button on:click={() => deleteJobLevel(job.id)} color="light" class="w-max"
                        >{$_("delete")}</Button
                      >
                    </div>
                  {/each}
                </div>
              </div>
            {:else if gettingDatas}
              <StateCard
                illustration={IllustrationProgress}
                description={$_("fetchingDataPleaseWait")}
              />
            {:else}
              <StateCard
                illustration={IllustrationEmpty}
                description={$_("youDontHaveAnyDataYet")}
              />
            {/if}
          </div>
        </div>
      {:else if tabDatas.selected === "JobPosition"}
        <div class="pt-6">
          {#if datas.length > 0}
            <div>
              <div class="flex justify-between">
                <Heading
                  tag="h6"
                  class="font-jakarta w-fit text-hc-graylight-900 dark:text-hc-graydark-50"
                  >{$_("jobPositionsList")}</Heading
                >
                <Button class="w-max" on:click={() => (addPositionsModal = true)} color="light">
                  {$_("addNewPosition")}
                </Button>
              </div>
              <DataTable bind:datas bind:configs bind:next bind:previous />
            </div>
          {:else if gettingDatas}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else}
            <StateCard
              illustration={IllustrationEmpty}
              description={$_("youDontHaveAnyDataYet")}
              actionText={$_("addNewPosition")}
              funcButton={() => (addPositionsModal = true)}
            />
          {/if}
        </div>
      {/if}
    </div>
  </div>
  <ToastContainer />
</div>

<Modal title={$_("addNewPosition")} bind:open={addPositionsModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5 relative">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("jobLevel")}*`} />
      <Button
        color="light"
        class="flex justify-between"
        on:click={() => (dropdownOpen = !dropdownOpen)}
      >
        {`${selectedJobLevel ? selectedJobLevel : $_("selectJobLevel")}`}
        <ChevronDownOutline class="w-6 h-6 text-fg-quaternary-light dark:text-quarterary-dark" />
      </Button>

      {#if jobLevels.length > 0}
        <div
          class={`z-50 ${dropdownOpen ? "absolute top-20" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}
        >
          <ul class="py-2 text-sm text-gray-700 dark:text-gray-200">
            {#each jobLevels as job (job.id)}
              <button
                class="flex p-2"
                onclick={() => {
                  selectedJobLevel = job.title;
                  dropdownOpen = false;
                }}>{job.title}</button
              >
            {/each}
          </ul>
        </div>
      {:else}
        <Alert color="red">
          <InfoCircleSolid slot="icon" class="w-5 h-5" />
          {$_("pleaseAddJobLevelsFirst")}
        </Alert>
      {/if}
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("positionName")}*`} />
      <Input bind:value={tmpJobPosition} placeholder={$_("placeholderPositionName")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={createJobPosition} disabled={buttonDisabled} class="w-full"
      >{$_("addNewPosition")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editPosition")} bind:open={editPositionModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5 relative">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("jobLevel")}*`} />
      <Button
        color="light"
        class="flex justify-between"
        on:click={() => (dropdownOpen = !dropdownOpen)}
      >
        {`${selectedJobLevel ? selectedJobLevel : $_("selectJobLevel")}`}
        <ChevronDownOutline class="w-6 h-6 text-fg-quaternary-light dark:text-quarterary-dark" />
      </Button>
      <div
        class={`z-50 ${dropdownOpen ? "absolute top-20" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}
      >
        <ul class="py-2 text-sm text-gray-700 dark:text-gray-200">
          {#if jobLevels.length > 0}
            {#each jobLevels as job (job.id)}
              <button
                class="flex p-2"
                onclick={() => {
                  selectedJobLevel = job.title;
                  dropdownOpen = false;
                }}>{job.title}</button
              >
            {/each}
          {:else}
            <Alert color="yellow">
              {$_("pleaseAddJobLevelsFirst")}
            </Alert>
          {/if}
        </ul>
      </div>
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("positionName")}*`} />
      <Input bind:value={tmpJobPosition} placeholder={$_("placeholderPositionName")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={editJobPosition} disabled={buttonDisabled} class="w-full"
      >{$_("editPosition")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteJobPosition} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>
