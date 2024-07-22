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
  import { TabBar, InformationForm, Typography } from "$lib";
  import { SearchOutline } from "flowbite-svelte-icons";
  import DataTable from "$lib/components/DataTable.svelte";
  import upload from "$lib/icon/upload-cloud-02.svg";
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
        name: "Model TROCR",
      },
    ],
  };
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
        return;
      }
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
      console.log(resp);
      return;
    }

    console.log("ayam", resp);
  };

  const changeTabbar = (name) => {
    tabDatas.selected = name;
  };

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
  }

  async function updateModelTrocr() {
    let token = sessionStorage.getItem("token");
    if (token) {
      let me = JSON.parse(sessionStorage.getItem("user"));
      const response = await call.call("client-dashboard", "update-user", {
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
      console.log("ayam", response);
      sessionStorage.setItem("model_trocr", selectedModelTrocr);
    }
  }

  onMount(() => {
    getUser();
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
    <TabBar {tabDatas} {changeTabbar} />

    <!-- Content -->
    {#if isAdmin}
      {#if tabDatas.selected === "My Profile"}
        <div class="mx-4 mt-8">
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
        <div class="mx-4 mt-8">
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
      {:else if tabDatas.selected === "Model TROCR"}
        <div class="mx-4 mt-8">
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
            } else if (tabDatas.selected === "Model TROCR") {
              updateModelTrocr();
            }
          }}
        >
          {#if tabDatas.selected === "My Profile" || tabDatas.selected === "Model TROCR"}
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
