<script>
  import { _ } from "svelte-i18n";
  import { Button, Modal } from "flowbite-svelte";
  import { getIcons } from "$lib";
  import { Typography, InputRw, Dropdown, DataTable } from "$lib";

  const Edit05Icon = getIcons("edit05");
  const PlusIcon = getIcons("plus");

  let {
    stateCreate = $bindable("formApi"),
    data = $bindable({}),
    paramsData = $bindable([]),
    isBody = $bindable(false),
  } = $props();

  const type = [
    { name: "INT", value: "INT" },
    { name: "FLOAT", value: "FLOAT" },
    { name: "DOUBLE", value: "DOUBLE" },
    { name: "CHAR", value: "CHAR" },
    { name: "VARCHAR", value: "VARCHAR" },
    { name: "TEXT", value: "TEXT" },
    { name: "BINARY", value: "BINARY" },
    { name: "BLOB", value: "BLOB" },
    { name: "DATE", value: "DATE" },
    { name: "DATETIME", value: "DATETIME" },
    { name: "TIMESTAMP", value: "TIMESTAMP" },
    { name: "TIME", value: "TIME" },
    { name: "YEAR", value: "YEAR" },
    { name: "BOOLEAN", value: "BOOLEAN" },
    { name: "JSON", value: "JSON" },
    { name: "XML", value: "XML" },
  ];
  const bodyType = [
    { name: "Raw", value: "raw" },
    { name: "Form-data", value: "form-data" },
    // { name: "x-www-form-urlencoded", value: "x-www-form-urlencoded" },
    // { name: "binary", value: "binary" },
    // { name: "GraphQL", value: "GraphQL" },
  ];

  let addParamModal = $state(false);
  let dropdownColumnTypeOpen = $state(false);
  let dropdownBodyTypeOpen = $state(false);
  let newData = $state({
    id: 1,
    key: "",
    type: "",
    description: "",
    actions: {
      objects: [
        {
          onclick: () => {},
          name: "Delete",
          icon: "trash",
        },
      ],
    },
  });

  const addParam = () => {
    newData.id = paramsData.length + 1;
    newData.actions.objects[0].onclick = () => deleteParam(newData.id);
    paramsData = [...paramsData, newData];

    newData = {
      id: 1,
      key: "",
      type: "",
      description: "",
      actions: {
        objects: [
          {
            onclick: () => {},
            name: "Delete",
            icon: "trash",
          },
        ],
      },
    };

    addParamModal = false;
  };

  const deleteParam = (index) => {
    paramsData = paramsData.filter((item) => item.id !== index);
  };
</script>

<div class="flex flex-col m-4">
  <div class="flex items-start justify-between mb-4 md:mb-8">
    <Typography
      preset="Display sm/Semibold"
      color="primary"
      text={isBody ? $_("addBody") : $_("addParams")}
    />
    <div class="flex gap-2 items-center">
      <Button color="light" on:click={() => (stateCreate = "formApi")}>{$_("cancel")}</Button>
      <Button class="gap-1.5" on:click={() => (stateCreate = "formApi")}>
        <Edit05Icon />
        {isBody ? $_("saveBody") : $_("saveParams")}
      </Button>
    </div>
  </div>

  <div class="flex flex-col gap-4 pt-4 md:pt-6 md:gap-6">
    {#if isBody}
      <div class="flex">
        <div class="w-1/4">
          <Typography preset="Text sm/medium" color="secondary" text={$_("bodyType")} />
        </div>
        <div class="flex flex-col gap-2 w-3/4">
          <Dropdown
            dropdownOpen={dropdownBodyTypeOpen}
            bind:selectedItem={data.body_type}
            items={bodyType}
            isOverlay={true}
            placeholder={$_("selectType")}
          />
          <Typography preset="Text sm/regular" text={$_("bodyTypeDesc")} />
        </div>
      </div>
    {/if}
    <DataTable bind:datas={paramsData} />
    <div>
      <Button
        outline
        class="gap-1.5"
        on:click={() => {
          addParamModal = true;
        }}
      >
        <PlusIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
        <div class="hidden md:block">{isBody ? $_("addBody") : $_("addParams")}</div>
      </Button>
    </div>
  </div>
</div>

<Modal title={$_("addNewRow")} bind:open={addParamModal} size="xs">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("key")}*`} />
      <InputRw bind:value={newData.key} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("types")}*`} />
      <Dropdown
        dropdownOpen={dropdownColumnTypeOpen}
        bind:selectedItem={newData.type}
        items={type}
        isOverlay={true}
        placeholder={$_("selectType")}
      />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <InputRw bind:value={newData.description} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={addParam} disabled={!newData.key || !newData.type} class="w-full">
      {$_("addNewRow")}
    </Button>
  </svelte:fragment>
</Modal>
