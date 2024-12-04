<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import { Input } from "flowbite-svelte";
  import { EyeOutline, EyeSlashOutline } from 'flowbite-svelte-icons';

  let { type = "text", onEdit = false, value = $bindable(), disabled = $bindable(false), placeholder = "", oninput = () => {}, show = false } = $props();

  let disabledInput = $state(false);

  $effect(() => {
    if (onEdit && disabled) {
      disabledInput = true;
    } else {
      disabledInput = false;
    }
  })
</script>

{#if type === "password" && onEdit}
  <Input
    type={show ? 'text' : 'password'}
    readonly={!onEdit}
    bind:value={value}
    bind:disabled={disabledInput}
    placeholder={placeholder}
    on:input={oninput}
    class={`bg-primary-light dark:bg-secondary-dark border-primary-light dark:border-primary-dark ${!onEdit ? "cursor-default focus:ring-0 focus:border-gray-300 dark:focus:border-gray-600" : "focus:border-primary-500 focus:ring-primary-500 dark:focus:border-primary-500 dark:focus:ring-primary-500"}`}
  >
    <button slot="left" onclick={() => (show = !show)} class="pointer-events-auto">
      {#if show}
        <EyeOutline class="w-6 h-6" />
      {:else}
        <EyeSlashOutline class="w-6 h-6" />
      {/if}
    </button>
  </Input>
{:else}
<Input
  type={type}
  readonly={!onEdit}
  bind:value={value}
  bind:disabled={disabledInput}
  placeholder={placeholder}
  on:input={oninput}
  class={`bg-primary-light dark:bg-secondary-dark border-primary-light dark:border-primary-dark ${!onEdit ? "cursor-default focus:ring-0 focus:border-gray-300 dark:focus:border-gray-600" : "focus:border-primary-500 focus:ring-primary-500 dark:focus:border-primary-500 dark:focus:ring-primary-500"}`}
/>
{/if}
