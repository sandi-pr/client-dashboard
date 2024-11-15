<script>
  import { onMount } from "svelte";
  import { Heading } from "flowbite-svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { Button } from "flowbite-svelte";
  import PlusIcon from "$lib/icon/plus.svg";
  import { _ } from "svelte-i18n";
  import { Typography, StateCard, check, call } from "$lib";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { goto } from "$app/navigation";

  let pagination = $state({ page: 1, pageSize: 10, pageCount: 1, total: 0 });

  let datas = $state([]);

  let gettingDatas = $state(false);

  let configs = $state({});

  async function getRecruitment(page) {
    if (page > pagination.pageCount) return;
    let token = localStorage.getItem("token");
    if (token !== null) {
      gettingDatas = true;
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
        gettingDatas = false;
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
                check("businessManagement.hr.recruitment:view") &&
                {
                  href:
                    `${await $call.getPage("client-dashboard", "hr-recruitment-detail")}` +
                    `?recruitmentId=${recruitment.id}&edit=false`,
                  name: "View",
                },
                check("businessManagement.hr.recruitment:edit") &&
                {
                  href:
                    `${await $call.getPage("client-dashboard", "hr-recruitment-detail")}` +
                    "?recruitmentId=" +
                    recruitment.id +
                    "&edit=true",
                  name: "Edit",
                },
              ].filter(Boolean),
            },
          };
        });

        datas = await Promise.all(promises2);
      }

      gettingDatas = false;
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
    if (check("businessManagement.hr.recruitment:view") === false) {
      console.log("no access");
      window.history.back();
      return;
    }
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
          {#if check("businessManagement.hr.recruitment:create")}
            <Button onclick={addRecruitment} icon={PlusIcon}
              ><img class="w-4 h-4 me-1" src={PlusIcon} alt="plus-icon" /> {$_("addNew")}</Button
            >
          {/if}
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
