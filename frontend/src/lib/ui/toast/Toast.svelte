<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: 
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->

<script>
  import { Toast } from "flowbite-svelte";
  import { CheckCircleSolid, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { fly } from "svelte/transition";
  import { _ } from "svelte-i18n";

  let { toastCreated = false, toaastUpdated = false, toastDeleted = false, toastError = false, toastNotAvailable = false } = $props();
</script>

<div class="absolute bottom-5 z-50 end-5">
  <Toast color="green" transition={fly} params={{ x: 200 }} bind:open={toastCreated}>
    <svelte:fragment slot="icon">
      <CheckCircleSolid class="w-5 h-5" />
    </svelte:fragment>
    {$_("successfullyCreated")}
  </Toast>
  <Toast color="green" transition={fly} params={{ x: 200 }} bind:open={toaastUpdated}>
    <svelte:fragment slot="icon">
      <CheckCircleSolid class="w-5 h-5" />
    </svelte:fragment>
    {$_("successfullyUpdated")}
  </Toast>
  <Toast color="green" transition={fly} params={{ x: 200 }} bind:open={toastDeleted}>
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
    bind:open={toastNotAvailable}
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
    bind:open={toastError}
  >
    <svelte:fragment slot="icon">
      <ExclamationCircleSolid class="w-5 h-5" />
    </svelte:fragment>
    {$_("sometingWentWrongPleaseTryAgain")}
  </Toast>
</div>