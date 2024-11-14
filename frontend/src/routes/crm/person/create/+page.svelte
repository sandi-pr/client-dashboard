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
  import { P, Heading, Input, Button, Helper } from "flowbite-svelte";
  import { ToastContainer, triggerToast } from "$lib";
  import Divider from "$lib/components/Divider.svelte";
  import { _ } from "svelte-i18n";
  import { createEventDispatcher } from "svelte";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  const dispatch = createEventDispatcher();

  onMount(() => {
    if (check("businessManagement.person:create") === false) {
      console.log("no access");
      window.history.back();
      return;
    }
  });

  import { CR3Call } from "hexacomponents-svelte";

  import { call } from "$lib";

  let data = $state({
    fullname: "",
    date_of_birth: "",
    email: "",
    phone: "",
    country: "",
    address: "",
    city: "",
    state: "",
    zip_code: "",
    is_active: "Yes",
  });

  let emailValid = $state(false);

  let buttonDisable = $state(true);

  let toastCreated = $state(false);
  let toastError = $state(false);
  let counter = $state(6);

  function validateEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  function handleEmailChange(event) {
    const newEmail = event.target.value;
    emailValid = validateEmail(newEmail);
    let email = newEmail;
    dispatch("emailChange", { email, emailValid });
  }

  async function createPerson() {
    console.log(data);
    let token = localStorage.getItem("token");
    try {
      let response = await $call.call("client-dashboard", "create-person", {
        URLQuery: {},
        URLParams: {},
        Params: {
          fullname: data.fullname,
          date_of_birth: data.date_of_birth,
          email: data.email,
          phone: data.phone.toString(),
          country: data.country,
          address: data.address,
          city: data.city,
          state: data.state,
          zip_code: data.zip_code.toString(),
          is_active: data.is_active,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });
      console.log(response);
      if (response.status === 200) {
        triggerToast("created");
        // setTimeout(async () => {
        goto(await $call.getPage("client-dashboard", "crm-person"));
        // }, 3000);
        return;
      } else {
        triggerToast("error");
        console.log(response);
      }
    } catch (error) {
      triggerToast("error");
      console.log(error);
    }
  }

  async function backPage() {
    goto(await $call.getPage("client-dashboard", "crm-person"));
    return;
  }

  $effect(() => {
    buttonDisable =
      data.fullname.length < 1 ||
      data.date_of_birth.length < 1 ||
      data.email.length < 1 ||
      data.phone.length < 1 ||
      data.country.length < 1 ||
      data.address.length < 1 ||
      data.city.length < 1 ||
      data.state.length < 1 ||
      data.zip_code.length < 1 ||
      emailValid === false;
  });

  // let dropdownOpen = false;
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex items-center pb-8">
        <div class="flex flex-col gap-2 grow">
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {$_("createNewPersons")}
            </Heading>
          </div>
        </div>
        <div class="flex gap-2">
          <Button color="light" class="px-3 focus-within:!ring-1" onclick={backPage}
            >{$_("cancel")}</Button
          >
          <Button
            disabled={buttonDisable}
            onclick={createPerson}
            class={`px-3 focus-within:!ring-1 ${
              !buttonDisable ? "cursor-pointer" : "cursor-not-allowed"
            }`}>{$_("create")}</Button
          >
        </div>
      </div>

      <Divider />
      <div class="flex flex-col gap-4 pt-8">
        <!-- general information -->
        <div class="flex gap-6">
          <div class="w-60">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              {$_("generalInformation")}
            </P>
          </div>
          <div class="flex flex-col grow gap-6">
            <div class="flex gap-5">
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("fullname")}*
                </P>
                <Input bind:value={data.fullname} class="bg-primary-light dark:bg-secondary-dark" />
              </div>
              <!-- <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("personId")}
                </P>
                <Input
                  bind:value={data.personId}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div> -->
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("email")}*</P
                >
                <Input
                  bind:value={data.email}
                  on:input={handleEmailChange}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
                {#if !emailValid && data.email.length > 0}
                  <Helper class="text-sm !text-red-400 pt-2">{$_("emailNotValid")}</Helper>
                {/if}
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("phone")}*</P
                >
                <Input
                  bind:value={data.phone}
                  type="number"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="flex flex-col grow gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
              >
                {$_("address")}*
              </P>
              <Input bind:value={data.address} class="bg-primary-light dark:bg-secondary-dark" />
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("city")}*</P
                >
                <Input bind:value={data.city} class="bg-primary-light dark:bg-secondary-dark" />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("state")}*</P
                >
                <Input bind:value={data.state} class="bg-primary-light dark:bg-secondary-dark" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("zipCode")}*</P
                >
                <Input
                  bind:value={data.zip_code}
                  type="number"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("country")}*</P
                >
                <Input bind:value={data.country} class="bg-primary-light dark:bg-secondary-dark" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("dateOfBirth")}*</P
                >
                <Input
                  type="date"
                  bind:value={data.date_of_birth}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <ToastContainer />
</div>
