<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: 
  - Azzis Arswendo <azzis@hexacode.org>
  - Dean Ramhan <dean@hexacode.org>
  - Taqie Fadlillah <taqief@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->
<script>
  import { P, Heading, Input, Button, Toggle, Helper } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import { onMount, createEventDispatcher } from "svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import Badge from "$lib/ui/badge/Badge.svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import { goto } from "$app/navigation";
  import { triggerToast } from "$lib";

  import { call } from "$lib";

  const params = new URLSearchParams(window.location.search);
  const personId = params.get("personId");
  let onEdit = $state(true);
  let dropdownOpen = $state(false);

  let data = $state({
    fullname: "",
    username: "",
    email: "",
    password: "",
    createdAt: "",
    authority: "User",
    is_active: true,
  });

  function deleteResource(id) {
    // resources = resources.filter(resource => resource.id !== id);
  }

  async function backPage() {
    goto(await $call.getPage("client-dashboard", "users"));
    return;
  }

  async function createPerson() {
    console.log(data);
    let token = sessionStorage.getItem("token");
    try {
      let response = await $call.call("client-dashboard", "create-user", {
        URLQuery: {},
        URLParams: {},
        Params: {
          fullname: data.fullname,
          email: data.email,
          username: data.username,
          password: data.password,
          authority: data.authority,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });
      console.log("ayam", response);
      if (response.status === 200) {
        triggerToast("created");
        goto(await $call.getPage("client-dashboard", "users"));
        return;
      } else {
        console.log(response);
      }
    } catch (error) {
      console.log(error);
    }
  }

  const dispatch = createEventDispatcher();
  let emailValid = $state(false);
  let passwordValid = $state(false);

  function validateEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  function validatePassword(password) {
    const passwordRegex = /^(?=.*[!@#$%^&*])(?=.*[A-Z])(?=.*\d).{8,}$/;
    return passwordRegex.test(password);
  }

  function handleEmailChange(event) {
    const newEmail = event.target.value;
    emailValid = validateEmail(newEmail);
    dispatch("emailChange", { email: newEmail, emailValid });
  }

  function handlePasswordChange(event) {
    const newPassword = event.target.value;
    passwordValid = validatePassword(newPassword);
    dispatch("passwordChange", { password: newPassword, passwordValid });
  }

  onMount(async () => {
    // getPerson();
  });

  // let dropdownOpen = false;
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- Header -->
      <div class="flex items-center pb-8">
        <div class="flex flex-col gap-2 grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              Create New Users
            </Heading>
          </div>
        </div>
        <div class="flex gap-2">
          {#if onEdit}
            <Button onclick={backPage} color="light" class="px-3 focus-within:!ring-1"
              >{$_("cancel")}</Button
            >
            <Button onclick={createPerson} class="px-3 focus-within:!ring-1">{$_("create")}</Button>
          {/if}
        </div>
      </div>
      <div class="flex flex-col gap-4 pt-8">
        <!-- general information -->
        <div class="flex gap-6">
          <div class="w-60">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              {$_("generalInformation")}
            </P>
          </div>
          <div class="flex flex-col grow gap-6">
            <div class="grid grid-cols-6 gap-5">
              <div class="flex flex-col col-span-3 gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("fullname")}
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.fullname}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col col-span-3 gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("username")}
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.username}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col col-span-3 gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("email")}
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.email}
                  on:input={handleEmailChange}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
                {#if !emailValid && data.email.length > 0}
                  <Helper class="text-sm !text-red-400 pt-2">{$_("emailNotValid")}</Helper>
                {/if}
              </div>
              <div class="flex flex-col col-span-3 gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("password")}
                </P>
                <Input
                  type="password"
                  readonly={!onEdit}
                  bind:value={data.password}
                  on:input={handlePasswordChange}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
                {#if !passwordValid && data.password.length > 0}
                  <Helper class="text-sm !text-red-400 pt-2">{$_("passwordInfo")}</Helper>
                {/if}
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("authority")}
                </P>
                <div class="my-auto">
                  <Button
                    color="light"
                    class="py-[9px] px-3"
                    on:click={() => (dropdownOpen = !dropdownOpen)}
                    id="dropdownDefaultButton"
                  >
                    <Badge
                      text={data.authority}
                      type={data.authority === "Super Admin" || data.authority === "Admin"
                        ? "brand"
                        : "info"}
                    />
                    <ChevronDownOutline class="w-3 h-3 ms-2 dark:text-white" />
                  </Button>
                  <div
                    id="dropdown"
                    class={`z-10 ${dropdownOpen && onEdit ? "fixed" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}
                  >
                    <ul
                      class="py-2 text-sm text-gray-700 dark:text-gray-200"
                      aria-labelledby="dropdownDefaultButton"
                    >
                      <!-- <button
                        class="flex p-2"
                        onclick={() => (data.authority = "Super Admin")((dropdownOpen = false))}
                      >
                        <Badge text="Super Admin" type="info" />
                      </button> -->
                      <button
                        class="flex p-2"
                        onclick={() => (data.authority = "Admin")((dropdownOpen = false))}
                      >
                        <Badge text="Admin" />
                      </button>
                      <button
                        class="flex p-2"
                        onclick={() => (data.authority = "Employee")((dropdownOpen = false))}
                      >
                        <Badge text="Employee" type="info" />
                      </button>
                      <button
                        class="flex p-2"
                        onclick={() => (data.authority = "User")((dropdownOpen = false))}
                      >
                        <Badge text="User" type="info" />
                      </button>
                    </ul>
                  </div>
                </div>
                <!-- <div class="my-auto">
                    <Badge text={data.authority} />
                  </div>  -->
              </div>
              <div class="flex flex-col gap-2 w-max">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("active")}</P
                >
                <div class="my-auto">
                  <Toggle bind:checked={data.is_active}></Toggle>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
