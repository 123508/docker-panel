import { ref } from 'vue'

const TOKEN_KEY = 'docker-panel-token'

export const token = ref(localStorage.getItem(TOKEN_KEY) || '')

export const setToken = (newToken: string) => {
  token.value = newToken
  if (newToken) {
    localStorage.setItem(TOKEN_KEY, newToken)
  } else {
    localStorage.removeItem(TOKEN_KEY)
  }
}

export const clearToken = () => {
  setToken('')
}
