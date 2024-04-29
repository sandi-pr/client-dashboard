<script lang="ts">
  import { Radio } from 'flowbite-svelte';
  let { group, options = [] }: {
    group: string;
    options: string[];
  } = $props();
</script>

<ul class="gap-2">
  {#each options as option}
    <li><Radio bind:group={group} value={option}>{option}</Radio></li>
  {/each}
</ul>
