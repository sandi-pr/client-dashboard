<!--
  CR3 Client Dashboard
  Package: gitlab.hexacode.org/public-repo/js/hexacomponents-svelte
  Maintainer:
  * Sandi Permana <sandy@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>

  Copyright (C) 2024 Hexacode Teknologi Indonesia
  All Rights Reserved.
-->

<script lang="ts">
  import { P } from "flowbite-svelte";
  import { getTextColors } from "$lib";

  let { preset = "Text md/Regular", text = "", color = "primary", propClass = "" } = $props();

  type Preset = {
    size: string;
    leading: string;
  };

  const presets: Record<string, Preset> = {
    "Display xl": { size: "6xl", leading: "72px" },
    "Display lg": { size: "5xl", leading: "60px" },
    "Display md": { size: "4xl", leading: "44px" },
    "Display sm": { size: "3xl", leading: "38px" },
    "Display xs": { size: "2xl", leading: "loose" },
    "Text xl": { size: "xl", leading: "30px" },
    "Text lg": { size: "lg", leading: "7" },
    "Text md": { size: "", leading: "" },
    "Text sm": { size: "sm", leading: "tight" },
    "Text xs": { size: "xs", leading: "18px" }
  };

  let [group, weight] = preset.split("/") as [string, string];
  let { size, leading } = presets[group] || {};

  let classes = propClass + leading.includes("px") ? ` leading-[${leading}]` : ` leading-${leading}`;
</script>

<P {size} color={getTextColors(color)} weight={weight.toLowerCase()} class={classes}>
  {text}
</P>

<!--
@component
## Config
@prop preset: string; // Preset of the text (size and weight).
@example preset: '{Display/Text} {xl-xs}/{Reguler/Medium/Semibold/Bold}' default: 'Text md/Regular'
@prop text: string; // Text to be displayed.
@prop color: string; // Color of the text.
@example color: '{primary/secondary/tertiary/quaternary/disabled/placeholder/brand/success/warning/error}'
-->
