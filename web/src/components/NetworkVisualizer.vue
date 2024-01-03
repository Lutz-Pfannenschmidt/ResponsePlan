<script setup lang="ts">
import { start, resetPositionAndScale } from './NetworkVisualizer';
import { onMounted } from 'vue';
import CrosshairIcon from '@/components/icons/CrosshairIcon.vue'
import CloseIcon from '@/components/icons/CloseIcon.vue'
import NetworkIcon from '@/components/icons/NetworkIcon.vue'
import { dev } from '@/socket';

async function requestScan() {
  let resp: Response;

  if (dev) {
    resp = await fetch("http://127.0.0.1:1337/api/scan")
  } else {
    resp = await fetch("/api/scan")
  }

  let id = await resp.text()
  alert(id)
}

onMounted(() => {
  start()
})

</script>

<template>
  <div class="canvas-wrapper z-[-1]"></div>

  <div class="tool-wrapper">
    
    <button class="device-info-close absolute right-[-10px] btn btn-square rounded-sm btn-accent btn-outline z-10 bg-base-100">
      <CloseIcon/>
    </button>

    <div class="side-panel device-info w-1/5
    z-1 border-y-4 border-l-4 overflow-scroll p-2 mb-2 bg-base-300 shadow-lg
      rounded-l border-base-content opacity-90 absolute right-[-3rem] whitespace-break-spaces"
    >
  </div>
    <div class="toolbar z-1 bottom-0 left-0 absolute m-4">
      <div class="tooltip tooltip-right tooltip-info" data-tip="reset canvas position">
        <button class="btn btn-square btn-accent btn-outline mr-2" @click="resetPositionAndScale()">
          <CrosshairIcon/>
        </button>
      </div>
      <div class="tooltip tooltip-right tooltip-info" data-tip="new network scan">
        <button class="btn btn-square btn-accent btn-outline mr-2" @click="requestScan()">
          <NetworkIcon/>
        </button>
      </div>
    </div>
  </div>
</template>