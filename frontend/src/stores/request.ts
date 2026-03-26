import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useRequestStore = defineStore('request', () => {
  const pendingRequests = ref(0)
  const error = ref<string | null>(null)

  const isLoading = computed(() => pendingRequests.value > 0)

  const startRequest = () => {
    pendingRequests.value++
    error.value = null
  }

  const endRequest = () => {
    pendingRequests.value = Math.max(0, pendingRequests.value - 1)
  }

  const setError = (message: string) => {
    error.value = message
  }

  const clearError = () => {
    error.value = null
  }

  return { isLoading, error, startRequest, endRequest, setError, clearError }
})