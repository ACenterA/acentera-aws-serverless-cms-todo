<template>
  <div>
    <AppHeader/>

    <div class="content-container build-container build-body">
      <div v-if="showAddForm" class="add-form" @click="addFormElements()">
        <div class="add-form-text">
          Add Form Element
        </div>
        <div class="circle-normal"><span class="glyphicon glyphicon-plus"/></div>

        <div class="circle-ripple circle-ripple-1"/>
        <div class="circle-ripple circle-ripple-2"/>
      </div>

      <div :class="{ 'form-elements-show': showAddForm === false }" class="form-elements">
        <div class="element-main-header">Form Elements
          <span class="glyphicon glyphicon-remove pull-right form-elements-remove" @click="removeFormElements()"/>
        </div>

        <div id="header" class="element-container">
          <div class="element-icon">
            <span class="glyphicon glyphicon-header"/>
          </div>
          <div class="element-text">
            Header
          </div>
        </div>
        <div id="name" class="element-container">
          <div class="element-icon">
            <span class="glyphicon glyphicon-user"/>
          </div>
          <div class="element-text">
            Full Name
          </div>
        </div>
        <div id="email" class="element-container">
          <div class="element-icon">
            <span class="glyphicon glyphicon-envelope"/>
          </div>
          <div class="element-text">
            Email
          </div>
        </div>
        <div id="address" class="element-container">
          <div class="element-icon">
            <span class="glyphicon glyphicon-map-marker"/>
          </div>
          <div class="element-text">
            Address
          </div>
        </div>
        <div id="input" class="element-container">
          <div class="element-text">
            Input
          </div>
        </div>
        <div id="textarea" class="element-container">
          <div class="element-text">
            Textarea
          </div>
        </div>
        <div id="checkboxes" class="element-container">
          <div class="element-text">
            Checkboxes
          </div>
        </div>
        <div id="radio_buttons" class="element-container">
          <div class="element-text">
            Radio Buttons
          </div>
        </div>
        <div id="select" class="element-container">
          <div class="element-text">
            Select
          </div>
        </div>
      </div>

      <div class="sortable-container">
        <div :class="{ 'sortable-border': fields.length === 0 }" class="sortable">
          <div v-for="(field, index) in fields" :id="'list-' + index" :key="field.id" :class="{ 'focused-element': field.isFocused === true }" tabindex="-1" class="form-group form-element-container" @click="elementFocus(index)">
            <div :class="{ hide: field.isFocused === false }" class="action-circles">

              <div class="action-circle properties-circle" @click="editElementProperties(index)">
                <span class="glyphicon glyphicon-cog properties-cog"/> <span class="properties-text">Properties</span>
              </div>
              <div class="action-circle delete-circle" @click="deleteElement(index)">
                <span class="glyphicon glyphicon-trash delete-trash"/> <span class="delete-text">Remove</span>
              </div>
            </div>

            <HeaderElement
              v-if="field.type === 'header'"
              :class="field.textalign"
              :field="field"
              :index="index"
              :fields="fields"
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
              :fields="fields"
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

            <div v-if="field.visibility === 'hidden'" class="element-not-visible">
              <span class="glyphicon glyphicon-exclamation-sign"/> This field is hidden and will not be seen on the form.
            </div>
          </div>
        </div>
      </div>

      <div :class="{ 'element-properties-show': showElementProperties === true }" class="element-properties">
        <div class="element-main-header">
          <span v-if="type === 'header'">Header</span>
          <span v-if="type === 'name'">Full Name</span>
          <span v-if="type === 'address'">Address</span> Properties
          <span class="glyphicon glyphicon-remove pull-right form-elements-remove" @click="removeElementProperties()"/>
        </div>
        <div v-if="type === 'header'" class="element-property">
          <div class="form-group">
            <label>Heading Text</label>
            <input v-model="label" type="text" class="form-control" @keyup="editLabel()">
          </div>
        </div>
        <div v-if="type !== 'header'" class="element-property">
          <div class="form-group">
            <label>Label Text</label>
            <input v-model="label" type="text" class="form-control" @keyup="editLabel()">
          </div>
        </div>
        <div v-if="type === 'header'" class="element-property">
          <div class="form-group">
            <label>Sub-Heading Text</label>
            <input v-model="subheader" type="text" class="form-control" @keyup="editSubHeader()">
            <small class="form-text text-muted">Small text below the heading.</small>
          </div>
        </div>
        <div v-if="type === 'header'" class="element-property">
          <div class="form-group">
            <label>Heading Size</label>
            <div class="radio-wrapper">
              <label :class="{ 'label-active': tagname === null || tagname === 'h1' }" class="radio-inline">
                <input v-model="tagname" type="radio" name="optradio" value="h1" @click="editTagName()">H1
              </label>
              <label :class="{ 'label-active': tagname === 'h2' }" class="radio-inline">
                <input v-model="tagname" type="radio" name="optradio" value="h2" @click="editTagName()">H2
              </label>
              <label :class="{ 'label-active': tagname === 'h3' }" class="radio-inline">
                <input v-model="tagname" type="radio" name="optradio" value="h3" @click="editTagName()">H3
              </label>
              <label :class="{ 'label-active': tagname === 'h4' }" class="radio-inline">
                <input v-model="tagname" type="radio" name="optradio" value="h4" @click="editTagName()">H4
              </label>
              <label :class="{ 'label-active': tagname === 'h5' }" class="radio-inline">
                <input v-model="tagname" type="radio" name="optradio" value="h5" @click="editTagName()">H5
              </label>
              <label :class="{ 'label-active': tagname === 'h6' }" class="radio-inline">
                <input v-model="tagname" type="radio" name="optradio" value="h6" @click="editTagName()">H6
              </label>
            </div>
          </div>
        </div>
        <div v-if="type === 'header'" class="element-property">
          <div class="form-group">
            <label>Text Alignment</label>
            <div class="radio-wrapper">
              <label :class="{ 'label-active': textalign === null || textalign === 'text-left' }" class="radio-inline">
                <input v-model="textalign" type="radio" name="optradio" value="text-left" @click="editTextAlign()">Left
              </label>
              <label :class="{ 'label-active': textalign === 'text-center' }" class="radio-inline">
                <input v-model="textalign" type="radio" name="optradio" value="text-center" @click="editTextAlign()">Center
              </label>
              <label :class="{ 'label-active': textalign === 'text-right' }" class="radio-inline">
                <input v-model="textalign" type="radio" name="optradio" value="text-right" @click="editTextAlign()">Right
              </label>
            </div>
          </div>
        </div>
        <div v-if="type === 'checkboxes' || type === 'radio_buttons' || type === 'select'" class="element-property">
          <div class="form-group">
            <label>Options</label>
            <textarea v-model="options" class="form-control" rows="5" @keyup="editOptions()"/>
          </div>
        </div>
        <div class="element-property">
          <div class="form-group">
            <label>Duplicate Question</label>
            <div class="radio-wrapper">
              <label class="radio-inline single-button">
                <input type="radio" name="optradio" @click="duplicate()">+ Duplicate
              </label>
              <small class="form-text text-muted">Duplicate this question with all saved settings.</small>
            </div>
          </div>
        </div>
        <div v-if="type === 'name' && typeof fields[activeIndex] !== 'undefined'">
          <div class="element-property">
            <div v-for="subfield in activeSubFields(fields[activeIndex].subfields)" :key="subfield.id" class="row">
              <div class="col-sm-6">{{ subfield.label_display }}</div>
              <div class="col-sm-6 col-padding">
                <input v-model="subfield.label" type="text" class="form-control">
              </div>
            </div>
          </div>

          <div v-for="subfield in subfieldsNameToggle(subfields)" :key="subfield.id" class="element-property">
            <label v-if="subfield.type === 'middle_name'">Middle Name</label>
            <label v-if="subfield.type === 'prefix'">Prefix</label>
            <label v-if="subfield.type === 'suffix'">Suffix</label>
            <div>
              <label class="switch">
                <input v-if="subfield.type === 'middle_name'" v-model="subfield.active" type="checkbox" @click="nameToggle(2)">
                <input v-if="subfield.type === 'prefix'" v-model="subfield.active" type="checkbox" @click="nameToggle(0)">
                <input v-if="subfield.type === 'suffix'" v-model="subfield.active" type="checkbox" @click="nameToggle(4)">

                <div class="slider">
                  <div :class="{'switch-on-active': subfield.active === 1}" class="switch-on">
                    ON
                  </div>
                  <div :class="{'switch-off-active': subfield.active === 1}" class="switch-off">
                    OFF
                  </div>
                </div>
              </label>
            </div>
          </div>
        </div>
        <div v-if="type === 'address'">
          <div class="element-property">
            <div v-for="subfield in activeSubFields(subfields)" :key="subfield.id" class="row">
              <div class="col-sm-6">{{ subfield.label_display }}</div>
              <div class="col-sm-6 col-padding">
                <input v-model="subfield.label" type="text" class="form-control">
              </div>
            </div>
          </div>
          <div v-for="(subfield, index) in subfields" :key="subfield.id" class="element-property">
            <label>{{ subfield.label_display }}</label>
            <div>
              <label class="switch">
                <input v-model="subfield.active" type="checkbox" @click="addressToggle(index)">

                <div class="slider">
                  <div :class="{'switch-on-active': subfield.active === 1}" class="switch-on">
                    ON
                  </div>
                  <div :class="{'switch-off-active': subfield.active === 1}" class="switch-off">
                    OFF
                  </div>
                </div>
              </label>
            </div>
          </div>
        </div>
        <div class="element-property">
          <label>Hide field</label>
          <div>
            <label class="switch">
              <input v-model="visibility" type="checkbox" value="hidden" @click="switchToggle()">
              <div class="slider">
                <div :class="{'switch-on-active': visibility === 'hidden'}" class="switch-on">
                  ON
                </div>
                <div :class="{'switch-off-active': visibility === 'hidden'}" class="switch-off">
                  OFF
                </div>
              </div>
            </label>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script>
import AppHeader from './AppHeader'
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
import 'jquery-ui/ui/widgets/draggable.js'
import 'jquery-ui/ui/widgets/sortable.js'
import { mapState } from 'vuex'

export default {
  components: {
    AppHeader,
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
      elements: {
        'header': {
          name: 'header',
          label: 'Header',
          type: 'header',
          tagname: 'h1',
          textalign: 'text-left',
          subfields: []

        },
        'name': {
          name: 'name',
          label: 'Name',
          type: 'name',
          subfields: [
            {
              name: 'prefix',
              label: 'Prefix',
              label_display: 'Prefix',
              type: 'prefix',
              placeholder: 'prefix',
              active: 0
            },
            {
              name: 'first_name',
              label: 'First Name',
              label_display: 'First Name',
              type: 'first_name',
              placeholder: 'first name',
              active: 1
            },
            {
              name: 'middle_name',
              label: 'Middle Name',
              label_display: 'Middle Name',
              type: 'middle_name',
              placeholder: 'middle name',
              active: 0
            },
            {
              name: 'last_name',
              label: 'Last Name',
              label_display: 'Last Name',
              type: 'last_name',
              placeholder: 'last name',
              active: 1
            },
            {
              name: 'suffix',
              label: 'Suffix',
              label_display: 'Suffix',
              type: 'suffix',
              placeholder: 'suffix',
              active: 0
            }
          ]
        },
        'email': {
          name: 'email',
          label: 'Email',
          type: 'email',
          placeholder: 'email',
          tagname: 'input',
          subfields: []
        },
        'address': {
          name: 'address',
          label: 'Address',
          type: 'address',
          subfields: [
            {
              name: 'street_address',
              label: 'Street Address',
              label_display: 'Street Address',
              type: 'street_address',
              placeholder: 'street address',
              active: 1
            },
            {
              name: 'street_address_line_2',
              label: 'Street Address Line 2',
              label_display: 'Street Address Line 2',
              type: 'street_address_line_2',
              placeholder: 'street address line 2',
              active: 1
            },
            {
              name: 'city',
              label: 'City',
              label_display: 'City',
              type: 'city',
              placeholder: 'city',
              active: 1
            },
            {
              name: 'state',
              label: 'State',
              label_display: 'State',
              type: 'state',
              placeholder: 'state',
              active: 1
            },
            {
              name: 'zip_code',
              label: 'Zip Code',
              label_display: 'Zip Code',
              type: 'zip_code',
              placeholder: 'zip code',
              active: 1
            },
            {
              name: 'country',
              label: 'Country',
              label_display: 'Country',
              type: 'country',
              placeholder: 'country',
              active: 1
            }
          ]
        },
        'input': {
          name: 'input',
          label: 'Input',
          type: 'text',
          tagname: 'input',
          subfields: []
        },
        'textarea': {
          name: 'textarea',
          label: 'Textarea',
          type: 'textarea',
          tagname: 'textarea',
          subfields: []
        },
        'checkboxes': {
          name: 'checkboxes',
          label: 'Checkboxes',
          type: 'checkboxes',
          tagname: 'input',
          options: 'Option 1\nOption 2\nOption 3',
          subfields: []
        },
        'radio_buttons': {
          name: 'radio_buttons',
          label: 'Radio Buttons',
          type: 'radio_buttons',
          tagname: 'input',
          options: 'Option 1\nOption 2\nOption 3',
          subfields: []
        },
        'select': {
          name: 'select',
          label: 'Select',
          type: 'select',
          tagname: 'select',
          options: 'Option 1\nOption 2\nOption 3',
          subfields: []
        }
      },
      activeIndex: null,
      fields: [],
      hasFields: false,
      label: null,
      middleName: null,
      options: '',
      showAddForm: true,
      showElementProperties: false,
      subfields: [],
      subheader: null,
      tagname: null,
      textalign: 'text-left',
      type: null,
      visibility: null,
      activeSubFields: function(subfields) {
        return subfields.filter(function(subfield) {
          return subfield.active === 1
        })
      },
      activeIndexSubFields: function() {
        return this.fields[this.activeIndex].subfields.filter(function(subfield) {
          return subfield.active === 1
        })
      },
      addFormElements: function() {
        this.showAddForm = false
        this.showElementProperties = false
      },
      addressToggle: function(num) {
        if (this.fields[this.activeIndex].subfields[num].active === true) {
          this.fields[this.activeIndex].subfields[num].active = 1
        } else {
          this.fields[this.activeIndex].subfields[num].active = 0
        }
      },
      // delete field by deleting element from page, array, and db
      deleteElement: function(index) {
        this.fields.splice(index, 1)
      },
      duplicate: function() {
        this.receiveElement(JSON.parse(JSON.stringify(this.fields[this.activeIndex])), this.activeIndex + 1)

        this.elementFocus(this.activeIndex + 1)
      },
      editElementProperties: function(index) {
        this.showAddForm = true
        this.showElementProperties = true
      },
      editLabel: function() {
        this.fields[this.activeIndex].label = this.label
      },
      editOptions: function() {
        this.fields[this.activeIndex].options = this.options
      },
      editSubHeader: function() {
        this.fields[this.activeIndex].subheader = this.subheader
      },
      editTagName: function() {
        this.fields[this.activeIndex].tagname = this.tagname
      },
      editTextAlign: function() {
        this.fields[this.activeIndex].textalign = this.textalign
      },
      elementFocus: function(index) {
        if (this.fields[index] !== undefined) {
          this.activeIndex = index
          this.label = this.fields[index].label
          this.options = this.fields[index].options
          this.type = this.fields[index].type
          this.tagname = this.fields[index].tagname
          this.textalign = this.fields[index].textalign
          this.subfields = this.fields[index].subfields
          this.subheader = this.fields[index].subheader
          this.visibility = this.fields[index].visibility

          this.fields.forEach(function(field) {
            field.isFocused = false
          })

          this.fields[index].isFocused = true
        }
      },
      nameToggle: function(num) {
        if (this.fields[this.activeIndex].subfields[num].active === true) {
          this.fields[this.activeIndex].subfields[num].active = 1
        } else {
          this.fields[this.activeIndex].subfields[num].active = 0
        }
      },
      receiveElement: function(element, newIndex) {
        this.fields.splice(newIndex, 0, {
          id: this.fields.length,
          name: element.name,
          type: element.type,
          label: element.label,
          options: element.options,
          subfields: element.subfields,
          subheader: element.subheader,
          subheader_update: true,
          placeholder: element.placeholder,
          tagname: element.tagname,
          textalign: element.textalign,
          visibility: element.visibility,
          isFocused: true,
          order_rank: newIndex
        })

        this.$store.commit('updateFields', { fields: this.fields })
      },
      removeElementProperties: function() {
        this.showElementProperties = false
      },
      removeFormElements: function() {
        this.showAddForm = true
      },
      subfieldsNameToggle: function(subfields) {
        return subfields.filter(function(subfield) {
          return subfield.type === 'prefix' ||
                        subfield.type === 'middle_name' ||
                        subfield.type === 'suffix'
        })
      },
      switchToggle: function() {
        if (this.visibility === true) {
          this.visibility = 'hidden'
        } else {
          this.visibility = null
        }

        this.fields[this.activeIndex].visibility = this.visibility
      },
      updateLabel: function(index) {
        var text = $("[contenteditable='true']").eq(index).text()
        console.log(index)

        this.fields[index].label = text
      },
      updateSubHeader: function(index) {
        var text = $('.editable-sub-' + index).eq(0).text()

        // this.fields[this.activeIndex].subheader = text;

        // reupdate text to deal with bug of vue being updated
        // and rendering text twice
        $('.editable-sub-' + index).eq(0).text(text)
      }
    }
  },
  computed: mapState([
    // map this.count to store.state.count
    'count'
  ]),
  mounted() {
    var that = this

    $('body').click(function(evt) {
      if (evt.target.className === 'form-element-container' ||
                    evt.target.className === 'element-properties') {
        return
      }

      // For descendants of "form-element-container" being clicked, remove this check if you do not want to put constraint on descendants.
      if ($(evt.target).closest('.form-element-container').length ||
                    $(evt.target).closest('.element-properties').length) {
        return
      }

      // Do processing of click event here for every element except with classname 'form-element-container'
      that.fields.forEach(function(field) {
        field.isFocused = false
        that.showElementProperties = false
      })
    })

    function setHeight() {
      var height = $(window).height()
      var offset = $('.sortable-container').offset().top

      height = height - offset
      $('.sortable-container').css({ 'height': height })
      $('.form-elements').css({ 'height': height })
      $('.element-properties').css({ 'height': height })
    }

    setHeight()

    $(window).resize(function() {
      setHeight()
    })

    $('.element-container').draggable({
      opacity: 0.7,
      helper: 'clone',
      connectToSortable: '.sortable'
    })

    $('.sortable').sortable({
      axis: 'y',
      cancel: '.editable',
      start: function(e, ui) {
        // creates a temporary attribute on the element with the old index
        $(this).attr('data-previndex', ui.item.index())
      },
      receive: function(event, ui) {
        if (ui.item.attr('id')) {
          var newIndex = parseInt($(this).data('ui-sortable').currentItem.index())
          var element = $.extend(true, {}, that.elements[ui.item.attr('id')])

          $(this).removeAttr('data-previndex')
          $(ui.helper).replaceWith('')

          that.receiveElement(element, newIndex)

          that.elementFocus(newIndex)
        }
      },
      update: function(event, ui) {
        if (ui.item.index() !== -1) {
          var newIndex = ui.item.index()

          var oldIndex = parseInt($(this).attr('data-previndex'))
          $(this).removeAttr('data-previndex')

          $(ui.helper).replaceWith('')

          that.fields.splice(newIndex, 0, that.fields.splice(oldIndex, 1)[0])

          that.$store.commit('updateFields', { fields: that.fields })

          that.elementFocus(newIndex)
        }
      }

    })
  }
}

</script>
