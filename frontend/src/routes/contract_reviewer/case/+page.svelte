<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer:
  - Dean Ramhan <dean@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
  - Taqie Fadlillah <taqief@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->
<script>
  import { onMount } from "svelte";
  import { Modal, Button, Input, Pagination } from "flowbite-svelte";
  import { Typography, triggerToast, StateCard, triggerLoading } from "$lib";
  import * as Icon from "flowbite-svelte-icons";
  import CalendarPlus02 from "$lib/icon/calendar-plus-02.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { _ } from "svelte-i18n";
  import { CR3Call } from "hexacomponents-svelte";
  import { format } from "date-fns";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import { goto } from "$app/navigation";

  import { call } from "$lib";

  let tmpRulesName = $state("");
  let addRulesModal = $state(false);

  let configs = $state({});

  let datas = $state([]);

  let page = $state(1);
  let pageCount = $state(1);
  let searchTerm = $state("");

  async function createRules() {
    addRulesModal = false;
    triggerLoading(true);
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "create-contract-rules", {
          URLQuery: {},
          URLParams: {},
          Params: {
            rule_name: tmpRulesName,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          addRulesModal = true;
          triggerToast("error");
          console.log(response);
          return;
        }

        triggerLoading(false);
        triggerToast("created");
        tmpRulesName = "";
        addRulesModal = false;
        getRules();
        // goto(
        //   (await $call.getPage("client-dashboard", "contract-case-detail")) +
        //     "?caseId=" +
        //     response.result.data.id,
        // );
      } catch (error) {
        triggerLoading(false);
        addRulesModal = true;
        triggerToast("error");
        console.log(error);
      }
    }
  }

  const handleSearch = async (value) => {
    page = 1;
    searchTerm = value;

    await getRules();
  };

  async function getRules() {
    triggerLoading(true);
    let token = sessionStorage.getItem("token");
    if (token) {
      let rules = [];
      try {
        const response = await $call.call("client-dashboard", "get-contract-rules", {
          URLQuery: {},
          URLParams: {},
          Params: {
            search: searchTerm,
            page: page,
            page_size: 10,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          console.log(response);
          return;
        }

        rules = response.result.data;

        if (rules.length > 0) {
          pageCount = response.result.meta.pagination.pageCount;
          configs = {
            pagination: response.result.meta.pagination,
            search: true,
            searchHandler: handleSearch,
          };
          let link = await $call.getPage("client-dashboard", "contract-case-detail");
          let promise = rules.map((rule) => {
            return {
              id: rule.id,
              case_name: rule.attributes.rule_name,
              actions: {
                objects: [
                  {
                    href: link + "?caseId=" + rule.id,
                    name: "View",
                  },
                ],
              },
            };
          });
          datas = await Promise.all(promise);
          triggerLoading(false);
        }
        triggerLoading(false);
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function nextPage() {
    if (page >= pageCount) return;
    page++;
    await getRules();
  }

  async function prevPage() {
    if (page <= 1) return;
    page--;
    await getRules();
  }

  onMount(async () => {
    getRules();
  });

  let buttonDisabled = $state(false);

  $effect(() => {
    buttonDisabled = tmpRulesName.length === 0;
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex justify-between">
          <Typography preset="Display sm/Semibold" text={$_("caseList")} />
          <Button
            onclick={() => (addRulesModal = true)}
            class="flex gap-2 items-center"
            color="primary"
            variant="contained"
            size="sm"
          >
            <img class="w-4 h-4" src={PlusIcon} alt="plus-icon" />
            {$_("addNew")}
          </Button>
        </div>

        <!-- table latest -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs next={nextPage} prev={prevPage} />
        {:else}
          <div class="pt-8">
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<Modal title={$_("createNewCases")} bind:open={addRulesModal}>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("caseName")}*`} />
      <Input bind:value={tmpRulesName} placeholder={$_("whatIsYourCaseTitle")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={createRules} class="w-full" disabled={buttonDisabled}
      >{$_("addNewCase")}</Button
    >
  </svelte:fragment>
</Modal>
