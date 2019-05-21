<template>
  <div class="full-width center-content">
    <img src="/static/img/logo.png" >
    <welcome-message name="World" />

    Welcome

    <button @click="testclick">CLick me</button>
  </div>
</template>

<script>

import WelcomeMessage from 'components/Home/WelcomeMessage'

export default {
  components: {
    WelcomeMessage
  },
  data() {
    return {
      store: this.$store
    }
  },
  mounted() {
    console.error('NICE IT WORKS')
    console.error(this)
    console.error(this.$store)
  },
  methods: {
    testclick() {
      console.error('NICECCECEC')
      console.error(this)
      console.error(this.$store)
      this.$store.dispatch('testSubmitFunction', 'aaa')
    }
  }
}
</script>
