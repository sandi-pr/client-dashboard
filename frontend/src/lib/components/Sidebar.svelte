<script>
  import { Drawer, Button, CloseButton, P } from "flowbite-svelte";
  import SidebarList from "$lib/components/SidebarList.svelte";
  import LogoIcon from "$lib/icon/Logomark.svg";
  import MenuIcon from "$lib/icon/menu-02.svelte";
  import { sineIn } from "svelte/easing";
  let hidden2 = true;
  let transitionParams = {
    x: -320,
    duration: 200,
    easing: sineIn,
  };

  export let name;
  export let email;
</script>

<Drawer transitionType="fly" {transitionParams} bind:hidden={hidden2} class="relative !p-0">
  <div class="fixed z-10 top-3 left-[280px]">
    <CloseButton on:click={() => (hidden2 = true)} class=" dark:text-white" />
  </div>
  <SidebarList {name} {email} isHidden={false} />
</Drawer>

<SidebarList {name} {email} />

<div class="lg:hidden pl-4 pr-2 py-3 justify-between items-center flex">
  <div class="justify-between items-center gap-2 flex">
    <div class="w-8 h-8 justify-center items-center flex">
      <img alt="logo" src={LogoIcon} class="w-8 h-8" />
    </div>
    <P size="lg" weight="bold" class="leading-7 text-fg-primary-light dark:text-fg-primary-dark"
      >Crossroads</P
    >
  </div>
  <div class="p-2 rounded-lg justify-center items-center flex">
    <Button on:click={() => (hidden2 = false)} color="light">
      <MenuIcon class="w-6 h-6 text-fg-secondary-light dark:text-fg-secondary-dark" />
    </Button>
  </div>
</div>
