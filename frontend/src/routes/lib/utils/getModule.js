import FolderIcon from "$lib/icon/folder.svelte";
import UsersIcon from "$lib/icon/users-03.svelte";
import MessageChatCircle from "$lib/icon/message-chat-circle.svelte";
import MessageCircleIcon from "$lib/icon/message-chat-circle.svelte";

const currentDate = new Date();
const twelveMonthsBefore = new Date();
twelveMonthsBefore.setMonth(currentDate.getMonth() - 12);

const formatDate = (date) => {
  const options = { month: "short", year: "numeric" };
  return new Intl.DateTimeFormat("en-US", options).format(date);
};

export const getModule = async () => {
  return [
    // {
    //   href: "/document_signing",
    //   icon: FolderIcon,
    //   name: "Document Signing",
    //   desc: "descModulDocSign",
    //   status: "coming soon",
    //   graph: {
    //     dataName: "Total Documents",
    //     value: "10",
    //     change: "3",
    //     changeType: "increase",
    //     dataChart: {
    //       chart: {
    //         height: "150px",
    //         width: "180px",
    //         type: "area",
    //         fontFamily: "Inter, sans-serif",
    //         dropShadow: {
    //           enabled: false,
    //         },
    //         toolbar: {
    //           show: false,
    //         },
    //       },
    //       tooltip: {
    //         enabled: true,
    //         x: {
    //           show: false,
    //         },
    //       },
    //       fill: {
    //         type: "gradient",
    //         gradient: {
    //           opacityFrom: 0.55,
    //           opacityTo: 0,
    //           shade: "#17B169",
    //           gradientToColors: ["#17B169"],
    //         },
    //       },
    //       dataLabels: {
    //         enabled: false,
    //       },
    //       stroke: {
    //         width: 2,
    //       },
    //       grid: {
    //         show: false,
    //         strokeDashArray: 4,
    //         padding: {
    //           left: 2,
    //           right: 2,
    //           top: 0,
    //         },
    //       },
    //       series: [
    //         {
    //           name: "New users",
    //           data: [3, 6, 1, 8, 2, 7],
    //           color: "#17B169",
    //         },
    //       ],
    //       xaxis: {
    //         categories: Array.from({ length: 12 }, (_, i) => {
    //           const date = new Date(twelveMonthsBefore);
    //           date.setMonth(twelveMonthsBefore.getMonth() + i + 1);
    //           return formatDate(date);
    //         }),
    //         labels: {
    //           show: false,
    //         },
    //         axisBorder: {
    //           show: false,
    //         },
    //         axisTicks: {
    //           show: false,
    //         },
    //       },
    //       yaxis: {
    //         show: false,
    //       },
    //     },
    //   },
    // },
    {
      href: "/person",
      icon: UsersIcon,
      name: "CRM",
      desc: "descCrm",
      status: "coming soon",
    },
    {
      href: "/hexa-bots",
      icon: MessageChatCircle,
      name: "Hexa Bots",
      desc: "descHexaBots",
      status: "coming soon",
    },
    {
      href: "/deployed_bots",
      icon: MessageCircleIcon, 
      name: "Deployed Bots" ,
      desc: "descDeployedBots",
      status: "coming soon",
    },
  ];
};
