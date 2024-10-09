import HomeLineIcon from '../icon/home-line.svelte';
import BrandIcon from "../icon/logomark.svelte";
import SettingIcon from '../icon/settings-01.svelte';
import HexagonIcon from '../icon/hexagon-02.svelte';
import FolderIcon from "../icon/folder.svelte";
import UsersIcon from "../icon/users-03.svelte";
import MessageChatIcon from "../icon/message-chat-circle.svelte";
import BlockchainIcon from '../icon/blockchain.svelte';
import BellIcon from "../icon/bell-01.svelte";
import FileShieldIcon from "../icon/file-shield-02.svelte";
import FIleAttachmentIcon from "../icon/file-attachment-04.svelte";
import UserSquareIcon from "../icon/user-square.svelte";
import UserPlusIcon from "../icon/user-plus-02.svelte";
import PackagePlusIcon from "../icon/package-plus.svelte";
import PackageSearchIcon from "../icon/package-search.svelte";
import FolderPlusIcon from "../icon/folder-plus.svelte";
import FilePlusIcon from "../icon/file-plus-02.svelte";
import BoxIcon from "../icon/box.svelte";
import BuildingIcon from "../icon/building-08.svelte";
import BookOpenIcon from "../icon/book-open-02.svelte";
import DistributeVertIcon from "../icon/distribute-spacing-vertical.svelte";
import ChevronDownIcon from "../icon/chevron-down.svelte";
import DataIcon from "../icon/data.svelte";
import MessagePlus from "../icon/message-plus-square.svelte";
import BookClosedIcon from "../icon/book-closed.svelte";
import BriefcaseIcon from "../icon/briefcase-02.svelte";
import UserIcon from "../icon/user-02.svelte";
import FaceHappyIcon from "../icon/face-happy.svelte";
import PackageIcon from "../icon/package.svelte";
import ReceiptIcon from "../icon/receipt.svelte";
import CalendarIcon from "../icon/calendar.svelte";
import StarsIcon from "../icon/stars-01.svelte";
import ArchiveIcon from "../icon/archive.svelte";
import FileCheckIcon from "../icon/file-check-02.svelte";

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
    case "bell":
      return BellIcon;
    case "fileShield":
      return FileShieldIcon;
    case "fileAttachment":
      return FIleAttachmentIcon;
    case "userSquare":
      return UserSquareIcon;
    case "userPlus":
      return UserPlusIcon;
    case "packagePlus":
      return PackagePlusIcon;
    case "packageSearch":
      return PackageSearchIcon;
    case "folderPlus":
      return FolderPlusIcon;
    case "filePlus":
      return FilePlusIcon;
    case "box":
      return BoxIcon;
    case "building":
      return BuildingIcon;
    case "bookOpen":
      return BookOpenIcon;
    case "distributeVert":
      return DistributeVertIcon;
    case "chevronDown":
      return ChevronDownIcon;
    case "data":
      return DataIcon;
    case "messagePlus":
      return MessagePlus;
    case "bookClosed":
      return BookClosedIcon;
    case "briefcase":
      return BriefcaseIcon;
    case "user":
      return UserIcon;
    case "faceHappy":
      return FaceHappyIcon;
    case "package":
      return PackageIcon;
    case "receipt":
      return ReceiptIcon;
    case "calendar":
      return CalendarIcon;
    case "stars":
      return StarsIcon;
    case "archive":
      return ArchiveIcon;
    case "fileCheck":
      return FileCheckIcon;
    default:
      return "";
  }
}

export default getIcons;