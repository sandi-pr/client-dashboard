<script>
  import { onMount } from "svelte";
  import { Input } from "flowbite-svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { Button } from "flowbite-svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import { _ } from "svelte-i18n";
  import { StateCard } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";

  const call = new CR3Call(globalThis.location.origin);

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });

  let datas = $state([
    // {
    //   id: 1,
    //   employee_name: "John Doe",
    //   job_position: "Software Engineer",
    //   checkin_time: "09:00",
    //   checkout_time: "18:00",
    //   attendance: {
    //     status: getStatus("Checkin"),
    //     actions: {
    //       objects: [
    //         {
    //           name: "Checkin",
    //           on_click: () => { console.log("Checkin"); },
    //         },
    //         {
    //           name: "Checkout",
    //           on_click: () => { console.log("Checkout"); },
    //         },
    //       ],
    //     }
    //   },
    // },
    // {
    //   id: 2,
    //   employee_name: "John Doe",
    //   job_position: "Software Engineer",
    //   checkin_time: "09:00",
    //   checkout_time: "18:00",
    //   attendance: {
    //     status: "",
    //     action: {
    //       objects: [
    //         {
    //           name: "Checkin",
    //           on_click: () => { console.log("Checkin"); },
    //         },
    //         {
    //           name: "Checkout",
    //           on_click: () => { console.log("Checkout"); },
    //         },
    //       ],
    //     }
    //   },
    // },
  ]);

  let configs = $state({});

  let selectedDate = $state();

  function getStatus(attendance) {
    if (attendance === "Checkin") {
      return {
        type: "success",
        name: $_("checkin"),
      };
    }

    if (is_active === "Checkout") {
      return {
        type: "brand",
        name: $_("checkout"),
      };
    }

    return;
  }

  async function getEmployeeAttendance(page) {
    if (page > pagination.pageCount) return;
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-employee-attendances", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: page,
          page_size: pagination.pageSize,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      if (response.result.data.length > 0) {
        pagination = response.result.meta.pagination;
        configs = {
          search: true,
          pagination: response.result.meta.pagination,
        };
        let promises2 = response.result.data.map(async (employee, index) => {
          return {
            id: employee.id,
            employee_name: employee.attributes.employee_name,
            job_position: employee.attributes.job_position,
            checkin_time: employee.attributes.checkin_time,
            checkout_time: employee.attributes.checkout_time,
            attendance: {
              status: getStatus(employee.attributes.attendance_status),
              actions: {
                objects: [
                  {
                    name: "Checkin",
                    on_click: () => { console.log("Checkin"); },
                  },
                  {
                    name: "Checkout",
                    on_click: () => { console.log("Checkout"); },
                  },
                ],
              }
            },
          };
        });

        datas = await Promise.all(promises2);
      }
    }
  }

  async function next() {
    getEmployeeAttendance(pagination.page + 1);
  }

  async function previous() {
    getEmployeeAttendance(pagination.page - 1);
  }

  async function addEmployee() {
    await call.goToPage("client-dashboard", "hr-employee-create");
    return;
  }

  onMount(async () => {
    // getEmployeeAttendance(pagination.page);
  });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text={$_("employeeAttendance")} />
          <!-- <Button onclick={addEmployee} icon={PlusIcon}
            ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" /> {$_("addNew")}</Button
          > -->
        </div>

        <div class="mt-8 mb-2">
          <Input
            type="date"
            bind:value={selectedDate}
            class="w-min bg-primary-light dark:bg-secondary-dark"
          />
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
