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
  import { P, Input, Avatar, Button, Select } from "flowbite-svelte";
  import { SelectedCard, InformationForm, Typography, triggerToast } from "$lib";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { SearchOutline } from "flowbite-svelte-icons";
  import DataTable from "$lib/components/DataTable.svelte";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import MarketIcon from "$lib/icon/building-02.svelte";
  import OfficeIcon from "$lib/icon/building-06.svelte";
  import BankIcon from "$lib/icon/bank.svelte";
  import UserIcon from "$lib/icon/user-03.svelte";
  import SignalIcon from "$lib/icon/signal-01.svelte";
  import OilGas from "$lib/icon/oil-gas.svelte";
  import SlidersIcon from "$lib/icon/sliders-04.svelte";
  import { datas, informationFormDatas, configs } from "./dummy";
  import { CR3Call } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

  let screenWidth;
  let user = informationFormDatas;
  let userCopy = {};
  let isAdmin = true;
  let fileImage = null;
  let fileImageName = null;
  let fileImageResult = null;
  let dropzoneActive = false;
  let currentPasswordMatch = false;
  let selectedModelTrocr = "";
  let options = [];
  let selectedCompanyType = null;
  let selectedIndustryType = null;
  let tabDatas = {
    selected: "My Profile",
    data: [
      {
        name: "My Profile",
      },
      {
        name: "Password",
      },
      {
        name: "Company Type",
      },
      {
        name: "Model TROCR",
      },
    ],
  };

  let companyTypes = [
    {
      id: 1,
      icon: MarketIcon,
      name: "Start-ups and Small Businesses (1-50 employees)",
      desc: "Gain a competitive edge with tools that help you manage your growing business efficiently",
    },
    {
      id: 2,
      icon: OfficeIcon,
      name: "Medium Businesses (51-500 employees)",
      desc: "Enhance coordination and productivity across various departments and teams.",
    },
    {
      id: 3,
      icon: BankIcon,
      name: "Large Enterprises  (>500 employees)",
      desc: "Maximize your impact with streamlined operations and resource management.",
    },
    {
      id: 4,
      icon: UserIcon,
      name: "Freelancers and Consultants",
      desc: "Organize your projects, finances, and client interactions effectively in one place.",
    },
  ];

  let industryType = [
    {
      id: 1,
      icon: SignalIcon,
      name: "Tech",
    },
    {
      id: 2,
      icon: OilGas,
      name: "Oil & Gas",
    },
    {
      id: 3,
      icon: SlidersIcon,
      name: "Other",
    },
  ];

  let passwordObject = {
    currentPassword: "",
    newPassword: "",
    confirmPassword: "",
  };

  function handleDragOver(event) {
    event.preventDefault();
    dropzoneActive = true;
  }

  function handleDragLeave() {
    dropzoneActive = false;
  }

  async function handleDrop(event) {
    event.preventDefault();
    dropzoneActive = false;

    const files = event.dataTransfer.files;

    if (files[0] && (files[0].type === "image/jpeg" || files[0].type === "image/png")) {
      fileImageResult = await new Promise((resolve, reject) => {
        const reader = new FileReader();

        reader.onload = () => resolve(reader.result);
        reader.onerror = reject;

        reader.readAsDataURL(files[0]);
      });

      fileImage = files[0];
      fileImageName = files[0].name;
      // console.log(fileImageResult);
    } else {
      alert("Invalid file. Please select a valid .jpg or .png file.");
    }
  }

  const handleInputFile = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = "image/jpeg, image/png";

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file && (file.type === "image/jpeg" || file.type === "image/png")) {
        fileImageResult = await new Promise((resolve, reject) => {
          const reader = new FileReader();

          reader.onload = () => resolve(reader.result);
          reader.onerror = reject;

          reader.readAsDataURL(file);
        });

        // conver to base64
        fileImage = await readFileAsBase64(file);
        fileImageName = file.name;
        // console.log(fileImageResult);
      } else {
        alert("Invalid file. Please select a valid .jpg or .png file.");
      }
    });

    inputElement.click();
  };

  async function readFileAsBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        const base64Data = reader.result.split(",")[1];
        resolve(base64Data);
      };
      reader.onerror = reject;

      reader.readAsDataURL(file);
    });
  }

  const getUser = async () => {
    let me = JSON.parse(sessionStorage.getItem("user"));
    user = {
      grid: 4,
      title: "My Details",
      titleGrid: 1,
      innerGrid: 5,
      data: [
        {
          name: "Full Name",
          value: me.fullname,
          grid: 2,
        },
        {
          name: "Username",
          value: me.username,
          grid: 2,
        },
        {
          name: "User ID",
          value: me.id,
          grid: 1,
        },
        {
          name: "Email",
          value: me.email,
          grid: 5,
        },
      ],
    };
  };

  const updatePhotoUser = async () => {
    let token = sessionStorage.getItem("token");
    if (token) {
      let me = JSON.parse(sessionStorage.getItem("user"));
      const response = await call.call("client-dashboard", "update-photo-user", {
        URLQuery: {},
        URLParams: {
          id: me.id.toString(),
        },
        Params: {
          photo: fileImage,
          filename: fileImageName,
        },
        Headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        triggerToast("error");
        return;
      }

      triggerToast("updated");
      console.log("ayam", response);
    }
  };

  const currentPasswordCheck = async () => {
    let me = JSON.parse(sessionStorage.getItem("user"));
    var resp = await call.call("client-login", "login", {
      URLQuery: {},
      URLParams: {},
      Params: {
        email: me.email,
        password: passwordObject.currentPassword,
      },
      Headers: {},
    });

    if (resp === null) {
      throw new Error("Server sibuk");
    }

    if (resp.status === 200) {
      // login success
      if (resp.result.token !== undefined) {
        currentPasswordMatch = true;
      } else {
        currentPasswordMatch = false;
      }
    } else {
      currentPasswordMatch = false;
    }
    // currentPasswordMatch = true;
  };

  const passwordCheck = () => {
    const password = passwordObject.newPassword;

    // Check if password is at least 8 characters long
    if (password.length < 8) {
      return false;
    }

    // Check if password contains at least one digit
    if (!/\d/.test(password)) {
      return false;
    }

    // Check if password contains at least one special character
    if (!/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
      return false;
    }

    // Check if password contains at least one lowercase letter
    if (!/[a-z]/.test(password)) {
      return false;
    }

    // Check if password contains at least one uppercase letter
    if (!/[A-Z]/.test(password)) {
      return false;
    }

    // If all checks pass, return true
    return true;
  };

  const updatePassword = async () => {
    if (!passwordCheck()) {
      return;
    }

    if (!currentPasswordMatch) {
      return;
    }

    let me = JSON.parse(sessionStorage.getItem("user"));
    let token = sessionStorage.getItem("token");
    var resp = await call.call("client-dashboard", "update-user", {
      URLQuery: {},
      URLParams: {
        id: me.id.toString(),
      },
      Params: {
        fullname: me.fullname,
        username: me.username,
        password: passwordObject.newPassword,
        email: me.email,
        company: me.company,
        blocked: me.is_active,
      },
      Headers: {
        Authorization: "Bearer " + token,
      },
    });

    console.log("telur", resp);

    if (resp.status !== 200) {
      triggerToast("error");
      console.log(resp);
      return;
    }

    console.log("ayam", resp);

    triggerToast("updated");
  };

  async function getCompanyType() {
    let token = sessionStorage.getItem("token");
    if (token) {
      const response = await call.call("client-dashboard", "get-company-type", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      selectedCompanyType = parseInt(response.result.data[0].attributes.company_size);
      selectedIndustryType = parseInt(response.result.data[0].attributes.company_industry);
    }
  }

  async function updateCompanyType() {
    let token = sessionStorage.getItem("token");
    if (token) {
      const response = await call.call("client-dashboard", "update-company-type", {
        URLQuery: {},
        URLParams: {},
        Params: {
          company_size: selectedCompanyType + "",
          company_industry: selectedIndustryType + "",
        },
        Headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (response.status !== 200) {
        console.log(response);
        return;
      }

      sessionStorage.setItem("industry_id", selectedIndustryType);

      triggerToast("updated");
    }
  }

  async function getModels() {
    const response = await call.call("hexa-ai-rag", "trocr-get-models", {
      URLQuery: {},
      URLParams: {},
      Params: {},
      Headers: {},
    });

    console.log("trocr-get-models", response);

    if (response.status !== 200) {
      console.log(response);
      return;
    }

    options = response.result.models.map((model) => {
      return {
        value: model.name,
        name: model.name,
      };
    });

    selectedModelTrocr = sessionStorage.getItem("model_trocr") || "trocr:latest";
  }

  async function updateModelTrocr() {
    let token = sessionStorage.getItem("token");
    if (token) {
      const response = await call.call("client-dashboard", "update-model-trocr", {
        URLQuery: {},
        URLParams: {},
        Params: {
          model_trocr: selectedModelTrocr,
        },
        Headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }
      triggerToast("updated");
      sessionStorage.setItem("model_trocr", selectedModelTrocr);
    }
  }

  onMount(() => {
    getUser();
    getCompanyType();
    getModels();
  });
</script>

<svelte:window bind:innerWidth={screenWidth} />

{#if Object.keys(user).length !== 0}
  <div class="m-4">
    <!-- Header -->
    <div class="mt-4 flex justify-between">
      <P size="3xl" weight="bold">User Menu</P>
      <!-- <P size="3xl" weight="bold">{isAdmin && "Super"} User Menu</P> -->
      <div class="relative">
        <!-- <Button
          class="bg-white font-bold mr-2 border border-hc-graylight-300 hover:bg-hc-graylight-50 text-gray-900 dark:bg-white dark:hover:bg-hc-graylight-50"
          disabled={JSON.stringify(user) === JSON.stringify(userCopy)}
          on:click={getUser}>Cancel</Button
        >
        <Button
          class="bg-hc-brand-600 hover:bg-hc-brand-700"
          disabled={JSON.stringify(user) === JSON.stringify(userCopy)}
          on:click={updatePhotoUser}>Save</Button
        > -->
        <div class="flex absolute inset-y-0 start-0 items-center ps-3 pointer-events-none">
          <SearchOutline class="w-4 h-4 text-primary-light dark:text-primary-dark" />
        </div>
        <Input class="ps-10 bg-primary-light dark:bg-primary-dark" placeholder="Search" />
      </div>
    </div>

    <!-- Submenu -->
    <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

    <!-- Content -->
    {#if isAdmin}
      {#if tabDatas.selected === "My Profile"}
        <div class="mx-4 mt-4">
          <div class="pt-4 pb-8 border-b border-secondary-light dark:border-secondary-dark">
            <Typography preset="Text lg/Semibold" text="Personal info" />
            <Typography
              preset="Text sm/Reguler"
              text="Update your photo and personal details here."
            />
          </div>
          <div class="py-8 border-b border-secondary-light dark:border-secondary-dark">
            <InformationForm datas={user} {screenWidth} />
          </div>
          <div class="py-8">
            <div class="grid grid-cols-4">
              <div class="col-span-1">
                <Typography preset="Text lg/Semibold" text="Your photo" />
              </div>
              <div class="col-span-3 flex w-full">
                <Avatar class="size-16" src={fileImageName ? fileImageResult : ""} />
                <div class="w-full ml-8">
                  <Button
                    ondragover={handleDragOver}
                    ondragleave={handleDragLeave}
                    ondrop={handleDrop}
                    onclick={handleInputFile}
                    color="light"
                    class={`w-[50%] rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1 py-8 flex flex-col justify-start items-center`}
                  >
                    <div
                      class="rounded-xl shadow border border-primary-light dark:border-primary-dark justify-center items-center inline-flex mb-4 pointer-events-none"
                    >
                      <img
                        src={upload}
                        alt="mail.svg"
                        class="p-4 dark:invert pointer-events-none"
                      />
                    </div>
                    <div class="flex space-x-1">
                      <P class="text-hc-brand-700 pointer-events-none">Click to upload</P><P
                        class="pointer-events-none">or drag and drop</P
                      >
                    </div>
                    <P class="pointer-events-none">PNG, JPG file (max. 800x400px)</P>
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>
      {:else if tabDatas.selected === "Password"}
        <div class="mx-4 mt-4">
          <!-- <DataTable {datas} {configs} /> -->
          <div class="pt-4 pb-8 border-b border-secondary-light dark:border-secondary-dark">
            <Typography preset="Text lg/Semibold" text="Password" />
            <Typography
              preset="Text sm/Reguler"
              text="Please enter your current password to change your password."
            />
          </div>
          <div class="py-8">
            <div class="flex flex-col w-full">
              <Typography preset="Text sm/Semibold" text="Current Password" />
              <Input
                type="password"
                bind:value={passwordObject.currentPassword}
                on:input={currentPasswordCheck}
                class="mt-2 bg-primary-light dark:bg-primary-dark"
              />
              {#if !currentPasswordMatch && passwordObject.currentPassword.length > 0}
                <Typography
                  preset="Text sm/Reguler"
                  text="Your current password did not match with our records. Please try again."
                />
              {/if}
            </div>
            <div class="flex flex-col w-full mt-4">
              <Typography preset="Text sm/Semibold" text="New Password" />
              <Input
                type="password"
                disabled={!currentPasswordMatch}
                bind:value={passwordObject.newPassword}
                on:change={() => {}}
                class="mt-2 bg-primary-light dark:bg-primary-dark"
              />
              {#if !passwordCheck() && passwordObject.newPassword.length > 0}
                <Typography
                  preset="Text sm/Reguler"
                  text="Your new password must be at least than 8 characters, contains at least one digit number, contains at least one special character, contains at least one lowercase letter, and contains at least one uppercase letter."
                />
              {/if}
            </div>
            <div class="flex flex-col w-full mt-4">
              <Typography preset="Text sm/Semibold" text="Confirm Password" />
              <Input
                type="password"
                disabled={!currentPasswordMatch}
                bind:value={passwordObject.confirmPassword}
                class="mt-2 bg-primary-light dark:bg-primary-dark"
              />
            </div>
          </div>
        </div>
      {:else if tabDatas.selected === "Company Type"}
        <div class="mx-4 mt-4">
          <div class="pt-4 pb-5 border-b border-secondary-light dark:border-secondary-dark">
            <Typography preset="Text lg/Semibold" text="Company Size" />
            <Typography preset="Text sm/Reguler" text="Please select your companies current size" />
          </div>
          <Divider />
          <div class="pb-4"></div>
          <SelectedCard
            selectionData={companyTypes}
            bind:selectedData={selectedCompanyType}
            numGrid={2}
            isRadio
          />
          <div class="pt-8 pb-5 border-b border-secondary-light dark:border-secondary-dark">
            <Typography preset="Text lg/Semibold" text="Industry" />
            <Typography
              preset="Text sm/Reguler"
              text="Please select your companies current industry"
            />
          </div>
          <Divider />
          <div class="pb-4"></div>
          <SelectedCard
            selectionData={industryType}
            bind:selectedData={selectedIndustryType}
            isFlex
            isRadio
            widthFlex="w-max"
          />
        </div>
      {:else if tabDatas.selected === "Model TROCR"}
        <div class="mx-4 mt-4">
          <div class="pt-4 pb-8 border-b border-secondary-light dark:border-secondary-dark">
            <Typography preset="Text lg/Semibold" text="Model TROCR" />
            <Select items={options} class="display-block" bind:value={selectedModelTrocr} />
          </div>
        </div>
      {/if}

      <!-- action buttons -->
      <div class="flex w-full justify-end">
        {#if tabDatas.selected !== "Model TROCR"}
          <Button
            class="bg-white font-bold mr-2 border border-hc-graylight-300 hover:bg-hc-graylight-50 text-gray-900 dark:bg-white dark:hover:bg-hc-graylight-50"
            on:click={() => {
              if (tabDatas.selected === "My Profile") {
                getUser();
              } else if (tabDatas.selected === "Password") {
                passwordObject = {
                  currentPassword: "",
                  newPassword: "",
                  confirmPassword: "",
                };
              }
            }}>Cancel</Button
          >
        {/if}
        <Button
          class="bg-hc-brand-600 hover:bg-hc-brand-700"
          on:click={() => {
            if (tabDatas.selected === "My Profile") {
              updatePhotoUser();
            } else if (tabDatas.selected === "Password") {
              updatePassword();
            } else if (tabDatas.selected === "Company Type") {
              updateCompanyType();
            } else if (tabDatas.selected === "Model TROCR") {
              updateModelTrocr();
            }
          }}
        >
          {#if tabDatas.selected === "My Profile" || tabDatas.selected === "Company Type" || tabDatas.selected === "Model TROCR"}
            Save
          {:else if tabDatas.selected === "Password"}
            Update Password
          {/if}
        </Button>
      </div>
    {:else}
      hello user
    {/if}
  </div>
{/if}
