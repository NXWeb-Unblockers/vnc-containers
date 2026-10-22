<script setup>
import { ref, onMounted } from "vue";
import RFB from "@novnc/novnc";

const status = ref(null);
const image = ref([]);
const containers = ref([]);
const clipboardText = ref("");
const clipboardOpen = ref(false);

async function connect(path) {
  const url =
    (location.protocol === "https:" ? "wss" : "ws") +
    "://" +
    location.host +
    `/websockify/${path}`;

  console.log(url);

  const rfb = new RFB(document.getElementById("screen"), url);

  rfb.scaleViewport = true;
  rfb.resizeSession = true;

  rfb.addEventListener("connect", () => {
    status.value = "Connected";
  });
  rfb.addEventListener("disconnect", () => {
    status.value = "Disconnected";
  });
  rfb.addEventListener("clipboard", (event) => {
    clipboardText.value = event.detail.text;
  });
}

// sends the textarea contents to the remote session whenever the user edits it
function sendClipboard() {
  rfb.clipboardPasteFrom(clipboardText.value);
}

function CreateContainer(name) {
  status.value = "Creating container";

  fetch("/api/createContainer", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ name })
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      }
      throw new Error("Network response was not ok");
    })
    .then((data) => {
      console.log(data);
      setTimeout(() => {
        connect(data.id);
      }, 5000);
    })
    .catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
    });
}

function StartContainer(name) {
  status.value = "Starting container";

  fetch("/api/startContainer", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ name })
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      }
      throw new Error("Network response was not ok");
    })
    .then((data) => {
      console.log(data);
      setTimeout(() => {
        connect(data.id);
      }, 5000);
    })
    .catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
    });
}

function DeleteContainer(name) {
  status.value = "Deleting container";

  fetch("/api/deleteContainer", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ name })
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      }
      throw new Error("Network response was not ok");
    })
    .then((data) => {
      console.log(data);
      getContainers();
      status.value = null;
    })
    .catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
    });
}

function getImages() {
  fetch("/api/getImages")
    .then((response) => {
      if (response.ok) {
        return response.json();
      }
      throw new Error("Network response was not ok");
    })
    .then((data) => {
      console.log(data);
      image.value = data;
    })
    .catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
    });
}
function getContainers() {
  fetch("/api/getContainers")
    .then((response) => {
      if (response.ok) {
        return response.json();
      }
      throw new Error("Network response was not ok");
    })
    .then((data) => {
      console.log(data);
      containers.value = data;
    })
    .catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
    });
}

onMounted(() => {
  getImages();
  getContainers();
});
</script>

<template>
  <div v-if="image.length > 0 && !status ">
    <h2>Select an image to create a new container</h2>
    <div v-for="img in image" :key="img">
      <button @click="CreateContainer(img)">{{ img }}</button>
    </div>
  </div>

  <div v-if="containers.length > 0 && !status ">
    <h2>Start a Container</h2>
    <div v-for="container in containers" :key="container">
      <button @click="StartContainer(container)">{{ container }}</button>
      <button @click="DeleteContainer(container)">Delete</button>
    </div>
  </div>


  <div v-if="status && status != 'Connected'" class="top_bar">
    <div class="status">{{ status }}</div>
  </div>

  <div v-if="status === 'Connected'" class="clipboard-drawer" :class="{ open: clipboardOpen }">
    <div class="clipboard-panel">
      <label for="clipboard">Clipboard (shared with VNC session)</label>
      <textarea
        id="clipboard"
        v-model="clipboardText"
        @input="sendClipboard"
        placeholder="Paste here to send to the VNC session, or copy text copied inside the session"
      ></textarea>
    </div>
    <button
      class="clipboard-tab"
      type="button"
      :aria-expanded="clipboardOpen"
      @click="clipboardOpen = !clipboardOpen"
    >
      Clipboard
    </button>
  </div>

  <div id="screen"></div>
</template>

<style scoped>
.status {
  text-align: center;
}

#screen {
  flex: 1;
  overflow: hidden;
  height: 100%;
  width: 100%;
}

.top_bar {
  background-color: #6e84a3;
  color: white;
  font: bold 12px Helvetica;
  padding: 6px 5px 4px 5px;
  border-bottom: 1px outset;
}

.clipboard-drawer {
  --clipboard-tab-width: 22px;
  position: fixed;
  top: 50%;
  right: 0;
  transform: translate(calc(100% - var(--clipboard-tab-width)), -50%);
  transition: transform 0.25s ease;
  display: flex;
  flex-direction: row-reverse;
  align-items: stretch;
  z-index: 1000;
}

.clipboard-drawer.open {
  transform: translateY(-50%);
}

.clipboard-tab {
  box-sizing: border-box;
  width: var(--clipboard-tab-width);
  writing-mode: vertical-rl;
  transform: rotate(180deg);
  align-self: center;
  background-color: #6e84a3;
  color: white;
  border: none;
  border-radius: 0 6px 6px 0;
  padding: 10px 6px;
  cursor: pointer;
  font: bold 12px Helvetica;
}

.clipboard-panel {
  display: flex;
  flex-direction: column;
  width: 260px;
  padding: 10px;
  background: white;
  border: 1px solid #6e84a3;
  border-right: none;
  box-shadow: -2px 0 6px rgba(0, 0, 0, 0.2);
}

.clipboard-panel label {
  font: bold 12px Helvetica;
  margin-bottom: 6px;
}

.clipboard-panel textarea {
  width: 100%;
  height: 120px;
  resize: vertical;
}
</style>
