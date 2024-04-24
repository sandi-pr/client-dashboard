<script>
  import { onMount } from "svelte";
  import { P, Card, Avatar, Button, Heading, A } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import UploadIcon from "$lib/icon/upload-01.svelte";
  import LightningIcon from "$lib/icon/lightning-01.svelte";
  import FileSearchIcon from "$lib/icon/file-search-01.svelte";
  import Chart from "$lib/components/Chart.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  // import { userDataSession } from "$lib/utils/auth";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  // import { PUBLIC_STRAPI_URL_DEV, PUBLIC_STRAPI_URL_PROD } from "$env/static/public";

  // const isDev = process.env.NODE_ENV === "development";

  // const strapiUrl = isDev ? PUBLIC_STRAPI_URL_DEV : PUBLIC_STRAPI_URL_PROD;

  // export let data;
  let moduleMonitoringData = [
    {
      href: "/signing",
      icon: UploadIcon,
      moduleName: "initiate New Signing",
      desc: "desc Initiate New Signing",
    },
    {
      href: "/workflow/create",
      icon: LightningIcon,
      moduleName: "create A New Workflow",
      desc: "desc Create A New Workflow",
    },
    {
      href: "/workflow",
      icon: FileSearchIcon,
      moduleName: "view All Workflows",
      desc: "desc View All Workflows",
    },
  ];

  let documentsLatest = [];
  let documentsReady = false;
  // let documents = [];
  // let documentsNeed = [];
  // let documentsComplete = [];
  let allUrls = [];
  // let chart = true;
  let monthlyDocumentCounts = [3, 6, 1, 8, 2, 7, 4, 9, 5, 0, 6, 9];
  let monthlyDocumentNeedCounts = [9, 2, 7, 4, 0, 6, 3, 5, 8, 1, 7, 4];
  let monthlyDocumentCompleteCounts = [0, 6, 3, 5, 8, 1, 7, 2, 9, 4, 3, 6];
  let dataChartAll = {};
  let dataChartNeed = {};
  let dataChartComplete = {};
  let chartReady = false;

  // const getAllDocument = async () => {
  //   const result = userDataSession();
  //   const address = result.user.identity.address;
  //   let page = 1;
  //   do {
  //     try {
  //       const response = await fetch(
  //         `/api/document?jwt=${result.jwt}&address=${address}&chart=true&page=${page}`,
  //         {
  //           method: "GET",
  //           headers: {
  //             "Content-Type": "application/json",
  //           },
  //         },
  //       );

  //       if (response.ok) {
  //         const result2 = await response.json();
  //         if (result2.data.length === 0) {
  //           break;
  //         }
  //         if (chart) {
  //           for (var i = 0; i < 5; i++) {
  //             documentsLatest.push(result2.data[i]);
  //           }
  //           for (const document of documentsLatest) {
  //             const signerData = document.attributes.signer_data;

  //             const urls = [];

  //             for (const item of signerData) {
  //               const address = Object.keys(item)[0];

  //               const response2 = await fetch(
  //                 `/api/identity?jwt=${result.jwt}&address=${address}`,
  //                 {
  //                   method: "GET",
  //                   headers: {
  //                     "Content-Type": "application/json",
  //                   },
  //                 },
  //               );

  //               if (response2.ok) {
  //                 const result3 = await response2.json();
  //                 if (result3.data.length > 0) {
  //                   let url = "/";
  //                   if (
  //                     result3.data[0].attributes.personalize.data !== null &&
  //                     result3.data[0].attributes.personalize.data.length > 0
  //                   ) {
  //                     url = `${strapiUrl}${result3.data[0].attributes.personalize.data[0].attributes.url}`;
  //                     urls.push({ url: url });
  //                   } else {
  //                     urls.push({ url: url });
  //                   }
  //                 }
  //               }
  //             }
  //             allUrls.push(urls);
  //           }
  //           chart = false;
  //         }
  //         documents = documents.concat(result2.data);
  //         const filteredDocumentsNeed = result2.data.filter((document) => {
  //           if (document.attributes.signer_data && document.attributes.signer_data.length > 0) {
  //             return document.attributes.signer_data.some(
  //               (item) => item[address] && item[address].approved === false,
  //             );
  //           }
  //           return false;
  //         });
  //         documentsNeed = documentsNeed.concat(filteredDocumentsNeed);
  //         const filteredDocumentsComplete = result2.data.filter((document) => {
  //           const hasStatus = document.attributes.status && document.attributes.status.length > 0;

  //           if (hasStatus) {
  //             const lastStatus = document.attributes.status[document.attributes.status.length - 1];
  //             return lastStatus.status === "Finalize Done";
  //           }

  //           return false;
  //         });
  //         documentsComplete = documentsComplete.concat(filteredDocumentsComplete);
  //         page++;
  //       }
  //     } catch (e) {
  //       console.log(e);
  //     }
  //   } while (true);
  //   groupAndCountDocumentsByMonth(documents, "all");
  //   groupAndCountDocumentsByMonth(documentsNeed, "need");
  //   groupAndCountDocumentsByMonth(documentsComplete, "complete");
  //   console.log("ini all document", documents);
  //   console.log("ini documentLatest", documentsLatest);
  //   documentsReady = true;
  // };

  // function groupAndCountDocumentsByMonth(docs, type) {
  //   let documentCountsByMonth = {};
  //   docs.forEach((document) => {
  //     const date = new Date(document.attributes.createdAt);
  //     const monthYearKey = `${date.getMonth() + 1}-${date.getFullYear()}`;

  //     if (documentCountsByMonth[monthYearKey]) {
  //       documentCountsByMonth[monthYearKey]++;
  //     } else {
  //       documentCountsByMonth[monthYearKey] = 1;
  //     }
  //   });

  //   const currentDate = new Date();
  //   const monthYearList = Array.from({ length: 12 }, (_, index) => {
  //     const date = new Date(currentDate);
  //     date.setMonth(currentDate.getMonth() - index);
  //     return `${date.getMonth() + 1}-${date.getFullYear()}`;
  //   }).reverse();

  //   let index = 0;

  //   if (type === "all") {
  //     index = 0;
  //     monthlyDocumentCounts = monthYearList.map(
  //       (monthYear) => documentCountsByMonth[monthYear] || 0
  //     );
  //     docData[index].change =
  //       monthlyDocumentCounts[monthlyDocumentCounts.length - 1] -
  //       monthlyDocumentCounts[monthlyDocumentCounts.length - 2];
  //   }
  //   if (type === "need") {
  //     index = 1;
  //     monthlyDocumentNeedCounts = monthYearList.map(
  //       (monthYear) => documentCountsByMonth[monthYear] || 0
  //     );
  //     docData[index].change =
  //       monthlyDocumentNeedCounts[monthlyDocumentNeedCounts.length - 1] -
  //       monthlyDocumentNeedCounts[monthlyDocumentNeedCounts.length - 2];
  //   }
  //   if (type === "complete") {
  //     index = 2;
  //     monthlyDocumentCompleteCounts = monthYearList.map(
  //       (monthYear) => documentCountsByMonth[monthYear] || 0
  //     );
  //     docData[index].change =
  //       monthlyDocumentCompleteCounts[monthlyDocumentCompleteCounts.length - 1] -
  //       monthlyDocumentCompleteCounts[monthlyDocumentCompleteCounts.length - 2];
  //   }

  //   docData[index].value = docs.length;
  //   if (docData[index].change >= 0) {
  //     docData[index].changeType = "increase";
  //   } else {
  //     docData[index].changeType = "decrease";
  //   }

  //   generateChartData();
  // }

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
          name: "New Document",
          data: monthlyDocumentCounts,
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
          name: "New Document",
          data: monthlyDocumentNeedCounts,
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
          name: "New Document",
          data: monthlyDocumentCompleteCounts,
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
      href: "/documents",
      dataName: "Documents signed",
      value: "10",
      change: "3",
      changeType: "increase",
    },
    {
      href: "/documents",
      dataName: "Needs Signing",
      value: "7",
      change: "3",
      changeType: "decrease",
    },
    {
      ihref: "/documents",
      dataName: "Completed",
      value: "3",
      change: "3",
      changeType: "increase",
    },
  ];

  onMount(async () => {
    // await getAllDocument();
    generateChartData();
  });

  // dami
  // const dataChartDami = {
  //   chart: {
  //     height: "150px",
  //     width: "180px",
  //     type: "area",
  //     fontFamily: "Inter, sans-serif",
  //     dropShadow: {
  //       enabled: false,
  //     },
  //     toolbar: {
  //       show: false,
  //     },
  //   },
  //   tooltip: {
  //     enabled: true,
  //     x: {
  //       show: false,
  //     },
  //   },
  //   fill: {
  //     type: "gradient",
  //     gradient: {
  //       opacityFrom: 0.55,
  //       opacityTo: 0,
  //       shade: "#17B169",
  //       gradientToColors: ["#17B169"],
  //     },
  //   },
  //   dataLabels: {
  //     enabled: false,
  //   },
  //   stroke: {
  //     width: 2,
  //   },
  //   grid: {
  //     show: false,
  //     strokeDashArray: 4,
  //     padding: {
  //       left: 2,
  //       right: 2,
  //       top: 0,
  //     },
  //   },
  //   series: [
  //     {
  //       name: "New users",
  //       data: [50, 22, 35, 42, 60, 29, 15, 55, 48, 33, 18, 63],
  //       color: "#17B169",
  //     },
  //   ],
  //   xaxis: {
  //     categories: [
  //       "02-2023",
  //       "03-2023",
  //       "04-2023",
  //       "05-2023",
  //       "06-2023",
  //       "07-2023",
  //       "08-2023",
  //       "09-2023",
  //       "10-2023",
  //       "11-2023",
  //       "12-2023",
  //       "01-2024",
  //     ],
  //     labels: {
  //       show: false,
  //     },
  //     axisBorder: {
  //       show: false,
  //     },
  //     axisTicks: {
  //       show: false,
  //     },
  //   },
  //   yaxis: {
  //     show: false,
  //   },
  // };

  documentsLatest = [
    {
      attributes: {
        updatedAt: "2024-02-07T08:00:00Z",
        name: "Document 1",
        emailOwner: "abc@gmail.com",
        signerData: [
          {
            "abc@gmail.com": {
              approve: false,
            },
          },
        ],
        status: [
          {
            status: "Needs Signing",
            updatedAt: "2024-02-07T08:00:00Z",
          },
        ],
      },
    },
    {
      attributes: {
        updatedAt: "2024-02-06T10:30:00Z",
        name: "Document 2",
        emailOwner: "bcd@gmail.com",
        signerData: [
          {
            "bcd@gmail.com": {
              approve: false,
            },
          },
        ],
        status: [
          {
            status: "Approval Error",
            updatedAt: "2024-02-06T10:30:00Z",
          },
        ],
      },
    },
    {
      attributes: {
        updatedAt: "2024-02-05T15:45:00Z",
        name: "Document 3",
        emailOwner: "bcd@gmail.com",
        signerData: [
          {
            "abc@gmail.com": {
              approve: true,
            },
          },
        ],
        status: [
          {
            status: "Needs Finalize",
            updatedAt: "2024-02-05T15:45:00Z",
          },
        ],
      },
    },
    {
      attributes: {
        updatedAt: "2024-02-04T12:20:00Z",
        name: "Document 4",
        emailOwner: "abc@gmail.com",
        signerData: [
          {
            "abc@gmail.com": {
              approve: true,
            },
          },
        ],
        status: [
          {
            status: "Done",
            updatedAt: "2024-02-04T12:20:00Z",
          },
        ],
      },
    },
  ];

  allUrls = [
    [{ url: "" }, { url: "" }, { url: "" }, { url: "" }, { url: "" }, { url: "" }],
    [{ url: "" }, { url: "" }, { url: "" }],
    [{ url: "" }, { url: "" }, { url: "" }, { url: "" }],
    [{ url: "" }, { url: "" }, { url: "" }, { url: "" }],
  ];

  const getStatusDoc = (status) => {
    let statusType = status.status;
    let colorType = "brand";
    if (statusType === "Initiate") {
      colorType = "brand";
    }
    if (statusType === "Needs Signing") {
      colorType = "warning";
    }
    if (statusType === "Approval Error") {
      colorType = "error";
    }
    if (statusType === "Needs Finalize") {
      colorType = "brand";
    }
    if (statusType === "Done") {
      colorType = "success";
    }
    return {
      name: statusType,
      type: colorType,
      updatedAt: status.updatedAt,
    };
  };

  let email = "abc@gmail.com";

  const getActionBotton = (document) => {
    let approve = false;
    let reject = false;
    let recall = false;
    let finalize = false;
    for (const user of document.signerData) {
      if (!user[email]?.approve) {
        approve = true;
        reject = true;
      }
    }
    if (document.emailOwner === email) {
      recall = true;
      let allSigned = true;
      for (const user of document.signerData) {
        if (!user?.email?.approve) {
          allSigned = false;
          break;
        }
      }
      if (allSigned) {
        if (document.status[document.status.length - 1].status !== "Done") finalize = true;
      }
    }
    return {
      approve,
      reject,
      recall,
      finalize,
    };
  };

  let datas = documentsLatest.map((document, index) => {
    let lastStatus = document.attributes.status[document.attributes.status.length - 1];
    return {
      id: index + 1,
      updatedAt: document.attributes.updatedAt,
      name: document.attributes.name,
      status: getStatusDoc(lastStatus),
      recipients: allUrls[index],
      actions: {
        objects: [
          {
            iconOnly: true,
            actionButton: getActionBotton(document.attributes),
          },
        ],
      },
    };
  });

  let configs = {
    headers: "Documents",
    seperator: true,
    documentManagement: true,
    actions: [
      {
        name: "View All",
        href: "/documents",
        type: "info",
      },
    ],
  };
</script>

<div class="flex bg-primary-light dark:bg-primary-dark">
  <div class="w-full pb-10 lg:ml-80">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Heading tag="h2" class="font-jakarta text-primary-light dark:text-primary-dark"
          >{$_("documentManagement")}</Heading
        >

        <!-- Module Monitoring Data -->
        <div>
          <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
            <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark"
              >{$_("action")}</Heading
            >
          </div>
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 w-full items-center justify-center">
            {#each moduleMonitoringData as { href, icon, moduleName, desc }}
              <Card
                {href}
                horizontal="true"
                class="items-center w-full !p-4 min-h-28 bg-primary-light dark:bg-primary-dark border-brand-light dark:border-brand-dark shadow-sm hover:no-underline hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark hover:shadow-md"
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
                    >{$_(toCamelCase(moduleName))}</P
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

        <!-- Document Monitoring Data -->
        <div>
          <div
            class="flex items-center justify-between border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4"
          >
            <Heading tag="h6" class="w-fit font-jakarta text-primary-light dark:text-primary-dark">
              {$_("reports")}
            </Heading>
            <Button
              href="/documents"
              color="light"
              class="bg-secondary-light dark:bg-secondary-dark border-secondary-light dark:border-secondary-dark hover:bg-secondary-hover-light hover:dark:bg-secondary-hover-dark focus-within:ring-0 hover:no-underline"
              >{$_("viewAll")}</Button
            >
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
                <div class="space-y-2 col-span-6 sm:col-span-7 md:col-span-6">
                  <div class="flex w-3/4 items-center gap-2">
                    <P
                      size="base"
                      weight="semibold"
                      color="text-secondary-light dark:text-secondary-dark"
                      class="font-jakarta min-w-fit max-w-fit">{$_(toCamelCase(dataName))}</P
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
                    class="font-jakarta leading-[30px] text-primary-light dark:text-primary-dark"
                    >{value}</P
                  >
                  <A
                    {href}
                    class="text-sm font-semibold text-fg-brand-primary-light dark:text-fg-brand-primary-dark font-jakarta leading-tight"
                    >{$_("seeAll")}</A
                  >
                </div>
                <div class="col-span-6 sm:col-span-5 md:col-span-6">
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

        <!-- Document Table Data -->
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
