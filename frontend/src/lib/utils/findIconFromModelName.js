export const findIconFromModelName = (modelName) => {
  let icon = "logo";
  const lowerCaseModelName = modelName.toLowerCase();

  // console.log(modelName);

  if (lowerCaseModelName.includes("llama")) {
    icon = "meta";
  } else if (lowerCaseModelName.includes("deepseek")) {
    icon = "deepseek";
  } else if (lowerCaseModelName.includes("chatgpt")) {
    icon = "openai";
  } else if (lowerCaseModelName.includes("granite")) {
    icon = "granite";
  } else if (lowerCaseModelName.includes("mistral")) {
    icon = "mistral";
  } else if (lowerCaseModelName.includes("gpt oss")) {
    icon = "openAi";
  }

  return icon;
};
