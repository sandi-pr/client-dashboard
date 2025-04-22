<script>
  import { P, Checkbox } from "flowbite-svelte";
  import { getIcons, toCamelCase } from "$lib";
  import { _ } from "svelte-i18n";

  import { ChevronDownOutline, ChevronUpOutline } from "flowbite-svelte-icons";

  import BookOpen2Icon from "$lib/icon/book-open-02.svelte";
  import Stars01Icon from "$lib/icon/stars-01.svelte";
  import MessageChatCircleIcon from "$lib/icon/message-chat-circle.svelte";

  export let permissionListDatas;
  export let edit;
  export let change;

  let mainColSpan = 12 - permissionListDatas[0].data.length * 2;

  const handleHideContent = (id) => {
    // console.log("id", id);
    permissionListDatas.map((permission) => {
      if (permission.parentId === id || permission.id === id) {
        permission.hide = !permission.hide;
      }
      return permission;
    });
  };

  const borderInsideGrid = (index, permission) => {
    if (permission.type === "content" && permission.hide) return " ";

    let borderClass = "border-secondary-light dark:border-secondary-dark ";
    if (index % 2 !== 0) return borderClass + "border-y-2 ";
    return " ";
  };

  const backgroundFromType = (type) => {
    if (type === "header") return "bg-primary-light dark:bg-primary-dark";
    else if (type === "section") return "bg-secondary-light dark:bg-secondary-dark";
    else if (type === "content") return "bg-primary-light dark:bg-primary-dark";
  };

  const leftPaddingFromLevel = (level) => {
    if (level === 0) return "";
    else if (level === 1) return "pl-4";
    else if (level === 2) return "pl-8";
  };
</script>

<div class="w-full rounded-lg border-2 border-secondary-light dark:border-secondary-dark">
  {#each permissionListDatas as permission, index (permission.id)}
    <div
      class={"flex px-8 w-full " +
        borderInsideGrid(index, permission) +
        backgroundFromType(permission.type)}
    >
      {#if permission.type === "header"}
        <div class={`py-4 self-center grow`}>
          <P class="text-lg text-tertiary-light dark:text-tertiary-dark"
            >{$_(toCamelCase(permission.title))}</P
          >
        </div>
        {#each permission.data as data (data.name)}
          <div class="py-4 flex justify-center w-20">
            <P class="font-semibold text-tertiary-light dark:text-tertiary-dark"
              >{$_(toCamelCase(data.name))}</P
            >
          </div>
        {/each}
      {:else if permission.type === "section"}
        <div class={`py-4 self-center grow ${leftPaddingFromLevel(permission.level)}`}>
          <div class="flex gap-3 items-center">
            {#if permission.icon}
              <svelte:component
                this={getIcons(permission.icon)}
                svgClass="w-5 h-5 text-hc-graylight-500 dark:text-hc-graydark-400"
              />
            {/if}
            <div class="w-full flex items-center justify-between">
              <P class="text-lg font-bold text-tertiary-light dark:text-tertiary-dark"
                >{$_(toCamelCase(permission.title))}</P
              >
              {#if permission.hide}
                <div on:click={() => handleHideContent(permission.id)}>
                  <ChevronDownOutline
                    class="w-8 h-8 text-tertiary-light dark:text-tertiary-dark cursor-pointer"
                  />
                </div>
              {:else}
                <div on:click={() => handleHideContent(permission.id)}>
                  <ChevronUpOutline
                    class="w-8 h-8 text-tertiary-light dark:text-tertiary-dark cursor-pointer"
                  />
                </div>
              {/if}
            </div>
          </div>
        </div>
      {:else if permission.type === "content" && !permission.hide}
        <div class="w-full flex items-center justify-between">
          <div class={`py-6 self-center grow`}>
            <P class="text-tertiary-light dark:text-tertiary-dark"
              >{$_(toCamelCase(permission.title))}</P
            >
          </div>
          {#each permission.data as data (data.name)}
            <div class="py-6 flex justify-center w-20">
              <Checkbox
                disabled={!edit}
                checked={data.granted}
                on:click={() => {
                  data.granted = !data.granted;
                  console.log(data.granted);
                  change();
                }}
              />
            </div>
          {/each}
        </div>
      {:else}
        <!-- else content here -->
      {/if}
    </div>
  {/each}
</div>
