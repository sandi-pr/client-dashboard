<script>
  import { Button, Heading, P } from "flowbite-svelte";
  import { timelineStep } from "../utils/timelineStep.js";

  export let titleHeading = "";
  export let steps = [];
  export let currentStep = 0;
  export let onboarding = false;
  export let validateCurrentStep = () => console.log("not function");
  export let endStep = () => console.log("not function");

  function increment() {
    timelineStep.update((value) => value + 1);
  }

  function decrement() {
    timelineStep.update((value) => value - 1);
  }

  export function nextStep() {
    if (currentStep < steps.length) {
      if (onboarding) {
        if (validateCurrentStep()) {
          currentStep += 1;
          increment();
        }
      } else {
        currentStep += 1;
        increment();
      }
      console.log(currentStep, steps.length);
    } else {
      endStep();
    }
  }

  export function prevStep() {
    if (currentStep > 0) {
      if (currentStep >= 2) {
        currentStep -= 1;
        decrement();
      } else {
        window.history.back();
      }
    }
  }
</script>

<div
  class="flex justify-center items-center px-20 border-b border-secondary-light dark:border-secondary-dark"
  class:px-36={onboarding}
>
  <Button pill={true} color="light" class="!p-[10px] !border-0" on:click={prevStep}
    ><svg
      class="w-4 h-4 text-hc-graylight-600 dark:text-white"
      aria-hidden="true"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 8 14"
    >
      <path
        stroke="currentColor"
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M7 1 1.3 6.326a.91.91 0 0 0 0 1.348L7 13"
      ></path>
    </svg></Button
  >

  <enhanced:img src="../icon/Logomark.svg" class="w-10 h-10 ml-3" alt="An alt text" />
  {#if onboarding}
    <Heading tag="h4" class="ml-2 w-1/3">Set up your company</Heading>
  {/if}

  {#if titleHeading}
    <P size="base" weight="bold" class="min-w-max pl-4 text-primary-light dark:text-primary-dark"
      >{titleHeading}</P
    >
  {/if}

  <div class="flex w-full justify-between px-6 py-6">
    {#each steps as step, index (index)}
      {#if index !== steps.length - 1}
        <div class="flex items-center justify-start w-full">
          <div
            class="flex-shrink-0 w-6 h-6 rounded-full dark:bg-hc-brand-400 dark:ring-gray-900 text-center py-0.5"
            class:bg-hc-brand-600={index < currentStep}
            class:bg-hc-graylight-100={index >= currentStep}
          >
            <p
              class="text-sm text-white dark:text-white"
              class:text-white={index < currentStep}
              class:text-hc-graylight-500={index >= currentStep}
            >
              {step}
            </p>
          </div>
          <div
            class="sm:flex flex-grow bg-gray-200 h-1 dark:bg-gray-700"
            class:bg-hc-brand-600={index < currentStep}
            class:bg-hc-graylight-100={index >= currentStep - 1}
          />
        </div>
      {:else}
        <div
          class="flex-shrink-0 w-6 h-6 rounded-full dark:bg-hc-brand-400 dark:ring-gray-900 text-center py-0.5"
          class:bg-hc-brand-600={index < currentStep}
          class:bg-hc-graylight-100={index >= currentStep}
        >
          <p
            class="text-sm text-white dark:text-white"
            class:text-white={index < currentStep}
            class:text-hc-graylight-500={index >= currentStep}
          >
            {step}
          </p>
        </div>
      {/if}
    {/each}
  </div>
  {#if currentStep === steps.length}
    <Button
      class="ml-4 bg-hc-brand-600 hover:bg-hc-brand-700 focus-within:ring-4 focus-within:ring-hc-brand-300"
      size="sm"
      on:click={nextStep}>Finish</Button
    >
  {:else}
    <Button
      class="ml-4 bg-hc-brand-600 hover:bg-hc-brand-700 focus-within:ring-4 focus-within:ring-hc-brand-300"
      size="sm"
      on:click={nextStep}>Next</Button
    >
  {/if}
</div>
