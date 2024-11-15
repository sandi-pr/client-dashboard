<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: 
  - Dean Ramhan <dean@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->
<script>
  import { Typography, StateCard, Dropdown, triggerToast, triggerLoading, call, check } from "$lib";
  import { Button, Modal, Input, Textarea } from "flowbite-svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import CalendarPlus02 from "$lib/icon/calendar-plus-02.svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import { _ } from "svelte-i18n";
  import { format } from "date-fns";
  import { DatePicker } from "@svelte-plugins/datepicker";
  import { onMount } from "svelte";

  let tmpProjectName = $state("");
  let tmpProjectDescription = $state("");
  let selectedRules = $state("");
  let selectedRole = $state("Buyer");
  let dropdownIndustryOpen = $state(false);
  let dropdownRulesOpen = $state(false);
  let dropdownRole = $state(false);
  let addProjectModal = $state(false);
  let isOpen = $state(false);
  let startDate = $state(null);
  let endDate = $state(null);
  let dateFormat = $state("d MMM yyyy");
  let formattedStartDate = $state("");
  let formattedEndDate = $state("");
  let rulesTmp = $state([]);
  let page = $state(1);
  let pageCount = $state(1);

  let industries = $state([]);
  let rules = $state([]);

  let role = $state([
    { name: "Buyer", value: "Buyer" },
    { name: "Seller", value: "Seller" },
  ]);

  let configs = $state({});

  let datas = $state([]);

  let gettingDatas = $state(false);

  async function getProjects() {
    gettingDatas = true;
    // triggerLoading(true);
    await getRules();

    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "get-contract-projects", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 10,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          gettingDatas = false;
          console.log(response);
          return;
        }

        if (response.status === 200) {
          pageCount = response.result.meta.pagination.pageCount;
          configs = {
            pagination: response.result.meta.pagination,
          };
          let url = await $call.getPage("client-dashboard", "contract-project-detail");
          let promise = response.result.data.map((project) => {
            return {
              id: project.id,
              project_name: project.attributes.project_name,
              start_date: format(new Date(project.attributes.duration_start), "dd/MM/yyyy"),
              end_date: format(new Date(project.attributes.duration_end), "dd/MM/yyyy"),
              actions: {
                objects: [
                  check("aiManagement.contractreviewer.project:view") && {
                    href: url + `?projectId=${project.id}&title=${project.attributes.project_name}`,
                    name: "View",
                  },
                ].filter(Boolean),
              },
            };
          });

          datas = await Promise.all(promise);
          console.log("datas", datas);

          triggerLoading(false);
          gettingDatas = false;
        }
      } catch (error) {
        triggerLoading(false);
        gettingDatas = false;
        triggerToast("error");
        console.log(error);
      }
    }
    gettingDatas = false;
  }

  async function getRules() {
    let token = localStorage.getItem("token");
    if (token) {
      rulesTmp = [];
      let page = 1;
      try {
        while (true) {
          const response = await $call.call("client-dashboard", "get-contract-rules", {
            URLQuery: {},
            URLParams: {},
            Params: {
              page: page,
              page_size: 25,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response.status !== 200) {
            console.log(response);
            return;
          }

          if (response.result.data.length === 0) {
            break;
          }
          rulesTmp = rulesTmp.concat(response.result.data);
          page++;
        }

        if (rulesTmp.length > 0) {
          let promise = rulesTmp.map((rule) => {
            return {
              value: rule.id,
              name: rule.attributes.rule_name,
            };
          });
          rules = await Promise.all(promise);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  function nextPage() {
    if (page >= pageCount) return;
    page++;
    getProjects();
  }

  function previousPage() {
    if (page <= 1) return;
    page--;
    getProjects();
  }

  const toggleDatePicker = () => {
    isOpen = !isOpen;
  };

  const formatDate = (dateString) => (dateString && format(new Date(dateString), dateFormat)) || "";

  $effect(() => {
    formattedStartDate = formatDate(startDate);
    formattedEndDate = formatDate(endDate);
  });

  async function createProject() {
    addProjectModal = false;
    triggerLoading(true);
    let token = localStorage.getItem("token");
    let industry_id = sessionStorage.getItem("industry_id");
    if (industry_id === null) {
      triggerLoading(false);

      return;
    }
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "create-contract-project", {
          URLQuery: {},
          URLParams: {},
          Params: {
            project_name: tmpProjectName,
            industry_id: parseInt(industry_id),
            contract_rule_id: selectedRules,
            project_as: selectedRole,
            duration_start: format(new Date(startDate), "yyyy-MM-dd"),
            duration_end: format(new Date(endDate), "yyyy-MM-dd"),
            description: tmpProjectDescription,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          addProjectModal = true;
        }

        triggerLoading(false);
        tmpProjectName = "";
        tmpProjectDescription = "";
        selectedRole = "Buyer";
        startDate = null;
        endDate = null;
        triggerToast("created");
        getProjects();
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        addProjectModal = false;
        console.log(error);
      }
    }
  }

  onMount(async () => {
    getProjects();
  });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text={$_("projectList")} />
          {#if check("aiManagement.contractreviewer.project:create")}
            <Button onclick={() => (addProjectModal = true)}>
              <img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" />
              {$_("addNew")}
            </Button>
          {/if}
        </div>

        <!-- Document Table Data -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs bind:next={nextPage} bind:previous={previousPage} />
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

<Modal title={$_("createNewProject")} bind:open={addProjectModal} size="lg">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("projectName")}*`} />
      <Input bind:value={tmpProjectName} placeholder={$_("whatIsYourProjectTitle")} />
    </div>
    <div class="flex gap-4">
      <!-- <div class="flex flex-col grow gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={`${$_("industryType")}*`} />
        <Dropdown
          dropdownOpen={dropdownIndustryOpen}
          bind:selectedItem={selectedIndustry}
          items={industries}
          placeholder={$_("selectIndustry")}
          isOverlay={true}
        />
      </div> -->
      <div class="flex flex-col grow gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={`${$_("case")}*`} />
        <Dropdown
          dropdownOpen={dropdownRulesOpen}
          bind:selectedItem={selectedRules}
          items={rules}
          placeholder={$_("chooseCase")}
          isOverlay={true}
        />
      </div>
    </div>
    <div class="flex gap-4">
      <!-- <div class="flex flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={$_("as")} />
        <Dropdown
          width="w-min"
          dropdownOpen={dropdownRole}
          bind:selectedItem={selectedRole}
          items={role}
          isOverlay={true}
        />
      </div> -->
      <div class="flex grow flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={`${$_("projectDuration")}*`} />
        <DatePicker bind:isOpen bind:startDate bind:endDate isRange isMultipane enableFutureDates>
          <button
            onclick={toggleDatePicker}
            class:open={isOpen}
            class="flex w-full p-2.5 border rounded-lg border-gray-300 dark:border-gray-500 gap-2 justify-center items-center"
          >
            <svelte:component
              this={CalendarPlus02}
              class="w-3 h-3 ms-2 text-gray-700 dark:text-gray-400"
            />
            <div class="date">
              {#if startDate != endDate}
                {formattedStartDate} - {formattedEndDate}
              {:else if startDate}
                {formattedStartDate}
              {:else}
                {$_("pickADate")}
              {/if}
            </div>
          </button>
        </DatePicker>
      </div>
    </div>

    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Textarea bind:value={tmpProjectDescription} placeholder={$_("placholderDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={createProject} disabled={!tmpProjectName || !endDate} class="w-full"
      >{$_("addNewProject")}</Button
    >
  </svelte:fragment>
</Modal>
