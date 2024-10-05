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
  import { P, Heading, Input, Button, Select, Textarea, Checkbox } from "flowbite-svelte";
  import { ToastContainer, triggerToast } from "$lib";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { format } from "date-fns";
  import { goto } from "$app/navigation";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";

  import { call } from "$lib";

  let dropzoneActive = false;
  let onEdit = $state(true);
  let files = $state([]);
  let fileTextNames = $state([]);
  let fileTextName = $state(null);
  let fileTextResult = $state(null);
  let base64Data = $state("");

  let personClient = $state([]);

  let data = $state({
    document_title: "",
    document_upload: null,
    description: "",
    document_type: "",
    person_client: "",
    contract_start: "",
    contract_end: "",
  });

  let resources = $state([]);

  let index = $state(resources.length);

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
      if (file.type === "text/plain") {
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

  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  const handleFileUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = ".pdf, .doc, .docx";
    inputElement.multiple = true;

    inputElement.addEventListener("change", async (event) => {
      const uploadedFiles = event.target.files;

      if (uploadedFiles.length > 0) {
        for (const file of uploadedFiles) {
          const base64Data = await readFileAsBase64(file);
          files.push({
            name: file.name,
            size: formatBytes(file.size),
            base64: base64Data,
          });
          fileTextNames.push(file.name);
        }
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

  async function createCollection() {
    console.log(files);
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response;
      if (files.length > 0) {
        response = await $call.call("client-dashboard", "create-document-collection", {
          URLQuery: {},
          URLParams: {},
          Params: {
            collection_name: data.document_title,
            // description: data.description,
            documents: JSON.stringify(files),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
      } else {
        response = await $call.call("client-dashboard", "create-document-collection", {
          URLQuery: {},
          URLParams: {},
          Params: {
            collection_name: data.document_title,
            // description: data.description,
            documents: "-",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
      }
      if (response.status !== 200) {
        console.log(response);
        return;
      }

      console.log(response);

      let response2 = await $call.call("hexa-ai-rag", "add-document-by-collection-id", {
        URLQuery: {},
        URLParams: {},
        Params: {
          collection_id: response.result.data.id.toString(),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response2.status !== 200) {
        console.log(response2);
        return;
      }

      console.log(response2);

      goto(await $call.getPage("client-dashboard", "document-collection"));
    }
  }
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- Header -->
      <div class="flex items-center pb-8">
        <div class="flex flex-col gap-2 grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {$_("addNewDocumentCollection")}
            </Heading>
          </div>
        </div>
        <div class="flex gap-2">
          <Button
            onclick={async function () {
              goto(await $call.getPage("client-dashboard", "document-collection"));
              return;
            }}
            color="light"
            class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
          >
          <Button onclick={createCollection} class="px-3 focus-within:!ring-1"
            >{$_("create")}</Button
          >
        </div>
      </div>

      <Divider />

      <div class="flex flex-col gap-4 pt-8">
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
                  {$_("documentCollectionName")}*
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.document_title}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <!-- <div class="grid grid-cols-1">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("description")}</P
                >
                <Textarea
                  readonly={!onEdit}
                  bind:value={data.description}
                  rows="3"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div> -->
            <!-- <div class="flex flex-col col-span-5 gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >{$_("documentUpload")}*</P
              >
              <Button
                onclick={() => {
                  handleFileUploadWrapper();
                }}
                color="light"
                class={`w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1 py-8 flex flex-col justify-start ${fileTextName ? "items-start cursor-default" : "items-center"}`}
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
                  <P class="text-hc-brand-700 pointer-events-none"
                    >{$_("clickToUpload")}</P
                  >
                  <P class="pointer-events-none">
                  </P>
                </div>
                <P class="pointer-events-none">PDF, DOC, DOCX (Max. 15MB)</P>
              </Button>
            </div> -->
            <!-- Display uploaded file names -->
            {#if fileTextNames.length > 0}
              <div class="flex flex-col gap-3">
                {#each files as file}
                  <div
                    class="flex justify-between p-4 rounded-xl border border-secondary-light dark:border-secondary-dark"
                  >
                    <div class="flex gap-3">
                      {#if file.name.includes(".pdf")}
                        <img src={PDFSolidIcon} alt="" />
                      {:else if file.name.includes(".docx")}
                        <img src={DOCXSolidIcon} alt="" />
                      {:else if file.name.includes(".doc")}
                        <img src={DOCSolidIcon} alt="" />
                      {/if}
                      <div class="flex flex-col justify-center">
                        <P class="text-secondary-light dark:text-secondary-dark">{file.name}</P>
                        <P class="text-secondary-light dark:text-secondary-dark">{file.size}</P>
                      </div>
                    </div>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>
  </div>
  <ToastContainer />
</div>
