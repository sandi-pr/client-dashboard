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
  import { Button, Input, P, Modal } from "flowbite-svelte";
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";

  import { StateCard, call, check, InputRw, Toggle, generateUniqueId, triggerToast } from "$lib";
  import { Typography } from "hexacomponents-svelte";
  import DataTableUser from "$lib/components/DataTableUsersList.svelte";
  import GridListPermissions from "$lib/components/GridListPermissions.svelte";
  import TabBar from "$lib/components/TabBar.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { page } from "$app/stores";
  import { format } from "date-fns";

  let metadata = $state({ page: 1, pageSize: 25, pageCount: 0, total: 0 });

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

  let userListDatas = $state([]);

  let tmpUserListDatas = $state([]);
  let dropdownOpenStatus = $state(false);

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

  let permissionDocumentAi = $state([
    {
      id: 1,
      type: "header",
      icon: "",
      level: 0,
      title: "Actions",
      data: [{ name: "allow" }],
    },
  ]);

  let permissionAiManagement = $state([
    {
      id: 1,
      type: "header",
      icon: "",
      level: 0,
      title: "Actions",
      data: [{ name: "allow" }],
    },
  ]);

  let permissionBusinessManagement = $state([
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
      // {
      //   name: "Documents Ai",
      // },
      // {
      //   name: "Ai Management",
      // },
      // {
      //   name: "Business Management",
      // },
      {
        name: "Apps Permissions",
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

        console.log("response", response);

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

        console.log("response2", response2);

        // {
        //   id: 18,
        //   type: "section",
        //   icon: "box",
        //   level: 0,
        //   title: "Products Services",
        //   data: [],
        // },

        // {
        //   id: 33,
        //   type: "content",
        //   icon: "",
        //   level: 0,
        //   title: "Change availability slots",
        //   data: [{ name: "allow", granted: false }],
        // },

        if (response.result.role.role === "admin") {
          editRoles = false;
        }

        permissionGeneral = [];
        permissionAll = [];
        let parentId = 0;
        // permissionDocumentAi = [];
        // permissionAiManagement = [];
        // permissionBusinessManagement = [];

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
          }
          // else if (response2.result.permissions[i].key.includes("documentAi")) {
          //   if (response2.result.permissions[i]?.icon) {
          //     permissionDocumentAi = [
          //       ...permissionDocumentAi,
          //       {
          //         id: permissionDocumentAi.length + 1,
          //         type: "section",
          //         icon: "",
          //         level: 0,
          //         title: response2.result.permissions[i].name,
          //         data: [],
          //       },
          //     ];
          //   } else {
          //     permissionDocumentAi = [
          //       ...permissionDocumentAi,
          //       {
          //         id: permissionDocumentAi.length + 1,
          //         type: "content",
          //         icon: "",
          //         level: 0,
          //         title: response2.result.permissions[i].name,
          //         data: [
          //           {
          //             name: "allow",
          //             granted:
          //               response.result.role.role === "admin"
          //                 ? true
          //                 : response.result.role.permissions.some(
          //                     (permission) => permission === response2.result.permissions[i].key,
          //                   ),
          //           },
          //         ],
          //         key: response2.result.permissions[i].key,
          //       },
          //     ];
          //   }
          // } else if (response2.result.permissions[i].key.includes("aiManagement")) {
          //   if (response2.result.permissions[i]?.icon) {
          //     permissionAiManagement = [
          //       ...permissionAiManagement,
          //       {
          //         id: permissionAiManagement.length + 1,
          //         type: "section",
          //         icon: "",
          //         level: 0,
          //         title: response2.result.permissions[i].name,
          //         data: [],
          //       },
          //     ];
          //   } else {
          //     permissionAiManagement = [
          //       ...permissionAiManagement,
          //       {
          //         id: permissionAiManagement.length + 1,
          //         type: "content",
          //         icon: "",
          //         level: 0,
          //         title: response2.result.permissions[i].name,
          //         data: [
          //           {
          //             name: "allow",
          //             granted:
          //               response.result.role.role === "admin"
          //                 ? true
          //                 : response.result.role.permissions.some(
          //                     (permission) => permission === response2.result.permissions[i].key,
          //                   ),
          //           },
          //         ],
          //         key: response2.result.permissions[i].key,
          //       },
          //     ];
          //   }
          // } else if (response2.result.permissions[i].key.includes("businessManagement")) {
          //   if (response2.result.permissions[i]?.icon) {
          //     permissionBusinessManagement = [
          //       ...permissionBusinessManagement,
          //       {
          //         id: permissionBusinessManagement.length + 1,
          //         type: "section",
          //         icon: "",
          //         level: 0,
          //         title: response2.result.permissions[i].name,
          //         data: [],
          //       },
          //     ];
          //   } else {
          //     permissionBusinessManagement = [
          //       ...permissionBusinessManagement,
          //       {
          //         id: permissionBusinessManagement.length + 1,
          //         type: "content",
          //         icon: "",
          //         level: 0,
          //         title: response2.result.permissions[i].name,
          //         data: [
          //           {
          //             name: "allow",
          //             granted:
          //               response.result.role.role === "admin"
          //                 ? true
          //                 : response.result.role.permissions.some(
          //                     (permission) => permission === response2.result.permissions[i].key,
          //                   ),
          //           },
          //         ],
          //         key: response2.result.permissions[i].key,
          //       },
          //     ];
          //   }
          // }
          else {
            if (response2.result.permissions[i]?.icon) {
              permissionAll = [
                ...permissionAll,
                {
                  id: permissionAll.length + 1,
                  type: "section",
                  icon: "",
                  level: 0,
                  title:
                    response2.result.permissions[i].name === "Hexa Bots"
                      ? "Configuration Bots"
                      : response2.result.permissions[i].name === "Datasets"
                        ? "Configuration Datasources"
                        : response2.result.permissions[i].name,
                  data: [],
                  hide: false,
                },
              ];
              parentId = permissionAll.length;
            } else {
              permissionAll = [
                ...permissionAll,
                {
                  id: permissionAll.length + 1,
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
          }
        }

        tmpPermissions = JSON.parse(
          JSON.stringify([
            ...permissionGeneral,
            ...permissionAll,
            // ...permissionDocumentAi,
            // ...permissionAiManagement,
            // ...permissionBusinessManagement,
          ]),
        );

        // console.log("permissionGeneral", permissionGeneral);
        // console.log("permissionDocumentAi", permissionDocumentAi);
        // console.log("permissionAiManagement", permissionAiManagement);
        // console.log("permissionBusinessManagement", permissionBusinessManagement);
      } catch (error) {
        console.log(error);
      }
    } else {
      console.log("no token");
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
          };
        });

        tmpUserListDatas = await Promise.all(promises);

        metadata.total = tmpUserListDatas.length;
        metadata.pageCount = Math.ceil(metadata.total / metadata.pageSize);

        let start = (metadata.page - 1) * metadata.pageSize;
        let end = metadata.page * metadata.pageSize;
        userListDatas = tmpUserListDatas.slice(start, end);
      }
      gettingDatas = false;
    }
  }

  const nextPage = () => {
    if (metadata.page < metadata.pageCount) {
      metadata.page++;
      let start = (metadata.page - 1) * metadata.pageSize;
      let end = metadata.page * metadata.pageSize;
      userListDatas = tmpUserListDatas.slice(start, end);
    }
  };

  const prevPage = () => {
    if (metadata.page > 1) {
      metadata.page--;
      let start = (metadata.page - 1) * metadata.pageSize;
      let end = metadata.page * metadata.pageSize;
      userListDatas = tmpUserListDatas.slice(start, end);
    }
  };

  const handleSearch = (evt) => {
    const searchTerm = evt.target.value;
    if (searchTerm.trim() === "") {
      // Jika pencarian kosong, tampilkan data pada halaman saat ini
      let start = (metadata.page - 1) * metadata.pageSize;
      let end = metadata.page * metadata.pageSize;
      userListDatas = tmpUserListDatas.slice(start, end);

      metadata.total = tmpUserListDatas.length;
      metadata.pageCount = Math.ceil(metadata.total / metadata.pageSize);
    } else {
      // Filter data berdasarkan term pencarian
      userListDatas = tmpUserListDatas.filter(
        (user) =>
          user.fullname.toLowerCase().includes(searchTerm.toLowerCase()) ||
          user.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
          user.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
          user.role.toLowerCase().includes(searchTerm.toLowerCase()),
      );

      // Update metadata untuk data yang difilter
      metadata.total = userListDatas.length;
      metadata.pageCount = Math.ceil(metadata.total / metadata.pageSize);

      // Reset ke halaman pertama setelah pencarian
      metadata.page = 1;
      userListDatas = userListDatas.slice(0, metadata.pageSize);
    }
  };

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

    let allData = [
      ...permissionGeneral,
      ...permissionAll,
      // ...permissionDocumentAi,
      // ...permissionAiManagement,
      // ...permissionBusinessManagement,
    ];

    dataSend = allData
      .filter((permission) =>
        permission.data.some((item) => item.name === "allow" && item.granted === true),
      )
      .map((permission) => permission.key)
      .join(",");

    if (token !== null) {
      let response = await $call.call("module-manager", "set-role-permissions", {
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
      if (response.status !== 200) {
        triggerToast("error");
        return;
      }
      triggerToast("updated");
      getRole();
    } else {
      console.log("no token");
    }
  }

  function change() {
    let allData = [
      ...permissionGeneral,
      ...permissionAll,
      // ...permissionDocumentAi,
      // ...permissionAiManagement,
      // ...permissionBusinessManagement,
    ];
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
                  goto(await $call.getPage("client-dashboard", "users"));
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
                      // onEdit = true;
                      // console.log("save", permissionGeneral);
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
