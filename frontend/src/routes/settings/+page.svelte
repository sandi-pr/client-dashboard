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
  import { P, Input, Avatar, Button, Select, Indicator } from "flowbite-svelte";
  import { SelectedCard, InformationForm, Typography, triggerToast, InputRw, Dropdown } from "$lib";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { SearchOutline } from "flowbite-svelte-icons";
  import DataTable from "$lib/components/DataTable.svelte";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import MarketIcon from "$lib/icon/building-02.svelte";
  import OfficeIcon from "$lib/icon/building-06.svelte";
  import BankIcon from "$lib/icon/bank.svelte";
  import BanksIcon from "$lib/icon/banks.svelte";
  import UserIcon from "$lib/icon/user-03.svelte";
  import SignalIcon from "$lib/icon/signal-01.svelte";
  import OilGas from "$lib/icon/oil-gas.svelte";
  import BoxesIcon from "$lib/icon/boxes.svelte";
  import EducationIcon from "$lib/icon/education.svelte";
  import FinancialServicesIcon from "$lib/icon/financial-services.svelte";
  import HealthCareIcon from "$lib/icon/health-care.svelte";
  import HospitalityIcon from "$lib/icon/hospitality.svelte";
  import ManufactureIcon from "$lib/icon/manufacture.svelte";
  import RealEstateIcon from "$lib/icon/real-estate.svelte";
  import RetailIcon from "$lib/icon/retail.svelte";
  import TechIcon from "$lib/icon/tech.svelte";
  import TelecomIcon from "$lib/icon/telecom.svelte";
  import SlidersIcon from "$lib/icon/sliders-04.svelte";
  import { datas, informationFormDatas, configs } from "./dummy";
  import { setPrimaryColor } from "$lib/utils/colorTheme";
  import { purpleColors, greenColors, blueColors } from "$lib/utils/presetColors";
  import { CR3Call } from "hexacomponents-svelte";
  import { _ } from "svelte-i18n";

  import { call } from "$lib";

  let screenWidth;
  let user = informationFormDatas;
  let userCopy = {};
  let isAdmin = true;
  let sidebarTitle = null;
  let dropdownThemeOpen = false;
  let dropdownModelTrocrOpen = false;
  let selectedTheme = "purple";
  let themes = [
    {
      name: "Purple",
      value: "purple",
    },
    {
      name: "Green",
      value: "green",
    },
    {
      name: "Blue",
      value: "blue",
    },
  ];
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
        name: "Theme",
      },
      // {
      //   name: "Configuration AI",
      // },
    ],
  };
  let onEditPassword = false;
  let onEditCompanyType = false;
  let onEditTheme = false;

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
      icon: BanksIcon,
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
      icon: TechIcon,
      name: "Tech",
    },
    {
      id: 2,
      icon: BankIcon,
      name: "Bank",
    },
    {
      id: 3,
      icon: OilGas,
      name: "Oil & Gas",
    },
    {
      id: 4,
      icon: RetailIcon,
      name: "Retail",
    },
    {
      id: 5,
      icon: HealthCareIcon,
      name: "Healthcare",
    },
    {
      id: 6,
      icon: EducationIcon,
      name: "Education",
    },
    {
      id: 7,
      icon: FinancialServicesIcon,
      name: "Financial Services",
    },
    {
      id: 8,
      icon: HospitalityIcon,
      name: "Hospitality",
    },
    {
      id: 9,
      icon: ManufactureIcon,
      name: "Manufacturing",
    },
    {
      id: 10,
      icon: BoxesIcon,
      name: "Logistics",
    },
    {
      id: 11,
      icon: TelecomIcon,
      name: "Telecommunication",
    },
    {
      id: 12,
      icon: RealEstateIcon,
      name: "Real Estate",
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

      fileImage = await readFileAsBase64(files[0]);
      fileImageName = files[0].name;
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
    let me = JSON.parse(localStorage.getItem("user"));
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
          name: "User Id",
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
    let token = localStorage.getItem("token");
    if (token) {
      let me = JSON.parse(localStorage.getItem("user"));
      const response = await $call.call("client-dashboard", "update-photo-user", {
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

      fileImage = null;
      fileImageResult = null;
      fileImageName = null;

      triggerToast("updated");
    }
  };

  const currentPasswordCheck = async () => {
    let me = JSON.parse(localStorage.getItem("user"));
    var resp = await $call.call("client-login", "login", {
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

    let me = JSON.parse(localStorage.getItem("user"));
    let token = localStorage.getItem("token");
    var resp = await $call.call("client-dashboard", "update-user", {
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

    if (resp.status !== 200) {
      triggerToast("error");
      console.log(resp);
      return;
    }

    triggerToast("updated");
  };

  async function getCompanyType() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "get-company-type", {
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

      selectedCompanyType = parseInt(response.result.data.attributes.company_size);
      selectedIndustryType = parseInt(response.result.data.attributes.company_industry);

      sidebarTitle = response.result.data.attributes.sidebar_title || "Hexapod";

      if (response.result.data.attributes.company_logo.data !== null) {
        let ext = response.result.data.attributes.company_logo.data.attributes.ext;

        let fileType;
        if (ext === ".png") {
          fileType = "image/png";
        } else if (ext === ".jpg") {
          fileType = "image/jpg";
        } else if (ext === ".jpeg") {
          fileType = "image/jpeg";
        } else if (ext === ".svg") {
          fileType = "image/svg+xml";
        } else {
          fileType = "image/png";
        }

        fileImageName = response.result.data.attributes.company_logo.data?.attributes.name;
        if (response.result.company_logo !== "") {
          fileImage = response.result.company_logo;
          let blob = base64toBlob(fileImage, fileType);
          fileImageResult = URL.createObjectURL(blob);
        }
      }

      selectedTheme = response.result.data.attributes.theme_color || "purple";
    }
  }

  function base64toBlob(base64Data, fileType) {
    const byteCharacters = atob(base64Data);
    const byteArrays = [];

    for (let offset = 0; offset < byteCharacters.length; offset += 512) {
      const slice = byteCharacters.slice(offset, offset + 512);

      const byteNumbers = new Array(slice.length);
      for (let i = 0; i < slice.length; i++) {
        byteNumbers[i] = slice.charCodeAt(i);
      }

      const byteArray = new Uint8Array(byteNumbers);
      byteArrays.push(byteArray);
    }

    const blob = new Blob(byteArrays, { type: fileType });
    return blob;
  }

  async function updateCompanyType() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "update-company-type", {
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

      window.location.reload();
    }
  }

  async function updateTheme() {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "update-theme", {
          URLQuery: {},
          URLParams: {},
          Params: {
            sidebar_title: sidebarTitle + "",
            theme_color: selectedTheme + "",
            company_logo: JSON.stringify({
              base64Data: fileImage,
              name: fileImageName,
            }),
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

        window.location.reload();
      } catch (error) {
        console.log(error);
      }
    }
  }

  // $: {
  //   if (selectedTheme) {
  //     updateTheme(selectedTheme);
  //   }
  // }

  async function getModels() {
    try {
      const response = await $call.call("hexa-ai-rag", "trocr-get-models", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {},
      });

      if (response === null) {
        options = [
          {
            value: "trocr:latest",
            name: "trocr:latest",
          },
        ];

        selectedModelTrocr = sessionStorage.getItem("model_trocr") || "trocr:latest";
        return;
      }

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
    } catch (error) {
      console.log(error);
    }
  }

  async function updateModelTrocr() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "update-model-trocr", {
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
        <!-- <div class="hidden flex absolute inset-y-0 start-0 items-center ps-3 pointer-events-none">
          <SearchOutline class="w-4 h-4 text-primary-light dark:text-primary-dark" />
        </div>
        <Input class="ps-10 bg-primary-light dark:bg-primary-dark" placeholder="Search" /> -->
      </div>
    </div>

    <!-- Submenu -->
    <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

    <!-- Content -->
    {#if isAdmin}
      {#if tabDatas.selected === "My Profile"}
        <div class="mx-4 mt-4">
          <div class="pt-4 pb-5 border-b border-secondary-light dark:border-secondary-dark">
            <Typography preset="Text lg/Semibold" text={$_("personalInfo")} />
            <Typography preset="Text sm/Reguler" text={$_("descPersonalInfo")} />
          </div>
          <div class="py-4 border-b border-secondary-light dark:border-secondary-dark">
            <InformationForm datas={user} {screenWidth} />
          </div>
          <!-- <div class="py-8">
            <div class="grid grid-cols-4">
              <div class="col-span-1">
                <Typography preset="Text lg/Semibold" text="Your photo" />
              </div>
              <div class="col-span-3 flex w-full">
                <Avatar class="size-16" src={fileImageName ? fileImageResult : ""} />
                <div class="grow ml-8">
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
          </div> -->
        </div>
      {:else if tabDatas.selected === "Password"}
        <div class="mx-4 mt-4">
          <!-- <DataTable {datas} {configs} /> -->
          <div class="pt-4 pb-5 border-b border-secondary-light dark:border-secondary-dark">
            <Typography preset="Text lg/Semibold" text={$_("changePassword")} />
            <Typography preset="Text sm/Reguler" text={$_("descPassword")} />
          </div>
          <div class="flex flex-col py-4 gap-4">
            <div class="flex flex-col w-full gap-1.5">
              <Typography preset="Text sm/Semibold" text={$_("currentPassword")} />
              <InputRw type="password" bind:value={passwordObject.currentPassword} onEdit={onEditPassword} oninput={currentPasswordCheck} />
              {#if !currentPasswordMatch && passwordObject.currentPassword.length > 0}
                <Typography preset="Text sm/Reguler" text={$_("passwordNotMatch")} color="warning-primary" />
              {/if}
            </div>
            <div class="flex flex-col w-full gap-1.5">
              <Typography preset="Text sm/Semibold" text={$_("newPassword")} />
              <InputRw type="password" bind:value={passwordObject.newPassword} onEdit={onEditPassword} disabled={!currentPasswordMatch} />
              {#if !passwordCheck() && passwordObject.newPassword.length > 0}
                <Typography preset="Text sm/Reguler" text={$_("hintPassword")} />
              {/if}
            </div>
            <div class="flex flex-col w-full gap-1.5">
              <Typography preset="Text sm/Semibold" text={$_("confirmPassword")} />
              <InputRw type="password" bind:value={passwordObject.confirmPassword} onEdit={onEditPassword} disabled={!currentPasswordMatch} />
            </div>
          </div>
          <div class="flex w-full justify-end">
            {#if onEditPassword}
              <Button
                class="bg-white font-bold mr-2 border border-hc-graylight-300 hover:bg-hc-graylight-50 text-gray-900 dark:bg-white dark:hover:bg-hc-graylight-50"
                on:click={() => {
                  passwordObject = {
                    currentPassword: "",
                    newPassword: "",
                    confirmPassword: "",
                  };
                  currentPasswordMatch = false;
                  onEditPassword = false;
                }}>{$_("cancel")}</Button
              >
              <Button
                class="bg-hc-brand-600 hover:bg-hc-brand-700"
                on:click={() => {updatePassword()}}
              >
                {$_("save")}
              </Button>
            {:else}
              <Button
                onclick={() => {
                  onEditPassword = true;
                }}
                class="bg-hc-brand-600 hover:bg-hc-brand-700">{$_("changePassword")}</Button
              >
            {/if}
          </div>
        </div>
      {:else if tabDatas.selected === "Company Type"}
        <div class="mx-4 mt-4">
          <div class="pt-4 pb-5 border-b border-secondary-light dark:border-secondary-dark">
            <Typography preset="Text lg/Semibold" text={$_("companySize")} />
            <Typography preset="Text sm/Reguler" text={$_("descCompanySize")} />
          </div>
          <div class="flex flex-col py-4 gap-4">
            <SelectedCard
              disabled={!onEditCompanyType}
              selectionData={companyTypes}
              bind:selectedData={selectedCompanyType}
              numGrid={2}
              isRadio
            />
            <div class="pt-4 pb-4 border-b border-secondary-light dark:border-secondary-dark">
              <Typography preset="Text lg/Semibold" text={$_("industryType")} />
              <Typography preset="Text sm/Reguler" text={$_("decIndustry")} />
            </div>
            <SelectedCard
              disabled={!onEditCompanyType}
              selectionData={industryType}
              bind:selectedData={selectedIndustryType}
              isFlex
              isRadio
              widthFlex="w-max"
            />
          </div>
          <div class="flex w-full justify-end">
            {#if onEditCompanyType}
              <Button
                class="bg-white font-bold mr-2 border border-hc-graylight-300 hover:bg-hc-graylight-50 text-gray-900 dark:bg-white dark:hover:bg-hc-graylight-50"
                on:click={() => {
                  onEditCompanyType = false;
                }}>{$_("cancel")}</Button
              >
              <Button
                class="bg-hc-brand-600 hover:bg-hc-brand-700"
                on:click={() => {updateCompanyType();}}
              >
                {$_("save")}
              </Button>
            {:else}
              <Button
                onclick={() => {
                  onEditCompanyType = true;
                }}
                class="bg-hc-brand-600 hover:bg-hc-brand-700">{$_("editCompany")}</Button
              >
            {/if}
          </div>
        </div>
      {:else if tabDatas.selected === "Theme"}
        <div class="mx-4 mt-4">
          <div class="pt-4 pb-8 border-b border-secondary-light dark:border-secondary-dark">
            <div class="grid grid-cols-4">
              <div class="col-span-1">
                <Typography preset="Text lg/Semibold" text={$_("sidebarTitle")} />
              </div>
              <div class="col-span-3 flex w-full">
                <InputRw bind:value={sidebarTitle} onEdit={onEditTheme} />
              </div>
            </div>
          </div>
          <div class="py-8 border-b border-secondary-light dark:border-secondary-dark">
            <div class="grid grid-cols-4">
              <div class="col-span-1">
                <Typography preset="Text lg/Semibold" text={$_("companyLogo")} />
              </div>
              <div class="col-span-3 flex w-full">
                <Avatar class="size-12" src={fileImageName ? fileImageResult : ""} />
                <div class="grow ml-3">
                  <Button
                    disabled={!onEditTheme}
                    ondragover={handleDragOver}
                    ondragleave={handleDragLeave}
                    ondrop={handleDrop}
                    onclick={handleInputFile}
                    color="light"
                    class={`w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark py-4 flex flex-col justify-start items-center`}
                  >
                    <div
                      class="rounded-xl shadow border border-primary-light dark:border-primary-dark justify-center items-center inline-flex mb-3 pointer-events-none"
                    >
                      <img
                        src={upload}
                        alt="mail.svg"
                        class="p-4 dark:invert pointer-events-none"
                      />
                    </div>
                    <div class="flex space-x-1">
                      <P class="text-hc-brand-700 pointer-events-none">{$_("clickToUpload")}</P><P
                        class="pointer-events-none">{$_("orDragAndDrop")}</P
                      >
                    </div>
                    <P class="pointer-events-none">PNG or SVG file (max. 1MB)</P>
                  </Button>
                </div>
              </div>
            </div>
          </div>
          <div class="py-8">
            <div class="grid grid-cols-4">
              <div class="col-span-1">
                <Typography preset="Text lg/Semibold" text={$_("themeColors")} />
              </div>
              <div class="col-span-3 flex w-1/3">
                <Dropdown
                  dropdownOpen={dropdownThemeOpen}
                  bind:selectedItem={selectedTheme}
                  items={themes}
                  isOverlay={true}
                  widget={Indicator}
                  disabled={!onEditTheme}
                />
              </div>
            </div>
          </div>
          <div class="flex w-full justify-end">
            {#if onEditTheme}
              <Button
                class="bg-white font-bold mr-2 border border-hc-graylight-300 hover:bg-hc-graylight-50 text-gray-900 dark:bg-white dark:hover:bg-hc-graylight-50"
                on:click={() => {
                  onEditTheme = false;
                }}>{$_("cancel")}</Button
              >
              <Button
                class="bg-hc-brand-600 hover:bg-hc-brand-700"
                on:click={() => {updateTheme();}}
              >
                {$_("save")}
              </Button>
            {:else}
              <Button
                onclick={() => {
                  onEditTheme = true;
                }}
                class="bg-hc-brand-600 hover:bg-hc-brand-700">{$_("editTheme")}</Button
              >
            {/if}
          </div>
        </div>
      {:else if tabDatas.selected === "Configuration AI"}
        <div class="mx-4 mt-4">
          <div class="pt-4 pb-8 border-b border-secondary-light dark:border-secondary-dark">
            <Typography preset="Text lg/Semibold" text={$_("configurationAI")} />
            <Typography preset="Text sm/Reguler" text={$_("descConfigurationAI")} />
          </div>
          <Divider />
          <div class="py-8">
            <div class="grid grid-cols-4">
              <div class="col-span-1">
                <Typography preset="Text lg/Semibold" text="Model TROCR" />
              </div>
              <div class="col-span-3 flex w-1/3">
                <Dropdown
                  dropdownOpen={dropdownModelTrocrOpen}
                  bind:selectedItem={selectedModelTrocr}
                  items={options}
                  isOverlay={true}
                />
              </div>
            </div>
          </div>
        </div>
      {/if}

      <!-- action buttons -->
      <!-- {#if tabDatas.selected !== "My Profile"}
        <div class="flex w-full px-4 justify-end">
          {#if tabDatas.selected !== "Configuration AI"}
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
            disabled={tabDatas.selected === "Theme" ? (sidebarTitle === "" ? true : false) : false}
            class="bg-hc-brand-600 hover:bg-hc-brand-700"
            on:click={() => {
              if (tabDatas.selected === "My Profile") {
                updatePhotoUser();
              } else if (tabDatas.selected === "Password") {
                updatePassword();
              } else if (tabDatas.selected === "Company Type") {
                updateCompanyType();
              } else if (tabDatas.selected === "Configuration AI") {
                updateModelTrocr();
              } else if (tabDatas.selected === "Theme") {
                updateTheme();
              }
            }}
          >
            {#if tabDatas.selected === "My Profile" || tabDatas.selected === "Company Type" || tabDatas.selected === "Configuration AI" || tabDatas.selected === "Theme"}
              Save
            {:else if tabDatas.selected === "Password"}
              Update Password
            {/if}
          </Button>
        </div>
      {/if} -->
    {:else}
      hello user
    {/if}
  </div>
{/if}
