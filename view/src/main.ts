import {createApp} from 'vue';
import App from './App.vue';
import "./style.pcss";
import {createRouter, createWebHistory} from "vue-router";
import NProgress from "nprogress";
import "nprogress/nprogress.css";
import axios from "axios";

// axios请求拦截器
axios.interceptors.request.use(
    config => {
        NProgress.start()
        return config
    },
    error => {
        NProgress.done()
        return Promise.reject(error)
    }
)
// axios响应拦截器
axios.interceptors.response.use(
    function (response) {
        NProgress.done()
        return response
    },
    function (error) {
        NProgress.done()
        return Promise.reject(error)
    }
)

const app = createApp(App);

const routes = [
    {
        path: "/",
        name: "Home",
        component: () => import("./pages/Home.vue"),
    },
    {
        path: "/login",
        name: "Login",
        component: () => import("./pages/Login.vue"),
    },
    {
        path: "/:pathMatch(.*)*",
        name: "NotFound",
        component: () => import("./pages/NotFound.vue"),
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeEach((to, from, next) => {
    NProgress.start()
    next()
})

router.afterEach(() => {
    NProgress.done()
})

app.use(router)

app.mount('#app')
