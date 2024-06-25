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
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import DataTable from "$lib/components/DataTable.svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { StateCard } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { format } from "date-fns";
  import { Button, Heading, Card, A, P } from "flowbite-svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import CopyIcon from "$lib/icon/copy-03.svelte";

  import { goto } from "$app/navigation";

  const call = new CR3Call(globalThis.location.origin);

  let documentsReady = $state(false);

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });

  let datas = $state([]);

  let collections = $state([]);
  let documents = $state([]);

  let statusRAG = $state("");

  async function getCollections() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      let page = 1;
      while (true) {
        let response = await call.call(
          "client-dashboard",
          "get-document-collections",
          {
            URLQuery: {},
            URLParams: {},
            Params: {
              page: page,
              page_size: 25,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          },
        );

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        if (response.result.data.length > 0) {
          collections = collections.concat(response.result.data);
          page++;
        } else {
          break;
        }
      }

      let page2 = 1;
      while (true) {
        let response2 = await call.call("client-dashboard", "get-documents", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page2,
            page_size: 25,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          console.log(response2);
          return;
        }

        if (response2.result.data.length > 0) {
          documents = documents.concat(response2.result.data);
          page2++;
        } else {
          break;
        }
      }

      datas = collections.map((collection) => {
        let totalDocuments = documents.filter(
          (document) => document.attributes.collection_id === collection.id,
        ).length;
        return {
          collectionID: collection.id,
          activationTime: format(
            collection.attributes.createdAt,
            "dd-MM-yyyy HH:mm",
          ),
          collectionName: collection.attributes.collection_name,
          totalDocuments: totalDocuments,
          deployStatus: statusRAG,
        };
      });
    }
  }

  const classStatusType = (type) => {
    if (type === "RAG Ready") {
      return "bg-hc-success-50 border-hc-success-200 text-hc-success-700 dark:bg-hc-success-900 dark:border-hc-success-700 dark:text-hc-success-300";
    } else if (type === "Not Ready") {
      return "bg-hc-error-50 border-hc-error-200 text-hc-error-700 dark:bg-hc-error-900 dark:border-hc-error-700 dark:text-hc-error-300";
    } else if (type === "On Progress") {
      return "bg-hc-warning-50 border-hc-warning-200 text-hc-warning-700 dark:bg-hc-warning-900 dark:border-hc-warning-700 dark:text-hc-warning-300";
    }
  };

  // async function next() {
  //   getUsers(pagination.page + 1);
  // }

  // async function previous() {
  //   getUsers(pagination.page - 1);
  // }

  async function getStatusRag() {
    let token = sessionStorage.getItem("token");
    const callStatusRag = async () => {
      try {
        let rag = await call.call("hexa-ai-rag", "get-rag-status", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        let status = rag.result.rag_status;
        if (statusRAG.includes("!")) {
          let status2 = status.split("!")[0].trim();
          if (status2.includes("(NOVERB)") || status2.includes("(INOVERB)") || status2.includes("(!INOVERB)")) {
            let status3 = status2.replace("(NOVERB)", "").replace("(INOVERB)", "").replace("(!INOVERB)", "").replace("(!NOVERB)", "");
            statusRAG = status3;
          } else {
            statusRAG = status2;
          }
        } else {
          statusRAG = status;
        }
        console.log(rag); // Mengganti statusRAG dengan rag
      } catch (error) {
        console.error("Error:", error);
      }
      setTimeout(callStatusRag, 2000);
    };
    callStatusRag();
  }

  onMount(async () => {
    await getCollections();
    await getStatusRag();
  });

  async function addDocument() {
    goto(await call.getPage("client-dashboard", "document-collection-create"));
    return;
  }

  async function getDocumentCollectionDetail(id) {
    return (
      `${await call.getPage("client-dashboard", "document-collection-list")}` +
      "?collectionId=" +
      id
    );
  }

  async function askToAI(id) {
    return (
      `${await call.getPage("client-dashboard", "document-collection-chat")}` +
      "?collectionId=" +
      id
    );
  }
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography
            preset="Display sm/Semibold"
            color="primary"
            text={$_("documentCollectionList")}
          />
          <Button onclick={addDocument} icon={PlusIcon}
            ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" />
            {$_("addNew")}</Button
          >
        </div>

        <!-- Person Table Data
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs bind:next bind:previous />
        {:else}
          <div class="pt-8">
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          </div>
        {/if} -->

        <!-- Bots List Data -->
        <div>
          <!-- <Search
            class="bg-white dark:bg-hc-graydark-950 mb-6 rounded-lg"
            placeholder={$_("searchBotName")}
            bind:value={searchTerm}
          /> -->
          <div class="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4">
            {#if datas.length > 0}
              {#each datas as doc (doc.collectionID)}
                <!-- {#if doc.deployStatus === "RAG Ready"} -->
                <div class="span-col-1">
                  <Card padding="none" size="md">
                    <div
                      class="bg-active-light dark:bg-active-dark border-b border-secondary-light dark:border-secondary-dark px-4 py-3"
                    >
                      <div class="flex justify-between items-center">
                        <div class="flex flex-col">
                          <div class="flex items-center gap-2">
                            <Heading
                              tag="h6"
                              class="!font-semibold text-primary-light dark:text-primary-dark"
                              >{`${doc.collectionName}`}</Heading
                            >
                            <div
                              class={`px-2 w-full py-0.5 rounded-full border justify-start items-center inline-flex text-center text-xs font-medium leading-[18px] 
                              ${classStatusType(doc.deployStatus)}`}
                            >
                              {doc.deployStatus}
                            </div>
                          </div>
                          <P
                            size="sm"
                            class="text-secondary-light dark:text-secondary-dark"
                            >{`ID: ${doc.collectionID}`}</P
                          >
                        </div>
                        <div class="flex items-center pl-2 gap-2">
                          {#await getDocumentCollectionDetail(doc.collectionID) then href}
                            <A
                              {href}
                              class="text-fg-brand-primary-light dark:text-fg-brand-primary-dark text-sm font-semibold hover:no-underline"
                            >
                              {$_("view")}
                            </A>
                          {/await}
                          <div class=" flex items-center">
                            <Button
                              class="focus-within:ring-0 !p-0"
                              color="none"
                              size="xs"
                              ><DotsVerticalIcon
                                class="text-fg-quinary-light dark:text-fg-quinary-dark"
                              /></Button
                            >
                            <!-- <Dropdown>
                              <DropdownItem>
                                <div class="flex">
                                  <MessageChatSquare
                                    svgClass="text-fg-quaternary-light dark:text-fg-quaternary-dark mr-2"
                                  />{$_("conversationLogs")}
                                </div>
                              </DropdownItem>
                              <DropdownItem>
                                <div class="flex">
                                  <PencilLine
                                    svgClass="text-fg-quaternary-light dark:text-fg-quaternary-dark mr-2"
                                  />{$_("editBotName")}
                                </div>
                              </DropdownItem>
                              <DropdownItem>
                                <div class="flex">
                                  <FileAttachment_02
                                    svgClass="text-fg-quaternary-light dark:text-fg-quaternary-dark mr-2"
                                  />{$_("viewDocuments")}
                                </div>
                              </DropdownItem>
                            </Dropdown> -->
                          </div>
                        </div>
                      </div>
                    </div>
                    <div
                      class="pb-4 bg-primary-light dark:bg-primary-dark border-b border-secondary-light dark:border-secondary-dark p-4"
                    >
                      <div class="flex justify-between items-center">
                        <div>
                          <P
                            size="sm"
                            weight="semibold"
                            class="text-primary-light dark:text-primary-dark"
                            >{$_("totalDocuments")}</P
                          >
                          <P
                            size="sm"
                            weight="normal"
                            class="text-secondary-light dark:text-secondary-dark line-clamp-1"
                            >{doc.totalDocuments}</P
                          >
                        </div>
                        <!-- <Button
                          on:click={() => copyToClipboard(bot.totalDocuments)}
                          class="p-0 bg-transparent dark:bg-transparent hover:bg-transparent hover:dark:bg-transparent focus-within:ring-0"
                        >
                          <div>
                            <CopyIcon
                              svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                            />
                          </div>
                        </Button> -->
                      </div>
                    </div>
                    <div
                      class="pb-4 bg-primary-light dark:bg-primary-dark border-b border-secondary-light dark:border-secondary-dark p-4"
                    >
                      <div class="flex justify-between items-center">
                        <div>
                          <P
                            size="sm"
                            weight="semibold"
                            class="text-primary-light dark:text-primary-dark"
                            >{$_("activation")}</P
                          >
                          <P
                            size="sm"
                            weight="normal"
                            class="text-secondary-light dark:text-secondary-dark line-clamp-1"
                            >{doc.activationTime}</P
                          >
                        </div>
                        <!-- <Button
                          on:click={() => copyToClipboard(bot.activationTime)}
                          class="p-0 bg-transparent dark:bg-transparent hover:bg-transparent hover:dark:bg-transparent focus-within:ring-0"
                        >
                          <div>
                            <CopyIcon
                              svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                            />
                          </div>
                        </Button> -->
                      </div>
                    </div>
                    <div
                      class="pb-4 bg-primary-light dark:bg-primary-dark border-b border-secondary-light dark:border-secondary-dark p-4"
                    >
                      {#await askToAI(doc.collectionID) then href}
                        <Button {href} color="light" class="w-full">
                          {$_("askToAI")}
                        </Button>
                      {/await}
                    </div>
                  </Card>
                </div>
                <!-- {/if} -->
              {/each}
            {/if}
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
