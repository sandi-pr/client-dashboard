<script>
  import { showAI, showNotif, Card, toCamelCase } from "$lib";
  import { _ } from "svelte-i18n";

  let { actionItems, onclick = null, numGrid = "3" } = $props();
</script>

<div
  class={`grid ${$showAI || $showNotif ? "grid-cols-1" : "grid-cols-1 lg:grid-cols-" + numGrid} gap-2 md:gap-4 w-full items-center justify-center`}
>
  {#each actionItems as item}
    {#if !item?.disabled}
      <Card
        href={item.onclick ? null : item.href ? item.href : null}
        onclick={item.onclick ?? null}
        type={"reportNoButton"}
        icon={item.icon}
        iconColor={item.iconColor}
        title={$_(toCamelCase(item.actionName))}
        description={$_(toCamelCase(item.desc))}
        subtext={item.subtext}
        changeType={item.changeType}
      />
    {/if}
  {/each}
</div>

<!--
@component
## Config
@prop actionItems: { href: string, icon: string, actionName: string, desc: string }[]; // Array of action items.
@prop onclick: () => void; // Function to be called when the card is clicked.
-->
