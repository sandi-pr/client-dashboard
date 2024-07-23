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
  import { Typography, StateCard, Dropdown } from "$lib";
  import { Button, Modal, Input, Textarea } from "flowbite-svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import CalendarPlus02 from "$lib/icon/calendar-plus-02.svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import { _ } from "svelte-i18n";
  import { format } from "date-fns";
  import { DatePicker } from "@svelte-plugins/datepicker";

  let tmpProjectName = $state("");
  let tmpProjectDescription = $state("");
  let selectedIndustry = $state("Oil & Gas");
  let selectedRole = $state("Buyer");
  let dropdownIndustryOpen = $state(false);
  let dropdownRole = $state(false);
  let addProjectModal = $state(false);
  let isOpen = $state(false);
  let startDate = $state(null);
  let endDate = $state(null);
  let dateFormat = $state("d MMM yyyy");
  let formattedStartDate = $state("");
  let formattedEndDate = $state("");

  let industries = $state([
    { name: "Oil & Gas", value: "oil-gas" },
    { name: "Other", value: "other" },
  ]);

  let role = $state([
    { name: "Buyer", value: "buyer" },
    { name: "Seller", value: "seller" },
  ]);

  let configs = $state({
    search: true,
    pagination: { page: 1, pageSize: 10, pageCount: 1, total: 1 },
  });

  let datas = $state([
    {
      project_name: "Aramco - Gasoline 90 RON",
      industry: "Oil & Gas",
      start_date: "01/12/2024",
      end_date: "22/05/2025",
      actions: {
        objects: [
          {
            href: `contract_reviewer/project` + `?projectId=1`,
            name: "View",
          },
        ],
      },
    }
  ]);

  function next() {
    console.log("next");
  }

  function previous() {
    console.log("previous");
  }

  function createProject() {
    addProjectModal = false;
  }

  const toggleDatePicker = () => {
    isOpen = !isOpen;
  };

  const formatDate = (dateString) => (dateString && format(new Date(dateString), dateFormat)) || "";

  $effect(() => {
    formattedStartDate = formatDate(startDate);
    formattedEndDate = formatDate(endDate);
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
          <Button onclick={() => (addProjectModal = true)}>
            <img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" />
            {$_("addNew")}
          </Button>
        </div>

        <!-- Document Table Data -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs bind:next bind:previous />
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
      <Input bind:value={tmpProjectName} placeholder={$_("whatIsYourProjectTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("industryType")}*`} />
      <Dropdown dropdownOpen={dropdownIndustryOpen} selectedItem={selectedIndustry} items={industries} placeholder={$_("selectIndustry")} />
    </div>
    <div class="flex gap-4">
      <div class="flex flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={$_("as")} />
        <Dropdown width="w-min" dropdownOpen={dropdownRole} selectedItem={selectedRole} items={role} />
      </div>
      <div class="flex grow flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={`${$_("projectDuration")}*`} />
        <DatePicker
          bind:isOpen
          bind:startDate
          bind:endDate
          isRange
          isMultipane
          enableFutureDates
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