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

  let tmpRulesName = $state("");
  let addRulesModal = $state(false);

  let configs = $state({});

  let datas = $state([]);

  async function createRules() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "create-contract-rules", {
          URLQuery: {},
          URLParams: {},
          Params: {
            rule_name: tmpRulesName,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerToast("error");
          console.log(response);
          return;
        }

        triggerToast("created");
        addRulesModal = false;
        tmpRulesName = "";
        getRules();
        goto(
          (await call.getPage("client-dashboard", "contract-rule-detail")) +
            "?ruleId=" +
            response.result.id,
        );
      } catch (error) {
        triggerToast("error");
        addRulesModal = false;
        console.log(error);
      }
    }
  }

  const handleSearch = (evt) => {
    console.log(evt.target.value);
  };

  async function getRules() {
    let token = sessionStorage.getItem("token");
    if (token) {
      let rules = [];
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
          rules = rules.concat(response.result.data);
          page++;
        }

        if (rules.length > 0) {
          let link = await call.getPage("client-dashboard", "contract-rule-detail");
          let promise = rules.map((rule) => {
            return {
              id: rule.id,
              rule_name: rule.attributes.rule_name,
              actions: {
                objects: [
                  {
                    href: link + "?ruleId=" + rule.id,
                    name: "View",
                  },
                ],
              },
            };
          });
          datas = await Promise.all(promise);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  onMount(async () => {
    getRules();
  });

  let buttonDisabled = $state(false);

  $effect(() => {
    buttonDisabled = tmpRulesName.length === 0;
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
      <Input bind:value={tmpRulesName} placeholder={$_("whatIsYourRulesName")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={createRules} class="w-full" disabled={buttonDisabled}
      >{$_("addNewRules")}</Button
    >
  </svelte:fragment>
</Modal>
