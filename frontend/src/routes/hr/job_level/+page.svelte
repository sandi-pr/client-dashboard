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
  import {
    P,
    Heading,
    Input,
    Button,
    Select,
    Textarea,
    Toggle,
    Toast,
    Helper,
  } from "flowbite-svelte";
  import { CheckCircleSolid, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { fly } from "svelte/transition";
  import { _ } from "svelte-i18n";
  import { onMount } from "svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { StateCard } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import DotsGridIcon from "$lib/icon/dots-grid.svelte";
  import { goto } from "$app/navigation";
  import { dndzone } from "svelte-dnd-action"

  const call = new CR3Call(globalThis.location.origin);

  let tmpJobLevel = $state("");

  let toastCreated = $state(false);
  let toastResourceEmpty = $state(false);
  let toastError = $state(false);
  let counter = $state(6);

  let jobLevels = $state([]);

  function addJobLevel() {
    const newJobLevel = {
      // dummy id
      id: tmpJobLevel,
      title: tmpJobLevel
    };
    jobLevels = [...jobLevels, newJobLevel];

    tmpJobLevel = "";
  }

  function deleteJobLevel(id) {
    jobLevels = jobLevels.filter((jobLevel) => jobLevel.id !== id);
  }

  function triggerToast(target) {
    if (target === "created") {
      toastCreated = true;
      counter = 6;
      timeout();
    } else if (target === "resource_empty") {
      toastResourceEmpty = true;
      counter = 6;
      timeout();
    } else if (target === "error") {
      toastError = true;
      counter = 6;
      timeout();
    } else {
      console.log("target not found");
    }
  }

  function timeout() {
    if (--counter > 0) return setTimeout(timeout, 1000);
    toastCreated = false;
    toastError = false;
  }

  let buttonDisable = $state(true);

  onMount(async () => {
    // getCategories();
  });

  // let dropdownOpen = false;
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- Header -->
      <div class="flex items-center pb-8 border-b border-secondary-light dark:border-secondary-dark">
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
          <Button class="px-3 focus-within:!ring-1"
            >{$_("importCsv")}</Button
          >
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
          <Input
            bind:value={tmpJobLevel}
            class="bg-primary-light dark:bg-secondary-dark"
          />
          <Button color="light" on:click={addJobLevel}>{$_("add")}</Button>
        </div>
        <div class="flex flex-col grow gap-6">
          
            {#if jobLevels.length > 0}
              <section use:dndzone="{{ items: jobLevels }}">
                <div class="bg-secondary-light dark:bg-secondary-dark rounded-lg py-10 pl-4">
                  {#each jobLevels as job (job.id)}
                    <div class="px-4 py-2 bg-white rounded-lg backdrop-blur-[26px] justify-between items-center flex w-max">
                      <div class="min-w-[172px] self-stretch items-center gap-4 flex">
                        <!-- <div class="w-6 h-6 relative"></div> -->
                        <DotsGridIcon class="w-6 h-6 text-fg-quaternary-light dark:text-fg-quaternary-dark" />
                        <P size="lg" weight="semibold">{job.title}</P>
                      </div>
                      <Button
                        on:click={() => deleteJobLevel(job.id)}
                        color="light"
                        class="w-max">{$_("delete")}</Button
                      >
                    </div>
                  {/each}
                </div>
              </section>
            {:else}
              <div class="pt-8">
                <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
              </div>
            {/if}
        </div>
      </div>
    </div>
  </div>
  <div class="fixed bottom-5 end-5">
    <Toast color="green" transition={fly} params={{ x: 200 }} bind:open={toastCreated}>
      <svelte:fragment slot="icon">
        <CheckCircleSolid class="w-5 h-5" />
      </svelte:fragment>
      {$_("successfullyCreated")}
    </Toast>
    <Toast color="red" transition={fly} params={{ x: 200 }} bind:open={toastResourceEmpty}>
      <svelte:fragment slot="icon">
        <ExclamationCircleSolid class="w-5 h-5" />
      </svelte:fragment>
      {$_("resourceNameEmpty")}
    </Toast>
    <Toast
      class="!max-w-sm"
      color="red"
      transition={fly}
      params={{ x: 200 }}
      bind:open={toastError}
    >
      <svelte:fragment slot="icon">
        <ExclamationCircleSolid class="w-5 h-5" />
      </svelte:fragment>
      {$_("sometingWentWrongPleaseTryAgain")}
    </Toast>
  </div>
</div>
