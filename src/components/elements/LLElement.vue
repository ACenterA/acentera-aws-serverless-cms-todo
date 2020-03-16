<template>
  <div>
    <label class="editable editable-label" contenteditable="true">{{ field.label }}</label>
    <div class="row">
      <div v-for="subfield in activeSubFields(field.subfields)" :key="subfield" class="col-xs-6">
        <input :name="subfield.name" :placeholder="subfield.placeholder" type="text" class="form-control">
        <span class="help-block">{{ subfield.label }}</span>
      </div>
    </div>
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
    activeSubFields: function(subfields) {
      return subfields.filter(function(subfield) {
        return subfield.active === 1
      })
    }
  }
}
</script>
