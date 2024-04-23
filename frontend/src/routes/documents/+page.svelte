<script>
  import { onMount } from "svelte";
  import { Heading, Button } from "flowbite-svelte";
  // import { userDataSession } from "../lib/utils/auth";
  import DataTable from "../lib/components/DataTable.svelte";
  import { _ } from "svelte-i18n";
  // import { PUBLIC_STRAPI_URL_DEV, PUBLIC_STRAPI_URL_PROD } from "$env/static/public";

  // const isDev = process.env.NODE_ENV === "development";

  // const strapiUrl = isDev ? PUBLIC_STRAPI_URL_DEV : PUBLIC_STRAPI_URL_PROD;

  // export let data;
  let documents = [];
  let page = 1;
  let pagination = { page: 1, pageSize: 25, pageCount: 1, total: 25 };
  let documentsReady = false;
  let allUrls = [];

  // const getAllDocument = async () => {
  //   const result = userDataSession();
  //   const address = result.user.identity.address;
  //   try {
  //     const response = await fetch(
  //       `/api/document?jwt=${result.jwt}&address=${address}&page=${page}`,
  //       {
  //         method: "GET",
  //         headers: {
  //           "Content-Type": "application/json",
  //         },
  //       },
  //     );

  //     if (response.ok) {
  //       const result2 = await response.json();
  //       documents = result2.data;
  //       pagination = result2.meta.pagination;
  //       for (const document of documents) {
  //         const signerData = document.attributes.signer_data;

  //         const urls = [];

  //         for (const item of signerData) {
  //           const address = Object.keys(item)[0];

  //           const response2 = await fetch(`/api/identity?jwt=${result.jwt}&address=${address}`, {
  //             method: "GET",
  //             headers: {
  //               "Content-Type": "application/json",
  //             },
  //           });

  //           if (response2.ok) {
  //             const result3 = await response2.json();
  //             if (result3.data.length > 0) {
  //               let url = "/";
  //               if (
  //                 result3.data[0].attributes.personalize.data !== null &&
  //                 result3.data[0].attributes.personalize.data.length > 0
  //               ) {
  //                 url = `${strapiUrl}${result3.data[0].attributes.personalize.data[0].attributes.url}`;
  //                 urls.push({ url: url });
  //               } else {
  //                 urls.push({ url: url });
  //               }
  //             }
  //           }
  //         }
  //         allUrls.push(urls);
  //       }
  //     }
  //   } catch (e) {
  //     console.log(e);
  //   }
  //   documentsReady = true;
  // };

  onMount(() => {
    // getAllDocument();
  });

  documents = [
    {
      attributes: {
        updatedAt: "2024-02-07T08:00:00Z",
        name: "Document 1",
        emailOwner: "abc@gmail.com",
        signerData: [
          {
            "abc@gmail.com": {
              approve: false,
            },
          },
        ],
        status: [
          {
            status: "Needs Signing",
            updatedAt: "2024-02-07T08:00:00Z",
          },
        ],
      },
    },
    {
      attributes: {
        updatedAt: "2024-02-06T10:30:00Z",
        name: "Document 2",
        emailOwner: "bcd@gmail.com",
        signerData: [
          {
            "bcd@gmail.com": {
              approve: false,
            },
          },
        ],
        status: [
          {
            status: "Approval Error",
            updatedAt: "2024-02-06T10:30:00Z",
          },
        ],
      },
    },
    {
      attributes: {
        updatedAt: "2024-02-05T15:45:00Z",
        name: "Document 3",
        emailOwner: "bcd@gmail.com",
        signerData: [
          {
            "abc@gmail.com": {
              approve: true,
            },
          },
        ],
        status: [
          {
            status: "Needs Finalize",
            updatedAt: "2024-02-05T15:45:00Z",
          },
        ],
      },
    },
    {
      attributes: {
        updatedAt: "2024-02-04T12:20:00Z",
        name: "Document 4",
        emailOwner: "abc@gmail.com",
        signerData: [
          {
            "abc@gmail.com": {
              approve: true,
            },
          },
        ],
        status: [
          {
            status: "Done",
            updatedAt: "2024-02-04T12:20:00Z",
          },
        ],
      },
    },
  ];

  allUrls = [
    [{ url: "" }, { url: "" }, { url: "" }, { url: "" }, { url: "" }, { url: "" }],
    [{ url: "" }, { url: "" }, { url: "" }],
    [{ url: "" }, { url: "" }, { url: "" }, { url: "" }],
    [{ url: "" }, { url: "" }, { url: "" }, { url: "" }],
  ];

  const getStatusDoc = (status) => {
    let statusType = status.status;
    let colorType = "brand";
    if (statusType === "Initiate") {
      colorType = "brand";
    }
    if (statusType === "Needs Signing") {
      colorType = "warning";
    }
    if (statusType === "Approval Error") {
      colorType = "error";
    }
    if (statusType === "Needs Finalize") {
      colorType = "brand";
    }
    if (statusType === "Done") {
      colorType = "success";
    }
    return {
      name: statusType,
      type: colorType,
      updatedAt: status.updatedAt,
    };
  };

  let email = "abc@gmail.com";

  const getActionBotton = (document) => {
    let approve = false;
    let reject = false;
    let recall = false;
    let finalize = false;
    for (const user of document.signerData) {
      if (!user[email]?.approve) {
        approve = true;
        reject = true;
      }
    }
    if (document.emailOwner === email) {
      recall = true;
      let allSigned = true;
      for (const user of document.signerData) {
        if (!user?.email?.approve) {
          allSigned = false;
          break;
        }
      }
      if (allSigned) {
        if (document.status[document.status.length - 1].status !== "Done") finalize = true;
      }
    }
    return {
      approve,
      reject,
      recall,
      finalize,
    };
  };

  let datas = documents.map((document, index) => {
    let lastStatus = document.attributes.status[document.attributes.status.length - 1];
    return {
      id: index + 1,
      updatedAt: document.attributes.updatedAt,
      name: document.attributes.name,
      status: getStatusDoc(lastStatus),
      recipients: allUrls[index],
      actions: {
        objects: [
          {
            iconOnly: true,
            actionButton: getActionBotton(document.attributes),
          },
        ],
      },
    };
  });

  let configs = {
    documentManagement: true,
    search: true,
    pagination: pagination,
  };
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <Heading
          tag="h3"
          class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
          >{$_("documentManagement")}</Heading
        >

        <!-- Document Table Data -->
        {#if documentsReady}
          <DataTable {datas} {configs} />
        {:else}
          <DataTable {datas} {configs} />
        {/if}
      </div>
    </div>
  </div>
</div>
