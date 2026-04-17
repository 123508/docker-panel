import { createRouter, createWebHistory } from 'vue-router'

import AppLayout from '@/layouts/AppLayout.vue'

import Dashboard from '@/views/Dashboard.vue'
import Containers from '@/views/Containers.vue'
import Images from '@/views/Images.vue'
import Volumes from '@/views/Volumes.vue'
import Networks from '@/views/Networks.vue'
import Compose from '@/views/Compose.vue'

const routes = [
    {
        path: '/',
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
