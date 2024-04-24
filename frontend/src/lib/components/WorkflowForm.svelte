<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Dean Ramhan <dean@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import AutoComplete from "simple-svelte-autocomplete";
  import { page } from "$app/stores";
  import { dndzone } from "svelte-dnd-action";
  import { flip } from "svelte/animate";
  import { goto } from "$app/navigation";
  // import axios from "axios";
  // import {
  //   PUBLIC_API_URL_DEV,
  //   PUBLIC_API_URL_PROD,
  //   PUBLIC_STRAPI_URL_DEV,
  //   PUBLIC_STRAPI_URL_PROD,
  // } from "$env/static/public";

  import * as Icon from "flowbite-svelte-icons";
  import { onMount } from "svelte";
  import {
    P,
    Hr,
    Card,
    Label,
    Input,
    Modal,
    Select,
    Button,
    Textarea,
    ButtonGroup,
    Avatar,
  } from "flowbite-svelte";
  import LoadingSpinner from "../components/LoadingSpinner.svelte";
  import logoMark from "../icon/Logomark.svg";
  import backChevron from "../icon/back-chevron.svg";
  import drag from "../icon/drag.svg";
  import userCheck from "../icon/user-check.svg";
  import userUp from "../icon/user-up.svg";
  // import { userDataSession } from "$lib/utils/auth";

  // const API_URL = process.env.NODE_ENV === "development" ? PUBLIC_API_URL_DEV : PUBLIC_API_URL_PROD;
  // const STRAPI_URL =
  //   process.env.NODE_ENV === "development" ? PUBLIC_STRAPI_URL_DEV : PUBLIC_STRAPI_URL_PROD;
  // const sessionOrganization = JSON.parse(sessionStorage.getItem("organization"));

  let workflowObject = {
    essentials: {
      name: "",
      description: "",
      organization: "",
      unit: "",
    },
    options: [false, false, false, false],
    signerReviewers: [],
  };
  let units = [];

  let inviteType = [
    { name: "By User", value: "User" },
    { name: "By Role", value: "Role" },
    { name: "By Group", value: "Group" },
  ];
  let inviteTypeSelected = "User";
  let inviteUnit = [];
  let inviteUnitSelected = "";

  let currentStep = 1;
  let showModal = false;
  const flipDurationMs = 100;
  let modalType = "Signer";
  let searchInput = "";
  let cannotCreate = true;
  let isLoading = false;

  const handleConsider = (evt) => {
    workflowObject.signerReviewers = evt.detail.items;
  };

  const handleFinalize = (evt) => {
    workflowObject.signerReviewers = evt.detail.items;
  };

  const handleCheckChange = (idx, value) => {
    workflowObject.options[idx] = value;
  };

  const handleGetWorkflow = async () => {
    isLoading = true;
    console.log("not yet implemented");
    // try {
    //   const response = await axios.get(`${API_URL}/api/workflows/${workflowId}`, {
    //     headers: {
    //       Authorization: "Bearer " + userDataSession().jwt,
    //     },
    //   });

    //   // console.log(response);

    //   workflowObject.essentials.name = response.data.data.attributes.name;
    //   workflowObject.essentials.description = response.data.data.attributes.description || "";
    //   workflowObject.essentials.organization = response.data.data.attributes.organization.data.id;
    //   workflowObject.essentials.unit = response.data.data.attributes.unit.data.id;

    //   workflowObject.options[0] = response.data.data.attributes.signature_type === "PSRE";
    //   workflowObject.options[1] = response.data.data.attributes.process_type === "Sequential";
    //   workflowObject.options[2] = response.data.data.attributes.review_type === "Review";
    //   workflowObject.options[3] = response.data.data.attributes.meterai_type === "Meterai";

    //   workflowObject.signerReviewers = await Promise.all(
    //     response.data.data.attributes.steps.map(async (signerReviewer, index) => {
    //       const responseUser = await axios.get(
    //         `${API_URL}/api/users?type=email&value=${response.data.data.attributes.steps[index].user.email}`,
    //         {
    //           headers: {
    //             Authorization: "Bearer " + userDataSession().jwt,
    //           },
    //         }
    //       );
    //       // console.log(responseUser);
    //       return {
    //         id: `${responseUser.data[0].id}/${signerReviewer.step}`,
    //         identityId: responseUser.data[0].identity.id,
    //         type: signerReviewer.step,
    //         unit: responseUser.data[0].unit.name,
    //         email: responseUser.data[0].email,
    //         username: responseUser.data[0].username,
    //         photoUrl: STRAPI_URL + responseUser.data[0].photo.url,
    //       };
    //     })
    //   );
    //   // console.log(workflowObject.signerReviewers);

    //   cannotCreate = false;
    // } catch (error) {
    //   console.log(error);
    // }
    isLoading = false;
  };

  const handleGetUnits = async () => {
    isLoading = true;
    console.log("not yet implemented");
    // try {
    //   const response = await axios.get(
    //     `${API_URL}/api/units?organizationId=${
    //       JSON.parse(sessionStorage.getItem("organization")).id
    //     }`,
    //     {
    //       headers: {
    //         Authorization: "Bearer " + userDataSession().jwt,
    //       },
    //     }
    //   );
    //   // console.log(response);

    //   if (response?.data?.status === 401) {
    //     alert("Session expired, please login again");
    //     goto("/");
    //   } else if (response?.data?.data && response.data.data.length > 0) {
    //     let mappedData = response.data.data.map((unit) => {
    //       return {
    //         name: unit.attributes.name,
    //         value: unit.id,
    //       };
    //     });

    //     units = mappedData.slice();

    //     mappedData.unshift({ name: "Units", value: "" });

    //     inviteUnit = [...mappedData];
    //   } else {
    //     units = [];
    //   }
    // } catch (error) {
    //   console.log(error);
    // }
    isLoading = false;
  };

  const handleAddSignerReviewers = async (type) => {
    console.log("not yet implemented");
    // try {
    //   if (
    //     workflowObject.signerReviewers.some(
    //       (signerReviewer) => signerReviewer.email === searchInput.email
    //     )
    //   ) {
    //     alert("User already invited");
    //     return;
    //   }

    //   const newSignerReviewer = {
    //     id: `${searchInput.id}/${type}`,
    //     identityId: searchInput.identity.id,
    //     type: type,
    //     unit: searchInput.unit.name,
    //     email: searchInput.email,
    //     username: searchInput.username,
    //     photoUrl: STRAPI_URL + searchInput.photo.url,
    //   };

    //   searchInput = "";
    //   modalType = "";
    //   showModal = false;
    //   cannotCreate = false;

    //   workflowObject.signerReviewers = [...workflowObject.signerReviewers, newSignerReviewer];

    //   // console.log(workflowObject);
    // } catch (error) {
    //   console.log(error);
    // }
  };

  const handleRemoveSignerReviewers = (id) => {
    workflowObject.signerReviewers = workflowObject.signerReviewers.filter(
      (signerReviewer) => signerReviewer.id !== id
    );

    if (workflowObject.signerReviewers.length === 0) cannotCreate = true;
  };

  const handleSearchInvitees = async (value) => {
    console.log("not yet implemented");
    // try {
    //   const response = await axios.get(
    //     `${API_URL}/api/users?type=search&value=${value}&role=${inviteTypeSelected}&unit=${inviteUnitSelected}`,
    //     {
    //       headers: {
    //         Authorization: "Bearer " + userDataSession().jwt,
    //       },
    //     }
    //   );
    //   // console.log(response);

    //   if (response?.data?.status === 401) {
    //     alert("Session expired, please login again");
    //     goto("/");
    //   }

    //   if (response?.data) {
    //     return response.data;
    //   } else {
    //     return [];
    //   }
    // } catch (error) {
    //   console.log(error);
    // }
  };

  const handleSubmitWorkflow = async () => {
    isLoading = true;
    console.log("not yet implemented");
    // try {
    //   let response = null;
    //   if (workflowId) {
    //     response = await axios.put(`${API_URL}/api/workflows/${workflowId}`, workflowObject, {
    //       headers: {
    //         Authorization: "Bearer " + userDataSession().jwt,
    //       },
    //     });
    //   } else {
    //     response = await axios.post(`${API_URL}/api/workflows`, workflowObject, {
    //       headers: {
    //         Authorization: "Bearer " + userDataSession().jwt,
    //       },
    //     });
    //   }

    //   console.log(response);

    //   if (response?.status === 401) {
    //     alert("Session expired, please login again");
    //     goto("/");
    //   } else if (response?.data?.id) {
    //     if (workflowId) alert("Workflow updated successfully");
    //     else alert("Workflow created successfully");
    //     goto("/workflow");
    //   }
    // } catch (error) {
    //   console.log(error);
    // }
    isLoading = false;
  };

  onMount(() => {
    if (workflowId) {
      handleGetWorkflow();
    }

    handleGetUnits();

    if (JSON.parse(sessionStorage.getItem("organization"))) {
      workflowObject.essentials.organization = JSON.parse(
        sessionStorage.getItem("organization")
      ).id;
    } else {
      workflowObject.essentials.organization = "";
    }
  });

  $: buttonDisable = () => {
    if (currentStep === 1) {
      return workflowObject.essentials.name === "";
    } else if (currentStep === 3) {
      return cannotCreate;
    }
  };

  $: workflowId = $page.params.strapiId;
  $: stepOneClear = workflowObject.essentials.name !== "" && workflowObject.essentials.unit !== "";
</script>

<LoadingSpinner {isLoading} />

<main class="">
  <div class="relative isolate overflow-hidden">
    <!-- Header -->
    <div
      class="flex items-center gap-x-8 gap-y-4 px-4 py-4 sm:flex-row sm:items-center sm:px-6 lg:px-8 my-4 rounded-md"
    >
      <button
        class="flex items-center gap-x-5 cursor-pointer"
        on:click={() => {
          if (currentStep > 1) {
            currentStep = currentStep - 1;
          } else {
            goto("/workflow");
          }
        }}
      >
        <img src={backChevron} alt="Back Icon" />
        <img src={logoMark} alt="Brand Logo" />
      </button>

      <P size="lg" weight="semibold">{workflowId ? "Update" : "Create"} Workflow</P>

      <div class="flex-grow text-center">
        <ol class="flex items-center w-full justify-center">
          <li
            class={`flex w-full items-center after:w-full after:h-1 after:border-b ${
              currentStep > 1
                ? "after:border-hc-brand-600 dark:after:border-hc-brand-600"
                : "after:border-hc-graylight-200 dark:after:border-hc-graydark-200"
            } after:border-4 after:inline-block`}
          >
            <span
              class="flex items-center justify-center w-[24px] h-[24px] rounded-full lg:h-[24px] lg:w-[24px] dark:shrink-0"
            >
              <button
                class={`w-[24px] h-[24px] bg-hc-brand-600 rounded-full cursor-pointer`}
                on:click={() => {
                  currentStep = 1;
                }}
              >
                <p
                  class="text-xs font-semibold leading-[24px] lg:text-sm lg:leading-[24px] text-white"
                >
                  1
                </p>
              </button>
            </span>
          </li>
          <li
            class={`flex w-full items-center after:w-full after:h-1 after:border-b ${
              currentStep > 2
                ? "after:border-hc-brand-600 dark:after:border-hc-brand-600"
                : "after:border-hc-graylight-200 dark:after:border-hc-graydark-200"
            } after:border-4 after:inline-block`}
          >
            <span
              class="flex items-center justify-center w-[24px] h-[24px] rounded-full lg:h-[24px] lg:w-[24px] dark:shrink-0"
            >
              <button
                class={`w-[24px] h-[24px] bg-hc-brand-600 rounded-full ${
                  stepOneClear ? "cursor-pointer" : "cursor-not-allowed"
                }`}
                on:click={() => {
                  if (stepOneClear) currentStep = 2;
                }}
              >
                <p
                  class="text-xs font-semibold leading-[24px] lg:text-sm lg:leading-[24px] text-white"
                >
                  2
                </p>
              </button>
            </span>
          </li>
          <li>
            <span
              class="flex items-center justify-center w-[24px] h-[24px] rounded-full lg:h-[24px] lg:w-[24px] dark:shrink-0"
            >
              <button
                class={`w-[24px] h-[24px] bg-hc-brand-600 rounded-full ${
                  stepOneClear ? "cursor-pointer" : "cursor-not-allowed"
                }`}
                on:click={() => {
                  if (stepOneClear) currentStep = 3;
                }}
              >
                <p
                  class="text-xs font-semibold leading-[24px] lg:text-sm lg:leading-[24px] text-white"
                >
                  3
                </p>
              </button>
            </span>
          </li>
        </ol>
      </div>

      <div class="ml-auto">
        <Button
          disabled={buttonDisable()}
          on:click={() => {
            if (currentStep === 3) {
              handleSubmitWorkflow();
            } else {
              currentStep = currentStep + 1;
            }
          }}
        >
          {currentStep === 3 ? (workflowId ? "Update" : "Create") : "Next"}
        </Button>
      </div>
    </div>

    <Hr />

    <!-- Content -->
    <!-- Step 1 -->
    {#if currentStep === 1}
      <div class="flex flex-col w-[50%] mx-auto mt-10 mb-5">
        <P size="4xl" weight="semibold">Workflow Information</P>

        <div class="flex flex-col">
          <div class="my-2">
            <Label class="text-gray-600 text-lg font-semibold mb-1">Workflow name*</Label>
            <Input
              class="bg-white dark:bg-hc-graydark-700"
              placeholder="Enter your workflow name"
              type="text"
              name="name"
              id="name"
              bind:value={workflowObject.essentials.name}
              required
            />
          </div>
          <Label class="text-gray-600 text-lg font-semibold mb-1">
            Select an Unit*
            <Select
              class="mt-2 font-normal"
              items={units}
              bind:value={workflowObject.essentials.unit}
            />
          </Label>
          <div class="my-2">
            <Label class="text-gray-600 text-lg font-semibold mb-1">Description</Label>
            <Textarea
              class="bg-white resize-none"
              rows="4"
              placeholder="Enter your workflow description"
              name="description"
              id="description"
              bind:value={workflowObject.essentials.description}
            />
          </div>
        </div>
      </div>
    {:else if currentStep === 2}
      <div class="grid grid-cols-12 gap-4 w-[80%] mx-auto mt-10 mb-24">
        <P size="4xl" weight="semibold" class="col-span-12 mb-2">Define your settings</P>

        <!-- first row -->
        <P size="xl" weight="semibold" class="col-span-2">Signature Type</P>
        <button
          class={"px-4 col-span-5 flex items-center justify-between bg-hc-graylight-50 rounded-lg border-[1px] border-hc-graylight-100 " +
            (!workflowObject.options[0]
              ? "bg-hc-brand-50 border-hc-brand-700 text-hc-brand-600"
              : "")}
          on:click={() => {
            handleCheckChange(0, false);
          }}
        >
          <div class="flex items-center">
            <Icon.FileLinesOutline size="xl" strokeWidth="1" />
            <div class="p-4">
              <p class="text-lg font-semibold">Non-PSRE</p>
              <p>
                Any documents that doesn't come from PSRE (Penyelengara Sertifikat Elektronik)
                Indonesia
              </p>
            </div>
          </div>
          <input
            type="radio"
            name="signature"
            id="signature_false"
            class="size-5 text-hc-brand-600"
            checked={!workflowObject.options[0]}
            readOnly
          />
        </button>
        <button
          class={"px-4 col-span-5 flex items-center justify-between bg-hc-graylight-50 rounded-lg border-[1px] border-hc-graylight-100 " +
            (workflowObject.options[0]
              ? "bg-hc-brand-50 border-hc-brand-700 text-hc-brand-600"
              : "")}
          on:click={() => {
            handleCheckChange(0, true);
          }}
        >
          <div class="flex items-center">
            <Icon.FileLinesOutline size="xl" strokeWidth="1" />
            <div class="p-4">
              <p class="text-lg font-semibold">PSRE</p>
              <p>A document that comes from PSRE (Penyelengara Sertifikat Elektronik) Indonesia</p>
            </div>
          </div>
          <input
            type="radio"
            name="signature"
            id="signature_true"
            class="size-5 text-hc-brand-600"
            checked={workflowObject.options[0]}
            readOnly
          />
        </button>
        <Hr hrClass="col-span-12 my-4" />

        <!-- second row -->
        <P size="xl" weight="semibold" class="col-span-2">Signing Type</P>
        <button
          class={"px-4 col-span-5 flex items-center justify-between bg-hc-graylight-50 rounded-lg border-[1px] border-hc-graylight-100 " +
            (!workflowObject.options[1]
              ? "bg-hc-brand-50 border-hc-brand-700 text-hc-brand-600"
              : "")}
          on:click={() => {
            handleCheckChange(1, false);
          }}
        >
          <div class="flex items-center">
            <Icon.ShuffleOutline size="xl" strokeWidth="1" />
            <div class="p-4">
              <p class="text-lg font-semibold">Parallel</p>
              <p>
                Each signer must sign the document one after another, without following a
                predetermined order
              </p>
            </div>
          </div>
          <input
            type="radio"
            name="signing"
            id="signing_false"
            class="size-5 text-hc-brand-600"
            checked={!workflowObject.options[1]}
            readOnly
          />
        </button>
        <button
          class={"px-4 col-span-5 flex items-center justify-between bg-hc-graylight-50 rounded-lg border-[1px] border-hc-graylight-100 " +
            (workflowObject.options[1]
              ? "bg-hc-brand-50 border-hc-brand-700 text-hc-brand-600"
              : "")}
          tabindex="0"
          on:click={() => {
            handleCheckChange(1, true);
          }}
        >
          <div class="flex items-center">
            <Icon.ArrowsRepeatOutline size="xl" strokeWidth="1" />
            <div class="p-4">
              <p class="text-lg font-semibold">Sequential</p>
              <p>
                Each signer must sign the document one after another, following a predetermined
                order
              </p>
            </div>
          </div>
          <input
            type="radio"
            name="signing"
            id="signing_true"
            class="size-5 text-hc-brand-600"
            checked={workflowObject.options[1]}
            readOnly
          />
        </button>
        <Hr hrClass="col-span-12 my-4" />

        <!-- third row -->
        <P size="xl" weight="semibold" class="col-span-2">Review Type</P>
        <button
          class={"px-4 col-span-5 flex items-center justify-between bg-hc-graylight-50 rounded-lg border-[1px] border-hc-graylight-100 " +
            (!workflowObject.options[2]
              ? "bg-hc-brand-50 border-hc-brand-700 text-hc-brand-600"
              : "")}
          on:click={() => {
            handleCheckChange(2, false);
          }}
        >
          <div class="flex items-center">
            <Icon.UserRemoveOutline size="xl" strokeWidth="1" />
            <div class="p-4">
              <p class="text-lg font-semibold">No need to review</p>
              <p>The document doesn't need to be reviewed</p>
            </div>
          </div>
          <input
            type="radio"
            name="review"
            id="review_false"
            class="size-5 text-hc-brand-600"
            checked={!workflowObject.options[2]}
            readOnly
          />
        </button>
        <button
          class={"px-4 col-span-5 flex items-center justify-between bg-hc-graylight-50 rounded-lg border-[1px] border-hc-graylight-100 " +
            (workflowObject.options[2]
              ? "bg-hc-brand-50 border-hc-brand-700 text-hc-brand-600"
              : "")}
          on:click={() => {
            handleCheckChange(2, true);
          }}
        >
          <div class="flex items-center">
            <Icon.UserAddOutline size="xl" strokeWidth="1" />
            <div class="p-4">
              <p class="text-lg font-semibold">Need Review</p>
              <p>The document needs to be reviewed by specific user</p>
            </div>
          </div>
          <input
            type="radio"
            name="review"
            id="review_true"
            class="size-5 text-hc-brand-600"
            checked={workflowObject.options[2]}
            readOnly
          />
        </button>
        <Hr hrClass="col-span-12 my-4" />

        <!-- fourth row -->
        <P size="xl" weight="semibold" class="col-span-2">E-Meterai</P>
        <button
          class={"px-4 col-span-5 flex items-center justify-between bg-hc-graylight-50 rounded-lg border-[1px] border-hc-graylight-100 " +
            (!workflowObject.options[3]
              ? "bg-hc-brand-50 border-hc-brand-700 text-hc-brand-600"
              : "")}
          on:click={() => {
            handleCheckChange(3, false);
          }}
        >
          <div class="flex items-center">
            <Icon.UserRemoveOutline size="xl" strokeWidth="1" />
            <div class="p-4">
              <p class="text-lg font-semibold">No need E-Meterai</p>
              <p>The document doesn't need E-Meterai</p>
            </div>
          </div>
          <input
            type="radio"
            name="meterai"
            id="meterai_false"
            class="size-5 text-hc-brand-600"
            checked={!workflowObject.options[3]}
            readOnly
          />
        </button>
        <button
          class={"px-4 col-span-5 flex items-center justify-between bg-hc-graylight-50 rounded-lg border-[1px] border-hc-graylight-100 " +
            (workflowObject.options[3]
              ? "bg-hc-brand-50 border-hc-brand-700 text-hc-brand-600"
              : "")}
          on:click={() => {
            handleCheckChange(3, true);
          }}
        >
          <div class="flex items-center">
            <Icon.UserAddOutline size="xl" strokeWidth="1" />
            <div class="p-4">
              <p class="text-lg font-semibold">Need E-Meterai</p>
              <p>The document needs E-Meterai</p>
            </div>
          </div>
          <input
            type="radio"
            name="meterai"
            id="meterai_true"
            class="size-5 text-hc-brand-600"
            checked={workflowObject.options[3]}
            readOnly
          />
        </button>
      </div>
    {:else if currentStep === 3}
      <div class="flex flex-col w-[50%] mx-auto mt-10 mb-24">
        <div class="flex justify-between items-center">
          <P size="4xl" weight="semibold">Assign Signer</P>
          <Button
            outline
            on:click={() => {
              showModal = true;
              modalType = "Invite";
            }}><Icon.UserAddOutline size="md" strokeWidth="1" />&nbsp;Invite</Button
          >
        </div>
        <p class="flex mb-4">
          For sequential, adjust & move order by dragging &nbsp;<span
            ><img src={drag} alt="Drag Indicator" /></span
          > &nbsp; icon
        </p>
        {#if workflowObject.signerReviewers.length > 0}
          <section
            class="px-4 pt-4 w-full rounded-lg border-[1px] border-hc-graylight-300"
            use:dndzone={{ items: workflowObject.signerReviewers, flipDurationMs }}
            on:consider={handleConsider}
            on:finalize={handleFinalize}
          >
            {#each workflowObject.signerReviewers as signerReviewer (signerReviewer.id)}
              <div
                class="flex justify-between p-4 bg-hc-graylight-50 rounded-lg mb-4"
                animate:flip={{ duration: flipDurationMs }}
              >
                <div class="flex items-center">
                  <img src={drag} alt="Drag Indicator" class="mr-4" />
                  <Avatar
                    class="w-12 h-12 mr-4"
                    src={signerReviewer.photoUrl}
                    alt={signerReviewer.username}
                  />
                  <div>
                    <P size="xl" weight="semibold">{signerReviewer.username}</P>
                    <P size="lg" weight="medium">{signerReviewer.unit}</P>
                  </div>
                </div>
                <div class="flex items-center pr-8">
                  <P
                    size="sm"
                    weight="semibold"
                    class="py-1 px-3 rounded-3xl border-[1px] text-hc-brand-700 bg-hc-brand-100 border-hc-brand-700 mr-4"
                    >{signerReviewer.type}</P
                  >
                  <button
                    on:click={() => {
                      handleRemoveSignerReviewers(signerReviewer.id);
                    }}
                  >
                    <P size="3xl" class="text-hc-graylight-400">X</P>
                  </button>
                </div>
              </div>
            {/each}
          </section>
        {:else}
          <div
            class="w-full h-[20rem] border-2 border-hc-graylight-200 flex justify-center items-center rounded-lg"
          >
            <p class="text-gray-600 mt-4 font-semibold">No one invited yet</p>
          </div>
        {/if}
      </div>
    {/if}
  </div>
</main>

<Modal
  title="Invite Signer & Reviewer"
  class="w-[75%]"
  bind:open={showModal}
  autoclose
  outsideclose
>
  {#if modalType === "Invite"}
    {#each [0, 1, 2, 3] as index}
      <Card
        padding="none"
        class="p-2 cursor-pointer bg-hc-graylight-50 hover:bg-hc-graylight-100 max-w-full"
        on:click={() => {
          if (index === 0) {
            modalType = "Signer";
          } else if (index === 1) {
            modalType = "Reviewer";
          } else if (index === 2) {
            modalType = "Third Party Signer";
          } else {
            modalType = "PSRE Signer";
          }
        }}
      >
        <div class="flex p-1 items-center justify-between">
          <div class="flex items-center">
            {#if index === 0}
              <Icon.UserOutline size="md" strokeWidth="2" />
              <P>&nbsp; Signer</P>
            {:else if index === 1}
              <img src={userCheck} alt="Reviewer Icon" />
              <P>&nbsp; Reviewer</P>
            {:else if index === 2}
              <img src={userUp} alt="Third Party Icon" />
              <P>&nbsp; Third Party Signer</P>
            {:else}
              <Icon.TicketOutline size="md" strokeWidth="2" />
              <P>&nbsp; PSRE Signer</P>
            {/if}
          </div>
          <!-- <P class="text-hc-brand-600">+1</P> -->
        </div>
      </Card>
    {/each}
    <!-- <Button class="w-full" disabled={workflowObject.signerReviewers.length === 0}>Invite</Button> -->
  {:else}
    <div class="flex w-full p-1 items-center justify-between">
      <ButtonGroup class="grid grid-cols-12">
        <div
          class={`col-span-6 [&>*]:w-full [&>*]:!h-[2.75rem] [&>*]:m-0
          [&>*>div]:rounded-lg
        `}
        >
          <AutoComplete
            searchFunction={handleSearchInvitees}
            delay="1000"
            localFiltering={false}
            labelFieldName="username"
            valueFieldName="id"
            bind:selectedItem={searchInput}
            class="h-[2.75rem] w-full m-0 border-hc-graylight-300 rounded-l-lg focus:ring-hc-brand-500 focus:border-hc-brand-500"
          />
        </div>
        <!-- <div class="col-span-5"></div> -->
        <Select
          class={inviteTypeSelected === "Group"
            ? "rounded-none rounded-r-lg col-span-3"
            : "rounded-none col-span-3"}
          placeholder="Types"
          items={inviteType}
          bind:value={inviteTypeSelected}
        />
        <Select
          class={inviteTypeSelected === "Group" ? "hidden" : "rounded-none rounded-r-lg col-span-3"}
          placeholder="Units"
          items={inviteUnit}
          bind:value={inviteUnitSelected}
        />
      </ButtonGroup>
    </div>
    <div class="w-full grid grid-cols-2 gap-4">
      <div></div>
      <Button
        class="col-span-1"
        disabled={searchInput.length === 0}
        on:click={() => {
          handleAddSignerReviewers(modalType);
        }}>Add {modalType}</Button
      >
    </div>
  {/if}
</Modal>
