<script>
  import { onMount } from "svelte";
  import { P, Card, Avatar, Heading, A } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import UserPlusIcon from "../lib/icon/user-plus-02.svelte";
  // import lightningIcon from "../lib/icon/lightning-01.svg";
  // import fileSearchIcon from "../lib/icon/file-search-01.svg";
  import Chart from "../lib/components/Chart.svelte";
  import DataTable from "../lib/components/DataTable.svelte";
  // import TablePersons from "../lib/components/TablePersons.svelte";
  import { toCamelCase } from "../lib/utils/toCamelCase.js";
  // import { userDataSession } from "$lib/utils/auth";
  // import { PUBLIC_STRAPI_URL_DEV, PUBLIC_STRAPI_URL_PROD } from "$env/static/public";

  // const isDev = process.env.NODE_ENV === "development";

  // const strapiUrl = isDev ? PUBLIC_STRAPI_URL_DEV : PUBLIC_STRAPI_URL_PROD;

  // export let data;

  // let documentsLatest = [];
  let personsLatest = [];
  let documentsReady = false;
  // let documents = [];
  // let documentsNeed = [];
  // let documentsComplete = [];
  // let allUrls = [];
  // let chart = true;
  let monthlyPersonsCounts = [3, 6, 1, 8, 2, 7, 4, 9, 5, 0, 6, 3];
  let monthlyActivePersonsCounts = [9, 2, 7, 4, 0, 6, 3, 5, 8, 1, 7, 2];
  let monthlyInactivePersonsCounts = [0, 6, 3, 5, 8, 1, 7, 2, 9, 4, 0, 6];
  let dataChartAll = {};
  let dataChartNeed = {};
  let dataChartComplete = {};
  let chartReady = false;

  let actionItems = [
    {
      href: "#",
      icon: UserPlusIcon,
      actionName: "Create new persons",
      desc: "desc Create New Person",
    },
  ];

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
          data: monthlyPersonsCounts,
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

    dataChartNeed = {
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
          data: monthlyActivePersonsCounts,
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

    dataChartComplete = {
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
          data: monthlyInactivePersonsCounts,
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

  let docData = [
    {
      href: "/person",
      dataName: "Total Persons",
      value: "10",
      change: "3",
      changeType: "increase",
    },
    {
      href: "/person",
      dataName: "Active Persons",
      value: "7",
      change: "3",
      changeType: "decrease",
    },
    {
      ihref: "/person",
      dataName: "Inactive Persons",
      value: "3",
      change: "3",
      changeType: "increase",
    },
  ];

  personsLatest = [
    {
      attributes: {
        customerId: "19846829",
        fullName: "Azzisi Bassem",
        dateOfBirth: "01/12/1994",
        email: "azzis144635@gmail.com",
        phone: "+44 33771021",
        country: "UK",
        address: "John Smith 42, Primrose Avenue, London EC2M 7EB, United Kingdom",
        city: "London",
        state: "London",
        zipCode: "EC2M 7EB",
        createdAt: "01-12-2024 03:02",
        active: true,
      },
    },
    {
      attributes: {
        customerId: "34678291",
        fullName: "John Doe",
        dateOfBirth: "03/15/1985",
        email: "johndoe@example.com",
        phone: "+1 123456789",
        country: "USA",
        address: "123 Main St, Anytown USA",
        city: "Anytown",
        state: "AnyState",
        zipCode: "12345",
        createdAt: "01-12-2024 03:02",
        active: true,
      },
    },
    {
      attributes: {
        customerId: "56473829",
        fullName: "Alice Johnson",
        dateOfBirth: "06/22/1978",
        email: "alice@example.com",
        phone: "+61 987654321",
        country: "Australia",
        address: "123 Main St, Anytown AU",
        city: "Anytown",
        state: "AnyState",
        zipCode: "12345",
        createdAt: "01-12-2024 03:02",
        active: false,
      },
    },
    {
      attributes: {
        customerId: "97261835",
        fullName: "Mohammed Ahmed",
        dateOfBirth: "12/04/1990",
        email: "mohammed@example.com",
        phone: "+91 9876543210",
        country: "India",
        address: "123 Main St, Anytown IN",
        city: "Anytown",
        state: "AnyState",
        zipCode: "12345",
        createdAt: "01-12-2024 03:02",
        active: true,
      },
    },
  ];

  let datas = personsLatest.map((person, index) => {
    return {
      id: index + 1,
      customerId: person.attributes.customerId,
      fullName: person.attributes.fullName,
      email: person.attributes.email,
      phone: person.attributes.phone,
      country: person.attributes.country,
      actions: {
        objects: [
          {
            href: "/person/" + person.attributes.customerId,
            name: "Edit",
          },
        ],
      },
    };
  });

  let configs = {
    headers: "Persons",
    seperator: true,
    actions: [
      {
        name: "View All",
        href: "/person/list",
        type: "info",
      },
    ],
  };

  onMount(async () => {
    // await getAllDocument();
    generateChartData();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Heading tag="h2" class="font-jakarta text-primary-light dark:text-primary-dark">
          CRM
        </Heading>

        <!-- Action Items -->
        <div>
          <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
            <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
              {$_("action")}
            </Heading>
          </div>
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 w-full items-center justify-center">
            {#each actionItems as { href, icon, actionName, desc }}
              <Card
                {href}
                horizontal="true"
                class="items-center w-full !p-4 min-h-28 bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark border-brand-light dark:border-brand-dark shadow-sm hover:no-underline"
                size="none"
              >
                <Avatar
                  size="sm"
                  border
                  class="!ring-0 bg-brand-secondary-light dark:bg-brand-secondary-dark !p-2 !h-10 !w-10"
                  ><svelte:component
                    this={icon}
                    svgClass="w-5 h-5 text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark"
                  /></Avatar
                >
                <div class="ml-4 gap[2px]">
                  <P
                    class="font-jakarta text-base font-semibold leading-normal text-primary-light dark:text-primary-dark"
                    >{$_(toCamelCase(actionName))}</P
                  >
                  <P
                    class="font-jakarta text-sm font-normal leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_(toCamelCase(desc))}</P
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
            class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4 items-center justify-center"
          >
            {#each docData as { href, dataName, value, change, changeType }, index ({ href, dataName, value, change, changeType })}
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
                          : "text-success-primary-light dark:text-success-primary-dark"
                      }`}
                      class="font-jakarta leading-tight"
                      >{`${changeType === "decrease" ? "-" : "+"}${change}`}</P
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
                      <Chart dataChart={dataChartNeed} />
                    {:else}
                      <Chart dataChart={dataChartComplete} />
                    {/if}
                  {:else}
                    <div></div>
                  {/if}
                </div>
              </Card>
            {/each}
          </div>
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
