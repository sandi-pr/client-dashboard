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
  import { on } from "svelte/events";
  import { onMount } from "svelte";
  let { classesSvg, color = "brand", iconColor = "" } = $props();

  const colors = {
    brand: getBgColors("brand-secondary"),
    primary: getBgColors("primary"),
    alt: getBgColors("brand-solid"),
    error: getBgColors("error-secondary"),
  };

  const colorsIcon = {
    brand: getTextColors("fg-brand-primary-alt"),
    alt: getTextColors("primary-onbrand"),
    error: getTextColors("fg-error-primary"),
  };

  const classes =`!ring-0 !m-0 !p-2 !h-10 !w-10 ${color === "alt" ? "rounded-lg" : color === "primary" ? "rounded-lg border border-primary dark:border-primary-dark" : "rounded-full"} ` + colors[color];
  const classesIcon = `w-5 h-5 ${iconColor !== "" ? getTextColors(iconColor) : colorsIcon[color]} `;
</script>

<Avatar size={color === "alt" ? "md" : "sm"} border rounded={color === "alt" || color === "primary"} class={classes}>
  <svelte:component this={classesSvg} svgClass={classesIcon} />
</Avatar>

<!--
@component
## Config
@prop classesSvg; // Import file of the icon from svg convert to Svelte component.
-->
