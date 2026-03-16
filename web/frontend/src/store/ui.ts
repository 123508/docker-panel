import { ref, watch } from 'vue'

type Theme = 'light' | 'dark'

const saved = localStorage.getItem('theme') as Theme | null

export const uiStyle = ref<Theme>(saved ?? 'light')

watch(
    uiStyle,
    (v) => {

        const html = document.documentElement

        if (v === 'dark') {
            html.classList.add('dark')
        } else {
            html.classList.remove('dark')
        }

        localStorage.setItem('theme', v)

    },
    { immediate: true }
)