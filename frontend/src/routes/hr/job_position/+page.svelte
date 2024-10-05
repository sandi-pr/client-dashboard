<script>
  import { onMount } from "svelte";
  import { Heading } from "flowbite-svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { Button } from "flowbite-svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import { _ } from "svelte-i18n";
  import { StateCard } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { goto } from "$app/navigation";

  import { call } from "$lib";

  let pagination = $state({ page: 1, pageSize: 5, pageCount: 1, total: 0 });

  let datas = $state([]);

  let configs = $state({});

  async function getPositions(page) {
    // if (page > pagination.pageCount) return;
    // let token = localStorage.getItem("token");
    // if (token !== null) {
    //   let response = await $call.call("client-dashboard", "get-positions", {
    //     URLQuery: {},
    //     URLParams: {},
    //     Params: {
    //       page: page,
    //       page_size: pagination.pageSize,
    //     },
    //     Headers: {
    //       Authorization: "Bearer " + token,
    //     },
    //   });

    //   if (response.status !== 200) {
    //     console.log(response);
    //     return;
    //   }

    //   if (response.result.data.length > 0) {
    //     pagination = response.result.meta.pagination;
    //     configs = {
    //       search: true,
    //       pagination: response.result.meta.pagination,
    //     };
    //     let promises2 = response.result.data.map(async (position, index) => {
    //       return {
    //         id: position.id,
    //         join_date: position.attributes.start_date,
    //         employee_name: position.attributes.employee_name,
    //         email: position.attributes.email,
    //         job_position: position.attributes.job_position,
    //         employment_status: position.attributes.employment_status,
    //         actions: {
    //           objects: [
    //             {
    //               href:
    //                 `${await $call.getPage("client-dashboard", "hr-position-detail")}` +
    //                 `?employeeId=${position.id}&edit=false`,
    //               name: "View",
    //             },
    //             {
    //               href:
    //                 `${await $call.getPage("client-dashboard", "hr-position-detail")}` +
    //                 "?employeeId=" +
    //                 position.id +
    //                 "&edit=true",
    //               name: "Edit",
    //             },
    //           ],
    //         },
    //       };
    //     });

    //     datas = await Promise.all(promises2);
    //   }
    // }
    datas = [
      {
        id: 1,
        position_id: 1,
        job_position: "Software Engineer",
        employment_status: "Full Time",
        actions: {
          objects: [
            {
              href: `localhost:3000/hr/job_position/detail?positionId=1&edit=false`,
              name: "View",
            },
            {
              href: `localhost:3000/hr/job_position/detail?positionId=1&edit=true`,
              name: "Edit",
            },
          ],
        },
      },
      {
        id: 2,
        position_id: 2,
        job_position: "Software Engineer",
        employment_status: "Part Time",
        actions: {
          objects: [
            {
              href: `localhost:3000/hr/job_position/detail?positionId=2&edit=false`,
              name: "View",
            },
            {
              href: `localhost:3000/hr/job_position/detail?positionId=2&edit=true`,
              name: "Edit",
            },
          ],
        },
      },
      {
        id: 3,
        position_id: 3,
        job_position: "Marketing Manager",
        employment_status: "Contractual",
        actions: {
          objects: [
            {
              href: `localhost:3000/hr/job_position/detail?positionId=3&edit=false`,
              name: "View",
            },
            {
              href: `localhost:3000/hr/job_position/detail?positionId=3&edit=true`,
              name: "Edit",
            },
          ],
        },
      },
      {
        id: 4,
        position_id: 4,
        job_position: "Marketing Manager",
        employment_status: "Internship",
        actions: {
          objects: [
            {
              href: `localhost:3000/hr/job_position/detail?positionId=4&edit=false`,
              name: "View",
            },
            {
              href: `localhost:3000/hr/job_position/detail?positionId=4&edit=true`,
              name: "Edit",
            },
          ],
        },
      },
    ];
  }

  async function next() {
    getPositions(pagination.page + 1);
  }

  async function previous() {
    getPositions(pagination.page - 1);
  }

  async function addPosition() {
    goto(await $call.getPage("client-dashboard", "hr-position-create"));
    return;
  }

  onMount(async () => {
    getPositions(pagination.page);
  });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text={$_("positionList")} />
          <Button onclick={addPosition} icon={PlusIcon}
            ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" /> {$_("addNew")}</Button
          >
        </div>

        <!-- Document Table Data -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs bind:next bind:previous />
        {:else}
          <div class="pt-8">
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
