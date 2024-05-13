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
  import { P, Card, Avatar, Heading, A } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import UserPlusIcon from "$lib/icon/user-plus-02.svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import Chart from "$lib/components/Chart.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import SectionHeader from "$lib/ui/header/SectionHeader.svelte";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { Button } from "flowbite-svelte";
  import { StateCard } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";

  const call = new CR3Call(globalThis.location.origin);

  let personsLatest = [];
  let documentsReady = $state(false);
  let actionItems = $state([]);

  let pagination = { page: 1, pageSize: 25, pageCount: 1, total: 25 };

  let datas = $state([]);

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
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-users", {
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
      console.log("ayam", response);

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
            // authority: getAuthority(user.authority) || "No Authority",
            // status: user.is_active === "Yes" ? "Admin" : "User",
            actions: {
              objects: [
                {
                  href:
                    `${await call.getPage("client-dashboard", "user-detail")}` +
                    "?userId=" +
                    user.id +
                    "&edit=" +
                    false,
                  name: "View",
                },
                {
                  href:
                    `${await call.getPage("client-dashboard", "user-detail")}` +
                    "?userId=" +
                    user.id +
                    "&edit=" +
                    true,
                  name: "Edit",
                },
              ],
            },
          };
        });

        datas = await Promise.all(promises);
        documentsReady = true;
      }
    }
  }

  onMount(async () => {
    // getActionItems();
    getUsers();
  });

  async function addUser() {
    await call.goToPage("client-dashboard", "user-create");
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
          <Typography preset="Display sm/Semibold" color="primary" text="Users List" />
          <Button onclick={addUser} icon={PlusIcon}
            >Add new <img class="w-4 h-4" src={PlusIcon} alt="plus-icon" /></Button
          >
        </div>

        <!-- Person Table Data -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs />
        {:else}
          <div class="pt-8">
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
