export const truncateText = (text, maxWords, maxFirstWordLength = 14) => {
  const words = text.split(" ");

  // Check if the first word is too long, or if the second word is long
  if (words[0] && words[0].length > maxFirstWordLength) {
    return words[0].slice(0, maxFirstWordLength) + "...";
  } else if (words[1] && words[1].length > maxFirstWordLength) {
    return words[0] + " " + words[1].slice(0, maxFirstWordLength) + "...";
  }

  // Limit to maxWords and append "..." if necessary
  if (words.length > maxWords) {
    return words.slice(0, maxWords).join(" ") + " ...";
  }

  return words.join(" ");
};
