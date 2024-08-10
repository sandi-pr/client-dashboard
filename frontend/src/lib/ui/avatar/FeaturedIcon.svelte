<!--
  CR3 Client Dashboard
  Package: gitlab.hexacode.org/public-repo/js/hexacomponents-svelte
  Maintainer:
  * Sandi Permana <sandy@hexacode.org>

  Copyright (C) 2024 Hexacode Teknologi Indonesia
  All Rights Reserved.
-->

<script lang="ts">
  import { Avatar } from "flowbite-svelte";
  let { classesSvg }: { classesSvg: any } = $props();
</script>

<Avatar size="sm" border class="!ring-0 !m-0 bg-brand-secondary-light dark:bg-brand-secondary-dark !p-2 !h-10 !w-10">
  <svelte:component this={classesSvg} svgClass="w-5 h-5 text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark" />
</Avatar>

<!--
@component
## Config
@prop classesSvg; // Import file of the icon from svg convert to Svelte component.
-->
