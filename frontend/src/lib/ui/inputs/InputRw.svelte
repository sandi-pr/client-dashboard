<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import { Input } from "flowbite-svelte";
  import { EyeOutline, EyeSlashOutline } from "flowbite-svelte-icons";
  import { twMerge } from "tailwind-merge";

  let {
    type = "text",
    onEdit = true,
    value = $bindable(),
    disabled = $bindable(false),
    placeholder = "",
    oninput = () => {},
    show = false,
  } = $props();

  let disabledInput = $state(false);

  $effect(() => {
    if (onEdit && disabled) {
      disabledInput = true;
    } else {
      disabledInput = false;
    }
  });

  const inputClass = () => {
    const disabledClass = onEdit
      ? ""
      : "cursor-default focus:ring-0 border-disabled-subtle-light dark:border-disabled-subtle-dark focus:border-disabled-subtle-light focus:dark:border-disabled-subtle-dark bg-disabled-subtle-light dark:bg-disabled-subtle-dark text-disabled-light dark:text-disabled-dark";

    return twMerge(
      "text-primary-light dark:text-primary-dark focus:ring-0 border-primary-light dark:border-primary-dark focus:border-primary-light focus:dark:border-primary-dark bg-primary-light dark:bg-primary-dark",
      disabledClass,
    );
  };
</script>

{#if type === "password" && onEdit}
  <Input
    type={show ? "text" : "password"}
    readonly={!onEdit}
    bind:value
    bind:disabled={disabledInput}
    {placeholder}
    on:input={oninput}
    class={inputClass()}
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
    {type}
    readonly={!onEdit}
    bind:value
    bind:disabled={disabledInput}
    {placeholder}
    on:input={oninput}
    class={inputClass()}
  />
{/if}
