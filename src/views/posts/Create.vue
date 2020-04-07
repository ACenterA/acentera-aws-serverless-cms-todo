<template>
  <div v-if="!isLoading" class="minvh">
    <post-detail :is-edit="false" />
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
import PostDetail from './components/PostDetail'

export default {
  name: 'CreatePost',
  components: {
    PostDetail
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
