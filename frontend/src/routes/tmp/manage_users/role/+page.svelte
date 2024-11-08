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
  import { Button, Input, P } from "flowbite-svelte";
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";
  import { SearchOutline } from "flowbite-svelte-icons";

  import { StateCard, call, checkOr, InputRw } from "$lib";
  import { Typography } from "hexacomponents-svelte";
  import DataTableUser from "$lib/components/DataTableUsersList.svelte";
  import GridListPermissions from "$lib/components/GridListPermissions.svelte";
  import TabBar from "$lib/components/TabBar.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { ChevronDownOutline, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import PlusIcon from "$lib/icon/plus.svg";

  let metadata = $state({ page: 1, pageSize: 25, pageCount: 0, total: 0 });

  let role = $state(
    {
      id: 1,
      rolename: "Admin",
      status: true,
      createdDate: "04/11/2024",
    }
  );

  let editRoles = $state(true);

  let gettingDatas = $state(false);

  let userListDatas = $state([]);

  let tmpUserListDatas = $state([]);
  let dropdownOpenStatus = $state(false);

  let onEdit = $state(false);

  let permissionGeneral = $state([
    {
      id: 1,
      type: "header",
      icon: "",
      level: 0,
      title: "Actions",
      data: [{ name: "allow" }],
    },
    {
      id: 2,
      type: "section",
      icon: "userUp",
      level: 0,
      title: "Manage User",
      data: [],
    },
    {
      id: 3,
      type: "content",
      icon: "",
      level: 0,
      title: "Create new user",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 4,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete existing user",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 5,
      type: "content",
      icon: "",
      level: 0,
      title: "View user information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 6,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit user information",
      data: [
        { name: "allow", granted: false },
      ],
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
    {
      id: 2,
      type: "section",
      icon: "bookOpen",
      level: 0,
      title: "Document Library",
      data: [],
    },
    {
      id: 3,
      type: "content",
      icon: "",
      level: 0,
      title: "Create new library",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 4,
      type: "content",
      icon: "",
      level: 0,
      title: "View library",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 5,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit libraries information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 6,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete libraries",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 7,
      type: "content",
      icon: "",
      level: 0,
      title: "Create new collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 8,
      type: "content",
      icon: "",
      level: 0,
      title: "View collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 9,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit collection information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 10,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 11,
      type: "content",
      icon: "",
      level: 0,
      title: "Upload documents into collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 12,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit documents information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 13,
      type: "content",
      icon: "",
      level: 0,
      title: "View documents information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 14,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete documents from collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 15,
      type: "section",
      icon: "bookClosed",
      level: 0,
      title: "Sop Compliance",
      data: [],
    },
    {
      id: 16,
      type: "content",
      icon: "",
      level: 0,
      title: "Create new collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 17,
      type: "content",
      icon: "",
      level: 0,
      title: "View collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 18,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit collection information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 19,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 20,
      type: "content",
      icon: "",
      level: 0,
      title: "Upload documents into collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 21,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit documents information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 22,
      type: "content",
      icon: "",
      level: 0,
      title: "View documents information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 23,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete documents from collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 24,
      type: "section",
      icon: "bookClosed",
      level: 0,
      title: "Contracts Agreement",
      data: [],
    },
    {
      id: 25,
      type: "content",
      icon: "",
      level: 0,
      title: "Create new collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 26,
      type: "content",
      icon: "",
      level: 0,
      title: "View collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 27,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit collection information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 28,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 29,
      type: "content",
      icon: "",
      level: 0,
      title: "Upload documents into collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 30,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit documents information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 31,
      type: "content",
      icon: "",
      level: 0,
      title: "View documents information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 32,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete documents from collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 33,
      type: "section",
      icon: "bookClosed",
      level: 0,
      title: "Customer Data Privacy",
      data: [],
    },
    {
      id: 34,
      type: "content",
      icon: "",
      level: 0,
      title: "Create new collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 35,
      type: "content",
      icon: "",
      level: 0,
      title: "View collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 36,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit collection information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 37,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 38,
      type: "content",
      icon: "",
      level: 0,
      title: "Upload documents into collection",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 39,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit documents information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 40,
      type: "content",
      icon: "",
      level: 0,
      title: "View documents information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 41,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete documents from collection",
      data: [
        { name: "allow", granted: false },
      ],
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
    {
      id: 2,
      type: "section",
      icon: "messageChat",
      level: 0,
      title: "Hexa Bots",
      data: [],
    },
    {
      id: 3,
      type: "content",
      icon: "",
      level: 0,
      title: "Create new bots",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 4,
      type: "content",
      icon: "",
      level: 0,
      title: "Change bots information knowledge source and other configurations",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 5,
      type: "content",
      icon: "",
      level: 0,
      title: "View bots information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 6,
      type: "content",
      icon: "",
      level: 0,
      title: "View conversation logs",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 7,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete deployed bots",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 8,
      type: "section",
      icon: "fileCheck",
      level: 0,
      title: "Contract Reviewer",
      data: [],
    },
    {
      id: 9,
      type: "content",
      icon: "",
      level: 0,
      title: "Add new projects",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 10,
      type: "content",
      icon: "",
      level: 0,
      title: "Add new documents to project",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 11,
      type: "content",
      icon: "",
      level: 0,
      title: "View projects information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 12,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit projects information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 13,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete projects",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 14,
      type: "content",
      icon: "",
      level: 0,
      title: "Create new cases",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 15,
      type: "content",
      icon: "",
      level: 0,
      title: "Add delete new rules",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 16,
      type: "content",
      icon: "",
      level: 0,
      title: "Add delete referenced documents",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 17,
      type: "content",
      icon: "",
      level: 0,
      title: "View cases information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 18,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit cases rules",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 19,
      type: "section",
      icon: "data",
      level: 0,
      title: "Datasets",
      data: [],
    },
    {
      id: 20,
      type: "content",
      icon: "",
      level: 0,
      title: "Create new datasets",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 21,
      type: "content",
      icon: "",
      level: 0,
      title: "View datasets",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 22,
      type: "content",
      icon: "",
      level: 0,
      title: "Edit datasets information",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 23,
      type: "content",
      icon: "",
      level: 0,
      title: "Delete datasets",
      data: [
        { name: "allow", granted: false },
      ],
    },
    {
      id: 24,
      type: "content",
      icon: "",
      level: 0,
      title: "Process datasets",
      data: [
        { name: "allow", granted: false },
      ],
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
        name: "Documents Ai",
      },
      {
        name: "Ai Management",
      },
      {
        name: "Business Management",
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
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text="Manager" />
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
                    <div class="flex flex-col w-1/5 gap-2">
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
                          <Typography preset="Text sm/Regular" text={role.status ? $_("active") : $_("deactive")} color="secondary" />
                          <ChevronDownOutline class="w-6 h-6 ms-2 dark:text-white" />
                        </Button>
                        <div
                          id="dropdown"
                          class={`z-10 ${dropdownOpenStatus ? "fixed" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-auto dark:bg-gray-700`}
                        >
                          <ul
                            class="py-2 text-sm text-gray-700 dark:text-gray-700"
                            aria-labelledby="dropdownDefaultButton"
                          >
                            {#if role.status}
                              <button
                                class="flex p-2"
                                onclick={() => {
                                  role.status = false;
                                  dropdownOpenStatus = false;
                                }}
                              >
                                <Typography preset="Text sm/Regular" text={$_("deactive")} color="secondary" />
                              </button>
                            {:else}
                              <button
                                class="flex p-2"
                                onclick={() => {
                                  role.status = true;
                                  dropdownOpenStatus = false;
                                }}
                              >
                                <Typography preset="Text sm/Regular" text={$_("active")} color="secondary" />
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
          {:else if tabDatas.selected === "General"}
            <GridListPermissions bind:permissionListDatas={permissionGeneral} bind:edit={editRoles} />
          {:else if tabDatas.selected === "Documents Ai"}
            <GridListPermissions bind:permissionListDatas={permissionDocumentAi} bind:edit={editRoles} />
          {:else if tabDatas.selected === "Ai Management"}
            <GridListPermissions bind:permissionListDatas={permissionAiManagement} bind:edit={editRoles} />
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
