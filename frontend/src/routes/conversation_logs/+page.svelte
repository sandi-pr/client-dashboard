<script>
  import { P } from "flowbite-svelte";
  import UserIcon from "../lib/icon/users-01.svg";
  // import VerifiedIcon from "../lib/icon/verified.svg";

  let dataConversatonLogs = [
    {
      role: "user",
      message:
        "Hi there, I'm currently in [Location] and I'm in need of some services. Specifically, I'm looking for a reliable wood flooring company. Do you have any recommendations?",
    },
    {
      role: "assistant",
      message:
        'Of course! London has several reputable wood flooring companies. One notable option is "Zway Flooring Solutions" located in the heart of the city. They are well-regarded for their premium materials and skilled craftsmanship. Would you like more information about them or any other suggestions?',
    },
  ];
</script>

<div class="flex flex-col">
  <div class="mx-4 py-6 border-b border-secondary-light dark:border-secondary-dark">
    <P size="3xl" weight="semibold" class="text-primary-light dark:text-primary-dark leading-[38px]"
      >Conversation Logs 4c64986428240cn948929284
    </P>
  </div>
  {#each dataConversatonLogs as qaPair}
    <div
      class="flex bg-secondary-light dark:bg-secondary-dark rounded-2xl px-8 py-6 mx-4 my-6 items-start gap-4"
    >
      <div
        class="w-12 h-12 relative flex justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
      >
        <img src={UserIcon} alt="User" class="" />
        {#if qaPair.role === "assistant"}
          <div class="absolute bottom-0 right-0">
            <!-- <img src={VerifiedIcon} alt="User" style={{ width: 16, height: 16 }} /> -->
            \/
          </div>
        {/if}
      </div>
      <div class="flex flex-col w-full gap-1">
        <div class="flex items-center space-x-2 rtl:space-x-reverse">
          <P weight="bold" class="text-primary-light dark:text-primary-dark"
            >{qaPair.role === "user" ? "Customer" : "Assistant"}</P
          >
        </div>
        <P weight="medium" class="text-tertiary-light dark:text-tertiary-dark">{qaPair.message}</P>
      </div>
    </div>
  {/each}
</div>
