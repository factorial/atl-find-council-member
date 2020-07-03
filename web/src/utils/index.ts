export function checkResponse(response: Response): Response {
  if (!response.ok) {
    throw new Error("Network response was not ok");
  }
  return response;
}

export function clearElement(el) {
  while (el.firstChild) {
    el.removeChild(el.firstChild);
  }
}
