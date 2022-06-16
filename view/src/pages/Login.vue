<script setup lang="ts">
import {Button, Input, message} from "ant-design-vue";
import "ant-design-vue/es/input/style/css";
import "ant-design-vue/es/button/style/css";
import "ant-design-vue/es/message/style/css";
import {LockOutlined, UserOutlined} from "@ant-design/icons-vue";
import {ref} from "vue";
import {useRouter} from "vue-router";
import axios, {AxiosError, AxiosResponse} from "axios";
import {ApiUserLoginResp} from "../defs"

const router = useRouter()

let Username = ref<string>("");
let Password = ref<string>("");
let ConfirmBtnLoading = ref<boolean>(false);

const ConfirmBtnCb = () => {
  ConfirmBtnLoading.value = true;
  // build form
  let formData = new FormData();
  formData.append("username", Username.value);
  formData.append("password", Password.value);
  // do request
  axios.post("/api/user/login", formData).then((res: AxiosResponse<ApiUserLoginResp>) => {
    localStorage.setItem("token", res.data.token);
    axios.defaults.headers.common["Authorization"] = "Bearer " + res.data.token;
    router.push("/");
    message.success("登录成功！", 1);
    ConfirmBtnLoading.value = false;
  }).catch((err: AxiosError<ApiUserLoginResp>) => {
    if (err.response?.data.msg == "user not found or wrong password") {
      message.error("用户名不存在或密码错误！", 1);
    } else {
      message.error("服务端出错！请等待修复或联系开发人员！");
    }
    ConfirmBtnLoading.value = false;
  })
}
</script>

<template>
  <div class="absolute w-full min-h-screen px-10">
    <div class="relative space-y-6 top-36">
      <div class="relative text-3xl">登录</div>
      <Input v-model:value="Username">
        <template #prefix>
          <UserOutlined/>
        </template>
      </Input>
      <Input type="password" v-model:value="Password">
        <template #prefix>
          <LockOutlined></LockOutlined>
        </template>
      </Input>
      <Button block type="primary" :loading="ConfirmBtnLoading" @click="ConfirmBtnCb">确认</Button>
    </div>
  </div>
</template>