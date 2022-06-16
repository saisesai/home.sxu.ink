<script setup lang="ts">
import {onMounted, ref} from "vue";
import axios, {AxiosError, AxiosResponse} from "axios";
import {ApiDeviceInfoRelayResp} from "../defs";
import {HandleServerInternalError} from "../utils";
import {Switch} from "ant-design-vue";
import "ant-design-vue/es/switch/style/css";

const props = defineProps<{
  client_id: string
}>();

const softwareVersion = ref<string>("0.0");
const hardwareVersion = ref<string>("0.0");
const online = ref<boolean>(false);
const on = ref<boolean>(false);

const swClkCb = () => {
  axios.post("/api/device/cmd",
      {"type": "relay", "client_id": props.client_id, "cmd": on.value ? "ON" : "OFF"})
      .catch(() => {
        HandleServerInternalError();
      })
}

onMounted(() => {
  axios.post("/api/device/info", {"type": "relay", "client_id": props.client_id})
      .then((res: AxiosResponse<ApiDeviceInfoRelayResp>) => {
        softwareVersion.value = res.data.sv;
        hardwareVersion.value = res.data.hv;
        online.value = res.data.online;
        on.value = res.data.on;
      })
      .catch((err: AxiosError<ApiDeviceInfoRelayResp>) => {
        HandleServerInternalError();
      })
})
</script>

<template>
  <div class="relative text-left mt-2">
    <div class="relative inline-block w-1/3 text-center">
      <div class="relative">软件版本：{{ softwareVersion }}</div>
      <div class="relative">硬件版本：{{ hardwareVersion }}</div>
    </div>
    <div class="relative inline-block w-1/3 text-center">
      <div class="relative inline-block w-3 h-3 rounded-full bg-gray-500"
           :class="{'relay-online': online}"></div>
      <div class="relative">{{ online ? "在线" : "离线" }}</div>
    </div>
    <div class="relative inline-block h-full w-1/3 text-center bottom-3">
      <Switch
          v-model:checked="on"
          :disabled="!online"
          checked-children="开"
          un-checked-children="关"
          @click="swClkCb"
      />
    </div>
  </div>
</template>

<style scoped lang="postcss">
.relay-online {
  @apply bg-green-500
}
</style>