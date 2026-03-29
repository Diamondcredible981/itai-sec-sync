<template>
  <div class="app">
    <nav class="sidebar">
      <div class="logo">
        <svg viewBox="0 0 32 32" class="logo-icon">
          <defs>
            <linearGradient id="logoGrad" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stop-color="#00d4aa"/>
              <stop offset="100%" stop-color="#00a8cc"/>
            </linearGradient>
          </defs>
          <circle cx="16" cy="16" r="8" fill="none" stroke="url(#logoGrad)" stroke-width="2"/>
          <circle cx="16" cy="8" r="2.5" fill="url(#logoGrad)"/>
          <circle cx="22" cy="20" r="2.5" fill="url(#logoGrad)"/>
          <circle cx="10" cy="20" r="2.5" fill="url(#logoGrad)"/>
        </svg>
        <span class="logo-text">安全分析</span>
      </div>
      <div class="nav-items">
        <router-link to="/" class="nav-item" :class="{ active: isActive('/') }">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <path d="M12 6v6l4 2"/>
          </svg>
          <span>概览</span>
        </router-link>
        <router-link to="/topology" class="nav-item" :class="{ active: isActive('/topology') && !isActive('/topology-manage') }">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="6" cy="6" r="3"/>
            <circle cx="18" cy="6" r="3"/>
            <circle cx="6" cy="18" r="3"/>
            <circle cx="18" cy="18" r="3"/>
            <path d="M6 9v6M18 9v6M9 6h6M9 18h6"/>
          </svg>
          <span>网络拓扑</span>
        </router-link>
        <router-link to="/topology-manage" class="nav-item" :class="{ active: isActive('/topology-manage') }">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5z"/>
            <path d="M2 17l10 5 10-5"/>
            <path d="M2 12l10 5 10-5"/>
          </svg>
          <span>拓扑管理</span>
        </router-link>
        <router-link to="/analysis" class="nav-item" :class="{ active: isActive('/analysis') }">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 3v18h18"/>
            <path d="M7 16l4-4 4 4 6-6"/>
          </svg>
          <span>能力分析</span>
        </router-link>
        <router-link to="/suggest" class="nav-item" :class="{ active: isActive('/suggest') }">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5z"/>
            <path d="M2 17l10 5 10-5"/>
            <path d="M2 12l10 5 10-5"/>
          </svg>
          <span>优化建议</span>
        </router-link>
        <router-link to="/manage" class="nav-item" :class="{ active: isActive('/manage') }">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="18" height="18" rx="2"/>
            <path d="M9 9h6M9 15h6M9 12h6"/>
          </svg>
          <span>数据管理</span>
        </router-link>
      </div>
      <div class="nav-footer">
        <label class="theme-switch">
          <input type="checkbox" :checked="isDark" @change="toggleTheme">
          <span class="switch-track">
            <span class="track-bg">
              <span class="track-icon sun-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="5"/>
                  <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
                </svg>
              </span>
              <span class="track-text light-label">浅色</span>
              <span class="track-icon moon-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
                </svg>
              </span>
              <span class="track-text dark-label">深色</span>
            </span>
            <span class="switch-knob">
              <svg v-if="isDark" class="knob-icon moon-knob" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
              </svg>
              <svg v-else class="knob-icon sun-knob" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="5"/>
                <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
              </svg>
              <span class="knob-text">{{ isDark ? '深色' : '浅色' }}</span>
            </span>
          </span>
        </label>
        <div class="version">v1.0.0</div>
      </div>
    </nav>
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { useTheme } from './composables/useTheme'

const route = useRoute()
const { isDark, toggleTheme } = useTheme()

function isActive(path) {
  if (path === '/') {
    return route.path === '/'
  }
  return route.path.startsWith(path)
}
</script>

<style scoped>
.app {
  display: flex;
  min-height: 100vh;
  background: var(--bg-primary);
}

.sidebar {
  width: 220px;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 20px;
  border-bottom: 1px solid var(--border-color);
}

.logo-icon {
  width: 32px;
  height: 32px;
}

.logo-text {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: 0.5px;
}

.nav-items {
  flex: 1;
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 8px;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.nav-item:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.nav-item.active {
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-secondary));
  color: white;
}

.nav-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.nav-footer {
  padding: 16px 20px;
  border-top: 1px solid var(--border-color);
}

.theme-switch {
  display: block;
  cursor: pointer;
}

.theme-switch input {
  display: none;
}

.switch-track {
  position: relative;
  display: flex;
  align-items: center;
  width: 140px;
  height: 40px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 22px;
  overflow: hidden;
}

.track-bg {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0 16px;
  z-index: 1;
}

.track-icon {
  width: 18px;
  height: 18px;
  opacity: 0.35;
  transition: opacity 0.3s ease;
}

.track-icon svg {
  width: 100%;
  height: 100%;
}

.track-text {
  font-size: 12px;
  font-weight: 500;
  opacity: 0.35;
  transition: opacity 0.3s ease;
  color: var(--text-secondary);
}

/* Switch knob */
.switch-knob {
  position: absolute;
  top: 4px;
  left: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
  width: 64px;
  height: 30px;
  background: var(--accent-primary);
  border-radius: 16px;
  z-index: 2;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.knob-icon {
  width: 16px;
  height: 16px;
  margin-left: 8px;
  color: white;
}

.knob-icon svg {
  width: 100%;
  height: 100%;
}

.knob-text {
  font-size: 11px;
  font-weight: 600;
  color: white;
  margin-right: 8px;
}

/* Light mode (unchecked) - knob on right */
.theme-switch input:not(:checked) + .switch-track {
  background: var(--bg-secondary);
  border-color: #e0e0e0;
}

.theme-switch input:not(:checked) + .switch-track .track-bg {
  color: var(--text-tertiary);
}

.theme-switch input:not(:checked) + .switch-track .light-label {
  opacity: 0.6;
}

.theme-switch input:not(:checked) + .switch-track .moon-icon {
  opacity: 0.5;
}

.theme-switch input:not(:checked) + .switch-track .switch-knob {
  left: calc(100% - 68px);
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.theme-switch input:not(:checked) + .switch-track .knob-icon {
  color: #333;
}

.theme-switch input:not(:checked) + .switch-track .knob-text {
  color: #333;
}

/* Dark mode (checked) - knob on left */
.theme-switch input:checked + .switch-track {
  background: var(--bg-secondary);
  border-color: var(--accent-primary);
}

.theme-switch input:checked + .switch-track .dark-label {
  opacity: 0.6;
}

.theme-switch input:checked + .switch-track .sun-icon {
  opacity: 0.5;
}

.theme-switch input:checked + .switch-track .switch-knob {
  left: 4px;
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-secondary));
}

.theme-switch input:checked + .switch-track .knob-icon,
.theme-switch input:checked + .switch-track .knob-text {
  color: white;
}

.version {
  font-size: 12px;
  color: var(--text-tertiary);
  font-family: 'JetBrains Mono', monospace;
  margin-top: 12px;
}

.main-content {
  flex: 1;
  margin-left: 220px;
  min-height: 100vh;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
