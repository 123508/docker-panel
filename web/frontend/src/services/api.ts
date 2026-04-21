import axios from 'axios'
import { token, clearToken } from '@/store/auth'
import router from '@/router'

export const http = axios.create({
    baseURL: '',
    headers: { 'Content-Type': 'application/json' }
})

http.interceptors.request.use(
    (config) => {
        if (token.value) {
            config.headers.Authorization = `Bearer ${token.value}`
        }
        return config
    },
    (error) => Promise.reject(error)
)

http.interceptors.response.use(
    (response) => {
        return response
    },
    (error) => {
        if (error.response && error.response.status === 401) {
            clearToken()
            router.push('/login')
        }
        return Promise.reject(error)
    }
)


