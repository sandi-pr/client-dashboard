<script>
  import { Button, Dropdown, DropdownItem } from "flowbite-svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import SendIcon from "$lib/icon/send-01.svg";
  import { sidebarLargeClosed } from "$lib/utils/sidebarLargeClosed";

  export let textarea;
  export let messageInput;
  export let handleKeyDown;
  export let isLoading;
  export let disableInput;
  export let showSendMessage;
  export let onclickSend;
  export let selectedRole;
  export let hideSidebar;
  let dropdownOpen = false;

  let dropdownItem = [
    {
      name: "As Bots",
      value: "As Bots",
      selected: false,
    },
    {
      name: "As Human",
      value: "As Human",
      selected: true,
    },
  ];

  const handleDropdown = (evt) => {
    // console.log(evt.target.value);
    dropdownItem = dropdownItem.map((item) => {
      if (item.value === evt.target.value) {
        item.selected = true;
        selectedRole = evt.target.value;
      } else {
        item.selected = false;
      }
      return item;
    });
    dropdownOpen = false;
  };

  const dropdownCurrentlySelected = () => {
    selectedRole = dropdownItem.find((item) => item.selected === true).value;
    return dropdownItem.find((item) => item.selected === true).name;
  };

  $: dropdownItem = dropdownItem.map((item) => {
    if (item.value === selectedRole) {
      item.selected = true;
    } else {
      item.selected = false;
    }
    return item;
  });
</script>

<div
  class={`-z-0 bg-gradient-to-t from-white dark:from-utility-bg-primary-dark from-80% fixed flex items-center bottom-0 left-0 right-0 ${!$sidebarLargeClosed ? hideSidebar ? "pl-0 lg:pl-[670px]" : "pl-[350px] lg:pl-[670px]" : hideSidebar ? "pl-0 lg:pl-[430px]" : "pl-[350px] lg:pl-[430px]"} py-12 md:px-8`}
>
  <div
    class="border-2 border-brand-solid-light dark:border-brand-solid-dark w-full py-3 mx-4 md:mx-0 px-6 rounded-3xl flex"
  >
    <div class="">
      <Button
        class="w-full !px-0 bg-transparentv dark:bg-transparent hover:bg-transparent text-tertiary-light dark:text-tertiary-dark font-bold focus:ring-0"
        >{dropdownCurrentlySelected()}
        &nbsp;
        <ChevronDownOutline class="w-6 h-6 text-tertiary-light dark:text-tertiary-dark" />
      </Button>
      <Dropdown bind:open={dropdownOpen}>
        {#each dropdownItem as item}
          <DropdownItem on:click={handleDropdown} value={item.value}>{item.name}</DropdownItem>
        {/each}
      </Dropdown>
    </div>
    <div class="grow">
      <textarea
        bind:this={textarea}
        bind:value={messageInput}
        on:keydown={handleKeyDown}
        disabled={isLoading || disableInput}
        on:focus={() => (showSendMessage = true)}
        on:focusout={() => (showSendMessage = false)}
        id="chat"
        rows="1"
        class="block p-2.5 w-full text-base text-primary-light dark:text-primary-dark bg-white dark:bg-transparent rounded-lg border-0 focus:ring-transparent dark:placeholder-gray-400 resize-none"
        placeholder="Send a message..."
      ></textarea>
    </div>
    <div class="">
      <Button
        on:click={onclickSend}
        class="w-max p-3 bg-fg-brand-primary-light dark:bg-fg-brand-primary-dark rounded-full"
      >
        <img src={SendIcon} alt="send" class="w-6 h-6" />
      </Button>
    </div>
  </div>
</div>
