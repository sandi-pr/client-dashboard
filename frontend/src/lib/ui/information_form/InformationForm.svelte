<script lang="ts">
  import { P, Textarea } from "flowbite-svelte";
  import { Typography } from "$lib";

  export let screenWidth: number;
  // export let user: {
  //   id: string;
  //   username: string;
  //   fullname: string;
  //   email: string;
  // };
  // syntax reference
  export let datas: {
    grid: number;
    title: string;
    titleGrid: number;
    innerGrid: number;
    data: {
      name: string;
      value: string;
      grid: number;
    }[];
  };
</script>

<!-- <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mx-4 mt-8">
  <P size="lg" weight="semibold" class="col-span-1">General Information</P>

  <div class="col-span-3 grid grid-cols-1 md:grid-cols-4 gap-4">
    <div class="col-span-4 [&>*]:text-hc-graylight-700 [&>*]:bg-white [&>*]:dark:bg-hc-graydark-950 [&>textarea]:text-hc-graylight-500">
      <P>User ID</P>
      <Textarea class="resize-none" rows={screenWidth > 1024 ? "1" : "2"} disabled value={user.id} />
    </div>
    <div class="col-span-4 [&>*]:text-hc-graylight-700 [&>*]:bg-white [&>*]:dark:bg-hc-graydark-950 [&>textarea]:text-hc-graylight-500">
      <P>Username</P>
      <Textarea class="resize-none" rows={screenWidth > 1024 ? "1" : "2"} disabled value={user.username} />
    </div>
    <div class="col-span-4 sm:col-span-2 [&>*]:text-hc-graylight-700 [&>*]:bg-white [&>*]:dark:bg-hc-graydark-950 [&>textarea]:text-hc-graylight-500">
      <P>Full name</P>
      <Textarea class="resize-none" rows={screenWidth > 1024 ? "1" : "2"} disabled value={user.fullname} />
    </div>

    <div class="col-span-4 sm:col-span-2 [&>*]:text-hc-graylight-700 [&>*]:bg-white [&>*]:dark:bg-hc-graydark-950 [&>textarea]:text-hc-graylight-500">
      <P>Email</P>
      <Textarea class="resize-none" rows={screenWidth > 1024 ? "1" : "2"} disabled value={user.email} />
    </div>

    <div class="col-span-4 sm:col-span-2 [&>*]:text-hc-graylight-700 [&>*]:bg-white [&>*]:dark:bg-hc-graydark-950 [&>textarea]:text-hc-graylight-500">
      <P>Role</P>
      <Textarea class="resize-none" rows={screenWidth > 1024 ? "1" : "2"} disabled value="Admin" />
    </div>
  </div>
</div> -->
<!-- syntax reference -->

<div class={`grid gap-4 ${screenWidth > 1024 ? "grid-cols-" + datas.grid : "grid-cols-1"}`}>
  <Typography preset="Text lg/Semibold" text={datas.title} propClass={`${screenWidth > 1024 ? "col-span-" + datas.titleGrid : "col-span-1"}`} />

  <div class={`col-span-${datas.grid - datas.titleGrid} grid gap-4 ${screenWidth > 1024 ? "grid-cols-" + datas.innerGrid : "grid-cols-1"}`}>
    {#each datas.data as data (data)}
      <div class={`${datas.innerGrid > data.grid ? (screenWidth > 1024 ? "col-span-" + data.grid : "col-span-" + datas.innerGrid) : "col-span-" + data.grid} [&>*]:text-hc-graylight-700 [&>*]:bg-white [&>*]:dark:bg-hc-graydark-950 [&>textarea]:text-hc-graylight-500`}>
        <P>{data.name}</P>
        <Textarea class="resize-none" rows={screenWidth > 1024 ? "1" : "2"} disabled value={data.value} />
      </div>
    {/each}
  </div>
</div>

<!--
@component
## Config
@prop screenWidth: number; // Screen width (Needed for responsive design).
@prop datas: object; // Data for the information form.
-->
