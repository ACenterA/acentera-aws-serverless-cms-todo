<template>
  <div class="row">
    <div class="col-xs-12">
      <label class="editable editable-label" contenteditable="true" @focusout="updateLabel(index)">{{ field.label }}</label>
      <input :placeholder="field.placeholder" type="email" name="email" class="form-control">
    </div>
  </div>
</template>

<script>
import $ from 'jquery'

export default {
  props: {
    field: {
      type: Object,
      default: null
    },
    index: {
      type: Number,
      default: 0
    },
    fields: {
      type: Object,
      default: null
    }
  },
  data: function() {
    return { fieldsRef: this.fields }
  },
  methods: {
    updateLabel: function(index) {
      var text = $("[contenteditable='true']").eq(index).text()

      this.fieldsRef[index].label = text
    }
  }
}
</script>
