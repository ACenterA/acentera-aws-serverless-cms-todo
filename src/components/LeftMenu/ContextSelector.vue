<template>
  <li role="menuitem" tabindex="-1" class="el-menu-item submenu-title-noDropdown" style="color: rgb(191, 203, 217); background-color: rgb(48, 65, 86);">
    <span>
      <el-select
        v-model="value9"
        reserve-keyword
        @change="onChange"
      >
        <el-option
          v-for="item in list"
          :key="item.id"
          :label="item.title"
          :value="item.id"
        />
      </el-select>
    </span>
  </li>
</template>

<script>
import { mapGetters } from 'vuex'
import { ALL_PROJECTS } from '@/gql/queries/projects.gql'

export default {
  data() {
    return {
      options4: [],
      value9: null,
      list: [],
      loading: false,
      states: []
    }
  },
  computed: {
    ...mapGetters([
      'projectid',
      'project',
      'isLoading'
    ]),
    activeClient() { return this.$store.getters.apollo }
  },
  apollo: {
    $client() {
      return this.activeClient
    },
    projects: {
      query: () => ALL_PROJECTS,
      update(data) {
        this.list = data.listAllProjects.items
        this.total = data.listAllProjects.items.length
        this.listLoading = false
        this.$store.dispatch('setProjects', this.list)
        return data.listAllProjects
      },
      /*
      // example if we want to have live updates
      subscribeToMore: {
        document: onCreateProject,
        // Variables passed to the subscription. Since we're using a function,
        // they are reactive
        //variables () {
        //  return {
        //    param: this.param,
        //  }
        // },
        // Mutate the previous result
        updateQuery: (previousResult, { subscriptionData }) => {
          console.error('recive de vent ehre')
          const previousItems = previousResult.listAllProjects.items

          // Here, return the new result from the previous with the new data
          var newData = subscriptionData.data.onCreateProject
          previousItems.push(newData)
        },
      },
      */
      fetchPolicy: 'cache-and-network' // 'network-only', // skip the cache
    }
  },
  watch: {
    projectid() {
      console.log('projectid Changed!')
      if (!this.value9) {
        this.value9 = this.projectid
      }
    }
  },
  mounted() {
    this.list = this.states.map(item => {
      return { value: item, label: item }
    })
    if (!this.value9) {
      this.value9 = this.projectid
    }
  },
  methods: {
    onChange: function(v) {
      console.error('on select of ')
      console.error(v)
      this.$store.dispatch('setProject', v).then(res => {
        console.error('received of')
        console.error(v)
      })
    }
  }
}
</script>
