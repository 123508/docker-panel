import { createRouter, createWebHistory } from 'vue-router'
import { token } from '@/store/auth'

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
        path: '/login',
        component: Login
    },
    {
        path: '/dashboard',
        component: AppLayout,
        meta: { requiresAuth: true },
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

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, from, next) => {
    if (to.matched.some(record => record.meta.requiresAuth)) {
        if (!token.value) {
            next({ path: '/login' })
        } else {
            next()
        }
    } else {
        if (token.value && (to.path === '/' || to.path === '/login')) {
            next({ path: '/dashboard' })
        } else {
            next()
        }
    }
})

export default router
