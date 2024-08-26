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

  let { illustration, description, actionText, funcButton }: { illustration: any, description: string, actionText: string, funcButton?: () => void } = $props();
</script>

<div class="pt-4 lg:pt-6">
  <div class="flex flex-col items-center gap-3 lg:gap-5">
    <svelte:component this={illustration} svgClass="text-utility-brand-50-light dark:text-gray-800" />
    {#if description}
    <Typography preset="Text md/Semibold" text={description} color="primary" />
    {/if}
  </div>
  {#if funcButton || actionText }
  <Button color="primary" text={actionText} {funcButton} />
  {/if}
</div>