import { ref, watch, onMounted } from 'vue'

const STORAGE_KEY = 'app-theme'
const DARK = 'dark'
const LIGHT = 'light'

const isDark = ref(true)

function initTheme() {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved) {
    isDark.value = saved === DARK
  } else {
    isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  applyTheme()
}

function applyTheme() {
  const html = document.documentElement
  if (isDark.value) {
    html.classList.remove('light')
  } else {
    html.classList.add('light')
  }
}

function toggleTheme() {
  isDark.value = !isDark.value
  localStorage.setItem(STORAGE_KEY, isDark.value ? DARK : LIGHT)
  applyTheme()
}

function setTheme(theme) {
  isDark.value = theme === DARK
  localStorage.setItem(STORAGE_KEY, theme)
  applyTheme()
}

export function useTheme() {
  onMounted(() => {
    initTheme()
  })

  watch(isDark, () => {
    applyTheme()
  })

  return {
    isDark,
    toggleTheme,
    setTheme,
    DARK,
    LIGHT
  }
}
