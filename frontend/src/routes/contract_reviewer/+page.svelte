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
  import { Modal, Button, Input, Textarea, Heading } from "flowbite-svelte";
  import { DatePicker } from "@svelte-plugins/datepicker";
  import { _ } from "svelte-i18n";
  import { format } from "date-fns";
  import {
    call,
    check,
    getIcons,
    getAllRules,
    getProjects,
    triggerToast,
    triggerLoading,
  } from "$lib";
  import { Dropdown, DataTable, StateCard, Typography, ActionCard } from "$lib";

  const IllustrationEmpty = getIcons("illustrationEmpty");
  const IllustrationProgress = getIcons("illustrationProgress");
  const CalendarPlus02 = getIcons("calendarPlus02");

  const isFrontend = process.env.NODE_ENV === "development";

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

  let datasConfigs = $state({});
  let rulesConfigs = $state({});

  let datasTable = $state([]);
  let rules = $state([]);
  let rulesTable = $state([]);

  let gettingDatas = $state(false);
  let gettingRules = $state(false);

  let buttonDisabled = $state(false);

  const getActionItems = async () => {
    actionItems = [
      {
        disabled: !check("aiManagement.contractreviewer.case:view"),
        href: isFrontend
          ? "/contract_reviewer/case"
          : await $call.getPage("client-dashboard", "contract-case"),
        icon: "bookOpen",
        actionName: "Manage cases",
        desc: "desc Set your referenced document",
      },
      {
        disabled: !check("aiManagement.contractreviewer.project:create"),
        onclick: () => (addProjectModal = true),
        icon: "packagePlus",
        actionName: "Add new project",
        desc: "desc Add new project",
      },
    ];
  };

  const getCompanyType = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "get-company-type", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status !== 200) {
          const error = new Error("An error occurred while fetching the data.");
          error.response = response;
          throw error;
        }

        // console.log("response", response);
        return response?.result?.data?.attributes?.company_industry || 0;
      } catch (error) {
        console.log("error", error);
        return 0;
      }
    }
  };

  const createProject = async () => {
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
        handleGetProjects();
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        addProjectModal = true;
        console.log(error);
      }
    }
  };

  const handleGetRules = async () => {
    gettingRules = true;
    const allrules = await getAllRules($call);
    if (allrules.length > 0) {
      let link = isFrontend
        ? "/contract_reviewer/case/detail"
        : await $call.getPage("client-dashboard", "contract-case-detail");

      rules = allrules.map((rule) => ({
        value: rule.id,
        name: rule.attributes.rule_name,
      }));

      rulesConfigs = {
        headers: "casesList",
        actions: [
          {
            name: "viewAll",
            href: isFrontend
              ? "/contract_reviewer/case"
              : await $call.getPage("client-dashboard", "contract-case"),
            type: "info",
          },
        ],
        seperator: true,
      };

      const promise = allrules.map(async (rule) => ({
        id: rule.id,
        rule_name: rule.attributes.rule_name,
        actions: {
          objects: [
            check("businessManagement.hr.recruitment:view") && {
              href: link + "?caseId=" + rule.id,
              name: "View",
              type: "primary",
            },
          ].filter(Boolean),
        },
      }));
      rulesTable = await Promise.all(promise);
    } else {
      rulesTable = [];
    }
    gettingRules = false;
  };

  const handleGetProjects = async () => {
    gettingDatas = true;

    const projects = await getProjects($call, { page: 1, pageSize: 5 });
    if (projects.length > 0) {
      datasConfigs = {
        headers: "projectList",
        actions: [
          {
            name: "viewAll",
            href: isFrontend
              ? "/contract_reviewer/project"
              : await $call.getPage("client-dashboard", "contract-project-list"),
            type: "info",
          },
        ],
        seperator: true,
      };

      let url = isFrontend
        ? "/contract_reviewer/project/detail"
        : await $call.getPage("client-dashboard", "contract-project-detail");

      const promise = projects.map((project) => {
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
                type: "primary",
              },
            ],
          },
        };
      });

      datasTable = await Promise.all(promise);
    } else {
      datasTable = [];
    }

    gettingDatas = false;
  };

  const toggleDatePicker = () => {
    isOpen = !isOpen;
  };

  const formatDate = (dateString) => (dateString && format(new Date(dateString), dateFormat)) || "";

  onMount(async () => {
    await getActionItems();
    await handleGetRules();
    await handleGetProjects();
    await getCompanyType();
  });

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
        <Typography
          preset="Display sm/Semibold"
          text={$_("contractReviewerOverview")}
          propClass="fluid-xl"
        />

        <!-- Action Items -->
        <div>
          {#if actionItems.some((item) => item.disabled !== true)}
            <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
              <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
                {$_("actions")}
              </Heading>
            </div>
          {/if}
          <ActionCard {actionItems} />
        </div>

        {#if datasTable.length > 0}
          <DataTable bind:datas={datasTable} bind:configs={datasConfigs} />
        {:else if gettingDatas}
          <StateCard
            illustration={IllustrationProgress}
            description={$_("fetchingDataPleaseWait")}
          />
        {:else}
          <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
        {/if}

        {#if rulesTable.length > 0}
          <DataTable bind:datas={rulesTable} bind:configs={rulesConfigs} />
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

<Modal title={$_("createNewProject")} bind:open={addProjectModal} size="lg" outsideclose>
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
          propClass="w-full disabled:cursor-not-allowed disabled:opacity-50 rtl:text-right p-2.5 text-gray-900 dark:text-white dark:placeholder-gray-400 border-gray-300 dark:border-gray-500 text-sm rounded-lg bg-primary-light dark:bg-gray-800"
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
            <CalendarPlus02 svgClass="w-3 h-3 ms-2 text-gray-700 dark:text-gray-400" />
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
