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
  import { Modal, Button, Input, Textarea, Heading, P } from "flowbite-svelte";
  import { Typography, Dropdown, triggerToast, StateCard } from "$lib";
  import * as Icon from "flowbite-svelte-icons";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import { _ } from "svelte-i18n";
  import { CR3Call } from "hexacomponents-svelte";
  import { format } from "date-fns";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";

  const call = new CR3Call(globalThis.location.origin);

  let files = $state([]);
  let formData = $state({ title: "", prompt: "", sample: "" });
  let tmpRulesName = $state("");
  let tmpRulesDescription = $state("");
  let selectedIndustry = $state("Oil & Gas");
  let selectedRole = $state("Buyer");
  let addRulesSectionsModal = $state(false);
  let startDate = $state(null);
  let endDate = $state(null);
  let dateFormat = $state("d MMM yyyy");
  let formattedStartDate = $state("");
  let formattedEndDate = $state("");

  let datas = $state({});

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
          addRulesSectionsModal = false;
        }
      } catch (error) {
        triggerToast("error");
        addRulesSectionsModal = false;
        console.log(error);
      }
    }
  }

  async function getRulesSections() {
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
      //     let promise = response.result.pageData.map((project) => {
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
      datas = {
        title: "Document Title",
        createdAt: "2024-01-01",
        data: [
          {
            id: 1,
            title: "buyer A",
            prompt: "Who is the buyer?",
            sample: "the buyer is {name}",
            createdAt: "2024-01-01",
          },
          {
            id: 2,
            title: "seller A",
            prompt: "Who is the seller?",
            sample: "the seller is {name}",
            createdAt: "2024-01-01",
          },
          {
            id: 3,
            title: "buyer B",
            prompt: "Who is the buyer?",
            sample: "the buyer is {name}",
            createdAt: "2024-01-01",
          },
          {
            id: 4,
            title: "seller B",
            prompt: "Who is the seller?",
            sample: "the seller is {name}",
            createdAt: "2024-01-01",
          },
        ],
      };
    }
  }

  onMount(async () => {
    getRulesSections();
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
        <div class="flex flex-col grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
              >{datas.title}<Button
                outline
                color="red"
                class="hover:!bg-transparent !border-hc-graylight-200 dark:!border-hc-graydark-600 focus-within:!ring-1 !p-2"
              >
                <CopyIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
              </Button></Heading
            >
          </div>
          <P
            class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
            >{$_("createOn")}{` ${datas.createdAt}`}</P
          >
        </div>

        <div class="w-full border-t"></div>

        <!-- table latest -->
        {#if Object.keys(datas).length > 0}
          {#each datas.data as data, index}
            <div class="grid grid-cols-12 gap-4">
              <div class="col-span-3">
                <P class="font-semibold text-tertiary-light">Section {index + 1}</P>
              </div>
              <div class="col-span-1 text-tertiary-light"><P>Order</P></div>
              <div class="col-span-8 text-tertiary-light"><P>Name</P></div>
              <div class="col-span-3"></div>
              <div class="col-span-1">
                <Input value={"#" + (index + 1)} readonly />
              </div>
              <div class="col-span-8">
                <Input value={data.title} readonly />
              </div>
              <div class="col-span-3"></div>
              <div class="col-span-9">
                <P class="text-tertiary-light">Extract Data Prompt</P>
              </div>
              <div class="col-span-3"></div>
              <div class="col-span-9">
                <Textarea value={data.prompt} readonly />
              </div>
              <div class="col-span-3"></div>
              <div class="col-span-9">
                <P class="text-tertiary-light">Extract Output Prompt</P>
              </div>
              <div class="col-span-3"></div>
              <div class="col-span-9">
                <Textarea value={data.sample} readonly />
              </div>
            </div>
            <div class="w-full border-t"></div>
          {/each}
          <Button
            class="self-center"
            color="primary"
            on:click={() => (addRulesSectionsModal = true)}
            ><Icon.PlusOutline class="w-5 h-5" />
            {$_("addNewSection")}</Button
          >
        {:else}
          <div class="pt-8">
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<Modal title={$_("createNewRulesSection")} bind:open={addRulesSectionsModal} size="lg">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("rulesSectionTitle")}*`} />
      <Input bind:value={formData.title} placeholder={$_("whatIsYourRulesSectionsTitle")} />
    </div>
  </div>
  <div class="flex flex-col gap-1.5">
    <Typography preset="Text sm/medium" color="secondary" text={`${$_("rulesSectionPrompt")}*`} />
    <Input bind:value={formData.prompt} placeholder={$_("whatIsYourRulesSectionsPrompt")} />
  </div>
  <div class="flex flex-col gap-1.5">
    <Typography preset="Text sm/medium" color="secondary" text={`${$_("rulesSectionSample")}*`} />
    <Input bind:value={formData.sample} placeholder={$_("whatIsYourRulesSectionsSample")} />
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={createRules} class="w-full">{$_("addNewRulesSection")}</Button>
  </svelte:fragment>
</Modal>
