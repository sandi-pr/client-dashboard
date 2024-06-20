<script>
  import { onMount } from "svelte";
  import { writable } from "svelte/store";
  import {
    P,
    Span,
    Modal,
    Avatar,
    Button,
    Select,
    Badge,
    CloseButton,
  } from "flowbite-svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  // import PdfViewer from "svelte-pdf";
  // import PlusIcon from "$lib/icon/plus.svg";
  import { _ } from "svelte-i18n";
  import XCloseIcon from "$lib/icon/x-close.svg";
  import QRCodeIcon from "$lib/image/qr_code.png";
  import Timeline from "$lib/components/Timeline.svelte";
  import { timelineStep } from "$lib/utils/timelineStep.js";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import UsersPlusIcon from "$lib/icon/users-plus.svelte";
  import UserIcon from "$lib/icon/user-03.svelte";
  import UserCheckIcon from "$lib/icon/user-check-02.svelte";
  import UsersUpIcon from "$lib/icon/users-up.svelte";
  import DotsGridIcon from "$lib/icon/dots-grid.svelte";
  import { dndzone, TRIGGERS } from "svelte-dnd-action";
  // import PDFViewer from "pdf.svelte";
  // import {
  //   PUBLIC_API_URL_DEV,
  //   PUBLIC_API_URL_PROD,
  //   PUBLIC_STRAPI_URL_DEV,
  //   PUBLIC_STRAPI_URL_PROD,
  // } from "$env/static/public";
  // import axios from "axios";
  // import { date } from "svelte-i18n";

  // const API_URL = process.env.NODE_ENV === "development" ? PUBLIC_API_URL_DEV : PUBLIC_API_URL_PROD;
  // const STRAPI_URL =
  //   process.env.NODE_ENV === "development" ? PUBLIC_STRAPI_URL_DEV : PUBLIC_STRAPI_URL_PROD;

  // Step 1
  let currentStep = 1;
  const steps = ["1", "2", "3"];

  // Step 2
  let files = [];
  // $: console.log(files);
  let pdf = null;
  let base64 = null;
  let fileUrl = null;

  // Create writable stores for QR code position
  const qrCodeX = writable(0);
  const qrCodeY = writable(0);
  // Variables to store initial mouse position and movement state
  // let initialMouseX, initialMouseY;
  // let isMoving = false;

  let dropdownOpen = false;
  let formModal = false;
  let addSignerModal = false;
  let selectRoleModal = false;
  let attributes = [];
  let name = "";
  let value = "";
  let personClient = [];
  let selectedPersonClient = "";
  let selectedRole = "";

  let isDragging = false;
  let imagePosition = { x: 0, y: 0 };

  // // Function to handle QR code movement
  // function handleQrCodeMove(event) {
  //   if (isMoving) {
  //     const deltaX = event.clientX - initialMouseX;
  //     const deltaY = event.clientY - initialMouseY;

  //     qrCodeX.update((x) => x + deltaX);
  //     qrCodeY.update((y) => y + deltaY);

  //     initialMouseX = event.clientX;
  //     initialMouseY = event.clientY;
  //   }
  // }

  // // Function to start QR code movement
  // function startMove(event) {
  //   console.log("on mouse down");
  //   console.log(event.clientX, event.clientY);
  //   isMoving = true;
  //   initialMouseX = event.clientX;
  //   console.log(initialMouseX);
  //   initialMouseY = event.clientY;
  //   console.log(initialMouseY);
  // }

  // // Function to stop QR code movement
  // function stopMove() {
  //   isMoving = false;
  //   saveQrCodePosition();
  // }

  // Function to save QR code position
  function saveQrCodePosition() {
    const position = { x: $qrCodeX, y: $qrCodeY };
    console.log("QR Code Position:", position);
    // You can save the position to a database or use it as needed
  }

  function removeAttribute(index) {
    attributes = attributes.filter((_, i) => i !== index);
  }

  // Subscribe to timelineStep changes
  const unsubscribe = timelineStep.subscribe((value) => {
    currentStep = value;
  });

  // Unsubscribe on component destruction
  onMount(() => {
    return () => {
      unsubscribe();
      files = [];
    };
  });

  // Step 1
  // Current Workflow
  // let currentWorkflow_Step1 = 0;

  // Step 1
  // Change Workflow
  // function changeWorkflow_Step1(newValue) {
  //   currentWorkflow_Step1 = newValue;
  //   // console.log(currentWorkflow_Step1);
  // }

  // Step 2
  // Convert file size
  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  // export let data;
  let data = {
    workflows: {
      data: [
        {
          id: 1,
          attributes: {
            name: "Workflow 1",
            description: "Description 1",
            unit: {
              data: {
                attributes: {
                  name: "Unit 1",
                },
              },
            },
          },
        },
        {
          id: 2,
          attributes: {
            name: "Workflow 2",
            description: "Description 2",
            unit: {
              data: {
                attributes: {
                  name: "Unit 2",
                },
              },
            },
          },
        },
        {
          id: 3,
          attributes: {
            name: "Workflow 3",
            description: "Description 3",
            unit: {
              data: {
                attributes: {
                  name: "Unit 3",
                },
              },
            },
          },
        },
      ],
    },
  };

  let workflows = data.workflows.data ?? [];

  let selectedWorkflow = null;
  let workflowSteps = null;
  let selectedSigner;
  let dropdownSigner = false;
  let signers = [];
  let showQr = [];
  let useMeterai = false;
  let emeterai = [];
  let fileTextName = null;
  let fileTextResult = null;
  let base64Data = "";

  // async function selectWorkflow(item) {
  //   selectedWorkflow = item;
  //   const response = await axios.get(`${API_URL}/api/workflows/${item.id}`, {
  //     headers: {
  //       Authorization: "Bearer " + data.data.jwt,
  //     },
  //   });
  //   workflowSteps = response.data.data.attributes.steps;
  //   selectedSigner = workflowSteps[0];
  //   const meterai = response.data.data.attributes.meterai_type;
  //   if (meterai === "Meterai") {
  //     useMeterai = true;
  //   }
  //   for (const step of workflowSteps) {
  //     const newSigner = {
  //       page: 1,
  //       address: step.identity.data.attributes.address ?? "",
  //       email: step.user.email,
  //       position: { x: 0, y: 0 },
  //     };
  //     signers.push(newSigner);
  //     showQr.push(false);
  //     if (meterai === "Meterai") {
  //       emeterai.push(false);
  //     }
  //   }
  //   console.log(workflowSteps);
  // }

  // Step 2
  // function handleFileChange(event) {
  //   const selectedFiles = event.target.files;

  //   for (const file of selectedFiles) {
  //     if (file.size > 1 * 1024 * 1024) {
  //       files = [];
  //       alert(`File "${file.name}" size exceeds the 10 MB limit. Please choose a smaller file.`);
  //       return;
  //     }
  //     if (file) {
  //       const reader = new FileReader();
  //       reader.onloadend = () => {
  //         pdf = reader.result;
  //         base64 = reader.result.split(",")[1];
  //       };
  //       reader.readAsDataURL(file);
  //     }
  //   }
  // }

  function removeFile(index) {
    // files = files.filter((_, i) => i !== index);
    let arrayFileList = [...files];
    files = arrayFileList.filter((_, i) => i !== index);
  }

  function handleSubmit(event) {
    event.preventDefault();
    addAttribute();
    formModal = false;
  }

  function addAttribute() {
    if (name && value) {
      attributes = [...attributes, { name, value }];
      name = "";
      value = "";
    }
  }

  function selectSigner(signer) {
    selectedSigner = signer;
    dropdownSigner = false;
  }

  const handleMouseDown = (e) => {
    isDragging = true;

    const offsetX = e.clientX - imagePosition.x;
    const offsetY = e.clientY - imagePosition.y;

    const handleMouseMove = (e) => {
      if (isDragging) {
        const newLeft = e.clientX - offsetX;
        const newTop = e.clientY - offsetY;

        imagePosition = { x: newLeft, y: newTop };
      }
    };

    const handleMouseUp = () => {
      const existingSignerIndex = signers.findIndex(
        (signer) => signer.email === selectedSigner.user.email
      );

      if (existingSignerIndex !== -1) {
        signers[existingSignerIndex].position = imagePosition;
        console.log("ini: ", signers[existingSignerIndex]);
      }

      isDragging = false;
      window.removeEventListener("mousemove", handleMouseMove);
      window.removeEventListener("mouseup", handleMouseUp);
      document.body.style.cursor = "auto";
    };

    window.addEventListener("mousemove", handleMouseMove);
    window.addEventListener("mouseup", handleMouseUp);
    document.body.style.cursor = "grabbing";
  };

  function placeSignature() {
    const existingSignerIndex = signers.findIndex(
      (signer) => signer.email === selectedSigner.user.email
    );
    showQr[existingSignerIndex] = true;
  }

  function placeMeterai() {
    const existingSignerIndex = signers.findIndex(
      (signer) => signer.email === selectedSigner.user.email
    );
    emeterai[existingSignerIndex] = true;
  }

  const validateCurrentStep = () => {
    if (currentStep === 1 && selectedWorkflow === null) {
      return false;
    }
    if (currentStep === 2 && pdf === null) {
      return false;
    }
    if (currentStep === 4 && !showQr.every((value) => value === true)) {
      return false;
    }
    return true;
  };

  async function readFileAsText(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        const text = reader.result.replace(/\r\n/g, "\n");
        const base64 = btoa(text);
        resolve({ text, base64 });
      };
      reader.onerror = reject;

      reader.readAsText(file);
    });
  }

  async function readFileAsBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.readAsDataURL(file);
      reader.onload = () => {
        const base64Data = reader.result.split(",")[1];
        resolve(base64Data);
      };
      reader.onerror = reject;

    });
  }

  const base64toBlob = (base64Data, contentType) => {
    const byteCharacters = atob(base64Data);
    const byteNumbers = new Array(byteCharacters.length);
    for (let i = 0; i < byteCharacters.length; i++) {
      byteNumbers[i] = byteCharacters.charCodeAt(i);
    }
    const byteArray = new Uint8Array(byteNumbers);
    return new Blob([byteArray], { type: contentType });
  };

  const handleFileUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    // inputElement.accept = ".txt";
    inputElement.accept = ".pdf, .doc, .docx";

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      files = event.target.files;

      if (file) {
        base64Data = await readFileAsBase64(file);

        const blob = base64toBlob(base64Data, file.type);
        fileUrl = URL.createObjectURL(blob);

        console.log(fileUrl);

        fileTextName = file.name;
      } else {
        alert("Invalid file. Please select a valid file.");
      }
    });

    inputElement.click();
  };

  function handleDndConsider(e) {
    const {
      items: newItems,
      info: { trigger },
    } = e.detail;
    if (trigger === TRIGGERS.DRAG_STARTED) {
      // Optional: Handle logic when drag starts
    }
    signers = newItems;
  }

  function handleDndFinalize(e) {
    const { items: newItems } = e.detail;
    signers = newItems;
  }

  function selectRole(role) {
    selectRoleModal = false;
    selectedRole = role;
    addSignerModal = true;
  }

  function addSigner() {
    const newSigner = {
      id: signers.length + 1,
      name: selectedPersonClient,
      role: selectedRole,
    };
    signers = [...signers, newSigner];

    selectedPersonClient = "";
    selectedRole = "";
  }

  function deleteSigner(id) {
    signers = signers.filter((signer) => signer.id !== id);
  }

  async function getPersons() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let page = 1;
      while (true) {
        let response = await call.call("client-dashboard", "get-persons", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 25,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        if (response.result.data.length > 0) {
          let promise = response.result.data.map((person) => {
            return {
              value: person.id,
              name:
                person.attributes.fullname !== null
                  ? person.attributes.fullname
                  : `<${person.attributes.email}>`,
            };
          });

          let data = await Promise.all(promise);
          personClient = personClient.concat(data);
          page++;
        } else {
          break;
        }
      }
    }
  }

  // const endStep = () => {
  //   let signerData = [];
  //   for (const signer of workflowSteps) {
  //   }
  //   const data = {
  //     organization: data.id_organization,
  //     attributes: attributes,
  //     name: files[0].name,
  //     description: "",
  //     user_address: data.identity.address,
  //   };
  // };

  onMount(async () => {
    // getPersons();
  });
</script>

<Timeline {currentStep} {steps} {validateCurrentStep}></Timeline>
<div class="mt-16">

  <!-- Step 1 -->
  {#if currentStep == 1}
    <div class="flex flex-col justify-center lg:px-72 md:px-10 sm:px-4 gap-4">
      <div class="flex justify-between items-start gap-3">
        <P size="3xl" weight="semibold">{$_("uploadContract")}</P>
      </div>
      <Button
        onclick={() => {
          handleFileUploadWrapper();
        }}
        color="light"
        class={`w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark gap-1 py-8 flex flex-col justify-start ${fileTextName ? "items-start cursor-default" : "items-center"}`}
      >
        {#if fileTextName}
          {#if fileTextName.endsWith(".txt")}
            <!-- {#each fileTextResult.split("\n") as paragraphs, i}
              <P>{paragraphs}</P>
            {/each} -->
            <P>{fileTextName}</P>
          {:else if fileTextName.endsWith(".pdf")}
            <P>{fileTextName}</P>
          {:else if fileTextName.endsWith(".doc") || fileTextName.endsWith(".docx")}
            <P>{fileTextName}</P>
          {:else}
            <P>{fileTextName}</P>
          {/if}
        {:else}
          <div
            class="rounded-xl shadow border border-primary-light dark:border-primary-dark justify-center items-center inline-flex mb-4 pointer-events-none"
          >
            <img src={upload} alt="mail.svg" class="p-4 dark:invert pointer-events-none" />
          </div>
          <div class="flex space-x-1">
            <P class="text-hc-brand-700 pointer-events-none">{$_("clickToUpload")}</P>
            <P class="pointer-events-none">
              <!-- {$_("orDragAndDrop")} -->
            </P>
          </div>
          <!-- <P class="pointer-events-none">SVG, PNG, JPG, or GIF (max. 800x400px)</P> -->
          <P class="pointer-events-none">PDF, DOC, DOCX (Max. 15MB)</P>
        {/if}
      </Button>

      <!-- {#if files && files.length > 0}
        <div class="flex flex-col gap-2">
          {#each files as file, index (file.name)}
            <div
              class="flex justify-between items-center bg-white px-3 py-3 rounded-lg border-2 border-hc-graylight-200"
            >
              <div class="flex justify-start items-center gap-3">
                <enhanced:img src="$lib/icon/pdf-icon.svg" class="w-10 h-10" alt="An alt text" />
                <div class="flex flex-col justify-start items-start truncate">
                  <p class="text-hc-graylight-600 text-sm font-semibold dark:text-gray-400">
                    {file.name}
                  </p>
                  <P
                    size="sm"
                    weight="light"
                    class="text-hc-graylight-600 dark:text-gray-400 text-sm"
                    >{formatBytes(file.size)}</P
                  >
                </div>
              </div>
              <div class="flex justify-center items-center">
                <CheckCircleOutline class="w-5 h-5 ml-3 text-hc-brand-700" />
                <TrashBinOutline
                  class="w-5 h-5 ml-3 text-hc-error-700 cursor-pointer"
                  on:click={() => removeFile(index)}
                />
              </div>
            </div>
          {/each}
        </div>
      {/if} -->
    </div>
  {/if}

  <!-- Step 2 -->
  {#if currentStep == 2}
  <div class="flex flex-col justify-center lg:px-72 md:px-10 sm:px-4 gap-4">
    <div class="flex justify-between items-start gap-3">
      <P size="3xl" weight="semibold">{$_("assignSigner")}</P>
      <Button outline on:click={() => (selectRoleModal = true)} class="gap-1 text-utility-brand-700-light dark:text-utility-brand-700-dark hover:!bg-transparent hover:text-utility-brand-700-light dark:hover:text-utility-brand-700-dark">
        <UsersPlusIcon svgClass="w-5 h-5 text-utility-brand-700-light dark:text-utility-brand-700-dark" />
        {$_("invite")}
      </Button>
    </div>
    <div class="flex flex-col grow gap-6">
      {#if signers.length > 0}
        <div
          class="bg-secondary-light dark:bg-secondary-dark rounded-lg py-10 px-4 border border-primary-light dark:border-primary-dark"
        >
          <div
            use:dndzone={{ items: signers, flipDurationMs: 150 }}
            on:consider={handleDndConsider}
            on:finalize={handleDndFinalize}
            class="flex flex-col gap-4"
          >
            {#each signers as signer (signer.id)}
              <div
                class="px-6 py-4 bg-white dark:bg-tertiary-dark rounded-lg backdrop-blur-[26px] justify-between items-center flex w-full"
              >
                <div
                  class="w-full self-stretch items-center justify-between flex gap-4"
                >
                  <div class="flex items-center gap-4">
                    <DotsGridIcon
                      svgClass="w-6 h-6 text-fg-quaternary-light dark:text-fg-quaternary-dark"
                    />
                    <Avatar border />
                    <P size="lg" weight="semibold">{signer.name}</P>
                  </div>
                  <Badge border rounded color={signer.role == "signer" ? "primary" : "indigo"}>{$_(signer.role)}</Badge>
                </div>
                <CloseButton on:click={() => deleteSigner(signer.id)} class="text-fg-primary-light dark:text-fg-primary-dark" />
              </div>
            {/each}
          </div>
        </div>
      {:else}
        <div class="flex grow h-60 items-center justify-center p-4 rounded-lg border border-primary-light dark:border-primary-dark">
          <P size="md" weight="semibold" class="text-gray-500">{$_("noOneInvitedYet")}</P>
        </div>
      {/if}
    </div>
  </div>
  <Modal title={$_("whoDoYouWantToInvite")} bind:open={selectRoleModal} autoclose outsideclose>
    <div class="flex flex-col grow gap-3">
      <Button on:click={() => selectRole("signer")} color="light" class="flex grow p-6 rounded-lg border border-secondary-light dark:border-secondary-dark justify-start items-center gap-4">
        <UserIcon svgClass="w-7 h-7 text-utility-gray-500-light dark:text-utility-gray-500-dark" />
        <P weight="semibold" class="text-primary-light dark:text-primary-dark">{$_("signer")}</P>
      </Button>
    </div>
    <div class="flex flex-col grow gap-3">
      <Button on:click={() => selectRole("reviewer")} color="light" class="flex grow p-6 rounded-lg border border-secondary-light dark:border-secondary-dark justify-start items-center gap-4">
        <UserCheckIcon svgClass="w-7 h-7 text-utility-gray-500-light dark:text-utility-gray-500-dark" />
        <P weight="semibold" class="text-primary-light dark:text-primary-dark">{$_("reviewer")}</P>
      </Button>
    </div>
    <div class="flex flex-col grow gap-3">
      <Button on:click={() => selectRole("thirdPartySigner")} color="light" class="flex grow p-6 rounded-lg border border-secondary-light dark:border-secondary-dark justify-start items-center gap-4">
        <UsersUpIcon svgClass="w-7 h-7 text-utility-gray-500-light dark:text-utility-gray-500-dark" />
        <P weight="semibold" class="text-primary-light dark:text-primary-dark">{$_("thirdPartySigner")}</P>
      </Button>
    </div>
  </Modal>
  <Modal title={$_("inviteSigner")} bind:open={addSignerModal} autoclose outsideclose>
    <Select
      placeholder={$_("selectPerson")}
      class="bg-primary-light dark:bg-secondary-dark"
      items={personClient}
      bind:value={selectedPersonClient}
    />
    <svelte:fragment slot="footer">
      <Button disabled={!selectedPersonClient} on:click={addSigner} class="w-full">{$_("addSigner")}</Button>
    </svelte:fragment>
  </Modal>
  {/if}

  <!-- Step 3 -->
  {#if currentStep == 3}
    {#if files}
      <div class="px-20">
        <div class="flex justify-between items-center pb-4">
          <P class="text-3xl leading-9 font-semibold text-hc-graylight-900 dark:text-hc-graydark-50"
            >Signature Placement</P
          >
          <div class="flex items-center gap-2">
            <Button color="light">
              {files[0].name}
              <ChevronDownOutline class="w-3 h-3 ms-2 dark:text-white" />
            </Button>
            <!-- <Dropdown bind:open={dropdownOpen} class="w-48 overflow-y-auto py-1">
              <DropdownItem on:click={() => (dropdownOpen = false)}>
                {files[0].name}
              </DropdownItem>
            </Dropdown> -->
          </div>
        </div>
        <div class="flex gap-12">
          <div class="flex flex-col gap-4">
            <P class="font-semibold text-hc-graylight-900 dark:text-hc-graydark-50"
              >Choose which signer to place the sign</P
            >
            <Button outline color="light" class="focus-within:!ring-1">
              <Avatar />
              <div class="flex-col grow pl-2 pr-4">
                <Span
                  class="font-jakarta text-slate-700 dark:text-hc-graydark-300 text-sm leading-tight"
                  >{selectedSigner.user.fullname}</Span
                >
                <P
                  class="font-jakarta text-slate-600 dark:text-hc-graydark-400 text-sm leading-tight"
                  >{selectedSigner.user.unit.name}</P
                >
              </div>
              <Button
                pill
                class={`
                  px-2 py-0 bg-hc-warning-50 dark:bg-hc-warning-950 text-hc-warning-700 dark:text-hc-warning-300 border border-hc-warning-200 dark:border-hc-warning-800 hover:bg-hc-warning-50 dark:hover:bg-hc-warning-950
                  rounded-md focus-within:ring-0 font-medium
                  `}
              >
                {selectedSigner.step}
              </Button>
              <ChevronDownOutline class="w-3 h-3 ms-2 dark:text-white" />
            </Button>
            <!-- <Dropdown bind:open={dropdownSigner}>
              {#each workflowSteps as signer}
                <DropdownItem>
                  <Button
                    on:click={selectSigner(signer)}
                    outline
                    color="light"
                    class="focus-within:!ring-1"
                  >
                    <Avatar />
                    <div class="flex-col grow pl-2 pr-4">
                      <Span
                        class="font-jakarta text-slate-700 dark:text-hc-graydark-300 text-sm leading-tight"
                        >{signer.user.fullname}</Span
                      >
                      <P
                        class="font-jakarta text-slate-600 dark:text-hc-graydark-400 text-sm leading-tight"
                        >{signer.user.unit.name}</P
                      >
                    </div>
                    <Button
                      pill
                      class={`
                        px-2 py-0 bg-hc-warning-50 dark:bg-hc-warning-950 text-hc-warning-700 dark:text-hc-warning-300 border border-hc-warning-200 dark:border-hc-warning-800 hover:bg-hc-warning-50 dark:hover:bg-hc-warning-950
                        rounded-md focus-within:ring-0 font-medium
                        `}
                    >
                      {signer.step}
                    </Button>
                    <ChevronDownOutline class="w-3 h-3 ms-2 dark:text-white" />
                  </Button>
                </DropdownItem>
              {/each}
            </Dropdown> -->
            <div class="flex gap-4 justify-center items-center">
              {#if useMeterai}
                <Button
                  outline
                  on:click={placeMeterai}
                  color="purple"
                  class="w-full hover:!bg-transparent hover:text-hc-brand-600 !border-hc-brand-300 dark:!border-hc-brand-400 focus-within:!ring-1 px-3 py-1 gap-1"
                >
                  Place e-Meterai
                </Button>
              {/if}
              <Button
                on:click={placeSignature}
                class="w-full px-3 !py-1 gap-1 focus-within:!ring-1"
              >
                Place Signature
              </Button>
            </div>
          </div>
          <div class="flex justify-center grow py-10 bg-gray-50 rounded-lg dark:bg-gray-800">
            <!-- <PDFViewer {pdf} /> -->
            <!-- <PdfViewer showBorder={false} showButtons={false} scale={1} data={atob(base64)} /> -->
            <div class="relative">
              <!-- <div
                class="absolute cursor-move"
                role="button"
                aria-label="Move QR Code"
                tabindex="0"
                style={`left: ${qrCodeX}px; top: ${qrCodeY}px;`}
                on:mousedown={startMove}
                on:mousemove={handleQrCodeMove}
                on:mouseup={stopMove}
              >
                <img src={QRCodeIcon} class="w-20 h-20" alt="QR Code" />
              </div> -->
              {#each signers as signer, index (signer)}
                {#if showQr[index]}
                  <div
                    on:mousedown={(e) =>
                      selectedSigner.user.email === signer.email ? handleMouseDown(e, signer) : ""}
                    role="button"
                    aria-label="Move QR Code"
                    tabindex="0"
                  >
                    <img
                      src={QRCodeIcon}
                      alt=""
                      width="70"
                      height="70"
                      style="position: absolute; top: {signer.position.y}px; left: {signer.position
                        .x}px; z-index: 2; cursor: {isDragging ? 'grabbing' : 'grab'}"
                    />
                  </div>
                {/if}
                {#if emeterai[index]}
                  {#if showQr[index]}
                    <div
                      on:mousedown={(e) =>
                        selectedSigner.user.email === signer.email
                          ? handleMouseDown(e, signer)
                          : ""}
                      role="button"
                      aria-label="Move QR Code"
                      tabindex="0"
                    >
                      <img
                        src={QRCodeIcon}
                        alt=""
                        width="70"
                        height="70"
                        style="position: absolute; top: {signer.position.y}px; left: {signer
                          .position.x}px; z-index: 2; cursor: {isDragging ? 'grabbing' : 'grab'}"
                      />
                    </div>
                  {/if}
                {/if}
              {/each}

              <div class="flex justify-center grow bg-gray-50 rounded-lg dark:bg-gray-800">
                <div>pdf viewer here</div>
                <!-- <PdfViewer
                  showBorder={false}
                  showButtons={false}
                  scale={1}
                  data={atob(base64)}
                  on:mouseleave={saveQrCodePosition}
                /> -->
              </div>
            </div>
          </div>
        </div>
      </div>
    {/if}
  {/if}
</div>
