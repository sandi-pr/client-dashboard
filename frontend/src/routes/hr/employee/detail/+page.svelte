<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: 
  - Azzis Arswendo <azzis@hexacode.org>
  - Dean Ramhan <dean@hexacode.org>
  - Taqie Fadlillah <taqief@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->
<script>
  import {
    P,
    Heading,
    Input,
    Button,
    Select,
    Helper,
    Toast,
    Modal,
    Toggle,
  } from "flowbite-svelte";
  import {
    CheckCircleSolid,
    ExclamationCircleSolid,
  } from "flowbite-svelte-icons";
  import { ToastContainer, triggerToast } from "$lib";
  import { fly } from "svelte/transition";
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { onMount, createEventDispatcher } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { goto } from "$app/navigation";

  const dispatch = createEventDispatcher();

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  const employeeId = params.get("employeeId");
  const edit = params.get("edit");
  let onEdit = $state(false);

  let data = $state({
    employeeId: "",
    employee_name: "",
    date_of_birth: "",
    email: "",
    phone: "",
    country: "",
    address: "",
    city: "",
    state: "",
    zip_code: "",
    is_active: false,
    job_position: "",
    branch: "",
    employment_status: "",
    start_date: "",
    end_date: "",
  });

  let emailValid = $state(true);

  let popupModal = $state(false);
  let jobLevels = $state([]);

  async function getJobLevels() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "get-job-level", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          let promise = response.result.job_level.map((job) => {
            return { name: job.title, value: job.title };
          });
          jobLevels = await Promise.all(promise);
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  function formatDate(inputDate) {
    const date = new Date(inputDate);

    const formattedDate = `${(date.getMonth() + 1).toString().padStart(2, "0")}/${date
      .getDate()
      .toString()
      .padStart(2, "0")}/${date.getFullYear()}, ${date
      .getHours()
      .toString()
      .padStart(2, "0")}:${date
      .getMinutes()
      .toString()
      .padStart(2, "0")} ${date.getHours() >= 12 ? "PM" : "AM"}`;

    return formattedDate;
  }

  function validateEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  function handleEmailChange(event) {
    const newEmail = event.target.value;
    emailValid = validateEmail(newEmail);
    let email = newEmail;
    dispatch("emailChange", { email, emailValid });
  }

  async function getEmployee() {
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "get-employee", {
        URLQuery: {},
        URLParams: { id: employeeId },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status === 200) {
        data = {
          employeeId: response.result.data[0].id,
          employee_name: response.result.data[0].attributes.employee_name,
          date_of_birth: response.result.data[0].attributes.date_of_birth,
          email: response.result.data[0].attributes.email,
          phone: response.result.data[0].attributes.phone,
          country: response.result.data[0].attributes.country,
          address: response.result.data[0].attributes.address,
          city: response.result.data[0].attributes.city,
          state: response.result.data[0].attributes.state,
          zip_code: response.result.data[0].attributes.zip_code,
          createdAt: formatDate(response.result.data[0].attributes.createdAt),
          is_active:
            response.result.data[0].attributes.is_active === "Yes"
              ? true
              : false,
          job_position: response.result.data[0].attributes.job_position,
          branch: response.result.data[0].attributes.branch,
          employment_status:
            response.result.data[0].attributes.employment_status,
          start_date: response.result.data[0].attributes.start_date,
          end_date: response.result.data[0].attributes.end_date,
        };
      }
    } catch (error) {
      triggerToast("error");
      console.log(error);
    }
  }

  async function updateEmployee() {
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "update-employee", {
        URLQuery: {},
        URLParams: { id: employeeId },
        Params: {
          employee_name: data.employee_name,
          date_of_birth: data.date_of_birth,
          email: data.email,
          phone: data.phone.toString(),
          country: data.country,
          address: data.address,
          city: data.city,
          state: data.state,
          zip_code: data.zip_code.toString(),
          is_active: data.is_active === true ? "Yes" : "No",
          job_position: data.job_position,
          branch: data.branch,
          employment_status: data.employment_status,
          start_date: data.start_date,
          end_date: data.end_date,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status === 200) {
        triggerToast("updated");
        // setTimeout(async () => {
          goto(await call.getPage("client-dashboard", "hr-employee"));
        // }, 3000);
        return;
      }
    } catch (error) {
      triggerToast("error");
      console.log(error);
    }
  }

  async function deleteEmployee() {
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "delete-employee", {
        URLQuery: {},
        URLParams: { id: employeeId },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });
      if (response.status === 200) {
        triggerToast("deleted");
        // setTimeout(async () => {
          goto(await call.getPage("client-dashboard", "hr-employee"));
        // }, 3000);
        return;
      }
    } catch (error) {
      triggerToast("error");
      console.log(error);
    }
  }

  async function backPage() {
    goto(await call.getPage("client-dashboard", "hr-employee"));
    return;
  }

  onMount(async () => {
    if (edit === "true") onEdit = true;
    getEmployee();
    getJobLevels();
  });

  // let dropdownOpen = false;
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex items-center pb-8">
        <div class="flex flex-col gap-2 grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {data.employee_name}
              <Button
                outline
                color="red"
                class="hover:!bg-transparent !border-hc-graylight-200 dark:!border-hc-graydark-600 focus-within:!ring-1 !p-2"
              >
                <CopyIcon
                  svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
                />
              </Button>
            </Heading>
          </div>
          <P
            class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
          >
            {$_("createOn")}{` ${data.createdAt}`}
          </P>
        </div>
        <div class="flex gap-2">
          {#if onEdit}
            <Button
              onclick={() => {
                onEdit = false;
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button onclick={updateEmployee} class="px-3 focus-within:!ring-1"
              >{$_("save")}</Button
            >
          {:else}
            <Button
              onclick={backPage}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("back")}</Button
            >
            <Button
              onclick={() => (popupModal = true)}
              outline
              color="red"
              class="px-3 focus-within:!ring-1">{$_("delete")}</Button
            >
            <Button
              onclick={() => {
                onEdit = true;
              }}
              class="px-3 focus-within:!ring-1">{$_("edit")}</Button
            >
            <!-- <Button
              onclick={() => {
                onEdit = true;
              }}
              class="px-3 focus-within:!ring-1">{$_("edit")}</Button
            > -->
          {/if}
        </div>
      </div>

      <Divider />
      <div class="flex flex-col gap-6 pt-8">
        <!-- general information -->
        <div class="flex gap-6">
          <div class="w-60">
            <P
              class="font-semibold text-secondary-light dark:text-secondary-dark"
            >
              {$_("generalInformation")}
            </P>
          </div>
          <div class="flex flex-col grow gap-6">
            <div class="flex gap-5">
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("employeeName")}
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.employee_name}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("employeeId")}
                </P>
                <Input
                  readonly
                  bind:value={data.employeeId}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("email")}</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.email}
                  on:input={handleEmailChange}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
                {#if onEdit}
                  {#if !emailValid && data.email.length > 0}
                    <Helper class="text-sm !text-red-400 pt-2"
                      >{$_("emailNotValid")}</Helper
                    >
                  {/if}
                {/if}
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("phone")}</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.phone}
                  type="number"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="flex flex-col grow gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
              >
                {$_("address")}
              </P>
              <Input
                readonly={!onEdit}
                bind:value={data.address}
                class="bg-primary-light dark:bg-secondary-dark"
              />
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("city")}</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.city}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("state")}</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.state}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("zipCode")}</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.zip_code}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("country")}</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.country}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("dateOfBirth")}</P
                >
                <Input
                  readonly={!onEdit}
                  type="date"
                  bind:value={data.date_of_birth}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="pb-2 text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("active")}</P
                >
                <!-- <Input readonly={!onEdit} bind:value={data.active} /> -->
                <Toggle disabled={!onEdit} bind:checked={data.is_active}
                ></Toggle>
                <!-- <Select
                  readonly={!onEdit}
                  class="bg-primary-light dark:bg-secondary-dark"
                  placeholder={$_("chooseOption")}
                  items={[
                    { value: "Yes", name: "Yes" },
                    { value: "No", name: "No" },
                  ]}
                  bind:value={data.is_active}
                /> -->
                <!-- <Dropdown
                    placement="bottom-start"
                    bind:open={dropdownOpen}
                    class="w-48 overflow-y-auto py-1 dark:bg-secondary-dark"
                  >
                    <DropdownItem on:click={() => (dropdownOpen = false)}>true</DropdownItem>
                    <DropdownItem on:click={() => (dropdownOpen = false)}>false</DropdownItem>
                  </Dropdown> -->
              </div>
            </div>
          </div>
        </div>
        <Divider />
        <div class="flex gap-6 pt-2">
          <div class="w-60">
            <P
              class="font-semibold text-secondary-light dark:text-secondary-dark"
            >
              {$_("employeeInformation")}
            </P>
          </div>
          <div class="flex flex-col grow gap-6">
            <div class="grid grid-cols-3 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("jobPosition")}*</P
                >
                <Select
                  disabled={!onEdit}
                  class="bg-primary-light dark:bg-secondary-dark"
                  items={jobLevels}
                  bind:value={data.job_position}
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("branch")}*</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.branch}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("employmentStatus")}*</P
                >
                <Input
                  readonly={!onEdit}
                  bind:value={data.employment_status}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("joinDate")}*</P
                >
                <Input
                  readonly={!onEdit}
                  type="date"
                  bind:value={data.start_date}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("endDate")}*</P
                >
                <Input
                  readonly={!onEdit}
                  type="date"
                  bind:value={data.end_date}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <ToastContainer />
</div>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteEmployee} color="red" class="me-2"
      >{$_("yesImSure")}</Button
    >
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>
