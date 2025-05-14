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
  import { Button, Input, Modal } from "flowbite-svelte";
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";
  import { SearchOutline } from "flowbite-svelte-icons";

  import { StateCard, call, check, triggerToast, isBakrie } from "$lib";
  import { Typography } from "hexacomponents-svelte";
  import DataTableUser from "$lib/components/DataTableUsersList.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import GridListPermissions from "$lib/components/GridListPermissions.svelte";
  import TabBar from "$lib/components/TabBar.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import PlusIcon from "$lib/icon/plus.svg";

  const isFrontend = process.env.NODE_ENV === "development";

  let metadata = { page: 1, pageSize: 25, pageCount: 1, total: 25 };

  let datas = $state([]);

  let editRoles = $state(false);

  let gettingDatas = $state(false);

  let userListDatas = $state([]);
  let filteredUsers = $state([]);
  let addRoleModal = $state(false);
  let tmpRoleName = $state("");

  let tabDatas = $state({
    selected: "Users List",
    data: [
      {
        name: "Users List",
      },
      !$isBakrie
        ? {
            name: "Roles",
          }
        : false,
    ].filter(Boolean),
  });

  let configs = $state({
    search: true,
    widthActions: "w-20",
    widthStatus: "w-32",
  });

  function getAuthority(authority) {
    if (authority === "Admin") {
      return {
        type: "brand",
        name: "Admin",
      };
    }

    if (authority === "User") {
      return {
        type: "brand",
        name: "User",
      };
    }

    if (authority === "Super Admin") {
      return {
        type: "brand",
        name: "Brand",
      };
    }
  }

  async function getUsers() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      gettingDatas = true;
      let response = await $call.call("client-dashboard", "get-users", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        gettingDatas = false;
        return;
      }

      let link = isFrontend
        ? "/manage_users/detail"
        : await $call.getPage("client-dashboard", "user-detail");

      if (response.result.data.length > 0) {
        let promises = response.result.data.map(async (user, index) => {
          return {
            no: index + 1,
            id: user.user_id !== null ? user.user_id : "null",
            fullname: user.fullname || "No Name",
            username: user.username,
            email: user.email,
            role: user.authority || "No Authority",
            status: user.blocked === true ? "Inactive" : "Active",
            href: link + "?userId=" + user.id,
          };
        });

        userListDatas = await Promise.all(promises);
        filteredUsers = userListDatas;
      }
      gettingDatas = false;
    }
  }

  async function getRoles() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      gettingDatas = true;
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
        gettingDatas = false;
        return;
      }

      console.log("response", response);

      if (response.result.roles.length > 0) {
        configs = {
          search: true,
          widthActions: "w-20",
          widthStatus: "w-32",
          // pagination: response.result.meta.pagination,
        };

        let link = isFrontend
          ? "/manage_users/role"
          : await $call.getPage("client-dashboard", "user-role");

        let promises = response.result.roles.map(async (role, index) => {
          return {
            id: index + 1,
            rolename: role.name,
            status: {
              type: "success",
              // name: role.is_active === true ? "Inactive" : "Active",
              name: "Active",
            },
            actions: {
              objects: [
                {
                  href: link + "?id=" + role.role,
                  name: "Detail",
                },
              ],
            },
          };
        });

        datas = await Promise.all(promises);
      }
      gettingDatas = false;
    }
  }

  const nextPage = (evt) => {
    console.log("next page");
  };

  const prevPage = (evt) => {
    console.log("prev page");
  };

  const handleSearch = (evt) => {
    console.log("clicked");
    const searchTerm = evt.target.value.toLowerCase();

    filteredUsers = userListDatas.filter((user) => {
      return (
        user.fullname.toLowerCase().includes(searchTerm) ||
        user.username.toLowerCase().includes(searchTerm) ||
        user.email.toLowerCase().includes(searchTerm) ||
        user.role.toLowerCase().includes(searchTerm) ||
        user.status.toLowerCase().includes(searchTerm)
      );
    });
  };

  onMount(async () => {
    // getActionItems();
    getUsers();
    getRoles();

    // // test
    // const test = await $call.call("module-manager", "get-roles", {
    //   URLQuery: {},
    //   URLParams: {},
    //   Params: {},
    //   Headers: {
    //     Authorization: "Bearer " + localStorage.getItem("token"),
    //   },
    // });
    // console.log("===================", test);
  });

  async function addUser() {
    goto(await $call.getPage("client-dashboard", "user-create"));
    return;
  }

  async function addRole() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await $call.call("module-manager", "create-role", {
          URLQuery: {},
          URLParams: {},
          Params: {
            role: tmpRoleName.toLowerCase().replace(/\s+/g, "_"),
            name: tmpRoleName,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        triggerToast("created");
        getRoles();
        tmpRoleName = "";
      } catch (error) {
        console.log(error);
      }
    } else {
      console.log("token not found");
    }
  }
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        {#if tabDatas.selected === "Users List"}
          <div class="flex items-center justify-between">
            <Typography preset="Display sm/Semibold" color="primary" text={$_("usersList")} />
            <Button
              class="gap-1"
              disabled={!check("general.user:create")}
              onclick={addUser}
              icon={PlusIcon}
              ><img class="w-4 h-4" src={PlusIcon} alt="plus-icon" />{$_("addNewUser")}</Button
            >
          </div>
        {:else if tabDatas.selected === "Roles"}
          <div class="flex items-center justify-between">
            <Typography preset="Display sm/Semibold" color="primary" text={$_("roles")} />
            <Button
              class="gap-1"
              disabled={!check("general.role:create")}
              onclick={() => {
                addRoleModal = true;
              }}
              icon={PlusIcon}
              ><img class="w-4 h-4" src={PlusIcon} alt="plus-icon" />{$_("addNewRoles")}</Button
            >
          </div>
        {:else}
          <!-- else content here -->
        {/if}

        <div class="mt-2 mb-4">
          <TabBar {tabDatas} />
        </div>

        <!-- Person Table Data -->
        <!-- {#if filteredUsers.length > 0} -->
        {#if tabDatas.selected === "Users List"}
          <div class="items-center">
            <Input
              id="search"
              type="search"
              class="bg-white dark:bg-hc-graydark-950 my-4"
              placeholder={$_("search")}
              on:input={handleSearch}
            >
              <SearchOutline slot="left" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
            </Input>
          </div>
          <DataTableUser
            userListDatas={filteredUsers}
            {metadata}
            next={nextPage}
            previous={prevPage}
          />
        {:else if tabDatas.selected === "Roles"}
          <!-- <GridListPermissions bind:permissionListDatas bind:edit={editRoles} /> -->
          {#if datas.length > 0}
            <DataTable bind:datas bind:configs />
          {:else if gettingDatas}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else}
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          {/if}
        {:else}
          <!-- else content here -->
        {/if}
        <!-- {:else if gettingDatas}
          <StateCard
            illustration={IllustrationProgress}
            description={$_("fetchingDataPleaseWait")}
          /> -->
        <!-- {:else}
          <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
        {/if} -->
      </div>
    </div>
  </div>
</div>

<Modal title={$_("createNewRole")} bind:open={addRoleModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("roleName")}*`} />
      <Input bind:value={tmpRoleName} placeholder={$_("inputRoleName")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={() => {
        addRole();
      }}
      disabled={!tmpRoleName}
      class="w-full">{$_("addNewRole")}</Button
    >
  </svelte:fragment>
</Modal>
