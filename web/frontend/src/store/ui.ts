import { ref, watch } from 'vue'

type ThemePreference = 'system' | 'light' | 'dark'

const saved = localStorage.getItem('theme') as ThemePreference | null
export const themePreference = ref<ThemePreference>(saved ?? 'dark')

const mql = window.matchMedia('(prefers-color-scheme: dark)')

function applyTheme(pref: ThemePreference) {
  const html = document.documentElement
  let isDark: boolean

  if (pref === 'system') {
    isDark = mql.matches
  } else {
    isDark = pref === 'dark'
  }

  if (isDark) {
    html.classList.remove('light')
  } else {
    html.classList.add('light')
  }
}

mql.addEventListener('change', () => {
  if (themePreference.value === 'system') {
    applyTheme('system')
  }
})

watch(
  themePreference,
  (v) => {
    applyTheme(v)
    localStorage.setItem('theme', v)
  },
  { immediate: true }
)

export function setThemePreference(pref: ThemePreference) {
  themePreference.value = pref
}
