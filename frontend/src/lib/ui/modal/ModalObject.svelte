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

  let {
    modalObjectOpen = $bindable(false),
    handleFileUploadWrapper,
    files = $bindable([]),
    isProcess = $bindable(false),
    addDocumentByCollection,
    alltags = [],
    tags = $bindable([]),
  } = $props();

  let tabBarSelected = $state("Upload");
  let linkImportLink = $state("");
  let linkChecked = $state(false);
  // let alltags = $state([]);
  // let tags = $state([]);
  // let isProcess = $state(false);
  // let files = $state([]);

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
        text: $_("documentType"),
        color: "secondary",
        preset: "Text sm/Medium",
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
      uploadFile: {
        title: $_("uploadDocuments"),
      },
    },
    {
      text: {
        text: "Label",
        color: "primary",
        preset: "Text sm/Medium",
      },
    },
    {
      multiSelect: {
        // selected: tags,
        // options: alltags,
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
        preset: "Text sm/Medium",
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
        preset: "Text sm/Medium",
      },
    },
    {
      multiSelect: {
        // selected: tags,
        // options: alltags,
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
        onClick: addDocumentByCollection,
        buttonDisabled: true,
      },
    },
  ];

  let footerLink = [
    {
      button: {
        text: "checkLink",
        onClick: async () => {
          // console.log(linkImportLink);
          let urlValid = await urlChecker(linkImportLink, [
            "application",
            "application/pdf",
            "application/msword",
            "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
          ]);
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
    mainContent: mainContentUpload,
    footer: footerUpload,
  });

  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  async function readFileAsBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        const base64Data = reader.result.split(",")[1];
        resolve(base64Data);
      };
      reader.onerror = reject;

      reader.readAsDataURL(file);
    });
  }

  const urlChecker = async (url, limitation = ["application", "application/pdf"]) => {
    let errorObject;
    try {
      let response;
      if (url.includes("drive.google.com")) {
        const fileId = url.match(/[-\w]{25,}/);
        if (!fileId) {
          throw new Error("Invalid Google Drive URL");
        }
        const gdriveUrl = `https://www.googleapis.com/drive/v3/files/${fileId[0]}?alt=media&key=AIzaSyDi62RQsNNjA-_IfGyaDfWc7L11dVm75YI`;
        response = await fetch(gdriveUrl);
      } else {
        response = await fetch(url);
      }

      if (!response.ok) {
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.includes("application/json")) {
          errorObject = await response.json();
        } else {
          errorObject = { error: { message: "Non-JSON response received" } };
        }
        throw new Error("Network response was not ok");
      }

      const contentType = response.headers.get("content-type");
      console.log("contentType", contentType);
      if (limitation.some((type) => contentType.includes(type))) {
        const blob = await response.blob();
        const base64data = await readFileAsBase64(blob);
        const fileType = contentType.split("/")[1];
        const fileName = `${url.split("/")[3].slice(0, 10)}.${fileType}`;
        files.push({
          id: files.length < 1 ? 1 : files.length + 1,
          name: fileName,
          size: formatBytes(blob.size),
          base64: base64data,
          ext: `.${fileType}`,
          type: contentType.includes("image")
            ? "image"
            : contentType.includes("text/plain")
              ? "text"
              : contentType.includes("application/pdf")
                ? "pdf"
                : contentType.includes("application/msword") ||
                    contentType.includes(
                      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
                    )
                  ? "document"
                  : "unknown",
        });
        return true;
      } else {
        configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text =
          $_("fileNotSupported");
        throw new Error("URL does not contain a valid resource");
      }
    } catch (error) {
      console.log("errorObject", errorObject);
      // console.log("error", error);
      if (errorObject && errorObject.error) {
        if (errorObject.error.code === 403) {
          configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text = $_(
            "accessRestrictedCheckSharingSettings",
          );
        } else if (errorObject.error.code === 404) {
          configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text = $_(
            "fileNotFoundMakeSureTheLinkIsCorrectAndAccessible",
          );
        }
      } else if (error instanceof TypeError && error.message === "Failed to fetch") {
        configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text = $_(
          "accessRestrictedCheckSharingSettings",
        );
      } else {
        configModal.mainContent[5].card.cardData.contentCenter[0][0].text.text =
          $_("failedToFetchResource");
      }
      console.error("Error fetching the URL:", error);
      return false;
    }
  };

  $effect(() => {
    configModal.footer[0].button.buttonDisabled = files.length < 1 ? true : false;
  });
</script>

<Modal
  bind:open={modalObjectOpen}
  bodyClass="p-4 md:p-5 !space-y-4 flex-1 overflow-y-auto overscroll-contain"
>
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
            files = [];
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
      {#if content?.uploadFile}
        <UploadFile
          title={content.uploadFile?.title}
          handleFileUploadWrapper={() => handleFileUploadWrapper()}
          file={files[0]}
          {files}
          bind:isProcess
          multiple={true}
          actionButton={true}
          actionView={true}
          actionDelete={() => {
            files = [];
          }}
        />
      {/if}
      {#if content?.multiSelect}
        <MultiSelect {...content.multiSelect} bind:selected={tags} options={alltags} />
      {/if}
      {#if content?.card && content.card?.cardData.contentCenter[0][0].text.text}
        <Card {...content.card} />
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
