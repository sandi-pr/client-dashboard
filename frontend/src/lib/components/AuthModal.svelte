<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Dean Ramhan <dean@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import { Button, Modal, Label, Input, Checkbox } from "flowbite-svelte";
  let formModal = false;
</script>

<Button on:click={() => (formModal = true)}>Form modal</Button>

<Modal bind:open={formModal} size="xs" autoclose={false} class="w-full">
  <form class="flex flex-col space-y-6" action="#">
    <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Sign in to our platform</h3>
    <Label class="space-y-2">
      <span>Email</span>
      <Input type="email" name="email" placeholder="name@company.com" required />
    </Label>
    <Label class="space-y-2">
      <span>Your password</span>
      <Input type="password" name="password" placeholder="•••••" required />
    </Label>
    <div class="flex items-start">
      <Checkbox>Remember me</Checkbox>
      <a href="/" class="ms-auto text-sm text-primary-700 hover:underline dark:text-primary-500">
        Lost password?
      </a>
    </div>
    <Button type="submit" class="w-full1">Login to your account</Button>
    <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
      Not registered? <a href="/" class="text-primary-700 hover:underline dark:text-primary-500">
        Create account
      </a>
    </div>
  </form>
</Modal>
