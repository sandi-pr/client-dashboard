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
  import { P, Heading, Input, Button, Toast } from "flowbite-svelte";
  import { CheckCircleSolid, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { ToastContainer, triggerToast } from "$lib";
  import { fly } from "svelte/transition";
  import { _ } from "svelte-i18n";
  import { onMount } from "svelte";
  import { dndzone, TRIGGERS } from "svelte-dnd-action";
  import { StateCard } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import DotsGridIcon from "$lib/icon/dots-grid.svelte";
  import { CR3Call } from "hexacomponents-svelte";

  import { call } from "$lib";

  let tmpJobLevel = "";

  let jobLevels = [];

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

    tmpJobLevel = "";
  }

  function deleteJobLevel(id) {
    jobLevels = jobLevels.filter((jobLevel) => jobLevel.id !== id);
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
  }

  function handleDndFinalize(e) {
    const { items: newItems } = e.detail;
    jobLevels = newItems;
  }

  async function saveJobLevels() {
    if (jobLevels.length === 0) {
      triggerToast("error");
      return;
    } else {
      let token = sessionStorage.getItem("token");
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
            triggerToast("created");
          }
        } catch (error) {
          console.log(error);
          triggerToast("error");
        }
      }
    }
  }

  async function getJobLevels() {
    let token = sessionStorage.getItem("token");
    if (token) {
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
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  $: buttonDisable = jobLevels.length < 1;

  onMount(async () => {
    getJobLevels();
  });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- Header -->
      <div
        class="flex items-center pb-8 border-b border-secondary-light dark:border-secondary-dark"
      >
        <div class="flex flex-col gap-2 grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {$_("jobLevel")}
            </Heading>
          </div>
        </div>
        <div class="flex gap-2">
          <Button
            class="px-3 focus-within:!ring-1"
            disabled={buttonDisable}
            on:click={saveJobLevels}>{$_("save")}</Button
          >
          <Button class="px-3 focus-within:!ring-1">{$_("importCsv")}</Button>
        </div>
      </div>

      <!-- Submenu -->
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
                on:consider={handleDndConsider}
                on:finalize={handleDndFinalize}
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
          {:else}
            <div class="pt-8">
              <StateCard
                illustration={IllustrationEmpty}
                description={$_("youDontHaveAnyDataYet")}
              />
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
  <ToastContainer />
</div>
