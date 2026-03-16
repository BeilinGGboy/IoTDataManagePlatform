<template>
  <router-view />
</template>

<script setup>
// 根组件：仅作为路由出口，所有页面通过 router-view 渲染
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
html, body, #app {
  width: 100%;
  height: 100%;
}
</style>
