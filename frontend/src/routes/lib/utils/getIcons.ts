import HomeLineIcon from '$lib/icon/home-line.svelte';
import BrandIcon from '$lib/icon/Logomark.svg';
import SettingIcon from '$lib/icon/settings-01.svelte';
import HexagonIcon from '$lib/icon/hexagon-02.svelte';
import FolderIcon from "$lib/icon/folder.svelte";
import UsersIcon from "$lib/icon/users-03.svelte";
import MessageChatIcon from "$lib/icon/message-chat-circle.svelte";
import BlockchainIcon from '$lib/icon/blockchain.svelte';

let getIcons = (presetIcon: String) => {
  switch (presetIcon) {
    case "brand":
      return BrandIcon;
    case "setting":
      return SettingIcon;
    case "hexagon":
      return HexagonIcon;
    case "folder":
      return FolderIcon;
    case "users":
      return UsersIcon;
    case "messageChat":
      return MessageChatIcon;
    case "home":
      return HomeLineIcon;
    case "blockchain":
      return BlockchainIcon;
    default:
      return HomeLineIcon;
  }
}

export default getIcons;