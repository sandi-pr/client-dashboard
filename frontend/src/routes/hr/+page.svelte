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
  import { P, Card, Avatar, Heading, A, Button } from "flowbite-svelte";
  import { StateCard } from "$lib";
  import { _ } from "svelte-i18n";
  import UserPlusIcon from "$lib/icon/user-plus-02.svelte";
  import PackagePlusIcon from "$lib/icon/package-plus.svelte";
  import PackageSearchIcon from "$lib/icon/package-search.svelte";
  import Chart from "$lib/components/Chart.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { format } from "date-fns";

  const call = new CR3Call(globalThis.location.origin);

  let { showAI = false } = $props();

  let latestEmployee = $state([
    // {
    //   id: 1,
    //   join_date: "2022-01-01",
    //   employee_name: "John Doe",
    //   email: "john@doe.com",
    //   job_position: "Software Engineer",
    //   employee_status: "Permanent",
    //   actions: {
    //     objects: [
    //       {
    //         href: "hr/employee/detail?employeeId=" + 1 + "&edit=false",
    //         name: "View",
    //       },
    //       {
    //         href: "hr/employee/detail?employeeId=" + 1 + "&edit=true",
    //         name: "Edit",
    //       },
    //     ],
    //   },
    // }
  ]);
  let latestAttendance = $state([
  //   {
  //     id: 1,
  //     employee_name: "John Doe",
  //     job_position: "Software Engineer",
  //     checkin_time: "09:00",
  //     checkout_time: "18:00",
  //     status: getStatus("Checkout"),
  //   }
  ]);
  let latestRecruitment = $state([
  //   {
  //     id: 1,
  //     fullname: "John Doe",
  //     job_position: "Software Engineer",
  //     experience: "3+ YoE",
  //     email: "john@doe.com",
  //     actions: {
  //       objects: [
  //         {
  //           href: "hr/recruitment/detail?recruitmentId=" + 1 + "&edit=true",
  //           name: "Edit",
  //         },
  //       ],
  //     },
  //   }
  ]);
  let configs = $state({});
  let dataReady = $state(false);
  let personsIsEmpty = $state(true);
  // let documents = [];
  // let documentsNeed = [];
  // let documentsComplete = [];
  // let allUrls = [];
  // let chart = true;
  let monthlyEmployeeCounts = $state([0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
  let monthlyCheckinEmployeeCounts = $state([0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
  let monthlyLeaveEmployeeCounts = $state([0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
  let dataChartAll = $state({});
  let dataChartCheckin = $state({});
  let dataChartLeave = $state({});
  let chartReady = $state(false);
  let actionItems = $state([]);
  let chartAction = $state([]);

  async function getActionItems() {
    actionItems = [
      {
        // href: await call.getPage("client-dashboard", "hr-employee-attendance"),
        href: "/hr/employee/attendance",
        icon: UserPlusIcon,
        actionName: "Employee Attendance",
        desc: "desc Employee Attendance",
      },
      {
        // href: await call.getPage("client-dashboard", "hr-employee-create"),
        href: "/hr/employee/create",
        icon: PackagePlusIcon,
        actionName: "Add New Employee",
        desc: "desc Add New Employee",
      },
      {
        // href: await call.getPage("client-dashboard", "hr-recruitment"),
        href: "/hr/recruitment",
        icon: PackageSearchIcon,
        actionName: "Add New Recruitment",
        desc: "desc Add New Recruitment",
      },
    ];
  }

  let tabDatas = $state({
    selected: "Attendance",
    data: [
      {
        name: "Attendance",
      },
      {
        name: "Employee",
      },
      {
        name: "Recruitment",
      }
    ],
  });

  function generateChartData() {
    const currentDate = new Date();
    const twelveMonthsBefore = new Date();
    twelveMonthsBefore.setMonth(currentDate.getMonth() - 12);

    const formatDate = (date) => {
      const options = { month: "short", year: "numeric" };
      return new Intl.DateTimeFormat("en-US", options).format(date);
    };

    dataChartAll = {
      chart: {
        type: "area",
        fontFamily: "Inter, sans-serif",
        dropShadow: {
          enabled: false,
        },
        toolbar: {
          show: false,
        },
      },
      tooltip: {
        enabled: true,
        x: {
          show: false,
        },
      },
      fill: {
        type: "gradient",
        gradient: {
          opacityFrom: 0.55,
          opacityTo: 0,
          shade: "#17B169",
          gradientToColors: ["#17B169"],
        },
      },
      dataLabels: {
        enabled: false,
      },
      stroke: {
        width: 2,
      },
      grid: {
        show: false,
        strokeDashArray: 4,
        padding: {
          left: 2,
          right: 2,
          top: 0,
        },
      },
      series: [
        {
          name: "New users",
          data: monthlyEmployeeCounts,
          color: "#17B169",
        },
      ],
      xaxis: {
        categories: Array.from({ length: 12 }, (_, i) => {
          const date = new Date(twelveMonthsBefore);
          date.setMonth(twelveMonthsBefore.getMonth() + i + 1);
          return formatDate(date);
        }),
        labels: {
          show: false,
        },
        axisBorder: {
          show: false,
        },
        axisTicks: {
          show: false,
        },
      },
      yaxis: {
        show: false,
      },
    };

    dataChartCheckin = {
      chart: {
        type: "area",
        fontFamily: "Inter, sans-serif",
        dropShadow: {
          enabled: false,
        },
        toolbar: {
          show: false,
        },
      },
      tooltip: {
        enabled: true,
        x: {
          show: false,
        },
      },
      fill: {
        type: "gradient",
        gradient: {
          opacityFrom: 0.55,
          opacityTo: 0,
          shade: "#17B169",
          gradientToColors: ["#17B169"],
        },
      },
      dataLabels: {
        enabled: false,
      },
      stroke: {
        width: 2,
      },
      grid: {
        show: false,
        strokeDashArray: 4,
        padding: {
          left: 2,
          right: 2,
          top: 0,
        },
      },
      series: [
        {
          name: "New users",
          data: monthlyCheckinEmployeeCounts,
          color: "#17B169",
        },
      ],
      xaxis: {
        categories: Array.from({ length: 12 }, (_, i) => {
          const date = new Date(twelveMonthsBefore);
          date.setMonth(twelveMonthsBefore.getMonth() + i + 1);
          return formatDate(date);
        }),
        labels: {
          show: false,
        },
        axisBorder: {
          show: false,
        },
        axisTicks: {
          show: false,
        },
      },
      yaxis: {
        show: false,
      },
    };

    dataChartLeave = {
      chart: {
        type: "area",
        fontFamily: "Inter, sans-serif",
        dropShadow: {
          enabled: false,
        },
        toolbar: {
          show: false,
        },
      },
      tooltip: {
        enabled: true,
        x: {
          show: false,
        },
      },
      fill: {
        type: "gradient",
        gradient: {
          opacityFrom: 0.55,
          opacityTo: 0,
          shade: "#17B169",
          gradientToColors: ["#17B169"],
        },
      },
      dataLabels: {
        enabled: false,
      },
      stroke: {
        width: 2,
      },
      grid: {
        show: false,
        strokeDashArray: 4,
        padding: {
          left: 2,
          right: 2,
          top: 0,
        },
      },
      series: [
        {
          name: "New users",
          data: monthlyLeaveEmployeeCounts,
          color: "#17B169",
        },
      ],
      xaxis: {
        categories: Array.from({ length: 12 }, (_, i) => {
          const date = new Date(twelveMonthsBefore);
          date.setMonth(twelveMonthsBefore.getMonth() + i + 1);
          return formatDate(date);
        }),
        labels: {
          show: false,
        },
        axisBorder: {
          show: false,
        },
        axisTicks: {
          show: false,
        },
      },
      yaxis: {
        show: false,
      },
    };
    chartReady = true;
  }

  async function getChartAction(totalEmployee, checkinEmployee, leaveEmployee) {
    chartAction = [
      {
        // href: await call.getPage("client-dashboard", "hr-employee"),
        href: "hr/employee",
        dataName: "Total Employee",
        value: totalEmployee.total,
        change: totalEmployee.change,
        changeType: totalEmployee.changeType,
      },
      {
        // href: await call.getPage("client-dashboard", "hr-employee"),
        href: "hr/employee",
        dataName: "Employee checked-in",
        value: checkinEmployee.total,
        change: checkinEmployee.change,
        changeType: checkinEmployee.changeType,
      },
      {
        // href: await call.getPage("client-dashboard", "hr-employee"),
        href: "hr/employee",
        dataName: "Employee leave",
        value: leaveEmployee.total,
        change: leaveEmployee.change,
        changeType: leaveEmployee.changeType,
      },
    ];
  }

  let latestUsers = $state([]);

  async function viewAll(title) {
    if (title === "Attendance") {
      let href = await call.getPage("client-dashboard", "hr-employee-attendance");
      globalThis.location.href = href;
      return;
    } else if (title === "Employee") {
      let href = await call.getPage("client-dashboard", "hr-employee");
      globalThis.location.href = href;
      return;
    } else if (title === "Recruitment") {
      let href = await call.getPage("client-dashboard", "hr-recruitment");
      globalThis.location.href = href;
    }
  }

  function getStatus(is_active) {
    if (is_active === "Checkin") {
      return {
        type: "brand",
        name: "Checked-In",
      };
    } else if (is_active === "Checkout") {
      return {
        type: "success",
        name: "Checked Out",
      };
    } else if (is_active === "Permanent") {
      return {
        type: "success",
        name: "Permanent",
      };
    }
    return;
  }

  async function getEmployee() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-employees", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: 1,
          page_size: 25,
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
        response.result.data.sort(
          (a, b) => new Date(b.attributes.updatedAt) - new Date(a.attributes.updatedAt),
        );
        let latestData = response.result.data.slice(0, 5);

        let promises = latestData.map(async (employee) => {
          return {
            id: employee.id,
            join_date: employee.attributes.join_date,
            employee_name: employee.attributes.fullname,
            email: employee.attributes.email,
            job_position: employee.attributes.job,
            employee_status: employee.attributes.is_active,
            actions: {
              objects: [
                {
                  href:
                    `${await call.getPage("client-dashboard", "hr-employee-detail")}` +
                    "?employeeId=" +
                    employee.id +
                    "&edit=false",
                  name: "View",
                },
                {
                  href:
                    `${await call.getPage("client-dashboard", "hr-employee-detail")}` +
                    "?employeeId=" +
                    employee.id +
                    "&edit=true",
                  name: "Edit",
                },
              ],
            },
          };
        });

        latestEmployee = await Promise.all(promises);
        dataReady = true;
      }
    }
  }

  async function getAttendance() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-attendances", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: 1,
          page_size: 5,
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
        let latestData = response.result.data.slice(0, 5);
        let promises = latestData.map(async (attendance, index) => {
          return {
            id: attendance.id,
            employee_name: attendance.attributes.fullname,
            job_position: attendance.attributes.job,
            checkin_time: attendance.attributes.checkin_time,
            checkout_time: attendance.attributes.checkout_time,
            status: getStatus(attendance.attributes.is_active),
          };
        });

        latestAttendance = await Promise.all(promises);
        dataReady = true;
      }
    }
  }

  async function getRecruitment() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-recruitments", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: 1,
          page_size: 25,
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
        response.result.data.sort(
          (a, b) => new Date(b.attributes.updatedAt) - new Date(a.attributes.updatedAt),
        );
        let latestData = response.result.data.slice(0, 5);

        let promises = latestData.map(async (recruitment) => {
          return {
            id: recruitment.id,
            fullname: recruitment.attributes.fullname,
            job_position: recruitment.attributes.job,
            experience: recruitment.attributes.experience,
            email: recruitment.attributes.email,
            actions: {
              objects: [
                {
                  href:
                    `${await call.getPage("client-dashboard", "hr-recruitment-detail")}` +
                    "?recruitmentId=" +
                    recruitment.id +
                    "&edit=true",
                  name: "Edit",
                },
              ],
            },
          };
        });

        latestRecruitment = await Promise.all(promises);
        dataReady = true;
      }
    }
  }

  async function getCountPersonsPerMonth() {
    const today = new Date();
    let earlyDateCurrentMonth = new Date(today.getFullYear(), today.getMonth(), 1);
    let lastDateCurrentMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0);
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      for (let i = 0; i < 12; i++) {
        let response = await call.call("client-dashboard", "get-count-employee-per-month", {
          URLQuery: {},
          URLParams: {},
          Params: {
            earlyDate: earlyDateCurrentMonth.toISOString().split("T")[0],
            lastDate: lastDateCurrentMonth.toISOString().split("T")[0],
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          console.log(response.status);
          console.log(response);
          return;
        }
        // console.log(
        //   "telor",
        //   format(earlyDateCurrentMonth, "yyyy-MM-dd"),
        //   format(lastDateCurrentMonth, "yyyy-MM-dd"),
        // );

        monthlyEmployeeCounts.push(response.result.employee.meta.pagination.total);
        monthlyCheckinEmployeeCounts.push(response.result.checkinEmployee.meta.pagination.total);
        monthlyLeaveEmployeeCounts.push(response.result.leaveEmployee.meta.pagination.total);

        earlyDateCurrentMonth = new Date(
          earlyDateCurrentMonth.getFullYear(),
          earlyDateCurrentMonth.getMonth() - 1,
          1,
        );
        lastDateCurrentMonth = new Date(
          lastDateCurrentMonth.getFullYear(),
          lastDateCurrentMonth.getMonth(),
          0,
        );
      }

      monthlyEmployeeCounts.reverse();
      monthlyCheckinEmployeeCounts.reverse();
      monthlyLeaveEmployeeCounts.reverse();

      let totalEmployee = groupingDataChart(monthlyEmployeeCounts);
      let checkinEmployee = groupingDataChart(monthlyCheckinEmployeeCounts);
      let leaveEmployee = groupingDataChart(monthlyLeaveEmployeeCounts);

      getChartAction(totalEmployee, checkinEmployee, leaveEmployee);
    } else {
      console.log("no token");
    }
  }

  function groupingDataChart(data) {
    let total = 0;
    for (let i = 0; i < data.length; i++) {
      total += data[i];
    }

    let change = data[data.length - 1] - data[data.length - 2];

    let changeType = change > 0 ? "increase" : change < 0 ? "decrease" : "-";

    return {
      total: total,
      change: change,
      changeType: changeType,
    };
  }

  const damyData = $state({
    total: 0,
    change: 0,
    changeType: "-",})

  onMount(async () => {
    // getCountPersonsPerMonth();
    getChartAction(damyData, damyData, damyData);
    generateChartData();
    getActionItems();
    // getEmployee();
    // getAttendance();
    // getUsers();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Heading tag="h2" class="font-jakarta text-primary-light dark:text-primary-dark">
          {$_("hrOverview")}
        </Heading>

        <!-- Action Items -->
        <div>
          <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
            <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
              {$_("action")}
            </Heading>
          </div>
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 w-full items-center justify-center">
            {#each actionItems as item}
              <Card
                href={item.href}
                horizontal="true"
                class="cursor-pointer items-center w-full !p-4 min-h-28 bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark border-brand-light dark:border-brand-dark shadow-sm hover:no-underline"
                size="none"
              >
                <Avatar
                  size="sm"
                  border
                  class="!ring-0 bg-brand-secondary-light dark:bg-brand-secondary-dark !p-2 !h-10 !w-10"
                  ><svelte:component
                    this={item.icon}
                    svgClass="w-5 h-5 text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark"
                  /></Avatar
                >
                <div class="ml-4 gap[2px]">
                  <P
                    class="font-jakarta text-base font-semibold leading-normal text-primary-light dark:text-primary-dark"
                    >{$_(toCamelCase(item.actionName))}</P
                  >
                  <P
                    class="font-jakarta text-sm font-normal leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_(toCamelCase(item.desc))}</P
                  >
                </div>
              </Card>
            {/each}
          </div>
        </div>

        <!-- Reports Monitoring -->
        <div>
          <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
            <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
              {$_("reports")}
            </Heading>
          </div>
          <div
            class={`grid ${showAI ? "grid-cols-1 xl:grid-cols-2" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-3"} gap-4 items-center justify-center`}
          >
            {#each chartAction as { href, dataName, value, change, changeType }, index ({ href, dataName, value, change, changeType })}
              <Card
                horizontal="true"
                class="grid grid-cols-12 items-center !px-4 !py-0 min-h-28 justify-between bg-primary-light dark:bg-primary-dark border-secondary-light dark:border-secondary-dark shadow-sm"
                size="none"
              >
                <div class="space-y-2 col-span-6 sm:col-span-7 md:col-span-5">
                  <div class="flex items-center gap-2">
                    <P
                      size="base"
                      weight="semibold"
                      color="text-secondary-light dark:text-secondary-dark"
                      class="font-jakarta">{$_(toCamelCase(dataName))}</P
                    >
                    <P
                      size="sm"
                      weight="bold"
                      color={`${
                        changeType === "decrease"
                          ? "text-error-primary-light dark:text-error-primary-dark"
                          : changeType === "increase"
                            ? "text-success-primary-light dark:text-success-primary-dark"
                            : ""
                      }`}
                      class="font-jakarta leading-tight"
                      >{`${changeType === "decrease" ? "" : changeType === "increase" ? "+" : ""}${change}`}</P
                    >
                  </div>
                  <P
                    size="xl"
                    weight="semibold"
                    color="text-gray-900"
                    class="font-jakarta leading-[30px] text-primary-light dark:text-primary-dark"
                    >{value}</P
                  >
                  <A
                    {href}
                    class="text-sm font-semibold text-fg-brand-primary-light dark:text-fg-brand-primary-dark font-jakarta leading-tight"
                    >{$_("seeAll")}</A
                  >
                </div>
                <div class="col-span-6 sm:col-span-5 md:col-span-7">
                  {#if chartReady}
                    {#if index === 0}
                      <Chart dataChart={dataChartAll} />
                    {:else if index === 1}
                      <Chart dataChart={dataChartCheckin} />
                    {:else}
                      <Chart dataChart={dataChartLeave} />
                    {/if}
                  {:else}
                    <div></div>
                  {/if}
                </div>
              </Card>
            {/each}
          </div>
        </div>

        <div class="flex flex-col">
          <div class="flex justify-between">
            <Typography
              preset="Text lg/Semibold"
              color="primary"
              text={$_(toCamelCase(tabDatas.selected))}
            />
            <div></div>
            <!-- <Button on:click={() => viewAll(tabDatas.selected)} color="light">{$_("seeAll")}</Button> -->
          </div>
          <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

          <div>
            {#if tabDatas.selected === "Employee"}
              <!-- Person Table Data -->
              {#if latestEmployee.length > 0}
                <DataTable bind:datas={latestEmployee} bind:configs />
              {:else}
                <div class="pt-6">
                  <StateCard
                    illustration={IllustrationEmpty}
                    description={$_("youDontHaveAnyDataYet")}
                  />
                </div>
              {/if}
            {:else if tabDatas.selected === "Attendance"}
              {#if latestAttendance.length > 0}
                <DataTable bind:datas={latestAttendance} bind:configs />
              {:else}
                <div class="pt-6">
                  <StateCard
                    illustration={IllustrationEmpty}
                    description={$_("youDontHaveAnyDataYet")}
                  />
                </div>
              {/if}
              {:else if tabDatas.selected === "Recruitment"}
              {#if latestRecruitment.length > 0}
                <DataTable bind:datas={latestRecruitment} bind:configs />
              {:else}
                <div class="pt-6">
                  <StateCard
                    illustration={IllustrationEmpty}
                    description={$_("youDontHaveAnyDataYet")}
                  />
                </div>
              {/if}
            {/if}
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
