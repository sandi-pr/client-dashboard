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
    MultiSelect,
    Select,
  } from "flowbite-svelte";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { _ } from "svelte-i18n";
  import { CR3Call } from "hexacomponents-svelte";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  const call = new CR3Call(globalThis.location.origin);
  let fileTextResult = $state(null);
  let base64Data = $state(null);
  let fileTextName = $state(null);

  let data = $state({
    fullname: "",
    email: "",
    phone: "",
    address: "",
    education_history: "",
    job_position: "",
    skillsOrKeywords: "",
    experience: "",
  });

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

  let jobLevels = $state([]);

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

    if (files[0]) {
      if (files.type === "text/plain") {
        const { text, base64 } = await readFileAsText(file);
        fileTextResult = text;
        base64Data = base64;
      } else {
        base64Data = await readFileAsBase64(file);
        console.log(base64Data); // For checking base64 data of other file types
        // Perform operations specific to other file types here
      }

      fileTextName = files[0].name;
    } else {
      alert("Invalid file. Please select a valid file.");
    }
  }

  const handleTextUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    // inputElement.accept = ".txt";
    inputElement.accept = ".txt, .pdf, .doc, .docx, .csv";

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file) {
        if (file.type === "text/plain") {
          const { text, base64 } = await readFileAsText(file);
          fileTextResult = text;
          base64Data = base64;
        } else {
          base64Data = await readFileAsBase64(file);
        }

        fileTextName = file.name;
      } else {
        alert("Invalid file. Please select a valid file.");
      }
    });

    inputElement.click();
  };

  async function readFileAsText(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        const text = reader.result.replace(/\r\n/g, "\n");
        const base64 = btoa(text);
        resolve({ text, base64 });
      };
      reader.onerror = reject;

      reader.readAsText(file);
    });
  }

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

  async function createRecruitment() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        console.log(typeof data.phone)
        console.log(typeof data.phone.toString())
        let response = await call.call(
          "client-dashboard",
          "create-recruitment",
          {
            URLQuery: {},
            URLParams: {},
            Params: {
              fullname: data.fullname,
              email: data.email,
              phone: data.phone.toString(),
              address: data.address,
              job_position: data.job_position,
              education_history: data.education_history,
              skills: JSON.stringify(data.skillsOrKeywords),
              experience: data.experience,
              resume_name: fileTextName,
              resume: base64Data,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          },
        );

        console.log(response);

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        if (response.status === 200) {
          goto(await call.getPage("client-dashboard", "hr-recruitment"));
        }
      } catch (error) {
        console.log(error);
      }
      return;
    } else {
      console.log("token not found");
      return;
    }
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
        console.log(error);
      }
    }
  }

  onMount(() => {
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
            >
              {$_("createNewRecruitment")}
            </Heading>
          </div>
        </div>

        <div class="flex gap-2">
          <Button
            onclick={async () => {
              goto(await call.getPage("client-dashboard", "hr-recruitment"));
            }}
            color="light"
            class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
          >
          <Button onclick={createRecruitment} class="px-3 focus-within:!ring-1"
            >{$_("create")}</Button
          >
        </div>
      </div>
      <div class="">
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
                    type="email"
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
                    type="number"
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
                {$_("resume")}
              </P>
            </div>
            <div class="flex flex-col grow gap-6">
              <Button
                ondragover={handleDragOver}
                ondragleave={handleDragLeave}
                ondrop={handleDrop}
                onclick={() => {
                  handleTextUploadWrapper();
                }}
                color="light"
                class={`w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1 py-8 flex flex-col justify-start ${fileTextName ? "items-start cursor-default" : "items-center"}`}
              >
                {#if fileTextName}
                  {#if fileTextName.endsWith(".txt")}
                    <!-- {#each fileTextResult.split("\n") as paragraphs, i}
                      <P>{paragraphs}</P>
                    {/each} -->
                    <P>{fileTextName}</P>
                  {:else if fileTextName.endsWith(".pdf")}
                    <P>{fileTextName}</P>
                  {:else if fileTextName.endsWith(".doc") || fileTextName.endsWith(".docx")}
                    <P>{fileTextName}</P>
                  {:else}
                    <P>{fileTextName}</P>
                  {/if}
                {:else}
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
                    <P class="text-hc-brand-700 pointer-events-none"
                      >{$_("clickToUpload")}</P
                    ><P class="pointer-events-none">{$_("orDragAndDrop")}</P>
                  </div>
                  <!-- <P class="pointer-events-none">SVG, PNG, JPG, or GIF (max. 800x400px)</P> -->
                  <P class="pointer-events-none"
                    >PDF, DOC, DOCX, CSV, TXT (Max. 15MB)</P
                  >
                {/if}
              </Button>
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
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("skillsOrKeywords")}</P
                >
                <MultiSelect
                  items={skillsOrKeywords}
                  bind:value={data.skillsOrKeywords}
                />
              </div>
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("experience")} (YoE)</P
                >
                <Input
                  type="number"
                  bind:value={data.experience}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
