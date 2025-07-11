export const errMsg = (e: any) =>
  e?.response?.data?.message ||
  e?.response?.data ||
  e?.message ||
  "Something went wrong";
