import axios from 'axios'

export const http = axios.create({
    baseURL: '',
    headers: { 'Content-Type': 'application/json' }
})

http.interceptors.request.use(
    (config) => {

        return config
    },
    (error) => Promise.reject(error)
)

http.interceptors.response.use(
    (response) => {
        return response
    },
    (error) => {
        return Promise.reject(error)
    }
)


