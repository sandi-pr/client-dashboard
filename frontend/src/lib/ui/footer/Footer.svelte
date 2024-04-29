<script lang="ts">
  import { Footer, FooterCopyright, FooterLink, FooterLinkGroup } from 'flowbite-svelte';
  import { onMount } from 'svelte';

  // Ekstrak properti langsung dari $props()
  let {
    name,
    year,
    links,
    placement = 'center',
    config = {}
  }: {
    name?: string;
    year?: number;
    links?: { href: string; name: string }[];
    placement?: 'left' | 'right' | 'center';
    config?: any;
  } = $props();

  let placementState = $state('center');
  let ready = $state(false);

  onMount(() => {
    if (config?.placement === 'left') {
      placementState = 'start';
    }
    if (config?.placement === 'right') {
      placementState = 'end';
    }
    if (config?.placement === 'center') {
      placementState = 'center';
    }
    ready = true;
  });
</script>

{#if ready}
  <div class={`flex justify-${placementState} items-center`}>
    <Footer>
      <FooterCopyright href="/" by={name} {year} />
      <FooterLinkGroup ulClass="flex flex-wrap items-center mt-3 text-sm text-gray-500 dark:text-gray-400 sm:mt-0">
        {#if links}
          {#each links as link}
            <FooterLink href={link.href}>{link.name}</FooterLink>
          {/each}
        {/if}
      </FooterLinkGroup>
    </Footer>
  </div>
{/if}

<!--
@component
## Config
@prop name: string; // Name displayed in the footer.
@prop year: number; // Year displayed in the footer.
@prop links: { href: string, name: string }[]; // Array of objects representing links in the footer.
@prop placement: 'left' | 'right' | 'center'; // Placement of the footer content.
-->
