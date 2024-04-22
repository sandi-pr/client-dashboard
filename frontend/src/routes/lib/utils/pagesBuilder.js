export const pagesBuilder = (pagination) => {
  let pages = [];
  for (let i = 1; i <= pagination.pageCount; i++) {
    pages.push({ name: i });
  }
  return pages;
};
