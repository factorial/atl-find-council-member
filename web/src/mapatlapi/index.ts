import { checkResponse } from "../utils";

export function fetchAddressCandidates(address) {
  const body = { address };
  const request = new Request(`https://mapatlapi.abrie.dev/address`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });

  return fetch(request);
}

export function fetchCandidate(candidate) {
  const {
    attributes: { Ref_ID },
  } = candidate;
  const body = { Ref_ID };
  const request = new Request(`https://mapatlapi.abrie.dev/record`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });
  return fetch(request);
}

export function fetchRepresentative(district: number) {
  const body = { district: `District ${district}` };
  const request = new Request(`https://mapatlapi.abrie.dev/council`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });

  return fetch(request);
}
