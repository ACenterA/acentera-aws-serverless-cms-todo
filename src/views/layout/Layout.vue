<!-- TOOD: We could totally override the app if needed? -->
<template>
  <div :class="classObj" class="bbadwornapp-wrapper">
    OK
    <div v-if="device==='mobile'&&sidebar.opened" class="drawer-bg" @click="handleClickOutside"/>
    <sidebar class="sidebar-container"/>
    <div :class="mainContainerClassObj" class="main-container">
      <navbar/>
      <!--<tags-view/>-->
      <app-main/>
    </div>
  </div>
</template>

<script>
import { Navbar, Sidebar, AppMain, TagsView } from './components'
import ResizeMixin from './mixin/ResizeHandler'

export default {
  name: 'Layout',
  components: {
    Navbar,
    Sidebar,
    AppMain,
    TagsView
  },
  mixins: [ResizeMixin],
  computed: {
    sidebar() {
      try {
        return this.$store.state.testapp.sidebar
      } catch (e) {
        return true
      }
    },
    device() {
      try {
        return this.$store.state.testapp.device
      } catch (ex) {
        return 'desktop'
      }
    },
    classObj() {
      return {
        visible: true
      }
    },
    mainContainerClassObj() {
      return {
        hiddenSidebar: !this.sidebar.visible
      }
    }
  },
  methods: {
    handleClickOutside() {
      this.$store.dispatch('closeSideBar', { withoutAnimation: false })
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
  @import "src/styles/mixin.scss";
  .app-wrapper {
    @include clearfix;
    position: relative;
    height: 100%;
    width: 100%;
    &.mobile.openSidebar{
      position: fixed;
      top: 0;
    }
  }
  .drawer-bg {
    background: #000;
    opacity: 0.3;
    width: 100%;
    top: 0;
    height: 100%;
    position: absolute;
    z-index: 800;
  }
</style>
