import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useAuth = defineStore('auth', () => {
  const isAuth = ref(false)
  function setAuth(auth: boolean) {
    isAuth.value = auth
  }

  return { isAuth, setAuth }
})
