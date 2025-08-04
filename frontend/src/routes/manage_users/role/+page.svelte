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
  import { onMount } from "svelte";
  import { Button, P, Modal } from "flowbite-svelte";
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";

  import { StateCard, call, check, InputRw, Toggle, generateUniqueId, triggerToast } from "$lib";
  import { Typography } from "hexacomponents-svelte";
  import GridListPermissions from "$lib/components/GridListPermissions.svelte";
  import TabBar from "$lib/components/TabBar.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { page } from "$app/stores";
  import { format } from "date-fns";

  const isFrontend = process.env.NODE_ENV === "development";

  let roleId = $state(null);

  let role = $state({
    id: 0,
    rolename: "",
    status: false,
    createdDate: "",
  });

  let editRoles = $state(true);

  let gettingDatas = $state(false);
  let tmpPermissions = $state([]);
  let disabledButton = $state(true);

  let onEdit = $state(false);
  let popupModal = $state(false);

  let permissionGeneral = $state([
    {
      id: 1,
      type: "header",
      icon: "",
      level: 0,
      title: "Actions",
      data: [{ name: "allow" }],
    },
  ]);

  let permissionAll = $state([
    {
      id: 1,
      type: "header",
      icon: "",
      level: 0,
      title: "Actions",
      data: [{ name: "allow" }],
    },
  ]);

  let tabDatas = $state({
    selected: "Information",
    data: [
      {
        name: "Information",
      },
      {
        name: "General",
      },
      {
        name: "Apps Permissions",
      },
    ],
  });

  async function getRole() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await $call.call("module-manager", "get-role", {
          URLQuery: {},
          URLParams: {
            role: roleId,
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        // console.log("role", response);

        role = {
          id: generateUniqueId(
            response.result.role?.id === undefined ? 0 : response.result.role?.id,
            "role",
          ),
          rolename: response.result.role.name,
          status: response.result.role?.status === undefined ? true : response.result.role?.status,
          createdDate: format(
            response.result.role?.created === undefined
              ? new Date()
              : new Date(response.result.role?.created),
            "dd/MM/yyyy",
          ),
        };

        let response2 = await $call.call("module-manager", "get-permissions", {
          URLQuery: {},
          URLParams: {
            role: roleId,
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          console.log(response2);
          return;
        }

        // console.log("permissions", response2);

        if (response.result.role.role === "admin") {
          editRoles = false;
        }

        permissionGeneral = [];
        permissionAll = [];
        let parentId = 0;

        for (let i = 0; i < response2.result.permissions.length; i++) {
          if (response2.result.permissions[i].key.includes("general")) {
            if (response2.result.permissions[i]?.icon) {
              permissionGeneral = [
                ...permissionGeneral,
                {
                  id: permissionGeneral.length + 1,
                  type: "section",
                  icon: "",
                  level: 0,
                  title: response2.result.permissions[i].name,
                  data: [],
                  hide: false,
                },
              ];
              parentId = permissionGeneral.length;
            } else {
              permissionGeneral = [
                ...permissionGeneral,
                {
                  id: permissionGeneral.length + 1,
                  type: "content",
                  icon: "",
                  level: 0,
                  title: response2.result.permissions[i].name,
                  data: [
                    {
                      name: "allow",
                      granted:
                        response.result.role.role === "admin"
                          ? true
                          : response.result.role.permissions.some(
                              (permission) => permission === response2.result.permissions[i].key,
                            ),
                    },
                  ],
                  key: response2.result.permissions[i].key,
                  parentId: parentId,
                  hide: false,
                },
              ];
            }
          } else {
            let newPermissionTitle = "";
            if (response2.result.permissions[i]?.icon) {
              if (response2.result.permissions[i].name === "Hexa Bots") {
                newPermissionTitle = "Configuration Bots";
              } else if (response2.result.permissions[i].name === "Datasets") {
                newPermissionTitle = "Finetunings";
              } else if (response2.result.permissions[i].name === "Document Library") {
                newPermissionTitle = "Documents";
              } else {
                newPermissionTitle = response2.result.permissions[i].name;
              }

              permissionAll = [
                ...permissionAll,
                {
                  id: permissionAll.length + 1,
                  type: "section",
                  icon: "",
                  level: 0,
                  title: newPermissionTitle,
                  data: [],
                  hide: false,
                },
              ];
              parentId = permissionAll.length;
            } else {
              if (response2.result.permissions[i].name === "Create new datasets") {
                newPermissionTitle = "Create new finetunings";
              } else if (response2.result.permissions[i].name === "View datasets") {
                newPermissionTitle = "View finetunings";
              } else if (response2.result.permissions[i].name === "Edit datasets information") {
                newPermissionTitle = "Edit finetunings information";
              } else if (response2.result.permissions[i].name === "Delete datasets") {
                newPermissionTitle = "Delete finetunings";
              } else if (response2.result.permissions[i].name === "Process datasets") {
                newPermissionTitle = "Process finetunings";
              } else {
                newPermissionTitle = response2.result.permissions[i].name;
              }

              permissionAll = [
                ...permissionAll,
                {
                  id: permissionAll.length + 1,
                  type: "content",
                  icon: "",
                  level: 0,
                  title: newPermissionTitle,
                  data: [
                    {
                      name: "allow",
                      granted:
                        response.result.role.role === "admin"
                          ? true
                          : response.result.role.permissions.some(
                              (permission) => permission === response2.result.permissions[i].key,
                            ),
                    },
                  ],
                  key: response2.result.permissions[i].key,
                  parentId: parentId,
                  hide: false,
                },
              ];
            }
          }
        }

        tmpPermissions = JSON.parse(JSON.stringify([...permissionGeneral, ...permissionAll]));
      } catch (error) {
        console.log(error);
      }
    } else {
      console.log("no token");
    }
  }

  async function deleteRole() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      let response = await $call.call("module-manager", "delete-role", {
        URLQuery: {},
        URLParams: {},
        Params: {
          role: roleId,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });
      if (response.status !== 200) {
        triggerToast("error");
        return;
      }
      triggerToast("deleted");
      goto(await $call.getPage("client-dashboard", "users"));
    } else {
      console.log("no token");
    }
  }

  async function savePermission() {
    let token = localStorage.getItem("token");

    let dataSend = "";

    let allData = [...permissionGeneral, ...permissionAll];

    const filteredPermissions = allData.filter((permission) =>
      permission.data.some((item) => item.name === "allow" && item.granted === true),
    );

    dataSend = filteredPermissions.map((permission) => permission.key).join(",");

    if (token !== null) {
      let response;
      try {
        response = await $call.call("module-manager", "set-role-permissions", {
          URLQuery: {},
          URLParams: {},
          Params: {
            role: roleId,
            permissions: dataSend,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        console.log("API response:", response);
      } catch (error) {
        console.error("API call failed:", error);
        triggerToast("error");
        return;
      }

      if (response.status !== 200) {
        console.warn("Unexpected response status:", response.status);
        triggerToast("error");
        return;
      }

      console.log("Permissions updated successfully");
      triggerToast("updated");
      getRole();
    } else {
      console.warn("No token found in localStorage");
    }
  }

  function change() {
    let allData = [...permissionGeneral, ...permissionAll];
    disabledButton = JSON.stringify(tmpPermissions) === JSON.stringify(allData);
    console.log(JSON.stringify(tmpPermissions) === JSON.stringify(allData));
    console.log(JSON.stringify(tmpPermissions[3]));
    console.log(JSON.stringify(allData[3]));
    console.log(disabledButton);
  }

  onMount(async () => {
    // getUsers();
    const unsubscribe = page.subscribe(async ($page) => {
      roleId = $page.url.searchParams.get("id");
      if (roleId === null) {
        goto(await $call.getPage("client-dashboard", "users"));
        return;
      }
      getRole();
    });
    return () => unsubscribe();
  });

  async function addUser() {
    goto(await $call.getPage("client-dashboard", "user-create"));
    return;
  }
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text={role.rolename} />
          <div class="flex gap-2">
            {#if onEdit}
              <Button
                onclick={() => {
                  onEdit = false;
                }}
                color="light"
                class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
              >
              <Button onclick={() => (onEdit = false)} class="px-3 focus-within:!ring-1"
                >{$_("save")}</Button
              >
            {:else}
              <Button
                onclick={async function () {
                  let link = isFrontend
                    ? "/manage_users"
                    : await $call.getPage("client-dashboard", "users");
                  goto(link);
                  return;
                }}
                color="light"
                class="px-3 focus-within:!ring-1">{$_("back")}</Button
              >
              {#if role.rolename !== "Admin"}
                {#if check("general.user:delete")}
                  <Button
                    onclick={() => (popupModal = true)}
                    outline
                    color="red"
                    class="px-3 focus-within:!ring-1">{$_("delete")}</Button
                  >
                {/if}
                {#if check("general.user:update")}
                  <Button
                    disabled={disabledButton}
                    onclick={() => {
                      savePermission();
                    }}
                    class="px-3 focus-within:!ring-1">{$_("save")}</Button
                  >
                {/if}
              {/if}
            {/if}
          </div>
        </div>

        <div class="mt-2 mb-4">
          <TabBar {tabDatas} />
        </div>

        <!-- Person Table Data -->
        {#if permissionGeneral.length > 0}
          {#if tabDatas.selected === "Information"}
            <div class="flex flex-col gap-4 pt-4 md:pt-6">
              <!-- general information -->
              <div class="flex flex-col md:flex-row gap-3 md:gap-6">
                <div class="w-60">
                  <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                    {$_("generalInformation")}
                  </P>
                </div>
                <div class="flex flex-col grow gap-2 md:gap-6">
                  <div class="flex gap-5">
                    <div class="flex flex-col grow gap-2">
                      <P
                        class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >
                        {$_("roleName")}
                      </P>
                      <InputRw {onEdit} bind:value={role.rolename} />
                    </div>
                    <div class="flex flex-col w-40 gap-2">
                      <P
                        class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >
                        {$_("roleId")}
                      </P>
                      <InputRw {onEdit} bind:value={role.id} />
                    </div>
                  </div>
                  <div class="flex gap-5">
                    <div class="flex flex-col w-2/5 gap-2">
                      <P
                        class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >
                        {$_("createdDate")}
                      </P>
                      <InputRw {onEdit} bind:value={role.createdDate} />
                    </div>
                    <div class="flex flex-col gap-2 w-max">
                      <P
                        class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                        >{$_("status")}</P
                      >
                      <div class="my-auto">
                        <Toggle {onEdit} bind:checked={role.status} />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          {:else if tabDatas.selected === "General"}
            <GridListPermissions
              bind:permissionListDatas={permissionGeneral}
              bind:edit={editRoles}
              bind:change
            />
          {:else}
            <GridListPermissions
              bind:permissionListDatas={permissionAll}
              bind:edit={editRoles}
              bind:change
            />
            <!-- else content here -->
          {/if}
        {:else if gettingDatas}
          <StateCard
            illustration={IllustrationProgress}
            description={$_("fetchingDataPleaseWait")}
          />
        {:else}
          <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
        {/if}
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
    <Button onclick={deleteRole} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>
