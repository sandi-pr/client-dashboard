<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: 
  - Taqie Fadlillah <taqief@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->

<script>
  import { fly } from "svelte/transition";
  import { Toast } from "flowbite-svelte";
  import {
    toastCreated,
    toastUpdated,
    toastDeleted,
    toastError,
    toastNotAvailable,
    toastDayNotAvailable,
    toastTimeNotAvailable,
    toastCopyToClipboard,
    customToasts,
  } from "$lib/utils/toast";
  import { CheckCircleSolid, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { _ } from "svelte-i18n";
</script>

<div class="fixed bottom-5 z-50 end-5">
  <Toast color="green" transition={fly} params={{ x: 200 }} bind:open={$toastCreated}>
    <svelte:fragment slot="icon">
      <CheckCircleSolid class="w-5 h-5" />
    </svelte:fragment>
    {$_("successfullyCreated")}
  </Toast>
  <Toast color="green" transition={fly} params={{ x: 200 }} bind:open={$toastUpdated}>
    <svelte:fragment slot="icon">
      <CheckCircleSolid class="w-5 h-5" />
    </svelte:fragment>
    {$_("successfullyUpdated")}
  </Toast>
  <Toast color="green" transition={fly} params={{ x: 200 }} bind:open={$toastDeleted}>
    <svelte:fragment slot="icon">
      <CheckCircleSolid class="w-5 h-5" />
    </svelte:fragment>
    {$_("successfullyDeleted")}
  </Toast>
  <Toast
    class="!max-w-sm"
    color="red"
    transition={fly}
    params={{ x: 200 }}
    bind:open={$toastNotAvailable}
  >
    <svelte:fragment slot="icon">
      <ExclamationCircleSolid class="w-5 h-5" />
    </svelte:fragment>
    {$_("notAvailable")}
  </Toast>
  <Toast
    class="!max-w-sm"
    color="red"
    transition={fly}
    params={{ x: 200 }}
    bind:open={$toastDayNotAvailable}
  >
    <svelte:fragment slot="icon">
      <ExclamationCircleSolid class="w-5 h-5" />
    </svelte:fragment>
    {$_("dayNotAvailable")}
  </Toast>
  <Toast
    class="!max-w-sm"
    color="red"
    transition={fly}
    params={{ x: 200 }}
    bind:open={$toastTimeNotAvailable}
  >
    <svelte:fragment slot="icon">
      <ExclamationCircleSolid class="w-5 h-5" />
    </svelte:fragment>
    {$_("timeNotAvailable")}
  </Toast>
  <Toast class="!max-w-sm" color="red" transition={fly} params={{ x: 200 }} bind:open={$toastError}>
    <svelte:fragment slot="icon">
      <ExclamationCircleSolid class="w-5 h-5" />
    </svelte:fragment>
    {$_("sometingWentWrongPleaseTryAgain")}
  </Toast>
  <Toast
    class="!max-w-sm"
    color="green"
    transition={fly}
    params={{ x: 200 }}
    bind:open={$toastCopyToClipboard}
  >
    <svelte:fragment slot="icon">
      <CheckCircleSolid class="w-5 h-5" />
    </svelte:fragment>
    {$_("copiedToClipboard")}
  </Toast>

  {#each Object.keys($customToasts) as key}
    <Toast
      class="!max-w-sm"
      color="green"
      transition={fly}
      params={{ x: 200 }}
      bind:open={$customToasts[key]}
    >
      <svelte:fragment slot="icon">
        <CheckCircleSolid class="w-5 h-5" />
      </svelte:fragment>

      {key}
      <!-- {$_("copiedToClipboard")} -->
    </Toast>
  {/each}
</div>
