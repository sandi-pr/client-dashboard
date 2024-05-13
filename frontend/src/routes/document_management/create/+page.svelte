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
  import { P, Heading, Input, Button, Select, Textarea, Toggle } from "flowbite-svelte";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import TabBar from "$lib/ui/tab_bar/TabBar.svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

  let dropzoneActive = false;
  let onEdit = $state(true);
  let uploadedImage = $state(null);
  let uploadedImageName = $state(null);
  let fileImage = $state(null);
  let fileText = $state(null);
  let fileTextName = $state(null);
  let fileTextResult = $state(null);
  let base64Data = $state(null);
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
    // for (let i = 0; i < files.length; i++) {
    //   console.log(files[i].name);
    if (files[0] && isValidTextFile(files[0])) {
      fileTextResult = await new Promise((resolve, reject) => {
        const reader = new FileReader();

        reader.onload = () => resolve(reader.result.replace(/\r\n/g, "\n"));
        reader.onerror = reject;

        reader.readAsText(files[0]);
      });

      fileText = files[0];
      fileTextName = files[0].name;
      // console.log(fileTextResult);
    } else {
      alert("Invalid file. Please select a valid .txt file.");
    }
    // }
  }

  function isValidTextFile(file) {
    const fileType = file.type;
    return fileType === "text/plain";
  }

  const handleTextUploadWrapper = (index) => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = "text/plain";

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file && isValidTextFile(file)) {
        fileTextResult = await new Promise((resolve, reject) => {
          const reader = new FileReader();

          reader.onload = () => {
            const text = reader.result.replace(/\r\n/g, "\n");
            resolve(text);

            // Encode text as Base64
            base64Data = btoa(text);
          };
          reader.onerror = reject;

          reader.readAsText(file);
        });

        fileText = file;
        fileTextName = file.name;
        // console.log(fileTextResult);
      } else {
        alert("Invalid file. Please select a valid .txt file.");
      }
    });

    inputElement.click();
  };

  async function createProduct() {
    const fileTextResultBlob = new Blob([fileTextResult], { type: "text/plain" });

    console.log(data);
    console.log(fileText);
    console.log(fileTextName);
    console.log(fileTextResult);
    console.log(fileTextResultBlob);

    let token = sessionStorage.getItem("token");

    if (token !== null) {
      if (data.document_type === "GENERAL") {
        try {
          let response = await call.call("client-dashboard", "create-document", {
            URLQuery: {},
            URLParams: {},
            Params: {
              document_title: data.document_title,
              document: base64Data,
              document_name: fileTextName,
              description: data.description,
              document_type: data.document_type,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          console.log(response);

          if (response.status !== 200) {
            console.log(response);
            return;
          }

          if (response.status === 200) {
            await call.goToPage("client-dashboard", "document-management-overview");
          }
        } catch (error) {
          console.log(error);
        }
      }
      return;
    }
  }

  async function getPersons() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let page = 1;
      while (true) {
        let response = await call.call("client-dashboard", "get-persons", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 25,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        if (response.result.data.length > 0) {
          let promise = response.result.data.map((person) => {
            return {
              value: person.id,
              name:
                person.attributes.fullname !== null
                  ? person.attributes.fullname
                  : `<${person.attributes.email}>`,
            };
          });

          let data = await Promise.all(promise);
          personClient = personClient.concat(data);
          page++;
        } else {
          break;
        }
      }
    }
  }

  let buttonDisable = $state(true);

  $effect(() => {
    if (
      data.document_title !== "" &&
      data.description !== "" &&
      data.description !== "" &&
      data.is_active !== ""
    ) {
      buttonDisable = false;
    } else {
      buttonDisable = true;
    }
  });

  onMount(async () => {
    getPersons();
  });
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
              Add new document
            </Heading>
          </div>
        </div>
        <div class="flex gap-2">
          <Button
            onclick={async function () {
              await call.goToPage("client-dashboard", "crm-product-service");
              return;
            }}
            color="light"
            class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
          >
          <Button onclick={createProduct} class="px-3 focus-within:!ring-1">{$_("create")}</Button>
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
                  {$_("Document Title")}*
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.document_title}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="flex flex-col col-span-5 gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >{$_("Document Upload")}*</P
              >
              <Button
                ondragover={handleDragOver}
                ondragleave={handleDragLeave}
                ondrop={handleDrop}
                onclick={() => {
                  if (!fileTextName) {
                    handleTextUploadWrapper(index);
                  }
                }}
                color="light"
                class={`w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1 py-8 flex flex-col justify-start ${fileTextName ? "items-start cursor-default" : "items-center"}`}
              >
                {#if fileTextName}
                  {#each fileTextResult.split("\n") as paragraphs, i}
                    <P>{paragraphs}</P>
                  {/each}
                {:else}
                  <div
                    class="rounded-xl shadow border border-primary-light dark:border-primary-dark justify-center items-center inline-flex mb-4 pointer-events-none"
                  >
                    <img src={upload} alt="mail.svg" class="p-4 dark:invert pointer-events-none" />
                  </div>
                  <div class="flex space-x-1">
                    <P class="text-hc-brand-700 pointer-events-none">Click to upload</P><P
                      class="pointer-events-none">or drag and drop</P
                    >
                  </div>
                  <!-- <P class="pointer-events-none">SVG, PNG, JPG, or GIF (max. 800x400px)</P> -->
                  <P class="pointer-events-none">Txt file</P>
                {/if}
              </Button>
            </div>
            <div class="grid grid-cols-1">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("Description")}</P
                >
                <Textarea
                  readonly={!onEdit}
                  bind:value={data.description}
                  rows="3"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
          </div>
        </div>
        <Divider />
        <!-- product detail -->
        <div class="flex gap-6">
          <div class="w-60">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              {$_("Document Type")}
            </P>
          </div>
          <div class="flex flex-col grow gap-6">
            <div class="flex flex-col grow gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >{$_("Document Type")}*</P
              >
              <Select
                items={categories}
                bind:value={data.document_type}
                placeholder="Select document type"
              />
              <!-- <Input
                readonly={!onEdit}
                bind:value={data.document_type}
                class="bg-primary-light dark:bg-secondary-dark"
              />   -->
            </div>
          </div>
        </div>
        {#if data.document_type === "contract"}
          <Divider />
          <!-- product detail -->
          <div class="flex gap-6">
            <div class="w-60">
              <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                {$_("contractInformation")}
              </P>
            </div>
            <div class="flex grow flex-col gap-6">
              <div class="flex flex-col grow gap-6">
                <div class="flex flex-col gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("personsOrClient")}*</P
                  >
                  <Select
                    class="bg-primary-light dark:bg-secondary-dark"
                    items={personClient}
                    bind:value={data.person_client}
                  />
                </div>
              </div>
              <div class="flex justify-center items-center gap-5">
                <div class="flex flex-col grow gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >
                    {$_("contractStartPeriod")}
                  </P>
                  <Input
                    readonly={!onEdit}
                    type="date"
                    bind:value={data.contract_start}
                    class="bg-primary-light dark:bg-secondary-dark"
                  />
                </div>
                <div class="mt-6">
                  <P>-</P>
                </div>
                <div class="flex flex-col grow gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >
                    {$_("contractEndPeriod")}
                  </P>
                  <Input
                    readonly={!onEdit}
                    type="date"
                    bind:value={data.contract_end}
                    class="bg-primary-light dark:bg-secondary-dark"
                  />
                </div>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
