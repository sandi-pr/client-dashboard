export const findIconFromModelName = (modelName) => {
  let icon = "gemini";
  const lowerCaseModelName = modelName.toLowerCase();
  console.log(lowerCaseModelName);

  if (/llama/.test(lowerCaseModelName)) {
    icon = "meta";
  } else if (/deepseek/.test(lowerCaseModelName)) {
    icon = "deepseek";
  } else if (/chatgpt/.test(lowerCaseModelName)) {
    icon = "openai";
  }

  return icon;
};
