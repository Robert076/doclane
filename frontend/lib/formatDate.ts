export function formatDate(date: string) {
  return new Date(date).toLocaleDateString("us-US", {
    year: "numeric",
    month: "long",
    day: "2-digit",
  });
}
