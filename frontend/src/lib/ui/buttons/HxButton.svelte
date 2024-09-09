<!-- <script>
  import { Button } from "flowbite-svelte";
  import { Typography, getIcons, getTextColors, getBgColors } from "$lib";

  let { 
    onClick = () => {},
    iconName = "",
    iconColor = "brand-secondary",
    iconSize = "w-5 h-5",
    bgColor = "",
    spacing = "",
    sizing = "",
    text = "",
    presetText = "Text md/Regular",
    colorText = "primary",
    propText = "",
    outline = false
  } = $props();
</script>

<Button
  outline={outline}
  on:click={onClick}
  class={`${sizing} ${spacing} ${getBgColors(bgColor)}`}
>
  {#if iconName !== ""}
    <svelte:component this={getIcons(iconName)} svgClass={`${iconSize} ${getTextColors(iconColor)}`} />
  {/if}
  <Typography
    preset={presetText}
    {text}
    color={colorText}
    propClass={propText}
  />
</Button> -->