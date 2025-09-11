<script>
  import { Typography, InputRw, Dropdown } from "$lib";
  import { Button, Modal } from "flowbite-svelte";
  import Edit05Icon from "$lib/icon/edit-05.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import PlusIcon from "$lib/icon/plus.svelte";
  import { _, t } from "svelte-i18n";

  let { stateCreate = $bindable("formApi"), paramsData = $bindable([]) } = $props();

  let addParamModal = $state(false);
  let dropdownTypeOpen = $state(false);
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

  let type = $state([
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
  ]);

  async function addParam() {
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

    // console.log("cek paramsData", paramsData);
    addParamModal = false;
  }

  async function deleteParam(index) {
    // TODO: delete paramsData
    paramsData = paramsData.filter((item) => item.id !== index);
  }
</script>

<div class="flex flex-col m-4">
  <div class="flex items-start justify-between mb-4 md:mb-8">
    <Typography preset="Display sm/Semibold" color="primary" text={$_("addParams")} />
    <div class="flex gap-2 items-center">
      <Button color="light" on:click={() => (stateCreate = "formApi")}>{$_("cancel")}</Button>
      <Button class="gap-1.5" on:click={() => (stateCreate = "formApi")}>
        <Edit05Icon />
        {$_("saveParams")}
      </Button>
    </div>
  </div>

  <div class="flex flex-col gap-4 pt-4 md:pt-6 md:gap-6">
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
        <div class="hidden md:block">{$_("addParams")}</div>
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
        dropdownOpen={dropdownTypeOpen}
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
    <Button onclick={addParam} disabled={!newData.key || !newData.type} class="w-full"
      >{$_("addNewRow")}</Button
    >
  </svelte:fragment>
</Modal>
