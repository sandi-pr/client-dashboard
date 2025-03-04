<script>
  import { Modal, Button } from "flowbite-svelte";
  import { Typography, FeaturedIcon, getIcons } from "$lib";
  import TabBar from "$lib/components/TabBar.svelte";
  import { _ } from "svelte-i18n";

  let { modalObjectOpen = $bindable()} = $props();
  let configModal = {
    header: {
      iconLeft: {
        icon: "filePlus",
      },
      contentCenter: [
        {
          text: {
            text: "Add Document",
            color: "primary",
            preset: "Text lg/Semibold",
          },
        },
        {
          text: {
            text: "Upload & Organize Your Files Effortlessly",
            color: "tertiary",
            preset: "Text sm/Regular",
          },
        },
      ],
    },
    mainContent: [
      {
        tabBar: {
          tabDatas: {
            selected: "Upload",
            data: [
              {
                name: "Upload",
              },
              {
                name: "Link",
              }
            ],
            changeTabbar: () => {}
          }
        }
      },
      {
        text: {
          text: "Home",
          color: "primary",
          preset: "Text md/Semibold",
        },
      },
      {
        input: {
          type: "text",
          onEdit: false,
          value: "",
          oninput: () => {}
        }
      },
      {
        text: {
          text: "Label",
          color: "primary",
          preset: "Text md/Semibold",
        }
      },
      {
        uploadFile: {
          onUpload: () => {},
          onDrop: () => {},
          files: [],
        }
      },
      {
        text: {
          text: "Label",
          color: "primary",
          preset: "Text md/Semibold",
        },
      },
      {
        multiSelect: {
          alltags: [],
          tags: [],
        }
      },
    ],
    footer: [
      {
        button: {
          text: "addNewDocument",
          onClick: () => {},
          buttonDisabled: false,
        },
      }
    ]
  };

  function changeTabbar(selected) {
    configModal?.mainContent?.forEach((content) => {
      if (content?.tabBar?.tabDatas?.selected) {
        content.tabBar.tabDatas.selected = selected;
      }
    });
  }
</script>

<Modal bind:open={modalObjectOpen} size="xs">
  <svelte:fragment slot="header">
    <div class="flex items-center gap-4">
      {#if configModal?.header}
        {#if configModal.header?.iconLeft}
          <FeaturedIcon classesSvg={getIcons(configModal.header.iconLeft?.icon)} />
        {/if}
        <div>
          {#each configModal.header?.contentCenter as config}
            <Typography {...config.text}/>
          {/each}
        </div>
      {/if}
    </div>
  </svelte:fragment>

  {#if configModal?.mainContent}
    {#each configModal?.mainContent as content}
      {#if content?.tabBar}
        <TabBar tabDatas={content.tabBar?.tabDatas} on:tabChange={(e) => {}} />
        <!-- <p>{content.tabBar?.tabDatas.selected}</p> -->
      {/if}
    {/each}
  {/if}

  <svelte:fragment slot="footer">
    {#if configModal?.footer}
    {#each configModal?.footer as item}
      {#if item?.button}
        <Button
          onclick={item.button?.onClick}
          disabled={item.button?.buttonDisabled || false}
          class="w-full">{$_(item.button?.text)}</Button
        >
      {/if}
    {/each}
    {/if}
  </svelte:fragment>
</Modal>