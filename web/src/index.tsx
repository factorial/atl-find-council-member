import debounce from "lodash.debounce";
import h from "hyperscript";
import { checkResponse, clearElement } from "./utils";
import {
  fetchAddressCandidates,
  fetchCandidate,
  fetchRepresentative,
} from "./mapatlapi";

async function searchAddress(address) {
  return fetchAddressCandidates(address)
    .then((resp) => checkResponse(resp))
    .then((resp) => resp.json())
    .then(({ candidates }) => candidates);
}

async function getRecord(candidate) {
  return fetchCandidate(candidate)
    .then((resp) => checkResponse(resp))
    .then((resp) => resp.json());
}

async function getRepresentative(district: number) {
  return fetchRepresentative(district)
    .then((resp) => checkResponse(resp))
    .then((resp) => resp.json());
}

function buildRepresentativeCard(representative) {
  return (
    <div className="py-3 flex justify-between">
      <div className="w-1/3 flex items-center justify-left">
        <div className="text-xl">{representative.District}</div>
      </div>
      <div className="w-1/3 flex items-center justify-center">
        <a className="text-xl" href={representative.Href}>
          {representative.Name}
        </a>
      </div>
      <div className="w-1/3 flex items-center justify-end">
        <img
          className="w-24"
          src={`https://citycouncil.atlantaga.gov${representative.Image}`}
        />
      </div>
    </div>
  );
}

function clearAddressInput() {
  const input = document.getElementById("address-input") as HTMLInputElement;
  input.value = "";
}

function clearCandidateList() {
  const list = document.getElementById("candidates");
  clearElement(list);
}

function showRepresentativeCard(card) {
  const profile = document.getElementById("selected-candidate");
  clearElement(profile);
  profile.appendChild(card);
}

async function selectCandidate(candidate) {
  const [record] = await getRecord(candidate);

  clearAddressInput();
  clearCandidateList();

  const representative = await getRepresentative(record.COUNCIL_DIST);
  const card = buildRepresentativeCard(representative);

  showRepresentativeCard(card);
}

function showCandidates(candidates) {
  const list = document.getElementById("candidates");
  clearElement(list);

  if (candidates) {
    const elements = candidates.map((candidate) => (
      <div
        className="cursor-pointer text-gray-600 hover:text-gray-700 truncate"
        onclick={() => selectCandidate(candidate)}
      >
        {candidate.address}
      </div>
    ));

    elements.forEach((el) => list.appendChild(el));
  }
}

function run() {
  const debouncedSearchAddress = debounce(searchAddress, 350);

  const addressInputElement = document.getElementById(
    "address-input"
  ) as HTMLInputElement;

  addressInputElement.addEventListener("input", async (evt) => {
    const target = evt.currentTarget as HTMLInputElement;
    const candidates = await debouncedSearchAddress(target.value);
    showCandidates(candidates);
  });
}

run();
