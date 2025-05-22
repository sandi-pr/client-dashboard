import { base64toBlob } from "./base64toBlob";

export const downloadFileBase64 = (base64Data, fileName, ext) => {
  const fileTypeMapping = {
    ".pdf": "application/pdf",
    ".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    ".doc": "application/msword",
    ".txt": "text/plain",
    ".jpg": "image/jpeg",
    ".jpeg": "image/jpeg",
    ".png": "image/png",
    ".gif": "image/gif",
    ".bmp": "image/bmp",
    ".webp": "image/webp",
    ".mp3": "audio/mpeg",
    ".wav": "audio/wav",
    ".ogg": "audio/ogg",
    ".mp4": "video/mp4",
    ".webm": "video/webm",
    ".avi": "video/x-msvideo",
    ".mov": "video/quicktime",
  };

  // Normalize the extension to lowercase for case-insensitive matching
  const normalizedExt = ext.toLowerCase();
  const fileType = fileTypeMapping[normalizedExt] || "application/octet-stream"; // Default to binary stream if unknown
  const a = document.createElement("a");
  a.href = URL.createObjectURL(base64toBlob(base64Data, fileType));
  a.download = fileName;
  a.click();
};

export const downloadFileBlob = (blob, fileName) => {
  const a = document.createElement("a");
  a.href = URL.createObjectURL(blob);
  a.download = fileName;
  a.click();
};
