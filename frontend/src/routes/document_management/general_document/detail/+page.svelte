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
  import { P, Heading, Input, Textarea, Button, Select, Modal } from "flowbite-svelte";
  import { ToastContainer, triggerToast } from "$lib";
  import { CheckCircleSolid, ExclamationCircleSolid } from "flowbite-svelte-icons";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import DownloadIcon from "$lib/icon/file-download-03.svg";
  import Divider from "$lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { format } from "date-fns";
  import { goto } from "$app/navigation"; 
  import { createEventDispatcher } from 'svelte';

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  let documentId = params.get("documentId");
  const edit = params.get("edit");
  const dispatch = createEventDispatcher();

  let onEdit = $state(false);

  let popupModal = $state(false);

  let contentFile = $state(null);
  let base64 = $state(null);
  let fileUrl = $state("null");

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

  let personClient = $state([]);

  async function getDocument() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      try {
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
          triggerToast("error");
          console.log(response);
          return;
        }

        if (response.result.document.data.length > 0) {
          let documentData = response.result.document.data[0];
          data = {
            id: documentData.id,
            createdAt: format(documentData.attributes.createdAt, "yyyy-MM-dd HH:mm"),
            document_title: documentData.attributes.document_title,
            description: documentData.attributes.description,
            document_type: documentData.attributes.document_type,
            ext: documentData.attributes.document.data.attributes.ext,
            person_client: documentData.attributes.person.data?.id,
            contract_start:
            documentData.attributes.contract_start !== null
                ? format(new Date(documentData.attributes.contract_start), "yyyy-MM-dd")
                : "",
            contract_end:
            documentData.attributes.contract_end !== null
                ? format(new Date(documentData.attributes.contract_end), "yyyy-MM-dd")
                : "",
          };

          base64 = response.result.base64;
          contentFile = atob(base64);

          let fileType = "";
          if (data.ext === ".pdf") {
            fileType = "application/pdf";
          } else if (data.ext === ".docx") {
            fileType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
          } else if (data.ext === ".doc") {
            fileType = "application/msword";
          } else {
            fileType = "text/plain";
          }

          const blob = base64toBlob(base64, fileType);
          fileUrl = URL.createObjectURL(blob);
          console.log(fileUrl);
          
          pdfReady = true;
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    } else {
      triggerToast("error");
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
          triggerToast("deleted");
          setTimeout(async () => {
            goto(await call.getPage("client-dashboard", "document-management-general"));
          }, 3000);
          return;
        } else {
          triggerToast("error");
          console.log(response);
          return;
        }
      } catch (error) {
        triggerToast("error");
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
    const a = document.createElement('a');
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

  function getDocumentId() {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get('documentid');
  }

  $effect(async () => {
    if (documentId) {
      `${await call.getPage("client-dashboard", "document-management-general-detail")}` + "?documentId=" + document.id + "&edit=false"
    }
  })

  onMount(() => {
    const handleUrlChange = () => {
      const newDocumentId = getDocumentId();
      if (newDocumentId !== documentId) {
        documentId = newDocumentId;
      }
    };
    // if (edit === "true") onEdit = true;

    window.addEventListener('popstate', handleUrlChange);
    window.addEventListener('pushstate', handleUrlChange); // In case you use pushState for navigation

    getPersons();
    getDocument();

    return () => {
      window.removeEventListener('popstate', handleUrlChange);
      window.removeEventListener('pushstate', handleUrlChange);
    };
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
          <Button
            onclick={() => (popupModal = true)}
            outline color="red"
            class="px-3 focus-within:!ring-1">{$_("delete")}</Button
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
                onclick={() => downloadFile(base64, data.document_title, data.ext)}
                class="w-max px-3 py-[10px] gap-1 focus-within:!ring-1"
              >
                <img alt="refresh" src={DownloadIcon} class="w-5 h-5" />
                {$_("downdloadDocument")}
              </Button>
            </div>
          </div>
          <div class="flex flex-col justify-center w-full py-10 px-6 bg-gray-50 rounded-lg dark:bg-gray-800">
            {#if pdfReady}
              {#if data.ext === ".pdf"}
                <iframe class="w-full h-screen" src={`${fileUrl}`} ></iframe>
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
  <ToastContainer />
</div>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteDocument} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>