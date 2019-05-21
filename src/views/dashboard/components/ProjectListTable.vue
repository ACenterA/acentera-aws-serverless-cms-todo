<template>
  <div>
    <el-table :data="list" style="width: 100%;padding-top: 15px;">
      <el-table-column label="Clusters" min-width="200">
        <template slot-scope="scope">
          {{ scope.row.title | orderNoFilter }}
        </template>
      </el-table-column>
      <el-table-column label="Creation Date" width="195" align="center">
        <template slot-scope="scope">
          {{ scope.row.created }}
        </template>
      </el-table-column>

      <el-table-column :label="$t('table.action')" min-width="150px">
        <template slot-scope="scope">

          <el-button type="primary" size="mini" @click="handleSelect(scope.row)">
            {{ $t('table.select') }}
          </el-button>

        </template>
      </el-table-column>
    </el-table>
    <div class="el-table el-table--fit el-table--enable-row-hover el-table--enable-row-transition el-table--medium" style="padding-top:20px;padding-bottom:20px;padding-left:10px;padding-right:50px;text-align:center;">
      <div class="min-width:395px">
        &nbsp;
      </div>
      <el-button @click="createNew()">New</el-button>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        success: 'success',
        pending: 'danger'
      }
      return statusMap[status]
    },
    orderNoFilter(str) {
      return str.substring(0, 30)
    }
  },
  data() {
    return {
      list: []
    }
  },
  computed: {
    ...mapGetters([
      'getClusterProjects',
      'getClusterProjectsLst'
    ])
  },
  watch: {
    getAmiProjects(valT) {
      console.error('NEW AMI CHANGE')
      console.error(valT)
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.list = this.getClusterProjectsLst
    },
    handleSelect(row) {
      console.error('clicked on element')
      console.error(row)
      this.$emit('onRowSelect', row)
    },
    createNew() {
      this.$router.push({
        name: 'createcluster'
      })
    },
    renderHeader(h, { column }) {
      return h('div', { }, [
        h('div', { 'class': 'cell' }, column.label)
      ])
    }
  }
}
</script>
