<script setup lang="ts">
import {Button, Input, message, Modal, Popconfirm, Radio} from "ant-design-vue";
import "ant-design-vue/es/button/style/css";
import "ant-design-vue/es/popconfirm/style/css";
import "ant-design-vue/es/message/style/css";
import "ant-design-vue/es/modal/style/css";
import "ant-design-vue/es/input/style/css";
import "ant-design-vue/es/radio/style/css";
import {useRouter} from "vue-router";
import {onMounted, ref} from "vue";
import {ApiUserDeviceGetResp, Device, ApiUserDeviceAddResp} from "../defs";
import axios, {AxiosError, AxiosResponse} from "axios";
import {HandleInvalidToken, HandleServerInternalError} from "../utils";
import DeviceCard from "../components/DeviceCard.vue";

const router = useRouter();

const devices = ref<Device[]>();

const addModVis = ref<boolean>(false);

const newDevTypeRelay = ref<boolean>(true);
const newDevId = ref<string>("");
const newDevRem = ref<string>("");

const addModOkCb = () => {
  if (newDevId.value == "") {
    message.error("设备id不能为空！", 1);
    return
  }
  if (newDevRem.value == "") {
    message.error("设备名称不能为空！", 1);
    return;
  }
  let deviceType = "";
  if (newDevTypeRelay.value == true) deviceType = "relay";
  axios.post("/api/user/device/add",
      {"type": deviceType, "client_id": newDevId.value, "remark": newDevRem.value})
      .then(() => {
        message.success("添加设备成功！")
        location.reload();
      })
      .catch((err:AxiosError<ApiUserDeviceAddResp>) => {
        if(err.response?.data.msg == "user device already exist"){
          message.warning("设备已存在，请勿重复添加！");
          return
        }
        console.log(err.response)
        HandleServerInternalError();
      })
  newDevId.value = "";
  newDevRem.value = "";
  addModVis.value = false;
}

const ExitBtnCb = () => {
  localStorage.clear();
  router.push("/login");
}

onMounted(() => {
  axios.get("/api/user/device/get").then((res: AxiosResponse<ApiUserDeviceGetResp>) => {
    devices.value = res.data.devices;
  }).catch((err: AxiosError<ApiUserDeviceGetResp>) => {
    if (err.response?.data.msg == "invalid token") {
      HandleInvalidToken(router);
    } else {
      HandleServerInternalError();
    }
  })
})
</script>

<template>
  <div class="absolute w-full">
    <div class="relative text-left w-full h-12 p-2 bg-white">
      <span class="relative text-xl">设备列表</span>
      <span class="relative float-right space-x-2">
        <Button @click="addModVis=true">新增</Button>
        <Modal v-model:visible="addModVis" title="编辑设备名称" cancelText="取消" okText="确认" @ok="addModOkCb">
          <div class="space-y-5">
            <div class="space-y-2">
              <div class="font-bold">设备类型</div>
              <Radio v-model:checked="newDevTypeRelay">relay(家用继电器设备)</Radio>
            </div>
            <div class="space-y-2">
              <div class="font-bold">设备ID</div>
              <Input v-model:value="newDevId"></Input>
            </div>
            <div class="space-y-2">
              <div class="font-bold">设备名称</div>
              <Input v-model:value="newDevRem"></Input>
            </div>
          </div>
        </Modal>
        <Popconfirm title="确认要退出吗？" ok-text="确认" cancel-text="取消" @confirm="ExitBtnCb">
          <Button danger>退出</Button>
        </Popconfirm>
      </span>
    </div>
    <div class="relative p-3">
      <DeviceCard v-for="(item, index) in devices" :key="index" :device="item"></DeviceCard>
    </div>
  </div>
</template>