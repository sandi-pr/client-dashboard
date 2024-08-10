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
  import { onMount } from "svelte";
  import { Modal, Button, Input, Textarea } from "flowbite-svelte";
  import {
    Typography,
    SectionHeader,
    Card,
    Dropdown,
    triggerToast,
    StateCard,
    triggerLoading,
  } from "$lib";
  import CalendarPlus02 from "$lib/icon/calendar-plus-02.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { DatePicker } from "@svelte-plugins/datepicker";
  import { _ } from "svelte-i18n";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { goto } from "$app/navigation";
  import { CR3Call } from "hexacomponents-svelte";
  import { format } from "date-fns";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { tr } from "date-fns/locale";

  const call = new CR3Call(globalThis.location.origin);

  let files = $state([]);
  let actionItems = $state([]);
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

  let industries = $state([]);

  let role = $state([
    { name: "Buyer", value: "Buyer" },
    { name: "Seller", value: "Seller" },
  ]);

  let configs = $state({});

  let datas = $state([]);
  let rules = $state([]);

  async function getActionItems() {
    actionItems = [
      {
        href: "#",
        icon: "packagePlus",
        actionName: "Add new project",
        desc: "desc Add new project",
      },
      {
        href: await call.getPage("client-dashboard", "contract-case"),
        icon: "bookOpen",
        actionName: "Manage cases",
        desc: "desc Set your referenced document",
      },
    ];
  }

  async function getCompanyType() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "get-company-type", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return 0;
        }

        return parseInt(response.result.data[0].attributes.company_industry);
      } catch (error) {
        console.log(error);
        return 0;
      }
    }
  }

  async function createProject() {
    addProjectModal = false;
    triggerLoading(true);
    let token = sessionStorage.getItem("token");
    let industry_id = sessionStorage.getItem("industry_id");
    if (industry_id === null) {
      triggerLoading(false);
      alert("Please select an industry first");
      return;
    }
    if (token) {
      try {
        const response = await call.call("client-dashboard", "create-contract-project", {
          URLQuery: {},
          URLParams: {},
          Params: {
            project_name: tmpProjectName,
            industry_id: parseInt(industry_id),
            project_as: selectedRole,
            contract_rule_id: selectedRules,
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
        getProjectsLatest();
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        addProjectModal = true;
        console.log(error);
      }
    }
  }

  async function getRules() {
    let token = sessionStorage.getItem("token");
    if (token) {
      rulesTmp = [];
      let page = 1;
      try {
        while (true) {
          const response = await call.call("client-dashboard", "get-contract-rules", {
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

  const toggleDatePicker = () => {
    isOpen = !isOpen;
  };

  async function getProjectsLatest() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "get-contract-projects", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: 1,
            page_size: 5,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        console.log("response", response);
        if (response.status === 200) {
          configs = {
            headers: "projectList",
            actions: [
              {
                name: "viewAll",
                href: await call.getPage("client-dashboard", "contract-project-list"),
                type: "info",
              },
            ],
            seperator: true,
          };
          let url = await call.getPage("client-dashboard", "contract-project-detail");
          let promise = response.result.data.map((project) => {
            return {
              id: project.id,
              project_name: project.attributes.project_name,
              start_date: format(new Date(project.attributes.duration_start), "dd/MM/yyyy"),
              end_date: format(new Date(project.attributes.duration_end), "dd/MM/yyyy"),
              actions: {
                objects: [
                  {
                    href: url + `?projectId=${project.id}`,
                    name: "View",
                  },
                ],
              },
            };
          });

          datas = await Promise.all(promise);
          console.log("datas", datas);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  onMount(async () => {
    getActionItems();
    getProjectsLatest();
    await getRules();
  });

  let buttonDisabled = $state(false);

  const formatDate = (dateString) => (dateString && format(new Date(dateString), dateFormat)) || "";

  $effect(() => {
    formattedStartDate = formatDate(startDate);
    formattedEndDate = formatDate(endDate);
  });

  $effect(() => {
    buttonDisabled = files.length === 0;
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Typography preset="Display sm/Semibold" text={$_("contractReviewerOverview")} />

        <!-- Action Items -->
        <div class="flex flex-col gap-4">
          <SectionHeader title={$_("action")} />
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 w-full items-center justify-center">
            {#each actionItems as item}
              <Card
                onclick={() => {
                  if (item.href === "#") addProjectModal = true;
                  else goto(item.href);
                }}
                icon={item.icon}
                title={$_(toCamelCase(item.actionName))}
                description={$_(toCamelCase(item.desc))}
              />
            {/each}
          </div>
        </div>

        <!-- table latest -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs />
        {:else}
          <div class="pt-8">
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<Modal title={$_("createNewProject")} bind:open={addProjectModal} size="lg">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("projectName")}*`} />
      <Input
        bind:value={tmpProjectName}
        placeholder={$_("whatIsYourProjectTitle")}
        class="bg-primary-light dark:bg-gray-800"
      />
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
        <DatePicker
          bind:isOpen
          bind:startDate
          bind:endDate
          isRange
          isMultipane
          enableFutureDates
          theme="custom-datepicker"
        >
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

<!-- <LoadingContainer /> -->

<!-- <style>
  :global(.datepicker[data-picker-theme="custom-datepicker"]) {
    --datepicker-container-background: #fff;
  }
</style> -->
