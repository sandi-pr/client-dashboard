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
  import { P, Heading, Input, Button, Select } from "flowbite-svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  const personId = params.get("personId");
  let onEdit = $state(false);

  let data = $state({
    personId: "",
    fullname: "",
    date_of_birth: "",
    email: "",
    phone: "",
    country: "",
    address: "",
    city: "",
    state: "",
    zip_code: "",
    createdAt: "",
    is_active: "",
  });

  async function getPerson() {
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "get-person", {
        URLQuery: {},
        URLParams: { id: personId },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status === 200) {
        data = {
          personId: response.result.data.id,
          fullname: response.result.data.attributes.fullname,
          date_of_birth: response.result.data.attributes.date_of_birth,
          email: response.result.data.attributes.email,
          phone: response.result.data.attributes.phone,
          country: response.result.data.attributes.country,
          address: response.result.data.attributes.address,
          city: response.result.data.attributes.city,
          state: response.result.data.attributes.state,
          zip_code: response.result.data.attributes.zip_code,
          createdAt: response.result.data.attributes.createdAt,
          is_active: response.result.data.attributes.is_active,
        };
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function updatePerson() {
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "update-person", {
        URLQuery: {},
        URLParams: { id: personId },
        Params: {
          fullname: data.fullname,
          date_of_birth: data.date_of_birth,
          email: data.email,
          phone: data.phone,
          country: data.country,
          address: data.address,
          city: data.city,
          state: data.state,
          zip_code: data.zip_code,
          is_active: data.is_active,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status === 200) {
        await call.goToPage("client-dashboard", "crm-person");
        return;
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function deletePerson() {
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "delete-person", {
        URLQuery: {},
        URLParams: { id: personId },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });
      if (response.status === 200) {
        await call.goToPage("client-dashboard", "crm-person");
        return;
      }
    } catch (error) {
      console.log(error);
    }
  }

  onMount(async () => {
    getPerson();
  });

  // let dropdownOpen = false;
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex items-center pb-8">
        <div class="flex flex-col gap-2 grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {data.fullname}
              <Button
                outline
                color="red"
                class="hover:!bg-transparent !border-hc-graylight-200 dark:!border-hc-graydark-600 focus-within:!ring-1 !p-2"
              >
                <CopyIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
              </Button>
            </Heading>
          </div>
          <P
            class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
          >
            {$_("createOn")}{` ${data.createdAt}`}
          </P>
        </div>
        <div class="flex gap-2">
          {#if onEdit}
            <Button
              onclick={() => {
                onEdit = false;
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button onclick={updatePerson} class="px-3 focus-within:!ring-1">{$_("save")}</Button>
          {:else}
            <Button onclick={deletePerson} color="light" class="px-3 focus-within:!ring-1"
              >{$_("delete")}</Button
            >
            <Button
              onclick={() => {
                onEdit = true;
              }}
              class="px-3 focus-within:!ring-1">{$_("edit")}</Button
            >
          {/if}
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
                  {$_("fullname")}
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.fullname}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("personId")}
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.personId}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("email")}</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.email}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("phone")}</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.phone}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="flex flex-col grow gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
              >
                {$_("address")}
              </P>
              <Input
                readonly={!onEdit}
                bind:value={data.address}
                class="bg-primary-light dark:bg-secondary-dark"
              />
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("city")}</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.city}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("state")}</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.state}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("zipCode")}</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.zip_code}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("country")}</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.country}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("dateOfBirth")}</P
                >
                <Input
                  readonly={!onEdit}
                  type="date"
                  bind:value={data.date_of_birth}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("active")}</P
                >
                <!-- <Input readonly={!onEdit} bind:value={data.active} /> -->
                <Select
                  readonly={!onEdit}
                  class="bg-primary-light dark:bg-secondary-dark"
                  items={[
                    { value: "Yes", name: "Yes" },
                    { value: "No", name: "No" },
                  ]}
                  bind:value={data.is_active}
                />
                <!-- <Dropdown
                    placement="bottom-start"
                    bind:open={dropdownOpen}
                    class="w-48 overflow-y-auto py-1 dark:bg-secondary-dark"
                  >
                    <DropdownItem on:click={() => (dropdownOpen = false)}>true</DropdownItem>
                    <DropdownItem on:click={() => (dropdownOpen = false)}>false</DropdownItem>
                  </Dropdown> -->
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
