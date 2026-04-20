import { createRouter, createWebHistory } from 'vue-router'

import AppLayout from '@/layouts/AppLayout.vue'

import Dashboard from '@/views/Dashboard.vue'
import Containers from '@/views/Containers.vue'
import ContainerDetail from '@/views/ContainerDetail.vue'
import ContainerCreate from '@/views/ContainerCreate.vue'
import Images from '@/views/Images.vue'
import Volumes from '@/views/Volumes.vue'
import Networks from '@/views/Networks.vue'
import Compose from '@/views/Compose.vue'
import Login from '@/views/Login.vue'

const routes = [
    {
        path: '/',
        component: Login
    },
    {
        path: '/dashboard',
        component: AppLayout,
        children: [
            {
                path: '',
                component: Dashboard
            },
            {
                path: 'containers',
                component: Containers
            },
            {
                path: 'containers/create',
                component: ContainerCreate
            },
            {
                path: 'containers/:id',
                component: ContainerDetail
            },
            {
                path: 'images',
                component: Images
            },
            {
                path: 'volumes',
                component: Volumes
            },
            {
                path: 'networks',
                component: Networks
            },
            {
                path: 'compose',
                component: Compose
            }
        ]
    }
]

export default createRouter({
    history: createWebHistory(),
    routes
})
