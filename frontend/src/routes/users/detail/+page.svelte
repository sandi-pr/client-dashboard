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
  import { P, Heading, Input, Button, Select, Textarea, Toggle } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import { onMount } from "svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import Badge from "$lib/ui/badge/Badge.svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  const userId = params.get("userId");
  const edit = params.get("edit");
  let onEdit = $state(false);
  let dropdownOpen = $state(false);

  let data = $state({
    fullname: "",
    username: "",
    password: "",
    email: "",
    createdAt: "",
    authority: "",
    is_active: false,
    company: 0,
  });

  function formatDate(inputDate) {
    const date = new Date(inputDate);

    const formattedDate = `${(date.getMonth() + 1).toString().padStart(2, "0")}/${date
      .getDate()
      .toString()
      .padStart(2, "0")}/${date.getFullYear()}, ${date
      .getHours()
      .toString()
      .padStart(2, "0")}:${date
      .getMinutes()
      .toString()
      .padStart(2, "0")} ${date.getHours() >= 12 ? "PM" : "AM"}`;

    return formattedDate;
  }

  async function getUser() {
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "get-user", {
        URLQuery: {},
        URLParams: { id: userId },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      console.log("ayam", response);

      data = {
        userId: response.result.id,
        fullname: response.result.fullname !== null ? response.result.fullname : "",
        username: response.result.username,
        email: response.result.email,
        createdAt: formatDate(response.result.createdAt),
        authority: response.result.authority === "Admin" ? "Admin" : "User",
        is_active: !response.result.Blocked,
        company: response.result.company,
      };
    } catch (error) {
      console.log(error);
    }
  }

  async function updateUser() {
    let token = sessionStorage.getItem("token");
    console.log("send", data);
    try {
      let response;
      if (data.password === "") {
        response = await call.call("client-dashboard", "update-user", {
          URLQuery: {},
          URLParams: { id: userId },
          Params: {
            fullname: data.fullname,
            username: data.username,
            email: data.email,
            company: data.company,
            blocked: !data.is_active,
            // authority: data.authority === "Admin" ? "admin" : "user",
            // is_active: data.is_active ? "Yes" : "No",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
      } else {
        response = await call.call("client-dashboard", "update-user", {
          URLQuery: {},
          URLParams: { id: userId },
          Params: {
            fullname: data.fullname,
            username: data.username,
            password: data.password,
            email: data.email,
            company: data.company,
            blocked: !data.is_active,
            // authority: data.authority === "Admin" ? "admin" : "user",
            // is_active: data.is_active ? "Yes" : "No",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
      }

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      await call.goToPage("client-dashboard", "users");
    } catch (error) {
      console.log(error);
    }
  }

  async function deleteUser() {
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "delete-user", {
        URLQuery: {},
        URLParams: { id: userId },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      console.log("ayam", response);
    } catch (error) {
      console.log(error);
    }
  }

  onMount(async () => {
    if (edit === "true") onEdit = true;
    getUser();
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
            <Button onclick={updateUser} class="px-3 focus-within:!ring-1">{$_("save")}</Button>
          {:else}
            <Button onclick={deleteUser} color="light" class="px-3 focus-within:!ring-1"
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
              <div class="flex flex-col col-span-4 gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("email")}
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.email}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col col-span-4 gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("password")}
                </P>
                <Input
                  type="password"
                  readonly={!onEdit}
                  bind:value={data.password}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <!-- <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("authority")}
                </P> -->
              <!-- <div class="my-auto">
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
                      <button
                        class="flex p-2"
                        onclick={() => (data.authority = "Super Admin")((dropdownOpen = false))}
                      >
                        <Badge text="Super Admin" />
                      </button>
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
                    </ul>
                  </div>
                </div> -->
              <!-- <div class="my-auto">
                    <Badge text={data.authority} />
                  </div> -->
              <!-- </div> -->
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
