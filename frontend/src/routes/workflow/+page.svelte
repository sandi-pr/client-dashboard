<script>
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  // import {
  //   PUBLIC_API_URL_DEV,
  //   PUBLIC_API_URL_PROD,
  //   PUBLIC_STRAPI_URL_DEV,
  //   PUBLIC_STRAPI_URL_PROD,
  // } from "$env/static/public";
  // import axios from "axios";
  import {
    P,
    Table,
    Modal,
    Input,
    Avatar,
    Button,
    TableBody,
    TableHead,
    Pagination,
    TableBodyRow,
    TableBodyCell,
    TableHeadCell,
  } from "flowbite-svelte";
  import * as Icon from "flowbite-svelte-icons";
  import LoadingSpinner from "../lib/components/LoadingSpinner.svelte";
  // import { pagesBuilder } from "../lib/utils/pagesBuilder";
  // import { userDataSession } from "../lib/utils/auth";

  // const API_URL =
  //   process.env.NODE_ENV === "development"
  //     ? PUBLIC_API_URL_DEV
  //     : PUBLIC_API_URL_PROD;
  // const STRAPI_URL =
  //   process.env.NODE_ENV === "development"
  //     ? PUBLIC_STRAPI_URL_DEV
  //     : PUBLIC_STRAPI_URL_PROD;

  let workflowObject = [];
  let pagination = {};
  let pages = [];
  let showModal = false;
  let selectedWorkflowId = "";
  let isLoading = false;

  const handleGetWorkflows = async (page, pageSize, name = "") => {
    isLoading = true;
    console.log("not yet implemented");
    // const response = await axios.get(
    //   `${API_URL}/api/workflows?page=${page}&pageSize=${pageSize}&name=${name}`,
    //   {
    //     headers: {
    //       Authorization: "Bearer " + userDataSession().jwt,
    //     },
    //   }
    // );
    // // console.log(response);

    // if (response?.data?.status === 401) {
    //   alert("Session expired, please login again");
    //   goto("/");
    // } else if (
    //   response?.data &&
    //   response.data?.data &&
    //   response.data.data.length !== 0
    // ) {
    //   workflowObject = response.data.data;
    //   pagination = response.data.meta.pagination;
    //   pages = pagesBuilder(pagination);
    // } else {
    //   workflowObject = [];
    //   pagination = {};
    //   pages = [];
    // }
    isLoading = false;
  };

  const handleDeleteWorkflow = async () => {
    isLoading = true;
    showModal = false;
    console.log("not yet implemented");
    // const response = await axios.delete(
    //   `${API_URL}/api/workflows/${selectedWorkflowId}`,
    //   {
    //     headers: {
    //       Authorization: "Bearer " + userDataSession().jwt,
    //     },
    //   }
    // );
    // // console.log(response);

    // if (response?.data?.status === 401) {
    //   alert("Session expired, please login again");
    //   goto("/");
    // } else if (response?.data?.id) {
    //   alert("Workflow berhasil dihapus");
    //   handleGetWorkflows(1, 10);
    // }
    isLoading = false;
  };

  const handleSearch = (evt) => {
    // console.log(evt.target.value);
    handleGetWorkflows(1, 10, evt.target.value);
  };

  const previous = () => {
    if (pagination.page === 1) {
      alert("Cannot go to previous page, you are on the first page");
    } else {
      handleGetWorkflows(pagination.page - 1, pagination.pageSize);
    }
  };

  const next = () => {
    const activePage = pages.find((page) => page.active === true);
    if (pagination.page === activePage.name) {
      alert("Cannot go to next page, you are on the last page");
    } else {
      handleGetWorkflows(pagination.page + 1, pagination.pageSize);
    }
  };

  const handleSpecificPage = (evt) => {
    // console.log(evt.target.innerText);
    handleGetWorkflows(evt.target.innerText, pagination.pageSize);
  };

  onMount(async () => {
    // handleGetWorkflows(1, 10); dinyalakan lagi setelah ada datanya
  });

  $: {
    pages.forEach((page) => {
      if (page.name == pagination.page) {
        page.active = true;
      } else {
        page.active = false;
      }
    });
    pages = pages;
  }

  // dami data
  workflowObject = [
    {
      id: 1,
      attributes: {
        createdAt: "2024-02-07T08:00:00Z",
        name: "Workflow 1",
        steps: [
          { id: 1, photo: { url: "url1" } },
          { id: 2, photo: { url: "url2" } },
          { id: 3, photo: { url: "url3" } },
        ],
      },
    },
    {
      id: 2,
      attributes: {
        createdAt: "2024-02-06T10:30:00Z",
        name: "Workflow 2",
        steps: [
          { id: 4, photo: { url: "url4" } },
          { id: 5, photo: { url: "url5" } },
        ],
      },
    },
    {
      id: 3,
      attributes: {
        createdAt: "2024-02-05T15:45:00Z",
        name: "Workflow 3",
        steps: [
          { id: 6, photo: { url: "url6" } },
          { id: 7, photo: { url: "url7" } },
          { id: 8, photo: { url: "url8" } },
        ],
      },
    },
    {
      id: 4,
      attributes: {
        createdAt: "2024-02-04T12:20:00Z",
        name: "Workflow 4",
        steps: [
          { id: 9, photo: { url: "url9" } },
          { id: 10, photo: { url: "url10" } },
          { id: 11, photo: { url: "url11" } },
          { id: 12, photo: { url: "url12" } },
        ],
      },
    },
  ];

  showModal = false;
  selectedWorkflowId = null;

  pagination = {
    page: 1,
    pageSize: 10,
    total: workflowObject.length,
  };

  // pages = [
  //   { name: 1, active: false },
  //   { name: 2, active: false },
  //   { name: 3, active: false },
  //   { name: 4, active: false },
  //   { name: 5, active: false },
  // ];
</script>

<LoadingSpinner {isLoading} />

<div class="ml-0">
  <!-- <div class="flex flex-col m-4">
    <svelte:component this={Header} />
  </div> -->

  <div class="mx-4 mt-4 flex justify-between">
    <P size="3xl" weight="bold">Workflows</P>
    <Button
      class="bg-hc-brand-600 hover:bg-hc-brand-700"
      on:click={() => {
        goto("/workflow/create");
      }}>Add new workflow</Button
    >
  </div>

  {#if workflowObject.length === 0}
    <div class="mx-4 w-full h-96 flex flex-col justify-center items-center">
      <P size="2xl" weight="medium">No workflow found</P>
      <P size="2xl" weight="medium">Please create a new workflow</P>
    </div>
  {:else}
    <div class="mb-8 mt-6 mx-4">
      <Input id="search" type="search" placeholder="Search" on:change={handleSearch}>
        <Icon.SearchOutline slot="left" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
      </Input>
    </div>

    <div
      class="mx-4 border-2 border-hc-graylight-200 dark:border-hc-graydark-700 border-tool rounded-xl"
    >
      <Table>
        <TableHead class="text-hc-graylight-500 dark:text-hc-graydark-500">
          <TableHeadCell class="rounded-ss-xl">Upload Dates</TableHeadCell>
          <TableHeadCell>Workflow Names</TableHeadCell>
          <TableHeadCell>Recipients</TableHeadCell>
          <TableHeadCell class="rounded-se-xl">Action</TableHeadCell>
        </TableHead>
        <TableBody>
          {#each workflowObject as workflow}
            <TableBodyRow class="last:rounded-3xl">
              <!-- <TableBodyRow class="[&>last-child]:bg-pink-500"> -->
              <TableBodyCell
                >{new Date(workflow.attributes.createdAt)
                  .toLocaleString("en-GB")
                  .replace(",", "")}</TableBodyCell
              >
              <TableBodyCell>{workflow.attributes.name}</TableBodyCell>
              <TableBodyCell>
                <div class="flex items-center">
                  <P class="mr-6">{workflow.attributes.steps.length} users</P>
                  {#each workflow.attributes.steps as step (step.id)}
                    <!-- <Avatar stacked src={STRAPI_URL + step.photo.url} /> -->
                    <Avatar stacked />
                  {/each}
                </div>
              </TableBodyCell>
              <TableBodyCell class="flex">
                <Icon.EyeOutline
                  class="w-5 h-5 my-1 mx-2 text-hc-graylight-500 cursor-pointer"
                  on:click={() => {
                    goto("/workflow/" + workflow.id);
                  }}
                />
                <Icon.TrashBinOutline
                  class="w-5 h-5 my-1 mx-2 text-hc-graylight-500 cursor-pointer"
                  on:click={() => {
                    showModal = true;
                    selectedWorkflowId = workflow.id;
                  }}
                />
                <Icon.PenOutline
                  class="w-5 h-5 my-1 mx-2 text-hc-graylight-500 cursor-pointer"
                  on:click={() => {
                    goto("/workflow/update/" + workflow.id);
                  }}
                />
              </TableBodyCell>
            </TableBodyRow>
          {/each}
        </TableBody>
      </Table>
    </div>

    <div class="flex justify-between w-[95%]">
      <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
        Showing <span
          >{pagination.page === 1 ? 1 : (pagination.page - 1) * pagination.pageSize + 1}</span
        >
        to
        <span
          >{pagination.page * pagination.pageSize >= pagination.total
            ? pagination.total
            : pagination.page * pagination.pageSize}</span
        >
        of
        <span>{pagination.total}</span>
        Entries
      </div>

      <div class="m-4">
        <Pagination
          {pages}
          on:previous={previous}
          on:next={next}
          on:click={handleSpecificPage}
          activeClass="bg-hc-graylight-100 dark:bg-hc-graydark-700 dark:border-[#374151] dark:text-hc-graydark-300"
        />
      </div>
    </div>
  {/if}

  <Modal title="Delete Workflow" class="ml-80" bind:open={showModal} autoclose outsideclose>
    <p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
      Are you sure you want to delete this workflow?
    </p>
    <svelte:fragment slot="footer">
      <Button on:click={handleDeleteWorkflow}>Delete</Button>
      <Button color="alternative">Back</Button>
    </svelte:fragment>
  </Modal>
</div>
