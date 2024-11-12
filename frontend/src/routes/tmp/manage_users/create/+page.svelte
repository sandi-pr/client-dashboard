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
  import { triggerToast, check } from "$lib";

  import { call } from "$lib";
  import { da } from "date-fns/locale";

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
    role: check("general.user:create") ? "User" : "",
    is_active: true,
  });

  let roles = $state([]);

  async function getRoles() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      let response = await $call.call("module-manager", "get-roles", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      console.log("response", response);

      if (response.result.roles.length > 0) {
        let promises = response.result.roles.map(async (role) => {
          return {
            name: role.name,
          };
        });

        roles = await Promise.all(promises);
      }
    }
  }

  function deleteResource(id) {
    // resources = resources.filter(resource => resource.id !== id);
  }

  async function backPage() {
    goto(await $call.getPage("client-dashboard", "users"));
    return;
  }

  async function createPerson() {
    console.log(data);
    let token = localStorage.getItem("token");
    try {
      let response = await $call.call("client-dashboard", "create-user", {
        URLQuery: {},
        URLParams: {},
        Params: {
          fullname: data.fullname,
          email: data.email,
          username: data.username,
          password: data.password,
          authority: data.role,
          blocked: !data.is_active,
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
        triggerToast("error");
        console.log("status not 200", response);
      }
    } catch (error) {
      console.log("error", error);
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

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
  };

  onMount(async () => {
    // getPerson();
    if (!check("general.user:create")) {
      goto(await $call.getPage("client-dashboard", "users"));
    } else {
      getRoles();
    }
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
            <Button
              disabled={data.fullname === "" ||
                data.username === "" ||
                data.email === "" ||
                data.password === "" ||
                data.role === "" ||
                !check("general.user:create")}
              onclick={createPerson}
              class="px-3 focus-within:!ring-1">{$_("create")}</Button
            >
          {/if}
        </div>
      </div>
      <div class="flex flex-col gap-4 pt-8">
        <!-- general information -->
        <div class="flex flex-col md:flex-row gap-6">
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
              <div class="flex flex-col col-span-2 lg:col-span-1 gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("role")}
                </P>
                <div class="my-auto w-full" onfocusout={handleDropdownFocusLoss}>
                  <Button
                    color="light"
                    class="py-[9px] px-3 w-full justify-between"
                    on:click={() => (dropdownOpen = !dropdownOpen)}
                    id="dropdownDefaultButton"
                  >
                    <Typography preset="Text sm/Regular" text={$_(data.role)} color="secondary" />
                    <!-- <Badge
                      text={data.role}
                      type={data.role === "Super Admin" ||
                      data.role === "Admin" ||
                      data.role === "Manager"
                        ? "brand"
                        : "info"}
                    /> -->
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
                      {#each roles as role}
                        <button
                          class="flex p-2"
                          onclick={() => {
                            data.role = role.name;
                            dropdownOpen = false;
                          }}
                        >
                          <Typography preset="Text sm/Regular" text={role.name} color="secondary" />
                        </button>
                      {/each}
                    </ul>
                  </div>
                </div>
                <!-- <div class="my-auto">
                    <Badge text={data.role} />
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
