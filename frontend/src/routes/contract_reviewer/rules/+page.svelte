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
  import { Modal, Button, Input } from "flowbite-svelte";
  import { Typography, triggerToast, StateCard } from "$lib";
  import * as Icon from "flowbite-svelte-icons";
  import CalendarPlus02 from "$lib/icon/calendar-plus-02.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { _ } from "svelte-i18n";
  import { CR3Call } from "hexacomponents-svelte";
  import { format } from "date-fns";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";

  const call = new CR3Call(globalThis.location.origin);

  let files = $state([]);
  let formData = $state({ name: "" });
  let tmpRulesName = $state("");
  let tmpRulesDescription = $state("");
  let selectedIndustry = $state("Oil & Gas");
  let selectedRole = $state("Buyer");
  let addRulesModal = $state(false);
  let isOpen = $state(false);
  let startDate = $state(null);
  let endDate = $state(null);
  let dateFormat = $state("d MMM yyyy");
  let formattedStartDate = $state("");
  let formattedEndDate = $state("");

  let configs = $state({});

  let datas = $state([]);

  function createRules() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = call.call("client-dashboard", "create-contract-project", {
          URLQuery: {},
          URLParams: {},
          Params: {
            project_name: tmpRulesName,
            industry_type: selectedIndustry,
            project_as: selectedRole,
            project_start_date: format(new Date(startDate), "yyyy-MM-dd"),
            project_end_date: format(new Date(endDate), "yyyy-MM-dd"),
            project_description: tmpRulesDescription,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status === 200) {
          triggerToast("created");
          addRulesModal = false;
        }
      } catch (error) {
        triggerToast("error");
        addRulesModal = false;
        console.log(error);
      }
    }
  }

  const toggleDatePicker = () => {
    isOpen = !isOpen;
  };

  const handleSearch = (evt) => {
    console.log(evt.target.value);
  };

  async function getRules() {
    // let token = sessionStorage.getItem("token");
    let token = true;
    if (token) {
      // try {
      //   const response = await call.call("client-dashboard", "get-contract-projects", {
      //     URLQuery: {},
      //     URLParams: {},
      //     Params: {
      //       page: 1,
      //       page_size: 5,
      //     },
      //     Headers: {
      //       Authorization: "Bearer " + token,
      //     },
      //   });
      //   console.log("response", response);
      //   if (response.status === 200) {
      //     configs = {
      //       headers: "projectList",
      //       actions: [
      //         {
      //           name: "viewAll",
      //           href: await call.getPage("client-dashboard", "contract-project-list"),
      //           type: "info",
      //         },
      //       ],
      //       seperator: true,
      //     };
      //     let url = await call.getPage("client-dashboard", "contract-project-detail");
      //     let promise = response.result.data.map((project) => {
      //       return {
      //         id: project.id,
      //         project_name: project.attributes.project_name,
      //         industry: project.attributes.industry_type,
      //         start_date: format(new Date(project.attributes.project_start_date), "dd/MM/yyyy"),
      //         end_date: format(new Date(project.attributes.project_end_date), "dd/MM/yyyy"),
      //         actions: {
      //           objects: [
      //             {
      //               href: url + `?projectId=${project.id}`,
      //               name: "View",
      //             },
      //           ],
      //         },
      //       };
      //     });

      //     datas = await Promise.all(promise);
      //     console.log("datas", datas);
      //   }
      // } catch (error) {
      //   console.log(error);
      // }
      datas = [
        {
          id: 1,
          project_name: "Project 1",
          actions: {
            objects: [
              {
                href: "#",
                name: "View",
              },
            ],
          },
        },
        {
          id: 2,
          project_name: "Project 2",
          actions: {
            objects: [
              {
                href: "#",
                name: "View",
              },
            ],
          },
        },
        {
          id: 3,
          project_name: "Project 3",
          actions: {
            objects: [
              {
                href: "#",
                name: "View",
              },
            ],
          },
        },
        {
          id: 4,
          project_name: "Project 4",
          actions: {
            objects: [
              {
                href: "#",
                name: "View",
              },
            ],
          },
        },
        {
          id: 5,
          project_name: "Project 5",
          actions: {
            objects: [
              {
                href: "#",
                name: "View",
              },
            ],
          },
        },
      ];
    }
  }

  onMount(async () => {
    getRules();
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
        <div class="flex justify-between">
          <Typography preset="Display sm/Semibold" text={$_("rulesList")} />
          <Button
            onclick={() => (addRulesModal = true)}
            class="flex gap-2 items-center"
            color="primary"
            variant="contained"
            size="sm"
          >
            <svelte:component this={CalendarPlus02} class="w-4 h-4" />
            {$_("addNewRules")}
          </Button>
        </div>

        <Input
          id="search"
          type="search"
          class="bg-white dark:bg-hc-graydark-950 py-3"
          placeholder={$_("search")}
          on:input={handleSearch}
        >
          <Icon.SearchOutline slot="left" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
        </Input>

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

<Modal title={$_("createNewRules")} bind:open={addRulesModal} size="lg">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("rulesName")}*`} />
      <Input bind:value={formData.name} placeholder={$_("whatIsYourRulesName")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={createRules} class="w-full">{$_("addNewRules")}</Button>
  </svelte:fragment>
</Modal>
