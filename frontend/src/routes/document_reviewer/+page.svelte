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
  import { call, check } from "$lib";
  import { Typography, ActionCardAlt } from "$lib";

  const isFrontend = process.env.NODE_ENV === "development";
  let actionItems = $state([]);

  const getActionItems = async () => {
    actionItems = [
      {
        disabled: !check("aiManagement.contractreviewer.case:view"),
        href: isFrontend
          ? "/document_reviewer/case"
          : await $call.getPage("client-dashboard", "contract-case"),
        icon: "bookOpen",
        actionName: "Manage cases",
        desc: "desc Set your referenced document",
      },
      {
        href: isFrontend
          ? "/document_reviewer/report"
          : await $call.getPage("client-dashboard", "contract-report"),
        icon: "clipboardCheck",
        actionName: "Audit Reporting",
        desc: "desc Audit Reporting",
      },
      {
        disabled: !check("aiManagement.contractreviewer.project:create"),
        href: isFrontend
          ? "/document_reviewer/project"
          : await $call.getPage("client-dashboard", "contract-project-list"),
        icon: "packagePlus",
        actionName: "Manage Projects",
        desc: "desc Manage Projects",
      },
      {
        disabled: !check("aiManagement.contractreviewer.project:create"),
        href: isFrontend
          ? "/document_reviewer/compare"
          : await $call.getPage("client-dashboard", "contract-compare-list"),
        icon: "searchRefraction",
        actionName: "Compare Documents",
        desc: "desc Compare Documents",
      },
      {
        disabled: !check("aiManagement.contractreviewer.project:create"),
        href: isFrontend
          ? "/document_reviewer/audit_scan"
          : await $call.getPage("client-dashboard", "contract-audit-scan"),
        icon: "feature",
        actionName: "Audit Scan",
        desc: "desc Audit Scan",
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
          <ActionCardAlt {actionItems} />
        </div>
      </div>
    </div>
  </div>
</div>
