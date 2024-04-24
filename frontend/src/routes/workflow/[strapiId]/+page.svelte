<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Dean Ramhan <dean@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import { onMount } from "svelte";
  // import axios from "axios";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { P, Button, Textarea, Avatar, Modal } from "flowbite-svelte";
  import * as Icon from "flowbite-svelte-icons";
  // import {
  //   PUBLIC_API_URL_DEV,
  //   PUBLIC_API_URL_PROD,
  //   PUBLIC_STRAPI_URL_DEV,
  //   PUBLIC_STRAPI_URL_PROD,
  // } from "$env/static/public";
  // import { userDataSession } from "$lib/utils/auth";
  import LoadingSpinner from "$lib/components/LoadingSpinner.svelte";
  import TabBar from "$lib/components/TabBar.svelte";

  // const API_URL = process.env.NODE_ENV === "development" ? PUBLIC_API_URL_DEV : PUBLIC_API_URL_PROD;
  // const STRAPI_URL =
  //   process.env.NODE_ENV === "development" ? PUBLIC_STRAPI_URL_DEV : PUBLIC_STRAPI_URL_PROD;

  let workflow = {};
  let organization = {};
  let organizationPhoto = "";
  let steps = [];
  let formattedDate;
  $: workflowId = $page.params.strapiId;

  let tabDatas = {
    selected: "information",
    data: [
      {
        name: "information",
      },
      {
        name: "recipients",
      },
    ],
  };
  let showModal = false;
  let isLoading = false;

  const handleGetWorkflow = async () => {
    isLoading = true;
    console.log("not yet implemented");
    // try {
    //   const response = await axios.get(`${API_URL}/api/workflows/${workflowId}`, {
    //     headers: {
    //       Authorization: "Bearer " + userDataSession().jwt,
    //     },
    //   });

    //   const responseUser = await Promise.all(
    //     response.data.data.attributes.steps.map(async (step, index) => {
    //       const responseUser = await axios.get(
    //         `${API_URL}/api/users?type=email&value=${response.data.data.attributes.steps[index].user.email}`,
    //         {
    //           headers: {
    //             Authorization: "Bearer " + userDataSession().jwt,
    //           },
    //         }
    //       );
    //       return responseUser.data[0].unit;
    //     })
    //   );

    //   if (response?.data?.status === 401) {
    //     alert("Session expired, please login again");
    //     goto("/");
    //   }

    //   workflow = response.data.data.attributes;
    //   organization = response.data.data.attributes.organization.data.attributes;
    //   organizationPhoto = STRAPI_URL + organization.photo.data.attributes.url;
    //   steps = response.data.data.attributes.steps.map((signerReviewer, index) => {
    //     return {
    //       ...signerReviewer,
    //       unit: responseUser[0],
    //     };
    //   });
    //   console.log(steps);

    //   const date = new Date(workflow.createdAt);
    //   const options = {
    //     year: "numeric",
    //     month: "2-digit",
    //     day: "2-digit",
    //     hour: "2-digit",
    //     minute: "2-digit",
    //     hour12: true,
    //   };
    //   formattedDate = `Created on ${date.toLocaleString("en-US", options)}`;

    // } catch (error) {
    //   console.log(error);
    // }
    isLoading = false;
  };

  const handleDeleteWorkflow = async () => {
    // isLoading = true;
    // try {
    //   const response = await axios.delete(`${API_URL}/api/workflows/${workflowId}`, {
    //     headers: {
    //       Authorization: "Bearer " + userDataSession().jwt,
    //     },
    //   });
    //   // console.log(response);
    //   if (response?.data?.status === 401) {
    //     alert("Session expired, please login again");
    //     goto("/");
    //   } else if (response?.data?.id) {
    //     alert("Workflow berhasil dihapus");
    //     goto("/workflow");
    //   }
    // } catch (error) {
    //   console.log(error);
    //   isLoading = false;
    // }
  };

  onMount(async () => {
    // handleGetWorkflow();
  });

  //data dami
  workflow = {
    name: "Nama Workflow",
    description: "Deskripsi Workflow",
  };
  formattedDate = "2024-02-09";
  organization = {
    name: "Nama Organisasi",
    email: "email@organisasi.com",
  };
  steps = [
    {
      id: 1,
      user: {
        username: "Username Pengguna 1",
        photo: {
          url: "URL_foto_pengguna_1",
        },
      },
      step: "Langkah 1",
      unit: {
        name: "Unit Pengguna 1",
      },
    },
    {
      id: 2,
      user: {
        username: "Username Pengguna 2",
        photo: {
          url: "URL_foto_pengguna_2",
        },
      },
      step: "Langkah 2",
      unit: {
        name: "Unit Pengguna 2",
      },
    },
    {
      id: 3,
      user: {
        username: "Username Pengguna 3",
        photo: {
          url: "URL_foto_pengguna_3",
        },
      },
      step: "Langkah 3",
      unit: {
        name: "Unit Pengguna 3",
      },
    },
  ];
</script>

<LoadingSpinner {isLoading} />

<div class="m-4">
  <div class="flex justify-between">
    <P size="3xl" weight="semibold">{workflow.name}</P>
    <div class="w-[60%]"></div>
    <div class="flex">
      <Button
        outline
        color="red"
        class="h-10"
        on:click={() => {
          showModal = true;
        }}><Icon.TrashBinOutline size="sm" />&nbsp;Delete</Button
      >
      <Button
        on:click={() => {
          goto("/workflow/update/" + workflowId);
        }}
        class="mx-2 h-10 bg-hc-brand-600 hover:bg-hc-brand-700 dark:bg-hc-brand-600 dark:hover:bg-hc-brand-700"
        ><Icon.EditOutline size="sm" />&nbsp;Edit</Button
      >
    </div>
  </div>
  <P class="mt-4 text-hc-graylight-500 dark:text-hc-graydark-500">{formattedDate}</P>

  <!-- Submenu -->
  <TabBar {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

  <!-- Content -->
  {#if tabDatas.selected === "information"}
    <div class="grid grid-cols-12 gap-4 mt-8">
      <P size="lg" weight="semibold" class="col-span-3">Workflow Information</P>
      <P class="col-span-9">Workflow Name</P>
      <!-- <P class="col-span-2">Unit</P> -->
      <div class="col-span-3"></div>
      <Textarea class="col-span-9 resize-none" rows="1" disabled value={workflow.name} />
      <!-- <Textarea class="col-span-2 resize-none" rows="1" disabled value="?" /> -->
      <div class="col-span-3"></div>
      <P class="col-span-9">Description</P>
      <div class="col-span-3"></div>
      <Textarea class="col-span-9 resize-none" disabled value={workflow.description} />
    </div>
    <div class="mt-4 border-t-[1px] border-hc-graylight-300 dark:border-hc-graydark-700"></div>
    <div class="grid grid-cols-12 gap-4 mt-4">
      <P size="lg" weight="semibold" class="col-span-3">Owner</P>
      <div class="col-span-9 flex items-center">
        <Avatar class="size-14" src={organizationPhoto} />
        <div>
          <P class="ml-2 font-semibold">{organization.name}</P>
          <P class="ml-2 text-sm">{organization.email ? organization.email : "-"}</P>
        </div>
      </div>
    </div>
  {:else if tabDatas.selected === "recipients"}
    <div>
      {#each steps as user, index (user.id)}
        <div class="flex items-center mt-8">
          <!-- <Avatar rounded class="size-14" src={STRAPI_URL + user.user.photo.url} /> -->
          <Avatar rounded class="size-14" />
          <div>
            <div class="flex items-center">
              <P class="ml-2 font-semibold">{user.user.username}</P>
              <P
                class={`ml-2 text-xs px-2 py-1 border-[1px] font-semibold rounded-lg ${
                  user.step.includes("Signer")
                    ? "bg-hc-error-50 border-hc-error-200 text-hc-error-600 dark:bg-hc-error-600 dark:border-hc-error-600 dark:text-white"
                    : ""
                }`}>{user.step}</P
              >
            </div>
            <P class="ml-2 text-sm">{user.unit.name}</P>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<Modal title="Delete Workflow" class="lg:ml-80" bind:open={showModal} autoclose outsideclose>
  <p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
    Are you sure you want to delete this workflow?
  </p>
  <svelte:fragment slot="footer">
    <Button on:click={handleDeleteWorkflow}>Delete</Button>
    <Button color="alternative">Back</Button>
  </svelte:fragment>
</Modal>
