<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: 
  - Azzis Arswendo <azzis@hexacode.org>
  - Dean Ramhan <dean@hexacode.org>
  - Taqie Fadlillah <taqief@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->

<script>
  import { SectionHeader, Card, Typography, CR3Call } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);
  let userData = $state({});

  let isFirstTime = true;
</script>

<div class="flex flex-col gap-1 px-4 pt-6">
  <Typography preset="Display sm/Semibold" text="Welcome to Crossroads!" color="primary" />
  <Typography
    text="{userData?.data?.username || 'Hexastaff'} - {userData?.data?.email ||
      'staff@hexacode.org'}"
    color="tertiary"
  />
</div>
<div class="flex flex-col gap-4 pt-8 px-4">
  <!-- {#if isFirstTime}
    <StateCard
      illustration={IllustrationFirst}
      description="Welcome to Crossroads! We need you to set up your product."
      actionText="Product Setup"
      funcButton={() => (isFirstTime = false)}
    />
  {:else}
    <SelectedCard {companyType} />
  {/if} -->
  <SectionHeader title="Available Modules" />
  <Card
    icon="users"
    title="CRM"
    description="Manage customer contacts, and personalize your customer service"
    onclick={async () => {
      await call.goToPage("client-dashboard", "crm-overview");
      return;
    }}
  />
  <Card
    icon="messageChat"
    title="Hexa Bot"
    description="Manage bots and automate your workflows"
    onclick={async () => {
      await call.goToPage("client-dashboard", "hexabot-overview");
      return;
    }}
  />
</div>
