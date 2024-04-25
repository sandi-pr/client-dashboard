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
  // import {
  //   PUBLIC_API_URL_DEV,
  //   PUBLIC_API_URL_PROD,
  //   PUBLIC_STRAPI_URL_DEV,
  //   PUBLIC_STRAPI_URL_PROD,
  // } from "$env/static/public";
  import { P, Button } from "flowbite-svelte";
  import { TabBar, InformationForm } from "hexacomponents-svelte";
  import { onMount } from "svelte";
  // import axios from "axios";
  import { goto } from "$app/navigation";
  import DataTable from "$lib/components/DataTable.svelte";
  // import TabBar from "$lib/components/TabBar.svelte";
  // import InformationForm from "$lib/components/InformationForm.svelte";
  import { datas, informationFormDatas, configs } from "./dummy";

  // let apiUrl = process.env.NODE_ENV === "development" ? PUBLIC_API_URL_DEV : PUBLIC_API_URL_PROD;
  // let strapiUrl =
  //   process.env.NODE_ENV === "development" ? PUBLIC_STRAPI_URL_DEV : PUBLIC_STRAPI_URL_PROD;

  let screenWidth;
  // let user = {};
  // let user = singularData;
  let user = informationFormDatas;
  let userCopy = {};
  let isAdmin = true;
  let tabDatas = {
    selected: "My Profile",
    data: [
      {
        name: "My Profile",
      },
      {
        name: "Roles & Admin",
      },
    ],
  };

  const getUser = async () => {
    try {
      console.log("not working");
      // const cookieResponse = await axios.get(`${apiUrl}/api/cookies`);
      // const cookies = cookieResponse.data.find((cookie) => cookie.name === "crd_cookies");
      // const jwt = JSON.parse(cookies.value).jwt;

      // const response = await axios.get(`${strapiUrl}/api/users/me?populate=*`, {
      //   headers: {
      //     Authorization: `Bearer ${jwt}`,
      //   },
      // });

      // // isAdmin = response.data?.user_roles[0]?.admin || false;
      // user = response.data;
      // userCopy = { ...user };
      // // console.log(user);
    } catch (error) {
      if (error?.response?.status === 401) {
        console.log("unauthorized");
        window.location = "/";
      } else {
        console.log("error");
        window.location = "/";
      }
    }
  };

  const updateUser = async () => {
    console.log("nunggu Microservice");
  };

  onMount(() => {
    // getUser();
  });
</script>

<svelte:window bind:innerWidth={screenWidth} />

{#if Object.keys(user).length !== 0}
  <div class="m-4">
    <!-- Header -->
    <div class="mt-4 flex justify-between">
      <P size="3xl" weight="bold">{isAdmin && "Super"} User Menu</P>
      <div>
        <Button
          class="bg-white font-bold mr-2 border border-hc-graylight-300 hover:bg-hc-graylight-50 text-gray-900 dark:bg-white dark:hover:bg-hc-graylight-50"
          disabled={JSON.stringify(user) === JSON.stringify(userCopy)}
          on:click={getUser}>Cancel</Button
        >
        <Button
          class="bg-hc-brand-600 hover:bg-hc-brand-700"
          disabled={JSON.stringify(user) === JSON.stringify(userCopy)}
          on:click={updateUser}>Save</Button
        >
      </div>
    </div>

    <!-- Submenu -->
    <TabBar {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

    <!-- Content -->
    {#if isAdmin}
      {#if tabDatas.selected === "My Profile"}
        <div class="mx-4 mt-8">
          <InformationForm datas={user} {screenWidth} />
        </div>
      {:else if tabDatas.selected === "Roles & Admin"}
        <div class="mt-8">
          <DataTable {datas} {configs} />
        </div>
      {/if}
    {:else}
      hello user
    {/if}
  </div>
{/if}
