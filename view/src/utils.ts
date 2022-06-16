import {message} from "ant-design-vue";
import "ant-design-vue/es/message/style/css";
import {Router} from "vue-router";

export function HandleInvalidToken(pRouter: Router) {
    localStorage.clear();
    pRouter.push("/login").then(r => {
    });
    message.error("登录失效，请重新登录", 1).then(r => {
    })
}

export function HandleServerInternalError() {
    message.error("服务端出错，请等待修复或联系开发人员！").then(r => {
    })
}