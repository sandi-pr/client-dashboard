<script>
  import {
    P,
    Table,
    Button,
    TableBody,
    TableHead,
    Pagination,
    TableBodyRow,
    TableBodyCell,
    TableHeadCell,
  } from "flowbite-svelte";

  import { _ } from "svelte-i18n";
  import { goto } from "$app/navigation";

  import ArrowLeft from "$lib/icon/arrow-left.svelte";
  import ArrowRight from "$lib/icon/arrow-right.svelte";

  export let userListDatas;
  export let metadata;
  export let next;
  export let previous;

  const pillStyle = (style) => {
    if (style === "status") {
      return "py-[2px] px-2 border-[1px] rounded-xl text-center w-min bg-utility-success-50-light dark:bg-utility-success-50-dark border-utility-success-200-light dark:border-utility-success-200-dark text-utility-success-700-light dark:text-utility-success-700-dark";
    } else if (style === "super admin") {
      return "py-[2px] px-2 border-[1px] rounded-xl text-center w-min bg-utility-brand-50-light dark:bg-utility-brand-50-dark border-utility-brand-200-light dark:border-utility-brand-200-dark text-utility-brand-700-light dark:text-utility-brand-700-dark";
    } else if (style === "admin") {
      return "py-[2px] px-2 border-[1px] rounded-xl text-center w-min bg-utility-blue-50-light dark:bg-utility-blue-50-dark border-utility-blue-200-light dark:border-utility-blue-200-dark text-utility-blue-700-light dark:text-utility-blue-700-dark";
    } else if (style === "manager") {
      return "py-[2px] px-2 border-[1px] rounded-xl text-center w-min bg-utility-orange-50-light dark:bg-utility-orange-50-dark border-utility-orange-200-light dark:border-utility-orange-200-dark text-utility-orange-700-light dark:text-utility-orange-700-dark";
    } else if (style === "user") {
      return "py-[2px] px-2 border-[1px] rounded-xl text-center w-min bg-utility-tertiary-50-light dark:bg-utility-tertiary-50-dark border-utility-tertiary-200-light dark:border-utility-tertiary-200-dark text-utility-tertiary-700-light dark:text-utility-tertiary-700-dark";
    } else {
      return "py-[2px] px-2 border-[1px] rounded-xl text-center w-min bg-utility-warning-50-light dark:bg-utility-warning-50-dark border-utility-warning-200-light dark:border-utility-warning-200-dark text-utility-warning-700-light dark:text-utility-warning-700-dark";
    }
  };

  const pagesGenerator = () => {
    let length = metadata.total / metadata.pageSize;
    length = Math.ceil(length);
    let pages = [];
    if (length > 6) {
      for (let i = 0; i < 3; i++) {
        pages.push({
          name: `${i + 1}`,
        });
      }
      pages.push({
        name: "...",
      });
      // push the last 3 pages
      for (let i = length - 3; i < length; i++) {
        pages.push({
          name: `${i + 1}`,
        });
      }
    } else {
      for (let i = 0; i < length; i++) {
        pages.push({
          name: `${i + 1}`,
        });
      }
    }
    return pages;
  };

  let pages = pagesGenerator();

  // console.log(userListDatas);
</script>

<div
  class="mt-4 rounded-lg [&>div]:sm:rounded-lg bg-gray-50 dark:bg-gray-700 [&>div]:border [&>div]:border-primary-light [&>div]:dark:border-primary-dark"
>
  <Table hoverable={true}>
    <TableHead>
      <!-- <TableHeadCell>{$_("userId")}</TableHeadCell> -->
      <TableHeadCell>{$_("fullname")}</TableHeadCell>
      <TableHeadCell>{$_("username")}</TableHeadCell>
      <TableHeadCell>{$_("email")}</TableHeadCell>
      <TableHeadCell>{$_("role")}</TableHeadCell>
      <TableHeadCell>{$_("status")}</TableHeadCell>
      <TableHeadCell></TableHeadCell>
    </TableHead>
    <TableBody tableBodyClass="divide-y">
      {#each userListDatas as userListData (userListData.id)}
        <TableBodyRow>
          <!-- <TableBodyCell>
            <P>{userListData.id}</P>
          </TableBodyCell> -->
          <TableBodyCell><P class="font-semibold">{userListData.fullname}</P></TableBodyCell>
          <TableBodyCell><P>{userListData.username}</P></TableBodyCell>
          <TableBodyCell><P>{userListData.email}</P></TableBodyCell>
          <TableBodyCell>
            {#if userListData.role === "Super Admin"}
              <div class={pillStyle("super admin")}>
                {userListData.role}
              </div>
            {:else if userListData.role === "Admin"}
              <div class={pillStyle("admin")}>
                {userListData.role}
              </div>
            {:else if userListData.role === "Manager"}
              <div class={pillStyle("manager")}>
                {userListData.role}
              </div>
            {:else if userListData.role === "User"}
              <div class={pillStyle("user")}>
                {userListData.role}
              </div>
            {:else}
              <div class={pillStyle("other")}>
                {userListData.role}
              </div>
            {/if}
          </TableBodyCell>
          <TableBodyCell
            ><div class={pillStyle("status")}>
              {userListData.status}
            </div></TableBodyCell
          >
          <TableBodyCell>
            <Button
              on:click={() => {
                goto(userListData.href);
              }}
              color="white"
              class="!p-2 focus-within:!ring-1 text-brand-tertiary-light dark:text-brand-tertiary-dark"
              size="lg"
            >
              &nbsp;{$_("detail")}
            </Button>
          </TableBodyCell>
        </TableBodyRow>
      {/each}
    </TableBody>
  </Table>
</div>

<div class="flex justify-between items-center !rounded-none !rounded-b-lg !border-t-0">
  <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
    {$_("showing")}
    <span>{metadata.page === 1 ? 1 : (metadata.page - 1) * metadata.pageSize + 1}</span>
    {$_("to")}
    <span
      >{metadata.page * metadata.pageSize >= metadata.total
        ? metadata.total
        : metadata.page * metadata.pageSize}</span
    >
    {$_("of")}
    <span>{metadata.total}</span>
    {$_("entries")}
  </div>

  <div class="m-4">
    <Pagination
      on:next={next}
      on:previous={previous}
      {pages}
      activeClass="bg-hc-graylight-100 dark:bg-hc-graydark-700 dark:border-[#374151] dark:text-hc-graydark-300"
    >
      <span slot="prev"
        ><div class="flex items-center">
          <ArrowLeft svgClass="w-5 h-5" />&nbsp;{$_("previous")}
        </div></span
      >
      <!--  -->
      <span slot="next"
        ><div class="flex items-center">
          {$_("next")}&nbsp;<ArrowRight svgClass="w-5 h-5" />
        </div></span
      >
    </Pagination>
  </div>
</div>
