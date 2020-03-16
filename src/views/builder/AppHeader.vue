<template>
  <div>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
    <header class="site-header">
      <div class="container">
        <div class="row">
          <div class="col-sm-4">
            <a href="#" class="logo">Bootstrap Form Builder test</a>
          </div>
        </div>
      </div>
    </header>
    <nav class="site-nav">
      <div class="container">
        <div class="row">
          <div class="col-sm-3">
            &nbsp;
          </div>
          <div class="col-sm-6 text-center">
            <ul>
              <li><a href="#" class="show-embed" @click="showEmbedCode()">Get Code</a></li>
            </ul>
          </div>
          <div class="col-sm-3">
            &nbsp;
          </div>
        </div>
      </div>
    </nav>

    <div class="modal">
      <div class="modal-content">
        <div class="embed-code">
          <h1>Embed Code</h1>
          <textarea style="width: 100%;" class="embed-code-box">
            Test
          </textarea>
        </div>
      </div>
    </div>

    <div class="code">
      <div v-for="(field,index) in this.$store.state.fields" :key="field" class="form-group">
        <HeaderElement
          v-if="field.type === 'header'"
          :class="field.textalign"
          :field="field"
          :index="index"
        />

        <NameElement
          v-if="field.type === 'name'"
          :field="field"
        />

        <InputElement
          v-if="field.type === 'text'"
          :field="field"
        />

        <EmailElement
          v-if="field.type === 'email'"
          :field="field"
          :index="index"
        />

        <AddressElement
          v-if="field.type === 'address'"
          :field="field"
        />

        <TextareaElement
          v-if="field.type === 'textarea'"
          :field="field"
        />

        <CheckboxesElement
          v-if="field.type === 'checkboxes'"
          :field="field"
        />

        <RadioButtonsElement
          v-if="field.type === 'radio_buttons'"
          :field="field"
        />

        <SelectElement
          v-if="field.type === 'select'"
          :field="field"
        />
      </div>
    </div>
  </div>
</template>

<script>
import 'assets/css/builder.styl'
import HeaderElement from '@/components/elements/HeaderElement'
import NameElement from '@/components/elements/NameElement'
import EmailElement from '@/components/elements/EmailElement'
import AddressElement from '@/components/elements/AddressElement'
import InputElement from '@/components/elements/InputElement'
import TextareaElement from '@/components/elements/TextareaElement'
import CheckboxesElement from '@/components/elements/CheckboxesElement'
import RadioButtonsElement from '@/components/elements/RadioButtonsElement'
import SelectElement from '@/components/elements/SelectElement'
import $ from 'jquery'
// import htmlBeautify from 'html-beautify'
// import selectText from 'select-text'
import pretty from 'pretty'

export default {
  components: {
    HeaderElement,
    NameElement,
    EmailElement,
    AddressElement,
    InputElement,
    TextareaElement,
    CheckboxesElement,
    RadioButtonsElement,
    SelectElement
  },
  data() {
    return {
      code: ''
    }
  },

  mounted() {
    $('body').click(function(evt) {
      if (evt.target.className === 'modal-content' ||
            evt.target.className === 'show-embed') {
        return
      }

      if ($(evt.target).closest('.modal-content').length ||
            $(evt.target).closest('.show-embed').length) {
        return
      }

      if ($('.modal').is(':visible')) {
        $('.modal').hide()
      }
    })

    $('.embed-code-box').click(function() {
      $(this).select()
    })
  },
  methods: {
    activeSubFields: function(subfields) {
      return subfields.filter(function(subfield) {
        return subfield.active === 1
      })
    },
    formatCode(code) {
      /*
      // remove br tags
      code = code.replace(/<br>/g, '');
      // remove span tags
      code = code.replace(/<span\s+>/g, '');
      code = code.replace(/<\/span>/g, '');

      code = code.replace(/(?!^)&#60;/g, '<br />&#60;');

      code = code.replace(/</gi, '&#60;');

      code = code.replace(/>/gi, "&#62;<br />");*/

      // remove content editable attribute
      code = code.replace(/contenteditable="true"/g, '')

      // remove class editable
      code = code.replace(/class="editable editable-label"/g, '')

      // remove <!----> tags
      code = code.replace(/<!---->/g, '')

      return code
    },
    showEmbedCode() {
      $('.modal').show()
      var code = $('.code').html()
      $('.embed-code-box').html(pretty(this.formatCode(code)))
    },
    hideEmbedCode() {
      this.$modal.hide('embed-code')
    }
  }
}

</script>
