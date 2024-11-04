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
  import { Button, Input } from "flowbite-svelte";
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";
  import { SearchOutline } from "flowbite-svelte-icons";

  import { StateCard, call, checkOr } from "$lib";
  import { Typography } from "hexacomponents-svelte";
  import DataTableUser from "$lib/components/DataTableUsersList.svelte";
  import GridListPermissions from "$lib/components/GridListPermissions.svelte";
  import TabBar from "$lib/components/TabBar.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import PlusIcon from "$lib/icon/plus.svg";

  let metadata = { page: 1, pageSize: 25, pageCount: 1, total: 25 };

  let datas = $state([]);

  let editRoles = $state(false);

  let gettingDatas = $state(false);

  let userListDatas = $state([]);

  let permissionListDatas = $state([
    {
      id: 1,
      type: "header",
      icon: "",
      level: 0,
      title: "Actions",
      data: [{ name: "User" }, { name: "Manager" }, { name: "Admin" }],
    },
    {
      id: 2,
      type: "section",
      icon: "book",
      level: 0,
      title: "Documents AI",
      data: [],
    },

    {
      id: 3,
      type: "content",
      icon: "",
      level: 0,
      title: "Create new document collection",
      data: [
        { name: "user", granted: false },
        { name: "manager", granted: true },
        { name: "admin", granted: true },
      ],
    },
    {
      id: 4,
      type: "content",
      icon: "",
      level: 0,
      title: "Read document collections",
      data: [
        { name: "user", granted: true },
        { name: "manager", granted: true },
        { name: "admin", granted: true },
      ],
    },
    {
      id: 5,
      type: "content",
      icon: "",
      level: 0,
      title: "Update document collection",
      data: [
        { name: "user", granted: false },
        { name: "manager", granted: true },
        { name: "admin", granted: true },
      ],
    },
    {
      id: 6,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete document collection",
      data: [
        { name: "user", granted: false },
        { name: "manager", granted: false },
        { name: "admin", granted: true },
      ],
    },
    {
      id: 7,
      type: "content",
      icon: "",
      level: 0,
      title: "Create new document inside document collection",
      data: [
        { name: "user", granted: false },
        { name: "manager", granted: true },
        { name: "admin", granted: true },
      ],
    },
    {
      id: 8,
      type: "content",
      icon: "",
      level: 0,
      title: "Read documents inside document collection",
      data: [
        { name: "user", granted: true },
        { name: "manager", granted: true },
        { name: "admin", granted: true },
      ],
    },
    {
      id: 9,
      type: "content",
      icon: "",
      level: 0,
      title: "Update document inside document collection",
      data: [
        { name: "user", granted: false },
        { name: "manager", granted: true },
        { name: "admin", granted: true },
      ],
    },
    {
      id: 10,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete document inside document collection",
      data: [
        { name: "user", granted: false },
        { name: "manager", granted: false },
        { name: "admin", granted: true },
      ],
    },
  ]);

  let tabDatas = $state({
    selected: "Users List",
    data: [
      {
        name: "Users List",
      },
      {
        name: "Roles & Permissions",
      },
    ],
  });

  let configs = {
    search: true,
  };

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

      if (response.result.data.length > 0) {
        configs = {
          search: true,
          // pagination: response.result.meta.pagination,
        };
        let promises = response.result.data.map(async (user, index) => {
          return {
            id: user.id,
            fullname: user.fullname || "No Name",
            username: user.username,
            email: user.email,
            role: user.authority || "No Authority",
            status: user.blocked === true ? "Inactive" : "Active",
            href:
              `${await $call.getPage("client-dashboard", "user-detail")}` + "?userId=" + user.id,
            // actions: {
            //   objects: [
            //     {
            //       href:
            //         `${await $call.getPage("client-dashboard", "user-detail")}` +
            //         "?userId=" +
            //         user.id +
            //         "&edit=" +
            //         false,
            //       name: "View",
            //     },
            //     {
            //       href:
            //         `${await $call.getPage("client-dashboard", "user-detail")}` +
            //         "?userId=" +
            //         user.id +
            //         "&edit=" +
            //         true,
            //       name: "Edit",
            //     },
            //   ],
            // },
          };
        });

        userListDatas = await Promise.all(promises);
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
  };

  onMount(async () => {
    // getActionItems();
    getUsers();

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
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        {#if tabDatas.selected === "Users List"}
          <div class="flex items-center justify-between">
            <Typography preset="Display sm/Semibold" color="primary" text="Users List" />
            <Button disabled={!checkOr(["Admin", "Manager"])} onclick={addUser} icon={PlusIcon}
              >{$_("addNewUser")}<img class="w-4 h-4" src={PlusIcon} alt="plus-icon" /></Button
            >
          </div>
        {:else if tabDatas.selected === "Roles & Permissions"}
          <div class="flex items-center justify-between">
            <Typography preset="Display sm/Semibold" color="primary" text="Roles & Permissions" />
            <!-- <div class="flex gap-3">
              {#if !editRoles}
                <Button
                  outline
                  color="light"
                  onclick={() => {
                    editRoles = true;
                  }}>{$_("edit")}</Button
                >
              {:else}
                <Button
                  onclick={() => {
                    console.log(permissionListDatas);
                    // editRoles = false;
                  }}>{$_("save")}</Button
                >
              {/if}
            </div> -->
          </div>
        {:else}
          <!-- else content here -->
        {/if}

        <div class="mt-2 mb-4">
          <TabBar {tabDatas} />
        </div>

        <!-- Person Table Data -->
        {#if userListDatas.length > 0 || permissionListDatas.length > 0}
          {#if tabDatas.selected === "Users List"}
            <div class="items-center">
              <Input
                id="search"
                type="search"
                class="bg-white dark:bg-hc-graydark-950 my-4"
                placeholder={$_("search")}
                on:change={handleSearch}
              >
                <SearchOutline slot="left" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
              </Input>
            </div>
            <DataTableUser {userListDatas} {metadata} {nextPage} {prevPage} />
          {:else if tabDatas.selected === "Roles & Permissions"}
            <GridListPermissions bind:permissionListDatas bind:edit={editRoles} />
          {:else}
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
