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
  import { Heading } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import {
    call,
    check,
  } from "$lib";
  import { Typography, ReportCard, ColumnChart } from "$lib";

  const isFrontend = process.env.NODE_ENV === "development";
  let actionItems = $state([]);

  const getActionItems = async () => {
    actionItems = [
      {
        disabled: !check("aiManagement.contractreviewer.case:view"),
        href: isFrontend
          ? "/document_reviewer/case"
          : await $call.getPage("client-dashboard", "contract-case"),
        icon: "checkCircle",
        iconColor: "fg-success-primary",
        actionName: "1,273",
        desc: "compliantData",
        subtext: "20%",
        changeType: "increase",
      },
      {
        href: isFrontend
          ? "/document_reviewer/report"
          : await $call.getPage("client-dashboard", "contract-report"),
        icon: "clipboardCheck",
        iconColor: "fg-warning-primary",
        actionName: "6,256",
        desc: "nonCompliantData",
        subtext: "50%",
        changeType: "decrease",
      },
      {
        disabled: !check("aiManagement.contractreviewer.project:create"),
        href: isFrontend
          ? "/document_reviewer/project"
          : await $call.getPage("client-dashboard", "contract-project-list"),
        icon: "packagePlus",
        iconColor: "fg-secondary",
        actionName: "11,890",
        desc: "totalData",
        subtext: "40%",
        changeType: "increase",
      },
    ];
  };

  onMount(async () => {
    await getActionItems();
  });

</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-32">
    <div class="flex flex-col m-4">

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Typography
          preset="Display sm/Semibold"
          text={$_("auditReporting")}
          propClass="fluid-xl"
        />

        <ColumnChart />

        <!-- Action Items -->
        <div>
          {#if actionItems.some((item) => item.disabled !== true)}
            <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
              <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
                {$_("recentProjectSummary")}
              </Heading>
            </div>
          {/if}
          <ReportCard {actionItems} />
        </div>
      </div>
    </div>
  </div>
</div>
