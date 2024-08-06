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
  import { onMount, onDestroy } from "svelte";
  import { Modal, Button, Input, Textarea, Heading, P } from "flowbite-svelte";
  import { Typography, Dropdown, triggerToast, StateCard, triggerLoading } from "$lib";
  import * as Icon from "flowbite-svelte-icons";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import { _ } from "svelte-i18n";
  import { CR3Call } from "hexacomponents-svelte";
  import { format } from "date-fns";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { page } from "$app/stores";

  const call = new CR3Call(globalThis.location.origin);

  let formData = $state({ title: "", prompt: "", sample: "" });
  let addRulesSectionsModal = $state(false);
  let deleteModal = $state(false);
  let edit = $state(false);

  let ruleId = $state(null);

  let tmpSectionId = $state(null);

  let datas = $state({});

  async function getRulesSections() {
    triggerLoading(true);
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "get-contract-rule", {
          URLQuery: {},
          URLParams: {
            id: ruleId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          console.log(response);
          return;
        }

        const response2 = await call.call("client-dashboard", "get-contract-rule-section", {
          URLQuery: {},
          URLParams: {},
          Params: {
            contract_rule_id: ruleId.toString(),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          console.log(response2);
          return;
        }

        datas = {
          title: response.result.data[0].attributes.rule_name,
          createdAt: format(response.result.data[0].attributes.createdAt, "yyyy-MM-dd HH:mm:ss"),
          data:
            response2.result.data.length < 1
              ? []
              : response2.result.data.map((item) => {
                  return {
                    id: item.id,
                    title: item.attributes.section_title,
                    prompt: item.attributes.section_prompt,
                    sample: item.attributes.section_sample,
                    createdAt: format(item.attributes.createdAt, "yyyy-MM-dd"),
                  };
                }),
        };

        triggerLoading(false);
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function createRulesSections() {
    addRulesSectionsModal = false;
    triggerLoading(true);
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "create-contract-rule-section", {
          URLQuery: {},
          URLParams: {},
          Params: {
            contract_rule_id: ruleId.toString(),
            section_title: formData.title,
            section_prompt: formData.prompt,
            section_sample: formData.sample,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          addRulesSectionsModal = true;
          console.log(response);
          return;
        }

        triggerLoading(false);
        formData = { title: "", prompt: "", sample: "" };
        triggerToast("created");
        addRulesSectionsModal = false;
        getRulesSections();
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function deleteSection() {
    deleteModal = false;
    triggerLoading(true);
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "delete-contract-rule-section", {
          URLQuery: {},
          URLParams: {
            id: tmpSectionId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          console.log(response);
          return;
        }

        triggerLoading(false);
        triggerToast("deleted");
        getRulesSections();
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function editSection() {
    addRulesSectionsModal = false;
    triggerLoading(true);
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "update-contract-rule-section", {
          URLQuery: {},
          URLParams: {
            id: tmpSectionId.toString(),
          },
          Params: {
            section_title: formData.title,
            section_prompt: formData.prompt,
            section_sample: formData.sample,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          addRulesSectionsModal = true;
          console.log(response);
          return;
        }

        triggerLoading(false);
        triggerToast("updated");
        formData = { title: "", prompt: "", sample: "" };
        edit = false;
        getRulesSections();
      } catch (error) {
        console.log(error);
      }
    }
  }

  onMount(() => {
    const unsubscribe = page.subscribe(async ($page) => {
      ruleId = $page.url.searchParams.get("ruleId");
      if (ruleId !== null) {
        getRulesSections();
      } else {
        goto(await call.getPage("client-dashboard", "contract-rules"));
      }
    });
    return () => unsubscribe();
  });

  $effect(async () => {
    const unsubscribe = page.subscribe(async ($page) => {
      ruleId = $page.url.searchParams.get("ruleId");
      if (ruleId !== null) {
        getRulesSections();
      } else {
        goto(await call.getPage("client-dashboard", "contract-rules"));
      }
    });
    onDestroy(() => unsubscribe());
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
                <div class="flex">
                  <button
                    class="text-primary-light dark:text-primary-dark mr-4"
                    onclick={() => {
                      tmpSectionId = data.id;
                      deleteModal = true;
                    }}>Delete</button
                  >
                  <button
                    class="text-brand-secondary-light dark:text-brand-secondary-dark"
                    onclick={() => {
                      edit = true;
                      tmpSectionId = data.id;
                      formData = {
                        title: data.title,
                        prompt: data.prompt,
                        sample: data.sample,
                      };
                      addRulesSectionsModal = true;
                    }}>Edit</button
                  >
                </div>
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
            on:click={() => {
              formData = { title: "", prompt: "", sample: "" };
              edit = false;
              addRulesSectionsModal = true;
            }}
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

<Modal
  title={edit ? $_("editRulesSection") : $_("createNewRulesSection")}
  bind:open={addRulesSectionsModal}
  size="lg"
>
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
    <Button
      onclick={edit ? editSection : createRulesSections}
      class="w-full"
      disabled={formData.title.length < 1 ||
        formData.prompt.length < 1 ||
        formData.sample.length < 1}>{edit ? $_("saveChanges") : $_("addNewRulesSection")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <Icon.ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteSection} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button
      color="alternative"
      onclick={() => {
        deleteModal = false;
        tmpSectionId = null;
      }}>{$_("noCancel")}</Button
    >
  </div>
</Modal>
