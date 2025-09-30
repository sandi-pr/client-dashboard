<script>
  import { onMount } from "svelte";
  import { Button, Helper, Modal } from "flowbite-svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { page } from "$app/stores";
  import { _ } from "svelte-i18n";
  import { goto } from "$app/navigation";
  import { call, getIcons, triggerToast } from "$lib";
  import { Typography, InputRw, TextareaRw, Divider, Dropdown, DataTable, TabBar } from "$lib";

  const Save01Icon = getIcons("save01");
  const PlusIcon = getIcons("plus");
  const Copy01Icon = getIcons("copy01");
  const Edit05Icon = getIcons("edit05");
  const DotsVerticalIcon = getIcons("dotsVertical");

  const isFrontend = process.env.NODE_ENV === "development";

  let {
    data = $bindable({
      apiName: "",
      description: "",
      createOn: "",
      urlLink: "",
      method: "GET",
      authorization: "Basic Auth",
      username: "",
      password: "",
      token: "",
    }),
    stateCreate = $bindable("formApi"),
    paramsData = $bindable(null),
    bodyData = $bindable(null),
    isBody = $bindable(false),
    onEdit = $bindable(false),
    saveApi = () => {},
    deleteApi = () => {},
  } = $props();

  let apiId = $state(null);
  let libraryId = $state(null);
  let collectionId = $state(null);

  let deleteModal = $state(false);
  let dropdownMore = $state(false);
  let dropdownMethodOpen = $state(false);
  let dropdownAuthorization = $state(false);

  let bodyConfigured = $state(0);
  let paramsConfigured = $state(0);

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
        name: "Body",
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
    { name: "Inherit Auth from Parent", value: "Inherit Auth from Parent" },
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
      apiId = $page.url.searchParams.get("apiId");
      if (apiId === null) {
        onEdit = true;
      }
      if (collectionId === null || collectionId === undefined) {
        goto(await $call.getPage("client-dashboard", "api-connector"));
        return;
      }
    });

    return () => {
      unsubscribe();
    };
  });

  $effect(() => {
    paramsConfigured = paramsData.length;
    bodyConfigured = bodyData.length;
  });
</script>

<div class="flex flex-col m-4 pb-8 md:pb-24">
  <div class="flex justify-between">
    {#if apiId !== null && onEdit === false}
      <div class="flex flex-col gap-2">
        <div class="flex gap-2 items-center">
          <Typography preset="Display sm/Semibold" color="primary" text={"GetUserAPI"} />
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
            <Copy01Icon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
          </button>
        </div>
        <Typography
          preset="Text xs/Medium"
          color="primary"
          text={`${$_("createOn")}: ${data.createOn}`}
        />
      </div>
      <div class="flex gap-2 items-center">
        <Button
          color="light"
          class="gap-1.5 text-fg-button-secondary-color-light dark:text-fg-button-secondary-color-dark"
          on:click={() => {
            onEdit = !onEdit;
            tabDatas.selected = "API Information";
            // data = dummyData;
          }}
        >
          <Edit05Icon
            svgClass="w-5 h-5 text-fg-button-secondary-color-light dark:text-fg-button-secondary-color-dark"
          />
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
                  deleteModal = true;
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
        <Typography preset="Display sm/Semibold" color="primary" text={$_("createNewApi")} />
      {:else}
        <Typography preset="Display sm/Semibold" color="primary" text={$_("editApi")} />
      {/if}
      <div class="flex gap-2 items-center">
        <Button
          color="light"
          on:click={async () => {
            const link = isFrontend
              ? "/api_connector/collections/api" +
                `?collectionId=${collectionId}&libraryId=${libraryId}`
              : (await $call.getPage("client-dashboard", "api-connector-collections-api")) +
                `?collectionId=${collectionId}&libraryId=${libraryId}`;
            apiId === null ? goto(link) : (onEdit = !onEdit);
          }}
        >
          {$_("cancel")}
        </Button>
        <Button
          class="gap-1.5"
          on:click={() => {
            onEdit = false;
            saveApi(data);
          }}
        >
          <Save01Icon />
          {$_("saveApi")}
        </Button>
      </div>
    {/if}
  </div>
  {#if apiId !== null && onEdit === false}
    <TabBar
      {tabDatas}
      changeTabbar={(tabName) => {
        tabDatas.selected = tabName;
      }}
    />
  {/if}
  {#if tabDatas.selected === "API Information"}
    <div class="flex flex-col gap-4 pt-4 md:pt-6 md:gap-6">
      <div class="flex flex-col md:flex-row gap-3 md:gap-6">
        <div class="w-60">
          <Typography preset="Text md/Semibold" color="secondary" text={$_("generalInformation")} />
        </div>
        <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
          <div class="flex flex-col gap-1.5">
            <Typography preset="Text sm/medium" color="secondary" text={$_("apiName")} />
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
          <Typography preset="Text md/Semibold" color="secondary" text={$_("method")} />
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
          <Typography preset="Text md/Semibold" color="secondary" text={$_("urlLink")} />
        </div>
        <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
          <div class="flex flex-col gap-1.5">
            <Typography preset="Text sm/medium" color="secondary" text={$_("apiUrlLink")} />
            <InputRw {onEdit} bind:value={data.urlLink} />
          </div>
        </div>
      </div>

      {#if apiId === null || onEdit === true}
        <Divider />
        <div class="flex flex-col md:flex-row gap-3 md:gap-6">
          <div class="w-60">
            <Typography preset="Text md/Semibold" color="secondary" text={$_("params")} />
          </div>
          <div class="flex flex-col gap-3 md:gap-6 rounded-lg">
            <div class="flex flex-col gap-1.5">
              <Button
                color="light"
                class="gap-1.5"
                on:click={async () => {
                  stateCreate = "formParams";
                  isBody = false;
                }}
              >
                <PlusIcon
                  svgClass="w-5 h-5 text-fg-button-secondary-color-light dark:text-fg-button-secondary-color-dark"
                />
                <div
                  class="hidden text-fg-button-secondary-color-light dark:text-fg-button-secondary-color-dark md:block"
                >
                  {$_("addParams")}
                </div>
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
            <Typography preset="Text md/Semibold" color="secondary" text={$_("body")} />
          </div>
          <div class="flex flex-col gap-3 md:gap-6 rounded-lg">
            <div class="flex flex-col gap-1.5">
              <Button
                color="light"
                class="gap-1.5"
                on:click={async () => {
                  stateCreate = "formBody";
                  isBody = true;
                }}
              >
                <PlusIcon
                  svgClass="w-5 h-5 text-fg-button-secondary-color-light dark:text-fg-button-secondary-color-dark"
                />
                <div
                  class="hidden text-fg-button-secondary-color-light dark:text-fg-button-secondary-color-dark md:block"
                >
                  {$_("addBody")}
                </div>
              </Button>
              <Helper>
                <span class="font-medium">{`${bodyConfigured} ${$_("bodyConfigured")}`}</span>
              </Helper>
            </div>
          </div>
        </div>
        <Divider />
        <div class="flex flex-col md:flex-row gap-3 md:gap-6">
          <div class="w-60">
            <Typography preset="Text md/Semibold" color="secondary" text={$_("authType")} />
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
              <Typography preset="Text md/Semibold" color="secondary" text={$_("token")} />
            </div>
            {#if data.authorization === "Basic Auth"}
              <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
                <div class="flex flex-col gap-1.5">
                  <Typography preset="Text sm/medium" color="secondary" text={$_("username")} />
                  <InputRw {onEdit} bind:value={data.username} />
                </div>
                <div class="flex flex-col gap-1.5">
                  <Typography preset="Text sm/medium" color="secondary" text={$_("password")} />
                  <InputRw {onEdit} bind:value={data.password} />
                </div>
              </div>
            {:else if data.authorization === "Bearer Token"}
              <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
                <div class="flex flex-col gap-1.5">
                  <Typography preset="Text sm/medium" color="secondary" text={$_("token")} />
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
  {:else if tabDatas.selected === "Body"}
    <div class="pt-2 md:pt-4">
      <DataTable bind:datas={bodyData} />
    </div>
  {:else if tabDatas.selected === "Authorization"}
    <div class="flex flex-col md:flex-row pt-4 md:pt-6 gap-3 md:gap-6">
      <div class="w-60">
        <Typography preset="Text md/Semibold" color="secondary" text={$_("authType")} />
      </div>
      <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/medium" color="secondary" text={$_("authorizationType")} />
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
          <Typography preset="Text md/Semibold" color="secondary" text={$_("token")} />
        </div>
        {#if data.authorization === "Basic Auth"}
          <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
            <div class="flex flex-col gap-1.5">
              <Typography preset="Text sm/medium" color="secondary" text={$_("username")} />
              <InputRw {onEdit} bind:value={data.username} />
            </div>
            <div class="flex flex-col gap-1.5">
              <Typography preset="Text sm/medium" color="secondary" text={$_("password")} />
              <InputRw {onEdit} bind:value={data.password} />
            </div>
          </div>
        {:else if data.authorization === "Bearer Token"}
          <div class="flex flex-col grow gap-3 md:gap-6 rounded-lg">
            <div class="flex flex-col gap-1.5">
              <Typography preset="Text sm/medium" color="secondary" text={$_("token")} />
              <InputRw {onEdit} bind:value={data.token} />
            </div>
          </div>
        {/if}
      </div>
    {/if}
  {/if}
</div>

<Modal
  bind:open={deleteModal}
  size="xs"
  autoclose
  backdropClass="fixed inset-0 z-50 bg-gray-900 bg-opacity-50 dark:bg-opacity-80"
  dialogClass="fixed top-0 start-0 end-0 h-modal md:inset-0 md:h-full z-50 w-full p-4 flex"
>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      Are you sure you want to delete {data.name}?
    </h3>
    <Button
      onclick={() => {
        deleteApi();
      }}
      color="red"
      class="me-2">Yes, I'm sure</Button
    >
    <Button color="alternative" onclick={() => (deleteModal = false)}>Cancel</Button>
  </div>
</Modal>
