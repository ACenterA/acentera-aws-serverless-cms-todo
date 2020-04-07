<template>
  <div v-if="!isLoading" class="minvh">
    <div v-if="!hasProjects">
      <div v-if="checkPermission(['admin'])" label="Admin">
        You must create projects from the administrator panels.
      </div>
      <div v-if="!checkPermission(['admin'])" label="NonAdmin">
        No projects. Please contact admins.
      </div>
    </div>
    <div v-if="hasProjects">
      <div v-if="!projectid">
        Select a Project.
        <list-projects />
      </div>
      <div v-else>
        Current Project: {{ project.title }}
        <list-tasks v-if="projectid" :project="projectid" />
      </div>
    </div>
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
// import ListTable from '@/views/admin/Projects/ListTable'
import checkPermission from '@/utils/permission'
import ListTasks from '@/components/Tasks/ListTasks'
import ListProjects from '@/components/List/Projects'

export default {
  components: {
    ListTasks,
    ListProjects
  },
  data() {
    return {
      store: this.$store,
      sum: 0
    }
  },
  computed: {
    ...mapGetters([
      'hasProjects',
      'projectid',
      'project',
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
    checkPermission
  }
}
</script>

<style lang="scss" module>
.todos {
  float: right;
}
</style>
