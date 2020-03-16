<template>
  <div>
    <label class="editable editable-label" contenteditable="true">{{ field.label }}</label>
    <select class="form-control">
      <option v-for="option in splitStringByLine(field.options)" :key="option">{{ option }}</option>
    </select>
  </div>
</template>

<script>
export default {
  props: {
    field: {
      type: Object,
      default: null
    }
  },
  methods: {
    splitStringByLine: function(str) {
      return str.split('\n')
    }
  }
}
</script>
