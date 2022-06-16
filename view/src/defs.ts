export type ApiUserLoginResp = {
    msg: string
    token: string
}

export type Device = {
    type: string
    client_id: string
    remark: string
}

export type ApiUserDeviceGetResp = {
    msg: string
    devices: Device[]
}

export type ApiUserDeviceAddResp = {
    msg: string
}

export type ApiUserDeviceDelResp = {
    msg: string
}

export type ApiDeviceInfoRelayResp = {
    msg: string
    sv: string
    hv: string
    online: boolean
    on: boolean
}

