<template>
  <div v-if="!isLoading" class="minvh">
    ok
  </div>
</template>
<script>
import { mapGetters } from 'vuex'

export default {
  components: {
  },
  data() {
    return {
      store: this.$store
    }
  },
  computed: {
    ...mapGetters([
      'projectid',
      'project',
      'isLoading'
    ])
  },
  mounted() {
  },
  methods: {
  }
}
</script>

<style lang="scss" module>
.todos {
  float: right;
}
</style>
