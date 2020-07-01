import debounce from "lodash.debounce";

import { checkResponse, createElement } from "./utils";
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

function councilMemberNameToSlug(name: string): string {
  const slug = name
    .replace(/[^\w\s]/g, "")
    .toLowerCase()
    .replace(/ /g, "-");
  return `https://citycouncil.atlantaga.gov/council-members/${slug}`;
}

function selectCandidate(candidate) {
  getRecord(candidate).then(([record]) => {
    const input = <HTMLInputElement>document.getElementById("address-input");
    input.value = candidate.address;
    const list = document.getElementById("candidates");
    list.innerHTML = "";

    const profile = document.getElementById("selected-candidate");
    profile.innerHTML = "";

    profile.appendChild(
      createElement({
        tagName: "div",
        className: "w-1/3 p-5 bg-green-100 text-center text-xl",
        text: `${record.COUNCIL_DIST}`,
      })
    );

    const link = {
      tagName: "a",
      className: "underline text-blue-500",
      text: `${record.COUNCIL_MEMBER}`,
      attributes: {
        href: councilMemberNameToSlug(record.COUNCIL_MEMBER),
      },
    };

    profile.appendChild(
      createElement({
        tagName: "div",
        className: "w-2/3 bg-green-100 py-5 text-left",
        childs: [link],
      })
    );
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
        "cursor-pointer",
        "truncate"
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
