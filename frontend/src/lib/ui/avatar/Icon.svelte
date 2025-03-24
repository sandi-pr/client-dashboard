<!--
  CR3 Client Dashboard
  Package: gitlab.hexacode.org/public-repo/js/hexacomponents-svelte
  Maintainer:
  * Sandi Permana <sandy@hexacode.org>

  Copyright (C) 2024 Hexacode Teknologi Indonesia
  All Rights Reserved.
-->

<script>
  import { twMerge } from "tailwind-merge";
  import { getIcons, getTextColors } from "$lib";

  let { name, size = "w-5 h-5", color = "secondary", propClass = "" } = $props();

  let classes = twMerge(size, getTextColors(color), propClass);
</script>

<svelte:component this={getIcons(name)} svgClass={classes} />

<!--
@component
## Config
@prop name; // Variable name of the icon.
@prop size = "w-5 h-5"; // Size of the icon.
@prop color = "secondary"; // Color of the icon.
-->
