<!--
  CR3 Client Dashboard
  Package: gitlab.hexacode.org/public-repo/js/hexacomponents-svelte
  Maintainer:
  * Sandi Permana <sandy@hexacode.org>

  Copyright (C) 2024 Hexacode Teknologi Indonesia
  All Rights Reserved.
-->

<script>
  import { Avatar } from "flowbite-svelte";
  import { getBgColors, getTextColors } from "$lib";
  let { classesSvg, color = "brand" } = $props();

  const colors = {
    brand: getBgColors("brand-primary"),
    primary: getBgColors(color),
    error: getBgColors("error-secondary"),
  };

  const colorsIcon = {
    brand: getTextColors("fg-brand-primary-alt"),
    error: getTextColors("fg-error-primary"),
  };

  const classes = "!ring-0 !m-0 !p-2 !h-10 !w-10 " + colors[color];
  const classesIcon = "w-5 h-5 " + colorsIcon[color];
</script>

<Avatar size="sm" border class={classes}>
  <svelte:component this={classesSvg} svgClass={classesIcon} />
</Avatar>

<!--
@component
## Config
@prop classesSvg; // Import file of the icon from svg convert to Svelte component.
-->
