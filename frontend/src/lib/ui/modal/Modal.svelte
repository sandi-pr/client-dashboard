<script lang="ts">
  import { Button, Modal, type ModalPlacementType } from 'flowbite-svelte';
  import { onMount } from 'svelte';

  let id = $state('');
  let placement: ModalPlacementType = $state('center');
  let open = $state(false);

  const setPlacement = () => {
    placement = (config?.placement as ModalPlacementType) ?? 'center';
    id = `${placement}-modal`;
    open = !open;
  };

  let { nameButton, title, content, footerButtons, placement: configPlacement, config = {} }: {
    nameButton: string;
    title: string;
    content: string;
    footerButtons?: { color: "red" | "yellow" | "green" | "purple" | "blue" | "light" | "dark" | "primary" | "none" | "alternative", onclick: () => void, text: string }[];
    placement?: 'left' | 'right' | 'center';
    config?: any;
  } = $props();

  onMount(() => {
    if (configPlacement === 'left') {
      placement = 'start' as ModalPlacementType;
    } else if (configPlacement === 'right') {
      placement = 'end' as ModalPlacementType;
    }
  });
</script>

<div class="inline-grid grid-cols-3 grid-rows-3 gap-4">
  <Button on:click={setPlacement}>{nameButton}</Button>
</div>

<Modal {id} title={title} bind:open placement={placement} autoclose>
  <p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">{content}</p>
  <svelte:fragment slot="footer">
    {#if footerButtons}
      {#each footerButtons as button}
        <Button color={button.color} on:click={button.onclick}>{button.text}</Button>
      {/each}
    {/if}
  </svelte:fragment>
</Modal>

<!--
@component
## Config
@prop nameButton: string; // Text displayed on the button to trigger the modal.
@prop title: string; // Title of the modal.
@prop content: string; // Content of the modal.
@prop footerButtons: { color: "red" | "yellow" | "green" | "purple" | "blue" | "light" | "dark" | "primary" | "none" | "alternative", onclick: () => void, text: string }[]; // Array of objects representing buttons in the modal footer.
@prop placement: 'left' | 'right' | 'center'; // Placement of the modal.
-->
