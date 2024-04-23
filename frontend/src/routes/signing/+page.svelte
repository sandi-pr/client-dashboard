<script>
  import { onMount } from "svelte";
  import { writable } from "svelte/store";
  import {
    P,
    Span,
    Card,
    Modal,
    Label,
    Badge,
    Input,
    // Search,
    Avatar,
    Button,
    // Checkbox,
    // Dropdown,
    // DropdownItem,
  } from "flowbite-svelte";
  import {
    SearchOutline,
    UserRemoveSolid,
    FilterOutline,
    ChevronDownOutline,
    CheckCircleOutline,
    TrashBinOutline,
  } from "flowbite-svelte-icons";
  // import PdfViewer from "svelte-pdf";
  // import PlusIcon from "$lib/icon/plus.svg";
  import XCloseIcon from "../lib/icon/x-close.svg";
  import QRCodeIcon from "../lib/image/qr_code.png";
  import Timeline from "../lib/components/Timeline.svelte";
  import { timelineStep } from "../lib/utils/timelineStep.js";
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
  const steps = ["1", "2", "3", "4"];

  // Step 2
  let files = [];
  // $: console.log(files);
  let pdf = null;
  let base64 = null;

  // Create writable stores for QR code position
  const qrCodeX = writable(0);
  const qrCodeY = writable(0);
  // Variables to store initial mouse position and movement state
  // let initialMouseX, initialMouseY;
  // let isMoving = false;

  let dropdownOpen = false;
  let formModal = false;
  let attributes = [];
  let name = "";
  let value = "";

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

  onMount(() => {
    console.log(data);
  });

  let selectedWorkflow = null;
  let workflowSteps = null;
  let selectedSigner;
  let dropdownSigner = false;
  let signers = [];
  let showQr = [];
  let useMeterai = false;
  let emeterai = [];

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
</script>

<Timeline {currentStep} {steps} {validateCurrentStep}></Timeline>
<div class="mt-16">
  <!-- Step 1 -->
  {#if currentStep == 1}
    <div class="flex flex-col justify-center lg:px-72 md:px-10 sm:px-4 gap-4">
      <div class="flex justify-between items-start gap-3">
        <P size="3xl" weight="semibold">Choose workflow</P>

        <Button color="light" class="text-hc-graylight-700">
          <FilterOutline class="w-3 h-3 ms-2 text-hc-graylight-700 dark:text-white" />Filter</Button
        >
        <!-- <Dropdown class="overflow-y-auto px-3 pb-3 text-sm h-44">
          <div slot="header" class="p-3">
            <Search size="md" />
          </div>
          <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">
            <Checkbox>Robert Gouth</Checkbox>
          </li>
          <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">
            <Checkbox>Jese Leos</Checkbox>
          </li>
          <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">
            <Checkbox checked>Bonnie Green</Checkbox>
          </li>
          <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">
            <Checkbox>Jese Leos</Checkbox>
          </li>
          <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">
            <Checkbox>Robert Gouth</Checkbox>
          </li>
          <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">
            <Checkbox>Bonnie Green</Checkbox>
          </li>
          <a
            slot="footer"
            href="/"
            class="flex items-center p-3 -mb-1 text-sm font-medium text-red-600 bg-gray-50 hover:bg-gray-100 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-red-500 hover:underline"
          >
            <UserRemoveSolid class="w-4 h-4 me-2 text-primary-700 dark:text-primary-700" />Delete
            user
          </a>
        </Dropdown> -->
      </div>
      <Input
        id="email"
        size="lg"
        color="base"
        class="focus:border-hc-brand-600 focus:ring-hc-brand-600"
        type="email"
        placeholder="name@flowbite.com"
      >
        <SearchOutline slot="left" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
      </Input>
      <div class="flex flex-col justify-center mt-3 gap-3">
        <div class="flex text-start">
          <P class="text-hc-graylight-600 text-base font-normal">Recent workflow</P>
        </div>

        <div class="flex flex-col gap-6">
          {#each workflows as item, index (item)}
            <!-- <Card
              class="{selectedWorkflow === item
                ? 'border-hc-brand-600 ring-hc-brand-600'
                : ''}bg-hc-graylight-50 hover:bg-hc-brand-50 hover:border-hc-brand-600 hover:ring-hc-brand-600 gap-3 cursor-pointer"
              on:click={() => selectWorkflow(item)}
              size="xl"
            > -->
            <Card
              class="{selectedWorkflow === item
                ? 'border-hc-brand-600 ring-hc-brand-600'
                : ''}bg-hc-graylight-50 hover:bg-hc-brand-50 hover:border-hc-brand-600 hover:ring-hc-brand-600 gap-3 cursor-pointer"
              size="xl"
            >
              <h5 class="text-base font-semibold tracking-tight text-hc-graylight-900">
                {item.attributes.name}
              </h5>
              <p class="text-sm font-normal text-hc-graylight-900">
                {item.attributes.description}
              </p>
              <div class="flex justify-between">
                <Badge
                  rounded
                  color="primary"
                  class="bg-hc-brand-200 text-hc-brand-700 border border-hc-brand-300"
                  >{item.attributes.unit.data.attributes.name}</Badge
                >
              </div>
            </Card>
          {/each}
        </div>
      </div>
    </div>
  {/if}

  <!-- Step 2 -->
  {#if currentStep == 2}
    <div class="flex flex-col justify-center lg:px-72 md:px-10 sm:px-4 gap-4">
      <div class="flex justify-between items-start gap-3">
        <P size="3xl" weight="semibold">Upload your documents</P>
      </div>
      <div class="flex items-center justify-center w-full">
        <label
          for="dropzone-file"
          class="flex flex-col items-center justify-center w-full h-fit border-2 border-hc-graylight-300 border-dashed rounded-lg cursor-pointer bg-white dark:hover:bg-bray-800 dark:bg-gray-700 hover:bg-gray-100 dark:border-gray-600 dark:hover:border-gray-500 dark:hover:bg-gray-600"
        >
          <div class="flex flex-col items-center justify-center pt-5 pb-6">
            <enhanced:img
              src="../lib/icon/upload-icon.svg"
              class="w-10 h-10 ml-3"
              alt="An alt text"
            />
            <p class="mb-2 text-sm text-hc-graylight-600 dark:text-gray-400">
              <span class="font-semibold text-hc-brand-700">Click to upload</span>
              or drag and drop
            </p>
            <p class="text-xs text-gray-500 dark:text-gray-400">PDF File (Max.200KB)</p>
          </div>
          <!-- <input
            id="dropzone-file"
            type="file"
            class="hidden"
            accept=".pdf"
            bind:files
            multiple
            on:change={handleFileChange}
          /> -->
        </label>
      </div>

      {#if files && files.length > 0}
        <div class="flex flex-col gap-2">
          {#each files as file, index (file.name)}
            <div
              class="flex justify-between items-center bg-white px-3 py-3 rounded-lg border-2 border-hc-graylight-200"
            >
              <div class="flex justify-start items-center gap-3">
                <enhanced:img src="../lib/icon/pdf-icon.svg" class="w-10 h-10" alt="An alt text" />
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
                  <!-- 16 MB – 70% uploaded -->
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
            <!-- <li>{file.name}</li> -->
          {/each}
        </div>
      {/if}
    </div>
  {/if}

  <!-- Step 3 -->
  {#if currentStep == 3}
    {#if pdf}
      <div class="px-20">
        <div class="flex justify-between items-center pb-4">
          <P class="text-3xl leading-9 font-semibold text-hc-graylight-900 dark:text-hc-graydark-50"
            >Document Preview</P
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
          <div>
            <P class="font-semibold text-hc-graylight-900 dark:text-hc-graydark-50"
              >Additional Attributes</P
            >
            <P
              class="text-base font-medium leading-normal text-hc-graylight-600 dark:text-hc-graydark-400"
              >Add additional descriptive details to your documents</P
            >
            <Button
              outline
              on:click={() => (formModal = true)}
              color="purple"
              class="mt-4 w-max hover:!bg-transparent hover:text-hc-brand-600 !border-hc-brand-300 dark:!border-hc-brand-400 focus-within:!ring-1 px-3 py-1 gap-1"
            >
              Add Attributes
            </Button>
            <Modal bind:open={formModal} size="xs" autoclose={false} class="w-full">
              <form class="flex flex-col space-y-4" on:submit={handleSubmit}>
                <P
                  class="text-2xl leading-loose font-semibold text-hc-graylight-900 dark:text-hc-graydark-50"
                  >Add Attributes</P
                >
                <Label class="space-y-2">
                  <Input
                    bind:value={name}
                    type="text"
                    name="name"
                    placeholder="Attributes title (e.g Scores)"
                    required
                  />
                </Label>
                <Label class="space-y-2">
                  <Input
                    bind:value
                    type="text"
                    name="value"
                    placeholder="Attributes value (e.g C2 Proficient)"
                    required
                  />
                </Label>
                <Button type="submit" class="w-full">Add attributes</Button>
              </form>
            </Modal>
            {#if attributes.length > 0}
              <div class="flex flex-col gap-4 mt-4">
                {#each attributes as { name, value }, index (name)}
                  <div
                    class="grow px-6 py-3 bg-gray-50 rounded-lg backdrop-blur-[26px] justify-start items-center gap-4 flex"
                  >
                    <div class="grow shrink basis-0 justify-start items-center gap-4 flex-col">
                      <P class="font-semibold text-hc-graylight-900 dark:text-hc-graydark-50"
                        >{name}</P
                      >
                      <P
                        class="text-base font-medium leading-normal text-hc-graylight-600 dark:text-hc-graydark-400"
                        >{value}</P
                      >
                    </div>
                    <Button
                      on:click={() => removeAttribute(index)}
                      color="light"
                      class="bg-gray-50 !border-0 focus-within:!ring-1 !p-2"
                    >
                      <img alt="refresh" src={XCloseIcon} class="w-5 h-5" />
                    </Button>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
          <div class="flex justify-center grow py-10 bg-gray-50 rounded-lg dark:bg-gray-800">
            <!-- <PDFViewshowQrer {pdf} /> -->
            <!-- <PdfViewer showBorder={false} showButtons={false} scale={1} data={atob(base64)} /> -->
            <div>pdf viewer here</div>
          </div>
        </div>
      </div>
    {/if}
  {/if}

  <!-- Step 4 -->
  {#if currentStep == 4}
    {#if pdf}
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
