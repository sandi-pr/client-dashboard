<script>
  import {
    P,
    Span,
    Heading,
    Input,
    Textarea,
    Button,
    ButtonGroup,
    Dropdown,
    DropdownItem,
    Avatar,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    Modal,
    Select,
  } from "flowbite-svelte";
  import AutoComplete from "simple-svelte-autocomplete";
  import { ChevronDownSolid } from "flowbite-svelte-icons";
  import PDFViewer from "svelte-pdf";
  import CheckWhiteIcon from "$lib/icon/check-white.svg";
  import XIcon from "$lib/icon/x.svelte";
  import ShareIcon from "$lib/icon/share-06.svelte";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import Check from "$lib/icon/check-2.png";
  import RefreshIcon from "$lib/icon/refresh.svg";
  import DownloadIcon from "$lib/icon/file-download-03.svg";
  import Divider from "../../../lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { userDataSession } from "$lib/utils/auth";
  import { _ } from "svelte-i18n";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import {
    PUBLIC_STRAPI_URL_DEV,
    PUBLIC_STRAPI_URL_PROD,
    PUBLIC_API_URL_DEV,
    PUBLIC_API_URL_PROD,
  } from "$env/static/public";
  import { add, format } from "date-fns";
  import axios from "axios";
  import TabDetail from "$lib/components/TabDetail.svelte";
  let screenWidth;

  const isDev = process.env.NODE_ENV === "development";

  const API_URL = process.env.NODE_ENV === "development" ? PUBLIC_API_URL_DEV : PUBLIC_API_URL_PROD;

  const strapiUrl = isDev ? PUBLIC_STRAPI_URL_DEV : PUBLIC_STRAPI_URL_PROD;

  let dropdownOpen = false;

  let popupModal = false;

  let popupModalDelegate = false;

  let approved = false;

  let tabDatas = {
    selected: "Documents",
    data: [
      {
        name: "Documents",
      },
      {
        name: "Initial Documents",
      },
      {
        name: "History",
      },
      {
        name: "Information",
      },
      {
        name: "Recipients",
      },
    ],
  };

  // export let data;
  let inviteType = [
    { name: "Types", value: "" },
    { name: "By User", value: "User" },
    { name: "By Role", value: "Role" },
    { name: "By Group", value: "Group" },
  ];
  let inviteUnit = [];
  let modalType = "Signer";
  let searchInput = "";
  let inviteTypeSelected = "";
  let inviteUnitSelected = "";
  // let docId;
  let document;
  // let pdf;
  let filename;
  let signerData;
  let pdfReady = false;
  let status = [];
  let users = {};
  let ownerAddress;
  let workflow;
  let steps;

  // const getDocument = async () => {
  //   const result = userDataSession();
  //   try {
  //     const response = await fetch(`/api/document?jwt=${result.jwt}&id=${docId}`, {
  //       method: "GET",
  //       headers: {
  //         "Content-Type": "application/json",
  //       },
  //     });

  //     if (response.ok) {
  //       var address = result.user.identity.address;
  //       const result2 = await response.json();
  //       document = result2.data;
  //       pdf = `${strapiUrl}${result2.data.attributes.document_initial.data.attributes.url}`;
  //       filename = result2.data.attributes.document_initial.data.attributes.name;
  //       ownerAddress = result2.data.attributes.user_address;
  //       signerData = result2.data.attributes.signer_data;
  //       for (const signer of signerData) {
  //         if (signer[address] && signer[address].approved) {
  //           approved = true;
  //         }
  //       }
  //       status = result2.data.attributes.status;
  //       for (const item of status) {
  //         const address = item.address;

  //         if (users[address]) return;

  //         const response2 = await fetch(`/api/identity?jwt=${result.jwt}&address=${address}`, {
  //           method: "GET",
  //           headers: {
  //             "Content-Type": "application/json",
  //           },
  //         });

  //         if (response2.ok) {
  //           const result3 = await response2.json();
  //           if (result3.data.length > 0) {
  //             users[address] = result3.data[0];
  //           } else {
  //             users[address] = { error: "user not found" };
  //           }
  //         }
  //       }
  //       workflow = result2.data.attributes.workflow.data;
  //       const response3 = await axios.get(`${API_URL}/api/workflows/${workflow.id}`, {
  //         headers: {
  //           Authorization: "Bearer " + userDataSession().jwt,
  //         },
  //       });

  //       const responseUser = await Promise.all(
  //         response3.data.data.attributes.steps.map(async (step, index) => {
  //           const responseUser = await axios.get(
  //             `${API_URL}/api/users?type=email&value=${response3.data.data.attributes.steps[index].user.email}`,
  //             {
  //               headers: {
  //                 Authorization: "Bearer " + userDataSession().jwt,
  //               },
  //             },
  //           );
  //           return responseUser.data[0];
  //         }),
  //       );
  //       // console.log(response);
  //       // console.log(responseUser);

  //       if (response3?.data?.status === 401) {
  //         alert("Session expired, please login again");
  //         goto("/");
  //       }
  //       steps = response3.data.data.attributes.steps.map((signerReviewer, index) => {
  //         return {
  //           ...signerReviewer,
  //           unit: responseUser[0],
  //         };
  //       });
  //       console.log("ini", steps);
  //       pdfReady = true;
  //     }
  //   } catch (e) {
  //     console.log("error: ", e);
  //   }
  // };

  const handleSearchInvitees = async (value) => {
    try {
      console.log(value, inviteTypeSelected, inviteUnitSelected);
      const response = await axios.get(
        `${API_URL}/api/users?type=search&value=${value}&role=${inviteTypeSelected}&unit=${inviteUnitSelected}`,
        {
          headers: {
            Authorization: "Bearer " + userDataSession().jwt,
          },
        }
      );

      // console.log(response);

      if (response?.data?.status === 401) {
        alert("Session expired, please login again");
        goto("/");
      }

      if (response?.data) {
        return response.data;
      } else {
        return [];
      }
    } catch (error) {
      console.log(error);
    }
  };

  // const classes = {
  //   overall: null,
  //   controls: null,
  //   container: null,
  // };

  // let totalPage;

  onMount(() => {
    // docId = data.props.data;
    // console.log(docId);
    // getDocument();
  });

  //dami data
  pdfReady = true;
  import pdf from "$lib/template.pdf";
  document = {
    attributes: {
      name: "Document Name",
      marking: [
        { page: 1, signer: "John Doe" },
        { page: 2, signer: "Jane Smith" },
        { page: 3, signer: "Michael Johnson" },
      ],
      createdAt: "2023-07-02T20:51:00Z",
    },
  };

  ownerAddress = "user_1"; // ID pemilik dokumen

  workflow = {
    attributes: {
      name: "Workflow Name",
      unit: {
        data: {
          attributes: {
            name: "Workflow Unit",
          },
        },
      },
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
    },
  };

  status = [
    {
      timestamp: "2024-02-09T08:00:00",
      address: "user_1",
      status: "Document created",
    },
    {
      timestamp: "2024-02-10T09:30:00",
      address: "user_2",
      status: "Signed",
    },
    {
      timestamp: "2024-02-11T10:15:00",
      address: "user_3",
      status: "Finalized",
    },
    // Data history lainnya jika ada
  ];

  users = {
    user_1: {
      attributes: {
        users_permissions_user: {
          data: {
            attributes: {
              fullname: "John Doe",
              user_roles: {
                data: [
                  {
                    attributes: {
                      name: "Reviewer",
                    },
                  },
                ],
              },
            },
          },
        },
      },
    },
    user_2: {
      attributes: {
        users_permissions_user: {
          data: {
            attributes: {
              fullname: "Jane Smith",
              user_roles: {
                data: [
                  {
                    attributes: {
                      name: "Signer",
                    },
                  },
                ],
              },
            },
          },
        },
      },
    },
    user_3: {
      attributes: {
        users_permissions_user: {
          data: {
            attributes: {
              fullname: "Michael Johnson",
              user_roles: {
                data: [
                  {
                    attributes: {
                      name: "Signer",
                    },
                  },
                ],
              },
            },
          },
        },
      },
    },
    // Data pengguna lainnya jika ada
  };

  steps = [
    {
      user: {
        fullname: "John Doe",
        photo: {
          url: "/john_doe_photo.jpg",
        },
      },
      step: "Reviewer",
      identity: {
        data: {
          attributes: {
            address: "user_1",
          },
        },
      },
      unit: {
        user_roles: [{ name: "Reviewer" }],
      },
    },
    {
      user: {
        fullname: "Jane Smith",
        photo: {
          url: "/jane_smith_photo.jpg",
        },
      },
      step: "Signer",
      identity: {
        data: {
          attributes: {
            address: "user_2",
          },
        },
      },
      unit: {
        user_roles: [{ name: "Signer" }],
      },
    },
    {
      user: {
        fullname: "Michael Johnson",
        photo: {
          url: "/michael_johnson_photo.jpg",
        },
      },
      step: "Signer",
      identity: {
        data: {
          attributes: {
            address: "user_3",
          },
        },
      },
      unit: {
        user_roles: [{ name: "Signer" }],
      },
    },
    // Data penerima lainnya jika ada
  ];

  signerData = [
    {
      user_1: {
        approved: false,
      },
    },
    {
      user_2: {
        approved: true,
      },
    },
    {
      user_3: {
        approved: false,
      },
    },
    // Data lainnya jika ada
  ];
</script>

<svelte:window bind:innerWidth={screenWidth} />

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10 lg:ml-80">
    <div class="flex flex-col m-4">
      <div class="flex items-center pb-8">
        <div class="flex flex-col grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
              >Document Name <Button
                outline
                color="red"
                class="hover:!bg-transparent !border-hc-graylight-200 dark:!border-hc-graydark-600 focus-within:!ring-1 !p-2"
              >
                <CopyIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
              </Button></Heading
            >
          </div>
          <P
            class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
            >{$_("createOn")}{` uploadDate`}</P
          >
        </div>
        {#if approved}
          <div class="flex gap-2">
            <Button
              color="green"
              Approve
              class="w-max px-3 py-[10px] gap-1 focus-within:!ring-1 bg-hc-success-600 dark:bg-hc-success-500 hover:bg-hc-success-700 hover:dark:bg-hc-success-600 border-hc-success-500 dark:border-hc-success-400"
            >
              <img alt="refresh" src={CheckWhiteIcon} class="w-5 h-5" />
              {#if screenWidth > 1024}
                {$_("approved")}
              {/if}
            </Button>
          </div>
        {:else}
          <div class="flex gap-2">
            <Button
              on:click={() => (popupModal = true)}
              color="green"
              type="submit"
              Approve
              class="w-max px-3 py-[10px] gap-1 focus-within:!ring-1 bg-hc-success-600 dark:bg-hc-success-500 hover:bg-hc-success-700 hover:dark:bg-hc-success-600 border-hc-success-500 dark:border-hc-success-400"
            >
              <img alt="refresh" src={CheckWhiteIcon} class="w-5 h-5" />
              {#if screenWidth > 1024}
                {$_("approve")}
              {/if}
            </Button>
            <Button
              outline
              color="red"
              class=" w-max hover:!bg-transparent !text-sm !font-semibold !text-error-primary-light dark:!text-error-primary-dark !border-error-light dark:!border-error-dark focus-within:!ring-1 px-3 py-[10px] gap-1"
            >
              <XIcon
                svgClass="w-5 h-5 text-fg-error-primary-light dark:text-fg-error-primary-dark"
              />
              {#if screenWidth > 1024}
                {$_("reject")}
              {/if}
            </Button>
            <Button
              on:click={() => (popupModalDelegate = true)}
              outline
              color="purple"
              class=" w-max hover:!bg-transparent !text-sm !font-semibold !text-utility-brand-600-light dark:!text-utility-brand-600-dark !border-brand-light dark:!border-brand-dark focus-within:!ring-1 px-3 py-[10px] gap-1"
            >
              <ShareIcon
                svgClass="w-5 h-5 text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
              />
              {#if screenWidth > 1024}
                {$_("delegate")}
              {/if}
            </Button>
          </div>
        {/if}
      </div>

      <!-- Tab Bar -->
      <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

      <!-- Content -->
      <div class="m-4">
        {#if tabDatas.selected === "Documents"}
          <div class="flex justify-between items-center pb-4">
            <P class="font-semibold text-hc-graylight-700 dark:text-hc-graydark-300"
              >{$_("documentName")}</P
            >
            <div class="flex items-center gap-2">
              <Button color="light">
                {filename}
                <ChevronDownSolid class="w-3 h-3 ms-2 dark:text-white" />
              </Button>
              <Dropdown bind:open={dropdownOpen} class="w-48 overflow-y-auto py-1">
                <DropdownItem on:click={() => (dropdownOpen = false)}>
                  {filename}
                </DropdownItem>
              </Dropdown>
              <Button type="submit" class="w-max px-3 py-[10px] gap-1 focus-within:!ring-1">
                <img alt="refresh" src={DownloadIcon} class="w-5 h-5" />
                {$_("downdloadDocument")}
              </Button>
            </div>
          </div>
          <div class="flex justify-center w-full py-10 bg-gray-50 rounded-lg dark:bg-gray-800">
            {#if pdfReady}
              <PDFViewer url={pdf} />
            {/if}
          </div>
        {:else if tabDatas.selected === "Initial Documents"}
          <div class="flex justify-between items-center pb-4">
            <P class="font-semibold text-hc-graylight-700 dark:text-hc-graydark-300"
              >{$_("documentName")}</P
            >
            <div class="flex items-center gap-2">
              <Button color="light">
                {filename}
                <ChevronDownSolid class="w-3 h-3 ms-2 dark:text-white" />
              </Button>
              <Dropdown bind:open={dropdownOpen} class="w-48 overflow-y-auto py-1">
                <DropdownItem on:click={() => (dropdownOpen = false)}>
                  {filename}
                </DropdownItem>
              </Dropdown>
              <Button type="submit" class="w-max px-3 py-[10px] gap-1 focus-within:!ring-1">
                <img alt="refresh" src={DownloadIcon} class="w-5 h-5" />
                {$_("documentName")}
              </Button>
            </div>
          </div>
          <div class="flex justify-center w-full py-10 bg-gray-50 rounded-lg dark:bg-gray-800">
            {#if pdfReady}
              <PDFViewer url={pdf} />
            {/if}
          </div>
        {:else if tabDatas.selected === "History"}
          <Table shadow class="mt-2 font-jakarta dark:border-white">
            <TableHead class="bg-hc-graylight-50 dark:bg-hc-graydark-900">
              <TableHeadCell>TIMESTAMP</TableHeadCell>
              <TableHeadCell>USER</TableHeadCell>
              <TableHeadCell>UPDATE</TableHeadCell>
            </TableHead>
            <TableBody class="divide-y">
              {#if status.length > 0}
                {#each status as row}
                  <TableBodyRow
                    class=" dark:bg-hc-graydark-950 border-hc-graylight-200 dark:border-hc-graydark-800"
                  >
                    <TableBodyCell class="w-72 text-hc-graylight-600 dark:text-hc-graydark-400"
                      >{format(new Date(row.timestamp), "MM/dd/yyyy HH:mm")}</TableBodyCell
                    >
                    <TableBodyCell class="flex text-hc-graylight-900 dark:text-hc-graydark-50">
                      <Avatar
                        stacked
                        class="bg-gray-100 text-slate-600 text-sm"
                        style="width: 40px; height: 40px; overflow: hidden;"
                      >
                        <!-- nyalakan lagi klo udah ada datanya -->
                        <!-- <img
                          src={`${strapiUrl}${
                            users[row.address].attributes.personalize.data[0]
                              .attributes.url
                          }`}
                          alt="img"
                          style="width: 100%; height: 100%; object-fit: cover; border-radius: 50%;"
                        /> -->
                      </Avatar>
                      <div class="flex w-full h-full items-center ml-4 justify-between">
                        <div>
                          <Span
                            class="font-jakarta text-slate-700 dark:text-hc-graydark-300 text-sm leading-tight"
                            >{users[row.address].attributes.users_permissions_user.data.attributes
                              .fullname}</Span
                          >
                          <P
                            class="font-jakarta text-slate-600 dark:text-hc-graydark-400 text-sm leading-tight"
                            >{users[row.address].attributes.users_permissions_user.data.attributes
                              .user_roles.data[0].attributes.name}</P
                          >
                        </div>
                      </div>
                    </TableBodyCell>
                    <TableBodyCell class="w-60">
                      <Button
                        pill
                        class={`
                        px-1.5 py-0.5
                        ${
                          row.status === "Document created"
                            ? "bg-hc-brand-50 dark:bg-hc-brand-950 text-hc-brand-700 dark:text-hc-brand-300 border border-hc-brand-200 dark:border-hc-brand-800 hover:bg-hc-brand-50 dark:hover:bg-hc-brand-950"
                            : ""
                        }
                        ${
                          row.status === "Signed"
                            ? "bg-hc-success-50 dark:bg-hc-success-950 text-hc-success-700 dark:text-hc-success-300 border border-hc-success-200 dark:border-hc-success-800 hover:bg-hc-success-50 dark:hover:bg-hc-success-950"
                            : ""
                        }
                        ${
                          row.status === "Finalized"
                            ? "bg-hc-graylight-50 dark:bg-hc-graydark-950 text-hc-graylight-700 dark:text-hc-graydark-300 border border-hc-graylight-200 dark:border-hc-graydark-800 hover:bg-hc-graylight-50 dark:hover:bg-hc-graydark-950"
                            : ""
                        }
                        rounded-md focus-within:ring-0 font-medium
                        `}
                      >
                        {row.status}
                      </Button>
                    </TableBodyCell>
                  </TableBodyRow>
                {/each}
              {/if}
            </TableBody>
          </Table>
        {:else if tabDatas.selected === "Information"}
          <div class="flex flex-col gap-4">
            <!-- general information -->
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-hc-graylight-700 dark:text-hc-graydark-300"
                  >{$_("generalInformation")}</P
                >
              </div>
              <div class="flex flex-col grow gap-2 mt-1">
                <P
                  class="text-sm font-medium leading-tight text-hc-graylight-700 dark:text-hc-graydark-300"
                  >{$_("documentName")}</P
                >
                <Input readonly value={document.attributes.name} />
              </div>
              <div class="flex flex-col gap-2 mt-1">
                <P
                  class="text-sm font-medium leading-tight text-hc-graylight-700 dark:text-hc-graydark-300"
                ></P>
                <Input readonly value={`total page`} />
              </div>
            </div>
            <Divider />
            <!-- attributes -->
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-hc-graylight-700 dark:text-hc-graydark-300"
                  >{$_("attributes")}</P
                >
              </div>
              <div class="grow grid grid-cols-2 gap-3">
                <!-- {#each data.docDetail.attributes as attr}
                  <div class="flex flex-col grow gap-2 mt-1">
                      <P class="text-sm font-medium leading-tight text-hc-graylight-700 dark:text-hc-graydark-300">{attr.name}</P>
                      <Input readonly value={attr.value} />
                  </div>
                {/each} -->
              </div>
            </div>
            <Divider />
            <!-- owner -->
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-hc-graylight-700 dark:text-hc-graydark-300"
                  >{$_("owner")}</P
                >
              </div>
              <div class="flex grow gap-2 mt-1">
                <Avatar
                  stacked
                  class="bg-gray-100 text-slate-600 text-sm"
                  style="width: 40px; height: 40px; overflow: hidden;"
                >
                  <!-- dinyalakan lagi setelah ada datanya -->
                  <!-- <img
                    src={`${strapiUrl}${users[ownerAddress].attributes.personalize.data[0].attributes.url}`}
                    alt="img"
                    style="width: 100%; height: 100%; object-fit: cover; border-radius: 50%;"
                  /> -->
                </Avatar>
                <div class="flex-col items-center ml-1 justify-between">
                  <Span
                    class="font-jakarta text-slate-700 dark:text-hc-graydark-300 text-sm leading-tight"
                    >{users[ownerAddress].attributes.users_permissions_user.data.attributes
                      .fullname}</Span
                  >
                  <P
                    class="font-jakarta text-slate-600 dark:text-hc-graydark-400 text-sm leading-tight"
                    >{users[ownerAddress].attributes.users_permissions_user.data.attributes
                      .email}</P
                  >
                </div>
              </div>
            </div>
            <Divider />
            <!-- signing -->
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-hc-graylight-700 dark:text-hc-graydark-300"
                  >{$_("signing")}</P
                >
              </div>
              <div class="grow grid md:grid-cols-2 2xl:grid-cols-4 gap-3">
                <div class="flex flex-col grow gap-2 mt-1">
                  <P
                    class="text-sm font-medium leading-tight text-hc-graylight-700 dark:text-hc-graydark-300"
                    >{$_("pageNeedSigned")}</P
                  >
                  <Input
                    readOnly
                    value={document.attributes.marking.map((item) => item.page).join(", ")}
                  />
                </div>
                <div class="flex flex-col grow gap-2 mt-1">
                  <P
                    class="text-sm font-medium leading-tight text-hc-graylight-700 dark:text-hc-graydark-300"
                    >{format(new Date(document.attributes.createdAt), "MM/dd/yyyy HH:mm")}</P
                  >
                  <Input readonly value="2023-07-02 08:51 PM" />
                </div>
              </div>
            </div>
            <Divider />
            <!-- workflow information -->
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-hc-graylight-700 dark:text-hc-graydark-300"
                  >{$_("workflowInformation")}</P
                >
              </div>
              <div class="flex grow flex-col gap-6">
                <div class="flex gap-6">
                  <div class="flex flex-col grow gap-2 mt-1">
                    <P
                      class="text-sm font-medium leading-tight text-hc-graylight-700 dark:text-hc-graydark-300"
                      >{$_("workflowName")}</P
                    >
                    <Input readonly value={workflow.attributes.name} />
                  </div>
                  <div class="flex flex-col gap-2 mt-1">
                    <P
                      class="text-sm font-medium leading-tight text-hc-graylight-700 dark:text-hc-graydark-300"
                      >{$_("unit")}</P
                    >
                    <Input readonly value={workflow.attributes.unit.data.attributes.name} />
                  </div>
                </div>
                <div class="flex flex-col gap-2 mt-1">
                  <P
                    class="text-sm font-medium leading-tight text-hc-graylight-700 dark:text-hc-graydark-300"
                    >{$_("description")}</P
                  >
                  <Textarea readonly value={workflow.attributes.description} />
                </div>
              </div>
            </div>
          </div>
        {:else if tabDatas.selected === "Recipients"}
          <div class="flex flex-col gap-6 mt-2">
            {#each steps as recipient}
              <div class="flex grow gap-2">
                <Avatar
                  stacked
                  class="bg-gray-100 text-slate-600 text-sm ml-4"
                  style="width: 40px; height: 40px; overflow: hidden;"
                >
                  <!-- nyalakan lagi saat sudah ada datanya -->
                  <!-- <img
                    src={`${strapiUrl}${recipient.user.photo.url}`}
                    alt="img"
                    style="width: 100%; height: 100%; object-fit: cover; border-radius: 50%;"
                  /> -->
                </Avatar>
                <div class="flex-col items-center ml-1 justify-between">
                  <div class="flex items-center gap-2 mb-1">
                    <Span
                      class="font-jakarta text-slate-700 dark:text-hc-graydark-300 text-sm leading-tight"
                      >{recipient.user.fullname}</Span
                    >
                    <Button
                      pill
                      class={`
                        px-1.5 py-0
                        ${
                          recipient.step === "Reviewer"
                            ? "bg-hc-brand-50 dark:bg-hc-brand-950 text-hc-brand-700 dark:text-hc-brand-300 border border-hc-brand-200 dark:border-hc-brand-800 hover:bg-hc-brand-50 dark:hover:bg-hc-brand-950"
                            : ""
                        }
                        ${
                          recipient.step === "Signer"
                            ? "bg-hc-warning-50 dark:bg-hc-warning-950 text-hc-warning-700 dark:text-hc-warning-300 border border-hc-warning-200 dark:border-hc-warning-800 hover:bg-hc-warning-50 dark:hover:bg-hc-warning-950"
                            : ""
                        }
                        rounded-md focus-within:ring-0 font-medium
                        `}
                    >
                      {$_(toCamelCase(recipient.step))}
                    </Button>
                    <Button
                      pill
                      class={`
                          px-1.5 py-0
                          ${
                            signerData.find(
                              (item) =>
                                item[recipient.identity.data.attributes.address] &&
                                item[recipient.identity.data.attributes.address].approved === false
                            )
                              ? "bg-hc-warning-50 dark:bg-hc-warning-950 text-hc-warning-700 dark:text-hc-warning-300 border border-hc-warning-200 dark:border-hc-warning-800 hover:bg-hc-warning-50 dark:hover:bg-hc-warning-950"
                              : ""
                          }
                          ${
                            signerData.find(
                              (item) =>
                                item[recipient.identity.data.attributes.address] &&
                                item[recipient.identity.data.attributes.address].approved === true
                            )
                              ? "bg-hc-success-50 dark:bg-hc-success-950 text-hc-success-700 dark:text-hc-success-300 border border-hc-success-200 dark:border-hc-success-800 hover:bg-hc-success-50 dark:hover:bg-hc-success-950"
                              : ""
                          }
                          rounded-md focus-within:ring-0 font-medium
                        `}
                    >
                      {signerData.find(
                        (item) =>
                          item[recipient.identity.data.attributes.address] &&
                          item[recipient.identity.data.attributes.address].approved === false
                      )
                        ? $_("needSigning")
                        : $_("signed")}
                    </Button>
                  </div>
                  <P
                    class="font-jakarta text-slate-600 dark:text-hc-graydark-400 text-sm leading-tight"
                    >{recipient.unit.user_roles[0].name}</P
                  >
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div>
    <img src={Check} class="mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" alt="check.png" />
    <Heading tag="h3" class="text-lg font-normal">{$_("documentHasBeenApproved")}</Heading>
    <P class="mb-5">{$_("descDocHasBeenApproved")}</P>
    <Button color="light" class="me-2">
      <img src={RefreshIcon} alt="refresh.svg" class="mr-2" width="12" />
      {$_("callback")}</Button
    >
    <Button color="purple">{$_("confirm")}</Button>
  </div>
</Modal>

<Modal
  title={$_("chooseSomeoneToDelegate")}
  class="w-[75%]"
  bind:open={popupModalDelegate}
  autoclose
  outsideclose
>
  <div class="flex w-full p-1 items-center justify-between">
    <ButtonGroup class="grid grid-cols-12">
      <div
        class={`col-span-6 [&>*]:w-full [&>*]:!h-[2.75rem] [&>*]:m-0
          [&>*>div]:rounded-lg
        `}
      >
        <AutoComplete
          searchFunction={handleSearchInvitees}
          delay="1000"
          localFiltering={false}
          labelFieldName="username"
          valueFieldName="id"
          bind:selectedItem={searchInput}
          class="h-[2.75rem] w-full m-0 border-hc-graylight-300 rounded-l-lg focus:ring-hc-brand-500 focus:border-hc-brand-500"
        />
      </div>
      <!-- <div class="col-span-5"></div> -->
      <Select
        class={inviteTypeSelected === "Group"
          ? "rounded-none rounded-r-lg col-span-3"
          : "rounded-none col-span-3"}
        placeholder="Types"
        items={inviteType}
        bind:value={inviteTypeSelected}
      />
      <Select
        class={inviteTypeSelected === "Group" ? "hidden" : "rounded-none rounded-r-lg col-span-3"}
        placeholder="Units"
        items={inviteUnit}
        bind:value={inviteUnitSelected}
      />
    </ButtonGroup>
  </div>
  <div class="w-full">
    <Button
      class="w-full"
      disabled={searchInput.length === 0}
      on:click={() => {
        // handleAddSignerReviewers(modalType);
      }}>{$_("delegate")} {modalType}</Button
    >
  </div>
</Modal>
