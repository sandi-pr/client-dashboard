<script>
  import { showAI, showNotif, Card, toCamelCase } from "$lib";
  import { _ } from "svelte-i18n";

  let { actionItems, onclick = null, numGrid = "2" } = $props();
</script>

<div
  class={`grid ${$showAI || $showNotif ? "grid-cols-1" : "grid-cols-1 lg:grid-cols-" + numGrid} gap-2 md:gap-4 w-full items-center justify-center`}
>
  {#each actionItems as item}
    {#if !item?.disabled}
      <Card
        href={item.onclick ? null : item.href}
        onclick={item.onclick ?? null}
        type={"alt"}
        icon={item.icon}
        title={$_(toCamelCase(item.actionName))}
        description={$_(toCamelCase(item.desc))}
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
