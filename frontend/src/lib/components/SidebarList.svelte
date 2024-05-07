<script>
  import { page } from "$app/stores";

  import {
    Sidebar,
    SidebarGroup,
    SidebarItem,
    SidebarWrapper,
    SidebarBrand,
    Search,
    Button,
    Avatar,
    Span,
    P,
  } from "flowbite-svelte";
  import { ArrowRightToBracketOutline } from "flowbite-svelte-icons";
  import HomeLineIcon from "$lib/icon/home-line.svelte";
  import LogoIcon from "$lib/icon/Logomark.svg";
  import SettingIcon from "$lib/icon/settings-01.svelte";
  import HexagonIcon from "$lib/icon/hexagon-02.svelte";
  import BlockchainIcon from "$lib/icon/blockchain.svelte";
  import { getModule } from "$lib/utils/getModule";
  import { onMount } from "svelte";
  import { sineIn } from "svelte/easing";
  import Me from "$lib/utils/me";

  let hidden2 = true;
  let spanClass = "flex-1 ms-3 whitespace-nowrap";
  let transitionParams = {
    x: -320,
    duration: 200,
    easing: sineIn,
  };

  export let name;
  export let email;

  export let isHidden = true;

  let site = {
    name: "Crossroads",
    href: "/",
    img: LogoIcon,
    class: "pl-0",
  };

  let sidebarItems = [
    { name: "Home", href: "/", icon: HomeLineIcon },
    { name: "Apps", href: "#", icon: HexagonIcon },
  ];

  const sidebarBuilder = async () => {
    let moduleClient = await getModule();
    const sidebarHalf = moduleClient.map((item) => {
      return {
        name: item.name,
        href: item.href,
        icon: item.icon,
      };
    });

    sidebarHalf.push({ name: "Blockchain", href: "/blockchain", icon: BlockchainIcon });

    return sidebarHalf;
  };

  $: activeUrl = $page.url.pathname;

  // let name = "";
  // let email = "";

  // if ($Me === null) {
  //   name = data?.data?.user?.fullname;
  //   email = data?.data?.user?.email;
  // } else {
  //   name = $Me?.user?.fullname;
  //   email = $Me?.user?.email;
  // }

  onMount(async () => {
    const sidebarAddition = await sidebarBuilder();
    sidebarItems = [...sidebarItems, ...sidebarAddition];
  });
</script>

<Sidebar
  {activeUrl}
  class={` ${
    isHidden ? "hidden h-screen" : "fixed h-full"
  } lg:border-r border-hc-graylight-200 dark:border-hc-graydark-800 lg:fixed lg:inset-y-0 lg:z-10 lg:flex w-80`}
>
  <SidebarWrapper
    class="h-full w-full lg:pt-6 lg:px-4 flex flex-col justify-between bg-white dark:bg-hc-graydark-950"
  >
    <SidebarGroup>
      <SidebarBrand {site} aClass="font-jakarta flex items-center ps-0 mb-5 text-gray-900" />
      <Search class="bg-white dark:bg-hc-graydark-950 mb-6" />
      {#each sidebarItems as item}
        <SidebarItem
          label={item.name}
          href={item.href}
          nonActiveClass="flex items-center px-3 py-2 text-base font-semibold hover:no-underline text-secondary-light rounded-lg dark:text-secondary-dark hover:bg-gray-50 dark:hover:bg-hc-graydark-800"
          activeClass="flex items-center px-3 py-2 text-base hover:no-underline font-semibold text-secondary-light bg-gray-50 dark:bg-hc-graydark-800 rounded-lg dark:text-secondary-dark hover:bg-gray-50 dark:hover:bg-hc-graydark-800"
        >
          <svelte:fragment slot="icon">
            <svelte:component
              this={item.icon}
              svgClass="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark"
            />
          </svelte:fragment>
        </SidebarItem>
      {/each}
    </SidebarGroup>
    <div>
      <SidebarGroup>
        <SidebarItem
          label="Super Admin Setting"
          href="/admin_setting"
          active={activeUrl === "/admin_setting"}
          nonActiveClass="flex items-center px-3 py-2 text-base font-normal hover:no-underline text-gray-900 rounded-lg dark:text-hc-graydark-300 hover:bg-gray-50 dark:hover:bg-hc-graydark-800"
          activeClass="flex items-center px-3 py-2 text-base hover:no-underline font-normal text-gray-900 bg-gray-50 dark:bg-hc-graydark-800 rounded-lg dark:text-hc-graydark-300 hover:bg-gray-50 dark:hover:bg-hc-graydark-800"
        >
          <svelte:fragment slot="icon">
            <svelte:component
              this={SettingIcon}
              svgClass="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark"
            />
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem
          label="Settings"
          href="/settings"
          active={activeUrl === "/settings"}
          nonActiveClass="flex items-center px-3 py-2 text-base font-normal hover:no-underline text-gray-900 rounded-lg dark:text-hc-graydark-300 hover:bg-gray-50 dark:hover:bg-hc-graydark-800"
          activeClass="flex items-center px-3 py-2 text-base hover:no-underline font-normal text-gray-900 bg-gray-50 dark:bg-hc-graydark-800 rounded-lg dark:text-hc-graydark-300 hover:bg-gray-50 dark:hover:bg-hc-graydark-800"
        >
          <svelte:fragment slot="icon">
            <svelte:component
              this={SettingIcon}
              svgClass="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark"
            />
          </svelte:fragment>
        </SidebarItem>
      </SidebarGroup>
      <SidebarGroup border>
        <SidebarItem
          label=""
          nonActiveClass="flex items-center pl-3 py-2 text-base font-normal hover:no-underline text-gray-900 rounded-lg dark:text-hc-graydark-300 hover:bg-gray-50 dark:hover:bg-hc-graydark-800"
          activeClass="flex items-center pl-3 py-2 text-base hover:no-underline font-normal text-gray-900 bg-gray-50 dark:bg-hc-graydark-800 rounded-lg dark:text-hc-graydark-300 hover:bg-gray-50 dark:hover:bg-hc-graydark-800"
        >
          <svelte:fragment slot="icon">
            <Avatar />
            <div class="flex grow h-full items-center ml-4 justify-between">
              <div>
                <Span
                  class="font-jakarta text-slate-700 dark:text-hc-graydark-300 text-sm leading-tight"
                  >{name}</Span
                >
                <P
                  class="font-jakarta text-slate-600 dark:text-hc-graydark-400 text-sm leading-tight"
                  >{email}</P
                >
              </div>
              <Button
                href="/logout"
                color="light"
                class="focus-within:ring-0 border-0 p-2 bg-transparent dark:bg-transparent"
              >
                <ArrowRightToBracketOutline
                  class="w-5 h-5 text-fg-quaternary-light dark:text-fg-quaternary-dark transition duration-75 group-hover:text-gray-900 dark:group-hover:text-white"
                />
              </Button>
            </div>
          </svelte:fragment>
        </SidebarItem>
      </SidebarGroup>
    </div>
  </SidebarWrapper>
</Sidebar>
