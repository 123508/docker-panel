import { createRouter, createWebHistory } from 'vue-router'
import { token } from '@/store/auth'

import AppLayout from '@/layouts/AppLayout.vue'

import Dashboard from '@/views/Dashboard.vue'
import Containers from '@/views/Containers.vue'
import ContainerDetail from '@/views/ContainerDetail.vue'
import ContainerCreate from '@/views/ContainerCreate.vue'
import Images from '@/views/Images.vue'
import Volumes from '@/views/Volumes.vue'
import VolumeCreate from '@/views/VolumeCreate.vue'
import VolumeDetail from '@/views/VolumeDetail.vue'
import Networks from '@/views/Networks.vue'
import NetworkCreate from '@/views/NetworkCreate.vue'
import NetworkDetail from '@/views/NetworkDetail.vue'
import ImageDetail from '@/views/ImageDetail.vue'
import ImagePull from '@/views/ImagePull.vue'
import Compose from '@/views/Compose.vue'
import ComposeCreate from '@/views/ComposeCreate.vue'
import ComposeContent from '@/views/ComposeContent.vue'
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
                path: 'images/pull',
                component: ImagePull
            },
            {
                path: 'volumes',
                component: Volumes
            },
            {
                path: 'volumes/create',
                component: VolumeCreate
            },
            {
                path: 'volumes/:name',
                component: VolumeDetail
            },
            {
                path: 'networks',
                component: Networks
            },
            {
                path: 'networks/create',
                component: NetworkCreate
            },
            {
                path: 'networks/:id',
                component: NetworkDetail
            },
            {
                path: 'images/:id',
                component: ImageDetail
            },
            {
                path: 'compose/create',
                component: ComposeCreate
            },
            {
                path: 'compose/:name',
                component: ComposeContent
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
