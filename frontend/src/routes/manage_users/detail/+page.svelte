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
  import {
    P,
    Heading,
    Input,
    Button,
    Select,
    Textarea,
    Toggle,
    Modal,
    Avatar,
  } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import { onMount } from "svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import Badge from "$lib/ui/badge/Badge.svelte";
  import { ChevronDownOutline, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { goto } from "$app/navigation";
  import { triggerToast, check, generateUniqueId } from "$lib";

  import { call } from "$lib";

  const params = new URLSearchParams(window.location.search);
  const edit = params.get("edit");
  let userId = $state("");
  let onEdit = $state(false);
  let dropdownOpenRole = $state(false);
  let dropdownOpenStatus = $state(false);
  let popupModal = $state(false);

  let data = $state({});

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
    let token = localStorage.getItem("token");
    try {
      let response = await $call.call("client-dashboard", "get-user", {
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
        // userId: generateUniqueId(response.result.user.id, "user"),
        userId: response.result.user.user_id,
        fullname: response.result.user.fullname !== null ? response.result.user.fullname : "",
        username: response.result.user.username,
        email: response.result.user.email,
        createdAt: formatDate(response.result.user.createdAt),
        role: response.result.user.authority,
        status: !response.result.user.blocked,
        company: response.result.user.company,
      };
    } catch (error) {
      console.log(error);
    }
  }

  async function updateUser() {
    let token = localStorage.getItem("token");
    console.log("send", data);
    try {
      let response;
      if (data.password === "") {
        response = await $call.call("client-dashboard", "update-user", {
          URLQuery: {},
          URLParams: { id: userId },
          Params: {
            fullname: data.fullname,
            username: data.username,
            email: data.email,
            company: data.company,
            blocked: data.status,
            authority: data.role,
            // status: data.status ? "Yes" : "No",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
      } else {
        response = await $call.call("client-dashboard", "update-user", {
          URLQuery: {},
          URLParams: { id: userId },
          Params: {
            fullname: data.fullname,
            username: data.username,
            password: data.password,
            email: data.email,
            company: data.company,
            blocked: !data.status,
            authority: data.role,
            // status: data.status ? "Yes" : "No",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
      }

      if (response.status !== 200) {
        triggerToast("error");
        console.log("status not 200", response);
        return;
      }

      triggerToast("updated");
      onEdit = false;
      getUser();

      // goto(await $call.getPage("client-dashboard", "users"));
    } catch (error) {
      console.log("error", error);
    }
  }

  async function deleteUser() {
    let token = localStorage.getItem("token");
    try {
      let response = await $call.call("client-dashboard", "delete-user", {
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

      triggerToast("deleted");
      goto(await $call.getPage("client-dashboard", "users"));
    } catch (error) {
      console.log(error);
    }
  }

  let me = $state({ id: 0 });

  onMount(async () => {
    me = JSON.parse(localStorage.getItem("user"));
    userId = params.get("userId");
    if (userId === null || userId === undefined) {
      triggerToast("error");
      goto(await $call.getPage("client-dashboard", "users"));
    }
    // if (edit === "true") onEdit = true;
    getUser();
    getRoles();
  });

  // let dropdownOpenRole = false;
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- Header -->
      <div class="flex items-center pb-8">
        <Avatar src={data.photo} alt={data.fullname} class="w-16 h-16 me-4" />
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
              onclick={async function () {
                await getUser();
                onEdit = false;
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button
              disabled={data.fullname === "" ||
                data.username === "" ||
                data.email === "" ||
                data.password === "" ||
                data.role === "" ||
                !check("general.user:update")}
              onclick={updateUser}
              class="px-3 focus-within:!ring-1">{$_("save")}</Button
            >
          {:else}
            {#if check("general.user:delete") && me?.id.toString() !== userId}
              <Button
                onclick={() => {
                  popupModal = true;
                }}
                color="red"
                class="px-3 focus-within:!ring-1">{$_("delete")}</Button
              >
            {/if}
            {#if check("general.user:update")}
              <Button
                onclick={() => {
                  onEdit = true;
                }}
                color="light"
                class="px-3 focus-within:!ring-1">{$_("edit")}</Button
              >
            {/if}
          {/if}
        </div>
      </div>

      <div
        class="mx-4 sm:mx-0 border-t-[1px] border-hc-graylight-300 dark:border-hc-graydark-700"
      ></div>
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
              <div class="flex flex-col col-span-2 gap-2">
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
              <div class="flex flex-col col-span-1 gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("User ID")}
                </P>
                <Input
                  readonly
                  bind:value={data.userId}
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
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("role")}
                </P>
                <div class="my-auto">
                  <Button
                    color="light"
                    class="py-[9px] px-3 w-full justify-between"
                    on:click={() => (dropdownOpenRole = !dropdownOpenRole)}
                    id="dropdownDefaultButton"
                  >
                    <!-- <Badge
                      text={data.role}
                      type={data.role === "Super Admin" || data.role === "Admin" ? "brand" : "info"}
                    /> -->
                    <P class="font-semibold">{data.role}</P>
                    <ChevronDownOutline class="w-6 h-6 ms-2 dark:text-white" />
                  </Button>
                  <div
                    id="dropdown"
                    class={`z-10 ${dropdownOpenRole && onEdit ? "fixed" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}
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
                      <!-- <button
                            class="flex p-2"
                            onclick={() => {
                              data.role = "Admin";
                              dropdownOpenRole = false;
                            }}
                          >
                            <P>Admin</P>
                          </button>
                          <button
                            class="flex p-2"
                            onclick={() => {
                              data.role = "Manager";
                              dropdownOpenRole = false;
                            }}
                          >
                            <P>Manager</P>
                          </button>
                        <button
                          class="flex p-2"
                          onclick={() => {
                            data.role = "User";
                            dropdownOpenRole = false;
                          }}
                        >
                          <P>User</P>
                        </button> -->
                    </ul>
                  </div>
                </div>
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("status")}</P
                >
                <div class="my-auto">
                  <Button
                    color="light"
                    class="py-[9px] px-3 w-full justify-between"
                    on:click={() => (dropdownOpenStatus = !dropdownOpenStatus)}
                    id="dropdownDefaultButton"
                  >
                    <!-- <Badge
                      text={data.status}
                      type={data.status === "Active" ? "brand" : "warning"}
                    /> -->
                    <P class="font-semibold">{data.status ? "Active" : "Deactive"}</P>
                    <ChevronDownOutline class="w-6 h-6 ms-2 dark:text-white" />
                  </Button>
                  <div
                    id="dropdown"
                    class={`z-10 ${dropdownOpenStatus && onEdit ? "fixed" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}
                  >
                    <ul
                      class="py-2 text-sm text-gray-700 dark:text-gray-200"
                      aria-labelledby="dropdownDefaultButton"
                    >
                      {#if data.status}
                        <button
                          class="flex p-2"
                          onclick={() => {
                            data.status = false;
                            dropdownOpenStatus = false;
                          }}
                        >
                          <P>Deactive</P>
                          <!-- <Badge text="Deactive" /> -->
                        </button>
                      {:else}
                        <button
                          class="flex p-2"
                          onclick={() => {
                            data.status = true;
                            dropdownOpenStatus = false;
                          }}
                        >
                          <P>Active</P>
                          <!-- <Badge text="Active" /> -->
                        </button>
                      {/if}
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteUser} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>
