import HomeLineIcon from "../icon/home-line.svelte";
import BrandIcon from "../icon/logomark.svelte";
import SettingIcon from "../icon/settings-01.svelte";
import HexagonIcon from "../icon/hexagon-02.svelte";
import FolderIcon from "../icon/folder.svelte";
import UsersIcon from "../icon/users-03.svelte";
import MessageChatIcon from "../icon/message-chat-circle.svelte";
import BlockchainIcon from "../icon/blockchain.svelte";
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
import ChevronUpIcon from "../icon/chevron-up.svelte";
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
import UserUpIcon from "../icon/user-up.svelte";
import EyeIcon from "../icon/eye.svelte";
import Presentation from "../icon/presentation-chart-02.svelte";
import Lightbulb from "../icon/lightbulb-04.svelte";
import Marketplace from "../icon/marketplace-01.svelte";
import MessageChatSquare from "../icon/message-chat-square.svelte";
import MessageHeartSquare from "../icon/message-heart-square.svelte";
import MessageSmileSquare from "../icon/message-smile-square.svelte";
import Star04Icon from "../icon/star-04.svelte";
import PenToolIcon from "../icon/pen-tool-02.svelte";
import Image03Icon from "../icon/image-03.svelte";
import VideoRecorderIcon from "../icon/video-recorder.svelte";
import Recording02Icon from "../icon/recording-02.svelte";
import AiDocumentIcon from "../icon/ai-document-01.svelte";
import FileLockIcon from "../icon/file-lock-02.svelte";
import Settings03Icon from "../icon/settings-03.svelte";
import MessagePlusSquare from "../icon/message-plus-square.svelte";
import GraduationHat01 from "../icon/graduation-hat-01.svelte";
import GraduationHat02 from "../icon/graduation-hat-02.svelte";
import CodingIcon from "../icon/coding.svelte";
import Globe01 from "../icon/globe-01.svelte";
import DashboardIcon from "../icon/dashboard-01.svelte";
import Dataflow04Icon from "../icon/dataflow-04.svelte";
import MetaIcon from "../icon/meta.svelte";
import GeminiIcon from "../icon/gemini.svelte";
import DeepseekIcon from "../icon/deepseek.svelte";
import PhoneCallIcon from "../icon/phone-call-02.svelte";
import ArrowLeftIcon from "../icon/arrow-left.svelte";
import RefreshCw05Icon from "../icon/refresh-cw-05.svelte";
import UploadCloud02Icon from "../icon/upload-cloud-02.svelte";
import OpenAiIcon from "../icon/open-ai.svelte";
import ArrowUpRightIcon from "../icon/arrow-up-right.svelte";
import Check2 from "../icon/check2.svelte";
import Lightbulb04 from "../icon/lightbulb-04.svelte";
import PenToolPlusIcon from "../icon/pen-tool-plus.svelte";
import Code01Icon from "../icon/code-01.svelte";
import CodeSquare01Icon from "../icon/code-square-01.svelte";
import Image01 from "../icon/image-01.svelte";
import InfoCircle from "../icon/info-circle.svelte";
import FolderDownload from "../icon/folder-download.svelte";
import GoogleDriveIcon from "../icon/google-drive.svelte";
import TrashIcon from "../icon/trash-01.svelte";
import UploadCloud01 from "../icon/upload-cloud-01.svelte";
import VerifiedTickIcon from "../icon/verified-tick.svelte";
import OnedriveIcon from "../icon/onedrive.svelte";
import Trash03Icon from "../icon/trash-03.svelte";
import LineChartUp01 from "../icon/line-chart-up-01.svelte";
import ThumbsDown from "../icon/thumbs-down.svelte";
import Send03 from "../icon/send-03.svelte";

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
    case "chevronUp":
      return ChevronUpIcon;
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
    case "userUp":
      return UserUpIcon;
    case "eye":
      return EyeIcon;
    case "presentation":
      return Presentation;
    case "lightbulb":
      return Lightbulb;
    case "marketplace":
      return Marketplace;
    case "messageChatSquare":
      return MessageChatSquare;
    case "messageHeartSquare":
      return MessageHeartSquare;
    case "messageSmileSquare":
      return MessageSmileSquare;
    case "star04":
      return Star04Icon;
    case "penTool":
      return PenToolIcon;
    case "image03":
      return Image03Icon;
    case "videoRecorder":
      return VideoRecorderIcon;
    case "recording02":
      return Recording02Icon;
    case "aiDocument":
      return AiDocumentIcon;
    case "fileLock":
      return FileLockIcon;
    case "settings03":
      return Settings03Icon;
    case "messagePlusSquare":
      return MessagePlusSquare;
    case "graduationHat01":
      return GraduationHat01;
    case "graduationHat02":
      return GraduationHat02;
    case "coding":
      return CodingIcon;
    case "globe01":
      return Globe01;
    case "dashboard01":
      return DashboardIcon;
    case "dataflow04":
      return Dataflow04Icon;
    case "meta":
      return MetaIcon;
    case "gemini":
      return GeminiIcon;
    case "deepseek":
      return DeepseekIcon;
    case "phoneCall":
      return PhoneCallIcon;
    case "arrowLeft":
      return ArrowLeftIcon;
    case "refreshCw05":
      return RefreshCw05Icon;
    case "uploadCloud02":
      return UploadCloud02Icon;
    case "openAi":
      return OpenAiIcon;
    case "arrowUpRight":
      return ArrowUpRightIcon;
    case "check2":
      return Check2;
    case "lightbulb04":
      return Lightbulb04;
    case "penToolPlus":
      return PenToolPlusIcon;
    case "code01":
      return Code01Icon;
    case "codeSquare01":
      return CodeSquare01Icon;
    case "image01":
      return Image01;
    case "infoCircle":
      return InfoCircle;
    case "folderDownload":
      return FolderDownload;
    case "googleDrive":
      return GoogleDriveIcon;
    case "trash":
      return TrashIcon;
    case "uploadCloud01":
      return UploadCloud01;
    case "verifiedTick":
      return VerifiedTickIcon;
    case "onedrive":
      return OnedriveIcon;
    case "trash03":
      return Trash03Icon;
    case "lineChartUp01":
      return LineChartUp01;
    case "thumbsDown":
      return ThumbsDown;
    case "send03":
      return Send03;
    default:
      return "";
  }
};

export default getIcons;
