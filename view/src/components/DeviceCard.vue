<script setup lang="ts">
import {ApiUserDeviceDelResp, Device} from "../defs";
import RelayContent from "./DeviceCardRelay.vue";
import {Button, Input, message, Modal, Popconfirm} from "ant-design-vue";
import "ant-design-vue/es/button/style/css";
import "ant-design-vue/es/modal/style/css";
import "ant-design-vue/es/input/style/css";
import "ant-design-vue/es/message/style/css";
import "ant-design-vue/es/popconfirm/style/css";
import {DeleteOutlined, EditOutlined} from "@ant-design/icons-vue";
import axios, {AxiosError} from "axios";
import {ref} from "vue";
import {HandleServerInternalError} from "../utils";

const props = defineProps<{
  device: Device | any
}>();

let content: any | null = null
if (props.device.type == "relay") {
  content = RelayContent;
}

const modVis = ref<boolean>(false);

const remarkNew = ref<string>("");

const modOkCb = () => {
  if (remarkNew.value == "") {
    message.error("输入不能为空！", 1);
    return
  }
  axios.post("/api/user/device/mod",
      {"client_id": props.device.client_id, "remark": remarkNew.value})
      .then(() => {
        message.success("更新设备名称成功");
        location.reload();
      })
      .catch(() => {
        HandleServerInternalError();
      })
  modVis.value = false
}

const deleteConfirmCb = () => {
  axios.post("/api/user/device/del",
      {"type": props.device.type, "client_id": props.device.client_id})
      .then(() => {
        message.success("删除设备成功！")
        location.reload();
      })
      .catch((err: AxiosError<ApiUserDeviceDelResp>) => {
        HandleServerInternalError();
      })
}
</script>

<template>
  <div class="relative p-2 w-full bg-white border border-gray-100 drop-shadow-md rounded mb-3">
    <div class="relative px-1 text-xl text-left  border-b-2 border-black">
      {{ props.device.remark }}
      <Button shape="circle" size="small" class="bottom-0.5 ml-1" @click="modVis=true">
        <template #icon>
          <EditOutlined/>
        </template>
      </Button>
      <Popconfirm title="确认删除设备？" ok-text="确认" cancel-text="取消" @confirm="deleteConfirmCb">
        <Button danger shape="circle" size="small" class="bottom-0.5 mx-3">
          <template #icon>
            <DeleteOutlined/>
          </template>
        </Button>
      </Popconfirm>
      <Modal v-model:visible="modVis" title="编辑设备名称" cancelText="取消" okText="确认" @ok="modOkCb">
        <div class="space-y-5">
          <div><span class="font-bold">设备ID: </span>{{ props.device.client_id }}</div>
          <div class="space-y-2">
            <div class="font-bold">新的设备名称：</div>
            <Input v-model:value="remarkNew"></Input>
          </div>
        </div>
      </Modal>
      <span class="relative text-sm float-right bottom-0">{{ props.device.client_id }}</span>
    </div>
    <component :is="content" :client_id="props.device.client_id"></component>
  </div>
</template>