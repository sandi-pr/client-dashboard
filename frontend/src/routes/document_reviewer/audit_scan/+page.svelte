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
  import { _ } from "svelte-i18n";
  import { format } from "date-fns";
  import { DatePicker } from "@svelte-plugins/datepicker";
  import { Button, Modal, Input, Textarea } from "flowbite-svelte";
  import { triggerToast, triggerLoading, call, check, getIcons } from "$lib";
  import { Typography, StateCard, Dropdown, DataTable } from "$lib";

  const IllustrationEmpty = getIcons("illustrationEmpty");
  const IllustrationProgress = getIcons("illustrationProgress");
  const CalendarPlus02 = getIcons("calendarPlus02");
  const PlusIcon = getIcons("plus");

  const isFrontend = process.env.NODE_ENV === "development";

  let isOpen = $state(false);
  let gettingDatas = $state(false);
  let addAuditScanModal = $state(false);
  let dropdownRulesOpen = $state(false);

  let page = $state(1);
  let pageCount = $state(1);

  let selectedRules = $state("");
  let tmpAuditScanName = $state("");
  let formattedEndDate = $state("");
  let selectedRole = $state("Buyer");
  let formattedStartDate = $state("");
  let dateFormat = $state("d MMM yyyy");
  let tmpProjectDescription = $state("");

  let configs = $state({});

  let datas = $state([]);
  let rules = $state([]);
  let rulesTmp = $state([]);

  let endDate = $state(null);
  let startDate = $state(null);

  const getAuditScan = async () => {
    try {
      triggerLoading(true);
      gettingDatas = true;

      let token = localStorage.getItem("token");
      const response = await $call.call("client-dashboard", "get-contract-projects", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: page,
          page_size: 10,
          is_compare: false,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        let error = new Error("Failed to fetch audit scan projects");
        error.response = response;
        throw error;
      }

      if (response.status === 200) {
        pageCount = response.result.meta.pagination.pageCount;
        configs = {
          pagination: response.result.meta.pagination,
        };

        let url = isFrontend
          ? "/document_reviewer/audit_scan/detail"
          : await $call.getPage("client-dashboard", "contract-audit-scan-detail");

        datas = response.result.data.map((project) => {
          return {
            id: project.id,
            project_name: project.attributes.project_name,
            start_date: format(new Date(project.attributes.duration_start), "dd/MM/yyyy"),
            end_date: format(new Date(project.attributes.duration_end), "dd/MM/yyyy"),
            actions: {
              objects: [
                check("aiManagement.contractreviewer.project:view") && {
                  href: url + `?projectId=${project.id}`,
                  name: "View",
                },
              ].filter(Boolean),
            },
          };
        });
      }
    } catch (error) {
      triggerToast("error");
      console.log(error);
    } finally {
      triggerLoading(false);
      gettingDatas = false;
    }
  };

  const getRules = async () => {
    try {
      triggerLoading(true);
      gettingDatas = true;

      let token = localStorage.getItem("token");
      rulesTmp = [];
      let page = 1;
      let error = [];

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
          error.push(response);
        }

        if (response.result.data.length === 0) {
          break;
        }

        rulesTmp = rulesTmp.concat(response.result.data);
        page++;
      }

      if (rulesTmp.length > 0) {
        rules = rulesTmp.map((rule) => {
          return {
            value: rule.id,
            name: rule.attributes.rule_name,
          };
        });
      }

      if (error.length > 0) {
        console.log(error);
      }
    } catch (error) {
      console.log(error);
    }
  };

  const createAuditScan = async () => {
    addAuditScanModal = false;
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
            project_name: tmpAuditScanName,
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
          addAuditScanModal = true;
        }

        triggerLoading(false);
        tmpAuditScanName = "";
        tmpProjectDescription = "";
        selectedRole = "Buyer";
        startDate = null;
        endDate = null;
        triggerToast("created");
        getAuditScan();
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        addAuditScanModal = false;
        console.log(error);
      }
    }
  };

  const nextPage = () => {
    if (page >= pageCount) return;
    page++;
    getAuditScan();
  };

  const previousPage = () => {
    if (page <= 1) return;
    page--;
    getAuditScan();
  };

  const toggleDatePicker = () => {
    isOpen = !isOpen;
  };

  const formatDate = (dateString) => (dateString && format(new Date(dateString), dateFormat)) || "";

  $effect(() => {
    formattedStartDate = formatDate(startDate);
    formattedEndDate = formatDate(endDate);
  });

  onMount(async () => {
    await getRules();
    await getAuditScan();
  });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text={$_("auditScanList")} />
          {#if check("aiManagement.contractreviewer.project:create")}
            <Button onclick={() => (addAuditScanModal = true)}>
              <PlusIcon svgClass="w-4 h-4 me-1" />
              {$_("addNew")}
            </Button>
          {/if}
        </div>

        <!-- Document Table Data -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs next={nextPage} previous={previousPage} />
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

<Modal title={$_("createNewProject")} bind:open={addAuditScanModal} size="lg">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("projectName")}*`} />
      <Input bind:value={tmpAuditScanName} placeholder={$_("whatIsYourProjectTitle")} />
    </div>
    <div class="flex gap-4">
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
      <div class="flex grow flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={`${$_("projectDuration")}*`} />
        <DatePicker bind:isOpen bind:startDate bind:endDate isRange isMultipane enableFutureDates>
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
    <Button onclick={createAuditScan} disabled={!tmpAuditScanName || !endDate} class="w-full"
      >{$_("addNewProject")}</Button
    >
  </svelte:fragment>
</Modal>
