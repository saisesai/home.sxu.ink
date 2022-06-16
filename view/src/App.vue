<script setup lang="ts">
import {onMounted} from "vue";
import {useRouter} from "vue-router";
import axios from "axios";

const router = useRouter()
let token = localStorage.getItem("token")

onMounted(() => {
  if (token == null) {
    localStorage.clear();
    router.push("/login/");
    return;
  }
  axios.defaults.headers.common["Authorization"] = "Bearer " + token;
});
</script>

<template>
  <div class="absolute container left-0 right-0 mx-auto bg-white max-w-md min-h-screen text-center">
    <router-view></router-view>
  </div>
</template>