export function checkResponse(response: Response): Response {
  if (!response.ok) {
    throw new Error("Network response was not ok");
  }
  return response;
}
