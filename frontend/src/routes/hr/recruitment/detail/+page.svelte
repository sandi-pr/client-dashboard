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
    Textarea,
    Button,
    Modal,
    MultiSelect,
    Select,
  } from "flowbite-svelte";
  import {
    CheckCircleSolid,
    ExclamationCircleSolid,
  } from "flowbite-svelte-icons";
  import { ToastContainer, triggerToast } from "$lib";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import DownloadIcon from "$lib/icon/file-download-03.svg";
  import Divider from "$lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { format } from "date-fns";
  import { goto } from "$app/navigation";
  import { createEventDispatcher } from "svelte";
  import Badge from "$lib/ui/badge/Badge.svelte";

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  const recruitmentId = params.get("recruitmentId");
  const edit = params.get("edit");
  const dispatch = createEventDispatcher();

  let onEdit = $state(false);

  let popupModal = $state(false);

  let contentFile = $state(null);
  let base64 = $state(null);
  let fileUrl = $state("null");

  let tabDatas = $state({
    selected: "Resume",
    data: [
      {
        name: "Resume",
      },
      {
        name: "Information",
      },
    ],
  });

  let pdfReady = $state(false);

  let data = $state({});

  let categories = $state([
    {
      value: "GENERAL",
      name: "General",
    },
    {
      value: "CONTRACT",
      name: "Contract",
    },
  ]);

  let jobLevels = $state([]);

  let skillsOrKeywords = $state([
    {
      value: "Figma",
      name: "Figma",
    },
    {
      value: "JavaScript",
      name: "JavaScript",
    },
    {
      value: "HTML",
      name: "HTML",
    },
    {
      value: "CSS",
      name: "CSS",
    },
    {
      value: "Svelte",
      name: "Svelte",
    },
  ]);

  async function getRecruitment() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-recruitment", {
        URLQuery: {},
        URLParams: {
          id: recruitmentId,
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      console.log(response);

      if (response.status !== 200) {
        triggerToast("error");
        console.log(response);
        return;
      }

      if (response.result.recruitment.data.length > 0) {
        let recruitmentData = response.result.recruitment.data[0];
        data = {
          id: recruitmentData.id,
          fullname: recruitmentData.attributes.fullname,
          email: recruitmentData.attributes.email,
          phone: recruitmentData.attributes.phone,
          address: recruitmentData.attributes.address,
          education_history: recruitmentData.attributes.education_history,
          job_position: recruitmentData.attributes.job_position,
          experience: recruitmentData.attributes.experience,
          createdAt: format(
            new Date(recruitmentData.attributes.createdAt),
            "dd MMMM yyyy",
          ),
          ext: recruitmentData.attributes.resume.data.attributes.ext,
          resumeName: recruitmentData.attributes.resume.data.attributes.name,
          skillsOrKeywords: recruitmentData.attributes.skills,
        };

        base64 = response.result.base64;
        contentFile = atob(base64);

        const blob = base64toBlob(base64, data.ext);
        fileUrl = URL.createObjectURL(blob);
        console.log(fileUrl);

        pdfReady = true;
      }
    }
  }

  async function updateRecruitment() {
    let token = sessionStorage.getItem("token");
    if (
      recruitmentId !== "" &&
      recruitmentId !== null &&
      recruitmentId !== undefined
    ) {
      try {
        let response = await call.call(
          "client-dashboard",
          "update-recruitment",
          {
            URLQuery: {},
            URLParams: {
              id: recruitmentId.toString(),
            },
            Headers: {
              Authorization: `Bearer ${token}`,
            },
            Params: {
              fullname: data.fullname,
              email: data.email,
              phone: data.phone.toString(),
              address: data.address,
              education_history: data.education_history,
              skills: JSON.stringify(data.skillsOrKeywords),
              job_position: data.job_position,
              experience: data.experience,
            },
          },
        );
        if (response.status === 200) {
          triggerToast("updated");
          // setTimeout(async () => {
            goto(await call.getPage("client-dashboard", "hr-recruitment"));
          // }, 3000);
          return;
        } else {
          triggerToast("error");
          console.log(response);
          return;
        }
      } catch (error) {
        console.log("error");
        console.log(error);
      }
    }
  }

  async function deleteRecruitment() {
    let token = sessionStorage.getItem("token");
    if (
      recruitmentId !== "" &&
      recruitmentId !== null &&
      recruitmentId !== undefined
    ) {
      try {
        let response = await call.call(
          "client-dashboard",
          "delete-recruitment",
          {
            URLQuery: {},
            URLParams: {
              id: recruitmentId.toString(),
            },
            Headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );
        if (response.status === 200) {
          triggerToast("deleted");
          // setTimeout(async () => {
            goto(await call.getPage("client-dashboard", "hr-recruitment"));
          // }, 3000);
          return;
        } else {
          triggerToast("error");
          console.log(response);
          return;
        }
      } catch (error) {
        console.log("error");
        console.log(error);
      }
    }
  }

  function downloadFile(content, fileName, ext) {
    let fileType = "";
    if (ext === ".pdf") {
      fileType = "application/pdf";
    } else if (ext === ".docx") {
      fileType =
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
    } else if (ext === ".doc") {
      fileType = "application/msword";
    } else {
      fileType = "text/plain";
    }
    const a = document.createElement("a");
    a.href = URL.createObjectURL(base64toBlob(base64, fileType));
    a.download = fileName;
    a.click();
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

  onMount(() => {
    if (edit === "true") onEdit = true;
    console.log("id", recruitmentId);
    getRecruitment();
    getJobLevels();
  });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex items-center justify-between pb-8">
        <div class="flex flex-col grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
              >{data.fullname}<Button
                outline
                color="red"
                class="hover:!bg-transparent !border-hc-graylight-200 dark:!border-hc-graydark-600 focus-within:!ring-1 !p-2"
              >
                <CopyIcon
                  svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
                />
              </Button></Heading
            >
          </div>
          <P
            class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
            >{$_("createOn")}{` ${data.createdAt}`}</P
          >
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
            <Button
              onclick={updateRecruitment}
              class="px-3 focus-within:!ring-1">{$_("save")}</Button
            >
          {:else}
            <Button
              onclick={async function () {
                goto(await call.getPage("client-dashboard", "hr-recruitment"));
                return;
              }}
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

      <!-- Tab Bar -->
      <TabDetail
        {tabDatas}
        on:tabChange={(e) => (tabDatas.selected = e.detail)}
      />

      <!-- Content -->
      <div class="">
        {#if tabDatas.selected === "Resume"}
          <div class="flex justify-between items-center pt-6 pb-4">
            <P
              class="font-semibold text-hc-graylight-700 dark:text-hc-graydark-300"
              >{data.resumeName}</P
            >
            <div class="flex items-center gap-2">
              <Button
                onclick={() => downloadFile(base64, data.resume, data.ext)}
                class="w-max px-3 py-[10px] gap-1 focus-within:!ring-1"
              >
                <img alt="refresh" src={DownloadIcon} class="w-5 h-5" />
                {$_("downdloadDocument")}
              </Button>
            </div>
          </div>
          <div
            class="flex flex-col justify-center w-full py-10 px-6 bg-gray-50 rounded-lg dark:bg-gray-800"
          >
            {#if pdfReady}
              {#if data.ext === ".pdf"}
                <iframe title="" class="w-full h-screen" src={`${fileUrl}`}
                ></iframe>
              {:else}
                {#each contentFile.split("\n") as paragraphs, i}
                  <p class="dark:text-white">{paragraphs}</p>
                {/each}
              {/if}
            {/if}
          </div>
        {:else if tabDatas.selected === "Information"}
          <div class="flex flex-col gap-4 pt-6">
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
                      {$_("fullName")}
                    </P>
                    <Input
                      readonly={!onEdit}
                      bind:value={data.fullname}
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("jobPosition")}
                    </P>
                    <Select
                      disabled={!onEdit}
                      class="bg-primary-light dark:bg-secondary-dark"
                      items={jobLevels}
                      bind:value={data.job_position}
                    />
                  </div>
                </div>
                <div class="flex gap-5">
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("email")}
                    </P>
                    <Input
                      readonly={!onEdit}
                      bind:value={data.email}
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("phone")}
                    </P>
                    <Input
                      readonly={!onEdit}
                      bind:value={data.phone}
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                </div>
                <div class="grid grid-cols-1">
                  <div class="flex flex-col gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >{$_("address")}</P
                    >
                    <Textarea
                      readonly={!onEdit}
                      bind:value={data.address}
                      rows="4"
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                </div>
                <div class="flex gap-5">
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("educationHistory")}
                    </P>
                    <Input
                      readonly={!onEdit}
                      bind:value={data.education_history}
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                </div>
              </div>
            </div>
            <Divider />
            <div class="flex gap-6">
              <div class="w-60">
                <P
                  class="font-semibold text-secondary-light dark:text-secondary-dark"
                >
                  {$_("skillsOrKeywords")}
                </P>
              </div>
              <div class="flex flex-col grow gap-6">
                <!-- {#if onEdit} -->
                  <div class="flex grow gap-2">
                  <MultiSelect
                    disabled={!onEdit}
                    items={skillsOrKeywords}
                    bind:value={data.skillsOrKeywords}
                  />
                </div>
                <!-- {:else}
                  <div class="flex flex-wrap gap-2">
                  {#each data.skillsOrKeywords as skill}
                    <Badge text={skill} type="info" class={slugify(skill)} />
                  {/each}
                </div>
                {/if} -->
                <div class="flex gap-5">
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("experience")}
                    </P>
                    <Input
                      readonly={!onEdit}
                      bind:value={data.experience}
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
  <ToastContainer />
</div>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteRecruitment} color="red" class="me-2"
      >{$_("yesImSure")}</Button
    >
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>
