<script>
  import { Modal, Button } from "flowbite-svelte";
  import { Typography, FeaturedIcon, getIcons, InputRw as Input, CardFeature as Card } from "$lib";
  import { getBgColors } from "$lib";
  import MultiSelect from "svelte-multiselect";
  import TabBar from "$lib/components/TabBar.svelte";
  import { _ } from "svelte-i18n";

  let { modalObjectOpen = $bindable() } = $props();

  let tabBarSelected = $state("Link");
  let linkImportLink = $state("");
  let linkChecked = $state(false);
  let allTags = $state([]);
  let tags = $state([]);
  let imageBlob = $state(null);

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
        placeholder: "Paste image link (e.g. https://plus.unsplash.com/...)",
        oninput: (e) => {
          linkImportLink = e.target.value;
          configModal.footer[0].button.buttonDisabled = e.target.value ? false : true;
        },
      },
    },
    {
      card: {
        cardData: {
          // iconData: { icon: "filePlus" },
          contentCenter: [
            [
              {
                icon: {
                  icon: "infoCircle",
                  color: "error",
                },
                text: {
                  text: "",
                  color: "fg-error-primary",
                  preset: "Text sm/Regular",
                },
                justify: "justify-start",
              },
            ],
          ],
        },
        propClass: getBgColors("error-primary") + " border-0 rounded-xl",
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
        onClick: async () => {
          // console.log(linkImportLink);
          let urlValid = await urlChecker(linkImportLink);
          if (urlValid) {
            linkChecked = true;
            configModal.footer = footerUpload;
            configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text = "";
          } else {
            linkChecked = false;
            linkImportLink = "";
          }
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

  const urlChecker = async (url) => {
    try {
      let response;
      if (url.includes("drive.google.com")) {
        const fileId = url.match(/[-\w]{25,}/);
        if (!fileId) {
          throw new Error("Invalid Google Drive URL");
        }
        const gdriveUrl = `https://www.googleapis.com/drive/v3/files/${fileId[0]}?alt=media&key=YOUR_API_KEY`;
        response = await fetch(gdriveUrl);
      } else {
        response = await fetch(url);
      }

      if (!response.ok) {
        // errorMessage = $_("networkResponseWasNotOk");
        configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text =
          $_("networkResponseWasNotOk");
        throw new Error("Network response was not ok");
      }

      const contentType = response.headers.get("content-type");
      if (contentType.includes("image") || contentType.includes("application")) {
        imageBlob = await response.blob();
        // Save the resource to a variable
        // For example: this.resource = resource;
        console.log("Resource fetched successfully:", imageBlob);
        return true;
      } else {
        // errorMessage = $_("urlDoesNotContainAValidResource");
        configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text = $_(
          "urlDoesNotContainAValidResource",
        );
        throw new Error("URL does not contain a valid resource");
      }
    } catch (error) {
      console.error("Error fetching the URL:", error);
      return false;
    }
  };
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
            configModal.footer =
              tabBarSelected === "Upload" ? footerUpload : linkChecked ? footerUpload : footerLink;
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
      {#if content?.card && content.card?.cardData.contentCenter[0][0].text.text}
        <Card {...content.card} />
      {/if}
    {/each}
  {/if}

  <!-- {#if imageBlob}
    <div class="image-viewer">
      <img src={URL.createObjectURL(imageBlob)} alt="Fetched Resource" />
    </div>
  {/if} -->

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
