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
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { format } from "date-fns";

  import { call } from "$lib";

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });

  let datas = $state([]);

  let gettingDatas = $state(false);

  let configs = $state({});
  // let employees = $state([]);
  let idAttendance = $state("");
  let attendance = $state([]);

  let selectedDate = $state(format(new Date(), "yyyy-MM-dd"));
  let currentDate = $state(format(new Date(), "yyyy-MM-dd"));

  function getStatus(attendance) {
    if (attendance === "Checkin") {
      return {
        type: "success",
        name: $_("Checkin"),
      };
    }

    if (attendance === "Checkout") {
      return {
        type: "brand",
        name: $_("Checkout"),
      };
    }

    if (attendance === "Absent") {
      return {
        type: "danger",
        name: $_("Absent"),
      };
    }

    if (attendance === "-") {
      return {
        type: "warning",
        name: $_("Uncheckin"),
      };
    }

    return;
  }

  async function getEmployeeAttendance() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await $call.call("client-dashboard", "get-employee-attendances", {
        URLQuery: {},
        URLParams: {},
        Params: {
          date: selectedDate,
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
        idAttendance = response.result.data[0].id.toString();
        return response.result.data[0].attributes.payload;
      } else {
        return [];
      }
    }
  }

  async function next() {
    getEmployeeAttendance();
  }

  async function previous() {
    getEmployeeAttendance();
  }

  async function getEmployee() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      gettingDatas = true;
      let page = 1;
      let employees = [];
      while (true) {
        let response = await $call.call("client-dashboard", "get-employees", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 25,
          },
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
          employees = [...employees, ...response.result.data];
          page++;
        } else {
          break;
        }
      }

      datas = [];

      for (let i = 0; i < employees.length; i++) {
        datas.push({
          id: employees[i].id,
          employee_name: employees[i].attributes.employee_name,
          job_position: employees[i].attributes.job_position,
          checkin_time: "-",
          checkout_time: "-",
          status: "",
          actions: {
            objects: [
              {
                attendance: true,
                actionButton: {
                  CheckIn: () => {
                    checkin_time(employees[i].id);
                  },
                  Absent: () => {
                    absent(employees[i].id);
                  },
                },
              },
            ],
          },
        });
      }

      attendance = await getEmployeeAttendance();

      if (attendance.length > 0) {
        for (let i = 0; i < attendance.length; i++) {
          for (let j = 0; j < datas.length; j++) {
            if (attendance[i].employee_id === datas[j].id) {
              datas[j].checkin_time = attendance[i].checkin_time;
              datas[j].checkout_time = attendance[i].checkout_time;

              console.log("status", attendance[j].status);

              if (selectedDate === format(new Date(), "yyyy-MM-dd")) {
                datas[j].actions.objects[0].actionButton = {};
              } else {
                if (attendance[i].status === "Checkin") {
                  console.log("checkin");
                  datas[j].actions.objects[0].actionButton = {
                    CheckOut: () => {
                      checkout_time(datas[j].id);
                    },
                  };
                }

                if (attendance[i].status === "Checkout") {
                  console.log("checkout");
                  datas[j].actions.objects[0].actionButton = {};
                }

                if (attendance[i].status === "Absent") {
                  console.log("absent");
                  datas[j].actions.objects[0].actionButton = {
                    // CheckIn: () => {
                    //   checkin_time(datas[j].id);
                    // },
                  };
                }
              }

              datas[j].status = getStatus(attendance[i].status);

              console.log("attendance", datas[j]);
            } else {
              datas[j].status = getStatus("-");
            }
          }
        }
      } else {
        for (let i = 0; i < datas.length; i++) {
          datas[i].status = getStatus("-");
        }
      }

      console.log("datas", datas);

      gettingDatas = false;

      configs = {
        search: true,
      };
    }
  }

  async function checkin_time(id) {
    let datas2 = datas;
    for (let i = 0; i < datas2.length; i++) {
      if (datas2[i].id === id) {
        datas2[i].checkin_time = format(new Date(), "HH:mm");
        datas2[i].status = getStatus("Checkin");
        datas2[i].actions.objects[0].actionButton = {
          CheckOut: () => {
            checkout_time(datas2[i].id);
          },
        };
        break;
      }
    }

    let promise = datas2.map((data) => {
      return {
        employee_id: data.id,
        employee_name: data.employee_name,
        job_position: data.job_position,
        checkin_time: data.checkin_time,
        checkout_time: data.checkout_time,
        status: "Checkin",
      };
    });

    let payload = await Promise.all(promise);

    console.log("payload", payload);

    let token = sessionStorage.getItem("token");

    if (token) {
      let response;
      if (attendance.length > 0) {
        response = await $call.call("client-dashboard", "update-attendance", {
          URLQuery: {},
          URLParams: {},
          URLParams: {
            id: idAttendance,
          },
          Params: {
            payload: JSON.stringify(payload),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
      } else {
        response = await $call.call("client-dashboard", "create-attendance", {
          URLQuery: {},
          URLParams: {},
          Params: {
            date: selectedDate,
            payload: JSON.stringify(payload),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
      }

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      console.log("success", response);
      datas = datas2;
    }
  }

  async function checkout_time(id) {
    let datas2 = datas;
    for (let i = 0; i < datas2.length; i++) {
      if (datas2[i].id === id) {
        datas2[i].checkout_time = format(new Date(), "HH:mm");
        datas2[i].status = getStatus("Checkout");
        datas2[i].actions.objects[0].actionButton = {};
        break;
      }
    }

    let promise = datas2.map((data) => {
      return {
        employee_id: data.id,
        employee_name: data.employee_name,
        job_position: data.job_position,
        checkin_time: data.checkin_time,
        checkout_time: data.checkout_time,
        status: "Checkout",
      };
    });

    let payload = await Promise.all(promise);

    let token = sessionStorage.getItem("token");

    if (token) {
      let response = await $call.call("client-dashboard", "update-attendance", {
        URLQuery: {},
        URLParams: {
          id: idAttendance,
        },
        Params: {
          payload: JSON.stringify(payload),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      console.log("success", response);

      datas = datas2;
    }
  }

  async function absent(id) {
    let datas2 = datas;
    for (let i = 0; i < datas2.length; i++) {
      if (datas2[i].id === id) {
        datas2[i].status = getStatus("Absent");
        datas2[i].actions.objects[0].actionButton = {
          // CheckIn: () => {
          //   checkin_time(datas2[i].id);
          // },
        };
      }
    }
    let promise = datas2.map((data) => {
      return {
        employee_id: data.id,
        employee_name: data.employee_name,
        job_position: data.job_position,
        checkin_time: data.checkin_time,
        checkout_time: data.checkout_time,
        status: "Absent",
      };
    });

    let payload = await Promise.all(promise);

    let token = sessionStorage.getItem("token");

    if (token) {
      let response;

      if (attendance.length > 0) {
        response = await $call.call("client-dashboard", "update-attendance", {
          URLQuery: {},
          URLParams: {
            id: idAttendance,
          },
          Params: {
            payload: JSON.stringify(payload),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
      } else {
        response = await $call.call("client-dashboard", "create-attendance", {
          URLQuery: {},
          URLParams: {},
          Params: {
            date: selectedDate,
            payload: JSON.stringify(payload),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
      }

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      console.log("success", response);

      datas = datas2;
    }
  }

  onMount(async () => {
    getEmployee();
  });

  $effect(() => {
    if (selectedDate !== currentDate) {
      currentDate = selectedDate;
      getEmployee();
    }
  });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography
            preset="Display sm/Semibold"
            color="primary"
            text={$_("employeeAttendance")}
          />
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
