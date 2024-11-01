<script>
  import { P, Checkbox } from "flowbite-svelte";
  import { _ } from "svelte-i18n";

  import BookOpen2Icon from "$lib/icon/book-open-02.svelte";
  import Stars01Icon from "$lib/icon/stars-01.svelte";
  import MessageChatCircleIcon from "$lib/icon/message-chat-circle.svelte";

  export let permissionListDatas;

  let mainColSpan = 12 - permissionListDatas[0].data.length * 2;

  const borderInsideGrid = (index) => {
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
  // console.log(userListDatas);
</script>

<div class="w-full rounded-lg border-2 border-secondary-light dark:border-secondary-dark">
  {#each permissionListDatas as permission, index (permission.id)}
    <div
      class={"grid px-8 grid-cols-12 " +
        borderInsideGrid(index) +
        backgroundFromType(permission.type)}
    >
      {#if permission.type === "header"}
        <div class={"py-4 self-center col-span-" + mainColSpan}>
          <P class="text-lg text-tertiary-light dark:text-tertiary-dark">{permission.title}</P>
        </div>
        {#each permission.data as data (data.name)}
          <div class="py-4 self-center place-self-center col-span-2">
            <P class="font-semibold text-tertiary-light dark:text-tertiary-dark">{data.name}</P>
          </div>
        {/each}
      {:else if permission.type === "section"}
        <div
          class={"py-4 self-center col-span-" +
            mainColSpan +
            " " +
            leftPaddingFromLevel(permission.level)}
        >
          <div class="flex gap-3">
            {#if permission.icon === "book"}
              <BookOpen2Icon svgClass="w-7 h-7 text-tertiary-light dark:text-tertiary-dark" />
            {:else if permission.icon === "star"}
              <Stars01Icon svgClass="w-7 h-7 text-tertiary-light dark:text-tertiary-dark" />
            {:else if permission.icon === "chat"}
              <MessageChatCircleIcon
                svgClass="w-7 h-7 text-tertiary-light dark:text-tertiary-dark"
              />
            {:else}
              <!-- else content here -->
            {/if}
            <P class="text-lg font-bold text-tertiary-light dark:text-tertiary-dark"
              >{permission.title}</P
            >
          </div>
        </div>
      {:else if permission.type === "content"}
        <div class={"py-6 self-center col-span-" + mainColSpan}>
          <P class="text-tertiary-light dark:text-tertiary-dark">{permission.title}</P>
        </div>
        {#each permission.data as data (data.name)}
          <div class="py-6 self-center place-self-center col-span-2">
            <Checkbox bind:checked={data.granted} />
          </div>
        {/each}
      {:else}
        <!-- else content here -->
      {/if}
    </div>
  {/each}
</div>
