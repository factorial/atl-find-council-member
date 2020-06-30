import { checkResponse } from "../utils";
import L from "leaflet";
import "leaflet/dist/leaflet.css";

function filterCityCouncilDistrictFeatures(geojson) {
  return geojson.features.filter(
    (feature) => feature.properties.GEOTYPE === "City Council District"
  );
}

function onEachFeature(feature, layer) {
  layer.on({
    click: (evt) => console.log(feature.properties.NAME),
  });
}

function mapCityCouncilDistrictFeatures(map, features) {
  features.forEach((feature) =>
    L.geoJSON(feature, { onEachFeature: onEachFeature }).addTo(map)
  );
}

export function attachMap(elementId: string) {
  const map = L.map(elementId).setView([33.749, -84.388], 13);

  var osm_mapnik = L.tileLayer(
    "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",
    {
      maxZoom: 19,
      attribution:
        '&copy; OSM Mapnik <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
    }
  ).addTo(map);
  const request = new Request(
    `https://opendata.arcgis.com/datasets/5ce01aea8d4046fe8659a8e25958c2bb_2.geojson`
  );

  fetch(request)
    .then((resp) => checkResponse(resp))
    .then((resp) => resp.json())
    .then((json) => filterCityCouncilDistrictFeatures(json))
    .then((features) => mapCityCouncilDistrictFeatures(map, features));
}
