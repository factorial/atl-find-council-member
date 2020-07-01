import debounce from "lodash.debounce";

import { checkResponse } from "./utils";
import { fetchAddressCandidates, fetchCandidate } from "./mapatlapi";

function searchAddress(address) {
  return fetchAddressCandidates(address)
    .then((resp) => checkResponse(resp))
    .then((resp) => resp.json())
    .then(({ candidates }) => candidates);
}

function getRecord(candidate) {
  return fetchCandidate(candidate)
    .then((resp) => checkResponse(resp))
    .then((resp) => resp.json());
}

function selectCandidate(candidate) {
  getRecord(candidate).then(([record]) => {
    const input = <HTMLInputElement>document.getElementById("address-input");
    input.value = candidate.address;
    const list = document.getElementById("candidates");
    list.innerHTML = `District ${record.COUNCIL_DIST} : ${record.COUNCIL_MEMBER}`;
  });
}

function displayCandidates(candidates) {
  const list = document.getElementById("candidates");
  list.innerHTML = "";

  if (candidates) {
    const elements = candidates.map((candidate) => {
      const element = document.createElement("div");
      element.classList.add(
        "text-gray-500",
        "hover:text-gray-600",
        "cursor-pointer"
      );
      element.innerHTML = `${candidate.address}`;
      element.addEventListener("click", () => selectCandidate(candidate));
      return element;
    });

    elements.forEach((el) => list.appendChild(el));
  }
}

function run() {
  const debouncedSearchAddress = debounce(searchAddress, 250);

  const addressInputElement = <HTMLInputElement>(
    document.getElementById("address-input")
  );

  addressInputElement.addEventListener("input", (evt) => {
    const target = <HTMLInputElement>evt.currentTarget;
    searchAddress(target.value).then((candidates) =>
      displayCandidates(candidates)
    );
  });
}

run();
