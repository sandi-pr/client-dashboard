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
  import { P, Heading, Input, Textarea, Button, Select } from "flowbite-svelte";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import DownloadIcon from "$lib/icon/file-download-03.svg";
  import Divider from "../../../lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { format } from "date-fns";
  import { goto } from "$app/navigation";

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  const documentId = params.get("documentId");
  const edit = params.get("edit");

  let onEdit = $state(false);

  let contentFile = $state("");

  let tabDatas = $state({
    selected: "Documents",
    data: [
      {
        name: "Documents",
      },
      {
        name: "Information",
      },
    ],
  });

  let pdfReady = $state(true);

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

  let personClient = $state([]);

  async function getDocument() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-document", {
        URLQuery: {},
        URLParams: {
          id: documentId.toString(),
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      if (response.result.document.data.length > 0) {
        let document = response.result.document.data[0];
        data = {
          id: document.id,
          createdAt: format(document.attributes.createdAt, "yyyy-MM-dd HH:mm"),
          document_title: document.attributes.document_title,
          description: document.attributes.description,
          document_type: document.attributes.document_type,
          ext: document.attributes.document.data.attributes.ext,
          person_client: document.attributes.person.data?.id,
          contract_start:
            document.attributes.contract_start !== null
              ? format(new Date(document.attributes.contract_start), "yyyy-MM-dd")
              : "",
          contract_end:
            document.attributes.contract_end !== null
              ? format(new Date(document.attributes.contract_end), "yyyy-MM-dd")
              : "",
        };

        let base64 = response.result.base64;
        contentFile = atob(base64);
      }
    }
  }

  async function deleteDocument() {
    let token = sessionStorage.getItem("token");
    if (documentId !== "" && documentId !== null && documentId !== undefined) {
      try {
        let response = await call.call("client-dashboard", "delete-document", {
          URLQuery: {},
          URLParams: {
            id: documentId.toString(),
          },
          Headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        if (response.status === 200) {
          goto(await call.getPage("client-dashboard", "document-management-list"));
          return;
        } else {
          console.log(response);
          return;
        }
      } catch (error) {
        console.log(error);
      }
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

  function downloadFile(content, fileName, ext) {
    let fileType = "";
    if (ext === ".pdf") {
      fileType = "application/pdf";
    } else if (ext === ".docx") {
      fileType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
    } else if (ext === ".doc") {
      fileType = "application/msword";
    } else {
      fileType = "text/plain";
    }
    const blob = new Blob([content], { type: fileType });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = fileName;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
  }

  onMount(() => {
    if (edit === "true") onEdit = true;
    getPersons();
    getDocument();
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
              >{data.document_title}<Button
                outline
                color="red"
                class="hover:!bg-transparent !border-hc-graylight-200 dark:!border-hc-graydark-600 focus-within:!ring-1 !p-2"
              >
                <CopyIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
              </Button></Heading
            >
          </div>
          <P
            class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
            >{$_("createOn")}{` ${data.createdAt}`}</P
          >
        </div>

        <div class="flex gap-2">
          <Button color="light" class="px-3 focus-within:!ring-1" onclick={deleteDocument}
            >{$_("delete")}</Button
          >
        </div>
      </div>

      <!-- Tab Bar -->
      <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

      <!-- Content -->
      <div class="">
        {#if tabDatas.selected === "Documents"}
          <div class="flex justify-between items-center pt-6 pb-4">
            <P class="font-semibold text-hc-graylight-700 dark:text-hc-graydark-300"
              >{$_("documentName")}</P
            >
            <div class="flex items-center gap-2">
              <Button
                onclick={() => downloadFile(contentFile, data.document_title, data.ext)}
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
              {#each contentFile.split("\n") as paragraphs, i}
                <P class="dark:text-white">{paragraphs}</P>
              {/each}
            {/if}
          </div>
        {:else if tabDatas.selected === "Information"}
          <div class="flex flex-col gap-4 pt-6">
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
                      {$_("documentName")}
                    </P>
                    <Input
                      readonly={!onEdit}
                      bind:value={data.document_title}
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                </div>
                <div class="grid grid-cols-1">
                  <div class="flex flex-col gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >{$_("description")}</P
                    >
                    <Textarea
                      readonly={!onEdit}
                      bind:value={data.description}
                      rows="4"
                      class="bg-primary-light dark:bg-secondary-dark"
                    />
                  </div>
                </div>
              </div>
            </div>
            <Divider />
            <div class="flex gap-6">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {$_("documentType")}*
                </P>
              </div>
              <div class="flex flex-col grow gap-6">
                <div class="flex flex-col gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("documentType")}*</P
                  >
                  <Select
                    disabled={!onEdit}
                    class="bg-primary-light dark:bg-secondary-dark"
                    items={categories}
                    bind:value={data.document_type}
                  />
                </div>
              </div>
            </div>
            {#if data.document_type === "CONTRACT"}
              <Divider />
              <!-- general information -->
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
                        disabled={!onEdit}
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
                        bind:value={data.contract_end}
                        class="bg-primary-light dark:bg-secondary-dark"
                      />
                    </div>
                  </div>
                </div>
              </div>
            {/if}
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
