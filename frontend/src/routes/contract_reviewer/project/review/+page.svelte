<script>
  import { onMount, onDestroy } from "svelte";
  import { Button, Modal, P } from "flowbite-svelte";
  import { Typography, triggerToast, triggerLoading } from "$lib";
  import PdfViewer from "svelte-pdf";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import FileDownloadIcon from "$lib/icon/file-download-03.svelte";
  import DownloadIcon from "$lib/icon/download-01.svelte";
  import { ChevronDownOutline, RefreshOutline } from "flowbite-svelte-icons";
  import templatePdf from "$lib/template.pdf";
  // import PdfViewer from 'svelte-pdf-scrollable';
  import DataTable from "$lib/components/DataTable.svelte";
  import DataTableAiReview from "$lib/components/DataTableAIReview.svelte";
  // import { Svroller, Svrollbar } from 'svrollbar';
  import { CR3Call } from "hexacomponents-svelte";
  import { page } from "$app/stores";
  import Markdown from "svelte-exmarkdown";
  import { gfmPlugin } from "svelte-exmarkdown/gfm";
  import { goto } from "$app/navigation";
  import { ListMarkdown } from "$lib";
  // import jsPDF from "jspdf";

  // import html2canvas from "html2canvas";

  import { call } from "$lib";

  const plugins = [gfmPlugin()];

  let projectDocumentId;
  let projectName = $state("");
  let configs = $state({});
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });
  let dropdownToggle = $state([true, true, false]);
  let documentViewModal = $state(false);
  let editModal = $state(false);
  let tmpResult = $state([]);
  let tmpTitle = $state("");
  let tmpText = $state("");
  let selectedFile = $state(null);
  let pdfData = $state([null, null, null, null]);
  let base64Data = $state();
  let fileUrl = $state(null);
  let datas = $state([]);
  let datas2 = $state([]);
  let markdown = $state("");
  let markdownDummy = $state(`

  Jakarta, 15 June 2023  
  Number: R-021/PNC400000/2023-S3  
  Enclosure: 1 (one) set of files  
  Subject: **Confirmation of Gasoline 90 RON for July-December 2023 (Term) – Aramco (DAP/CFR)**

  Mr/Mrs/Ms  
  Kelvin Lee / Joseph Lim  
  Aramco Trading Singapore Pte. Ltd.

  With reference to our email correspondences dated 11-30 May 2023 captioned OFFER: Competitive Negotiation of Gasoline 90 RON Supply for H2 2023 Delivery CFR/DAP (Term) - REVISE 2 and email correspondences dated 31 May-9 June 2023 captioned Recap: Competitive Negotiation of Gasoline 90 RON Supply for H2 2023 Delivery CFR/DAP (Term) (PTP23SS00012), we hereby accept and confirm your proposal with Deal Confirmation as follows:

  | Kolom 1   | Kolom 2   |
  |-----------|-----------|
  | Data 1    | Data 2    |
  | Data 3    | Data 4    |
  | Data 5    | Data 6    |

  The provisions in this Deal Confirmation shall govern the transaction between the Buyer and the
  Seller until execution of an Agreement is confirmed or signed by an authorized person from each
  Party. In the event of conflict between provision(s) in such Deal Confirmation and provision(s) in
  the non-confirmed Agreement, then the provision(s) in such Deal Confirmation shall prevail.

  For the avoidance of doubt to the validity of the Deal Confirmation, any typographical error in the
  Deal Confirmation, for the sake of operability of the deal, may be amended by the Buyer in writing.

  In the event of the Seller would like to have additional provision(s) under the Agreement which
  this GTC does not confirm, the Buyer has the right to counter, defend, or reject such provision(s).
  In the absence of stating counter, defense, or rejection to such provision(s) shall not be construed
  as acceptance to all or part of provision(s) written in the Seller's Agreement.

  Please be acknowledged.

  VP Trading Operations

  **Maya Kusmaya**

  Carbon copy:
  1. Manager Product Operation
    PT Pertamina Patra Niaga
  2. Manager Contract & Settlement
    PT Pertamina Patra Niaga
`);

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

  function base64toBlob(base64Data, fileType) {
    const byteCharacters = atob(base64Data);
    const byteArrays = [];

    for (let offset = 0; offset < byteCharacters.length; offset += 512) {
      const slice = byteCharacters.slice(offset, offset + 512);

      const byteNumbers = new Array(slice.length);
      for (let i = 0; i < slice.length; i++) {
        byteNumbers[i] = slice.charCodeAt(i);
      }

      const byteArray = new Uint8Array(byteNumbers);
      byteArrays.push(byteArray);
    }

    const blob = new Blob(byteArrays, { type: fileType });
    return blob;
  }

  // async function createPdfUrl(file) {
  //   const response = await fetch(templatePdf);
  //   const blob = await response.blob();
  //   base64Data = await readFileAsBase64(blob);

  //   const pdfBlob = base64toBlob(base64Data, blob.type);
  //   fileUrl = URL.createObjectURL(pdfBlob);
  // }

  // onMount(async () => {
  //   // await getDocument(1);
  //   await createPdfUrl();
  // });

  let dummyDatas = $state([
    {
      id: 1,
      no: 1,
      section: {
        name: "Buyer",
        prompt: "What is the name of company as Buyer?",
        result:
          "Contacts:    \nOperations : Fax:+65 6225 5388; Tel: +65 6224  2228;  M: +65 9139 5056 ; keemeng.wee@aramcotrading.sg    \nAccounting  /Invoicing:  Fax;+966 13 875 -8484; Tel:+966 13 880 -1499  or1488 ; accounting@aramcotrading.com  /ATS-Settlement@aramcotrading.sg     \nTreasury/Letter of Credit:  Fax: :+966 13 875 -8484 Tel: +966 13 880 -1588, M: +966 55 -5269 418; treasury@aramcotradin g.com    \nNew demurrage claim shall be directed  To: Tel: +966  13 880-8720/ Tel: +966 13 872 -0147 ; NewDemurrageClaim@aramcotrading.com  \nFor all other related  correspondence demurrage email :  claims@aramcotrading.com   \nPage 4 of 4 \n DEAL NO.  PTP23SS00012  \ncoefficient value is 4.0 0 USD/ Barrel . For the avoidance of doubt, this estimation is solely for \nour customs clearance purpose. Number of clearance document must be align, for example \nif there are Two (2) B/L documents, then the packing list and custom invoice must be made \nin 2 (T wo) documents which represent each B/L document. All clearance documents (B/L or \nendorsed B/L, Packing List, and Customs Invoice) must be sent to Buyer at least Six (6) hours \nafter vessel depart from loading port. For CFR or DAP basis, vessel clearance mus t comply with \nIndonesia customs regulation such as but not limited to (RKSP (transcript BC1.0), inward \nmanifest (transcript BC1.1), and discharged permit from local customs office). All documents \nabove must be submitted immediately or Twenty -Four (24) hour s after vessel departure from \nloading port or Three (3) days before vessel arrive at discharging port, whichever is earlier, \nexcept for loading port with less than Twenty -Four (24) hours sea days to discharge port.  \n \nIn the event that the Seller fails to s ubmit such documents in good order and/or in \nconformance with delivered quantity and/or within the time limit, all costs in relation to the \ncustom clearance charges and/or custom fine, including but not limited to demurrage \nincurred by late custom clearanc e shall be entirely on the Seller’s account.  \n19. Law \u0026 \nArbitration:  English Law and Singapore Arbitration (SIAC)  \n20. Title \u0026 Risks:  CFR = load port  \nDAP = disport  \n21. Others:  Pricing to be finalized to 3 decimal places  \n22. Contact Details : • Aramco Trading Singapore PTE.  LTD. (ATS) will not be held responsible for the delays and \nno claim for demurrage shall be entertained if email claim does not reach \nNewDemurrageClaim@aramcotrading.com  \n• All correspondence relating to demurrage or demurrage claim (if any) shall be directed to \nthe following department in charge:  \nEmail:  claims@aramcotrading.com  \nTelephone no.:  +966 -13 880-8720 / +966 13 872 0147  \n• As afore -stipulated, resulting from the correspondence sent to any other place.  \nWe are pleased to have concluded this business with your company.  \n \nBest R egards,  \nNada Soqair  \nAramco Trading Singapore PTE. LTD .  \n******** ************* ********* ****** IMPORTANT NOTICE  ***************** ****** *******************  \nARAMCO TRADING COMPANY (AND ITS SUBSIDIARIES)  WILL NOT BE HELD RESPONSIBLE FOR DELAYS \nRESULTING FROM CORRESPONDENCE SENT TO ANY OTHER PLACE.  \nTHE ABSENCE OF A RESPONSE AND/OR DULY EXECUTED TRADE CONFIRMATION ON TH E PART OF EITHER \nPARTY OR BOTH PARTIES SHALL NOT AFFECT THE VALIDITY OF AND BINDING NATURE OF THE AGREEMENT.  \nPLEASE NOTE THAT ALL NOTICES FOR THE ATTENTION OF ARAMCO TRADING SINGAPORE PTE. LTD. MUST BE \nSENT TO THE OPERATIONS DEPARTMENT, ACCOUNTING SERVICES /INVOICES OR TREASURY LETTERS OF \nCREDIT DEPARTMENT RESPECTIVELY AT THE CONTACT ADDRESS GIVEN IN THE FOOTER BELOW  \n*************** ******** ******************* *******************************************************",
      },
      reference_documents: {
        id: 3,
        name: "GTC_version_feb2024.pdf",
        prompt: "What is the name of company as Buyer?",
        result:
          "PT Pertamina Patra Niaga GTC – Purchases of Petroleum Products                    \n \n                           Page 11 of 35 \n represent each B/L document. All clearance documents (B/L, Packing List, and Customs \nInvoice) have to be sent to Buyer at least six  (6) hours after vessel depart from loading port. For \nCFR or DAP basis, vessel clearance must comply with Indonesia customs regulation such as \nbut not limited to (RKSP (transcript BC1.0), inward manifest (transcript BC1.1), and discharged  \nperm it from local customs office).  All documents above must be submitted immediately or \ntwenty -four (24) hours after vessel departure from loading port or three (3) days before vessel \narrive at discharging port, whichever is earlier , except for loading port with less than twenty -four \n(24) hours sea days to discharge port.  \nIn the event that the Seller fails to submit such doc uments in good order  and/or in conformance \nwith delivered quantity and/or within the time limit, all costs in relation to the custom clearance \ncharges and/or custom fine, including but not limited to demurrage incurred by late custom clearance shall be ent irely on the Seller’s account.  \n17.8.  For CFR  or DAP basis, the shipper must comply with all Indonesian regulations. In the event \nthe Seller fails to tender delivery due to any violation caused by the shipper, Seller without any \nreasons whatsoever has to commit to deliver the cargo within the agreed ADD. All associated \ncost thereof shall be solely borne by the Seller.  \n17.9.  For CFR  or DAP basis, the shipping agent has to process and submit the discharging \npermission letter to the customs office in discharging port area. Seller must make sure that the \nshipping agent at discharging port submit discharge permit to local customs office at discharging \nport. \nFor CFR  or DAP basis, in the event of quantity discrepancy at the discharging port, buyer shall \nnot be liable f or all customs fines and penalties caused by any excess or shortage related to \nthe discharging custom tolerance regulation.  \n17.10.  The Buyer reserves all rights to claim for loss or damages including but not limited to change(s) \nmade to declared loading port. For  CFR or DAP delivery, custom clearance charges and/or \ncustom fine and/or paid government taxes and/or random inspection(s) pertaining to Red channel policy that may arise, including but not limited to, in the event of (a) wrong submission \nof custom documents, (b) error in custom documents submitted, (c) quantity or quality \ndiscrepancy attributable to the Seller including but not limited to outturn losses, and (d) rejection \nof cargo(es) due to off specification issue, including but not limited to any losses or damages \nsuffered by the Buyer resulting from such failure. In the event of any such claims, the Buyer \nshall notify the Seller of its claim in writing, and the Seller shall endeavor full settlement of all \nsuch claims.  \n17.11.  The Seller shall settle any other costs, expenses or charges (including but not limited to port \ncharges, agent fee, tug boat fee, inspector and /or any mutually agreed shared inspection costs, \ncharges incurred for acquiring LTBP rate, claim for quantity and/or quality losses or shortages, Delay of Delivery charge(s) and /or for FOB delivery, demurrage claim etc) which may arise \nunder the terms of this GTC  and/or the Agreement . \n17.12.  In case of any discrepancy or conflict among the terms and conditions in purchase order , the \nterms and conditions in this GTC, the terms and conditions in the Agreement and the terms and \nconditions in Incoterms for CFR Basis (latest edition)  or DAP Basis (latest edition) or FOB Basis \n(latest edition), then the terms and conditions shall prevail with the following sequenc e: \n1) Dealdone/Confirmation Letter along with the general terms and conditions agreed;  \n2) Agreement ; \n3) Purchase Order  (if any) ; \n4) Incoterms for CFR Basis (latest edition)  or DAP Basis (latest edition)  or FOB Basis (latest \nedition);  \n17.13.  CFR, DAP  and FOB shall have the meaning given to that term in Incoterms (or any subsequent \namendment thereto) as modified by the terms and conditions of this GTC  and/or the Agreement . \nFurther, if there is any inconsistency or conflict between this GTC and Incoterms, this GTC shall prevail.  \n17.14.  The Seller shall send the draft Agreement via email to trade.support@pertamina.com , unless \nspecified otherwise by the Buyer.",
      },
      ai_review: "\n**Buyer** : PT Pertamina Patra Niaga Jakarta, Indonesia\n",
      ai_suggestion:
        "Based on the provided contract, here are some answers to questions related to it:\n\n\n1. How many days before vessel arrival at discharging port, whichever is earlier, except for loading port with less than twenty-four (24) hours sea days to discharge port must Seller submit documents in good order and/or in conformance with delivered quantity and/or within the time limit?\n\n\n**Answer:** 3 days\n\n\n2. If the Seller fails to submit such documents as per point above, who shall be entirely on the account for all costs relating to custom clearance charges and/or custom fine and including but not limited to demurrage incurred by late custom clearance?\n\n\n**Answer:** Seller's account\n\n\n3. For CFR or DAP basis, in the event of quantity discrepancy at discharging port, who shall not be liable for all customs fines and penalties caused by any excess or shortage related to discharging custom tolerance regulation?\n\n\n**Answer:** Buyer\n\n\n4. The parties agreed that CFR, DAP, and FOB terms should have their meanings given to those terms as per Incoterms (or subsequent amendment) modified by the terms and conditions of GTC and/or Agreement; if there's inconsistency between this GTC and Incoterms, which one shall prevail?\n\n\n**Answer:** This GTC",
    },
  ]);

  async function getProjectResult() {
    triggerLoading(true);
    datas = [];
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "get-contract-project-result", {
          URLQuery: {},
          URLParams: {
            id: projectDocumentId,
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          console.log(response);
          return;
        }

        let projectId = response.result.data[0].attributes.project_id;

        const response2 = await $call.call("client-dashboard", "get-contract-project", {
          URLQuery: {},
          URLParams: {
            id: projectId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          if (init) triggerLoading(false);
          triggerToast("error");
          console.log(response2);
          return;
        }

        projectName = response2.result.data[0].attributes.project_name;

        let referenceDocumentNames = {};

        markdown = response.result.data[0].attributes.result_markdown;
        tmpResult = response.result.data[0].attributes.result;

        const referenceDocumentIds = tmpResult
          .flatMap((section) =>
            section.generated.results_query.flatMap((item) => item.reference_document_ids),
          )
          .filter((value, index, self) => self.indexOf(value) === index);

        // console.log("referenceDocumentIds", referenceDocumentIds);

        for (const referenceDocumentId of referenceDocumentIds) {
          const response3 = await $call.call("client-dashboard", "get-reference-document-by-id", {
            URLQuery: {},
            URLParams: {
              id: referenceDocumentId.toString(),
            },
            Params: {},
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response3.status !== 200) {
            triggerLoading(false);
            triggerToast("error");
            console.log(response3);
            return;
          }

          if (response3.result.data.length === 0) {
            continue;
          } else {
            referenceDocumentNames[referenceDocumentId] =
              response3.result.data[0].attributes.document_name;
          }
        }

        let idCounter = 1; // Initialize global counter

        let promise = response.result.data[0].attributes.result.flatMap((item) => {
          return item.generated.results_query.map((generatedItem) => {
            const documentId = generatedItem.reference_document_ids[0];

            return {
              no: idCounter,
              id: idCounter++, // Increment counter for each new item
              section: {
                name: item.title,
                prompt: generatedItem.section_prompt,
                result: generatedItem.section_results[0],
              },
              reference_documents: {
                id: documentId,
                name:
                  referenceDocumentNames[documentId] || "Reference Document " + documentId + ".pdf",
                // dummy
                // name: "Reference Document " + documentId + ".pdf",
                prompt: generatedItem.reference_prompt,
                result: generatedItem.reference_results[0],
              },
              ai_review: removeCodeBlocks(generatedItem.ai_review),
              ai_suggestion: removeCodeBlocks(generatedItem.ai_suggestion),
            };
          });
        });

        datas = await Promise.all(promise);

        // let alreadyAdded = new Set();

        // let promise2 = datas
        //   .filter((data) => {
        //     if (alreadyAdded.has(data.section.name)) {
        //       return false;
        //     } else {
        //       alreadyAdded.add(data.section.name);
        //       return true;
        //     }
        //   })
        //   .map((data) => {
        //     return {
        //       name: data.section.name,
        //       result: data.ai_review,
        //     };
        //   });

        // datas2 = await Promise.all(promise2);

        let promise2 = response.result.data[0].attributes.result.map((item) => {
          return {
            name: item.title,
            result: removeCodeBlocks(item.generated.final_result),
          };
        });

        datas2 = await Promise.all(promise2);

        triggerLoading(false);
      } catch (error) {
        console.log(error);
      }
    }
  }

  const removeCodeBlocks = (text) => {
    return text
      .replace(/```/g, "") // Hapus semua backticks triple
      .replace(/“/g, "") // Hapus semua tanda kutip buka
      .replace(/”/g, "") // Hapus semua tanda kutip tutup
      .replace(/\s{3,}/g, "  ") // Gantikan lebih dari tiga spasi dengan dua spasi
      .trim();
  };

  async function updateResult() {
    editModal = false;
    triggerLoading(true);
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let index = tmpResult.findIndex((data) => data.title === tmpTitle);
      if (index !== -1) {
        tmpResult[index].text = tmpText;
      }

      try {
        let response = await $call.call("client-dashboard", "update-contract-project-result", {
          URLQuery: {},
          URLParams: {
            id: projectDocumentId,
          },
          Params: {
            result: JSON.stringify(tmpResult),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          editModal = true;
          triggerToast("error");
          console.log(response);
          return;
        }

        editModal = false;
        await getProjectResult();

        triggerToast("updated");
        triggerLoading(false);
      } catch (error) {
        console.log(error);
      }
    }
  }

  // function downloadFile() {
  //   // const element = document.getElementById("content-to-download");

  //   // Menggunakan skala untuk memastikan gambar pas di halaman PDF
  //   // const scale = 2; // Nilai ini bisa diatur sesuai kebutuhan
  //   // const options = {
  //   //   width: element.clientWidth * scale,
  //   //   height: element.clientHeight * scale,
  //   //   style: {
  //   //     transform: `scale(${scale})`,
  //   //     transformOrigin: "top left",
  //   //     width: `${element.clientWidth}px`,
  //   //     height: `${element.clientHeight}px`,
  //   //   },
  //   // };

  //   // domtoimage
  //   //   .toPng(element, options)
  //   //   .then((dataUrl) => {
  //   //     const pdf = new jsPDF({
  //   //       orientation: "portrait",
  //   //       unit: "px",
  //   //       format: [options.width / scale, options.height / scale],
  //   //     });
  //   //     pdf.addImage(dataUrl, "PNG", 0, 0, options.width / scale, options.height / scale);
  //   //     pdf.save("output.pdf");
  //   //   })
  //   //   .catch((error) => {
  //   //     console.error("Oops, something went wrong!", error);
  //   //   });

  //   const element = document.getElementById("content-to-download");
  //   const pageHeight = 295; // Tinggi halaman dalam mm (A4)
  //   const pageWidth = 210; // Lebar halaman dalam mm (A4)
  //   const margin = 10; // Margin halaman dalam mm
  //   const scale = 2; // Skala untuk kualitas gambar

  //   html2canvas(element, {
  //     scale: scale,
  //     logging: true,
  //   })
  //     .then((canvas) => {
  //       const imgData = canvas.toDataURL("image/png");
  //       const imgWidth = pageWidth - 2 * margin;
  //       const imgHeight = (canvas.height * imgWidth) / canvas.width;
  //       const pdf = new jsPDF({
  //         orientation: "portrait",
  //         unit: "mm",
  //         format: "a4",
  //       });

  //       let position = 0;
  //       const pageCount = Math.ceil(imgHeight / (pageHeight - 2 * margin));

  //       for (let i = 0; i < pageCount; i++) {
  //         const yOffset = -position * (pageHeight - 2 * margin);
  //         pdf.addImage(imgData, "PNG", margin, yOffset + margin, imgWidth, imgHeight - 2 * margin);
  //         if (i < pageCount - 1) {
  //           pdf.addPage();
  //         }
  //         position++;
  //       }

  //       pdf.save(projectName + ".pdf");
  //     })
  //     .catch((error) => {
  //       console.error("Oops, something went wrong!", error);
  //     });
  // }

  onMount(() => {
    const unsubscribe = page.subscribe(async ($page) => {
      projectDocumentId = $page.url.searchParams.get("projectDocumentId");
      if (projectDocumentId !== null) {
        getProjectResult();
      } else {
        triggerToast("error");
        goto(await $call.getPage("client-dashboard", "contract-project-list"));
      }
    });
    return () => unsubscribe();
  });

  $effect(async () => {
    const unsubscribe = page.subscribe(async ($page) => {
      projectDocumentId = $page.url.searchParams.get("projectDocumentId");
      if (projectDocumentId !== null) {
        getProjectResult();
      } else {
        goto(await $call.getPage("client-dashboard", "contract-project-list"));
      }
    });
    onDestroy(() => unsubscribe());
  });
</script>

<div class="flex flex-col px-4 pt-6 gap-6 pb-12">
  <!-- section header -->
  <div class="flex justify-between">
    <div class="flex items-center gap-3">
      <div class="flex flex-col gap-1.5">
        <P class="text-3xl text-primary-light dark:text-primary-dark font-semibold"
          >{$_("aiReview")}</P
        >
        <Typography preset="Text sm/Regular" text={projectName} color="tertiary" />
        <!-- <Button
          onclick={() => {
            downloadFile();
          }}
          class="h-max"><DownloadIcon svgClass="w-5 h-5 text-white" />Download</Button
        > -->
      </div>
    </div>
    <div class="flex items-center gap-3">
      <!-- <Button outline on:click={() => {}} color="light" class="h-max">
        <FileDownloadIcon svgClass="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark" />
        &nbsp;
        {$_("seeReferencedDoc")}
      </Button> -->
      <!-- <Button outline on:click={() => {}} color="light" class="h-max">
        <FileDownloadIcon svgClass="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark" />
        &nbsp;
        {$_("generateNewDocs")}
      </Button> -->
      <!-- <Button
        on:click={() => {
          downloadFile();
        }}
        class="h-max"
      >
        <DownloadIcon svgClass="w-5 h-5 text-white" />
        &nbsp;
        {$_("download")}</Button
      > -->
    </div>
  </div>
  <Divider />
  <div class="flex gap-3">
    <div class="flex flex-col grow gap-6">
      <!-- <div onfocusout={handleDropdownFocusLoss} class="flex flex-col gap-1 w-full"> -->
      <div class="flex flex-col gap-6 w-full">
        <!-- document section -->
        <!-- <Button
          on:click={() => (dropdownToggle[0] = !dropdownToggle[0])}
          class="flex justify-between border-primary-light dark:border-primary-dark bg-secondary-light dark:bg-secondary-dark"
          color="light"
        >
          <div class="flex items-center gap-3">
            <P class="text-lg font-semibold text-tertiary-light dark:text-tertiary-dark"
              >{$_("generatedDocument")}</P
            >
            <Button outline on:click={() => (documentViewModal = true)} color="light" class="h-max">
              <FileDownloadIcon svgClass="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark" />
              &nbsp;
              {$_("seeReferencedDoc")}
            </Button>
          </div>
          <ChevronDownOutline
            class="w-8 h-8 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
          />
        </Button>
        {#if dropdownToggle[0]}
          <div class="flex justify-center mt-6">
            <div class="bg-tertiary-light dark:bg-tertiary-dark rounded-lg p-4">
              <PdfViewer
                showBorder={false}
                showButtons={["navigation", "zoom", "rotate", "autoflip", "pageInfo"]}
                scale={1}
                data={atob(base64Data)}
              />
            </div>
          </div>
        {/if}

        <Divider /> -->

        <!-- first table section -->
        <!-- <Button
          on:click={() => (dropdownToggle[1] = !dropdownToggle[1])}
          class="flex justify-between border-primary-light dark:border-primary-dark bg-secondary-light dark:bg-secondary-dark"
          color="light"
        >
          <P class="text-lg font-semibold text-tertiary-light dark:text-tertiary-dark"
            >{$_("summaryOfFindings")}</P
          >
          <ChevronDownOutline
            class="w-8 h-8 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
          />
        </Button> -->
        {#if dropdownToggle[0]}
          <div
            class="flex flex-col justify-center mt-6 bg-tertiary-light dark:bg-tertiary-dark rounded-lg p-4"
          >
            <ListMarkdown id="content-to-download" bind:datas={datas2} />
          </div>
        {/if}
        {#if dropdownToggle[1]}
          <DataTableAiReview bind:datas />
          <!-- <DataTable bind:datas /> -->
        {/if}
      </div>
    </div>
  </div>
</div>

<Modal title={$_("documentViewer")} bind:open={documentViewModal} autoclose outsideclose>
  {#if selectedFile && (selectedFile.name.includes(".pdf") || selectedFile.name.includes(".PDF"))}
    <div class="flex justify-center">
      <PdfViewer
        showBorder={false}
        showButtons={["navigation", "zoom", "rotate", "autoflip", "pageInfo"]}
        scale={1}
        data={atob(selectedFile.base64)}
      />
    </div>
  {/if}
</Modal>

<Modal title={$_("editResult")} bind:open={editModal} size="lg">
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={tmpTitle} />
      <Input bind:value={tmpText} placeholder={$_("whatIsYourChangeResult")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={() => {
        updateResult();
      }}
      class="w-full"
      disabled={buttonDisabled}>{$_("change")}</Button
    >
  </svelte:fragment>
</Modal>
