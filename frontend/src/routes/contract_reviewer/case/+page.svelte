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
  import { _ } from "svelte-i18n";
  import { Modal, Button, Input } from "flowbite-svelte";
  import { triggerToast, triggerLoading, call, check, getIcons } from "$lib";
  import { Typography, StateCard, DataTable } from "$lib";

  const IllustrationEmpty = getIcons("illustrationEmpty");
  const IllustrationProgress = getIcons("illustrationProgress");
  const PlusIcon = getIcons("plus");

  const isFrontend = process.env.NODE_ENV === "development";

  let gettingDatas = $state(false);
  let addRulesModal = $state(false);
  let buttonDisabled = $state(false);

  let page = $state(1);
  let pageCount = $state(1);

  let searchTerm = $state("");
  let tmpRulesName = $state("");

  let configs = $state({});

  let datas = $state([]);

  const getRules = async () => {
    // triggerLoading(true);
    let token = localStorage.getItem("token");
    if (token) {
      gettingDatas = true;
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
          gettingDatas = false;
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
          let link = isFrontend
            ? "/contract_reviewer/case/detail"
            : await $call.getPage("client-dashboard", "contract-case-detail");
          let promise = rules.map((rule) => {
            return {
              id: rule.id,
              case_name: rule.attributes.rule_name,
              actions: {
                objects: [
                  check("businessManagement.hr.recruitment:view") && {
                    href: link + "?caseId=" + rule.id,
                    name: "View",
                  },
                ].filter(Boolean),
              },
            };
          });
          datas = await Promise.all(promise);
          triggerLoading(false);
          gettingDatas = false;
        }
        triggerLoading(false);
        gettingDatas = false;
      } catch (error) {
        console.log(error);
        gettingDatas = false;
      }
    }
  };

  const createRules = async () => {
    addRulesModal = false;
    triggerLoading(true);
    let token = localStorage.getItem("token");
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
  };

  const handleSearch = async (value) => {
    page = 1;
    searchTerm = value;

    await getRules();
  };

  const nextPage = async () => {
    if (page >= pageCount) return;
    page++;
    await getRules();
  };

  const prevPage = async () => {
    if (page <= 1) return;
    page--;
    await getRules();
  };

  onMount(async () => {
    getRules();
  });

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
          {#if check("aiManagement.contractreviewer.case:create")}
            <Button
              onclick={() => (addRulesModal = true)}
              class="flex gap-2 items-center"
              color="primary"
              variant="contained"
              size="sm"
            >
              <PlusIcon svgClass="w-4 h-4" />
              {$_("addNew")}
            </Button>
          {/if}
        </div>

        <!-- table latest -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs next={nextPage} prev={prevPage} />
        {:else if gettingDatas}
          <StateCard
            illustration={IllustrationProgress}
            description={$_("fetchingDataPleaseWait")}
          />
        {:else}
          <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
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
