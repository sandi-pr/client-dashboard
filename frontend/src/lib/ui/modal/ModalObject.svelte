<script>
  import { Modal, Button } from "flowbite-svelte";
  import { Typography, FeaturedIcon, getIcons, InputRw as Input } from "$lib";
  import MultiSelect from "svelte-multiselect";
  import TabBar from "$lib/components/TabBar.svelte";
  import { _ } from "svelte-i18n";

  let { modalObjectOpen = $bindable() } = $props();

  let tabBarSelected = $state("Link");
  let linkImportLink = $state("");
  let linkChecked = $state(false);
  let allTags = $state([]);
  let tags = $state([]);

  let mainContentUpload = [
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
            },
          ],
        },
      },
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
        oninput: () => {},
      },
    },
    {
      text: {
        text: "Label",
        color: "primary",
        preset: "Text md/Semibold",
      },
    },
    {
      uploadFile: {
        onUpload: () => {},
        onDrop: () => {},
        files: [],
      },
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
      },
    },
  ];

  let mainContentLink = [
    {
      tabBar: {
        tabDatas: {
          selected: "Link",
          data: [
            {
              name: "Upload",
            },
            {
              name: "Link",
            },
          ],
        },
      },
    },
    {
      text: {
        text: $_("documentType"),
        color: "primary",
        preset: "Text md/Semibold",
      },
    },
    {
      input: {
        type: "text",
        onEdit: false,
        value: "Documents (TXT, PDF, DOC, DOCX)",
      },
    },
    {
      text: {
        text: $_("importLink"),
        color: "primary",
        preset: "Text md/Semibold",
      },
    },
    {
      input: {
        type: "text",
        onEdit: true,
        placeholder: "Paste image link (e.g. http://drive.google.com/..)",
        oninput: (e) => {
          linkImportLink = e.target.value;
          configModal.footer[0].button.buttonDisabled = e.target.value ? false : true;
        },
      },
    },
    {
      text: {
        text: $_("labels"),
        color: "primary",
        preset: "Text md/Semibold",
      },
    },
    {
      multiSelect: {
        alltags: allTags,
        tags: tags,
        allowUserOptions: "append",
        autoScroll: true,
        noMatchingOptionsMsg: "",
        createOptionMsg: null,
        ulSelectedClass: "text-primary-light dark:text-primary-dark",
        liSelectedClass: "border border-primary-light dark:border-primary-dark",
        ulOptionsClass: "bg-secondary-light dark:bg-tertiary-dark",
        outerDivClass:
          "bg-primary-light dark:bg-secondary-dark border rounded-lg border-gray-300 dark:border-gray-600 focus:none focus:ring-0 focus:border-gray-300 dark:focus:border-gray-600 py-2",
      },
    },
  ];

  let footerUpload = [
    {
      button: {
        text: "addNewDocument",
        onClick: () => {},
        buttonDisabled: false,
      },
    },
  ];

  let footerLink = [
    {
      button: {
        text: "checkLink",
        onClick: () => {
          // TODO: check link logic
          linkChecked = true;
          configModal.footer = footerUpload;
        },
        buttonDisabled: true,
      },
    },
  ];

  let configModal = $state({
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
    mainContent: mainContentLink,
    footer: footerLink,
  });
</script>

<Modal bind:open={modalObjectOpen}>
  <svelte:fragment slot="header">
    <div class="flex items-center gap-4">
      {#if configModal?.header}
        {#if configModal.header?.iconLeft}
          <FeaturedIcon classesSvg={getIcons(configModal.header.iconLeft?.icon)} />
        {/if}
        <div>
          {#each configModal.header?.contentCenter as config}
            <Typography {...config.text} />
          {/each}
        </div>
      {/if}
    </div>
  </svelte:fragment>

  {#if configModal?.mainContent}
    {#each configModal?.mainContent as content}
      {#if content?.tabBar}
        <TabBar
          tabDatas={content.tabBar?.tabDatas}
          on:tabChange={(e) => {
            tabBarSelected = e.detail;
            configModal.mainContent =
              tabBarSelected === "Upload" ? mainContentUpload : mainContentLink;
          }}
        />
      {/if}
      {#if content?.text}
        <Typography {...content.text} />
      {/if}
      {#if content?.input}
        <Input {...content.input} />
      {/if}
      {#if content?.multiSelect}
        <MultiSelect {...content.multiSelect} />
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
