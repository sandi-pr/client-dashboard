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
    Dropdown,
    DropdownItem,
    Input,
    Card,
  } from "flowbite-svelte";
  import { ChevronDownOutline, SearchOutline, FilterOutline } from "flowbite-svelte-icons";
  import PdfViewer from "svelte-pdf";
  import { _ } from "svelte-i18n";
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
  import { CR3Call } from "hexacomponents-svelte";
  import { goto } from "$app/navigation";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";

  import { call } from "$lib";

  // Step 1
  let currentStep = 1;
  const steps = ["1", "2", "3"];

  // Step 2
  let files = [];
  let documents = [
    {
      name: "Document1.pdf",
      url: "",
    },
    {
      name: "Document2.docx",
      url: "",
    },
    {
      name: "Document3.doc",
      url: "",
    },
  ];
  let pdf = null;
  let fileUrl = null;
  let pdfFile = null;

  const qrCodeX = writable(0);
  const qrCodeY = writable(0);
  const isMoving = writable(false);

  let initialMouseX, initialMouseY;
  let addSignerModal = false;
  let selectRoleModal = false;
  let personClient = [];
  let selectedPersonClient = "";
  let selectedRole = "";
  let selectedDocument = "";

  function selectDocument(doc) {
    selectedDocument = doc;
  }

  // // Function to handle QR code movement
  function handleQrCodeMove(event) {
    isMoving.update((moving) => {
      if (moving) {
        const deltaX = event.clientX - initialMouseX;
        const deltaY = event.clientY - initialMouseY;

        qrCodeX.update((x) => x + deltaX);
        qrCodeY.update((y) => y + deltaY);

        initialMouseX = event.clientX;
        initialMouseY = event.clientY;
      }
      return moving;
    });
  }

  function startMove(event) {
    console.log("startMove");
    isMoving.set(true);
    initialMouseX = event.clientX;
    initialMouseY = event.clientY;

    document.addEventListener("mousemove", handleQrCodeMove);
    document.addEventListener("mouseup", stopMove);
  }

  function stopMove() {
    isMoving.set(false);
    saveQrCodePosition();

    document.removeEventListener("mousemove", handleQrCodeMove);
    document.removeEventListener("mouseup", stopMove);
  }

  function saveQrCodePosition() {
    qrCodeX.subscribe(
      (x) => (signers.find((signer) => signer.email === selectedSigner.email).position.x = x),
    );
    qrCodeY.subscribe(
      (y) => (signers.find((signer) => signer.email === selectedSigner.email).position.y = y),
    );
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

  let selectedSigner = null;
  let dropdownSigner = true;
  let signers = [];
  let showQr = [];
  let useMeterai = false;
  let emeterai = [];
  let fileTextName = null;
  let base64Data = "";

  function selectSigner(signer) {
    selectedSigner = signer;

    if (showQr.length < signers.length) {
      for (let i = 0; i < signers.length; i++) {
        showQr[i] = false;
      }
    }

    qrCodeX.set(signer.position.x);
    qrCodeY.set(signer.position.y);
    console.log("selectedSigner", selectedSigner);
  }

  function placeSignature() {
    if (selectedSigner === null) {
      selectedSigner = signers[0];
    }
    const existingSignerIndex = signers.findIndex(
      (signer) => signer.email === selectedSigner.email,
    );

    qrCodeX.set(signers[existingSignerIndex].position.x);
    qrCodeY.set(signers[existingSignerIndex].position.y);

    if (showQr.length < signers.length) {
      for (let i = 0; i < signers.length; i++) {
        showQr[i] = false;
      }
    }

    showQr[existingSignerIndex] = true;

    console.log("selectedSigner", selectedSigner);
    console.log("showQr", showQr);
  }

  function placeMeterai() {
    const existingSignerIndex = signers.findIndex(
      (signer) => signer.email === selectedSigner.email,
    );
    emeterai[existingSignerIndex] = true;
  }

  const validateCurrentStep = () => {
    if (currentStep === 1 && fileUrl === null) {
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

  const handleFileUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    // inputElement.accept = ".txt";
    inputElement.accept = ".pdf, .docx";

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      files = event.target.files;

      pdfFile = file;

      // console.log(file);

      if (file) {
        base64Data = await readFileAsBase64(file);

        const blob = base64toBlob(base64Data, file.type);
        fileUrl = URL.createObjectURL(blob);

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
      id: selectedPersonClient,
      name: personClient.find((person) => person.value === selectedPersonClient).name,
      email: personClient.find((person) => person.value === selectedPersonClient).email,
      role: selectedRole,
      position: { x: 0, y: 0 },
    };
    signers = [...signers, newSigner];

    selectedPersonClient = "";
    selectedRole = "";
  }

  function deleteSigner(id) {
    signers = signers.filter((signer) => signer.id !== id);
  }

  async function getPersons() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      let page = 1;
      while (true) {
        let response = await $call.call("client-dashboard", "get-persons", {
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
              name: person.attributes.fullname,
              email: person.attributes.email,
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

  const endStep = async () => {
    console.log(fileTextName);
    console.log(base64Data);
    console.log(signers);
    console.log("marking");
    const status = [
      {
        status: "INITIATE",
        timestamp: new Date().toISOString(),
      },
    ];
    console.log(status);

    let token = localStorage.getItem("token");
    if (token !== null) {
      let response = await $call.call("client-dashboard", "create-contract", {
        URLQuery: {},
        URLParams: {},
        Params: {
          contract_name: fileTextName,
          contract: base64Data,
          signers: JSON.stringify(signers),
          status: JSON.stringify(status),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      console.log(response);
      goto(await $call.getPage("client-dashboard", "document-contract"));
    }
  };

  onMount(async () => {
    getPersons();
  });
</script>

<Timeline {currentStep} {steps} {validateCurrentStep} {endStep}></Timeline>
<div class="mt-16">
  <!-- Step 1 -->
  {#if currentStep == 1}
    <div class="flex flex-col justify-center lg:px-72 md:px-10 sm:px-4 gap-4">
      <div class="flex grow justify-start items-start gap-3">
        <P size="3xl" weight="semibold">{$_("chooseDocument")}</P>
        <!-- <Button color="light" class="text-hc-graylight-700">
          <FilterOutline class="w-3 h-3 ms-2 text-hc-graylight-700 dark:text-white" />{$_("filter")}</Button
        > -->
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
        placeholder={$_("searchDocument")}
      >
        <SearchOutline slot="left" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
      </Input>
      <div class="flex flex-col justify-center mt-3 gap-3">
        <div class="flex text-start">
          <P class="text-hc-graylight-600 text-base font-normal">{$_("recentDocument")}</P>
        </div>

        <div class="flex flex-col gap-6">
          {#each documents as doc, index (doc)}
            <!-- <Card
              class="{selectedWorkflow === doc
                ? 'border-hc-brand-600 ring-hc-brand-600'
                : ''}bg-hc-graylight-50 hover:bg-hc-brand-50 hover:border-hc-brand-600 hover:ring-hc-brand-600 gap-3 cursor-pointer"
              on:click={() => selectWorkflow(doc)}
              size="xl"
            > -->
            <Card
              on:click={() => selectDocument(doc)}
              class="{selectedDocument === doc
                ? 'border-hc-brand-600 ring-hc-brand-600'
                : ''}bg-hc-graylight-50 hover:bg-hc-brand-50 hover:border-hc-brand-600 hover:ring-hc-brand-600 gap-3 cursor-pointer"
              size="xl"
            >
              <div class="flex gap-3">
                {#if doc.name.includes(".pdf")}
                  <img src={PDFSolidIcon} alt="" />
                {:else if doc.name.includes(".docx")}
                  <img src={DOCXSolidIcon} alt="" />
                {:else if doc.name.includes(".doc")}
                  <img src={DOCSolidIcon} alt="" />
                {/if}
                <div class="flex flex-col justify-center">
                  <P class="text-secondary-light dark:text-secondary-dark">{doc.name}</P>
                </div>
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
        <P size="3xl" weight="semibold">{$_("assignSigner")}</P>
        <Button
          outline
          on:click={() => (selectRoleModal = true)}
          class="gap-1 text-utility-brand-700-light dark:text-utility-brand-700-dark hover:!bg-transparent hover:text-utility-brand-700-light dark:hover:text-utility-brand-700-dark"
        >
          <UsersPlusIcon
            svgClass="w-5 h-5 text-utility-brand-700-light dark:text-utility-brand-700-dark"
          />
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
                  <div class="w-full self-stretch items-center justify-between flex gap-4">
                    <div class="flex items-center gap-4">
                      <DotsGridIcon
                        svgClass="w-6 h-6 text-fg-quaternary-light dark:text-fg-quaternary-dark"
                      />
                      <Avatar border />
                      <P size="lg" weight="semibold">{signer.name}</P>
                    </div>
                    <Badge border rounded color={signer.role == "signer" ? "primary" : "indigo"}
                      >{$_(signer.role)}</Badge
                    >
                  </div>
                  <CloseButton
                    on:click={() => deleteSigner(signer.id)}
                    class="text-fg-primary-light dark:text-fg-primary-dark"
                  />
                </div>
              {/each}
            </div>
          </div>
        {:else}
          <div
            class="flex grow h-60 items-center justify-center p-4 rounded-lg border border-primary-light dark:border-primary-dark"
          >
            <P size="md" weight="semibold" class="text-gray-500">{$_("noOneInvitedYet")}</P>
          </div>
        {/if}
      </div>
    </div>
    <Modal title={$_("whoDoYouWantToInvite")} bind:open={selectRoleModal} autoclose outsideclose>
      <div class="flex flex-col grow gap-3">
        <Button
          on:click={() => selectRole("signer")}
          color="light"
          class="flex grow p-6 rounded-lg border border-secondary-light dark:border-secondary-dark justify-start items-center gap-4"
        >
          <UserIcon
            svgClass="w-7 h-7 text-utility-gray-500-light dark:text-utility-gray-500-dark"
          />
          <P weight="semibold" class="text-primary-light dark:text-primary-dark">{$_("signer")}</P>
        </Button>
      </div>
      <div class="flex flex-col grow gap-3">
        <Button
          on:click={() => selectRole("reviewer")}
          color="light"
          class="flex grow p-6 rounded-lg border border-secondary-light dark:border-secondary-dark justify-start items-center gap-4"
        >
          <UserCheckIcon
            svgClass="w-7 h-7 text-utility-gray-500-light dark:text-utility-gray-500-dark"
          />
          <P weight="semibold" class="text-primary-light dark:text-primary-dark">{$_("reviewer")}</P
          >
        </Button>
      </div>
      <div class="flex flex-col grow gap-3">
        <Button
          on:click={() => selectRole("thirdPartySigner")}
          color="light"
          class="flex grow p-6 rounded-lg border border-secondary-light dark:border-secondary-dark justify-start items-center gap-4"
        >
          <UsersUpIcon
            svgClass="w-7 h-7 text-utility-gray-500-light dark:text-utility-gray-500-dark"
          />
          <P weight="semibold" class="text-primary-light dark:text-primary-dark"
            >{$_("thirdPartySigner")}</P
          >
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
        <Button disabled={!selectedPersonClient} on:click={addSigner} class="w-full"
          >{$_("addSigner")}</Button
        >
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
        </div>
        <div class="flex gap-12">
          <div class="flex flex-col gap-4">
            <P class="font-semibold text-hc-graylight-900 dark:text-hc-graydark-50"
              >Choose which signer to place the sign</P
            >
            <Button outline color="light" class="focus-within:!ring-1 acs">
              <Avatar />
              <div class="flex grow pl-2 pr-4">
                <Span
                  class="font-jakarta text-start text-slate-700 dark:text-hc-graydark-300 text-sm leading-tight"
                  >{selectedSigner !== null ? selectedSigner.name : signers[0].name}</Span
                >
              </div>
              <Button
                pill
                class={`
                  px-2 py-0 bg-hc-warning-50 dark:bg-hc-warning-950 text-hc-warning-700 dark:text-hc-warning-300 border border-hc-warning-200 dark:border-hc-warning-800 hover:bg-hc-warning-50 dark:hover:bg-hc-warning-950
                  rounded-md focus-within:ring-0 font-medium
                  `}
              >
                {selectedSigner !== null ? selectedSigner.role : signers[0].role}
              </Button>
              <ChevronDownOutline class="w-3 h-3 ms-2 dark:text-white" />
            </Button>
            <Dropdown triggeredBy=".acs" placement="bottom">
              {#each signers as signer}
                {#if signer.role === "signer"}
                  <DropdownItem>
                    <Button
                      on:click={() => {
                        selectSigner(signer);
                      }}
                      outline
                      color="light"
                      class="focus-within:!ring-1"
                    >
                      <Avatar />
                      <div class="flex-col grow pl-2 pr-4">
                        <Span
                          class="font-jakarta text-slate-700 dark:text-hc-graydark-300 text-sm leading-tight"
                          >{signer.name}</Span
                        >
                        <P
                          class="font-jakarta text-slate-600 dark:text-hc-graydark-400 text-sm leading-tight"
                          >{signer.email}</P
                        >
                      </div>
                      <Button
                        pill
                        class={`
                        px-2 py-0 bg-hc-warning-50 dark:bg-hc-warning-950 text-hc-warning-700 dark:text-hc-warning-300 border border-hc-warning-200 dark:border-hc-warning-800 hover:bg-hc-warning-50 dark:hover:bg-hc-warning-950
                        rounded-md focus-within:ring-0 font-medium
                        `}
                      >
                        {signer.role}
                      </Button>
                      <ChevronDownOutline class="w-3 h-3 ms-2 dark:text-white" />
                    </Button>
                  </DropdownItem>
                {/if}
              {/each}
            </Dropdown>
            <div class="flex gap-4 justify-center items-center">
              <Button
                on:click={placeSignature}
                class="w-full px-3 !py-2 gap-1 focus-within:!ring-1"
              >
                Place Signature
              </Button>
            </div>
          </div>
          <div class="flex justify-center grow py-10 bg-gray-50 rounded-lg dark:bg-gray-800">
            <!-- <PDFViewer {pdf} /> -->
            <div class="relative">
              <PdfViewer showBorder={false} showButtons={false} scale={1} data={atob(base64Data)} />
              <div class="">
                {#each signers as signer, index (signer)}
                  {#if showQr[index]}
                    {#if signer.email === selectedSigner.email}
                      <div
                        on:mousedown={selectedSigner.email === signer.email ? startMove : null}
                        role="button"
                        aria-label="Move QR Code"
                        tabindex="0"
                        style="position: absolute; top: {$qrCodeY}px; left: {$qrCodeX}px; cursor: {selectedSigner.email !==
                        signer.email
                          ? ''
                          : $isMoving
                            ? 'grabbing'
                            : 'grab'};"
                      >
                        <img src={QRCodeIcon} alt="" width="70" height="70" />
                      </div>
                    {:else}
                      <div
                        role="button"
                        aria-label="QR Code"
                        tabindex="0"
                        style="position: absolute; top: {signer.position.y}px; left: {signer
                          .position.x}px;"
                      >
                        <img src={QRCodeIcon} alt="" width="70" height="70" />
                        <span class="text-xs">{signer.name}</span>
                      </div>
                    {/if}
                  {/if}
                {/each}
              </div>
            </div>
          </div>
        </div>
      </div>
    {/if}
  {/if}
</div>
