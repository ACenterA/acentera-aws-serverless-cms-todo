<template>
  <div v-if="!isLoading" class="minvh">
    In Users ...

    <list-table/>

  </div>
</template>

<script>
import permission from '@/directive/permission/index.js'
import { mapGetters } from 'vuex'
import moment from 'moment'
import ListTable from './Users/ListTable'

export default{

  components: {
    ListTable
  },
  directives: { permission },
  filters: {
    dateToStr(str) {
      if (str) {
        var localTime = moment(str, 'YYYY-MM-DDTHH24:MI:ss.SSSZ').local()
        return localTime.format('YYYY/MM/DD hh:mm')
      }
    }
  },
  data() {
    return {
      loading: false,
      limit: 10,
      nextToken: ''
    }
  },
  computed: {
    ...mapGetters([
      'isLoading',
      'name',
      'avatar',
      'roles'
    ])
    // activeClient() { return this.$store.getters.apollo }
  },
  mounted() {
  },
  created() {
  },
  methods: {
  }
}
</script>

<style>
 #bad_nprogress {
    opacity: 0.4;
    background-color: #a3a3a3;
    pointer-events: none;
    position: fixed;
    top: 0px;
    left: 0px;
    min-height: 100vh;
    width: 100%;
 }
</style>

<style rel="stylesheet/scss" lang="scss" scoped>

</style>
