// Function to encode URL parameters
export const encodeUrlParams = (statusCode, prevUrl) => {
  const encodedPrevUrl = encodeURIComponent(prevUrl);
  const newUrl = `/errors?statusCode=${statusCode}&prevUrl=${encodedPrevUrl}`;
  window.history.pushState({}, "", newUrl);
};

// function to decode URL parameters
export const decodeUrlParams = (url) => {
  const urlParams = new URLSearchParams(url);
  const statusCode = urlParams.get("statusCode");
  const prevUrl = urlParams.get("prevUrl");
  return { statusCode, prevUrl };
};
