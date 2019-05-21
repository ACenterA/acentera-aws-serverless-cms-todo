<template>
  <section class="serverless-cms-container">
    <transition name="fade-transform" mode="out-in">
      <keep-alive :include="cachedViews">
        <router-view :key="key"/>
      </keep-alive>
    </transition>
  </section>
</template>

<script>
export default {
  name: 'AppMain',
  computed: {
    cachedViews() {
      if (this.$store && this.$store.state && this.$store.state.tagsView) {
        console.error('got tag views..')
        return this.$store.state.tagsView.cachedViews
      } else {
        return {}
      }
    },
    key() {
      return this.$route.fullPath
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss">
.serverless-cms-container {
  background-color: rgb(240, 242, 245);
}
</style>
