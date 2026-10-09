<script setup>
import { ref, onMounted } from "vue";
import RFB from "@novnc/novnc";

const status = ref(null);
const image = ref([]);
const containers = ref([]);

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
    </div>
  </div>


  <div v-if="status && status != 'Connected'" class="top_bar">
    <div class="status">{{ status }}</div>
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
</style>
