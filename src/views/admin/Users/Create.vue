<template>
  <div v-if="!isLoading" class="minvh">
    <user-detail :is-edit="true"/>
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
import UserDetail from './components/UserDetail'

export default {
  name: 'CreateUser',
  components: {
    UserDetail
  },
  data() {
    return {
      store: this.$store
    }
  },
  computed: {
    ...mapGetters([
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
