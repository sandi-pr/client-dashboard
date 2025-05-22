<script>
  import { Typography, InputRw, TextareaRw, Divider, Dropdown, call, triggerToast } from "$lib";
  import { Button, Helper } from "flowbite-svelte";
  import { page } from "$app/stores";
  import { _ } from "svelte-i18n";
  import { goto } from "$app/navigation";
  import Save01Icon from "$lib/icon/save-01.svelte";
  import PlusIcon from "$lib/icon/plus.svelte";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  import Edit05Icon from "$lib/icon/edit-05.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import TabBar from "$lib/components/TabBar.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { onMount } from "svelte";
  import { on } from "svelte/events";

  const isFrontend = process.env.NODE_ENV === "development";

  let collectionId = $state(null);
  let libraryId = $state(null);
  let apiId = $state(null);
  let dropdownMore = $state(false);

  let data = $state({
    apiName: "",
    description: "",
    createOn: "",
    apiLink: "",
    method: "GET",
    authorization: "Basic Auth",
    username: "",
    password: "",
    token: "",
  });

  let dummyData = $state({
    apiName: "GetUserAPI",
    description: "API uses to get a User Name",
    createOn: "2023-07-02 08:51 PM",
    apiLink: "https://hexacode.ai/api/character?name",
    method: "GET",
    authorization: "Inherit Auth from Parent",
    username: "",
    password: "",
    token: "",
  })
  let paramsConfigured = $state(0);
  let onEdit = $state(true);
  let dropdownMethodOpen = $state(false);
  let dropdownAuthorization = $state(false);

  let { stateCreate = $bindable("formApi"), paramsData = $bindable(null), saveApi = () => {} } = $props();

  let tabDatas = $state({
    selected: "API Information",
    data: [
      {
        name: "API Information",
      },
      {
        name: "Params",
      },
      {
        name: "Authorization",
      },
    ],
  });

  let method = $state([
    { name: "GET", value: "GET" },
    { name: "POST", value: "POST" },
    { name: "PUT", value: "PUT" },
    { name: "PATCH", value: "PATCH" },
    { name: "DELETE", value: "DELETE" },
    { name: "HEAD", value: "HEAD" },
    { name: "OPTIONS", value: "OPTIONS" },
  ]);

  let authorization = $state([
    { name: "Inherit Auth from Parent", value: "Inherit Auth from Parent"},
    { name: "No Auth", value: "No Auth" },
    { name: "Basic Auth", value: "Basic Auth" },
    { name: "Bearer Token", value: "Bearer Token" },
  ]);

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownMore = false;
  };

  onMount(() => {
    const unsubscribe = page.subscribe(async ($page) => {
      collectionId = $page.url.searchParams.get("collectionId");
      libraryId = $page.url.searchParams.get("libraryId");
      apiId = $page.url.searchParams.get("apiId");
      if (apiId !== null) {
        onEdit = false;
        data = dummyData;
      }
      if (
        collectionId === null ||
        collectionId === undefined ||
        libraryId === null ||
        libraryId === undefined
      ) {
        goto(await $call.getPage("client-dashboard", "api-connector"));
        return;
      }
    });

    return () => {
      unsubscribe();
    };
  })

</script>

<div class="flex flex-col m-4 pb-8 md:pb-24">
  <div class="flex justify-between"> 
    {#if apiId !== null && onEdit === false}
      <div class="flex flex-col gap-2">
        <div class="flex gap-2 items-center">
          <Typography
            preset="Display sm/Semibold"
            color="primary"
            text={"GetUserAPI"}
          />
          <button
            onclick={() => {
              navigator.clipboard
                .writeText(data.apiName)
                .then(() => {
                  // console.log("Text copied to clipboard");
                  triggerToast("copyToClipboard");
                })
                .catch((err) => {
                  console.error("Failed to copy: ", err);
                });
            }}
            class="p-2 border-2 rounded-lg border-secondary-light dark:border-secondary-dark"
          >
            <Copy01Icon
              svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
            />
          </button>
        </div>
        <Typography
          preset="Text xs/Medium"
          color="primary"
          text={`${$_("createOn")}: ${data.createOn}`}
        />
      </div>
      <div class="flex gap-2 items-center">
        <Button color="light" class="gap-1.5 text-fg-button-secondary-color-light dark:text-fg-button-secondary-color-dark" 
          on:click={() => {
            onEdit = !onEdit;
            data = dummyData;
          }}>
          <Edit05Icon svgClass="w-5 h-5 text-fg-button-secondary-color-light dark:text-fg-button-secondary-color-dark" />
          {$_("edit")}
        </Button>
        <div onfocusout={handleDropdownFocusLoss} class=" flex items-center relative">
          <Button on:click={() => (dropdownMore = !dropdownMore)} color="light" class="p-2.5">
            <DotsVerticalIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
          </Button>
          <div
            class={`z-10 ${dropdownMore ? "absolute top-10 right-0" : "hidden"} bg-primary-light divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-primary-dark`}
          >
            <ul class="py-2 grow text-sm">
              <button
                  class="flex p-2 w-full text-error-primary-light dark:text-error-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                  onclick={() => {
                    dropdownMore = !dropdownMore;
                  }}
                >
                  {$_("deleteApi")}
                </button>
            </ul>
          </div>
        </div>
      </div>
    {:else}
      {#if apiId === null}
        <Typography
          preset="Display sm/Semibold"
          color="primary"
          text={$_("createNewApi")}
        />
      {:else}
        <Typography
          preset="Display sm/Semibold"
          color="primary"
          text={$_("editApi")}
        />
      {/if}
      <div class="flex gap-2 items-center">
        <Button color="light" on:click={async() => {
          const link =
          isFrontend
            ? "/api_connector/collections/api" + `?collectionId=${collectionId}&libraryId=${libraryId}`
            : await $call.getPage("client-dashboard", "api-connector-collections-api") + `?collectionId=${collectionId}&libraryId=${libraryId}`;
          onEdit ? onEdit = !onEdit : goto(link);
        }}>{$_("cancel")}</Button>
        <Button class="gap-1.5" on:click={saveApi}>
          <Save01Icon />
          {$_("saveApi")}
        </Button>
      </div>
    {/if}
  </div>
  {#if apiId !== null && onEdit === false}
   <TabBar {tabDatas} />
  {/if}
  {#if tabDatas.selected === "API Information"}
    <div class="flex flex-col gap-4 pt-4 md:pt-6 md:gap-6">
      <div class="flex flex-col md:flex-row gap-3 md:gap-6">
        <div class="w-60">
          <Typography
              preset="Text md/Semibold"
              color="secondary"
              text={$_("generalInformation")}
            />
        </div>
        <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
          <div class="flex flex-col gap-1.5">
            <Typography
              preset="Text sm/medium"
              color="secondary"
              text={$_("apiName")}
            />
            <InputRw {onEdit} bind:value={data.apiName} />
          </div>
          <div class="flex flex-col gap-1.5">
            <Typography
              preset="Text sm/medium"
              color="secondary"
              text={$_("descriptionOptional")}
            />
            <TextareaRw {onEdit} bind:value={data.description} />
          </div>
        </div>
      </div>
      <Divider />
      <div class="flex flex-col md:flex-row gap-3 md:gap-6">
        <div class="w-60">
          <Typography
              preset="Text md/Semibold"
              color="secondary"
              text={$_("method")}
            />
        </div>
        <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
          <div class="flex flex-col gap-1.5">
            <Dropdown
              dropdownOpen={dropdownMethodOpen}
              bind:selectedItem={data.method}
              items={method}
              width="w-min"
              isOverlay={true}
              disabled={!onEdit}
            />
          </div>
        </div>
      </div>
      <Divider />
      <div class="flex flex-col md:flex-row gap-3 md:gap-6">
        <div class="w-60">
          <Typography
              preset="Text md/Semibold"
              color="secondary"
              text={$_("urlLink")}
            />
        </div>
        <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
          <div class="flex flex-col gap-1.5">
            <Typography
              preset="Text sm/medium"
              color="secondary"
              text={$_("apiUrlLink")}
            />
            <InputRw {onEdit} bind:value={data.urlLink} />
          </div>
        </div>
      </div>
      
      {#if apiId === null || onEdit === true}
        <Divider />
        <div class="flex flex-col md:flex-row gap-3 md:gap-6">
          <div class="w-60">
            <Typography
                preset="Text md/Semibold"
                color="secondary"
                text={$_("params")}
              />
          </div>
          <div class="flex flex-col gap-3 md:gap-6 rounded-lg">
            <div class="flex flex-col gap-1.5">
              <Button
                color="light"
                class="gap-1.5"
                on:click={async () => {
                  stateCreate = "formParams";
                }}
              >
                <PlusIcon svgClass="w-5 h-5 text-fg-button-secondary-color-light dark:text-fg-button-secondary-color-dark" />
                <div class="hidden text-fg-button-secondary-color-light dark:text-fg-button-secondary-color-dark md:block">{$_("addParams")}</div>
              </Button>
              <Helper>
                <span class="font-medium">{`${paramsConfigured} ${$_("paramsConfigured")}`}</span>
              </Helper>
            </div>
          </div>
        </div>
        <Divider />
        <div class="flex flex-col md:flex-row gap-3 md:gap-6">
          <div class="w-60">
            <Typography
                preset="Text md/Semibold"
                color="secondary"
                text={$_("authType")}
              />
          </div>
          <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
            <div class="flex flex-col gap-1.5">
              <Typography
                preset="Text sm/medium"
                color="secondary"
                text={$_("authorizationType")}
              />
              <Dropdown
                dropdownOpen={dropdownAuthorization}
                bind:selectedItem={data.authorization}
                items={authorization}
                isOverlay={true}
                disabled={!onEdit}
              />
              <Helper>
                <span class="font-medium">{$_("descAuthorizationType")}</span>
              </Helper>
            </div>
          </div>
        </div>
        {#if data.authorization === "Basic Auth" || data.authorization === "Bearer Token"}
          <Divider />
          <div class="flex flex-col md:flex-row gap-3 md:gap-6">
            <div class="w-60">
              <Typography
                  preset="Text md/Semibold"
                  color="secondary"
                  text={$_("token")}
                />
            </div>
            {#if data.authorization === "Basic Auth"}
              <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
                <div class="flex flex-col gap-1.5">
                  <Typography
                    preset="Text sm/medium"
                    color="secondary"
                    text={$_("username")}
                  />
                  <InputRw {onEdit} bind:value={data.username} />
                </div>
                <div class="flex flex-col gap-1.5">
                  <Typography
                    preset="Text sm/medium"
                    color="secondary"
                    text={$_("password")}
                  />
                  <InputRw {onEdit} bind:value={data.password} />
                </div>
              </div>
            {:else if data.authorization === "Bearer Token"}
              <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
                <div class="flex flex-col gap-1.5">
                  <Typography
                    preset="Text sm/medium"
                    color="secondary"
                    text={$_("token")}
                  />
                  <InputRw {onEdit} bind:value={data.token} />
                </div>
              </div>
            {/if}
          </div>
        {/if}
      {/if}
    </div>
  {:else if tabDatas.selected === "Params"}
    <div class="pt-2 md:pt-4">
      <DataTable bind:datas={paramsData} />
    </div>
  {:else if tabDatas.selected === "Authorization"}
    <div class="flex flex-col md:flex-row pt-4 md:pt-6 gap-3 md:gap-6">
      <div class="w-60">
        <Typography
            preset="Text md/Semibold"
            color="secondary"
            text={$_("authType")}
          />
      </div>
      <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
        <div class="flex flex-col gap-1.5">
          <Typography
            preset="Text sm/medium"
            color="secondary"
            text={$_("authorizationType")}
          />
          <Dropdown
            dropdownOpen={dropdownAuthorization}
            bind:selectedItem={data.authorization}
            items={authorization}
            isOverlay={true}
            disabled={!onEdit}
          />
          <Helper>
            <span class="font-medium">{$_("descAuthorizationType")}</span>
          </Helper>
        </div>
      </div>
    </div>
    {#if data.authorization === "Basic Auth" || data.authorization === "Bearer Token"}
      <Divider />
      <div class="flex flex-col md:flex-row gap-3 md:gap-6">
        <div class="w-60">
          <Typography
              preset="Text md/Semibold"
              color="secondary"
              text={$_("token")}
            />
        </div>
        {#if data.authorization === "Basic Auth"}
          <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
            <div class="flex flex-col gap-1.5">
              <Typography
                preset="Text sm/medium"
                color="secondary"
                text={$_("username")}
              />
              <InputRw {onEdit} bind:value={data.username} />
            </div>
            <div class="flex flex-col gap-1.5">
              <Typography
                preset="Text sm/medium"
                color="secondary"
                text={$_("password")}
              />
              <InputRw {onEdit} bind:value={data.password} />
            </div>
          </div>
        {:else if data.authorization === "Bearer Token"}
          <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
            <div class="flex flex-col gap-1.5">
              <Typography
                preset="Text sm/medium"
                color="secondary"
                text={$_("token")}
              />
              <InputRw {onEdit} bind:value={data.token} />
            </div>
          </div>
        {/if}
      </div>
    {/if}
  {/if}
</div>