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
  import { Button } from "flowbite-svelte";
  import PlusIcon from "$lib/icon/plus.svg";

  import { goto } from "$app/navigation";

  const params = new URLSearchParams(window.location.search);
  let collectionId = params.get("collectionId");

  const call = new CR3Call(globalThis.location.origin);

  let documentsReady = $state(false);

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });

  let datas = $state([]);

  let configs = $state({});

  async function getDocuments(page) {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call(
        "client-dashboard",
        "get-documents-by-collectionid",
        {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: pagination.pageSize,
            collection_id: collectionId,
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
        configs = {
          search: true,
          pagination: response.result.meta.pagination,
        };
        let promises = response.result.data.map(async (document) => {
          return {
            id: document.id,
            upload_date: format(
              document.attributes.createdAt,
              "dd-MM-yyyy HH:mm",
            ),
            document_name: document.attributes.document_title,
            type:
              document.attributes.document_type === "GENERAL"
                ? "General"
                : "Contract",
            // status: {
            //   type: "brand",
            //   name: document.attributes.status,
            // },
            actions: {
              objects: [
                {
                  href:
                    `${await call.getPage("client-dashboard", "document-collection-detail")}` +
                    "?documentId=" +
                    document.id +
                    "&collectionId=" + collectionId,
                  name: "View",
                },
              ],
            },
          };
        });

        datas = await Promise.all(promises);
      }
    }

    documentsReady = true;
  }

  async function next() {
    getUsers(pagination.page + 1);
  }

  async function previous() {
    getUsers(pagination.page - 1);
  }

  onMount(async () => {
    getDocuments(pagination.page);
  });

  async function addDocument() {
    goto(
      (await call.getPage("client-dashboard", "document-management-create")) +
        "?collectionId=" +
        collectionId,
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
          <Typography
            preset="Display sm/Semibold"
            color="primary"
            text={$_("documentsList")}
          />
          <Button onclick={addDocument} icon={PlusIcon}
            ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" />
            {$_("addNew")}</Button
          >
        </div>

        <!-- Person Table Data -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs bind:next bind:previous />
        {:else}
          <div class="pt-8">
            <StateCard
              illustration={IllustrationEmpty}
              description={$_("youDontHaveAnyDataYet")}
            />
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
