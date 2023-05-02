import http from "k6/http";

export let options = {
  vus: 100,
  duration: "1s",
};

export default function () {
  let month = Math.floor(Math.random() * 12) + 1;
  let day = Math.floor(Math.random() * 28) + 1;
  http.get(
    `http://host.docker.internal/api/v1/events?month=${month}&day=${day}`
  );
}
