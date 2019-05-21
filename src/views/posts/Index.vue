<template>
  <div v-if="!isLoading" class="minvh">
    <div v-if="!projectid">
      <list-projects/>
    </div>
    <div v-else>
      <list-posts v-if="projectid" :project="projectid"/>
    </div>
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
import ListPosts from '@/components/Posts/ListPosts'
import ListProjects from '@/components/List/Projects'

export default {
  components: {
    ListPosts,
    ListProjects
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
