// import { socket } from "@/socket";
import { socket } from "@/socket";
import type { Host, NmapRun } from "./NmapRun";

let svg: SVGSVGElement;
let wrapper: HTMLDivElement;
let isDragging = false;
let startX = 0;
let startY = 0;
let translateX = 0;
let translateY = 0;
let scale = 1;

let data: NmapRun;

export async function start() {
  const existingSVGs = document.querySelectorAll("svg.canvas");
  existingSVGs.forEach((svg) => {
    svg.remove();
  });

  svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
  svg.setAttribute("class", "canvas");
  svg.style.position = "relative";
  svg.style.cursor = "default"

  document.body.style.overflow = "hidden"

  wrapper = document.querySelector("div.canvas-wrapper") as HTMLDivElement;
  wrapper.appendChild(svg);
  wrapper.style.width = `${window.innerWidth}px`
  wrapper.style.height = `${window.innerHeight - (document.querySelector(".navbar") as HTMLElement).clientHeight}px`

  wrapper.addEventListener("mousedown", (e) => {
    console.log(e.target)
    if (e.target instanceof SVGRectElement) {
      const index = e.target.dataset.index || "0";
      if (index !== undefined) {
        displayDeviceInfo(data.hosts.find(host => host.tcp_sequence.values == index) as Host);
      }
    }

    isDragging = true;
    startX = e.clientX;
    startY = e.clientY;
  });

  wrapper.addEventListener("mousemove", (e) => {
    if (isDragging) {
      const dx = e.clientX - startX;
      const dy = e.clientY - startY;
      translateX += dx;
      translateY += dy;
      startX = e.clientX;
      startY = e.clientY;
      updateSvg();

      window.getSelection()?.removeAllRanges();
    }
  });

  wrapper.addEventListener("mouseup", () => {
    isDragging = false;
  });

  wrapper.addEventListener("mouseleave", () => {
    isDragging = false;
  });

  wrapper.addEventListener("wheel", (e) => {
    const delta = e.deltaY > 0 ? -0.1 : 0.1;
    const scaleFactor = 1 + delta;

    scale *= scaleFactor;
    updateSvg();
  });

  closeSidePanel();
  // setTimeout(() => socket.emit("defaultScan"), 1000)
}

function update() {
  updateSvg();
  requestAnimationFrame(update);
}

function updateSvg() {
  if (!svg) return;

  while (svg.firstChild) {
    svg.firstChild.remove();
  }

  svg.style.top = `${translateY}px`
  svg.style.left = `${translateX}px`
  svg.style.transform = `scale(${scale})`

  // const subnetMap = new Map<string, Device[]>();
  // data.forEach((device, i) => {
  //   const subnet = device.ip.split(".").slice(0, 3).join(".");
  //   if (!subnetMap.has(subnet)) {
  //     subnetMap.set(subnet, []);
  //   }
  //   subnetMap.get(subnet)?.push({ ...device, index: i });
  // });

  // subnetMap.forEach((devices) => {
  //   devices.sort((a, b) => {
  //     const aLast = parseInt(a.ip.split(".")[3]);
  //     const bLast = parseInt(b.ip.split(".")[3]);
  //     return aLast - bLast;
  //   });
  // });

  let width = 1000;
  let height = 40;

  // subnetMap.forEach((devices, subnet) => {
  // const subnetHeading = document.createElementNS("http://www.w3.org/2000/svg", "text");
  // subnetHeading.setAttribute("x", `80`);
  // subnetHeading.setAttribute("y", `${height}`);
  // subnetHeading.setAttribute("font-family", "monospace");
  // subnetHeading.setAttribute("font-size", "30px");
  // subnetHeading.setAttribute("fill", "accentColor");
  // subnetHeading.textContent = "subnet";
  // svg.appendChild(subnetHeading);

  // const subnetSubheading = document.createElementNS("http://www.w3.org/2000/svg", "text");
  // subnetSubheading.setAttribute("x", `80`);
  // subnetSubheading.setAttribute("y", `${height + 20}`);
  // subnetSubheading.setAttribute("font-family", "monospace");
  // subnetSubheading.setAttribute("font-size", "20px");
  // subnetSubheading.setAttribute("fill", "currentColor");
  // subnetSubheading.textContent = `${"devices[0].osNmap"} (${"devices.length"} devices)`;
  // svg.appendChild(subnetSubheading);

  // const routerRect = document.createElementNS("http://www.w3.org/2000/svg", "rect");
  // routerRect.setAttribute("x", `${10}`);
  // routerRect.setAttribute("y", `${height - 30}`);
  // routerRect.setAttribute("width", `${60}`);
  // routerRect.setAttribute("height", `${60}`);
  // routerRect.setAttribute("fill", "transparent");
  // routerRect.setAttribute("stroke", "currentColor");
  // routerRect.setAttribute("stroke-width", "2");
  // routerRect.setAttribute("rx", "10");
  // routerRect.setAttribute("ry", "10");
  // routerRect.dataset.index = `${"devices[0].index"}`;
  // svg.appendChild(routerRect);

  // height += 40;

  // if (devices[0].ip.endsWith(".1")) {
  //   devices.shift();
  // }

  data.hosts.forEach((host, i) => {
    const rect = document.createElementNS("http://www.w3.org/2000/svg", "rect");
    rect.setAttribute("x", `0`);
    rect.setAttribute("y", `${20 + height + 40 * i}`);
    rect.setAttribute("width", `${10}`);
    rect.setAttribute("height", `${10}`);
    rect.setAttribute("fill", `${false ? "red" : "black"}`);
    rect.dataset.index = `${host.tcp_sequence.values}`;
    svg.appendChild(rect);

    const text = document.createElementNS("http://www.w3.org/2000/svg", "text");
    text.setAttribute("x", `0`);
    text.setAttribute("y", `${20 + height + 40 * i}`);
    text.setAttribute("fill", "black");
    text.textContent = host.os.os_matches[0].name + " (" + host.os.os_matches[0].accuracy + "% accuracy)" || "Not detected";
    svg.appendChild(text);
  });

  height += 50 * data.hosts.length;
  
  // });
  
  svg.style.width = `${width}px`;
  svg.style.height = `${height}px`;
}

socket.onmessage = (msg) => {
  data = JSON.parse(msg.data) as NmapRun;
  console.log(data);
  updateSvg();
}

function openSidePanel() {
  const sidePanel = document.querySelector('.side-panel.device-info') as HTMLElement;
  const navbar = document.querySelector(".navbar") as HTMLElement

  const closeButton = document.querySelector("button.device-info-close") as HTMLButtonElement;
  closeButton.addEventListener("click", () => { closeSidePanel() });

  sidePanel.style.width = wrapper.clientWidth * 0.3 + "px";
  sidePanel.style.height = wrapper.clientHeight * 0.9 + "px";
  sidePanel.style.top = (wrapper.clientHeight * 0.05) + navbar.clientHeight + "px";
  sidePanel.style.right = "0px";

  closeButton.style.top = wrapper.clientHeight * 0.05 + navbar.clientHeight - closeButton.clientHeight / 2 + "px";
  closeButton.style.right = sidePanel.clientWidth - closeButton.clientWidth / 2 + "px";

  sidePanel.innerHTML = "";
}

function closeSidePanel() {
  const sidePanel = document.querySelector('.side-panel.device-info') as HTMLElement;
  const closeButton = document.querySelector("button.device-info-close") as HTMLButtonElement;
  sidePanel.innerHTML = "";

  sidePanel.style.transition = "right 0.1s ease";
  closeButton.style.transition = "right 0.1s ease";

  sidePanel.style.right = - sidePanel.clientWidth - 10 - (closeButton.clientWidth / 2) + "px";
  closeButton.style.right = - closeButton.clientWidth - 10 + "px";

}

function displayDeviceInfo(device: Host) {
  openSidePanel();
  const sidePanel = document.querySelector('.side-panel.device-info') as HTMLElement;
  const navbar = document.querySelector(".navbar") as HTMLElement

  const addressStrings: string[] = device.addresses.map((address) => address.addr);
  const commaSeparatedAddresses: string = addressStrings.join(', ');
  const heading = document.createElement("h1");
  heading.textContent = commaSeparatedAddresses;
  heading.classList.add("text-accent", "text-center", "text-2xl", "font-bold", "pt-4");
  sidePanel.appendChild(heading);

  const subHeading = document.createElement("h2");
  subHeading.textContent = device.os.os_matches[0].name || "OS not detected";
  subHeading.classList.add("text-center", "text-xl", "font-semibold");
  sidePanel.appendChild(subHeading);

  const hr = document.createElement("hr");
  hr.classList.add("my-4", "w-2/3", "mx-auto");
  sidePanel.appendChild(hr);

  const table = document.createElement("table");
  table.classList.add("mx-auto", "w-2/3", "table");
  sidePanel.appendChild(table);

  const tbody = document.createElement("tbody");
  table.appendChild(tbody);

  const tr1 = document.createElement("tr");
  tbody.appendChild(tr1);

  const td1 = document.createElement("td");
  td1.textContent = "MAC Address";
  td1.classList.add("font-semibold");
  tr1.appendChild(td1);

  const td2 = document.createElement("td");
  td2.textContent = device.addresses.filter(address => address.addr_type == "mac").join(", ") || "Unknown";
  tr1.appendChild(td2);

  const tr2 = document.createElement("tr");
  tbody.appendChild(tr2);

  const td3 = document.createElement("td");
  td3.textContent = "Status";
  td3.classList.add("font-semibold");
  tr2.appendChild(td3);

  const td4 = document.createElement("td");
  td4.textContent = `${device.status.state} (${device.status.reason})` || "Unknown";
  tr2.appendChild(td4);

  const tr3 = document.createElement("tr");
  tbody.appendChild(tr3);

  const td5 = document.createElement("td");
  td5.textContent = "Hostname";
  td5.classList.add("font-semibold");
  tr3.appendChild(td5);

  const td6 = document.createElement("td");
  td6.textContent = device.hostnames.map(hostname => hostname.name).join(", ") || "Unknown";
  tr3.appendChild(td6);

  const tr4 = document.createElement("tr");
  tbody.appendChild(tr4);

  const td7 = document.createElement("td");
  td7.textContent = "Uptime";
  td7.classList.add("font-semibold");
  tr4.appendChild(td7);

  const td8 = document.createElement("td");
  td8.textContent = `${device.uptime.last_boot} (${device.uptime.seconds}s)` || "Unknown";
  tr4.appendChild(td8);

  const hr2 = document.createElement("hr");
  hr2.classList.add("my-4", "w-2/3", "mx-auto");
  sidePanel.appendChild(hr2);

  const portsTable = document.createElement("table");
  portsTable.classList.add("mx-auto", "w-2/3", "table");
  sidePanel.appendChild(portsTable);

  const portsTableBody = document.createElement("tbody");
  portsTable.appendChild(portsTableBody);

  const portsTableHeading = document.createElement("tr");
  portsTableBody.appendChild(portsTableHeading);

  const portsTableHeading1 = document.createElement("th");
  portsTableHeading1.textContent = "Port";
  portsTableHeading1.classList.add("font-semibold");
  portsTableHeading.appendChild(portsTableHeading1);

  const portsTableHeading2 = document.createElement("th");
  portsTableHeading2.textContent = "Service";
  portsTableHeading2.classList.add("font-semibold");
  portsTableHeading.appendChild(portsTableHeading2);

  const portsTableHeading3 = document.createElement("th");
  portsTableHeading3.textContent = "Protocol";
  portsTableHeading3.classList.add("font-semibold");
  portsTableHeading.appendChild(portsTableHeading3);

  device.ports?.forEach((port) => {
    const portRow = document.createElement("tr");
    portsTableBody.appendChild(portRow);

    const portRow1 = document.createElement("td");

    portRow1.textContent = port.id.toString();

    // if (
    //   port.protocol === "http" || port.protocol === "https" ||
    //   port.service === "http" || port.service === "https"
    //   ) {
    //   const portLink = document.createElement("a");
    //   portLink.href = `http://${device.ip}:${port.port}`;
    //   portLink.classList.add("text-accent", "underline");
    //   portLink.target = "_blank";
    //   portLink.textContent = port.port.toString();
    //   portRow1.innerHTML = "";
    //   portRow1.appendChild(portLink);
    // }

    portRow.appendChild(portRow1);

    const portRow2 = document.createElement("td");
    portRow2.textContent = port.service.name;
    portRow.appendChild(portRow2);

    const portRow3 = document.createElement("td");
    portRow3.textContent = port.protocol;
    portRow.appendChild(portRow3);
  });

  const helpText = document.createElement("p");
  helpText.textContent = "Click to copy json data to clipboard";
  helpText.classList.add("text-center", "text-sm", "text-accent", "cursor-pointer");
  sidePanel.appendChild(helpText);
  helpText.addEventListener("click", () => {
    navigator.clipboard.writeText(JSON.stringify(device));
  });
}

export function resetPositionAndScale() {
  startX = 0;
  startY = 0;
  translateX = 0;
  translateY = 0;
  scale = 1;
  updateSvg();
}