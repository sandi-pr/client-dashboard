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
  import { Button, CR3Call, Typography } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

  let personsLatest = [];
  let documentsReady = $state(false);
  let actionItems = $state([]);

  let datas = $state([
    {
      productId: "19846829",
      prroduct_name: "Solid Hardwood Flooring",
      price: "$5/hr",
      category: {
        type: "brand",
        name: "Service",
      },
      // is_active
      status: {
        type: "success",
        name: "Active",
      },
      actions: {
        objects: [
          {
            href: "/crm/product_service/detail",
            name: "View",
          },
        ],
      },
    }
  ]);

  let configs = {
    search: true,
  };

  // async function getPerson() {
  //   let token = sessionStorage.getItem("token");
  //   if (token !== null) {
  //     let response = await call.call("client-dashboard", "get-persons", {
  //       URLQuery: {},
  //       URLParams: {},
  //       Params: {},
  //       Headers: {
  //         Authorization: "Bearer " + token,
  //       },
  //     });

  //     if (response.status !== 200) {
  //       console.log(response);
  //       return;
  //     }

  //     if (response.result.data.length > 0) {
  //       let promises = response.result.data.map(async (person, index) => {
  //         return {
  //           id: person.id,
  //           fullname: person.attributes.fullname,
  //           email: person.attributes.email,
  //           phone: person.attributes.phone,
  //           address: person.attributes.address,
  //           city: person.attributes.city,
  //           state: person.attributes.state,
  //           zip_code: person.attributes.zip_code,
  //           country: person.attributes.country,
  //           date_of_birth: person.attributes.date_of_birth,
  //           is_active: person.attributes.is_active,
  //           actions: {
  //             objects: [
  //               {
  //                 href:
  //                   `${await call.getPage("client-dashboard", "crm-person-detail")}` +
  //                   "?personId=" +
  //                   person.id,
  //                 name: "View",
  //               },
  //             ],
  //           },
  //         };
  //       });

  //       datas = await Promise.all(promises);
  //       documentsReady = true;
  //     }
  //   }
  // }

  onMount(async () => {
    // getActionItems();
    // getPerson();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text="Product & Service" />
          <Button text="Add new" icon={PlusIcon} spaceY="mt-0" />
        </div>

        <!-- Person Table Data -->
        <div>
          {#if documentsReady}
            <DataTable {datas} {configs} />
          {:else}
            <DataTable {datas} {configs} />
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>
