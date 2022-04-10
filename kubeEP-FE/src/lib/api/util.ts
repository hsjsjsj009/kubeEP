import axios, {AxiosPromise, AxiosRequestHeaders, Method} from "axios";
import { envVariables } from "$lib/env-variables";

export interface Response<Data> {
    status: string,
    data: Data,
}

const axiosAPI = axios.create({
    baseURL: envVariables.apiBasePath as string
});

// implement a method to execute all the request from here.
const apiRequest = (method: Method, url: string, request?: unknown, headers?: AxiosRequestHeaders) => {
    //using the axios instance to perform the request that received from each http method
    return axiosAPI({
        method,
        url,
        data: request,
        headers,
    })
};

// function to execute the http get request
const get = <R>(url: string, headers?:AxiosRequestHeaders): AxiosPromise<Response<R>> => apiRequest("get", url, headers)

// function to execute the http delete request
const deleteRequest = <Request, R>(url: string, request: Request, headers?: AxiosRequestHeaders): AxiosPromise<Response<R>> =>  apiRequest("delete", url, request, headers);

// function to execute the http post request
const post = <Request, R>(url: string, request: Request, headers?: AxiosRequestHeaders): AxiosPromise<Response<R>> => apiRequest("post", url, request, headers);

// function to execute the http put request
const put = <Request, R>(url: string, request: Request, headers?: AxiosRequestHeaders): AxiosPromise<Response<R>> => apiRequest("put", url, request, headers);

// function to execute the http path request
const patch = <Request, R>(url: string, request: Request, headers?: AxiosRequestHeaders): AxiosPromise<Response<R>> =>  apiRequest("patch", url, request, headers);

// expose your method to other services or actions
const API ={
    get,
    delete: deleteRequest,
    post,
    put,
    patch
};


export default API;
