import debounce from "lodash.debounce";

import { checkResponse, createElement } from "./utils";
import {
  fetchAddressCandidates,
  fetchCandidate,
  fetchDistrict,
} from "./mapatlapi";

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

function getDistrict(district: number) {
  return fetchDistrict(district)
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

    getDistrict(record.COUNCIL_DIST).then((rep) => {
      const district = {
        tagName: "div",
        className: "w-full text-center text-xl",
        text: `${rep.District}`,
      };

      profile.appendChild(
        createElement({
          tagName: "div",
          className: "w-1/3 flex items-center",
          childs: [district],
        })
      );

      const link = {
        tagName: "a",
        className: "underline text-blue-500",
        text: `${rep.Name}`,
        attributes: {
          href: rep.Href,
        },
      };

      const name = {
        tagName: "div",
        className: "w-full text-center",
        childs: [link],
      };

      profile.appendChild(
        createElement({
          tagName: "div",
          className: "w-1/3 text-left flex items-center",
          childs: [name],
        })
      );

      const image = {
        tagName: "img",
        className: "w-12 mx-auto",
        attributes: {
          src: `https://citycouncil.atlantaga.gov/${rep.Image}`,
        },
      };

      profile.appendChild(
        createElement({
          tagName: "div",
          className: "w-1/3 text-left",
          childs: [image],
        })
      );
    });
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
