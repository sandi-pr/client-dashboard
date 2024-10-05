<script>
  import { onMount } from "svelte";
  import { Heading } from "flowbite-svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { Button } from "flowbite-svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import { _ } from "svelte-i18n";
  import { StateCard } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { goto } from "$app/navigation";

  import { call } from "$lib";

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });

  let datas = $state([]);

  let gettingDatas = $state(false);

  let configs = $state({});

  function getStatus(is_active) {
    if (is_active === "Yes") {
      return {
        type: "success",
        name: "Active",
      };
    }

    if (is_active === "No") {
      return {
        type: "danger",
        name: "Inactive",
      };
    }

    return;
  }

  async function getPerson(page) {
    if (page > pagination.pageCount) return;
    let token = localStorage.getItem("token");
    if (token !== null) {
      gettingDatas = true;
      let response = await $call.call("client-dashboard", "get-persons", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: page,
          page_size: pagination.pageSize,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        gettingDatas = false;
        return;
      }

      if (response.result.data.length > 0) {
        pagination = response.result.meta.pagination;
        configs = {
          search: true,
          pagination: response.result.meta.pagination,
        };
        let promises2 = response.result.data.map(async (person, index) => {
          return {
            id: person.id,
            fullname: person.attributes.fullname,
            email: person.attributes.email,
            phone: person.attributes.phone,
            // address: person.attributes.address,
            // city: person.attributes.city,
            // state: person.attributes.state,
            // zip_code: person.attributes.zip_code,
            country: person.attributes.country,
            // date_of_birth: person.attributes.date_of_birth,
            status: getStatus(person.attributes.is_active),
            actions: {
              objects: [
                {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-person-detail")}` +
                    `?personId=${person.id}&edit=false`,
                  name: "View",
                },
                {
                  href:
                    `${await $call.getPage("client-dashboard", "crm-person-detail")}` +
                    "?personId=" +
                    person.id +
                    "&edit=true",
                  name: "Edit",
                },
              ],
            },
          };
        });

        datas = await Promise.all(promises2);
      }

      gettingDatas = false;
    }
  }

  async function next() {
    getPerson(pagination.page + 1);
  }

  async function previous() {
    getPerson(pagination.page - 1);
  }

  async function addPerson() {
    goto(await $call.getPage("client-dashboard", "crm-person-create"));
    return;
  }

  onMount(async () => {
    getPerson(pagination.page);
  });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text={$_("personsList")} />
          <Button onclick={addPerson} icon={PlusIcon}
            ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" /> {$_("addNew")}</Button
          >
        </div>

        <!-- Document Table Data -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs bind:next bind:previous />
        {:else if gettingDatas}
          <StateCard
            illustration={IllustrationProgress}
            description={$_("fetchingDataPleaseWait")}
          />
        {:else}
          <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
        {/if}
      </div>
    </div>
  </div>
</div>
