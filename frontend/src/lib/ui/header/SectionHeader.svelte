<!--
  CR3 Client Dashboard
  Package: gitlab.hexacode.org/public-repo/js/hexacomponents-svelte
  Maintainer:
  * Sandi Permana <sandy@hexacode.org>

  Copyright (C) 2024 Hexacode Teknologi Indonesia
  All Rights Reserved.
-->

<script lang="ts">
  import { Typography, Button } from "$lib";
  import Divider from "$lib/components/Divider.svelte";
  import { goto } from '$app/navigation';

  let { title = 'Section Header', hrefAction = '' }: { title: string, hrefAction?: string } = $props();
</script>

<div class="flex flex-col gap-1.5 md:gap-5">
  <div class="flex justify-between items-center">
    <Typography preset="Text lg/Semibold" text={title} color="primary" />
    {#if hrefAction}
    <Button color="light" text="View All" funcButton={() => {goto(hrefAction)}} spaceY="mt-0" />
    {/if}
  </div>
  <Divider />
</div>

<!--
@component
## Config
@prop title: string; // Title of the section header.
@prop hrefAction: string; // Link to the action button of view all.
-->