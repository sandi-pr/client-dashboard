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
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

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
  let pagination = $state({
    page: 1,
    pageSize: 10,
    pageCount: 1,
    total: 1,
  });

  let industries = $state([
    { name: "Oil & Gas", value: "oil-gas" },
    { name: "Other", value: "other" },
  ]);

  let role = $state([
    { name: "Buyer", value: "buyer" },
    { name: "Seller", value: "seller" },
  ]);

  let configs = $state({});

  let datas = $state([]);

  async function getProjects() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "get-contract-projects", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: pagination.page,
            page_size: pagination.pageSize,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        console.log("response", response);
        if (response.status === 200) {
          configs = {
            search: true,
            pagination: response.result.meta.pagination,
          };
          let url = await call.getPage("client-dashboard", "contract-project-detail");
          let promise = response.result.data.map((project) => {
            return {
              id: project.id,
              project_name: project.attributes.project_name,
              industry: project.attributes.industry_type,
              start_date: format(new Date(project.attributes.project_start_date), "dd/MM/yyyy"),
              end_date: format(new Date(project.attributes.project_end_date), "dd/MM/yyyy"),
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

  function next() {
    if (pagination.page < pagination.pageCount) {
      pagination.page = pagination.page + 1;
      getProjects();
    }
  }

  function previous() {
    if (pagination.page > 1) {
      pagination.page = pagination.page - 1;
      getProjects();
    }
  }

  const toggleDatePicker = () => {
    isOpen = !isOpen;
  };

  const formatDate = (dateString) => (dateString && format(new Date(dateString), dateFormat)) || "";

  $effect(() => {
    formattedStartDate = formatDate(startDate);
    formattedEndDate = formatDate(endDate);
  });

  function createProject() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = call.call("client-dashboard", "create-contract-project", {
          URLQuery: {},
          URLParams: {},
          Params: {
            project_name: tmpProjectName,
            industry_type: selectedIndustry,
            project_as: selectedRole,
            project_start_date: format(new Date(startDate), "yyyy-MM-dd"),
            project_end_date: format(new Date(endDate), "yyyy-MM-dd"),
            project_description: tmpProjectDescription,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status === 200) {
          triggerToast("created");
          addProjectModal = false;
        }
      } catch (error) {
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
      <Dropdown
        dropdownOpen={dropdownIndustryOpen}
        selectedItem={selectedIndustry}
        items={industries}
        placeholder={$_("selectIndustry")}
      />
    </div>
    <div class="flex gap-4">
      <div class="flex flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={$_("as")} />
        <Dropdown
          width="w-min"
          dropdownOpen={dropdownRole}
          selectedItem={selectedRole}
          items={role}
        />
      </div>
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
