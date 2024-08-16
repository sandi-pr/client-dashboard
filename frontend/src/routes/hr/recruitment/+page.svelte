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
  import { goto } from "$app/navigation";

  import { call } from "$lib";

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });

  let datas = $state([]);

  let configs = $state({});

  async function getRecruitment(page) {
    if (page > pagination.pageCount) return;
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await $call.call("client-dashboard", "get-recruitments", {
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
        return;
      }

      if (response.result.data.length > 0) {
        pagination = response.result.meta.pagination;
        configs = {
          search: true,
          pagination: response.result.meta.pagination,
        };
        let promises2 = response.result.data.map(async (recruitment, index) => {
          return {
            id: recruitment.id,
            fullname: recruitment.attributes.fullname,
            job_position: recruitment.attributes.job_position,
            experience: recruitment.attributes.experience + "+ YoE",
            email: recruitment.attributes.email,
            actions: {
              objects: [
                {
                  href:
                    `${await $call.getPage("client-dashboard", "hr-recruitment-detail")}` +
                    `?recruitmentId=${recruitment.id}&edit=false`,
                  name: "View",
                },
                {
                  href:
                    `${await $call.getPage("client-dashboard", "hr-recruitment-detail")}` +
                    "?recruitmentId=" +
                    recruitment.id +
                    "&edit=true",
                  name: "Edit",
                },
              ],
            },
          };
        });

        datas = await Promise.all(promises2);
      }
    }
  }

  async function next() {
    getRecruitment(pagination.page + 1);
  }

  async function previous() {
    getRecruitment(pagination.page - 1);
  }

  async function addRecruitment() {
    goto(await $call.getPage("client-dashboard", "hr-recruitment-create"));
    return;
  }

  onMount(async () => {
    getRecruitment(pagination.page);
  });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4 pt-2">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between">
          <Typography preset="Display sm/Semibold" color="primary" text={$_("recruitmentList")} />
          <Button onclick={addRecruitment} icon={PlusIcon}
            ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" /> {$_("addNew")}</Button
          >
        </div>

        <!-- Document Table Data -->
        {#if datas.length > 0}
          <DataTable bind:datas bind:configs bind:next bind:previous />
        {:else}
          <div class="pt-8">
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
