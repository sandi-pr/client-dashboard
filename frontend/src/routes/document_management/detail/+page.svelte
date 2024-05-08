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
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import PDFViewer from "svelte-pdf";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import DownloadIcon from "$lib/icon/file-download-03.svg";
  import Divider from "../../../lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { add, format } from "date-fns";
  import TabDetail from "$lib/components/TabDetail.svelte";
  let screenWidth;

  let dropdownOpen = false;

  let popupModal = false;

  let popupModalDelegate = false;

  let approved = false;

  let onEdit = $state(false);

  let tabDatas = $state({
    selected: "Information",
    data: [
      {
        name: "Documents",
      },
      {
        name: "Information",
      },
    ],
  });

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
  let docId;
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
  let contract_start = "Jan 6, 2024";
  let contract_end = "Jan 7, 2024";
  let document_types = $state([
    { name: "General", value: "general" },
    { name: "Contract", value: "contract" },
  ]);
  let person_client = $state([
    {
      value: "1",
      name: "Azizi",
    },
    {
      value: "2",
      name: "Habibi",
    }
  ]);
  let document_type;
  let personClient;

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
  <div class="w-full pb-10">
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
      </div>

      <!-- Tab Bar -->
      <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

      <!-- Content -->
      <div class="">
        {#if tabDatas.selected === "Documents"}
          <div class="flex justify-between items-center pb-4">
            <P class="font-semibold text-hc-graylight-700 dark:text-hc-graydark-300"
              >{$_("documentName")}</P
            >
            <div class="flex items-center gap-2">
              <Button type="submit" class="w-max px-3 py-[10px] gap-1 focus-within:!ring-1">
                <img alt="refresh" src={DownloadIcon} class="w-5 h-5" />
                {$_("downdloadDocument")}
              </Button>
            </div>
          </div>
          <div class="flex justify-center w-full py-10 bg-gray-50 rounded-lg dark:bg-gray-800">
            {#if pdfReady}
              <PDFViewer showBorder={false} showButtons={false} scale={1.2} url={pdf} />
            {/if}
          </div>
        {:else if tabDatas.selected === "Information"}
          <div class="flex flex-col gap-4 pt-6">
            <!-- general information -->
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {$_("generalInformation")}
                </P>
              </div>
              <div class="flex flex-col grow gap-6">
                <div class="flex gap-5">
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("documentName")}
                    </P>
                    <Input
                      readonly={!onEdit}
                      bind:value={document.attributes.name}
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                </div>
                <div class="grid grid-cols-1">
                  <div class="flex flex-col gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >{$_("description")}</P
                    >
                    <Textarea
                      readonly={!onEdit}
                      bind:value={workflow.attributes.description}
                      rows="4"
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                </div>
              </div>
            </div>
            <Divider />
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {$_("documentType")}*
                </P>
              </div>
              <div class="flex flex-col grow gap-6">
                <div class="flex flex-col gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("documentType")}*</P
                  >
                  <Select class="bg-primary-light dark:bg-secondary-dark" items={document_types} bind:value={document_type} />
                </div>
              </div>
            </div>
            <Divider />
            <!-- general information -->
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {$_("contractInformation")}
                </P>
              </div>
              <div class="flex grow flex-col gap-6">
                <div class="flex flex-col grow gap-6">
                  <div class="flex flex-col gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >{$_("personsOrClient")}*</P
                    >
                    <Select class="bg-primary-light dark:bg-secondary-dark" items={person_client} bind:value={personClient} />
                  </div>
                </div>
                <div class="flex justify-center items-center gap-5">
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("contractStartPeriod")}
                    </P>
                    <Input
                      readonly={!onEdit}
                      bind:value={contract_start}
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                  <div class="mt-6">
                    <P>-</P>
                  </div>
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("contractEndPeriod")}
                    </P>
                    <Input
                      readonly={!onEdit}
                      bind:value={contract_end}
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
