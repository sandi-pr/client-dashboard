<script>
  import { Typography, InputRw, SelectedCard, triggerToast } from "$lib";
  import { Button, Textarea } from "flowbite-svelte";
  import { DatePicker } from "@svelte-plugins/datepicker";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import CalendarPlus02 from "$lib/icon/calendar-plus-02.svelte";
  import { format } from "date-fns";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { goto } from "$app/navigation";

  const call = new CR3Call(globalThis.location.origin);

  let data = $state({});
  let industries = $state([]);

  let projectAs = $state([
    {
      value: "1",
      name: "Buyer",
    },
    {
      value: "2",
      name: "Seller",
    },
  ])

  let isOpen = $state(false);
  let startDate = $state(null);
  let endDate = $state(null);
  let dateFormat = $state("d MMM yyyy");
  let formattedStartDate = $state("");
  let formattedEndDate = $state("");
  let selectedIndustry = $state("");
  let selectedRole = $state("");

  async function getIndustries() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      let page = 1;
      let tmpIndustries = [];
      while (true) {
        const response = await call.call("client-dashboard", "get-industries", {
          URLQuery: {},
          URLParams: {},
          Params: {
            industry_id: selectedIndustry,
            page: page,
            page_size: 25,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          if (response.result.data.length === 0) {
            break;
          } else {
            tmpIndustries = tmpIndustries.concat(response.result.data);
            page++;
          }
        }
      }

      let promise = tmpIndustries.map((industry) => {
        return {
          value: industry.id,
          name: industry.attributes.industry_name,
          description: industry.attributes.industry_description,
        };
      });

      industries = await Promise.all(promise);
    }
  }

  function createProject() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = call.call("client-dashboard", "create-contract-project", {
          URLQuery: {},
          URLParams: {},
          Params: {
            project_name: data.name,
            industry_type: selectedIndustry,
            project_as: selectedRole,
            project_start_date: format(new Date(startDate), "yyyy-MM-dd"),
            project_end_date: format(new Date(endDate), "yyyy-MM-dd"),
            project_description: data.description,
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
        console.log(error);
      }
    }
  }

  const formatDate = (dateString) => (dateString && format(new Date(dateString), dateFormat)) || "";

  $effect(() => {
    formattedStartDate = formatDate(startDate);
    formattedEndDate = formatDate(endDate);
  });

  const toggleDatePicker = () => {
    isOpen = !isOpen;
  };

  async function backPage() {
    goto(await call.getPage("client-dashboard", "contract-reviewer-overview"));
    return;
  }

  let buttonDisable = true;

  onMount(() => {
    getIndustries();
  });
</script>

<div class="flex flex-col px-4 pt-6 gap-6 pb-12">
  <!-- section header -->
  <div class="flex justify-between">
    <div class="flex flex-col gap-2">
      <div class="flex flex-col">
        <Typography preset="Display sm/Semibold" text={$_("createNewProject")} color="primary" />
      </div>
    </div>
    <div class="flex gap-2">
      <Button color="light" class="px-3 focus-within:!ring-1" onclick={backPage}
        >{$_("cancel")}</Button
      >
      <Button
        disabled={buttonDisable}
        onclick={createProject}
        class={`px-3 focus-within:!ring-1 ${
          !buttonDisable ? "cursor-pointer" : "cursor-not-allowed"
        }`}>{$_("create")}</Button
      >
    </div>
  </div>
  <Divider />
  <!-- section form -->
  <div class="flex gap-3">
    <div class="w-60">
      <Typography preset="Text md/Semibold" text={`${$_("industry")}*`} color="secondary" />
    </div>
    <div class="flex flex-col grow gap-3">
      <!-- form dropdown -->
      <div class="flex flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={`${$_("projectName")}*`} />
        <InputRw onEdit bind:value={data.name} />
      </div>
      <div class="flex grow flex-col gap-1.5">
        <Typography
          preset="Text sm/medium"
          color="secondary"
          text={`${$_("projectDuration")}*`}
        />
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
      <div class="flex flex-col gap-3">
        <div class="flex grow flex-col gap-1.5">
          <Typography preset="Text sm/Medium" text={`${$_("industryType")}*`} color="secondary" />
          <SelectedCard selectionData={industries} numGrid={3} />
        </div>
        <div class="flex grow flex-col gap-1.5">
          <Typography preset="Text sm/medium" color="secondary" text={$_("as")} />
          <SelectedCard selectionData={projectAs} numGrid={2} />
        </div>
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
          <Textarea bind:value={data.description} placeholder={$_("placholderDescription")} />
        </div>
      </div>
    </div>
  </div>
  <Divider />
</div>
