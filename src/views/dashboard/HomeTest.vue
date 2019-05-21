<template>
  <div v-if="!isLoading" class="minvh">
    <!--
    <div class="full-width center-content">
    </div>
    -->
    <list-tasks v-if="projectid" :project="projectid"/>
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
// import ListTable from '@/views/admin/Projects/ListTable'
import ListTasks from '@/components/Tasks/ListTasks'

export default {
  components: {
    ListTasks
  },
  data() {
    return {
      store: this.$store
    }
  },
  computed: {
    ...mapGetters([
      'projectid',
      'isLoading'
    ])
  },
  /*
  apollo: {
    $client() {
      return this.activeClient // window.Apollo // this.activeClient
    },
    notes: {
      query: () => listTasks,
      update(data) {
        console.log('data: ', data)
        return data.listTasks
      }
    }
  },
  */
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
