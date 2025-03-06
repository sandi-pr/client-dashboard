<script>
  import { Modal, Button } from "flowbite-svelte";
  import {
    Typography,
    FeaturedIcon,
    getIcons,
    InputRw as Input,
    CardFeature as Card,
    UploadFile,
  } from "$lib";
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
  let isProcess = $state(false);
  let files = $state([]);

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
        id: "link",
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

  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  const urlChecker = async (url) => {
    try {
      let response;
      console.log("url", url);
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

      console.log("response", response);

      if (!response.ok) {
        configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text =
          $_("networkResponseWasNotOk");
        throw new Error("Network response was not ok");
      }

      const contentType = response.headers.get("content-type");
      if (contentType.includes("image") || contentType.includes("application")) {
        const blob = await response.blob();
        const reader = new FileReader();
        reader.readAsDataURL(blob);
        reader.onloadend = () => {
          const base64data = reader.result;
          const fileType = contentType.split("/")[1];
          const fileName = `${url.split("/")[3].slice(0, 10)}.${fileType}`;
          files.push({
            id: files.length < 1 ? 1 : files.length + 1,
            name: fileName,
            document_title: "Document Title", // Replace with actual title if available
            description: "Document Description", // Replace with actual description if available
            size: formatBytes(blob.size),
            base64: base64data,
            ext: `.${fileType}`,
            type: contentType.includes("image") ? "image" : "document",
          });
          console.log("Resource fetched successfully:", base64data);
        };
        return true;
      } else {
        configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text = $_(
          "urlDoesNotContainAValidResource",
        );
        throw new Error("URL does not contain a valid resource");
      }
    } catch (error) {
      configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text =
        $_("failedToFetchResource");
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
      {#if content?.input && content.input.id === "link" && linkChecked}
        <UploadFile
          title={`${$_("uploadDocuments")}*`}
          handleFileUploadWrapper={() => {}}
          file={files[0]}
          {files}
          bind:isProcess
          multiple={false}
          actionView={() => {
            console.log("hello");
          }}
          actionButton={true}
          actionDelete={() => {
            files = [];
            linkChecked = false;
            linkImportLink = "";
            configModal.footer = footerLink;
          }}
        />
      {:else if content?.input}
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
