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
  import { onMount, onDestroy } from "svelte";
  import { _ } from "svelte-i18n";
  import DataTableListCollection from "$lib/components/DataTableListCollection.svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { StateCard, TabBar } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { format } from "date-fns";
  import { Button, P } from "flowbite-svelte";
  import PlusIcon from "$lib/icon/plus.svg";

  import { goto } from "$app/navigation";

  const params = new URLSearchParams(window.location.search);
  let collectionId = params.get("collectionId");

  const call = new CR3Call(globalThis.location.origin);

  let documentsReady = $state(false);

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });

  let datas = $state([]);

  let datasChecked = $state([]);

  let collection = $state({});

  let configs = $state({});

  let tabDatas = {
    selected: "Documents",
    data: [
      {
        name: "Documents",
      },
      {
        name: "Information",
      },
    ],
  };

  let timer;

  async function getDocuments(page) {
    let token = sessionStorage.getItem("token");
    token = true;
    if (token !== null) {
      // let response = await call.call(
      //   "client-dashboard",
      //   "get-documents-by-collectionid",
      //   {
      //     URLQuery: {},
      //     URLParams: {},
      //     Params: {
      //       page: page,
      //       page_size: pagination.pageSize,
      //       collection_id: collectionId,
      //     },
      //     Headers: {
      //       Authorization: "Bearer " + token,
      //     },
      //   },
      // );

      // if (response.status !== 200) {
      //   console.log(response);
      //   return;
      // }

      let response = {
        result: {
          data: [
            {
              id: 1,
              attributes: {
                createdAt: new Date(),
                document_title: "Document 1",
                document_type: "GENERAL",
                status: "PENDING",
              },
            },
            {
              id: 2,
              attributes: {
                createdAt: new Date(),
                document_title: "Document 2",
                document_type: "CONTRACT",
                status: "APPROVED",
              },
            },
          ],
          meta: {
            pagination: {
              page: 1,
              pageSize: 10,
              pageCount: 1,
              total: 2,
            },
          },
        },
      };

      if (response.result.data.length > 0) {
        configs = {
          search: true,
          headers: "File Uploaded",
          pagination: response.result.meta.pagination,
          actions: [
            {
              id: "add_new",
              name: "Add New",
              icon: "plus",
              type: "brand-inverse",
              href: "#",
            },
          ],
        };
        let promises = response.result.data.map(async (document) => {
          return {
            id: document.id,
            file_name: { name: document.attributes.document_title, type: "custom", size: "200Kb" },
            type: {
              name: document.attributes.document_type === "GENERAL" ? "General" : "Contract",
              type: "pill",
              color: "brand",
            },
            file_size: "200Kb",
            upload_date: format(document.attributes.createdAt, "dd-MM-yyyy HH:mm"),
            uploaded_by: "Admin",
            actions: {
              objects: [
                {
                  href: "localhost:5173/client-dashboard/document-collection-detail?documentId=1&collectionId=1",
                  // `${await call.getPage("client-dashboard", "document-collection-detail")}` +
                  // "?documentId=" +
                  // document.id +
                  // "&collectionId=" +
                  // collectionId,
                  name: "View",
                },
                {
                  // href: "localhost:5173/client-dashboard/document-collection-detail?documentId=1&collectionId=1",
                  // `${await call.getPage("client-dashboard", "document-collection-detail")}` +
                  // "?documentId=" +
                  // document.id +
                  // "&collectionId=" +
                  // collectionId,
                  function: handleTitikTiga,
                  name: "Titik tiga",
                },
              ],
            },
          };
        });

        datas = await Promise.all(promises);
        datasChecked = datas.map((data) => {
          return { id: data.id, checked: false };
        });
      }

      // timer = setTimeout(getDocuments(pagination.page), 2000);
    }

    documentsReady = true;
  }

  async function getCollection() {
    // search collection using params collectionId here
    collection = {
      id: 1,
      name: "Collection 1",
      createdAt: format(new Date(), "dd-MM-yyyy HH:mm"),
      description: "Collection 1 Description",
    };
  }

  const changeTabbar = (name) => {
    tabDatas.selected = name;
  };

  async function next() {
    getDocuments(pagination.page + 1);
  }

  async function previous() {
    getDocuments(pagination.page - 1);
  }

  async function handleCheckAll() {
    datasChecked = datasChecked.map((data) => {
      return { id: data.id, checked: !data.checked };
    });
  }

  async function handleCheck(id) {
    let datasCheckedCopy = JSON.parse(JSON.stringify(datasChecked));
    datasCheckedCopy = datasCheckedCopy.map((data) => {
      if (data.id === id) {
        return { id: data.id, checked: !data.checked };
      }
      return data;
    });
    datasChecked = datasCheckedCopy;
  }

  async function handleTitikTiga(id) {
    console.log("Titik tiga ", id);
  }

  onMount(async () => {
    getDocuments(pagination.page);
    collectionId ? getCollection() : null;
  });

  onDestroy(() => {
    if (timer) clearTimeout(timer);
  });

  async function addDocument() {
    goto(
      (await call.getPage("client-dashboard", "document-management-create")) +
        "?collectionId=" +
        collectionId
    );
    return;
  }
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <div>
            <Typography
              preset="Display sm/Semibold"
              color="primary"
              text={collection.name ? collection.name : $_("loading")}
            />
            <P
              class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
            >
              {$_("createOn")}{` ${collection.createdAt}`}
            </P>
          </div>
          <Button onclick={addDocument} icon={PlusIcon}
            ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" />
            {$_("startReading")}</Button
          >
        </div>

        <!-- Submenu -->
        <TabBar {tabDatas} {changeTabbar} />

        <!-- Person Table Data -->
        {#if datas.length > 0}
          <DataTableListCollection
            bind:datas
            bind:datasChecked
            bind:configs
            bind:next
            bind:previous
            bind:handleCheckAll
            bind:handleCheck
          />
        {:else}
          <div class="pt-8">
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
